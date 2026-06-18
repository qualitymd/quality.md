---
type: Change
title: Report regression coverage
description: Add focused regression tests for high-risk evaluation report behavior.
status: Done
tags: [evaluation, report, tests]
timestamp: 2026-06-18T00:00:00Z
---

# Report regression coverage

A unit of work that turns repeated report-rendering experiment findings into
automated coverage. Detail lives in the child:

- [Functional spec](0024-report-regression-coverage/spec.md) - what report
  regressions must be covered.
- [Design doc](0024-report-regression-coverage/design.md) - how the tests avoid
  committed fixture churn.

## Motivation

The experiment program repeatedly found report-rendering regressions around
seeded safety cases, prompt-injection handling, not-assessed propagation,
dotted-path limitation extraction, structural roots, and empty recommendation
arrays. These are now accepted behavior and should be protected by focused
tests.

## Scope

Covered: add tests that construct temporary evaluation runs and assert the
high-risk rendered `report.md` / `report.json` properties.

Deferred: committed benchmark fixture snapshots, browser/visual report tests,
and end-to-end real-repo fixture harnesses.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches
`In-Review`:

- [x] `internal/evaluation` tests - add focused report-regression coverage.

## Status

`Done`. Implemented and archived after adding focused temp-run tests for secret-style, prompt-injection-style, not-assessed, dotted-path, structural-root, and empty-recommendation report behavior.
