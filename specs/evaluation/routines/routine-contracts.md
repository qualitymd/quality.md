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
NOT** assign a Rating Level.

Requirement assessment **MAY** record finding-local **candidate actions** —
non-binding remediation leads carried on a finding — and **MUST NOT** synthesize,
aggregate, deduplicate, or prioritize them across findings. Cross-finding
synthesis is recommendation generation, which Evaluation v0 forbids.

Requirement rating **MUST** produce `RequirementRatingResult` and **MUST NOT**
inspect new source evidence.

Factor and Area analysis **MUST** preserve lower-level drivers that prevent a
higher rating.

Area analysis **MUST** produce Area Findings in `AreaAnalysisResult.findings`
when it identifies material observations that should be visible at Area or
Factor report level. Area Findings **MUST** summarize observations for the
containing Area, cite their source routine payloads through non-empty
`inputRefs`, and relate only to Factors declared in that Area. Area analysis
**MUST NOT** record recommendation, impact, priority, effort, benefit, ROI, or
global ranking fields on Area Findings.

## Report Routine

`generateEvaluationReports` **MUST** be deterministic projection only.

The report routine **MUST NOT** create new findings, ratings, evidence, limits,
analysis, or recommendations.
