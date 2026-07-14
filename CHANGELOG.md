# Changelog

User-facing release notes for `qualitymd`, the `/quality` skill, and the
QUALITY.md specification.

## Unreleased

### CLI

- Harness-backed evaluation runs are no longer pinned to concurrency 1: the
  runner now keeps a rolling window of dependency-ready work requests
  outstanding — up to the run's resolved `evaluation.concurrency` — topping
  the window up as results are accepted, so the invoking harness can judge
  independent requests in parallel or fan them out to subagents.
  Omitted-concurrency harness runs now use the shared automatic default
  (`max(2, NumCPU*2)`) instead of forced 1; receipts, dry-run previews, and
  progress output state the resolved concurrency and window width.
- `awaiting_evaluator` receipts carry `evaluatorRequests` — the outstanding
  bounded work-request set — replacing the singular `evaluatorRequest`, and
  `--evaluator-result` accepts one result envelope or a JSON array covering
  any subset of the outstanding requests. Not-yet-submitted requests stay
  outstanding at no retry cost; a schema-invalid or failed member re-emits
  for its retry attempt (with its `lastFailure` named) without touching other
  members' accepted results.
- `evaluation.json` bumps to schema version 7: the harness checkpoint state
  is now the plural `pendingEvaluatorCalls`. Clean break: an in-flight
  awaiting run created before the upgrade cannot resume across it — re-run
  it; completed runs are unaffected. `qualitymd evaluation status` lists
  every outstanding request for an awaiting run.

### /quality skill

- The evaluate workflow's checkpoint loop services the outstanding request
  set: it may delegate independent requests to subagents, submits results as
  they become ready (one envelope or several per call), and names the window
  width on the first windowed receipt.

## v0.30.0 - 2026-07-11

### Specification

- The format now commits to non-filesystem source selectors: `source` stays a
  single string, and its kind is detected from the selector itself — glob
  metacharacters make a glob, an existing filesystem entry is a path, and any
  other selector describes the evaluated material in prose (for example
  "open tickets in the support queue"), resolved by evaluating tools. The
  frontmatter shape and `quality.schema.json` are unchanged. The specification
  version moves to `0.11 (Draft)`: what a `source` selector means in a
  conforming document has changed.

### CLI

- Evaluation runs resolve every source selector through a per-kind resolver.
  Path and glob selectors keep the deterministic workspace walk; a prose
  selector is dispatched to the invoking harness as a `resolveSource` work
  request on the existing checkpoint transport, and the returned material is
  validated, capped, hashed, and captured into the run artifact before any
  dependent judgment. Selector kinds are detected and pinned at run creation
  and honored on resume.
- A selector kind the selected evaluator cannot resolve now fails the run at
  plan time with the new `selector_unsupported` category — naming the
  selector, its detected kind, and the remedy — distinct from
  `source_unavailable`, which keeps meaning the named material is missing.
- `evaluation.json` bumps to schema version 6 and gains a per-area `sources`
  provenance record: each area's selector, detected kind, serving resolver
  (and harness runtime when harness-resolved), bundle hash, capture time, and
  per-file hashes — with captured content kept for harness-gathered evidence
  so resume re-judges against the exact evidence of record.
- `qualitymd evaluation run --dry-run --json` and run receipts now include
  the per-area source dispatch plan (selector, detected kind, resolver).
- Evaluation runs now resolve area sources as the format specification
  defines: a source-less root area evaluates the QUALITY.md file's directory,
  a source-less child area inherits the nearest ancestor's source, and glob
  selectors (for example `docs/**/*.md`) are expanded. Previously these were
  silently packaged as empty evidence.
- An area source that resolves to no readable files now fails the run loudly
  with `source_unavailable`, naming the unresolved selector, instead of
  judging the area against an empty source bundle.
- Source packaging skips symlinked directories and files, sockets, and
  devices instead of crashing on them (`is a directory`), so repositories with
  committed symlinks evaluate and resume cleanly.
- `qualitymd lint` no longer reports spec-permitted extension frontmatter as
  `invalid-frontmatter`: a key that names no model property is now the
  warning-severity `unknown-key` advisory, so conforming documents carrying
  extension properties lint valid and load. Misplaced root keys on nested
  areas remain errors.
- `quality.schema.json` now accepts any non-empty scalar (string, number, or
  boolean) for content scalars such as `assessment`, rating-level
  `criterion`, and `ratings` override values, matching `lint`; its `$comment`
  no longer claims a rating-level ordering check that no tool performs.

### /quality skill

- The harness checkpoint loop now serves `resolveSource` resolution requests
  alongside judgment requests: the skill gathers exactly the material the
  request's selector describes and returns it verbatim, or reports
  `source_unavailable` when it does not exist.
- `/quality` skill metadata now declares version `0.30.0` and requires the
  `qualitymd` CLI `0.30.x` line.

### Compatibility / migration

- The QUALITY.md specification moves to `0.11 (Draft)`. A `source` selector is
  now defined for non-filesystem material: a selector that is neither a glob
  nor an existing filesystem path is prose describing the evaluated material,
  where before only filesystem selectors were defined. Existing path and glob
  selectors are unchanged.
- Runner `evaluation.json` artifacts now use schema version `6` and add the
  per-area `sources` provenance record; in-flight version 5 runner runs should
  be started fresh.
- A selector kind the selected evaluator cannot resolve now fails the run at
  plan time with the new `selector_unsupported` category. `source_unavailable`
  keeps its narrower meaning: the named material is missing.

Compatibility:

- CLI: `v0.30.0`
- QUALITY.md specification: `0.11 (Draft)`
- /quality skill: `0.30.0`, requires `qualitymd >=0.30.0 <0.31.0`

## v0.29.0 - 2026-07-10

### CLI

- Evaluation runs can now execute dependency-ready evaluator-backed steps
  concurrently. Configure the runner with `evaluation.concurrency` in
  `.quality/config.yaml`; omit it to use the automatic default
  `max(2, runtime.NumCPU()*2)`, set `1` for sequential execution, or set any
  positive integer to choose the active evaluator-call limit.
- `qualitymd evaluation run --dry-run --json`, run receipts, and
  `evaluation.json` now report resolved `concurrency` directly. The previous
  public `executionStrategy` and `strategyFallbacks` fields are removed.
- Harness-backed evaluation remains sequential (`concurrency: 1`) while the
  runner supports one pending harness checkpoint at a time.

### /quality skill

- Runtime guidance now documents `evaluation.concurrency` as the runner
  configuration knob and describes dry-run previews as reporting resolved
  concurrency.
- `/quality` skill metadata now declares version `0.29.0` and requires the
  `qualitymd` CLI `0.29.x` line.

### Compatibility / migration

- Replace `evaluation.executionStrategy` in `.quality/config.yaml` with
  `evaluation.concurrency`. Use `concurrency: 1` for the old sequential
  behavior.
- Machine consumers of evaluation run receipts and `evaluation.json` should read
  `concurrency`; `executionStrategy` and `strategyFallbacks` are no longer
  emitted.
- Runner `evaluation.json` artifacts now use schema version `5`; in-flight
  version 4 runner runs should be started fresh.
- No QUALITY.md specification version change; the specification remains
  `0.10 (Draft)`.

Compatibility:

- CLI: `v0.29.0`
- QUALITY.md specification: `0.10 (Draft)`
- /quality skill: `0.29.0`, requires `qualitymd >=0.29.0 <0.30.0`

## v0.28.0 - 2026-07-10

### CLI

- New deterministic evaluation runner: `qualitymd evaluation run` now executes
  the complete evaluation — run creation, the dependency-ordered work graph,
  evaluator invocation for bounded judgment, result validation, atomic
  persistence, and Markdown report generation — with one authoritative
  `evaluation.json` run artifact per run, run-local logs, `--resume`, and a
  `--dry-run --json` preview. Evaluators are pluggable: `codex` and `claude`
  CLI subprocesses, direct `openai` and `anthropic` API profiles, and
  configured profiles in `.quality/config.yaml`.
- New `harness` evaluator: an invoking agent (such as Claude Code or Codex)
  supplies evaluation judgment directly through checkpointed work requests —
  no nested agent process and no provider API key. Runs checkpoint with the
  stable status `awaiting_evaluator` (exit `0`), results are submitted with
  `qualitymd evaluation run --resume <run> --evaluator-result -`, and the
  first accepted result binds the run to one harness runtime.
- `qualitymd evaluation status`, `qualitymd evaluation list` (new
  `--state awaiting` filter), and `qualitymd status` now report the runner
  lifecycle, identify runs awaiting harness judgment, and name the exact
  continuation command.
- `auto` evaluator discovery is readiness-aware: a CLI candidate is selected
  only after verifying its executable, authentication state (where a
  documented non-interactive probe exists), and required structured-output
  capabilities, with per-candidate evidence in the dry-run preview. CLI
  evaluator adapters use native JSON Schema output and no-persistence flags
  when the installed CLI supports them.
- Evaluation runs spend far fewer input tokens with unchanged results:
  evaluator prompts render a cacheable stable prefix (task, schema, packaged
  source, shared area context) before the per-work-unit delta, the Anthropic
  adapter applies provider prompt caching, each requirement is assessed and
  rated in one evaluator call instead of two, an area's source is packaged
  once per run, and `logs/evaluator-calls.jsonl` records cached input tokens.

### /quality skill

- `/quality evaluate` is now an agent-mediated wrapper around
  `qualitymd evaluation run`: the skill keeps intent parsing, scope
  resolution, preflight validation, evaluator-selection explanation, and
  result summaries, while the runner owns evaluation execution and artifacts.
- Inside a capable agent harness, evaluation now defaults to
  `--evaluator harness`: the session's own agent judges each runner-emitted
  bounded work request and submits the typed result, preserving explicit user
  and configured evaluator choices and never silently switching providers.
- `/quality` skill metadata now declares version `0.28.0` and requires the
  `qualitymd` CLI `0.28.x` line.

### Documentation

- New automation guides: run recurring `/quality evaluate` in Claude Code
  routines and Codex scheduled tasks, including skill/CLI availability,
  artifact persistence, permissions, unattended behavior, and fallback
  credential boundaries.
- Install guide documents evaluator defaults: harness-backed evaluation for
  agent surfaces needs no provider API key; direct CLI `auto` discovery and
  API-profile credentials are the fallback paths.

- The LedgerLite report gallery is now an authoring exemplar: the example
  QUALITY.md demonstrates the authoring guide family (full body sections with
  unknowns and review provenance, a model-wide agent harnessability factor with
  seven sub-factors, a normative service-contract area, a veto requirement, a
  measured rating override, and sensor-grounded assessments), and its generated
  evaluation includes all four finding types, a not-assessed result, and a
  synthetic quality changelog showing how the model matured.
- The LedgerLite exemplar now uses a deeper 39-requirement model with a
  `codebase` area, expanded security/auditability/capacity/compatibility
  coverage, quality-named `quality-md` factors, reused sensor-catalog
  assessments, and a recommendation that matures an inferential model-body
  review into a repeatable drift detector.

### Compatibility / migration

- Requirement judgment work units changed shape (`assessRateRequirement`
  replaces the `assessRequirement`/`rateRequirement` pair) and evaluator input
  hashes widened, so in-flight runs started with an earlier `qualitymd` do not
  reuse prior results on resume; start a fresh run. Persisted payload kinds,
  identities, schemas, and generated reports are unchanged.
- No QUALITY.md specification version change; the specification remains
  `0.10 (Draft)`.
- `/quality` skill version `0.28.0` requires the `qualitymd` CLI `0.28.x`
  line: harness checkpoints need `--evaluator-result`, the
  `awaiting_evaluator` lifecycle, and the readiness-aware dry-run preview.

Compatibility:

- CLI: `v0.28.0`
- QUALITY.md specification: `0.10 (Draft)`
- /quality skill: `0.28.0`, requires `qualitymd >=0.28.0 <0.29.0`

## v0.27.1 - 2026-07-02

### Documentation

- Model vocabulary now uses normal English capitalization across the
  specification, docs, skill, and CLI output; capitalization no longer marks
  terms of art.

### /quality skill

- `/quality` skill metadata now declares version `0.27.1` for the patched
  `qualitymd` `0.27.x` line.

### Packaging

- The install smoke workflow now verifies every CLI install method documented on
  the Mintlify quickstart: hosted installers, npm, pnpm, Homebrew, and Go source
  installs.

### Compatibility / migration

- No QUALITY.md specification version change; the specification remains
  `0.10 (Draft)`.
- `/quality` skill version `0.27.1` requires the `qualitymd` CLI `0.27.x` line.

Compatibility:

- CLI: `v0.27.1`
- QUALITY.md specification: `0.10 (Draft)`
- /quality skill: `0.27.1`, requires `qualitymd >=0.27.0 <0.28.0`

## v0.27.0 - 2026-07-02

### Specification

- QUALITY.md specification version is now `0.10 (Draft)`. The format
  specification now defines document and Model semantics only: document
  conformance, source resolution, Requirement scope, Factor connection, Rating
  Scale semantics, extensions, and body semantics.
- Evaluation methods, aggregation approaches, report formats, Advice, and
  tool-specific process contracts are no longer normative format conformance
  classes. The `/quality` skill evaluation method is now described as a
  non-normative illustrative appendix in the format specification.

### /quality skill

- `/quality` evaluation specs now own the skill's evaluation process contract
  directly while binding the skill's model reading to the format spec's Model
  semantics.
- `/quality` runtime prerequisites now route agents to the specification's
  Model semantics section for source resolution, Requirement scope, Factor
  connection, and Rating Scale meaning.

### Documentation

- Mintlify documentation now describes the specification as a document and Model
  semantics contract rather than an evaluation/reporting conformance contract.
- The Mintlify home page now adds FAQ context on quality-model lineage and why
  agent-mediated evaluation makes tailored quality models more practical.

### Packaging

- Hosted installer URLs now use Mintlify redirects to the GitHub-hosted
  `install/` scripts so `https://getquality.md/install.sh`, `/install.ps1`, and
  `/install.cmd` resolve in production.
- The install smoke workflow now verifies the hosted installer URLs in addition
  to the repository checkout scripts.

### Compatibility / migration

- QUALITY.md specification version is now `0.10 (Draft)`. The specification no
  longer defines evaluator or report-renderer conformance classes; evaluation
  behavior belongs to the evaluating method, such as the `/quality` skill.
- `/quality` skill version `0.27.0` requires the `qualitymd` CLI `0.27.x` line.

Compatibility:

- CLI: `v0.27.0`
- QUALITY.md specification: `0.10 (Draft)`
- /quality skill: `0.27.0`, requires `qualitymd >=0.27.0 <0.28.0`

## v0.26.3 - 2026-07-01

### CLI

- Managed standalone `qualitymd update` now fetches hosted installer
  entrypoints from `https://getquality.md/install.sh` and
  `https://getquality.md/install.ps1` instead of GitHub raw source paths.

### Packaging

- Installer entrypoints now live in the Mintlify docs site root so the hosted
  docs domain serves `/install.sh`, `/install.ps1`, and `/install.cmd` directly.
- The install smoke workflow now tests the relocated Mintlify-root installer
  files on macOS, Linux, and Windows.

### Documentation

- The Mintlify quickstart now shows npm and pnpm CLI install options, and the
  release guide now verifies the docs-hosted installer entrypoints.

### Compatibility / migration

- No QUALITY.md specification version change; the specification remains
  `0.9 (Draft)`.
- `/quality` skill version `0.26.3` requires the `qualitymd` CLI `0.26.x` line.

Compatibility:

- CLI: `v0.26.3`
- QUALITY.md specification: `0.9 (Draft)`
- /quality skill: `0.26.3`, requires `qualitymd >=0.26.0 <0.27.0`

## v0.26.2 - 2026-06-30

### CLI

- Generated Evaluation report titles, section headings, and Contents labels now
  use sentence case for fixed report labels such as `Quality evaluation`, `Key
details`, `Model evaluation`, and `Primary source data`, while preserving
  model-provided titles, table headers, and structured metadata.
- Generated Evaluation run reports no longer render the standalone Finding
  Summary table near Key details. The Key details total and the full
  `findings.md` link's type/severity count summary remain.
- Generated Evaluation run reports now place both full Findings and full
  Recommendations report links above their capped preview tables, add semantic
  markers to full-list count summaries, label recommendation count groups by
  `impact`, and render Strength Findings as `💪 Strength`.
- Evaluation Finding `severity` is now concern-only: required for `gap` and
  `risk` Findings and rejected for `strength` and `note` Findings. Generated
  reports render `—` for non-concern severity cells and summarize full Findings
  links with type and concern-severity counts.
- Generated Evaluation reports now link to `report.md`, `findings.md`,
  `recommendations.md`, and the new workspace-root `glossary.md` through one
  H1-adjacent `Evaluation links:` blockquote, and no longer repeat local
  `Legend` blocks in each report artifact.
- Generated `Model evaluation` and `Area / Factor breakdown` tables now label
  their first column `▦ Area / □ Factor`, making Area and Factor row markers
  visible in the table header.
- Evaluation runs now persist `data/evaluation-manifest.json` with
  `kind: "EvaluationManifest"`, a durable `evaluationId`, and nested local
  `run` metadata. Generated `report.md` frontmatter now exposes
  `evaluationId`, `created`, `model`, and `run`, and no longer duplicates scope
  as frontmatter.
- Generated Evaluation reports now use canonical labels such as
  `Quality rating`, `Finding type`, `Finding severity`, and
  `Recommendation impact`, and the internal fixed enum catalogs now carry
  type-level and value-level descriptions for glossary/help surfaces.
- Generated recommendation list tables now show a single recommendation number
  column (`#`) instead of duplicate `Rank` and `#` columns.
- Generated Evaluation run reports now emphasize the full Findings and
  Recommendations report links under the capped Top findings and Top
  recommendations tables, and show the complete ranked count for each linked
  report.
- Recommendation results now use opaque `qrec_...` IDs for JSON data paths and
  structured ranking/coverage refs. Generated reports reserve recommendation
  numbers for ranked order, so "recommendation #1" means the first ranked
  recommendation.
- Generated Evaluation reports now use standard `Contents` sections and
  `Evaluation links:` navigation instead of compact `Jump to:` navigation or
  bottom `Legend` sections, while keeping marker/icon values paired with text
  labels in report tables.
- Evaluation data validation, schemas, and generated Markdown reports now share
  typed enum catalogs for fixed values such as statuses, confidence, finding
  type/severity/basis status, recommendation impact, ranking tier, coverage
  disposition, and report kind. Invalid enum inputs remain rejected at
  `evaluation data set` / `verify`, while reports render known values with
  consistent labels and markers.
- `qualitymd status` now reports workspace status instead of project state,
  including a `schemaVersion: 2` JSON `workspace` block with relative workspace
  paths for config, `.quality/`, evaluations, changelog, and workflow logs.

### /quality skill

- `/quality` runtime guidance, guide, workflow, and resource headings now follow
  sentence case while preserving QUALITY.md, formal model vocabulary, and other
  proper names.
- `/quality evaluate` guidance now writes `severity` only for `gap` and `risk`
  Findings, omitting it from `strength` and `note` Findings.
- `/quality evaluate` guidance now writes recommendation payloads before ranking,
  reads assigned `qrec_...` IDs, and uses ranked order as the user-facing
  recommendation number.
- `/quality` runtime guidance now treats `qualitymd status` as workspace status
  while preserving project language for modeled value and setup intent.

### Documentation

- Active README/docs/guides/specs, generated Mintlify specification output, and
  checked-in report-gallery artifacts now use sentence-case headings while
  preserving proper names, formal model terms, model-provided titles, and
  historical archive records.
- The release guide now treats local `release-check` as the pre-tag dry-run gate
  and relies on GitHub Actions preflight to validate npm and Homebrew publishing
  credentials from repository secrets.

### Compatibility / migration

- No QUALITY.md specification version change; the specification remains
  `0.9 (Draft)`.
- `/quality` skill version `0.26.2` requires the `qualitymd` CLI `0.26.x` line.

Compatibility:

- CLI: `v0.26.2`
- QUALITY.md specification: `0.9 (Draft)`
- /quality skill: `0.26.2`, requires `qualitymd >=0.26.0 <0.27.0`

## v0.26.1 - 2026-06-29

### CLI

- Generated Evaluation run reports now remove redundant post-title navigation,
  recommended-next-action prose, Scope, Coverage, and Report Details sections;
  move the Area/Factor table up as `Model Evaluation`; and use
  `Primary Source Data` sections with report-local primary payloads.
- Generated Evaluation reports now render recommendation impact with subtle
  shape markers (`⬥⬥ Very high`, `⬥ High`, `● Medium`, `○ Low`), restore impact
  to `report.md` Top Recommendations, and use quieter Area / Factor Breakdown
  row markers.
- Generated Evaluation run reports now open with `Summary`, `Key Details`, and
  `Contents` sections, move run metadata into frontmatter and lower report
  details, and omit the visible `Limits & Incomplete Inputs` section for now.
- Generated Evaluation run reports now title `report.md` as
  `Quality Evaluation - <Area>` and include factor filters in parentheses, for
  example `Quality Evaluation - Public API (Reliability, Correctness)`.

### /quality skill

- `/quality` skill metadata now declares version `0.26.1` for the patched
  `qualitymd` `0.26.x` line.

### Compatibility / Migration

- No QUALITY.md specification version change; the specification remains
  `0.9 (Draft)`.
- `/quality` skill version `0.26.1` requires the `qualitymd` CLI `0.26.x` line.

Compatibility:

- CLI: `v0.26.1`
- QUALITY.md specification: `0.9 (Draft)`
- /quality skill: `0.26.1`, requires `qualitymd >=0.26.0 <0.27.0`

## v0.26.0 - 2026-06-29

### CLI

- Evaluation runs now persist a globally-unique `RunManifest.id` alongside the
  local run number. Recommendations use per-run `number` values and reports
  render typed references such as
  `evaluation:<run-id>/recommendation/<number>`; ranked findings no longer carry
  synthetic artifact IDs.

### /quality skill

- `/quality evaluate` guidance now authors Advice against CLI-assigned
  recommendation numbers and references Requirement Findings by `findingRef`
  selectors only.

### Specification

- QUALITY.md specification version is now `0.9 (Draft)`, reflecting the updated
  Evaluation run identity, recommendation numbering, and finding reference
  semantics.

### Compatibility / Migration

- Breaking Evaluation data/report artifact contract change: pre-`v0.26.0`
  Evaluation runs that use `QEVAL`, `QREC`, or `QFIND` artifact IDs are
  historical data and are not migrated.
- `/quality` skill version `0.26.0` requires the `qualitymd` CLI `0.26.x` line.

Compatibility:

- CLI: `v0.26.0`
- QUALITY.md specification: `0.9 (Draft)`
- /quality skill: `0.26.0`, requires `qualitymd >=0.26.0 <0.27.0`

## v0.25.6 - 2026-06-29

### CLI

- `qualitymd init` now adds a concise `QUALITY.md` pointer to local agent
  instruction files by default, with `--no-agent-instructions` to opt out.
- `qualitymd init` scaffolds now point direct CLI users to `/quality setup` when
  the file was initialized outside the setup workflow.
- Evaluation Advice now uses CLI-assigned, user-citable artifact IDs:
  recommendations persist as `QREC-<run>-<seq>`, ranked findings persist as
  `QFIND-<run>-<seq>`, and generated reports render those IDs in finding and
  recommendation tables/detail pages.

### /quality skill

- `/quality setup` now scaffolds missing model files with
  `qualitymd init --no-agent-instructions` and explicitly distinguishes
  scaffold-only, partially authored, and mature existing `QUALITY.md` files.

### Compatibility / Migration

- `/quality` skill version `0.25.6` continues to require the `qualitymd` CLI
  `0.25.x` line.
- No QUALITY.md specification version change; the specification remains
  `0.8 (Draft)`.

Compatibility:

- CLI: `v0.25.6`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.25.6`, requires `qualitymd >=0.25.0 <0.26.0`

## v0.25.5 - 2026-06-27

### CLI

- Evaluation Advice now uses the simplified `RecommendationResult` shape:
  `description`, `background`, `expectedValue`, and `doneCriterion` replace the
  older recommendation prose fields. Generated recommendation reports render
  from persisted Advice data only, with `report.md` Top Recommendations showing
  `Rank`, `Recommendation`, `Area / Factors`, and `Reason`.
- Evaluation report build now generates a full ranked `findings.md` index,
  links `report.md` Top Findings to exact Requirement finding details, and shows
  Advice rank context inside Requirement finding detail sections.
- Generated Evaluation report Markdown authoring is now centralized around
  shared table, link, code-span, and empty-cell helpers, including table-cell
  escaping for separators and multiline text.

### Compatibility / Migration

- `/quality` skill version `0.25.5` continues to require the `qualitymd` CLI
  `0.25.x` line.
- No QUALITY.md specification version change; the specification remains
  `0.8 (Draft)`.

Compatibility:

- CLI: `v0.25.5`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.25.5`, requires `qualitymd >=0.25.0 <0.26.0`

## v0.25.4 - 2026-06-27

### Packaging

- The release workflow's final publish job now passes an explicit GitHub
  repository to `gh release edit`, fixing the `v0.25.3` publish step when the job
  has no checked-out git repository.

### Compatibility / Migration

- `/quality` skill version `0.25.4` continues to require the `qualitymd` CLI
  `0.25.x` line.
- No QUALITY.md specification version change; the specification remains
  `0.8 (Draft)`.

Compatibility:

- CLI: `v0.25.4`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.25.4`, requires `qualitymd >=0.25.0 <0.26.0`

## v0.25.3 - 2026-06-27

### Packaging

- Release verification now reads the Homebrew cask from the public raw tap URL
  instead of the cross-repo GitHub contents API, fixing the `v0.25.2`
  draft-stage verifier failure under the workflow token.

### Compatibility / Migration

- `/quality` skill version `0.25.3` continues to require the `qualitymd` CLI
  `0.25.x` line.
- No QUALITY.md specification version change; the specification remains
  `0.8 (Draft)`.

Compatibility:

- CLI: `v0.25.3`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.25.3`, requires `qualitymd >=0.25.0 <0.26.0`

## v0.25.2 - 2026-06-27

### Packaging

- Release Homebrew and verification scripts now discover draft GitHub releases
  by listing releases when GitHub's tag-specific release endpoint returns 404,
  fixing the `v0.25.1` draft-stage Homebrew failure.
- The Homebrew updater now clones the public tap without using the write token
  and keeps the authenticated push command from exposing token-bearing command
  lines on local failures.

### Compatibility / Migration

- `/quality` skill version `0.25.2` continues to require the `qualitymd` CLI
  `0.25.x` line.
- No QUALITY.md specification version change; the specification remains
  `0.8 (Draft)`.

Compatibility:

- CLI: `v0.25.2`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.25.2`, requires `qualitymd >=0.25.0 <0.26.0`

## v0.25.1 - 2026-06-27

### /quality Skill

- `/quality evaluate` closeouts now point first to the completed run's
  `report.md` as the decision-ready evaluation result, and describe
  `recommendations.md` as the action-planning report instead of presenting a
  generic report list.

### Packaging

- Release publishing now runs credential and target preflight before tagging,
  creates draft GitHub releases, publishes Homebrew and npm through separate
  idempotent jobs, verifies all release channels before publishing the GitHub
  Release, and provides a repair helper for safe-to-rerun release channels.

### Compatibility / Migration

- `/quality` skill version `0.25.1` continues to require the `qualitymd` CLI
  `0.25.x` line.
- No QUALITY.md specification version change; the specification remains
  `0.8 (Draft)`.

Compatibility:

- CLI: `v0.25.1`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.25.1`, requires `qualitymd >=0.25.0 <0.26.0`

## v0.25.0 - 2026-06-27

### CLI

- Evaluation runs now persist CLI-owned `RunManifest` scope data. New runs use
  `qualitymd evaluation create [model] [--area <area-id>] [--factor <factor-id>...]`;
  the previous `--narrowing` surface is removed.
- Evaluation reports now render `report.md` as the scoped Area report for the
  run and no longer select a headline subject from agent-authored payload
  ordering.
- Evaluation data Finding Core objects now use `basis` instead of `cause` for
  the finding-local explanation or support posture. Generated reports now render
  `Basis` and `Basis Evidence` labels instead of `Cause` and `Cause Evidence`.
- Evaluation now requires Advice payloads before a run is reportable:
  `FindingRankingResult`, one or more `RecommendationResult` payloads, and
  `RecommendationRankingResult` with finding coverage accounting.
- Evaluation report build now renders Top Findings and Top Recommendations in
  `report.md`, plus `recommendations.md` and per-recommendation detail pages.

### /quality Skill

- `/quality evaluate` now resolves scoped runs to canonical `area:` and
  `factor:` IDs before creating the run and relies on CLI-owned planned scope.
- `/quality evaluate` now writes Requirement Findings with `basis`,
  `basis.status`, and `basis.rationale` instead of `cause` fields.
- `/quality evaluate` now produces required, domain-agnostic recommendations as
  part of evaluation Advice. Recommendations may be concrete work or a
  recommended review of whether to raise, clarify, or confirm the next quality
  bar.

### Compatibility / Migration

- `/quality` skill version `0.25.0` requires the `qualitymd` CLI `0.25.x` line.
- Existing Evaluation run data that uses Finding Core `cause` is incompatible
  with current evaluation data validation and report generation. There is no
  migration or compatibility mode.
- Existing Evaluation runs without Advice payloads are no longer reportable
  under the current evaluation report builder.

Compatibility:

- CLI: `v0.25.0`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.25.0`, requires `qualitymd >=0.25.0 <0.26.0`

## v0.24.1 - 2026-06-27

### CLI

- Evaluation reports now label immediate descendant Areas as `Child Areas` and
  immediate descendant Factors as `Sub-Factors`; the old `Sub-Areas` and
  `Child Factors` report labels are no longer generated.

### /quality Skill

- `/quality` runtime outputs now adopt labeled, scan-friendly templates for
  direct model-authoring checkpoints, setup review, evaluation/review/improve/update
  closeouts, recommendation results, model-review signals, and next-workflow
  choices.
- The quality model-change history now lives under `.quality/changelog/` with
  timestamped entry filenames such as
  `2026-06-27T184233Z-adjust-rating-scale.md`. `.quality/logs/` remains the
  flat home for workflow/process logs such as setup and evaluate feedback logs.
- Direct `QUALITY.md` edit checkpoints now state the inferred purpose for the
  change and ask the user to react to the consequential scope or risk assumption
  instead of defaulting to a generic adjustment prompt.
- `/quality review` and `/quality improve` are now public workflow stubs.
  `review` is focus-routed and read-only by default; `improve` confirms focus and
  mutation surface before delegating to existing model-authoring or compatible
  recommendation follow-up routes.

### Compatibility / Migration

- `/quality` skill version `0.24.1` requires the `qualitymd` CLI `0.24.x` line.
- No QUALITY.md specification version change; the specification remains
  `0.8 (Draft)`.

Compatibility:

- CLI: `v0.24.1`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.24.1`, requires `qualitymd >=0.24.0 <0.25.0`

## v0.24.0 - 2026-06-27

### CLI

- `qualitymd lint` now treats Requirement `factors` references as Area-local:
  references to ancestor-only Factors are `unknown-factor` errors. It also warns
  when the same Factor name appears more than once inside one Area's Factor tree.
- Evaluation data schema version is now `3`. `qualitymd evaluation data set` and
  `data verify` reject Area-analysis findings, require every rated Requirement
  result to have a paired Requirement Assessment with at least one Finding, and
  require rated Requirement, Factor, and Area results to carry non-empty
  `ratingDrivers` whose `inputRefs` resolve in the effective run data set.

### /quality Skill

- `/quality` now requires `qualitymd >=0.24.0 <0.25.0`.
- `/quality evaluate` now treats Requirement Findings as the only Evaluation
  finding layer. Rated Requirements must be backed by Requirement Findings and
  rating drivers, while Factor and Area roll-ups explain judgment through
  `ratingDrivers`, rationale, confidence, limits, and incomplete-input fields
  without synthesizing new findings.
- Direct `QUALITY.md` edit requests now use a lightweight authoring checkpoint:
  the skill infers intent, asks follow-up only when it materially affects the
  edit, states "Here's what I'm planning to do" in simple prose with a visible
  value prop, and accepts `looks good` as confirmation when the mutation is
  clear. The checkpoint explicitly waits for that feedback before editing, and
  the skill acknowledges long direct-authoring reads before inspecting model
  context.
- Authoring guidance now keeps Requirement `factors` references local to the
  Requirement's declaring Area and models root cross-cutting judgment with root
  Requirements instead of descendant opt-in Factor references.

### Specification

- QUALITY.md specification version is now `0.8 (Draft)`.
- Evaluation Findings are now Requirement-level only. Factor and Area reports no
  longer render `Findings` sections, and roll-up analysis must cite lower-level
  evidence through rating drivers rather than inventing another finding layer.
- Requirement Factor references are now Area-local. Same-named Factors on
  different Areas remain distinct, while duplicate Factor names within one
  Area's recursive Factor tree are discouraged because scalar references become
  ambiguous to readers.

### Compatibility / Migration

- Existing Evaluation schema version 2 run data is incompatible historical data
  for current `evaluation data set`, `data verify`, and report generation. There
  is no migration or compatibility mode.
- `/quality` skill version `0.24.0` requires the `qualitymd` CLI `0.24.x` line.

Compatibility:

- CLI: `v0.24.0`
- QUALITY.md specification: `0.8 (Draft)`
- /quality skill: `0.24.0`, requires `qualitymd >=0.24.0 <0.25.0`

## v0.23.0 - 2026-06-27

### CLI

- `qualitymd evaluation report build` now writes `report.md` as the run-level
  Evaluation report. The root Area detail report is `root-area.md` when present,
  and scoped Area/Factor runs can build reports without root Area analysis when
  their recorded headline scope is complete.

### /quality Skill

- `/quality` now requires `qualitymd >=0.23.0 <0.24.0`.
- `/quality evaluate` closeouts now treat `report.md` as the run report and name
  the headline subject report separately when it differs.

### Specification

- QUALITY.md specification version is now `0.7 (Draft)`.
- `EvaluationOutputResult` now includes `runReportRef`, `headlineResultRef`,
  and `headlineReportRef`; `rootAreaAnalysisRef` is present only when the root
  Area was evaluated.

### Compatibility / Migration

- Consumers that opened the root Area detail at `report.md` should follow
  `EvaluationOutputResult.headlineReportRef`, `rootAreaReportMd`, or
  `root-area.md` when the root Area was evaluated. `report.md` is now the run
  entrypoint.
- `/quality` skill version `0.23.0` requires the `qualitymd` CLI `0.23.x` line.

Compatibility:

- CLI: `v0.23.0`
- QUALITY.md specification: `0.7 (Draft)`
- /quality skill: `0.23.0`, requires `qualitymd >=0.23.0 <0.24.0`

## v0.22.0 - 2026-06-27

### CLI

- `qualitymd evaluation data example <kind>` now emits fuller representative
  payloads across supported Evaluation data kinds, including populated limits,
  stop conditions, canonical Model references, and report references.
- Workspace config and artifacts now resolve relative to the selected
  `QUALITY.md`, with the Git root retained as the containment boundary.
  Evaluation list/latest/data/status/report commands can use `--model` to anchor
  nested workspaces and model-relative run paths.
- Evaluation Requirement and Area Findings now use a shared structured Finding
  Core: `statement`, `condition`, `criteria`, `cause`, `effect`, and `evidence`.
  Legacy finding `description`, `summary`, and top-level `rationale` fields are
  rejected.
- Requirement Finding remediation leads are now written as `candidateActions`
  with local candidate action IDs. The legacy finding-local `actions` field is
  rejected, and candidate actions remain out of Evaluation v0 reports.

### /quality Skill

- `/quality` now requires `qualitymd >=0.22.0 <0.23.0`.
- `/quality evaluate` now writes structured Finding Core payloads, uses
  `candidateActions` for Requirement Finding-local remediation leads, and keeps
  Area Findings free of candidate actions.
- `/quality` workspace guidance now treats config, Evaluation history, quality
  logs, and workflow feedback logs as model-relative to the selected
  `QUALITY.md`.

### Specification

- QUALITY.md specification version is now `0.6 (Draft)`.
- The Evaluation contract now defines structured Finding Core semantics and
  names `candidateActions` as finding-local raw material for future advice.
- The CLI specifications now define model-relative workspace artifact
  resolution and model-anchored Evaluation run discovery.

### Compatibility / Migration

- Existing Evaluation data that uses legacy finding `description`, `summary`,
  top-level `rationale`, or finding-local `actions` must be rewritten before it
  will pass current `evaluation data set` or `evaluation data verify`
  validation.
- Existing repo-root `.quality/` workspace artifacts for nested `QUALITY.md`
  files are no longer discovered as the default workspace location; move them
  under the selected model's directory or pass the appropriate model-relative
  override.
- `/quality` skill version `0.22.0` requires the `qualitymd` CLI `0.22.x` line.

Compatibility:

- CLI: `v0.22.0`
- QUALITY.md specification: `0.6 (Draft)`
- /quality skill: `0.22.0`, requires `qualitymd >=0.22.0 <0.23.0`

## v0.21.0 - 2026-06-26

### CLI

- Evaluation finding severity no longer accepts `info`. Requirement Finding and
  Area Finding severity is now limited to `critical`, `high`, `medium`, and
  `low`, with `type: note` carrying informational observations.
- `qualitymd evaluation data schema [<kind>]`, `evaluation data set`, and
  `evaluation data verify` all reflect and enforce the reduced severity set.

### /quality Skill

- `/quality` now requires `qualitymd >=0.21.0 <0.22.0`.
- `/quality evaluate` now treats informational Area Findings as `type: note`
  rather than `severity: info`.

### Compatibility / Migration

- Existing Evaluation data that uses `severity: "info"` must be rewritten before
  it will pass current `evaluation data set` or `evaluation data verify`
  validation.
- `/quality` skill version `0.21.0` requires the `qualitymd` CLI `0.21.x` line.

Compatibility:

- CLI: `v0.21.0`
- QUALITY.md specification: `0.5 (Draft)`
- /quality skill: `0.21.0`, requires `qualitymd >=0.21.0 <0.22.0`

## v0.20.0 - 2026-06-26

### CLI

- `AreaAnalysisResult` now accepts validated `findings` for Area-level
  Evaluation findings. Area Findings use closed `type`, `severity`,
  `confidence`, and Factor relationship enums, require non-empty `inputRefs`,
  reject advice/ranking fields, and validate that related Factors belong to the
  same Area.
- Evaluation Area reports now render Area Findings, and Factor reports render
  the matching Area Findings related to that Factor.

### /quality Skill

- `/quality` now requires `qualitymd >=0.20.0 <0.21.0`.
- `/quality evaluate` now records material Area Findings during Area analysis so
  Area and Factor reports can show synthesized findings before any future
  recommendation phase.

### Compatibility / Migration

- Existing `AreaAnalysisResult` payloads without `findings` remain valid and
  render as no Area Findings. New Area Finding objects are closed: fields such as
  `impact`, `priority`, `effort`, `benefit`, `ROI`, `actions`, and
  `recommendations` are rejected.
- `/quality` skill version `0.20.0` requires the `qualitymd` CLI `0.20.x` line.

Compatibility:

- CLI: `v0.20.0`
- QUALITY.md specification: `0.5 (Draft)`
- /quality skill: `0.20.0`, requires `qualitymd >=0.20.0 <0.21.0`

## v0.19.0 - 2026-06-26

### CLI

- `qualitymd evaluation data schema <kind>` now emits a self-contained schema for
  the requested Evaluation payload kind. Required fields, enum value sets, and
  model-reference patterns are visible in the emitted document without following
  a top-level `$ref` into `$defs`; the no-argument full-surface schema remains
  available for all kinds.
- `qualitymd evaluation data schema` now uses the same terminal rendering as
  `qualitymd schema`: redirects stay verbatim JSON, while terminal output may be
  highlighted and paged.

### /quality Skill

- `/quality` now requires `qualitymd >=0.19.0 <0.20.0`.
- The skill now treats `qualitymd evaluation data schema <kind>` as the source
  for required fields and allowed enum values, with
  `qualitymd evaluation data example <kind>` used as one concrete valid instance.

### Documentation

- Change-case and functional-spec authoring guidance now requires
  per-requirement `Durable spec:` annotations in change-case specs and clarifies
  the split between runtime skill files under `skills/` and durable skill specs
  under `specs/skills/`.

### Compatibility / Migration

- Consumers that read `qualitymd evaluation data schema <kind>` should read the
  kind schema directly from the document root. The no-argument
  `qualitymd evaluation data schema` output remains the full-surface schema.
- `/quality` skill version `0.19.0` requires the `qualitymd` CLI `0.19.x` line.

Compatibility:

- CLI: `v0.19.0`
- QUALITY.md specification: `0.5 (Draft)`
- /quality skill: `0.19.0`, requires `qualitymd >=0.19.0 <0.20.0`

## v0.18.0 - 2026-06-26

### /quality Skill

- `/quality` now requires `qualitymd >=0.18.0 <0.19.0`.
- `/quality evaluate` now has one best-quality workflow: evaluation rigor
  (`quick`/`standard`/`deep`, `--rigor`, and `/quality evaluate deep`) is removed,
  every in-scope Requirement must be assessed or explicitly marked not assessed,
  independent collection and QC fan out to subagents when available, and an
  always-on QC phase verifies roll-up-binding findings while sweeping for missed
  gaps before ratings are finalized.

- The setup opening now shows the run frame first, distinguishes `QUALITY.md`
  changes from local workflow feedback-log writes, and keeps setup feedback-log
  timing consistent. Read-only orientation now has a status-first output shape,
  and recommendation follow-up now opens with a frame before inspection or
  mutation.

### Documentation

- The agent-mediated UX guide's evaluate run-frame example no longer includes a
  `Rigor:` field.

### Compatibility / Migration

- `/quality evaluate` no longer accepts evaluation-rigor invocations such as
  `/quality evaluate --rigor quick`, `/quality evaluate --rigor deep`, or
  `/quality evaluate deep`. Narrow the evaluation scope by Area or Factor instead.
- `/quality` skill version `0.18.0` requires the `qualitymd` CLI `0.18.x` line.

Compatibility:

- CLI: `v0.18.0`
- QUALITY.md specification: `0.5 (Draft)`
- /quality skill: `0.18.0`, requires `qualitymd >=0.18.0 <0.19.0`

## v0.17.0 - 2026-06-26

### CLI

- Added the read-only `qualitymd model` command group:
  `qualitymd model tree`, `qualitymd model list`, and `qualitymd model get`
  project a Model's Area, Factor, Requirement, and Rating Level identities,
  labels, containment, and canonical reference IDs in human and JSON forms.
- `qualitymd evaluation data set` now reads a non-empty JSON array of payloads
  from stdin instead of one bare JSON object. The command validates the whole
  batch before writing, rejects duplicate derived paths, reports indexed
  diagnostics for invalid elements, writes all-or-nothing, and emits a batch
  receipt with `count`, `writes[]`, `dryRun`, and `nextActions`.
- Evaluation finding `actions` is now a typed candidate-action field: an array of
  objects with a required `description` and optional `rationale`, validated when
  persisting a `RequirementAssessmentResult` (previously an untyped, unvalidated
  array). Candidate actions are non-binding, finding-local remediation leads kept
  out of the v0 report; the generated finding detail no longer renders an
  `Actions` row.
- Evaluation routine references now validate their `kind` member against the
  supported Evaluation payload-kind vocabulary, and report references validate
  their `kind` member against the report-kind vocabulary, so misspelled or
  invented reference kinds fail at `evaluation data set` time.

### /quality Skill

- `/quality` now requires `qualitymd >=0.17.0 <0.18.0`.
- The evaluate workflow now queries in-scope canonical model IDs from the run's
  `model-snapshot.md` with `qualitymd model list --json` instead of hand-deriving
  Area, Factor, Requirement, and Rating Level references.
- The evaluate workflow now batches routine Evaluation payload persistence:
  payloads are accumulated into a JSON array, validated once with
  `evaluation data set --dry-run`, and written once with `evaluation data set`.
- The CLI workflow resource is now `cli-workflow-conventions.md`. It keeps the
  skill-specific workflow conventions the CLI cannot describe and routes command,
  flag, and payload discovery to CLI help, `--json`, schema, and example outputs
  instead of an embedded command listing.
- The evaluate workflow now records non-binding candidate actions on `gap` and
  `risk` findings — short remediation leads captured where the evidence is
  richest — as raw material for a future Advise phase, and omits them on
  `strength` findings. They are not recommendations and are not surfaced in the
  Evaluation v0 report or closeout.

- Agent-collaborated models that still use the old `harnessability` factor no
  longer count that factor as current Agent Harnessability coverage. The skill now
  treats it as stale legacy naming to fix by renaming to `agent-harnessability`
  and adding any missing current sub-factors.

- Decision gates (setup write/update, update apply-plan, recommendation apply and
  issue creation) now lead with the question and present choices as a visually
  separated block instead of a flat stack of labels, so the action to take is
  clear at a glance and stays readable when formatting is stripped. Setup
  discovery questions show why a choice matters before the answer line, and the
  evaluate workflow reports progress before run creation and the per-requirement
  phase.

- Questions, closed choices, and confirmation gates now render through a native
  interaction affordance — an option picker, a confirm/approve gate, the harness's
  own authorization prompt — when the agent runtime offers one that fits, falling
  back to the established numbered-option and `y`/`n` text forms when it does not.
  The skill keeps the explanation in the message rather than in small widget
  labels, stays independent of any specific assistant, and no longer stacks its own
  confirmation on top of a mutation the harness already prompts to authorize.

### Specification

- _Assess Requirements_ now notes that a Finding MAY carry non-binding candidate
  actions — finding-local remediation leads — distinct from the recommendations
  the Advice phase produces.
- The CLI spec now defines `qualitymd model tree|list|get` as the deterministic
  structure/identity projection surface for Models.
- The Evaluation data CLI spec now defines the `data set` stdin envelope as a
  non-empty JSON array, with whole-batch validation, all-or-nothing writes,
  indexed diagnostics, duplicate-path rejection, dry-run parity, and the new
  batch receipt shape.
- Evaluation reference JSON conventions now constrain routine-reference and
  report-reference `kind` fields to their closed vocabularies.

### Documentation

- The README now foregrounds the early-alpha status and points quality-workflow
  users to the `/quality` skill and `qualitymd update`.

### Compatibility / Migration

- `qualitymd evaluation data set` no longer accepts a bare payload object on
  stdin. Wrap one payload in a one-element array, e.g. `[payload]`, and submit
  multi-payload runs as one array.
- `qualitymd evaluation data set --json` now returns a batch receipt
  (`count`, `writes[]`) instead of a single `path`/`kind` receipt. Consumers must
  iterate `writes[]`.
- Existing persisted Evaluation data files are not migrated. New writes through
  `evaluation data set` must satisfy the closed reference-kind vocabularies.
- `/quality` skill version `0.17.0` requires the `qualitymd` CLI `0.17.x` line.

Compatibility:

- CLI: `v0.17.0`
- QUALITY.md specification: `0.5 (Draft)`
- /quality skill: `0.17.0`, requires `qualitymd >=0.17.0 <0.18.0`

## v0.16.0 - 2026-06-26

### CLI

- Evaluation routine and report JSON now encode `areaId`, `factorId`,
  `requirementId`, and `ratingLevelId` as canonical qualified model-reference
  strings such as `area:root`, `factor:root::reliability`,
  `requirement:root::has-tests`, and `rating:target`. Repeated identity fields
  use arrays of the same strings.
- Evaluation data `schemaVersion` is now `2`; old structured identity payloads
  are rejected rather than migrated.
- Generated Evaluation Markdown reports now open with kind-prefixed H1 titles,
  place navigation trails below the title, omit redundant `Path:` / `Name:`
  header lines, render empty scalar table cells as `—`, and include a static
  legend defining the marker.
- `qualitymd lint` now rejects `root` as a reserved Area name.

### /quality Skill

- `/quality` now requires `qualitymd >=0.16.0 <0.17.0` and instructs agents to
  write Evaluation payload identities as qualified model-reference strings.

### Specification

- QUALITY.md specification version is now `0.5 (Draft)`.
- Area names may no longer be `root`, preserving `area:root` as the unambiguous
  root Area reference.
- Evaluation data specs now define persisted model identities as canonical
  qualified reference strings.

### Documentation

- Report-tree specs now document title-first, kind-prefixed report headers, the
  empty-cell em dash, and the static report legend.

### Compatibility / Migration

- Existing completed Evaluation runs are not migrated. Recreate or rewrite run
  data with schema version `2` before rebuilding reports with `qualitymd
evaluation report build`.
- Consumers of Evaluation JSON must read identity fields as qualified strings,
  not arrays or objects.

Compatibility:

- CLI: `v0.16.0`
- QUALITY.md specification: `0.5 (Draft)`
- /quality skill: `0.16.0`, requires `qualitymd >=0.16.0 <0.17.0`

## v0.15.0 - 2026-06-26

### CLI

- `qualitymd evaluation data set` now rejects unknown or misspelled fields,
  invalid field types, invalid enum values, and unresolved Area, Factor,
  Requirement, and Rating Level IDs against the run's `model-snapshot.md`.
  `--dry-run` performs the same validation without writing.
- Added `qualitymd evaluation data schema [<kind>]` for discovering the
  structured Evaluation payload contract and `qualitymd evaluation data verify
<run>` for checking persisted run data.
- `qualitymd evaluation data example <kind>` now includes representative nested
  entries for findings, rating drivers, unknowns, inputs, and limits instead of
  leaving those arrays empty.
- Evaluation Requirement reports now show a plural `Factors:` context line under
  the `Area:` trail and no longer duplicate attached Factors in the header
  summary table.

### /quality Skill

- `/quality evaluate` now discovers Evaluation payload shapes through
  `qualitymd evaluation data schema <kind>`, uses examples as populated samples,
  and treats `data set --dry-run` as authored-payload validation.
- Public workflow run frames are now required as the first output before any tool
  call, with provisional `resolving…` values allowed when a field needs tool
  inspection.

### Specification

- The active Evaluation spec bundle is now named `specs/evaluation/` and the
  live docs/specs use plain "Evaluation" instead of the vestigial "Evaluation
  v2" name. The `schemaVersion` field and payload shapes are unchanged.
- Evaluation data specs now document schema discovery, run-data verification,
  strict payload validation, model binding, populated examples, and dry-run
  parity.

### Documentation

- Agent-mediated UX guidance now requires workflows to open with an immediate
  run frame that confirms intent and previews the path.

### Compatibility / Migration

- Payloads accepted by earlier releases may now fail validation if they contain
  misspelled fields or references to model nodes that are absent from the run
  snapshot. Use `qualitymd evaluation data verify <run>` to inspect existing run
  data before rebuilding reports.

Compatibility:

- CLI: `v0.15.0`
- QUALITY.md specification: `0.4 (Draft)`
- /quality skill: `0.15.0`, requires `qualitymd >=0.15.0 <0.16.0`

## v0.14.1 - 2026-06-26

### CLI

- `qualitymd evaluation create` now names new run folders
  `NNNN-full-eval` or `NNNN-<scope-path>-eval`, shortening the old
  `quality-eval` tag.

### /quality Skill

- `/quality evaluate` now passes `--narrowing` as the scope's full structural
  path for Area and Factor scoped evaluations.

Compatibility:

- CLI: `v0.14.1`
- QUALITY.md specification: `0.4 (Draft)`
- /quality skill: `0.14.1`, requires `qualitymd >=0.14.1 <0.15.0`

## v0.14.0 - 2026-06-26

### CLI

- Evaluation v2 keeps the root Area report at `report.md`, but descendant Area,
  Factor, and Requirement Markdown reports now use short subject-aware filenames
  such as `<area>-area.md`, `<factor>-factor.md`, and
  `<requirement>-requirement.md`. Generated Markdown links and
  `EvaluationOutputResult` report refs use the new paths.

### /quality Skill

- `/quality evaluate` now describes the generated Markdown report tree instead
  of implying the only human report artifact is `report.md`.

### Documentation

- Evaluation v2 report layout specs, the `/quality` reporting contract, and the
  v2 sketch now document the short descendant report filenames.

### Compatibility / Migration

- Existing completed evaluation runs are not migrated. Rebuild a report with
  `qualitymd evaluation report build <run>` to generate the new descendant
  report filenames for a run.
- Callers that hard-code descendant `report.md` paths should follow generated
  report refs from `data/evaluation-output-result.json` or use the documented
  structural path rules.

Compatibility:

- CLI: `v0.14.0`
- QUALITY.md specification: `0.4 (Draft)`
- /quality skill: `0.14.0`, requires `qualitymd >=0.14.0 <0.15.0`

## v0.13.1 - 2026-06-26

### CLI

- Evaluation v2 Markdown report headers are cleaner and more navigable: reports
  now use labeled `Area:` trails, Factor reports add Factor-only `Factor:`
  trails, redundant `Breadcrumb:` / parent header links are gone, and the top
  status block is a compact report-specific summary.
- Evaluation v2 Area and Factor report tables now link generated human report
  pages from the Factor, Area, or Requirement subject cell instead of repeating a
  generic `Details` column on every row. Structured `Data` links remain
  explicitly labeled and machine-readable report refs are unchanged.
- Evaluation v2 Markdown reports now render Rating Level titles from the run's
  model snapshot and render CLI-owned statuses, confidence levels, booleans,
  report kinds, limits, unknowns, and known finding classifications as
  human-readable titles while keeping raw values in routine JSON,
  `EvaluationOutputResult`, and build receipts.

### /quality Skill

- `/quality setup` now opens with a short explanation of what QUALITY.md gives
  teams and agents, then performs a read-only context scan and shows a
  project-specific setup preview before discovery questions or file writes.
- `/quality setup`, `evaluate`, `update`, and recommendation follow-up guidance
  now use clearer CTAs, shortest answer paths, numbered ambiguity and outcome
  choices, decision briefs before writes or external handoff, progress updates,
  closeout `Next` fields, and precise code spans for literal artifacts.
- `/quality setup` closed-choice prompts now use numbered options with the
  recommended answer first and `1` as the shortest confirmation. Setup risk
  discovery presents cost labels while keeping the existing risk-tolerance
  meaning internally.
- `/quality evaluate` feedback logs now use workflow-process `outcome` values
  such as `completed-reportable`, `stopped-model`, and `interrupted`, keeping
  `.quality/logs/<timestamp>-evaluate-feedback-log.md` clearly separate from
  report, rating, and recommendation semantics.

### Documentation

- The functional-spec and Change Case guides now include a set-level requirement
  check, assumptions/dependencies guidance, normative vs. informational reference
  classification, an explicit unambiguous-requirement bar item, and an optional
  EARS statement template.
- The durable Evaluation v2 report-tree spec and v2 sketch now match the current
  generated Markdown report navigation and table-link shape.

Compatibility:

- CLI: `v0.13.1`
- QUALITY.md specification: `0.4 (Draft)`
- /quality skill: `0.13.1`, requires `qualitymd >=0.13.0 <0.14.0`

## v0.13.0 - 2026-06-26

### CLI

- Evaluation runs are now a clean Evaluation v2 runtime: new runs seed only
  `model.md` and `data/`, old assessment/analysis/recommendation run shapes are
  rejected as unsupported, and report builds emit
  `data/evaluation-output-result.json` plus the v2 report tree without legacy
  `report-summary.md` or `report.json` artifacts.
- `qualitymd evaluation data set` now reads one JSON payload from stdin and no
  longer accepts `--file`. Use
  `qualitymd evaluation data set <run> < payload.json`.
- Evaluation status, list, and project status now report v2 data artifacts and
  gaps instead of previous-runtime record counts or active recommendation counts.

### /quality Skill

- `/quality evaluate` guidance now persists routine outputs through stdin-only
  `qualitymd evaluation data set`, treats unsupported old run shapes as a reason
  to create a fresh Evaluation v2 run, and no longer routes from active
  recommendation counts.

### Specification

- Active specs now point to Evaluation v2 as the only runtime evaluation
  contract. Previous evaluation-record and legacy report artifact specs and
  old checked-in skill examples were removed from the active spec tree.

Compatibility:

- CLI: `v0.13.0`
- QUALITY.md specification: `0.4 (Draft)`
- /quality skill: `0.13.0`, requires `qualitymd >=0.13.0 <0.14.0`

## v0.12.0 - 2026-06-25

### QUALITY.md Format

- Requirements now use stable id-like names as map keys and carry the
  human-facing statement in required `title`. Requirement `assessment` remains
  required, `description` is optional, and tools now support qualified
  `requirement:<area-path>::<requirement-name>` model references. This is a
  breaking format change for existing statement-key Requirements.

### CLI

- Evaluation runs now have an Evaluation v2 data surface:
  `qualitymd evaluation data set/list/get/kinds/example`. Agents persist
  structured routine JSON under `data/`, `evaluation report build` assembles
  `data/evaluation-output-result.json`, and the old public
  `evaluation assessment`, `evaluation analysis`, `evaluation recommendation`,
  and `evaluation report gate` command surfaces are removed.
- `qualitymd init` and `qualitymd init --minimal` now seed the standard
  four-level Rating Scale with emoji-prefixed human titles (`🟢 Outstanding`,
  `🔵 Target`, `🟡 Minimum`, `🔴 Unacceptable`) while keeping stable `level` IDs
  plain.

### /quality Skill

- `/quality evaluate` now follows the Evaluation v2 framing, assessment, rating,
  Factor analysis, Area analysis, and deterministic reporting protocol. New
  evaluations persist structured routine outputs through
  `qualitymd evaluation data set` and no longer generate recommendation
  artifacts as part of the v0 evaluation flow.
- `/quality setup` now stays focused on producing a useful initial `QUALITY.md`:
  it no longer asks for future recommendation handling, handoff destination,
  review cadence, recurring review, or automation preferences, and its closeout
  reports lint validation plus important model gaps instead of a
  maturity/evaluation-ready label.
- The authoring guide now defines the agent harness holistically as the whole
  engineered system around the model, then scopes the agent-harness area to
  checked-in steering and project-owned control artifacts. Setup now checks for
  owned runtime harness machinery, and the Top 10 checks flag instructions-only
  harness areas or unmodeled owned harness controls.
- The `/quality` skill now treats agent-mediated UX as part of its workflow
  contract: setup, evaluate, update, and recommendation follow-up guidance now
  emphasize the primary question or call to action, use scannable labels for
  recommendations/evidence/validation/next steps, keep emoji semantic, and close
  with clearer status and next-action summaries.
- `/quality setup` and the authoring guide now recommend emoji-prefixed Rating
  Level titles for the standard four-level scale as a human scanning aid, without
  treating emoji as rating semantics or a format requirement.
- The authoring guide now teaches three area decomposition shapes —
  primary-subject, collection, and composite — that apply recursively at every
  area, names the two recurring use-context constituents (the agent harness and
  the QUALITY.md self-check), and
  scopes the ~ten-factor coverage aim to a primary-subject node rather than the
  root. The Top 10 QUALITY.md checks gain matching routing findings for a
  composite entity flattened into one root and for a missing expected
  constituent.
- The authoring and evaluation guidance now treats the `quality-md` self-check as
  an ordinary in-scope area: full `/quality evaluate` covers it when present, and
  its factor, local, and aggregate ratings roll up like any other modeled area.
  Quality-log entries remain reserved for meaningful confirmed model changes.
- The authoring guide now teaches enumerating a composite root's domain
  constituents by _constituent kind_ — drawn from the entity's stewardship
  concerns (a lifecycle band plus the protective pair _secure_ and _safeguard_)
  and an audience×purpose axis (Diátaxis) — so high-leverage constituents such as
  tests, specs, docs, and a threat model are modeled or consciously accounted for
  instead of silently missing. Setup and the Top 10 QUALITY.md checks gain
  matching guidance and a routing finding.
- The authoring guide's stewardship-concern guidance is now grounded in the
  phenomenology of care: a constituent kind is a prompt for what the entity asks
  to be cared for (the claim comes from a Need or Risk, not the list), its artifact
  is read as evidence of tending rather than proof of quality, the protective pair
  _secure_/_safeguard_ answers to who is vulnerable, and the lifecycle concerns
  recur rather than completing. Framing only — the concerns and the earn-it
  guardrail are unchanged.
- The authoring guide and setup now keep the stewardship/care language in its
  motivation register so it never modifies or replaces a model term: recurring
  root factors are named _model-wide factors_ (which may trace to stewardship
  concerns) rather than "stewardship factors" or "stewardship lenses." Wording
  discipline only — the concerns, axes, and three-projections rule are unchanged.
- The authoring guide, setup, and the Top 10 QUALITY.md checks now **model a
  composite root's constituents by default** instead of treating each area as
  something to earn. A constituent is left without its own area only when it has no
  distinct concerns (folded into a sibling) or is not germane to the entity. A
  germane concern is never dropped to prose: a thin or missing artifact is surfaced
  as a ratable gap — a minimal area with a missing-anchor finding, or a requirement
  on an existing area. Deferral is now a narrow, blocker-recorded exception, and a
  composite model that leaves a germane constituent unmodeled or merely deferred no
  longer passes setup without an important-gap note. This makes a first-pass model
  as full as the evidence supports.
- The authoring guide now teaches Agent Harnessability (`agent-harnessability`) as
  the model-wide factor projection of the agent-collaboration concern for
  agent-collaborated composite roots, defined around project materials and tooling
  that equip agent work while preserving human direction, review, and
  accountability. It now uses seven sub-factors:
  agent-accessibility, task-specifiability, agent-operability, continuity,
  self-verifiability, enforcement-of-standards, and containment-of-action. Setup
  proposes it by default, the Top 10 checks flag missing coverage, legacy
  six-sub-factor `harnessability` can still count as prior semantic coverage, and
  thin or absent harnesses are treated as rating evidence rather than omission
  reasons.
- The authoring guide now gives the agent-harness area its own domain-agnostic
  steering-materials template, with illustrative factors such as completeness,
  accuracy, currentness, understandability, coherence, selectivity,
  discoverability/triggerability, maintainability, trustworthiness, and
  assessability. The Top 10 checks now flag a germane harness area carried by only
  one or two thin factors, and harness requirements that assume a software
  toolchain outside a software served domain.
- The authoring guide now teaches encoding a concern's _projection boundary_ in the
  model: when a model carries two or more projections of one concern (factor,
  constituent/area, audience), each node gets a YAML comment naming its sibling
  projection and the distinction, and — when both projections are rated nodes that
  surface in an evaluation report — a short disambiguating clause in each
  description. The Agent Harnessability factor vs. the agent-harness area is the
  canonical instance, and the Top 10 checks flag two same-rooted projections carried
  with no boundary note.
- The `/quality` reference examples now pair the existing software-service fixture
  with a complete non-software data-product fixture,
  `0002-city-bike-stations-quality-eval`, using the same model, record, report, and
  recommendation shape. The example index marks the corpus as domain-illustrative,
  and the Top 10 checks now include a data-product constituent bracket alongside the
  software illustration.
- The `/quality` authoring guidance is now split into an OKF-shaped guide family:
  `guides/authoring.md` remains the entry point and router, while routed sub-guides
  cover body authoring, model structure, factors, requirements, rating scale,
  Agent Harnessability, the agent-harness area, and quality-log judgment. The
  durable skill specs mirror that tree so each concern can be reviewed and improved
  independently.

### Packaging

### Specification

### Documentation

- The domain-agnostic modeling guide now distinguishes explicit guidance for
  recurring use-context constituents (the agent harness and QUALITY.md self-check)
  from modeled-domain defaults, and adds a served-domain guardrail so harness and
  self-check requirements stay agnostic to what the project evaluates.
- Re-scoped the README opening so QUALITY.md's modeled domain is broad while the
  `/quality` skill remains the primary experience. The Agent Harnessability section
  now presents that factor family as illustrative and earned for agent-collaborated
  entities, not a default for every QUALITY.md.
- The specification's minimal example now notes that model shape is invariant across
  domains and links to the worked non-software guide example; the non-normative
  lineage section clarifies that QUALITY.md borrows boundaries and vocabulary from
  prior traditions, not their characteristic lists as default Factors.
- Added a contributor guide, _Modeling quality across domains_, that consolidates
  the quality-domain-agnosticism doctrine, names the stress axes and a canonical
  set of secondary knowledge-work domains (documentation, data sets, research
  reports, services and operations), and carries a full worked non-software
  example. Standardized the domain enumerations across the README, specification,
  and skill guidance so they name a consistent illustrative set.

Compatibility:

- CLI: `v0.12.0`
- QUALITY.md specification: `0.3 (Draft)`
- /quality skill: `0.12.0`, requires `qualitymd >=0.12.0 <0.13.0`

## v0.11.0 - 2026-06-23

### CLI

- `qualitymd evaluation create` no longer seeds `debug-log.md` in new evaluation
  run folders. Historical runs that already contain `debug-log.md` remain
  compatible as legacy process artifacts.

### /quality Skill

- `/quality evaluate` now uses a setup-style workflow feedback log at
  `.quality/logs/<timestamp>-evaluate-feedback-log.md` for material
  workflow-experience events such as scope friction, retries, redaction
  decisions, slow phases, and report recovery. The log is local-only, never
  transmitted automatically, and must not duplicate ratings, findings, evidence,
  raw command output, recommendations, secrets, or raw prompt-injection text.
- Reworked `/quality setup` guidance into an explicit setup workflow with a
  setup brief, concrete discovery questions, confidence-labeled defaults, and
  clearer workflow terminology while preserving the `QUALITY.md`-only mutation
  boundary.
- Made setup discovery agent-agnostic: the workflow presents all discovery questions
  using the agent's own interaction surface — iterating one question at a time
  when no structured question tool is available, or paging through one when it
  is — instead of assuming a fixed question UI.
- Setup now reads the `qualitymd init` scaffold before authoring it, avoiding a
  wasted read-before-write round-trip.
- The Top 10 checklist now keeps lifecycle `readiness` from `qualitymd status`
  distinct from model-usefulness findings instead of blending lifecycle state
  with setup closeout judgment.
- Renamed the skill's `modes/` folder to `workflows/` (runtime skill and spec
  mirror); each workflow is still dispatched as a mode.
- Made `/quality setup` discovery teaching-first: each setup question now
  carries authored copy explaining what it shapes in the model, setup asks every
  discovery question every run (no accept-all-and-skip shortcut; a per-question
  fast confirm and show-all-at-once remain), and confidence is now reported as
  `Low`/`Med`/`High` with an evidence note instead of `strongly inferred` /
  `weakly inferred` / `assumed`.
- Refined `/quality setup` discovery: removed the modeling-rigor and
  review-posture questions, added a Rating Scale confirmation question that
  teaches Rating Levels are configurable while recommending
  `outstanding`/`target`/`minimum`/`unacceptable`, and made the final recap a
  review gate that requires an explicit user response before setup writes or
  edits `QUALITY.md`.
- Made the `/quality setup` final recap more open-ended and conversational:
  users can still say `"looks good"` to proceed, but the prompt now explicitly
  invites priorities, worries, wording, edge cases, repo-invisible context, or
  anything else that should shape `QUALITY.md` before it is written.
- Replaced `/quality setup`'s final four open-ended human-context questions with
  a compact context checkpoint: setup now asks users to correct a draft covering
  primary users/outcomes, maintainers/collaborators, other stakeholders, and
  missing or not-agent-accessible context, while recording omitted
  low-confidence gaps as Unknown rather than confirmed fact.
- Tightened `/quality setup` missing-context discovery so low/no-evidence
  product, operational, stakeholder, telemetry, security/compliance, incident,
  SLA, or production-metric context is recorded as unknown/not-agent-accessible
  unless the user explicitly provides it or points setup to missed
  agent-accessible evidence.
- Added an always-on setup feedback log: after CLI support is verified and the
  run frame is emitted, setup creates
  `.quality/logs/<timestamp>-setup-feedback-log.md`, updates the current run's
  file for material workflow-experience events, and finalizes it at close so
  clean and interrupted runs leave explicit local feedback records. It is
  recorded locally and never transmitted — sharing is an explicit user action —
  and never contains secrets or raw prompt-injection text. `.quality/logs/`
  (plural) is distinct from the quality log's `.quality/log/`, which setup still
  does not write.
- Bumped the bundled skill release metadata to `0.11.0`; the supported
  `qualitymd` CLI compatibility range is now `>=0.11.0 <0.12.0`.

### Packaging

- No packaging changes.

### Specification

- No QUALITY.md specification changes.

Compatibility:

- CLI: `v0.11.0`
- QUALITY.md specification: `0.3 (Draft)`
- /quality skill: `0.11.0`, requires `qualitymd >=0.11.0 <0.12.0`

## v0.10.0 - 2026-06-23

### CLI

- Added `qualitymd init --minimal`, which writes a lint-valid starter
  `QUALITY.md` without guided template prose for agent or human authoring passes
  that replace the model content wholesale.
- Changed `qualitymd status --json` source coverage to report implicit document
  directory resolution as `sourceState: "default"` instead of a defect-like
  missing source for lint-valid models.

### /quality Skill

- Reworked `/quality setup` into a context-informed setup flow that asks a few
  discovery questions with recommended defaults, writes only `QUALITY.md`,
  validates/readiness-checks the model, and offers next-step choices without
  running evaluation, writing the quality log, creating issues, or configuring
  automation.
- Removed `wizard` from the documented `/quality` public contract. Bare
  `/quality` remains a read-only orientation entrypoint, while public workflows
  are `setup`, `evaluate`, `update`, and recommendation follow-up.
- Made natural Area and Factor labels the primary documented `/quality evaluate`
  scoped input, with qualified `area:` and `factor:` references retained for
  exact addressing.
- Clarified setup and authoring guidance for quality-factor coverage, bounded
  context discovery, domain identification, agent-accessible evidence, and
  quality-domain agnostic examples.
- Bumped the bundled skill release metadata to `0.10.0`; the supported
  `qualitymd` CLI compatibility range is now `>=0.10.0 <0.11.0`.

### Packaging

- No packaging changes.

### Specification

- Clarified that QUALITY.md is quality-domain agnostic and adjusted
  non-normative examples so software product quality is illustrative rather than
  the default.

### Documentation

- Clarified agentic use context: AI assistants, coding agents, harnesses,
  agent-accessible evidence, and `/quality` skill workflows are valid use
  context while concrete model content remains quality-domain agnostic.

### Compatibility / Migration

Compatibility:

- CLI: `v0.10.0`
- QUALITY.md specification: `0.3 (Draft)`
- /quality skill: `0.10.0`, requires `qualitymd >=0.10.0 <0.11.0`

Existing `QUALITY.md` files do not require migration. Consumers of
`qualitymd status --json` that distinguish source coverage states should accept
the new `default` state for lint-valid areas that resolve to the document
directory.

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
