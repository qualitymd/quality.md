---
type: Change Case
title: Assessment superseding
description: Let corrected assessment records supersede stale assessments.
status: Done
tags: [evaluation, assessments, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Assessment superseding

A unit of work that turns the remaining correction-workflow gap into a
deterministic assessment status capability. Detail lives in the child:

- [Functional spec](0023-assessment-superseding/spec.md) - what assessment
  superseding must do.
- [Design doc](0023-assessment-superseding/design.md) - how status and reports
  keep corrected assessments coherent with analyses.

## Motivation

Duplicate assessment detection prevents contradictory reports, but it leaves no
ergonomic way to correct an assessment inside a run. E28 solved this for
recommendations with explicit superseding. Assessments need a stricter variant:
a corrected assessment may supersede a stale one, but analyses must reference
the active assessment so roll-ups cannot silently rely on stale judgment.

## Scope

Covered: allow assessment records to declare older same-requirement assessments
they supersede; make status validate superseding references and stale analysis
references; make reports distinguish active assessments from superseded
assessments.

Deferred: deleting old assessment files, in-place replacement, superseding
analysis records, and automatic duplicate-assessment resolution without
explicit evaluator intent.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches
`In-Review`:

- [x] [`specs/evaluation-records.md`](../../specs/evaluation-records.md) - define
      optional assessment superseding metadata.
- [x] [`specs/cli/evaluation-add-record.md`](../../specs/cli/evaluation-add-record.md) - document assessment payload support for superseding.
- [x] [`specs/cli/evaluation-show-status.md`](../../specs/cli/evaluation-show-status.md) - document invalid superseding references and stale analysis references.
- [x] [`specs/cli/evaluation-build-report.md`](../../specs/cli/evaluation-build-report.md) - document active/superseded assessment rendering.
- [x] [quality skill spec](../../specs/skills/quality-skill/quality-skill.md) -
      guide correction workflows to update analyses when superseding
      assessments.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - teach the skill
      how to correct assessment records without leaving stale roll-ups.

## Status

`Done`. Implemented and archived after implementing assessment `supersedes` metadata, superseding status gaps, active/superseded report rendering, durable specs, and skill guidance.
