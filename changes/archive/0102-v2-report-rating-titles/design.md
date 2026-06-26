---
type: Design Doc
title: Evaluation v2 report rating titles - design
description: Design for resolving Rating Level display titles in Evaluation v2 Markdown reports.
tags: [cli, evaluation, reports, rating-scale]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report rating titles - design

## Context

This design answers the
[functional spec](spec.md) for restoring Rating Level title rendering in
Evaluation v2 Markdown reports. The current report renderer already passes the
run's snapshotted `model.Spec` into the Area, Factor, and Requirement report
renderers, and already resolves Area, Factor, and Requirement titles from that
snapshot. Rating labels are the outlier: the helpers receive only raw result
maps, so they return `ratingLevelId` directly.

## Approach

Keep the fix local to the Evaluation v2 renderer:

1. Add a Rating Level title lookup helper near the existing Area, Factor, and
   Requirement title helpers.
2. Thread `*model.Spec` into the rating-label helpers that render Markdown
   cells.
3. Preserve the current status handling: produced ratings resolve through the
   title helper; non-rating statuses render as status markers.
4. Leave `v2ReceiptRating` and `v2OutputResult` unchanged so JSON outputs keep
   stable IDs.
5. Update the report-tree durable spec and the focused report test.

The helper should scan `spec.RatingScale` by `level`, return the non-empty
`title`, and fall back to the `level` ID. It should tolerate a nil or incomplete
model snapshot because report generation already handles historical and
hand-edited runs more gracefully than the main model authoring path.

## Spec response

Title rendering is handled at the final Markdown projection boundary, so
structured routine data remains untouched. Because the report builder loads
`model.md` from the run folder before rendering, resolving through `model.Spec`
uses the historical run snapshot rather than the current workspace model.

The fallback behavior is explicit in one helper and reused across Area, Factor,
Requirement, nested list, and table contexts. Tests should use a title that is
visibly different from the ID so a regression to `ratingLevelId` fails.

## Alternatives

### Store titles in routine JSON

Rejected. Routine JSON is a machine-readable evaluation record and already uses
structural IDs for Areas, Factors, Requirements, and Rating Levels. Duplicating
titles there would create drift when the run snapshot already contains the
authoritative display vocabulary.

### Change `RatingDisplay(level)` to use titles

Rejected for this change. `RatingDisplay(level)` has no model snapshot argument,
so it cannot resolve a title. Changing that helper's signature would broaden a
localized report-rendering fix into the model-reference API without an active
caller need.

### Render `title (id)`

Rejected. The established display-title contract says the title is the primary
human label, with IDs retained in machine outputs and explicit reference
contexts. Adding IDs to every rating cell would make reports noisier without
solving a traceability gap.

## Trade-offs & risks

This keeps the behavior localized but leaves two display concepts in the code:
reference rendering helpers that return IDs and report rendering helpers that
resolve titles. The risk is manageable because the report helpers are in the
renderer next to the existing Area, Factor, and Requirement title helpers.

The main regression risk is missing a Markdown call site. A source search for
`v2RatingLabel`, `v2RequirementRatingLabel`, and `ratingLevelId` in
`report_v2.go` is enough to identify the relevant call sites.

## Open questions

None.
