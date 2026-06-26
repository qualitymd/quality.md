---
type: Change Case
title: Evaluation v2 report header navigation
description: Simplify Evaluation v2 report headers with labeled Area and Factor trails, compact status summaries, and no redundant parent links.
status: Done
tags: [cli, evaluation, reports, navigation]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report header navigation

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0104-evaluation-v2-report-header-navigation/spec.md) - what
  the change must do.
- [Design doc](0104-evaluation-v2-report-header-navigation/design.md) - how it
  should be built, and why.

## Motivation

Evaluation v2 reports currently begin with a generic `Breadcrumb:` line and then
often repeat the same navigation with `Parent Area:` or `Parent Factor:` links.
That shape is technically navigable, but it exposes implementation-ish UI
language and makes the top of each report busier than the reader needs.

The report header should answer three questions quickly: where am I in the Area
tree, what is this report about, and how is it doing? Labeled navigation trails
and compact one-row summaries can make that answer clearer without changing the
report tree, routine JSON, ratings, or analysis semantics.

## Scope

Covered:

- Replace the generic `Breadcrumb:` label with an `Area:` navigation trail in
  every Evaluation v2 Markdown report.
- Add a `Factor:` navigation trail to Factor reports only.
- Remove standalone `Parent Area:` and `Parent Factor:` header links because the
  labeled trails provide the same upward navigation with more context.
- Reshape the top identity/status block for Area, Factor, and Requirement
  reports into a compact report-specific summary while preserving required path,
  name, factor-link, confidence, and data-link information.
- Update the durable Evaluation v2 report-tree spec, implementation, and focused
  report tests.

Deferred / non-goals:

- No changes to report paths, report refs, routine JSON, generated
  `EvaluationOutputResult`, or report-build receipt JSON.
- No changes to rating, confidence, status, or display-title semantics.
- No new generated HTML/CSS surface; reports remain Markdown.
- No Factor breadcrumb on Requirement reports. Requirement reports keep an
  owning Area trail and linked attached Factors in their summary.

## Affected artifacts

### Code

- [x] `internal/evaluation/report_v2.go` - render labeled Area and Factor trails,
      remove parent-link header lines, and emit compact summary headers.
- [x] `internal/evaluation/evaluation_test.go` - cover the new Markdown header
      shape and the absence of old `Breadcrumb:` / `Parent` lines.

### Format spec

- [x] None - `SPECIFICATION.md` defines Evaluation Report semantics, not the
      concrete Evaluation v2 Markdown header layout.

### Durable specs

- [x] `specs/evaluation-v2/reports/report-tree.md` - replace the breadcrumb and
      parent-link navigation requirements with labeled Area/Factor trail and
      compact header requirements.

### Durable docs / bundled skill

- [x] `evaluation-v2-sketch.md` - reconcile the older starter Markdown examples
      that still show `Breadcrumb:` and `Parent:` headers, or explicitly mark the
      sketch as superseded by the durable report-tree spec.
- [x] None for the bundled `/quality` skill - this is generated report rendering,
      not agent workflow guidance.

### Suggested new durable specs

- None.

## Status

`Done`. Implemented the renderer changes, updated focused report tests, synced
the durable report-tree spec and reference sketch, and archived. See the
[status lifecycle](../index.md#status-lifecycle).
