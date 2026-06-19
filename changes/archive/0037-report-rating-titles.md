---
type: Change Case
title: Render display titles in evaluation reports
description: Make the human evaluation reports display model, target, factor, and rating-level titles while preserving stable ids in report.json and gates.
status: Done
tags: [cli, evaluation, report, rating-scale]
timestamp: 2026-06-19T00:00:00Z
---

# Render display titles in evaluation reports

A **Change Case** for making `qualitymd`'s human evaluation reports display
Model, Target, Factor, and Rating Level `title` values rather than only stable
ids, so that authors can put human labels and emojis in the model and have them
improve report readability. The detail lives in its
[functional spec](0037-report-rating-titles/spec.md). No design doc: the
approach is a localized renderer change with no design alternatives worth
recording.

## Motivation

Human-facing titles are required on Models, Targets, Factors, and Rating Levels
because they make reports easier to scan — a reviewer reads `API Service` or
`🔴 Unacceptable` faster than `api-service` or `unacceptable`. The format already
permits decorated rating-level titles: a Rating Level's `title` is only
constrained to be a "non-empty scalar," and
[`specs/cli/evaluation-build-report.md`](../../specs/cli/evaluation-build-report.md)
already says human-facing renderers **SHOULD** use Model, Target, Factor, and
Rating Level `title` values as primary display labels.

But the renderer did not consistently implement that contract. `report.md` and
`report-summary.md` printed stable target keys, factor keys, and rating `level`
ids in places where human titles were available. Emojis placed in rating titles
therefore never reached a report, and target/factor titles were left unused. This
case closes the conformance gap and dogfoods decorated rating titles in the
project's own [`QUALITY.md`](../../QUALITY.md).

## Scope

Covered:

- The human report renderers (`report.md` and `report-summary.md`) display
  Model, Target, Factor, and Rating Level titles as primary labels, falling back
  to stable ids when a run's snapshot has no title. Stable target keys, factor
  keys, target paths, and rating `level` ids stay in `report.json` and in the
  gate/`BuildResult` so machine consumers and `--fail-at-or-below` are
  unaffected.
- Emoji titles added to the four rating levels in this repo's `QUALITY.md`, as
  the encouraged-convention dogfood.

Deferred / non-goals:

- No emoji documentation push: `SPECIFICATION.md` and the authoring/runtime
  guides are not amended to promote emoji titles (the format already permits
  them; the convention can be documented later if desired).
- No change to `report.json` output or the JSON contract; ratings remain `level`
  ids there.
- No new schema, lint rule, or constraint on title characters.

## Affected specs & docs

- [x] [`specs/cli/evaluation-build-report.md`](../../specs/cli/evaluation-build-report.md)
      — already prescribes title-as-label for human renderers; clarify that this
      covers Model, Target, Factor, and Rating Level labels and that
      `report.json` keeps stable ids for gates/machine consumers.
- [x] [`specs/evaluation-records.md`](../../specs/evaluation-records.md) —
      clarify that record payloads and `report.json` preserve stable model ids
      while human Markdown reports resolve display titles from the model
      snapshot.
- [x] `internal/evaluation/report.go` — render the rating-level title in
      `report.md` and `report-summary.md` via Model, Target, Factor, and Rating
      Level title lookups.
- [x] `internal/evaluation/evaluation_test.go` — update the report assertions
      that expect stable ids in human output to expect display titles while
      preserving stable ids in JSON.
- [x] [`QUALITY.md`](../../QUALITY.md) — add emoji titles to the four rating levels.

No `SPECIFICATION.md` update is expected: emoji titles are already valid under
the existing "non-empty scalar" rule for `title`, and this case changes report
rendering, not the document format.

## Children

- [Functional spec](0037-report-rating-titles/spec.md) — what the report
  renderers must display.

## Status

`Done`. Implementation, durable artifact synchronization, review, and archival
are complete; `go test ./...`, `go vet ./...`, and `dprint check` pass. See the
[status lifecycle](../index.md#status-lifecycle).
