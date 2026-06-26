---
type: Functional Specification
title: Evaluation v2 clean break - functional spec
description: Requirements for removing legacy evaluation compatibility and closing Evaluation v2 data, status, report, CLI, skill, and spec gaps.
tags: [evaluation, cli, reports, skill, cleanup]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 clean break - functional spec

Companion to
[Evaluation v2 clean break](../0097-evaluation-v2-clean-break.md). This spec
states the punchlist for making Evaluation v2 the only active runtime evaluation
contract. The format itself is governed by
[`SPECIFICATION.md`](../../SPECIFICATION.md).

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The Evaluation v2 replacement landed with transitional leftovers. A real
generated run exposed the cost: the report rendered Area and Factor ratings as
`rated`, collapsed finding detail, omitted rating drivers, dropped some
Requirement titles, and produced status guidance framed around legacy counts and
next actions. The v2 graph also accepted shape mismatches that the renderer could
not faithfully display.

The intended product shape is now clear enough to remove the transition layer
instead of maintaining it. This case closes the gap by choosing a single v2 data
contract, deleting previous evaluation runtime surfaces, and updating the
durable specs and bundled skill so agents do not keep producing old artifacts.

## Scope

Covered:

- New and current evaluation runs are v2-only.
- The CLI, Go package, status surface, report builder, durable specs, bundled
  skill, examples, and tests stop advertising or depending on previous
  assessment/analysis/recommendation records.
- `qualitymd evaluation data set` reads a single JSON payload from stdin and no
  longer exposes `--file`.
- Report rendering matches the v2 sketch and v2 spec for ratings, findings,
  evidence, titles, drivers, ordering, and completeness validation.
- `evaluation-v2-sketch.md` is reconciled with durable v2 specs so it no longer
  acts as an untracked competing source of truth.
- The evaluation feedback log findings from the first real run are reflected in
  either requirements or explicit deferrals.

Deferred / non-goals:

- No migration, mixed-version load support, or compatibility transform for old
  evaluation runs.
- No batch or directory data import. The feedback-log request for fewer CLI
  invocations is retained as a later opportunity after stdin-only input lands.
- No recommendation-generation replacement in this case.
- No v2 report gate unless designed as a fresh v2 command in a later case.
- No `qualitymd evaluation data schema <kind>` command in this case.
- No new `evaluation list --state/--limit` or `evaluation data list`
  Area/Factor/Requirement filters in this case unless implementation explicitly
  chooses to carry them forward; unimplemented future CLI filters from the
  sketch remain deferred.
- No judgment QC routine in this case, though v2 data shapes should leave room
  for later QC results.

## Requirements

### Runtime contract

Evaluation commands **MUST** treat Evaluation v2 as the only active runtime
contract.

New evaluation runs **MUST NOT** create `assessments/`, `analysis/`, or
`recommendations/` directories.

New evaluation runs **MUST NOT** seed legacy planning coverage such as
`assessmentResults` or `analyses`.

Commands **MUST NOT** load, repair, migrate, or interpret previous-workflow runs
as compatible v2 runs. When a command encounters an old run shape, it **MUST**
return an unsupported-run diagnostic that directs users to create a new v2 run
instead of attempting compatibility.

Active runtime code **MUST NOT** expose old record writer APIs such as
`AddRecord`, `WriteRecords`, assessment-result records, analysis records, or
recommendation records.

Active runtime code **MUST NOT** produce legacy report artifacts
`report-summary.md` or `report.json`. The old report JSON contract **MUST NOT**
remain the machine interface for v2 reports.

Any report gate behavior based on old `report.json` **MUST** be removed unless a
separate v2 gate is explicitly specified and implemented.

### Run creation and next actions

`qualitymd evaluation create` **MUST** create only the v2 run files and
directories required by the current protocol.

The created run's next action **MUST** point to v2 data persistence through
stdin, for example:

```sh
qualitymd evaluation data set <run> < payload.json
```

The created run **MUST NOT** tell users to run `qualitymd evaluation assessment
add`, `qualitymd evaluation analysis set`, `qualitymd evaluation recommendation
add`, or `qualitymd evaluation data set --file`.

### Data persistence CLI

`qualitymd evaluation data set <run>` **MUST** read exactly one structured JSON
payload from stdin.

`qualitymd evaluation data set` **MUST NOT** accept a `--file` flag.

`qualitymd evaluation data set` **MUST** retain `--dry-run` and `--json`.

When stdin is empty, invalid JSON, multiple JSON documents, or a payload that
does not match a supported v2 kind, the command **MUST** fail without writing
data and **SHOULD** report the failing path or top-level kind when it can do so
safely.

The command **MUST** keep credential and secret redaction behavior in diagnostics
and examples. The feedback-log observation that secret-like command output was
sanitized is a behavior to preserve.

### V2 payload contract

The implementation and durable specs **MUST** define one authoritative status
vocabulary for every v2 payload kind.

Area and Factor scoped analysis statuses **MUST** use the protocol vocabulary
`analyzed`, `empty`, `not_analyzed`, and `blocked`. They **MUST NOT** use
`rated` as an Area or Factor analysis status.

Requirement rating statuses **MUST** remain distinct from Area and Factor
analysis statuses. A Requirement rating may be `rated`; that value **MUST NOT**
be accepted as a scoped Area or Factor analysis status.

The data writer and report builder **MUST** validate v2 payloads deeply enough to
reject schema-incompatible values before they are rendered. This includes
finding shape, rating driver refs, evidence refs, status values, required ids,
and required model references.

Rating drivers **MUST** use canonical structured references to findings,
unknowns, evaluation limits, or other supported driver sources. Bare string
driver refs **MUST** be rejected unless the v2 spec explicitly defines them.

Finding payloads **MUST** use a schema the report renderer can faithfully render:
stable id, type, severity when applicable, summary or description, evidence
refs, impact or rationale when present, and any supported location detail.

Requirement lookup **MUST** include Requirements declared directly under Factors,
not only root-level or Area-level Requirements.

Evaluation v2 JSON **MUST** use direct routine payloads with `schemaVersion` and
`kind`, not a shared legacy record envelope.

`schemaVersion` **MUST** remain a payload-shape marker only. It **MUST NOT**
imply migrations, compatibility transforms, or mixed-version run support.

Structured payloads **MUST** use structured `AreaId`, `FactorId`,
`RequirementId`, and `RatingLevelId` values as primary persisted identities.
Rendered references such as `area:api` or `factor:api::reliability` **MUST NOT**
replace those structured identity fields inside routine JSON.

Optional JSON fields **SHOULD** be omitted when absent unless the field explicitly
defines `null` as meaningful. Repeated fields **SHOULD** default to `[]`.

Payload-local IDs **MUST** be scoped clearly enough that `SourceRef`,
`RoutineOutputRef`, `ReportRef`, `ArtifactRef`, `FindingRef`, `UnknownRef`,
`EvidenceRef`, `EvidenceTargetRef`, and `EvaluationLimitRef` can resolve without
relying on display text.

`qualitymd evaluation data set` **MUST** overwrite the derived routine-output
path by default and write canonical JSON for the accepted payload.

`qualitymd evaluation data list`, `qualitymd evaluation data get`,
`qualitymd evaluation data kinds`, and `qualitymd evaluation data example`
**MUST** remain mechanical v2 data inspection and discovery commands. `data get`
and `data example` **MUST** emit artifact JSON directly without a second JSON
wrapper mode.

`qualitymd evaluation data kinds` **MUST** include every agent-writable payload
kind accepted by `data set` and **MUST** identify CLI-owned payload kinds such as
`EvaluationOutputResult` as not agent-writable.

Unlinked Requirements **MUST** be treated as invalid model input for Evaluation
v2. Evaluation v2 **MUST** stop on that invalid model state rather than defining
roll-up behavior for unlinked Requirements.

When one Requirement is connected to multiple Factors, the same
`RequirementRatingResult` **MUST** contribute to each connected Factor. V0
**MUST NOT** infer link weight, role, or strength unless a later format
extension defines it.

Cross-cutting and model-wide Factors **MUST** use normal Factor IDs in
Evaluation v2 payloads.

### Sketch protocol coverage

The durable Evaluation v2 specs and runtime guidance **MUST** absorb every
current still-relevant contract from `evaluation-v2-sketch.md`, or explicitly
mark it deferred or superseded.

Routine specs and skill instructions **MUST** cover these agent-run prompt
contracts:

- `frameEvaluation`;
- `frameAreaEvaluation`;
- `frameRequirementEvaluation`;
- `assessRequirement`;
- `rateRequirement`;
- `frameFactorAnalysis`;
- `analyzeFactor`;
- `frameAreaAnalysis`;
- `analyzeArea`;
- `assembleEvaluationOutputResult`;
- `generateEvaluationReports`.

Each agent-run routine contract **MUST** identify role or task, allowed inputs,
required output, constraints, stop rules, and self-check expectations.

Frame payloads **MUST** keep the sketch's shared `subject`, `inputs`, and
`derivedContext` shape unless the durable v2 spec explicitly supersedes it.

`EvaluationFrame` **MUST** account for requested scope, resolved scope, Rating
Scale references, in-scope Area/Factor/Requirement identity, evaluation
policies, rigor, and expected run-level evaluation limits.

`AreaEvaluationFrame` **MUST** account for Area-local source refs, direct local
Requirement IDs, root Factor IDs, direct child Area IDs, Area-local scope, and
expected Area limits.

`RequirementEvaluationFrame` **MUST** account for evidence targets, applied
Rating Level criteria, stop conditions, expected evaluation limits,
Requirement-specific rating overrides, and connected Factors.

`FactorAnalysisFrame` and `AreaAnalysisFrame` **MUST** include
`synthesisGuidanceRef`, `emptySignalPolicy`, stop conditions, expected
evaluation limits, and only direct lower-level input refs.

For v0, Factor synthesis **MUST** use
`protocol:factor-synthesis-default-v0`, and Area synthesis **MUST** use
`protocol:area-synthesis-default-v0`. The default synthesis policy **MUST**
preserve rating-binding drivers, surface incomplete inputs, and use a
worst-bound default unless explicit override rationale is recorded.

Synthesis-policy customization from model, Area, or Factor sources **MUST** stay
deferred unless this case explicitly adds it to the durable specs and runtime
guidance.

Requirement Assessment Results **MUST** include `EvidenceTargetCoverage` for the
framed evidence targets. Requirement Rating Results **MUST** include one
`CriterionResult` for every applied criterion and derive `MissingEvidence` from
evidence-target coverage.

Factor and Area Analysis Results **MUST** preserve the sketch's
`localAnalysis` and `localAndDescendantAnalysis` split. Local analysis **MUST**
use only direct inputs for that node; local-and-descendant analysis **MUST**
synthesize local analysis with direct child analyses. Empty input scope **MUST**
be represented as `empty`, not `not_analyzed`.

Factor and Area rating drivers **MUST** preserve lower-level binding drivers
through roll-up and **MUST** distinguish direct Requirement, child Factor, local
analysis, root Factor, child Area, and evaluation-limit refs as applicable.

`IncompleteInput` **MUST** be represented explicitly when referenced ratings,
Factors, Areas, or child analyses are missing, blocked, not rated, not analyzed,
or too low-confidence to rely on fully.

### Status and list behavior

`qualitymd evaluation status` **MUST** inspect the v2 graph and report whether
the run is complete enough for v2 report generation.

`qualitymd evaluation status --json` **MUST NOT** expose legacy record counts
such as `assessmentResults`, `analyses`, or `recommendations`.

Status gaps **MUST** cover missing, malformed, schema-incompatible, stale, and
structurally incomplete v2 data, not only the presence of root frame, Area frame,
and root Area analysis files.

`qualitymd evaluation status` and failed `qualitymd evaluation report build`
**MUST** use the same typed v2 gap model for the same run state.

Status next actions **MUST** distinguish at least:

- more v2 data is required;
- existing v2 data is invalid and must be replaced;
- the report can be built;
- the latest report tree is current.

`qualitymd evaluation list` and project-level `qualitymd status` **MUST** remove
legacy record counts and active legacy recommendation projections. Any
replacement summary **MUST** be based on v2 run status and v2 report outputs.

### Report build and output assembly

`qualitymd evaluation report build` **MUST** fail when required v2 data is
missing, malformed, schema-incompatible, stale, or structurally incomplete.

Report build **MUST** assemble an `EvaluationOutputResult` from validated v2
routine outputs before rendering reports.

Report build **MUST** render a complete output plan before writing generated
report files, and **SHOULD** avoid partial report writes when required
structured data is missing or invalid.

`EvaluationOutputResult` **MUST** include the root Area analysis ref, one
`AreaOutput` index per evaluated Area, and the generated Markdown `ReportRef`s.
Each `AreaOutput` **MUST** reference its Area frame, Area analysis result, local
Factor analysis results, local Requirement assessment/rating results, and
generated Area/Factor/Requirement reports.

The renderer **MUST** consume the assembled `EvaluationOutputResult` and
referenced validated outputs. It **MUST NOT** bypass that assembly by rendering
directly from arbitrary collected artifacts.

Report generation **MUST** write the v2 report tree and
`data/evaluation-output-result.json`. It **MUST NOT** write legacy
`report-summary.md` or legacy `report.json`.

### Report rendering gaps

Area, Factor, and root rating labels **MUST** render selected Rating Level labels
when a rating is present, not status values such as `rated`.

Area tables **MUST** distinguish local rating from aggregate rating instead of
using boolean `+ Sub-Areas` or `+ Sub-Factors` columns where a rating belongs.

Factor child tables **MUST** include child/descendant aggregate rating context
where the v2 sketch or durable report spec requires it.

Requirement tables **MUST** render both assessment status and rating status in a
compact, unambiguous form.

Requirement titles **MUST** be resolved for Factor-declared Requirements.

Rating drivers **MUST** render in report detail, with their target findings,
unknowns, or evaluation limits traceable by id or label.

Finding summaries and finding detail **MUST** render from the canonical v2
finding schema, including confidence, impact/rationale when available, evidence
refs, and severity for gaps.

Evidence refs **MUST** be preserved and expanded enough that a reader can follow
why a finding or rating driver exists. Missing or unknown evidence **MUST** be
reported as an evaluation limit or unknown, not silently dropped.

Report output **MUST** use one deterministic ordering rule for model nodes,
Requirements, findings, drivers, and evidence. The durable spec **MUST** choose
and document the rule, reconciling the sketch's model-order language with the
current structural-id ordering language.

Report tree generation **MUST** preserve the sketch's navigable tree shape unless
explicitly superseded: root `report.md`, non-root Area reports under
`areas/**/report.md`, Area-local Requirement reports, and Area-local Factor
reports with nested child Factor reports.

Every generated report **MUST** include linked breadcrumbs. Every non-root report
**MUST** include a parent link. Area, Factor, and Requirement reports **MUST**
link to their related local reports as described by the v2 report tree contract.

Report generation **MUST** write reports only to paths recorded in
`EvaluationOutputResult`.

Area reports **MUST** include Area field data, Summary, Rating Drivers, Factors,
Sub-Areas, Requirements, and Limits & Incomplete Inputs sections.

Factor reports **MUST** include Factor field data, Summary, Rating Drivers,
Direct Requirements, Sub-Factors, and Limits & Incomplete Inputs sections.

Requirement reports **MUST** include Requirement field data, Summary, Findings
Summary, Finding Details, and Unknowns & Missing Evidence sections.

Reports **MUST** render explicit empty-state rows for empty Factors,
Sub-Areas, Requirements, drivers, limits, direct Requirements, Sub-Factors,
unknowns, and missing evidence.

Reports **MUST** preserve secret-handling boundaries: they may name a locator and
credential or unsafe-content type, but **MUST NOT** reproduce secret values or
unsafe raw content.

The implementation **MUST** reconcile `evaluation-v2-sketch.md` with the durable
v2 specs and generated report behavior. If the sketch remains in the repo after
this case lands, it **MUST** be clearly marked historical or superseded so it
does not compete with the durable specs.

The implementation **MUST** maintain a sketch reconciliation ledger in this
Change Case, durable specs, or implementation checklist until every sketch
heading is implemented, deferred, or superseded.

### Skill, examples, and feedback log

Bundled `/quality evaluate` guidance **MUST** use stdin for data persistence and
**MUST NOT** instruct agents to call `qualitymd evaluation data set --file`.

Bundled skill guidance and durable skill specs **MUST NOT** teach previous
runtime artifacts as active examples: `assessments/`, `analysis/`,
`recommendations/`, `report-summary.md`, legacy `report.json`, `assessment add`,
`analysis set`, or recommendation creation.

The skill **MUST** treat v2 report status and report build diagnostics as the
source of truth for whether a report can be produced.

The skill **SHOULD** preserve the feedback-log lesson that missing package
manager tooling, such as unavailable `yarn`, belongs in workflow feedback or
evidence context rather than being hidden.

The skill **SHOULD** stop creating placeholder `design.md` and `plan.md`
content if those files are not useful evaluation artifacts, or should specify
their v2 purpose clearly enough that a report does not depend on placeholders.

The feedback-log request for bulk data persistence **SHOULD** be recorded as a
future opportunity, but this case **MUST** land the simpler stdin-only
single-payload contract first.

The skill **MUST** preserve the sketch's source-as-data, evidence locator,
secret-handling, confidence, stop-condition, resume, and failure guidance in the
agent-facing Evaluation v2 workflow.

The skill **MUST** treat failed work units as producing either no persisted
output or a valid structured output with the routine-appropriate stopped status
such as `blocked`, `not_assessed`, `not_rated`, or `not_analyzed`.

### Verification

Tests **MUST** cover a v2-only happy path from run creation through data set,
status, report build, and report rendering.

Tests **MUST** include an acquire-shaped regression fixture that would have
caught the observed report gaps: `rated` scoped analysis, string rating drivers,
missing Factor-declared Requirement titles, finding schema mismatch, missing
confidence, and evidence refs.

Tests **MUST** cover report navigation links, explicit empty-state rows,
`EvaluationOutputResult` report path ownership, and report sections for drivers,
limits, incomplete inputs, findings, unknowns, and missing evidence.

Tests **MUST** cover v2 settled decisions from the sketch: unlinked Requirements
are invalid, multi-Factor Requirements contribute the same Requirement rating to
each connected Factor, empty scopes render as `empty`, and lower-level
rating-binding drivers survive Factor and Area roll-up.

Tests **MUST** cover unsupported old run layouts and ensure they are rejected
instead of interpreted.

Search checks or equivalent tests **SHOULD** confirm active runtime docs and
help text no longer contain previous-workflow commands, `--file` data-set usage,
or legacy report artifacts. Historical archived Change Cases and append-only
logs may still mention past behavior.

## Durable spec changes

### To add

None expected. If implementation reveals that the v2 report tree or
`EvaluationOutputResult` deserves a separate 1:1 artifact spec, this case may
add it before moving past `Draft`.

### To modify

- `SPECIFICATION.md` - ensure evaluation semantics describe the v2-only
  runtime, deterministic v2 report tree, and no legacy compatibility promise.
- `specs/evaluation-v2/index.md` and child specs - tighten the v2 protocol,
  payload kinds, routine records, status vocabulary, validation depth,
  `EvaluationOutputResult` assembly, report rendering, report tree, and
  ordering requirements (per the payload, sketch protocol coverage, status,
  report build, and rendering requirements above).
- `specs/evaluation-records.md` and `specs/evaluation-records/**` - stop
  presenting previous assessment/analysis/recommendation records as retained
  active runtime contracts (per the runtime contract requirement above).
- `specs/reports/index.md` and `specs/reports/**` - stop presenting
  `report-summary.md`, legacy `report.md`, and legacy `report.json` as active
  generated artifacts for current evaluation runs (per the report build
  requirement above).
- `specs/cli.md` - align shared CLI artifact JSON behavior with stdin-only data
  persistence and v2-only evaluation output where needed.
- `specs/cli/evaluation-create.md` - remove legacy folders, legacy planning
  files, and `--file` next-action examples (per run creation requirements).
- `specs/cli/evaluation-data.md` - replace `--file <path|->` with stdin-only
  single-payload input while retaining `--dry-run` and `--json` (per data
  persistence CLI requirements).
- `specs/cli/evaluation-status.md` - replace legacy counts and partial v2 gap
  checks with v2 graph completeness and v2 next actions (per status
  requirements).
- `specs/cli/evaluation-list.md` - remove legacy record counts from list output
  (per status and list requirements).
- `specs/cli/evaluation-report.md` - specify v2 output assembly, validation
  failures, and v2 report tree artifacts only (per report build requirements).
- `specs/skills/quality-skill/evaluation.md`,
  `specs/skills/quality-skill/reporting.md`, and
  `specs/skills/quality-skill/quality-skill.md` - align durable skill behavior
  with stdin-only data persistence, v2-only status/report commands, and no
  previous-runtime compatibility.
- `specs/skills/quality-skill/examples/**` - replace or re-scope legacy example
  runs so active examples do not teach old artifacts as current behavior.

### To rename

None expected.

### To delete

- Delete or archive active previous-workflow report specs if "modify" would
  leave them looking like current runtime contracts:
  `specs/reports/report-summary-md.md`, `specs/reports/report-md.md`, and
  `specs/reports/report-json.md`.
- Delete or archive active previous-workflow evaluation record specs if they
  cannot be reframed as historical non-runtime material without implying
  compatibility: `specs/evaluation-records.md` and
  `specs/evaluation-records/**`.
