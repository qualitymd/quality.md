---
type: Change Case
title: Model-relative workspace paths
description: Make qualitymd workspace config and artifact paths resolve relative to the selected QUALITY.md.
status: Done
tags: [workspace, config, evaluation, cli, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Model-relative workspace paths

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0134-model-relative-workspace-paths/spec.md) - what the case must do.
- [Design doc](0134-model-relative-workspace-paths/design.md) - how it's built, and why.

## Motivation

QUALITY.md already makes Area `source` paths relative to the containing
`QUALITY.md`. Tooling paths currently use a different default orientation:
workspace config, `.quality/`, evaluation history, and logs are resolved from
the Git repository root. That split forces users and agents to remember two
coordinate systems for one selected model file, and it makes nested or multiple
`QUALITY.md` files awkward.

This case makes the selected `QUALITY.md` the default anchor for tooling paths
too. The Git repository root remains a safety boundary, but no longer acts as
the default coordinate system for workspace artifacts.

## Scope

Covered:

- Model-relative workspace resolution for config, data, evaluation, quality-log,
  and workflow-feedback-log paths.
- CLI support for a selected model anchor on evaluation history discovery,
  including `evaluation list` and `--latest`-based commands.
- Path validation, display, receipts, and next-action command alignment for
  nested `QUALITY.md` files.
- Durable specs, docs, and bundled skill wording that currently says
  repository-relative where the new contract is model-relative.

Deferred:

- A migration command for existing repo-root `.quality/` directories.
- Automatic fallback discovery from old repo-root `.quality/` locations.
- A public `dataDir` config key.
- Non-Git workspace containment.

## Affected artifacts

Derived by sweeping for `repository-relative`, `repo root`, workspace resolver
defaults, `.quality/config.yaml`, `.quality/evaluations`, `.quality/log`,
`.quality/logs`, `evaluationDir`, `--evaluation-dir`, `--latest`, `ResolveRun`,
`ListRuns`, and `status` history resolution across code, specs, docs, and skill
content.

**Code**

- [x] `internal/workspace/workspace.go` - introduce the workspace root as the
      directory containing the selected model, resolve workspace paths from it,
      retain repository root as the containment boundary, and expose path forms
      needed by commands.
- [x] `internal/evaluation/path.go` and `internal/evaluation/list.go` - route
      evaluation directory, run listing, `--latest`, and run display through the
      model-relative workspace contract.
- [x] `internal/evaluation/create.go` - create runs under the model-relative
      evaluation directory and emit aligned receipts/next actions.
- [x] `internal/status/status.go` - inspect evaluation history from the selected
      model's workspace root.
- [x] `internal/cli/evaluation.go` - add `--model` where run discovery or
      `--latest` needs a selected model anchor, and align help text.
- [x] `internal/cli/status.go` - keep `status [path]` as the selected model
      anchor and align output/help language.
- [x] `internal/lint/rules.go` and `internal/lint/result.go` - update root
      `config` validation wording from repository-relative to model-relative.
- [x] Focused tests under `internal/evaluation`, `internal/status`,
      `internal/cli`, and `internal/lint` - add nested-model coverage and update
      path expectations.

**Durable specs** (substance in the [functional spec](0134-model-relative-workspace-paths/spec.md))

- [x] `SPECIFICATION.md` - no normative Model schema change; review only to
      keep Area `source` path wording aligned with the broader model-relative
      rationale if needed.
- [x] `specs/cli/evaluation-create.md` - make config and evaluation directory
      resolution model-relative.
- [x] `specs/cli/evaluation-list.md` - make run listing model-anchored.
- [x] `specs/cli/evaluation-status.md` - require `--model` support for
      `--latest` and model-relative run path handling when a model is supplied.
- [x] `specs/cli/evaluation-report.md` - same for report build `--latest` and
      model-relative run path handling when a model is supplied.
- [x] `specs/cli/evaluation-data.md` - same for data commands that accept
      `--latest` or run paths.
- [x] `specs/cli/status.md` - make status evaluation history model-relative.
- [x] `specs/cli/lint-rules.md` - update `invalid-config` from
      repository-relative to model-relative.
- [x] `specs/skills/quality-skill/reporting.md`,
      `specs/skills/quality-skill/quality-log.md`,
      `specs/skills/quality-skill/workflow-feedback-log.md`, and relevant
      workflow specs - align skill artifact paths and feedback logs with the
      selected model's workspace root.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md`, `skills/quality/resources/cli-workflow-conventions.md`,
      and setup/evaluate workflow files - teach model-relative workspace
      resolution and pass `--model` for history/latest commands when needed.
- [x] `install.md` - update config and `evaluationDir` wording.
- [x] `README.md` - review for workspace wording; update if it mentions the
      default artifact base.
- [x] `docs/log.md`, `specs/log.md`, and `changes/log.md` - append-only logs
      are historical; update only with new entries if required by the landing
      work, not by rewriting old entries.

No planned impact: `qualitymd init` scaffold content, JSON Schema generation, or
Evaluation data payload schema.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Implemented
and archived. Workspace config and artifact paths now resolve relative to the
selected `QUALITY.md`; evaluation history/latest commands accept `--model` for
nested workspaces; durable specs, docs, and skill guidance are aligned.
