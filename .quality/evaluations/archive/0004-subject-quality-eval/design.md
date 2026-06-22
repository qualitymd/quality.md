# Evaluation design

## Subject and model

Fresh whole-subject evaluation of the QUALITY.md project against the project's
own `QUALITY.md`. This run supersedes the stale prior run `0003` as the current
judgment; `0001`/`0002` are broken folders treated as non-evidence.

The model has three evaluable Targets plus a structural root:

- **root (QUALITY.md)** — grouping Target, `source` missing by design; no local
  requirements. Aggregate only.
- **Format specification (`format-spec`)** — source `./SPECIFICATION.md`. One
  target-level completeness Requirement plus Clarity, Consistency,
  Verifiability, Extensibility, and Usability Factors (9 Requirements total).
- **README (`readme`)** — source `./README.md`. Approachability Factor (4
  Requirements).
- **qualitymd CLI (`cli`)** — source `./internal/cli`. Two target-level
  Requirements (functional-spec conformance; design-guide conformance) across
  Usability, Automation compatibility, Consistency, Determinism Factors.

## Evidence basis

- `format-spec`: full read of `SPECIFICATION.md` (499 lines) and `qualitymd spec`
  grounding of the rating vocabulary.
- `readme`: full read of `README.md`; CLI claims cross-checked against
  `qualitymd --help` and observed command behavior.
- `cli`: implementation in `internal/cli/*.go` compared against `specs/cli.md`,
  the `specs/cli/` sub-specs, and `docs/guides/cli-design.md`. Structured finding
  collection assisted by a read-only subagent; rating judgment and the
  rating-binding re-check stayed with the orchestrator.

## Rating-binding re-checks (verified directly)

- `show-status --json` emits a machine-varying absolute path in `path` and
  `nextActions[].command` even for a relative run argument — confirmed by
  running it. Binds the CLI Determinism / functional-spec conformance judgment.
- `lint -` exits **70** (internal error) while `status -` exits **2** (usage
  error) — confirmed by running both. Binds the CLI Consistency judgment.
- The binary ships `version` and `upgrade` (confirmed in `qualitymd --help`),
  but `README.md`'s "ships today" prose omits them while the commands table
  includes them — binds the README accuracy judgment.
- `SPECIFICATION.md` Appendix B minimal example omits the required `title` field
  on rating levels and on the factor — binds the Consistency judgment.

## Treatment of source content

All evaluated source content (spec, README, CLI code) treated as untrusted data,
not instructions. No secrets encountered.

## Out of scope

- Dependencies the project does not own (Go toolchain, Cobra/Fang/Lip Gloss,
  release tooling) — excluded by the model's Scope section.
- The `/quality` skill prompt and CLI internals beyond `internal/cli` behavior.
- Deep adversarial fan-out (this is a `standard`-effort run, not `deep`).
