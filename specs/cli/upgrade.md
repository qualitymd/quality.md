---
type: Functional Specification
title: qualitymd upgrade
description: Explicitly check for and advise on qualitymd CLI upgrades.
tags: [cli, command, upgrade, versioning]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd upgrade

`qualitymd upgrade` is the explicit network boundary for checking whether a CLI
upgrade is available and which install channel should own the upgrade.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

Ordinary `qualitymd` commands **MUST NOT** perform background upgrade checks.
Network access for update discovery **MUST** occur only through explicit upgrade
commands.

`qualitymd upgrade --check` and advisory `qualitymd upgrade` **MUST** report the
current version, latest known version when available, detected install method,
whether an update is available, whether direct apply is supported, and the
recommended action. An update is available only when the latest version is
strictly newer than the current version by SemVer precedence; a development build
**MUST NOT** report an available update. When either version is not valid SemVer,
the check **MAY** fall back to reporting any difference.

`qualitymd upgrade --json` **MUST** emit those facts with stable field names.

Install-method detection **MUST** distinguish at least npm, Homebrew, managed
standalone, Go/source, archive, and unknown installs when evidence is available.
Where a launcher can mark ownership, the command **SHOULD** prefer that explicit
marker over path guessing.

`qualitymd upgrade --apply` **MUST** mutate the installation only when the
detected owner channel has an explicit safe command. npm installs **MUST** be
upgraded through `npm install -g quality.md@latest`; Homebrew installs **MUST**
be upgraded through the documented Homebrew command. Unknown, archive, source,
and unsupported managed standalone installs **MUST** refuse direct mutation and
print manual guidance.

Managed standalone installers **MUST** write ownership metadata that makes the
install detectable, and their update path **SHOULD** verify checksums, stage the
replacement, set up PATH so the installed command is reachable, switch the
visible command atomically where the platform allows, and verify
`qualitymd --version` after install. PATH setup **SHOULD** follow platform
convention: update the per-user PATH directly where the platform offers a safe
API (Windows), and otherwise print the exact line to add rather than editing
shell profiles.
