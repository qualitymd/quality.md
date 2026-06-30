---
type: Change Case
title: Report Gallery
description: Add a generated example Evaluation report gallery for fast report-design iteration and public inspection.
status: Done
tags: [evaluation, reports, examples]
timestamp: 2026-06-27T00:00:00Z
---

# Report Gallery

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0156-report-gallery/spec.md) - what the case must do.
- [Design doc](0156-report-gallery/design.md) - how it is built, and why.

## Motivation

Evaluation report design is currently slow to iterate on because the practical
preview loop tends to run through a release, a fresh `/quality evaluate`, and
then manual report inspection. The report builder is already a deterministic
projection of persisted Evaluation data, so report design should be previewable
from a representative run without cutting a release or re-running agent
judgment.

The same representative run can also serve as a public example. Users and
contributors need a realistic sample `QUALITY.md` and report tree they can
browse without running an evaluation, while still knowing the judgment data is
synthetic and illustrative.

## Scope

Covered:

- Add a checked-in report gallery example under `examples/` with a sample
  software-service `QUALITY.md`.
- Generate the example's `.quality/evaluations/0001-full-eval/` run from code,
  including persisted `data/` payloads and generated reports.
- Validate synthetic payloads through the current Evaluation data write path
  before persisting them.
- Add a deterministic regeneration task and a check task that fails when the
  checked-in gallery is stale.
- Document that the example uses synthetic routine outputs and intentionally
  omits a concrete source system/code tree.

Deferred:

- Adding a fictional source-code or documentation tree for the synthetic
  evidence references.
- Adding a non-software paired gallery example.
- Adding a new public `qualitymd` command for galleries.
- Supporting historical payload contracts for the gallery.

## Affected artifacts

Derived by sweeping for report generation, `examples/`, `mise` tasks, and
Evaluation report contracts.

**Code**

- [x] `scripts/report-gallery/` - add the deterministic generator.
- [x] `mise.toml` - add regeneration and freshness-check tasks.

**Durable specs**

- [x] No durable spec changes. This adds a repository example and development
      task, not a new `qualitymd` runtime contract.

**Durable docs / examples**

- [x] `CONTRIBUTING.md` - list the report-gallery task with development tasks.
- [x] `examples/report-gallery/software-service/README.md` - explain the
      illustrative example and link to generated reports.
- [x] `examples/report-gallery/software-service/QUALITY.md` - sample model.
- [x] `examples/report-gallery/software-service/.quality/evaluations/` -
      generated example Evaluation run.
- [x] `changes/index.md`, `changes/archive/index.md`, and `changes/log.md` -
      Change Case lifecycle.

## Status

`Done`. Implemented and archived. The generated report gallery, regeneration
task, freshness check, contributor note, and synthetic-data disclosure are in
place; `go test ./...`, `mise run fmt-md-check`, and
`mise run report-gallery-check` pass.
