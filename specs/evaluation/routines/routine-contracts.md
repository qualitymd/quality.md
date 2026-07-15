---
type: Functional Specification
title: Evaluation routine contracts
description: Prompt-contract rules for agent-run evaluation routines.
tags: [evaluation, routines, prompts]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation routine contracts

Evaluation routines are prompt contracts for an agent-run workflow. They are
not ordinary CLI functions, although the CLI validates and persists their
outputs.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Prompt contract shape

Each agent-run routine prompt **SHOULD** define:

- role;
- task;
- inputs;
- required output;
- constraints;
- stop rules; and
- self-check.

The prompt contract **MUST** make clear what inputs the agent may use, what
output shape is required, and when the agent must stop rather than inventing
precision.

## Routine set

Evaluation **MUST** define prompt contracts for:

- `frameEvaluation`
- `frameAreaEvaluation`
- `frameRequirementEvaluation`
- `assessRequirement`
- `rateRequirement`
- `frameFactorAnalysis`
- `analyzeFactor`
- `frameAreaAnalysis`
- `analyzeArea`
- `rankFindings`
- `recommend`
- `accountForFindingCoverage`
- `rankRecommendations`
- `summarizeEvaluation`
- `assembleEvaluationOutputResult`
- `generateEvaluationReports`

Helper routines such as `prepareAreaContext`, `validateEvaluationFrame`,
`validateAreaContext`, `frameAreaRequirementEvaluations`,
`assessAreaRequirements`, `rateRequirementAssessments`,
`indexRequirementRatingsByFactor`, `evaluateAreaFactorForest`, and
`evaluateFactorNode` **MAY** be documented in runtime instructions, but their
outputs **MUST NOT** become required persisted payloads unless a later spec adds
them.

## Framing routines

Framing routines **MUST** produce structured frame payloads before judgment
routines run.

Framing routines **MUST NOT** inspect evidence to decide findings, assign rating
levels, synthesize factor or area ratings, or write report prose.

## Judgment routines

Judgment routines **MUST** use only their frames and declared inputs.

Requirement assessment **MUST** produce `RequirementAssessmentResult` and **MUST
NOT** assign a rating level. Each requirement finding **MUST** use the shared
finding core: type, confidence, statement, condition, criteria, basis, effect,
and evidence. `gap` and `risk` findings **MUST** include severity; `strength`
and `note` findings **MUST NOT** include severity.

Requirement assessment **MUST** classify findings by type using these semantics:
`gap` falls short of declared criteria, `risk` could plausibly cause future
quality loss, `strength` supports or exceeds criteria, and `note` preserves
relevant non-driving context. Ambiguous current-state evidence that constrains a
rating is a `gap`. Missing or insufficient evidence that prevents rating belongs
in assessment/rating status, `unknowns`, or `missingEvidence`, not in a finding
type.

Requirement assessment **MUST NOT** record `basis.status: verified` unless
finding evidence directly supports the basis statement. When a `gap` or `risk`
has enough evidence for condition and effect but not basis, the finding **MUST**
use `basis.status: not_assessed` rather than inventing basis.

For a `strength`, `basis.status: verified` means the positive condition's basis
is directly supported by cited evidence; `basis.status: not_applicable` means no
separate basis beyond the cited evidence is claimed.

Requirement assessment **MAY** record finding-local `candidateActions` —
non-binding remediation leads carried on a finding — and **MUST NOT** synthesize,
aggregate, deduplicate, or prioritize them across findings. Cross-finding
synthesis belongs in the advice routines.

Requirement rating **MUST** produce `RequirementRatingResult` and **MUST NOT**
inspect new source evidence.

Requirement rating **MUST NOT** produce `status: rated` unless the paired
requirement assessment has one or more requirement findings sufficient to justify
the selected configured rating level. Requirement rating **MUST** remain
scale-agnostic and **MUST NOT** assume fixed rating level meanings such as
target, sub-target, pass, or fail.

Rated requirement results **MUST** include non-empty `ratingDrivers` that cite
the paired requirement assessment through `inputRefs`; drivers **SHOULD** select
the specific requirement findings that drove the rating.

Factor and area analysis **MUST** preserve lower-level drivers that prevent a
higher rating.

Factor and area analysis **MUST NOT** produce findings. Rated factor and area
analysis scopes **MUST** include non-empty `ratingDrivers` that cite lower-level
requirement rating results, factor analysis results, or area analysis results
through `inputRefs`. Analysis drivers **MUST NOT** introduce new evidence or
claims absent from the referenced lower-level outputs.

## Advice routines

Advice routines **MUST** run after area analysis and before report generation.

`rankFindings` **MUST** produce `FindingRankingResult` from persisted
requirement findings. Ranking **MUST** use judgment about quality-bar relevance,
concern severity when present, binding effect on ratings, confidence, affected
scope, and whether the finding changes next quality-management action.

`recommend` **MUST** produce one or more `RecommendationResult` payloads. A
recommendation **MUST** be quality-domain agnostic: it should name the quality
improvement or review move without assuming software, product, operations, or
any other modeled domain unless that domain is present in the evaluated model.

`accountForFindingCoverage` **MUST** ensure every persisted finding is accounted
for before recommendation ranking closes. A finding may be covered by one or
more recommendations, referenced by `RecommendationResult.id`, or marked
`not_advice_driving` with rationale.

`rankRecommendations` **MUST** produce `RecommendationRankingResult` after
coverage accounting. Ranking **MUST** use judgment about expected quality impact,
quality-bar relevance, trace strength, confidence, and relationship to binding
constraints. Ranking entries **MUST** reference recommendations by
`RecommendationResult.id`. It **MUST NOT** require effort, ROI, quick-win status,
backlog priority, or numeric score fields.

`summarizeEvaluation` **MUST** produce one `EvaluationSummaryResult` after
recommendation ranking. It **MUST** synthesize only the scoped root analysis and
persisted ranked advice supplied to it, use the model's rating-scale labels, and
write concise stakeholder-facing prose that names the concrete areas and risks
driving the overall standing. It **MUST NOT** inspect new evidence, introduce
unsupported claims or counts, or use roll-up mechanism terms as reader-facing
explanation.

Recommendations **MAY** be concrete work or recommended review. When the
evaluation meets the current bar and has no gap/risk requiring work, the advice
phase **MUST** still recommend whether to review, raise, clarify, or confirm the
next quality bar.

## Report routine

`generateEvaluationReports` **MUST** be deterministic projection only.

The report routine **MUST NOT** create new findings, ratings, evidence, limits,
analysis, or recommendations.
