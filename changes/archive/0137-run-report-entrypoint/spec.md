---
type: Functional Specification
title: Run Report Entrypoint - functional spec
description: What the change must do to make report.md the run report and root-area.md the root Area report.
tags: [evaluation, reports, cli, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Run Report Entrypoint - functional spec

Companion to the
[Run Report Entrypoint](../0137-run-report-entrypoint.md) change case. This
spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

`report.md` is the natural human entrypoint for an Evaluation run, but the
current file is specifically the root Area report. Scoped Evaluation exposes the
mismatch: a run narrowed to a sub-Area or Factor still asks the reader to open a
file whose content is centered on the root Area. The run needs a separate
summary and navigation report, while the root Area report should carry an Area
report filename like every other Area report.

## Scope

Covered: generated Markdown report paths, `EvaluationOutputResult`, report
build receipts, reportability for scoped runs, the root format specification,
durable Evaluation/CLI/skill specs, runtime `/quality` guidance, generated
schema/examples, and tests.

Deferred:

- compatibility copies at old root Area report paths;
- migration of existing evaluation runs;
- report formats outside Markdown and the existing `EvaluationOutputResult`;
- Advice, recommendations, and issue handoff artifacts; and
- UI or viewer work.

## Requirements

### Report paths and content

- `qualitymd evaluation report build <run>` **MUST** generate `report.md` as a
  run-level Evaluation report for every reportable run.

  > Rationale: `report.md` should answer "what did this run evaluate, what was
  > the result, and where are the details?" before sending the reader into a
  > subject-specific Area or Factor report.
  >
  > Durable spec: modify `specs/evaluation/reports/report-tree.md`,
  > `specs/evaluation/records/data-layout.md`, `SPECIFICATION.md`, and
  > `specs/cli/evaluation-report.md` - define the run report as the Markdown
  > entrypoint and headline outcome.

- The root Area report **MUST** be generated as `root-area.md` when the root
  Area has an Area Analysis Result in the run.

  > Rationale: the file should describe its subject. Reserving `report.md` for
  > the run avoids overloading a root Area detail report with entrypoint
  > semantics.
  >
  > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
  > `specs/evaluation/records/data-layout.md` - rename the root Area report
  > path.

- Non-root Area, Factor, and Requirement report paths **MUST NOT** change.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
  > `specs/evaluation/records/data-layout.md` - preserve descendant path
  > contracts.

- The report builder **MUST NOT** write compatibility copies of the root Area
  report to `report.md` or descendant `report.md` paths.

  > Rationale: QUALITY.md is early alpha; a clean break prevents two files with
  > overlapping meaning from drifting.
  >
  > Durable spec: modify `specs/evaluation/reports/report-tree.md` - keep the
  > no-duplicate report path rule.

- The run-level `report.md` **MUST** state the evaluated scope, headline result,
  generated subject reports, data output link, coverage, limits or incomplete
  inputs, and an explicit root Area status when the root Area was not evaluated.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
  > `specs/skills/quality-skill/reporting.md` - define run report content and
  > closeout expectations.

- The run-level `report.md` **MUST NOT** introduce findings, ratings, evidence,
  limits, analysis, recommendations, or candidate actions that are absent from
  structured routine outputs.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
  > `specs/cli/evaluation-report.md` - preserve deterministic projection
  > boundaries.

### Output indexing and receipts

- `EvaluationOutputResult` **MUST** include `runReportRef` for the generated
  `report.md`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
  > `specs/evaluation/records/json-conventions.md`, and
  > `specs/evaluation/reports/report-tree.md` - add a run report reference.

- `EvaluationOutputResult` **MUST** include `headlineResultRef` and
  `headlineReportRef` for the subject that represents the run's headline
  outcome.

  > Rationale: scoped Area and Factor runs should not need to pretend the root
  > Area is the headline result.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
  > `specs/evaluation/protocol.md`, and `specs/evaluation/orchestration.md` -
  > define scoped headline result indexing.

- `EvaluationOutputResult.reportOutputs[]` **MUST** include the run report ref
  and all generated Area, Factor, and Requirement report refs.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` - make the
  > output index complete for generated Markdown reports.

- The report build receipt **MUST** expose the run report path as `reportMd` and
  **MAY** expose `headlineReportMd` and `rootAreaReportMd` when those reports
  exist.

  > Durable spec: modify `specs/cli/evaluation-report.md` - align CLI receipts
  > with the new report roles.

### Scoped reportability

- A run **MUST** be reportable when it contains a valid Evaluation Frame and the
  structured outputs needed for the recorded headline scope, even when the root
  Area is out of scope.

  > Rationale: targeted evaluations should not be forced to synthesize root Area
  > analysis just to generate a human report.
  >
  > Durable spec: modify `specs/evaluation/orchestration.md`,
  > `specs/evaluation/protocol.md`, and `specs/skills/quality-skill/evaluation.md`
  >
  > - make reportability scoped.

- When the Evaluation Frame records one or more `inputs.factorIds`, the first
  listed Factor with a Factor Analysis Result **MUST** be the headline subject.

  > Durable spec: modify `specs/evaluation/protocol.md` - define deterministic
  > headline selection.

- When the Evaluation Frame records no Factor scope and one or more
  `inputs.areaIds`, the first listed Area with an Area Analysis Result **MUST**
  be the headline subject.

  > Durable spec: modify `specs/evaluation/protocol.md` - define deterministic
  > headline selection.

- When no scoped headline input is recorded, the root Area Analysis Result
  **MUST** remain the headline subject.

  > Rationale: this preserves full-evaluation behavior for existing minimal
  > frames while still changing the report file roles.
  >
  > Durable spec: modify `specs/evaluation/protocol.md` - define fallback
  > headline selection.

## Verification

- Report tests **MUST** prove full runs generate `report.md` as a run report and
  `root-area.md` as the root Area report.
- Report tests **MUST** prove generated navigation links point to `root-area.md`
  for the root Area.
- Report tests **MUST** prove descendant Area, Factor, and Requirement report
  paths remain unchanged.
- Scoped reportability tests **MUST** prove sub-Area and Factor runs can build
  reports without root Area analysis when their scoped headline outputs exist.
- Output tests **MUST** prove `EvaluationOutputResult` contains `runReportRef`,
  `headlineResultRef`, `headlineReportRef`, and updated report paths.
- Receipt tests **MUST** prove `reportMd` points to the run report and
  `headlineReportMd` points to the headline subject report when present.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) - the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/evaluation/reports/report-tree.md` - define run report content,
  `root-area.md`, report links, and no compatibility copies.
- `SPECIFICATION.md` - define Evaluation Report semantics around the headline
  outcome instead of only the root Area outcome.
- `specs/evaluation/records/data-layout.md` - update generated Markdown report
  paths.
- `specs/evaluation/records/json-conventions.md` - add run report refs and
  headline report refs.
- `specs/evaluation/records/payload-kinds.md` - update
  `EvaluationOutputResult` shape.
- `specs/evaluation/protocol.md` - define scoped headline result selection.
- `specs/evaluation/orchestration.md` - make report generation dependencies
  scope-aware.
- `specs/cli/evaluation-report.md` - update report build behavior and receipt
  fields.
- `specs/skills/quality-skill/evaluation.md` - align scoped evaluation
  reportability.
- `specs/skills/quality-skill/reporting.md` - align report tree and closeout
  wording.

### To rename

None.

### To delete

None.
