---
type: Design Doc
title: Evaluation concurrency design
description: Coordinator-owned scheduling design for concurrency-bounded evaluation runs.
tags: [evaluation, runner, concurrency]
timestamp: 2026-07-10T00:00:00Z
---

# Evaluation concurrency design

## Context

This design answers the [evaluation concurrency functional spec](spec.md). The
runner already builds a dependency graph in deterministic model order and
persists accepted results into one authoritative `evaluation.json`. The missing
piece is a scheduler that can run independent evaluator-backed steps while
keeping artifact mutation, logs, and report output deterministic.

## Approach

### Concurrency as the only public control

Workspace config gains:

```yaml
evaluation:
  concurrency: 8
```

The resolver removes `executionStrategy` from the public path. If the config
omits concurrency, it computes `max(2, runtime.NumCPU()*2)`. If the selected
evaluator cannot run concurrently, the resolved value becomes `1`. Receipts and
`evaluation.json` record only the resolved `concurrency`.

Evaluator capabilities become a direct statement about concurrent calls, not a
strategy list. A boolean is enough for this slice because the user chose not to
add public fallback or maximum fields. API and CLI evaluators can declare
concurrent-call support; the harness evaluator declares no concurrent-call
support until the runner can persist multiple pending checkpoints.

### Coordinator-owned execution

The concurrent path uses a coordinator plus workers. The coordinator owns all
mutable runner state:

- dependency readiness;
- accepted result ordering;
- `evaluation.json` mutation and saves;
- run status and terminal failure;
- cancellation decisions; and
- report generation.

Workers execute one scheduled evaluation step and return a completion message.
Workers do not call `acceptUnit`, mutate `evaluation.json`, build reports, or
write run artifacts.

The scheduler-facing message uses evaluation language rather than generic work
language:

```go
type EvaluationStepCompletion struct {
    Step     *Unit
    Accepted *AcceptedEvaluationStep
    Rejected *RejectedEvaluationStep
    Err      error
}

type AcceptedEvaluationStep struct {
    Payloads  []map[string]any
    InputHash string
}

type RejectedEvaluationStep struct {
    Failure *Failure
}
```

The coordinator receives completions, accepts and persists them in deterministic
graph order, and then marks dependents ready. If a later step finishes before an
earlier graph-order step, its completion can wait in memory until all earlier
completed-or-skipped steps that affect output order are accepted.

### Step categories

The engine keeps the current distinction between:

- frame steps: deterministic context payloads;
- judgment steps: evaluator-backed assessment, analysis, ranking, and
  recommendation calls; and
- report step: deterministic final assembly.

Frame and report steps stay coordinator-executed. Judgment steps are the main
parallel target. This keeps deterministic state-building and report assembly
simple while gaining concurrency where the run spends time.

### Sequential baseline

The existing sequential loop remains as the execution path for
`concurrency == 1`, harness-backed runs, and as the correctness baseline for
tests. The concurrent path reuses the same lower-level request building,
validation, retry, logging, and acceptance helpers, but splits judgment
execution from artifact acceptance so the coordinator is the only writer.

### Cancellation and failure

Concurrent execution derives a child context from the run context. A terminal
failure cancels that child context, stops scheduling new steps, lets in-flight
workers return promptly through evaluator cancellation, persists the classified
failure through the coordinator, and returns the same status shape as the
sequential path.

User cancellation keeps the existing `cancelled` behavior: accepted results are
already durable, incomplete steps remain resumable, and the run can continue
with `--resume`.

## Spec response

The config and receipt requirements are satisfied by replacing
`EvaluationConfig.ExecutionStrategy` with `Concurrency`, removing strategy
fields from receipts and artifacts, and resolving omitted config before run
creation.

The scheduling requirements are satisfied by treating the graph as the source
of readiness and keeping all persistence in the coordinator. The coordinator's
deterministic acceptance order preserves result ordering and report output even
when evaluator calls finish out of order.

The evaluator-boundary requirements are satisfied by capability resolution and
by making worker completions in-memory values only. Harness remains sequential
because its persisted state has one `PendingEvaluatorCall`.

## Alternatives

### Keep `executionStrategy`

This was rejected because the enum duplicates what `concurrency` already says:
`1` is sequential; values above `1` permit concurrent scheduling. Keeping both
would create precedence questions and a larger public contract.

### Use `errgroup.SetLimit` over the graph slice

This was rejected as the main scheduler because the graph is not a flat set of
independent tasks. A coordinator can react to newly satisfied dependencies and
own persistence in one place. `errgroup` remains useful internally only for
worker lifecycle if it keeps the code simpler.

### Let workers write accepted results directly

This was rejected because it spreads artifact mutation across goroutines,
requires broad locking, and makes deterministic result ordering harder to
review. A single writer matches the runner's existing durability contract.

### Default to a small fixed value

This was rejected because evaluator calls are multi-second external work. The
default should scale with the host and be useful for ordinary full-model runs;
explicit config remains available for users who need less or more fanout.

## Trade-offs and risks

- `runtime.NumCPU()*2` can launch many evaluator calls on large machines when
  the user omits concurrency. This is intentional for this slice, but users can
  set `concurrency: 1` or a smaller value when cost or rate limits matter.
- CLI evaluators may have local subprocess or authentication contention under
  high concurrency. Tests should cover scheduler behavior with fake evaluators;
  live CLI behavior may need follow-up tuning.
- Holding completed out-of-order results in memory is acceptable because
  individual accepted payloads are small compared with source bundles and model
  responses, and the coordinator persists as soon as graph-order acceptance
  permits.
- Keeping the sequential path initially duplicates some control flow, but it
  gives a clear baseline while concurrent scheduling matures.

## Open questions

None for this implementation slice.
