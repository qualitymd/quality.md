---
type: Functional Specification
title: evaluation.json
description: Artifact contract for the authoritative run artifact written by the evaluation runner.
tags: [evaluation, records, artifact]
timestamp: 2026-07-11T00:00:00Z
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
  "schemaVersion": 7,
  "kind": "EvaluationRun",
  "manifest": {},
  "state": {},
  "sources": {},
  "results": {},
  "outputs": {}
}
```

`schemaVersion` **MUST** be `7` and is a payload-shape marker only; versions
1–3 belong to the historical multi-file data tree, version 4 belongs to the
strategy-named runner artifact, version 5 predates the per-area `sources`
record, and version 6 predates the multi-outstanding harness checkpoint
window (its singular pending call is not read back — an awaiting version-6
run is re-run, not migrated). `kind` **MUST** be `EvaluationRun`.

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
- `pendingEvaluatorCalls` — the awaiting harness checkpoint's correlation
  metadata (request identity, work-unit identity, input hash, correlation ID,
  and attempt) for every outstanding work request, in emission order, bounded
  by the manifest's resolved `concurrency`, present exactly while harness
  work requests await their results; and
- `harnessIdentity` — the harness runtime the run's judgment is bound to, set
  by the first accepted harness result.

`awaiting_evaluator` marks a resumable, incomplete checkpoint — never a
failure. `pendingEvaluatorCalls` is the single source of truth for what is
awaiting and **MUST NOT** carry raw prompt, source, or result bodies: each
pending request is rebuilt from the model snapshot, work-graph state, and
current source package per the
[runner harness contract](runner.md#harness-checkpoints).

`state` **MUST NOT** carry provider context identifiers or prompt-cache
status; those live only in run-local logs.

> Rationale: provider-retained identifiers expire outside the run's control.
> Keeping them out of the authoritative artifact keeps resume honest about what
> is reconstructible. — 0192

## Sources

`sources` is the per-area source provenance of record, keyed by area
reference. Each record **MUST** carry the area's effective `selector`, its
detected `kind` (`path`, `glob`, or `prose`), and the `resolver` that serves
it (`walk` for the deterministic filesystem walk, `harness` for
checkpoint-dispatched resolution), written at run creation per the
[runner detection contract](runner.md#selector-kind-detection). As each
area's bundle materializes, the record **MUST** be completed with the
`bundleHash`, `capturedAt`, a bundle `truncated` mark when a cap applied, and
per-file entries carrying `path`, `sha256`, and a `truncated` mark; a
harness-resolved record additionally carries the `harnessRuntime` that served
resolution.

```json
"sources": {
  "area:api": {
    "selector": "open tickets in the support queue",
    "kind": "prose",
    "resolver": "harness",
    "harnessRuntime": "claude-code",
    "bundleHash": "…",
    "capturedAt": "2026-07-11T00:00:00Z",
    "files": [{ "path": "…", "sha256": "…", "content": "…" }]
  }
}
```

Walked (path/glob) records **MUST NOT** carry file `content` — their material
is re-readable from the workspace, and resume re-packages it. Harness-resolved
records **MUST** carry file `content`: the captured bundle is the evidence of
record, and resume **MUST** rebuild dependent requests from it rather than
re-gather. Captured prose file paths are unique workspace-relative paths for
the gathered files, recorded and hashed verbatim. Absolute paths, URLs,
external identifiers, and paths that escape the workspace are invalid.

> Rationale: one record serves kind pinning, resume for harness-resolved
> areas, and audit provenance — a reviewer reads the same shape for walked and
> agent-gathered evidence and can tell which is which. Keeping captured
> content in the one authoritative artifact preserves the store's atomic
> write; the 512 KB bundle cap bounds the growth. Gathered material lands in
> the artifact verbatim — the same class of exposure as quoted evidence in
> reports. — 0197

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

Changing the implementation runtime does not by itself invalidate a current
artifact. A conforming runtime **MUST** accept every current-schema completed or
resumable run accepted by the immediately prior release. A real artifact-shape
change requires one explicit schema-version break and refusal of incompatible
resume, not a dual reader or migration shim.

For a harness-backed run, resume compatibility additionally covers the
pending requests and identity: a pending request whose rebuilt input hash no
longer matches its `pendingEvaluatorCalls` entry, or a result from a runtime
other than the bound `harnessIdentity`, **MUST** fail with
`run_state_invalid`.

Work-unit reuse on resume is decided by the
[orchestration resume rules](orchestration.md#resume).

## Consumers

`qualitymd evaluation status` and `qualitymd evaluation list` **MUST** read
`evaluation.json` for runner-created runs and derive reportability from the
same coverage rules as report build.

`qualitymd evaluation report build` **MUST** re-render reports from
`evaluation.json` for runner-created runs; generated Markdown is never
report-generation input, per the [report tree](reports/report-tree.md).
