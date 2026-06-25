# Changelog

User-facing release notes for `qualitymd`, the `/quality` skill, and the
QUALITY.md specification.

## Unreleased

### CLI

- `qualitymd init` and `qualitymd init --minimal` now seed the standard
  four-level Rating Scale with emoji-prefixed human titles (`🟢 Outstanding`,
  `🔵 Target`, `🟡 Minimum`, `🔴 Unacceptable`) while keeping stable `level` IDs
  plain.

### /quality Skill

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
  constituents by *constituent kind* — drawn from the entity's stewardship
  concerns (a lifecycle band plus the protective pair *secure* and *safeguard*)
  and an audience×purpose axis (Diátaxis) — so high-leverage constituents such as
  tests, specs, docs, and a threat model are modeled or consciously accounted for
  instead of silently missing. Setup and the Top 10 QUALITY.md checks gain
  matching guidance and a routing finding.
- The authoring guide's stewardship-concern guidance is now grounded in the
  phenomenology of care: a constituent kind is a prompt for what the entity asks
  to be cared for (the claim comes from a Need or Risk, not the list), its artifact
  is read as evidence of tending rather than proof of quality, the protective pair
  *secure*/*safeguard* answers to who is vulnerable, and the lifecycle concerns
  recur rather than completing. Framing only — the concerns and the earn-it
  guardrail are unchanged.
- The authoring guide and setup now keep the stewardship/care language in its
  motivation register so it never modifies or replaces a model term: recurring
  root factors are named *model-wide factors* (which may trace to stewardship
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
  longer reaches `evaluation-ready`. This makes a first-pass model as full as the
  evidence supports.
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
- The authoring guide now teaches encoding a concern's *projection boundary* in the
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
- Added a contributor guide, *Modeling quality across domains*, that consolidates
  the quality-domain-agnosticism doctrine, names the stress axes and a canonical
  set of secondary knowledge-work domains (documentation, data sets, research
  reports, services and operations), and carries a full worked non-software
  example. Standardized the domain enumerations across the README, specification,
  and skill guidance so they name a consistent illustrative set.

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
- Setup now reports model *maturity* (`starter`, `immature`, `evaluation-ready`)
  as distinct from the lifecycle `readiness` that `qualitymd status` owns, and
  the Top 10 checklist no longer blends the two axes.
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
