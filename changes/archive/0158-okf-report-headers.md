---
type: Change Case
title: OKF-compatible report headers
description: Add pointer-only OKF-compatible frontmatter and consistent navigation headers to generated Evaluation reports.
status: Done
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# OKF-compatible report headers

Generated Evaluation reports are the primary human review surface, but their
top-of-file context is inconsistent and has no lightweight metadata for agents
or editor previews. This case makes report files easier to scan and navigate
while keeping structured JSON under `data/` as the source of truth.

- [Functional spec](0158-okf-report-headers/spec.md) - what the reports must
  emit.
- [Design doc](0158-okf-report-headers/design.md) - how the renderer will do it.

## Motivation

Readers often land directly on a nested Area, Factor, Requirement, finding, or
recommendation report. They should be able to identify the report subject, the
run context, and the next useful report links within a few seconds. Agents should
also have a small OKF-compatible pointer layer without scraping generated prose
or duplicating Evaluation result data.

## Scope

Covered:

- generated Markdown reports for Evaluation runs;
- pointer-only frontmatter with `type`, subject-only `title`, and canonical
  `data` links;
- visible report header/navigation blocks for run, index, subject, and
  recommendation reports;
- durable report specs, focused tests, and regenerated report gallery outputs.

Deferred:

- full Evaluation run folders as OKF bundles with generated `index.md`,
  `schema.md`, or `log.md`;
- structured Evaluation JSON schema changes other than adding
  `RunManifest.createdAt`, or report reference payload changes;
- using generated report frontmatter as an input to report rendering.

## Affected artifacts

- Code:
  - `internal/markdown/markdown.go` and tests - frontmatter/YAML-safe helpers.
  - `internal/evaluation/report_tree.go` and tests - report metadata, headers,
    report type taxonomy, frontmatter, and navigation rendering.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - generated report frontmatter,
    report type taxonomy, header/navigation contract, and H1-after-frontmatter
    rule.
  - `specs/skills/quality-skill/reporting.md` - skill report tree expectations
    and runtime artifact framing.
  - `specs/cli/evaluation-report.md` - report build output contract, if needed
    for frontmatter/header behavior.
  - `specs/cli/evaluation-create.md` and
    `specs/evaluation/records/payload-kinds.md` - RunManifest `createdAt`
    timestamp used by visible report headers.
- Durable docs:
  - `docs/guides/reporting-design.md` - design principles and taxonomy guidance.
  - `docs/guides/index.md` - guide listing.
- Generated examples:
  - `examples/report-gallery/` - regenerated reports.
- Change Case lifecycle:
  - `changes/index.md`, `changes/log.md`, and this case; archived on completion.

## Status

`Done`. Implemented and archived. Generated Evaluation Markdown reports now
carry OKF-compatible pointer frontmatter, visible run/header navigation,
subject-shaped report types, subject-only frontmatter titles, and
`RunManifest.createdAt`-backed run freshness.
