---
type: Change Case
title: Planned coverage status
description: Let evaluation status compare optional planned coverage against written records.
status: Done
tags: [evaluation, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Planned coverage status

A unit of work that turns the experiment program's planned-coverage prototype
into a CLI-supported evaluation status capability. Detail lives in the child:

- [Functional spec](0020-planned-coverage-status/spec.md) - what planned
  coverage metadata and status comparison must do.
- [Design doc](0020-planned-coverage-status/design.md) - how the CLI stores and
  compares planned coverage.

## Motivation

The E11 interruption/resume experiment showed that a partially written run can
be resumed through normal `add-record` commands, but `show-status` only reports
coarse missing analysis. It cannot name which planned assessments are still
unwritten.

The planned-coverage prototype showed a small structured manifest can close that
gap by comparing planned assessment and analysis coverage to records actually
written. This should become an optional status input so scoped quick runs remain
valid while resumable standard/deep runs get better diagnostics.

## Scope

Covered: define optional planned coverage metadata for evaluation runs and make
`show-status` report missing planned assessments, missing planned analyses, and
unexpected records when that metadata exists.

Deferred: enforcing full-model coverage without a plan, changing record
payloads, adding replacement/superseding semantics, and changing report
rendering beyond the existing renderability gate.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches
`In-Review`:

- [x] [`specs/cli/evaluation-show-status.md`](../../specs/cli/evaluation-show-status.md) - document planned-coverage status gaps or warnings.
- [x] [`specs/cli/evaluation-set-planned-coverage.md`](../../specs/cli/evaluation-set-planned-coverage.md) - document the planned coverage writer command.
- [x] [`specs/cli.md`](../../specs/cli.md) and
      [`specs/cli/index.md`](../../specs/cli/index.md) - list the planned coverage
      writer command.
- [x] [`specs/evaluation-records.md`](../../specs/evaluation-records.md) - define
      the optional planned coverage runtime artifact.
- [x] [`README.md`](../../README.md) - list the planned coverage writer in the CLI
      surface.
- [x] [quality skill spec](../../specs/skills/quality-skill/quality-skill.md) - require the skill to record planned coverage when the run needs
      resume diagnostics.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - teach the skill
      when to create planned coverage metadata.

## Status

`Done`. Implemented and archived after implementing `qualitymd evaluation set-planned-coverage`, planned-coverage status gaps, durable specs/docs, and skill prompt guidance.
