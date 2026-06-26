---
type: Functional Specification
title: Evaluation v2 report enum display titles - functional spec
description: Requirements for rendering owned Evaluation v2 enum-like values as human titles in Markdown reports without changing machine data.
tags: [cli, evaluation, reports, display]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report enum display titles - functional spec

Companion to the
[Evaluation v2 report enum display titles](../0103-v2-report-enum-display-titles.md)
change case. This spec states the delta for Evaluation v2 Markdown report
rendering. It defers the cumulative Evaluation v2 report contract to
[`specs/evaluation-v2/reports/report-tree.md`](../../specs/evaluation-v2/reports/report-tree.md)
and the Rating Level display contract to
[`SPECIFICATION.md`](../../SPECIFICATION.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Evaluation v2 Markdown reports are the human review surface over structured
routine JSON. Raw enum-like values such as `not_analyzed`, `partially_assessed`,
and `medium` are useful stable data values, but they are harder to scan in
status-heavy Markdown tables than human titles. Rating Level titles already
establish the right boundary: machine data keeps stable IDs, while Markdown
reports render human labels.

## Scope

Covered:

- Evaluation v2 Markdown reports for Area, Factor, and Requirement pages.
- CLI-owned report-display vocabularies: scoped analysis statuses, requirement
  assessment statuses, requirement rating statuses, confidence levels, boolean
  display values, limit/unknown table types, report output kinds, and known
  finding severities/types.
- Fallback behavior for unknown values.
- Tests proving human Markdown titles and machine JSON stability.
- Durable report-tree spec sync.

Deferred / non-goals:

- No changes to routine JSON payload shapes, validation rules, generated
  `EvaluationOutputResult`, or report-build receipt JSON.
- No changes to Rating Level semantics or title resolution. Produced Rating
  Levels continue to resolve through the run's `model.md` snapshot.
- No closed enum contract for finding `type` or `severity`; those fields remain
  open unless a later change constrains them.
- No presentation change outside Evaluation v2 Markdown reports.

## Requirements

- Evaluation v2 Markdown reports **MUST** render CLI-owned status and confidence
  values through human titles with semantic emoji instead of raw enum strings.

  > Rationale: status and confidence values are high-signal scan targets in
  > reports; titles make the human report surface clearer without altering the
  > data contract. - 0103

- Evaluation v2 Markdown reports **MUST** keep produced Rating Level labels
  resolved through the snapshotted model Rating Level `title`, not through the
  enum display catalog.

  > Rationale: Rating Levels are model-owned display vocabulary, while this
  > change covers CLI-owned report states. - 0103

- Evaluation v2 Markdown reports **MUST** fall back to a readable humanized value
  when a report field contains an unknown or free-form string.

  > Rationale: historical, hand-edited, or future routine outputs can contain
  > values this release does not know; report generation should stay readable and
  > deterministic rather than fail for presentation-only decoration. - 0103

- Evaluation v2 report display catalogs **SHOULD** be typed for vocabularies the
  Evaluation v2 package owns.

  > Rationale: typed catalogs make it harder to mix payload kinds, report kinds,
  > statuses, and confidence values while still allowing loose JSON boundaries to
  > use safe fallback conversion. - 0103

- Evaluation v2 Markdown reports **MAY** apply display titles to known finding
  `type` and `severity` strings, but **MUST NOT** reject or erase unknown finding
  values.

  > Rationale: findings are currently open routine-output content, so display
  > polish must not silently tighten the schema. - 0103

- Evaluation v2 machine-readable outputs **MUST NOT** replace stable enum values
  with display titles in routine JSON, `data/evaluation-output-result.json`, or
  the `qualitymd evaluation report build --json` receipt.

  > Rationale: agents and tools join on stable values; emoji titles are only a
  > human Markdown projection. - 0103

- Report tests **MUST** cover titled Markdown statuses/confidence and at least
  one unchanged machine-readable stable value.

## Durable spec changes

### To add

None

### To modify

- [`specs/evaluation-v2/reports/report-tree.md`](../../specs/evaluation-v2/reports/report-tree.md)
  - require generated Markdown reports to render CLI-owned enum-like values with
    display titles while preserving raw values in structured outputs.

### To rename

None

### To delete

None
