---
type: Functional Specification
title: Evaluation report rating labels - functional spec
description: Requirements for labeling Evaluation v2 report rating columns explicitly and for rendering a local rating beside a descendant-inclusive sub rating in the Factor and Sub-Area breakdown tables.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation report rating labels - functional spec

Companion to the
[Evaluation report rating labels](../0111-evaluation-report-rating-labels.md)
change case. This spec states what the report rating-column change must do. The
durable contract it lands in is the
[Evaluation v2 report tree](../../../specs/evaluation-v2/reports/report-tree.md)
spec.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Evaluation v2 reports name ratings inconsistently. Header tables label their two
rating columns `Overall` and `Local` — bare adjectives sitting beside
self-describing columns (`Confidence`, `Status`, `Data`), so a reader has to
infer the missing noun. Worse, the breakdown tables that list a node's children
were meant to show each child's local rating beside its descendant-inclusive
roll-up rating, but the implementation puts the roll-up in the `Rating` column
and a `Yes`/`No` boolean in the `+ Sub-Factors` / `+ Sub-Areas` column, leaving
the local rating unshown and a rating slot occupied by a boolean. Clean-break
case 0097 already recorded the breakdown problem as an unmet requirement. This
change makes every rating column name what it rates and restores the intended
local-vs-roll-up breakdown.

## Scope

Covered: generated Evaluation v2 Area and Factor report header summary table
column labels; the Area report Factors and Sub-Areas tables and the Factor
report child Factors table, including their column labels and the rating each
renders; and the durable contracts and sketch that describe those tables.

Not covered: report content, ratings, analysis, findings, drivers, navigation,
links, paths; Requirement tables; structured routine data, `EvaluationOutputResult`,
and JSON field names; and migration of existing runs.

## Requirements

### Header rating labels

- The Area report header summary table **MUST** label its descendant-inclusive
  rating column `Overall Rating` and its local rating column `Local Rating`.

- The Factor report header summary table **MUST** label its descendant-inclusive
  rating column `Overall Rating` and its local rating column `Local Rating`.

  > Rationale: the adjacent header columns are self-describing nouns; bare
  > `Overall` / `Local` force the reader to supply the missing noun. — 0111

### Breakdown rating columns

The Area report Factors table, the Area report Sub-Areas table, and the Factor
report child Factors table each list a subject's children with one row per child.

- Each breakdown table **MUST** render a `Local Rating` column from the child's
  `localAnalysis` rating.

- The Area report Factors table and the Factor report child Factors table
  **MUST** render a `+ Sub-Factors Rating` column, and the Area report Sub-Areas
  table **MUST** render a `+ Sub-Areas Rating` column, from the child's
  `localAndDescendantAnalysis` rating.

- When a child has no descendant Factors (for a Factor row) or no descendant
  Areas (for an Area row), the `+ Sub-Factors Rating` / `+ Sub-Areas Rating`
  cell **MUST** render an em dash (`—`) rather than repeating the local rating.

  > Rationale: with no descendants there is no roll-up distinct from the local
  > rating; the em dash preserves the old boolean's "has children" signal
  > without presenting a redundant rating. — 0111

- No breakdown table **MUST** render a `Yes`/`No` boolean in a rating column.

  > Rationale: clean-break case 0097 required Area tables to distinguish local
  > from aggregate rating instead of using boolean `+ Sub-Areas` / `+ Sub-Factors`
  > columns where a rating belongs; this change satisfies that requirement. —
  > 0097, 0111

### Unaffected presentation

- Report content, ratings, links, navigation, paths, structured routine data,
  and `EvaluationOutputResult` **MUST NOT** change as part of this label change.

- The Requirements tables **MUST NOT** change; their `Rating` column has no
  local/aggregate split.

## Acceptance Criteria

- The Area report header table reads `| Overall Rating | Local Rating |
Confidence | Data |`.
- The Factor report header table reads `| Overall Rating | Local Rating |
Status | Confidence | Data |`.
- The Area report Factors table header reads `| Factor | Path | Local Rating |
  - Sub-Factors Rating | Sub-Factors |`, with`Local Rating`sourced from`localAnalysis`and`+ Sub-Factors Rating`from`localAndDescendantAnalysis`.
- The Area report Sub-Areas table header reads `| Area | Path | Local Rating |
  - Sub-Areas Rating | Factors |`, sourced the same way.
- The Factor report child Factors table header reads `| Factor | Path |
Local Rating | + Sub-Factors Rating |`, sourced the same way.
- A breakdown row for a child with no descendants renders `—` in its
  `+ Sub-X Rating` cell.
- No generated report renders a `✅ Yes` / `⬜ No` boolean in a rating column.
- Report paths, links, navigation labels, and structured `data/` outputs are
  unchanged.

## Durable spec changes

### To add

None.

### To modify

- [`specs/evaluation-v2/reports/report-tree.md`](../../../specs/evaluation-v2/reports/report-tree.md)
  - state the header rating column labels and require the Factors, Sub-Areas,
    and child Factors breakdown tables to render a local rating column beside a
    descendant-inclusive `+ Sub-X Rating` column. Driven by
    [Header rating labels](#header-rating-labels) and
    [Breakdown rating columns](#breakdown-rating-columns).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, each Evaluation v2 report rating column
names what it rates: header tables read `Overall Rating` / `Local Rating`, and
every breakdown row shows the child's local rating beside its descendant-inclusive
roll-up (or an em dash when there is no roll-up). That removes the labeling
ambiguity and the boolean-in-a-rating-column regression named in the motivation,
without touching report content, structured data, or navigation.
