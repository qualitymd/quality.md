---
type: Change Case
title: Recommendation Number Columns
description: Remove duplicate Rank/# columns from human recommendation list tables so # is the single visible recommendation number.
status: Done
tags: [evaluation, reports, recommendations]
timestamp: 2026-06-29T00:00:00Z
---

# Recommendation Number Columns

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0178-recommendation-number-columns/spec.md) - what the case must do.
- [Design doc](0178-recommendation-number-columns/design.md) - how it is built, and why.

## Motivation

0176 split opaque recommendation IDs from user-facing recommendation numbers, but
generated recommendation list tables still render both `Rank` and `#`. Because
the two columns now carry the same ranked-order value, the duplicate display
reintroduces the ambiguity the ID/number split removed.

Human recommendation list tables should show one visible ordinal: `#`. The
structured data may continue to call the ordering field `rank`.

## Scope

Covered:

- Remove the `Rank` column from generated recommendation list tables.
- Keep `#` as the user-facing recommendation number derived from ranking order.
- Preserve recommendation detail report filenames, opaque `qrec_...` IDs, and
  structured `rank` fields.

Non-goals:

- Changing Finding ranking tables, structured Evaluation payloads, or
  recommendation detail metadata.
- Renaming `RecommendationRankingResult.orderedRecommendations[].rank`.
- Changing recommendation ordering semantics.

## Affected artifacts

**Code**

- [x] `internal/evaluation/report_tree.go` - render recommendation list tables
      with `#` only.
- [x] `internal/evaluation/evaluation_test.go` - update table assertions.

**Durable specs**

- [x] `specs/evaluation/reports/report-tree.md` - specify recommendation list
      columns without duplicate `Rank`.

**Durable docs / examples**

- [x] `docs/guides/reporting-design.md` - update recommendation report examples
      that still show `Rank` beside `#`.
- [x] `CHANGELOG.md` - note the report polish.
- [x] `examples/report-gallery/software-service/.quality/evaluations/**` -
      regenerate report gallery output.
- [x] `changes/index.md`, `changes/log.md`, and `changes/archive/index.md` -
      Change Case lifecycle.

## Status

`Done`. Implemented and archived. Generated recommendation list tables now use
`#` as the single visible recommendation number column; durable report specs,
tests, docs, changelog, and report-gallery output are aligned.
