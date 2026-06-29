---
type: Design Doc
title: Recommendation Number Columns — design doc
description: Design for removing duplicate Rank/# columns from recommendation list tables.
tags: [evaluation, reports, recommendations]
timestamp: 2026-06-29T00:00:00Z
---

# Recommendation Number Columns — design doc

Design for
[Recommendation Number Columns](../0178-recommendation-number-columns.md) and
its [functional spec](spec.md).

## Context

The renderer already computes `rankedRecommendation.Rank` from
`RecommendationRankingResult.orderedRecommendations[].rank`. Both current table
renderers pass that same value twice: once as `Rank` and once as `#`.

## Approach

Change only the Markdown table rendering for recommendation list tables:

- `writeTopRecommendationsTable` drops the `Rank` header and its row cell.
- `writeRecommendationIndexTable` drops the `Rank` header and its row cell.
- Empty-state rows lose one cell to match the new column count.

The underlying `rankedRecommendation.Rank`, detail report paths, source data,
and structured ranking payloads remain unchanged.

## Spec response

The chosen renderer-only change satisfies the report-table requirements without
touching structured Evaluation data. Tests and report-gallery output verify the
human Markdown surface.

## Alternatives

- **Rename `Rank` to `Number` and drop `#`.** Rejected because `#` is compact and
  already used across generated reports as the recommendation number.
- **Keep both columns but hide one in CSS later.** Rejected because these are
  Markdown artifacts; the source table should be unambiguous without styling.

## Trade-offs & risks

None beyond fixture churn in generated reports.

## Open questions

None.
