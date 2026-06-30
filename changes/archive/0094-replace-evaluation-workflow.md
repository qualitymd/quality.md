---
type: Change Case
title: Replace evaluation workflow
description: Replace the current evaluation workflow with the Evaluation v2 protocol, records, reports, CLI surface, and skill runtime contract.
status: Done
tags: [evaluation, workflow, cli, skill, reports]
timestamp: 2026-06-25T00:00:00Z
---

# Replace evaluation workflow

This change case replaces the current QUALITY.md evaluation workflow with the
Evaluation v2 protocol captured in
[`evaluation-v2-sketch.md`](../../evaluation-v2-sketch.md). The new workflow treats
evaluation as an agent-orchestrated judgment protocol whose durable outputs are
structured routine JSON under `data/` plus deterministic Markdown reports.

Children:

- [Functional spec](0094-replace-evaluation-workflow/spec.md) - what the change must do.
- [Design doc](0094-replace-evaluation-workflow/design.md) - how the change is built.

## Motivation

The current evaluation workflow grew around CLI-managed assessment, analysis,
recommendation, and report records. That shape makes the record layout the mental
model for evaluation, even though the important work is judgment: framing a
Requirement before evidence review, assessing evidence, rating the assessment,
and synthesizing Factors and Areas without hiding the drivers that bind a
rating.

Evaluation v2 makes those judgment moves explicit. It records each routine output
as structured data, keeps the CLI mechanical, and makes deterministic reporting a
projection over completed structured results. The goal is a clearer protocol for
agents, stronger auditability, better resumability, and a report tree that lets a
reader navigate from the root Area down to Area, Factor, and Requirement detail.

## Scope

This case covers the replacement evaluation protocol, JSON routine output
contracts, run-folder data layout, deterministic report generation, CLI command
surface, agent-agnostic orchestration model, and `/quality` skill runtime
instructions needed to run Evaluation v2.

This case also covers retiring or replacing the current evaluation record and
report contracts where they conflict with Evaluation v2.

Deferred:

- Recommendation generation and recommendation follow-up changes beyond removing
  recommendation generation from the v0 evaluation protocol.
- QC routines around judgment outputs.
- Schema migration or compatibility transforms for old evaluation runs.
- Custom synthesis policy sources beyond the v0 protocol defaults.
- `qualitymd evaluation data schema <kind>`.

## Affected artifacts

### Format spec

- [x] `SPECIFICATION.md` - replace the current evaluation/report semantics with
      Evaluation v2 semantics and deferred recommendation posture.

### Durable specs

- [x] `specs/evaluation-v2/` - add the durable parent spec folder for the new
      protocol, routines, records, and report contracts.
- [x] `specs/evaluation-records.md` and `specs/evaluation-records/` - replace or
      retire the current assessment/analysis/recommendation/run-folder/report
      record contracts that conflict with Evaluation v2.
- [x] `specs/reports/` - replace the current `report.md`,
      `report-summary.md`, and `report.json` contracts with the Evaluation v2
      report tree contract, or retire artifacts no longer produced.
- [x] `specs/cli.md` - align the shared CLI contract with Evaluation v2
      artifact JSON behavior where needed.
- [x] `specs/cli/evaluation-*.md` - replace the current evaluation command
      surface with `evaluation create`, `evaluation data set`,
      `evaluation data list/get/kinds/example`, `evaluation status`,
      `evaluation list`, and
      `evaluation report build`.
- [x] `specs/skills/quality-skill/evaluation.md` and
      `specs/skills/quality-skill/reporting.md` - replace the durable skill
      evaluation and reporting contracts with Evaluation v2.
- [x] `specs/skills/quality-skill/examples/` - update or replace examples that
      depend on the current run layout, records, reports, or recommendations.

### Code

- [x] `cmd/qualitymd/` and `internal/cli/` - implement the new evaluation command
      tree and remove or redirect conflicting current commands.
- [x] `internal/evaluation/`, `internal/status/`, and report generation code -
      implement Evaluation v2 run creation, data persistence, status/gap
      inspection, output assembly, and deterministic report rendering.
- [x] `internal/model/`, `internal/lint/`, and related traversal helpers -
      expose any model identity, ordering, and reference helpers needed by
      Evaluation v2. No new helper API was needed for this slice; v2 paths use
      structural IDs and deterministic lexical order.
- [x] Tests under `internal/**` - cover data validation, path derivation,
      status/report gaps, deterministic reporting, and CLI behavior.

### Bundled skill and runtime guidance

- [x] `skills/quality/SKILL.md` - update CLI operating rules and artifact
      contract for Evaluation v2.
- [x] `skills/quality/workflows/evaluate.md` - replace the current evaluate
      workflow with the Evaluation v2 protocol and orchestration rules.
- [x] `skills/quality/resources/cli-quick-reference.md` and
      `skills/quality/resources/output-policy.md` - update command usage and
      output consumption rules.
- [x] Optional Evaluation v2 runtime instruction split under
      `skills/quality/workflows/evaluation-v2/` if implementation uses the
      mirrored spec structure. Deferred; the current workflow remains in
      `skills/quality/workflows/evaluate.md`.

### Public docs and examples

- [x] `README.md` - update evaluation and reporting prose when the workflow
      replacement lands.
- [x] `docs/guides/cli-design.md` and `specs/cli.md` - already have local edits
      for artifact JSON behavior; reconcile them with this case.
- [x] Other docs found by searching for current evaluation command names,
      record folders, report artifacts, recommendations, `report-summary.md`,
      `report.json`, `assessment add`, `analysis set`, and related examples.
      Remaining old artifact mentions are retained only in superseded legacy
      specs or explicitly legacy examples.

## Status

`Done`. Implemented, verified, and archived for release.
