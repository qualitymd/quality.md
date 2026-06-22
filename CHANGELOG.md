# Changelog

User-facing release notes for `qualitymd`, the `/quality` skill, and the
QUALITY.md specification.

## Unreleased

### CLI

### /quality Skill

### Packaging

### Specification

## v0.9.0 - 2026-06-22

### CLI

- Added strict Area name, Factor name, and Rating Level ID validation to
  `qualitymd lint`, added canonical `area:`, `factor:`, and `rating:` model
  references, and updated generated Area Breakdown tables to separate Area
  titles from stable Area references.
- Added fixed-type unqualified model references and changed human Area
  Breakdown paths to render without redundant `area:` prefixes while keeping
  `report.json` and evaluation records structured.
- Separated human path display from model references so generated Markdown
  reports render the root Area path as `/` while `area:root`, `root`, and
  structured report JSON identifiers remain unchanged.
- Moved the default evaluation data directory from `quality/evaluations/` to
  `.quality/evaluations/`, added shared QUALITY.md workspace resolution with an
  optional root `config` pointer to the workspace config file, and kept
  `--evaluation-dir` as the highest-precedence override.
- Updated `qualitymd lint` to accept a valid root `config` tooling pointer while
  keeping other unknown root and nested keys as errors.
- Made evaluation record writes self-describing and safer to script:
  `evaluation assessment add`, `evaluation analysis set`, and
  `evaluation recommendation add` now document their JSON payload fields in
  `--help`, support `-n/--dry-run` validation without writing records, and report
  aggregated JSON-key validation errors.

### /quality Skill

- Aligned scoped evaluation guidance with qualified model references while
  keeping unqualified references only at fixed-type human input edges.
- Clarified that generated report path display values, including `/` for the
  root Area, are human labels and must not be written into evaluation records.
- Bumped the bundled skill release metadata to `0.9.0`; the supported
  `qualitymd` CLI compatibility range is now `>=0.9.0 <0.10.0`.
- Moved the default quality log from `quality/log/` to `.quality/log/` under the
  workspace quality data directory and aligned setup, wizard, evaluate, and
  recommendation follow-up guidance with the new workspace paths.
- Updated evaluation record-writing guidance and quick-reference payload examples
  to discover schemas through CLI help and validate new payloads with dry-run.
- Tightened `/quality evaluate` planning guidance so `design.md` and the initial
  `plan.md` are authored before assessment begins, with later scope, coverage,
  rigor, or evidence-strategy changes recorded as plan amendments.

### Packaging

- Added a package check that fails when a relative link inside the shipped
  `skills/quality` bundle does not resolve.

### Specification

- Defined Area names, Area IDs, Factor names, Factor IDs, Rating Level IDs, and
  canonical model references; clarified that durable machine artifacts preserve
  structured `areaPath` and `factorPath` arrays.
- Distinguished display values from model references, including `/` as the
  human display value for the root Area path.
- Bumped the QUALITY.md specification version to `0.3 (Draft)`.

### Documentation

- No documentation-only changes.

### Compatibility / Migration

Compatibility:

- CLI: `v0.9.0`
- QUALITY.md specification: `0.3 (Draft)`
- /quality skill: `0.9.0`, requires `qualitymd >=0.9.0 <0.10.0`

Existing QUALITY.md files may need Area names, Factor names, or Rating Level IDs
renamed if they used spaces, dots, slashes, colons, or leading/trailing
separators in structural keys. New evaluation runs default to
`.quality/evaluations/`; use `--evaluation-dir` or a workspace config pointer
when continuing to write somewhere else. Existing evaluation records remain
structured, while generated human Markdown reports now show `/` for the root
Area path.

## v0.8.0 - 2026-06-22

### CLI

- Added `qualitymd schema`, which emits the companion structural JSON Schema for
  QUALITY.md frontmatter. The schema is also published as `quality.schema.json`
  for editor and integration tooling; semantic validation remains owned by
  `qualitymd lint`.
- Added a compact Area-by-Factor report breakdown to generated evaluation report
  surfaces, making scoped strengths and gaps easier to scan without changing the
  underlying evaluation record contract.
- Updated the starter `QUALITY.md` scaffold comments to better explain Area
  `source` versus Requirement `assessment`.

### /quality Skill

- Removed `/quality improve` as a public mode. Evaluation still emits
  recommendations, and recommendation follow-up is now a non-mode workflow for
  applying a confirmed option or handing the recommendation to an issue tracker.
- Added convention-first quality-log guidance for meaningful model changes under
  `quality/log/`, including setup-created inaugural entries and confirmed
  recommendation/model-authoring changes.
- Added setup guidance for modeling the active `QUALITY.md` itself as a
  `quality-md` Area when a project wants to evaluate its quality model artifact.
- Bumped the bundled skill release metadata to `0.8.0`; the supported
  `qualitymd` CLI compatibility range is now `>=0.8.0 <0.9.0`.

### Specification

- Added a non-normative note pointing to the companion JSON Schema and clarifying
  that the schema is structural-only and subordinate to the specification and
  `qualitymd lint`.
- No specification version change; the QUALITY.md specification remains
  `0.2 (Draft)`.

### Documentation

- Expanded the README's quality-loop onboarding with setup, evaluation, review,
  model-refinement, cadence, and issue-tracker handoff guidance.
- Split and aligned durable `/quality` skill specs and evaluation-record specs
  into parent/component/artifact contracts so mode behavior, reporting, quality
  logs, recommendation follow-up, and generated artifacts are easier to review.
- Updated contributor guidance for functional-spec granularity, durable-spec
  change accounting, and Change Case status conventions.

### Packaging

- No packaging changes.

### Compatibility / Migration

Compatibility:

- CLI: `v0.8.0`
- QUALITY.md specification: `0.2 (Draft)`
- /quality skill: `0.8.0`, requires `qualitymd >=0.8.0 <0.9.0`

For `/quality` users, replace `/quality improve` prompts with evaluation
follow-up prompts such as "apply recommendation 002" or "handoff recommendation
002." CLI users do not need to migrate existing `QUALITY.md` files or evaluation
records for this release.

## v0.7.2 - 2026-06-20

### CLI

- `qualitymd init --json` now emits only the JSON error receipt when scaffolding
  fails, instead of also rendering the styled error a second time. A `--json`
  consumer reading stderr sees a single failure report; the exit code is
  unchanged.
- Reworked internal evaluation package naming to drop the redundant `Evaluation`
  prefix (for example `evaluation.Run` and `evaluation.RunStatus`), clarified
  usage-error exit classification, and added contract documentation comments to
  follow the project Go style guide. Machine-readable JSON output and exit codes
  are unchanged.

### /quality Skill

- Bumped the bundled skill release metadata to `0.7.2`; the supported
  `qualitymd` CLI compatibility range remains `>=0.7.0 <0.8.0`.

### Specification

- No specification changes.

### Documentation

- Refined the Go style and package-design contributor guides: reorganized the
  headings, clarified which conventions the deterministic check gate already
  enforces, and tied the package-design guidance to the in-repo `receipt.Action`
  example.

### Packaging

- No packaging changes.

### Compatibility / Migration

Compatibility:

- CLI: `v0.7.2`
- QUALITY.md specification: `0.2 (Draft)`
- /quality skill: `0.7.2`, requires `qualitymd >=0.7.0 <0.8.0`

No migration is required from `v0.7.1`.

## v0.7.1 - 2026-06-21

### CLI

- Added documentation comments for exported evaluation, report, and status
  contract types so the implementation follows the project Go style guide
  without changing command behavior or machine-readable output.

### /quality Skill

- Bumped the bundled skill release metadata to `0.7.1`; the supported
  `qualitymd` CLI compatibility range remains `>=0.7.0 <0.8.0`.

### Specification

- No specification changes.

### Documentation

- Added a Go style guide for judgment-based conventions that the deterministic
  check gate does not enforce, linked it from contributor and agent guidance, and
  added it to the dogfood quality model as a CLI maintainability requirement.

### Packaging

- No packaging changes.

### Compatibility / Migration

Compatibility:

- CLI: `v0.7.1`
- QUALITY.md specification: `0.2 (Draft)`
- /quality skill: `0.7.1`, requires `qualitymd >=0.7.0 <0.8.0`

No migration is required from `v0.7.0`.

## v0.7.0 - 2026-06-21

### Breaking Changes

- Renamed the formal model-node concept from Target to Area throughout the draft
  format: `areas:` replaces `targets:`, evaluation records use `areaPath`
  instead of `targetPath`, reports render Area/root area labels, and current
  readers treat legacy Target-shaped records as incompatible historical records.
- Renamed `qualitymd evaluation create --subject` to `--model` and changed new
  evaluation run folders from `NNNN-subject[-<narrowing>]-quality-eval` to
  `NNNN[-<narrowing>]-quality-eval`.

### Documentation

- Updated the `/quality` skill, scaffold, README, specs, dogfood `QUALITY.md`,
  and maintained Sparrow example bundle to use Area terminology while preserving
  the default `target` / `Target` rating level.

### Compatibility / Migration

Compatibility:

- CLI: `v0.7.0`
- QUALITY.md specification: `0.2 (Draft)`
- /quality skill: `0.7.0`, requires `qualitymd >=0.7.0 <0.8.0`

This is a breaking draft-format change. Existing `targets:` models and
evaluation records that use `targetPath` must be updated to `areas:` and
`areaPath`; current readers treat old records as incompatible historical data
rather than translating them.

## v0.6.0 - 2026-06-21

### CLI

- Added `debug-log.md` to newly created evaluation runs as a process-only
  diagnostic artifact for evaluation orchestration, retries, recovery, coverage
  decisions, and tooling friction.

### /quality Skill

- Updated evaluation guidance so agents hand-author `debug-log.md` for notable
  evaluation-process events while keeping subject-quality findings, rating
  rationale, and raw project-command output in formal assessment, analysis, and
  recommendation records.

### Specification

- Added the `debug-log.md` runtime artifact contract to evaluation records and
  clarified that reports, ratings, findings, recommendations, and next actions
  remain derived from formal evaluation records.

### Documentation

- Updated the `/quality` usage guide and reference example evaluation run to
  show the process-only debug log alongside design, plan, records, and reports.

### Packaging

### Compatibility / Migration

Compatibility:

- CLI: `v0.6.0`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: `0.6.0`, requires `qualitymd >=0.6.0 <0.7.0`

Existing evaluation runs remain readable and reportable. New evaluation runs
include `debug-log.md`; older runs do not need to be migrated.

## v0.5.1 - 2026-06-21

### CLI

- Made evaluation report assembly use explicit typed states for rating results,
  local target ratings, run gaps, lifecycle, next steps, rigor, evaluation
  level, missing metadata, and path identities, so invalid or incomplete report
  inputs surface as clear reportability gaps.
- Made evaluation-history inspection tolerant of older, malformed, partial, or
  hand-edited run records while keeping current record writers strict and report
  generation explicit.

### /quality Skill

- Clarified authoring, setup, getting-started, and quick-check guidance so the
  Markdown body is treated as evaluable judgment context, including
  agent-accessible support and scoped unknowns/open questions.
- Aligned the skill contract with tolerant historical-run inspection: malformed
  or incompatible historical records are evaluation-history status, not subject
  quality evidence.

### Specification

- Clarified Markdown body semantics: the body supports building, interpreting,
  using, and evaluating the model; unknowns and open questions are section-scoped
  author-declared context distinct from `not assessed` evaluation results.
- Clarified evaluation-record and report contracts for typed report states and
  tolerant evaluation-history readers.

### Documentation

- Replaced the standalone Known gaps body pattern with per-section unknowns,
  open questions, and human/agent review state guidance across the authoring
  guide, scaffold, examples, and dogfood model.
- Documented the Markdown body as concise, self-explanatory judgment context
  whose completeness, recency, grounding, and agent-accessibility can itself be
  evaluated.
- Archived completed Change Cases 0042 through 0045.

### Packaging

- Updated pinned installer examples to `v0.5.1`.

### Compatibility / Migration

Compatibility:

- CLI: `v0.5.1`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: `0.5.1`, requires `qualitymd >=0.5.0 <0.6.0`

No migration is required from `v0.5.0`. Existing QUALITY.md files remain valid;
authors can move standalone Known gaps content into the relevant section's
unknowns or open questions when they next revise the body.

## v0.5.0 - 2026-06-20

### CLI

- Renamed `qualitymd upgrade` to `qualitymd update` with no `upgrade` alias.
  `qualitymd update` applies by default, while `qualitymd update --check`
  reports without mutating.
- Added managed standalone self-apply, release-readiness gating, post-apply
  visible-version verification, release-note references, and a cached ambient
  update notice controlled by `QUALITYMD_NO_UPDATE_CHECK=1`.

### /quality Skill

- Renamed `/quality upgrade` to `/quality update` and retargeted the maintenance
  workflow to `qualitymd update` / `qualitymd update --check`.

### Specification

### Documentation

### Packaging

### Compatibility / Migration

Compatibility:

- CLI: `v0.5.0`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: `0.5.0`, requires `qualitymd >=0.5.0 <0.6.0`

The `qualitymd upgrade` command and `/quality upgrade` mode were renamed to
`update`; update scripts and prompts to call `qualitymd update --check`,
`qualitymd update`, and `/quality update`.

## v0.4.1 - 2026-06-20

### CLI

- Fixed evaluation run status and record-write receipts so machine-readable
  paths stay repository-relative instead of leaking local absolute paths.
- Fixed empty evaluation runs to report the documented `missing-root-analysis`
  gap, and mapped `qualitymd lint -` to the usage-error exit category.
- Kept `report-summary.md` Top Issues focused on issue-like findings instead of
  labeling ordinary positive findings as issues.

### /quality Skill

- Bumped the bundled skill release metadata to `0.4.1`; the supported
  `qualitymd` CLI compatibility range remains `>=0.4.0 <0.5.0`.

### Specification

- Fixed Appendix B's minimal example to include required rating-level and Factor
  `title` fields.
- Added non-normative invalid counter-examples for missing rating-level titles,
  direct Target Requirements without Factors, and list-valued Assessments.

### Documentation

- Added README first-result output and evaluation report output examples so a
  newcomer can see what successful `qualitymd` commands produce.
- Archived pre-current-contract dogfood evaluation runs out of the live
  evaluation scan path and recorded a fresh deep evaluation at `target`.

### Packaging

- No packaging changes.

### Compatibility / Migration

Compatibility:

- CLI: `v0.4.1`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: `0.4.1`, requires `qualitymd >=0.4.0 <0.5.0`

No migration is required from `v0.4.0`.

## v0.4.0 - 2026-06-19

### CLI

- Reshaped the evaluation surface: `evaluation create`, `list`, `status`,
  record resources (`assessment add`, `analysis set`, `recommendation add`),
  and `report build` / `report gate` replace the older mixed verbs.
- Folded planned coverage into `plan.md` frontmatter, added batched record
  writes and explicit `--latest` run resolution, and removed the create receipt's
  altitude field.
- Added a concise no-argument welcome screen and reorganized root help into
  common workflow tasks and management commands.
- Refined terminal styling so help output stays mostly neutral while status and
  command output keep semantic color.

### /quality Skill

- Updated evaluation workflows and CLI references for the new evaluation command
  surface and `plan.md` coverage metadata.
- Expanded the authoring guide with stronger guidance for traceability, rating
  boundaries, target/source alignment, roll-up decisions, stakeholder concerns,
  measurable criteria, and requirement-set review.

### Specification

- Replaced the durable evaluation CLI sub-specs with the new noun/verb surface
  and updated the evaluation record contract to remove `planned-coverage.json`.

### Documentation

### Packaging

### Compatibility / Migration

Compatibility:

- CLI: `v0.4.0`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: `0.4.0`, requires `qualitymd >=0.4.0 <0.5.0`

Evaluation workflows and integrations that use the previous `qualitymd
evaluation` command names must migrate to the `v0.4.0` noun/verb command
surface.

## v0.3.1 - 2026-06-19

### CLI

- Fixed `qualitymd upgrade --check` so update detection uses SemVer precedence
  instead of treating any version difference as an upgrade.
- `qualitymd evaluation build-report` human Markdown now uses Model, Target,
  Factor, and Rating Level titles from the run's model snapshot as display
  labels while preserving stable ids in `report.json` and gates.

### /quality Skill

- Added a durable interaction contract for run frames, decision briefs, stop and
  reroute behavior, history-aware operation, improvement delta reports, and
  status-first output.
- Expanded the bundled CLI quick reference around structured version, status,
  evaluation, and upgrade workflows.

### Specification

- Clarified that evaluation record payloads and `report.json` preserve stable
  identifiers, while human Markdown reports resolve display labels from the
  recorded model snapshot.

### Documentation

- Documented installer PATH guidance, `--non-interactive` behavior, and the
  intentional Homebrew cask distribution path.
- Clarified the release runbook's Homebrew cask rationale.

### Packaging

- Hardened the Unix installer checksum verification fallback across `sha256sum`,
  `shasum`, and `openssl`, with explicit warnings when verification cannot run.
- Hardened the Windows installer for Windows PowerShell TLS 1.2 and per-user
  PATH setup.

### Compatibility / Migration

Compatibility:

- CLI: `v0.3.1`
- QUALITY.md specification: `0.1 (Draft)`
- /quality skill: `0.3.1`, requires `qualitymd >=0.3.0 <0.4.0`

No migration is required from `v0.3.0`.

## v0.3.0 - 2026-06-19

### CLI

- Added `qualitymd version` for structured local version metadata and
  `qualitymd upgrade` for explicit upgrade checks and guarded owner-channel
  apply behavior.
- `qualitymd evaluation build-report` now writes `report-summary.md` beside
  `report.md` and `report.json`.

### /quality Skill

- Added `/quality upgrade` to plan and orchestrate paired `/quality` skill and
  `qualitymd` CLI upgrades while delegating mutation to the skill installer and
  CLI owner channel.
- Updated CLI prerequisite checks and evaluation reporting guidance for
  `version --json`, `upgrade --check`, and `report-summary.md`.

### Specification

- Required `title` on Models, Targets, Factors, and Rating Levels; `lint` now
  reports missing required titles as errors.

### Documentation

- Documented `/quality upgrade` as the existing-install maintenance flow for the
  skill/CLI pair.
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

- No normative QUALITY.md format or evaluation semantics change.
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

- No normative QUALITY.md format or evaluation semantics change.
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

- Replaced the meta-model workflow with the bundled QUALITY.md authoring
  guide.
- Modularized setup, wizard, evaluate, and improve procedures into
  mode-specific skill files.
- Reframed the wizard as a read-only quality wayfinder.

### Specification

- Updated evaluation records, evaluation CLI commands, report behavior,
  recommendation and assessment superseding, and factor-reference terminology.
- Kept the QUALITY.md specification at `0.1 (Draft)`.

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
