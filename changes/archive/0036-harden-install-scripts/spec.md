---
type: Functional Specification
title: Harden install scripts and upgrade idiomatics — functional spec
description: Portability and convention requirements for the shell/PowerShell installers and the upgrade version check.
tags: [cli, install, upgrade, homebrew]
timestamp: 2026-06-19T00:00:00Z
---

# Harden install scripts and upgrade idiomatics — functional spec

Companion to the
[Harden install scripts and upgrade idiomatics](../0036-harden-install-scripts.md)
change case. This spec states *what* the installers and the upgrade check must
do; the [design doc](design.md) covers *how*. It defers the format and the broad
install/upgrade contract to [`specs/cli/upgrade.md`](../../../specs/cli/upgrade.md)
and the [versioning reference](../../../docs/reference/versioning.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The install/upgrade surfaces follow the right idioms (`curl | sh`, `irm | iex`,
per-platform npm packages, a Homebrew cask, a delegating `qualitymd upgrade`),
but a convention review found portability bugs in the details: checksum
verification that silently does nothing on most Linux, a Windows installer that
can fail on a stock TLS configuration, missing PATH integration, a string-compare
masquerading as a version check, and a non-interactive flag with no effect. Each
weakens the guarantee the surrounding idiom is supposed to provide. The fixes are
narrow and platform-idiomatic; they do not expand the install surface.

## Scope

Covered: shell-installer checksum robustness, Windows TLS, PATH guidance on both
installers, SemVer-correct update detection, honest non-interactive semantics,
and a durable record that the Homebrew cask is the intended path.

Deferred: signing/notarizing the macOS binary (would let the cask quarantine hook
be dropped). Non-goals: a `homebrew-core` source formula; new channels (winget,
Scoop, Linux repos); automatic Unix dotfile editing.

## Requirements

### Shell installer (`install/install.sh`)

- Checksum verification **MUST** compute the archive's SHA-256 using whichever
  tool is available, trying at least `sha256sum`, then `shasum -a 256`, then
  `openssl dgst -sha256`. When `checksums.txt` was downloaded but no digest tool
  is available, or the archive has no entry in it, the installer **MUST** report
  that verification was skipped on standard error rather than passing silently.

  > Rationale: gating on `shasum` alone verified nothing on stock Linux (which
  > ships `sha256sum`), and the skip was invisible — a silent no-op reads as
  > "verified." — 0036

- A verified checksum mismatch **MUST** abort the install with a non-zero exit.

- On success the installer **MUST** emit the exact line needed to put the binary
  on `PATH` (an `export PATH=...` for the resolved `bin` directory) when that
  directory is not already on `PATH`. It **MUST NOT** edit shell profiles or other
  dotfiles.

  > Rationale: piped `curl | sh` has no reliable TTY for consent and dotfile
  > editing is easy to get wrong; printing the line is the safe, idiomatic floor.
  > — 0036

- The installer **MUST** accept `--non-interactive` / `--yes` and the
  `QUALITYMD_NO_INPUT` environment variable, and they **MUST** have an observable
  effect: in that mode the installer suppresses the human-oriented PATH guidance
  and prints only a concise completion line. The flags **MUST** remain accepted
  so documented CI and smoke invocations keep working.

  > Rationale: the flag previously only swapped near-identical text, implying an
  > interactive mode that does not exist. — 0036

### PowerShell installer (`install/install.ps1`)

- Before any network call, the installer **MUST** ensure TLS 1.2 (or stronger) is
  enabled for the session so it works under Windows PowerShell 5.1, whose default
  protocols GitHub rejects.

- After installing, the installer **MUST** ensure the `bin` directory is on the
  user's persistent `PATH` (the per-user environment, not machine-wide) when it is
  absent, update the current session's `PATH`, and tell the user a new shell is
  needed for other sessions to see it.

  > Rationale: Windows exposes a clean per-user PATH API, so the idiomatic
  > installer (Scoop, rustup-init) updates it rather than only advising; the prior
  > script did neither. — 0036

- `-NonInteractive` **MUST** have an observable effect consistent with the shell
  installer: it suppresses human-oriented guidance while still performing the
  install and PATH update.

### Upgrade check (`qualitymd upgrade`)

- The command **MUST** report that an update is available only when the latest
  version is strictly newer than the current version by SemVer precedence. When
  either version is not valid SemVer, it **MAY** fall back to reporting a
  difference, and it **MUST NOT** report an update for a development build.

  > Rationale: a plain string inequality flagged downgrades and older prereleases
  > as upgrades. — 0036

### Homebrew distribution

- The project's durable docs **MUST** record that the Homebrew **cask** is the
  intended, idiomatic distribution path for the pre-built binary — GoReleaser
  deprecated the formula path and Homebrew steers binaries to casks — and that the
  quarantine post-install hook is the documented pattern for an unsigned binary,
  removable once the binary is signed and notarized. This change makes no
  functional change to the cask.

## Durable spec changes

### To add

None.

### To modify

- [`specs/cli/upgrade.md`](../../../specs/cli/upgrade.md) — state that "update
  available" means the latest version is strictly newer by SemVer precedence (per
  the upgrade-check requirement above), and add PATH setup to the managed-installer
  SHOULD list (per the installer PATH requirements above).

### To delete

None.
