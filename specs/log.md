# Specs Update Log

## 2026-06-17

- **Revision**: Reframed the [`/quality`](skills/quality-skill.md) skill's
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
- **Revision**: Fleshed out the stub [`/quality`](skills/quality-skill.md) skill
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
  assessments through the CLI and bundled `references/` remain deferred in step
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
  [`/quality`](skills/quality-skill.md) functional spec — the companion
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
