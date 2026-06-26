---
type: Functional Specification
title: Evaluation v2 report subject links - functional spec
description: Requirements for moving Evaluation v2 generated report links from Details columns into subject table cells.
tags: [cli, evaluation, reports, navigation]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report subject links - functional spec

Companion to the
[Evaluation v2 report subject links](../0105-evaluation-v2-report-subject-links.md)
change case. This spec states the delta for Evaluation v2 Markdown report
rendering. It defers the cumulative report contract to
[`specs/evaluation-v2/reports/report-tree.md`](../../specs/evaluation-v2/reports/report-tree.md)
and the command contract to
[`specs/cli/evaluation-report.md`](../../specs/cli/evaluation-report.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Evaluation v2 Markdown reports are the human review surface over the generated
report tree. In row tables where each Factor, Area, or Requirement has exactly
one generated report page, a separate `Details` column repeats the same link
label for every row. That weakens scanability and hides navigation behind a
generic word instead of the named subject a reader is trying to inspect.

## Scope

Covered:

- Area report `Factors`, `Sub-Areas`, and `Requirements` tables.
- Factor report `Requirements` and `Child Factors` tables.
- Empty-state rows for those tables.
- Durable Evaluation v2 report-tree spec and reference-sketch sync.
- Tests proving the Markdown shape and unchanged machine-readable outputs.

Deferred / non-goals:

- No changes to report tree paths or output refs.
- No changes to breadcrumbs, parent links, compact Factor lists, attached Factor
  links on Requirement reports, or `Data` rows.
- No changes to routine JSON, `data/evaluation-output-result.json`, or
  `qualitymd evaluation report build --json` receipts.
- No changes to Rating Level title resolution or CLI-owned enum display titles.
- No change to in-page `Finding Details` sections.

## Requirements

- Area report `Factors` tables **MUST** render each Factor title in the
  `Factor` column as a link to that Factor's generated report, and **MUST NOT**
  include a separate `Details` column for that row target.

  > Rationale: the Factor name is the meaningful navigation affordance; a
  > repeated `details` label adds width without adding information. - 0105

- Area report `Sub-Areas` tables **MUST** render each child Area title in the
  `Area` column as a link to that child Area's generated report, and **MUST NOT**
  include a separate `Details` column for that row target.

- Area report `Requirements` tables **MUST** render each Requirement title in
  the `Requirement` column as a link to that Requirement's generated report, and
  **MUST NOT** include a separate `Details` column for that row target.

- Factor report `Requirements` tables **MUST** render each Requirement title in
  the `Requirement` column as a link to that Requirement's generated report, and
  **MUST NOT** include a separate `Details` column for that row target.

- Factor report `Child Factors` tables **MUST** render each child Factor title
  in the `Factor` column as a link to that child Factor's generated report, and
  **MUST NOT** include a separate `Details` column for that row target.

- Updated report tables **MUST** keep their existing non-link columns, including
  path, rating, status, child-summary, and attached-Factor columns where those
  columns already exist.

  > Rationale: moving the link should improve navigation without removing the
  > disambiguating and summary data readers already use. - 0105

- Updated report tables **MUST** render explicit empty-state rows with the
  correct post-change column count.

- Generated report subject links **MUST** use the same relative report paths as
  the current `Details` links and the same title fallback behavior as the
  current subject cells.

  > Rationale: the change relocates link presentation only; it must not alter
  > report topology or model display vocabulary. - 0105

- Evaluation v2 reports **MUST** keep explicit `Data` links for machine-readable
  payloads and **MUST NOT** move those links onto subject labels.

  > Rationale: `analysis`, `assessment`, and `rating` links point to structured
  > data, not the generated human report for the row subject. - 0105

- Evaluation v2 machine-readable outputs **MUST NOT** change routine JSON,
  `data/evaluation-output-result.json`, or report-build receipt values because
  of this presentation change.

- Durable report references outside the implementation **MUST** be reconciled so
  live report-tree documentation and reference sketches do not keep prescribing
  `Details` columns for generated report row navigation.

- Report tests **MUST** cover linked subject cells, the absence of redundant
  `Details` headers in affected tables, and unchanged machine-readable report
  output values.

## Durable spec changes

### To add

None

### To modify

- [`specs/evaluation-v2/reports/report-tree.md`](../../specs/evaluation-v2/reports/report-tree.md)
  - require row subject cells to carry generated human report links when the row
    has exactly one generated human report target, and clarify that machine data
    links remain explicitly labeled.

### To rename

None

### To delete

None
