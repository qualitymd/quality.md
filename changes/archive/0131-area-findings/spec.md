---
type: Functional Specification
title: Area findings in evaluation reports — functional spec
description: What the change must do to add Area Findings to Area analysis results and reports.
tags: [evaluation, reports, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Area findings in evaluation reports — functional spec

Companion to the
[Area findings in evaluation reports](../0131-area-findings.md) change case.
This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

Requirement Findings capture direct assessment observations. Area analysis needs
its own durable finding layer for the synthesized conclusions a reader expects
on Area and Factor reports: the important gaps, risks, strengths, unknowns, and
notes that explain the Area's quality without jumping straight to
recommendations. Adding Area Findings first gives future advice a stable object
to address and keeps reports from inventing analysis-only prose.

## Scope

Covered: `AreaAnalysisResult.findings`, the `/quality` skill behavior that
authors those findings during Area analysis, CLI validation/schema/example
support for the new field, and deterministic Area/Factor report rendering.

Deferred: recommendations, recommendation impact/confidence, global top-findings
synthesis across the evaluated scope, impact/priority/ROI, Area-importance or
Factor-weight scores, and any report-only findings.

## Requirements

### Area Finding data contract

- `AreaAnalysisResult` **MUST** accept a `findings` array containing Area Finding
  objects. When absent, report generation and status inspection **MUST** treat it
  as an empty array.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/cli/evaluation-data.md` — add the `AreaAnalysisResult.findings`
  > payload field and empty-array semantics.

- Each Area Finding **MUST** be scoped to the containing
  `AreaAnalysisResult.areaId`; the finding object **MUST NOT** carry its own
  `areaId` field.

  > Rationale: containment keeps the Area association singular and prevents
  > sibling-Area spanning from entering this first slice.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — define
  > Area Finding scope by containment.

- Each Area Finding **MUST** include `id`, `type`, `severity`, `confidence`, and
  `summary`. The `id` **MUST** be unique within the containing
  `AreaAnalysisResult.findings` array.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — define
  > required Area Finding fields and local ID uniqueness.

- Area Finding `type` **MUST** use the same finding type vocabulary as
  Requirement Findings: `strength`, `gap`, `risk`, `unknown`, or `note`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `SPECIFICATION.md` — share the finding type vocabulary across Requirement
  > Findings and Area Findings.

- Area Finding `severity` **MUST** use the same severity vocabulary as
  Requirement Findings: `critical`, `high`, `medium`, `low`, or `info`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `SPECIFICATION.md` — share the finding severity vocabulary across
  > Requirement Findings and Area Findings.

- Area Finding `confidence` **MUST** use the evaluation confidence vocabulary:
  `high`, `medium`, `low`, or `none`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — add
  > confidence to Area Findings using the existing evaluation confidence set.

- The CLI schema for `AreaAnalysisResult` **MUST** expose Area Finding `type`,
  `severity`, and `confidence` as JSON Schema `enum` values, and `data set`
  **MUST** reject Area Findings whose values fall outside those closed sets.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/cli/evaluation-data.md` — require schema-visible enum sets and
  > matching write-time validation for Area Finding classification fields.

- Each Area Finding **MUST** include a non-empty `inputRefs` array using the
  existing routine-reference contract. The references **MUST** point to persisted
  evaluation data that supports the synthesis, such as Requirement assessment
  results, Requirement rating results, Factor analysis results, child Area
  analysis results, or relevant limits.

  > Rationale: Area Findings are synthesis over recorded evaluation outputs, not
  > fresh evidence collection.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/evaluation/routines/routine-contracts.md` — require traceable inputs
  > for Area Findings.

- An Area Finding **MAY** include `rationale` explaining why Area analysis
  synthesized the finding from its inputs.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — add the
  > optional rationale field.

### Factor relationships

- An Area Finding **MAY** include `factorRelationships`, an array of relationships
  to Factors declared in the same Area as the containing `AreaAnalysisResult`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — add the
  > optional Factor relationship array.

- Each Factor relationship **MUST** include `factorId` and `relationship`, and
  **MAY** include `rationale`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — define the
  > relationship object shape.

- A Factor relationship's `factorId` **MUST** resolve against the run's
  `model-snapshot.md` and **MUST** identify a Factor declared in the same Area as
  the containing `AreaAnalysisResult.areaId`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/cli/evaluation-data.md` — add same-Area model-binding validation for
  > Area Finding Factor relationships.

- Factor relationship `relationship` **MUST** use one of:
  `primary-driver`, `contributing-driver`, `evidence-limit`,
  `offsetting-strength`, or `related`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — define
  > the closed relationship vocabulary.

- The CLI schema for `AreaAnalysisResult` **MUST** expose Factor relationship
  `relationship` as a JSON Schema `enum` value set, and `data set` **MUST**
  reject Factor relationships whose values fall outside that closed set.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/cli/evaluation-data.md` — require schema-visible enum sets and
  > matching write-time validation for Area Finding Factor relationships.

- Area Findings **MUST NOT** carry Area importance, Factor weight, recommendation
  impact, recommendation priority, ROI, candidate actions, or recommendations.
  The Area Finding object and Factor relationship object **MUST** be closed
  objects: unknown or misspelled fields, including those forbidden advice/ranking
  fields, **MUST** be rejected by `data set` and omitted from the emitted schema.

  > Rationale: those judgments belong to later synthesis or advice phases; adding
  > them here would make Area analysis appear to know advice-phase impact before
  > recommendations exist.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
  > `specs/cli/evaluation-data.md`,
  > `specs/evaluation/routines/routine-contracts.md`, and `SPECIFICATION.md` —
  > keep Area Findings as findings, not advice, with the CLI enforcing the closed
  > object boundary.

### Skill workflow

- The `/quality` `analyzeArea` routine **MUST** produce Area Findings when Area
  analysis identifies synthesized gaps, risks, strengths, unknowns, or notes
  that materially explain the Area's quality. It **MUST** record those findings
  on `AreaAnalysisResult.findings`.

  > Durable spec: modify `specs/evaluation/routines/routine-contracts.md`,
  > `specs/skills/quality-skill/evaluation.md`, and
  > `specs/skills/quality-skill/workflows/evaluate.md` — make Area Finding
  > synthesis part of Area analysis.

- `analyzeArea` **MUST NOT** inspect new source evidence to create Area Findings;
  it **MUST** synthesize them only from its frame and persisted routine outputs
  already available to Area analysis.

  > Durable spec: modify `specs/evaluation/routines/routine-contracts.md` and
  > `specs/skills/quality-skill/evaluation.md` — preserve the phase boundary
  > between evidence assessment and analysis.

- The `/quality` evaluation QC phase **MUST** include Area Findings that bind an
  Area roll-up rating or carry low confidence in the same verification boundary
  as other roll-up-binding or low-confidence findings.

  > Durable spec: modify `specs/skills/quality-skill/evaluation.md` — include
  > Area Findings in the existing QC obligation.

### Reports

- Area reports **MUST** include a `Findings` section rendered from the containing
  `AreaAnalysisResult.findings` array.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md`,
  > `specs/skills/quality-skill/reporting.md`, and `SPECIFICATION.md` — add Area
  > report Findings.

- Factor reports **MUST** include a `Findings` section rendered from the owning
  Area's `AreaAnalysisResult.findings` array, filtered to findings whose
  `factorRelationships` include the current Factor.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md`,
  > `specs/skills/quality-skill/reporting.md`, and `SPECIFICATION.md` — add
  > Factor report Findings projected from Area Findings.

- Requirement reports **MUST** continue to render Requirement Assessment
  Findings from `RequirementAssessmentResult.findings`; this change **MUST NOT**
  replace or rename Requirement Findings.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` only if needed
  > to clarify that Requirement report Findings remain unchanged.

- Report generation **MUST NOT** create, synthesize, rewrite, or globally rank
  Area Findings. Reports **MUST** render only persisted `AreaAnalysisResult`
  findings.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
  > `SPECIFICATION.md` — preserve deterministic report projection and defer
  > global top-finding synthesis.

### Report ordering

- Area reports **MUST** sort Area Findings by finding type, severity, confidence,
  then original `AreaAnalysisResult.findings` order.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` — define Area
  > report finding order.

- Factor reports **MUST** sort matching Area Findings by finding type, severity,
  Factor relationship, confidence, then original `AreaAnalysisResult.findings`
  order.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` — define
  > Factor report finding order.

- Report sorting **MUST** order finding types as `risk`, `gap`, `unknown`,
  `note`, then `strength`; severities as `critical`, `high`, `medium`, `low`,
  then `info`; confidences as `high`, `medium`, `low`, then `none`; and Factor
  relationships as `primary-driver`, `contributing-driver`, `evidence-limit`,
  `offsetting-strength`, then `related`.

  > Rationale: the order keeps possible harm and unmet expectations visible
  > before neutral notes and strengths, without claiming an unimplemented
  > importance score.
  >
  > Durable spec: modify `specs/evaluation/reports/report-tree.md` — define the
  > deterministic sort orders.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) — the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `SPECIFICATION.md` — recognize Area Findings as analysis-level findings,
  preserve Requirement Findings, update report semantics for Area/Factor
  Findings, and state that global top-finding synthesis and recommendations are
  deferred.
- `specs/evaluation/routines/routine-contracts.md` — require `analyzeArea` to
  produce traceable Area Findings from existing evaluation inputs only.
- `specs/evaluation/records/payload-kinds.md` — define
  `AreaAnalysisResult.findings`, Area Finding fields, Factor relationships,
  same-Area Factor binding, schema-visible closed enums, closed object shapes,
  and forbidden advice/ranking fields.
- `specs/evaluation/reports/report-tree.md` — add Area and Factor Findings
  sections, filtering, empty states, and deterministic ordering.
- `specs/cli/evaluation-data.md` — reflect the new accepted payload validation
  and schema behavior for Area Findings, enum visibility, closed object
  rejection, and same-Area Factor binding.
- `specs/skills/quality-skill/evaluation.md` — update `/quality` evaluation
  behavior for Area Finding synthesis and QC.
- `specs/skills/quality-skill/reporting.md` — update report expectations and
  keep generated recommendations out of scope.
- `specs/skills/quality-skill/workflows/evaluate.md` — update evaluate workflow
  phase responsibilities for Area Findings.

### To rename

None.

### To delete

None.
