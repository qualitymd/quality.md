---
type: Functional Specification
title: Render display titles in evaluation reports — functional spec
description: Requirements for displaying model, target, factor, and rating-level titles in report.md and report-summary.md while keeping stable ids in report.json and gates.
tags: [cli, evaluation, report, rating-scale]
timestamp: 2026-06-19T00:00:00Z
---

# Render display titles in evaluation reports — functional spec

Companion to the
[Render display titles in evaluation reports](../0037-report-rating-titles.md)
change case. This spec states _what_ the report renderers must display. It defers
the broad `build-report` contract to
[`specs/cli/evaluation-build-report.md`](../../../specs/cli/evaluation-build-report.md)
and the rating-scale schema to [`SPECIFICATION.md`](../../../SPECIFICATION.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Display titles exist on Models, Targets, Factors, and Rating Levels (required
since the display-titles change) so authors can give reports readable labels.
Rating Level titles are constrained only to be a non-empty scalar, so authors can
decorate them — most usefully with emojis — to make reports scannable. The
`build-report` spec already says human renderers SHOULD use these `title` values
as primary display labels. The renderer did not consistently do that: it printed
target keys, factor keys, and rating `level` ids in places where titles were
available. This spec pins down displaying titles in the human reports while
keeping stable identifiers wherever machines and gates depend on them.

## Scope

Covered: how `report.md` and `report-summary.md` label Models, Targets, Factors,
and Rating Levels. Out of scope: changing `report.json` content, the CI gate,
and any new constraint on what characters a title may contain.

## Requirements

- In `report.md` and `report-summary.md`, every place that displays a Model,
  Target, Factor, or Rating Level label **MUST** show that entity's `title` as
  the primary label when the run's model snapshot contains one.

  > Rationale: a display title (e.g. `API Service` or `🔴 Unacceptable`) is the
  > readability payoff; printing the stable id discards it. — 0037

- The renderer **MUST** resolve Model, Target, and Factor titles from the run's
  recorded model snapshot by stable target path and factor key. When the recorded
  model has no title for a target or factor, it **MUST** fall back to the stable
  identifier rather than rendering an empty label.

- The renderer **MUST** resolve a rating's title from the run's recorded rating
  scale by `level` id. When the recorded scale has no title for a level (for
  example a run snapshotted before titles were required), it **MUST** fall back
  to the `level` id rather than rendering an empty label.

- Non-rating rating-cell states **MUST** be preserved verbatim: a not-assessed
  result still renders as `not assessed` and a structural result still renders as
  its structural marker. These are not rating levels and have no title.

- `report.json` **MUST** continue to carry stable target keys, factor keys,
  target paths, and rating `level` ids (not titles), and the `BuildResult` rating
  and the `--fail-at-or-below` gate **MUST** continue to operate on `level` ids.
  Title rendering is confined to the two human Markdown reports.

  > Rationale: ids are the stable join key for gates and machine consumers;
  > titles are presentation. Keeping ids in JSON preserves byte-stable,
  > diffable, gate-safe output. — 0037

- Rendering **MUST** remain deterministic and idempotent: unchanged records and
  scale produce byte-identical report files.

## Durable spec changes

### To add

None.

### To modify

- [`specs/cli/evaluation-build-report.md`](../../../specs/cli/evaluation-build-report.md)
  — the existing SHOULD that human renderers use Model, Target, Factor, and
  Rating Level `title` values stays; clarify that titles are the primary labels
  in `report.md` / `report-summary.md`, that the renderer falls back to stable
  ids when a snapshot lacks a title, and that `report.json` and the gate keep
  stable ids.
- [`specs/evaluation-records.md`](../../../specs/evaluation-records.md) —
  clarify that evaluation record payloads and `report.json` preserve stable
  model identifiers, while human Markdown reports resolve display labels from the
  run's model snapshot.

### To delete

None.
