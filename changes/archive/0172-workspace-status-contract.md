---
type: Change Case
title: Workspace Status Contract
description: Align status, workspace, project, and repository terminology across CLI output, JSON, specs, docs, and skill guidance.
status: Done
tags: [cli, status, terminology, workspace]
timestamp: 2026-06-29T00:00:00Z
---

# Workspace Status Contract

`qualitymd status` already resolves a QUALITY.md workspace from the selected
model file, but its command help, status specs, JSON type name, README quick
reference, and skill workflow guidance still call the result "project state."
This case makes `workspace` the technical term for the CLI status surface while
preserving `project` for modeled value and use context.

- [Functional spec](0172-workspace-status-contract/spec.md) - what the CLI,
  JSON, specs, docs, and skill guidance must change.
- [Design doc](0172-workspace-status-contract/design.md) - how the status
  implementation exposes workspace status without broadening project wording.

## Motivation

QUALITY.md needs a stable distinction between the project being modeled and the
workspace `qualitymd` operates on. A repository may contain multiple
`QUALITY.md` workspaces, and a project may span repositories. Status is a
mechanical routing surface, so calling it "project state" invites agents and
users to collapse project, repository, and workspace into one concept.

## Scope

Covered:

- `qualitymd status` help, human output, and JSON contract;
- status implementation type names and tests;
- CLI status durable specs and CLI command indexes;
- README/install quick references and terminology guidance;
- `/quality` runtime guidance and skill specs that consume status output;
- domain terminology guidance and agent upfront reading;
- release notes and OKF logs.

Deferred:

- renaming the `status` command itself;
- changing non-status evaluation run path semantics;
- changing `project` value-proposition language;
- changing the agent instruction pointer text created by `qualitymd init`;
- changing historical Change Cases, archived logs, or released changelog entries.

## Affected artifacts

- Code:
  - `internal/status/status.go` - expose a workspace status snapshot and JSON
    workspace block.
  - `internal/cli/status.go` - update command help and human output.
  - `internal/workspace/workspace.go` - align package comments.
  - `internal/cli/status_test.go`, `internal/status/status_test.go`, and
    related CLI tests - update expectations.
- Durable specs:
  - `specs/cli/status.md` - define the workspace-status contract and JSON
    schema version.
  - `specs/cli.md` and `specs/cli/index.md` - align status summaries.
  - `specs/skills/quality-skill/quality-skill.md` and status-consuming guide
    specs - align status-consuming language where relevant.
- Durable docs:
  - `domain.md` - terminology source for project/workspace/repository.
  - `AGENTS.md` - required upfront reading.
  - `README.md` and `install.md` - user-facing quick reference and setup/config
    language.
- Bundled skill:
  - `skills/quality/SKILL.md` and `skills/quality/resources/cli-workflow-conventions.md`
    - align workspace/status guidance.
- Release notes:
  - `CHANGELOG.md` - Unreleased CLI and skill notes.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md`, `docs/log.md`, and `skills/quality/log.md`.

## Status

`Done`. Implemented and archived. `qualitymd status` now presents workspace
status, emits status JSON `schemaVersion: 2` with relative workspace metadata,
and aligns CLI specs, docs, and skill guidance around the project/workspace
distinction.
