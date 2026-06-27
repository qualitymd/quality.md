---
type: Functional Specification
title: Model-relative workspace paths - functional spec
description: What must change so qualitymd workspace paths resolve relative to the selected QUALITY.md.
tags: [workspace, config, evaluation, cli, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Model-relative workspace paths - functional spec

Companion to the
[Model-relative workspace paths](../0134-model-relative-workspace-paths.md)
change case. This spec states *what* the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

QUALITY.md already uses the containing `QUALITY.md` as the base for Area
`source` paths. Tooling paths should follow the same orientation: users and
agents select a model file, then config, `.quality/`, evaluation history, and
logs are found beside that model by default. The Git repository root still
matters, but only as a containment boundary that prevents tooling paths from
escaping the repository.

This keeps the common root-level case visually unchanged while making nested or
multiple `QUALITY.md` files predictable.

## Scope

Covered: workspace resolution, config and artifact path validation, evaluation
run creation/listing/latest resolution, status history inspection, run-path
display and next actions, lint wording for root `config`, durable specs, user
docs, and bundled skill runtime instructions.

Deferred: migration from old repo-root `.quality/`, fallback readers for old
locations, a public `dataDir` key, non-Git containment, and changes to the
normative Model schema or Evaluation data payload schemas.

## Terminology

- **Selected model**: the explicit `QUALITY.md` path supplied to a command, or
  `QUALITY.md` in the current working directory when omitted.
- **Workspace root**: the directory containing the selected model.
- **Repository root**: the nearest Git root found from the selected model. It is
  the safety boundary for path validation, not the default path base.
- **Model-relative path**: a relative path resolved from the workspace root.

## Requirements

### Workspace resolution

- Tools that need workspace paths **MUST** resolve a QUALITY.md workspace from
  the selected model. The workspace **MUST** include the selected model path,
  workspace root, repository root, config file, quality data directory,
  evaluation directory, quality log directory, and workflow feedback-log
  directory.

  > Rationale: the workspace needs both anchors: workspace root for user-facing
  > path orientation, repository root for containment.
  >
  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/status.md`, and the relevant skill specs to add workspace root to
  > the shared workspace contract.

- When no model file is supplied to a command that accepts a selected model, the
  selected model **MUST** default to `QUALITY.md` in the current working
  directory.

  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/status.md`, `specs/cli/evaluation-list.md`,
  > `specs/cli/evaluation-status.md`, `specs/cli/evaluation-report.md`, and
  > `specs/cli/evaluation-data.md`.

- Workspace path resolution **MUST** find the repository root from the selected
  model before accepting config or artifact paths, and **MUST** reject any
  model-relative path that escapes the repository after normalization.

  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/status.md`, and `specs/cli/lint-rules.md`.

### Default artifact paths

- The default quality data directory **MUST** be `.quality/` under the workspace
  root.

  > Durable spec: modify `specs/cli/evaluation-create.md` and relevant skill
  > specs.

- The default config file **MUST** be `.quality/config.yaml` under the workspace
  root.

  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/status.md`, and relevant skill specs.

- The default evaluation directory **MUST** be `.quality/evaluations/` under the
  workspace root.

  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/status.md`, `specs/cli/evaluation-list.md`,
  > `specs/cli/evaluation-status.md`, `specs/cli/evaluation-report.md`,
  > `specs/cli/evaluation-data.md`, and relevant skill specs.

- The default quality log directory **MUST** be `.quality/log/` under the
  workspace root, and workflow feedback logs **MUST** default to `.quality/logs/`
  under the workspace root.

  > Durable spec: modify `specs/skills/quality-skill/quality-log.md`,
  > `specs/skills/quality-skill/workflow-feedback-log.md`, and relevant skill
  > workflow specs.

### Config and overrides

- The selected model's root `config` frontmatter key **MAY** point to the
  workspace config file. When present, `config` **MUST** be interpreted as a
  model-relative path, **MUST** be a non-empty scalar, **MUST NOT** be absolute,
  and **MUST NOT** escape the repository after normalization.

  > Rationale: keeping `config` model-relative preserves the one-rule mental
  > model; leaving it repository-relative would retain the old split.
  >
  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/status.md`, and `specs/cli/lint-rules.md`.

- `evaluationDir` in the resolved config file **MUST** be interpreted as a
  model-relative path, **MUST NOT** be absolute, and **MUST NOT** escape the
  repository after normalization.

  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/status.md`, and relevant evaluation command specs.

- `--evaluation-dir` **MUST** be interpreted as a model-relative path when used
  with a selected model, **MUST NOT** be absolute, and **MUST NOT** escape the
  repository after normalization.

  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/evaluation-list.md`, `specs/cli/evaluation-status.md`,
  > `specs/cli/evaluation-report.md`, and `specs/cli/evaluation-data.md`.

- Evaluation directory precedence **MUST** remain: explicit `--evaluation-dir`,
  then `evaluationDir` in the resolved config file, then
  `.quality/evaluations/`.

  > Durable spec: modify `specs/cli/evaluation-create.md` and
  > `specs/cli/status.md`.

### Evaluation history and run paths

- `qualitymd evaluation create --model <model>` **MUST** create numbered run
  folders under the selected model's resolved evaluation directory.

  > Durable spec: modify `specs/cli/evaluation-create.md`.

- Commands that discover evaluation history without a positional run path,
  including `qualitymd evaluation list` and `--latest` on evaluation data,
  status, and report commands, **MUST** support `--model <model>` as the
  selected model anchor.

  > Rationale: without a selected model, nested workspaces would fall back to the
  > process current directory or repository root and recreate the split
  > convention this change removes.
  >
  > Durable spec: modify `specs/cli/evaluation-list.md`,
  > `specs/cli/evaluation-status.md`, `specs/cli/evaluation-report.md`, and
  > `specs/cli/evaluation-data.md`.

- When a command receives both `--model <model>` and a relative positional
  `<run>` path, it **MUST** resolve that run path relative to the selected
  model's workspace root. When no `--model` is supplied and a positional `<run>`
  path is supplied, existing filesystem-path behavior **MAY** be preserved.

  > Durable spec: modify `specs/cli/evaluation-status.md`,
  > `specs/cli/evaluation-report.md`, and `specs/cli/evaluation-data.md`.

- Run creation receipts, status/list JSON path fields, report-build receipts,
  and generated next-action commands **SHOULD** use paths and flags that are
  directly reusable by follow-up commands for the same selected model.

  > Durable spec: modify `specs/cli/evaluation-create.md`,
  > `specs/cli/evaluation-list.md`, `specs/cli/evaluation-status.md`,
  > `specs/cli/evaluation-report.md`, and `specs/cli/evaluation-data.md`.

### Status, lint, and skill behavior

- `qualitymd status <model>` **MUST** inspect evaluation history from the
  selected model's resolved evaluation directory.

  > Durable spec: modify `specs/cli/status.md`.

- `qualitymd lint` **MUST** accept root `config` only when it satisfies the
  model-relative config-path requirements above, and its diagnostics/rule
  catalog **MUST** describe the pointer as model-relative rather than
  repository-relative.

  > Durable spec: modify `specs/cli/lint-rules.md`.

- The `/quality` skill **MUST** describe workspace artifacts as relative to the
  selected `QUALITY.md`, and **MUST** pass a selected model anchor when invoking
  CLI commands that discover history or latest runs for non-default model paths.

  > Durable spec: modify relevant `specs/skills/quality-skill/` specs.

## Acceptance Criteria

- With `packages/api/QUALITY.md` and no config file,
  `qualitymd evaluation create --model packages/api/QUALITY.md --json` creates a
  run under `packages/api/.quality/evaluations/`.
- With `packages/api/QUALITY.md` and no config file,
  `qualitymd status packages/api/QUALITY.md --json` reports evaluation history
  from `packages/api/.quality/evaluations`, not repo-root `.quality/evaluations`.
- With `packages/api/QUALITY.md` containing
  `config: .quality/custom-config.yaml`, the CLI reads
  `packages/api/.quality/custom-config.yaml`.
- With that config containing `evaluationDir: tmp/evals`, evaluation creation,
  status history, `evaluation list --model packages/api/QUALITY.md`, and
  `--latest --model packages/api/QUALITY.md` commands use
  `packages/api/tmp/evals`.
- With both `--evaluation-dir custom/evals` and config
  `evaluationDir: tmp/evals`, the command uses `packages/api/custom/evals`.
- With `config: /tmp/config.yaml`, lint reports an `invalid-config` error.
- With `config: ../../outside.yaml` from a nested model where normalization
  escapes the Git repository, lint reports an `invalid-config` error.
- With `config: ../shared/quality-config.yaml` from a nested model where
  normalization remains inside the Git repository, lint accepts the pointer and
  resolves it from the selected model's workspace root.
- `qualitymd evaluation list --model packages/api/QUALITY.md --json` lists runs
  from that model's workspace, while the same command without `--model` defaults
  to `QUALITY.md` in the current working directory.
- `qualitymd evaluation status --latest --model packages/api/QUALITY.md`,
  `qualitymd evaluation report build --latest --model packages/api/QUALITY.md`,
  and `qualitymd evaluation data list --latest --model packages/api/QUALITY.md`
  resolve latest from the same model-relative evaluation directory.
- Given `--model packages/api/QUALITY.md` and relative run path
  `.quality/evaluations/0001-full-eval`, evaluation data/status/report commands
  resolve the run under `packages/api/.quality/evaluations/0001-full-eval`.
- Root-level `QUALITY.md` behavior remains visibly unchanged because the
  workspace root and repository root are the same directory.
- Runtime skill instructions and user docs state one rule: relative workspace
  tooling paths are relative to the selected `QUALITY.md`; the Git root is only
  the containment boundary.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) - the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/cli/evaluation-create.md` - define workspace root, model-relative
  config/evaluation path resolution, and model-relative run creation receipts.
- `specs/cli/evaluation-list.md` - require `--model` support and model-relative
  history discovery.
- `specs/cli/evaluation-status.md` - require `--model` support for `--latest`
  and model-relative relative-run resolution when a model is supplied.
- `specs/cli/evaluation-report.md` - require the same selected-model behavior
  for report build.
- `specs/cli/evaluation-data.md` - require the same selected-model behavior for
  data commands that accept run paths or `--latest`.
- `specs/cli/status.md` - resolve status evaluation history from the selected
  model's workspace root.
- `specs/cli/lint-rules.md` - update root `config` validation from
  repository-relative to model-relative.
- `specs/skills/quality-skill/reporting.md` - align evaluation artifact
  defaults and selected-model CLI invocations.
- `specs/skills/quality-skill/quality-log.md` - make `.quality/log/`
  model-relative.
- `specs/skills/quality-skill/workflow-feedback-log.md` - make `.quality/logs/`
  model-relative.
- Relevant `specs/skills/quality-skill/workflows/*` specs - align setup and
  evaluate workflow artifact paths and run-frame wording.

### To rename

None.

### To delete

None.
