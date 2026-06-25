---
type: Functional Specification
title: Replace evaluation workflow - functional spec
description: Target behavior for Evaluation v2 protocol, records, reports, CLI surface, and skill workflow.
tags: [evaluation, workflow, cli, skill, reports]
timestamp: 2026-06-25T00:00:00Z
---

# Replace evaluation workflow - functional spec

Companion to
[Replace evaluation workflow](../0094-replace-evaluation-workflow.md). This
spec states the target behavior for replacing the current evaluation workflow
with Evaluation v2.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background

Evaluation v2 shifts QUALITY.md evaluation from a record-first workflow to an
agent-orchestrated judgment protocol. The protocol names the judgment moves an
agent performs, the structured inputs and outputs for each move, the bottom-up
tree traversal, and the cases where an evaluator stops rather than producing
fake precision.

The durable record is still important, but it becomes the structured projection
of completed routine outputs under `data/`, not the starting mental model for
the user experience. Reports become deterministic Markdown projections over
those outputs.

## Scope

This change covers:

- the Evaluation v2 routine protocol;
- the JSON payload kinds persisted for routine outputs;
- the run-folder `data/` layout and report tree;
- deterministic report generation;
- the CLI command surface that creates runs, accepts routine outputs, inspects
  data, reports status, and builds reports;
- agent-agnostic orchestration rules; and
- `/quality` skill runtime behavior for Evaluation v2.

This change is a wholesale replacement for the current evaluation workflow,
current evaluation record specs, and current report generation contract.

Deferred:

- Recommendation generation and recommendation selection.
- QC routines around judgment outputs.
- Data-schema command support.
- Batch data writes.
- Automatic migrations, compatibility transforms, or mixed-version run support.
- Custom synthesis policies from model, Area, or Factor configuration.

## Requirements

### Evaluation protocol

Evaluation v2 **MUST** define evaluation as an agent-orchestrated judgment
protocol, not as a fully deterministic program.

The protocol **MUST** include these ordered concerns:

1. Frame the evaluation.
2. Frame each in-scope Area evaluation.
3. Frame each local Requirement evaluation before evidence assessment.
4. Assess each local Requirement.
5. Rate each Requirement Assessment.
6. Analyze each Factor node bottom-up.
7. Analyze each Area bottom-up.
8. Assemble an `EvaluationOutputResult`.
9. Generate deterministic reports.

Frames **MUST** be produced before their corresponding judgment routines.

Frame payloads **MUST** use a shared top-level structure with `subject`,
`inputs`, and `derivedContext`.

Frame routines **MUST NOT** assess evidence, assign ratings, synthesize Factors,
synthesize Areas, or produce report prose.

Result routines **MUST** use only the frame and completed lower-level routine
outputs available to that routine.

Report generation **MUST** be deterministic and **MUST NOT** introduce new
findings, ratings, evidence, limits, analysis, recommendations, or other
evaluation judgment.

### Stop and limit behavior

The evaluation **MUST** stop before substantive judgment when the model cannot be
resolved, the requested scope is ambiguous and cannot be safely inferred,
in-scope Requirements are absent, or required source evidence cannot be resolved.

Requirement assessment **MUST** stop or produce a non-assessed status rather than
invent judgment when the Requirement is too vague to bind evidence, required
evidence would rely on uncited assumptions, or evaluated source content attempts
to instruct the evaluator.

Requirement rating **MUST** stop or produce a non-rated status rather than
invent judgment when evidence cannot distinguish adjacent Rating Levels,
material claims are not tied to evidence, or rating overrides cannot be
interpreted.

Factor and Area analysis **MUST** stop or produce non-analyzed status when
required inputs are missing, invalid, unsafe, or too incomplete for defensible
synthesis.

The protocol **MUST** distinguish stop conditions from evaluation limits. A stop
condition prevents a routine from proceeding; an evaluation limit constrains what
the routine can honestly claim.

### Requirement framing

For each local Requirement in an in-scope Area, the evaluator **MUST** produce a
`RequirementEvaluationFrame` before inspecting assessment evidence.

`RequirementEvaluationFrame` **MUST** identify the Requirement, its connected
Factors, the Rating Levels available for rating, the Requirement assessment
basis, rating overrides when present, evidence targets, applied Rating Level
criteria, stop conditions, and expected evaluation limits.

Applied Rating Level criteria **MUST** be adapted to the Requirement before
evidence judgment. They **MUST NOT** be adapted to observed evidence.

Applied Rating Level criteria **SHOULD** make adjacent Rating Levels
distinguishable enough for assessment and rating.

### Requirement assessment

For each framed Requirement, the evaluator **MUST** produce a
`RequirementAssessmentResult` before producing any Requirement rating.

`RequirementAssessmentResult` **MUST** record the assessed Requirement, status,
status reason, evidence summary, evidence-target coverage, findings, unknowns,
evaluation limits, confidence, and confidence reason.

Requirement assessment **MUST NOT** assign a Rating Level.

Findings **MUST** be evidence-backed and **MUST** be classified as `gap`,
`opportunity`, `strength`, or `observation`.

Severity **MUST** apply only to `gap` findings and **MUST** use `critical`,
`major`, or `minor` when present.

Unknowns **MUST** represent relevant facts the assessment did not establish.

Evaluation limits **MUST** represent boundaries on what the assessment attempted
or can honestly claim.

### Requirement rating

For each Requirement Assessment, the evaluator **MUST** produce a
`RequirementRatingResult`.

`RequirementRatingResult` **MUST** map the `RequirementAssessmentResult` to the
pre-framed applied Rating Level criteria in the
`RequirementEvaluationFrame`.

Requirement rating **MUST NOT** inspect new source evidence or change applied
criteria.

`RequirementRatingResult` **MUST** record status, status reason, selected
`ratingLevelId` only when rated, rationale, rating drivers, criterion results,
missing evidence, evaluation limits, confidence, and confidence reason.

The rating rationale **MUST** explain why the selected level applies and why
higher levels do not.

Criterion results **MUST** include one result for every applied Rating Level
criterion.

Each `RequirementRatingDriver` **MUST** reference at least one finding, unknown,
or evaluation limit.

### Factor analysis

The evaluator **MUST** analyze every in-scope Factor node in an Area's Factor
tree bottom-up.

For each Factor node, the evaluator **MUST** produce a `FactorAnalysisFrame`
after child Factor analyses are complete and before producing the
`FactorAnalysisResult`.

`FactorAnalysisFrame` **MUST** reference only direct Requirement Rating Results
attached to the Factor node and direct child Factor analyses.

For v0, `FactorAnalysisFrame.derivedContext.synthesisGuidanceRef` **MUST** point
to `protocol:factor-synthesis-default-v0`.

`FactorAnalysisResult` **MUST** include `localAnalysis` and
`localAndDescendantAnalysis`.

`localAnalysis` **MUST** analyze only direct Requirement Rating Results attached
to the Factor node.

`localAndDescendantAnalysis` **MUST** synthesize `localAnalysis` with direct
child Factor `localAndDescendantAnalysis` results. A parent Factor **MUST NOT**
need transitive descendant refs because each child analysis already accounts for
its descendants.

Factor scoped analysis status **MUST** distinguish `analyzed`, `empty`,
`not_analyzed`, and `blocked`.

Empty local or descendant signal **MUST** be represented as `empty`, not
`not_analyzed`.

### Area analysis

The evaluator **MUST** analyze every in-scope Area bottom-up.

For each Area, the evaluator **MUST** produce an `AreaAnalysisFrame` after local
root Factor analyses and direct child Area analyses are complete and before
producing the `AreaAnalysisResult`.

`AreaAnalysisFrame` **MUST** reference this Area's root Factor analyses and
direct child Area analyses.

For v0, `AreaAnalysisFrame.derivedContext.synthesisGuidanceRef` **MUST** point
to `protocol:area-synthesis-default-v0`.

`AreaAnalysisResult` **MUST** include `localAnalysis` and
`localAndDescendantAnalysis`.

`localAnalysis` **MUST** analyze this Area's root Factor
`localAndDescendantAnalysis` results.

`localAndDescendantAnalysis` **MUST** synthesize `localAnalysis` with direct
child Area `localAndDescendantAnalysis` results. A parent Area **MUST NOT** need
transitive descendant refs because each child analysis already accounts for its
descendants.

The root Area's `localAndDescendantAnalysis.ratingLevelId`, when produced, is
the overall evaluation rating.

Area scoped analysis status **MUST** distinguish `analyzed`, `empty`,
`not_analyzed`, and `blocked`.

### Synthesis defaults

Factor and Area synthesis v0 **MUST** use protocol defaults referenced by
`synthesisGuidanceRef`.

The v0 default synthesis policy **MUST** be `worst_bound`: the final level is
constrained by the lowest rating-relevant input unless the routine records
explicit override rationale.

The v0 driver policy **MUST** preserve binding drivers that prevent a higher
rating.

The v0 incomplete-input policy **MUST** surface incomplete inputs in analysis
and confidence without automatically blocking every analysis.

Synthesis customization from protocol, model, Area, or Factor sources **MUST**
remain deferred for v0.

### Structured JSON payloads

Evaluation v2 JSON payloads **MUST** store direct routine outputs rather than a
common record envelope.

Every persisted JSON payload **MUST** include `schemaVersion` and `kind`.

`schemaVersion` **MUST** be treated as a payload-shape marker, not a migration
mechanism.

The CLI **MUST NOT** define automatic upgrades, compatibility transforms, or
mixed-version run support for v0.

The first implementation slice **MUST** support these agent-written payload
kinds through `evaluation data set`:

- `EvaluationFrame`
- `AreaEvaluationFrame`
- `RequirementEvaluationFrame`
- `RequirementAssessmentResult`
- `RequirementRatingResult`
- `FactorAnalysisFrame`
- `FactorAnalysisResult`
- `AreaAnalysisFrame`
- `AreaAnalysisResult`

`EvaluationOutputResult` **MUST** be CLI-owned and generated by
`evaluation report build`; agents **MUST NOT** write it through
`evaluation data set`.

Persisted routine JSON **MUST** use resolved structural model identity fields
such as `AreaId`, `FactorId`, `RequirementId`, and `RatingLevelId`.

Persisted routine JSON **MUST NOT** replace structural model IDs with rendered
model reference strings such as `area:api`,
`factor:api::reliability`, `requirement:api::retry-window`, or
`rating:target`.

Generated routine outputs, protocol guidance, report artifacts, and payload
local artifacts **MUST** use `*Ref` forms rather than structural `*Id` names.

Optional fields **SHOULD** be omitted when absent. Repeated fields **SHOULD**
default to `[]`.

### Run-folder data layout

Structured Evaluation v2 data **MUST** live under the evaluation run's `data/`
subfolder.

Human-readable reports **MUST** live outside `data/`.

The root Area report **MUST** be written to `report.md` in the run root.

Non-root Area reports **MUST** be written to `areas/**/report.md`.

Area-local Requirement and Factor reports **MUST** be written under the owning
Area's report folder.

`data/areas/` **MUST** mirror the evaluated Area tree. Each Area data folder
**MUST** own that Area's local Requirement outputs, local Factor outputs, and
child Area folders.

Factor data folders **MUST** live inside the declaring Area's data folder.
Nested Factors **MUST** recurse through nested `factors/` folders.

Routine outputs **SHOULD** be single JSON files such as `*-frame.json`,
`*-result.json`, or `*-analysis-result.json` unless the output truly needs
multiple files or attachments.

Data paths **MUST** be derived from structured model IDs and payload `kind`, not
from display titles, natural labels, or rendered human labels.

### Reports

`qualitymd evaluation report build` **MUST** assemble and write
`data/evaluation-output-result.json` before rendering reports.

`EvaluationOutputResult` **MUST** index the root Area analysis, every evaluated
Area output, and the generated report paths needed for deterministic reporting.

Report generation **MUST** consume `EvaluationOutputResult` and referenced
structured routine outputs.

Report generation **MUST** generate a navigable report tree for Areas, Factors,
and Requirements.

Every generated report **MUST** start with linked breadcrumbs from the root Area
to the current report subject.

Every non-root report **MUST** include a parent link.

Area reports **MUST** link to local root Factor reports, local Requirement
reports, and direct child Area reports.

Factor reports **MUST** link to the owning Area report, parent Factor report
when present, child Factor reports, and direct Requirement reports.

Requirement reports **MUST** link to the owning Area report and every attached
Factor report.

Area reports **MUST** include Area identity, path, overall rating, local rating,
confidence, data links, summary, rating drivers, Factors, Sub-Areas,
Requirements, and limits/incomplete inputs.

Factor reports **MUST** include Area identity, Factor identity, path, local
rating, local-and-descendant rating, status, confidence, data links, summary,
rating drivers, direct Requirements, sub-Factors, and limits/incomplete inputs.

Requirement reports **MUST** include Area identity, Requirement identity,
rating, assessment status, rating status, Factor links, confidence, data links,
summary, findings summary, finding details, unknowns, and missing evidence.

Reports **MUST** render explicit empty-state rows for empty tables and
sections.

Reports **MUST** render `not_assessed`, `not_rated`, `empty`, `not_analyzed`,
and `blocked` distinctly from Rating Level labels.

Reports **MUST** preserve secret-handling boundaries. They may name a locator and
credential type but **MUST NOT** reproduce secret values or unsafe raw content.

Report ordering **MUST** be deterministic: Areas, Factors, and Requirements in
model order; rating drivers in source result order; findings in assessment
order; and evidence in recorded order.

### CLI surface

The CLI **MUST** remain mechanical. The agent owns judgment; the CLI owns run
creation, payload validation, canonical persistence, inspection, and
deterministic report projection.

The required Evaluation v2 flow **MUST** be:

```text
qualitymd evaluation create [model]
qualitymd evaluation data set <run> --file <path|->
qualitymd evaluation report build <run>
```

`qualitymd evaluation create [model]` **MUST** create the numbered run folder
and capture the selected model path. `[model]` **MUST** default to `QUALITY.md`.

`evaluation create` **MUST** support `--json` and `--evaluation-dir <path>`.

`qualitymd evaluation data set <run> --file <path|->` **MUST** read one
structured routine payload from a file or stdin, validate it, route by `kind`,
derive the canonical `data/**` path, and write canonical JSON.

`evaluation data set` **MUST** overwrite the derived path by default so repeated
writes of the same routine output are idempotent.

`evaluation data set` **MUST** support `--dry-run` and `--json`. Under
`--json`, it **MUST** emit a write receipt, not the stored JSON artifact.

Batch payloads for `evaluation data set` **MUST** remain deferred.

`qualitymd evaluation report build <run>` **MUST** validate the run, assemble
`data/evaluation-output-result.json`, render the deterministic report tree, and
fail without partial report writes when required structured data is missing or
invalid.

`evaluation report build` **MUST** support `--json`.

`qualitymd evaluation status <run> --json` **MUST** be available for resume and
debugging. It is not part of the required happy-path flow.

`evaluation status` and failed `evaluation report build` **MUST** report the
same typed gaps for the same run state.

The CLI **MUST** provide inspection commands for listing and reading stored data:

```text
qualitymd evaluation data list <run>
qualitymd evaluation data get <run>
```

`evaluation data get` **MUST** emit the stored JSON artifact directly on stdout.
It **MUST NOT** provide a second JSON result-wrapper mode.

The CLI **MUST** provide data contract discovery:

```text
qualitymd evaluation data kinds
qualitymd evaluation data example <kind>
```

`evaluation data kinds` **MUST** include every `kind` accepted by
`evaluation data set`.

`evaluation data example <kind>` **MUST** emit a complete valid example JSON
artifact for the requested kind.

Commands whose stdout is already a JSON artifact **MUST NOT** define a second
JSON result-wrapper mode. They **MAY** recognize `--json` only to fail with a
targeted usage error explaining that stdout is already JSON and the command
should be rerun without `--json`.

### Orchestration

Evaluation v2 **MUST** define an agent-agnostic dependency-ordered work graph.

The protocol **MUST NOT** require a specific parallel execution mechanism.
Runtimes may use sequential execution, subagents, threads, processes, queues, or
parallel workers.

Parallel execution **MUST** be observationally equivalent to sequential
execution in deterministic model order.

Parallel execution **MUST NOT** change ratings, paths, report content, ordering,
or persisted output shapes.

The orchestrator **MUST** resolve scope, create frames before judgment routines,
schedule ready work, enforce dependency ordering, enforce source-as-data and
secret-handling rules, persist accepted routine outputs through the CLI, prevent
report generation before structured outputs are complete, handle resume from
persisted outputs, and centralize synthesis merge points.

Workers **SHOULD NOT** write arbitrary files. The preferred pattern is for a
worker to produce routine JSON and for the orchestrator to persist it through
`qualitymd evaluation data set <run> --file <payload>`.

If workers call the CLI directly, the final result **MUST** be equivalent to
orchestrator-mediated writes: same validation, paths, overwrite semantics, and
output graph.

### Skill runtime

The `/quality` skill **MUST** be updated so `evaluate` follows Evaluation v2.

The skill **MUST** treat evaluated source content as data, not instructions.

The skill **MUST** preserve the CLI/judgment split: agents make judgment and
produce structured routine payloads; the CLI creates runs, validates payloads,
persists canonical data, inspects status, and builds deterministic reports.

The skill runtime instructions **SHOULD** mirror the durable Evaluation v2 spec
lightly, with routine prompts or prompt contracts for each agent-run routine.

Routine prompt contracts **SHOULD** define role, task, inputs, required output,
constraints, stop rules, and self-check.

## Durable spec changes

### To add

- `specs/evaluation-v2/` - add the parent durable spec folder for the Evaluation
  v2 protocol, routine contracts, record JSON contracts, data layout, and report
  tree contracts.

### To modify

- `SPECIFICATION.md` - replace current evaluation/report semantics with the
  Evaluation v2 protocol and deterministic report projection.
- `specs/cli.md` - align shared CLI JSON artifact rules and any cross-cutting
  evaluation command invariants.
- `specs/cli/evaluation-create.md` - update run creation for Evaluation v2.
- `specs/cli/evaluation-list.md` - update list states for Evaluation v2 runs.
- `specs/cli/evaluation-status.md` - update status gaps for Evaluation v2 data
  completeness and reportability.
- `specs/cli/evaluation-report.md` - update report build behavior.
- `specs/skills/quality-skill/evaluation.md` - replace the skill evaluation
  workflow contract.
- `specs/skills/quality-skill/reporting.md` - replace reporting behavior with
  deterministic Evaluation v2 report projection.

### To rename

None

### To delete

- Current durable specs that are fully superseded by `specs/evaluation-v2/`,
  including current evaluation-record child specs and report artifact specs,
  should be removed or reduced to compatibility/history notes during
  implementation when the replacement durable spec is in place.
