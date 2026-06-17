# Specs Update Log

## 2026-06-17

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
