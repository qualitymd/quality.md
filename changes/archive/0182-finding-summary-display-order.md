---
type: Change Case
title: Finding Summary Display Order
description: Rename the run report Finding Breakdown to Finding Summary, show every Finding type, and use concern-first ordering.
status: Done
tags: [evaluation, reports, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Finding Summary Display Order

Generated run reports currently render `Finding Breakdown` as a sparse table:
only Finding types present in the ranked Findings appear, the heading reads like
an occurrence-only breakdown, and the `Detail` column hides that the extra
information is severity detail.

- [Functional spec](0182-finding-summary-display-order/spec.md) - what the run
  report Finding Summary must change.
- [Design doc](0182-finding-summary-display-order/design.md) - how the enum
  catalog, renderer, durable specs, and gallery absorb the display cleanup.

## Motivation

The run report opening should make the full Finding taxonomy visible without
requiring readers to infer absent categories from the legend. A fixed
concern-first summary helps readers scan current issues before supporting
strengths or neutral notes, while the table name and severity column should
describe exactly what the table contains.

## Scope

Covered:

- run report Finding Summary heading and table columns;
- Finding type display order in generated reports and local keys;
- zero-count rows for absent Finding types in the run report summary;
- severity-count rendering for `gap` and `risk`;
- durable Evaluation report specs, CLI report docs, report design guidance,
  tests, generated examples, and release notes.

Deferred:

- changing structured Finding type values or validation;
- changing Finding ranking order;
- making Finding `severity` optional or type-specific;
- adding severity counts for strengths or notes;
- changing non-run report tables beyond shared Finding type ordering in legends
  and local keys.

## Affected artifacts

- Code:
  - `internal/evaluation/display.go` - reorder the Finding type enum catalog.
  - `internal/evaluation/report_tree.go` - render `Finding Summary`, show every
    Finding type row, and rename the `Detail` column to `Severity`.
  - `internal/evaluation/evaluation_test.go` - update report and legend
    expectations.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - update Finding type order and
    run report Finding Summary contract.
  - `specs/cli/evaluation-report.md` - keep the CLI report command contract in
    sync with the run report shape.
- Durable docs:
  - `docs/guides/reporting-design.md` - update run report examples and local-key
    order.
  - `CHANGELOG.md` - add the user-facing report change.
- Generated examples:
  - `examples/report-gallery/` - regenerate generated Markdown reports.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md`, `specs/evaluation/log.md`, and `docs/log.md`.

## Status

`Done`. Implemented and archived. Generated run reports now render `Finding
Summary`, show every Finding type in concern-first order including zero-count
rows, and label sparse concern severity counts as `Severity`.
