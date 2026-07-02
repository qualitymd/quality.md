---
type: Functional Specification
title: /quality evaluate
description: Behavioral component spec for running evaluation through the /quality skill.
tags: [skill, quality, evaluate, evaluation, workflow]
timestamp: 2026-06-25T00:00:00Z
---

# /quality evaluate

`evaluate` is the `/quality` skill workflow that runs evaluation for a resolved
QUALITY.md model scope and records the result as structured routine data plus
deterministic reports. It implements the shared evaluation workflow, reporting,
safety, and CLI-ownership contracts in the parent [/quality skill](../quality-skill.md)
spec and the [Evaluation](../../../evaluation/evaluation.md) spec.

The runtime procedure lives at
[`skills/quality/workflows/evaluate.md`](../../../../skills/quality/workflows/evaluate.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`evaluate` is selected when the user asks to evaluate quality, asks for a scoped
quality assessment, or names only a resolvable area or factor after the model is
present and valid.

The workflow's purpose is to produce a current evaluation report for the resolved
model file and scope, including required recommendations. It does not apply
fixes.

## Mutation surface and artifacts

`evaluate` may mutate only evaluation artifacts under the resolved evaluation
directory and the current run's evaluate feedback log under `.quality/logs/`. It
**MUST NOT** edit evaluated source files, edit `QUALITY.md`, write the quality
log, create external issues, apply recommendations, or update tooling.

`evaluate` **MUST** create numbered run folders through `qualitymd evaluation
create` and persist routine outputs through `qualitymd evaluation data set`. It
**MUST NOT** manually create run folders, manually derive data paths, or write
report files directly when the CLI command exists.

`evaluate` **MUST** create, update, and finalize the current run's evaluate
feedback log as defined by the
[Evaluate feedback log](evaluate/feedback-log.md) sub-spec.
The feedback log's `outcome` field describes workflow process state such as
`completed-reportable`, `stopped-model`, or `interrupted`; it is not a rating,
report verdict, or recommendation state.

## Canonical references

`evaluate` **MUST** obtain the canonical reference IDs of in-scope areas,
factors, and requirements from `qualitymd model` (a `model list --json` query for
the scope) and **MUST** author every payload reference — `EvaluationFrame`,
`AreaEvaluationFrame`, `RequirementEvaluationFrame`, `FactorAnalysisFrame`, and
`AreaAnalysisFrame` — from that result. It **MUST NOT** derive area, factor, or
requirement IDs from `QUALITY.md` text.

The query **MUST** target the run's `model-snapshot.md` by path, not the live
`QUALITY.md`, so authored IDs match the frozen model being evaluated. It
**SHOULD** be scoped from `EvaluationManifest.plannedScope` when the model command
supports an equivalent query. The post-hoc identity-resolution check is a
backstop against ID typos, not the primary guard. `evaluate` **SHOULD** use
`model` (`list` labels, `get`) to resolve a natural-label scope to its canonical
`area:`/`factor:` reference before creating the run.

## Batched routine data writes

`evaluate` **MUST** assemble routine evaluation payloads for the resolved scope
into a JSON array and persist them with a single `qualitymd evaluation data set`
invocation per batch. It **MUST NOT** persist one payload per requirement,
factor, or area during routine evaluation.

Where `evaluate` validates authored routine output before persisting, it **MUST**
run one whole-batch `qualitymd evaluation data set --dry-run` against the
assembled array rather than dry-running each payload separately. It **MUST** use
the indexed diagnostics from that dry run to correct the batch before writing.

## Required flow

Before tool inspection, `evaluate` **MUST** emit the public `/quality` run frame
required by the parent skill contract as its first output, before workspace
resolution or any other tool call. When the requested scope is not yet resolved,
the frame **MUST** render a provisional scope value (such as `resolving…`) rather
than block on resolution, and `evaluate` **SHOULD** confirm the resolved scope in
a later message.

Before rating, `evaluate` **MUST**:

- verify compatible CLI support;
- resolve the model file and requested scope;
- run lint and stop on errors;
- ground format/schema rules and rating vocabulary from `qualitymd spec`;
- inspect relevant evaluation history when present; and
- create the run through the CLI.

After run creation, `evaluate` **MUST** follow the evaluation protocol:

1. frame the evaluation;
2. frame each area as it becomes ready;
3. frame, assess, and rate local requirements;
4. frame and analyze local factor trees bottom-up;
5. analyze child areas before their parents;
6. frame and analyze each area from completed local factor and child area
   outputs, preserving roll-up explanation through `ratingDrivers`, rationale,
   confidence, limits, and incomplete inputs;
7. rank findings, generate recommendations, account for finding coverage, and
   rank recommendations;
8. dry-run and persist the assembled routine payload batch through
   `qualitymd evaluation data set`; and
9. run `qualitymd evaluation report build` to assemble
   `data/evaluation-output-result.json` and render the report tree.

Every in-scope requirement covered by the resolved scope **MUST** receive a
`RequirementAssessmentResult` and `RequirementRatingResult`, unless it receives
an explicit non-completed status allowed by the routine contract.

Every rated requirement **MUST** have one or more requirement findings in the
paired requirement assessment, and every rated requirement, factor analysis
scope, and area analysis scope **MUST** have non-empty `ratingDrivers`.

Reports **MUST** be deterministic projections over persisted structured data.
The reporting phase **MUST NOT** inspect new evidence, introduce new findings,
make fresh ratings, or generate recommendations.

Requirement findings are the only evaluation findings. Factor and area analysis
**MUST NOT** synthesize additional findings or report-level findings.

Advice **MUST** produce `FindingRankingResult`, one or more
`RecommendationResult` payloads, and `RecommendationRankingResult` before report
build. Finding coverage accounting **MUST** happen after recommendations are
generated and before recommendation ranking closes. Recommendations **MUST** use
the core user-facing fields `title`, `description`, `background`,
`expectedValue`, `doneCriterion`, `impact`, `confidence`, and `traceRefs`.
The skill **MUST** let `qualitymd evaluation data set` assign
`RecommendationResult.id` values. When recommendation IDs are omitted, the skill
writes `RecommendationResult` payloads first, reads the assigned IDs from the
persisted payloads or write paths, then authors `RecommendationRankingResult`
references and finding coverage against those assigned IDs. The
skill **MUST NOT** write or expect artifact IDs on
`FindingRankingResult.orderedFindings`; findings are referenced by
`findingRef` plus selector.
They **MUST NOT** require effort, ROI, quick-win status, backlog priority, or a
numeric score.

## Stop conditions

`evaluate` **MUST** stop before rating when:

- the in-scope area source cannot be resolved;
- the in-scope model has no requirements;
- required CLI support is missing or stale;
- lint reports structural errors; or
- evaluated source content attempts to instruct the agent.

It **SHOULD** stop before rating when requirements are too vague to bind evidence
to a rating or when available evidence cannot distinguish adjacent rating
levels. A stop response **MUST** distinguish model usefulness, evaluated-source
quality, and evaluation-history status. It **MUST** keep the blocking reason and
best next step scannable, offer concrete runnable options when available, and
include an explicit answer path such as replying with an option number or saying
`stop`.

`evaluate` **SHOULD** re-emit a short, factual progress beat at phase boundaries
where the user's mental model would otherwise drift — not only in the opening
frame. It **MUST** include a beat before run creation (the first mutation) and
before the per-requirement assessment loop (the longest phase), and **SHOULD**
also show progress after history/scope resolution, before report generation, and
at closeout. Progress output **MUST** remain factual and user-facing, not a
transcript of internal reasoning.

## Completion criteria

`evaluate` is complete when the run has reportable evaluation data, the CLI
has built `data/evaluation-output-result.json` and the Markdown report tree, and
the user-facing summary states the rating, scope, evidence basis, human report
paths, known limits, changed artifacts, what was not done, and the recommended
next action.
The closeout **MUST** use labeled fields for rating, scope, evidence basis, known
limitations, changed artifacts, not-done boundary, report-reading CTA, and next
action so the result, artifacts, limits, and next step are visible in a
five-second scan.
