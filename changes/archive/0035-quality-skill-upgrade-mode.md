---
type: Change Case
title: /quality upgrade mode
description: Add a /quality upgrade mode that plans and orchestrates compatible upgrades for both the /quality skill and the qualitymd CLI.
status: Done
tags: [skill, upgrade, versioning]
timestamp: 2026-06-19T00:00:00Z
---

# /quality upgrade mode

A **Change Case** capturing the _why_ and _status_ for adding an upgrade
orchestration mode to the `/quality` skill. The detail lives in its
[functional spec](0035-quality-skill-upgrade-mode/spec.md).

> **Done.** Implementation and durable artifact synchronization are complete;
> the change is archived.

## Motivation

The `/quality` skill and `qualitymd` CLI are distributed separately but must
remain compatible. Skill release metadata records the installed skill version
and required CLI range, while CLI managed-upgrade support gives the CLI an
explicit way to report and apply updates through supported owner channels. Users
still need one skill-facing workflow that can inspect the pair, explain what is
stale or incompatible, and coordinate the right upgrade actions.

`/quality upgrade` should be that orchestration layer. It keeps mechanics in the
owners that already have them: the Agent Skills installer upgrades the skill,
and `qualitymd upgrade` or package-manager guidance upgrades the CLI. The skill
owns diagnosis, sequencing, confirmation, and the session-reload warning that
follows a skill upgrade.

## Scope

Covered: a new `/quality upgrade` mode, installed skill and CLI version checks,
compatibility diagnosis, paired upgrade planning, explicit confirmation before
mutation, delegation to skill-installer and CLI upgrade mechanics, post-upgrade
verification, and restart/reload guidance.

Deferred / non-goals: no hand-editing installed skill files, no custom skill
package manager inside `qualitymd`, no background upgrade checks during normal
quality workflows, no automatic mutation without confirmation, and no claim that
the currently running agent session can hot-reload upgraded skill instructions.

Implementation is complete and archived.

## Affected specs & docs

- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md) - add `upgrade` to invocation, mode dispatch, and the skill's maintenance
      workflow contract.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) and new
      `skills/quality/modes/upgrade.md` - route and specify the runtime upgrade
      procedure.
- [x] [`skills/quality/modes/wizard.md`](../../skills/quality/modes/wizard.md) -
      allow the wizard to recommend `/quality upgrade` for stale or incompatible
      skill/CLI state.
- [x] [`docs/guides/use-quality-skill.md`](../../docs/guides/use-quality-skill.md) - document the upgrade mode as the skill-facing maintenance path.
- [x] [`install.md`](../../install.md) - document `/quality upgrade` for existing
      installs and clarify manual fallbacks.
- [x] [`docs/reference/versioning.md`](../../docs/reference/versioning.md) -
      explain paired skill/CLI upgrade behavior and session reload limits.

No `SPECIFICATION.md` update is expected: this changes skill maintenance
workflow, not the `QUALITY.md` document format or evaluation semantics.

## Children

- [Functional spec](0035-quality-skill-upgrade-mode/spec.md) - what the
  `/quality upgrade` mode must provide.
- [Design doc](0035-quality-skill-upgrade-mode/design.md) - how upgrade mode
  plans, delegates, verifies, and handles session reload limits.

## Status

`Done`. The change is archived.
