---
type: Change Case
title: Changelog Directory
description: Rename the model-change quality changelog to .quality/changelog/ and keep .quality/logs/ as a flat workflow-log directory.
status: Done
tags: [skill, workspace, logging, changelog]
timestamp: 2026-06-27T00:00:00Z
---

# Changelog Directory

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0146-changelog-directory/spec.md) - what the case must do.
- [Design doc](0146-changelog-directory/design.md) - how it's built, and why.

## Motivation

The quality workspace currently distinguishes `.quality/log/` from
`.quality/logs/` by singular/plural naming alone. That is too subtle for users
and agents: `.quality/log/` is the curated history of meaningful `QUALITY.md`
model changes, while `.quality/logs/` contains run/process logs such as setup
and evaluate feedback logs. The workspace should make those artifact classes
obvious from their paths and filenames.

## Scope

Covered:

- Rename the quality model-change log directory from `.quality/log/` to
  `.quality/changelog/`.
- Require quality changelog entries to use sortable UTC timestamps plus slugs.
- Keep workflow and process logs under `.quality/logs/` as a flat directory.
- Clarify workflow log filenames so feedback logs are one log kind, not the
  whole category.
- Update code constants, durable specs, runtime skill guidance, docs, examples,
  tests, dogfood data, logs, and release notes.

Deferred:

- A single `.quality/changelog.md` file.
- Subdirectories under `.quality/logs/`.
- Backward-compatibility readers, aliases, dual writes, or migration commands.
- A `qualitymd changelog` CLI surface or configurable `changelogDir`.

## Affected artifacts

Derived by sweeping for `.quality/log/`, `.quality/logs/`, quality changelog,
feedback log, changelog, `DefaultQualityLogDir`, and `logDir`/`changelogDir`.

**Code**

- [x] `internal/workspace/workspace.go` - default quality changelog directory.
- [x] Workspace/status tests if present or affected by the default path.

**Durable specs** (substance in the [functional spec](0146-changelog-directory/spec.md))

- [x] `specs/skills/quality-skill/quality-changelog.md` - renamed component
      contract for quality changelog location/naming.
- [x] `specs/skills/quality-skill/workflow-feedback-log.md` - keep
      `.quality/logs/` flat and workflow-log oriented.
- [x] `specs/skills/quality-skill/quality-skill.md` and related workflow specs -
      align cross-workflow mutation and path references.
- [x] `specs/cli/evaluation-create.md` - update workspace path vocabulary.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` - route model-change history to
      `.quality/changelog/` and keep `.quality/logs/` as flat workflow logs.
- [x] `skills/quality/resources/cli-workflow-conventions.md` - update workspace
      artifact conventions.
- [x] `skills/quality/workflows/*` and `skills/quality/guides/*` - align path,
      link, and mutation wording.
- [x] `README.md`, `npm/quality.md/README.md`, and project `QUALITY.md` where
      active path or artifact wording appears.
- [x] Bundle logs under `docs/`, `specs/`, and `skills/quality/`.
- [x] `CHANGELOG.md` - record the breaking workspace path rename.

**Dogfood data**

- [x] Move `.quality/log/` entries to `.quality/changelog/` with timestamped
      filenames.

## Status

`Done`. Implemented and archived after workspace defaults, runtime skill
guidance, durable specs, dogfood data, logs, release notes, and verification were
updated. `go test ./...`, `qualitymd status --json`, `mise run fmt-md-check`,
and `mise run check` pass.
