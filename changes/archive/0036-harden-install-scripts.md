---
type: Change Case
title: Harden install scripts and upgrade idiomatics
description: Fix portability and convention gaps in the shell/PowerShell installers and the upgrade version check, and record why the Homebrew cask is the idiomatic distribution path.
status: Done
tags: [cli, install, upgrade, homebrew]
timestamp: 2026-06-19T00:00:00Z
---

# Harden install scripts and upgrade idiomatics

A **Change Case** capturing the _why_ and _status_ for a focused hardening pass
on the `qualitymd` install and upgrade surfaces. The detail lives in its
[functional spec](0036-harden-install-scripts/spec.md) and
[design doc](0036-harden-install-scripts/design.md).

## Motivation

A review of the install/upgrade surfaces against current platform and CLI
conventions found the architecture sound — `curl | sh`, `irm | iex`, per-platform
npm optional dependencies, a Homebrew cask, `go install`, and a delegating
`qualitymd upgrade` — but surfaced concrete portability bugs and convention gaps
in the details that imitate those idioms:

1. **Checksum verification is dead on most Linux.** `install.sh` gates SHA-256
   verification on `command -v shasum`; stock Linux ships `sha256sum`, not
   `shasum`, so the downloaded `checksums.txt` is verified against _nothing_ on a
   typical Linux host — and the skip is silent.
2. **`install.ps1` can fail on stock Windows.** Windows PowerShell 5.1 defaults
   to TLS 1.0/1.1, which GitHub's endpoints reject, so `Invoke-WebRequest` can
   fail out of the box with no TLS 1.2 pin.
3. **PATH integration is missing/asymmetric.** `install.sh` prints terse PATH
   advice; `install.ps1` neither updates PATH nor advises, leaving a fresh
   install invisible with no hint.
4. **Upgrade "update available" is a string compare, not SemVer.** Any difference
   reads as an update, including a downgrade or an older prerelease.
5. **`--yes` / `--non-interactive` is effectively a no-op** in `install.sh`,
   implying an interactive mode that does not exist.

A separate review item — "convert the Homebrew cask to a formula" — was
investigated and **reversed**: GoReleaser deprecated the formula (`brews`) path
in v2.10 (removal in v3), Homebrew steers pre-built binary distributions to
casks, and the cask's quarantine post-install hook is the documented pattern for
unsigned binaries. The existing `homebrew_casks` setup is the idiomatic
self-published path; the durable record should say so, so the question is not
re-litigated.

## Scope

Covered: the five portability/convention fixes above across `install/install.sh`,
`install/install.ps1`, and `internal/cli/upgrade.go`, plus a durable note that
the Homebrew cask is intentional and idiomatic.

Deferred / non-goals: no signing/notarization of the macOS binary (the real fix
that would let the quarantine hook be dropped — recorded as deferred, not done
here); no pursuit of a `homebrew-core` source formula (notability-gated,
externally maintained); no new install channels (winget, Scoop, Linux package
repos); no automatic profile/dotfile editing on Unix; no background network
checks in ordinary commands.

## Affected specs & docs

- [x] [`specs/cli/upgrade.md`](../../specs/cli/upgrade.md) — states that "update
      available" means the latest version is strictly newer by SemVer precedence,
      and adds PATH setup to the managed-installer SHOULD list (per the SemVer and
      PATH requirements).
- [x] [`install.md`](../../install.md) — documents the PATH guidance the installers
      now emit, the honest meaning of `--non-interactive`, and the piped-argument
      form; keeps the existing `brew install` cask command.
- [x] [`CONTRIBUTING.md`](../../CONTRIBUTING.md) — notes why the Homebrew **cask**
      (not a formula) is the distribution choice.
- [x] [`docs/guides/cut-a-release.md`](../../docs/guides/cut-a-release.md) — frames the
      cask as the GoReleaser-recommended path and the quarantine strip as the
      documented unsigned-binary pattern.
- [x] `install/install.sh`, `install/install.ps1` — the installer fixes
      (checksum fallback, TLS 1.2, PATH, non-interactive semantics).
- [x] `internal/cli/upgrade.go` and `internal/cli/version_upgrade_test.go` —
      SemVer-aware update detection and its regression coverage.
- [x] `.goreleaser.yaml` — a clarifying comment that the cask is intentional and
      the formula path is deprecated.

No `SPECIFICATION.md` update is expected: this changes CLI distribution and
upgrade mechanics, not the `QUALITY.md` document format or evaluation semantics.

## Children

- [Functional spec](0036-harden-install-scripts/spec.md) — what the installers
  and upgrade check must do.
- [Design doc](0036-harden-install-scripts/design.md) — how the fixes are built,
  and why the cask stays.

## Status

`Done`. Implementation, durable artifact synchronization, review, and archival
are complete. See the [status lifecycle](../index.md#status-lifecycle).
