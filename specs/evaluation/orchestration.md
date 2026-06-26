---
type: Functional Specification
title: Evaluation orchestration
description: Agent-agnostic dependency graph, parallelism, persistence, resume, and failure rules.
tags: [evaluation, orchestration, agents]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation orchestration

Evaluation is a dependency-ordered work graph. This spec defines the graph and
runtime invariants without requiring a specific agent or concurrency mechanism.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Work Units

An orchestrator **MUST** treat these as the primary work-unit boundaries:

- `EvaluationWork`
- `AreaWork`
- `RequirementWork`
- `FactorWork`
- `ReportWork`

The orchestrator **MUST** enforce dependencies between work units even when it
delegates work to subagents, workers, threads, processes, queues, or another
runtime-specific mechanism.

## Dependencies

`EvaluationFrame` **MUST** exist before root `AreaWork` begins.

`AreaEvaluationFrame` **MUST** exist before local `RequirementWork`, local
Factor work, or child Area work for that Area begins.

`RequirementRatingResult`s **MUST** exist before a Factor node that depends on
those direct Requirements can be analyzed.

Direct child `FactorAnalysisResult`s **MUST** exist before a parent Factor is
analyzed.

Root `FactorAnalysisResult`s and direct child `AreaAnalysisResult`s **MUST**
exist before an Area is analyzed.

All required Area, Factor, Requirement Assessment, and Requirement Rating outputs
**MUST** exist and be valid before `EvaluationOutputResult` and reports are
generated.

## Parallelism

A runtime **MAY** execute ready work units concurrently.

Parallel execution **MUST** be observationally equivalent to deterministic
sequential execution in model order.

Parallel execution **MUST NOT** change ratings, report content, output ordering,
data paths, or persisted payload shapes.

Good v0 parallelism includes sibling Requirements, child Areas, and sibling
Factors whose dependencies are ready.

## Persistence

Workers **SHOULD** return structured routine JSON to the orchestrator. The
orchestrator should assemble accepted payloads into a JSON array and persist the
batch through:

```text
qualitymd evaluation data set <run> < payloads.json
```

Workers **SHOULD NOT** write arbitrary files in an evaluation run folder.

If a runtime allows workers to call the CLI directly, the resulting run **MUST**
be equivalent to orchestrator-mediated persistence.

## Resume

Before scheduling a work unit, the orchestrator **MAY** inspect persisted
outputs.

A work unit **MAY** be skipped when the expected output exists, is structurally
valid, its dependencies have not changed, and the runtime accepts reuse for the
current run.

A work unit **MUST** be rerun when required output is missing, malformed,
schema-incompatible, stale relative to dependencies, or not provably reusable.

## Failure

A failed work unit **SHOULD** produce either no persisted output or a valid
structured output whose status reflects the failure mode.

The orchestrator **SHOULD** continue independent work where possible, then rely
on `evaluation status` or failed `evaluation report build` to surface typed
gaps.
