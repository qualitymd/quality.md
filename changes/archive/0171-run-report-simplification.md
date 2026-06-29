---
type: Change Case
title: Run Report Simplification
description: Simplify generated run-level Evaluation reports around the decision-ready result.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Run Report Simplification

Generated run-level `report.md` has accumulated opening navigation, repeated
scope and traceability sections, an extra recommended-next-action sentence, and
long source-data lists that duplicate detail reports. This case simplifies the
run report around the decision-ready result: Summary, Key Details, Model
Evaluation, Top Findings, Top Recommendations, Legend, and Primary Source Data.

- [Functional spec](0171-run-report-simplification/spec.md) - what generated
  run reports and report source-data sections must change.
- [Design doc](0171-run-report-simplification/design.md) - how the renderer and
  report contract deliver the simplified shape.

## Motivation

The run report is the primary human entrypoint for an Evaluation. Its first
screen should get readers to the judgment, model evaluation, and ranked evidence
without forcing them through duplicated navigation, scope, coverage, and
provenance blocks. Granular reports already carry detailed traceability for
their subjects, so the run report should link into them instead of reproducing
their source-data payload lists.

## Scope

Covered:

- generated run-level `report.md` opening and section order;
- run-level removal of the post-H1 `Report:` and `Area:` lines;
- run-level removal of the recommended-next-action sentence;
- run-level removal of `Scope`, `Coverage`, and `Report Details` sections;
- run-level rename and reposition of `Area / Factor Breakdown` to
  `Model Evaluation`;
- generated report `Primary Source Data` section heading and report-local
  primary-source principle;
- durable report specs, `/quality` reporting specs, and report design guidance;
- generated report gallery output;
- focused report renderer tests and release notes.

Deferred:

- Evaluation structured data schema changes;
- finding and recommendation ranking behavior;
- report frontmatter fields;
- detail report navigation trails;
- global renaming of detail report `Area / Factor Breakdown` sections;
- non-Markdown or interactive report surfaces.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - render the simplified run report
    shape and primary source-data sections.
  - `internal/evaluation/evaluation_test.go` - update generated report
    assertions.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - define the simplified run
    report and `Primary Source Data` contract.
  - `specs/cli/evaluation-report.md` - mirror the generated report shape and
    primary source-data contract for `qualitymd evaluation report build`.
  - `specs/skills/quality-skill/reporting.md` - mirror the report-generation
    source-data and breakdown-section contract for `/quality`.
- Durable docs:
  - `docs/guides/reporting-design.md` - align report design principles,
    examples, and checklist.
- Generated examples:
  - `examples/report-gallery/` - regenerated report output.
- Release notes:
  - `CHANGELOG.md` - Unreleased CLI note.
- Release metadata:
  - `skills/quality/SKILL.md` - skill version metadata for the release tag.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md` and `docs/log.md` for durable spec/doc updates.

## Status

`Done`. Implemented and archived. Generated run-level `report.md` now opens
around Summary, Key Details, Model Evaluation, ranked findings, ranked
recommendations, Legend, and Primary Source Data.
