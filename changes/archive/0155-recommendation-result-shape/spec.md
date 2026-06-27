---
type: Functional Specification
title: Recommendation Result Shape
description: Requirements for simplifying RecommendationResult fields and rendering recommendation reports from persisted data.
tags: [evaluation, advice, recommendations, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Recommendation Result Shape

Companion to the [Recommendation Result Shape](../0155-recommendation-result-shape.md)
Change Case. This spec states *what* the change must do. It defers to
[Evaluation payload kinds](../../specs/evaluation/records/payload-kinds.md),
[Evaluation report tree](../../specs/evaluation/reports/report-tree.md),
[qualitymd evaluation report](../../specs/cli/evaluation-report.md), and
[Evaluation Advice](../archive/0150-evaluation-advice/spec.md) as normative
context for the existing Advice and report contracts.

The key words "MUST", "MUST NOT", and "MAY" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The existing recommendation prose fields over-separate action, reason, expected
benefit, and verification while still letting reason and benefit blur together.
The report entrypoint also needs a concise recommendation table that explains the
quality-management reason to act without duplicating every detail field.

This change makes the RecommendationResult core read as an action card:
headline, what to do, evidence-grounded background, expected value, done
criterion, metadata, and trace. It also makes the report projection boundary
explicit: generated reports render persisted Advice data and do not consume
frontmatter or prose from other generated Markdown reports.

## Scope

Covered: `RecommendationResult` fields, validation/schema/examples, report
tables, recommendation detail reports, durable specs, skill guidance, and tests.

Deferred / non-goals: historical run migration, compatibility aliases, fallback
readers for old recommendation fields, sorting/filtering UI, and new planning
metadata such as effort, ROI, quick-win status, backlog priority, or numeric
scores.

## Requirements

### RecommendationResult fields

R1. `RecommendationResult` **MUST** require `id`, `title`, `description`,
`background`, `expectedValue`, `doneCriterion`, `impact`, `confidence`, and
non-empty `traceRefs`.

> Rationale: these fields distinguish headline, action, origin, value, done
> state, metadata, and trace without forcing a reader to infer which "why" a
> field carries.
>
> Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
> `specs/skills/quality-skill/workflows/evaluate.md`.

R2. `RecommendationResult` **MUST NOT** accept `whyItMatters`,
`recommendedNextMove`, `expectedBenefit`, or `howToKnowItWorked`.

> Rationale: QUALITY.md is early alpha; keeping both field families would make
> the Advice contract harder to author and report from.
>
> Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
> `specs/skills/quality-skill/workflows/evaluate.md`.

R3. `RecommendationResult.expectedValue` **MUST** state the
quality-management value expected from completing the recommendation, and
**MUST NOT** require or imply effort, ROI, quick-win status, backlog priority, or
numeric score fields.

> Durable spec: modify `specs/evaluation/records/payload-kinds.md`.

### Report rendering

R4. The run-level `report.md` Top Recommendations table **MUST** render the
columns `Rank`, `Recommendation`, `Area / Factors`, and `Reason`, in that order.
`Reason` **MUST** render `RecommendationResult.expectedValue`.

> Rationale: the run entrypoint should stay compact while giving the reader the
> value-oriented reason to consider each recommendation.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` and
> `specs/cli/evaluation-report.md`.

R5. In Top Recommendations and `recommendations.md`, the `Recommendation` cell
**MUST** use `RecommendationResult.title` as link text and **MUST** link to the
generated recommendation detail report.

> Durable spec: modify `specs/evaluation/reports/report-tree.md`.

R6. In Top Recommendations and `recommendations.md`, the `Area / Factors` cell
**MUST** render linked Area and Factor names resolved from
`RecommendationResult.traceRefs` through persisted evaluation data and the model
snapshot. When no Area or Factor can be resolved, it **MUST** render `—`.

> Rationale: recommendation context should be navigable like ranked findings
> context and should not depend on prose parsing.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md`.

R7. `recommendations.md` **MUST** render a complete ranked recommendation index
from persisted `RecommendationResult` payloads and
`RecommendationRankingResult`, including all ranked recommendations, Area /
Factors, impact, confidence, ranking rationale, Reason, links to recommendation
detail reports, and a coverage summary from `findingCoverage`.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` and
> `specs/cli/evaluation-report.md`.

R8. Each recommendation detail report **MUST** render `title`, rank when ranked,
impact, confidence, `description`, `background`, `expectedValue`,
`doneCriterion`, trace refs, ranking rationale when ranked, and source data
links.

> Durable spec: modify `specs/evaluation/reports/report-tree.md`.

R9. Generated recommendation reports **MUST** render only persisted
`RecommendationResult`, `RecommendationRankingResult`, model snapshot, and
referenced evaluation data. They **MUST NOT** read YAML frontmatter or Markdown
body content from another generated report as recommendation source data.

> Rationale: generated reports are deterministic projections of structured
> evaluation data; report-to-report scraping would make stale or hand-edited
> Markdown a hidden data source.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` and
> `specs/cli/evaluation-report.md`.

R10. Recommendation ranking rationale **MUST** remain sourced from
`RecommendationRankingResult.orderedRecommendations[].rationale` and **MUST NOT**
be conflated with `RecommendationResult.background` or `expectedValue`.

> Durable spec: modify `specs/evaluation/reports/report-tree.md`.

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation/records/payload-kinds.md` - update `RecommendationResult`
  field names and expected value semantics (R1-R3).
- `specs/evaluation/reports/report-tree.md` - update Top Recommendations,
  recommendation index, detail reports, Area / Factors links, data-only report
  source boundary, and ranking rationale separation (R4-R10).
- `specs/cli/evaluation-report.md` - update `evaluation report build` report
  projection summary (R4, R7, R9).
- `specs/skills/quality-skill/workflows/evaluate.md` - update Advice authoring
  field names (R1-R2).

### To rename

None.

### To delete

None.
