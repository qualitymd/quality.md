---
type: Design Doc
title: Skill release metadata design
description: How /quality records skill version and CLI requirements in portable Agent Skills metadata and validates release drift.
tags: [skill, versioning, release, design]
timestamp: 2026-06-19T00:00:00Z
---

# Skill release metadata design

Design behind the
[Skill release metadata](../0034-skill-release-metadata.md) change and its
[functional spec](spec.md).

## Context

The `/quality` skill needs one local release contract that agents, release
authors, and docs can all point at. The Agent Skills spec already provides the
portable extension point: `metadata`, a string key/value map in `SKILL.md`
frontmatter. Current installers may ignore it, so the design needs local
validation without pretending installer dependency resolution exists.

## Approach

Implement the change in four small layers.

### Skill frontmatter

Add the skill SemVer and required CLI range to `skills/quality/SKILL.md`:

```yaml
compatibility: Requires qualitymd CLI >=0.3.0 <0.4.0.
metadata:
  version: "0.3.0"
  requires-qualitymd-cli: ">=0.3.0 <0.4.0"
```

Use `metadata.version` because the Agent Skills spec already shows `version` as
a metadata example. Use `requires-qualitymd-cli` instead of `requires` because
there is no standard dependency schema yet, and because the prerequisite is the
CLI artifact specifically.

### Release validation

Extend `scripts/check-release.mjs` with a narrow frontmatter parser for
`skills/quality/SKILL.md`. The check should extract only the needed fields and
validate them against the tag and curated `CHANGELOG.md` compatibility block.

Keep the validation intentionally shallow:

- exact tag-to-version comparison after stripping the leading `v`;
- exact range comparison between metadata and release notes;
- simple SemVer-range shape validation for the CLI requirement.

This preserves the release helper's current role: catch obvious drift before a
tag, not become a package-manager solver.

### Runtime and docs

Update the root skill prompt to name `metadata.requires-qualitymd-cli` as the
released-install range. Mode docs can keep referring back to the root prompt
rather than duplicating the range.

Update the versioning reference and release guide so `SKILL.md` metadata is the
source of truth. Changelog compatibility blocks and human docs mirror it.

### Installer boundary

Do not add installer or registry behavior in this change. The metadata is
present and mechanically checked at release time, but an installer that ignores
Agent Skills metadata remains conforming to its current behavior. A future
installer can adopt the same fields or migrate them into a standard dependency
schema.

## Alternatives

**Wait for installer metadata.** Rejected. The open Agent Skills spec already
supports arbitrary metadata, and waiting leaves the current release process with
no local source of truth.

**Use top-level `version` and `requires` fields.** Rejected. They would be
readable, but they are not defined by the Agent Skills spec and could conflict
with a future standard dependency schema.

**Use fully prefixed keys such as `qualitymd.skill.version`.** Rejected. They are
unique but too heavy for metadata that lives inside the `quality` skill itself.
`requires-qualitymd-cli` is specific where specificity matters.

**Use only `compatibility` prose.** Rejected. It is useful for display, but it is
not a stable machine-readable release contract.

**Add a sidecar manifest.** Rejected for now. `SKILL.md` frontmatter is already
the always-discovered skill metadata location, and a sidecar would create one
more artifact to keep aligned before installers need it.

## Trade-offs and risks

Project-owned metadata improves release discipline but does not create external
installer enforcement. Users can still install a skill beside an incompatible
CLI until a client or installer chooses to act on the fields.

Parsing YAML frontmatter in the release helper adds one more structured-read
path to a small script. Keeping it scoped to `SKILL.md` and a few scalar string
fields limits that cost.

The metadata version follows the repository release tag today. If the skill ever
ships independently from CLI tags, release tooling will need a mode that
validates a skill-only version without forcing the CLI tag to match.

## Open questions

- Should a future skill-only release use the same repository tag namespace, or a
  skill-specific tag prefix?
- Should installer enforcement consume these exact keys, or migrate to a future
  Agent Skills dependency schema while preserving these as compatibility aliases?
