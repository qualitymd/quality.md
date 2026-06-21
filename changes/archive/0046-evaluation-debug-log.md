---
type: Change Case
title: Evaluation debug log
description: Add a process-only debug log artifact to evaluation runs.
status: Done
tags: [evaluation, records, skill]
timestamp: 2026-06-21T00:00:00Z
---

# Evaluation debug log

This Change Case adds a process-only `debug-log.md` artifact to evaluation runs
so maintainers can understand evaluation orchestration, retries, recovery,
coverage decisions, and tooling friction without confusing those notes with
subject-quality evidence.

Detail lives in:

- [Functional spec](0046-evaluation-debug-log/spec.md) - what the change must do.
- [Design doc](0046-evaluation-debug-log/design.md) - how it is built, and why.

## Motivation

Evaluation records explain the subject's quality, but they do not consistently
explain the evaluator's process: ambiguous scope resolution, interrupted or
resumed work, discarded unsafe output, report build retries, or cases where a
project command was routed into a formal assessment instead of copied into a
process note. A small `debug-log.md` makes those process events auditable while
keeping assessments, analysis, recommendations, and reports authoritative for
ratings.

## Scope

Covered:

- Add `debug-log.md` to the evaluation run folder contract and run scaffold.
- Define the artifact's process-only boundary, including how subject commands
  used as evidence may be referenced without duplicating assessment evidence.
- Update `/quality` evaluation guidance so the skill hand-authors the log for
  notable evaluation-process events.
- Update durable specs, examples, tests, and release notes to match.

Deferred:

- A dedicated `qualitymd evaluation debug ...` command.
- Generated report content derived from `debug-log.md`.
- Structured machine parsing of debug-log events.

## Affected artifacts

### Code

- [x] `internal/evaluation/create.go` - seed `debug-log.md` in new evaluation
      runs.
- [x] `internal/evaluation/evaluation_test.go` - assert run creation and fixture
      behavior with `debug-log.md`.

### Durable specs

- [x] `specs/evaluation-records.md` - add the runtime artifact contract and
      process-only boundary.
- [x] `specs/cli/evaluation-create.md` - include seeded `debug-log.md`.
- [x] `specs/skills/quality-skill/quality-skill.md` - update evaluation
      workflow and artifact expectations.
- [x] `specs/skills/quality-skill/examples/` - update the example run and index
      to include the process log.
- [x] `specs/log.md` - record the durable spec update.

### Durable docs

- [x] `docs/guides/use-quality-skill.md` - update the evaluation artifact
      overview if it names run contents.
- [x] `docs/reference/versioning.md` - update the skill metadata example for
      the new compatibility line.
- [x] `docs/log.md` - record durable documentation updates if docs change.

### Bundled skill

- [x] `skills/quality/SKILL.md` - include `debug-log.md` in evaluation artifacts
      and hard rules.
- [x] `skills/quality/modes/evaluate.md` - instruct the evaluator to maintain
      the process log and preserve the assessment boundary.
- [x] `skills/quality/resources/cli-quick-reference.md` - update run artifact
      notes if needed.

### Release

- [x] `CHANGELOG.md` - add user-facing release notes.
- [x] `skills/quality/SKILL.md` - bump skill metadata for the release.

## Status

`Done`. Implemented, verified, and archived.
