---
type: Design Doc
title: OKF-compatible report headers — design doc
description: Design for shared frontmatter and report header rendering.
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# OKF-compatible report headers — design doc

Design for
[OKF-compatible report headers](../0158-okf-report-headers.md) and its
[functional spec](spec.md).

## Context

The current report renderer writes each report's opening Markdown directly in
`internal/evaluation/report_tree.go`. That keeps each report readable in code,
but it makes frontmatter, report-level navigation, and header conventions easy
to spell differently across report kinds. `internal/markdown` already owns
generated Markdown primitives such as table rows, relative links, and data links.

## Approach

Add a small shared report-header layer inside the Evaluation report renderer:

- a report metadata struct for `type`, frontmatter `title`, visible H1, data
  payload paths, navigation links, context lines, summary table columns, and
  optional jump links;
- report type constants for the subject-shaped taxonomy;
- a frontmatter helper in `internal/markdown` that renders simple YAML
  scalars/lists safely enough for generated titles and paths;
- report-specific header builders that compute data paths and summary rows near
  the existing report-specific logic;
- one shared renderer that writes frontmatter, H1, run/report navigation,
  hierarchy context lines, summary table, and jump links.

The existing body renderers keep ownership of report-specific sections such as
Rating Drivers, Top Findings, Findings Summary, and recommendation prose.
Frontmatter remains output only; no renderer reads generated report frontmatter.

## Spec response

- Pointer-only frontmatter is rendered once from report metadata, limiting drift
  across report kinds.
- The report type taxonomy is centralized as constants, avoiding spelling drift.
- Frontmatter title is computed separately from visible H1 so the YAML subject
  title can omit prefixes while the H1 keeps `Area:`, `Factor:`, and
  `Requirement:` scan labels.
- Data pointers are derived from the same data path helpers used by Data column
  links.
- Report-level navigation is built by one helper that knows the current report
  path and links to `report.md`, `findings.md`, and `recommendations.md`.
- Existing Area, Factor, and Requirement context lines are represented as header
  context lines, preserving the current hierarchy contract.
- Determinism comes from existing artifact ordering and explicit, stable data
  path order in each report header builder.

## Alternatives

- **Hand-author frontmatter in every report renderer.** Rejected because it
  repeats type names, YAML escaping, and data-path lists across many functions.
- **Make generated run folders full OKF bundles now.** Rejected for this change:
  `index.md`, `schema.md`, and `log.md` need a separate bundle contract, and
  users asked to implement the report headers first.
- **Generate frontmatter from `EvaluationOutputResult` after reports are
  rendered.** Rejected because report-local data pointers and visible H1 titles
  are already known in the report-specific renderers; a second post-process
  would be harder to keep local and deterministic.

## Trade-offs & risks

- YAML rendering must escape titles robustly. Keep the helper small but tested.
- Adding navigation to every report changes gallery output substantially; tests
  should assert semantics rather than only broad prefixes.
- Summary tables should stay report-specific; the shared header renderer must
  not grow report-domain judgment.

## Open questions

- Whether a future Change Case should make each Evaluation run folder a full OKF
  bundle with generated `index.md`, `schema.md`, or `log.md`.
