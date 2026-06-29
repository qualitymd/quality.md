---
type: Change Case
title: Recommendation IDs and Numbers
description: Split opaque recommendation identity from user-facing recommendation numbers so Recommendation #1 means the first ranked recommendation.
status: Done
tags: [evaluation, advice, recommendations, identifiers]
timestamp: 2026-06-29T00:00:00Z
---

# Recommendation IDs and Numbers

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0176-recommendation-ids-and-numbers/spec.md) - what the case must do.
- [Design doc](0176-recommendation-ids-and-numbers/design.md) - how it is built, and why.

## Motivation

Recommendation users naturally say "recommendation #1" to mean the first
recommendation in the ranked advice list. The current contract assigns
`RecommendationResult.number` before recommendation ranking, then also renders
ranked recommendation order. That gives humans and agents two small ordinal-ish
handles for the same recommendation, which makes "rec #1" ambiguous.

This case gives each recommendation one opaque artifact identity and reserves
the user-facing recommendation number for ranked order. The existing workflow
sequencing stays intact: recommendation payloads can still be persisted before
ranking, but they receive a non-ordinal `qrec_...` ID. Ranking then assigns the
number users see.

## Scope

Covered:

- Replace `RecommendationResult.number` with an opaque `RecommendationResult.id`
  using the `qrec_<token>` form.
- Store recommendation JSON under the opaque ID path.
- Reference recommendations by opaque ID in recommendation ranking and finding
  coverage data.
- Derive user-facing recommendation numbers from
  `RecommendationRankingResult.orderedRecommendations[].rank`.
- Label human report columns and follow-up prompts as recommendation numbers, not
  rank-first identity.
- Keep recommendation detail report filenames rank-numbered for the ranked
  advice list.

Deferred:

- A formal public command for resolving recommendation numbers or IDs outside
  recommendation follow-up.
- A migration command for pre-change evaluation runs.

Non-goals:

- Changing run identity, run numbering, finding identity, Model references,
  Requirement Finding selectors, or report anchors.
- Adding effort, ROI, quick-win, backlog-priority, priority, or score fields to
  recommendations.
- Preserving compatibility with pre-change `RecommendationResult.number` data.

## Affected artifacts

Derived by sweeping for `RecommendationResult.number`, `recommendationRef`,
`recommendationRefs`, `RecommendationRankingResult`, `findingCoverage`,
`recommendations/`, recommendation report tables, and recommendation follow-up
selection language.

**Code**

- [x] `internal/evaluation/data.go` - assign opaque `qrec_...` IDs, derive
      recommendation data paths from IDs, and validate ranking/coverage refs by
      ID.
- [x] `internal/evaluation/data_contract.go` - change recommendation result and
      ranking contracts from numeric refs to string IDs.
- [x] `internal/evaluation/report_tree.go` - render recommendation numbers from
      ranking order while keeping opaque IDs in source-data/detail context.
- [x] `internal/evaluation/evaluation-data.schema.json` - regenerate from the
      typed contract.
- [x] `internal/evaluation/evaluation_test.go` - update fixtures and assertions.
- [x] `scripts/report-gallery/main.go` and generated gallery output if report
      fixtures change.

**Durable specs**

- [x] `specs/evaluation/records/json-conventions.md` - define recommendation ID
      and recommendation number semantics.
- [x] `specs/evaluation/records/payload-kinds.md` - replace
      `RecommendationResult.number` with `id`; ranking and coverage refs use
      IDs.
- [x] `specs/evaluation/records/data-layout.md` - recommendation data path from
      opaque ID; detail report path from recommendation number.
- [x] `specs/evaluation/reports/report-tree.md` - report tables use
      recommendation number derived from ranking; opaque IDs are secondary.
- [x] `specs/evaluation/protocol.md`, `specs/evaluation/orchestration.md`, and
      `specs/evaluation/routines/routine-contracts.md` - keep advice sequencing
      and reference semantics aligned.
- [x] `specs/cli/evaluation-data.md` - data set assignment and validation use
      recommendation IDs.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` - skill workflow writes
      recommendation results before ranking but references IDs.

**Durable docs / runtime guidance**

- [x] `skills/quality/SKILL.md` - align Advice authoring guidance.
- [x] `skills/quality/workflows/evaluate.md` - align workflow runtime guidance.
- [x] `skills/quality/guides/recommendation-follow-up.md` - selection language
      treats recommendation numbers as ranked order and IDs as secondary.
- [x] `specs/skills/quality-skill/recommendation-follow-up.md` if durable
      follow-up wording needs the same distinction.
- [x] `CHANGELOG.md` - clean-break release note if the current release section
      tracks evaluation contract changes.
- [x] `docs/guides/reporting-design.md` - align recommendation detail path
      naming.
- [x] `examples/report-gallery/software-service/.quality/evaluations/**` -
      regenerate checked-in report gallery output.

## Status

`Done`. Implemented and archived. Evaluation recommendations now persist
opaque `qrec_...` IDs for JSON data paths and structured ranking/coverage refs;
generated reports derive user-facing recommendation numbers from ranked order;
runtime skill guidance, durable specs, schema, tests, changelog, and report
gallery output are aligned.
