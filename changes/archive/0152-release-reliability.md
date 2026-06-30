---
type: Change Case
title: Release Reliability
description: Harden release preflight, publish workflow, verification, and repair behavior after v0.25.0 exposed partial-publish failure modes.
status: Done
tags: [release, automation, packaging, reliability]
timestamp: 2026-06-27T00:00:00Z
---

# Release Reliability

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0152-release-reliability/spec.md) - what the case must do.
- [Design doc](0152-release-reliability/design.md) - how it is built, and why.

## Motivation

The `v0.25.0` release exposed a fragile release pipeline: GoReleaser created a
GitHub release and uploaded assets, then failed while updating the Homebrew tap
because the tap token was stale. npm publish was sequenced after that failure, so
npm stayed on the old version until the release workflow was repaired and rerun.
The partial GitHub release also briefly carried generated commit-list notes
instead of curated changelog notes.

Release publishing should fail before tagging when credentials are bad, keep
GitHub releases draft until downstream channels succeed, let npm and Homebrew
publish independently after artifacts are built, and provide a single verifier
and repair path for partial releases.

## Scope

Covered:

- Add a release preflight task that checks target availability and publish
  credentials before tagging.
- Split release workflow jobs so GitHub archive creation, Homebrew cask update,
  npm publish, and final verification are independently visible.
- Create GitHub releases as drafts and publish them only after Homebrew and npm
  succeed.
- Move Homebrew tap updates out of GoReleaser into an idempotent repo script.
- Make npm publishing skip package versions that are already published.
- Add a post-release verifier for GitHub assets, curated notes, npm package
  versions, Homebrew cask version/checksums, and a downloaded binary.
- Add a repair task that can reapply curated notes, update Homebrew, publish
  missing npm packages, and verify.
- Update the release guide and logs.
- Cut a patch release proving the hardened process works.

Deferred:

- Replacing the Homebrew tap token with a GitHub App installation token. The
  automation will document the token contract and verify it, but creating the
  GitHub App is external account work.
- Full installer smoke tests for every platform. The verifier checks one native
  downloaded archive on the runner platform.
- Changing package names or distribution channels.

## Affected artifacts

**Code / automation**

- [x] `.github/workflows/release.yml` - split publish jobs, draft release, final
      verify/publish.
- [x] `.goreleaser.yaml` - remove Homebrew publishing from GoReleaser.
- [x] `mise.toml` - add release preflight, verify, repair, and Homebrew update
      tasks.
- [x] `scripts/check-release.mjs` - integrate preflight or keep it explicitly
      paired with release-check.
- [x] `scripts/build-npm.mjs` - add skip-existing publish behavior.
- [x] `scripts/release-preflight.mjs` - new credential and target preflight.
- [x] `scripts/release-verify.mjs` - new post-release verifier.
- [x] `scripts/update-homebrew-cask.mjs` - new idempotent tap updater.
- [x] `scripts/repair-release.mjs` - new partial-release repair wrapper.

**Durable docs / specs**

- [x] `docs/guides/cut-a-release.md` - update release procedure, credential
      contract, failure playbooks, and verification.
- [x] `CHANGELOG.md` - release note for the release reliability improvements.
- [x] `changes/index.md`, `changes/log.md`, `changes/archive/index.md` -
      Change Case lifecycle.

**No planned impact**

- [x] `SPECIFICATION.md`, Evaluation specs, runtime `/quality` evaluation
      behavior, and CLI command behavior outside release-support scripts.

## Status

`Done`. Implemented, documented, shipped, and verified through the successful
`v0.25.4` release workflow.
