---
type: Design Doc
title: Batched harness checkpoints design
description: Generalize the single-request harness checkpoint to a bounded rolling window of dependency-ready requests, reusing the concurrent scheduler's ready-frontier computation to top up and stream results over discrete resume calls while the runner keeps sole ownership of scheduling, validation, and persistence.
tags: [evaluation, evaluator, agents, concurrency]
timestamp: 2026-07-11T00:00:00Z
---

# Batched harness checkpoints design

## Context

This design answers the
[batched harness checkpoints functional spec](spec.md). Today the harness
transport emits exactly one request per checkpoint: `runHarnessUnit` builds one
`WorkRequest`, `checkpointHarness` persists a single
`State.PendingEvaluatorCall` and sets `awaiting_evaluator`, and the engine's
sequential loop returns. The concurrent scheduler from
[0195](../../archive/0195-evaluation-concurrency.md) already computes the
dependency-ready frontier (`nextReadyEvaluationStep`) and runs it under a worker
pool, but harness runs are routed around it (`harnessBacked()` forces the
sequential loop) because the transport cannot hold more than one in-flight call.

The change keeps a bounded **rolling window** of outstanding requests — the ready
frontier, capped by the run's resolved concurrency — and streams results back
over ordinary discrete resume calls, refilling the window as members are
accepted. Every orchestration decision and authoritative mutation stays inside
the runner. It is a transport change, not a semantics change: the accepted
results of a windowed run must match a `concurrency: 1` run
([R3](spec.md#r3--parity-with-the-sequential-path)).

## Approach

### Pending state becomes a set

`State.PendingEvaluatorCall *PendingEvaluatorCall` becomes a slice, persisted in
`evaluation.json`, and `ArtifactSchemaVersion` bumps `6 → 7`:

```go
// State
PendingEvaluatorCalls []*PendingEvaluatorCall `json:"pendingEvaluatorCalls,omitempty"`
```

Each element keeps today's shape (`RequestID`, `WorkUnitID`, `InputHash`,
`CorrelationID`, `Attempt`) — correlation metadata only, never raw prompt,
source, or result bodies. `harnessRequestID(evaluationID, workUnitID, inputHash,
attempt)` stays per-unit, so a member's identity and retry accounting are
unchanged; only their number outstanding at once grows. This slice is the single
source of truth for "what is awaiting" and resumes as a unit.

Early-alpha clean break: schema 7 does not read a schema-6 singular pending
field. An awaiting run created before the upgrade cannot resume across it; it is
re-run, not migrated (consistent with the repo's no-compat-shim policy). Only
in-flight awaiting runs are affected — completed runs are unaffected.

### One control loop: apply, top up, emit

The engine stops special-casing harness into the one-at-a-time loop. Each runner
invocation (initial run or resume) runs the same three steps, reusing the
scheduler's frontier computation (`nextReadyEvaluationStep`) that `scheduleReady`
already uses to pull ready units up to the cap:

```text
on run / resume:
  1. apply submitted results (if any):
       for each envelope: bind to an outstanding PendingEvaluatorCall by
         (RequestID, InputHash); run through the existing validate/merge/persist
         path; accepted -> clear its pending entry; failed/invalid -> bump Attempt,
         keep pending; unknown correlation -> reject, touch nothing else
  2. top up the window:
       while len(outstanding) < resolvedConcurrency:
         unit = nextReadyEvaluationStep(units, scheduled, done)   // deterministic graph order
         if unit == nil: break
         if unit is deterministic: run inline, mark done, continue // frames, source, reports
         build WorkRequest; inputHash = workUnitInputHash(req)
         checkpointHarness(unit, req, inputHash, attempt)          // append to PendingEvaluatorCalls
  3. if outstanding is empty and graph complete: finish run
     else: persist atomically; emit awaiting receipt carrying the current
           outstanding set (its []EvaluatorRequest); exit 0
```

Because accepting results in step 1 frees capacity that step 2 immediately
refills from newly-ready units, the window stays saturated up to the cap without
draining a whole wave first — the rolling behavior
[R1](spec.md#r1--keep-the-ready-frontier-saturated-up-to-the-cap) requires.
Deterministic units on the frontier (frame, walkable source resolution, report
build) still run inline and never enter the window; only judgment (and harness
`resolveSource`) requests are emitted. Gathering in graph order keeps the emitted
set and persisted state stable across identical runs.

### Streamed result intake

`--evaluator-result` accepts one submission carrying one or more
`HarnessResultEnvelope` (a lone object stays valid as a single-member
submission, so a `concurrency: 1` loop is byte-identical). `consumeHarnessResult`
generalizes to iterate members and apply each through the **existing** per-result
path unchanged — normalization, schema validation, retry accounting,
accepted-result merge, logging, atomic persistence — so each accepted member is
durable the instant its unit completes
([R4](spec.md#r4--independent-member-acceptance-and-durability)). Members are
applied in correlation order, not submission order, so persisted state does not
depend on how the harness ordered its reply.

Each submission may cover any subset of the outstanding set
([R5](spec.md#r5--partial-submission-failure-and-resume)):

- **accepted** members clear their `PendingEvaluatorCall`;
- **failed** members (envelope with `Failure`, or schema-invalid) increment the
  unit's `Attempt` and stay outstanding, consuming that unit's existing retry
  budget;
- **not submitted** members simply stay outstanding — no retry cost, because
  "not yet judged" is not a failed attempt.

Resume without `--evaluator-result` re-emits the current outstanding set
unchanged. A resume whose evaluator differs from the recorded one stays refused.
Because each request is a self-contained evidence boundary (its own bounded
source bundle), judging members concurrently cannot cross-contaminate — the
property [R3](spec.md#r3--parity-with-the-sequential-path) rests on.

### Capability and cap

The harness evaluator advertises the subagent-delegation capability, and
`resolveConcurrency` no longer clamps harness to 1 — the resolved
`evaluation.concurrency` (default `max(2, NumCPU*2)` from 0195) becomes the
outstanding-window cap. At resolved concurrency 1 the top-up loop stops after one
member, so the receipt, pending state, and loop are identical to today
([R6](spec.md#r6--harness-concurrency-capability-and-cap)). Dry-run preview and
run receipts surface resolved concurrency and the current window width.

### Skill loop streams the window

The evaluate workflow's checkpoint step changes from "judge the request" to
"judge the outstanding requests, submitting results as they become ready":

```text
awaiting_evaluator? -> dispatch each outstanding request to self or a subagent
                    -> as results come ready, submit them (one or several per
                       --evaluator-result call); each call returns a receipt with
                       the window topped up with newly-ready requests
                    -> repeat until terminal receipt
```

Submitting several ready results per call trades a little latency for fewer
resume round trips; submitting each as it lands maximizes overlap. Either is
conformant. The workflow still gets no authority to schedule units, widen source,
or alter accepted output; the runner's outstanding set is the entire boundary for
the turn ([R7](spec.md#r7--evaluate-workflow-judges-the-window)).

## Spec response

- **Saturated window and cap ([R1](spec.md#r1--keep-the-ready-frontier-saturated-up-to-the-cap), [R6](spec.md#r6--harness-concurrency-capability-and-cap)):** apply-then-top-up refills freed capacity from the reused frontier each resume; concurrency 1 degenerates to the current single-request path.
- **Correlation and validation ([R2](spec.md#r2--accept-and-correlate-streamed-results)):** per-member `(RequestID, InputHash)` binding precedes the existing validation path, exactly as the singular check does today.
- **Parity ([R3](spec.md#r3--parity-with-the-sequential-path)):** self-contained per-request evidence plus correlation-order application make accepted results independent of window size; verified by running a fixture model at concurrency 1 and >1 and diffing accepted results and reports.
- **Durability and partial submission ([R4](spec.md#r4--independent-member-acceptance-and-durability), [R5](spec.md#r5--partial-submission-failure-and-resume)):** each member reuses the atomic per-unit persistence; the outstanding set is recomputed and re-emitted each resume, so an interruption or a partial reply resumes the exact remainder.
- **Skill ([R7](spec.md#r7--evaluate-workflow-judges-the-window)):** the loop dispatches the outstanding set and streams envelopes back; no new runner authority is granted to the harness.

## Alternatives

### Barrier per checkpoint (drain each wave before the next)

Rejected. Emitting the whole ready frontier, then requiring all of it back before
emitting the next wave, has a simpler intake (expect the full array) but the
slowest member idles the freed capacity until the wave completes — the classic
barrier-latency loss. The rolling window keeps utilization high for only a modest
increase in intake complexity, which the recomputed-outstanding-set logic needs
anyway for partial-reply resume. Streaming was chosen for this reason.

### Persistent long-lived transport (JSON-RPC / JSONL subprocess)

Rejected for this slice, deferred in the spec. A durable process streaming
requests and results without exiting would cut resume round trips further, but
some harnesses expose only discrete command calls (0194 rejected it on the same
grounds). The rolling window over discrete resume calls reuses the existing
exit/persist/resume machinery verbatim; a persistent transport can be added later
under the same envelopes without changing request identity.

### A separate harness-specific window setting

Rejected. Reusing `evaluation.concurrency` keeps one user-facing knob (0195's
whole point) and one mental model. If receipt size ever forces a lower bound,
that is a runner-internal clamp, not a new public setting.

### Treat any not-submitted member as a failed attempt

Rejected. Under streaming, partial submission is the normal path; collapsing
"not yet submitted" into "failed" would burn a unit's retry budget on every
partial reply and eventually fail a judgeable unit. Distinguishing them keeps
retries meaningful.

## Trade-offs and risks

- **More resume round trips.** Streaming submits results in smaller groups than a
  single per-wave array, so a windowed run makes more `evaluation run --resume`
  calls. Each is cheap and the harness may still batch several ready results into
  one call; the win is that no member idles the window. Net faster than
  one-at-a-time and than a barrier.
- **Default concurrency now applies to harness.** Omitted-concurrency harness runs
  jump from forced-1 to the 0195 default. Intended, but a visible behavior change:
  receipts and dry-run must state resolved concurrency so it is never silent, and
  the skill should name the window width on the first windowed run.
- **Larger awaiting receipts.** A receipt carries up to `concurrency` bounded
  source bundles at once. Output stays machine-oriented JSON under existing source
  limits; a receipt-size clamp below the concurrency cap can be added internally
  without a public knob.
- **Partial-reply accounting.** The not-submitted-vs-failed distinction is the
  subtle correctness point; it must be covered by tests (partial submission,
  all-failed submission, duplicate envelope, unknown correlation, resume mid-
  window).
- **Clean-break schema bump.** In-flight awaiting runs at upgrade cannot resume.
  Acceptable under early-alpha policy; note it in the changelog.

## Open questions

- **Receipt-size clamp.** `usageLogFields` aside, the one open implementation
  question: a receipt carries up to `concurrency` bounded source bundles, which a
  large model could make big. If a fixture shows it, cap the outstanding window
  below the resolved concurrency as a runner-internal clamp — no public knob.
  Left unclamped until measured.

Resolved during Design:

- **Usage aggregation.** No change. Usage is logged per evaluator call today
  (`engine.go` `usageLogFields` → the evaluator-call log), with no accumulated
  run total in the artifact. Each accepted member reuses that same per-result
  path, so a windowed run logs the same per-member usage a sequential run does.
  A run-total or per-window rollup would be a separate observability change, out
  of scope here.
- **Default concurrency for harness.** Keep 0195's shared default; do not invent
  a harness-specific default (that would reintroduce the special-casing this
  change removes). Resolved concurrency and window width are surfaced in receipts
  and dry-run, and the skill names the window width on the first windowed run, so
  the jump from forced-1 is never silent.
