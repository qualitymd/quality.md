---
type: Functional Specification
title: Evaluation Links Blockquote
description: Requirements for rendering generated report Evaluation links as an H1-adjacent blockquote and keeping Area/Factor row markers visible in breakdown headers.
tags: [evaluation, reports, navigation]
timestamp: 2026-06-30T00:00:00Z
---

# Evaluation Links Blockquote

This change alters only visible Markdown presentation in generated reports:
`Evaluation links:` navigation moves into an H1-adjacent blockquote, and
Area/Factor breakdown tables expose their row-marker key in the first column
header. It does not change linked artifact targets, link ordering, report
frontmatter, structured Evaluation data, or glossary behavior.

The key words "MUST" and "MUST NOT" are to be interpreted as described in BCP 14
when, and only when, they appear in all capitals.

## Background / motivation

The `Evaluation links:` line introduced by the glossary/report-links work gives
readers a stable path to the overview, Findings, Recommendations, and glossary.
Rendered as a plain paragraph below summary or key-details content, it can read
like ordinary report prose. A blockquote immediately under the H1 makes the
links visible as header-level navigation while staying lightweight and portable
across Markdown renderers.

## Requirements

1. Generated Markdown reports **MUST** render `Evaluation links:` as a single
   blockquote line.

   > Rationale: The links need visual distinction, but a blockquote is lighter
   > than horizontal rules or a table and remains ordinary Markdown.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

2. The generated `Evaluation links:` blockquote **MUST** render immediately
   after the visible H1 in every generated Markdown report.

   > Rationale: The navigation belongs to the report header, not the report body;
   > placing it before run context, key details, and summaries keeps it
   > discoverable without interrupting body sections.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

3. The generated `Evaluation links:` blockquote **MUST** preserve the existing
   filename link text, link order, and artifact-relative link targets for
   `report.md`, `findings.md`, `recommendations.md`, and `glossary.md`.

   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

4. Generated reports **MUST NOT** add horizontal rules, one-cell tables, or a
   new `Evaluation links` section for this navigation treatment.

   > Rationale: Those alternatives add heavier report chrome or a body section
   > for a utility link cluster.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

5. Generated `Model Evaluation` and `Area / Factor Breakdown` tables **MUST**
   label their first column `▦ Area / □ Factor`.

   > Rationale: The first-column cells use `▦` for Area rows and `□` for Factor
   > rows; including both markers in the header gives the table a compact
   > in-place key without adding a separate legend.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define `Evaluation links:` as a
  blockquote immediately below each report H1 while preserving link text, order,
  and relative targets, and define the breakdown first-column label as
  `▦ Area / □ Factor`.
- `specs/cli/evaluation-report.md` - align run-report opening expectations with
  the generated report blockquote placement.

### To rename

None

### To delete

None
