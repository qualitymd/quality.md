---
type: Change Case
title: Evaluation v2 report enum display titles
description: Render owned Evaluation v2 status and display vocabularies with typed human titles and semantic emoji in Markdown reports while preserving machine values.
status: Done
tags: [cli, evaluation, reports, display]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report enum display titles

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0103-v2-report-enum-display-titles/spec.md) - what the
  change must do.
- [Design doc](0103-v2-report-enum-display-titles/design.md) - how it is built,
  and why.

## Motivation

Evaluation v2 Markdown reports contain several enum-like values that are meant
for agent and tool stability: scoped analysis statuses, requirement assessment
statuses, requirement rating statuses, confidence levels, and generated report
references. Those raw values are useful in JSON, but in human reports they make
status-heavy tables harder to scan and less consistent with Rating Level titles.

The report renderer needs a small display-title layer for vocabularies the CLI
owns, while preserving the split introduced for Rating Levels: stable values in
machine data, human titles at the Markdown projection boundary.

## Scope

Covered:

- Evaluation v2 Markdown reports render owned status, confidence, boolean, limit,
  unknown, finding-severity, and report-kind vocabularies through human titles
  with semantic emoji.
- Rating Level labels continue to resolve from the run's `model.md` snapshot
  rather than from the enum display catalog.
- Persisted routine JSON, generated `EvaluationOutputResult`, and report-build
  receipts keep the stable raw values.
- Display catalogs are typed where the vocabulary is owned by the Evaluation v2
  package, with fallback humanization for unknown or free-form string values.
- The durable Evaluation v2 report-tree spec and focused tests cover the display
  boundary.

Deferred:

- No new routine JSON fields, schema constraints, lint rules, or model
  vocabulary.
- No requirement that finding `type` or finding `severity` become closed enums;
  reports may title known values but must still tolerate free-form values.
- No broader report layout redesign.

## Affected artifacts

### Code

- [x] `internal/evaluation/display.go` - typed display catalogs and fallback
      rendering helpers for Evaluation v2 report display values.
- [x] `internal/evaluation/report_v2.go` - render report status, confidence,
      boolean, limits, unknowns, findings, and report refs through display
      helpers.
- [x] `internal/evaluation/evaluation_test.go` - assert titled Markdown values
      and stable machine values.

### Format spec

- [x] None - `SPECIFICATION.md` does not govern Evaluation v2 report presentation
      titles.

### Durable specs

- [x] `specs/evaluation-v2/reports/report-tree.md` - require Markdown reports to
      use display titles for CLI-owned enum-like values while keeping raw values
      in structured data.

### Durable docs / bundled skill

- [x] None - this is CLI report rendering behavior, not user-facing authoring or
      `/quality` skill prompt behavior.

### Suggested new durable specs

- None.

## Status

`Done`. Spec and design settled, Evaluation v2 report display helpers
implemented, durable report-tree spec updated, tests and formatting verified, and
archived. See the [status lifecycle](../index.md#status-lifecycle).
