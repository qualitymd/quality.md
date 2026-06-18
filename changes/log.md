# Changes Update Log

## 2026-06-18

- **Done**: Archived change
  [0011 — CLI human output polish](archive/0011-cli-human-output-polish.md)
  after styling `models list`, adding lint next actions to JSON and human
  output, making dev version output include a short VCS revision when available,
  syncing the durable CLI specs, and adding focused output-gate tests.

- **Implementation**: Created and advanced change
  [0011 — CLI human output polish](archive/0011-cli-human-output-polish.md) to
  `In-Progress`. The change covers terminal styling for `models list`, lint
  next actions in JSON and human output, dev version output that includes a VCS
  revision, and broader output-gate tests, with durable CLI specs to be synced
  during implementation.

## 2026-06-17

- **Done**: Archived change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md)
  after landing the skill artifact, `qualitymd models` CLI surface, durable
  specs/docs, raw runtime example bundle, and verification. Removed it from the
  open-changes index and added it to the archive index.

- **Implementation**: Implemented change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md):
  added
  `skills/quality/SKILL.md`, implemented `qualitymd models list/view` with Markdown
  and JSON output plus `--source`, moved the bundled quality meta-model under
  `internal/models`, added skill-first install/docs, synced the durable CLI and
  skill specs, re-captured the example as raw runtime artifacts with JSON
  assessment/analysis/report files, ignored default dogfood runs under
  `quality/evaluations/`, and verified the skill/CLI surfaces locally.

- **Implementation**: Advanced change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) from
  `Design` to `In-Progress`. The functional spec and design doc are settled, so
  implementation files and durable specs/docs can now be updated: the
  `skills/quality/SKILL.md` artifact, `qualitymd models` CLI surface,
  `.quality/config.yaml` behavior, raw JSON example artifacts, and related durable
  documentation.

- **Refinement**: Added a comprehensive acceptance checklist to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md), covering
  skill packaging/install, CLI prerequisite handling including dev builds,
  `qualitymd models` Markdown/JSON behavior, `.quality/config.yaml` validation,
  default dogfood-output ignoring, quick model-altitude dogfooding, JSON artifact
  parsing and shape, minimal `report.json` finding summaries, example re-capture,
  and durable spec/doc sync before **Done**. Optional installer/UI metadata such as
  `agents/openai.yaml` is explicitly deferred until the installer or target agent
  docs require it.

- **Refinement**: Settled the final `SKILL.md` description text for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md): "Use when
  a user wants setup, wizard guidance, evaluation, or improvement for quality
  management of a project/entity or one of its components/targets. Trigger for
  requests about quality factors, characteristics, attributes, criteria, Targets,
  Factors, Requirements, improving a quality factor such as
  security/reliability/usability, evaluating a subject against quality criteria, or
  evaluating/improving the QUALITY.md model itself."

- **Refinement**: Added evaluation-directory configuration to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  now reads repository-local `.quality/config.yaml` with `evaluationDir` to choose
  the parent directory for numbered evaluation runs, defaulting to
  `quality/evaluations/` when absent; the config is framed as shared qualitymd
  system config that future CLI evaluation commands should also honor. The path must
  be repository-relative, normalized, and unable to escape the repository; unknown
  keys are warned and ignored. Broader configuration ideas (default target file,
  effort, output formats, thresholds, retention, install commands) are deferred
  until they have a concrete need.

- **Refinement**: Added trigger-description requirements for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  description must now cover quality management/evaluation/improvement prompts even
  when the user does not mention `QUALITY.md` (for example, improving security
  quality), include mode trigger terms (`setup`, `wizard`, `evaluate`, `improve`)
  and generic quality vocabulary (factors, characteristics, attributes, criteria),
  and stay bounded away from generic copyediting or one-off "higher quality"
  requests. The design records initial `SKILL.md` description text that names
  Targets, Factors, Requirements, subject evaluation, model evaluation/improvement,
  and broad project/entity plus component/target quality, while leaving CLI
  implementation details to the skill body. The design now records the criteria
  used to derive that description and the durable spec sync explicitly carries those
  criteria into `quality-skill.md`'s Frontmatter and metadata section.

- **Refinement**: Added dogfooding guidance to change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The design
  now requires an In-Progress verification pass that installs the skill from the
  working tree, accepts a local development `qualitymd` binary when it exposes the
  required commands, runs a quick model-altitude evaluation against this repo's
  `QUALITY.md`, checks the generated artifact shape, and avoids committing ad hoc
  `quality/evaluations/` output unless deliberately re-captured as a durable
  example.

- **Refinement**: Resolved the remaining open questions for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). Root
  `install.md` now uses a verification-first install flow (`qualitymd --version`,
  documented CLI install/upgrade, verify again, `npx skills add qualitymd/quality.md`,
  `npx skills list`) with the exact package-manager command filled in when the first
  public CLI release channel exists. `report.json` now inlines only minimal generic
  finding summaries by reference for single-file gate/dashboard consumers, while
  full finding detail remains in `assessments/*.json`. Future
  `qualitymd outline QUALITY.md --json`, CLI evaluation record/gate commands,
  and deep-run subagent fanout are explicitly deferred rather than open design
  questions.

- **Refinement**: Generalized the structured finding shape for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). Replaced
  the sample's bespoke top-level `credentialType` with a generic finding object:
  `locator`, `observation`, open `category`, optional `severity`, supporting
  `evidence`, and optional `attributes` for domain-specific metadata. Added the
  broader JSON-shape rule that public top-level fields stay tied to the evaluation
  workflow, while factor- or requirement-specific details live under
  `attributes`.

- **Refinement**: Updated change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) so
  `qualitymd models view <name>` supports `--json`. The default output remains
  Markdown with the same terminal-rendered vs plain/verbatim split as
  `qualitymd spec`, preserving byte-for-byte `model.md` snapshots while giving
  humans a readable TTY view. `--json` emits the same source-rewritten bundled model
  as structured data (`model` plus `bodyMarkdown`) for agents and gates that should
  not have to reparse Markdown/YAML. Updated the functional spec, design doc, and
  design log wording.

- **Refinement**: Corrected the onboarding model for change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md). The skill
  is now the primary entry point, installed from this repo with
  `npx skills add qualitymd/quality.md`; the `qualitymd` CLI is a prerequisite that
  `setup` and `wizard` detect, version-check, and help install or upgrade before
  running CLI-dependent work. Updated the functional spec, design doc, parent
  change, and indexes to remove the plugin-first assumption, added root
  `install.md` to affected docs, and kept Claude plugin/marketplace packaging as a
  possible secondary channel rather than this change's primary distribution path.

- **Design**: Advanced change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0010-implement-quality-skill/design.md). Confirmed the three
  **blocking** open items at their recommended resolutions and worked out the *how*:
  the skill ships from `skills/quality/SKILL.md` as an Agent Skills artifact
  installable with `npx skills add qualitymd/quality.md`, while `setup`/`wizard`
  verify the separate `qualitymd` CLI prerequisite before doing CLI-dependent work.
  Specified the `qualitymd models` surface (`list` with `--json`; `view <name>
  [--source]` as verbatim Markdown by default and structured JSON with `--json`,
  reusing the `lint --fix` node-rewrite to re-point the meta-model's apex
  `source`), homed in a new `internal/models` package for the bundled-model
  catalog. Settled the raw, non-OKF evaluation
  artifacts — JSON `assessments/`/`analysis/` source-of-record records and a
  `report.json` rendered over them beside the human `report.md`, altitude-first
  folder naming, and `improve`'s new-folder re-evaluation. Recorded the alternatives
  (plugin-marketplace-first/CLI-installer distribution, inline vs referencing
  `report.json`, meta-model embed home) and planned the In-Progress durable sync.
  Updated the change [index](archive/0010-implement-quality-skill/index.md), bundle
  [index](index.md), and the parent concept's status.

- **Refinement**: Tightened change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md) after
  review. Reconciled the open items' conflicting lifecycle timing — they are now
  **surfaced in Draft**, with the **blocking** ones resolved before **Design** and
  the rest during **In-Progress**, all before **Done** (replacing contradictory
  "settle before Design" / "while Draft" / "before Done" wording across the spec and
  parent). Settled that **evaluation artifacts are raw runtime outputs, not OKF
  concepts**: JSON assessment/analysis records plus a `report.json`, with the worked
  example re-captured raw and its now-unused concept types retired from
  [`specs/schema.md`](../specs/schema.md). Brought the deterministic
  `qualitymd models` command into the change's **Covered** scope, since the model
  altitude depends on it and the skill drives the CLI rather than reimplementing it.
  Renamed the open-items anchor and aligned the parent's `Q1` references to the
  spec's item numbering. Updated the [log](log.md).

- **Creation**: Added change
  [0010 — Implement the /quality skill](archive/0010-implement-quality-skill.md)
  (`status: Draft`) with its
  [functional spec](archive/0010-implement-quality-skill/spec.md) to build the
  specified-but-unimplemented
  [`/quality` skill](../specs/skills/quality-skill/quality-skill.md). The spec
  **defers the behavioral contract** to the durable skill spec and states only the
  delta — package an invocable skill that conforms to it and drives the `qualitymd`
  CLI for every mechanical step — plus the open items and gaps a review of the
  skill spec surfaced for this change to settle: where the skill is packaged; where the
  **model** altitude draws its criteria (the built-in
  [meta-model](../internal/models/quality-meta-model.md) is
  neither referenced nor CLI-exposed); what `setup` does beyond `init`; how the
  default target file resolves (a CLI convention still "to be specified"); and
  whether `improve`'s post-apply re-evaluation writes a new evaluation folder.
  Records [`quality-skill.md`](../specs/skills/quality-skill/quality-skill.md) and
  user docs as affected durable artifacts, to sync once the questions resolve.
  Updated the bundle [index](index.md).

- **Completion**: Implemented and archived
  [0009 — Diagnose rating-scale soundness in the meta-model](archive/0009-rating-scale-diagnostic.md),
  adding the *rating scale and any overrides are well-formed and meaningful*
  requirement to the [meta-model](../internal/models/quality-meta-model.md)'s
  Functionality factor. The meta-model previously assessed the rating scale only
  structurally (lint's "well-shaped" check) and as one clause in a conformance
  list, despite the scale being what turns assessments into verdicts and despite
  per-requirement `ratings` overrides giving authors room to miscalibrate a
  threshold or quietly redefine a level. Synced the Functionality summary and the
  diagnostic coverage checklist, and confirmed `qualitymd lint` still reports the
  model valid. No durable specs/docs were affected — the requirement traces to the
  rating-scale semantics already in [`SPECIFICATION.md`](../SPECIFICATION.md),
  which is unchanged.

- **Creation**: Added change
  [0009 — Diagnose rating-scale soundness in the meta-model](archive/0009-rating-scale-diagnostic.md)
  (`status: Draft`) with its
  [functional spec](archive/0009-rating-scale-diagnostic/spec.md). Prompted by the
  meta-model's coverage asymmetry — six requirements for the Markdown body
  sections, but the rating scale assessed only for structural shape. The change
  adds a single Functionality requirement covering level meaning, band
  separability, floor placement against the model's needs and risks, and sound
  per-requirement criterion overrides, written to pass trivially for a model that
  inherits the suggested scale unchanged. Omits a design doc as a one-requirement
  content change; records no affected durable specs/docs.

- **Completion**: Implemented and archived
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md),
  adding `title` and `description` to Target, `description` to Model, and the
  matching durable [`SPECIFICATION.md`](../SPECIFICATION.md) prose and
  [`lint`](../specs/cli/lint.md) rule row. The structural schema now accepts
  target display fields, `misplaced-root-key` flags only nested `ratingScale`,
  and focused tests cover accepted nested target `title`/`description` plus the
  still-rejected nested `ratingScale`.

- **Implementation**: Advanced change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  from `Design` to `In-Progress` so the schema, linter, and durable
  specification updates can land.

- **Design**: Refined change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md):
  made `Model.description` **Optional** (was `Recommended`), matching
  `Target.description`, so `description` reads uniformly across the tree. Updated
  the [functional spec](archive/0008-target-display-fields/spec.md) (Model schema now shows
  `description` as Optional) and the
  [design doc](archive/0008-target-display-fields/design.md): the `OptionalPresence`
  addition, the `# Optional` Model snippet, the composition alternative (now closer
  since `title`/`description` presence matches, still rejected on the `model-content`
  `RequiredAny` group and the mid-list `ratingScale` splice), and the trade-off note.

- **Design**: Advanced change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0008-target-display-fields/design.md). Reading the code showed the
  change is almost entirely schema + prose: three property additions in
  [`internal/schema`](../internal/schema/schema.go) (`Target` gains `title`/
  `description`, `Model` gains `description`) drive everything, because
  `misplaced-root-key` already fires on "a key valid on `Model` but not on this
  target" and so **self-narrows** to `ratingScale` with no rule-logic change, and
  the data-driven spec↔schema consistency test passes once the declarations and
  `SPECIFICATION.md` snippets move in lockstep. The only test edit flips
  `rules_test.go`'s "nested target title" case from a finding to a valid model.
  Recorded the decisions: keep `Model`/`Target` as explicit property lists rather
  than composing `Model` from `Target` (their `title` presence diverges —
  `Recommended` on the root, `Optional` on a target — and the consistency test
  already guards drift); and **not** add a `missing-target-description` warning
  (`RecommendedPresence` is documentary, not auto-enforced), leaving it as a noted
  follow-up. Updated the change [index](archive/0008-target-display-fields/index.md) and
  bundle [index](index.md).

- **Creation**: Added change
  [0008 — Describe targets with title and description](archive/0008-target-display-fields.md)
  (`status: Draft`) with its
  [functional spec](archive/0008-target-display-fields/spec.md). A target's only
  human-facing label today is its map key; the change lets every target carry an
  optional `title` (display name) and a recommended `description` (what the target
  *is*), and adds `description` to the model root. It also reframes the root as a
  **Model** — the model-wide `ratingScale` plus all **Target** properties — so
  the difference from a nested target is a type distinction (`ratingScale` is the
  one Model-only property) rather than the awkward "a non-root target MUST NOT
  declare …" prohibition, matching how `internal/schema` already models the two
  as distinct nodes. Records the
  [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`lint` sub-spec](../specs/cli/lint.md) deltas plus the `internal/schema` and
  `internal/lint` conformance (the `misplaced-root-key` rule and its
  "nested target title" test case) as affected. Updated the bundle
  [index](index.md). A design doc follows at **Design** for how
  `misplaced-root-key` narrows to `ratingScale` alone.

- **Completion**: Implemented and archived
  [0007 — Delightful human CLI output](archive/0007-delightful-cli-output.md),
  giving the human surface a single brand palette shared with the Fang harness, a
  styled [`lint`](../specs/cli/lint.md) finding list (severity glyphs, color,
  clickable `file:line`, colored summary) and [`init`](../specs/cli/init.md)
  confirmation, runnable `--help` examples on all three commands,
  [`spec`](../specs/cli/spec.md) paging through `$PAGER`/`less`, and an
  informative `--version` recovered from the Go toolchain's embedded build info.
  All of it sits behind the TTY/`NO_COLOR` gate, so the agent-facing plain and
  `--json` paths are byte-for-byte unchanged. Added the **Human output styling**
  and **Binary version** conventions to the [CLI spec](../specs/cli.md), the
  paging clause to the [`spec` sub-spec](../specs/cli/spec.md), the shared
  `internal/cli/style.go`, and focused tests; the styling consolidates onto one
  `colorEnabled` gate, retiring `spec`'s bespoke `shouldRenderSpec`.

- **Completion**: Implemented and archived
  [0006 — Specify and implement the spec command](archive/0006-spec-command.md),
  replacing the placeholder [`spec` sub-spec](../specs/cli/spec.md), adding the
  [design doc](archive/0006-spec-command/design.md), registering
  `qualitymd spec`, embedding [`SPECIFICATION.md`](../SPECIFICATION.md) in the
  binary, rendering Markdown for TTY output with Glamour while preserving
  byte-for-byte Markdown for redirected/plain output, and updating the
  [`README.md`](../README.md) command status.

- **Completion**: Implemented and archived
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md),
  adding `internal/schema` as the typed structural schema declaration, deriving
  lint's unknown-key, shape, required-property, model-content, and rating-scale
  minimum checks from it, and adding tests that compare the declaration against
  [`SPECIFICATION.md`](../SPECIFICATION.md). Reconciled the public format
  snippet's `title` presence to `Recommended`, updated
  [`lint`](../specs/cli/lint.md) to record schema-derived structural validation,
  and added the [design doc](archive/0005-schema-source-of-truth/design.md).

- **Refinement**: Added a human-readable rendering to change
  [0006 — Specify and implement the spec command](0006-spec-command.md). `spec`
  now **SHOULD** render the specification formatted (via the stack's terminal
  renderer) when stdout is a terminal, while still writing **verbatim Markdown**
  when output must be plain (non-terminal or `NO_COLOR`) — so a redirect still
  reproduces the artifact byte-for-byte and the `--json` carve-out is unaffected.
  The rendered/verbatim split needs no flag: it rides the
  [baseline](../specs/cli.md#baseline)'s existing terminal-detection rule, exactly
  as color does. Updated the [functional spec](0006-spec-command/spec.md) and the
  change's scope.

- **Creation**: Added change
  [0006 — Specify and implement the spec command](0006-spec-command.md)
  (`status: Draft`) with its
  [functional spec](0006-spec-command/spec.md). Picks up the `spec` command that
  [0004](archive/0004-specify-agent-accessibility.md) deferred as "a separate
  change that inherits this baseline." The change settles the still-stub
  [`specs/cli/spec.md`](../specs/cli/spec.md) — whose open questions predate and
  now conflict with the agent-accessibility work (it floats a `--format json`
  form the settled [`--json` convention](../specs/cli.md#conventions) forbids) —
  and lands the command: emit the bundled format specification verbatim as
  Markdown to stdout, no arguments, no `--json` (the verbatim-artifact carve-out),
  full baseline conformance. Records [`specs/cli/spec.md`](../specs/cli/spec.md)
  and [`README.md`](../README.md) as affected; structured forms, sub-views, and
  `spec`-specific flags are deferred. A design doc follows at **Design** for how
  the root-level specification is embedded. Updated the bundle [index](index.md).

- **Refinement**: Recorded the schema-source direction for change
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md):
  a **typed Go declaration** the linter derives from directly (over an embedded
  data file or a `specs/` concept), with spec/linter consistency enforced by a
  test checking [`SPECIFICATION.md`](../SPECIFICATION.md) against it rather than by
  generating docs — lowest-machinery path that meets the spec's requirements, with
  data-file/generation revisited only if a second consumer appears. Left
  unknown-key typo suggestions as a deferred follow-up, untouched. Detail lands in
  the design doc at **Design**.

- **Creation**: Added change
  [0005 — Single source of truth for the structural schema](archive/0005-schema-source-of-truth.md)
  (`status: Draft`) with its
  [functional spec](archive/0005-schema-source-of-truth/spec.md). Prompted by reviewing
  design.md's linter, which derives its structural rules from one schema artifact:
  our structural schema is encoded twice — implicitly in `internal/lint/rules.go`
  and again in prose in [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`specs/`](../specs/index.md) — so the two can drift. The change requires a
  single authoritative definition of valid keys, required properties, and the
  rating-scale shape that the linter derives from, as a behavior-preserving
  refactor; records [`SPECIFICATION.md`](../SPECIFICATION.md) and
  [`specs/cli/lint.md`](../specs/cli/lint.md) as affected durable docs; and defers
  doc generation, runtime configuration, and unknown-key typo suggestions. Updated
  the bundle [index](index.md).

- **Completion**: Implemented and archived
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md),
  adding the durable [CLI spec](../specs/cli.md) agent-accessibility contract,
  categorized exit codes (`0`/`1`/`2`/`70`), the broadened `--json` convention,
  and the [`init --json`](../specs/cli/init.md) receipt contract. The
  implementation maps exit categories through Fang, suppresses duplicate stderr
  for already-reported lint findings, adds the neutral `internal/receipt` action
  type, and tests the exit categories plus receipt and overwrite-refusal shapes.
- **Design**: Advanced change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  from `Draft` to `Design` and added its
  [design doc](archive/0004-specify-agent-accessibility/design.md). Reading
  `fang@v1.0.0` confirmed `fang.Execute` returns the command error and never
  exits, so categorized exit codes ride Fang's intended model: a thin boundary
  mapping in `cli.Execute` (`errors.As` → category, default `ExitInternal`), a
  `CodedError` carrying the category plus a `Silent` flag, a `WithErrorHandler`
  that suppresses the already-reported `lint` found-problems error, and
  Cobra-native `FlagErrorFunc`/`Args` hooks to tag usage errors at their source —
  with only an unknown-subcommand string fallback left as an explicit open
  decision. Picked `0`/`1`/`2`/`70` for success / found-problems / usage /
  internal, broadened *internal error* to "could not complete the requested
  action" so guarded refusals (e.g. `init` overwrite without `--force`) have a
  home, and ruled `-` plus `--json` a usage error. Settled the one open
  design decision — keep a thin, owned prefix check so an unknown subcommand maps
  to `ExitUsage` (option a), failing safe to `ExitInternal` and pinned by a test.
  Reconciled the
  [functional spec](archive/0004-specify-agent-accessibility/spec.md)'s internal-error
  definition with that broadening, and updated the change
  [index](archive/0004-specify-agent-accessibility/index.md) and bundle
  [index](index.md).

- **Scope**: Broadened the `--json` convention within change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  after comparing against agent-first CLIs (e.g. Basecamp, where nearly every
  command accepts `--json`). The change now revises the
  [CLI spec](../specs/cli.md)'s `--json` convention from a narrow should-offer
  gate to a SHOULD-by-default: commands SHOULD offer `--json`, human rendering
  stays the default (no auto-JSON), a command MAY omit it only when its output is
  a verbatim artifact that is the payload (e.g. `spec`), and under `--json` a
  side-effecting command emits a **result receipt** rather than human prose. The
  conformance work gains `init --json` (a receipt of the written path / created
  flag / `nextActions`), and [`specs/cli/init.md`](../specs/cli/init.md) joins
  the affected durable docs — its "offers no `--json`" statement is replaced by
  the receipt contract. `spec` stays the deliberate carve-out. Updated the
  functional [spec](archive/0004-specify-agent-accessibility/spec.md), the change
  [index](archive/0004-specify-agent-accessibility/index.md), and the bundle
  [index](index.md).

- **Scope**: Expanded change
  [0004 — Specify and enforce agent accessibility](archive/0004-specify-agent-accessibility.md)
  from a spec-only change to spec **plus** conformance after auditing the shipped
  commands. `internal/cli` exits `1` on every error path, so a `lint` that *found
  problems* is indistinguishable from a usage or internal error — a baseline
  violation. The change now settles the **exit-code categories** (success,
  ran-but-found-problems, usage error, internal error) concretely, removes both
  *Agent-accessibility and CI requirements* and *Exit-code semantics* from the
  [CLI spec](../specs/cli.md)'s "To be specified" list, and brings `init` and
  `lint` into compliance with tests. Threading distinct codes through Fang is a
  real design question, so the change now carries a forthcoming
  [design doc](archive/0004-specify-agent-accessibility/design.md) (added at **Design**);
  the unimplemented [`spec`](../specs/cli/spec.md) command is scoped out as a
  separate change that inherits the baseline. Retitled accordingly and updated
  the change [index](archive/0004-specify-agent-accessibility/index.md) and bundle
  [index](index.md).

- **Creation**: Added change
  [0004 — Specify agent accessibility](archive/0004-specify-agent-accessibility.md)
  (`status: Draft`) with its
  [functional spec](archive/0004-specify-agent-accessibility/spec.md). The change settles
  the *Agent-accessibility and CI requirements* item on the
  [CLI spec](../specs/cli.md)'s "To be specified" list by adding an **Agent
  accessibility** section framed as two tiers: a baseline every in-scope command
  owes (non-interactivity, stdout-is-payload/stderr-is-everything-else,
  determinism, categorized exit codes, plain non-TTY output) and the opt-in
  capabilities (`--json`, `nextActions`, quiet/verbosity) gated by criteria and
  cross-referenced to the existing conventions. Records
  [`specs/cli.md`](../specs/cli.md) as the only affected durable doc; no command
  behavior changes. Omits a design doc as spec-only work. Updated the bundle
  [index](index.md).

- **Completion**: Implemented and archived
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md),
  adding `qualitymd lint`, the shared lint rule catalog, JSON and human output,
  deterministic finding ordering, in-place `--fix` repair for fixable findings,
  parser/render/write support in `internal/spec`, and focused tests for the rule
  set, output shape, exit behavior, and repair behavior. Updated the README
  status and moved the change into [`archive/`](archive/).

- **Revision**: Hardened [0003's design doc](archive/0003-implement-lint-command/design.md)
  after review. Gave `internal/spec` a one-way dependency — it owns the document
  layer and no longer imports `internal/lint`, which now owns the rule catalog and
  the valid-model convenience (`lint.Load` replacing `spec.Load`) — removing a
  `spec`↔`lint` import cycle and a duplicate validator. Routed misplaced root-only
  keys (`title`/`ratingScale` on a non-root target) to `misplaced-root-key`
  instead of `invalid-frontmatter`; added the original Markdown body to the
  `Document` model so `Render` preserves it byte-for-byte; noted that `Load`'s
  acceptance tightens under the required-frontmatter parser; had `lint` reject a
  bare `-` this phase; clarified that post-repair `summary` counts (including
  `fixable`) reflect the re-lint; and reframed Resolved Questions as Open
  questions with the parent-CLI invocation as the one genuinely-open item.
  Recorded the provisional `lint [path]` shape as deliberately not durably specced
  in the [change](archive/0003-implement-lint-command.md).
- **Revision**: Worked down the open questions and risks in
  [0003's design doc](archive/0003-implement-lint-command/design.md): kept the shared
  document/model code in `internal/spec`, assigned rule-level repair operations
  to `internal/lint` and rendering/atomic replacement to `internal/spec`,
  resolved unknown frontmatter keys as `invalid-frontmatter` in this phase,
  confirmed `lint [path]` as the local invocation shape, and added mitigations
  for YAML round-tripping, deterministic ordering, atomic replacement, and
  symlink paths.
- **Revision**: Scoped `--fix` into change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md) after
  reviewing fixable-rule behavior. The durable lint spec, implementation spec,
  and design now require deterministic in-place repair of fixable findings,
  transactional per-file writes, post-repair linting, and JSON repair reporting,
  while keeping suppression, rule selection, and patch/full-file repair output
  modes deferred.
- **Design**: Advanced change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md) from
  `Draft` to `Design` and added its
  [design doc](archive/0003-implement-lint-command/design.md): `lint` parses once into a
  shared document/model graph with stable `modelPath` locations and optional
  source positions, runs narrow rule visitors from `internal/lint`, exposes the
  traversal primitives needed by current rules and future query commands, and
  adds a narrow repair writer for `lint --fix`. The design uses `lint [path]`,
  defaulting to `QUALITY.md`, as the minimum invocation shape while the parent
  CLI spec continues to own the broader file/stdin convention. Updated the
  change [index](archive/0003-implement-lint-command/index.md).

- **Creation**: Added change
  [0003 — Implement the lint command](archive/0003-implement-lint-command.md)
  (`status: Draft`) with a child
  [functional spec](archive/0003-implement-lint-command/spec.md). The change defers
  command-specific behavior to the completed durable
  [`lint` sub-spec](../specs/cli/lint.md), records README status updates as the
  durable docs work before Done, and calls out the remaining cross-cutting CLI
  invocation/file-argument convention as a dependency to settle before Design.
  Updated the bundle [index](index.md).

- **Archival**: Retired the placeholder [0001 — Example change](archive/0001-example-change.md)
  into [`archive/`](archive/) now that the bundle has real changes to follow,
  keeping it as the reference template the
  [propose-a-change guide](../docs/guides/work-with-changes.md) points to. Set its
  status to `Done`, fixed the relative links for the deeper path, and updated the
  bundle [index](index.md) and the [archive index](archive/index.md).

- **Completion**: Implemented and archived
  [0002 — Specify the init command](archive/0002-init-command.md), adding
  `qualitymd init`, replacing the durable [`init` sub-spec](../specs/cli/init.md),
  and updating the README status.

- **Refinement**: Tightened change [0002 — Specify the init command](archive/0002-init-command.md)
  after review: framed implementation as the change's own **In-Progress** phase
  rather than deferred work, specified that a successful `init` writes its
  confirmation to standard error (keeping stdout clean for `-` piping), recorded
  the non-atomic `--force` overwrite as a [design](archive/0002-init-command/design.md)
  risk, and trimmed the `--json` note in the
  [functional spec](archive/0002-init-command/spec.md) to a pointer to the
  [CLI spec](../specs/cli.md) convention.

- **Design**: Advanced change [0002 — Specify the init command](archive/0002-init-command.md)
  from `Draft` to `Design` and added its [design doc](archive/0002-init-command/design.md):
  the scaffold ships as a static `//go:embed` asset (comments and body prose can't
  round-trip through YAML struct marshalling), overwrite protection rides on an
  atomic `O_CREATE|O_EXCL` open, and a conformance test runs the embedded skeleton
  through `spec.Load`. Updated the change [index](archive/0002-init-command/index.md).

- **Creation**: Added change [0002 — Specify the init command](archive/0002-init-command.md)
  (`status: Draft`) with its [functional spec](archive/0002-init-command/spec.md), settling
  the "To be specified" list on the [`init` sub-spec](../specs/cli/init.md): the
  scaffold contents (seeded rating scale, a commented target → factor → requirement
  skeleton, recommended body sections as headed stubs), the output target and
  stdout (`-`) piping, and `--force` overwrite protection. Records
  [`specs/cli/init.md`](../specs/cli/init.md) and [`README.md`](../README.md) as
  affected. Updated the bundle [index](index.md).

- **Process**: Defined the relationship between `changes/` and the enduring
  [`specs/`](../specs/index.md) bundle (replacing the "independent for now"
  note) — a change states a *delta* and is archived, while `specs/` and
  [`SPECIFICATION.md`](../SPECIFICATION.md) hold the *cumulative* source of
  truth. Added an **Affected specs & docs** section to the
  [Change concept](archive/0001-example-change.md) so each change records the durable
  specs and docs it creates or updates, brought into sync before `Done`.

## 2026-06-16

- **Initialization**: Created the `changes/` OKF bundle — a home for incremental
  work, independent of [`specs/`](../specs/index.md) for now. Added the bundle
  [index](index.md), [`schema.md`](schema.md) (`type: Schema`) registering the
  `Change`, `Functional Specification`, and `Design Doc` types, and an
  [`archive/`](archive/) folder for completed changes.
- **Creation**: Added a placeholder [Example change](archive/0001-example-change.md)
  (`status: Draft`) with child [spec](archive/0001-example-change/spec.md) and
  [design](archive/0001-example-change/design.md) concepts showing the intended shape.
