---
type: Functional Specification
title: CLI managed upgrades - functional spec
description: Define structured version metadata, explicit upgrade checks, safe upgrade application, and the managed standalone install path.
tags: [cli, install, upgrade, versioning]
timestamp: 2026-06-19T00:00:00Z
---

# CLI managed upgrades - functional spec

Companion to [CLI managed upgrades](../0032-cli-managed-upgrades.md). This spec
states the CLI distribution and upgrade-behavior delta for `qualitymd`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The CLI is the runtime compatibility boundary for the `/quality` skill. When it
is missing or stale, the skill needs reliable facts and a concrete remediation
path before it can run setup, lint, or evaluation workflows. Today those facts
are split between human `--version` output and installation docs.

An upgrade surface should keep the mechanical work in the CLI while preserving
`qualitymd`'s deterministic, agent-friendly defaults. Ordinary commands should
not perform ambient network checks or prompt for updates. Explicit version and
upgrade commands should provide structured data, clear human guidance, and safe
application only when the installation channel is known.

## Scope

Covered: structured version metadata, latest-version checks, install-method
detection, advisory and applied upgrade flows, JSON output, npm launcher
marking, GitHub-hosted Unix/macOS and Windows installer scripts, an optional
Windows command launcher, and the long-term managed standalone installer
contract.

Deferred / non-goals: no automatic update checks during ordinary commands, no
interactive TUI prompt, no self-replacement for unknown binaries, no
package-manager replacement when the package manager should own the install,
and no change to format-spec or evaluation semantics.

## Requirements

`qualitymd --version` **MUST** remain available and continue to report the
human-oriented binary version through the existing invocation harness.

`qualitymd version --json` **MUST** emit structured version metadata including
at least the CLI version, commit when known, development-build status, and the
bundled `SPECIFICATION.md` version.

`qualitymd version` without `--json` **SHOULD** provide the same facts in a
concise human-readable form. It **MUST NOT** contact the network.

`qualitymd upgrade --check` **MUST** explicitly check whether a newer compatible
CLI release is available. It **MUST NOT** run as a side effect of unrelated
commands.

`qualitymd upgrade --check --json` **MUST** emit a stable machine-readable
result with the current version, latest known version when available, detected
install method, update availability, recommended action, and whether automatic
application is supported for this installation.

`qualitymd upgrade` without `--apply` **SHOULD** behave as an advisory workflow:
it may perform the same latest-version check as `--check`, then print the exact
recommended command or managed-upgrade action without mutating the install.

`qualitymd upgrade --apply` **MUST** apply an update only when the detected
install method has a safe, explicit update action. When the method is unknown,
source-built, archive-extracted, or otherwise unmanaged, it **MUST** refuse to
mutate the binary and print manual upgrade guidance.

> Rationale: arbitrary binaries might be owned by a package manager, a system
> administrator, a pinned toolchain, or a test fixture. Replacing them from
> inside the binary would be surprising and unsafe. - 0032

Install-method detection **MUST** distinguish at least managed standalone, npm,
Homebrew, Go/source/custom, and unknown installs when enough evidence is
available.

The npm launcher **SHOULD** mark npm-managed invocations so the native binary
can detect that npm owns the install. Detection **MUST NOT** rely only on
fragile path guesses when the launcher can provide an explicit marker.

For npm-managed installs, the upgrade workflow **MUST** recommend or execute the
npm package update command for `quality.md` rather than replacing the native
binary directly.

For Homebrew-managed installs, the upgrade workflow **MUST** recommend or
execute the Homebrew upgrade command for the cask/tap rather than replacing the
staged binary directly.

Latest-version checks **SHOULD** compare against the channel that owns the
installation when that channel can lag the GitHub Release. For example, a
Homebrew-managed install should be compared against the Homebrew-published
version when available.

A managed standalone installer **SHOULD** become the preferred long-term install
path and **SHOULD** be hosted from GitHub release or repository URLs until the
project has a first-party hosting site that serves installer assets directly.
Its install layout **MUST** make ownership detectable by `qualitymd`, and its
upgrade path **MUST** support checksum verification, atomic replacement, PATH
setup or verification, conflict detection, post-install `qualitymd --version`
verification, and non-interactive execution.

The managed standalone installer **MUST** provide a POSIX-compatible
`install.sh` entrypoint for macOS and Linux and a PowerShell `install.ps1`
entrypoint for Windows. Both scripts **MUST** support a non-interactive mode and
an explicit version or release selector, with `latest` as the default.

The managed standalone installer **MAY** provide an `install.cmd` launcher for
Windows users who start from `cmd.exe`. If present, the command launcher
**MUST** delegate to `install.ps1` rather than carrying separate install logic.
The project **SHOULD NOT** maintain a separate `.bat` implementation unless a
real supported environment cannot run `.cmd`.

> Rationale: PowerShell is the maintainable Windows implementation surface, but
> `cmd.exe` remains a common copy/paste context. A thin `.cmd` launcher improves
> ergonomics without creating a second Windows installer to keep in sync. - 0032

Install documentation and upgrade guidance **MUST NOT** use a vanity project
domain as the canonical installer URL unless that domain hosts the installer
assets directly. A redirect-only website is acceptable as a marketing link, but
the copyable install command must use a GitHub URL.

The top-level `install.md` **MUST** remain agent-centered: it should keep the
skill-first flow, then verify or install the deterministic CLI prerequisite,
then run `/quality setup` or `/quality wizard`.

When the managed standalone installer exists, `install.md` **MUST** present the
GitHub-hosted managed installer as the preferred CLI install path. It **MUST**
show separate macOS/Linux and Windows commands, and it **MUST** include the
non-interactive and explicit-version/pinned forms needed by agents and CI.

`install.md` **MUST** keep npm, Homebrew, and `go install` as alternatives or
fallbacks, and it **MUST** explain that `qualitymd upgrade --check` can be used
to inspect current/latest version and the recommended upgrade action.

When a managed standalone install is detected, `qualitymd upgrade --apply`
**MAY** run the managed updater directly. If the managed installer is not yet
available, the command **MUST** report that direct application is unsupported and
provide the best available manual install guidance.

Upgrade commands **MUST** use the established stdout/stderr and exit-code
contract: payload and human result on stdout, diagnostics on stderr, usage
errors as usage failures, and failed checks or failed update application as
internal/precondition failures rather than silent success.

Upgrade checks and version output **MUST** be safe for agents and CI: no
interactive prompts, no secret output, stable field names under `--json`, and
deterministic output for the same local state and fetched release metadata.

## Example JSON shape

```json
{
  "currentVersion": "0.4.1",
  "commit": "abc1234",
  "specificationVersion": "0.1",
  "developmentBuild": false,
  "installMethod": "npm",
  "latestVersion": "0.5.0",
  "updateAvailable": true,
  "applySupported": true,
  "recommendedCommand": "npm install -g quality.md@latest"
}
```

The final schema can differ, but it must carry equivalent facts with stable
machine-readable field names.

## Durable spec changes

### To add

- `specs/cli/version.md` - specify structured version metadata and its human and
  JSON forms (per the version requirements above).
- `specs/cli/upgrade.md` - specify upgrade checks, install-method detection,
  advisory output, `--apply`, JSON output, and safe refusal behavior (per the
  upgrade requirements above).

### To modify

- `specs/cli.md` - register the new version and upgrade command contracts in the
  CLI overview and reconcile them with existing `--version` behavior (per the
  version and upgrade requirements above).
- `specs/cli/index.md` - list the new command sub-specs (per the added durable
  specs above).

### To delete

None
