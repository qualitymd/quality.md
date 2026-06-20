---
type: Change Case
title: Update command and improvements
description: Rename upgrade to update with apply-by-default and a --check advisory, rename the /quality upgrade skill mode to update, self-apply managed standalone installs, gate availability on release readiness, surface release notes, and add an ambient cached update nudge.
status: Design
tags: [cli, update, upgrade, install, versioning, skill]
timestamp: 2026-06-20T00:00:00Z
---

# Update command and improvements

A Change Case that renames `qualitymd`'s self-update surface from `upgrade` to
`update` and folds in a set of improvements: a `qualitymd update` command that
applies by default, a `--check` advisory mode, and an ambient update notice
driven by a local cache. The rename also reaches the paired skill — the
`/quality upgrade` mode becomes `/quality update`. It also carries the three
refinements from the original comparison — managed standalone self-apply,
release-readiness gating, and a release-notes reference. Detail lives in the
[functional spec](0041-update-command/spec.md) and
[design doc](0041-update-command/design.md).

> **Design.** Spec and design rewritten for the apply-by-default `update` shape
> and the paired skill-mode rename; ready to advance to **In-Progress**.

## Motivation

[0032](archive/0032-cli-managed-upgrades.md) shipped `qualitymd upgrade` with a
deliberately offline-by-default posture: ordinary commands never touch the
network, the command is advisory unless `--apply` is passed, and only npm and
Homebrew self-update. That kept the CLI deterministic, but it also means a user
only ever learns a new release exists if they think to run a check, and the one
channel the project fully owns — managed standalone — cannot apply.

The decision for this change is to adopt the conventional update-command shape
instead:

- a single `qualitymd update` verb that **applies by default**, with `--check`
  for the advisory path (replacing the `upgrade` + `--apply` split);
- an **ambient update notice**: ordinary commands surface a one-line "update
  available" hint from a locally cached check, so users discover releases without
  running anything.

Because the CLI verb is the user-facing name the `/quality` skill exposes, the
skill's maintenance mode is renamed in lockstep: `/quality upgrade` →
`/quality update`. Keeping the skill mode and the CLI command on the same word
avoids a confusing split where the skill says "upgrade" but drives `qualitymd
update`.

The ambient notice deliberately reverses 0032's "ordinary commands MUST NOT
perform background upgrade checks" rule. What makes that safe — rather than a
regression of the agent-first posture — is that the notice is fenced by strict
rails: it is written to stderr only, never to stdout or `--json` payloads; it is
suppressed on non-interactive terminals, in CI, and behind an opt-out; and it
reads from a cache so no ordinary command ever blocks on the network. Because
qualitymd is agent- and skill-first and agents invoke the CLI non-interactively
(often with `--json`), the notice will rarely fire for agents by design — it is a
courtesy for interactive human users, while the `/quality` skill remains the
primary update-awareness path for agents.

This change also folds in the three earlier refinements so they land on the new
command shape: managed standalone self-apply (now the default `update` path),
release-readiness gating so the tool never advises or applies a not-yet-published
release, and a release-notes reference in advisory and applied output.

## Scope

Covered: rename `upgrade` → `update` with a deprecated `upgrade` alias for one
release cycle; flagless apply-by-default; `--check` advisory and `--json`;
managed standalone self-apply; release-readiness gating on availability and
apply, including resolving Homebrew latest/readiness from the published tap cask
rather than the GitHub release tag; a release-notes reference; an ambient cached
update notice on ordinary
commands with a local cache, bounded background refresh, and a documented
opt-out; suppression of the notice under `--json`, non-interactive/CI, and
development builds; and renaming the `/quality upgrade` skill mode to
`/quality update`, including the skill spec, runtime mode file, routing, and
references.

Deferred / non-goals: no persisted dismissal flow ("skip this version" /
"skip until next version" — that was the full-parity option, not chosen); no
interactive modal or prompt (qualitymd is not a persistent TUI — the notice is a
single stderr line); no rollback command; no pre-release/beta channel selection;
no change to `QUALITY.md` format or evaluation semantics.

## Affected artifacts

### Code

- [ ] [`internal/cli/upgrade.go`](../internal/cli/upgrade.go) — rename to the
      `update` command (file → `internal/cli/update.go`): apply by default,
      `--check` advisory, drop `--apply`; widen the latest-version provider to
      carry readiness and a release-notes reference; gate availability and apply
      on readiness; resolve Homebrew latest/readiness from the published tap cask
      instead of the GitHub release tag (fixing the 0032 bug where brew compared
      against GitHub), with a post-apply `--version` verify on every apply path;
      add managed standalone to the apply path; add a deprecated `upgrade` alias.
- [ ] `internal/cli` root command wiring — a persistent post-run hook that emits
      the ambient notice from cache under the gating rules, and a bounded,
      non-blocking cache refresh.
- [ ] `internal/cli` update-cache helper (new) — read/write the cached
      latest-release record under `$QUALITYMD_HOME`, with a freshness timestamp.
- [ ] [`internal/cli/version.go`](../internal/cli/version.go) — reuse version
      metadata for the notice and dev-build suppression.
- [ ] [`internal/cli/version_upgrade_test.go`](../internal/cli/version_upgrade_test.go)
      → `internal/cli/version_update_test.go` (renamed alongside `update.go`) —
      cover apply-by-default, `--check`, managed standalone apply, readiness
      gating, release-notes output, the alias, post-apply `--version` verify, and
      notice gating (TTY/CI/`--json`/opt-out/dev-build/cache states).

### Durable specs

See the functional spec's
[Durable spec changes](0041-update-command/spec.md#durable-spec-changes)
for the per-requirement breakdown.

- [ ] `specs/cli/update.md` (new, rename of `specs/cli/upgrade.md`) — the
      `update` command contract, readiness, release notes, and the ambient notice.
- [ ] [`specs/cli/upgrade.md`](../specs/cli/upgrade.md) — removed (content carried
      forward into `update.md`).
- [ ] [`specs/cli.md`](../specs/cli.md) — reverse the "unrelated commands MUST NOT
      contact the network for upgrade checks" line; register `update`, the cached
      ambient notice, and the opt-out.
- [ ] [`specs/cli/index.md`](../specs/cli/index.md) — rename the command sub-spec
      entry.
- [ ] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      — rename the `upgrade` skill mode to `update` throughout the durable skill
      contract (mode list, routing, the `Upgrade` section, examples, prerequisite
      commands), and point its CLI references at `qualitymd update`.

No `SPECIFICATION.md` change: this is CLI distribution and update mechanics, not
the document format or evaluation semantics.

### Durable docs

- [ ] [`docs/reference/versioning.md`](../docs/reference/versioning.md) — replace
      the "ordinary commands do not contact the network" guidance with the cached
      ambient-notice model and opt-out; document `update` and managed standalone
      self-apply.
- [ ] [`skills/quality/modes/upgrade.md`](../skills/quality/modes/upgrade.md) —
      rename the runtime mode file to `skills/quality/modes/update.md`; call
      `qualitymd update --check` / `qualitymd update`; reflect that managed
      standalone now self-applies.
- [ ] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) — rename the
      `upgrade` mode to `update` in the router, mode list, routing table, and
      examples; update the command references and the `requires-qualitymd-cli`
      range for the renamed command.
- [ ] [`skills/quality/modes/wizard.md`](../skills/quality/modes/wizard.md) —
      recommend `/quality update` (was `/quality upgrade`) and reference
      `qualitymd update --check`.
- [ ] [`skills/quality/resources/cli-quick-reference.md`](../skills/quality/resources/cli-quick-reference.md)
      — `update` / `update --check`, dropped `--apply`, the alias, and the notice.
- [ ] [`skills/quality/guides/top-10-quality-md-checks.md`](../skills/quality/guides/top-10-quality-md-checks.md)
      and [`specs/skills/quality-skill/guides/top-10-quality-md-checks.md`](../specs/skills/quality-skill/guides/top-10-quality-md-checks.md)
      — rename the `upgrade` route token to `update` in the route list.
- [ ] [`install.md`](../install.md) — `qualitymd update --check`, self-applying
      channels, and how to opt out of the ambient notice.
- [ ] [`CHANGELOG.md`](../CHANGELOG.md) — record the rename (with alias), the
      `/quality update` mode rename, and the ambient notice as notable changes.

### Install/scaffold

- [ ] [`install/install.sh`](../install/install.sh),
      [`install/install.ps1`](../install/install.ps1) — confirm the
      non-interactive contract `update` relies on to drive a managed standalone
      apply without prompts.

## Children

- [Functional spec](0041-update-command/spec.md) — what `update`, the renamed
  `/quality update` skill mode, the readiness gate, the release-notes reference,
  and the ambient notice must provide.
- [Design doc](0041-update-command/design.md) — how the command, cache, and
  notice are built inside the existing `internal/cli` seams.

## Status

`Design`. Spec and design rewritten for the apply-by-default `update` shape and
the paired skill-mode rename; ready to advance to **In-Progress** for
implementation.
