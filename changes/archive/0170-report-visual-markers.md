---
type: Change Case
title: Report Visual Markers
description: Refine generated report visual markers for recommendation impact and Area / Factor breakdown rows.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Report Visual Markers

Generated Evaluation reports currently use plain text for recommendation impact,
omit impact from the compact run-report Top Recommendations table, and mark
Factor rows in the Area / Factor Breakdown with a loud puzzle-piece emoji while
Area rows have no matching marker. This case makes the report markers calmer and
more semantically specific: recommendation impact gets subtle shape markers, Top
Recommendations regains impact, and Area / Factor rows use quiet box-like type
markers.

- [Functional spec](0170-report-visual-markers/spec.md) - what report display
  markers and Top Recommendations columns must change.
- [Design doc](0170-report-visual-markers/design.md) - how the renderer shares
  impact display and localizes Area / Factor Breakdown markers.

## Motivation

Reports should be fast to scan without making visual markers compete with rating,
severity, or confidence semantics. Red/orange impact colors would imply badness,
and the current puzzle-piece Factor marker distracts in dense breakdown tables.
Subtle non-color shape markers preserve scanability while avoiding accidental
quality-level or concern-level connotations.

## Scope

Covered:

- recommendation impact display labels in generated human reports;
- run-level `report.md` Top Recommendations table columns;
- Area / Factor Breakdown row markers in run and Area reports;
- durable report specs and report design guidance;
- generated report gallery output;
- focused report renderer tests.

Deferred:

- Evaluation structured data schema changes;
- recommendation impact enum changes or ranking semantics;
- rating, confidence, severity, and finding type display labels;
- generated report kind labels outside the Area / Factor Breakdown table;
- non-Markdown or interactive report surfaces.

## Affected artifacts

- Code:
  - `internal/evaluation/display.go` and `internal/evaluation/report_tree.go` -
    render impact display labels and Area / Factor Breakdown markers.
  - `internal/evaluation/evaluation_test.go` - update report rendering and
    helper assertions.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - define impact display labels,
    Top Recommendations impact column, and Area / Factor Breakdown markers.
- Durable docs:
  - `docs/guides/reporting-design.md` - align report design examples and marker
    guidance.
- Generated examples:
  - `examples/report-gallery/` - regenerated report output.
- Release notes:
  - `CHANGELOG.md` - Unreleased CLI note.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md` and `docs/log.md` for durable spec/doc updates.

## Status

`Done`. Implemented and archived. Generated Evaluation reports now use subtle
shape markers for recommendation impact, restore impact to run-report Top
Recommendations, and use box-like Area / Factor Breakdown row markers.
