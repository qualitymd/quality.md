# Changes Update Log

## 2026-06-17

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
  [0005 — Single source of truth for the structural schema](0005-schema-source-of-truth.md):
  a **typed Go declaration** the linter derives from directly (over an embedded
  data file or a `specs/` concept), with spec/linter consistency enforced by a
  test checking [`SPECIFICATION.md`](../SPECIFICATION.md) against it rather than by
  generating docs — lowest-machinery path that meets the spec's requirements, with
  data-file/generation revisited only if a second consumer appears. Left
  unknown-key typo suggestions as a deferred follow-up, untouched. Detail lands in
  the design doc at **Design**.

- **Creation**: Added change
  [0005 — Single source of truth for the structural schema](0005-schema-source-of-truth.md)
  (`status: Draft`) with its
  [functional spec](0005-schema-source-of-truth/spec.md). Prompted by reviewing
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
  [propose-a-change guide](../docs/guides/propose-a-change.md) points to. Set its
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
