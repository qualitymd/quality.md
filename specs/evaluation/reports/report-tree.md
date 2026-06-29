---
type: Functional Specification
title: Evaluation report tree
description: Deterministic run, Area, Factor, and Requirement Markdown reports for Evaluation.
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
analysis, or recommendations. It **MUST** render persisted Advice outputs when
the run is otherwise renderable.

Generated report frontmatter and Markdown body content **MUST NOT** be read as
report-generation input. Generated reports are output conveniences for readers,
agents, and editor previews; structured Evaluation data remains the source of
truth.

## Report Paths

The run-level Evaluation report **MUST** be generated as `report.md` at the run
root.

The full ranked findings index **MUST** be generated as `findings.md` at the run
root.

The root Area report **MUST** be generated as `root-area.md` at the run root
when the root Area has an Area Analysis Result in the run.

Non-root Area, Factor, and Requirement reports **MUST** use short
subject-aware filenames derived from structural model IDs:

- `areas/<area>/<area>-area.md`
- `requirements/<requirement>/<requirement>-requirement.md`
- `factors/<factor>/<factor>-factor.md`
- `factors/<factor>/factors/<sub-factor>/<sub-factor>-factor.md`

Report filenames **MUST NOT** be derived from display titles, natural labels, or
rendered human labels.

The report builder **MUST NOT** write duplicate compatibility copies using old
root Area or descendant `report.md` filenames.

The report builder **MUST NOT** rename report content files such as
`findings.md` or `recommendations.md` to `index.md` for OKF compatibility.
Generated `index.md`, `schema.md`, and `log.md` files for Evaluation run folders
are deferred.

The recommendation index **MUST** be generated as `recommendations.md` at the
run root. Recommendation detail reports **MUST** be generated under
`recommendations/` with rank-prefixed filenames:

- `recommendations/<NNN>-<slug>.md`

The `<NNN>` prefix **MUST** follow `RecommendationRankingResult` ordering. The
slug **SHOULD** derive from the recommendation title and fall back to the
recommendation ID when needed.

> Rationale: `report.md` is the run entrypoint. The root Area detail report uses
> `root-area.md` so its filename names the subject it contains, while descendant
> filenames keep enough local subject context for editor and browser tabs.
> Structural IDs keep paths stable; the existing directory tree carries full
> identity. — 0108, 0137

## Report Frontmatter

Every generated Markdown report **MUST** begin with YAML frontmatter containing
only `type` and `title`.

The `type` field **MUST** use this report-subject taxonomy:

| Report output                  | `type`                          |
| ------------------------------ | ------------------------------- |
| `report.md`                    | `Evaluation Overview Report`    |
| root and non-root Area reports | `Area Evaluation Report`        |
| Factor reports                 | `Factor Evaluation Report`      |
| Requirement reports            | `Requirement Evaluation Report` |
| `findings.md`                  | `Finding Index Report`          |
| `recommendations.md`           | `Recommendation Index Report`   |
| recommendation detail reports  | `Recommendation Report`         |

Report frontmatter **MUST NOT** use Model concept names such as `Area`,
`Factor`, or `Requirement` as the report `type`.

The `title` field **MUST** equal the plain-text content of the report's H1 title
line, with the leading Markdown `#` marker removed. For example, a Requirement
report frontmatter title uses `Requirement: <title>`, matching the visible H1
`# Requirement: <title>`.

Report frontmatter **MUST NOT** duplicate generated time, run identity, model
snapshot, subject identity, scope, ratings, confidence, summaries, rating
drivers, findings, recommendations, limits, evidence, or rendered display
labels when those values are available from structured Evaluation payloads or
the visible Markdown body.

> Rationale: `type` records the report artifact kind, while `title` records the
> generated Markdown document title a reader sees first. Keeping frontmatter
> title equal to the H1 avoids a second subject-only title system without
> turning the first screen of every report into a source-data manifest.
> Report-local source-data pointers live in the visible bottom section instead.
> — 0158, 0162, 0167

## Source Data Section

Every generated Markdown report **MUST** end with a `## Source Data` section.

The `Source Data` section **MUST** list the run-root-relative structured
Evaluation payload paths used as source data for that specific Markdown report
artifact.

Each source-data list item **MUST** render as a Markdown link whose label is the
run-root-relative payload path and whose target is relative to the report file
that contains the section.

Reports that render run number, run ID, creation time, or requested scope from
the run manifest **MUST** include `data/run-manifest.json`.

The `Source Data` section **MUST NOT** include
`data/evaluation-output-result.json` solely because that generated output index
exists. A report **MAY** list `data/evaluation-output-result.json` only if that
report is directly rendered from it.

> Rationale: The bottom section keeps source payloads visible and parseable for
> people and agents without making frontmatter noisy. Path labels make the
> source files discoverable in plain text, while report-relative targets keep
> nested reports navigable. `EvaluationOutputResult` indexes generated outputs
> after report build; it is not source data for those reports unless a renderer
> explicitly consumes it. — 0159, 0162

## Run Report

The run-level `report.md` **MUST** render as the scoped Area report described by
`RunManifest.plannedScope`. It **MUST** include:

- scoped Area title and rating;
- top 10 ranked findings;
- top 10 ranked recommendations;
- link to the full findings index;
- link to the full recommendation index;
- summary from the scoped Area result;
- Area / Factor Breakdown for the scoped Area;
- requested and planned Evaluation scope;
- root Area coverage status; and
- limits and incomplete inputs from the scoped Area result.

The run-level `report.md` **MUST NOT** render a standalone `Rating Drivers`
section or `Driver | Effect | Inputs` table. Rating drivers remain structured
source data available through the payloads listed in the report's `Source Data`
section.

When `plannedScope.factorFilter` is non-empty, `report.md` **MUST** identify the
filtered Factors and **MUST** avoid presenting the result as a complete Area
roll-up unless the structured analysis records an appropriate limit.

The run report **MUST** state when the root Area was not evaluated in the run.

The run report **MUST NOT** introduce report-only findings, ratings, evidence,
limits, analysis, recommendations, candidate actions, or source claims.

The Top Findings table **MUST** render rows from
`FindingRankingResult.orderedFindings` ordered by rank and capped at 10 rows. It
**MUST** render the columns `Rank`, `Finding`, `Area`, `Factors`, `Type`, and
`Severity`, in that order. The `Finding` cell **MUST** use the finding
`statement` as link text and link to the exact finding detail section in the
Requirement report. The table **MUST NOT** render a finding artifact-ID column.
The `Area` cell **MUST** link the declaring Area title to the Area report. The
`Factors` cell **MUST** render comma-separated attached Factor title links, or
`—` when no Factor link can be rendered. `Type` and `Severity` **MUST** render
existing display labels, including their emoji, for known finding type and
severity enum values.

The top finding and recommendation sections **MUST** be omitted only when the
persisted Advice payloads contain no rows to render. `report.md` **MUST** always
link to `findings.md` and `recommendations.md` when the report tree is built.

The Top Recommendations table **MUST** render rows from
`RecommendationRankingResult.orderedRecommendations` ordered by rank and capped
at 10 rows. It **MUST** render the columns `Rank`, `#`, `Recommendation`,
`Area / Factors`, and `Reason`, in that order. The `#` cell **MUST** render the
assigned `RecommendationResult.number`. The `Recommendation` cell **MUST** use
`RecommendationResult.title` as link text and link to the generated
recommendation detail report. The `Area / Factors` cell **MUST** render linked
Area and Factor names resolved from `RecommendationResult.traceRefs` through
persisted evaluation data and the model snapshot, or `—` when no Area or Factor
can be resolved. The `Reason` cell **MUST** render
`RecommendationResult.expectedValue`.

## Finding Reports

`findings.md` **MUST** render a complete ranked findings index from
`FindingRankingResult`. It **MUST** include:

- all ranked findings ordered by rank;
- the same columns and link behavior as the run report Top Findings table.

## Recommendation Reports

`recommendations.md` **MUST** render a complete recommendation index from
persisted `RecommendationResult` payloads and `RecommendationRankingResult`.
It **MUST** include:

- all ranked recommendations;
- Area / Factors links resolved from `RecommendationResult.traceRefs`;
- impact;
- confidence;
- Reason from `RecommendationResult.expectedValue`;
- ranking rationale;
- links to recommendation detail reports; and
- a coverage summary from `findingCoverage`.

Each recommendation detail report **MUST** include:

- recommendation title;
- assigned recommendation ID;
- typed recommendation artifact reference;
- rank when ranked;
- impact;
- confidence;
- description;
- background;
- expected value;
- done criterion;
- trace refs.

Recommendation Markdown reports **MUST** remain human-first and **MUST NOT**
require YAML frontmatter for machine readability.
Generated recommendation reports **MUST** render only persisted
`RecommendationResult`, `RecommendationRankingResult`, model snapshot, and
referenced evaluation data. They **MUST NOT** read YAML frontmatter or Markdown
body content from another generated report as recommendation source data.
Recommendation ranking rationale **MUST** remain sourced from
`RecommendationRankingResult.orderedRecommendations[].rationale` and **MUST NOT**
be conflated with `RecommendationResult.background` or `expectedValue`.

## Navigation

Every report **MUST** render its H1 title line as the first Markdown content
after frontmatter. The H1 **MUST** prefix the subject display title with the
report kind: `Area:` for root and non-root Area reports, `Factor:` for Factor
reports, `Requirement:` for Requirement reports, and `Recommendation:` for
recommendation detail reports. The run-level H1 **MUST** identify the report as
an Evaluation Report. The H1 title line and frontmatter `title` **MUST** use the
same plain-text title.

Every report **MUST** render a run context line and a report navigation line near
the H1. The report navigation line **MUST** link to the run overview
`report.md`, full findings index `findings.md`, and full recommendation index
`recommendations.md` when the current report is not that target.

The `Area:` navigation trail **MUST** render after the H1. Its elements **MUST**
link to generated Area reports from the root Area through the current Area
report or owning Area report. When an ancestor Area report was not generated
because the run was scoped below it, the trail **MUST** render that ancestor as
plain text. The root trail element **MUST** render the Model `title` when
present.

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
when that row has exactly one generated human report target. Generated Markdown
report bodies **MUST NOT** duplicate report-level source-data links in `Data`
columns or equivalent header source-data lines; the `Source Data` section owns
the visible source-data manifest.

> Rationale: labeled trails expose the Model hierarchy directly, and subject-cell
> links make report navigation land on the named thing readers naturally open.
> Machine data links target structured payloads, not generated human report
> pages, so they live in a dedicated bottom section rather than summary tables.
> Keeping those links out of visible summary tables makes the report header
> easier to scan without hiding the source-data manifest from agents or
> secondary tooling. — 0104, 0105, 0109, 0159, 0162

## Area Reports

Area reports **MUST** include:

- kind-prefixed Area title;
- Area navigation trail;
- overall and local ratings;
- overall and local confidence;
- summary;
- Area / Factor Breakdown for the reported Area;
- local Requirements; and
- limits and incomplete inputs.

Area reports **MUST NOT** render standalone `Rating Drivers` sections or
`Driver | Effect | Inputs` tables. Rating drivers remain available in the
structured Area Analysis Result payloads listed in the report's `Source Data`
section.

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
- summary;
- direct Requirements;
- direct Sub-Factors; and
- limits and incomplete inputs.

Factor reports **MUST NOT** render standalone `Rating Drivers` sections or
`Driver | Effect | Inputs` tables. Rating drivers remain available in the
structured Factor Analysis Result payloads listed in the report's `Source Data`
section.

## Requirement Reports

Requirement reports **MUST** include:

- owning Area link;
- attached Factor links in a plural `Factors:` context line;
- kind-prefixed Requirement title;
- Requirement rating status and selected rating when present;
- assessment status;
- confidence;
- summary;
- findings summary;
- finding detail sections; and
- unknowns and missing evidence.

Requirement report finding detail sections **MUST** provide stable anchors
derived from the finding ID. Ranked findings links **MUST NOT** depend on
statement wording.

Requirement report finding detail sections **MUST** render Advice ranking
context when the finding appears in `FindingRankingResult`: Advice rank as
`<rank> / <total ranked findings>`, assigned finding ID, tier, and ranking
rationale. When no matching ranking entry exists, the section **MUST** render an
explicit not-ranked state.

Finding detail sections **MUST NOT** render finding-local `candidateActions`.
Candidate actions remain finding-local raw material; selected next moves belong
in `RecommendationResult` and generated recommendation reports.

## Rendering Rules

Reports **MUST** render empty tables with explicit empty-state rows.

Requirement report Finding sections **MUST** render the same
list columns: `ID`, `Statement`, `Type`, `Severity`, `Confidence`, `Effect`, and
`Basis`.

Finding detail sections **MUST** render the Finding Core in this order:
condition, criteria, basis, effect, and evidence. Requirement Finding details
**MUST NOT** render `candidateActions`.

Area and Factor reports **MUST NOT** render `Findings` sections. Their
human-facing roll-up explanation belongs in summary, ratings, confidence,
limits, incomplete inputs, and breakdown tables. Structured `ratingDrivers`
remain available through report `Source Data` links and routine JSON payloads,
not standalone Markdown body sections.

Run and Area reports **MUST** render an `Area / Factor Breakdown` section before
scope or Requirement detail sections. The breakdown table **MUST** use the
columns `Area / Factor`, `Overall Rating`, `Local Rating`, `Findings`, and
`Recommendations`, in that order. The `Area / Factor` cell **MUST** render the
row subject as the generated human report link when that report exists, and the
table **MUST NOT** render a separate `Report` column.

The run report's Area / Factor Breakdown **MUST** list the scoped Area as the
first row, followed by in-scope descendant Areas and Factors in deterministic
model order. An Area report's Area / Factor Breakdown **MUST** list the reported
Area as the first row, followed by its evaluated descendant Areas and Factors in
deterministic model order. The first row **MUST** emphasize only the table's
root Area in the `Area / Factor` cell. Factor rows **MUST** carry the Factor
report-kind marker inline in the `Area / Factor` cell instead of using a
separate Kind column.

The `Findings` column **MUST** count ranked findings that resolve to each row's
Area or Factor. The `Recommendations` column **MUST** count ranked
recommendations that resolve to each row's Area or Factor. A ranked
recommendation with multiple trace refs **MUST** count at most once for a given
breakdown row.

> Rationale: `Subject Reports` was a generated-file manifest rather than a
> quality overview, and separate Area `Factors` / `Child Areas` tables forced
> readers to assemble the local model shape by kind. A single narrow breakdown
> keeps navigation and quality signals together while leaving the machine report
> manifest in `EvaluationOutputResult.reportOutputs`. — 0161

Report headers **SHOULD** use report-specific summary tables instead of a
generic `Field | Value` key-value table. Run reports should summarize
`Overall Rating`, `Scope`, and `Confidence`; Area headers should summarize
`Overall Rating`, `Local Rating`, and `Confidence`; Factor headers should
summarize `Overall Rating`, `Local Rating`, `Status`, and `Confidence`;
Requirement headers should summarize `Rating`, `Assessment`, and `Confidence`;
findings and recommendations indexes should summarize index-specific counts and
priority signals; attached Factors belong in the plural `Factors:` context line,
not in the summary table.

> Rationale: the title identifies the report subject, so the header table should
> prioritize state and navigation rather than repeat the subject kind as
> metadata. The subject kind now rides the H1 title; location rides the
> navigation trail, so separate `Path:` / `Name:` header lines would be
> redundant. — 0104, 0119

Run report frontmatter `title` and H1 text **MUST** render as
`Quality Evaluation - <Area title>` for Area-only planned scopes. When the
planned scope has a factor filter, the run report frontmatter `title` and H1
text **MUST** render as
`Quality Evaluation - <Area title> (<Factor title list>)`, where
`<Factor title list>` contains every planned factor filter as comma-separated
Factor titles in `RunManifest.plannedScope.factorFilter` order. The run report
title **MUST NOT** include `Evaluation Report`, `Area:`, raw Area references, or
raw Factor references; stable scope references belong in the Scope section.

> Rationale: `report.md` and the report `type` already identify the artifact as
> a report. The H1 should name the quality-evaluation scope, while
> factor-scoped evaluations preserve both the Area context and the user's
> requested Factors. — 0168

Long reports **SHOULD** include a compact `Jump to:` line after the header when
local section navigation materially improves scanning. Short reports may omit
the line when it would add noise.

When a report table cell would otherwise render an empty scalar value, including
one component of a paired Confidence or Status cell, the cell **MUST** render an
em dash (`—`) instead of a blank segment. Empty whole-section placeholder rows
such as `(no findings)` and `(none recorded)` **MUST** remain worded empty-state
rows rather than being replaced by the cell marker. Generated report table cells
**MUST** escape Markdown table separators and normalize multiline scalar content
so persisted Evaluation text cannot alter the table column shape.

Each generated report **MUST** include exactly one static legend at the foot of
the report defining `—` as "not applicable or not recorded". The legend **MUST**
render regardless of whether the report contains an em-dash cell.

> Rationale: blank cells are ambiguous in committed Markdown reports. A neutral
> em dash makes absence visible without overclaiming `N/A`; escaping table
> separators and normalizing multiline text prevents persisted structured data
> from corrupting the generated Markdown table; and a static legend avoids
> data-dependent footnote churn across re-runs. — 0118, 0157

Every rating column **MUST** name what it rates. A header summary table **MUST**
label its descendant-inclusive rating column `Overall Rating` and its local
rating column `Local Rating`.

> Rationale: the adjacent header columns are self-describing nouns, so bare
> `Overall` / `Local` made a reader supply the missing noun. — 0111

The Factor report Sub-Factors table lists a Factor's immediate descendant
Factors, one row per child. It **MUST** render a `Local Rating` column from the
child's `localAnalysis` rating and a descendant-inclusive `+ Sub-Factors Rating`
column from the child's `localAndDescendantAnalysis` rating. It **MUST NOT**
render a boolean in a rating column. When a row's subject has no descendant
Factors, its `+ Sub-Factors Rating` cell **MUST** render an em dash (`—`) rather
than repeating the local rating.

> Rationale: these breakdown tables previously rendered the aggregate rating in
> the local `Rating` column and a `Yes`/`No` boolean where the roll-up rating
> belonged, leaving the local rating unshown — the unmet distinction clean-break
> case 0097 required. The em dash preserves the old boolean's "has children"
> signal without presenting a redundant rating. — 0097, 0111

Factor reports **MUST** render the immediate descendant-Factor section heading
as `Sub-Factors` and its empty-state row as `(no Sub-Factors)`. Reports
**MUST NOT** use `Sub-Areas` or `Child Factors` for generated human-facing
labels.

> Rationale: the Model vocabulary names Factor descendants as sub-factors.
> Generated reports should not make the same relationship look like a different
> concept. — 0147

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
- findings in Requirement Assessment Result order; and
- evidence in recorded order.
