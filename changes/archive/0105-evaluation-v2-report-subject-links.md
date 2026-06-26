---
type: Change Case
title: Evaluation v2 report subject links
description: Move generated report navigation links from generic Details columns into the subject cells of Evaluation v2 Markdown report tables.
status: Done
tags: [cli, evaluation, reports, navigation]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report subject links

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0105-evaluation-v2-report-subject-links/spec.md) - what the
  change must do.
- [Design doc](0105-evaluation-v2-report-subject-links/design.md) - how it is
  built, and why.

## Motivation

Evaluation v2 Area and Factor reports currently put generated report navigation
behind a repeated `Details` column. Every row's `details` link points to the row
subject's own report, so the column adds width without adding information. It
also makes the most meaningful cell - the Factor, Area, or Requirement name -
look like inert text even though that is the thing readers naturally want to
open.

The human Markdown report tree should keep the same navigation targets while
placing each generated report link on the row subject itself. Machine data links
and in-page finding details remain explicitly labeled because they are different
kinds of navigation.

## Scope

Covered:

- Evaluation v2 Area report tables for local root Factors, direct child Areas,
  and local Requirements.
- Evaluation v2 Factor report tables for direct Requirements and direct child
  Factors.
- Empty-state rows and table headers affected by removing `Details` columns.
- The durable Evaluation v2 report-tree spec, reference sketch, and focused
  report tests.

Deferred / non-goals:

- No change to breadcrumbs, parent links, compact linked Factor lists, or `Data`
  rows.
- No change to Requirement report finding sections; `Finding Details` is an
  in-page content section, not a generated child report link.
- No change to routine JSON, generated `EvaluationOutputResult`, report-build
  receipts, report file paths, Rating Level display, or enum display titles.
- No broader Markdown table redesign.

## Affected artifacts

### Code

- [x] `internal/evaluation/report_v2.go` - render generated report navigation on
      subject cells and remove redundant `Details` columns from affected tables.
- [x] `internal/evaluation/evaluation_test.go` - assert linked subject cells,
      removed `Details` headers, and unchanged machine-readable outputs.

### Format spec

- [x] None - `SPECIFICATION.md` does not govern Evaluation v2 report table
      presentation. (Deliberate.)

### Durable specs

- [x] `specs/evaluation-v2/reports/report-tree.md` - require subject-cell links
      for row subjects that have exactly one generated human report target, and
      distinguish those from explicit machine data links.

### Durable docs / bundled skill

- [x] `evaluation-v2-sketch.md` - reconcile older starter report table examples
      that still show `Details` columns, or explicitly mark the sketch as
      superseded by the durable report-tree spec.
- [x] None for the bundled `/quality` skill - this is generated report
      rendering, not agent workflow guidance.

### Suggested new durable specs

- None.

## Status

`Done`. Implemented linked subject cells, removed redundant generated-report
`Details` columns, updated focused report tests, synced the durable report-tree
spec and reference sketch, and archived. See the
[status lifecycle](../index.md#status-lifecycle).
