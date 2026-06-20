---
type: Functional Specification
title: Update command and improvements - functional spec
description: Rename upgrade to update with apply-by-default and a --check advisory, rename the /quality upgrade skill mode to update, self-apply managed standalone, gate availability on release readiness, surface release notes, and add an ambient cached update notice.
tags: [cli, update, upgrade, install, versioning, skill]
timestamp: 2026-06-20T00:00:00Z
---

# Update command and improvements - functional spec

Companion to [Update command and improvements](../0041-update-command.md).
This spec restates the self-update contract from
[`specs/cli/upgrade.md`](../../../specs/cli/upgrade.md) under the apply-by-default
`update` shape, and adds the ambient-notice, readiness, and release-notes
behavior, plus the paired rename of the `/quality upgrade` skill mode to
`/quality update`. When this change lands, the durable spec is renamed to
`specs/cli/update.md`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", and "MAY" are to be
interpreted as described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

[0032](../archive/0032-cli-managed-upgrades.md) made the self-update surface
advisory-by-default and forbade any network access from ordinary commands, so a
user learns about a release only by explicitly running a check. This change
adopts the conventional update-command shape instead: one `update` verb that
applies by default, and an ambient notice that surfaces an available release from
a local cache during ordinary commands. The CLI verb is also the name the
`/quality` skill exposes, so the skill's `upgrade` maintenance mode is renamed to
`update` in lockstep — the skill word and the CLI command stay the same.

The ambient notice reverses 0032's prohibition on background checks. It is made
safe — not a regression of the agent-first posture — by strict fencing: stderr
only, never in stdout or `--json`; suppressed off a terminal, in CI, and behind an
opt-out; and served from a cache so no ordinary command blocks on the network.
Because qualitymd is agent- and skill-first and agents drive the CLI
non-interactively, the notice rarely fires for agents by design; it is a courtesy
for interactive users while the `/quality` skill stays the primary update path for
agents.

## Scope

Covered: the `update` command (apply-by-default, `--check`, `--json`); a
deprecated `upgrade` alias; managed standalone self-apply; readiness gating on
availability and apply; a release-notes reference; an ambient cached update
notice with bounded refresh and an opt-out; and renaming the `/quality upgrade`
skill mode to `/quality update`.

Deferred / non-goals: no persisted dismissal ("skip this version"); no
interactive prompt or modal; no rollback; no pre-release channel selection; no
format or evaluation-semantics change.

## Requirements

### Command shape

`qualitymd update` invoked without a mode flag **MUST** attempt to apply an
update through the install channel that owns the binary, subject to apply support
and the readiness gate below. When the detected install method has no safe apply
action (unknown, archive, or Go/source), `update` **MUST** refuse to mutate the
install and print manual guidance instead.

> Rationale: applying is the common intent, so it is the default; refusing to
> mutate an unowned binary is still mandatory, because a package manager, admin,
> or pinned toolchain may own it. — 0041, carrying 0032

`qualitymd update --check` **MUST** report, without mutating the install, the
current version, latest known version when available, detected install method,
whether an update is available, whether apply is supported for this install, the
recommended action, and the release-notes reference when known.
`qualitymd update --json` **MUST** emit those facts with stable field names. An
update is available only when the latest version is strictly newer than the
current version by SemVer precedence and the readiness gate is satisfied; a
development build **MUST NOT** report an available update.

The CLI **SHOULD** retain a deprecated `upgrade` alias that behaves as `update`
for one release cycle and prints a deprecation notice to stderr.

> Rationale: the `/quality` skill and existing scripts call `qualitymd upgrade`;
> an alias lets the rename land without breaking the paired skill-and-CLI flow
> mid-transition. — 0041

### Managed standalone self-apply

`qualitymd update` **MUST** support a detected managed standalone install by
invoking that channel's managed installer non-interactively, rather than refusing.
On success it **MUST** verify the running version through the post-install
`qualitymd --version` check the managed standalone update path already requires.
The installer invocation **MUST** run without interactive prompts so `update`
works for the `/quality` skill and CI.

> Rationale: 0032 refused to apply for every channel except npm and Homebrew to
> avoid mutating unowned binaries. That reasoning does not apply to the managed
> standalone channel, whose layout qualitymd owns and whose installer is already
> scripted, checksum-verified, and non-interactive — so it was the one channel
> the project controls yet could not update. — 0041

npm installs **MUST** apply through `npm install -g quality.md@latest` and
Homebrew installs through the documented Homebrew command, unchanged from 0032.

### Release readiness

For the channel that owns the installation, an update **MUST NOT** be reported
available, and `update` **MUST NOT** apply one, unless the target release is
actually retrievable for that channel — for the GitHub-backed channels, that the
release publishes the platform archive and its checksums; for a package-manager
channel, that the registry reports the version as the installable latest. When
readiness cannot be confirmed, `update --check` **MUST** report no available
update and **SHOULD** make the unconfirmed state legible rather than conflating it
with "up to date"; `update` **MUST** fail before mutating, with a diagnostic that
the target release is not yet available.

> Rationale: a release tag can appear before its archive and checksum assets
> finish uploading. Advising — or, now that apply is the default, running — an
> update to a not-yet-downloadable release fails confusingly mid-update. — 0041

Readiness checking **MUST** stay within the explicit `update` invocation and the
bounded cache refresh below; it **MUST NOT** add a synchronous network fetch to
any other command.

### Release-notes reference

`update --check`, an applied `update`, and the ambient notice **SHOULD** include a
reference to the target release's notes when one is known for the owning channel.
`update --json` **MUST** carry that reference under a stable field, omitted when
not known. The reference is advisory: its presence or absence **MUST NOT** change
update availability or apply behavior.

> Rationale: versions alone don't tell a user whether to update; a link to what
> changed keeps that decision in the tool. — 0041

### Ambient update notice

Ordinary `qualitymd` commands **MAY** surface a one-line update-available notice
when a local cache indicates a newer, ready release exists. This is the deliberate
exception to the otherwise-offline default; every other ordinary command behavior
**MUST** remain deterministic and offline.

The notice **MUST** be written to stderr and **MUST NOT** appear in stdout or in
any `--json` / machine-readable output. It **MUST** be suppressed when the output
is not an interactive terminal, when a CI environment is detected, when the
documented opt-out is set, and for development builds. The notice **MUST** name the
current and latest versions and the exact command to run, and **SHOULD** include
the release-notes reference when known. It **MUST NOT** mutate anything.

> Rationale: stdout and `--json` are consumed by agents and pipelines; a notice
> there would corrupt parsed output. Terminal/CI/opt-out gating keeps the notice
> to interactive humans. — 0041

The notice **MUST** be served from a local cache; an ordinary command **MUST NOT**
block on or perform a synchronous network fetch to produce it, and its exit code
and primary output **MUST NOT** be affected by the notice's presence, absence, or
any failure to produce it. The cache **MAY** be refreshed by a bounded,
best-effort check no more frequently than a documented interval; a failed refresh
**MUST** be silent. A refreshed value **MAY** surface on a later invocation rather
than the one that triggered the refresh.

> Rationale: a freshly fetched value need not appear on the triggering command —
> using the previous cache and refreshing for next time is what keeps ordinary
> commands non-blocking. — 0041

### Configuration and opt-out

The CLI **MUST** provide a documented way to disable update checks and the ambient
notice entirely, for centrally managed, air-gapped, or noise-sensitive installs.
Checks and the notice are **enabled by default**.

> Rationale: an always-on network check and stderr notice must be escapable for
> managed fleets and offline environments without disabling the rest of the CLI.
> — 0041

### Paired `/quality` skill mode rename

The `/quality` skill's maintenance mode currently named `upgrade` **MUST** be
renamed to `update`, so the skill invocation is `/quality update`. The rename
**MUST** cover the durable skill contract and the runtime skill files: the mode
list and router, the routing table, the mode's own procedure file, the wizard's
recommendation of it, the CLI quick reference, and the route token in the
top-10-checks guide. Every CLI reference inside the renamed mode **MUST** target
`qualitymd update` (and `qualitymd update --check`) rather than `qualitymd
upgrade`. The mode's behavior — diagnosing the installed skill/CLI pair, planning
before mutation, asking before applying, and delegating CLI mechanics to the
binary — is unchanged; only the name and the CLI commands it drives change.

> Rationale: the skill verb is the name users type, and it should match the CLI
> command it drives. Leaving the skill on `upgrade` while the CLI moves to
> `update` splits the vocabulary the project asks users and agents to learn. — 0041

## Example

```text
$ qualitymd lint QUALITY.md          # stdout: normal lint output
update available: 0.4.1 -> 0.5.0 (run `qualitymd update`)   # stderr, interactive only
                                                            #   https://github.com/qualitymd/quality.md/releases/tag/v0.5.0

$ qualitymd update --check --json
{
  "currentVersion": "0.4.1",
  "installMethod": "managed-standalone",
  "latestVersion": "0.5.0",
  "updateAvailable": true,
  "applySupported": true,
  "recommendedAction": "qualitymd update",
  "releaseNotesURL": "https://github.com/qualitymd/quality.md/releases/tag/v0.5.0"
}
```

Field names are illustrative; the implementation must carry equivalent facts with
stable machine-readable names.

## Durable spec changes

### To add

- `specs/cli/update.md` — the `update` command contract: apply-by-default and
  `--check`/`--json`, the deprecated `upgrade` alias, managed standalone
  self-apply, readiness gating, the release-notes reference, the ambient cached
  notice with its stderr/TTY/CI/`--json`/dev-build gating and non-blocking
  refresh, and the opt-out (per the requirements above). Carries forward the
  surviving content of `specs/cli/upgrade.md`.

### To modify

- `specs/cli.md` — reverse the rule that unrelated commands must not contact the
  network for update checks, and register `update`, the cached ambient notice,
  and the opt-out (per the command-shape and ambient-notice requirements above).
- `specs/cli/index.md` — rename the command sub-spec entry from `upgrade` to
  `update` (per the added durable spec above).
- `specs/skills/quality-skill/quality-skill.md` — rename the `upgrade` skill
  mode to `update` throughout the durable skill contract (mode list, router,
  routing table, the `Upgrade` section, examples, and prerequisite commands), and
  retarget its CLI references to `qualitymd update` (per the paired skill-mode
  rename requirement above).

### To delete

- `specs/cli/upgrade.md` — renamed to `specs/cli/update.md`; no behavior is
  dropped silently, only moved and extended.
