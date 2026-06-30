---
type: Functional Specification
title: Release Reliability
description: Requirements for hardened release preflight, publish, verification, and repair behavior.
tags: [release, automation, packaging, reliability]
timestamp: 2026-06-27T00:00:00Z
---

# Release Reliability

Companion to the
[Release Reliability](../0152-release-reliability.md) Change Case. This spec
states _what_ the change must do. It defers to
[Cut a release](../../../docs/guides/cut-a-release.md),
the release workflow, npm package scripts, and GoReleaser configuration as the
current release surfaces.

The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The `v0.25.0` release partially published because downstream credentials were
not verified before tagging and because Homebrew, npm, release notes, and GitHub
asset publication were sequenced inside one fragile path. A reliable release
pipeline should fail early on credential or target conflicts, avoid public
partial releases, support safe reruns, and provide one verifier that proves every
distribution channel reached the same tag.

## Scope

Covered: release preflight, GitHub Actions release job topology, Homebrew cask
update automation, npm skip-existing publish behavior, release verification,
partial-release repair, and release guide updates.

Deferred / non-goals: creating external GitHub App credentials, changing package
names, supporting additional distribution channels, or testing every platform
installer in this slice.

## Requirements

R1. `mise run release-preflight -- <tag>` **MUST** verify before tagging that the
target tag/release and npm package versions do not already exist, that GitHub
release credentials can access the main repository, that the Homebrew tap token
can write to `qualitymd/homebrew-tap`, and that the npm token authenticates.

> Rationale: stale Homebrew credentials caused `v0.25.0` to fail after GitHub
> assets were already public.
>
> Durable docs: modify `docs/guides/cut-a-release.md` to require preflight
> before `release-check`/tagging.

R2. The release workflow **MUST** create the GitHub release as a draft with
curated changelog notes supplied at creation time, and **MUST** publish the
release only after Homebrew, npm, and final verification succeed.

> Rationale: public users should not see a partial release or generated
> commit-list notes while downstream publish steps are still pending.
>
> Durable docs: modify `docs/guides/cut-a-release.md` to describe draft-first
> release publication.

R3. The release workflow **MUST** separate GitHub archive publication, Homebrew
cask update, npm publish, and verification into distinct jobs or steps whose
status makes the partial state obvious.

> Rationale: Homebrew failure should not hide whether npm ran, and npm should not
> be blocked by a tap update once release assets exist.
>
> Durable docs: modify `docs/guides/cut-a-release.md` to describe channel-level
> status and recovery.

R4. Homebrew tap updates **MUST** be performed by an idempotent repository script
that reads release checksums, updates `Casks/qualitymd.rb`, commits only when the
cask is stale, and pushes with `HOMEBREW_TAP_GITHUB_TOKEN`.

> Rationale: moving the tap update out of GoReleaser makes the failure mode
> visible, testable, and rerunnable without rebuilding archives.
>
> Durable docs: modify `docs/guides/cut-a-release.md` to name the tap-token
> contract and repair command.

R5. npm publishing **MUST** skip package versions that already exist when
requested by release automation.

> Rationale: a rerun after partial npm publication should publish missing
> packages without failing on already-published ones.
>
> Durable docs: none beyond release guide recovery steps.

R6. `mise run release-verify -- <tag>` **MUST** verify the GitHub release state,
expected assets, curated notes, npm launcher and platform package versions,
Homebrew cask version and checksums, and one downloaded archive binary reporting
the expected version.

> Rationale: release completion should be proven by distribution-channel state,
> not inferred from a green workflow alone.
>
> Durable docs: modify `docs/guides/cut-a-release.md` to make the verifier the
> final local and workflow gate.

R7. `mise run release-repair -- <tag>` **MUST** be safe to run after partial
publication: it must reapply curated release notes, update Homebrew if stale,
publish missing npm packages while skipping existing versions, and run the
release verifier.

> Rationale: the `v0.25.0` recovery required manual release deletion, secret
> repair, workflow rerun, manual tap repair, and ad hoc checks.
>
> Durable docs: modify `docs/guides/cut-a-release.md` failure playbooks to prefer
> `release-repair` over tag rewriting.

R8. The release guide **MUST** document credential ownership and rotation
expectations for `HOMEBREW_TAP_GITHUB_TOKEN` and `NPM_TOKEN`.

> Rationale: credentials are operational release dependencies; a stale tap token
> should be discoverable before release day.
>
> Durable docs: modify `docs/guides/cut-a-release.md`.

## Durable spec changes

### To add

None.

### To modify

- `docs/guides/cut-a-release.md` - release preflight, draft-first publication,
  separated channel publish, verification, repair, and credential contracts.

### To rename

None.

### To delete

None.
