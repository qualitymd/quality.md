---
type: Functional Specification
title: Evaluation v2 protocol
description: Phase ordering, traversal, and stop behavior for Evaluation v2.
tags: [evaluation, protocol, workflow]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation v2 protocol

This spec defines the Evaluation v2 protocol order. It inherits shared
invariants from [Evaluation v2](evaluation-v2.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Protocol Order

An Evaluation v2 run **MUST** perform these protocol moves:

1. `frameEvaluation`
2. `frameAreaEvaluation` for each in-scope Area
3. `frameRequirementEvaluation` for each local Requirement
4. `assessRequirement` for each framed Requirement
5. `rateRequirement` for each Requirement Assessment
6. `frameFactorAnalysis` and `analyzeFactor` for each Factor node bottom-up
7. `frameAreaAnalysis` and `analyzeArea` for each Area bottom-up
8. `assembleEvaluationOutputResult`
9. `generateEvaluationReports`

The protocol **MUST NOT** require a specific execution engine. Sequential
execution and parallel worker execution are both valid when they satisfy the same
dependency order and produce the same persisted outputs.

## Area Traversal

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

## Factor Traversal

The evaluator **MUST** walk each Area's Factor tree bottom-up.

For each Factor node, the evaluator **MUST** analyze child Factors before the
parent Factor.

A Factor Analysis Frame **MUST** include direct Requirement Rating Results for
that exact Factor node and direct child Factor Analysis Results. It **MUST NOT**
include transitive descendant Factor refs because each child result already
accounts for its descendants.

## Requirement Flow

A Requirement **MUST** be framed before evidence assessment.

A Requirement **MUST** be assessed before it is rated.

Requirement rating **MUST** map the Requirement Assessment Result to the
pre-framed applied Rating Level criteria. It **MUST NOT** inspect new evidence or
change the criteria after evidence is observed.

## Stop Conditions And Limits

The protocol **MUST** distinguish stop conditions from evaluation limits.

A stop condition blocks a routine from producing ordinary judgment output. An
evaluation limit constrains what the routine can honestly claim while still
allowing the routine to continue when the limit is non-blocking.

When a routine cannot continue, it **SHOULD** produce no persisted output or a
valid structured output with the appropriate non-complete status, such as
`blocked`, `not_assessed`, `not_rated`, or `not_analyzed`.

## Synthesis Defaults

Factor synthesis v0 **MUST** use
`protocol:factor-synthesis-default-v0`.

Area synthesis v0 **MUST** use `protocol:area-synthesis-default-v0`.

Both defaults **MUST** apply `worst_bound`: the synthesized level is constrained
by the lowest rating-relevant input unless explicit rationale justifies an
override.

Both defaults **MUST** preserve binding drivers and surface incomplete inputs.
