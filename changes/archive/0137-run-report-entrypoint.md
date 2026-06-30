---
type: Change Case
title: Run Report Entrypoint
description: Make report.md the Evaluation run entrypoint and move the root Area report to root-area.md.
status: Done
tags: [evaluation, reports, cli, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Run Report Entrypoint

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0137-run-report-entrypoint/spec.md) - what the case must do.
- [Design doc](0137-run-report-entrypoint/design.md) - how it's built, and why.

## Motivation

Evaluation currently writes the root Area report to `report.md`. That was a
reasonable single-entrypoint shortcut while every run implied the whole Model,
but scoped evaluations make the name misleading: a user who evaluates a
sub-Area or Factor still opens `report.md` and lands on the root Area report.
The run needs a true run-level report that states the evaluated scope, headline
subject, coverage, limits, and links to the generated subject reports. The root
Area report should be named like the Area report it is.

## Scope

Covered:

- Generate `report.md` as the Evaluation run-level report.
- Rename the root Area report from `report.md` to `root-area.md`.
- Keep non-root Area, Factor, and Requirement report paths unchanged.
- Make reportability and headline output work for scoped Area and Factor runs
  without requiring root Area analysis when root is out of scope.
- Index the run report and headline report explicitly in
  `EvaluationOutputResult` and build receipts.
- Update durable specs, CLI tests, generated data schema/examples, and
  `/quality` skill contracts.

Deferred:

- Backward-compatible duplicate root Area report paths.
- Automatic migration of old evaluation runs.
- A separate Advice or recommendation report.
- Interactive report viewers.

## Affected artifacts

Derived by sweeping for `report.md`, root Area report paths, report refs,
`EvaluationOutputResult`, report build receipts, and reportability across code,
durable specs, runtime skill content, generated schemas, tests, examples, and
release notes.

**Code**

- [x] `internal/evaluation/report_tree.go` - render a run-level `report.md`,
      move root Area reports to `root-area.md`, and emit explicit run/headline
      output refs.
- [x] `internal/evaluation/report.go` - update the build receipt fields and
      reportability gate behavior.
- [x] `internal/evaluation/display.go` - add a run report kind.
- [x] `internal/evaluation/data.go` - update generated examples and report ref
      helpers.
- [x] `internal/evaluation/data_contract.go` - update `EvaluationOutputResult`
      shape.
- [x] `internal/evaluation/evaluation-data.schema.json` - regenerate the schema
      from the typed data contract.
- [x] `internal/evaluation/evaluation_test.go` - update root path assertions and
      add scoped Area/Factor reportability coverage.
- [x] `internal/cli/evaluation.go` - update human report build output.

**Durable specs** (substance in the [functional spec](0137-run-report-entrypoint/spec.md))

- [x] `SPECIFICATION.md` - update Evaluation Report semantics from root-only to
      headline outcome.
- [x] `specs/evaluation/reports/report-tree.md` - define the run report and
      root Area report path.
- [x] `specs/evaluation/records/data-layout.md` - update generated report tree
      layout.
- [x] `specs/evaluation/records/json-conventions.md` - add the run report kind
      and headline report ref conventions.
- [x] `specs/evaluation/records/payload-kinds.md` - update
      `EvaluationOutputResult` ownership and shape.
- [x] `specs/evaluation/protocol.md` - align report generation with scoped
      headline subjects.
- [x] `specs/evaluation/orchestration.md` - align reportability dependencies
      with scoped outputs.
- [x] `specs/cli/evaluation-report.md` - update report build behavior and
      receipt expectations.
- [x] `specs/skills/quality-skill/evaluation.md` - align scoped evaluation
      workflow and report build expectations.
- [x] `specs/skills/quality-skill/reporting.md` - align `/quality` report tree
      and closeout wording.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` - update runtime report path contract.
- [x] `skills/quality/workflows/evaluate.md` - update evaluate workflow
      closeout/report path guidance.
- [x] `CHANGELOG.md` - note the report entrypoint break.

No planned impact: `QUALITY.md` format authoring rules, `qualitymd init`, setup
workflow, recommendation follow-up behavior, or install docs.

## Status

`Done`. Implemented and archived after code, durable specs, generated schema,
runtime skill guidance, changelog, and tests were updated.
