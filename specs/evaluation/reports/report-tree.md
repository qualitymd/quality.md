---
type: Functional Specification
title: Evaluation report tree
description: Deterministic Area, Factor, and Requirement Markdown reports for Evaluation.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation report tree

Evaluation reports are deterministic Markdown projections over completed
structured routine outputs.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Source Of Truth

Report generation **MUST** consume `EvaluationOutputResult` and referenced
structured routine outputs.

Report generation **MUST NOT** inspect new source evidence.

Report generation **MUST NOT** introduce new findings, ratings, evidence, limits,
analysis, or recommendations.

## Report Paths

The root Area report **MUST** be generated as `report.md` at the run root.

Non-root Area, Factor, and Requirement reports **MUST** use short
subject-aware filenames derived from structural model IDs:

- `areas/<area>/<area>-area.md`
- `requirements/<requirement>/<requirement>-requirement.md`
- `factors/<factor>/<factor>-factor.md`
- `factors/<factor>/factors/<sub-factor>/<sub-factor>-factor.md`

Report filenames **MUST NOT** be derived from display titles, natural labels, or
rendered human labels.

The report builder **MUST NOT** write duplicate compatibility copies using old
descendant `report.md` filenames.

> Rationale: the root `report.md` remains the single run entrypoint, while
> descendant filenames need enough local subject context for editor and browser
> tabs. Structural IDs keep paths stable; the existing directory tree carries
> full identity. — 0108

## Navigation

Every report **MUST** render its H1 title line as the first content of the
report. The H1 **MUST** prefix the subject display title with the report kind:
`Area:` for root and non-root Area reports, `Factor:` for Factor reports, and
`Requirement:` for Requirement reports.

The `Area:` navigation trail **MUST** render after the H1. Its elements **MUST**
link from the root Area report through the current Area report or owning Area
report. The root trail element **MUST** render the Model `title` when present.

Factor reports **MUST** include a `Factor:` navigation trail after the `Area:`
trail. The `Factor:` trail **MUST** link each Factor ancestor and the current
Factor to its generated Factor report.

Requirement reports **MUST** include a plural `Factors:` context line after the
`Area:` trail. The line **MUST** list every attached Factor as a link to its
generated Factor report, joined with `;` as a flat set. When no Factors are
attached, it **MUST** render an explicit empty-state marker.

Requirement reports **MUST NOT** render a singular `Factor:` breadcrumb, use the
`/` nesting separator for the `Factors:` line, or choose one attached Factor as
a navigation parent.

Reports **MUST NOT** render standalone `Breadcrumb:`, `Parent Area:`,
`Parent Factor:`, `Parent:`, `Path:`, or `Name:` header lines.

Area reports **MUST** link to local root Factor reports, local Requirement
reports, and direct child Area reports.

Factor reports **MUST** link to their owning Area report, parent Factor report
when present, child Factor reports, and direct Requirement reports.

Requirement reports **MUST** link to their owning Area report and every attached
Factor report.

Report tables **MUST** render the row subject as the generated human report link
when that row has exactly one generated human report target. Reports **MUST**
keep explicit `Data` links for machine-readable payloads instead of moving those
links onto subject labels. Each `Data` link **MUST** use the linked payload's
base filename as its link text (for example `area-analysis-result.json`), not a
generic word.

> Rationale: labeled trails expose the Model hierarchy directly, and subject-cell
> links make report navigation land on the named thing readers naturally open.
> Machine data links remain explicit because they target structured payloads, not
> generated human report pages. The payload filename is the one detail a generic
> label omits, and it matches the file a reader opens. — 0104, 0105, 0109

## Area Reports

Area reports **MUST** include:

- kind-prefixed Area title;
- Area navigation trail;
- overall and local ratings;
- overall and local confidence;
- data links;
- summary;
- rating drivers;
- local root Factors;
- direct child Areas;
- local Requirements; and
- limits and incomplete inputs.

## Factor Reports

Factor reports **MUST** include:

- owning Area link;
- Factor navigation trail;
- kind-prefixed Factor title;
- overall and local ratings, where `Overall Rating` is the Factor
  `localAndDescendantAnalysis` rating and `Local Rating` is its `localAnalysis`
  rating;
- local and local-and-descendant statuses;
- confidence;
- data links;
- summary;
- rating drivers;
- direct Requirements;
- direct child Factors; and
- limits and incomplete inputs.

## Requirement Reports

Requirement reports **MUST** include:

- owning Area link;
- attached Factor links in a plural `Factors:` context line;
- kind-prefixed Requirement title;
- Requirement rating status and selected rating when present;
- assessment status;
- confidence;
- data links;
- summary;
- findings summary;
- finding detail sections; and
- unknowns and missing evidence.

## Rendering Rules

Reports **MUST** render empty tables with explicit empty-state rows.

Report headers **SHOULD** use report-specific summary tables instead of a
generic `Field | Value` key-value table. Area headers should summarize
`Overall Rating`, `Local Rating`, `Confidence`, and `Data`; Factor headers should
summarize `Overall Rating`, `Local Rating`, `Status`, `Confidence`, and `Data`;
Requirement headers should summarize `Rating`, `Assessment`, `Confidence`, and
`Data`; attached Factors belong in the plural `Factors:` context line, not in
the summary table.

> Rationale: the title identifies the report subject, so the header table should
> prioritize state and navigation rather than repeat the subject kind as
> metadata. The subject kind now rides the H1 title; location rides the
> navigation trail, so separate `Path:` / `Name:` header lines would be
> redundant. — 0104, 0119

When a report table cell would otherwise render an empty scalar value, including
one component of a paired Confidence or Status cell, the cell **MUST** render an
em dash (`—`) instead of a blank segment. Empty whole-section placeholder rows
such as `(no findings)`, `(no rating drivers)`, and `(none recorded)` **MUST**
remain worded empty-state rows rather than being replaced by the cell marker.

Each generated report **MUST** include exactly one static legend at the foot of
the report defining `—` as "not applicable or not recorded". The legend **MUST**
render regardless of whether the report contains an em-dash cell.

> Rationale: blank cells are ambiguous in committed Markdown reports. A neutral
> em dash makes absence visible without overclaiming `N/A`, and a static legend
> avoids data-dependent footnote churn across re-runs. — 0118

Every rating column **MUST** name what it rates. A header summary table **MUST**
label its descendant-inclusive rating column `Overall Rating` and its local
rating column `Local Rating`.

> Rationale: the adjacent header columns are self-describing nouns, so bare
> `Overall` / `Local` made a reader supply the missing noun. — 0111

The Area report Factors table, the Area report Sub-Areas table, and the Factor
report child Factors table each list a subject's children, one row per child.
Each **MUST** render a `Local Rating` column from the child's `localAnalysis`
rating, and a descendant-inclusive sub-rating column — `+ Sub-Factors Rating` for
a Factor row, `+ Sub-Areas Rating` for an Area row — from the child's
`localAndDescendantAnalysis` rating. These tables **MUST NOT** render a
boolean in a rating column. When a row's subject has no descendant Factors (for a
Factor row) or no descendant Areas (for an Area row), its `+ Sub-Factors Rating`
/ `+ Sub-Areas Rating` cell **MUST** render an em dash (`—`) rather than
repeating the local rating.

> Rationale: these breakdown tables previously rendered the aggregate rating in
> the local `Rating` column and a `Yes`/`No` boolean where the roll-up rating
> belonged, leaving the local rating unshown — the unmet distinction clean-break
> case 0097 required. The em dash preserves the old boolean's "has children"
> signal without presenting a redundant rating. — 0097, 0111

Reports **MUST** render selected Rating Levels with the Rating Level `title`
resolved from the run's `model-snapshot.md` snapshot, falling back to the stable Rating
Level ID only when a title is unavailable.

> Rationale: Markdown reports are the human review surface, and the model
> snapshot is the historical source for display vocabulary. Structured routine
> data and machine receipts keep stable Rating Level IDs. — 0102

Reports **MUST** render `not_assessed`, `not_rated`, `empty`, `not_analyzed`,
and `blocked` distinctly from Rating Level labels.

Reports **MUST** render CLI-owned enum-like report values, including statuses,
confidence levels, boolean values, report kinds, limits/incomplete-input types,
unknown/missing-evidence types, and known finding classifications, with
human-readable display titles in Markdown while preserving the raw values in
routine JSON, `EvaluationOutputResult`, and report-build receipts.

> Rationale: Markdown reports are optimized for human review and scanning, but
> agents and tools need stable values in the structured data. Unknown or
> free-form values should remain readable through fallback title-casing rather
> than turning presentation decoration into schema validation. — 0103

Reports **MUST** omit Rating Level values when the source result status says the
rating or scoped analysis was not produced.

Reports **MUST** preserve secret-handling boundaries. They may name the locator
and credential type but **MUST NOT** reproduce secret values or unsafe raw
content.

Ordering **MUST** be deterministic:

- Areas by canonical Area identity, with the root Area first;
- Factors by canonical Factor identity;
- Requirements by canonical Requirement identity;
- rating drivers in source result order;
- findings in Requirement Assessment Result order; and
- evidence in recorded order.
