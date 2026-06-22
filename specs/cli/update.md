---
type: Functional Specification
title: qualitymd update
description: Apply or check for qualitymd CLI updates through managed install channels.
tags: [cli, command, update, versioning]
timestamp: 2026-06-20T00:00:00Z
---

# qualitymd update

`qualitymd update` is the explicit update surface for the CLI install channel
that owns the visible binary. It applies by default and uses `--check` for
advisory, non-mutating checks.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Command Shape

`qualitymd update` invoked without a mode flag **MUST** attempt to apply an
available update through the install channel that owns the binary. When no newer
ready release is available, it succeeds without mutating. When the detected
install method has no safe apply action, it **MUST** refuse to mutate the install
and print manual guidance.

`qualitymd update --check` **MUST** report, without mutating, the current
version, latest known version when available, detected install method, whether
the latest version is ready for that channel, whether an update is available,
whether apply is supported, the recommended action, and the release-notes
reference when known.

`--json` selects machine-readable output and **MUST NOT** change whether
`update` applies. `qualitymd update --json` applies by default and emits a
structured result of that apply; `qualitymd update --check --json` reports facts
without mutating. Both forms **MUST** include stable field names and a stable
boolean field reporting whether an update was applied.

An update is available only when both current and latest versions are valid
SemVer releases, the latest is strictly newer by SemVer precedence, the current
build is not a development build, and the release readiness gate is satisfied.
When either version is not valid SemVer, either version is a prerelease, or the
build is a development build, `update` **MUST** report no available update and
**MUST NOT** apply.

There is no `upgrade` alias.

## Install Ownership and Apply

Install-method detection **MUST** distinguish at least npm, Homebrew, managed
standalone, Go/source, archive, and unknown installs when evidence is available.
Where a launcher can mark ownership, the command **SHOULD** prefer that explicit
marker over path guessing.

npm installs **MUST** apply through `npm install -g quality.md@latest`.
Homebrew installs **MUST** apply through the documented Homebrew command.
Managed standalone installs **MUST** apply by invoking the managed installer
non-interactively. Unknown, archive, and source installs **MUST** refuse direct
mutation and print manual guidance.

After applying, `update` **MUST** verify the visible `qualitymd --version` and
**MUST** report honestly when the visible version did not advance to the target
release.

Managed standalone installers **MUST** write ownership metadata that makes the
install detectable, and their update path **SHOULD** verify checksums, stage the
replacement, set up PATH so the installed command is reachable, switch the
visible command atomically where the platform allows, and verify
`qualitymd --version` after install. PATH setup **SHOULD** follow platform
convention: update the per-user PATH directly where the platform offers a safe
API (Windows), and otherwise print the exact line to add rather than editing
shell profiles.

## Release Readiness and Notes

An update **MUST NOT** be reported available, and `update` **MUST NOT** apply,
unless the area release is retrievable for the owning channel. For the
GitHub-backed managed standalone channel, readiness requires the platform
archive and checksums. For registry-backed channels, the owning registry is the
source of both latest version and readiness: npm uses the npm registry, and
Homebrew uses the published tap cask rather than a GitHub release tag.

When readiness cannot be confirmed, `update --check` **MUST** report no available
update and **SHOULD** make the unconfirmed state legible. `update` **MUST** fail
before mutating when a newer version is known but not ready.

`update --check`, an applied `update`, and the ambient notice **SHOULD** include
a reference to the area release's notes when one is known. `update --json`
**MUST** carry that reference under a stable field, omitted when not known. The
reference is advisory and **MUST NOT** change availability or apply behavior.

## Ambient Notice

The cross-command ambient update notice contract lives in
[`qualitymd update notice`](update-notice.md). The explicit `qualitymd update`
command remains the owner of apply and `--check` behavior; ordinary commands may
only surface cached, non-blocking notices under the shared notice contract.
