---
type: Design Doc
title: Report Visual Markers - design doc
description: Design for generated report impact labels and Area / Factor Breakdown markers.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Report Visual Markers - design doc

Design for [Report Visual Markers](../0170-report-visual-markers.md) and its
[functional spec](spec.md).

## Context

Recommendation reports currently call `impactTitle` from the report renderer,
which returns plain text for known impact values. Confidence and finding severity
already have emoji-prefixed display catalogs, so using color for impact would
make adjacent report columns harder to interpret. The run report's compact Top
Recommendations table also omits impact, even though the full recommendation
index and detail reports include it.

Area / Factor Breakdown Factor rows currently derive their marker from the
Factor report kind title, which yields `🧩`. Area rows use only links and
indentation. That asymmetry makes Factors visually loud and Areas visually
unmarked in a table whose heading already names both row types.

## Approach

Move recommendation impact display into the shared display helper layer near the
other report display labels. Keep the known-value mapping small and explicit:

- `very_high` -> `◆ Very high`;
- `high` -> `▲ High`;
- `medium` -> `● Medium`;
- `low` -> `○ Low`.

Unknown impact values continue to flow through the existing humanization helper.
All report locations that render recommendation impact keep calling
`impactTitle`, so the recommendation index, recommendation detail reports, and
the restored run-report Top Recommendations column stay aligned.

Restore the run report Top Recommendations `Impact` column by adding the display
label to `writeTopRecommendationsTable` between Area / Factors and Reason.
Adjust the empty row so it preserves the new column count.

Localize Area / Factor Breakdown markers to the breakdown-title helpers rather
than changing global report kind labels. Add an Area marker helper returning
`▦` and replace `factorReportMarker` with a Factor marker helper returning `□`.
The row title builders prepend the marker after indentation and inside the
visible label content, preserving the current link targets, ordering, and root
row bolding.

Regenerate the report gallery after renderer changes so public examples remain
the stale-output check for report design.

## Spec response

- Shared `impactTitle` gives every recommendation impact report location one
  display contract.
- `writeTopRecommendationsTable` restores impact to the compact run report.
- Area / Factor Breakdown title helpers apply `▦` and `□` without changing
  global report kind labels.
- Existing fallback humanization handles unknown impact values.
- Tests and generated gallery output verify the visible Markdown contract.

## Alternatives

- **Use red/orange/yellow/blue impact markers.** Rejected because red reads as
  badness or urgent concern, while impact is expected quality-management value.
- **Reuse confidence colors for impact.** Rejected because impact and confidence
  often appear side by side and need distinct semantics.
- **Omit Area / Factor markers entirely.** Rejected because a subtle type marker
  helps readers scan mixed Area and Factor rows without adding a separate Type
  column.
- **Change global Area and Factor report kind labels.** Rejected because the
  current problem is the dense breakdown table, not every report-kind display.

## Trade-offs & risks

The chosen markers are typographic glyphs rather than colorful emoji. They are
calmer, but their meanings are learned from context rather than universally
obvious. Keeping them local to the table and preserving the `Area / Factor`
heading limits that risk.

Using `□` for Factors can read as an unchecked box. That connotation is
acceptable here because Factors are quality lenses that collect checkable
Requirements, but the report does not use the marker to indicate completion.

## Open questions

None.
