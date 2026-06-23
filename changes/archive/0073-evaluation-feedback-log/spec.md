---
type: Functional Specification
title: Evaluation feedback log - functional spec
description: Requirements for replacing evaluation debug logging with an evaluate workflow feedback log aligned with setup feedback logs.
tags: [skill, evaluation, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Evaluation feedback log - functional spec

Companion to the
[Evaluation feedback log](../0073-evaluation-feedback-log.md). This spec defines
the delta from evaluation's run-local `debug-log.md` concept to a workflow
feedback log aligned with setup's `.quality/logs/` feedback artifact.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background

The evaluation `debug-log.md` added by
[0046](../archive/0046-evaluation-debug-log.md) was process-only by design, but
the artifact name and sparse seed do not communicate the intended maintainer
feedback purpose. Setup's feedback log has since established a clearer pattern:
a local, shareable, workflow-experience artifact with lifecycle metadata,
redaction rules, and sections that turn friction into actionable workflow
improvements.

Evaluation should use that same feedback-log family. Its process notes are not
rating evidence or report content; they are improvement signals about running the
evaluation workflow.

## Requirements

### Shared feedback-log contract

A shared durable spec **MUST** define the workflow feedback-log artifact for all
`/quality` workflows that adopt it.

The shared contract **MUST** live outside a specific workflow, at
`specs/skills/quality-skill/workflow-feedback-log.md`.

The shared contract **MUST** own the workflow-neutral rules for purpose,
location, naming, lifecycle status, frontmatter metadata, body sections,
redaction, local-only/no-transmission posture, and non-authority boundaries.

The shared contract **MUST NOT** be setup-shaped. Setup and evaluate must be
workflow adopters of the contract rather than each owning separate generic
feedback-log rules.

### Evaluate feedback-log adoption

The evaluate workflow **MUST** write a workflow feedback log under
`.quality/logs/<timestamp>-evaluate-feedback-log.md`, using the shared
workflow-feedback-log contract.

The evaluate-specific durable spec **MUST** live at
`specs/skills/quality-skill/workflows/evaluate/feedback-log.md`.

The evaluate feedback log **MUST** record the experience of running evaluation,
including material events such as scope ambiguity, history inspection friction,
coverage adjustment, interruption or resume, retries, record corrections,
tooling failures, slow phases, redaction decisions, prompt-injection handling,
report generation recovery, UX/AX observations, what worked well, and suggested
workflow improvements.

The evaluate feedback log **MUST NOT** be authoritative for ratings, findings,
recommendations, next actions, generated reports, or QUALITY.md model content.

The evaluate feedback log **MUST NOT** duplicate assessment evidence, rating
rationale, raw project-command output, or recommendation prose as primary
content. When a project command is exercised as evaluation evidence, the feedback
log **MAY** record the routing decision and point to the formal assessment
record after that record exists.

The evaluate feedback log **MUST NOT** contain secret values or raw
prompt-injection text. When a secret or hostile instruction affects evaluation,
the log **MAY** record a sanitized locator and type.

### Evaluate lifecycle

`/quality evaluate` **MUST** create the evaluate feedback log after CLI support
is verified and the run frame is emitted, before assessment evidence collection
or record writes begin.

The initial evaluate feedback log **MUST** include all shared frontmatter fields
and body sections with `status: in-progress`.

After the numbered evaluation run is created, the evaluate feedback log **MUST**
record the run path in workflow-specific metadata or timeline content.

`/quality evaluate` **MUST** update the current run's evaluate feedback log when
material workflow-experience information appears, refreshing `updated-at` for
material updates. It **SHOULD** avoid noisy churn for ordinary happy-path steps.

At normal close, `/quality evaluate` **MUST** finalize the feedback log with
`status: completed`, set `completed-at`, record the evaluation outcome when
known, and make empty sections explicit with a no-notable-content note.

When evaluation stops after the log exists because lint fails, source cannot be
resolved, requirements are not assessable, CLI support fails, user input is
needed, or another non-success condition occurs, the skill **SHOULD** finalize
the log with `status: failed` or `status: interrupted` when it can do so without
masking the stop condition.

Writing, updating, or finalizing the evaluate feedback log **MUST NOT** change
evaluation stop rules, rating semantics, record authority, reportability, or
next-step routing.

### Evaluation run scaffold

New evaluation runs **MUST NOT** seed `debug-log.md`.

`qualitymd evaluation create` **MUST** continue to seed the run artifacts that
belong to the evaluation record itself: `model.md`, `design.md`, `plan.md`, and
the record directories.

Historical runs that contain `debug-log.md` **MUST** remain readable as
historical evaluation runs. No migration of historical `debug-log.md` artifacts
is required.

Live docs and specs **SHOULD** describe `debug-log.md` only as a legacy artifact
for historical runs, not as the current process-log contract for new evaluation
runs.

### Setup adopter refactor

Setup's existing feedback-log durable spec **MUST** retain setup-specific
behavior while deferring shared artifact identity, metadata, body shape,
redaction, and no-transmission rules to the shared feedback-log spec.

Setup's runtime behavior **MUST NOT** regress: setup still creates, updates, and
finalizes `.quality/logs/<timestamp>-setup-feedback-log.md`.

## Acceptance criteria

- A shared feedback-log spec exists at
  `specs/skills/quality-skill/workflow-feedback-log.md`.
- An evaluate feedback-log spec exists at
  `specs/skills/quality-skill/workflows/evaluate/feedback-log.md`.
- Setup's feedback-log spec references the shared contract instead of owning all
  generic feedback-log rules itself.
- `/quality evaluate` runtime instructions create, maintain, and finalize
  `.quality/logs/<timestamp>-evaluate-feedback-log.md`.
- The evaluate run frame and artifact wording name a feedback log rather than a
  debug log.
- New evaluation run scaffolds no longer include `debug-log.md`.
- Historical `debug-log.md` is documented as legacy-compatible where live specs
  still need to mention it.
- Evaluation feedback logs include lifecycle frontmatter and setup-style body
  sections adapted for evaluation.
- Evaluation feedback logs record workflow experience only; they do not
  duplicate ratings, findings, assessment evidence, raw command output, or
  recommendation content.
- Secret and prompt-injection redaction rules apply to evaluation feedback logs.
- Examples and links under `specs/skills/quality-skill/examples/` no longer
  present `debug-log.md` as a current artifact.
- Tests covering evaluation run creation are updated to the new scaffold.
- `CHANGELOG.md`, relevant `index.md` files, and relevant `log.md` files are
  updated.
- Repository verification passes with the relevant formatting and Go test gates.

## Durable spec changes

### To add

- `specs/skills/quality-skill/workflow-feedback-log.md` - define the shared
  workflow feedback-log artifact contract.
- `specs/skills/quality-skill/workflows/evaluate/feedback-log.md` - define
  evaluate-specific feedback-log adoption.
- `specs/skills/quality-skill/workflows/evaluate/index.md` - list the evaluate
  child spec.

### To modify

- `specs/skills/quality-skill/workflows/setup/feedback-log.md` - defer shared
  artifact rules to the shared spec and keep setup-specific behavior.
- `specs/skills/quality-skill/workflows/setup/index.md` - keep the setup child
  listing accurate after the shared/adopter split.
- `specs/skills/quality-skill/workflows/setup.md` - update references to the
  shared/adopter split if needed.
- `specs/skills/quality-skill/workflows/index.md` - list the evaluate child spec
  folder.
- `specs/skills/quality-skill/evaluation.md` - replace debug-log workflow
  guidance with evaluate feedback-log guidance.
- `specs/skills/quality-skill/reporting.md` - remove `debug-log.md` from the
  current run-folder artifact contract and document historical compatibility.
- `specs/skills/quality-skill/quality-skill.md` - update parent feedback-log and
  evaluation artifact guidance.
- `specs/evaluation-records/run-folder.md` - remove `debug-log.md` from new run
  folders and document legacy compatibility.
- `specs/evaluation-records/debug-log-md.md` - revise as a legacy artifact
  contract or remove if all live references are eliminated.
- `specs/evaluation-records/index.md`, `specs/evaluation-records/log.md`, and
  `specs/evaluation-records.md` - keep evaluation-record listings and history
  current.
- `specs/cli/evaluation-create.md` - update the scaffold contract for new
  evaluation runs.
- `specs/log.md` and `specs/skills/quality-skill/workflows/log.md` - record the
  durable spec updates.

### To rename

None

### To delete

None
