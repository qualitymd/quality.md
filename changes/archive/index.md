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
