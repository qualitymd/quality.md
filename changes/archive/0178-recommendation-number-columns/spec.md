---
type: Functional Specification
title: Recommendation Number Columns — functional spec
description: Requirements for removing duplicate Rank/# columns from generated recommendation list tables.
tags: [evaluation, reports, recommendations]
timestamp: 2026-06-29T00:00:00Z
---

# Recommendation Number Columns — functional spec

Companion to
[Recommendation Number Columns](../0178-recommendation-number-columns.md). This
spec states the delta for generated recommendation list tables. The durable
source of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md).

The key words **MUST** and **MUST NOT** are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Recommendation users read "recommendation #1" as the first recommendation in
the ranked list. After 0176, generated list tables still show `Rank` and `#`
side by side, with both values derived from the same ranking order. That duplicate
ordinal makes the report noisier and weakens the intended distinction between
user-facing recommendation numbers and opaque recommendation IDs.

## Requirements

1. The run report Top Recommendations table **MUST** render one user-facing
   recommendation number column labeled `#`, followed by `Recommendation`, `Area
/ Factors`, `Impact`, and `Reason`.

   > Rationale: `#` is the common-parlance recommendation number; a second `Rank`
   > column repeats the same value and invites users to distinguish two numbers
   > that no longer differ.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

2. The Recommendations report Ranked Recommendations table **MUST** render one
   user-facing recommendation number column labeled `#`, followed by
   `Recommendation`, `Area / Factors`, `Impact`, `Confidence`, `Reason`, and
   `Ranking Rationale`.

   > Rationale: The complete recommendation list should use the same visible
   > identity as the capped run-report list while preserving ranking rationale.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

3. The `#` value in both tables **MUST** be derived from
   `RecommendationRankingResult.orderedRecommendations[].rank`.

   > Rationale: This preserves the data contract while presenting ranked order
   > as the recommendation number.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

4. Generated recommendation list tables **MUST NOT** render a separate `Rank`
   column when they already render `#`.

   > Rationale: The whole purpose of the change is to keep one visible ordinal.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

## Acceptance criteria

- `report.md` Top Recommendations tables use `# | Recommendation | Area /
Factors | Impact | Reason`.
- `recommendations.md` Ranked Recommendations tables use `# | Recommendation |
Area / Factors | Impact | Confidence | Reason | Ranking Rationale`.
- Generated rows no longer duplicate the recommendation number in adjacent
  `Rank` and `#` cells.
- Recommendation detail reports still show `#` and `ID`; no structured payload
  field is renamed.
- Report-gallery output is regenerated.
- Focused Evaluation tests and `mise run check` pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - remove duplicate `Rank` columns
  from recommendation list table contracts.

### To rename

None

### To delete

None
