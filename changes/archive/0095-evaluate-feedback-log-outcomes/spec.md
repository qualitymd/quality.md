---
type: Functional Specification
title: Evaluate feedback log outcomes - functional spec
description: Requirements for keeping the evaluate feedback log as a process-only workflow artifact and making its outcome field use workflow-process outcomes.
tags: [skill, evaluation, logging, feedback]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluate feedback log outcomes - functional spec

Companion to
[Evaluate feedback log outcomes](../0095-evaluate-feedback-log-outcomes.md).
This spec states what the evaluate feedback-log contract must say. The format
itself is governed by [`SPECIFICATION.md`](../../../SPECIFICATION.md). This case
changes bundled skill guidance and durable skill spec mirrors only; it adds no
normative format rule and no CLI behavior.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The evaluate feedback log is the process-only companion to formal evaluation
artifacts. It improves the `/quality evaluate` workflow by recording execution
friction and recovery without turning those notes into evaluation evidence.
Because evaluation outputs already have authoritative homes, the feedback log's
metadata must avoid evaluation-sounding semantics. Its `outcome` field should say
how the workflow ended, not what rating, finding, recommendation, or report
result was produced.

## Requirements

### Keep the evaluate feedback log current and process-only

`/quality evaluate` **MUST** create and maintain a current-run feedback log under
`.quality/logs/<timestamp>-evaluate-feedback-log.md` after CLI support is
verified and the run frame is emitted.

The evaluate feedback log **MUST** record workflow-experience information only:
scope-resolution friction, history inspection, coverage adjustment, interruption
or resume, retries, record corrections, tooling failures, slow phases, redaction
decisions, prompt-injection handling, report generation recovery, UX/AX
observations, what worked well, and suggested workflow improvements.

The evaluate feedback log **MUST NOT** be authoritative for ratings, findings,
recommendations, next actions, generated reports, QUALITY.md model content,
evaluation evidence, or rating rationale.

### Make outcome a workflow-process value

The evaluate feedback log `outcome` field **MUST** describe how the evaluate
workflow ended. It **MUST NOT** describe a rating result, report verdict,
recommendation state, or evaluated-source quality.

Allowed terminal `outcome` values **MUST** be:

- `completed-reportable` - evaluation completed and report artifacts were built.
- `stopped-lint` - evaluation stopped because the model did not pass lint.
- `stopped-model` - evaluation stopped because model requirements, scope, or
  criteria were not usable enough to rate.
- `stopped-source` - evaluation stopped because evaluated source could not be
  resolved or safely inspected.
- `stopped-tooling` - evaluation stopped because required CLI/tooling support was
  missing, stale, or failed before a reportable run.
- `failed` - evaluation failed for another error after the feedback log existed.
- `interrupted` - evaluation was interrupted or user input was required before
  completion.

Runtime guidance **MUST** refer to these values as workflow outcomes, not report
outcomes.

> Rationale: `reported` and "report outcome" blur process status with evaluation
> judgment. The explicit values keep the feedback log useful for workflow
> improvement without creating a shadow report channel. - 0095

### Keep debug-log compatibility separate

Current evaluate feedback **MUST** remain under `.quality/logs/`. New evaluation
runs **MUST NOT** seed or require run-root `debug-log.md`. Historical
`debug-log.md` compatibility remains a reader concern only.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/workflows/evaluate/feedback-log.md` - define
  allowed workflow-process `outcome` values and replace loose close-outcome
  wording (per the workflow-outcome requirement above).
- `specs/skills/quality-skill/workflows/evaluate.md` - align runtime workflow
  finalization wording and frontmatter example with workflow-process outcomes
  (per the keep-current and workflow-outcome requirements above).
- `specs/skills/quality-skill/quality-skill.md` - align the parent workflow
  feedback-log summary if needed (per the keep-current requirement above).
- `specs/skills/quality-skill/workflows/log.md` and `specs/log.md` - record the
  durable spec updates.

### To rename

None

### To delete

None
