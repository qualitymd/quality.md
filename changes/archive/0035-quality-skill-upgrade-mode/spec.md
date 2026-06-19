---
type: Functional Specification
title: /quality upgrade mode - functional spec
description: Define a /quality skill mode that plans and orchestrates compatible skill and CLI upgrades.
tags: [skill, upgrade, versioning]
timestamp: 2026-06-19T00:00:00Z
---

# /quality upgrade mode - functional spec

Companion to [/quality upgrade mode](../0035-quality-skill-upgrade-mode.md).
This spec defines an upgrade orchestration mode for keeping the `/quality` skill
and `qualitymd` CLI compatible.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The `/quality` skill and `qualitymd` CLI are distributed separately but must be
kept compatible. Skill release metadata records the skill version and required
CLI range, while CLI upgrade work gives `qualitymd` a structured way to inspect
and update the CLI. Users still need one agent-facing workflow that can answer:
what versions are installed, whether the skill and CLI are compatible, what
should be upgraded, and what happens after an upgrade.

`/quality upgrade` exists to orchestrate that pair. The skill owns the diagnosis
and coordination; the CLI and skill installer own mutation.

## Scope

Covered: a new `/quality upgrade` mode, installed skill and CLI version checks,
compatibility diagnosis, upgrade planning for both artifacts, explicit
confirmation before mutation, delegation to CLI and skill-installer upgrade
mechanics, post-upgrade verification, and restart/reload guidance.

Deferred / non-goals: no hand-editing installed skill files, no custom skill
package manager inside `qualitymd`, no background upgrade checks during normal
quality workflows, no automatic mutation without confirmation, and no assumption
that the currently running skill can reload itself after upgrade.

## Requirements

The `/quality` skill **MUST** recognize `upgrade` as a first-class mode.

`/quality upgrade` **MUST** inspect the current `/quality` skill release
metadata, including `metadata.version` and `metadata.requires-qualitymd-cli`,
when that metadata is available.

`/quality upgrade` **MUST** inspect the installed `qualitymd` CLI version before
building an upgrade plan.

When `qualitymd version --json` is available, `/quality upgrade` **SHOULD** use
it as the source for CLI version facts. Otherwise it **MAY** fall back to
`qualitymd --version`.

`/quality upgrade` **MUST** diagnose at least these states:

- skill missing release metadata;
- CLI missing;
- CLI present but outside the skill's required range;
- CLI present and inside the skill's required range;
- newer CLI available;
- newer `/quality` skill available, when installer support can report it;
- skill upgrade applied but current agent session still uses the previous
  loaded skill instructions.

`/quality upgrade` **MUST** produce an upgrade plan before applying changes. The
plan must identify which artifacts need action:

- `/quality` skill;
- `qualitymd` CLI;
- both;
- neither.

`/quality upgrade` **MUST** ask for explicit user confirmation before applying
any upgrade action.

`/quality upgrade` **MUST NOT** edit installed skill files directly. Skill
upgrades must be delegated to the Agent Skills installer or package manager.

`/quality upgrade` **MUST NOT** replace the `qualitymd` binary directly. CLI
upgrades must be delegated to `qualitymd upgrade --apply`, package-manager
commands, or documented install guidance.

When `qualitymd upgrade --check` is available, `/quality upgrade` **SHOULD** use
it to obtain the CLI upgrade recommendation.

When `qualitymd upgrade --apply` is available and the user confirms mutation,
`/quality upgrade` **MAY** invoke it for the CLI upgrade path.

When the skill installer supports checking or upgrading skills, `/quality
upgrade` **SHOULD** use that installer for skill version discovery and mutation.

When the skill installer cannot check or upgrade skills, `/quality upgrade`
**MUST** report that skill upgrade automation is unavailable and provide the
documented manual install or reinstall command.

After applying a CLI upgrade, `/quality upgrade` **MUST** verify the visible
`qualitymd` version and compatibility against the target skill's required CLI
range.

After applying a skill upgrade, `/quality upgrade` **MUST** tell the user that
the current agent session may still be running the previously loaded skill and
that a restart, reload, or new session may be required.

`/quality upgrade` **MUST** stop before running quality setup, evaluation, or
improvement workflows. It is maintenance orchestration, not a quality evaluation
mode.

`wizard` **MAY** recommend `/quality upgrade` when it detects a missing, stale,
or incompatible CLI/skill pair.

The root `SKILL.md` **MUST** list `/quality upgrade` in invocation variants and
route it to a mode-specific `modes/upgrade.md` file.

The mode-specific upgrade instructions **MUST** keep deterministic mechanics in
the CLI or installer and keep the skill responsible for diagnosis,
confirmation, sequencing, and user-facing explanation.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/quality-skill.md` - add `upgrade` to invocation,
  mode dispatch, setup/prerequisite behavior, and mode contracts (per the mode
  and orchestration requirements above).

### To delete

None
