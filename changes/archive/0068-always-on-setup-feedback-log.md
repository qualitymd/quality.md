---
type: Change Case
title: Always-on setup feedback log
description: Make /quality setup create a workflow feedback log at run start, update it as the workflow progresses, and finalize it at close.
status: Done
tags: [skill, setup, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Always-on setup feedback log

A **Change Case** to make the setup feedback log an always-created, progressively
updated run artifact instead of an optional close-only artifact. The workflow
feedback log remains a skill-written, local Markdown artifact under
`.quality/logs/`; this change makes its presence unambiguous and preserves useful
partial feedback when a setup run is interrupted.

Detail lives in:

- [Functional spec](0068-always-on-setup-feedback-log/spec.md) - what the change
  must do.
- [Design doc](0068-always-on-setup-feedback-log/design.md) - how it is shaped,
  and why.

## Motivation

The current setup feedback-log contract allows `setup` to omit the log when
nothing notable occurred. That makes absence ambiguous: it can mean a clean run,
an older skill version, an interrupted run, or a missed close-step write. The log
is intended to improve `/quality` workflows from real runs, but a close-only,
optional artifact loses signal precisely when setup is interrupted or when a
"boring" run would be useful as a baseline.

Making the log always-on resolves that ambiguity. Creating it early records the
environment and run boundary while they are fresh. Updating it through the
workflow gives friction, latency, retry, and UX/AX notes a durable place to land
without turning the final user-facing setup summary into an internal feedback
report. Finalizing it at close lets a clean run say so explicitly.

## Scope

Covered:

- Change `setup` from optional close-step feedback-log authoring to mandatory
  run-start creation, progressive update, and close finalization.
- Define stable feedback-log metadata and body sections for the always-on format.
- Define status values for in-progress, completed, interrupted, and failed runs.
- Preserve the existing `.quality/logs/<timestamp>-setup-feedback-log.md`
  location and local/no-transmission posture.
- Preserve the boundary that feedback logs describe workflow experience, not
  model content, evaluation findings, or quality-log history.
- Update the runtime skill and durable specs to match.

Deferred / non-goals:

- No QUALITY.md format change.
- No CLI or Go code change; the artifact remains skill-written.
- No automatic external transmission, upload endpoint, telemetry service, or
  consent gate.
- No structured machine parser, validator, or report renderer for feedback logs.
- No change to evaluation's per-run `debug-log.md` contract.
- No requirement that evaluate or update adopt always-on feedback logging in this
  case, though the shared workflow feedback-log contract should not block later
  adoption.

## Affected artifacts

### Code

- [x] **None.** The artifact remains skill-written; no `qualitymd`
      command, flag, config field, or Go package is planned.

### Durable specs

- [x] `specs/skills/quality-skill/workflows/setup/feedback-log.md` - revise the
      artifact contract from optional close-step creation to mandatory run-start
      creation, progressive updates, and close finalization; define the stable
      frontmatter fields and body sections.
- [x] `specs/skills/quality-skill/workflows/setup.md` - update setup's workflow
      steps, mutation surface, completion behavior, and feedback-log close
      behavior to match the always-on contract.
- [x] `specs/skills/quality-skill/quality-skill.md` - update the shared workflow
      feedback-log description so it no longer implies optional, close-only
      setup logging.
- [x] `specs/log.md` and
      `specs/skills/quality-skill/workflows/log.md` - record the durable spec
      updates when implementation lands.

### Durable docs

- [x] **None.** No public durable docs changed; runtime skill and specs carry the
      behavior.

### Bundled skill

- [x] `skills/quality/SKILL.md` - update hard rules and run-frame wording so
      setup's feedback log is expected rather than optional.
- [x] `skills/quality/workflows/setup.md` - create the log during preflight,
      update it during notable workflow steps, and finalize it at close.
- [x] `skills/quality/resources/cli-quick-reference.md` - update wording if it
      says workflow feedback logs are only created on demand or optional.

### Release

- [x] `CHANGELOG.md` - add the user-facing `/quality Skill` note under
      `Unreleased`.

## Children

- [Functional spec](0068-always-on-setup-feedback-log/spec.md) - required
  behavior for always-on setup feedback logs.
- [Design doc](0068-always-on-setup-feedback-log/design.md) - how the always-on
  log lifecycle is shaped, and why.

## Status

`Done`. Implemented skill-only with no CLI/Go change. Durable specs, the runtime
skill, quick reference, spec logs, and `CHANGELOG.md` are updated; no public docs
changed. This case supersedes the optional close-only setup behavior introduced
by [0066 - Setup feedback log](0066-setup-feedback-log.md) while preserving that
artifact's local, skill-only, no-transmission boundary. Moved to
[`archive/`](index.md).
