---
type: Change Case
title: Evaluation report rating labels
description: Make Evaluation v2 report rating columns name what they rate (Overall Rating / Local Rating) and fix the Factor and Sub-Area breakdown tables to show a local rating plus a descendant-inclusive sub rating instead of an aggregate rating beside a Yes/No boolean.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation report rating labels

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0111-evaluation-report-rating-labels/spec.md) - what the
  change must do.
- [Design doc](0111-evaluation-report-rating-labels/design.md) - how it is
  built, and why.

## Motivation

Evaluation v2 reports refer to ratings inconsistently. Two problems compound:

1. **Ambiguous header labels.** Area and Factor report header tables label their
   two rating columns `Overall` and `Local`. The neighboring columns
   (`Confidence`, `Status`, `Data`) are self-describing nouns, but `Overall` and
   `Local` are bare adjectives with no noun — a reader has to infer "overall
   *what*?". Spelling them `Overall Rating` and `Local Rating` removes the guess
   and makes every rating column name what it rates.

2. **Wrong breakdown columns.** The Area report's Factors and Sub-Areas tables,
   and the Factor report's child Factors table, were meant to show a node's
   local rating beside its descendant-inclusive ("+ Sub-X") rating — the design
   the v2 sketch still documents
   ([`evaluation-v2-sketch.md`](../evaluation-v2-sketch.md) lines 894–901). The
   implementation regressed: the `Rating` column renders the
   `localAndDescendantAnalysis` (aggregate) rating, and the `+ Sub-Factors` /
   `+ Sub-Areas` column renders a `Yes`/`No` boolean where a rating belongs. The
   clean-break case 0097 already recorded this as an unmet **MUST**
   ([`changes/archive/0097-evaluation-v2-clean-break/spec.md`](archive/0097-evaluation-v2-clean-break/spec.md)
   line 332): "Area tables **MUST** distinguish local rating from aggregate
   rating instead of using boolean `+ Sub-Areas` or `+ Sub-Factors` columns
   where a rating belongs."

This case makes the report vocabulary consistent: wherever a `localAnalysis`
rating appears it is labeled `Local Rating`; the descendant-inclusive rating is
labeled `Overall Rating` in a subject's own header and `+ Sub-Factors Rating` /
`+ Sub-Areas Rating` in a breakdown row, where naming the roll-up scope is
clearer.

## Scope

Covered:

- Evaluation v2 Area and Factor report header summary table column labels;
- Evaluation v2 Area report Factors and Sub-Areas breakdown tables, and Factor
  report child Factors table — column labels and the rating each column renders;
- the durable Evaluation v2 report-tree contract for those tables; and
- the v2 report sketch's descriptions of those tables.

Deferred / non-goals:

- no change to report content, ratings, findings, analysis, evidence, limits,
  drivers, or navigation;
- no change to the Requirements tables (their `Rating` column has no
  local/aggregate split);
- no change to structured routine data, `EvaluationOutputResult`, or any JSON
  field names — only generated Markdown presentation;
- no change to report paths, filenames, or links; and
- no migration of existing completed evaluation runs.

## Affected artifacts

### Code

- [x] `internal/evaluation/report_v2.go` - rename header columns to
      `Overall Rating` / `Local Rating`; in the Factors, Sub-Areas, and child
      Factors tables render the `Local Rating` column from `localAnalysis` and a
      `+ Sub-Factors Rating` / `+ Sub-Areas Rating` column from
      `localAndDescendantAnalysis`, rendering an em dash when the node has no
      descendants; remove the now-unused `v2BoolLabel` wrapper.
- [x] `internal/evaluation/evaluation_test.go` - update report assertions to the
      new headers and cells; differentiate a node's local and aggregate ratings
      in the navigation fixture so the column split is verified.

### Format spec

- [x] None - `SPECIFICATION.md` does not govern Evaluation v2 generated report
      table presentation. (Deliberate.)

### Durable specs

- [x] `specs/evaluation-v2/reports/report-tree.md` - state the header column
      labels and require the Factors/Sub-Areas/child-Factors tables to show a
      local rating column and a descendant-inclusive sub-rating column.

### Durable docs / bundled skill

- [x] `evaluation-v2-sketch.md` - align the report table sketches and field
      meanings with the new column labels.

### Suggested new durable specs

- None.

## Status

`Done`. Implemented, verified with `mise run check`, and archived. See the
[status lifecycle](../index.md#status-lifecycle).
