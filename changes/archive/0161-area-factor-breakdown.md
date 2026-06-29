---
type: Change Case
title: Area / Factor Breakdown Reports
description: Replace flat subject-report manifest tables with compact Area / Factor breakdowns in generated Evaluation reports.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Area / Factor Breakdown Reports

Generated Evaluation reports currently use a run-level `Subject Reports` table
as a mixed report manifest and separate Area report `Factors` and `Child Areas`
tables as local structure summaries. This case replaces those surfaces with a
single reader-facing Area / Factor breakdown table.

- [Functional spec](0161-area-factor-breakdown/spec.md) - what the reports must
  change.
- [Design doc](0161-area-factor-breakdown/design.md) - how the renderer will do
  it.

## Motivation

`Subject Reports` is useful as an artifact list, but it is not the overview
readers need from `report.md`: it mixes Areas, Factors, Requirements, findings
indexes, recommendation indexes, and recommendation detail reports in one flat
table. Area reports then use separate tables for Factors and Child Areas, so a
reader has to stitch together the model shape and current Evaluation state from
multiple surfaces.

The report overview should foreground the evaluated Area / Factor shape and its
rating, finding, and recommendation signals. Machine-readable report manifests
can remain in `data/evaluation-output-result.json`.

## Scope

Covered:

- generated run-level `report.md`;
- generated root and non-root Area reports;
- durable report specs, `/quality` reporting spec, report design guidance,
  focused tests, and regenerated report gallery output.

Deferred:

- changing Factor detail report `Sub-Factors` sections;
- adding Requirement rows to the breakdown table;
- adding confidence, status, finding-type mix, recommendation-impact mix, or
  Requirement counts to the breakdown table;
- changing `EvaluationOutputResult.reportOutputs` or generated report paths.

## Affected artifacts

Derived by sweeping for `Subject Reports`, `Child Areas`, `Factors`,
`reportOutputs`, report-gallery output, and generated report design guidance.

- Code:
  - `internal/evaluation/report_tree.go` - render the shared Area / Factor
    breakdown and remove the run-level subject-report body table.
  - `internal/evaluation/evaluation_test.go` - update report assertions for the
    new table and removed old sections.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - replace the `Subject Reports`,
    Area `Factors`, and Area `Child Areas` body contracts with the Area /
    Factor breakdown contract.
  - `specs/skills/quality-skill/reporting.md` - align the skill reporting
    contract with the new generated report shape.
- Durable docs:
  - `docs/guides/reporting-design.md` - update generated report examples and
    design guidance for the breakdown table.
- Generated examples:
  - `examples/report-gallery/` - regenerated generated reports.
- Change Case lifecycle:
  - `changes/index.md`, `changes/log.md`, and this case; archived on
    completion.

No planned impact: Evaluation data schema version, `EvaluationOutputResult`
shape, report paths, Recommendation report contracts, Requirement report
contracts, setup authoring, install docs, bundled skill runtime files, or
historical archived Change Cases.

## Status

`Done`. Implemented and archived. Generated run and Area reports now render
Area / Factor Breakdown tables with Area/Factor links, overall and local
ratings, ranked finding counts, and ranked recommendation counts. The old
run-level `Subject Reports` section and split Area `Factors` / `Child Areas`
sections are removed from report bodies. `go test ./...`, `mise run
fmt-md-check`, and report gallery regeneration pass.
