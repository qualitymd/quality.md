---
type: Change Case
title: Run Report Opening
description: Reshape generated report.md around Summary, Key Details, Contents, and non-judgmental frontmatter metadata.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Run Report Opening

Generated run-level `report.md` currently opens with repeated metadata, an
unlabeled summary table, a bare `Summary:` label, and a `Jump to:` line. This
case makes the run report read like a primary report artifact: concise summary
first, key decision details second, first-class contents navigation third, with
route-oriented metadata in frontmatter and lower report details.

- [Functional spec](0169-run-report-opening/spec.md) - what the run report
  opening and metadata contract must change.
- [Design doc](0169-run-report-opening/design.md) - how the renderer separates
  frontmatter, opening content, contents, and report details.

## Motivation

`report.md` is the evaluation run's decision-ready entrypoint. Its first screen
should help a reader understand the quality result and next action without
decoding duplicated run metadata or reconstructing the document structure. Run
identity and stable references remain useful, but they are routing metadata and
traceability details, not the main human answer.

## Scope

Covered:

- generated run-level `report.md` frontmatter metadata;
- generated run-level `report.md` opening section order;
- `## Summary`, `## Key Details`, and `## Contents` sections;
- lower `## Report Details` traceability;
- removal of visible run-report `Limits & Incomplete Inputs` for now;
- durable report specs and report design guidance;
- generated report gallery output;
- focused report renderer tests.

Deferred:

- Area, Factor, Requirement, finding-index, recommendation-index, and
  recommendation-detail report opening layouts;
- Evaluation data schema changes;
- rating, roll-up, Advice, coverage, and source-data semantics;
- visible limits handling in detail reports.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - render the new run-report opening,
    frontmatter metadata, contents, report details, and visible section set.
  - `internal/evaluation/evaluation_test.go` - update and expand run-report
    rendering assertions.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - generated run report
    frontmatter, opening, contents, details, and visible limits contract.
  - `specs/cli/evaluation-report.md` - align report build contract with the new
    run report shape.
- Durable docs:
  - `docs/guides/reporting-design.md` - report design guidance for opening
    stack, contents, and frontmatter metadata.
- Generated examples:
  - `examples/report-gallery/` - regenerated run report output.
- Release notes:
  - `CHANGELOG.md` - Unreleased CLI note.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md` and `docs/log.md` for durable spec/doc updates.

## Status

`Done`. Implemented and archived. Generated run-level `report.md` now carries
non-judgmental run metadata in frontmatter, opens with `Summary`, `Key Details`,
and `Contents`, moves traceability to `Report Details`, and omits the visible
limits section for now.
