---
type: Reference
title: Versioning
description: Versioned surfaces and compatibility policy for the CLI, /quality skill, and QUALITY.md specification.
tags: [versioning, release, cli, skill, specification]
timestamp: 2026-06-19T00:00:00Z
---

# Versioning

This project has three separately versioned surfaces:

- the `qualitymd` CLI;
- the `/quality` skill; and
- the QUALITY.md format specification.

The CLI and skill are distributed separately. They can be released from the same
repository state, but the skill is not bundled with the CLI binary.

## Where Versions Are Recorded

- CLI release version: the git release tag and build metadata reported by
  `qualitymd --version`.
- Skill release version and supported CLI range:
  `skills/quality/SKILL.md` frontmatter metadata, mirrored in release notes and
  install documentation.
- Specification version: the version line near the top of `SPECIFICATION.md`.

## CLI Version

The `qualitymd` CLI uses SemVer and reports its version through:

```sh
qualitymd --version
qualitymd version --json
```

The CLI version is the runtime compatibility boundary for the skill. It covers
the deterministic surface the skill depends on:

- commands and subcommands;
- flags and arguments;
- exit-code categories;
- machine-readable output shapes;
- evaluation artifact mechanics; and
- the `SPECIFICATION.md` version bundled into the binary and emitted by
  `qualitymd spec` and `qualitymd version --json`.

Update checks are explicit:

```sh
qualitymd update --check
```

Ordinary commands may show an ambient update notice from a local cache and may
refresh that cache in a bounded, best-effort background process. They do not
block on network update checks, and the notice is suppressed outside interactive
stderr, in CI, for `--json`, for development builds, and when
`QUALITYMD_NO_UPDATE_CHECK=1` is set. The same environment variable disables
explicit update checks for centrally managed or air-gapped installs.

The update surface reports the detected install method, latest known version,
release readiness, whether apply is supported, and the recommended action.
`qualitymd update` applies by default for managed standalone, npm, and Homebrew
installs, then verifies the visible `qualitymd --version`. Source, archive, and
unknown installs keep distinct ownership boundaries so the CLI does not replace
binaries it does not own.

While the CLI is in `0.x`, minor versions define compatibility lines. A breaking
change to the skill-facing CLI surface should bump the minor version; patch
versions should preserve compatibility within the same minor line. After `1.0`,
normal SemVer major-version compatibility rules apply.

## Skill Version

The `/quality` skill has its own SemVer because it is installed and upgraded
separately from the CLI.

Each skill release declares the `qualitymd` CLI SemVer range it supports in
`skills/quality/SKILL.md`:

During the `0.x` phase, skill releases should usually depend on one compatible
CLI minor line, for example:

```yaml
compatibility: Requires qualitymd CLI >=0.5.0 <0.6.0.
metadata:
  version: "0.5.0"
  requires-qualitymd-cli: ">=0.5.0 <0.6.0"
```

At runtime, the skill uses `metadata.requires-qualitymd-cli` as the released
install range, checks the installed CLI with `qualitymd version --json` when
available, and stops when the CLI is missing or outside the supported range.

Use a broader range only when the skill has been checked against every included
CLI minor line. Current skill installers may ignore this project-owned metadata;
release checks validate it, and release notes mirror it, but installer
enforcement is deferred until an installer/package contract supports it.

## Paired Skill and CLI Upgrades

Use `/quality update` to maintain an existing `/quality` installation. The mode
plans the skill and CLI pair together: it reads the installed skill metadata,
checks the visible CLI, uses `qualitymd update --check` for CLI owner-channel
guidance, and delegates skill changes to the Agent Skills installer when that
installer exposes an update path.

Compatibility is the hard gate: the visible CLI must satisfy the target skill's
`metadata.requires-qualitymd-cli` range before setup, evaluation, or improvement
work can proceed. Latest-version updates are applied only after the user confirms
the `/quality update` plan.

Skill updates may not take effect in the currently running agent session. After
the skill package changes, restart, reload, or start a new session when the
agent discovers skills only at startup or caches loaded skill instructions.

## Specification Version

`SPECIFICATION.md` carries a specification version for the QUALITY.md document
format, frontmatter schema, evaluation semantics, and required report semantics.

The specification version changes when the meaning of conforming QUALITY.md
documents or conforming evaluation/report behavior changes. Editorial fixes,
clarifications that do not change conformance, and non-normative examples do not
require a specification version change.

Each CLI release embeds one specification version. `qualitymd spec` emits the
specification bundled with that CLI release.

## Release Tags

Repository release tags coordinate the repository contents at release time. A
tag can publish CLI binaries, npm packages, Homebrew artifacts, and skill
contents from the same source state, but the tag does not make the skill part of
the CLI distribution.

When a release changes more than one versioned surface, update every affected
version and compatibility declaration in the same change.

## Compatibility Policy

The project does not currently maintain fine-grained capability versions. The
skill uses the CLI SemVer range as the compatibility contract.

If that contract becomes too coarse, the project can add machine-readable
capability reporting later. Until then, do not introduce per-command or
per-artifact capability versions.
