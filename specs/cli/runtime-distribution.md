---
type: Functional Specification
title: qualitymd runtime and distribution
description: One-runtime standalone executable, platform, install-channel, release, and conformance requirements.
tags: [cli, runtime, distribution, release]
timestamp: 2026-07-14T00:00:00Z
---

# qualitymd runtime and distribution

This document specifies the production runtime and distribution boundary for
`qualitymd`. It complements the runtime-independent [CLI contract](../cli.md)
and the procedural [release guide](../../docs/guides/cut-a-release.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## One-runtime boundary

The production CLI **MUST** be one TypeScript application compiled as a
standalone executable. Building, testing, generating, packaging, releasing,
installing, and running it **MUST NOT** require project Go source, a Go
toolchain, GoReleaser, or a project-owned Node or Bun sidecar.

The executable **MUST NOT** require Node.js, Bun, Go, `node_modules`, or other
project dependencies on the target machine. Provider SDKs **MAY** launch their
documented native agent runtime only when its evaluator is selected; that child
is an evaluator implementation detail and is not required by other commands.

The repository **MUST** pin the Bun builder, TypeScript compiler, Effect v4
packages, provider SDKs, formatter, linter, and test runner used by CI and
release.

## Platform and asset matrix

Every release **MUST** publish `checksums.txt` and directly executable archives
for this matrix:

| Target            | GitHub asset                        | Executable      |
| ----------------- | ----------------------------------- | --------------- |
| Darwin arm64      | `qualitymd_darwin_arm64.tar.gz`     | `qualitymd`     |
| Darwin x64        | `qualitymd_darwin_amd64.tar.gz`     | `qualitymd`     |
| Linux arm64 glibc | `qualitymd_linux_arm64.tar.gz`      | `qualitymd`     |
| Linux x64 glibc   | `qualitymd_linux_amd64.tar.gz`      | `qualitymd`     |
| Linux arm64 musl  | `qualitymd_linux_arm64_musl.tar.gz` | `qualitymd`     |
| Linux x64 musl    | `qualitymd_linux_amd64_musl.tar.gz` | `qualitymd`     |
| Windows arm64     | `qualitymd_windows_arm64.zip`       | `qualitymd.exe` |
| Windows x64       | `qualitymd_windows_amd64.zip`       | `qualitymd.exe` |

The original six asset names remain stable. The libc-specific musl assets
extend Linux coverage without changing glibc URLs. Installers and npm launchers
**MUST** select the host libc variant.

Archives **MUST** carry a checksum in `checksums.txt`. Managed installers
**MUST** verify that checksum before replacing a visible executable.

## Install ownership

Supported install owners are:

- npm, through `npm install -g quality.md`, `npx quality.md`, and the
  `qualitymd` and `quality.md` bins;
- Homebrew, through `brew install qualitymd/tap/qualitymd`;
- the managed shell and PowerShell installers; and
- a manually downloaded release archive.

The npm launcher **MUST** delegate to the matching platform package and mark
the child process as npm-owned. Managed installers **MUST** write ownership
metadata. `qualitymd update` uses those markers before path inference.

Go/source is not an install owner. Live guidance **MUST NOT** recommend
`go install`; development builds use the pinned repository toolchain.

## Release provenance and readiness

A release build **MUST** derive the version from the release tag and embed the
version plus source commit in every executable. `qualitymd --version` and
`qualitymd version --json` **MUST** expose that provenance.

The release workflow **MUST**:

1. run credential preflight;
2. build the complete target matrix;
3. publish archives and checksums to a draft GitHub release;
4. publish npm launcher and platform packages;
5. update Homebrew;
6. verify every channel; and
7. publish the draft only after cross-channel readiness succeeds.

Repair **MUST** verify and reuse already published, checksum-identical GitHub
artifacts rather than rebuild them. Missing npm or Homebrew channels may be
completed independently and **MUST** be verified before the release is made
ready.

Bun standalone executables are not assumed byte-reproducible across independent
compiles. Release metadata **MUST** record the builder version and this measured
nondeterminism; archive metadata and ordering **SHOULD** otherwise be normalized.

## Conformance gates

Release candidates **MUST** pass:

- type-check, format, lint, unit, integration, generated-doc, and report-gallery
  checks without Go;
- a complete cross-target build and archive/checksum verification;
- native smoke tests for each supported OS/architecture and both Linux libc
  families;
- npm, managed-installer, Homebrew, and direct-archive smokes;
- CLI compatibility fixtures, including stdout/stderr and exit categories;
- filesystem containment and edge-case source packaging;
- evaluator discovery and executable override;
- structured result streaming, cancellation, and resume; and
- release-repair verification from a partially published state.

## Telemetry and privacy

`qualitymd` **MUST NOT** add project telemetry. It **MUST NOT** enable provider
telemetry on a user's behalf. Selecting an SDK or authenticated provider runtime
may invoke that provider's own documented data behavior; evaluator documentation
must disclose that boundary without representing it as project telemetry.
