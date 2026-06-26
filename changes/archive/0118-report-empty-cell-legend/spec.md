---
type: Functional Specification
title: Report empty-cell marker and legend
description: Render every empty scalar cell in a generated Evaluation report as a single em dash, and add one static per-report legend defining the marker.
status: Draft
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Report empty-cell marker and legend

This spec governs how generated Evaluation reports render an **absent scalar
cell value** and how each report explains that rendering. It is the delta for
[0118](../0118-report-empty-cell-legend.md). The durable contract it lands in is
[`specs/evaluation/reports/report-tree.md`](../../specs/evaluation/reports/report-tree.md)
(normative — its Rendering Rules and the existing em-dash and status-label rules
bind here).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

A generated report renders an absent scalar value as a blank table cell, which
is ambiguous — a reader cannot distinguish "no value recorded" from "not
applicable" from a rendering fault. The report builder also carries several
uncoordinated empty-value treatments (blank cells, a bare `/` for an all-empty
confidence pair, an em dash for the descendant roll-up column, parenthetical
`(no …)` rows for empty sections), so "nothing here" reads differently in
different spots.

This change settles on one visible marker for an empty *cell* — the em dash
(`—`), already in use for the roll-up column — and one static legend per report
that defines it. The em dash is preferred over a literal `N/A` because the
absence is usually "not recorded," not "not applicable," and a neutral dash
states absence without overclaiming. A static legend is preferred over per-cell
asterisks because reports are committed and diffed: a "first instance" asterisk
has no stable anchor in a Markdown table and migrates as data changes, while a
fixed legend stays put. The not-assessed family of *outcomes* is out of scope —
it is an evaluation result, not an empty cell, and keeps its distinct labels.

## Scope

Covered: the em-dash marker for empty scalar cells (including composite pair
cells), and one static legend per generated report. Deferred / non-goals are
recorded in the [parent concept](../0118-report-empty-cell-legend.md#scope);
in short — no change to status or Rating Level labels, to parenthetical
empty-section rows, to structured data or raw values, or to paths, filenames,
trails, and links; and no literal `N/A` text.

## Requirements

### Empty scalar cells

When a report table cell would otherwise render an empty scalar value — because
its source field is absent, null, or empty, or because the cell does not apply
to its row — the report **MUST** render an em dash (`—`) in place of a blank
cell.

>> Rationale: a blank cell is indistinguishable from a rendering fault; one
>> visible marker makes the absence explicit and unmistakable. — 0118

When a cell pairs two scoped values — the Confidence and analysis-status cells
render an overall/local pair — and either component value is absent, the report
**MUST** render `—` for that component (for example, `— / —`), not an empty
segment.

>> Rationale: the pair string is assembled before the empty-cell guard sees it,
>> so without this an all-empty confidence cell renders a bare `/` that reads
>> as a formatting glitch rather than "no confidence recorded." — 0118

### Distinctness from outcomes and sections

The em-dash empty-cell marker **MUST** remain distinct from Rating Level labels
and from the `not_assessed`, `not_rated`, `empty`, `not_analyzed`, and `blocked`
status labels. A not-assessed, not-rated, not-analyzed, or blocked outcome
**MUST** continue to render with its existing distinct status label, never the
em dash.

>> Rationale: "not assessed" is an evaluation *outcome* the format requires be
>> shown distinctly from a Rating Level; collapsing it into a generic empty
>> marker would erase a meaningful result. The em dash marks an absent cell
>> *value*, not an absent rating. — 0118

Empty whole-section placeholder rows — such as `(no findings)`,
`(no rating drivers)`, and `(none recorded)` — **MUST** remain rendered as their
parenthetical empty-state rows and **MUST NOT** be replaced by the em-dash cell
marker.

>> Rationale: a section that produced nothing and a single missing value within
>> a populated row are different absences; the parenthetical names the empty
>> section in words, while the em dash marks one missing cell. — 0118

### Per-report legend

Every generated report — the root Area report, and every Area, Factor, and
Requirement report — **MUST** include exactly one legend, rendered once at the
foot of the report, that defines the em-dash marker as standing for a value that
is not applicable or not recorded.

A report's legend **MUST** render regardless of whether that report contains any
em-dash cells, so the legend's presence does not vary with evaluation data.

>> Rationale: the legend was chosen over per-cell asterisks specifically for
>> diff stability across re-runs; conditioning its presence on the data would
>> reintroduce the appear/disappear churn it exists to avoid. — 0118

The legend **MUST** define the em-dash marker only; the status labels and the
parenthetical empty-section rows are self-describing and **MUST NOT** be added to
the legend.

>> Rationale: a legend that re-explains self-evident worded labels is noise;
>> scope it to the one marker a reader cannot infer. — 0118

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` — in Rendering Rules, add the
  empty-scalar-cell em-dash rule (covering composite pair cells) and the
  per-report static-legend requirement, and reconcile both with the existing
  roll-up em-dash rule (0097/0111) and the `not_assessed`/`not_rated`/`empty`/
  `not_analyzed`/`blocked` distinctness rule (per the requirements above).

### To rename

None

### To delete

None
