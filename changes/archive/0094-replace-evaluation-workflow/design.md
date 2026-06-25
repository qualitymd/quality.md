---
type: Design Doc
title: Replace evaluation workflow - design
description: Implementation approach for replacing the current evaluation workflow with Evaluation v2.
tags: [evaluation, workflow, cli, skill, reports]
timestamp: 2026-06-25T00:00:00Z
---

# Replace evaluation workflow - design

Companion to the [functional spec](spec.md). This design records how the
Evaluation v2 replacement should be built.

## Context

The current evaluation implementation is organized around run artifacts:
assessment records, analysis records, recommendation records, and generated
summary/full/JSON reports. Evaluation v2 changes the organizing principle. The
agent follows a protocol of judgment routines, persists each routine output as
structured JSON under `data/`, and asks the CLI to assemble deterministic reports
from those completed outputs.

The implementation needs to replace rather than extend the old workflow because
the artifact model changes at every layer: command names, run layout, report
shape, skill instructions, and durable specs.

## Approach

### Replacement strategy

Build Evaluation v2 as the new evaluation implementation, not as a compatibility
layer over the existing assessment/analysis/recommendation record model.

The durable specs should move to a new `specs/evaluation-v2/` parent folder that
owns the protocol, routines, JSON records, data layout, and reports. The existing
evaluation-record and report specs should either be removed when fully replaced
or reduced to historical/compatibility notes if a reader still needs to
understand old runs.

The CLI should expose the new command surface directly:

```text
qualitymd evaluation create [model]
qualitymd evaluation data set <run> --file <path|->
qualitymd evaluation report build <run>
```

Existing commands that write old records should be removed or redirected only if
the redirect can give a targeted message. They should not silently write old
records into a v2 run.

### Responsibility split

Keep the existing product posture: the skill and agent own judgment, while the
CLI owns mechanics.

The agent:

- resolves the user's requested evaluation intent through the skill workflow;
- creates frames;
- inspects evidence;
- writes structured routine payloads;
- applies stop conditions and confidence judgment; and
- decides Factor and Area synthesis from structured lower-level outputs.

The CLI:

- creates and identifies run folders;
- validates structured payloads;
- derives canonical paths;
- writes canonical JSON;
- reports present/missing/invalid data;
- assembles `EvaluationOutputResult`; and
- renders reports deterministically.

This keeps model judgment out of the CLI while still making the run folder
auditable and resumable.

### Run and data model

Use a v2 run layout with two top-level concerns:

- human-readable reports rooted at `report.md`; and
- structured routine data under `data/`.

The `data/` tree mirrors the evaluated Area tree. Inside each Area, local
Requirements and local Factors own their own routine output files. Factor folders
nest under Factor folders because Factor identity is Area-scoped and recursive.

The data layer should model payloads as direct JSON documents with
`schemaVersion` and `kind`, not as records wrapped in a shared envelope. The
write path can still maintain internal metadata while deriving the canonical path
from `kind` plus the structured subject ID.

Path derivation should be centralized in one package so `data set`, `data get`,
`status`, and `report build` agree on where each payload lives.

### Payload validation

Implement validation in layers:

1. Common payload checks: parseable JSON, supported `schemaVersion`, known
   `kind`, required top-level fields, no impossible optional/null states.
2. Kind-specific structural checks: required fields, enum values, subject shape,
   scoped status/rating combinations, and required refs.
3. Model-aware checks: Area, Factor, Requirement, and Rating Level IDs resolve
   against the run's model snapshot.
4. Graph checks: referenced routine outputs exist when a command needs a
   reportable run.

`evaluation data set --dry-run` should run the same validation and path
derivation as a real write, then report the intended write without mutating the
run.

### Status and report build

Status and failed report build should share one gap collector. The collector
should inspect the model snapshot and the expected v2 routine output graph, then
return typed gaps such as missing payload, malformed payload, unsupported schema
version, unresolved model ID, missing dependency, or non-reportable routine
status.

`evaluation report build` should run the gap collector first. If report-blocking
gaps exist, it should fail before writing reports. If the run is reportable, it
should assemble `data/evaluation-output-result.json` and then render the report
tree.

Render the complete output plan in memory before writing generated files. The v0
implementation may write the generated report files directly after successful
rendering; a later hardening pass can add a staging tree if users need
all-or-nothing replacement semantics for very large report trees.

### Report rendering

Treat Markdown reports as pure projections over persisted structured data.

Report rendering should use deterministic structural ID order for Areas,
Factors, and Requirements, and preserve source order for rating drivers,
findings, and evidence. It should resolve display titles from the model snapshot
but keep data links to the JSON files that back each report section. Source-order
model traversal can be added later if the model layer exposes an order-preserving
view of the parsed QUALITY.md frontmatter.

The renderer should be split by report subject:

- Area report renderer;
- Factor report renderer;
- Requirement report renderer; and
- shared table/link/breadcrumb helpers.

Those helpers should enforce empty-state rows, status rendering, secret
redaction boundaries, and path/link generation consistently.

### Skill runtime

Update the `/quality evaluate` workflow to drive the new CLI shape.

The runtime skill instructions should not duplicate every durable record field in
`SKILL.md`. Instead, keep `SKILL.md` as the high-level operating contract and
place Evaluation v2 runtime instructions in workflow files. If the routine
prompts become large, mirror the durable spec shape under an
`evaluation-v2/` workflow folder:

```text
skills/quality/workflows/evaluation-v2.md
skills/quality/workflows/evaluation-v2/
  protocol.md
  routines/
    frame-requirement-evaluation.md
    assess-requirement.md
    ...
```

The runtime instructions should be specific enough for an agent to produce valid
payloads without memorizing the CLI implementation, and should route mechanical
discovery to `qualitymd evaluation data kinds` and
`qualitymd evaluation data example <kind>`.

### Orchestration

Represent the protocol as a dependency graph even if the first implementation is
mostly sequential.

A sequential runner can follow the same graph:

1. frame evaluation;
2. walk Areas recursively;
3. frame local Requirements;
4. assess and rate local Requirements;
5. analyze local Factor trees bottom-up;
6. wait for child Area analyses;
7. frame and analyze the Area;
8. build reports.

Parallel agent runtimes can execute ready work units concurrently as long as
they write the same payloads to the same derived paths and preserve deterministic
ordering in stored arrays and reports.

The implementation does not need a subagent manager in the CLI. The protocol
only needs stable work-unit boundaries and persisted data contracts.

## Alternatives

### Extend the current evaluation records

Extending the current assessment, analysis, recommendation, `report.json`, and
`report-summary.md` model was rejected. The v2 protocol changes the core unit of
work from numbered record types to routine outputs. Keeping the old model would
force new concepts such as frames, Requirement ratings, Factor frames, and Area
frames through artifacts designed for the previous workflow.

### Put reports inside the analysis protocol

Adding report framing or report inference was rejected for v0. Reports are an
orthogonal projection concern. The structured routine outputs should already
contain the summaries, rationales, drivers, findings, unknowns, and limits that
reports need.

### Let agents write files directly

Letting agents write arbitrary run files was rejected as the primary path. It
would scatter path derivation, canonical formatting, validation, overwrite
behavior, and resume semantics across runtimes. Agents should produce JSON
payloads; the CLI should persist them.

### Add recommendation generation now

Recommendation generation was deferred. The v0 replacement should first make the
basic protocol, structured data, synthesis, reporting, and CLI flow work. A later
case can select recommendations from the completed analysis graph without
duplicating the same issue at multiple layers.

### Add QC routines now

QC was deferred. The data shapes leave room for later QC results, but v0 should
not multiply every substantive routine before the core protocol is implemented.

## Trade-offs & risks

This is a large breaking workflow change. Existing evaluation runs should remain
inspectable as historical artifacts if current history readers support them, but
v2 should not pretend old records are v2 data.

The JSON payload set is broad. The implementation should start with shared ID,
ref, status, and confidence helpers so each kind-specific validator does not
invent its own conventions.

Deep Area and Factor trees can produce many files. The data layout is verbose by
design because it improves resume, focused report generation, and navigation.
Centralized path derivation is the main guard against drift.

All-or-nothing report rendering is important but may be awkward if implemented
late. The report builder should be structured around a complete output plan
before it writes files.

Parallel execution can expose hidden ordering assumptions. The deterministic
sequential model order should be treated as the canonical order even when
payloads arrive in a different order.

No migrations means users may need to start fresh v2 evaluations. That is
consistent with the sketch's `schemaVersion` marker decision, but release notes
and skill guidance need to be explicit.

## Open questions

Payload-local IDs are expected to be enough for v0 refs, but implementation may
surface a need for additional typed refs. If that happens, update the record
contract before coding around ad hoc paths.

The exact split between durable spec files and runtime skill files should be
settled while implementing the spec bundle. The runtime should mirror the spec
lightly, not create a second source of truth.
