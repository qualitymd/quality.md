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
source-data links in the visible `Primary source data` section at the end of the
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
3. `Key details`: compact table of scan-critical facts.
4. `Contents`: local navigation for the major sections in this report.
5. `Model evaluation`: the evaluated Model structure and ratings.
6. Ranked findings and recommendations.

Keep these roles distinct. Do not repeat the same fact in navigation, summary
prose, key-details table, and local jump links. The opening area should read as one
answer, not a pile of metadata.

For run reports, frontmatter owns routing metadata. Do not repeat
frontmatter-only facts such as Evaluation ID, creation time, and stable subject
reference in the visible body unless they materially help human judgment. Keep
scope visible when it changes how the rating should be read, usually in
`Key details` and the H1 title for factor-scoped runs.

### Deterministic projection

Generated reports are deterministic projections over completed structured
Evaluation data. Do not introduce report-only findings, ratings, evidence,
limits, recommendations, or analysis in the Markdown body or in frontmatter.

Structured data under `data/` remains the machine-readable source of truth.
The report body may point to that data in its bottom `Primary source data`
section, but the report must not become a second result format.

Keep rating drivers in structured Evaluation payloads instead of rendering
standalone `Rating Drivers` body sections. If a reader or agent needs the full
rating trace, the report's `Primary source data` section should point them to
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
- Findings report: `Findings`, `Highest Concern Severity`;
- Recommendations report: `Recommendations`, `Highest Impact`, `Coverage`.

For `report.md`, render this table under `## Key details` instead of leaving it
as an unlabeled header table. The key-details table is for facts: rating,
confidence, scope, total Finding count, and total Recommendation count. Do not
use the key-details table for summary prose, recommendation text, evidence,
limits, or ranking-state wording.

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
`## Summary` section before Top findings and Top recommendations. The summary is
prose: it should state the overall judgment, the main reason, the quality
consequence, and the best next move when a recommendation is available.

Use `## Key details` for scan-critical facts. This prevents the summary from
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

- `Model evaluation`;
- `Top findings`;
- `Top recommendations`;
- `Primary source data`;

Useful detail-report targets, when the corresponding sections are rendered:

- `Area / Factor breakdown`;
- `Limits and incomplete inputs`;
- `Findings summary`;
- `Finding details`;
- `Unknowns and missing evidence`.

For `report.md`, render Contents after the opening summary and key details, so
the reader sees the bottom line before document navigation. For detail and list
reports, render Contents after the opening context and key details.

### Frontmatter

Generated report YAML frontmatter should stay readable and non-judgmental. It is
an identity and indexing layer, not a second Evaluation result format,
metadata-summary table, or source-data manifest.

It is reasonable for frontmatter to carry document metadata such as stable report
kind, title, Evaluation ID, creation time, model path, and run folder label when
those fields help agents, static-site tooling, or editors find and route the
report without making the file harder for people to open. These fields are
metadata about the report document and Evaluation, not judgment about the
evaluated subject.

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
title: "Quality evaluation - LedgerLite Service"
evaluationId: 20260629T120000Z-0123456789ab
created: 2026-06-29T12:00:00Z
model: QUALITY.md
run: 0001-full-eval
---
```

Good frontmatter fields:

- `type`;
- `title`;
- `evaluationId`;
- `created`;
- `model`;
- `run`.

The frontmatter `title` should match the visible H1 title text without the
leading Markdown `#` marker. Detail reports keep prefixes such as
`Requirement:` or `Area:` in both the H1 and frontmatter title when they are
part of the report document title; the `type` field already carries the report
artifact taxonomy. The run report should title the user's quality-evaluation
scope directly as `Quality evaluation - <Area>` and append factor filters in
parentheses, for example `Quality evaluation - Public API (Reliability,
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
| `findings.md`                         | `Finding Index Report`          | Ranked findings                |
| `recommendations.md`                  | `Recommendation Index Report`   | Ranked recommendations         |
| `recommendations/<number>-<slug>.md`  | `Recommendation Report`         | Recommendation                 |

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

### Evaluation links and indicator labels

Generated reports may use markers or icons to make values scannable, but the
text label is authoritative. Report content must render every rating, status,
confidence, severity, finding type, recommendation impact, and priority-like
value as a text label, optionally preceded by a marker or icon. Markers are
supplemental scanning aids; do not rely on color or icon shape alone.

In `Model evaluation` and `Area / Factor breakdown` tables, use
`▦ Area / □ Factor` as the first column label. This keeps the Area and Factor row
markers visible as a compact key while row cells still carry linked text labels.

Use one compact `Evaluation links:` blockquote immediately below the H1 in every
generated report. The blockquote links to the run overview, Findings,
Recommendations, and the workspace-root glossary:

```markdown
> **Evaluation links:** [report.md](report.md) | [findings.md](findings.md) | [recommendations.md](recommendations.md) | [glossary.md](../../../glossary.md)
```

Use filename link text. Link targets are relative to the current report
artifact. Do not link every table cell or enum value to `glossary.md`; the one
navigation blockquote keeps definitions reachable without making generated
Markdown noisy.

Do not render local `Legend` blocks. The workspace-root `glossary.md` owns
definitions and fixed vocabulary tables. Reports must still render text labels
in cells, so the reader should not need the glossary to understand basic table
values.

### Primary source data at the bottom

Every generated report should end with a stable `## Primary source data`
section. The section lists the report-local primary structured Evaluation
payloads used to render that report artifact. This is the only standard footer
utility section.

```markdown
## Primary source data

- [data/evaluation-manifest.json](data/evaluation-manifest.json)
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
- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
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
# Quality evaluation - LedgerLite Service

> **Evaluation links:** [report.md](report.md) | [findings.md](findings.md) | [recommendations.md](recommendations.md) | [glossary.md](../../../glossary.md)

## Summary

LedgerLite is usable in the synthetic evaluation, but API idempotency, rollback rehearsal, and recovery ownership keep the overall service below target.

## Key details

| Overall Rating | Confidence | Scope                                 | Findings | Recommendations |
| -------------- | ---------- | ------------------------------------- | -------- | --------------- |
| 🟡 Minimum     | 🔵 Medium  | Full evaluation of LedgerLite Service | 7 total  | 3 total         |

## Contents

- [Model evaluation](#model-evaluation)
- [Top findings](#top-findings)
- [Top recommendations](#top-recommendations)
- [Primary source data](#primary-source-data)

## Model evaluation

| ▦ Area / □ Factor                                                      | Overall Rating | Local Rating | Findings | Recommendations |
| ---------------------------------------------------------------------- | -------------- | ------------ | -------- | --------------- |
| **[▦ LedgerLite Service](root-area.md)**                               | Minimum        | —            | 7        | 3               |
| ↳ [▦ Public API](areas/api/api-area.md)                                | Minimum        | Minimum      | 2        | 1               |
| ↳ [□ Correctness](areas/api/factors/correctness/correctness-factor.md) | Minimum        | Minimum      | 1        | 1               |

...

## Top findings

**Full findings report:** [findings.md](findings.md) (7 total: 🚩 2 Gaps: 🔴 1 High, 🟡 1 Medium; ⚠️ 1 Risk: 🟡 1 Medium; 💪 4 Strengths)

## Top recommendations

**Full recommendations report:** [recommendations.md](recommendations.md) (3 total; impact: ⬥ 1 High, ● 2 Medium)

...

## Primary source data

- [data/evaluation-manifest.json](data/evaluation-manifest.json)
- [data/areas/root/area-analysis-result.json](data/areas/root/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](data/advice/recommendation-ranking-result.json)
```

### Findings and recommendations

Finding and recommendation reports should keep their visible H1s simple while
using frontmatter `type` to carry the report artifact taxonomy.

```markdown
# Findings

> **Evaluation links:** [report.md](report.md) | [findings.md](findings.md) | [recommendations.md](recommendations.md) | [glossary.md](../../../glossary.md)

Run: 0001-full-eval - Generated: 2026-06-29 18:42 UTC - Scope: full evaluation

| Findings   | Highest Concern Severity |
| ---------- | ------------------------ |
| 7 findings | High                     |
```

### Area reports

Area reports should orient with both run navigation and Area hierarchy.

```markdown
# Area: Public API

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Generated: 2026-06-29 18:42 UTC - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md) / [Public API](api-area.md)

| Overall Rating | Local Rating | Confidence            |
| -------------- | ------------ | --------------------- |
| 🟡 Minimum     | 🟡 Minimum   | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Area / Factor breakdown](#area--factor-breakdown)
- [Requirements](#requirements)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

The API has predictable errors, but idempotency retry semantics need a tighter contract.

## Area / Factor breakdown

| ▦ Area / □ Factor                                            | Overall Rating | Local Rating | Findings | Recommendations |
| ------------------------------------------------------------ | -------------- | ------------ | -------- | --------------- |
| **[▦ Public API](api-area.md)**                              | Minimum        | Minimum      | 2        | 1               |
| ↳ [□ Correctness](factors/correctness/correctness-factor.md) | Minimum        | Minimum      | 1        | 1               |
| ↳ [□ Operability](factors/operability/operability-factor.md) | Target         | Target       | 1        | 0               |
```

### Factor reports

Factor reports should preserve Area context and Factor hierarchy.

```markdown
# Factor: Correctness

> **Evaluation links:** [report.md](../../../report.md) | [findings.md](../../../findings.md) | [recommendations.md](../../../recommendations.md) | [glossary.md](../../../../../../glossary.md)

Run: [0001-full-eval](../../../report.md) - Generated: 2026-06-29 18:42 UTC - Scope: full evaluation

Area: [LedgerLite Service](../../../root-area.md) / [Public API](../api-area.md)

Factor: [Correctness](correctness-factor.md)

| Overall Rating | Local Rating | Status                    | Confidence            |
| -------------- | ------------ | ------------------------- | --------------------- |
| 🟡 Minimum     | 🟡 Minimum   | ✅ Analyzed / ✅ Analyzed | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Requirements](#requirements)
- [Sub-factors](#sub-factors)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

Correctness follows its direct requirement signal.
```

### Requirement reports

Requirement reports are often the deep link target from findings. Their header
should make the owning Area, attached Factors, rating state, assessment state,
and confidence immediately visible.

```markdown
# Requirement: mutation endpoints are idempotent under retry

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [Run 0001](../../../../report.md) - Evaluation ID: `20260629T184200Z-0123456789ab` - Created: 2026-06-29T18:42:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [correctness](../../factors/correctness/correctness-factor.md)

| Rating     | Assessment  | Confidence            |
| ---------- | ----------- | --------------------- |
| 🟡 Minimum | ✅ Assessed | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)
```

### Recommendation detail reports

Recommendation reports are action-planning surfaces. Their headers should make
the recommendation number, opaque ID, impact, confidence, trace context, and
navigation visible before the recommendation prose.

```markdown
# Recommendation: Tighten the idempotency replay contract

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [Run 0001](../report.md) - Evaluation ID: `20260629T184200Z-0123456789ab` - Created: 2026-06-29T18:42:00Z - Scope: full evaluation

Trace: [Public API](../areas/api/api-area.md) / [Correctness](../areas/api/factors/correctness/correctness-factor.md)

| # | ID           | Impact | Confidence | Reference                                                            |
| - | ------------ | ------ | ---------- | -------------------------------------------------------------------- |
| 1 | qrec_example | ⬥ High | 🟢 High    | evaluation:20260629T184200Z-0123456789ab/recommendation/qrec_example |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary source data](#primary-source-data)
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
- `report.md` uses `## Summary`, `## Key details`, `## Contents`, and
  `## Model evaluation` before ranked lists.
- Frontmatter-only routing metadata is not repeated in the visible body unless
  it materially helps human judgment.
- Summary prose states judgment, reason, and consequence without becoming a
  metadata list or a recommendation table.
- Key details tables contain scan-critical facts, not prose, evidence, or
  recommendations.
- Generated report artifacts with multiple substantive top-level sections have
  `## Contents`; OKF `index.md` and other listing/index artifacts do not.
- Every generated report has an `Evaluation links:` blockquote immediately below
  the H1, with filename link text and a relative link to `glossary.md`.
- Report cells keep text labels; markers and icons are supplemental and never
  the only meaning for semantic values.
- Generated reports do not render local or bottom `Legend` sections.
- Frontmatter contains readable, non-judgmental document metadata only.
- No generated report introduces claims that are absent from structured data.
- The bottom `Primary source data` section lists report-local primary payloads,
  not every transitive payload used by linked detail reports.
- Empty values render visibly and consistently.
