---
type: Design Doc
title: Evaluation v2 report subject links - design
description: Design for moving Evaluation v2 report table navigation links from Details columns into subject cells.
tags: [cli, evaluation, reports, navigation]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report subject links - design

## Context

This design answers the [functional spec](spec.md) for moving Evaluation v2
generated report navigation from generic `Details` columns into the row subject
cells. The existing renderer already has a small, deterministic Markdown writer
in `internal/evaluation/report_v2.go`; the change is a presentation adjustment
inside that writer, not a data-model or report-tree change.

This case is separate from
[0104 - Evaluation v2 report header navigation](../0104-evaluation-v2-report-header-navigation.md):
0104 changes report headers and trails, while this case changes row links in
Area and Factor report body tables.

## Approach

Update the affected table renderers in `renderV2AreaReport` and
`renderV2FactorReport`:

1. Replace subject-cell plain text with `reportLink(...)` calls whose label is
   the same title currently rendered in that subject cell.
2. Remove the `Details` header cell and separator cell from affected tables.
3. Remove the trailing `details` cell from each affected data row.
4. Reduce empty-state rows to the updated column counts.
5. Leave `Data` field-table rows, breadcrumbs or 0104's successor trails,
   parent links or their 0104 replacements, compact Factor lists, Requirement
   attached-Factor links, and finding sections unchanged.

The current `reportLink` helper already computes relative paths from the current
report to the target report and escapes the link label through `markdownCell`.
Using it in the subject cell preserves the existing target path behavior while
changing only where the link appears.

## Spec response

The Area report requirements map directly to the three tables in
`renderV2AreaReport`: `Factors`, `Sub-Areas`, and `Requirements`. Each row
already computes the same target path needed for its `details` cell, so the
target can move to the subject cell without new lookup logic.

The Factor report requirements map to the `Requirements` and `Child Factors`
tables in `renderV2FactorReport`. These use the same helper paths as the Area
report tables, so the implementation stays local to table header and row
construction.

Tests should extend the focused Evaluation v2 report-build coverage with
fixtures that generate at least one affected link in an Area report and one in a
Factor report. Assertions should check for linked subject labels, absence of the
affected `Details` headers, and stable `EvaluationOutputResult` values.

## Alternatives

### Keep `Details` for explicitness

Rejected. A repeated `details` link is explicit but low-signal. The subject
title is already the object being opened, and the `Path` column remains
available for disambiguation.

### Link both the subject and `Details`

Rejected. Duplicating the same target in one row keeps the extra width and makes
keyboard/link traversal noisier without adding navigation value.

### Move links into the `Path` column

Rejected. The path is useful as an identifier, but it is less readable than the
title and would make the table feel more technical than the surrounding report
surface.

### Make all table cells link to the report

Rejected. Full-row links are not natural Markdown table output and would blur
the difference between navigation cells and summary/status cells.

## Trade-offs & risks

Removing `Details` makes the table narrower and clearer, but it also removes a
textual cue that a row has a child report. The linked subject label is the cue;
tests should assert the Markdown link shape so the affordance does not regress
back to plain text.

The main implementation risk is missing one affected table or leaving an
empty-state row with the old number of cells. The renderer has only five affected
tables, so targeted tests plus a search for `Details` in `report_v2.go` are
enough to catch this.

0104 may change the header area around these tables before this case is
implemented. This design avoids depending on header details so the two changes
can land in either order as long as report body tables keep the same semantic
sections.

## Open questions

None.
