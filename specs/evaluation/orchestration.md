---
type: Functional Specification
title: Evaluation orchestration
description: Runner-owned work graph, scheduling, persistence, resume, retry, and cancellation rules.
tags: [evaluation, orchestration, runner]
timestamp: 2026-07-09T00:00:00Z
---

# Evaluation orchestration

Evaluation is a dependency-ordered work graph. The deterministic
[evaluation runner](runner.md) owns the graph: it builds it from the model
snapshot and planned scope, schedules ready work units, applies the retry
policy, persists accepted results, and decides resume. This spec defines the
graph and its runtime invariants.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Work graph

The runner **MUST** build a deterministic work graph in model order:

1. `frameEvaluation`;
2. `frameAreaEvaluation` for each planned area;
3. `resolveSource` for each planned area whose pinned selector kind has no
   deterministic resolver (see the
   [runner's source packaging contract](runner.md#source-packaging));
4. `frameRequirementEvaluation`, then `assessRateRequirement` for each planned
   requirement;
5. `frameFactorAnalysis` and `analyzeFactor` for each in-scope factor node,
   bottom-up;
6. `frameAreaAnalysis` and `analyzeArea` for each in-scope area, bottom-up;
7. `rankFindings`;
8. `recommend`;
9. `rankRecommendations`; and
10. `buildReports`.

Work-unit IDs **MUST** be deterministic strings: `<kind>` for run-wide units
and `<kind>:<canonical-ref>` for subject-scoped units, for example
`assessRateRequirement:requirement:docs::links-resolve`.

Frame units and `buildReports` **MUST** be deterministic runner work; the
remaining units are evaluator-backed work dispatched under the
[evaluator contract](evaluator-contract.md). `resolveSource` is
evaluator-backed gathering, not judgment: its accepted effect is the captured
source bundle persisted in the artifact's
[`sources` record](evaluation-json.md#sources), never a result payload.

`assessRateRequirement` **MUST** execute the protocol's `assessRequirement` and
`rateRequirement` moves as one evaluator call and persist both the
`RequirementAssessmentResult` and `RequirementRatingResult` payloads, under
their unchanged kinds, identities, and schemas per
[payload kinds](records/payload-kinds.md#judgment-kinds).

> Rationale: a separate rate call re-ships the requirement's full context to
> score an assessment the same evaluator just produced, roughly doubling
> requirement calls for no new evidence. The merge is a call-shape change only;
> reports, roll-up, and resume cannot tell the difference. — 0193

The graph executes every [protocol](protocol.md) move: the protocol's
`assessRequirement` and `rateRequirement` run inside `assessRateRequirement`,
`accountForFindingCoverage` runs as runner acceptance validation of the
`rankRecommendations` result, and `assembleEvaluationOutputResult` plus
`generateEvaluationReports` run inside `buildReports`.

## Dependencies

The runner **MUST** enforce these dependencies even when work is delegated to
subagents, workers, threads, processes, or another runtime-specific mechanism:

`frameEvaluation` **MUST** complete before any area work begins.

An area's `frameAreaEvaluation` **MUST** complete before local requirement
work, local factor work, or child area work for that area begins.

An area's `resolveSource` unit, when present, **MUST** complete — its bundle
captured and persisted — before any of the area's `assessRateRequirement`
units are dispatched. Analysis and advice units consume prior results, not
source, and take no dependency on it.

`RequirementRatingResult`s **MUST** exist before a factor node that depends on
those direct requirements can be analyzed. The combined
`assessRateRequirement` unit satisfies this dependency: its accepted result
carries the requirement's valid rating.

Direct child `FactorAnalysisResult`s **MUST** exist before a parent factor is
analyzed.

Root `FactorAnalysisResult`s and direct child `AreaAnalysisResult`s **MUST**
exist before an area is analyzed.

All required in-scope area, factor, requirement assessment, requirement rating,
factor analysis, area analysis, and advice outputs **MUST** exist and be valid
before `EvaluationOutputResult` and reports are generated.

`FindingRankingResult` **MUST** exist before recommendation generation closes.

All `RecommendationResult` payloads **MUST** exist before final finding
coverage accounting and recommendation ranking close. Finding coverage and
recommendation ranking **MUST** reference recommendations by
`RecommendationResult.id`.

`RecommendationRankingResult` **MUST** exist before `EvaluationOutputResult`
and reports are generated.

Report generation **MUST** require a valid run manifest, a valid evaluation
frame, the analysis outputs required by the manifest's `plannedScope`, and the
required advice outputs. It **MUST NOT** select report scope from evaluation
frame ordering.

## Scheduling and parallelism

The runner **MAY** execute dependency-ready evaluator-backed work units
concurrently, up to the resolved [concurrency](runner.md#concurrency) cap.
Concurrency is a scheduling choice under the runner contract; it never becomes
an alternate orchestration engine.

Parallel execution **MUST** be observationally equivalent to deterministic
sequential execution in model order.

Parallel execution **MUST NOT** change ratings, report content, output
ordering, artifact paths, or persisted payload shapes.

Evaluator workers **MUST NOT** write evaluation run artifacts directly. Accepted
results reach disk only through the runner-owned persistence path, and persisted
payload order **MUST** remain deterministic graph order even when evaluator calls
finish out of order.

## Harness checkpoints

When the selected evaluator is harness-backed, a ready evaluator work unit
**MUST** checkpoint instead of dispatching a subprocess: the runner persists
the awaiting state atomically, returns the bounded work request, and schedules
nothing further until a correlated result is submitted. Deterministic units
**MUST** continue to execute on each invocation up to the next evaluator
checkpoint or the terminal receipt, so deterministic work never leaks into the
agent interface.

Harness-backed execution **MUST** run with resolved concurrency `1` until the
runner defines multiple pending evaluator checkpoints.

A submitted result **MUST** advance the graph only when it correlates with the
persisted pending request (request identity and input hash). A mismatched,
duplicate, or unsolicited result **MUST** be rejected with `run_state_invalid`
and leave the pending request recoverable. A pending request whose rebuilt
input hash no longer matches its checkpoint **MUST** fail with
`run_state_invalid` rather than bind judgment to changed evidence.

A schema-invalid or unparseable harness submission **MUST** consume the same
[retry budget](#retry-and-failure) as any other evaluator attempt: the runner
records the attempt, re-checkpoints with the next attempt's request, and fails
the run when the budget is exhausted.

## Persistence

The runner **MUST** merge each accepted work-unit result into
`evaluation.json` and persist it before treating the work unit as complete for
scheduling or resume. Write mechanics are the
[`evaluation.json` contract](evaluation-json.md#persistence).

> Rationale: an interrupted run must resume without repeating accepted
> judgment work; batched persistence silently discards completed evaluator
> calls. — 0192

Evaluators and subagents **MUST NOT** write files in an evaluation run folder;
results reach disk only through the runner, per the
[evaluator contract](evaluator-contract.md#boundaries).

## Resume

Before scheduling work, a resumed run **MUST** verify artifact compatibility
per the [`evaluation.json` contract](evaluation-json.md#resume-compatibility),
then rebuild the graph from the current model snapshot and compare it with
saved state.

A completed work unit **MUST** be reused when its recorded input hash matches
the recomputed hash of its resolved inputs.

A work unit **MUST** be rerun when its required output is missing, malformed,
schema-incompatible, or dependency-stale (its input hash no longer matches).

## Retry and failure

The runner **MUST** retry a failed work unit only when its failure category is
`rate_limited`, `timeout`, `invalid_evaluator_output`, or
`schema_invalid_output`, up to three attempts total per work unit. Any other
failure category **MUST** fail the run immediately.

An `assessRateRequirement` result that carries an assessment but no valid
rating, or the reverse, **MUST** fail the work unit as retryable evaluator
output; the runner **MUST NOT** persist a partial requirement result.

> Rationale: the combined call introduces the half-answered result as a new
> failure mode; persisting an unrated requirement would trip roll-up later and
> break the rating dependency above. — 0193

If a work unit exhausts its attempts or fails with a non-retryable category,
then the run **MUST** finish with status `failed` and remain resumable.

Failure categories are the [runner failure taxonomy](runner.md#failure-taxonomy).

## Cancellation

When a run is interrupted by user cancellation or a termination signal (SIGINT
or SIGTERM), the runner **MUST** cancel in-flight evaluator calls, leave
`evaluation.json` valid and resumable, record the interruption in run state and
the event log, and report the run as `cancelled` rather than failed.
Interrupted work units keep their attempt counts and stay incomplete.

> Rationale: stopping a long evaluation is an expected user action, not an
> internal error; accepted work must survive for `--resume`. — 0192
