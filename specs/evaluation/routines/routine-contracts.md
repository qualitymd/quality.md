---
type: Functional Specification
title: Evaluation routine contracts
description: Prompt-contract requirements for agent-run Evaluation routines.
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

## Prompt Contract Shape

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

## Routine Set

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
- `assembleEvaluationOutputResult`
- `generateEvaluationReports`

Helper routines such as `prepareAreaContext`, `validateEvaluationFrame`,
`validateAreaContext`, `frameAreaRequirementEvaluations`,
`assessAreaRequirements`, `rateRequirementAssessments`,
`indexRequirementRatingsByFactor`, `evaluateAreaFactorForest`, and
`evaluateFactorNode` **MAY** be documented in runtime instructions, but their
outputs **MUST NOT** become required persisted payloads unless a later spec adds
them.

## Framing Routines

Framing routines **MUST** produce structured frame payloads before judgment
routines run.

Framing routines **MUST NOT** inspect evidence to decide findings, assign Rating
Levels, synthesize Factor or Area ratings, or write report prose.

## Judgment Routines

Judgment routines **MUST** use only their frames and declared inputs.

Requirement assessment **MUST** produce `RequirementAssessmentResult` and **MUST
NOT** assign a Rating Level. Each Requirement Finding **MUST** use the shared
Finding Core: statement, condition, criteria, basis, effect, and evidence.

Requirement assessment **MUST** classify findings by type using these semantics:
`gap` falls short of declared criteria, `risk` could plausibly cause future
quality loss, `strength` supports or exceeds criteria, and `note` preserves
relevant non-driving context. Ambiguous current-state evidence that constrains a
rating is a `gap`. Missing or insufficient evidence that prevents rating belongs
in assessment/rating status, `unknowns`, or `missingEvidence`, not in a Finding
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
synthesis belongs in the Advice routines.

Requirement rating **MUST** produce `RequirementRatingResult` and **MUST NOT**
inspect new source evidence.

Requirement rating **MUST NOT** produce `status: rated` unless the paired
Requirement Assessment has one or more Requirement Findings sufficient to justify
the selected configured Rating Level. Requirement rating **MUST** remain
scale-agnostic and **MUST NOT** assume fixed Rating Level meanings such as
target, sub-target, pass, or fail.

Rated Requirement results **MUST** include non-empty `ratingDrivers` that cite
the paired Requirement Assessment through `inputRefs`; drivers **SHOULD** select
the specific Requirement Findings that drove the rating.

Factor and Area analysis **MUST** preserve lower-level drivers that prevent a
higher rating.

Factor and Area analysis **MUST NOT** produce findings. Rated Factor and Area
analysis scopes **MUST** include non-empty `ratingDrivers` that cite lower-level
Requirement Rating Results, Factor Analysis Results, or Area Analysis Results
through `inputRefs`. Analysis drivers **MUST NOT** introduce new evidence or
claims absent from the referenced lower-level outputs.

## Advice Routines

Advice routines **MUST** run after Area analysis and before report generation.

`rankFindings` **MUST** produce `FindingRankingResult` from persisted
Requirement Findings. Ranking **MUST** use judgment about quality-bar relevance,
finding severity, binding effect on ratings, confidence, affected scope, and
whether the finding changes next quality-management action.

`recommend` **MUST** produce one or more `RecommendationResult` payloads. A
recommendation **MUST** be quality-domain agnostic: it should name the quality
improvement or review move without assuming software, product, operations, or
any other modeled domain unless that domain is present in the evaluated Model.

`accountForFindingCoverage` **MUST** ensure every persisted Finding is accounted
for before recommendation ranking closes. A finding may be covered by one or
more recommendations, referenced by `RecommendationResult.id`, or marked
`not_advice_driving` with rationale.

`rankRecommendations` **MUST** produce `RecommendationRankingResult` after
coverage accounting. Ranking **MUST** use judgment about expected quality impact,
quality-bar relevance, trace strength, confidence, and relationship to binding
constraints. Ranking entries **MUST** reference recommendations by
`RecommendationResult.id`. It **MUST NOT** require effort, ROI, quick-win status,
backlog priority, or numeric score fields.

Recommendations **MAY** be concrete work or recommended review. When the
evaluation meets the current bar and has no gap/risk requiring work, the Advice
phase **MUST** still recommend whether to review, raise, clarify, or confirm the
next quality bar.

## Report Routine

`generateEvaluationReports` **MUST** be deterministic projection only.

The report routine **MUST NOT** create new findings, ratings, evidence, limits,
analysis, or recommendations.
