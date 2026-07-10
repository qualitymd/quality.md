---
type: Functional Specification
title: Evaluation concurrency
description: Requirements for replacing execution strategy with a concurrency setting and running dependency-ready evaluation steps concurrently.
tags: [evaluation, runner, concurrency]
timestamp: 2026-07-10T00:00:00Z
---

# Evaluation concurrency

This spec governs the delta from strategy-named evaluation execution to
concurrency-bounded evaluation execution. It inherits the binding work graph,
result-validation, persistence, retry, resume, and report semantics from the
durable [evaluation runner](../../../specs/evaluation/runner.md),
[orchestration contract](../../../specs/evaluation/orchestration.md), and
[evaluator contract](../../../specs/evaluation/evaluator-contract.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

Evaluation work units often wait on multi-second evaluator calls. Running every
ready unit sequentially leaves independent requirement assessments and analyses
idle even though the runner already has a dependency graph that can identify
safe concurrent work. At the same time, the public `executionStrategy` enum makes
users choose an implementation strategy instead of the control they actually
need: how many evaluator calls may run at once.

This change moves the user-facing contract to `concurrency`, keeps omitted
configuration automatic and CPU-correlated, and preserves the runner's existing
deterministic artifact and report guarantees.

## Scope

This change covers workspace configuration, runner dry-run and run receipts,
`evaluation.json` state, evaluator concurrency capability, dependency-aware
scheduling, cancellation, retry, persistence, and durable skill/spec guidance.

It defers multi-checkpoint harness concurrency, public fallback/adjustment
records, evaluator-specific maximums, and adaptive rate-limit control. The
existing retry policy remains the only rate-limit behavior changed by this case.

## Requirements

### Configuration and receipts

The workspace config **MUST** replace `evaluation.executionStrategy` with an
optional positive integer `evaluation.concurrency`.

> Rationale: concurrency is the behavior users can reason about; strategy names
> expose scheduler internals and duplicate the meaning of `1` versus `>1`.
>
> Durable spec: modify `specs/evaluation/runner.md`,
> `specs/evaluation/evaluator-contract.md`, and
> `specs/cli/evaluation-run.md` to make concurrency the public configuration
> surface and remove execution-strategy guidance.

When `evaluation.concurrency` is absent, the runner **MUST** request automatic
concurrency equal to `max(2, runtime.NumCPU()*2)`.

> Rationale: evaluator calls are multi-second external work, so the automatic
> default should be more aggressive than CPU-bound worker defaults while still
> tracking the host machine.
>
> Durable spec: modify `specs/evaluation/runner.md` to define the automatic
> concurrency default.

If `evaluation.concurrency` is present but not a positive integer, then
`qualitymd evaluation run` and `qualitymd evaluation run --dry-run` **MUST**
fail before creating or mutating an evaluation run.

> Durable spec: modify `specs/evaluation/runner.md` and
> `specs/cli/evaluation-run.md` to define invalid concurrency handling.

Dry-run JSON, run receipts, and `evaluation.json` state **MUST** report the
resolved `concurrency` value and **MUST NOT** expose `executionStrategy`,
`strategyFallbacks`, or requested-versus-resolved concurrency fields.

> Rationale: the user asked to keep the surface to a single concurrency value
> for now; downgrade explanation can stay out of the public contract until a
> concrete need appears.
>
> Durable spec: modify `specs/cli/evaluation-run.md` and
> `specs/evaluation/evaluation-json.md` to remove public strategy fields and
> retain the resolved concurrency field.

### Scheduling semantics

The runner **MUST** treat a resolved concurrency of `1` as sequential execution
and a resolved concurrency greater than `1` as permission to execute
dependency-ready evaluation steps concurrently up to that limit.

> Durable spec: modify `specs/evaluation/runner.md` and
> `specs/evaluation/orchestration.md` to define concurrency-bounded scheduling.

Concurrent execution **MUST** be observationally equivalent to deterministic
sequential execution: it must not change accepted payload ordering, report
content, rating roll-up, artifact paths, result schemas, or resume semantics.

> Rationale: concurrency is a scheduling optimization, not a new evaluation
> method.
>
> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/evaluation-json.md` to preserve deterministic outputs under
> concurrent execution.

The runner **MUST** persist accepted evaluator-backed results before treating
their evaluation step as complete for dependency scheduling or resume.

> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/runner.md` to preserve the existing durability boundary
> under concurrent execution.

Evaluator workers **MUST NOT** write evaluation run artifacts directly.

> Rationale: a single runner-owned persistence path keeps concurrent completion
> order from changing durable output order or corrupting the run artifact.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/evaluation/orchestration.md` to preserve evaluator boundaries under
> concurrent scheduling.

### Evaluator support and harness boundary

The runner **MUST** resolve concurrency against the selected evaluator's
declared support for concurrent calls before execution begins.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/evaluation/runner.md` to define evaluator concurrency support.

Harness-backed evaluation **MUST** resolve to `concurrency: 1` until the durable
runner and harness contracts define multiple pending evaluator checkpoints.

> Rationale: the current harness protocol has one pending request and one result
> submission path; widening it would be a separate checkpoint and UX contract.
>
> Durable spec: modify `specs/evaluation/runner.md` and
> `specs/evaluation/orchestration.md` to keep harness execution sequential.

### Failure, cancellation, and logging

When any evaluation step fails terminally under concurrent execution, the runner
**MUST** cancel unscheduled work, cancel in-flight evaluator calls through the
run context, persist the classified failure, and return the same failed or
cancelled run status that sequential execution would return.

> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/runner.md` to define terminal failure and cancellation under
> concurrent execution.

Run-local logs **MUST** remain safe under concurrent execution and **MUST NOT**
record raw prompts, raw source bundles, raw model responses, credentials, or
environment values.

> Durable spec: modify `specs/evaluation/runner.md` to preserve logging
> boundaries under concurrent execution.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/runner.md` — replace execution strategy with concurrency
  resolution, automatic default, evaluator capability resolution, durability,
  logging, and harness sequential behavior (per the configuration, scheduling,
  evaluator, and failure requirements).
- `specs/evaluation/orchestration.md` — define concurrency-bounded scheduling,
  observational equivalence, single persistence owner, harness sequential
  boundary, and cancellation behavior (per the scheduling, harness, and failure
  requirements).
- `specs/evaluation/evaluator-contract.md` — replace strategy capability with
  concurrent-call support and keep evaluators out of artifact writes (per the
  configuration and evaluator-boundary requirements).
- `specs/evaluation/evaluation-json.md` — remove public execution strategy and
  strategy fallback state while preserving resolved concurrency and deterministic
  result ordering (per the receipt and scheduling requirements).
- `specs/cli/evaluation-run.md` — remove execution-strategy receipt fields and
  document concurrency in dry-run and run JSON receipts (per the configuration
  and receipt requirements).

### To rename

None

### To delete

None
