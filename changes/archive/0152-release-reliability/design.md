---
type: Design Doc
title: Release Reliability Design
description: Implementation approach for hardened release preflight, publish, verification, and repair behavior.
tags: [release, automation, packaging, reliability]
timestamp: 2026-06-27T00:00:00Z
---

# Release Reliability Design

## Context

This design answers the [Release Reliability functional spec](spec.md). The
release system needs a better split between artifact creation, channel
publication, and verification so one stale credential cannot create a public
partial release.

## Approach

### Workflow topology

The release workflow becomes channel-aware:

```text
github-release
  -> homebrew
  -> npm
  -> verify
  -> publish-release
```

`github-release` checks out the tag, extracts curated changelog notes, deletes
only a stale draft release for that tag when rerunning, and runs GoReleaser with
`--draft --release-notes <file>`. GoReleaser no longer owns Homebrew.

`homebrew` runs the new `update-homebrew-cask.mjs` script against the draft
release checksums and the tap repository.

`npm` runs `build-npm.mjs <version> --publish --skip-existing`.

`verify` runs `release-verify.mjs <tag>` while the release is still draft. It
uses `GITHUB_TOKEN` to inspect draft assets.

`publish-release` runs only after verification succeeds and flips the draft to a
public latest release with `gh release edit <tag> --draft=false --latest`.

### Scripts

`release-preflight.mjs` performs remote checks before tagging. It uses ordinary
GitHub and npm APIs, never prints token values, and proves Homebrew write access
by creating and deleting a temporary branch in `qualitymd/homebrew-tap`.

`update-homebrew-cask.mjs` downloads `checksums.txt` for the tag, clones the tap
with the tap token, rewrites `Casks/qualitymd.rb`, commits only if the file
changed, and pushes. It is safe to rerun when the tap is already current.

`build-npm.mjs` gains `--skip-existing`; before each `npm publish`, it checks
`npm view <package>@<version> version` and skips when the version is already
published.

`release-verify.mjs` verifies the externally visible release state:

- GitHub release exists and has the expected assets;
- release notes match `CHANGELOG.md`;
- npm packages are all at the tag version;
- Homebrew cask version/checksums match `checksums.txt`;
- one native archive downloads and its `qualitymd version --json` reports the
  tag version.

`repair-release.mjs` composes the existing scripts: reapply curated release
notes, run Homebrew update, run npm publish with skip-existing, then verify. It
does not move tags.

### Release check

`release-check` stays the local build and metadata gate. The guide tells
maintainers to run `release-preflight` before `release-check`; keeping them
separate avoids making every local release-check require publish credentials.

## Spec response

- R1 is satisfied by `release-preflight.mjs` and the new `mise` task.
- R2-R3 are satisfied by the workflow topology and GoReleaser draft release.
- R4 is satisfied by `update-homebrew-cask.mjs` and removing Homebrew from
  GoReleaser.
- R5 is satisfied by `build-npm.mjs --skip-existing`.
- R6 is satisfied by `release-verify.mjs` and the workflow `verify` job.
- R7 is satisfied by `repair-release.mjs`.
- R8 is satisfied by release guide updates.

## Alternatives

**Keep Homebrew in GoReleaser.** Rejected. That was the coupling that made a tap
credential failure interrupt the whole release after GitHub assets were already
created.

**Make `release-check` require live credentials.** Rejected. Local release checks
should remain usable without publish credentials; preflight is the explicit
credential gate.

**Publish GitHub releases immediately and repair notes later.** Rejected. Draft
releases avoid user-visible partial states and generated notes.

## Trade-offs & risks

- The workflow has more jobs and scripts, which increases surface area. The gain
  is clearer channel status and safer reruns.
- Draft releases require authenticated verification while still draft. The
  workflow has `GITHUB_TOKEN`; local verification can also verify public releases
  after publication.
- Preflight branch creation is a small, temporary mutation of the tap repository.
  It is deleted immediately and never touches `main`.
