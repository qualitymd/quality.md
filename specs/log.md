# Specs Update Log

## 2026-06-19

- **Revision**: Updated the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) for
  project-owned Agent Skills metadata: `metadata.version`,
  `metadata.requires-qualitymd-cli`, matching `compatibility` prose, and the
  released-install prerequisite range.

- **Revision**: Added [`qualitymd version`](cli/version.md) and
  [`qualitymd upgrade`](cli/upgrade.md) command specs, updated the
  [CLI overview](cli.md) and [command index](cli/index.md), and specified
  structured local version metadata plus explicit, opt-in upgrade checks.

- **Revision**: Updated [`evaluation build-report`](cli/evaluation-build-report.md),
  [Evaluation records](evaluation-records.md), and the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) for generated
  `report-summary.md` alongside `report.md` and `report.json`. Added a
  representative summary artifact to the
  [skill examples](skills/quality-skill/examples/index.md).

- **Revision**: Updated the root [`SPECIFICATION.md`](../SPECIFICATION.md),
  [`lint`](cli/lint.md), [`init`](cli/init.md), [`status`](cli/status.md),
  [`evaluation build-report`](cli/evaluation-build-report.md), and
  [`/quality` authoring guide spec](skills/quality-skill/authoring-guide.md) for
  required display titles on Models, Targets, Factors, and Rating Levels.

- **Revision**: Updated the durable
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) to align CLI
  prerequisite checks with the new versioning policy: released skill installs
  use the declared `qualitymd` SemVer range, while local development builds can
  still be accepted by command-surface probing.

- **Revision**: Reframed the durable
  [`/quality` wizard](skills/quality-skill/quality-skill.md#wizard) contract as a
  read-only quality wayfinder: it now probes setup/model/history state, classifies
  lifecycle readiness, recommends one next workflow, and offers concrete
  alternatives without producing ratings or editing files.

- **Revision**: Updated the durable
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) to settle the
  skill's bundled runtime resources and require deep-evaluation subagent prompts
  to carry scope, relevant requirements, secret-handling, source-as-data, and
  structured-findings-only instructions.

- **Revision**: Normalized active specs to the RFC 8174 / BCP 14 convention and
  reduced uppercase requirement keywords where the text was advisory,
  illustrative, or presentation guidance rather than a conformance switch.
  Updated the root format spec, CLI specs, evaluation record contract, and
  `/quality` skill specs to reserve BCP 14 wording for testable requirements.

## 2026-06-18

- **Creation**: Added the [QUALITY.md authoring guide](skills/quality-skill/authoring-guide.md) sub-spec — the 1:1 contract for
  the skill's [`quality-md-guide.md`](../skills/quality/resources/quality-md-guide.md)
  resource. States its purpose (a single comprehensive guide to understanding and
  working with `QUALITY.md` files), resolves the embed-vs-drift tension by scoping
  the guide as a deliberately self-contained reference resource that restates the
  format (an explicit exception to the prompt's no-drifting-copy rule) paid for by
  a conformance duty to [`SPECIFICATION.md`](../SPECIFICATION.md), and fixes its
  structural conventions (single-level concepts, reference + how-to per chapter,
  directive/job form, authoring order). Promoted the authoring guide out of the
  skill spec's **Deferred** resources note and added the one-line conformance note
  and purpose statement to the guide itself.

- **Revision**: Updated the
  [`lint`](cli/lint.md) spec for mandatory factor references: added
  `missing-factor-reference`, made `unknown-factor` apply to factor references
  generally, and aligned `empty-factor` terminology with requirement
  references.

- **Revision**: Updated the durable
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) to allow the
  installable skill to keep `SKILL.md` as the always-loaded router while moving
  mode-specific procedures into `skills/quality/modes/` and supporting docs into
  `skills/quality/resources/`.

- **Revision**: Removed the public bundled-model workflow. Deleted the
  `qualitymd models` command spec, removed it from the CLI command catalog,
  made `evaluation create-run` subject-only with no `--altitude` flag, and
  updated the evaluation record contract to treat `model` altitude as historical
  only. Synced the durable `/quality` skill spec around subject evaluation plus
  guide-backed `QUALITY.md` authoring.

- **Rationale capture**: Made durable the motivating learnings behind the
  in-review evaluation changes (0012–0024), so the *why* survives when the change
  records are archived. Added inline rationale to
  [Evaluation records](evaluation-records.md) (standalone-spec reason; CLI-owned
  scan-derived numbering and collision-as-error from the real run-number
  collision; CLI `model.md` snapshot ownership; duplicate-assessment correction
  hazard; assessment-vs-recommendation superseding asymmetry; explicit-intent
  superseding; no-new-schema-field constraint for evidence/route hints),
  [`build-report`](cli/evaluation-build-report.md) (renderer secret /
  prompt-injection trust boundary; summary-first shape provenance; the bounded
  `design.md`/`plan.md` summary-source contract; dotted-path and idempotency
  rationale), [`add-record`](cli/evaluation-add-record.md) (input/decode rules,
  CLI-owned-field rejection, exit-code mapping, run-snapshot rating validation,
  collision retry, deterministic generated rendering, subcommand-per-kind),
  [`create-run`](cli/evaluation-create-run.md) (partial-artifact prevalidation
  reason and prevalidate-over-rollback), and the
  [`/quality` skill spec](skills/quality-skill/quality-skill.md) (re-run vs.
  re-read re-check rationale; superseding correction rationale).

- **Revision**: Added assessment superseding for evaluation runs. Corrected
  assessment records may now carry `supersedes`; status reports dangling,
  cross-requirement, and stale-analysis superseding references; reports
  distinguish active from superseded assessments.

- **Revision**: Tightened `evaluation create-run --subject` validation. Invalid
  subject paths now fail before creating the evaluation directory or a numbered
  run folder, and subject paths must resolve to files rather than directories.

- **Revision**: Added recommendation superseding for evaluation runs. Corrected
  recommendation records may now carry `supersedes`, status reports dangling
  superseding references, and reports distinguish active from superseded advice
  when choosing Next Action.

- **Revision**: Added optional planned coverage metadata for evaluation runs.
  `qualitymd evaluation set-planned-coverage` writes run-root
  `planned-coverage.json`, and `show-status` now reports missing planned
  assessments, missing planned analyses, and unexpected records when that
  artifact exists.

- **Revision**: Tightened evaluation reportability for duplicate assessment
  records. Runs with more than one assessment for the same ordered `targetPath`
  and `requirement` now report a `duplicate-assessment` gap in `show-status`,
  and `build-report` refuses to render them.

- **Revision**: Tightened evaluation reportability: `show-status` and
  `build-report` now require exactly one in-scope root analysis record with an
  empty `targetPath`, so a child target cannot silently become the report
  headline.

- **Revision**: Updated the evaluation report contract for change
  [0018](../changes/archive/0018-evaluation-report-ux.md). Generated reports are now
  summary-first for human readers and expose non-null scope, target summaries,
  evidence basis, limitations, empty recommendation arrays, explicit rating
  objects, and structural grouping-target states in `report.json`.

- **Refinement**: Clarified that equivalent summary limitations are
  deterministically deduplicated across recorded context and rationale-derived
  constraints while preserving the first displayed wording.

- **Refinement**: Clarified that derived limitation summaries preserve
  locator-like text such as dotted file paths.

- **Creation**: Added the durable
  [Evaluation records](evaluation-records.md) contract and the evaluation CLI
  sub-specs for
  [`create-run`](cli/evaluation-create-run.md),
  [`add-record`](cli/evaluation-add-record.md),
  [`show-status`](cli/evaluation-show-status.md), and
  [`build-report`](cli/evaluation-build-report.md). Updated the CLI overview,
  command index, and schema type description for the shared runtime record
  contract.

- **Revision**: Synced the durable
  [`/quality` skill spec](skills/quality-skill/quality-skill.md), installable
  skill prompt, skill guide, and reference recommendation examples with the new
  evaluation CLI surface. The skill now delegates run scaffolding, record writes,
  renderability checks, and report rendering to `qualitymd evaluation ...`, while
  retaining judgment and the stricter evidence/re-check/effort rules.

- **CLI output polish**: Updated the CLI baseline to apply the `NO_COLOR` and
  non-terminal plain-output gate across stdout and stderr human surfaces. Updated
  `qualitymd lint` to require deterministic next actions for invalid results and
  a matching stderr footer, and updated `qualitymd models list` to allow
  terminal-only styling over the unchanged plain table.

## 2026-06-17

- **Revision**: Synced the [`/quality` skill spec](skills/quality-skill/quality-skill.md)
  with its implementation path and raw runtime artifact contract. The spec now
  points to [`skills/quality/SKILL.md`](../skills/quality/SKILL.md), uses the
  finalized trigger description, defines setup and model-altitude evaluation via
  `qualitymd models view quality-meta-model --source`, adds `.quality/config.yaml`
  `evaluationDir`, requires JSON assessment/analysis records plus `report.json`,
  and removes runtime artifact concept types from [`schema.md`](schema.md).

- **Creation**: Added the [`qualitymd models`](cli/models.md) command spec and
  listed it in the [CLI spec](cli.md). The command exposes bundled models through
  `models list` and `models view <name>`, with JSON forms for agents/tools and a
  `--source` rewrite for model-altitude evaluation with the bundled
  `quality-meta-model`.

- **Revision**: Strengthened the [`0001` evaluation report](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  against the elements ISO/IEC 25040:2024 §5.6.3.2 lists for an evaluation
  report, while staying within the format spec's "presents at least" Report
  contract. Added an **Evaluated** provenance line (source commit, model
  revision, evaluator, date, assessment inputs) so the verdict is reproducible; a
  **Limitations** section (effort ceiling, point-in-time secret scan, single-test
  confidence) qualifying the ratings without changing them, kept distinct from
  Scope exclusions and *not assessed* outcomes; and sharpened the committed-credential
  rationale to trace measure → applied criterion against the requirement's own
  `ratings` overrides. Declared the new `9f2c1ab` commit as a shared fictional
  locator in the [examples index](skills/quality-skill/examples/index.md).

- **Revision**: Extended the [`0001` example bundle](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  with a **two-level nested target** to exercise multi-level roll-up: a
  **Webhooks** child target (sibling of Ledger, source `./webhooks`, a signing
  requirement under a refined **Security** factor) with a **Delivery**
  sub-target (`./webhooks/delivery`, retry + redelivery-suppression under a
  refined **Reliability** factor) in
  [`model.md`](skills/quality-skill/examples/0001-subject-quality-eval/model.md).
  The Delivery deduplication requirement is *rated* **Minimum**, so the Webhooks
  aggregate (**Minimum**) falls below its own local rating (**Target**) — the
  bundle's first intermediate aggregate that differs from its local — and the
  root's counterfactual now layers: rotating the committed credential lifts the
  root only to **Minimum** (the webhook gap then binds), reaching **Target**
  once that closes too. Updated the
  [report](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  rationale/scope/advice, revised recommendation
  [001](skills/quality-skill/examples/0001-subject-quality-eval/recommendations/001-rotate-committed-gateway-key.md)'s
  done-criterion to match, added recommendation
  [003](skills/quality-skill/examples/0001-subject-quality-eval/recommendations/003-bound-webhook-dedup-window.md),
  and refreshed the [examples index](skills/quality-skill/examples/index.md).

- **Revision**: Polished the [`0001` example bundle](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  for readability and to track the format spec's new **target display fields**:
  gave the **Ledger** child target a `title` and `description` (the report
  already displayed "Ledger", which the bare map key did not provide); added an
  "At a glance" condensed model tree above the verbatim YAML in
  [`model.md`](skills/quality-skill/examples/0001-subject-quality-eval/model.md),
  mirroring the specification's sample-report tree; centralized the facts shared
  across the bundle (fictional subject/locators, the suggested four-level scale,
  and that `model.md` is the evaluated file rather than a runtime output) into
  the [examples index](skills/quality-skill/examples/index.md), trimming the
  repeated per-file admonitions to a pointer; stated an explicit
  model → report → recommendations reading order in the index; and added a
  cross-target **secondary factor** — the Ledger's double-entry requirement now
  tags the root **Reliability** factor, so the example exercises a factor lens
  ranging wider than the local rating (the spec's Analyze/Report secondary-factor
  notes).

- **Revision**: Updated [`lint`](cli/lint.md) for target display fields:
  `misplaced-root-key` now documents only `ratingScale` on a Target, matching the
  format spec's Model/Target distinction after target `title` and `description`
  became valid Target properties.

- **Restructure + Creation**: Gave the [`/quality`](skills/quality-skill/quality-skill.md)
  skill its own OKF folder — moved the spec into
  [`skills/quality-skill/`](skills/quality-skill/index.md) (as a named concept
  beside a new folder `index.md`, keeping `index.md` a frontmatter-free listing
  per OKF) and rewrote its relative links for the new depth. Added an
  [`examples/`](skills/quality-skill/examples/index.md) subfolder with the first
  worked reference instance of the skill's [Reporting](skills/quality-skill/quality-skill.md#reporting)
  contract: a whole-model **subject** evaluation of a fictional "Sparrow
  Payments" service held at **Unacceptable** by a committed live gateway
  credential — a [report](skills/quality-skill/examples/0001-subject-quality-eval/report.md)
  plus two standalone
  [recommendations](skills/quality-skill/examples/0001-subject-quality-eval/recommendations/001-rotate-committed-gateway-key.md)
  with done-criteria, plus the
  [model evaluated](skills/quality-skill/examples/0001-subject-quality-eval/model.md)
  reproduced alongside (an `Example Model` concept embedding the fictional
  `QUALITY.md` — root **Security** with **Secrets handling**/**Access control**
  sub-factors and **Reliability**, a **Ledger** child target with
  **Correctness**, and a `ratings` override on the secrets requirement the
  report's rationale and done-criterion lean on) so every finding traces to a
  declared requirement and `source`. The example exercises `file:line`
  evidence, the
  secret-by-reference rule (credential cited by location and type, value
  withheld, rotation recommended), a prompt-injection comment recorded as a
  finding and treated as data, and a *not assessed* requirement excluded from
  the roll-up but noted. Registered `Evaluation Report` and `Recommendation`
  concept types in [`schema.md`](schema.md) and linked the example from the
  skill's Reporting section.

- **Revision**: Reconciled the structural schema source of truth with
  [`SPECIFICATION.md`](../SPECIFICATION.md) and [`lint`](cli/lint.md): the public
  format spec now notes the typed schema declaration consumed by lint and the
  drift test that compares it to the YAML schema snippets; the model `title`
  snippet now matches its recommended status; and the lint sub-spec records that
  structural validation derives from the single schema declaration rather than a
  second valid-key list.
- **Revision**: Added the [CLI spec](cli.md)'s **Agent accessibility** section:
  a non-opt-in baseline for in-scope commands (non-interactivity,
  stdout-is-payload/stderr-is-diagnostics, determinism, plain non-TTY output,
  and categorized exit codes) plus opt-in capabilities (`--json`, `nextActions`,
  and deferred quiet/verbosity). Documented exit codes `0` success, `1`
  ran-but-found-problems, `2` usage error, and `70` internal/could-not-complete
  error; revised the `--json` convention to a SHOULD-by-default with a
  verbatim-artifact carve-out and result receipts for side-effecting commands;
  and updated [`init`](cli/init.md) to specify its JSON receipt, `--json -`
  usage error, and JSON overwrite-refusal object.
- **Revision**: Reframed the [`/quality`](skills/quality-skill/quality-skill.md) skill's
  relationship to the format spec's evaluation from **deference** to
  **conformance**: the skill now *owns and specifies* its evaluation process
  (this spec, its prompt, and the CLI together) rather than pointing at
  [`SPECIFICATION.md`](../SPECIFICATION.md) to "run." Added a **Conformance to
  the format spec** subsection stating that every evaluation MUST conform to the
  spec's Evaluation contract (assessment → finding → rating, *not assessed* over
  guessing, inferred weighted roll-up, required report contents) while the skill
  remains one *implementation* free to specify its own workflow, ordering,
  heuristics, effort levels, and artifacts — with the spec as the **conformance
  target** that governs on divergence. Scoped the existing "don't embed the
  format" rule to the *format/schema rules and rating vocabulary* (still grounded
  at runtime from `qualitymd spec`), explicitly excluding the *evaluation
  process*, which the skill carries.
- **Revision**: Tightened [`lint --fix`](cli/lint.md)'s in-place write contract
  to avoid ambiguous symbolic-link replacement: repair should refuse a linted
  symlink path until symlink write semantics are specified.
- **Revision**: Scoped in [`lint --fix`](cli/lint.md) for deterministic in-place
  repair of fixable findings. The durable lint spec now defines the `--fix`
  flag, transactional per-file repair behavior, preservation requirements,
  post-repair linting, exit-status semantics, and JSON repair reporting through
  `summary.fixed` and `repairs`, while leaving patch/full-file repair output
  modes deferred.
- **Revision**: Closed the lint-specific open items in
  [`lint`](cli/lint.md): added scope/deferred boundaries, specified that this
  phase has no command-specific flags beyond inherited `--json`, defined the
  `--json` finding document and per-finding fields, made `modelPath` the stable
  location contract with optional source positions, set the minimum human-output
  content, and fixed deterministic finding ordering plus downstream-rule
  blocking behavior; the sub-spec now carries a neutral inheritance note instead
  of a draft placeholder. Automatic repair, suppression, rule selection, severity
  overrides, and a lint-emitted rule catalog remain deferred.
- **Revision**: Added fixability criteria to the [`lint`](cli/lint.md) rule
  authoring contract — single mechanical edit, no presumed intent,
  content-preserving, strictly improving, idempotent, local/explainable, and
  stable formatting — and expanded the rule inventory with `Fixable` and
  `Fixable rationale` columns. The initial inventory marks only `empty-property`
  as fixable; placeholder insertion remains scaffold behavior owned by
  [`init`](cli/init.md), not a lint repair.
- **Revision**: Fleshed out the stub [`/quality`](skills/quality-skill/quality-skill.md) skill
  spec into a draft, working through the unblocked TODOs (inspired by the
  basecamp and shadcn/improve `SKILL.md` patterns). Promoted the use-case sketch
  to normative requirements — resolving the open question as *evaluate is
  read-only; improve recommends and applies only on explicit confirmation* — and
  added: a **Boundaries and hard rules** section (judgment-vs-CLI division,
  evaluated content is untrusted data not instructions, never reproduce secrets,
  scoped-result-is-not-a-whole-model-verdict, determinism); **Invocation**
  (skill frontmatter/metadata kept in sync with the format via `qualitymd spec`
  rather than hard-coded, plus five resolved arguments — mode
  `evaluate`/`improve`/`setup`/`wizard`, subject-vs-model altitude, target file,
  target/factor scope that composes with either altitude and resolves bare names
  against the grounded model, and effort — where a bare `/quality` runs a
  read-only `wizard` mode that inspects model state and suggests concrete next
  invocations);
  **Driving the CLI** (`init`/`lint`/`spec` for every mechanical step, introspect
  rather than hard-code the surface); an **Evaluation workflow** that wraps the
  format spec's five Evaluation phases with mechanical read → lint → ground →
  evaluate → report steps; **Grounding judgment** (rate against declared
  criteria, evidence per rating, *not assessed* over guessing, inferred weighted
  roll-up); **Effort levels** (`quick`/`standard`/`deep` coverage); and
  **Reporting** (the Evaluation Report, scoped, human or `--json`), plus an
  illustrative (non-normative) commented-examples block sketching invocation
  patterns in the style of the shadcn/improve README. Recording
  assessments through the CLI and bundled `resources/` remain deferred in step
  with the CLI's deferred record/gate surface.
- **Revision**: Firmed up [`SPECIFICATION.md`](../SPECIFICATION.md) to fully
  ground the [`lint`](cli/lint.md) rules — made the per-property YAML shapes
  normative, defined a null/empty *required* value as absent, and stated the
  root-only `title`/`ratingScale` constraint as an explicit MUST NOT — then
  synced `lint`: broadened `invalid-frontmatter` to structural-shape conformance,
  added an "empty is absent" note, renamed `unresolved-factor` to `unknown-factor`
  (verb-consistent with `unknown-rating-key`), added `RECOMMENDED` to the RFC 2119
  keywords, corrected the rule `description` guidance and the
  `empty-model`/`empty-factor`/`empty-target` entries, and recorded the body
  heading under "Not checked".
- **Creation**: Added the [`skills/`](skills/index.md) subfolder with a stub
  [`/quality`](skills/quality-skill/quality-skill.md) functional spec — the companion
  evaluation skill that carries judgment against a `QUALITY.md` and records
  results through the deterministic CLI. The concrete workflow, inputs, and
  recorded-assessment shape remain deferred (in step with the CLI's deferred
  record/log surface).
- **Revision**: Fleshed out the [`lint`](cli/lint.md) sub-spec with a "Rule
  scope" section (four inclusion criteria — conformance-grounded, deterministic,
  format-not-goodness, self-contained — plus the error/warning/info severity
  mapping) and a "Rule naming" section (kebab-case noun-phrase ids in two shapes,
  using QUALITY.md vocabulary, treated as a stable public identifier) and a
  "Rule authoring" section with criteria for the static rule `description`
  (generic, present-tense, severity-signalling) and the per-instance finding
  `message` (names the element and location, contrasts found-vs-required,
  deterministic, tone-matched, actionable, one-per-instance), with the principle
  that description and message must not duplicate each other.
- **Revision**: Added the concrete "Rules" section to [`lint`](cli/lint.md) —
  eleven `error` rules (one per format-spec MUST: `invalid-frontmatter`,
  `missing-rating-scale`, `too-few-levels`, `missing-level-name`,
  `duplicate-level`, `missing-criterion`, `empty-model`, `misplaced-root-key`,
  `invalid-assessment`, `unknown-factor`, `unknown-rating-key`) and six
  `warning` rules (one per mechanically checkable SHOULD/RECOMMENDED:
  `missing-title`, `missing-level-description`, `missing-factor-description`,
  `empty-factor`, `empty-target`, `empty-property`), each citing the clause it
  enforces, plus a "Not checked" note for rating-level order (not mechanically
  verifiable) and the as-yet-unused `info` severity. The finding schema and
  lint-specific flags remain deferred.
- **Revision**: Replaced the placeholder [`init`](cli/init.md) command sub-spec
  with the durable requirements for scaffold contents, output target handling,
  stdout piping, overwrite protection, `--force`, and stderr reporting.
- **Revision**: Added a "Technical requirements" section to the
  [CLI spec](cli.md) requiring that every functional requirement be satisfiable
  through the idiomatic capabilities of the chosen stack (Go + Cobra + Charm
  Fang + Lip Gloss), rather than working against the grain of those libraries.
- **Revision**: Added a "Conventions" section to the [CLI spec](cli.md)
  establishing `--json` as the spelling for machine-readable output wherever a
  command offers one (not a requirement that every command do so), with criteria
  for when a command should offer `--json` and worked examples across the current
  commands. Updated [`lint`](cli/lint.md) to reference `--json` in place of its
  earlier `--format json`.
- **Convention**: Added a "Suggested next actions" convention to the
  [CLI spec](cli.md): commands may close with advisory, deterministic next-action
  suggestions — a stderr footer in human output, an in-band `nextActions` array
  under `--json` — that never affect behavior or the exit code.

## 2026-06-16

- **Convention**: Added a bundle-root [`schema.md`](schema.md) (`type: Schema`)
  registering the bundle's concept types (`Schema`, `Functional Specification`)
  in frontmatter, and listed it from the [index](index.md). Retyped the CLI
  spec and command sub-specs from `Specification` / `Command Specification` to a
  single `Functional Specification` type.

- **Initialization**: Created the `specs/` OKF bundle with a root
  [CLI spec](cli.md) capturing the high-level CLI requirements (design
  properties, global flags, output formats, exit codes, agent accessibility).
- **Creation**: Added placeholder command sub-specs for
  [`init`](cli/init.md), [`lint`](cli/lint.md), and [`spec`](cli/spec.md), plus
  the [`cli/` index](cli/index.md).
- **Revision**: Reduced the [CLI spec](cli.md) to a stub. No requirements had
  actually been requested, so the speculative design properties, global flags,
  output formats, exit codes, and agent-accessibility requirements were stripped
  out, leaving scope, a command list, and a "to be specified" outline — matching
  the placeholder command sub-specs.
