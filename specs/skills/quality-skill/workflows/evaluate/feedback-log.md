---
type: Functional Specification
title: Evaluate feedback log
description: Evaluate-specific adoption rules for the shared workflow feedback log.
tags: [skill, evaluation, logging, feedback]
timestamp: 2026-07-09T00:00:00Z
---

# Evaluate feedback log

Sub-spec of the [/quality evaluate](../evaluate.md) workflow spec. It defines how
`evaluate` adopts the shared
[workflow feedback log](../../workflow-feedback-log.md) artifact contract.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Background

Evaluation's historical `debug-log.md` was process-only, but it used a separate
name and sparse shape from setup's feedback log. Evaluate now records workflow
experience in the same feedback-log family as setup, while formal evaluation
judgment remains in assessment, analysis, recommendation, and report artifacts.

## Evaluate behavior

`evaluate` **MUST** create a feedback log after it verifies CLI support and emits
the run frame, before it invokes `qualitymd evaluation run`.

`evaluate` **MUST** write the log to
`.quality/logs/<timestamp>-evaluate-feedback-log.md`, using the evaluate run's
start timestamp and following the shared feedback-log contract.

The evaluate feedback log frontmatter **MAY** include these evaluate-specific
fields when known:

- `evaluation-run` — model-relative path to the numbered evaluation run.
- `scope` — concise human-readable scope label.
- `outcome` — workflow-process outcome when known, blank until close. Allowed
  terminal values are `completed-reportable`, `stopped-lint`, `stopped-model`,
  `stopped-source`, `stopped-tooling`, `failed`, and `interrupted`.

After `qualitymd evaluation run` reports a numbered run path, `evaluate`
**MUST** record that path in `evaluation-run` or in the timeline.

`evaluate` **MUST** update the current run's feedback log when material
workflow-experience information appears, including scope ambiguity, history
inspection friction, evaluator-selection friction, interruption or resume,
retries, tooling failures, slow phases, redaction decisions, UX/AX
observations, unusually smooth affordances worth preserving, and suggested
workflow improvements.

The log **SHOULD** avoid noisy churn for routine internal steps already captured
by CLI receipts, the runner's run-local logs, or generated reports.

At normal close, `evaluate` **MUST** set `status: completed`, set
`completed-at`, record `outcome: completed-reportable`, update effort when
available, and ensure each body section has useful content or an explicit
no-notable-content note.

When evaluation stops after the log exists because lint fails, source cannot be
resolved, requirements are not assessable, CLI support fails, user input is
needed, or another non-success condition occurs, the skill **SHOULD** finalize
the log with `status: failed` or `status: interrupted` when it can do so without
masking the stop condition, and **SHOULD** set the closest workflow outcome:
`stopped-lint`, `stopped-model`, `stopped-source`, `stopped-tooling`, `failed`,
or `interrupted`. If finalization is impossible, the existing
`status: in-progress` log remains acceptable partial feedback.

## Evaluation boundary

The evaluate feedback log **MUST NOT** be authoritative for ratings, findings,
recommendations, next actions, generated reports, or QUALITY.md model content.

The evaluate feedback log **MUST NOT** duplicate assessment evidence, rating
rationale, raw project-command output, or recommendation prose as primary
content.

The evaluate feedback log `outcome` field **MUST** describe how the evaluate
workflow ended. It **MUST NOT** describe a rating result, report verdict,
recommendation state, or evaluated-source quality.

When a project command is exercised as evaluation evidence, the feedback log
**MAY** record the routing decision and point to the formal assessment record
after that record exists.

Writing, updating, or finalizing the feedback log **MUST NOT** change evaluation
stop rules, rating semantics, record authority, reportability, or next-step
routing.
