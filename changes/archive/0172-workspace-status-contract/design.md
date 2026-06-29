---
type: Design Doc
title: Workspace Status Contract - design doc
description: Design for exposing workspace status metadata and aligning terminology.
tags: [cli, status, terminology, workspace]
timestamp: 2026-06-29T00:00:00Z
---

# Workspace Status Contract - design doc

Design for
[Workspace Status Contract](../0172-workspace-status-contract.md) and its
[functional spec](spec.md).

## Context

`internal/workspace` already resolves the selected model file into the
workspace root, repository-relative root, config file, `.quality/` data
directory, evaluation directory, changelog directory, and feedback-log directory.
`qualitymd status` currently uses that resolver indirectly when it inspects
evaluation history, but the status snapshot does not expose those workspace
facts and the public prose still calls the result project state.

## Approach

Resolve the workspace once after the selected model is lint-valid. Convert the
existing `workspace.Workspace` into a small status-owned `WorkspaceStatus`
struct that carries only stable relative path facts:

- repository-relative workspace root;
- workspace-relative selected model path;
- workspace-relative config path and config presence;
- workspace-relative quality data directory;
- workspace-relative evaluation, changelog, and workflow-log directories.

Keep absolute paths private to command execution. The existing
`evaluations.path` field remains the workspace-relative evaluation directory so
existing status consumers still have the direct history path they already used.
The new `workspace` block is additional context, so existing `path`, `model`,
`evaluations`, `readiness`, and `nextActions` remain.

Bump status `schemaVersion` to `2` because agents may branch on the presence of
the new `workspace` block. Missing or invalid models continue to omit workspace
metadata because a trustworthy workspace cannot always be resolved before the
model and config are valid.

Rename the implementation contract type from `ProjectSnapshot` to
`WorkspaceSnapshot`, update package comments, and adjust CLI rendering. Human
status output should say `Workspace Status` and include `Workspace: <root>` when
metadata is available. The model validity line becomes `Model file:
<workspace-model-path>: <state>` for valid workspaces and falls back to the
input path when no workspace metadata is available.

Update durable specs, README/install guidance, and skill workflow conventions in
the same pass. Preserve project language in root CLI value-proposition text,
setup guidance, model body guidance, and agent-instruction pointers.

## Spec response

- Status help and command summaries change to workspace-status wording.
- The status snapshot gains first-class workspace metadata while keeping
  existing top-level JSON fields.
- `schemaVersion` increments to `2`.
- Plain output identifies workspace status and the resolved workspace root.
- Skill guidance that consumes status output now says workspace status or
  structured status output instead of project state.
- Project wording remains in value-proposition and setup contexts.

## Alternatives

- **Only update prose.** Rejected because agents consuming JSON would still need
  to infer workspace facts from `path` and `evaluations.path`, preserving the
  ambiguity for the main automation surface.
- **Rename top-level `path` to `modelPath`.** Rejected because it would force an
  unrelated consumer rewrite. The new `workspace.model` field clarifies the
  selected model path without removing the existing field.
- **Expose absolute workspace paths.** Rejected because status JSON should be
  deterministic and safe to share; absolute paths are host-specific.
- **Resolve workspace for missing or invalid models.** Rejected for now because
  config lint failures and missing files can make workspace facts untrustworthy.
  The readiness fields already identify those states.

## Trade-offs & risks

Adding a `workspace` block is a machine-contract change. The schema-version bump
mitigates that, and the old top-level fields remain available.

The human output changes a compact line that tests and users may recognize. That
is acceptable because the command is early alpha and the new wording teaches the
contract more accurately.

## Open questions

None.
