---
type: Design Doc
title: Model-relative workspace paths - design
description: How qualitymd workspace resolution moves from repository-relative to model-relative paths.
tags: [workspace, config, evaluation, cli, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Model-relative workspace paths - design

## Context

Answers the [functional spec](spec.md) for change case
[0134](../0134-model-relative-workspace-paths.md). The current resolver in
`internal/workspace` finds the Git repository root from the selected model and
resolves `.quality/`, `.quality/config.yaml`, `evaluationDir`, and `.quality/log`
from that root. Evaluation creation and status already consume this resolver,
but evaluation listing and `--latest` still start from the process current
directory unless handed an evaluation directory override.

The target design keeps one shared resolver, changes its default path base to
the selected model's containing directory, and extends run-discovery commands so
they always have a selected model anchor when they need one.

## Approach

### Add workspace root to the resolver

Extend `workspace.Workspace` with a `WorkspaceRoot PathRef`, where:

- `WorkspaceRoot.Abs` is `filepath.Dir(Model.Abs)`;
- `WorkspaceRoot.Rel` is the slash-normalized path from repository root to that
  directory, or `.` for a root-level model.

Keep `RepoRoot` in the workspace. It remains the containment boundary and the
stable base for repository-relative display when a command needs it.

Replace the current repository-relative path resolver with a helper that accepts
both bases:

```go
ResolveWorkspacePath(repoRoot, workspaceRoot, value) (abs, workspaceRel, repoRel string, error)
```

The helper rejects empty values, absolute paths, and normalized paths whose final
absolute location escapes `repoRoot`. It returns both model-relative and
repository-relative display forms so callers do not recompute them.

### Resolve config and artifacts from workspace root

Resolution order stays the same, but the base changes:

1. Resolve selected model path.
2. Find repository root from the selected model.
3. Set workspace root to the selected model's directory.
4. Parse root `config` from the selected model.
5. Resolve `config`, defaulting to `.quality/config.yaml`, from workspace root.
6. Read config if present.
7. Resolve `.quality/`, `evaluationDir`, `.quality/log/`, and `.quality/logs/`
   from workspace root.

The existing `PathRef` can continue to expose `Rel`, but the implementation
should make the chosen meaning explicit. The safest shape is either:

- rename `Rel` to `WorkspaceRel` and add `RepoRel`, or
- keep `Rel` as the command-facing workspace-relative form and add `RepoRel`
  for diagnostics that need repository context.

The first option is clearer but touches more call sites. The second is a smaller
implementation if tests lock in the public meaning.

### Model-anchor evaluation discovery

Add `Model string` to `evaluationRunFlags` and bind `--model` on every command
that can resolve a run through `--latest`:

- `qualitymd evaluation data set/list/get/verify`;
- `qualitymd evaluation status`;
- `qualitymd evaluation report build`.

Add `--model` to `qualitymd evaluation list` because it discovers history
without a positional run.

Update `evaluation.ResolveRun` and `evaluation.ListRuns` to accept a model path
instead of only `repoRoot`. When `--latest` is used, they resolve the workspace
from that model and list `Workspace.Evaluations`.

When a positional run path is supplied:

- if `--model` is supplied and the run path is relative, resolve it from the
  selected model's workspace root;
- if `--model` is absent, preserve existing filesystem-path behavior for direct
  run paths.

This preserves simple direct-path use while making model-relative operation
available and consistent for nested models.

### Receipts and next actions

Use model-relative run paths in receipts when a selected model is part of the
workspace resolution. For nested models, generated next actions should include
the selected model anchor, for example:

```sh
qualitymd evaluation data set --model packages/api/QUALITY.md .quality/evaluations/0001-full-eval < payloads.json
```

For root-level models this remains visually close to today's output. For nested
models it avoids emitting a `.quality/...` path that a follow-up command would
accidentally resolve from the caller's current directory.

### Status and lint

`qualitymd status <path>` already has a selected model argument. Route
evaluation history through the updated workspace resolver and report the
evaluation history path in the same display form used by evaluation commands.

Keep lint's root `config` rule in the lint package, but route scalar validation
through the updated workspace path validator or a sibling "clean
model-relative" validator. Diagnostics and rule catalog text should say
model-relative.

### Skill and docs

Update durable specs and runtime skill guidance after the code contract is
settled. The skill should carry the same mental model:

- selected model first;
- workspace artifacts beside that model by default;
- Git root as containment;
- `--model` passed whenever history or latest-run lookup is for a non-default
  model path.

## Spec response

- **Workspace resolution** - satisfied by adding `WorkspaceRoot` and resolving
  config/artifacts from it while retaining `RepoRoot` as the safety boundary.
- **Default artifact paths** - satisfied by changing default path resolution
  bases, not by changing path strings.
- **Config and overrides** - satisfied by a shared model-relative validator and
  by preserving existing precedence.
- **Evaluation history and run paths** - satisfied by adding `--model` to
  history/latest commands and resolving relative run paths from workspace root
  when a model is supplied.
- **Status, lint, and skill behavior** - satisfied by routing status through the
  updated workspace and updating diagnostics/specs/skill guidance.

## Alternatives

- **Keep repo-root defaults and add `dataDir`.** Rejected. It gives advanced
  users an escape hatch but leaves the default split convention in place.
- **Make only `.quality/` model-relative but keep `config` repository-relative.**
  Rejected. Users would still need two path bases for one workspace.
- **Require `--model` for every evaluation command, even direct run paths.**
  Rejected. Direct filesystem paths are useful and already work; `--model` is
  needed when the command discovers a run or when the caller wants
  model-relative run-path resolution.
- **Resolve nested paths against the workspace root without a repository
  boundary.** Rejected. The repository boundary is still the right safety limit
  for tooling-owned paths in this Git-oriented CLI.
- **Auto-discover the model from a run folder.** Deferred. It would require
  trusting run-local snapshots or adding metadata that is outside this scoped
  path-base change.

## Trade-offs & risks

- Existing nested-model users who intentionally relied on repo-root `.quality/`
  will see a breaking change. The project is early alpha, and the simpler
  current contract is worth the break.
- Path display can become confusing if some receipts are workspace-relative and
  others are repository-relative. Tests should lock the chosen public display
  form for create/list/status/report receipts and next actions.
- Adding `--model` broadly may make help text busier. The flag is still better
  than implicit cwd/Git-root behavior for latest-run discovery.
- `../shared/config.yaml` from a nested model can be valid if it remains inside
  the repository. This is consistent with the containment rule but should be
  documented so "model-relative" is not mistaken for "must stay under the model
  directory."

## Open questions

- Should `PathRef.Rel` be renamed to `WorkspaceRel` now, or should the
  implementation keep `Rel` and add `RepoRel` to reduce churn?
- Should receipts include both workspace-relative and repository-relative path
  fields under `--json`, or is one command-facing path plus next actions enough?
