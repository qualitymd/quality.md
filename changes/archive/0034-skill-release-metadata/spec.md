---
type: Functional Specification
title: Skill release metadata - functional spec
description: Define /quality skill release metadata, CLI prerequisite metadata, and release validation.
tags: [skill, versioning, release]
timestamp: 2026-06-19T00:00:00Z
---

# Skill release metadata - functional spec

Companion to [Skill release metadata](../0034-skill-release-metadata.md). This
spec states the release-metadata and validation delta for the `/quality`
installable skill.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The `/quality` skill depends on a compatible `qualitymd` CLI, and the project
already treats the skill as its own versioned surface. Without machine-readable
metadata in the skill artifact, the supported CLI range exists only as prose in
release notes and docs. That is enough for early releases, but it is easy for
release prep to drift and impossible for the skill to point at one local source
of truth.

Agent Skills frontmatter permits arbitrary string metadata. The `/quality` skill
can use that portable extension point now while explicitly leaving installer
dependency enforcement for a later package-manager or registry contract.

## Scope

Covered: required project-owned `metadata` keys in `skills/quality/SKILL.md`,
matching `compatibility` prose, release-note mirroring, release-check validation,
runtime prerequisite wording, and documentation updates.

Deferred / non-goals: no installer enforcement, no registry dependency
semantics, no top-level custom `version` or `requires` fields outside the Agent
Skills metadata map, no change to CLI upgrade behavior, and no separate
capability versions for individual commands or artifacts.

## Requirements

The `/quality` installable `SKILL.md` **MUST** continue to satisfy the Agent
Skills required frontmatter fields `name` and `description`.

The `/quality` installable `SKILL.md` **MUST** declare `metadata.version` as a
string containing the `/quality` skill release SemVer without a leading `v`.
Repository release tags keep the leading `v`; skill metadata does not.

The `/quality` installable `SKILL.md` **MUST** declare
`metadata.requires-qualitymd-cli` as a string containing the supported
`qualitymd` CLI SemVer range for released installs.

The `/quality` installable `SKILL.md` **MUST** declare `compatibility` prose that
names the required `qualitymd` CLI range. The range in `compatibility` **MUST**
match `metadata.requires-qualitymd-cli`.

> Rationale: `metadata.requires-qualitymd-cli` is the machine-readable contract;
> `compatibility` keeps the same prerequisite visible to agents and clients that
> display only standard Agent Skills fields. - 0034

The metadata keys **MUST** be project-owned Agent Skills metadata keys:
`version` and `requires-qualitymd-cli`. The project **MUST NOT** add custom
top-level `version`, `requires`, or dependency fields to `SKILL.md` unless a
future Agent Skills spec or installer contract defines them.

> Rationale: the Agent Skills spec defines `metadata` as the extension point.
> Staying inside it keeps the skill portable while avoiding a fake dependency
> standard. - 0034

Release notes for a release that publishes or changes the `/quality` skill
**MUST** mirror the skill metadata version and required CLI range in the curated
`CHANGELOG.md` compatibility block.

`mise run release-check -- <tag>` **MUST** fail when the skill metadata version
does not match the release tag without its leading `v`.

`mise run release-check -- <tag>` **MUST** fail when the changelog `/quality`
skill compatibility line is present but does not match `SKILL.md`
`metadata.version` and `metadata.requires-qualitymd-cli`.

`mise run release-check -- <tag>` **SHOULD** validate that
`metadata.requires-qualitymd-cli` is a plausible SemVer range. It need not
implement a complete SemVer solver in this change.

The `/quality` skill prompt and mode guidance **MUST** identify
`metadata.requires-qualitymd-cli` as the supported CLI range for released
installs. Local development builds **MAY** be accepted when the commands the
skill depends on are present.

When structured CLI version output from the CLI managed-upgrades work is
available, the skill **SHOULD** prefer that structured output for prerequisite
checks. Until then, `qualitymd --version` remains acceptable for the basic
version check.

Documentation **MUST** make clear that current skill metadata is project-owned
and release-validated, not installer-enforced. Installer or registry enforcement
is deferred until an installer/package contract supports it.

## Example frontmatter

```yaml
---
name: quality
description: "Setup or work with QUALITY.md files or the qualitymd CLI; model, evaluate, or improve project or harness quality, get wizard quality advice, anything concerning quality factors/attributes/characteristics relevant to project context"
compatibility: Requires qualitymd CLI >=0.3.0 <0.4.0.
metadata:
  version: "0.3.0"
  requires-qualitymd-cli: ">=0.3.0 <0.4.0"
---
```

The concrete version and range are release-prep values; the shape above is the
contract.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/quality-skill.md` - require the project-owned
  Agent Skills metadata keys, matching `compatibility` prose, and runtime
  prerequisite wording (per the metadata and runtime requirements above).

### To delete

None
