---
type: Functional Specification
title: evaluation.json
description: Artifact contract for the authoritative run artifact written by the evaluation runner.
tags: [evaluation, records, artifact]
timestamp: 2026-07-09T00:00:00Z
---

# evaluation.json

`evaluation.json` is the authoritative structured run artifact the
[evaluation runner](runner.md) writes at the evaluation run root. It carries
everything needed to validate, resume, review, and render the run. Run-local
logs stay in separate files per the
[runner logging contract](runner.md#logging); generated Markdown reports are
projections over this artifact per the [report tree](reports/report-tree.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Location and envelope

New runner-created evaluations **MUST** write one authoritative structured run
artifact at `<run>/evaluation.json`. The run folder layout around it is defined
by the [data layout](records/data-layout.md).

The document envelope is:

```json
{
  "schemaVersion": 5,
  "kind": "EvaluationRun",
  "manifest": {},
  "state": {},
  "results": {},
  "outputs": {}
}
```

`schemaVersion` **MUST** be `5` and is a payload-shape marker only; versions
1–3 belong to the historical multi-file data tree, and version 4 belongs to the
strategy-named runner artifact. `kind` **MUST** be `EvaluationRun`.

## Manifest

`manifest` carries the run's immutable identity and setup. It **MUST** include
`evaluationId`, `createdAt`, `model`, `requestedScope`, `plannedScope`, `run`
(the local run `number` and folder `label`), `evaluator` (the selected
evaluator or profile name), `evaluatorKind` (the selected evaluator's runtime
kind), and `concurrency` (the resolved concurrency cap).

## State

`state` carries the run's execution lifecycle. It **MUST** include:

- the run `status`, one of `running`, `awaiting_evaluator`, `completed`,
  `failed`, or `cancelled`;
- the classified `failure` when the run failed;
- per-work-unit entries carrying status, attempts, input hash, failure, and
  timestamps;
- a `cancelled` marker when a user interruption was observed mid-run;
- `pendingEvaluatorCall` — the awaiting harness checkpoint's correlation
  metadata (request identity, work-unit identity, input hash, correlation ID,
  and attempt), present exactly while a harness work request awaits its
  result; and
- `harnessIdentity` — the harness runtime the run's judgment is bound to, set
  by the first accepted harness result.

`awaiting_evaluator` marks a resumable, incomplete checkpoint — never a
failure. `pendingEvaluatorCall` **MUST NOT** carry raw prompt, source, or
result bodies: the pending request is rebuilt from the model snapshot,
work-graph state, and current source package per the
[runner harness contract](runner.md#harness-checkpoints).

`state` **MUST NOT** carry provider context identifiers or prompt-cache
status; those live only in run-local logs.

> Rationale: provider-retained identifiers expire outside the run's control.
> Keeping them out of the authoritative artifact keeps resume honest about what
> is reconstructible. — 0192

## Results

`results` **MUST** carry the accepted routine payloads in deterministic
work-graph order, each attributed to the work unit that produced it. The
payload kinds and field contracts are the shared vocabulary defined in
[payload kinds](records/payload-kinds.md).

## Outputs

`outputs` **MUST** be written at report build and include `reportMd` (the
generated run-level report reference), the CLI-owned `EvaluationOutputResult`
payload, and the scoped area rating carried into command receipts.

## Persistence

The runner **MUST** write `evaluation.json` atomically — a temp file plus
rename — when updating persisted run state. Writes happen after every accepted
work-unit result, before the unit counts as complete, per the
[orchestration persistence rules](orchestration.md#persistence).

While parallel execution is active, `evaluation.json` writes **MUST** be
serialized by a single store so concurrent merges cannot interleave.

## Resume compatibility

A run artifact is resume-compatible when `evaluation.json` is present at the
run root, its `schemaVersion` is supported by the running `qualitymd` version,
and its manifest `model` resolves to the selected model. If compatibility
verification fails, then resume **MUST** fail with `run_state_invalid` and
report that starting a new run is the remedy.

For a harness-backed run, resume compatibility additionally covers the pending
request and identity: a pending request whose rebuilt input hash no longer
matches `pendingEvaluatorCall`, or a result from a runtime other than the
bound `harnessIdentity`, **MUST** fail with `run_state_invalid`.

Work-unit reuse on resume is decided by the
[orchestration resume rules](orchestration.md#resume).

## Consumers

`qualitymd evaluation status` and `qualitymd evaluation list` **MUST** read
`evaluation.json` for runner-created runs and derive reportability from the
same coverage rules as report build.

`qualitymd evaluation report build` **MUST** re-render reports from
`evaluation.json` for runner-created runs; generated Markdown is never
report-generation input, per the [report tree](reports/report-tree.md).
