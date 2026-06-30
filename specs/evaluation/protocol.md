---
type: Functional Specification
title: Evaluation protocol
description: Phase ordering, traversal, and stop behavior for Evaluation.
tags: [evaluation, protocol, workflow]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation protocol

This spec defines the Evaluation protocol order. It inherits shared
invariants from [Evaluation](evaluation.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Protocol order

An Evaluation run **MUST** perform these protocol moves:

1. `frameEvaluation`
2. `frameAreaEvaluation` for each in-scope Area
3. `frameRequirementEvaluation` for each local Requirement
4. `assessRequirement` for each framed Requirement
5. `rateRequirement` for each Requirement Assessment
6. `frameFactorAnalysis` and `analyzeFactor` for each Factor node bottom-up
7. `frameAreaAnalysis` and `analyzeArea` for each Area bottom-up
8. `rankFindings`
9. `recommend`
10. `accountForFindingCoverage`
11. `rankRecommendations`
12. `assembleEvaluationOutputResult`
13. `generateEvaluationReports`

The protocol **MUST NOT** require a specific execution engine. Sequential
execution and parallel worker execution are both valid when they satisfy the same
dependency order and produce the same persisted outputs.

## Area traversal

The evaluator **MUST** walk the Area tree bottom-up for analysis.

For each Area, the evaluator **MUST**:

1. create an Area Evaluation Frame;
2. evaluate local Requirements;
3. analyze the Area's local Factor forest bottom-up;
4. evaluate child Areas;
5. create an Area Analysis Frame after root Factor and child Area analyses are
   complete; and
6. produce the Area Analysis Result.

Root Area `localAndDescendantAnalysis`, when analyzed, is the overall evaluation
result.

## Planned scope

Evaluation output assembly **MUST** use `EvaluationManifest.plannedScope` as the run's
authoritative scope. It **MUST NOT** select a headline subject from
agent-authored frame input ordering.

The planned expansion is derived from `plannedScope` and the model snapshot: the
planned Area, its descendant Areas, the filtered Factors when `factorFilter` is
non-empty or all in-scope Factors when it is empty, and the Requirements reached
from those in-scope Areas and Factors.

Report generation **MUST** fail when required analysis data for the planned
scope is missing. Coverage checks **MUST** compare the planned expansion against
the structured analysis artifacts actually produced.

## Factor traversal

The evaluator **MUST** walk each Area's Factor tree bottom-up.

For each Factor node, the evaluator **MUST** analyze child Factors before the
parent Factor.

A Factor Analysis Frame **MUST** include direct Requirement Rating Results for
that exact Factor node and direct child Factor Analysis Results. It **MUST NOT**
include transitive descendant Factor refs because each child result already
accounts for its descendants.

When a Factor has no direct Requirements (an umbrella Factor), its Factor
Analysis Result **MUST** record `localAnalysis` with the `empty` status and a
reason, because there is no local signal to analyze; the Factor's
`localAndDescendantAnalysis` carries the roll-up of its child Factor analyses.

## Requirement flow

A Requirement **MUST** be framed before evidence assessment.

A Requirement **MUST** be assessed before it is rated.

Requirement rating **MUST** map the Requirement Assessment Result to the
pre-framed applied Rating Level criteria. It **MUST NOT** inspect new evidence or
change the criteria after evidence is observed.

## Stop conditions and limits

The protocol **MUST** distinguish stop conditions from evaluation limits.

A stop condition blocks a routine from producing ordinary judgment output. An
evaluation limit constrains what the routine can honestly claim while still
allowing the routine to continue when the limit is non-blocking.

When a routine cannot continue, it **SHOULD** produce no persisted output or a
valid structured output with the appropriate non-complete status, such as
`blocked`, `not_assessed`, `not_rated`, or `not_analyzed`.

## Synthesis defaults

Factor synthesis v0 **MUST** use
`protocol:factor-synthesis-default-v0`.

Area synthesis v0 **MUST** use `protocol:area-synthesis-default-v0`.

Both defaults **MUST** apply `worst_bound`: the synthesized level is constrained
by the lowest rating-relevant input unless explicit rationale justifies an
override.

Both defaults **MUST** preserve binding drivers and surface incomplete inputs.
They **MUST NOT** synthesize new findings; roll-up explanation belongs in
`ratingDrivers`, rationale, confidence, limits, and incomplete inputs.

## Advice flow

Advice **MUST** run after Area analysis has completed and before output assembly
or report generation.

`rankFindings` **MUST** consider every Requirement Finding produced in the
effective run data and write one `FindingRankingResult`. The ranking's tiers and
order express relative priority; no Requirement Finding is omitted from the
ranking, so a low-priority finding is ranked at a low tier rather than dropped.
If the evaluation produced no findings, the ranking result **MUST** record an
empty ranking.

`recommend` **MUST** produce at least one `RecommendationResult`. When the
evaluation found no concrete improvement work, the recommendation **MUST** be
quality-management advice, such as reviewing whether the current quality bar
should be raised or clarified.

`accountForFindingCoverage` **MUST** verify every finding is represented in
`RecommendationRankingResult.findingCoverage` as either addressed by one or more
recommendations or explicitly not advice-driving. Coverage accounting happens
after recommendation generation and before recommendation ranking is closed.
Coverage recommendation refs **MUST** use `RecommendationResult.id`.

`rankRecommendations` **MUST** rank persisted recommendations by expected
quality impact, quality-bar relevance, trace strength, confidence, and whether
the advice addresses binding constraints. Ranking recommendation refs **MUST**
use `RecommendationResult.id`. It **MUST NOT** use effort, ROI,
quick-win status, backlog priority, or numeric score fields as required
Evaluation data.
