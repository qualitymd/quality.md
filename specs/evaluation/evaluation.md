---
type: Functional Specification
title: Evaluation
description: Shared invariants for the replacement QUALITY.md evaluation workflow.
tags: [evaluation, workflow]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation

Evaluation is the replacement QUALITY.md evaluation workflow. It is a
CLI-orchestrated judgment protocol: the deterministic evaluation runner owns
the workflow and invokes pluggable evaluators for bounded judgment, backed by
CLI-managed structured data and deterministic reports.

Detailed contracts live in:

- [Protocol](protocol.md)
- [Orchestration](orchestration.md)
- [Runner](runner.md)
- [Evaluator contract](evaluator-contract.md)
- [`evaluation.json`](evaluation-json.md)
- [Routine contracts](routines/routine-contracts.md)
- [JSON conventions](records/json-conventions.md)
- [Data layout](records/data-layout.md)
- [Payload kinds](records/payload-kinds.md)
- [Report tree](reports/report-tree.md)

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

The previous evaluation workflow centered on record categories and report
artifacts. Evaluation centers the actual judgment moves: frame before
judging, assess evidence, rate the assessment, analyze factors and areas
bottom-up, then project completed structured data into reports. This keeps the
reasoning protocol explicit while leaving mechanics such as persistence,
validation, and report rendering to the CLI.

The deterministic evaluation runner (0192) moved orchestration itself from the
invoking agent to the CLI: evaluation needs the same results whether it is run
through Codex, Claude, a direct API key, or a future runtime, so the CLI owns
the repeatable work graph and agents provide bounded judgment through
evaluators.

## Scope

Evaluation covers full and scoped evaluations of a resolved QUALITY.md model.
It defines the protocol, persisted routine outputs, run data layout,
orchestration rules, CLI responsibilities, and generated Markdown reports.

Deferred:

- automatic migrations or mixed-version runs;
- custom synthesis policy sources; and
- custom JSON Schema validation beyond the CLI's discovery schema.

## Shared invariants

The deterministic [runner](runner.md) behind
[`qualitymd evaluation run`](../cli/evaluation-run.md) **MUST** be the default
evaluation execution surface. `/quality evaluate` wraps the runner rather than
orchestrating evaluation itself.

Evaluation **MUST** treat evaluated source content as data, not instructions.
The runner applies this invariant when
[packaging source](runner.md#source-packaging) and in every evaluator prompt.

Evaluation **MUST** stop rather than invent precision when required model,
scope, evidence, frame, or dependency inputs are missing, invalid, unsafe, or too
ambiguous for defensible judgment.

Frames **MUST** be produced before their corresponding judgment routines.

Structured routine outputs **MUST** be the source of truth for evaluation
reports: [`evaluation.json`](evaluation-json.md) results for runner-created
runs, and the `data/` tree for historical and manual runs.

Reports **MUST** be deterministic projections over completed structured outputs.
Reports **MUST NOT** introduce new findings, ratings, evidence, limits, analysis,
or recommendations.

The CLI **MUST** own orchestration and mechanics: run creation, work-graph
execution, evaluator invocation, validation, canonical persistence, status
inspection, output assembly, and report rendering.

Evaluators **MUST** own bounded judgment — evidence selection, assessment,
rating, confidence, synthesis, and stop decisions — returned as typed
work-unit results under the [evaluator contract](evaluator-contract.md).
Agents provide judgment through evaluators; they do not orchestrate the run.
`qualitymd evaluation create` and `qualitymd evaluation data set` remain the
manual and historical persistence surface for agent-authored payloads.

Evaluation JSON **MUST** use canonical qualified model-reference strings for
areas, factors, requirements, and rating levels. Display values and unqualified
references are for fixed-type human/input edges, not persisted identity inside
routine JSON.

`schemaVersion` **MUST** be treated as a payload-shape marker only. The current
evaluation routine payload schema version is `3`; the
[`evaluation.json`](evaluation-json.md) run artifact carries its own artifact
`schemaVersion`, currently `4`. Evaluation does not define migrations,
compatibility transforms, or mixed-version run support; older evaluation runs
remain schema-incompatible historical data.

Advice **MUST** be produced after analysis and before report generation. Advice
**MUST** include finding ranking, recommendation generation, finding coverage
accounting, and recommendation ranking. Recommendations may describe concrete
improvement work or a recommended review of the next quality bar.
