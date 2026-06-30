---
type: Change Case
title: Duplicate assessment status
description: Make evaluation status reject duplicate assessment records for the same target requirement.
status: Done
tags: [evaluation, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Duplicate assessment status

A unit of work that turns an experiment-discovered correction workflow hazard
into deterministic reportability behavior. Detail lives in the child:

- [Functional spec](0019-duplicate-assessment-status/spec.md) - what duplicate
  assessment detection must do.
- [Design doc](0019-duplicate-assessment-status/design.md) - how status and
  report rendering should enforce it.

## Motivation

The experiment program found that re-adding a corrected assessment for the same
target and requirement appends a new numbered assessment record. `show-status`
still reports the run as reportable, and `build-report` renders conflicting
requirement entries while the headline roll-up may still come from an analysis
that references only the original assessment.

That makes correction and resume workflows ambiguous. The durable record
contract already says assessment records are one JSON file per assessed
requirement, so duplicate records should be visible before report rendering.

## Scope

Covered: detect duplicate assessment records that share the same target path and
requirement, report them as `duplicate-assessment` gaps, and make report
rendering fail through the existing renderability gate.

Deferred: replacing, deleting, or superseding existing records; adding an
`add-record --replace` flag; and changing recommendation duplicate behavior.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches
`In-Review`:

- [x] [`specs/cli/evaluation-show-status.md`](../../specs/cli/evaluation-show-status.md) - document the duplicate-assessment reportability gap.
- [x] [`specs/cli/evaluation-build-report.md`](../../specs/cli/evaluation-build-report.md) - clarify that non-renderable duplicate-assessment runs fail before
      writing reports.
- [x] [`specs/evaluation-records.md`](../../specs/evaluation-records.md) - clarify
      assessment record uniqueness for a target path and requirement.

## Status

`Done`. Implemented and archived after implementing duplicate-assessment renderability checks, syncing durable specs, and verifying the command-boundary duplicate trial.
