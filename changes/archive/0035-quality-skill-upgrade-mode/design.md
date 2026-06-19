---
type: Design Doc
title: /quality upgrade mode design
description: How the /quality skill plans paired skill and CLI upgrades while delegating mutation to the installer and CLI.
tags: [skill, upgrade, versioning, design]
timestamp: 2026-06-19T00:00:00Z
---

# /quality upgrade mode design

Design behind the
[/quality upgrade mode](../0035-quality-skill-upgrade-mode.md) change and its
[functional spec](spec.md).

## Context

`/quality upgrade` is maintenance orchestration for two separately distributed
artifacts: the `/quality` skill and the `qualitymd` CLI. The skill has release
metadata that names its required CLI range. The CLI has explicit upgrade checks
and guarded apply behavior. The missing piece is a skill workflow that can
inspect both sides, explain the plan, and delegate mutation without taking over
either installer.

The currently running skill also has a hard boundary: changing the installed
skill files does not guarantee the active agent session reloads those
instructions. The workflow must treat skill upgrade as an install-time change
that may require restart or a new session.

## Approach

Implement the mode as a new `skills/quality/modes/upgrade.md` procedure routed
from the root `SKILL.md`.

### Version and compatibility snapshot

The mode starts by collecting a snapshot:

- current skill metadata from the loaded `SKILL.md` frontmatter;
- installed CLI facts from `qualitymd version --json` when present, falling
  back to `qualitymd --version`;
- CLI upgrade recommendation from `qualitymd upgrade --check` when present;
- skill upgrade availability from the Agent Skills installer when it exposes a
  check/list/update surface.

If installer support cannot report the latest skill version, the plan should say
skill automation is unavailable and show the reinstall command as the manual
fallback.

### Plan before mutation

The skill renders a concise plan with four possible outcomes:

- no action needed;
- CLI action only;
- skill action only;
- both skill and CLI actions.

The plan should state why each action is needed: missing CLI, CLI outside the
required range, newer CLI available, skill update available, or unavailable
automation. This keeps the mode useful even when it cannot apply every action.

### Delegated apply

After explicit confirmation, the mode delegates:

- CLI mutation to `qualitymd upgrade --apply` when available and applicable;
- skill mutation to the Agent Skills installer or package manager when available;
- otherwise, documented manual commands.

The mode does not edit skill files and does not replace binaries. It runs only
the owner commands or reports the command the user should run.

### Verification and reload guidance

After CLI mutation, the mode re-runs the CLI version check and verifies the
visible CLI satisfies the target skill's required range.

After skill mutation, the mode reports that the active session may still be using
previously loaded instructions. The final message should include the concrete
restart/reload/new-session action appropriate to the agent surface when known,
or generic restart guidance otherwise.

## Alternatives

**Only document manual upgrade steps.** Rejected. Manual steps remain as
fallbacks, but users need one skill-facing workflow that reasons about the
skill/CLI pair.

**Make `qualitymd upgrade` update the skill too.** Rejected. The CLI should not
own Agent Skills installation or mutate skill packages; that belongs to the
skill installer.

**Make `/quality upgrade` directly edit installed skill files.** Rejected.
Direct file edits bypass installer ownership, provenance, and any future package
metadata or registry behavior.

**Run upgrade checks automatically in every mode.** Rejected. Normal quality
workflows should verify prerequisites, but network/latest-version checks stay
explicit so setup, wizard, evaluate, and improve remain predictable.

**Assume the upgraded skill is active immediately.** Rejected. Agents may load
skill instructions once per session. The workflow must tell the user when a
restart or new session is needed.

## Trade-offs and risks

The mode depends on installer capabilities that may not exist everywhere. The
design handles this by producing an advisory plan and manual fallback instead of
making upgrade mode all-or-nothing.

The workflow can become confusing if it mixes "compatible" with "latest."
Compatibility should be the hard gate; latest-version updates are advisory
unless the user asked to upgrade.

Applying both skill and CLI updates can change the required range during the
workflow. The mode should verify against the post-upgrade skill metadata when it
can observe it, and otherwise verify against the current skill metadata while
telling the user to restart for the installed skill update to take effect.

## Open questions

- Which Agent Skills installer command should be treated as the canonical
  non-interactive skill upgrade command once the installer exposes one?
- Should `/quality upgrade` support a check-only variant in user language, or is
  "show the plan and ask before applying" enough?
