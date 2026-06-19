# Changelog

User-facing release notes for `qualitymd`, the `/quality` skill, and the
`QUALITY.md` specification.

## Unreleased

### CLI

- Added `qualitymd version` for structured local version metadata and
  `qualitymd upgrade` for explicit upgrade checks and guarded owner-channel
  apply behavior.
- `qualitymd evaluation build-report` now writes `report-summary.md` beside
  `report.md` and `report.json`.

### /quality Skill

- Updated CLI prerequisite checks and evaluation reporting guidance for
  `version --json`, `upgrade --check`, and `report-summary.md`.

### Specification

- Required `title` on Models, Targets, Factors, and Rating Levels; `lint` now
  reports missing required titles as errors.

### Documentation

- Documented the GitHub-hosted managed installer entrypoints and explicit
  upgrade workflow.
- Clarified the release runbook after the `v0.2.2` release: push `main` and
  wait for hosted CI before tagging, treat pre-tag snapshot logs and trailing
  newline-only release-note diffs as non-blocking, and keep release-prep manual
  until repeated mistakes justify more mechanics.
- Replaced the release guide's open process checklist with explicit process
  support boundaries for manual release prep and future `/quality` skill package
  metadata.

### Packaging

### Compatibility / Migration

Compatibility:

- CLI: `v0.3.0`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: `0.3.0`, requires `qualitymd >=0.3.0 <0.4.0`

## v0.2.2 - 2026-06-19

### CLI

- No command, flag, exit-code, or machine-readable output behavior changes.

### /quality Skill

- No skill instruction or compatibility changes.

### Specification

- No normative `QUALITY.md` format or evaluation semantics change.
- Specification remains `0.1 (Draft)`.

### Documentation

- Added release-note extraction and release-check automation to the release
  guide.

### Packaging

- Wired GitHub Release bodies to the curated `CHANGELOG.md` release section.
- Aligned committed npm launcher optional dependency placeholders with the
  `npm-build` dry-run version so release checks leave the tree clean.
- Ignored generated `npm/platforms/` package directories produced by local npm
  dry runs.

### Compatibility / Migration

Compatibility:

- CLI: `v0.2.2`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: compatible with `qualitymd >=0.2.0 <0.3.0`

No migration is required from `v0.2.1`.

## v0.2.1 - 2026-06-19

### CLI

- Fixed lint issues so the release commit passes the full local and hosted CI
  gate.
- No command, flag, exit-code, or machine-readable output behavior changes.

### /quality Skill

- Documented versioning ownership and the skill/CLI compatibility boundary.
- No skill compatibility expansion or breaking skill-facing CLI change.

### Specification

- No normative `QUALITY.md` format or evaluation semantics change.
- Specification remains `0.1 (Draft)`.

### Documentation

- Added the release runbook and versioning reference.
- Updated install and contributor documentation to point at the release and
  compatibility policy.

### Packaging

- Generate the npm launcher README from the repository README during package
  build.
- Added an npm package check to catch README packaging drift.

### Compatibility / Migration

Compatibility:

- CLI: `v0.2.1`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: compatible with `qualitymd >=0.2.0 <0.3.0`

No migration is required from `v0.2.0`.

## v0.2.0 - 2026-06-19

### CLI

- Added the evaluation CLI surface for creating runs, adding records, checking
  status, and building reports.
- Expanded evaluation record lifecycle and reporting behavior, including
  planned coverage, superseding, duplicate-assessment detection, and report
  renderability checks.
- Added mandatory factor-reference linting for direct target-level requirements.

### /quality Skill

- Replaced the meta-model workflow with the bundled `QUALITY.md` authoring
  guide.
- Modularized setup, wizard, evaluate, and improve procedures into
  mode-specific skill files.
- Reframed the wizard as a read-only quality wayfinder.

### Specification

- Updated evaluation records, evaluation CLI commands, report behavior,
  recommendation and assessment superseding, and factor-reference terminology.
- Kept the `QUALITY.md` specification at `0.1 (Draft)`.

### Documentation

- Documented npm and Homebrew install commands.
- Archived completed change cases and recorded durable specification guidance.
- Tightened the authoring guide, skill runtime contract, and assessment
  reference terminology.

### Compatibility / Migration

Compatibility:

- CLI: `v0.2.0`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: compatible with `qualitymd >=0.2.0 <0.3.0`

`v0.2.0` is the compatibility floor for the modular `/quality` skill and the
evaluation CLI workflow.

## v0.1.2 - 2026-06-17

### Packaging

- Wired the Homebrew tap publishing token into the Goreleaser cask release path.

### Compatibility / Migration

No user-facing CLI, skill, or specification behavior changed from `v0.1.1`.

## v0.1.1 - 2026-06-17

### CLI

- Added the initial `qualitymd init`, `qualitymd lint`, and `qualitymd spec`
  command surfaces.
- Implemented structural schema validation, deterministic CLI output, and
  agent-accessibility checks.

### /quality Skill

- Added the initial installable `/quality` skill and example evaluation bundle.

### Specification

- Split rating levels into description and criterion fields.
- Added target display fields and aligned lint behavior with the format model.

### Packaging

- Added npm launcher version alignment and the project license.

### Compatibility / Migration

Compatibility:

- CLI: `v0.1.1`
- QUALITY.md specification: `0.1 (Draft)`

This is the first tagged public release represented in the changelog.
