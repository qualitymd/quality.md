---
type: Functional Specification
title: Evaluation v2 report tree
description: Deterministic Area, Factor, and Requirement Markdown reports for Evaluation v2.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation v2 report tree

Evaluation v2 reports are deterministic Markdown projections over completed
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

Every report **MUST** start with an `Area:` navigation trail whose elements link
from the root Area report through the current Area report or owning Area report.

Factor reports **MUST** include a `Factor:` navigation trail after the `Area:`
trail. The `Factor:` trail **MUST** link each Factor ancestor and the current
Factor to its generated Factor report.

Requirement reports **MUST NOT** render a `Factor:` breadcrumb or choose one
attached Factor as a navigation parent.

Reports **MUST NOT** render standalone `Breadcrumb:`, `Parent Area:`,
`Parent Factor:`, or `Parent:` header lines.

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

- Area title and path;
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
- Factor title and path;
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
- Requirement title and name;
- Requirement rating status and selected rating when present;
- assessment status;
- attached Factor links;
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
Requirement headers should summarize `Rating`, `Assessment`, `Factors`,
`Confidence`, and `Data`.

> Rationale: the title and path/name line identify the report subject, so the
> header table should prioritize state and navigation rather than repeat the
> subject kind as metadata. — 0104

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
resolved from the run's `model.md` snapshot, falling back to the stable Rating
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

- Areas by structural Area ID, with the root Area first;
- Factors by declaring Area ID and structural Factor path;
- Requirements by declaring Area ID and Requirement name;
- rating drivers in source result order;
- findings in Requirement Assessment Result order; and
- evidence in recorded order.
