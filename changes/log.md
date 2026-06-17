# Changes Update Log

## 2026-06-17

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
