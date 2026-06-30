---
type: Change Case
title: Evaluation Links Blockquote
description: Render generated report Evaluation links as an H1-adjacent blockquote and keep Area/Factor row markers visible in breakdown headers.
status: Done
tags: [evaluation, reports, navigation]
timestamp: 2026-06-30T00:00:00Z
---

# Evaluation Links Blockquote

Generated Evaluation reports currently render `Evaluation links:` as an inline
paragraph after the opening summary, key-details table, or report-specific
orientation. That keeps links available, but the links can visually blend into
nearby report content. The report header needs a light, consistent treatment
that distinguishes cross-artifact navigation without adding another section or a
heavy divider. The same report presentation pass also keeps Area and Factor row
markers visible in breakdown table headers so the glyphs in the first column
have an immediate key.

- [Functional spec](0184-evaluation-links-blockquote/spec.md) - what the report
  link presentation must change.
- [Design doc](0184-evaluation-links-blockquote/design.md) - why the renderer
  uses a blockquote directly below the H1.

## Motivation

Readers should notice the cross-artifact links immediately when they open any
generated report, while the report's own summary, metadata, and tables remain
separate. A blockquote under the H1 gives the link cluster a distinct visual
treatment in common Markdown renderers without creating table markup, horizontal
rule chrome, or a new body section.

## Scope

Covered:

- generated Markdown report `Evaluation links:` presentation;
- placement of the links immediately after each report H1;
- `Model Evaluation` and `Area / Factor Breakdown` first-column labels;
- durable report specs, CLI report spec alignment, report design guidance,
  tests, generated examples, release notes, and logs.

Deferred:

- changing the set or order of linked artifacts;
- changing relative-link resolution;
- adding CSS, HTML, tables, horizontal rules, or a dedicated navigation section;
- changing report frontmatter or structured Evaluation data.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - render `Evaluation links:` as a
    blockquote immediately below H1s and render the breakdown first column as
    `▦ Area / □ Factor`.
  - `internal/evaluation/evaluation_test.go` - assert blockquote text and
    placement, plus the breakdown first-column label.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - define blockquote shape and H1
    placement, plus the breakdown first-column label.
  - `specs/cli/evaluation-report.md` - keep run-report opening expectations
    aligned with the generated report contract.
- Durable docs:
  - `docs/guides/reporting-design.md` - update generated report examples and
    checklist guidance.
  - `CHANGELOG.md` - note the user-visible report presentation change.
- Generated examples:
  - `examples/report-gallery/` - regenerate generated Markdown reports.
- OKF logs and indexes:
  - `changes/log.md`, `changes/archive/index.md`, and this case.
  - `specs/log.md`, `specs/evaluation/log.md`, and `docs/log.md`.

## Status

`Done`. Implemented and archived with renderer behavior, durable specs, docs,
tests, generated examples, logs, and release notes aligned.
