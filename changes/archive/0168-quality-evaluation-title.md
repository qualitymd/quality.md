---
type: Change Case
title: Quality Evaluation Title
description: Rename generated run report titles from Evaluation Report labels to Quality Evaluation scope titles.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Quality Evaluation Title

Generated `report.md` titles currently read as stacked artifact and subject
labels such as `Evaluation Report: Area: LedgerLite Service`. This case renames
the run report title to `Quality Evaluation - <Area>` and appends factor filters
in parentheses when present.

- [Functional spec](0168-quality-evaluation-title/spec.md) - what the run report
  title contract must change.
- [Design doc](0168-quality-evaluation-title/design.md) - how the renderer
  derives the title from the scoped Area and factor filter.

## Motivation

`report.md` is already a generated report artifact, so putting `Report` in the
visible title is redundant. The current `Area:` segment also over-emphasizes the
implementation detail that the run report renders a scoped Area result. A
factor-scoped evaluation should keep Area context while still showing the user's
requested factor scope in the title.

## Scope

Covered:

- generated `report.md` frontmatter `title`;
- generated `report.md` H1;
- factor-filtered run report titles, including multiple factors;
- durable report specs and report design guidance;
- generated report gallery output;
- focused report renderer tests.

Deferred:

- generated detail report titles such as `Area:`, `Factor:`, `Requirement:`,
  and `Recommendation:`;
- generated report `type` values;
- report file paths, slugs, and navigation labels;
- new generated report frontmatter fields.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - derive the run report heading from
    scoped Area and factor-filter labels.
  - `internal/evaluation/evaluation_test.go` - update title assertions and cover
    multiple factor filters.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - generated run report title
    contract.
  - `specs/cli/evaluation-report.md` - CLI report build contract.
- Durable docs:
  - `docs/guides/reporting-design.md` - report header example and report title
    guidance.
- Generated examples:
  - `examples/report-gallery/` - regenerated run report title output.
- Release notes:
  - `CHANGELOG.md` - Unreleased CLI note.
- OKF logs and indexes:
  - `changes/log.md`, `changes/archive/index.md`, and this archived case.
  - `specs/log.md` and `docs/log.md` for durable spec/doc updates.

## Status

`Done`. Implemented and archived. Generated run report titles now use
`Quality Evaluation - <Area>` with factor filters rendered in parentheses,
durable report specs and design guidance are aligned, and focused tests cover
area-only, single-factor, and multiple-factor title content.
