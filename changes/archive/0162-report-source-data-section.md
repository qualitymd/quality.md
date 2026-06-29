---
type: Change Case
title: Report Source Data Section
description: Move generated report source-data pointers from YAML frontmatter into a visible bottom report section.
status: Done
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report Source Data Section

Generated Evaluation reports currently list report-local source payloads in
YAML frontmatter. That preserves traceability for agents and secondary tooling,
but it makes every report open with a noisy machine list before the human report
content.

- [Functional spec](0162-report-source-data-section/spec.md) - what the reports
  must change.
- [Design doc](0162-report-source-data-section/design.md) - how the renderer
  will do it.

## Motivation

Report source-data pointers should remain easy for people and agents to find,
but they do not need to occupy the first screen of every generated report.
Keeping frontmatter to identity fields and moving report-local source links into
a stable bottom `Source Data` section preserves traceability while making the
top of each report easier to scan.

## Scope

Covered:

- generated Evaluation Markdown reports;
- generated report frontmatter identity fields;
- visible report-local source-data sections;
- durable report specs, CLI and skill reporting specs, report design guidance,
  focused tests, and regenerated report gallery output.

Deferred:

- making Evaluation run folders full OKF bundles;
- adding generated `index.md`, `schema.md`, or `log.md` files to run folders;
- changing the structured Evaluation payload schema;
- using generated report Markdown or frontmatter as report-generation input.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - remove `data` frontmatter emission
    and append visible report-local source-data sections.
  - `internal/evaluation/evaluation_test.go` - update source-data report
    assertions.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - generated report frontmatter
    and source-data section contract.
  - `specs/skills/quality-skill/reporting.md` - `/quality` reporting
    expectations.
  - `specs/cli/evaluation-report.md` - CLI report build contract.
- Durable docs:
  - `docs/guides/reporting-design.md` - report identity frontmatter and bottom
    source-data section guidance.
- Generated examples:
  - `examples/report-gallery/` - regenerated reports.
- Change Case lifecycle:
  - `changes/index.md`, `changes/log.md`, and this case; archive on completion.

## Status

`Done`. Implemented and archived. Generated Evaluation Markdown reports now use
identity-only frontmatter and end with report-local `Source Data` sections;
durable specs, report design guidance, focused tests, and report-gallery output
are aligned. `go test ./...`, `mise run report-gallery-check`, and
`mise run check` pass.
