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
Finding Core: statement, condition, criteria, cause, effect, and evidence.

Requirement assessment **MUST** classify findings by type using these semantics:
`gap` falls short of declared criteria, `risk` could plausibly cause future
quality loss, `strength` supports or exceeds criteria, `unknown` records missing
or ambiguous evidence, and `note` preserves relevant non-driving context.

Requirement assessment **MUST NOT** record `cause.status: verified` unless
finding evidence directly supports the cause statement. When a `gap` or `risk`
has enough evidence for condition and effect but not cause, the finding **MUST**
use `cause.status: not_assessed` rather than inventing cause.

Requirement assessment **MAY** record finding-local `candidateActions` —
non-binding remediation leads carried on a finding — and **MUST NOT** synthesize,
aggregate, deduplicate, or prioritize them across findings. Cross-finding
synthesis is recommendation generation, which Evaluation v0 forbids.

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

## Report Routine

`generateEvaluationReports` **MUST** be deterministic projection only.

The report routine **MUST NOT** create new findings, ratings, evidence, limits,
analysis, or recommendations.
