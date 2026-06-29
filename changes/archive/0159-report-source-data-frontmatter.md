---
type: Change Case
title: Report source-data frontmatter
description: Make generated report frontmatter the source-data manifest and remove duplicate body data links.
status: Done
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report source-data frontmatter

Generated Evaluation reports now carry `data` frontmatter, but their Markdown
bodies still repeat source payload links in `Data` header columns. This case
makes frontmatter the canonical report-local source-data manifest and keeps
report bodies focused on human review.

- [Functional spec](0159-report-source-data-frontmatter/spec.md) - what the
  reports must change.
- [Design doc](0159-report-source-data-frontmatter/design.md) - how the renderer
  will do it.

## Motivation

The `data` frontmatter should have one clear responsibility: list the data files
used as source data for that reporting artifact. Repeating those links in the
visible report body makes headers wider without improving human review, and it
blurs the distinction between human report content and machine-readable
source-data pointers.

## Scope

Covered:

- generated Markdown reports for Evaluation runs;
- `data` frontmatter semantics as report-local source-data manifests;
- removal of visible `Data` summary columns and body source-data links;
- durable report specs, report design guidance, focused tests, and regenerated
  report gallery output.

Deferred:

- making Evaluation run folders full OKF bundles;
- adding new frontmatter fields beyond `type`, `title`, and `data`;
- using generated report frontmatter as report-generation input;
- listing non-`data/` sources such as `model-snapshot.md` in `data`.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` and tests - source-data frontmatter
    lists and removal of body `Data` columns.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - report frontmatter `data`
    semantics and body header contract.
  - `specs/skills/quality-skill/reporting.md` - skill report tree expectations.
  - `specs/cli/evaluation-report.md` - report build source-data boundary, if
    needed for clarity.
- Durable docs:
  - `docs/guides/reporting-design.md` - header patterns and pointer
    frontmatter guidance.
- Generated examples:
  - `examples/report-gallery/` - regenerated reports.
- Change Case lifecycle:
  - `changes/index.md`, `changes/log.md`, and this case; archived on
    completion.

## Status

`Done`. Implemented and archived. Generated report frontmatter `data` now lists
the structured source payloads used to render each report artifact, while report
bodies omit duplicate `Data` columns and report-level source-data links.
