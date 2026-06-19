# Evaluation design

## Subject and scope

Whole-model evaluation of `QUALITY.md` against itself: the three modeled
targets `format-spec` (`SPECIFICATION.md`), `readme` (`README.md`), and `cli`
(`internal/cli`, governed by `specs/cli.md` + `specs/cli/` and
`docs/guides/cli-design.md`). No narrowing; every in-scope requirement is
covered.

## Approach

Evidence collection was fanned out per target to structured-finding subagents
(format-spec read in full; readme read in full with CLI claims verified against
the running binary; cli verified by reading the Go implementation and exercising
the built binary). All rating judgment, factor and target roll-ups, and the
rating-binding re-check were performed by the orchestrating evaluator, not the
subagents.

## In-scope areas

- `SPECIFICATION.md` in full (frontmatter schema, terminology, evaluation and
  report semantics, extensions, appendices).
- `README.md` in full, with every command/flag/capability claim checked against
  the `qualitymd` 0.2.2 release binary (current HEAD `3db41f8`).
- `qualitymd` CLI: `init`, `lint`, `spec`, and the `evaluation` run-record
  surface (`create-run`, `add-record`, `set-planned-coverage`, `show-status`,
  `build-report`), checked against `specs/cli.md`, `specs/cli/*`, and
  `docs/guides/cli-design.md`.

## Evidence basis

- Static reads of `SPECIFICATION.md`, `README.md`, `internal/cli/*.go`,
  `specs/cli*`, and `docs/guides/cli-design.md`.
- Runtime exercise of the CLI: `--help` for every command, `--version`,
  `init`/`lint`/`spec` happy and error paths, `evaluation create-run`,
  `show-status --json`, and `build-report` (including `--fail-at-or-below`
  gating and byte-stability of `report.md`).
- Independent re-check of the two rating-binding CLI findings and the README
  payoff/first-result absence (see plan).

## Rating-binding findings (re-checked)

1. `show-status --json` emits an **absolute** `path` while `create-run --json`
   emits a **repository-relative** path — confirmed against the running binary
   and against `specs/cli/evaluation-create-run.md:33` (relative MUST) and
   `specs/cli.md:99-100` (deterministic, no machine-varying payload values).
2. An empty run reports gap kind `missing-analysis`, but
   `specs/cli/evaluation-show-status.md:29-32` mandates `missing-root-analysis`
   when no root-analysis record exists — confirmed MUST-level divergence.
3. `README.md` contains no rendered `qualitymd` output or report excerpt
   (only the conceptual `assessment -> findings -> rating result` diagram),
   confirming the payoff-by-produced-output and quick-first-result gaps.

## Out of scope

- Dependencies the project does not own (Go toolchain, Cobra/Fang, release
  tooling), per the model's Scope.
- The `/quality` skill prompts and modes themselves (not modeled targets).
- Exhaustive field-by-field diff of the README schema table against
  `SPECIFICATION.md` (spot-checked only; recorded as a limitation).

## Limitations

- `cli` runtime checks were run against the installed 0.2.2 release binary
  built from the current commit, not a fresh local build (the sandboxed local
  build produced a non-executable artifact). Behavior is expected to match
  HEAD source.
- README package-name claims (`npx skills add ...`, `npm install quality.md`)
  could not be verified against a registry; treated as provisional install
  placeholders, not rating-binding.
