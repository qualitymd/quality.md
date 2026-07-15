---
type: Design Doc
title: Transport-aware evaluator concurrency — design
description: Capability model, shared run preparation, coordinator-owned direct-call pool, harness delegation boundary, and deterministic persistence.
tags: [evaluation, evaluator, concurrency, agents, runner]
timestamp: 2026-07-15T00:00:00Z
---

# Transport-aware evaluator concurrency — design

Archived design behind the
[Transport-aware evaluator concurrency](../0204-transport-aware-evaluator-concurrency.md)
Change Case and its [functional spec](spec.md).

## Context

The current work graph already exposes a deterministic, model-ordered ready
frontier. Harness resume can accept a subset of outstanding requests and refill
the window. The gap is how transport capacity is represented and how direct SDK
calls are driven: capability booleans are ambiguous, automatic capacity comes
from host CPU count, provider creation patches a harness artifact after staging,
and `Effect.forEach` drains each request wave before the runner resumes.

The design keeps one graph and one acceptance path. It changes run preparation
so evaluator identity and capacity are inputs to the first artifact, then gives
the direct transport a coordinator that reacts to individual completions. The
harness transport continues to checkpoint a bounded window for an external
agent to service.

## Approach

### Explicit dispatch capability

Replace the evaluator-level `concurrent` and `subagents` booleans with one
structured capability:

```ts
interface EvaluatorDispatchCapability {
  readonly concurrentCalls: boolean
  readonly delegatedRequests: boolean
  readonly automaticConcurrency: number
  readonly maxConcurrency?: number
}
```

`EvaluatorCapabilities` carries this value as `dispatch`. Validation enforces
positive integers, `automaticConcurrency <= maxConcurrency` when a maximum is
present, and `automaticConcurrency: 1` plus `maxConcurrency: 1` when neither
dispatch mechanism is supported.

The built-ins declare:

```ts
harness: { concurrentCalls: false, delegatedRequests: true,  automaticConcurrency: 4 }
codex:   { concurrentCalls: true,  delegatedRequests: false, automaticConcurrency: 4 }
claude:  { concurrentCalls: false, delegatedRequests: false, automaticConcurrency: 1, maxConcurrency: 1 }
```

Configured profiles reuse the built-in capability for their `kind`. Model and
command overrides do not alter transport capacity. Codex keeps multi-agent
features disabled inside each fresh session; its parallelism is several direct
sessions. Harness delegation describes the invoking parent agent, not a child
agent tool inside the CLI.

### Pure concurrency resolution

A pure evaluation-domain function resolves the cap:

```ts
interface ConcurrencyResolution {
  readonly source: "configured" | "automatic"
  readonly requested: number
  readonly resolved: number
  readonly clamped: boolean
}

resolveConcurrency(configured, dispatch)
```

The requested value is the configured positive integer when present and the
dispatch capability's automatic value otherwise. The resolver clamps it to the
optional maximum. A sequential dispatch capability always resolves to one.
Workspace parsing gains one shared positive-integer validation boundary instead
of leaving partial validation in dry-run. `HostRuntime.hardwareConcurrency`
leaves evaluation resolution and is removed entirely if no other use remains.

Only `resolved` enters `evaluation.json` and public receipts as `concurrency`.
The complete resolution record belongs in run-local events so operators can
distinguish an automatic choice from an explicit cap without expanding the
authoritative evaluation payload.

### Shared selected-evaluator-first run preparation

Split today's harness-specific creation into a shared preparation operation:

```text
resolve workspace/model/scope
        ↓
select evaluator and validate capabilities
        ↓
resolve concurrency
        ↓
build plan, graph, frames, sources, and identity
        ↓
stage at most resolved-concurrency ready requests
        ↓
atomically write the selected evaluator's initial artifact
```

The prepared value contains the immutable manifest, graph inputs, initial ready
requests, and transport selection. A harness driver persists and returns the
awaiting receipt. A direct driver starts workers from the same prepared ready
set. There is no create-as-harness step and no post-write mutation of evaluator,
capabilities, or concurrency.

Pure model/plan/request/artifact construction remains under `src/domain/`.
Filesystem, time, provider sessions, and workflow order remain in the
application/service/adapter layers per the Effect TypeScript guide.

### Coordinator-owned direct-call pool

The direct provider path uses one application-level coordinator and scoped
workers:

1. The coordinator reads the deterministic ready frontier and starts work in
   graph order until `active.size === resolvedConcurrency` or no evaluator work
   is ready.
2. Each worker owns exactly one fresh evaluator call and emits a completion
   value—response or classified failure—through an Effect queue. It never reads
   or writes the run artifact.
3. The coordinator takes one completion, validates it through the shared
   acceptance path, seals evidence where applicable, writes the artifact
   atomically, and only then removes that active slot.
4. It executes newly ready deterministic units and starts the next evaluator
   units immediately, without waiting for unrelated active workers.
5. Terminal completion builds reports; terminal failure or cancellation stops
   new dispatch and interrupts every scoped worker.

The coordinator is the single mutable owner. Worker fibers and the completion
queue are scoped resources, so interruption propagates through provider SDK
streams and child runtimes. A stable work-unit key prevents late or duplicate
completion from being accepted after cancellation, retry, or supersession.

This replaces the current checkpoint-wave adapter for direct runs. Harness and
direct transports still share request building, result normalization,
validation, evidence sealing, retry classification, payload merge ordering, and
atomic artifact writes; they differ only in how a bounded request reaches its
worker and returns.

### Deterministic merge and per-result persistence

Completion order is deliberately nondeterministic. Acceptance order is the
order completions reach the coordinator, but persisted payload arrays are
reprojected into graph order after every merge. A requirement's assessment,
rating, and sealed evidence remain one acceptance unit.

Every accepted result produces an atomic artifact write before its completion
can release capacity or unblock a dependent unit. This is more write-intensive
than persisting one drained batch, but it gives resume a precise durable
boundary: completed work is reusable, active-but-unaccepted work is incomplete,
and no worker owns hidden state needed to reconstruct the run.

Retryable failures return the unit to runner scheduling with its incremented
attempt and do not affect sibling calls. A terminal failure records the unit and
run failure, interrupts other active calls, and leaves their unaccepted units
incomplete. Resume starts fresh sessions only for incomplete or retryable work.

### Harness delegation stays outside the CLI

The harness path remains a persisted outstanding set capped by resolved
concurrency. A resume may submit one result or any subset. The runner accepts
each submitted member, executes deterministic dependencies, and fills the
window from the ready frontier before returning the next receipt.

The `/quality` workflow treats the receipt as eligible work, not proof of active
parallelism. It may:

- judge a request in the parent session;
- assign one request to one native subagent, passing exactly the bounded request;
  or
- service fewer requests than the window and submit results incrementally.

The CLI does not import Codex-, Claude-, or harness-specific subagent APIs. A
worker never receives `evaluation.json`, the whole ready frontier, or permission
to choose another request. This keeps native worker availability and placement
an agent concern while graph progress remains a CLI concern.

### Observability and schema transition

Public dry-run and run receipts keep `concurrency` as the resolved cap. Harness
human/progress output says, for example, `2 outstanding of up to 4`; it does not
say four requests are running.

Run-local events record:

- dispatch mode (`direct`, `delegated`, or `sequential`);
- configured or automatic resolution source;
- requested, automatic, maximum when present, and resolved caps; and
- observed peak active direct calls or peak outstanding harness requests.

Evaluator-call logs retain per-call timing, usage, and cache metadata from 0203.
No scheduler decision depends on cache hits. This case does not introduce a
warm-up call or cohort barrier; those require separate evidence and design.

Because `manifest.evaluatorCapabilities` changes shape, the artifact advances to
schema version 9. Resume, report rebuild, status, and list consumers accept the
new current schema. Version-8 in-flight runs fail with the existing clean-break
remedy; completed version-8 artifacts remain historical records.

### Verification design

Use deterministic Effect test Layers and controllable evaluator services:

- a pure table test covers configured/automatic, maximum clamp, and sequential
  resolution for all built-ins and configured profiles;
- an initialization fixture with several ready requirements proves the first
  persisted pending/active set never exceeds the selected provider's cap;
- deferred worker controls prove a cap of two starts exactly two calls, then
  completing one starts the third while the second remains active;
- deliberately reversed completion order proves deterministic payload and
  report ordering;
- the same fixture at one and greater-than-one concurrency proves artifact and
  report parity apart from timestamps, capability/schema metadata, and logs;
- interruption after one accepted result proves that result is durable and
  active siblings are cancelled and resumable;
- retry fixtures prove one failure does not consume sibling budgets;
- harness fixtures preserve partial submission, no-cost unsubmitted requests,
  and rolling top-up;
- CLI tests protect dry-run, receipt, human progress, schema-version refusal,
  and the absence of a new `--concurrency` flag; and
- `mise run check` is the pre-review gate.

No repository smoke-test helper is added; integration fixtures remain inside
the existing test suite and temporary filesystem directories.

## Spec response

- R1 is answered by one graph/acceptance coordinator and workers that receive
  only runner-issued requests.
- R2–R3 are answered by the structured dispatch capability and pure resolver,
  with transport defaults instead of CPU-derived capacity.
- R4 is answered by shared selected-evaluator-first run preparation and removal
  of manifest patching.
- R5 is answered by scoped direct workers, a completion queue, immediate
  per-result acceptance, and top-up from the deterministic frontier.
- R6 is answered by retaining the checkpoint window while keeping subagent
  placement outside the CLI and passing exactly one request per worker.
- R7 is answered by serialized acceptance, an atomic write before slot release,
  scoped cancellation, and per-unit retry state.
- R8 is answered by cap-oriented public output, run-local resolution/peak
  diagnostics, schema version 9, and deterministic concurrency fixtures.

## Alternatives

### Give the whole graph to the parent agent

Rejected. The agent could choose subagents and dependencies dynamically, but it
would also become responsible for retry, persistence, resume, and ordering. That
restores the harness-specific orchestration problem the deterministic runner was
built to remove.

### Let the CLI spawn native subagents

Rejected. It couples `qualitymd` to harness-specific APIs and authentication,
duplicates the parent agent's worker manager, and makes the standalone runner
less portable. The CLI should expose bounded eligible work; the harness should
map it to its actual slots.

### Keep `concurrent` and `subagents` booleans

Rejected. Their current values already demonstrate the ambiguity: Codex can run
several SDK sessions while nested agents are explicitly disabled, and harness
can delegate without accepting simultaneous in-process calls. Explicit dispatch
fields make each claim independently testable.

### Keep CPU-derived automatic concurrency

Rejected. Local CPU count is neither a provider quota nor an agent-worker count.
Transport defaults are conservative, deterministic, and can evolve with
adapter evidence without pretending host compute is the limiting resource.

### Keep provider execution in checkpoint waves

Rejected. Reusing harness resume logic is compact, but waiting for every request
in a wave creates a barrier and the harness-shaped bootstrap stages state before
the actual provider cap is known. Shared request/acceptance functions retain the
useful reuse without sharing the transport driver.

### Warm one provider call before opening the pool

Deferred. A warm-up barrier might improve prefix-cache reuse or might add pure
latency when provider caches are already warm. Change 0203 supplies telemetry;
a cache-aware scheduler needs separate evidence and should not weaken the
completion-driven baseline by assumption.

## Trade-offs and risks

- A completion-driven coordinator is more complex than `Effect.forEach` over a
  batch. The design contains that complexity in one application component and
  protects it with controllable worker tests.
- Per-result writes increase filesystem activity, but evaluation calls are much
  slower than an atomic local JSON write and durability is worth the cost.
- Automatic cap `4` is still a policy choice rather than live quota discovery.
  It is intentionally conservative, visible, and overrideable downward or
  upward through the existing workspace cap.
- Direct concurrency can increase rate-limit and token-cost pressure. Existing
  classified retries and explicit configuration remain the controls; adaptive
  control is deferred.
- Capability and artifact shapes change cleanly in early alpha. In-flight
  version-8 runs must restart, and all current consumers must advance together.
- Agent harnesses may service fewer workers than the outstanding cap. That is
  valid; honest progress wording prevents eligibility from being mistaken for
  actual parallelism.

## Open questions

None. Adaptive, cache-aware, and dynamically negotiated capacity remain
explicitly deferred rather than unresolved inside this design.
