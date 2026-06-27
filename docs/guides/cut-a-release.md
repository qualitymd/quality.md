---
type: How-to Guide
title: Cut a release
description: How to prepare, verify, tag, publish, and check a QUALITY.md release.
tags: [release, changelog, versioning]
timestamp: 2026-06-27T00:00:00Z
---

# Cut a Release

Use this guide when publishing a tagged release of the `qualitymd` CLI and any
same-repo release notes for the `/quality` skill or QUALITY.md specification.

This guide is the authoritative release runbook. Other docs should link here
instead of restating release-prep, tag-publish, verification, or failure-handling
steps.

## Preconditions

Before cutting a release:

- The main branch contains the intended release state.
- CI is passing.
- The working tree is clean.
- `GITHUB_TOKEN` or `GH_TOKEN` is available locally for release preflight.
- `NPM_TOKEN` is configured locally and as a repository secret for npm
  publishing.
- `HOMEBREW_TAP_GITHUB_TOKEN` is configured locally and as a repository secret
  for the Homebrew tap.
- You can run the local dry-run tasks.

Check the current state:

```sh
git status --short
git fetch --tags
git tag --list --sort=-version:refname | head
```

The Homebrew token must be able to write to `qualitymd/homebrew-tap`. The npm
token must be able to publish the `quality.md` package and every package in the
`@qualitymd` scope. Rotate either token before tagging if preflight cannot prove
that access.

## Choose the Release Version

The CLI release version is the git tag, for example `v0.3.0`.

While the CLI is still `0.x`:

- Use a patch release for compatible fixes, docs, packaging fixes, and compatible
  command behavior changes.
- Use a minor release for new user-visible CLI capability.
- Use a minor release for breaking changes to the skill-facing CLI surface.
- Avoid broadening the `/quality` skill supported CLI range unless the skill has
  been checked against every included CLI minor line.

After `1.0`, use normal SemVer major/minor/patch rules.

## Check Versioned Surfaces

This repo has three separately versioned surfaces:

- `qualitymd` CLI
- `/quality` skill
- QUALITY.md specification

Before tagging, decide which surfaces changed.

Update the CLI release notes when the release changes:

- commands, flags, arguments, or exit behavior;
- machine-readable output;
- evaluation artifact behavior;
- packaging or installation behavior;
- embedded `qualitymd spec` output.

Update the `/quality` skill release notes or compatibility declaration when the
release changes:

- skill instructions;
- skill workflows;
- bundled skill resources;
- required CLI command surface;
- supported CLI SemVer range.

Update the `SPECIFICATION.md` specification version only when the meaning of
conforming QUALITY.md documents or conforming evaluation/report behavior
changes. Editorial fixes and non-normative examples do not require a spec version
change.

When a release changes more than one surface, update every affected version and
compatibility declaration in the same release-prep change.

## Prepare Release Notes

Use the bundle logs as source material:

- `changes/log.md`
- `docs/log.md`
- `specs/log.md`

Release notes are user-facing. Include outcomes and compatibility impact, not
every internal process step.

The root `CHANGELOG.md` is the canonical source for curated release notes. Keep
an `Unreleased` section at the top and one section per tagged release. During
release preparation, move the relevant `Unreleased` notes into a new tag section
and use that tag section as the GitHub Release body.

Suggested release-note groups:

```md
## v0.3.0 - YYYY-MM-DD

### CLI

### /quality Skill

### Specification

### Documentation

### Packaging

### Compatibility / Migration
```

Each release should include a compatibility block when relevant:

```md
Compatibility:

- CLI: v0.3.0
- QUALITY.md specification: 0.1
- /quality skill: compatible with qualitymd >=0.3.0 <0.4.0
```

Preview the exact notes that the release workflow will publish:

```sh
mise run release-notes -- v0.3.0
```

## Commit Release Preparation

Commit the release-prep changes before tagging.

Common release-prep changes include:

- `CHANGELOG.md`;
- `SPECIFICATION.md` version line, when applicable;
- skill compatibility metadata or docs, when applicable;
- install or release documentation, when applicable.

Use a direct message such as:

```sh
git commit -m "Prepare v0.3.0 release"
```

## Run Release Check

After committing the release-prep change and before tagging, run:

```sh
mise run release-check -- v0.3.0
```

The release check:

- requires a clean working tree before and after the check;
- refuses to continue if the local or remote tag already exists;
- verifies `CHANGELOG.md` has `Unreleased` and target release sections;
- verifies the target release notes are substantive;
- checks obvious compatibility-block drift against the tag and
  `SPECIFICATION.md` version;
- verifies `skills/quality/SKILL.md` `metadata.version` matches the tag without
  `v`, that `metadata.requires-qualitymd-cli` has the expected range shape, and
  that any `/quality skill` changelog compatibility line mirrors the metadata;
- runs `mise run fmt`, `mise run test`, `mise run vet`, `mise run snapshot`, and
  `mise run npm-build`.

Inspect generated artifacts enough to catch packaging mistakes before pushing the
tag.

Pre-tag Goreleaser snapshots can mention the latest existing tag in their log
while `release-check` is validating the candidate tag. That is expected before
the new tag exists; trust the final `release checks passed for <tag>` line and
the release-note preview for the candidate version.

Then run release preflight:

```sh
mise run release-preflight -- v0.3.0
```

The preflight check verifies that the target tag and GitHub Release do not
already exist, that npm does not already contain the candidate package versions,
that npm credentials authenticate, and that the Homebrew tap token can create
and delete a temporary branch.

After `release-check` passes, push the release-prep commit and wait for the
hosted `main` CI run for that exact commit to pass before creating the tag:

```sh
git push origin main
gh run list --branch main --limit 5
```

## Tag and Publish

Create and push the tag:

```sh
git tag v0.3.0
git push origin v0.3.0
```

The release workflow runs from `.github/workflows/release.yml`. It publishes:

- a draft GitHub Release with the curated body from the matching
  `CHANGELOG.md` section;
- GitHub release archives through Goreleaser;
- Homebrew cask updates through `scripts/update-homebrew-cask.mjs`;
- npm / npx packages through `scripts/build-npm.mjs`;
- a final verification pass through `scripts/release-verify.mjs`;
- publication of the GitHub Release only after Homebrew, npm, release notes,
  checksums, and a downloaded binary verify.

The workflow requires two repository secrets in addition to the built-in
`GITHUB_TOKEN`:

- `NPM_TOKEN` — npm automation token with publish access to the `quality.md`
  package and the `@qualitymd` scope.
- `HOMEBREW_TAP_GITHUB_TOKEN` — token with write access to
  `qualitymd/homebrew-tap`.

Homebrew distribution uses a **cask**, not a formula: that is the
recommended path for a self-published pre-built binary, so a binary-only formula
is not the goal here. The cask currently strips the macOS quarantine attribute
because the binaries are unsigned. Remove that step from the tap cask once the
binaries are signed and notarized.

Do not move a published tag to repair a release. If a published release is wrong,
fix forward with a new patch release unless no artifacts were published and the
tag is still only local or clearly failed before publication.

## Verify Published Artifacts

After the workflow finishes, verify:

- The GitHub Release exists for the tag.
- Release notes match the curated `CHANGELOG.md` release-note entry.
- Checksums are present.
- The Homebrew tap updated.
- npm packages published for the supported platforms.
- `npx quality.md --version` reports the expected version.
- `qualitymd version --json` reports the expected CLI version, commit when
  stamped, and bundled `SPECIFICATION.md` version.
- `qualitymd update --check --json` returns a structured check result for at
  least one fresh install channel.
- A downloaded archive binary reports the expected version.
- `qualitymd spec` emits the expected bundled specification.
- Public installer entrypoints under `install/` are present in the release
  source and use GitHub-hosted release assets.

Use fresh installs where practical rather than already-built local binaries.
A release-note comparison that differs only by a final trailing newline is
equivalent and does not require a release edit.

The mechanical channel verification is:

```sh
mise run release-verify -- v0.3.0
```

## Handle a Failed Release

If the workflow fails before publishing anything, fix the problem and rerun the
workflow or replace the tag only if it has not been consumed externally.

If the workflow created a draft GitHub Release and failed before publishing the
release, fix the issue and rerun the workflow for the same tag. The workflow
deletes a stale draft release before rebuilding it, but refuses to overwrite an
already-published release.

If the workflow partially publishes artifacts:

- Do not rewrite the tag.
- Record what published and what failed.
- Fix the release automation or package issue.
- Prefer `mise run release-repair -- <tag>` when the GitHub release, Homebrew
  cask, and npm package versions are already intended for that same tag.
- Publish a new patch tag when the wrong artifact content was published or when
  the version has been consumed in a state that cannot be repaired.
- Add a release-note entry if users may have consumed the partial release.

If npm published but GitHub or Homebrew failed, treat the version as consumed.
npm versions are immutable in normal practice. The repair helper skips npm
versions that already exist, refreshes the Homebrew cask, restores curated
release notes, and runs release verification.

## After Release

After verifying the release:

- Ensure `CHANGELOG.md` has a fresh `Unreleased` section.
- Check install docs still point to valid commands.
- If the `/quality` skill compatibility range changed, check the skill guide and
  versioning reference are still aligned.
- Close or archive any completed Change Cases whose work shipped in the release.

## Process Support Boundaries

The guide is enough to run the curated-release process. Keep release support
small and add mechanics only when they remove repeated mistakes or manual
comparison work.

Release preparation stays manual. The release-notes, release-check,
release-preflight, release-verify, and release-repair helpers are the automation
boundary for now; add another helper only after repeated release-prep mistakes
show what comparison or validation should become mechanical.

The `/quality` skill currently records project-owned release metadata in
`skills/quality/SKILL.md` under `metadata.version` and
`metadata.requires-qualitymd-cli`, mirrored by `compatibility` prose and curated
release notes. `release-check` validates that metadata, but installer enforcement
remains deferred until an installer/package contract supports dependency checks.
