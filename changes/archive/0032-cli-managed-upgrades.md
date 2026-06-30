---
type: Change Case
title: CLI managed upgrades
description: Add explicit version metadata and upgrade workflows so qualitymd can check, advise, and eventually apply updates through a managed install path.
status: Done
tags: [cli, install, upgrade, versioning]
timestamp: 2026-06-19T00:00:00Z
---

# CLI managed upgrades

A **Change Case** capturing the _why_ and _status_ for adding a long-term
`qualitymd` install and upgrade contract. The detail lives in its
[functional spec](0032-cli-managed-upgrades/spec.md).

> **Done.** Implementation and durable artifact synchronization are complete; the change is archived.

## Motivation

The `/quality` skill depends on a compatible `qualitymd` CLI, but today the
remediation path for a missing or stale CLI is documentation-driven: check
`qualitymd --version`, then choose npm, Homebrew, or `go install` manually. That
is acceptable for early distribution, but it leaves agents and users without a
single deterministic way to ask "what am I running, is it current, and what is
the right upgrade path for this installation?"

Codex's update flow shows a useful structure: separate version display, latest
version checks, install-method detection, and the action that applies an update.
`qualitymd` should adopt that separation, but with a stricter default posture:
normal commands stay deterministic and offline; update checks happen only when
explicitly requested; and self-modification is allowed only for installs whose
layout the project manages.

The long-term target is a GitHub-hosted managed standalone installer that gives
the project a safe update domain: platform detection, checksums, atomic
replacement, PATH setup, conflict detection, and post-install verification.
Package-manager, source, archive, and custom installs remain supported, but
`qualitymd` should delegate or advise for those channels instead of pretending
every binary can be replaced safely.

## Scope

Covered: structured version metadata, an explicit upgrade check/apply workflow,
install-method detection, npm launcher support for reliable detection,
GitHub-hosted `install.sh` and `install.ps1` scripts, an optional Windows
`.cmd` launcher, and the contract for a future managed standalone
installer/upgrader.

Deferred / non-goals: no background network checks during ordinary `qualitymd`
commands, no automatic prompts, no update action for unknown/custom binaries,
and no change to `QUALITY.md` format semantics. The first implementation may
ship advisory checks before the managed standalone installer exists.

Implementation is complete and archived.

## Affected specs & docs

- [x] [`specs/cli.md`](../../specs/cli.md) - extend the CLI overview and version
      conventions with structured version output and explicit upgrade behavior.
- [x] [`specs/cli/index.md`](../../specs/cli/index.md) and new `specs/cli/version.md`
      / `specs/cli/upgrade.md` sub-specs - specify version metadata and upgrade
      check/apply command behavior.
- [x] [`docs/reference/versioning.md`](../../docs/reference/versioning.md) - explain
      how CLI release compatibility, latest-version checks, and managed install
      channels relate to the skill's supported CLI range.
- [x] [`install.md`](../../install.md) and [`README.md`](../../README.md) - document the
      preferred GitHub-hosted install path, package-manager alternatives,
      `qualitymd upgrade --check`, non-interactive/pinned install forms, and
      manual upgrade fallbacks.
- [x] [`docs/guides/cut-a-release.md`](../../docs/guides/cut-a-release.md) - add
      release verification for version metadata, upgrade checks, and any managed
      installer artifacts.
- [x] `install/install.sh`, `install/install.ps1`, and optional
      `install/install.cmd` (new) - provide the GitHub-hosted Unix/macOS and
      Windows managed standalone installer entrypoints.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) and
      [`skills/quality/resources/cli-quick-reference.md`](../../skills/quality/resources/cli-quick-reference.md) - teach the skill to use structured version/upgrade checks before
      CLI-dependent workflows.
- [x] [`npm/quality.md/bin/qualitymd.js`](../../npm/quality.md/bin/qualitymd.js) -
      mark npm-launched binaries so install-method detection is reliable.

No `SPECIFICATION.md` update is expected: this changes CLI distribution and
upgrade mechanics, not the `QUALITY.md` document format or evaluation semantics.

## Children

- [Functional spec](0032-cli-managed-upgrades/spec.md) - what version metadata,
  upgrade checks, upgrade application, and managed-install behavior must provide.
- [Design doc](0032-cli-managed-upgrades/design.md) - how the version, upgrade,
  install-context, and installer pieces fit together.

## Status

`Done`. The change is archived.
