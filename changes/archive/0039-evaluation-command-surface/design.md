---
type: Design Doc
title: Evaluation command surface redesign - design doc
description: How the reshaped qualitymd evaluation surface is built — cobra tree, shared run/payload helpers, plan.md-folded coverage, and the report build/gate split.
tags: [cli, evaluation, surface]
timestamp: 2026-06-19T00:00:00Z
---

# Evaluation command surface redesign - design doc

Design behind the
[Evaluation command surface redesign](../0039-evaluation-command-surface.md)
change case and its [functional spec](spec.md). The spec says *what* the surface
must be; this doc says *how* the Go code delivers it and why these choices over
the alternatives. Behavior and requirements live in the spec — this doc does not
restate them.

## Context

The evaluation commands are wired in `internal/cli/evaluation.go` over logic in
`internal/evaluation/` (`create.go`, `write.go`, `planned_coverage.go`,
`path.go`, `types.go`, `load.go`, `report.go`). The reshape is almost entirely a
command-tree and artifact-location change: the judgment schemas, numbering,
rendering, and gap detection stay; what changes is how commands are named and
nested, where planned coverage is stored, and how the gate is invoked.

The work is breaking by design (pre-1.0, single in-repo consumer), so the CLI and
the bundled `/quality` skill change together in the implementation phase.

## Approach

### Command tree

Rebuild the cobra tree in `evaluation.go` so the parent `evaluation` command holds
the run-level verbs directly and a small set of noun subcommands each holds their
verbs:

```text
evaluation
├── create                 (NoArgs)
├── list                   (NoArgs)
├── status      <run|--latest>
├── assessment
│   ├── add     <run|--latest>
│   └── list    <run|--latest>
├── analysis
│   ├── set     <run|--latest>
│   └── list    <run|--latest>
├── recommendation
│   ├── add     <run|--latest>
│   └── list    <run|--latest>
└── report
    ├── build   <run|--latest>
    └── gate    <run|--latest>
```

The existing `newEvaluationAddRecordKindCmd` factory (which built the three
`add-record` kind subcommands) is replaced by three noun-builders, each
constructing its `add`/`set` and `list` subcommands. The noun parents are pure
groupers with no action of their own — the same pattern cobra already gives the
current `add-record` parent.

### Shared helpers

Three concerns repeat across the new verbs, so they become small shared functions
rather than per-command copies:

- **Run resolution.** `resolveRun(args, latest, dir)` returns the run path from
  either the positional argument or `--latest`, erroring when neither or both are
  supplied. `--latest` reuses the existing directory scan that backs
  `ListRunDirs`/`nextRunNumber` to pick the newest run. Every run-scoped command
  calls this instead of taking `cobra.ExactArgs(1)` directly.
- **Payload reading + batching.** The current `readPayload` already centralizes
  `--file`/`-`/stdin and terminal-stdin rejection. Extend the decode step to
  accept a single object *or* a JSON array of that kind: peek the first
  non-whitespace byte (`[` vs `{`), normalize to a slice, then validate and write
  each element in order with `DisallowUnknownFields` per element. The write
  receipt aggregates the written paths into an array.
- **Output streams.** A single convention enforced at the command layer: the
  requested data/result and every `--json` payload go to `cmd.OutOrStdout()`;
  side-effect confirmation lines go to `cmd.ErrOrStderr()`. In practice the
  writers already log "Wrote …" to stderr and `show-status` already reads to
  stdout; the change is to state the rule, route `list` to stdout, and make
  `report gate` emit only an exit code plus an optional stderr line.

### Planned coverage in `plan.md`

This is the one change that moves data, not just names.

- **Seeding.** `create` seeds `plan.md` as **body-only** (prose headings), with no
  `coverage:` frontmatter. Absent frontmatter — or an absent `coverage` key —
  means "no planned coverage declared," so no drift checks run. The skill *adds* a
  `coverage:` frontmatter block when it chooses to declare coverage. Seeding an
  empty `coverage:` is deliberately avoided, because an empty block would read as
  "planned nothing" and make every written record `unexpected`.
- **Reading.** `load.go` stops reading `planned-coverage.json` and instead parses
  `plan.md`'s YAML frontmatter using the same frontmatter parser the loader
  already applies to `model.md`/`QUALITY.md`. The `PlannedCoverage`,
  `PlannedCoverageAssessment`, and `PlannedCoverageAnalysis` types stay but gain
  YAML mapping (they were JSON-tagged for the old file). The validation,
  canonical sort, and identity functions in `planned_coverage.go`
  (`validatePlannedCoverage`, `sortPlannedCoverage`, `plannedAssessmentIdentity`,
  `plannedAnalysisIdentity`) are kept and now run at **read time**.
- **Drift gaps.** `plannedCoverageGaps` is unchanged in spirit — it still emits
  `missing-planned-*` and `unexpected-*` by comparing planned identities against
  written-record identities. Only its source changes (frontmatter, not the JSON
  file).
- **Malformed coverage.** Today a malformed `planned-coverage.json` is a hard load
  error. Under the new rule, malformed coverage frontmatter becomes a **gap**
  (e.g. `invalid-plan-coverage`) so `status` surfaces it without making the whole
  run unloadable. A `plan.md` with no parseable frontmatter is treated as
  body-only, not an error.
- **Removal.** The `SetPlannedCoverage` command path and `planned-coverage.json`
  writing go away; `planned_coverage.go` keeps only the read/validate/identity
  helpers.

### Report build/gate split

`report.go` already separates `BuildReport` from `Gate`/`ScaleLevels`; the split
becomes two cobra commands:

- `report build` calls `BuildReport`, rendering the three artifacts.
- `report gate` reads the already-rendered `report.json` from the run, extracts
  the root aggregate rating, applies the existing `Gate` decision against
  `--at-or-below <level>`, and exits `0` (better), `1` (at-or-below, including a
  not-assessed root), or `2` (level not in the run's scale). It writes nothing and
  errors if `report.json` is absent.

The key code change is that `gate` no longer renders or writes on a failing build
— the previous `build-report --fail-at-or-below` wrote files even when gating.
A small reader pulls the root rating from `report.json` rather than recomputing
from records, keeping `gate` cheap and side-effect-free.

### Altitude removal

Delete `Options.Altitude` (`types.go`), the unreachable non-`subject` guard in
`CreateRun` (`create.go`), and the `altitude` field on the create receipt. New
runs are always subject-altitude. `path.go`'s run-folder regex keeps accepting
historical `model-` folders so old runs still load.

## Alternatives

- **`record <kind>` verb instead of resource nouns.** Rejected. A single `record`
  verb keeps the kind as an argument, but it calls the analysis upsert "add" and
  leaves nowhere natural for per-kind `list`. Promoting the kinds to nouns lets
  each verb (`add` vs `set`) state the on-disk contract and gives `list` a home.
- **Keep planned coverage as its own command/artifact (`coverage`/`plan set`).**
  Rejected. It was a second structured file and a second write path for one
  concept — the run's plan. Folding it into `plan.md` frontmatter matches the
  `QUALITY.md` frontmatter+body shape and deletes a command and an artifact.
- **CLI-managed `plan.md` frontmatter (write-time validation).** Considered and
  rejected for this case. A command that rewrites only the frontmatter while
  preserving the hand-edited body keeps coverage on the validated-write path, but
  it reintroduces a verb we are removing and a merge-frontmatter-preserve-body
  operation. Hand-authoring with **read-time** validation in `status` is the
  larger simplification and is consistent with how `plan.md`'s body is already
  authored. The cost — coverage keys are no longer canonicalized at write — is
  absorbed by comparing on normalized identities at read time.
- **`report build --gate <level>` convenience flag.** Rejected. Re-adding the flag
  quietly recreates the render+gate overload this case removes. CI composes the
  two verbs (`build` then `gate`); a not-assessed root still fails the gate.
- **Implicit "newest run" default.** Rejected. A bare run-scoped command with no
  argument could silently target the wrong run. `--latest` stays an explicit
  opt-in alongside the positional argument.
- **Heterogeneous cross-kind batch** (one array mixing assessments, analyses, and
  recommendations). Rejected for now. Within-kind array batching covers the
  skill's "batch independent writes" need with simple, per-kind validation and
  clear error attribution; a mixed batch complicates numbering and error
  reporting for little gain.

## Trade-offs & risks

- **Hand-authored coverage can drift from record keys.** A fat-fingered
  `targetPath` or `requirement` in frontmatter would produce a false
  `missing-planned-*`/`unexpected-*` pair. Mitigated by comparing on the same
  normalized identity functions used for written records; residual risk is real
  but visible (it shows up as a drift gap, not a silent pass).
- **Larger command tree.** More cobra commands and help surfaces. Mitigated by the
  shared `resolveRun`/payload/stream helpers so the per-command code stays thin.
- **Breaking rename.** Any caller using the old names breaks. Acceptable: the only
  consumer is the bundled skill, updated in the same change; pre-1.0 has no
  external stability promise.
- **`plan.md` is now frontmatter-parsed.** A body that legitimately starts with a
  `---` rule could confuse a naive parser. Mitigated by treating unparseable or
  absent frontmatter as body-only rather than erroring.
- **Validation moves to read time.** `status` now does the duplicate-key and
  schema checks the write command used to. This is cheap and runs where the data
  is consumed, but it means an invalid plan surfaces at `status`, not at authoring
  time.

## Open questions

- **Flat specs vs a subfolder.** With ~7 evaluation command specs, the durable
  `specs/cli/` layout could stay flat (`evaluation-*.md`) or move to
  `specs/cli/evaluation.md` + `specs/cli/evaluation/`. The functional spec
  currently plans the flat layout; the subfolder is the open call to make before
  implementation.
- **`list --state` filter scope.** The spec marks state filtering `SHOULD`. Decide
  during implementation whether the first cut ships the filter or only the full
  enumeration.
- **`invalid-plan-coverage` gap naming.** The exact gap kind string for malformed
  coverage frontmatter is provisional and should be fixed alongside the other gap
  kinds in `load.go`.
