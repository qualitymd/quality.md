---
type: Functional Specification
title: Evaluation v2 report rating titles - functional spec
description: Requirements for rendering Rating Level titles in Evaluation v2 Markdown reports while preserving IDs in structured outputs.
tags: [cli, evaluation, reports, rating-scale]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report rating titles - functional spec

Companion to the
[Evaluation v2 report rating titles](../0102-v2-report-rating-titles.md) change
case. This spec states the delta for Evaluation v2 report rendering. It defers
the general Rating Level schema and display-value vocabulary to
[`SPECIFICATION.md`](../../SPECIFICATION.md), and the cumulative Evaluation v2
report contract to
[`specs/evaluation-v2/reports/report-tree.md`](../../specs/evaluation-v2/reports/report-tree.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Human reports are the review surface for evaluation results, so their rating
labels need to use the model's human-facing vocabulary. Evaluation v2 currently
prints stable Rating Level IDs in Markdown rating cells, even though every
Rating Level has a required `title` for display. Stable IDs remain the right
machine identity in routine JSON and build receipts; the gap is only the
Markdown presentation layer.

## Scope

Covered:

- Evaluation v2 Markdown rating labels in the generated report tree.
- Fallback behavior when an older or malformed model snapshot lacks a title for
  a selected Rating Level.
- Tests for the split between human Markdown titles and machine-readable IDs.
- Durable report-tree spec sync.

Deferred / non-goals:

- No changes to Rating Level semantics, ordering, schema validation, lint rules,
  or default scaffold titles.
- No changes to routine JSON payload shapes, `EvaluationOutputResult`, or build
  receipt JSON.
- No restoration of legacy `report-summary.md`, legacy `report.json`, report
  gates, or old report model APIs.

## Requirements

- Evaluation v2 Markdown reports **MUST** render the selected Rating Level's
  `title` as the primary label wherever an Area, Factor, or Requirement rating
  cell displays a produced Rating Level.

  > Rationale: the title is the display label authors maintain for readers;
  > printing only `target` or `minimum` discards the scanning aid and repeats the
  > report regression that display titles were introduced to prevent. - 0102

- Evaluation v2 Markdown reports **MUST** resolve Rating Level titles from the
  run's snapshotted `model.md`, keyed by Rating Level `level`.

  > Rationale: reports are historical projections of a run, so they must not
  > pick up later edits from the current workspace `QUALITY.md`. - 0102

- If the snapshotted Rating Scale has no title for a selected Rating Level, the
  Markdown renderer **MUST** fall back to the stable Rating Level ID rather than
  rendering an empty label or failing report generation.

  > Rationale: historical or hand-edited runs can lack title data; fallback keeps
  > the report readable without changing the recorded rating. - 0102

- Evaluation v2 Markdown reports **MUST** continue to render non-rating statuses
  such as `not_rated`, `not_analyzed`, `empty`, and `blocked` as statuses, not as
  Rating Level titles.

- Evaluation v2 machine-readable outputs **MUST NOT** replace Rating Level IDs
  with titles in routine JSON, `data/evaluation-output-result.json`, or the
  `qualitymd evaluation report build --json` receipt.

  > Rationale: IDs are the durable join key for agents and tools; titles are
  > presentation. - 0102

- The report tests **MUST** cover both sides of the split: Markdown shows the
  Rating Level title, and at least one machine-readable report-build output still
  carries the stable Rating Level ID.

## Durable spec changes

### To add

None

### To modify

- [`specs/evaluation-v2/reports/report-tree.md`](../../specs/evaluation-v2/reports/report-tree.md)
  - make the Rating Level title-resolution rule explicit for Markdown reports,
    driven by the requirements above.

### To rename

None

### To delete

None
