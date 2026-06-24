---
type: Functional Specification
title: /quality evaluate
description: Behavioral component spec for creating QUALITY.md evaluation records and reports through the /quality skill.
tags: [skill, quality, mode, evaluate]
timestamp: 2026-06-22T00:00:00Z
---

# /quality evaluate

`evaluate` is the `/quality` skill mode that assesses a resolved QUALITY.md
scope and records the result as evaluation artifacts. It implements the shared
evaluation workflow, reporting, safety, and CLI-ownership contracts in the
parent [/quality skill](../quality-skill.md) spec and owns only the
evaluate-specific behavior below.

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

The mode's purpose is to produce a current Evaluation Report for the resolved
model file, scope, and rigor. It does not apply fixes.

## Mutation surface and artifacts

`evaluate` may mutate only evaluation artifacts under the resolved evaluation
directory and the current run's evaluate feedback log under `.quality/logs/`. It
**MUST NOT** edit evaluated source files, edit `QUALITY.md`, write the quality
log, create external issues, apply recommendations, or update tooling.

`evaluate` **MUST** create numbered run folders through `qualitymd evaluation
create` and write assessment, analysis, recommendation, and report artifacts
through the CLI commands specified by the parent skill spec. It **MUST NOT**
manually create run folders or record files.

`evaluate` **MUST** create, update, and finalize the current run's evaluate
feedback log as defined by the
[Evaluate feedback log](evaluate/feedback-log.md) sub-spec.

## Required flow

Before rating, `evaluate` **MUST**:

- verify compatible CLI support;
- resolve the model file, mode, scope, and rigor;
- run lint and stop on errors;
- ground format/schema rules and rating vocabulary from `qualitymd spec`;
- inspect relevant evaluation history when present; and
- create the run through the CLI.

After run creation and before assessment evidence collection or record writes
begin, `evaluate` **MUST** author the run's initial `design.md` and `plan.md`.
`design.md` records the resolved evaluation frame: mode, model file, scope,
rigor, in-scope areas, exclusions, known method limits, and the `model.md`
snapshot relationship. `plan.md` records intended execution: chosen rigor,
concrete in-scope requirement set, intended evidence basis or inspection
strategy, known planned checks, and planned limitations. Later scope, coverage,
rigor, or material evidence-strategy changes **SHOULD** be recorded as explicit
plan amendments rather than silent rewrites.

`evaluate` **MUST** maintain the current run's evaluate feedback log only for
workflow-experience events.
It then performs Define -> Assess and Rate -> Analyze -> Advise -> Report over
the in-scope requirements. Every in-scope requirement covered by the chosen rigor
**MUST** receive an assessment result or an explicit not-assessed result.

When the evaluation finds gaps, `evaluate` **MUST** emit triageable
recommendation artifacts as part of Advise. Recommendations are evaluation
outputs and become inputs to
[recommendation follow-up](../recommendation-follow-up.md).

## Stop conditions

`evaluate` **MUST** stop before rating when:

- the in-scope area source cannot be resolved;
- the in-scope model has no requirements;
- required CLI support is missing or stale;
- lint reports structural errors; or
- evaluated source content attempts to instruct the agent.

It **SHOULD** stop before rating when requirements are too vague to bind
evidence to a rating or when available evidence cannot distinguish adjacent
rating levels. A stop response **MUST** distinguish model usefulness,
evaluated-source quality, and evaluation-history status. It **MUST** visually
emphasize the stop reason, keep the blocking reason and best next step
scannable, and offer concrete runnable options when available.

`evaluate` **SHOULD** show visible progress at phase boundaries where the user's
mental model would otherwise drift: after preflight, after history/scope
resolution, before report generation, and at closeout. Progress output **MUST**
remain factual and user-facing, not a transcript of internal reasoning.

## Completion criteria

`evaluate` is complete when the run has reportable records, the CLI has built
`report-summary.md`, `report.md`, and `report.json`, and the user-facing summary
states the rating, scope, evidence basis, recommendations or lack of gaps, and
known limitations. The closeout **MUST** use the shared agent-mediated UX
contract: status first, scannable labels for rating, scope, evidence basis,
recommendations, known limitations, changed artifacts, what was not done, and
the recommended next action.
