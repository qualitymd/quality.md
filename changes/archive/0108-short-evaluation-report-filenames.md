---
type: Change Case
title: Short evaluation report filenames
description: Keep the root Evaluation v2 report at report.md while renaming descendant Markdown reports to short subject-aware filenames for clearer editor and browser tabs.
status: Done
tags: [cli, evaluation, reports, navigation]
timestamp: 2026-06-26T00:00:00Z
---

# Short evaluation report filenames

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0108-short-evaluation-report-filenames/spec.md) - what the
  change must do.
- [Design doc](0108-short-evaluation-report-filenames/design.md) - how it is
  built, and why.

## Motivation

Evaluation v2 currently writes every human Markdown report as `report.md`.
Directory paths preserve the subject identity, but editor and browser tabs show
only the filename in many common views. A reviewer who opens several Area,
Factor, and Requirement reports at once can no longer tell which tab is which
without switching through them.

The root Area can stay the conventional run entrypoint at `report.md` because
there is only one root report per run. Descendant Area, Factor, and Requirement
reports should keep deterministic structural paths while using short local
filenames that name the subject kind and local structural identifier.

## Scope

Covered:

- generated Evaluation v2 human Markdown report filenames for Area, Factor, and
  Requirement reports;
- generated links, breadcrumbs, table links, report refs, and report output
  indexes that point to those Markdown reports;
- CLI report-path derivation and tests that assert generated report paths; and
- durable Evaluation v2 report layout and `/quality` reporting contracts.

Deferred / non-goals:

- no change to Area, Factor, or Requirement directory names;
- no use of mutable display titles, natural labels, or rendered human labels in
  persisted report paths;
- no compatibility copies using the old descendant `report.md` filenames;
- no change to report content, ratings, findings, analysis, evidence, limits, or
  navigation labels beyond link targets;
- no change to structured routine data paths under `data/`; and
- no migration of existing completed evaluation runs.

## Affected artifacts

### Code

- [x] `internal/evaluation/report_v2.go` - derive and write short subject-aware
      Markdown report filenames, and update all generated report links and
      `EvaluationOutputResult` report refs.
- [x] `internal/evaluation/data.go` - reviewed example/fallback report refs;
      only the still-correct root Area `report.md` ref remains.
- [x] `internal/evaluation/evaluation_test.go` and any affected CLI tests -
      assert the new generated paths and links.
- [x] `dprint.json` - keep generated descendant Markdown reports excluded from
      repository Markdown formatting.

### Format spec

- [x] None - `SPECIFICATION.md` does not govern Evaluation v2 generated report
      filenames. (Deliberate.)

### Durable specs

- [x] `specs/evaluation-v2/records/data-layout.md` - replace the `report.md`
      report-tree contract with short subject-aware filenames.
- [x] `specs/evaluation-v2/reports/report-tree.md` - require generated report
      links and navigation to use the new filenames while preserving labels and
      hierarchy.
- [x] `specs/skills/quality-skill/reporting.md` - align the `/quality`
      reporting artifact contract with the new generated filenames.
- [x] `specs/cli/status.md` - remove stale wording that names `report.md` bodies
      specifically while preserving the rule that status does not read generated
      Markdown reports.

### Durable docs / bundled skill

- [x] `evaluation-v2-sketch.md` - reconcile report-tree examples and any
      remaining `report.md` path examples for Evaluation v2.
- [x] `skills/quality/workflows/evaluate.md` - update the run-frame artifact
      summary that currently names `report.md`.

### Suggested new durable specs

- None.

## Status

`Done`. Implemented, verified, and archived. See the
[status lifecycle](../index.md#status-lifecycle).
