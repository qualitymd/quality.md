---
type: Change Case
title: Evaluation feedback log
description: Replace evaluation's debug-log concept with an evaluate workflow feedback log aligned with setup feedback logging.
status: Done
tags: [skill, evaluation, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Evaluation feedback log

A **Change Case** to align `/quality evaluate` process logging with setup's
workflow feedback log. Evaluation should record the experience of running the
workflow in a local, shareable feedback artifact under `.quality/logs/`, using
the same lifecycle, metadata, redaction posture, and improvement-oriented
sections as setup. The existing evaluation `debug-log.md` concept becomes
historical compatibility, not the primary artifact for new runs.

Detail lives in:

- [Functional spec](0073-evaluation-feedback-log/spec.md) - what the change must
  do.
- [Design doc](0073-evaluation-feedback-log/design.md) - how the change is
  shaped, and why.

## Motivation

Evaluation's `debug-log.md` was intended as process-only, but its name and sparse
shape make it feel like an internal scratchpad. Setup's feedback log is clearer:
it records workflow experience, carries enough environment metadata to make
friction actionable, and keeps a strict boundary away from model content,
ratings, findings, recommendations, and evidence.

Evaluation has the same improvement need. Scope ambiguity, history inspection,
redaction, retries, slow phases, report recovery, and confusing interaction
points are signals for improving the `/quality` workflow itself. They should
land in the same feedback-log family as setup rather than in a separate
"debug" artifact.

## Scope

Covered:

- Add a shared workflow feedback-log durable spec outside the setup workflow.
- Add an evaluate-specific feedback-log durable spec nested under the evaluate
  workflow.
- Refactor setup's feedback-log durable spec so shared rules live in the shared
  contract and setup keeps only setup-specific adoption behavior.
- Make `/quality evaluate` write
  `.quality/logs/<timestamp>-evaluate-feedback-log.md` as the primary workflow
  feedback artifact.
- Remove `debug-log.md` from new evaluation run scaffolds and update the
  evaluation run-folder contract.
- Preserve historical `debug-log.md` files as legacy-compatible evaluation
  artifacts; do not migrate historical runs.
- Update runtime skill instructions, durable specs, examples, tests, and release
  notes.

Deferred / non-goals:

- No QUALITY.md format change.
- No telemetry, upload endpoint, automatic sharing, external transmission, or
  consent gate.
- No structured parser, validator, append command, or report renderer for
  feedback logs.
- No change to evaluation rating semantics, report authority, findings,
  recommendations, or record payload contracts.
- No migration of historical evaluation runs.

## Affected artifacts

### Code

- [x] `internal/evaluation/create.go` - stop seeding `debug-log.md` in new
      evaluation run folders.
- [x] `internal/evaluation/evaluation_test.go` - update run creation assertions
      for the revised scaffold.

### Durable specs

- [x] `specs/skills/quality-skill/workflow-feedback-log.md` - add a shared
      workflow feedback-log contract.
- [x] `specs/skills/quality-skill/workflows/evaluate/feedback-log.md` - add the
      evaluate-specific feedback-log adopter contract.
- [x] `specs/skills/quality-skill/workflows/evaluate/index.md` - list the new
      evaluate child spec.
- [x] `specs/skills/quality-skill/workflows/setup/feedback-log.md` - refactor to
      reference the shared contract and retain setup-specific behavior.
- [x] `specs/skills/quality-skill/workflows/setup/index.md` - update wording if
      needed after the shared-spec extraction.
- [x] `specs/skills/quality-skill/workflows/setup.md` - update references to the
      shared/adopter split if needed.
- [x] `specs/skills/quality-skill/workflows/index.md` - list the evaluate child
      spec folder.
- [x] `specs/skills/quality-skill/evaluation.md` - replace debug-log workflow
      guidance with evaluate feedback-log guidance.
- [x] `specs/skills/quality-skill/reporting.md` - update the run-folder artifact
      contract and legacy `debug-log.md` treatment.
- [x] `specs/skills/quality-skill/quality-skill.md` - update the parent workflow
      feedback-log description and evaluation artifact list.
- [x] `specs/evaluation-records/run-folder.md` - remove `debug-log.md` from new
      run folders and document legacy compatibility.
- [x] `specs/evaluation-records/debug-log-md.md` - revise as a legacy artifact
      contract or remove if all live references are eliminated.
- [x] `specs/evaluation-records/index.md`,
      `specs/evaluation-records/log.md`, and `specs/evaluation-records.md` -
      update listings and revision history.
- [x] `specs/cli/evaluation-create.md` - update the create scaffold contract.
- [x] `specs/log.md` and `specs/skills/quality-skill/workflows/log.md` - record
      durable spec updates.

### Durable docs

- [x] `CHANGELOG.md` - add the user-facing unreleased note.

### Bundled skill

- [x] `skills/quality/SKILL.md` - update hard rules, run-frame artifact wording,
      and process-log guidance.
- [x] `skills/quality/workflows/evaluate.md` - create, maintain, and finalize the
      evaluate feedback log.
- [x] `skills/quality/workflows/setup.md` - update shared/adopter wording if
      needed.
- [x] `skills/quality/resources/cli-quick-reference.md` - replace debug-log
      guidance with evaluate feedback-log guidance.

### Examples

- [x] `specs/skills/quality-skill/examples/` - update example evaluation run
      structure and links away from `debug-log.md`.

### Release

- [x] `CHANGELOG.md` - add the user-facing `/quality Skill` and CLI notes under
      `Unreleased`.

## Children

- [Functional spec](0073-evaluation-feedback-log/spec.md) - required behavior
  and acceptance criteria for evaluation feedback logging.
- [Design doc](0073-evaluation-feedback-log/design.md) - artifact placement,
  compatibility, and implementation approach.

## Status

`Done`. Implemented and archived. Added the shared workflow feedback-log spec,
the evaluate adopter spec, runtime evaluate feedback-log guidance, CLI scaffold
change, durable spec updates, example updates, release notes, and tests.
Verification passed with `mise run check`.
