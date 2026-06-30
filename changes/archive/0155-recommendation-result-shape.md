---
type: Change Case
title: Recommendation Result Shape
description: Simplify RecommendationResult fields and make recommendation reports render from persisted Advice data only.
status: Done
tags: [evaluation, advice, recommendations, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Recommendation Result Shape

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0155-recommendation-result-shape/spec.md) - what the case must do.
- [Design doc](0155-recommendation-result-shape/design.md) - how it is built, and why.

## Motivation

The current `RecommendationResult` prose fields split recommendation intent
across `whyItMatters`, `recommendedNextMove`, `expectedBenefit`, and
`howToKnowItWorked`. In practice, `whyItMatters` and `expectedBenefit` overlap,
and the table/report surface needs a clearer distinction between what to do,
why the recommendation arose, what value acting should create, and how completion
is verified.

Reports must also stay deterministic projections of persisted evaluation data.
Generated recommendation tables and detail pages should never read YAML
frontmatter or prose from another generated report to fill recommendation cells.

## Scope

Covered:

- Replace `RecommendationResult` user-facing prose fields with `description`,
  `background`, `expectedValue`, and `doneCriterion`.
- Keep `id`, `title`, `impact`, `confidence`, and `traceRefs`.
- Remove the old field names without compatibility aliases, fallback readers, or
  dual writers.
- Render `report.md` Top Recommendations with `Rank`, `Recommendation`,
  `Area / Factors`, and `Reason`, where `Reason` displays `expectedValue`.
- Render `recommendations.md` and recommendation detail reports from persisted
  `RecommendationResult` and `RecommendationRankingResult` data only.
- Preserve ranking rationale as ranking data, not as the recommendation's own
  reason.

Deferred:

- Adding new scoring, effort, ROI, quick-win, backlog priority, or filtering
  fields.
- Adding interactive sorting or grouping to `recommendations.md`.
- Migrating historical evaluation runs.

## Affected artifacts

Derived by sweeping for `RecommendationResult`, `whyItMatters`,
`recommendedNextMove`, `expectedBenefit`, `howToKnowItWorked`,
`Top Recommendations`, and recommendation report rendering.

**Code**

- [x] `internal/evaluation/data_contract.go` - update `RecommendationResult`
      validation/schema fields.
- [x] `internal/evaluation/data.go` - update examples and payload semantics.
- [x] `internal/evaluation/report_tree.go` - update recommendation tables and
      detail reports from persisted data fields.
- [x] `internal/evaluation/evaluation_test.go` - update fixtures and report
      assertions.

**Durable specs**

- [x] `specs/evaluation/records/payload-kinds.md` - update
      `RecommendationResult` field contract.
- [x] `specs/evaluation/reports/report-tree.md` - update Top Recommendations,
      recommendation index, and detail report contracts.
- [x] `specs/cli/evaluation-report.md` - update the generated report contract
      summary.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` - update Advice
      authoring field names.

**Durable docs / runtime guidance**

- [x] `skills/quality/workflows/evaluate.md` - update runtime Advice authoring
      field names if present.
- [x] `skills/quality/SKILL.md` - update runtime Advice authoring field names if
      present.
- [x] `CHANGELOG.md` - release-note entry for the model/report shape change.
- [x] `changes/index.md` and `changes/log.md` - Change Case lifecycle.

## Status

`Done`. Implementation is complete across the Advice data contract, report
rendering, durable specs, runtime guidance, release notes, and tests. `mise run
check` passes.
