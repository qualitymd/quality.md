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

The full ranked Findings report **MUST** be generated as `findings.md` at the run
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

The Recommendations report **MUST** be generated as `recommendations.md` at the
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

Every generated Markdown report **MUST** begin with YAML frontmatter.

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

The run-level `report.md` frontmatter **MUST** include only these fields:
`type`, `title`, `evaluationId`, `created`, `model`, and `run`. `evaluationId`
and `created` **MUST** render from `EvaluationManifest.evaluationId` and
`EvaluationManifest.createdAt`. `model` **MUST** render from
`EvaluationManifest.model`. `run` **MUST** identify the run folder label when
available.

Non-run report frontmatter **MUST** contain only `type` and `title`.

Report frontmatter **MUST NOT** duplicate ratings, confidence, summaries, rating
drivers, findings, recommendations, limits text, evidence, source-data
manifests, or rendered display labels.

> Rationale: `type` records the report artifact kind, while `title` records the
> generated Markdown document title a reader sees first. Run report frontmatter
> carries non-judgmental routing metadata so the visible opening can focus on
> judgment and navigation. Report-local source-data pointers live in the visible
> bottom section instead. — 0158, 0162, 0167, 0169

## Primary Source Data Section

Every generated Markdown report **MUST** end with a `## Primary Source Data`
section.

The `Primary Source Data` section **MUST** list the run-root-relative
structured Evaluation payload paths used as primary source data for that
specific Markdown report artifact. Primary source data is report-local: it
includes payloads that establish the report's identity, scope, subject result,
ranking, or recommendation content. It **MUST NOT** list descendant Area,
Factor, Requirement, assessment, rating, finding, or recommendation payloads
solely because the report links to, summarizes, or counts data from more
granular reports.

Each source-data list item **MUST** render as a Markdown link whose label is the
run-root-relative payload path and whose target is relative to the report file
that contains the section.

Reports that render run number, Evaluation ID, creation time, model path, or
requested scope from the Evaluation manifest **MUST** include
`data/evaluation-manifest.json`.

The `Primary Source Data` section **MUST NOT** include
`data/evaluation-output-result.json` solely because that generated output index
exists. A report **MAY** list `data/evaluation-output-result.json` only if that
report is directly rendered from it.

> Rationale: The bottom section keeps source payloads visible and parseable for
> people and agents without making frontmatter noisy. Path labels make the
> source files discoverable in plain text, while report-relative targets keep
> nested reports navigable. `EvaluationOutputResult` indexes generated outputs
> after report build; it is not source data for those reports unless a renderer
> explicitly consumes it. Primary source data keeps high-level reports from
> duplicating granular report provenance while preserving a stable bridge to the
> structured inputs for the current report. — 0159, 0162, 0171

## Fixed Enum Display

Generated Markdown reports **MUST** render known fixed Evaluation enum values
with the shared marker-plus-label display for that vocabulary. Structured JSON,
schemas, and receipts **MUST** preserve the raw canonical enum value.

Known fixed enum report displays are:

| Vocabulary            | Value                         | Display                          |
| --------------------- | ----------------------------- | -------------------------------- |
| Analysis status       | `analyzed`                    | `✅ Analyzed`                    |
| Analysis status       | `empty`                       | `⬜ Empty`                       |
| Analysis status       | `not_analyzed`                | `⚪ Not Analyzed`                |
| Analysis status       | `blocked`                     | `⛔ Blocked`                     |
| Assessment status     | `assessed`                    | `✅ Assessed`                    |
| Assessment status     | `partially_assessed`          | `🟡 Partially Assessed`          |
| Assessment status     | `not_assessed`                | `⚪ Not Assessed`                |
| Assessment status     | `blocked`                     | `⛔ Blocked`                     |
| Rating status         | `rated`                       | `✅ Rated`                       |
| Rating status         | `not_rated`                   | `⚪ Not Rated`                   |
| Rating status         | `blocked`                     | `⛔ Blocked`                     |
| Confidence            | `high`                        | `🟢 High`                        |
| Confidence            | `medium`                      | `🔵 Medium`                      |
| Confidence            | `low`                         | `🟡 Low`                         |
| Confidence            | `none`                        | `⚪ None`                        |
| Finding type          | `gap`                         | `🚩 Gap`                         |
| Finding type          | `risk`                        | `⚠️ Risk`                         |
| Finding type          | `strength`                    | `✅ Strength`                    |
| Finding type          | `note`                        | `ℹ️ Note`                         |
| Finding severity      | `critical`                    | `🔴 Critical`                    |
| Finding severity      | `high`                        | `🔴 High`                        |
| Finding severity      | `medium`                      | `🟡 Medium`                      |
| Finding severity      | `low`                         | `🔵 Low`                         |
| Finding basis         | `verified`                    | `✅ Verified`                    |
| Finding basis         | `plausible`                   | `🟡 Plausible`                   |
| Finding basis         | `not_assessed`                | `⚪ Not Assessed`                |
| Finding basis         | `not_applicable`              | `⬜ Not Applicable`              |
| Recommendation impact | `very_high`                   | `⬥⬥ Very high`                   |
| Recommendation impact | `high`                        | `⬥ High`                         |
| Recommendation impact | `medium`                      | `● Medium`                       |
| Recommendation impact | `low`                         | `○ Low`                          |
| Finding rank          | `P1`                          | `🔴 P1 Highest`                  |
| Finding rank          | `P2`                          | `🟠 P2 High`                     |
| Finding rank          | `P3`                          | `🟡 P3 Medium`                   |
| Finding rank          | `P4`                          | `⚪ P4 Low`                      |
| Finding coverage      | `addressed_by_recommendation` | `✅ Addressed by Recommendation` |
| Finding coverage      | `not_advice_driving`          | `⬜ Not Advice Driving`          |
| Report kind           | `run`                         | `📄 Run`                         |
| Report kind           | `area`                        | `🗺️ Area`                         |
| Report kind           | `factor`                      | `🧩 Factor`                      |
| Report kind           | `requirement`                 | `📋 Requirement`                 |
| Report kind           | `findings`                    | `🔝 Findings`                    |
| Report kind           | `recommendations`             | `📚 Recommendations`             |
| Report kind           | `recommendation`              | `💡 Recommendation`              |

Finding severity ordering **MUST** use `critical`, `high`, `medium`, then
`low`. Recommendation impact ordering **MUST** use `very_high`, `high`,
`medium`, then `low`. Finding ranking tier ordering **MUST** use `P1`, `P2`,
`P3`, then `P4`.

Each fixed Evaluation enum catalog **MUST** carry a human key label and concise
catalog description. Each fixed Evaluation enum value **MUST** carry a concise
value description. Generated reports **MUST NOT** render fixed enum catalog or
value descriptions inline; those descriptions belong in `glossary.md`.

> Rationale: fixed Evaluation values are strict machine data, but Markdown
> reports are a human scanning surface. Keeping labels, markers, ordering, and
> key labels in one contract prevents validation and report presentation from
> drifting. Descriptions belong in catalog metadata and the shared glossary, not
> repeated in generated report bodies. — 0173, 0179, 0183

## Evaluation Links

Generated Markdown reports **MUST** render a compact cross-artifact navigation
line labeled exactly `Evaluation links:`.

The line **MUST** render inline links in this order, separated by `|`:

```markdown
**Evaluation links:** [report.md](report.md) | [findings.md](findings.md) | [recommendations.md](recommendations.md) | [glossary.md](../../../glossary.md)
```

The link text **MUST** be the filename: `report.md`, `findings.md`,
`recommendations.md`, and `glossary.md`.

The `report.md`, `findings.md`, and `recommendations.md` links **MUST** target
the generated artifacts for the same Evaluation run. The `glossary.md` link
**MUST** target the workspace-root `glossary.md`. All four targets **MUST** be
relative to the current report artifact.

Generated reports **MUST** render the `Evaluation links:` line after the report
opening summary, key-details table, or report-specific orientation, and before
`## Contents` when a Contents section is present.

Generated reports **MUST NOT** link every table cell or enum value to
`glossary.md`.

Generated reports **MUST NOT** rely on `glossary.md` to make marker-only content
acceptable. Semantic table cells **MUST** render text labels for ratings,
statuses, confidence, finding type, severity, recommendation impact, finding
ranking tiers, and other priority-like values, optionally preceded by markers.

Generated reports **MUST NOT** render local `Legend` blocks or a bottom
`## Legend` section.

> Rationale: generated reports should stay readable on their own while using one
> stable link line for cross-artifact movement and glossary access. A single
> glossary-backed navigation line avoids repeating partial legends in every
> report and avoids noisy cell-level links. — 0183

## Run Report

The run-level `report.md` **MUST** render as the scoped Area report described by
`EvaluationManifest.plannedScope`. It **MUST** include:

- scoped Area title and rating;
- `## Summary`, `## Key Details`, `## Contents`, and `## Model Evaluation`
  before Top Findings;
- top 10 ranked findings;
- top 10 ranked recommendations;
- link to the Findings report;
- link to the Recommendations report;
- summary from the scoped Area result;
- Model Evaluation table for the scoped Area;
- requested Evaluation scope in Key Details; and
- `Evaluation links:` navigation.

The run-level `report.md` **MUST NOT** render the visible top `Run:` context
line used by detail reports. It **MUST NOT** render the top `Area:` context line
used by detail reports.

The run report `## Summary` section **MUST** render the scoped Area summary.
It **MUST NOT** render a `Recommended next action:` sentence.

The run report `## Key Details` section **MUST** render a table with `Overall
Rating`, `Confidence`, `Scope`, `Findings`, and `Recommendations`, in that
order. `Confidence` **MUST** render the scoped Area confidence paired with the
visible `Overall Rating`, not a paired overall/local confidence value. `Scope`
**MUST** render a human-readable description of the evaluated Area and filtered
Factors when present. `Findings` and `Recommendations` **MUST** render total
ranked counts as `<N> total` and **MUST NOT** include the word `ranked`. The
section **MUST NOT** include limits or incomplete-input counts.

The run report **MUST** render a compact `Finding Summary` table near
`## Key Details`. The table **MUST** render one row for every Finding type,
including zero-count types, with columns `Finding Type`, `Count`, and
`Severity`. For `gap` and `risk`, the `Severity` cell **MUST** render only
observed severity counts ordered by the Finding severity catalog. It **MUST NOT**
render zero-count severity values. For `strength`, `note`, and concern types
with no observed severities, the `Severity` cell **MUST** render `—`.

Generated reports **MUST** render a `## Contents` section when they contain at
least two substantive top-level body sections. Generated `## Contents` sections
**MUST** render a simple bullet list of Markdown links to visible `##` sections
in the same report, excluding the `Contents` section itself. Generated Contents
**MUST NOT** include nested `###` or deeper headings.

Generated reports **MUST NOT** render compact `Jump to:` lines.

Generated reports **MUST NOT** render a `## Contents` section when the artifact
is an OKF `index.md`, another listing/index artifact whose primary purpose is
navigation, or a report with fewer than two substantive top-level body sections.
The `## Primary Source Data` section **MUST** be eligible for generated Contents
when it is one of multiple substantive top-level sections in a generated report.

> Rationale: report artifacts are reader-facing Markdown documents, and standard
> Contents sections give readers and agents a predictable way to scan
> multi-section reports. Index files are already navigation artifacts, so a
> Contents section would duplicate their purpose. — 0175

The run-level `report.md` **MUST NOT** render `## Scope`, `## Coverage`, or
`## Report Details` sections.

The run-level `report.md` **MUST NOT** render a `## Limits & Incomplete Inputs`
section.

The run-level `report.md` **MUST NOT** render a standalone `Rating Drivers`
section or `Driver | Effect | Inputs` table. Rating drivers remain structured
source data available through routine JSON payloads and granular report
`Primary Source Data` sections.

When `plannedScope.factorFilter` is non-empty, `report.md` **MUST** identify the
filtered Factors in visible report content such as the H1 title.

The run report **MUST** state when the root Area was not evaluated in the run,
but it **MUST NOT** use a `## Coverage` section for that signal.

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
Because the run-report tables are capped overview tables, those full-list links
**MUST** render as emphasized sentence-case labels followed by filename-as-label
links and the complete ranked count outside the link text:

```markdown
**Full findings report:** [findings.md](findings.md) (7 total)
**Full recommendations report:** [recommendations.md](recommendations.md) (3 total)
```

The count **MUST** reflect the complete ranked list for the linked report, not
the number of rows rendered in the capped `report.md` table.

The Top Recommendations table **MUST** render rows from
`RecommendationRankingResult.orderedRecommendations` ordered by rank and capped
at 10 rows. It **MUST** render the columns `#`, `Recommendation`, `Area /
Factors`, `Impact`, `Confidence`, and `Reason`, in that order. The `#` cell
**MUST** render the user-facing recommendation number derived from the ranking
entry's `rank`. The table **MUST NOT** render a separate `Rank` column. The
`Recommendation` cell **MUST** use `RecommendationResult.title` as link text and
link to the generated recommendation detail report. The `Area / Factors` cell
**MUST** render linked Area and Factor names resolved from
`RecommendationResult.traceRefs` through persisted evaluation data and the model
snapshot, or `—` when no Area or Factor can be resolved. The `Impact` cell
**MUST** render the shared recommendation impact display label. The `Confidence`
cell **MUST** render the shared confidence display label from the Recommendation
ranking entry. The `Reason` cell **MUST** render
`RecommendationResult.expectedValue`.

## Finding Reports

`findings.md` **MUST** render a complete ranked Findings report from
`FindingRankingResult`. It **MUST** include:

- a `## Ranked Findings` section;
- all ranked findings ordered by rank;
- the same columns and link behavior as the run report Top Findings table.

## Recommendation Reports

`recommendations.md` **MUST** render a complete Recommendations report from
persisted `RecommendationResult` payloads and `RecommendationRankingResult`.
It **MUST** include:

- a `## Ranked Recommendations` section;
- all ranked recommendations with a `#` column derived from
  `RecommendationRankingResult.orderedRecommendations[].rank` and no separate
  `Rank` column;
- Area / Factors links resolved from `RecommendationResult.traceRefs`;
- impact;
- confidence;
- Reason from `RecommendationResult.expectedValue`;
- ranking rationale;
- links to recommendation detail reports; and
- a coverage summary from `findingCoverage`.

Each recommendation detail report **MUST** include:

- recommendation title;
- recommendation number;
- assigned recommendation ID;
- typed recommendation artifact reference;
- impact;
- confidence;
- description;
- background;
- expected value;
- done criterion;
- ranking rationale when ranked;
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

Known recommendation impact values **MUST** render consistently across
recommendation Markdown report locations using the shared fixed enum display.
Unknown impact values **MUST** render as humanized plain labels without a
marker.

## Navigation

Every report **MUST** render its H1 title line as the first Markdown content
after frontmatter. The H1 **MUST** prefix the subject display title with the
report kind: `Area:` for root and non-root Area reports, `Factor:` for Factor
reports, `Requirement:` for Requirement reports, and `Recommendation:` for
recommendation detail reports. The run-level H1 **MUST** identify the report as
a Quality Evaluation. The H1 title line and frontmatter `title` **MUST** use the
same plain-text title.

Every non-run report **MUST** render a run context line near the H1. Every
non-run report **MUST** render a report navigation line near the H1. The report
navigation line **MUST** link to the run overview `report.md`, full findings
report `findings.md`, and Recommendations report `recommendations.md` when the
current report is not that target. The run-level `report.md` **MUST NOT** render
that report navigation line.

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
columns or equivalent header source-data lines; the `Primary Source Data`
section owns those pointers.

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
structured Area Analysis Result payloads listed in the report's
`Primary Source Data` section.

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
structured Factor Analysis Result payloads listed in the report's
`Primary Source Data` section.

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
remain available through report `Primary Source Data` links and routine JSON
payloads, not standalone Markdown body sections.

Run reports **MUST** render a `Model Evaluation` section before Top Findings.
Area reports **MUST** render an `Area / Factor Breakdown` section before
Requirement detail sections. Both sections use the same breakdown table columns:
`Area / Factor`, `Overall Rating`, `Local Rating`, `Findings`, and
`Recommendations`, in that order. The `Area / Factor` cell **MUST** render the
row subject as the generated human report link when that report exists, and the
table **MUST NOT** render a separate `Report` column.

The run report's Model Evaluation table **MUST** list the scoped Area as the
first row, followed by in-scope descendant Areas and Factors in deterministic
model order. An Area report's Area / Factor Breakdown **MUST** list the reported
Area as the first row, followed by its evaluated descendant Areas and Factors in
deterministic model order. The first row **MUST** emphasize only the table's
root Area in the `Area / Factor` cell. Area rows **MUST** carry the `▦` marker,
and Factor rows **MUST** carry the `□` marker, inline in the `Area / Factor`
cell instead of using a separate Kind column. The marker **MUST** be part of the
row subject's visible link text when a report link exists.

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
Findings and Recommendations reports should summarize list-specific counts and
priority signals; attached Factors belong in the plural `Factors:` context line,
not in the summary table.

Opening summary tables **MUST** render under `## Key Details` when they are part
of the report opening.

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
Factor titles in `EvaluationManifest.plannedScope.factorFilter` order. The run
report title **MUST NOT** include `Evaluation Report`, `Area:`, raw Area
references, or raw Factor references; stable scope references belong in
`data/evaluation-manifest.json` and human-readable scope belongs in Key Details.

> Rationale: `report.md` and the report `type` already identify the artifact as
> a report. The H1 should name the quality-evaluation scope, while
> factor-scoped evaluations preserve both the Area context and the user's
> requested Factors. — 0168

When a report table cell would otherwise render an empty scalar value, including
one component of a paired Confidence or Status cell, the cell **MUST** render an
em dash (`—`) instead of a blank segment. Empty whole-section placeholder rows
such as `(no findings)` and `(none recorded)` **MUST** remain worded empty-state
rows rather than being replaced by the cell marker. Generated report table cells
**MUST** escape Markdown table separators and normalize multiline scalar content
so persisted Evaluation text cannot alter the table column shape.

Generated reports **MUST NOT** render blank table cells for empty scalar values.
When an em dash appears as an empty-cell marker, the report **MUST NOT** define
it through a local Legend block.

> Rationale: blank cells are ambiguous in committed Markdown reports. A neutral
> em dash makes absence visible without overclaiming `N/A`; escaping table
> separators and normalizing multiline text prevents persisted structured data
> from corrupting the generated Markdown table. The shared glossary and the
> report's text labels now own marker definitions, so local keys no longer need
> to define the absence marker beside every table. — 0118, 0157, 0174, 0183

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
unknown/missing-evidence types, known finding classifications, finding basis
statuses, recommendation impacts, finding ranking tiers, and finding coverage
dispositions, with human-readable display titles in Markdown while preserving
the raw values in routine JSON, `EvaluationOutputResult`, schemas, and
report-build receipts.

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
