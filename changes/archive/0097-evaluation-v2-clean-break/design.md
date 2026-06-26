---
type: Design Doc
title: Evaluation v2 clean break - design
description: Implementation approach for removing legacy evaluation compatibility and making Evaluation v2 the only runtime evaluation path.
tags: [evaluation, cli, reports, skill, cleanup]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 clean break - design

Answers the [functional spec](spec.md). This design records the implementation
shape for turning the transitional Evaluation v2 implementation into the only
active evaluation runtime.

## Context

The Evaluation v2 implementation currently sits beside the previous evaluation
runtime. `evaluation data set` writes v2 payloads under `data/`, and
`report_v2.go` can build a v2 report tree, but the package still loads old
`assessments/`, `analysis/`, and `recommendations/` records, still exposes old
record writer APIs, still has legacy report rendering, and still uses legacy
counts in status surfaces.

The fix is not a compatibility layer. The user-facing decision is a clean break:
new and current runtime commands should understand v2 runs only, and old run
layouts should be rejected as unsupported instead of repaired, migrated, or
silently interpreted.

## Approach

### Implementation order

Implement the case in a removal-first sequence so compatibility does not survive
as an accidental fallback:

1. Update durable specs, skill mirrors, and examples enough to make v2-only the
   active source of truth.
2. Introduce a single v2 run/graph inspection path used by status and report
   build.
3. Change `evaluation data set` to stdin-only input.
4. Make `evaluation create`, `evaluation status`, `evaluation list`,
   `evaluation report build`, and project `status` depend only on the v2 graph.
5. Delete legacy record writers, loaders, report renderers, and tests after the
   runtime commands no longer reference them.

This ordering keeps each intermediate state understandable: old paths may exist
for a few commits, but no new code should be wired through them.

### V2 run boundary

Define one active run boundary in `internal/evaluation`: a v2 run is a run folder
with `model.md` plus v2 `data/` payloads and generated v2 reports. It is not
identified by old record folders.

Replace the tolerant old-run inspection model with explicit classification:

```text
v2 run          model.md plus v2 data layout
empty v2 run    model.md plus data/ but not enough routine outputs yet
unsupported     old record folders or old-only report artifacts
invalid         missing/unreadable model.md or malformed required v2 data
```

Commands can share this classification while still returning command-specific
next actions. Unsupported runs should produce a clear diagnostic and stop. They
should not populate legacy counts, compatibility gaps, or partial reports.

### Data graph and validation

Make one v2 graph collector the shared internal representation for status and
report build. The graph should be model-aware and should keep both payload data
and validation gaps:

```go
type V2Graph struct {
    Model *model.Spec
    EvaluationFrame map[string]any
    Areas map[string]*V2AreaNode
    Factors map[string]*V2FactorNode
    Requirements map[string]*V2RequirementNode
    Gaps []RunGap
}
```

The exact Go names can differ, but the boundary matters: `status` and
`report build` should not each rediscover required files with separate rules.
The graph collector should:

- walk expected model nodes, not only files present on disk;
- load and canonicalize payloads through the same payload decoder used by
  `data set`;
- validate kind, schema version, subject ids, statuses, rating refs, finding
  shape, evidence refs, and driver refs;
- record missing dependencies as typed gaps; and
- expose a reportable/non-reportable result without rendering.

`data set` should run payload validation in "single payload" mode before write.
Status/report should run graph validation in "complete graph" mode. The two
modes should share enum and structural helpers so the scoped status vocabulary
cannot diverge again.

### Stdin data-set input

Move file reading out of the command contract. `newEvaluationDataSetCmd` should
drop the `--file` flag and always read one JSON object from `cmd.InOrStdin()`.

The stdin reader should keep the existing strict single-document behavior from
`decodeDataPayload`: empty input, invalid JSON, non-object JSON, and multiple
JSON documents are usage errors. If stdin is an interactive terminal and no
payload is available, fail with a concise usage error instead of waiting
indefinitely.

`--dry-run` and `--json` stay on the command. The dry-run path should still
derive the target path and perform full validation without writing.

### Run creation

`CreateRun` should seed only files with a current v2 purpose:

- `model.md` snapshot;
- `data/` directory;
- optional `design.md` or `plan.md` only if their contents explain their v2
  role and are useful to the agent workflow.

If `design.md` and `plan.md` remain, remove the legacy coverage YAML and
placeholder prose. The next action should show stdin usage:

```text
qualitymd evaluation data set <run> < payload.json
```

### Status and list

`RunStatus` should become a v2 status receipt. Replace `RecordCounts` with v2
graph information, such as data artifact counts by kind, reportability, gaps,
and whether generated reports are current.

`evaluation status`, `evaluation list`, and project-level `qualitymd status`
should all read from the same v2 graph/status helper. Project status should stop
using legacy recommendation counts; if no v2 recommendation feature exists,
there is no active evaluation recommendation summary to display.

Next actions should be derived from graph state:

- write or replace v2 data when required data is missing or invalid;
- build the report when the graph is reportable and reports are missing/stale;
- report that the run is current when the report tree matches the validated
  data graph.

### Report build and rendering

`BuildReport` should become v2-only. Remove the old branch that calls
`Inspect`, `Run.Report`, `renderReportSummaryMarkdown`, and legacy `report.json`
generation.

Report build should work in three phases:

1. Build and validate the complete v2 graph.
2. Assemble `EvaluationOutputResult` in memory from that validated graph.
3. Render a complete output plan, then write `data/evaluation-output-result.json`
   and the report tree.

Rendering should receive the assembled output result plus the validated graph
and model snapshot. That gives the renderer one authoritative list of outputs
while still letting it resolve display titles, finding details, evidence, and
driver targets.

The existing `report_v2.go` can be kept as the starting point, but it should be
split along the real boundaries as it grows:

- graph collection and validation;
- output-result assembly;
- Area/Factor/Requirement report rendering;
- shared Markdown table/link/status helpers.

The v2 renderer should use one ordering helper for Areas, Factors,
Requirements, findings, drivers, and evidence. If the model package cannot
preserve frontmatter order reliably, use documented structural-id order and mark
the sketch's model-order language as superseded.

### Legacy code removal

Once commands are on v2 graph helpers, delete the old runtime surface instead of
leaving it uncalled:

- `RecordKind`, `WriteRecordReceipt`, `PlannedCoverage`,
  `AssessmentResultRecord`, `AnalysisRecord`, `RecommendationRecord`, and their
  helpers;
- `AddRecord`, `WriteRecords`, old validation helpers, and old canonical
  examples;
- old `Run` fields for design/plan coverage, compatibility gaps, legacy record
  slices, and `RecordCounts`;
- old `ReportDocument`, `ReportJSON`, legacy Markdown renderers, and
  `GateReport` unless a v2 gate is specified separately;
- tests whose only purpose is old-run compatibility.

After deletion, run `rg` over old command names, old artifact names, and old API
names. Remaining hits should be either this Change Case, archived historical
records, or intentionally historical notes.

### Specs, skill, and examples

Do the durable/spec mirror work in the same branch as the code change so agents
do not learn an obsolete contract from active files.

The root sketch should be handled explicitly: either fold its still-current
details into `specs/evaluation-v2/**` and mark the sketch historical, or remove
it if the project no longer needs it. It should not remain the best description
of report shape.

Run the sketch reconciliation as an explicit checklist, not as an informal read:

- protocol order and work-unit dependencies;
- all prompt contracts from `frameEvaluation` through `generateEvaluationReports`;
- shared frame shape;
- JSON conventions and shared ref/id types;
- v0 synthesis defaults and `emptySignalPolicy`;
- Requirement evidence-target coverage, criterion results, missing evidence,
  unknowns, confidence, and evaluation limits;
- Factor and Area local versus local-and-descendant semantics;
- roll-up driver preservation and incomplete inputs;
- `EvaluationOutputResult` fields and report-path ownership;
- report tree paths, breadcrumbs, parent links, empty-state rows, data links,
  section outlines, and secret-handling;
- orchestration resume/failure behavior; and
- deferred QC, recommendation generation, data schema, and bulk import.

Anything not implemented should be marked deferred or superseded in the durable
specs so the sketch no longer competes with the current contract.

Sketch reconciliation ledger:

| Sketch section                | Disposition in this case                                                                                                                                         |
| ----------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Core idea                     | Preserve as v2-only routine-output protocol.                                                                                                                     |
| Evaluation protocol           | Preserve; absorb into durable v2 protocol/routine specs.                                                                                                         |
| Frame the evaluation          | Preserve in `EvaluationFrame` payload/routine contract.                                                                                                          |
| Frame Area Evaluations        | Preserve in `AreaEvaluationFrame` payload/routine contract.                                                                                                      |
| Frame Requirement Evaluations | Preserve in `RequirementEvaluationFrame` payload/routine contract.                                                                                               |
| Frame Factor Analyses         | Preserve, including direct refs, v0 synthesis guidance, and `emptySignalPolicy`.                                                                                 |
| Assess Requirements           | Preserve, including evidence-target coverage, findings, unknowns, limits, and confidence.                                                                        |
| Rate Requirements             | Preserve, including criterion results, missing evidence, rating drivers, and no new evidence inspection.                                                         |
| Analyze Factors               | Preserve local versus local-and-descendant analysis, direct-input refs, incomplete inputs, and driver roll-up.                                                   |
| Analyze Areas                 | Preserve local versus local-and-descendant analysis, direct child Area refs, incomplete inputs, and root overall rating.                                         |
| Generate Reports              | Preserve deterministic report projection, report tree paths, navigation, empty states, data links, and no new judgment.                                          |
| Area report sections          | Preserve Area field table, Summary, Rating Drivers, Factors, Sub-Areas, Requirements, and Limits & Incomplete Inputs.                                            |
| Factor report sections        | Preserve Factor field table, Summary, Rating Drivers, Direct Requirements, Sub-Factors, and Limits & Incomplete Inputs.                                          |
| Requirement report sections   | Preserve Requirement field table, Summary, Findings Summary, Finding Details, and Unknowns & Missing Evidence.                                                   |
| Pseudocode shape              | Supersede with durable protocol/order and graph validation implementation; keep behavior, not pseudocode.                                                        |
| Orchestration model           | Preserve as skill/runtime guidance: work units, dependencies, allowed parallelism, determinism, persistence, resume, and failure rules.                          |
| Routine list                  | Preserve as the durable routine set for v2.                                                                                                                      |
| Routine prompt contracts      | Preserve in durable skill/routine specs with task, inputs, output, constraints, stop rules, and self-checks.                                                     |
| Spec and runtime organization | Preserve intent; exact file split may differ but durable v2 specs and skill mirrors must own the contracts.                                                      |
| Future CLI surface            | Partly preserve current mechanical commands; stdin-only `data set` supersedes `--file`; `data schema <kind>` and any unimplemented list filters remain deferred. |
| Rating drivers                | Preserve as first-class roll-up drivers.                                                                                                                         |
| Deferred recommendations      | Preserve as a non-goal for this case.                                                                                                                            |
| Persistence shape             | Preserve `data/`, report tree, canonical path derivation, and direct routine JSON files; exact filenames stay governed by durable specs/code.                    |
| JSON conventions              | Preserve direct payloads, `schemaVersion`, `kind`, structured IDs, refs, omitted optionals, and `[]` repeated fields.                                            |
| Evaluation output result      | Preserve as CLI-owned generated structured output and report-path index.                                                                                         |
| Shared JSON types             | Preserve or explicitly supersede in durable record/shared-type specs; no display-text refs as primary persisted identity.                                        |
| Future QC layer               | Defer, while leaving data shapes open for later QC.                                                                                                              |
| Settled decisions             | Preserve unless this case explicitly supersedes one in durable specs.                                                                                            |
| Open questions                | Resolve or carry forward; current open questions are local ID adequacy and skill-vs-CLI contract split.                                                          |

The active skill files should describe the same stdin and v2-only contract the
CLI implements. Legacy worked examples under the skill spec should either be
removed from active examples or reframed so they cannot be mistaken for current
runtime artifacts.

### Verification

Use tests and searches as complementary checks:

- unit tests for stdin-only `data set`, strict single-payload decoding, dry-run,
  and JSON receipt behavior;
- graph/status tests for missing, malformed, schema-incompatible, structurally
  incomplete, unsupported old-run, and reportable states;
- report regression tests using an acquire-shaped fixture that exercises scoped
  status vocabulary, finding shape, rating drivers, Factor-declared
  Requirements, evidence refs, and ordering;
- sketch-conformance tests or fixture assertions for empty-state rows,
  breadcrumbs/parent links, report refs from `EvaluationOutputResult`,
  evidence-target coverage, criterion results, missing evidence, incomplete
  inputs, and driver preservation;
- project status tests proving legacy recommendations/counts are gone; and
- repository searches for old live-surface terms after implementation.

## Alternatives

### Keep old loaders as hidden history support

Rejected. Hidden support is still support: it keeps old fields, tests, and
report artifacts alive, and it makes status/report behavior harder to reason
about. Historical archived files can remain readable as files, but current
commands should not treat them as active evaluation runs.

### Convert old runs to v2 on read

Rejected. A conversion layer would invent frames, Requirement ratings, Factor
frames, driver refs, and output results that the old records never captured. The
result would look like v2 data while lacking v2 judgment provenance.

### Keep `--file -` as a compatibility alias for stdin

Rejected. The requested CLI shape is stdin-only. Keeping the alias would leave
old skill instructions and examples looking acceptable, which is exactly the
drift this case is meant to stop.

### Patch only the renderer

Rejected. The bad generated report is a symptom of a deeper split: loose payload
validation, mismatched status vocabulary, legacy status counts, and a renderer
that bypasses a completed output result. Renderer fixes without graph and
contract cleanup would leave the same failure mode available.

### Generate a v2 `report.json`

Deferred. A machine-readable report may be useful later, but the old
`report.json` contract should not be carried forward by name or shape. For this
case, `EvaluationOutputResult` is the generated structured output.

## Trade-offs & risks

The clean break will make old local evaluation runs unsupported by current
commands. That is deliberate, but the diagnostic needs to be explicit enough
that users understand they should create a new v2 run rather than repair old
files.

Deleting legacy types will cause broad test churn. Keep the implementation
sequence narrow: move commands to v2 graph helpers first, then delete unreferenced
old code, then replace tests.

Deep graph validation can become verbose. Centralize enum, id, and reference
checks early so each payload kind does not develop its own ad hoc validator.

Rendering from `EvaluationOutputResult` plus the graph risks duplicating state.
The output result should be the output index and report contract, while the graph
remains the validated source used for display detail and links.

The skill and specs must land with the code. Otherwise agents may keep producing
`--file`, old report artifacts, or legacy example shapes even after the CLI
rejects them.

## Open questions

- Should `design.md` and `plan.md` remain in new v2 run folders with useful v2
  content, or should run creation remove them entirely?
- Should unsupported old-run diagnostics name the old folders that triggered the
  rejection, or keep the message generic to avoid over-explaining historical
  internals?
- Are payload-local IDs enough for all v0 refs, or does any payload-local
  artifact need a first-class typed ref beyond owner ref plus local ID?
- Which routine contract details belong in runtime skill instructions, and which
  belong in CLI-supported payload validation?
