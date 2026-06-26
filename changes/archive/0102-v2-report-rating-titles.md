---
type: Change Case
title: Evaluation v2 report rating titles
description: Restore Rating Level title rendering in Evaluation v2 human Markdown reports while preserving stable IDs in machine outputs.
status: Done
tags: [cli, evaluation, reports, rating-scale]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report rating titles

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0102-v2-report-rating-titles/spec.md) - what the change must
  do.
- [Design doc](0102-v2-report-rating-titles/design.md) - how it is built, and
  why.

## Motivation

Evaluation v2 reworked report generation around a Markdown report tree and
structured routine JSON. In that transition, Rating Level display titles stopped
rendering in human Markdown report rating cells: `target` appears where a reader
should see the model's `title`, such as `Target` or `🔵 Target`.

That loses the readability payoff of required Rating Level titles and regresses
the established display contract from earlier report work. Evaluation data and
machine receipts still need stable Rating Level IDs; only human report labels
need title resolution.

## Scope

Covered:

- Evaluation v2 Markdown report rating labels resolve Rating Level titles from
  the run's `model.md` snapshot, falling back to stable IDs only when a title is
  unavailable.
- Evaluation v2 machine-readable payloads and build receipts keep stable Rating
  Level IDs.
- The durable Evaluation v2 report-tree spec explicitly carries the title
  rendering rule.
- Report tests prove the human Markdown/title and machine ID split.

Deferred:

- No legacy `report-summary.md`, legacy `report.json`, report gate, or
  `--fail-at-or-below` behavior is restored.
- No new schema fields, lint rules, rating semantics, emoji requirements, or
  changes to structured routine JSON.
- No broad report layout redesign beyond the rating label bug.

## Affected artifacts

### Code

- [x] `internal/evaluation/report_v2.go` - resolve Rating Level titles for human
      Markdown report labels.
- [x] `internal/evaluation/evaluation_test.go` - assert Markdown uses titles
      while machine output keeps IDs.

### Format spec

- [x] None - `SPECIFICATION.md` already permits and defines Rating Level display
      values.

### Durable specs

- [x] `specs/evaluation-v2/reports/report-tree.md` - make Rating Level title
      resolution explicit for Evaluation v2 Markdown reports.

### Durable docs / bundled skill

- [x] None - existing `/quality` skill guidance already says required titles are
      the primary human-facing labels, and this change is confined to CLI report
      rendering.

### Suggested new durable specs

- None.

## Status

`Done`. Spec and design settled, Evaluation v2 report rendering fixed, durable
report-tree spec updated, tests and formatting verified, and archived. See the
[status lifecycle](../index.md#status-lifecycle).
