---
type: Functional Specification
title: Requirement report Factors line - functional spec
description: Requirements for adding a plural `Factors:` context line to the Evaluation v2 Requirement report header and removing the now-redundant `Factors` column from its summary table.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Requirement report Factors line - functional spec

Companion to the
[Requirement report Factors line](../0117-requirement-report-factors-line.md)
change case. This spec states what the Requirement report header change must do.
The durable contract it lands in is the
[Evaluation v2 report tree](../../../specs/evaluation-v2/reports/report-tree.md)
spec.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

A Requirement is placed in the Model along two independent axes: its declaring
Area and the set of Factors it serves. The Requirement report already surfaces
the Area axis as a labeled `Area:` trail at the top of the header, mirroring the
Factor report, which stacks an `Area:` trail and a `Factor:` trail. The Factor
axis of a Requirement, however, appears only as a `Factors` column inside the
header summary table, where it reads as a rating-adjacent cell rather than as the
placement context it is.

This change adds a parallel `Factors:` context line to the Requirement report
header so both placement axes are visible at a glance. Because a Requirement can
attach to several unrelated Factors, the report tree forbids a singular `Factor:`
breadcrumb and forbids choosing one attached Factor as a navigation parent; the
new line honors that by rendering a *flat set* of all attached Factors with no
nesting and no chosen parent, separated by `;` (the set separator already used
for Factor lists) rather than the `/` nesting separator of the `Area:` and
`Factor:` trails. With the Factors now in the header line, the summary-table
`Factors` column is redundant and is removed.

## Scope

Covered: the generated Evaluation v2 Requirement report header `Factors:` context
line; removal of the `Factors` column from the Requirement header summary table;
and the durable contract and sketch describing the Requirement report header and
navigation.

Not covered: Factor attachment/resolution; the `Area:` trail and the Factor
report `Factor:` trail; ratings, assessments, findings, confidence, summaries,
limits, or other report content; structured routine data, `EvaluationOutputResult`,
JSON field names, report paths, filenames, links; and migration of existing runs.

## Requirements

### The Factors context line

- The Requirement report **MUST** render a `Factors:` context line in its header,
  after the `Area:` trail and before the Requirement title.

- The `Factors:` line **MUST** list every Factor attached to the Requirement,
  each rendered as a link to that Factor's generated report.

- The `Factors:` line **MUST** join its entries with `;` and **MUST NOT** use
  the `/` separator reserved for nesting trails.

  > Rationale: a Requirement's Factors are a flat set, not a hierarchy; `/`
  > would falsely imply nesting, and `;` matches the existing Factor-list join.
  > — 0117

- When a Requirement has no attached Factors, the `Factors:` line **MUST** render
  an explicit empty-state marker (for example `Factors: (none)`) rather than a
  bare `Factors:` label with no value.

- The Requirement report **MUST NOT** render a singular `Factor:` breadcrumb or
  choose one attached Factor as a navigation parent.

  > Rationale: a Requirement can serve several unrelated Factors, so no single
  > Factor is its parent; the plural set-line is the correct rendering. This
  > preserves the existing report-tree prohibition. — 0104, 0117

### The header summary table

- The Requirement header summary table **MUST NOT** include a `Factors` column;
  the attached Factors are rendered by the `Factors:` context line instead.

  > Rationale: rendering the same Factor links in both the header line and a table
  > column would duplicate them. — 0117

- The Requirement header summary table **MUST** otherwise be unchanged, summarizing
  `Rating`, `Assessment`, `Confidence`, and `Data`.

### Unaffected presentation

- Report content, ratings, assessments, findings, confidence, summaries, links,
  navigation targets, paths, structured routine data, and `EvaluationOutputResult`
  **MUST NOT** change as part of this header change.

- The `Area:` trail, the Factor report `Factor:` trail, and all Area and Factor
  reports **MUST NOT** change.

## Acceptance Criteria

- A Requirement report renders, in order: an `Area:` trail line, a `Factors:`
  line, then the `#` title.
- The `Factors:` line lists each attached Factor as a report link, joined with
  `;` (for example `Factors: [reliability](../../factors/reliability/reliability-factor.md)`).
- A Requirement with multiple attached Factors shows all of them on the one line,
  `;`-separated; a Requirement with none shows an explicit empty-state marker.
- The Requirement header summary table reads `| Rating | Assessment | Confidence |
  Data |` with no `Factors` column.
- No Requirement report renders a line beginning `Factor:` (singular).
- Report paths, links, navigation labels, and structured `data/` outputs are
  unchanged.

## Durable spec changes

### To add

None.

### To modify

- [`specs/evaluation-v2/reports/report-tree.md`](../../../specs/evaluation-v2/reports/report-tree.md)
  - in Navigation, reconcile the rule that Requirement reports **MUST NOT** render
    a `Factor:` breadcrumb so it permits a plural `Factors:` set-line (flat, no
    chosen parent) while still forbidding a singular `Factor:` trail; in
    Requirement Reports, add the `Factors:` context line to the MUST-include list;
    and in Rendering Rules, drop `Factors` from the Requirement header
    summary-table columns. Driven by
    [The Factors context line](#the-factors-context-line) and
    [The header summary table](#the-header-summary-table).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, each Requirement report shows its two
placement axes as parallel header context lines — the `Area:` trail and a flat
`Factors:` set — with the attached Factors rendered exactly once and the singular
`Factor:` breadcrumb still prohibited. That achieves the case's motivation
(surface the Factor axis as navigational context, parallel to the Area axis)
without duplicating links, touching report content or structured data, or
implying a nesting relationship the Factor set does not have.
