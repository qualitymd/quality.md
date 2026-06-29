---
type: How-to Guide
title: Designing report outputs
description: How to design generated Evaluation report Markdown for fast human reading and light agent accessibility.
tags: [reports, evaluation, design, contributing]
timestamp: 2026-06-29T00:00:00Z
---

# Designing report outputs

Generated Evaluation reports are primarily human review surfaces. They should
help a reader land in any report file, understand the run context in a few
seconds, and move to the surrounding reports without reconstructing the tree from
paths.

Use this guide when designing or reshaping generated Markdown reports,
especially `report.md`, `findings.md`, `recommendations.md`, Area reports,
Factor reports, Requirement reports, and recommendation detail reports.

## Research basis

- Lead with the answer. Report-writing guidance commonly treats the summary as
  the part a busy reader may rely on most: it should state the purpose or
  context, principal findings, conclusion, and recommended next move without
  requiring a full read.
- People scan before they read. Nielsen Norman Group's web-writing research
  found higher measured usability from concise, scannable, objective content,
  and the strongest result when all three were combined:
  [Concise, SCANNABLE, and Objective](https://www.nngroup.com/articles/concise-scannable-and-objective-how-to-write-for-the-web/).
- Structure is navigation. W3C WAI guidance emphasizes meaningful page
  structure, headings, and regions:
  [Page Structure Tutorial](https://www.w3.org/WAI/tutorials/page-structure/).
  Multi-section generated reports benefit from a standard local Contents section
  when the headings below are stable and useful.
- Metadata should be useful, not exhaustive. Page metadata patterns commonly
  reserve metadata for facts such as date, type, owner, category, and reference
  number while including only what helps the reader:
  [Page metadata](https://designsystem.gov.scot/components/page-metadata).
- Breadcrumbs are useful when readers enter a nested structure directly. GOV.UK
  frames breadcrumbs as a way to show where users are and how to move up:
  [Breadcrumbs](https://design-system.service.gov.uk/components/breadcrumbs/).
- Frontmatter is a good place for document metadata and relationships, but it
  should not replace visible content. Hugo describes frontmatter as metadata for
  content; Docusaurus treats Markdown frontmatter as optional metadata:
  [Hugo front matter](https://gohugo.io/content-management/front-matter/),
  [Docusaurus docs](https://docusaurus.io/docs/create-doc).

## Principles

### Human first

The Markdown body is the report. A person opening a report in GitHub, an editor,
or a terminal pager should not need to read YAML or JSON to understand what they
are seeing.

Use visible Markdown for:

- report title and kind;
- run context;
- decision-relevant scope;
- navigation to sibling and parent reports;
- rating, confidence, status, and summary;
- findings and recommendations;
- limits and incomplete inputs.

Use frontmatter only for non-judgmental document metadata that makes report
identity, routing, and indexing cheap for agents or secondary tooling. Put
source-data links in the visible `Primary Source Data` section at the end of the
report.

### One-second orientation

The top of every report should answer:

- What report is this?
- What Evaluation run produced it?
- What subject or index does it cover?
- What was the planned scope?
- Where are the overview, findings, recommendations, and relevant parent or
  sibling reports?

Prefer short lines and compact tables over paragraphs in the header.

### Opening stack

The top of a primary run report should have one stable job per section:

1. H1: report kind and subject.
2. `Summary`: bottom line, main reason, and consequence.
3. `Key Details`: compact table of scan-critical facts.
4. `Contents`: local navigation for the major sections in this report.
5. `Model Evaluation`: the evaluated Model structure and ratings.
6. Ranked findings and recommendations.

Keep these roles distinct. Do not repeat the same fact in navigation, summary
prose, key-details table, and local jump links. The opening area should read as one
answer, not a pile of metadata.

For run reports, frontmatter owns routing metadata. Do not repeat
frontmatter-only facts such as run ID, creation time, and stable subject
reference in the visible body unless they materially help human judgment. Keep
scope visible when it changes how the rating should be read, usually in
`Key Details` and the H1 title for factor-scoped runs.

### Deterministic projection

Generated reports are deterministic projections over completed structured
Evaluation data. Do not introduce report-only findings, ratings, evidence,
limits, recommendations, or analysis in the Markdown body or in frontmatter.

Structured data under `data/` remains the machine-readable source of truth.
The report body may point to that data in its bottom `Primary Source Data`
section, but the report must not become a second result format.

Keep rating drivers in structured Evaluation payloads instead of rendering
standalone `Rating Drivers` body sections. If a reader or agent needs the full
rating trace, the report's `Primary Source Data` section should point them to
the relevant primary analysis payload; the visible body should foreground
summaries, ratings, findings, recommendations, model evaluation, limits, and
incomplete inputs.

### Report-specific headers

Use report-specific summary tables instead of a generic `Field | Value` table.
The title already identifies the subject; the header table should emphasize
state and scan-critical context.

Good header table examples:

- run report: `Overall Rating`, `Scope`, `Confidence`;
- Area report: `Overall Rating`, `Local Rating`, `Confidence`;
- Factor report: `Overall Rating`, `Local Rating`, `Status`, `Confidence`;
- Requirement report: `Rating`, `Assessment`, `Confidence`;
- Findings report: `Findings`, `Highest Severity`;
- Recommendations report: `Recommendations`, `Highest Impact`, `Coverage`.

For `report.md`, render this table under `## Key Details` instead of leaving it
as an unlabeled header table. The key-details table is for facts: rating,
confidence, scope, ranked-finding count, and ranked-recommendation count. Do not
use the key-details table for summary prose, recommendation text, evidence, or
limits.

### Navigation at the top

Detail reports should expose report-level navigation near the H1:

- `report.md` as the overview;
- `findings.md` as the Findings report;
- `recommendations.md` as the Recommendations report;
- parent Area or Factor reports when relevant;
- attached Factor reports for Requirement reports.

Keep the existing model context lines:

- `Area:` trail for Area, Factor, Requirement, and scoped run reports;
- `Factor:` trail for Factor reports;
- plural `Factors:` line for Requirement reports.

Use breadcrumb trails for hierarchy and a separate report-nav line for sibling
report-list artifacts. Do not overload breadcrumbs with every report in the run.

### Summary, key details, and contents

Put the bottom line before long tables. For the run report, use an explicit
`## Summary` section before Top Findings and Top Recommendations. The summary is
prose: it should state the overall judgment, the main reason, the quality
consequence, and the best next move when a recommendation is available.

Use `## Key Details` for scan-critical facts. This prevents the summary from
turning into a metadata dump and prevents the key-details table from trying to
explain judgment.

Use a `## Contents` section for generated report artifacts that contain multiple
substantive top-level sections. Keep Contents shallow: list visible `##`
sections in the report, do not include `Contents` itself, and do not nest `###`
headings.

Do not render compact `Jump to:` lines. A single Contents idiom keeps report
navigation predictable for people and agents.

Do not add Contents to OKF `index.md` files or other listing/index artifacts
whose primary job is navigation. Do not add Contents to a generated report with
fewer than two substantive top-level body sections.

Useful run-report targets, when the corresponding sections are rendered:

- `Model Evaluation`;
- `Top Findings`;
- `Top Recommendations`;
- `Primary Source Data`;

Useful detail-report targets, when the corresponding sections are rendered:

- `Area / Factor Breakdown`;
- `Limits & Incomplete Inputs`;
- `Findings Summary`;
- `Finding Details`;
- `Unknowns & Missing Evidence`.

For `report.md`, render Contents after the opening summary and key details, so
the reader sees the bottom line before document navigation. For detail and list
reports, render Contents after the opening context and key details.

### Frontmatter

Generated report YAML frontmatter should stay readable and non-judgmental. It is
an identity and indexing layer, not a second Evaluation result format,
metadata-summary table, or source-data manifest.

It is reasonable for frontmatter to carry document metadata such as stable report
kind, title, run number or slug, run ID, creation time, requested scope, and
subject reference when those fields help agents, static-site tooling, or editors
find and route the report without making the file harder for people to open.
These fields are metadata about the report document and run, not judgment about
the evaluated subject.

Do not repeat Evaluation judgment or evidence in frontmatter, including:

- summaries;
- ratings;
- confidence;
- rating drivers;
- findings;
- recommendations;
- limits text;
- evidence text;
- source-data manifests;
- rendered display labels.

Generated reports are runtime artifacts, so they do not yet require a report
bundle `index.md`, `schema.md`, or `log.md`. A later Change Case should decide
the full report-bundle contract and update the generated report specs.
Do not rename report files such as `findings.md` or `recommendations.md` to
`index.md` for OKF compatibility. In OKF, `index.md` is a bundle listing file;
`findings.md` is still a report concept whose subject is the ranked Findings
index.

Use a stable report `type` and a human-friendly `title` in every generated
report. Additional non-judgmental metadata fields are acceptable only when the
durable report spec for that artifact allows them:

```yaml
---
type: Evaluation Overview Report
title: "Quality Evaluation - LedgerLite Service"
run: 0001-full-eval
runId: 20260629T120000Z-0123456789ab
created: 2026-06-29T12:00:00Z
scope: full evaluation
---
```

Good frontmatter fields:

- `type`;
- `title`;
- `run`;
- `runId`;
- `created`;
- `scope`;
- `subject`.

The frontmatter `title` should match the visible H1 title text without the
leading Markdown `#` marker. Detail reports keep prefixes such as
`Requirement:` or `Area:` in both the H1 and frontmatter title when they are
part of the report document title; the `type` field already carries the report
artifact taxonomy. The run report should title the user's quality-evaluation
scope directly as `Quality Evaluation - <Area>` and append factor filters in
parentheses, for example `Quality Evaluation - Public API (Reliability,
Correctness)`.

Use report types that name the reported subject and keep the word `Report` in
the type so report artifacts do not collide with Model concepts:

| Output                                | `type`                          | Subject                        |
| ------------------------------------- | ------------------------------- | ------------------------------ |
| `report.md`                           | `Evaluation Overview Report`    | Evaluation run / planned scope |
| `root-area.md`                        | `Area Evaluation Report`        | Root Area                      |
| `areas/<area>/<area>-area.md`         | `Area Evaluation Report`        | Area                           |
| `factors/<factor>/<factor>-factor.md` | `Factor Evaluation Report`      | Factor                         |
| `requirements/<requirement>/...`      | `Requirement Evaluation Report` | Requirement                    |
| `findings.md`                         | `Finding Index Report`          | Ranked Findings                |
| `recommendations.md`                  | `Recommendation Index Report`   | Ranked Recommendations         |
| `recommendations/<NNN>-<slug>.md`     | `Recommendation Report`         | Recommendation                 |

Avoid `type: Area`, `type: Factor`, or `type: Requirement`; those name Model
concepts, not report artifacts. `report.md` should use
`Evaluation Overview Report` rather than `Area Evaluation Report` because it is
the run entrypoint and includes findings, recommendations, scope, coverage, and
subject report links.

This keeps the first lines of a report readable in editors that expose
frontmatter, makes frontmatter `title` the same document title as the H1, and
keeps the structured Evaluation data as the single source of truth for
judgment. If a future consumer needs richer machine access to ratings, findings,
recommendations, evidence, or limits, use `data/evaluation-output-result.json`
and the payloads it indexes instead of expanding generated report frontmatter.

The visible H1 remains the first Markdown content after report frontmatter.

### Local keys and indicator labels

Generated reports may use markers or icons to make values scannable, but the
text label is authoritative. Report content must render every rating, status,
confidence, severity, finding type, recommendation impact, and priority-like
value as a text label, optionally preceded by a marker or icon. Markers are
supplemental scanning aids; do not rely on color or icon shape alone.

Use short local keys immediately after the first table or section that uses an
indicator family in that artifact. Keys are notation-only: one family per line,
with the family label followed by the marker-plus-label set. Do not add term
definitions, rationale, or explanatory prose to generated keys. Do not collect
keys in a bottom `Legend` section.

Examples:

```markdown
Ratings: 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.
```

```markdown
Rows: `▦` Area, `□` Factor.
Empty: `—`.
```

```markdown
Type: ✅ Strength, ⚠️ Gap, ⚠️ Risk, ❓ Unknown, ℹ️ Note.
Severity: 🔴 Critical, 🔴 High, 🟡 Medium, 🔵 Low.

## Contents

- [Ranked Findings](#ranked-findings)
- [Primary Source Data](#primary-source-data)

## Ranked Findings

| Rank | Finding | Area | Factors | Type | Severity |
```

```markdown
Impact: ◆ Very high, ▲ High, ● Medium, ○ Low.
```

For Rating Levels, render the configured Rating Scale from the run's model
snapshot. For CLI-owned display catalogs such as confidence, finding type,
severity, recommendation impact, finding ranking tier, and coverage disposition,
render the known catalog set. Unknown free-form values remain readable through
their rendered text labels and do not need invented key entries.

### Primary Source Data at the bottom

Every generated report should end with a stable `## Primary Source Data`
section. The section lists the report-local primary structured Evaluation
payloads used to render that report artifact. This is the only standard footer
utility section; local keys belong near first use.

```markdown
## Primary Source Data

- [data/run-manifest.json](data/run-manifest.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json](data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json](data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](data/advice/finding-ranking-result.json)
```

Use path-as-label links. Human-friendly labels would read better locally, but
paths are more useful to agents, avoid another naming surface, and make the
target obvious in plain text. In nested reports, keep the label
run-root-relative and make the link target relative to the report file, for
example:

```markdown
- [data/run-manifest.json](../../data/run-manifest.json)
```

Keep the list report-local. Do not list every payload in the run, and do not
list every payload used by more granular linked reports. Do not list
`data/evaluation-output-result.json` merely because it exists: that file is a
generated output index, not report source data unless a future renderer directly
consumes it. For `report.md`, prefer primary run/report inputs such as the run
manifest, scoped Area analysis, and advice ranking payloads; let linked detail
reports own their own deeper source lists.

## Header patterns

### Run report

The run report is the primary decision-ready report. Its opening should make the
result, key facts, evaluated Model shape, and ranked evidence obvious without
duplicating detail-report navigation or provenance sections.

```markdown
# Quality Evaluation - LedgerLite Service

## Summary

LedgerLite is usable in the synthetic evaluation, but API idempotency, rollback rehearsal, and recovery ownership keep the overall service below target.

## Key Details

| Overall Rating | Confidence    | Scope           | Findings | Recommendations |
| -------------- | ------------- | --------------- | -------- | --------------- |
| Minimum        | Medium / None | full evaluation | 7 ranked | 3 ranked        |

Ratings: 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.
Empty: `—`.

## Contents

- [Model Evaluation](#model-evaluation)
- [Top Findings](#top-findings)
- [Top Recommendations](#top-recommendations)
- [Primary Source Data](#primary-source-data)

## Model Evaluation

| Area / Factor                                                          | Overall Rating | Local Rating | Findings | Recommendations |
| ---------------------------------------------------------------------- | -------------- | ------------ | -------- | --------------- |
| **[▦ LedgerLite Service](root-area.md)**                               | Minimum        | —            | 7        | 3               |
| ↳ [▦ Public API](areas/api/api-area.md)                                | Minimum        | Minimum      | 2        | 1               |
| ↳ [□ Correctness](areas/api/factors/correctness/correctness-factor.md) | Minimum        | Minimum      | 1        | 1               |

Rows: `▦` Area, `□` Factor.

...

## Top Findings

Type: ✅ Strength, ⚠️ Gap, ⚠️ Risk, ❓ Unknown, ℹ️ Note.
Severity: 🔴 Critical, 🔴 High, 🟡 Medium, 🔵 Low.

## Top Recommendations

Impact: ◆ Very high, ▲ High, ● Medium, ○ Low.

...

## Primary Source Data

- [data/run-manifest.json](data/run-manifest.json)
- [data/areas/root/area-analysis-result.json](data/areas/root/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](data/advice/recommendation-ranking-result.json)
```

### Findings and recommendations

Finding and recommendation reports should keep their visible H1s simple while
using frontmatter `type` to carry the report artifact taxonomy.

```markdown
# Findings

Run: 0001-full-eval - Generated: 2026-06-29 18:42 UTC - Scope: full evaluation

Report: [Overview](report.md) - Findings - [Recommendations](recommendations.md)

| Findings          | Highest Severity |
| ----------------- | ---------------- |
| 7 ranked findings | High             |

Severity: 🔴 Critical, 🔴 High, 🟡 Medium, 🔵 Low.
```

### Area reports

Area reports should orient with both run navigation and Area hierarchy.

```markdown
# Area: Public API

Run: [0001-full-eval](../../report.md) - Generated: 2026-06-29 18:42 UTC - Scope: full evaluation

Report: [Overview](../../report.md) - [Findings](../../findings.md) - [Recommendations](../../recommendations.md)

Area: [LedgerLite Service](../../root-area.md) / [Public API](api-area.md)

| Overall Rating | Local Rating | Confidence      |
| -------------- | ------------ | --------------- |
| Minimum        | Minimum      | Medium / Medium |

Ratings: 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.

## Contents

- [Summary](#summary)
- [Area / Factor Breakdown](#area--factor-breakdown)
- [Requirements](#requirements)
- [Limits & Incomplete Inputs](#limits--incomplete-inputs)
- [Primary Source Data](#primary-source-data)

## Summary

The API has predictable errors, but idempotency retry semantics need a tighter contract.

## Area / Factor Breakdown

| Area / Factor                                                | Overall Rating | Local Rating | Findings | Recommendations |
| ------------------------------------------------------------ | -------------- | ------------ | -------- | --------------- |
| **[▦ Public API](api-area.md)**                              | Minimum        | Minimum      | 2        | 1               |
| ↳ [□ Correctness](factors/correctness/correctness-factor.md) | Minimum        | Minimum      | 1        | 1               |
| ↳ [□ Operability](factors/operability/operability-factor.md) | Target         | Target       | 1        | 0               |

Rows: `▦` Area, `□` Factor.
```

### Factor reports

Factor reports should preserve Area context and Factor hierarchy.

```markdown
# Factor: Correctness

Run: [0001-full-eval](../../../report.md) - Generated: 2026-06-29 18:42 UTC - Scope: full evaluation

Report: [Overview](../../../report.md) - [Findings](../../../findings.md) - [Recommendations](../../../recommendations.md)

Area: [LedgerLite Service](../../../root-area.md) / [Public API](../api-area.md)

Factor: [Correctness](correctness-factor.md)

| Overall Rating | Local Rating | Status              | Confidence      |
| -------------- | ------------ | ------------------- | --------------- |
| Minimum        | Minimum      | Analyzed / Analyzed | Medium / Medium |

Ratings: 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable.
Status: ✅ Analyzed, ⬜ Empty, ⚪ Not Analyzed, ⛔ Blocked.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.

## Contents

- [Summary](#summary)
- [Requirements](#requirements)
- [Sub-Factors](#sub-factors)
- [Limits & Incomplete Inputs](#limits--incomplete-inputs)
- [Primary Source Data](#primary-source-data)

## Summary

Correctness follows its direct requirement signal.
```

### Requirement reports

Requirement reports are often the deep link target from findings. Their header
should make the owning Area, attached Factors, rating state, assessment state,
and confidence immediately visible.

```markdown
# Requirement: mutation endpoints are idempotent under retry

Run: [Run 0001](../../../../report.md) - Run ID: `20260629T184200Z-0123456789ab` - Created: 2026-06-29T18:42:00Z - Scope: full evaluation

Report: [Overview](../../../../report.md) - [Findings](../../../../findings.md) - [Recommendations](../../../../recommendations.md)

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [correctness](../../factors/correctness/correctness-factor.md)

| Rating  | Assessment | Confidence      |
| ------- | ---------- | --------------- |
| Minimum | Assessed   | Medium / Medium |

Ratings: 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable.
Assessment: ✅ Assessed, 🟡 Partially Assessed, ⚪ Not Assessed, ⛔ Blocked.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.

## Contents

- [Summary](#summary)
- [Findings Summary](#findings-summary)
- [Finding Details](#finding-details)
- [Unknowns & Missing Evidence](#unknowns--missing-evidence)
- [Primary Source Data](#primary-source-data)
```

### Recommendation detail reports

Recommendation reports are action-planning surfaces. Their headers should make
rank, impact, confidence, trace context, and navigation visible before the
recommendation prose.

```markdown
# Recommendation: Tighten the idempotency replay contract

Run: [Run 0001](../report.md) - Run ID: `20260629T184200Z-0123456789ab` - Created: 2026-06-29T18:42:00Z - Scope: full evaluation

Report: [Overview](../report.md) - [Findings](../findings.md) - [Recommendations](../recommendations.md)

Trace: [Public API](../areas/api/api-area.md) / [Correctness](../areas/api/factors/correctness/correctness-factor.md)

| # | Rank | Impact | Confidence | Reference                                                 |
| - | ---- | ------ | ---------- | --------------------------------------------------------- |
| 1 | 1    | ▲ High | High       | evaluation:20260629T184200Z-0123456789ab/recommendation/1 |

Impact: ◆ Very high, ▲ High, ● Medium, ○ Low.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary Source Data](#primary-source-data)
```

## Checklist

Before changing report output, check:

- The first visible Markdown content is one clear H1.
- A reader can identify the subject and decision-relevant scope in the opening
  area.
- Detail report navigation links to the overview, findings, and recommendations
  where those reports exist.
- Detail report hierarchical context is visible through Area, Factor, or Factors
  lines.
- `report.md` uses `## Summary`, `## Key Details`, `## Contents`, and
  `## Model Evaluation` before ranked lists.
- Frontmatter-only routing metadata is not repeated in the visible body unless
  it materially helps human judgment.
- Summary prose states judgment, reason, and consequence without becoming a
  metadata list or a recommendation table.
- Key details tables contain scan-critical facts, not prose, evidence, or
  recommendations.
- Generated report artifacts with multiple substantive top-level sections have
  `## Contents`; OKF `index.md` and other listing/index artifacts do not.
- Local keys appear after first use, render one family per line, and stay
  notation-only.
- Report cells keep text labels; markers and icons are supplemental and never
  the only meaning for semantic values.
- Generated reports do not render bottom `Legend` sections.
- Frontmatter contains readable, non-judgmental document metadata only.
- No generated report introduces claims that are absent from structured data.
- The bottom `Primary Source Data` section lists report-local primary payloads,
  not every transitive payload used by linked detail reports.
- Empty values render visibly and consistently.
