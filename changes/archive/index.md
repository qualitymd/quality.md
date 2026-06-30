# Archived change cases

Completed change cases, moved here from the bundle root when they reach **Done**.

# Change cases

- [0001 — Example change](0001-example-change.md) - placeholder retired as a
  reference template for the Change Case concept shape (`Done`).
- [0002 — Specify the init command](0002-init-command.md) - settled and shipped
  `qualitymd init` (`Done`).
- [0003 — Implement the lint command](0003-implement-lint-command.md) - built
  `qualitymd lint` from the completed durable lint sub-spec (`Done`).
- [0004 — Specify and enforce agent accessibility](0004-specify-agent-accessibility.md) - added the CLI agent-accessibility contract, broadened `--json`, and enforced categorized exit codes plus `init --json` (`Done`).
- [0005 — Single source of truth for the structural schema](0005-schema-source-of-truth.md) - extracted the structural schema into one typed declaration consumed by `lint` (`Done`).
- [0006 — Specify and implement the spec command](0006-spec-command.md) -
  settled and shipped `qualitymd spec`, emitting the bundled format specification
  from the binary (`Done`).
- [0007 — Delightful human CLI output](0007-delightful-cli-output.md) - gave the
  human surface a shared brand palette, styled `lint` and `init` output, `--help`
  examples, `spec` paging, and an informative `--version`, all behind the
  TTY/`NO_COLOR` gate so the plain and JSON paths are untouched (`Done`).
- [0008 — Describe targets with title and description](0008-target-display-fields.md) -
  lets every target carry a recommended `title` and optional `description`, and
  reframes the root as a Model (`ratingScale` + Target properties) so
  `ratingScale` is the one Model-only key (`Done`).
- [0009 — Diagnose rating-scale soundness in the meta-model](0009-rating-scale-diagnostic.md) -
  adds a meta-model Functionality requirement that judges a model's rating scale
  and per-requirement criterion overrides for meaning, not only structure
  (`Done`).
- [0010 — Implement the /quality skill](0010-implement-quality-skill.md) - ships
  the `/quality` skill artifact, the `qualitymd models` bundled-model surface,
  skill-first onboarding docs, raw JSON evaluation example artifacts, and durable
  spec sync (`Done`).
- [0011 — CLI human output polish](0011-cli-human-output-polish.md) - finishes
  the remaining styled-output, lint next-action, dev-version, and gate-coverage
  work (`Done`).
- [0012 — Evaluation record format](0012-evaluation-record-format.md) - lifted
  the evaluation artifact contract out of the skill prompt into the enduring
  `specs/evaluation-records.md` spec the CLI writes and the skill consumes
  (`Done`).
- [0013 — Evaluation run scaffold](0013-evaluation-run-scaffold.md) - added
  `qualitymd evaluation create-run` with deterministic shared run numbering and
  run-folder scaffolding (`Done`).
- [0014 — Evaluation record write](0014-evaluation-record-write.md) - added
  `qualitymd evaluation add-record assessment|analysis|recommendation` with
  schema validation and atomic numbered writes (`Done`).
- [0015 — Evaluation status and report build](0015-evaluation-report-build.md) -
  added `qualitymd evaluation show-status` and `build-report` over a shared
  renderability gate, with deterministic `report.md`/`report.json` and the
  `--fail-at-or-below` CI gate (`Done`).
- [0016 — Skill consumes evaluation CLI](0016-skill-consume-eval-cli.md) -
  switched the `/quality` skill to drive the evaluation CLI for scaffolding,
  record writes, and reports, replacing the inlined Artifact Contract with a
  reference (`Done`).
- [0017 — Skill rigor and efficiency](0017-skill-rigor-efficiency.md) -
  operationalized effort levels, evidence and pinned-locator rigor, the
  rating-binding re-check, batched writes, and confined deep fan-out (`Done`).
- [0018 — Evaluation report UX](0018-evaluation-report-ux.md) - made generated
  reports summary-first, scoped, and easier to scan, verified on copied ESLint
  and DataLoader runs (`Done`).
- [0019 — Duplicate assessment status](0019-duplicate-assessment-status.md) -
  made duplicate assessments for the same target requirement a reportability
  gap (`Done`).
- [0020 — Planned coverage status](0020-planned-coverage-status.md) - added
  `qualitymd evaluation set-planned-coverage` and planned-coverage status gaps so
  interrupted or resumed runs can name missing planned work (`Done`).
- [0021 — Recommendation superseding](0021-recommendation-superseding.md) - let
  corrected recommendation records supersede stale recommendations so reports
  choose the active Next Action deterministically (`Done`).
- [0022 — Create-run subject validation](0022-create-run-subject-validation.md) -
  validated `create-run --subject` before creating run folders so bad paths leave
  no partial evaluation artifacts (`Done`).
- [0023 — Assessment superseding](0023-assessment-superseding.md) - let corrected
  assessment records supersede stale assessments while requiring analyses to
  reference active records (`Done`).
- [0024 — Report regression coverage](0024-report-regression-coverage.md) - added
  focused tests for high-risk generated report behavior found by the experiment
  program (`Done`).
- [0025 — Durable spec rationale](0025-durable-spec-rationale.md) - taught the
  contributor guides so durable specs carry their *why* — a spec-level
  Background/Motivation section and per-requirement annotations — so a landing
  change absorbs its rationale instead of leaving it in the archive (`Done`).
- [0026 — Authoring guide replaces meta-model workflow](0026-authoring-guide-remove-meta-model.md) -
  replaced the bundled quality meta-model workflow with a practical `QUALITY.md`
  authoring guide, removed the public `qualitymd models` surface, and made
  evaluation run creation subject-only (`Done`).
- [0027 — Modularize quality skill modes](0027-modularize-quality-skill.md) -
  split setup, wizard, evaluate, and improve procedures into mode files while
  keeping `SKILL.md` as the root router and global contract (`Done`).
- [0028 — Require factor references](0028-require-characterized-requirements.md) -
  requires every requirement to reference at least one factor and aligns
  requirement-to-factor terminology (`Done`).
- [0029 — Sharpen assessment references and traceability](0029-sharpen-assessment-references.md) -
  frames a requirement's assessment as inline or a reference to another entity,
  reserves "source" for `Target.source`, and makes the model's traceability
  graph an authoring concern (`Done`).
- [0030 — CLI status snapshot command](0030-cli-status-command.md) - added
  `qualitymd status [path] [--json]` so agents and humans can route from a
  deterministic project-state snapshot instead of hand-parsing `QUALITY.md` or
  evaluation reports (`Done`).
- [0031 — Evaluation report summary artifact](0031-report-summary-artifact.md) -
  generated `report-summary.md` beside full evaluation reports for compact
  review and triage surfaces (`Done`).
- [0032 — CLI managed upgrades](0032-cli-managed-upgrades.md) - added structured
  version metadata, explicit upgrade checks, managed installer entrypoints, and
  install-method-aware upgrade guidance (`Done`).
- [0033 — Required display titles](0033-required-display-titles.md) - required
  human-facing titles on Models, Targets, Factors, and Rating Levels, with lint,
  scaffold, docs, and examples updated to match (`Done`).
- [0034 — Skill release metadata](0034-skill-release-metadata.md) - recorded the
  `/quality` skill version and required `qualitymd` CLI range in Agent
  Skills-compatible metadata with release-check validation (`Done`).
- [0035 — /quality upgrade mode](0035-quality-skill-upgrade-mode.md) - added a
  `/quality upgrade` mode that plans paired skill/CLI maintenance, delegates
  mutation to owner tools, and reports restart/reload needs after skill updates
  (`Done`).
- [0036 — Harden install scripts and upgrade idiomatics](0036-harden-install-scripts.md) -
  fixed portability/convention gaps in the shell/PowerShell installers and the
  upgrade version check, and recorded why the Homebrew cask is the idiomatic path
  (`Done`).
- [0037 — Render display titles in evaluation reports](0037-report-rating-titles.md) -
  made human reports display model, target, factor, and rating-level titles
  while preserving stable ids for machine reports and gates (`Done`).
- [0038 — /quality skill interaction UX](0038-quality-skill-interaction-ux.md) -
  added a durable user interaction contract for the `/quality` skill covering run
  frames, decision briefs, stop/reroute behavior, history-aware operation,
  improvement delta reports, and status-first output (`Done`).
- [0039 — Evaluation command surface redesign](0039-evaluation-command-surface.md) -
  reshaped the `qualitymd evaluation` surface around noun/verb resources, folded
  planned coverage into `plan.md`, split report build from gate, added run and
  record listing, batched writes, and removed create altitude residue (`Done`).
- [0040 — Readable report summary](0040-readable-report-summary.md) - reshaped
  `report-summary.md` into a decision-brief artifact with reader-facing
  vocabulary and prominent recommendation identifiers (`Done`).
- [0041 — Update command and improvements](0041-update-command.md) - renamed
  `qualitymd upgrade` to apply-by-default `qualitymd update`, added the cached
  ambient update notice and readiness-aware update checks, and renamed the paired
  `/quality` maintenance mode to `update` (`Done`).
- [0042 — Typed report model](0042-typed-report-model.md) - replaced stringly
  typed and implicit evaluation-report states with explicit typed concepts for
  ratings, local target state, next steps, lifecycle state, run gaps, rigor,
  evaluation level, missing metadata, and path identities (`Done`).
- [0043 — Evaluation history compatibility](0043-evaluation-history-compatibility.md) -
  made evaluation-history readers tolerant of old, malformed, partial, or
  hand-edited run records while keeping current writers strict and report
  generation explicit (`Done`).
- [0044 — Section unknowns and open questions](0044-section-unknowns-open-questions.md) -
  replaced the standalone Known gaps body section with per-section unknowns,
  open questions, and a human/agent review state line (`Done`).
- [0045 — Evaluable body context](0045-evaluable-body-context.md) - clarified
  the Markdown body as concise, self-explanatory, agent-accessible judgment
  context for building, justifying, and evaluating model quality (`Done`).
- [0046 — Evaluation debug log](0046-evaluation-debug-log.md) - added a
  process-only `debug-log.md` artifact to evaluation runs (`Done`).
- [0047 — Area terminology changeover](0047-area-terminology.md) - replaced the
  formal Target terminology with Area across the live schema, records, reports,
  CLI, skill, scaffold, examples, and docs (`Done`).
- [0048 — Area factor report breakdown](0048-area-factor-report-breakdown.md) -
  exposed a compact Area-by-Factor breakdown from a first-class report model
  across `report-summary.md`, `report.md`, and `report.json`, renamed the Area
  rating fields, and added the durable `specs/reports/` artifact specs (`Done`).
- [0049 — Companion JSON Schema](0049-companion-json-schema.md) - published a
  structural, non-normative JSON Schema for QUALITY.md frontmatter, generated
  from `internal/schema` and guarded against drift, embedded via a root
  `schema.go`, and emitted by a new `qualitymd schema` command (`Done`).
- [0050 — Quality log](0050-quality-log.md) - added the convention-first quality
  log: dated `quality/log/` entries the `/quality` skill writes (`setup` seeds,
  `improve` appends) to record meaningful, evidence-linked model changes, with the
  format contract in `SKILL.md`, the meaningful-change taxonomy in the authoring
  guide, `wizard` reconciliation, and a new `## Quality log` section in the
  durable skill spec (`Done`).
- [0051 — Setup quality-md Area](0051-setup-quality-md-area.md) - added a
  setup-authored `quality-md` Area pattern that evaluates the active `QUALITY.md`
  artifact against the active authoring guide, with YAML comments explaining
  `source` versus `assessment` and stronger authoring guidance for Factor names
  and referenced assessments (`Done`).
- [0052 — Durable spec alignment](0052-durable-spec-alignment.md) - aligned
  durable specs with artifact-spec versus behavioral-component guidance,
  splitting the `/quality` skill parent into child specs for modes, evaluation
  workflow, reporting, and quality log, and strengthening the general
  spec-splitting guidance (`Done`).
- [0053 — Align remaining durable specs](0053-align-remaining-durable-specs.md) -
  split remaining large durable specs for evaluation records, lint, and ambient
  update notices into parent and component/artifact contracts (`Done`).
- [0054 — Remove improve mode](0054-remove-improve-mode.md) - removed
  `/quality improve` as a public mode while preserving recommendation follow-up
  with apply-now and issue-tracker handoff outcomes (`Done`).
- [0055 — Self-describing evaluation record input](0055-evaluation-input-ergonomics.md) -
  made evaluation record payloads discoverable and validatable through
  payload-documenting help, no-persist dry-runs, aggregated key-named validation,
  synced skill surfaces, and a published-skill link guard (`Done`).
- [0056 — Prospective evaluation plan artifacts](0056-prospective-evaluation-plan-artifacts.md) -
  made `design.md` and the initial `plan.md` prospective `/quality evaluate`
  artifacts authored before assessment begins, with later plan changes recorded
  as amendments (`Done`).
- [0057 — Quality data directory](0057-quality-data-directory.md) - defined the
  QUALITY.md workspace envelope, moved live support artifacts under `.quality/`,
  added the root `config` tooling pointer, and kept strict lint handling for
  unknown keys (`Done`).
- [0058 — Model reference identifiers](0058-model-reference-identifiers.md) -
  defined strict Area/Factor/Rating names, canonical typed model references,
  edge-only shorthand boundaries, and clearer Area Breakdown report columns
  (`Done`).
- [0059 — Unqualified model references](0059-unqualified-model-references.md) -
  defined bounded unqualified Area, Factor, and Rating references for fixed-type
  contexts and used them for Area-only report summary paths (`Done`).
- [0060 — Friendly path display](0060-friendly-path-display.md) - separated
  human display values from model-reference grammar so root Area paths render as
  `/` in human reports while references stay stable (`Done`).
- [0061 — Natural scope labels](0061-natural-scope-labels.md) - made natural
  Area and Factor labels the primary documented scoped-evaluation input for the
  `/quality` skill while preserving qualified references for exact addressing
  and stable identifiers in durable artifacts (`Done`).
- [0062 — Remove wizard mode](0062-remove-wizard-mode.md) - removed `wizard`
  from the `/quality` skill's public contract while preserving safe read-only
  orientation for ambiguous requests (`Done`).
- [0063 — Contextual setup flow](0063-contextual-setup-flow.md) - reworked
  `/quality setup` into a short context-informed discovery flow that writes only
  `QUALITY.md` and routes next steps (`Done`).
- [0064 — Structured setup workflow](0064-structured-setup-workflow.md) - turned
  `/quality setup` guidance into an explicit workflow with concrete discovery
  questions and prompt framing (`Done`).
- [0065 — Setup discovery and close refinements](0065-setup-discovery-and-close-refinements.md) -
  made setup discovery agent-agnostic with one-question-at-a-time iteration, added
  read-before-author, disentangled model maturity from CLI lifecycle readiness,
  and renamed the skill `modes/` folder to `workflows/` (`Done`).
- [0066 — Setup feedback log](0066-setup-feedback-log.md) - added a hand-authored,
  skill-only workflow feedback log under `.quality/logs/`
  (`<timestamp>-setup-feedback-log.md`) recording the setup run *experience*,
  distinct from the quality log and evaluation `debug-log.md`, recorded locally
  and never transmitted, with secrets/prompt-injection excluded and sensitive
  context sanitized; no CLI/Go change (`Done`).
- [0067 — Setup discovery pedagogy](0067-setup-discovery-pedagogy.md) - made
  `/quality setup` discovery teaching-first: authored per-question background and
  how-to-change-later copy inline in the skill, ask every question every run (no
  accept-all-and-skip escape; per-question fast confirm and show-all-at-once
  remain), relabeled confidence to `Low`/`Med`/`High` with an evidence note, and
  added a final review recap before writing `QUALITY.md`; no CLI/Go change
  (`Done`).
- [0068 — Always-on setup feedback log](0068-always-on-setup-feedback-log.md) -
  made `/quality setup` create a workflow feedback log after CLI support and the
  run frame, update the current run's file for material workflow-experience
  events, and finalize it at close with stable metadata and explicit clean-run or
  interrupted-run feedback; no CLI/Go change (`Done`).
- [0069 — Setup review gate and discovery trim](0069-setup-review-gate-and-pedagogy-trim.md) -
  made `/quality setup` remove the modeling-rigor and review-posture discovery
  questions, add Rating Scale confirmation, trim per-question pedagogy to
  purpose/context, and require an explicit final recap response before authoring;
  no CLI/Go change (`Done`).
- [0070 — Setup missing-context provenance](0070-setup-missing-context-provenance.md) -
  tightened `/quality setup` missing-context discovery so low/no-evidence project
  facts are recorded as unknown/not-agent-accessible unless the user explicitly
  provides them or points setup to missed agent-accessible evidence; no CLI/Go
  change (`Done`).
- [0071 — Setup open-ended review gate](0071-setup-open-ended-review-gate.md) -
  made `/quality setup`'s final recap prompt friendlier and more open-ended while
  preserving the `"looks good"` fast path before authoring `QUALITY.md`; no
  CLI/Go change (`Done`).
- [0072 — Setup context checkpoint](0072-setup-context-checkpoint.md) - replaced
  `/quality setup`'s final four open-ended human-context discovery questions with
  a compact checkpoint that asks users to correct a draft and records omitted
  low-confidence gaps as Unknown rather than confirmed fact; no CLI/Go change
  (`Done`).
- [0073 — Evaluation feedback log](0073-evaluation-feedback-log.md) - aligned
  `/quality evaluate` with setup feedback logging by adding a shared workflow
  feedback-log spec, adding an evaluate adopter spec, moving current evaluate
  feedback to `.quality/logs/<timestamp>-evaluate-feedback-log.md`, and making
  `debug-log.md` legacy-compatible only for historical runs (`Done`).
- [0074 — Composite root areas and use-context constituents](0074-composite-root-areas.md) -
  taught the authoring guide three recursive, composable decomposition shapes
  (primary-subject, collection, composite) and the two recurring use-context
  constituents (agent harness and QUALITY.md self-check, the latter kept out of
  the entity roll-up), scoped the ~ten-factor coverage aim per primary-subject
  node, and added matching Top 10 routing findings; guide and guide-spec edits
  only, no CLI/Go change (`Done`).
- [0075 — Rating title emoji defaults](0075-rating-title-emoji-defaults.md) -
  made emoji-prefixed Rating Level titles the default for `qualitymd init`,
  `qualitymd init --minimal`, and `/quality setup`, while keeping stable
  `level` IDs plain and leaving `SPECIFICATION.md` neutral (`Done`).
- [0076 — Domain constituent kinds and stewardship concerns](0076-domain-constituent-kinds.md) -
  taught the authoring guide to enumerate a composite root's domain constituents
  by constituent kind, drawn from stewardship concerns (a lifecycle band plus the
  secure/safeguard protective pair) and an audience×purpose axis (Diátaxis), with
  a three-projections rule and the earn-it guardrail; added matching setup and
  Top 10 guidance; guide and guide-spec edits only, no CLI/Go change (`Done`).
- [0077 — Care-grounded stewardship concerns](0077-stewardship-care-grounding.md) -
  refined 0076's stewardship-concern generator with the phenomenology of care
  (Rousse & Spivak): the artifact is the *trace* of tending rather than the care
  itself, the claim that earns a constituent comes from a Need or Risk rather than
  the list, the protective pair answers to vulnerability, stewardship is
  holding-in-trust, and the lifecycle band is recurring; authoring guide and
  durable guide-spec edits only, no CLI/Go change (`Done`).
- [0079 — Stewardship vocabulary discipline](0079-stewardship-vocabulary-discipline.md) -
  kept the stewardship/care core language in its motivation register so it never
  modifies or replaces taxonomy nouns (factor, area, requirement, constituent,
  audience); removed the two "stewardship lenses" fusions in the authoring guide,
  added a register rule to AGENTS.md and the durable authoring-guide spec, and
  guarded the setup summary so factors are named as factors; documentation-only,
  no CLI/Go change (`Done`).
- [0080 — Model constituents by default](0080-model-constituents-by-default.md) -
  flipped the constituent-coverage guidance from earn-it/defer-freely to
  model-by-default: enumerate the implied kinds and model each as its own area
  unless a disqualifier holds (no distinct concerns → fold; not germane → out of
  scope), never omit a germane concern in prose (surface a thin/missing artifact as
  a missing-anchor area or a requirement on an existing area), demote deferral to a
  blocker-only exception, and bar an under-covered composite model from
  `evaluation-ready`; authoring guide, setup, Top 10, getting-started, and their
  spec mirrors, no CLI/Go change (`Done`).
- [0081 — Harnessability factor](0081-harnessability-factor.md) - added
  harnessability as the model-wide factor projection of the agent-collaboration
  concern for agent-collaborated composite roots, decomposed into six
  sub-factors, proposed by setup by default, and checked by the Top 10 guide;
  skill-guidance, spec-mirror, README, and CHANGELOG only, no CLI/Go change
  (`Done`).
- [0082 — Normalize QUALITY.md self-check roll-up](0082-normalize-quality-md-rollup.md) -
  removed the special out-of-roll-up treatment for the `quality-md` self-check
  area so it evaluates and rolls up like any other modeled area; skill-guidance,
  spec-mirror, and CHANGELOG only, with `QUALITY.md` assessed no-change and no
  CLI/Go change (`Done`).
- [0083 — Quality-domain agnosticism guide and secondary illustrations](0083-quality-domain-agnosticism.md) -
  added a contributor-doctrine guide on modeling quality across domains, with stress
  axes, the canonical secondary-domain set, range-finder illustrations, and a full
  worked documentation example; aligned `AGENTS.md`, README, `SPECIFICATION.md`
  lineage wording, bundled skill guidance, logs, and CHANGELOG; documentation and
  skill-guidance only, no code or format-schema change (`Done`).
- [0084 — Agent-mediated UX conformance](0084-agent-mediated-ux-conformance.md) -
  made the `/quality` skill's live workflow output follow the agent-mediated UX
  guide: status-first presentation, visually emphasized primary questions or
  calls to action, scannable labels, semantic emoji only, explicit mutation
  boundaries, and closeouts with clear next actions; skill-guidance, spec-mirror,
  AGENTS, and CHANGELOG only, no CLI/Go change (`Done`).
- [0085 — Agent Harnessability naming](0085-agent-harnessability-naming.md) -
  renamed the harnessability factor guidance to Agent Harnessability, using
  `agent-harnessability` as the recommended key, an accountability-preserving
  definition, and legacy `harnessability` recognition for semantic coverage;
  skill-guidance, spec-mirror, README, and CHANGELOG only, no CLI/Go change
  (`Done`).
- [0086 — Umbrella factor roll-up framing](0086-umbrella-factor-rollup-framing.md) -
  corrected guidance that overstated the Agent Harnessability umbrella factor as
  "not rated directly" / "does not roll up" to the accurate "carries no
  requirements of its own; its rating rolls up from its sub-factors"; bundled
  authoring guide and its spec mirror only, no code or format-spec change
  (`Done`).
- [0087 — Encode projection boundaries in the model](0087-encode-projection-boundaries.md) -
  added a general authoring rule that a concern's projection boundary is encoded in
  the emitted model when more than one projection is modeled (a YAML comment per
  node, plus a disambiguating description clause when both are rated nodes that
  surface in a report), with the Agent Harnessability factor vs. the agent-harness
  area as the canonical instance and a matching Top 10 missing-boundary-note check;
  bundled authoring guide and Top 10 checks, their spec mirrors, the guides log, and
  CHANGELOG only, no code or format-spec change (`Done`).
- [0088 — Domain-agnostic corpus alignment](0088-domain-agnostic-corpus-alignment.md) -
  added the `0002-city-bike-stations-quality-eval` non-software data-product
  fixture with the same reportable runtime artifact shape as `0001`, marked and
  cross-linked the example corpus as domain-illustrative, added the earned Factors
  rule to `AGENTS.md`, re-scoped the README modeled-domain framing while preserving
  the agent-first use context, and reinforced the format spec, Top 10, reporting,
  and report-summary specs; docs/spec logs and CHANGELOG only, no CLI or Go behavior
  change (`Done`).
- [0089 — Agent-harness modeling guidance](0089-agent-harness-modeling-guidance.md) -
  added `continuity` to Agent Harnessability, strengthened self-verifiability and
  adjacent harnessability sub-factor guidance, gave the agent-harness area a
  domain-agnostic steering-materials factor and requirement template, added the
  use-context-constituent served-domain guardrail, and extended the Top 10 checks
  for thinly factored or software-leaking harness areas; docs/spec/skill guidance
  only, no CLI/Go or format-spec change (`Done`).
- [0090 — Skill-content OKF authoring split](0090-skill-content-okf-authoring-split.md) -
  split `/quality` authoring guidance into an OKF-shaped entry guide plus routed
  sub-guides, mirrored that tree in durable skill sub-specs, added runtime skill
  indexes/schema/logs, and made routed authoring reads explicit in `SKILL.md`,
  setup, and recommendation follow-up; docs/spec/skill guidance only, no CLI/Go
  or format-spec change (`Done`).
- [0091 — Agent-harness holistic definition](0091-agent-harness-holistic-definition.md) -
  redefined the agent harness as the whole engineered system around the model,
  scoped the agent-harness area to checked-in steering and project-owned control
  artifacts, added the mixed-artifact scoping rule and owned-runtime-machinery
  checks, and mirrored the guidance across setup, Top 10, durable specs, logs,
  docs, and CHANGELOG; docs/spec/skill guidance only, no CLI/Go or format-spec
  change (`Done`).
- [0092 — Setup workflow scope trim](0092-setup-workflow-scope-trim.md) - trimmed
  `/quality setup` so it no longer asks about future recommendation handling,
  handoff destination, review cadence, recurring review, or automation
  preferences; setup closeout now reports validation plus important model gaps
  instead of maturity/evaluation-readiness labels, and setup feedback logs record
  workflow outcomes; docs/spec/skill guidance only, no CLI/Go or format-spec
  change (`Done`).
- [0093 — Named Requirement identity](0093-requirement-identity.md) - changed
  Requirements from natural-language statement keys to stable id-like names with
  required `title`, optional `description`, retained `assessment`, qualified
  Requirement model references, lint/schema enforcement, migrated scaffolds,
  examples, docs, and skill guidance; breaking format change (`Done`).
- [0094 — Replace evaluation workflow](0094-replace-evaluation-workflow.md) -
  replaces the evaluation workflow with Evaluation v2 structured routine data,
  deterministic Area/Factor/Requirement reports, a new `evaluation data` CLI
  surface, superseded legacy evaluation report/record specs, and updated
  `/quality` runtime guidance (`Done`).
- [0095 — Evaluate feedback log outcomes](0095-evaluate-feedback-log-outcomes.md) -
  keeps `/quality evaluate` feedback logging as a process-only workflow artifact
  under `.quality/logs/` and makes the feedback-log `outcome` field describe
  workflow terminal states rather than report, rating, or recommendation
  semantics; docs/spec/skill guidance only, no CLI/Go or format-spec change
  (`Done`).
- [0096 — Setup intro preview](0096-setup-intro-preview.md) - makes `/quality
  setup` open with educational orientation, show a project-specific setup preview
  before discovery, and create the setup feedback log after that preview when the
  run continues; docs/spec/skill guidance only, no CLI/Go or format-spec change
  (`Done`).
- [0097 — Evaluation v2 clean break](0097-evaluation-v2-clean-break.md) -
  removes previous evaluation record/report compatibility, makes Evaluation v2
  the only active runtime path, switches data persistence to stdin-only input,
  deletes old active specs/examples, and prepares the v0.13.0 release (`Done`).
- [0098 — Setup opening as first output](0098-setup-opening-first-output.md) -
  makes `/quality setup` emit its warm welcome, phase roadmap, and run frame as
  the first output before any tool call, decoupling the run frame from the CLI
  prerequisite gate so the user is oriented before the read-only scan runs;
  docs/spec/skill guidance only, no CLI/Go or format-spec change (`Done`).
- [0099 — Closed-choice setup UX](0099-closed-choice-setup-ux.md) - makes
  `/quality setup` closed-choice prompts use numbered options with the
  recommended answer first and `1` as the shortest confirmation, while presenting
  cost labels for risk discovery and mapping them to risk-tolerance values
  internally; docs/spec/skill guidance only, no CLI/Go or format-spec change
  (`Done`).
- [0100 — Strengthen spec requirement standards (29148 + EARS)](0100-strengthen-spec-requirement-standards.md) -
  patches the functional-spec and change-case guides with set-level requirement
  checks, assumptions/dependencies, normative/informational references,
  unambiguous requirement language, and an optional EARS statement template
  (`Done`).
- [0101 — Quality skill UX action clarity](0101-quality-skill-ux-action-clarity.md) -
  aligns `/quality` setup, evaluate, update, and recommendation follow-up prompt
  shapes with the current agent-mediated UX guide: explicit CTAs, shortest
  answer paths, numbered ambiguity and outcome choices, decision briefs before
  `QUALITY.md` writes and external issue creation, update run-frame/progress
  sequencing, result closeout `Next` fields, and code-span precision; docs/spec/
  skill guidance only, no CLI/Go or format-spec change (`Done`).
- [0102 — Evaluation v2 report rating titles](0102-v2-report-rating-titles.md) -
  restores Rating Level title rendering in Evaluation v2 Markdown reports while
  preserving stable Rating Level IDs in routine JSON, `EvaluationOutputResult`,
  and build receipts (`Done`).
- [0103 — Evaluation v2 report enum display titles](0103-v2-report-enum-display-titles.md) -
  renders owned Evaluation v2 statuses, confidence levels, booleans, report
  kinds, limits, unknowns, and known finding classifications through typed human
  display titles with semantic emoji in Markdown reports while preserving raw
  values in routine JSON, `EvaluationOutputResult`, and build receipts (`Done`).
- [0104 — Evaluation v2 report header navigation](0104-evaluation-v2-report-header-navigation.md) -
  replaces generic Evaluation v2 report breadcrumbs and parent links with labeled
  `Area:` trails, Factor-only `Factor:` trails, and compact report-specific
  summary headers (`Done`).
- [0105 — Evaluation v2 report subject links](0105-evaluation-v2-report-subject-links.md) -
  moves generated report navigation from repeated `Details` columns into linked
  row subject cells in Evaluation v2 Area and Factor report tables (`Done`).
- [0106 — Binary confirmation UX](0106-binary-confirmation-ux.md) - makes
  `/quality` true binary mutation confirmations use visible `y`/`n` answer paths
  while preserving numbered responses for multi-option prompts; docs/spec/skill
  guidance only, no CLI/Go or format-spec change (`Done`).
- [0107 — Durable spec alignment](0107-durable-spec-alignment.md) - aligned the
  active durable spec bundle with current functional-spec and OKF guidance by
  auditing `SPECIFICATION.md` and 49 active durable spec concepts, adding missing
  BCP 14 convention declarations to routed authoring guide specs, and recording
  the alignment in spec logs; docs/spec guidance only, no CLI/Go or runtime
  behavior change (`Done`).
- [0108 — Short evaluation report filenames](0108-short-evaluation-report-filenames.md)
  - kept the root Evaluation v2 report at `report.md` while renaming descendant
    reports to short subject-aware filenames (`Done`).
- [0109 — Filename text for evaluation data links](0109-evaluation-data-link-filenames.md)
  - rendered the Evaluation v2 `Data`-column links with their payload filename
    text instead of the generic words `analysis`/`assessment`/`rating`
    (`Done`).
- [0110 — Run frame title and workflow vocabulary](0110-run-frame-and-workflow-vocabulary.md)
  - retitled the `/quality` run frame header to `**Quality · <workflow>**`
    (dropping the fake `/quality run` command and the `Mode:` field), added a
    durable run-frame requirement forbidding a command-style header or `Mode:`
    label, and retired "mode" in favor of "workflow" for the public-surface
    concept across the skill specs, runtime skill, and docs (`Done`).
- [0111 — Evaluation report rating labels](0111-evaluation-report-rating-labels.md)
  - labeled Evaluation v2 report rating columns explicitly (`Overall Rating` /
    `Local Rating`) and fixed the Factor and Sub-Area breakdown tables to show a
    local rating beside a descendant-inclusive `+ Sub-X Rating` (em dash when no
    descendants) instead of an aggregate rating beside a `Yes`/`No` boolean
    (`Done`).
- [0112 — Evaluation model snapshot filename](0112-evaluation-model-snapshot-filename.md)
  - renamed the Evaluation v2 run-folder model snapshot from `model.md` to
    `model-snapshot.md` (read/written through one `evaluation.ModelSnapshotFile`
    constant, clean break with no old-name reader) so the filename signals a
    frozen point-in-time copy of the working-tree model (`Done`).
- [0113 — Evaluation run folder naming](0113-evaluation-run-folder-naming.md) -
  shortened new Evaluation v2 run-folder names to `NNNN-full-eval` /
  `NNNN-<scope-path>-eval`, kept legacy `-quality-eval` folders recognized and
  numbered against, and documented full structural `--narrowing` slugs for
  scoped `/quality evaluate` runs (`Done`).
- [0114 — Run frame as first output](0114-run-frame-first-output.md) - lifted the
  "run frame is the first output, before any tool call" timing rule and a
  provisional / `resolving…` value allowance to the shared homes (the `SKILL.md`
  dispatcher instruction and the durable `Run frames` spec section), and brought
  the lagging `evaluate` workflow and its spec into line (`Done`).
- [0115 — Type-safe, model-bound Evaluation v2 data](0115-evaluation-data-typed-contract.md)
  - added a typed Evaluation data contract for strict field/type/enum validation,
    model-binding against the run snapshot, populated examples, a generated JSON
    Schema artifact, and the `data schema` / `data verify` CLI surfaces (`Done`).
- [0116 — Drop the "Evaluation v2" naming](0116-drop-evaluation-v2-naming.md) -
  renamed the live Evaluation spec bundle and live references from
  `specs/evaluation-v2/` / "Evaluation v2" to plain `specs/evaluation/` /
  "Evaluation", removed the superseded sketch, and dropped private `v2` Go report
  symbols (`Done`).
- [0117 — Requirement report Factors line](0117-requirement-report-factors-line.md)
  - added a plural `Factors:` context line to Requirement reports and removed the
    duplicated `Factors` column from their summary table (`Done`).
- [0118 — Report empty-cell marker and legend](0118-report-empty-cell-legend.md)
  - renders empty scalar report cells as `—` and adds one static legend per
    generated report (`Done`).
- [0119 — Report header kind prefix and title-first layout](0119-report-header-kind-prefix.md)
  - makes report H1s title-first and kind-prefixed, moves trails below the title,
    and drops redundant `Path:` / `Name:` header lines (`Done`).
- [0120 — String model-identity fields in evaluation data](0120-string-model-identity-fields.md)
  - stores Evaluation data model identities as canonical qualified-reference
    strings, reserves `root` as an Area name, and bumps Evaluation data schema
    version to 2 (`Done`).
- [0121 — Scannable interaction hierarchy](0121-scannable-interaction-hierarchy.md)
  - reshapes agent-mediated decision gates, discovery questions, progress, and
    result blocks so the call to action leads and survives bold-stripping, across
    the agent-mediated UX guide, the bundled `/quality` skill, and durable skill
    specs (`Done`).
- [0122 — Finding-level candidate actions](0122-finding-candidate-actions.md) -
  types the finding `actions` field as non-binding, finding-local candidate
  actions (`description` + optional `rationale`), validates them in the CLI, has
  the skill record them on `gap`/`risk` findings as raw material for a future
  Advise phase, and keeps them out of the Evaluation v0 report and closeout
  (`Done`).
- [0123 — Render interactions through native affordances](0123-native-interaction-affordances.md)
  - makes progressive enhancement the default agent-mediated interaction contract:
    interactions are intents rendered through a fit-for-purpose native affordance
    when present, the numbered-option and `y`/`n` forms become the text fallback,
    and a harness-authorized mutation is not double-gated, across the
    agent-mediated UX guide, the shared `/quality` interaction contract, the
    bundled skill workflows, and durable skill specs (`Done`).
- [0124 — Constrain reference kind fields to closed kind vocabularies](0124-reference-kind-enum.md)
  - enum-constrains the `kind` member of Evaluation reference objects from typed
    sources — supported payload kinds for routine references (`*Ref` /
    `inputRefs[]`), report kinds for report references — so a misspelled or
    invented reference kind is rejected at `evaluation data set` time instead of
    persisting as a free-form string, with the constraint surfaced in the
    regenerated `evaluation-data.schema.json` and `json-conventions.md` (`Done`).
- [0125 — Model query commands](0125-model-query-commands.md) - adds the read-only
  `qualitymd model` group (`tree`/`list`/`get`) that projects a model's elements,
  canonical reference IDs, and containment; moves the reference grammar and a
  shared projection into `internal/model` (so `model` depends on neither `status`
  nor `evaluation`), folds `status`'s shape counts onto that projection, and wires
  the evaluate workflow to query in-scope IDs from the run's `model-snapshot.md`
  instead of hand-deriving them (`Done`).
- [0126 — Bulk data set](0126-bulk-data-set.md) - replaces
  `qualitymd evaluation data set`'s single-object stdin contract with a
  non-empty JSON array batch, validates and writes all-or-nothing, rejects
  duplicate derived paths, emits a batch receipt, and updates the evaluate
  workflow to persist routine payloads with one whole-batch dry-run and write
  (`Done`).
- [0127 — Introspection-first CLI workflow conventions](0127-introspection-first-cli-reference.md)
  - renames/refocuses the bundled skill's former `cli-quick-reference.md` as
    `cli-workflow-conventions.md`, removes embedded command/flag listings, keeps
    non-introspectable workflow conventions, and routes command, flag, and
    payload discovery to CLI introspection (`Done`).
- [0128 — Agent-mediated skill alignment](0128-agent-mediated-skill-alignment.md)
  - closes the remaining `/quality` skill alignment gaps against the
    agent-mediated UX guide: setup frame ordering, feedback-log disclosure,
    recommendation follow-up framing, and read-only orientation output (`Done`).
- [0129 — Evaluation orchestration overhaul](0129-evaluation-orchestration-overhaul.md)
  - removes the `/quality evaluate` rigor dial, makes exhaustive in-scope
    Requirement coverage mandatory, defaults independent collection and QC work
    to subagent fan-out when available, and requires an always-on two-pronged QC
    phase before roll-up (`Done`).
- [0130 — Self-contained per-kind data schema](0130-self-contained-data-schema.md)
  - makes `qualitymd evaluation data schema <kind>` emit a self-contained schema
    with required fields and enum value sets legible at the document root, keeps
    the no-argument full-surface schema intact, reuses the shared schema
    presenter, and aligns the `/quality` skill to treat schema as the constraint
    source and examples as concrete instances (`Done`).
- [0131 — Area findings in evaluation reports](0131-area-findings.md) - adds
  traceable Area Findings on `AreaAnalysisResult.findings`, validates same-Area
  Factor relationships, and renders Area Findings in Area reports plus matching
  Findings in Factor reports while leaving recommendations and global rankings
  deferred (`Done`).
- [0132 — Remove info finding severity](0132-remove-info-finding-severity.md) -
  removes `info` from the Evaluation finding severity vocabulary so active
  Requirement and Area Finding severity values are `critical`, `high`, `medium`,
  and `low`, with informational observations represented by finding `type: note`
  (`Done`).
- [0133 — Richer evaluation data examples](0133-richer-evaluation-data-examples.md) -
  makes `qualitymd evaluation data example` payloads fuller representative
  artifacts with typed Area, Factor, Requirement, and Rating Level references
  (`Done`).
- [0134 — Model-relative workspace paths](0134-model-relative-workspace-paths.md)
  - makes workspace config and artifact paths resolve relative to the selected
    `QUALITY.md`, with the Git root retained as containment and `--model`
    available for nested evaluation history/latest commands (`Done`).
- [0135 — Structured Finding Core](0135-structured-finding-core.md)
  - aligns Requirement and Area Findings around statement, condition, criteria,
    cause, effect, and evidence; rejects legacy finding summary/description
    fields; and renders Requirement, Area, and Factor Findings through one
    report structure (`Done`).
- [0136 — Candidate Actions Payload](0136-candidate-actions-payload.md)
  - renames finding-local `actions` to `candidateActions`, adds local candidate
    action IDs, and keeps candidate actions as non-rendered Requirement Finding
    raw material for later Advice (`Done`).
- [0137 — Run Report Entrypoint](0137-run-report-entrypoint.md) - makes
  `report.md` the Evaluation run entrypoint, moves the root Area report to
  `root-area.md`, and adds scoped headline report refs (`Done`).
- [0138 — Lightweight Authoring Checkpoint](0138-lightweight-authoring-checkpoint.md) -
  adds a conversational direct-authoring checkpoint for `QUALITY.md` edits
  (`Done`).
- [0139 — Real Review Gates](0139-real-review-gates.md) - makes feedback
  invitations real gates in the shared UX guide and `/quality` direct authoring
  (`Done`).
- [0140 — Casual Review Gate Wording](0140-casual-review-gate-wording.md) -
  makes direct `QUALITY.md` review gates state the planned change and value prop
  in simple conversational prose (`Done`).
- [0141 — Area-local Factor References](0141-area-local-factor-references.md) -
  makes Requirement Factor references resolve only within their declaring Area
  and warns on same-Area duplicate Factor names (`Done`).
- [0142 — Requirement Findings Only](0142-requirement-findings-only.md) - makes
  Requirement Findings the only Evaluation finding layer, removes Area Findings,
  and requires rated results to be finding-backed and driver-backed (`Done`).
- [0143 — Public Review and Improve Workflows](0143-public-review-improve-workflows.md)
  - makes `/quality review` and `/quality improve` public, focus-routed workflow
    stubs and simplifies README framing around evaluate, review, and improve
    (`Done`).
- [0144 — Pointed Review Gates](0144-pointed-review-gates.md) - makes review
  gates state inferred purpose and ask for reaction to the consequential
  assumption instead of defaulting to a generic adjustment prompt (`Done`).
- [0145 — Scannable Skill Output](0145-scannable-skill-output.md) - makes
  `/quality` runtime guidance and shared UX guidance use labeled,
  five-second-scan output shapes for multi-fact review gates, summaries, and
  closeouts (`Done`).
- [0146 — Changelog Directory](0146-changelog-directory.md) - renames the
  model-change history directory to `.quality/changelog/` with timestamped entry
  files and keeps `.quality/logs/` as a flat workflow-log directory (`Done`).
- [0147 — Report Descendant Terms](0147-report-descendant-terms.md) - aligns
  generated report labels around Child Areas for immediate Area descendants and
  Sub-Factors for immediate Factor descendants (`Done`).
- [0148 — Finding Basis](0148-finding-basis.md) - renames Finding Core `cause`
  to `basis` across evaluation records, generated reports, runtime skill
  guidance, durable specs, examples, and tests (`Done`).
- [0149 — Scope-driven evaluation runs](0149-scope-driven-evaluation-runs.md) -
  captures requested/planned scope in CLI-owned `RunManifest` data, replaces
  `--narrowing` with `--area`/`--factor`, and renders `report.md` as the scoped
  Area report without positional headline selection (`Done`).
- [0150 — Evaluation Advice](0150-evaluation-advice.md) - makes Advice required
  through Finding ranking, Recommendations, Recommendation ranking, coverage
  accounting, and generated recommendation reports (`Done`).
- [0151 — Evaluation Report CTA](0151-evaluation-report-cta.md) - makes
  `/quality evaluate` closeouts point first to the human `report.md` result and
  describe `recommendations.md` as the action-planning report (`Done`).
- [0152 — Release Reliability](0152-release-reliability.md) - hardens release
  preflight, draft-first publication, Homebrew/npm split jobs, verification, and
  repair behavior after the v0.25.0 release failures (`Done`).
- [0153 — Evaluate ranking and umbrella-factor authoring clarity](0153-evaluate-ranking-and-umbrella-clarity.md)
  - removes two `evaluate` authoring ambiguities a feedback log surfaced: the
    finding ranking now leads with completeness (tiers express priority, no
    finding dropped) in skill guidance, the example, and the protocol/payload-kinds
    rules, and the single `FactorAnalysisResult` example becomes the worked
    umbrella case (empty `localAnalysis`, child-Factor roll-up in
    `localAndDescendantAnalysis`) earned by a new Factor Traversal rule (`Done`).
- [0154 — Ranked Findings Reports](0154-ranked-findings-reports.md) - makes
  ranked findings easier to scan, navigate, and inspect across generated
  Evaluation reports (`Done`).
- [0155 — Recommendation Result Shape](0155-recommendation-result-shape.md) -
  simplifies `RecommendationResult` fields and makes recommendation reports
  render from persisted Advice data only (`Done`).
- [0156 — Report Gallery](0156-report-gallery.md) - adds a generated
  software-service report gallery with a sample `QUALITY.md`, synthetic
  Evaluation data, browsable reports, and regeneration checks (`Done`).
- [0157 — Report Markdown Authoring](0157-report-markdown-authoring.md) -
  centralizes generated Evaluation report Markdown authoring with shared table,
  link, code-span, and empty-cell helpers while preserving deterministic Go
  rendering (`Done`).
- [0158 — OKF-compatible report headers](0158-okf-report-headers.md) - adds
  pointer-only OKF-compatible frontmatter, visible run/header navigation,
  subject-shaped report types, and `RunManifest.createdAt`-backed freshness to
  generated Evaluation reports (`Done`).
- [0159 — Report source-data frontmatter](0159-report-source-data-frontmatter.md) -
  makes generated report frontmatter `data` the report-local source-data
  manifest and removes duplicate body `Data` columns from generated Evaluation
  reports (`Done`).
- [0160 — Report Body Rating Drivers](0160-report-body-rating-drivers.md) - keeps
  rating drivers in structured Evaluation data while removing standalone
  Markdown driver sections from generated Evaluation reports (`Done`).
- [0161 — Area / Factor Breakdown Reports](0161-area-factor-breakdown.md) -
  replaces flat generated subject-report tables and split Area structure tables
  with compact Area / Factor breakdowns in generated Evaluation reports
  (`Done`).
- [0162 — Report Source Data Section](0162-report-source-data-section.md) - moved
  generated report source-data pointers from YAML frontmatter into visible
  bottom `Source Data` sections (`Done`).
- [0163 — Report Artifact IDs](0163-report-artifact-ids.md) - defines
  handoff-ready Evaluation report artifact IDs for runs, recommendations, and
  ranked findings (`Done`).
- [0164 — Agent Instruction Init Pointer](0164-agent-instruction-init-pointer.md) -
  makes `qualitymd init` surface the new `QUALITY.md` to agent instruction files
  while preserving setup's narrow mutation boundary (`Done`).
- [0165 — Run IDs and Artifact Numbering](0165-run-id-artifact-numbering.md) -
  replaces local `QEVAL`/`QREC`/`QFIND` artifact IDs with a globally-unique run
  ID, per-run recommendation numbers, and finding references by selector
  (`Done`).
- [0166 — Setup Factor Proposal Checkpoint](0166-setup-factor-proposal-checkpoint.md) -
  teaches and applies factor desiderata during setup through a draft factor
  proposal checkpoint before final review (`Done`).
- [0167 — Report Frontmatter H1 Titles](0167-report-frontmatter-h1-titles.md) -
  aligns generated report frontmatter `title` values with visible H1 document
  titles while keeping `type` as the report artifact taxonomy (`Done`).
- [0168 — Quality Evaluation Title](0168-quality-evaluation-title.md) - renames
  generated run report titles to `Quality Evaluation - <Area>` and appends
  factor filters in parentheses (`Done`).
- [0169 — Run Report Opening](0169-run-report-opening.md) - reshapes generated
  `report.md` around Summary, Key Details, Contents, and non-judgmental
  frontmatter metadata while removing visible limits for now (`Done`).
- [0170 — Report Visual Markers](0170-report-visual-markers.md) - refines
  generated report visual markers for recommendation impact and Area / Factor
  breakdown rows (`Done`).
- [0171 — Run Report Simplification](0171-run-report-simplification.md) -
  simplifies generated run-level Evaluation reports around Summary, Key
  Details, Model Evaluation, ranked findings, ranked recommendations, and
  Primary Source Data (`Done`).
- [0172 — Workspace Status Contract](0172-workspace-status-contract.md) - aligns
  `qualitymd status` around workspace status, adds status JSON v2 workspace
  metadata, and preserves project wording for modeled value (`Done`).
- [0173 — Evaluation Enum Catalogs](0173-evaluation-enum-catalogs.md) -
  centralizes fixed Evaluation enum values, labels, markers, schema validation,
  and generated report ordering (`Done`).
- [0174 — Report Local Keys and Navigation](0174-report-local-keys.md) -
  replaces generated report contents and legend sections with compact navigation
  and local indicator keys (`Done`).
- [0175 — Report Contents Sections](0175-report-contents-sections.md) -
  restores generated report Contents sections for multi-section report artifacts
  and removes compact `Jump to:` lines (`Done`).
- [0176 — Recommendation IDs and Numbers](0176-recommendation-ids-and-numbers.md) -
  splits opaque recommendation identity from user-facing recommendation numbers
  so Recommendation #1 means the first ranked recommendation (`Done`).
- [0177 — Report Full List Links](0177-report-full-list-links.md) - makes
  run-report links to full findings and recommendations lists more scannable and
  shows total ranked counts (`Done`).
- [0178 — Recommendation Number Columns](0178-recommendation-number-columns.md) -
  removes duplicate `Rank`/`#` columns from recommendation list tables so `#` is
  the single visible recommendation number (`Done`).
- [0179 — Enum Catalog Metadata](0179-enum-catalog-metadata.md) - added
  type-level labels and descriptions plus value descriptions to Evaluation enum
  catalogs, and generated report local keys now render catalog labels (`Done`).
- [0180 — Finding Taxonomy and Report Details](0180-finding-taxonomy-report-details.md) -
  removes `unknown` as a Finding type, distinguishes gap/risk markers, and
  improves generated run report key details (`Done`).
- [0181 — Evaluation Identity Manifest](0181-evaluation-identity-manifest.md) -
  makes Evaluation identity primary with `evaluationId`, nests local run
  metadata in the manifest, and simplifies run report frontmatter (`Done`).
- [0182 — Finding Summary Display Order](0182-finding-summary-display-order.md) -
  renames the run report Finding Breakdown to Finding Summary, shows all Finding
  types in concern-first order, and labels sparse severity counts (`Done`).
- [0183 — Glossary and Report Links](0183-glossary-report-links.md) - adds a
  workspace-root glossary, replaces generated report legends with `Evaluation
  links:`, and seeds glossary terms and vocabularies (`Done`).
