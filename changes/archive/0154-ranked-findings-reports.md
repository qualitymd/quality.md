---
type: Change Case
title: Ranked Findings Reports
description: Make ranked findings easier to scan, navigate, and inspect across generated Evaluation reports.
status: Done
tags: [evaluation, reports, advice, findings]
timestamp: 2026-06-27T00:00:00Z
---

# Ranked Findings Reports

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0154-ranked-findings-reports/spec.md) - what the case must do.
- [Design doc](0154-ranked-findings-reports/design.md) - how it is built, and why.

## Motivation

Evaluation Advice now ranks findings, but the generated report surface does not
fully help readers use that ordering. `report.md` shows a compact Top Findings
table, yet it hides the Area and Factor context, points only to the Requirement
report rather than the exact finding detail, and has no full findings index
parallel to `recommendations.md`.

Readers should be able to start at `report.md`, scan the highest-ranked findings
with model context, jump to the exact finding detail, and then open a complete
ranked findings index when the top 10 is not enough. When they land on an
individual Requirement report, the finding detail should also explain the
Advice ranking context that brought them there.

## Scope

Covered:

- Redesign `report.md` Top Findings to use the columns `Rank`, `Finding`,
  `Area`, `Factors`, `Type`, and `Severity`.
- Link finding statements to stable finding-detail anchors on Requirement
  reports.
- Link Area and Factor display names to their generated reports.
- Generate a full `findings.md` ranked findings index using the same columns as
  Top Findings.
- Always link from `report.md` to `findings.md`.
- Show finding Advice rank, tier, and ranking rationale in Requirement finding
  detail sections.
- Render finding type and severity with existing display names and emoji.

Deferred:

- Adding new Advice payload fields or changing `FindingRankingResult`.
- Adding filtering, grouping, or sorting controls to `findings.md`.
- Showing finding confidence in the Top Findings or findings-index tables.
- Showing recommendation coverage status in the findings tables.

## Affected artifacts

Derived by sweeping for `Top Findings`, `FindingRankingResult`,
`Finding Details`, `recommendations.md`, `report.md`, generated report paths,
and report tree contracts.

**Code**

- [x] `internal/evaluation/report_tree.go` - render the new Top Findings table,
      generate `findings.md`, link finding statements to stable anchors, and add
      Advice ranking context to Requirement finding details.
- [x] `internal/evaluation/display.go` - add report kind display metadata if the
      new findings index becomes a first-class rendered report kind.
- [x] `internal/evaluation/data.go` - include the findings report in the
      `EvaluationOutputResult` example.
- [x] `internal/evaluation/evaluation-data.schema.json` - regenerate the
      companion schema for the `findings` report kind.
- [x] `internal/evaluation/evaluation_test.go` - cover the new report table,
      links, finding index, and Requirement detail ranking block.

**Format spec and durable specs** (substance in the [functional spec](0154-ranked-findings-reports/spec.md))

- [x] `specs/evaluation/reports/report-tree.md` - require the ranked findings
      report surfaces and stable finding-detail links.
- [x] `specs/cli/evaluation-report.md` - require `evaluation report build` to
      render persisted Advice into `findings.md`.
- [x] `specs/evaluation/records/data-layout.md` - include `findings.md` in the
      generated report tree paths.
- [x] `SPECIFICATION.md` - no change needed; the existing Advice and report
      semantics already allow this report projection.
- [x] `specs/skills/quality-skill/reporting.md` - no change needed; generated
      report reading guidance does not need to name the findings index.
- [x] `specs/log.md` and `specs/evaluation/log.md` - durable spec log entries.

**Durable docs / bundled skill runtime**

- [x] `CHANGELOG.md` - release-note entry for the new ranked findings reports.
- [x] `changes/index.md`, `changes/archive/index.md`, and `changes/log.md` -
      Change Case lifecycle.

## Status

`Done`. Implementation, durable specs, release notes, generated schema, tests,
and Change Case lifecycle artifacts are complete. `go test ./...` and
`mise run fmt-md-check` pass.
