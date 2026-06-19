---
type: Design Doc
title: CLI managed upgrades design
description: How qualitymd adds version metadata, upgrade checks, install-method detection, and managed installer scripts without ambient update behavior.
tags: [cli, install, upgrade, versioning, design]
timestamp: 2026-06-19T00:00:00Z
---

# CLI managed upgrades design

Design behind the [CLI managed upgrades](../0032-cli-managed-upgrades.md)
change and its [functional spec](spec.md).

## Context

`qualitymd` already reports a useful human `--version`, including release stamps
and development fallback metadata. The missing piece is an explicit, structured
way for humans, agents, and the `/quality` skill to answer:

- what binary is running;
- what format specification it embeds;
- whether a newer release exists; and
- which install channel should own the upgrade.

The design must preserve the project's agent-facing CLI discipline: ordinary
commands stay deterministic and offline, machine-readable facts are available
under `--json`, and self-modification is allowed only where ownership is clear.

## Approach

Implement the change as four separable layers, so each can land or be tested
without depending on the full managed installer.

### Version metadata

Add a `version` subcommand next to the existing harness-level `--version`.
`qualitymd --version` remains the compatibility-preserving human shortcut.
`qualitymd version` renders richer human metadata, and
`qualitymd version --json` emits the agent-facing receipt.

Use the existing `buildInfo()` path as the source for CLI version and commit.
Add a small internal version-info package only if more than one command needs
the same struct; otherwise keep the command-local implementation small until the
second consumer appears. The bundled specification version should be read from
the same embedded spec source used by `qualitymd spec`, not manually duplicated.

### Install context

Add install-method detection as a narrow internal helper with an enum-like
result:

- `managed-standalone`
- `npm`
- `homebrew`
- `go-source`
- `archive`
- `unknown`

The first implementation should favor positive signals over clever path
guessing:

- npm: the launcher sets an explicit environment marker before execing the
  native binary;
- Homebrew: macOS executable paths under Homebrew prefixes, and later any cask
  metadata signal if needed;
- managed standalone: executable path under the managed install root and a
  marker file in that root;
- Go/source/archive/custom: detected only where evidence is cheap and reliable;
  otherwise fall through to `unknown`.

The helper should be side-effect free. It classifies the current executable; it
does not fetch latest versions or run update commands.

### Upgrade check and action selection

Add `qualitymd upgrade --check` as the explicit network boundary. It gathers:

- current version metadata;
- install context;
- latest version for the relevant channel; and
- the recommended action.

The latest-version provider should be injectable so tests can run without
network access. At runtime:

- managed standalone and generic releases can check GitHub Release metadata;
- npm can check the npm package version or verify package readiness when needed;
- Homebrew should prefer the tap/cask-published version because it can lag a
  GitHub Release.

`qualitymd upgrade` without `--apply` should be advisory: print the exact
recommended command and emit the same result shape under `--json`. With
`--apply`, the command runs only actions selected from a fixed table. Unknown,
archive, and source/custom installs refuse with manual guidance.

For package-manager installs, apply mode delegates to the owner:

- npm: `npm install -g quality.md@latest`
- Homebrew: `brew upgrade qualitymd/tap/qualitymd`

Managed standalone can run the GitHub-hosted installer in non-interactive latest
mode once that installer exists. Until then, it should report apply unsupported
instead of synthesizing a replacement.

### Managed installer assets

Place public installer entrypoints under top-level `install/`, not `scripts/`,
because this repo uses `scripts/` for internal development and release helpers.

Add:

- `install/install.sh` for macOS/Linux;
- `install/install.ps1` for Windows; and
- optional `install/install.cmd` as a thin `cmd.exe` launcher that delegates to
  `install.ps1`.

The install scripts should be idempotent and serve both install and update. They
resolve a release version, download the platform archive and checksums from
GitHub, verify the archive, stage the new release under a managed root, switch a
stable command path atomically where the platform allows it, and verify the
visible `qualitymd --version` before reporting success.

Use a marker file in the managed root so `qualitymd` can detect standalone
ownership without brittle path inference alone. The marker should include a
layout version so future installers can migrate deliberately.

## Phasing

1. Add `qualitymd version` / `--json` and the shared version metadata shape.
2. Add install-context detection and the npm launcher marker.
3. Add `qualitymd upgrade --check` with injected latest-version providers and
   JSON output.
4. Add advisory `qualitymd upgrade` and guarded `--apply` for npm/Homebrew.
5. Add GitHub-hosted `install/` entrypoints and managed-standalone detection.
6. Enable `qualitymd upgrade --apply` for managed standalone installs.

This order gives the `/quality` skill useful structured facts early, before the
full installer exists.

## Alternatives

**Use `getquality.md` as the copyable installer URL.** Rejected for now. The
domain currently redirects to GitHub rather than hosting installer assets.
Copyable install commands should name the real asset host so agents, CI, and
security review see where code is fetched from.

**Put public installer scripts under `scripts/install/`.** Rejected. The current
`scripts/` directory is internal tooling (`build-npm`, release checks, README
sync). A top-level `install/` directory makes public entrypoints visible and
keeps internal helpers visually separate.

**Create separate update scripts.** Rejected for the first version. The install
scripts should be idempotent and support latest/non-interactive execution, so
`qualitymd upgrade --apply` can rerun the installer rather than maintaining a
parallel update path.

**Maintain separate `.bat` installer logic.** Rejected. PowerShell is the
Windows implementation surface; `.cmd` is enough as a compatibility launcher and
must delegate to `install.ps1`.

**Run background update checks from ordinary commands.** Rejected.
`qualitymd` is used by agents and CI, so network access and output changes need
to be explicit.

**Replace arbitrary current executables in place.** Rejected. The binary may be
managed by a package manager, pinned by a toolchain, or staged in a test
environment. Only the owner channel should mutate it.

## Trade-offs and risks

The staged approach means `qualitymd upgrade --check` may exist before direct
managed-standalone updates. That is acceptable because it still improves the
skill and user workflow with structured facts and concrete guidance.

Install-method detection can be wrong if it relies on paths alone. The design
mitigates this by adding explicit markers where we control the launcher or
installer, and by refusing mutation when ownership is uncertain.

Installer scripts add release-surface complexity: checksums, locks, PATH edits,
Windows shell behavior, and idempotency all need tests or release verification.
Keeping install and update in one script per platform limits the number of
surfaces that can drift.

Homebrew and npm can lag GitHub releases. Latest-version checks need
channel-aware providers so the tool does not tell users a package-manager install
is stale before that channel can actually provide the new version.

## Open questions

- Should the managed standalone root live under a project-specific home such as
  `$QUALITYMD_HOME`, under the existing user config convention if one emerges,
  or under a platform-standard data directory?
- Should uninstall scripts be part of this change or a follow-up once the exact
  managed root and PATH-edit behavior are settled?
- Should `qualitymd upgrade --apply` execute package-manager commands directly
  by default, or require an extra confirmation flag for package-manager installs
  outside non-interactive contexts?
