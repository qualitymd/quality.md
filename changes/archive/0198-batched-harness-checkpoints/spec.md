---
type: Functional Specification
title: Batched harness checkpoints
description: Requirements for keeping a bounded window of dependency-ready judgment requests outstanding to the invoking harness and accepting correlated results as they stream back, while preserving runner-owned orchestration and parity with sequential harness runs.
tags: [evaluation, evaluator, agents, concurrency]
timestamp: 2026-07-11T00:00:00Z
---

# Batched harness checkpoints

This spec governs the delta from a single-request harness checkpoint to a
bounded, rolling window of outstanding requests for
`qualitymd evaluation run --evaluator harness`. It inherits the binding work
graph, per-result validation, persistence, and report semantics from the durable
**(normative)** evaluation
[runner](../../../specs/evaluation/runner.md),
[evaluator contract](../../../specs/evaluation/evaluator-contract.md),
[checkpoint protocol](../../../specs/evaluation/protocol.md), and
[orchestration contract](../../../specs/evaluation/orchestration.md). It does not
change evaluation judgment semantics or result payload schemas.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

The harness transport stages exactly one pending work request per checkpoint and
exits, so the invoking harness never holds more than one request at a time and
harness-backed runs are pinned to concurrency 1. The dependency-aware concurrent
scheduler from [0195](../../archive/0195-evaluation-concurrency.md) already
computes the dependency-ready frontier (`nextReadyEvaluationStep`) and runs it
under a worker pool for in-process evaluators; the harness path is routed around
it because the transport cannot hold more than one in-flight call.
[0194](../../archive/0194-harness-native-evaluator-dispatch.md) reserved an
unused capability for exactly this internal judgment parallelism.

Keeping up to `concurrency` requests from the ready frontier **outstanding** at
once — emitting each newly-ready request as an earlier one is accepted, rather
than draining a whole wave before starting the next — lets the invoking harness
judge, or fan out to subagents, several independent self-contained requests at a
time while the runner keeps sole ownership of scheduling, validation, and
persistence. Requests are drawn only from the dependency-ready frontier, so
parallelism is bounded by the work-graph DAG: wide models gain, deep chains stay
serial. Results stream back over ordinary discrete resume calls, so no
long-lived transport is introduced. The change is a transport optimization; it
must not change what a run evaluates.

## Scope

Covered: a rolling window of concurrent outstanding requests, result intake with
per-member correlation, the resolved-concurrency cap on the number of
outstanding requests, parity with the sequential path, independent member
durability, partial-submission and resume behavior, the harness concurrency
capability, and the evaluate-workflow loop.

Deferred: a persistent long-lived request/result transport (for example one
JSON-RPC or JSONL subprocess) as an alternative to discrete resume calls;
adaptive window sizing beyond the resolved concurrency cap; runner-side subagent
scheduling; and non-harness evaluators.

Non-goal: changing judgment semantics, accepted payload schemas, rating roll-up,
generated reports, or run-folder layout; granting the harness authority to
schedule units, widen source, or persist results.

## Assumptions and dependencies

- [0194](../../archive/0194-harness-native-evaluator-dispatch.md) (harness
  transport) and [0195](../../archive/0195-evaluation-concurrency.md) (concurrent
  scheduler and `evaluation.concurrency`) have landed.
- The invoking harness can hold several outstanding requests and submit their
  results — with ordinary reasoning or native subagents — over successive resume
  calls, returning one result object per request.
- The runner can rebuild every outstanding request from its model snapshot,
  work-graph state, and current source package, and compare each request's input
  hash with the persisted pending metadata (the single-request rebuild from 0194,
  generalized to a set).

## Requirements

### R1 — Keep the ready frontier saturated up to the cap

While a harness run has dependency-ready judgment work that is not yet
outstanding, the runner **MUST** keep up to the run's resolved concurrency
judgment requests outstanding, emitting each newly-ready request in an
`awaiting_evaluator` receipt as capacity frees — without waiting for other
outstanding requests to be judged first. Each emitted request **MUST** carry the
same complete bounded work request a single-request checkpoint carries today.

> Durable spec: modify `specs/evaluation/orchestration.md` — the checkpoint
> returns the outstanding-request set (the checkpoint transport contract lives
> there, not in `protocol.md`); `specs/evaluation/runner.md` — rolling
> ready-frontier emission bounded by the cap; `specs/cli/evaluation-run.md` —
> receipt output shape.

### R2 — Accept and correlate streamed results

When a harness run is resumed with submitted results, the runner **MUST** accept
one or more result envelopes in a single submission and bind each to an
outstanding request by its request identifier and input hash; envelope order
**MUST NOT** affect binding. If an envelope matches no outstanding request by
both identifier and input hash, then the runner **MUST** reject that envelope
without discarding or altering any other member's accepted result.

> Rationale: binding is by the runner's rebuilt-request hash, so a stale or
> mismatched result cannot be accepted against evidence it was not judged
> against — the single-request integrity check from 0194, preserved per member.

> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/evaluator-contract.md` — streamed result acceptance and
> correlation; `specs/cli/evaluation-run.md` — result input accepts one or more
> envelopes per call.

### R3 — Parity with the sequential path

For identical model, scope, and source, a harness run at resolved concurrency
greater than 1 **MUST** produce the same accepted results, requirement and area
ratings, and generated reports as the same run at concurrency 1. The window size
**MUST NOT** be observable in any accepted judgment or payload.

> Rationale: this is the boundary that keeps the rolling window a transport
> optimization; a windowed run that scored differently would silently corrupt
> evaluations, so the constraint is MUST, not SHOULD.

> Durable spec: none.

### R4 — Independent member acceptance and durability

When the runner accepts a submitted result, it **MUST** accept each valid member
independently and persist it atomically with its work unit's completion, so that
an interruption leaves every accepted member durable and every still-outstanding
member recoverable for resume.

> Durable spec: modify `specs/evaluation/runner.md` — per-member atomic
> persistence within the window; `specs/evaluation/evaluation-json.md` — artifact
> schema version bump for multi-outstanding checkpoint state.

### R5 — Partial submission, failure, and resume

Each resume submission **MAY** carry results for any subset of the outstanding
requests. The runner **MUST** accept the valid members, free the capacity they
held, and leave the still-outstanding requests in place; a still-outstanding
request that was not submitted **MUST NOT** consume retry budget. If a submitted
member is schema-invalid or returns a classified failure, then the runner
**MUST** re-emit that member for a retry attempt, consuming that work unit's
retry budget as today, without altering run state or accepted results. When a
harness run is resumed without submitted results, the runner **MUST** re-emit the
current outstanding set unchanged, and it **MUST** continue to refuse a resume
whose evaluator differs from the one the run recorded.

> Rationale: partial submission is the normal path under a rolling window, so
> "not yet judged" (still outstanding, no cost) and "failed" (re-emit, retry
> cost) must be distinct — otherwise a partial reply would burn a judgeable
> unit's retry budget.

> Durable spec: modify `specs/evaluation/orchestration.md` — partial
> submission, re-emission, and resume-without-result over the outstanding set;
> the scheduler tops up and re-emits the outstanding set.

### R6 — Harness concurrency capability and cap

The harness evaluator **MUST** advertise that it can serve subagent-delegable
judgment requests, and the runner **MUST NOT** reduce a harness run's resolved
concurrency to 1. The resolved concurrency **MUST** bound the number of
outstanding (emitted, not-yet-accepted) requests at any time. At resolved
concurrency 1, the runner **MUST** keep exactly one request outstanding,
preserving current single-request behavior.

> Rationale: concurrency 1 must stay byte-identical to today so existing runs,
> automations, and interrupted-run resume are unaffected by the new path.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` — harness gains
> the subagent-delegation capability and is no longer clamped to 1;
> `specs/cli/evaluation-run.md` and `specs/cli/evaluation-status.md` — resolved
> concurrency and outstanding-window width in dry-run preview and run receipts.

### R7 — Evaluate workflow judges the window

When the evaluate workflow receives outstanding requests, it **MUST** judge each
within that request's own bounded boundary and submit one correlated result
envelope per request, **MAY** submit results as they become ready rather than
waiting for the whole window, and **MAY** delegate independent requests to
subagents. The workflow **MUST NOT** schedule work units, widen source beyond a
request's boundary, or alter accepted output.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/quality-skill.md` — rolling checkpoint loop,
> streamed submission, subagent delegation, one envelope per request.

## Open questions

- **Receipt-size clamp.** A soft implementation question, not a design one:
  whether a model with large bounded source bundles needs the outstanding window
  clamped below the resolved concurrency to bound receipt size. This is a
  runner-internal clamp with no public setting (see the design doc); leave it
  unclamped until a fixture shows it matters.

Resolved during Design:

- **Usage aggregation.** No change. Usage is reported and logged per result
  today (one entry per evaluator call), with no accumulated run total. Each
  accepted member flows through that same per-result path, so a windowed run
  logs the same per-member usage a sequential run does — parity extends to usage
  with no new aggregate. Any run-total or per-window rollup is a separate
  observability change, out of scope here.

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation/orchestration.md` — outstanding-set checkpoint, streamed
  result intake and correlation, partial submission, re-emission, and
  resume-without-result (R1, R2, R5). Resolved during implementation: the
  checkpoint transport contract lives in the orchestration and runner specs;
  `specs/evaluation/protocol.md` carries protocol moves and judgment
  semantics, which this case does not change, so it needs no edit.
- `specs/evaluation/runner.md` — rolling ready-frontier emission and per-member
  atomic persistence (R1, R4).
- `specs/evaluation/evaluator-contract.md` — streamed result acceptance and the
  harness subagent-delegation capability (R2, R6).
- `specs/evaluation/evaluation-json.md` — artifact schema bump for
  multi-outstanding checkpoint state (R4).
- `specs/cli/evaluation-run.md` — outstanding-set receipt output, streamed result
  input, window width in dry-run preview (R1, R2, R6).
- `specs/cli/evaluation-status.md` — awaiting-run inspection with multiple
  outstanding requests (R6).
- `specs/skills/quality-skill/evaluation.md`,
  `specs/skills/quality-skill/quality-skill.md`, and (found during the
  implementation sweep) `specs/skills/quality-skill/workflows/evaluate.md` —
  rolling checkpoint loop, streamed submission, subagent delegation (R7).

### To rename

None.

### To delete

None.
