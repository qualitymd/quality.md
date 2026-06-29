---
type: Change Case
title: Report Frontmatter H1 Titles
description: Align generated report frontmatter titles with the visible H1 document titles.
status: Done
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report Frontmatter H1 Titles

Generated Evaluation report frontmatter currently stores a subject-only
`title`, while the first Markdown H1 stores the actual report document title.
This case makes `title` the plain-text H1 content so generated Markdown reports
have one document-title contract.

- [Functional spec](0167-report-frontmatter-h1-titles/spec.md) - what the
  report title contract must change.
- [Design doc](0167-report-frontmatter-h1-titles/design.md) - how the shared
  report header renderer will enforce it.

## Motivation

Generated report artifacts already carry their artifact kind in `type` and their
reader-facing document title in the H1. Keeping a second subject-only
frontmatter title creates ambiguity for agents, editor previews, and OKF-style
discovery: the metadata title and visible document title disagree. The
frontmatter `title` should identify the generated Markdown document by the same
title the reader sees first.

## Scope

Covered:

- generated Evaluation Markdown report frontmatter;
- generated report H1 titles;
- durable report specs, CLI and skill reporting specs, report design guidance,
  focused tests, and regenerated report gallery output.

Deferred:

- adding new generated report frontmatter fields;
- changing report file paths or title-derived slugs;
- making generated report Markdown or frontmatter an input to report generation;
- changing `type` values for generated report artifacts.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - make the shared report header write
    frontmatter `title` from the H1 text.
  - `internal/evaluation/evaluation_test.go` - update report frontmatter title
    assertions.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - generated report frontmatter
    and H1 title contract.
  - `specs/cli/evaluation-report.md` - CLI report build contract.
  - `specs/skills/quality-skill/reporting.md` - `/quality` reporting
    expectations.
- Durable docs:
  - `docs/guides/reporting-design.md` - report identity frontmatter guidance.
- Generated examples:
  - `examples/report-gallery/` - regenerated reports.
- OKF logs and indexes:
  - `changes/index.md`, `changes/log.md`, and this case; archive on completion.
  - `changes/archive/index.md` after completion.
  - `specs/log.md` and `docs/log.md` for durable spec/doc updates.

## Status

`Done`. Implemented and archived. Generated report frontmatter `title` values
now match visible H1 document titles, durable report specs and design guidance
are aligned, tests are updated, and report-gallery output is regenerated.
