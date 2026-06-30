---
type: Change Case
title: Requirement report Factors line
description: Add a plural `Factors:` context line to the header of each Evaluation v2 Requirement report — a flat set of links to every attached Factor, parallel to the `Area:` trail — and drop the now-redundant `Factors` column from the Requirement header summary table.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Requirement report Factors line

This parent concept captures the why and status; the detail lives in its child:

- [Functional spec](0117-requirement-report-factors-line/spec.md) - what the
  change must do.

No design doc: the change is a small, self-contained report-header presentation
edit whose one non-obvious decision (reconciling the existing `Factor:`
breadcrumb prohibition) is settled in the spec's rationale annotations.

## Motivation

A Requirement sits in the Model along two independent axes: the Area that
declares it, and the set of Factors it serves. The Evaluation v2 Requirement
report already surfaces the first axis as a labeled `Area:` trail at the top of
the header — but the second axis is buried in a `Factors` column inside the
header summary table, where it reads as just another rating-adjacent cell rather
than as the navigational context it is.

The Factor report already establishes the pattern: it stacks an `Area:` trail and
a `Factor:` trail as header context lines
([`report-tree.md`](../specs/evaluation-v2/reports/report-tree.md) lines 56–58).
Giving the Requirement report a parallel `Factors:` line makes both placement
axes visible at a glance and matches the rest of the report tree.

One subtlety the report tree deliberately encodes: a Requirement can attach to
_several_ unrelated Factors, so the report **MUST NOT** render a singular
`Factor:` breadcrumb or pick one attached Factor as a navigation parent
([`report-tree.md`](../specs/evaluation-v2/reports/report-tree.md) lines 60–61).
The new line respects that exactly — it is a _plural_ `Factors:` set with no
nesting and no chosen parent. It uses a set separator (`;`, matching the
existing `requirementFactorLinks` join), not the `/` nesting separator the
`Area:` and `Factor:` trails use, so the flat-set semantics stay distinct from a
hierarchical breadcrumb. Adding the line makes the table's `Factors` column
redundant, so this case removes that column.

## Scope

Covered:

- the Evaluation v2 Requirement report header: a new `Factors:` context line
  listing links to every attached Factor as a flat, `;`-separated set;
- removal of the `Factors` column from the Requirement header summary table; and
- the durable report-tree contract and sketch describing the Requirement report
  header and its navigation, including the `Factor:`-breadcrumb prohibition that
  must be reconciled with the new plural line.

Deferred / non-goals:

- no change to which Factors a Requirement attaches to, or how they are resolved;
- no change to the `Area:` trail, the Factor report `Factor:` trail, or any Area
  or Factor report;
- no change to ratings, assessments, findings, confidence, summaries, limits, or
  any other report content;
- no change to structured routine data, `EvaluationOutputResult`, JSON field
  names, report paths, filenames, or links; and
- no migration of existing completed evaluation runs.

## Affected artifacts

### Code

- [ ] `internal/evaluation/report_v2.go` - in `renderV2RequirementReport`, write a
      `Factors:` context line after the `Area:` trail (reusing the attached-Factor
      links already produced by `requirementFactorLinks`, joined with `;`), and
      drop the `Factors` column from the Requirement header summary table so the
      links are not rendered twice.
- [ ] `internal/evaluation/evaluation_test.go` - update the Requirement report
      assertions: expect the new `Factors:` line, expect the header table without
      the `Factors` column, and keep the assertion that no singular `Factor:`
      line is rendered.

### Format spec

- [ ] None - `SPECIFICATION.md` does not govern Evaluation v2 generated report
      presentation. (Deliberate.)

### Durable specs

- [ ] `specs/evaluation-v2/reports/report-tree.md` - reconcile the Navigation
      rule that forbids a `Factor:` breadcrumb so it permits a plural `Factors:`
      set-line; add the `Factors:` line to the Requirement Reports MUST-include
      list; and drop `Factors` from the Requirement header summary-table columns
      in the Rendering Rules.

### Durable docs / bundled skill

- [ ] `evaluation-v2-sketch.md` - note the Requirement report `Factors:` context
      line in the navigation rules / report shape.
- [ ] `specs/skills/quality-skill/reporting.md` - reviewed; the
      "Requirement reports MUST link to every attached Factor report" rule still
      holds (the links move from the table to the header line). No change expected.

### Suggested new durable specs

- None.

## Status

`Done`. Implemented, verified, and archived. Requirement reports now show the
plural `Factors:` context line and no longer duplicate attached Factors in the
summary table.
