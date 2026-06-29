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

- People scan before they read. Nielsen Norman Group's web-writing research
  found higher measured usability from concise, scannable, objective content,
  and the strongest result when all three were combined:
  [Concise, SCANNABLE, and Objective](https://www.nngroup.com/articles/concise-scannable-and-objective-how-to-write-for-the-web/).
- Headings, lists, summaries, links, and tables make pages easier to scan.
  Apply that to reports with one clear H1, short context lines, summary tables,
  ranked tables, and jump links on long pages.
- Page structure is an accessibility feature. W3C WAI guidance emphasizes
  meaningful page structure, headings, and regions:
  [Page Structure Tutorial](https://www.w3.org/WAI/tutorials/page-structure/).
- Page headers should orient before they explain. Atlassian's page header
  component centers the page title with supporting breadcrumbs, metadata, and
  actions:
  [Page header](https://atlassian.design/components/page-header).
- Metadata should be useful, not exhaustive. The Scottish Government design
  system recommends page metadata for facts such as date, type, owner, category,
  and reference number, while including only what helps the reader:
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
- scope and coverage;
- navigation to sibling and parent reports;
- rating, confidence, status, and summary;
- findings and recommendations;
- limits and incomplete inputs.

Use frontmatter only to make report identity cheap for agents or secondary
tooling. Put source-data links in the visible `Source Data` section at the end
of the report.

### One-second orientation

The top of every report should answer:

- What report is this?
- What Evaluation run produced it?
- What subject or index does it cover?
- What was the planned scope?
- Where are the overview, findings, recommendations, and relevant parent or
  sibling reports?

Prefer short lines and compact tables over paragraphs in the header.

### Deterministic projection

Generated reports are deterministic projections over completed structured
Evaluation data. Do not introduce report-only findings, ratings, evidence,
limits, recommendations, or analysis in the Markdown body or in frontmatter.

Structured data under `data/` remains the machine-readable source of truth.
The report body may point to that data in its bottom `Source Data` section, but
the report must not become a second result format.

Keep rating drivers in structured Evaluation payloads instead of rendering
standalone `Rating Drivers` body sections. If a reader or agent needs the full
rating trace, the report's `Source Data` section should point them to the
relevant analysis payload; the visible body should foreground summaries,
ratings, findings, recommendations, Area / Factor breakdowns, limits, and
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
- findings index: `Findings`, `Highest Severity`;
- recommendations index: `Recommendations`, `Highest Impact`, `Coverage`.

### Navigation at the top

Every report should expose report-level navigation near the H1:

- `report.md` as the overview;
- `findings.md` as the full findings index;
- `recommendations.md` as the full recommendation index;
- parent Area or Factor reports when relevant;
- attached Factor reports for Requirement reports.

Keep the existing model context lines:

- `Area:` trail for Area, Factor, Requirement, and scoped run reports;
- `Factor:` trail for Factor reports;
- plural `Factors:` line for Requirement reports.

Use breadcrumb trails for hierarchy and a separate report-nav line for sibling
indexes. Do not overload breadcrumbs with every report in the run.

### Summaries before detail

Put the bottom line before long tables. For the run report, place the evaluation
summary before Top Findings and Top Recommendations. For detail reports, keep the
short subject summary close to the header table.

Long pages should add a `Jump to:` line after the header area. Useful targets:

- `Top Findings`;
- `Top Recommendations`;
- `Area / Factor Breakdown`;
- `Scope`;
- `Coverage`;
- `Limits & Incomplete Inputs`;
- `Findings Summary`;
- `Finding Details`;
- `Unknowns & Missing Evidence`.

Skip jump links on short pages where they add more noise than navigation value.

### Identity frontmatter

Generated report YAML frontmatter stays tiny. It is an identity layer, not a
metadata summary or source-data manifest. Keep it OKF-compatible by including
only a `type` and `title`. Do not repeat Evaluation result facts that already
live in the associated JSON files or visible body, including generated time, run
identity, subject identity, scope, ratings, confidence, findings,
recommendations, limits, or display labels.

Generated reports are runtime artifacts, so they do not yet require a report
bundle `index.md`, `schema.md`, or `log.md`. A later Change Case should decide
the full report-bundle contract and update the generated report specs.
Do not rename report files such as `findings.md` or `recommendations.md` to
`index.md` for OKF compatibility. In OKF, `index.md` is a bundle listing file;
`findings.md` is still a report concept whose subject is the ranked Findings
index.

Use only a stable report `type` and a human-friendly `title`:

```yaml
---
type: Requirement Evaluation Report
title: mutation endpoints are idempotent under retry
---
```

Good frontmatter fields:

- `type`;
- `title`.

The frontmatter `title` should name the report subject without repeating the
type prefix. Keep prefixes such as `Requirement:` or `Area:` in the visible H1,
where they help human scanning.

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

Avoid frontmatter fields for anything available from the linked JSON payloads or
visible Markdown header, especially:

- generated time;
- run identity;
- model snapshot;
- subject identity;
- scope;
- summaries;
- ratings;
- confidence;
- rating drivers;
- findings;
- recommendations;
- limits text;
- evidence text;
- rendered display labels.

This keeps the first lines of a report readable in editors that expose
frontmatter and keeps the structured Evaluation data as the single source of
truth. If a future consumer needs richer machine access, use
`data/evaluation-output-result.json` and the payloads it indexes instead of
expanding generated report frontmatter.

The visible H1 remains the first Markdown content after report frontmatter.

### Source Data at the bottom

Every generated report should end with a stable `## Source Data` section. The
section lists the structured Evaluation payloads used to render that report
artifact:

```markdown
## Source Data

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
list `data/evaluation-output-result.json` merely because it exists: that file is
a generated output index, not report source data unless a future renderer
directly consumes it.

## Header patterns

### Run report

The run report is the primary decision-ready report. Its header should make the
result and next navigation obvious.

```markdown
# Evaluation Report: LedgerLite Service

Run: 0001-full-eval - Generated: 2026-06-29 18:42 UTC - Model: [model-snapshot.md](model-snapshot.md)

Report: Overview - [Findings](findings.md) - [Recommendations](recommendations.md) - [Root Area](root-area.md)

| Overall Rating | Scope           | Confidence    |
| -------------- | --------------- | ------------- |
| Minimum        | full evaluation | Medium / None |

Summary:

LedgerLite is usable in the synthetic evaluation, but API idempotency, rollback rehearsal, and recovery ownership keep the overall service below target.

Jump to: [Top Findings](#top-findings) - [Top Recommendations](#top-recommendations) - [Area / Factor Breakdown](#area--factor-breakdown) - [Scope](#scope) - [Limits](#limits--incomplete-inputs)

## Area / Factor Breakdown

| Area / Factor                                                           | Overall Rating | Local Rating | Findings | Recommendations |
| ----------------------------------------------------------------------- | -------------- | ------------ | -------- | --------------- |
| **[LedgerLite Service](root-area.md)**                                  | Minimum        | —            | 7        | 3               |
| ↳ [Public API](areas/api/api-area.md)                                   | Minimum        | Minimum      | 2        | 1               |
| ↳ 🧩 [Correctness](areas/api/factors/correctness/correctness-factor.md) | Minimum        | Minimum      | 1        | 1               |
```

### Index reports

Finding and recommendation indexes should identify themselves as complete
indexes, not merely tables.

```markdown
# Findings

Run: 0001-full-eval - Generated: 2026-06-29 18:42 UTC - Scope: full evaluation

Report: [Overview](report.md) - Findings - [Recommendations](recommendations.md)

| Findings          | Highest Severity |
| ----------------- | ---------------- |
| 7 ranked findings | High             |
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

## Area / Factor Breakdown

| Area / Factor                                                 | Overall Rating | Local Rating | Findings | Recommendations |
| ------------------------------------------------------------- | -------------- | ------------ | -------- | --------------- |
| **[Public API](api-area.md)**                                 | Minimum        | Minimum      | 2        | 1               |
| ↳ 🧩 [Correctness](factors/correctness/correctness-factor.md) | Minimum        | Minimum      | 1        | 1               |
| ↳ 🧩 [Operability](factors/operability/operability-factor.md) | Target         | Target       | 1        | 0               |
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
```

### Requirement reports

Requirement reports are often the deep link target from findings. Their header
should make the owning Area, attached Factors, rating state, assessment state,
and confidence immediately visible.

```markdown
# Requirement: mutation endpoints are idempotent under retry

Run: [QEVAL-0001](../../../../report.md) - Created: 2026-06-29T18:42:00Z - Scope: full evaluation

Report: [Overview](../../../../report.md) - [Findings](../../../../findings.md) - [Recommendations](../../../../recommendations.md)

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [correctness](../../factors/correctness/correctness-factor.md)

| Rating  | Assessment | Confidence      |
| ------- | ---------- | --------------- |
| Minimum | Assessed   | Medium / Medium |

Jump to: [Findings Summary](#findings-summary) - [Finding Details](#finding-details) - [Unknowns & Missing Evidence](#unknowns--missing-evidence)
```

### Recommendation detail reports

Recommendation reports are action-planning surfaces. Their headers should make
rank, impact, confidence, trace context, and navigation visible before the
recommendation prose.

```markdown
# Recommendation: Tighten the idempotency replay contract

Run: [QEVAL-0001](../report.md) - Created: 2026-06-29T18:42:00Z - Scope: full evaluation

Report: [Overview](../report.md) - [Findings](../findings.md) - [Recommendations](../recommendations.md)

Trace: [Public API](../areas/api/api-area.md) / [Correctness](../areas/api/factors/correctness/correctness-factor.md)

| ID            | Rank | Impact | Confidence | Reference                                          |
| ------------- | ---- | ------ | ---------- | -------------------------------------------------- |
| QREC-0001-001 | 1    | High   | High       | evaluation:QEVAL-0001/recommendation:QREC-0001-001 |
```

## Checklist

Before changing report output, check:

- The first visible Markdown content is one clear H1.
- A reader can identify the run, generated time, scope, and subject in the
  header area.
- Report-level navigation links to the overview, findings, and recommendations
  where those reports exist.
- Hierarchical context is visible through Area, Factor, or Factors lines.
- The header uses a report-specific summary table.
- Long pages have useful jump links; short pages do not add noisy jump links.
- Frontmatter contains identity only.
- No generated report introduces claims that are absent from structured data.
- The bottom `Source Data` section lists the structured payloads used to render
  the report.
- Empty values render visibly and consistently.
