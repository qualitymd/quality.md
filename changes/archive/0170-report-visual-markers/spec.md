---
type: Functional Specification
title: Report Visual Markers - functional spec
description: Requirements for generated report impact labels and Area / Factor Breakdown markers.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Report Visual Markers - functional spec

Companion to [Report Visual Markers](../0170-report-visual-markers.md). This
spec states the delta for generated Evaluation Markdown reports. The durable
source of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Recommendation impact is a priority/value signal, not a quality rating,
confidence level, or severity warning. Human reports need it to scan quickly
without using colors that imply badness or reuse the existing rating,
confidence, and severity palettes. The Area / Factor Breakdown has a related
scanability problem: the puzzle-piece Factor marker draws more attention than
the row's rating and recommendation counts, while Area rows have no symmetric
type marker.

## Scope

This change covers generated Markdown reports:

- recommendation impact display labels;
- run-level `report.md` Top Recommendations columns;
- Area / Factor Breakdown row markers.

Deferred:

- Evaluation structured data schema changes;
- impact enum values and ranking behavior;
- rating, confidence, severity, and finding type display labels;
- report kind labels outside Area / Factor Breakdown rows;
- non-Markdown or interactive report views.

## Requirements

1. Generated Markdown reports **MUST** render known recommendation impact values
   with these display labels: `◆ Very high`, `▲ High`, `● Medium`, and
   `○ Low`.

   > Rationale: Shape markers make impact scannable without implying that high
   > impact is bad or reusing the color semantics already carried by ratings,
   > confidence, and severity. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

2. Generated Markdown reports **MUST** continue to humanize unknown
   recommendation impact values without adding a marker.

   > Rationale: Unknown values should remain legible during contract drift while
   > avoiding a misleading marker for values outside the known impact scale.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

3. The run-level `report.md` Top Recommendations table **MUST** render the
   columns `Rank`, `#`, `Recommendation`, `Area / Factors`, `Impact`, and
   `Reason`, in that order.

   > Rationale: The compact run report recommendation table should show why a
   > recommendation matters without forcing readers to open the full
   > recommendation index. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

4. Recommendation index and recommendation detail reports **MUST** use the same
   recommendation impact display labels as the run report.

   > Rationale: A single impact display contract prevents report pages from
   > disagreeing about the same recommendation field. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

5. Area / Factor Breakdown Area rows **MUST** use `▦` as the Area row marker and
   Factor rows **MUST** use `□` as the Factor row marker.

   > Rationale: Quiet box-like markers distinguish row type without the stronger
   > metaphor and visual weight of the puzzle-piece emoji. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

6. Area / Factor Breakdown markers **MUST** appear inside the linked row label's
   visible text while preserving the existing link target, indentation, row
   order, and bold treatment for the root row.

   > Rationale: Keeping the marker with the row label makes the table readable in
   > plain Markdown while preserving the existing hierarchy and navigation
   > contract. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

7. Generated report kind labels outside Area / Factor Breakdown rows **MUST NOT**
   change as part of this case.

   > Rationale: The distracting marker is local to dense breakdown rows; changing
   > report kind labels would broaden the visual language change beyond the
   > current need. Durable spec: none; existing durable specs already govern
   > report kind labels.

## Acceptance criteria

- Known recommendation impact values render as `◆ Very high`, `▲ High`,
  `● Medium`, and `○ Low` wherever recommendation impact appears in generated
  Markdown reports.
- Unknown recommendation impact values are humanized without shape markers.
- Run-level `report.md` Top Recommendations includes an `Impact` column between
  `Area / Factors` and `Reason`.
- Recommendation index and detail reports use the same impact display labels.
- Area / Factor Breakdown rows use `▦` for Areas and `□` for Factors.
- Area / Factor Breakdown indentation, links, row order, and root-row bolding
  are preserved.
- Durable report specs and report design guidance are aligned.
- Generated report gallery output is regenerated.
- Focused Go tests pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define recommendation impact
  display labels, Top Recommendations `Impact` column, and Area / Factor
  Breakdown row markers (requirements 1-6).

### To rename

None

### To delete

None
