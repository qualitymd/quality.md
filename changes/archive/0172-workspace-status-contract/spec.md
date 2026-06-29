---
type: Functional Specification
title: Workspace Status Contract - functional spec
description: Requirements for aligning project, workspace, and repository terminology in status surfaces.
tags: [cli, status, terminology, workspace]
timestamp: 2026-06-29T00:00:00Z
---

# Workspace Status Contract - functional spec

Companion to
[Workspace Status Contract](../0172-workspace-status-contract.md). This spec
states the delta for `qualitymd status` and the documentation that teaches the
project/workspace/repository distinction. The durable source of truth is
absorbed into [`qualitymd status`](../../../specs/cli/status.md), the
[`qualitymd CLI`](../../../specs/cli.md), and the `/quality` skill contracts.

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

`qualitymd status` is a mechanical routing surface. It validates the selected
model file, resolves workspace-owned paths, summarizes evaluation history, and
returns next actions. That is workspace status, not a judgment about the modeled
project. Calling the output "project state" is especially misleading in
repositories with nested `QUALITY.md` files or projects whose modeled boundary
is broader than the containing directory.

The change should teach the distinction without weakening the value proposition:
`project` remains the right word for what QUALITY.md helps users model and
improve, while `workspace` is the right word for what the CLI resolves,
inspects, reports, and mutates.

## Scope

Covered:

- status command help and human output;
- status JSON contract and schema version;
- status implementation terminology;
- status specs and CLI command summaries;
- status-consuming `/quality` runtime and skill-spec guidance;
- README/install/status quick references;
- domain and agent terminology guidance;
- release notes, tests, and OKF logs.

Deferred:

- renaming `qualitymd status`;
- changing evaluation run resolution outside the status JSON metadata;
- changing non-status uses of `project` that describe modeled value, setup
  intent, or root-area meaning;
- changing historical archive/log prose.

## Requirements

1. `qualitymd status` human help **MUST** describe the command as a QUALITY.md
   workspace status surface, not a project-state surface.

   > Rationale: Help is the most visible CLI teaching surface, and `status`
   > operates on the selected workspace's mechanical state. Durable spec: modify
   > `specs/cli/status.md` and `specs/cli.md`.

2. `qualitymd status` human output **MUST** identify the result as workspace
   status and show the resolved workspace root when a workspace is available.

   > Rationale: The default output should make nested workspaces legible without
   > forcing users or agents to parse JSON. Durable spec: modify
   > `specs/cli/status.md`.

3. `qualitymd status --json` **MUST** bump the status JSON `schemaVersion` and
   include a top-level `workspace` object when the selected model resolves to a
   workspace.

   > Rationale: Adding first-class workspace metadata changes the machine
   > contract in a way agents may branch on; a schema-version bump makes the
   > change explicit. Durable spec: modify `specs/cli/status.md`.

4. The status JSON `workspace` object **MUST** expose only stable,
   non-absolute, workspace/repository-relative path facts: workspace root,
   selected model path, config path and presence, quality data directory,
   evaluation directory, quality changelog directory, and workflow log
   directory.

   > Rationale: Agents need routing paths but not host-specific absolute paths.
   > The repository remains a containment boundary, not the project. Durable
   > spec: modify `specs/cli/status.md`.

5. Status JSON **MUST NOT** rename or remove the existing top-level `path`,
   `model`, `evaluations`, `readiness`, or `nextActions` fields as part of this
   case.

   > Rationale: The workspace block should clarify the contract without making
   > unrelated consumers rewrite stable status fields. Durable spec: modify
   > `specs/cli/status.md`.

6. Implementation-visible status type names and comments **SHOULD** use
   workspace status terminology where they name the status snapshot contract.

   > Rationale: Internal names shape future user-facing prose and tests, but this
   > requirement is not a public JSON-field rename. Durable spec: none.

7. Durable docs and skill guidance that instruct agents to consume
   `qualitymd status` **MUST** call that result workspace status or structured
   status output, not project state.

   > Rationale: The agent is the main user interface; stale skill wording would
   > keep recreating the ambiguity this case removes. Durable spec: modify
   > `specs/skills/quality-skill/quality-skill.md` and related skill workflow
   > specs.

8. User-facing docs **MUST** preserve `project` for the modeled value
   proposition and setup intent, and **MUST** use `workspace` only for technical
   CLI/path/artifact behavior.

   > Rationale: Replacing every `project` with `workspace` would make QUALITY.md
   > sound like a filesystem tool rather than a project-quality model. Durable
   > spec: none.

9. This change **MUST NOT** edit historical Change Case prose, archived log
   entries, or released changelog entries solely to update terminology.

   > Rationale: Historical records describe past contracts and should not be
   > rewritten for current terminology cleanup. Durable spec: none.

## Acceptance criteria

- `qualitymd status --help` says workspace status, not project-state.
- Root `qualitymd --help` lists status as a workspace status command.
- Plain `qualitymd status` prints a workspace-status heading and workspace root
  when the selected model is valid.
- `qualitymd status --json` emits `schemaVersion: 2` and a top-level
  `workspace` object for valid models.
- The JSON workspace object contains only relative path facts and config
  presence.
- Existing status JSON fields remain present.
- Nested `QUALITY.md` workspaces report the nested workspace root and
  workspace-relative evaluation paths.
- README/install/skill guidance consistently use workspace for status/config
  mechanics and project for modeled value.
- Focused Go tests and repository checks pass.

## Durable spec changes

### To add

None

### To modify

- `specs/cli/status.md` - define status as workspace status, bump JSON
  `schemaVersion`, require the `workspace` object, and update output language
  (requirements 1-5).
- `specs/cli.md` - align the status command summary with workspace-status
  terminology (requirement 1).
- `specs/cli/index.md` - align the status command summary with workspace-status
  terminology (requirement 1).
- `specs/skills/quality-skill/quality-skill.md` and status-consuming skill
  workflow/guide specs - align agent guidance that consumes status output
  (requirement 7).

### To rename

None

### To delete

None
