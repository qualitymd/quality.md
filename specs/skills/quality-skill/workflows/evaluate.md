---
type: Functional Specification
title: /quality evaluate
description: Behavioral component spec for running evaluation through the /quality skill as a wrapper around the deterministic evaluation runner.
tags: [skill, quality, evaluate, evaluation, workflow]
timestamp: 2026-07-11T00:00:00Z
---

# /quality evaluate

`evaluate` is the `/quality` skill workflow that produces a current evaluation
result for a resolved QUALITY.md model scope by invoking the CLI-owned
deterministic runner, `qualitymd evaluation run`. It implements the shared
[evaluation wrapper](../evaluation.md), safety, and CLI-ownership contracts in
the parent [/quality skill](../quality-skill.md) spec.

The runtime procedure lives at
[`skills/quality/workflows/evaluate.md`](../../../../skills/quality/workflows/evaluate.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`evaluate` is selected when the user asks to evaluate quality, asks for a scoped
quality assessment, or names only a resolvable area or factor after the model is
present and valid.

The workflow's purpose is to produce a current evaluation report for the
resolved model file and scope, including required recommendations, by wrapping
`qualitymd evaluation run` with the agent-mediated interface: intent parsing,
the run frame, scope resolution, preflight validation, evaluator-selection
explanation, invocation, and result summary. It does not apply fixes.

## Mutation surface and artifacts

`evaluate` may mutate only the current run's evaluate feedback log under
`.quality/logs/` directly. Evaluation artifacts under the resolved evaluation
directory — the numbered run folder, `model-snapshot.md`, the authoritative
`evaluation.json`, run-local logs, and the generated Markdown report tree — are
created by `qualitymd evaluation run`, which the workflow invokes. `evaluate`
**MUST NOT** edit evaluated source files, edit `QUALITY.md`, write the quality
changelog, create external issues, apply recommendations, or update tooling.

`evaluate` **MUST** create evaluation runs only through
`qualitymd evaluation run`. It **MUST NOT** manually create run folders,
hand-author or edit structured run data, write report files directly, or use
`qualitymd evaluation create` and `qualitymd evaluation data set` for new
evaluations.

> Rationale: a wrapper that re-evaluates the source or writes its own run data
> recreates the two-engine architecture the deterministic runner removes. —
> 0192

`evaluate` **MUST** create, update, and finalize the current run's evaluate
feedback log as defined by the
[Evaluate feedback log](evaluate/feedback-log.md) sub-spec.
The feedback log's `outcome` field describes workflow process state such as
`completed-reportable`, `stopped-model`, or `interrupted`; it is not a rating,
report verdict, or recommendation state.

## Canonical references

`evaluate` **MUST** resolve a natural-label or unqualified scope to canonical
`area:`/`factor:` references before invoking the runner, using `qualitymd model`
introspection (`model list --json`, `model get`) against the resolved model. It
**MUST NOT** derive area or factor references from `QUALITY.md` text, and it
**MUST** pass resolved references to the runner through `--area` and repeatable
`--factor`.

## Required flow

Before tool inspection, `evaluate` **MUST** emit the public `/quality` run frame
required by the parent skill contract as its first output, before workspace
resolution or any other tool call. When the requested scope is not yet resolved,
the frame **MUST** render a provisional scope value (such as `resolving…`) rather
than block on resolution, and `evaluate` **SHOULD** confirm the resolved scope in
a later message.

Before invoking the runner, `evaluate` **MUST**:

- verify compatible CLI support, including the evaluation runner surface;
- resolve the model file and requested scope to canonical references;
- run lint and stop on errors;
- inspect relevant evaluation history when present;
- open the evaluate feedback log; and
- explain evaluator selection.

For evaluator selection, `evaluate` **MUST** resolve and explain the transport
per the shared [selection precedence](../evaluation.md#evaluator-selection):
an explicit user request, then a non-`auto` configured
`evaluation.evaluator`, then `--evaluator harness` when the current agent can
service harness checkpoints, then CLI `auto` discovery. It **MAY** preview the
resolved model, scope, evaluator, and work-unit counts with
`qualitymd evaluation run --dry-run --json`, and **MAY** ask the user to choose
an evaluator when the CLI reports a missing or ambiguous evaluator, presenting
the CLI's remedies as the options. It **MUST NOT** silently cross to a
different provider after harness selection or failure.

When selection fails because an agent runtime or requested capability is
missing, `evaluate` **MUST** surface the runner's concrete installation,
authentication, configuration, or selector remedy. It **MUST NOT** silently
fall back after a run has accepted results.

`evaluate` **MUST** then invoke the runner with explicit flags:

```text
qualitymd evaluation run [--model <model>] [--area <area-ref>]
  [--factor <factor-ref>...] [--evaluator <name>] --json
```

While the receipt status is `awaiting_evaluator`, `evaluate` **MUST** service
each of the receipt's outstanding requests only from that request's
runner-supplied immutable area context — it **MAY** use the harness's ordinary
reasoning or delegate independent requests to native subagents, and **MAY**
submit results as they become ready rather than waiting for the whole set —
submit one typed result envelope per request (one or several per call) with
`qualitymd evaluation run --resume <run> --evaluator-result - --json`, and
repeat until a terminal receipt; each resume returns the window topped up
with newly-ready requests. In unattended automation the loop **MUST NOT**
add interactive gates: the run advances, returns a report, or stops with the
runner's classified remedy.

While the runner executes and after it returns, `evaluate` **MUST NOT**
independently collect evidence, run a parallel QC loop, second-guess the
runner's authoritative result, judge anything beyond a checkpoint's supplied
request, or write structured evaluation data. The evaluation protocol —
coverage, verification, roll-up, advice, and report generation — is
runner-owned.

A `resolveSource` checkpoint is the exception to judgment's no-exploration
rule: it may use bounded read-only workspace tools to identify the selector's
files, then returns only the finite workspace-relative file set. Requirement
judgment never receives that exploration transcript and never reads outside the
captured area context.

## Failure and resume

An `awaiting_evaluator` receipt is expected progress, not a failure. If the
checkpoint loop is interrupted, `evaluate` **MUST** recover by resuming the run
without a result to re-obtain the same outstanding requests; requests never
submitted stay outstanding at no retry cost.

When the runner reports `failed` or `cancelled`, `evaluate` **MUST** explain the
receipt's stable failure category in user terms and offer
`qualitymd evaluation run --resume <run>` as the recovery path when the run is
resumable. It **MUST NOT** combine `--resume` with an `--evaluator` that differs
from the run's recorded evaluator; re-evaluating with a different evaluator is a
new run.

## Progress

`evaluate` **SHOULD** re-emit a short, factual progress beat at phase boundaries
where the user's mental model would otherwise drift — not only in the opening
frame. It **MUST** include a beat before invoking the runner (the first
mutation), and **SHOULD** also show progress after scope resolution and at
closeout. The runner writes its own progress diagnostics to stderr; `evaluate`
summarizes rather than duplicates them. Progress output **MUST** remain factual
and user-facing, not a transcript of internal reasoning.

## Stop conditions

`evaluate` **MUST** stop before invoking the runner when:

- the in-scope area source cannot be resolved;
- the in-scope model has no requirements;
- required CLI support is missing or stale; or
- lint reports structural errors.

It **SHOULD** stop before invoking the runner when requirements are too vague to
bind evidence to a rating. A stop response **MUST** distinguish model
usefulness, evaluated-source quality, and evaluation-history status. It **MUST**
keep the blocking reason and best next step scannable, offer concrete runnable
options when available, and include an explicit answer path such as replying
with an option number or saying `stop`.

## Completion criteria

`evaluate` is complete when the runner's receipt reports a completed run and the
user-facing summary states the rating, scope, the `report.md` path, top findings
and recommendations drawn from the generated reports, known limits from the
receipt and reports, changed artifacts, what was not done, and the recommended
next action, and the feedback log is finalized.
The closeout **MUST** use labeled fields for rating, scope, evidence basis, known
limitations, changed artifacts, not-done boundary, report-reading CTA, and next
action so the result, artifacts, limits, and next step are visible in a
five-second scan.
