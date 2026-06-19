---
type: How-to Guide
title: Cut a release
description: How to prepare, verify, tag, publish, and check a QUALITY.md release.
tags: [release, changelog, versioning]
timestamp: 2026-06-19T00:00:00Z
---

# Cut a Release

Use this guide when publishing a tagged release of the `qualitymd` CLI and any
same-repo release notes for the `/quality` skill or `QUALITY.md` specification.

This guide summarizes the release-facing version checks. For the mechanical
tag-and-publish summary, see [Contributing](../../CONTRIBUTING.md).

## Preconditions

Before cutting a release:

- The main branch contains the intended release state.
- CI is passing.
- The working tree is clean.
- `NPM_TOKEN` is configured for npm publishing.
- `HOMEBREW_TAP_GITHUB_TOKEN` is configured for the Homebrew tap.
- You can run the local dry-run tasks.

Check the current state:

```sh
git status --short
git fetch --tags
git tag --list --sort=-version:refname | head
```

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
- `QUALITY.md` specification

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
- skill modes;
- bundled skill resources;
- required CLI command surface;
- supported CLI SemVer range.

Update the `SPECIFICATION.md` specification version only when the meaning of
conforming `QUALITY.md` documents or conforming evaluation/report behavior
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

The target release-note shape is a root `CHANGELOG.md` with an `Unreleased`
section and one section per tagged release. Until that file exists, write the
same curated notes directly into the GitHub Release body.

Suggested release-note groups:

```md
## v0.3.0 - YYYY-MM-DD

### CLI

### /quality Skill

### Specification

### Documentation

### Compatibility / Migration
```

Each release should include a compatibility block when relevant:

```md
Compatibility:

- CLI: v0.3.0
- QUALITY.md specification: 0.1
- /quality skill: compatible with qualitymd >=0.3.0 <0.4.0
```

## Run Local Verification

Run the normal gates:

```sh
mise run fmt
mise run test
mise run vet
```

For docs-only release preparation, run:

```sh
mise run fmt-md-check
```

Run release dry runs:

```sh
mise run snapshot
mise run npm-build
```

Inspect generated artifacts enough to catch packaging mistakes before pushing the
tag.

## Commit Release Preparation

Commit the release-prep changes before tagging.

Common release-prep changes include:

- `CHANGELOG.md`, once adopted;
- `SPECIFICATION.md` version line, when applicable;
- skill compatibility metadata or docs, when applicable;
- install or release documentation, when applicable.

Use a direct message such as:

```sh
git commit -m "Prepare v0.3.0 release"
```

## Tag and Publish

Create and push the tag:

```sh
git tag v0.3.0
git push origin v0.3.0
```

The release workflow runs from `.github/workflows/release.yml`. It publishes:

- GitHub release archives through Goreleaser;
- Homebrew cask updates through Goreleaser;
- npm / npx packages through `scripts/build-npm.mjs`.

Do not move a published tag to repair a release. If a published release is wrong,
fix forward with a new patch release unless no artifacts were published and the
tag is still only local or clearly failed before publication.

## Verify Published Artifacts

After the workflow finishes, verify:

- The GitHub Release exists for the tag.
- Release notes match the curated release-note entry.
- Checksums are present.
- The Homebrew tap updated.
- npm packages published for the supported platforms.
- `npx quality.md --version` reports the expected version.
- A downloaded archive binary reports the expected version.
- `qualitymd spec` emits the expected bundled specification.

Use fresh installs where practical rather than already-built local binaries.

## Handle a Failed Release

If the workflow fails before publishing anything, fix the problem and rerun the
workflow or replace the tag only if it has not been consumed externally.

If the workflow partially publishes artifacts:

- Do not rewrite the tag.
- Record what published and what failed.
- Fix the release automation or package issue.
- Publish a new patch tag.
- Add a release-note entry if users may have consumed the partial release.

If npm published but GitHub or Homebrew failed, treat the version as consumed.
npm versions are immutable in normal practice, so the repair should be a new
version.

## After Release

After verifying the release:

- Ensure `CHANGELOG.md` has a fresh `Unreleased` section, once adopted.
- Check install docs still point to valid commands.
- If the `/quality` skill compatibility range changed, check the skill guide and
  versioning reference are still aligned.
- Close or archive any completed Change Cases whose work shipped in the release.

## Open Process Support Items

The guide is enough to start a manual curated-release process. These follow-up
items would make the process more reliable once the first release or two has
exercised it:

- Add root `CHANGELOG.md` with `Unreleased` and per-tag sections, then make it
  the canonical source for curated release notes.
- Decide whether GitHub Release bodies should be populated from the
  `CHANGELOG.md` tag entry instead of Goreleaser's generated GitHub changelog.
- Add a lightweight release-prep task that checks for a target changelog entry,
  a fresh `Unreleased` section, version-surface consistency, and successful dry
  runs.
- Add skill package metadata for the `/quality` skill version and supported
  `qualitymd` CLI range when the installer supports it.
- Backfill old changelog entries only as far as useful; terse entries for older
  tags are acceptable.
- Keep the published-tag repair rule visible in contributor docs: once artifacts
  are published, fix forward with a new version.
