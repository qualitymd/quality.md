---
type: Change Case
title: Skill release metadata
description: Record the /quality skill version and required qualitymd CLI range in Agent Skills-compatible metadata, then validate release-note drift mechanically.
status: Done
tags: [skill, versioning, release]
timestamp: 2026-06-19T00:00:00Z
---

# Skill release metadata

A **Change Case** capturing the _why_ and _status_ for making the `/quality`
skill's release identity and CLI prerequisite explicit in the installable skill
artifact. The detail lives in its
[functional spec](0034-skill-release-metadata/spec.md).

> **Done.** Implementation and durable artifact synchronization are complete; the change is archived.

## Motivation

The `/quality` skill is a separately versioned surface, but its current version
and supported `qualitymd` CLI range are only described in release notes and
documentation. That leaves the release guide honest but still manual: a release
author has to remember to keep the skill prompt, versioning docs, and changelog
compatibility block aligned.

The Agent Skills specification already permits arbitrary string metadata, so the
project can record a project-owned skill version and CLI requirement without
waiting for installer-level dependency enforcement. This should make the
installable `SKILL.md` the source of truth for the skill release contract while
being clear that current installers may ignore the metadata.

## Scope

Covered: add `/quality` skill metadata for the skill SemVer and required
`qualitymd` CLI range; use `compatibility` for matching human-readable prose;
update the skill spec, release guide, and versioning docs; update release checks
to catch drift between `SKILL.md` metadata and curated release notes; and update
runtime skill wording to use the metadata range for released installs.

Deferred / non-goals: no installer or registry enforcement, no new upstream
Agent Skills dependency schema, no CLI self-upgrade behavior, no change to
`QUALITY.md` document format semantics, and no fine-grained command capability
versions.

Implementation is complete and archived.

## Affected specs & docs

- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md) - require the `/quality` installable artifact to carry project-owned
      Agent Skills metadata for the skill version and required `qualitymd` CLI
      range.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - add the
      metadata and update CLI prerequisite wording to name the metadata range.
- [x] [`docs/reference/versioning.md`](../../docs/reference/versioning.md) - make
      `SKILL.md` metadata the source of truth for the `/quality` skill release
      version and supported CLI range, with release notes as a mirror.
- [x] [`docs/guides/cut-a-release.md`](../../docs/guides/cut-a-release.md) -
      update release preparation and process-boundary guidance for skill
      metadata and deferred installer enforcement.
- [x] [`docs/guides/use-quality-skill.md`](../../docs/guides/use-quality-skill.md) - point prerequisite checks at the metadata-declared CLI range.
- [x] [`scripts/check-release.mjs`](../../scripts/check-release.mjs) and related
      test coverage - validate skill metadata and changelog compatibility drift.
- [x] [`CHANGELOG.md`](../../CHANGELOG.md) - record the user-facing release-note
      and compatibility-block convention.

No `SPECIFICATION.md` update is expected: this changes skill release metadata
and release mechanics, not the `QUALITY.md` document format or evaluation
semantics.

## Children

- [Functional spec](0034-skill-release-metadata/spec.md) - what skill release
  metadata and release validation must provide.
- [Design doc](0034-skill-release-metadata/design.md) - how the metadata,
  release-check validation, runtime wording, and installer boundary fit together.

## Status

`Done`. The change is archived.
