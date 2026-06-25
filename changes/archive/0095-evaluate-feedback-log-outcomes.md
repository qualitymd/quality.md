---
type: Change Case
title: Evaluate feedback log outcomes
description: Reaffirm the evaluate feedback log as the process-only workflow artifact for /quality evaluate and make its outcome field use workflow-process outcomes instead of report or rating semantics.
status: Done
tags: [skill, evaluation, logging, feedback]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluate feedback log outcomes

A **Change Case** to keep the evaluate feedback log in the `/quality evaluate`
workflow and tighten its semantics. The feedback log belongs under
`.quality/logs/<timestamp>-evaluate-feedback-log.md` and records the experience
of running evaluation — scope friction, retries, redaction, tooling failures,
recovery, and UX/AX observations — not ratings, findings, evidence, or
recommendations.

Detail lives in:

- [Functional spec](0095-evaluate-feedback-log-outcomes/spec.md) - what the
  evaluate feedback-log contract must say.
- [Design doc](0095-evaluate-feedback-log-outcomes/design.md) - why the log stays
  outside evaluation runs and how outcome values are framed.

## Motivation

Evaluation is long-running, judgment-heavy, and tool-dependent. It needs a
process feedback artifact more than setup does: scope resolution, history
inspection, record writes, report generation, redaction, and interruption all
surface workflow lessons that should not be buried in user-facing closeout or
formal evaluation records.

The existing evaluate feedback-log guidance already has the right location and
boundary, but its `outcome` wording is loose: runtime guidance says to record
"the report outcome," and the durable sub-spec allows values like `reported`.
That phrasing can make a workflow-process field sound like it carries report,
rating, or evaluation semantics. The field should instead describe how the
workflow ended.

## Scope

Covered:

- Keep the evaluate feedback log as a current required workflow artifact under
  `.quality/logs/<timestamp>-evaluate-feedback-log.md`.
- Clarify that `outcome` is a workflow-process outcome, not report, rating,
  finding, or recommendation state.
- Define allowed evaluate `outcome` values for common terminal states.
- Align runtime workflow guidance, durable evaluate feedback-log spec, parent
  skill spec/guidance, and logs/changelog.

Deferred / non-goals:

- Do not move feedback logging back into evaluation run folders or revive
  `debug-log.md` for new runs.
- Do not make the CLI create, validate, or transmit feedback logs.
- Do not change Evaluation v2 data schemas, report generation, rating semantics,
  recommendation behavior, or Go code.

## Affected artifacts

Derived by analysis: searched for evaluate feedback-log contracts, `outcome`
wording, `.quality/logs`, and `debug-log.md` compatibility across `skills/`,
`specs/`, and release notes. Grouped by kind; empty kinds are deliberate.
Checkboxes are reconciled before In-Review.

### Code

None expected - bundled skill guidance and durable spec mirrors only.

### Format spec (`SPECIFICATION.md`)

None - no QUALITY.md format, rating, roll-up, or evaluation-semantics change.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0095-evaluate-feedback-log-outcomes/spec.md#durable-spec-changes)
section is authoritative. In summary:

- [x] `specs/skills/quality-skill/workflows/evaluate/feedback-log.md` - define
      workflow-process outcome values and reinforce the evaluation boundary.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` - align evaluate
      workflow creation/finalization wording with workflow outcomes.
- [x] `specs/skills/quality-skill/quality-skill.md` - align parent workflow
      feedback-log summary if needed.
- [x] `specs/skills/quality-skill/workflows/log.md` and `specs/log.md` - record
      durable spec updates.

### Durable docs

- [x] `CHANGELOG.md` - add a user-facing `/quality` skill note.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/workflows/evaluate.md` - primary runtime workflow update.
- [x] `skills/quality/SKILL.md` - reviewed; no edit needed because root guidance
      already frames evaluate feedback logs as process-only and non-authoritative.
- [x] `skills/quality/workflows/log.md` - record the runtime workflow edit.

### Install / scaffold

None.

## Children

- [Functional spec](0095-evaluate-feedback-log-outcomes/spec.md) - what the
  evaluate feedback-log contract must say.
- [Design doc](0095-evaluate-feedback-log-outcomes/design.md) - why the log stays
  outside evaluation runs and how outcome values are framed.

## Status

`Done`. Implemented and archived. Evaluate feedback logs remain process-only
workflow artifacts under `.quality/logs/`, and their `outcome` values now
describe workflow terminal states instead of report, rating, or recommendation
semantics. Verified with `mise run check`. No CLI, Go, format-schema, rating,
roll-up, report, or evaluation-record behavior changed.
