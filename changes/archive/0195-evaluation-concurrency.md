---
type: Change Case
title: Evaluation concurrency
description: Replace evaluation execution strategies with a single concurrency setting and run dependency-ready evaluator steps concurrently under runner-owned coordination.
status: Done
tags: [evaluation, runner, concurrency]
timestamp: 2026-07-10T00:00:00Z
---

# Evaluation concurrency

## Motivation

The deterministic evaluation runner already models the evaluation as a work
graph, but it executes that graph one unit at a time. That makes full
evaluations slower than necessary because independent evaluator-backed steps can
spend many seconds waiting on model or CLI work while other ready steps sit idle.

The current `executionStrategy` setting also exposes an implementation choice
that users do not need to reason about. The real user-facing control is how many
evaluator calls the runner may keep active. A single `concurrency` setting can
express sequential and concurrent execution directly while keeping scheduling,
validation, persistence, and report output under the runner's deterministic
contract.

## Scope

Covered:

- replace workspace `evaluation.executionStrategy` with
  `evaluation.concurrency`;
- resolve omitted concurrency to an automatic default of
  `max(2, runtime.NumCPU()*2)`;
- preserve `concurrency: 1` as sequential execution and allow higher positive
  integers to run dependency-ready evaluation steps concurrently;
- keep the public run artifact, receipts, and dry-run output centered on the
  resolved `concurrency` value, without public strategy or fallback fields;
- implement concurrent scheduling with a coordinator that owns artifact mutation
  and deterministic persistence; and
- keep harness-backed evaluation sequential until multi-checkpoint harness work
  is designed.

Deferred:

- multi-checkpoint or batched harness execution;
- evaluator-specific public maximums or fallback records;
- adaptive rate-limit backoff beyond the existing per-step retry policy; and
- removing the sequential execution path as an internal baseline.

Non-goals:

- changing evaluation judgment semantics, accepted payload schemas, report
  content, rating roll-up, or run folder layout;
- allowing evaluator workers to write run artifacts directly; or
- adding a command-line `--concurrency` flag in this slice.

## Affected artifacts

- **Code:** `internal/workspace/workspace.go` and tests for
  `evaluation.concurrency`; `internal/evaluator/` for concurrency capability
  declarations; `internal/runner/` for concurrency resolution, receipts,
  scheduler/coordinator execution, artifact state, logs, dry-run output, and
  tests; `internal/cli/evaluation_run.go` and tests for changed JSON/human
  output; generated CLI docs if command output changes.
- **Durable specs:** modify `specs/evaluation/runner.md`,
  `specs/evaluation/orchestration.md`,
  `specs/evaluation/evaluator-contract.md`,
  `specs/evaluation/evaluation-json.md`, and
  `specs/cli/evaluation-run.md`. See the functional spec's
  [durable spec changes](0195-evaluation-concurrency/spec.md#durable-spec-changes).
- **Bundled skill:** update `skills/quality/SKILL.md`,
  `skills/quality/workflows/evaluate.md`, and
  `skills/quality/resources/cli-workflow-conventions.md` so examples and
  guidance use `evaluation.concurrency`, not `executionStrategy`; update their
  local logs/indexes where required.
- **Durable docs:** update `install.md` or other user-facing docs only if they
  mention evaluation runner strategy configuration; regenerate generated CLI or
  Mintlify docs if their sources change.
- **Release record:** add the user-visible concurrency behavior to
  `CHANGELOG.md` during implementation/release preparation.
- **Format specification and scaffold:** no impact. The QUALITY.md model schema
  and starter model content do not change.
- **Generated evaluation reports:** no content or layout change; report-gallery
  artifacts should remain stable for unchanged evaluation data.

## Children

- [Functional spec](0195-evaluation-concurrency/spec.md)
- [Design doc](0195-evaluation-concurrency/design.md)

## Status

`Done`. Implemented and landed: `evaluation.concurrency` replaces
`evaluation.executionStrategy`, omitted concurrency resolves to
`max(2, runtime.NumCPU()*2)`, directly callable evaluators can run
dependency-ready judgment steps concurrently under a coordinator-owned
scheduler, harness-backed runs resolve to `concurrency: 1`, public receipts and
`evaluation.json` expose only resolved `concurrency`, and runner artifacts use
schema version `5`. Verified with the full check suite and
`go test -race ./internal/runner`.
