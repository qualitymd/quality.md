---
type: Functional Specification
title: Evaluation debug log - functional spec
description: Requirements for adding a process-only debug log to evaluation runs.
tags: [evaluation, records, skill]
timestamp: 2026-06-21T00:00:00Z
---

# Evaluation debug log - functional spec

Companion to the [Evaluation debug log](../0046-evaluation-debug-log.md). This
spec defines the required behavior for adding `debug-log.md` to evaluation runs.
The [design doc](design.md) covers the implementation approach.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [BCP 14](../../docs/reference/rfc2119.md) when, and only when, they
appear in all capitals.

## Background

Evaluation runs already have formal judgment artifacts for subject quality:
assessment records, analysis records, recommendation records, and generated
reports. Those artifacts should not carry every process detail the evaluator
needed to reach a reportable run. A separate `debug-log.md` preserves
diagnostics for orchestration, recovery, ambiguity, coverage, and tooling
friction while preventing process notes from becoming a shadow assessment system.

## Scope

This change covers a hand-authored, run-root `debug-log.md` artifact seeded by
the CLI and maintained by the `/quality` skill during evaluation.

Deferred:

- A dedicated CLI append/list command for debug events.
- Structured parsing, validation, or report rendering from `debug-log.md`.
- Any change to rating semantics or report authority.

## Requirements

### Run artifact

New evaluation runs MUST contain a run-root `debug-log.md` file seeded by
`qualitymd evaluation create`.

`debug-log.md` MUST be a runtime evaluation artifact, not an OKF `log.md`, and
MUST NOT be interpreted as an OKF concept history.

`debug-log.md` SHOULD begin with a short purpose statement, a boundary statement,
and an `## Events` section suitable for append-oriented process notes.

### Authority boundary

`debug-log.md` MUST NOT be authoritative for ratings, findings,
recommendations, next actions, or generated report content.

Formal subject-quality evidence MUST remain in assessment records. Roll-up
judgment MUST remain in analysis records. User-facing remediation advice MUST
remain in recommendation records. Generated reports MUST continue to derive from
the formal evaluation record model, not from debug-log prose.

### Process-only content

Debug-log entries MUST describe events involving the evaluation process itself:
scope resolution, evaluation history inspection, coverage planning, interruptions
or resumes, retries, artifact corrections, CLI/tooling readiness and failures,
subagent coordination, redaction decisions, prompt-injection handling, and report
generation recovery.

Debug-log entries MUST NOT record subject-quality findings as primary content or
duplicate rating rationale from assessment or analysis records.

When a project command is exercised as evidence against the evaluated subject,
the debug log MAY record the routing fact and cite the formal assessment record,
but it MUST NOT copy the command's raw output or use that output as process-log
evidence.

> Rationale: The command boundary is the slippery case. A lint command can be
> either evaluator mechanics or subject evidence depending on intent; the log
> should explain routing decisions without preserving a second copy of the
> assessment evidence.

### Redaction

`debug-log.md` MUST NOT contain secret values.

`debug-log.md` MUST NOT reproduce raw prompt-injection text from evaluated source
content.

When a secret or hostile instruction affects the evaluation process, the log MAY
record a sanitized locator and type.

### Skill behavior

The `/quality` skill SHOULD hand-author `debug-log.md` entries for notable
evaluation-process events during `evaluate` and the evaluation phase of
`improve`.

The skill MUST keep debug-log prose separate from assessment, analysis,
recommendation, and report content.

The skill MAY omit debug-log entries for ordinary happy-path steps that are
already captured by `design.md`, `plan.md`, record receipts, or generated
reports.

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation-records.md` - add `debug-log.md` to the run-folder contract
  and specify its authority boundary, process-only content, and redaction rules.
- `specs/cli/evaluation-create.md` - require `qualitymd evaluation create` to
  seed `debug-log.md`.
- `specs/skills/quality-skill/quality-skill.md` - require the skill to preserve
  the debug-log boundary while hand-authoring notable evaluation-process events.

### To delete

None.
