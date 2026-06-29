---
type: Functional Specification
title: Report Full List Links — functional spec
description: Requirements for emphasized run-report links to full findings and recommendations lists with total counts.
tags: [evaluation, reports, recommendations, findings]
timestamp: 2026-06-29T00:00:00Z
---

# Report Full List Links — functional spec

Companion to
[Report Full List Links](../0177-report-full-list-links.md). This spec states
the delta for generated `report.md` links from capped top-list sections to the
complete Findings and Recommendations reports. The durable source of truth is
absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md).

The key words **MUST** and **MUST NOT** are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

`report.md` is the decision-ready entry point, while `findings.md` and
`recommendations.md` hold the complete ranked lists. Because the `report.md` Top
Findings and Top Recommendations tables are capped, the complete-list links are
not just supporting navigation; they are the path to rows not visible in the
overview. The link line should therefore be easy to scan and should state the
total ranked count.

The current labels also use `Full Findings report` and `Full Recommendations
report`, which makes a sentence read like a title. Sentence-case labels preserve
the artifact names without over-capitalizing ordinary prose.

## Scope

Covered: generated `report.md` full-list lines below Top Findings and Top
Recommendations, their Markdown emphasis, casing, link labels, and total counts.

Non-goals:

- changing the 10-row cap for Top Findings or Top Recommendations;
- changing ranked-list table columns, sorting, or filtering;
- changing `findings.md` or `recommendations.md` report contents;
- changing report filenames, report frontmatter, or top navigation.

## Requirements

1. Generated `report.md` **MUST** render the full Findings report line below the
   Top Findings table with a bold sentence-case label, the linked artifact
   filename as link text, and the total ranked finding count.

   > Rationale: The complete-list link is the escape hatch from a capped table,
   > so it should stand out and state how much content the full report contains.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

2. Generated `report.md` **MUST** render the full Recommendations report line
   below the Top Recommendations table with a bold sentence-case label, the
   linked artifact filename as link text, and the total ranked recommendation
   count.

   > Rationale: Recommendations are the action-planning path; the overview must
   > make the complete advice list visible without obscuring the artifact target.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

3. The full-list line link text **MUST** remain the literal target filename
   (`findings.md` or `recommendations.md`), and the total count **MUST NOT** be
   included inside the link text.

   > Rationale: Keeping filenames as link text makes the artifact target clear
   > in plain Markdown and avoids making counts part of the destination label.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

4. The total count on each full-list line **MUST** reflect the complete ranked
   list size for that report, not the number of rows rendered in the capped
   `report.md` table.

   > Rationale: The count exists to tell readers whether the full report has
   > additional rows beyond the overview table.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

## Acceptance criteria

- Generated run reports render:
  `**Full findings report:** [findings.md](findings.md) (<N> total)`.
- Generated run reports render:
  `**Full recommendations report:** [recommendations.md](recommendations.md) (<N> total)`.
- The count for each line is the total ranked list length, including when it is
  greater than the 10-row overview table cap.
- Focused Go tests and report-gallery checks pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - specify the emphasized
  sentence-case full-list lines and total ranked counts (per requirements 1-4).

### To rename

None

### To delete

None
