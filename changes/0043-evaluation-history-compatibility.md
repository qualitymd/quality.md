---
type: Change Case
title: Evaluation history compatibility
description: Make evaluation-history readers tolerant of old, malformed, partial, or hand-edited run records while keeping current writers strict and report generation explicit.
status: In-Review
tags: [evaluation, records, compatibility, status]
timestamp: 2026-06-20T00:00:00Z
---

# Evaluation history compatibility

A Change Case for hardening evaluation-history behavior as the evaluation
process and record implementation continue to change.

> **In-Review.** Implementation, durable spec and skill updates, regression
> coverage, and verification are complete and ready for review.

## Motivation

Evaluation records are runtime artifacts created by the `/quality` skill and
`qualitymd` CLI, but users and agents can encounter records that were produced
by older tooling, hand-edited, copied between repos, partially written, or
generated outside the current skill/CLI pair. As the evaluation process changes,
those records will frequently fail the current shape. That must not break normal
use cases such as checking project status, listing history, selecting a fresh
run, or starting a new evaluation.

The project does not want schema migrations or compatibility writers. The
practical contract should be: current writers stay strict, historical readers
inspect tolerantly, and report/gate commands fail clearly for only the selected
unreportable run.

## Scope

Covered: tolerant inspection of existing run folders and records, clear
run-level compatibility gaps, status/list/latest behavior around incompatible
runs, report/gate failure messages for selected incompatible runs, and durable
spec/skill guidance that encodes the compatibility posture.

Deferred / non-goals: no migration framework, no automatic rewriting of
historical runs, no compatibility aliases for invalid current records, no change
to QUALITY.md frontmatter semantics, and no support for manually creating new
evaluation records outside the CLI.

## Affected artifacts

### Code

- [x] `internal/evaluation/load.go` — split or extend loading so inspection can
      return partial run state with record problems instead of failing the whole
      run.
- [x] `internal/evaluation/types.go` — add compatibility/reportability gap kinds
      for malformed, unsupported, missing-version, or otherwise incompatible
      records.
- [x] `internal/evaluation/list.go` — keep `evaluation list` usable when one run
      cannot be strictly loaded, and expose incompatible runs as incomplete or
      problematic rather than aborting the listing.
- [x] `internal/evaluation/report.go` — ensure report build/gate fail with
      clear non-reportable-run diagnostics for selected incompatible runs.
- [x] `internal/status/status.go` — preserve tolerant project-status behavior and
      align summaries with the new compatibility gaps.
- [x] `internal/cli/evaluation.go` — render human and JSON status/list output so
      compatibility problems are actionable.
- [x] `internal/evaluation/evaluation_test.go`, `internal/status/status_test.go`,
      and `internal/cli/evaluation_test.go` — cover historical and malformed run
      records without making normal status/list workflows fail.

### Durable specs

See the functional spec's
[Durable spec changes](0043-evaluation-history-compatibility/spec.md#durable-spec-changes)
for the per-requirement breakdown.

- [x] [`specs/evaluation-records.md`](../specs/evaluation-records.md)
- [x] [`specs/cli/evaluation-list.md`](../specs/cli/evaluation-list.md)
- [x] [`specs/cli/evaluation-status.md`](../specs/cli/evaluation-status.md)
- [x] [`specs/cli/evaluation-report.md`](../specs/cli/evaluation-report.md)
- [x] [`specs/cli/status.md`](../specs/cli/status.md)
- [x] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
- [x] [`specs/log.md`](../specs/log.md)

No `SPECIFICATION.md` change is expected: this concerns runtime evaluation
records, not QUALITY.md document semantics.

### Durable docs and skill runtime

- [x] `skills/quality/SKILL.md` — encode the historical-run posture for wizard,
      evaluate, and improve workflows.
- [x] `skills/quality/resources/cli-quick-reference.md` — mention that
      incompatible historical runs are status/list problems rather than a reason
      to stop normal evaluation setup.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` — align review-history
      advice with compatibility gaps.

### Fixtures

- [x] Add or update targeted fixtures/tests for missing `schemaVersion`, future
      `schemaVersion`, malformed JSON, malformed recommendation frontmatter,
      missing required current fields, and stale/partial runs.

## Children

- [Functional spec](0043-evaluation-history-compatibility/spec.md) — what the
  change must do.
- [Design doc](0043-evaluation-history-compatibility/design.md) — how it will be
  built, and why.
