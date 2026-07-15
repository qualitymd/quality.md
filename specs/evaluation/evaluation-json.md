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
  "schemaVersion": 9,
  "kind": "EvaluationRun",
  "manifest": {},
  "state": {},
  "evidence": {},
  "results": {},
  "outputs": {}
}
```

`schemaVersion` **MUST** be `9` and is a payload-shape marker only; versions
1–3 belong to the historical multi-file data tree, version 4 belongs to the
strategy-named runner artifact, version 5 predates the per-area `sources`
record, and version 6 predates the multi-outstanding harness checkpoint
window. Version 7 contains runner-resolved per-area source bundles; it is not
read back or migrated. Version 8 predates structured direct/delegated dispatch
capabilities and selected-evaluator-first creation; an in-flight version-8 run
is not read back or migrated. Completed older artifacts remain historical
records. `kind` **MUST** be `EvaluationRun`.

## Manifest

`manifest` carries the run's immutable identity and setup. It **MUST** include
`evaluationId`, `createdAt`, `model`, `requestedScope`, `plannedScope`, `run`
(the local run `number` and folder `label`), `evaluator` (the selected
evaluator or profile name), `evaluatorKind` (the selected evaluator's runtime
kind), `evaluatorCapabilities` (including the structured dispatch capability),
`concurrency` (the resolved concurrency cap), and `areaSources` (each planned
area's effective selector and detected path, glob, or prose form).

The runner **MUST** select and validate the evaluator and resolve concurrency
before writing this manifest or staging requests. Evaluator name, kind,
capabilities, and concurrency are immutable after request issue; a provider run
**MUST NOT** be created as harness-backed and patched afterward.

## State

`state` carries the run's execution lifecycle. It **MUST** include:

- the run `status`, one of `running`, `awaiting_evaluator`, `completed`,
  `failed`, or `cancelled`;
- the classified `failure` when the run failed;
- per-work-unit entries carrying status, attempts, input hash, failure,
  timestamps, and the accepted evidence manifest hash for requirement units;
- a `cancelled` marker when a user interruption was observed mid-run;
- `pendingEvaluatorCalls` — correlation metadata (request identity, work-unit
  identity, input hash, correlation ID, and attempt) for every issued but
  unaccepted evaluator request, in emission order and bounded by the manifest's
  resolved `concurrency`; for direct transports it is the resumable boundary
  around active or completed-but-unaccepted calls, and for harness transport it
  is the outstanding checkpoint window; and
- `harnessIdentity` — the harness runtime the run's judgment is bound to, set
  by the first accepted harness result.

`awaiting_evaluator` marks a resumable harness checkpoint — never a failure.
Direct dispatch uses `running`; cancellation may persist `cancelled` while
leaving unaccepted pending calls resumable. `pendingEvaluatorCalls` is the
single source of truth for issued but unaccepted work and **MUST NOT** carry raw
prompt, source, or result bodies: each request is rebuilt from the model
snapshot, work-graph state, effective source selector, and inspection policy.
At every persisted boundary its length **MUST NOT** exceed manifest
`concurrency`.

`state` **MUST NOT** carry provider context identifiers or prompt-cache
status; those live only in run-local logs.

> Rationale: provider-retained identifiers expire outside the run's control.
> Keeping them out of the authoritative artifact keeps resume honest about what
> is reconstructible. — 0192

## Evidence

`evidence` is the accepted per-requirement evidence provenance, keyed by the
`assessRateRequirement` work-unit ID. Each record **MUST** carry the requirement
ID, effective source selector and form, ordered observations, evaluation limits,
capture time, and canonical manifest hash.

```json
"evidence": {
  "assessRateRequirement:requirement:api::authentication": {
    "requirementId": "requirement:api::authentication",
    "source": { "selector": "src/api", "kind": "path" },
    "observations": [
      {
        "id": "ev-001",
        "kind": "file",
        "role": "evaluated",
        "path": "src/api/auth.ts",
        "locator": { "startLine": 18, "endLine": 61 },
        "sha256": "…",
        "bytes": 2048,
        "capturedAt": "2026-07-14T00:00:00Z"
      },
      {
        "id": "ev-002",
        "kind": "file",
        "role": "supporting",
        "path": "docs/authentication.md",
        "locator": { "heading": "API authentication" },
        "sha256": "…",
        "bytes": 1024,
        "capturedAt": "2026-07-14T00:00:00Z"
      }
    ],
    "limits": [],
    "capturedAt": "2026-07-14T00:00:00Z",
    "manifestHash": "…"
  }
}
```

File paths **MUST** be workspace-relative regular UTF-8 files and contained
after real-path resolution. `evaluated` path and glob observations **MUST**
belong to the selected subject; `supporting` observations may be elsewhere in
the authorized workspace. Locators use a valid line range or Markdown heading.
Each file carries runner-computed bytes and SHA-256. The artifact **MUST NOT**
carry file bodies or tool transcripts.

The runner **MUST** seal and persist a requirement manifest atomically with its
assessment and rating. Every assessment evidence `sourceRef` **MUST** name an
observation in that manifest. Accepted manifests are immutable resume inputs;
the runner does not regather them.

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
serialized by a single store so concurrent merges cannot interleave. One
accepted result **MUST** be persisted before its slot is released or dependent
work is dispatched; persisted payload arrays **MUST** be reprojected in graph
order after every merge.

## Resume compatibility

A run artifact is resume-compatible when `evaluation.json` is present at the
run root, its `schemaVersion` is supported by the running `qualitymd` version,
and its manifest `model` resolves to the selected model. If compatibility
verification fails, then resume **MUST** fail with `run_state_invalid` and
report that starting a new run is the remedy.

Changing the implementation runtime does not by itself invalidate a current
artifact. A real artifact-shape change requires one explicit schema-version
break and refusal of incompatible resume, not a dual reader or migration shim.
The schema-9 runtime **MUST** refuse an in-flight schema-8 run with the existing
start-new-run remedy; it **MUST NOT** add a version-8 reader or migration path.

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
