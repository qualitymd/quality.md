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
