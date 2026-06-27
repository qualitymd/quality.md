---
type: Design Doc
title: Recommendation Result Shape Design
description: Implementation approach for the simplified RecommendationResult fields and recommendation report projections.
tags: [evaluation, advice, recommendations, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Recommendation Result Shape Design

## Context

This design answers the [Recommendation Result Shape functional spec](spec.md).
The existing Advice implementation already validates `RecommendationResult`,
loads persisted Recommendation and RecommendationRanking payloads, and renders
Top Recommendations, `recommendations.md`, and recommendation detail pages. The
change reshapes the RecommendationResult field contract and tightens the report
projection boundary.

## Approach

### Data contract

Update the single data contract source in `internal/evaluation/data_contract.go`
so `RecommendationResult` requires:

```text
id
title
description
background
expectedValue
doneCriterion
impact
confidence
traceRefs
```

Because the data contract rejects unknown fields, removing
`whyItMatters`, `recommendedNextMove`, `expectedBenefit`, and
`howToKnowItWorked` is enough to make old-shape payloads invalid. No fallback
reader, alias, or migration path is added.

Update `recommendationExample` in `internal/evaluation/data.go` and the test
fixtures to the new field names.

### Recommendation row model

Add a render-ready recommendation row helper parallel to the ranked finding row
helper:

- rank;
- recommendation ID, title, and generated detail path;
- linked Area / Factors context resolved from `traceRefs`;
- expected value for the Top Recommendations `Reason` column;
- impact and confidence display labels;
- ranking rationale from the ranking entry; and
- source recommendation/ranking payloads for detail rendering.

The helper dereferences `RecommendationResult.traceRefs` through existing
requirement/finding reference utilities. Requirement references provide the
declaring Area and attached Factors. Finding references first resolve to their
Requirement, then use the same Area and Factor context. If multiple trace refs
resolve to the same Area / Factor pair, the table de-duplicates display text
while preserving deterministic order.

### Report projection

`report.md` renders Top Recommendations as:

```text
Rank | Recommendation | Area / Factors | Reason
```

`Reason` displays `RecommendationResult.expectedValue`.

`recommendations.md` uses the same recommendation row model but includes the
planning metadata appropriate for the full index:

```text
Rank | Recommendation | Area / Factors | Impact | Confidence | Reason | Ranking Rationale
```

Recommendation detail reports render rank, impact, confidence, description,
background, expected value, done criterion, trace refs, ranking rationale when
present, source data links, and the legend. They read only loaded structured
payloads and model/evaluation data already present in `evaluationArtifacts`;
they do not parse generated Markdown recommendation files.

## Spec Response

- **RecommendationResult fields (R1-R3)** - the typed data contract, examples,
  fixtures, durable specs, and skill guidance use only the new field family.
- **Report rendering (R4-R8, R10)** - one row helper feeds the run report,
  recommendation index, and detail pages, keeping expected value separate from
  ranking rationale.
- **Data-only source boundary (R9)** - report rendering continues to consume
  `evaluationArtifacts`, which are loaded from persisted structured data, model
  snapshot, and run metadata; no generated Markdown frontmatter is read.

## Alternatives

**Keep `whyItMatters` and rename only table labels.** Rejected because the data
model would keep carrying the same ambiguity even if the table appeared simpler.

**Use `reason` instead of `background` plus `expectedValue`.** Rejected because
one field would again mix why the recommendation arose with why acting on it is
valuable.

**Use `rationale` on RecommendationResult.** Rejected because ranking entries
already use `rationale` to explain ordering; reusing the same name would make
reports and agent instructions less precise.

**Read legacy recommendation Markdown frontmatter for older runs.** Rejected
because early-alpha compatibility is not required and generated reports must not
become hidden data sources for other generated reports.

## Trade-offs & Risks

- Existing historical runs using the old Advice field names become
  non-reportable under the current data contract. This is an intentional clean
  break.
- `Area / Factors` resolution depends on trace refs being specific enough to
  reach Requirement context. If a recommendation traces only to a broad analysis
  artifact, the report will render `—` for that cell rather than inventing
  context.
- The full recommendation index becomes wider, but `report.md` stays compact.

## Open Questions

None.
