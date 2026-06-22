# Evaluation design

## Resolved parameters

- Altitude: subject
- Narrowing: Full evaluation
- Rigor: deep
- Scope: Full evaluation

## Scope

Full subject-quality evaluation of the root `QUALITY.md` model against current
repository sources after applying the confirmed fixes found during the run.

In scope:

- `format-spec` source: `SPECIFICATION.md`
- `readme` source: `README.md`
- `cli` source: `internal/cli`, checked against `specs/cli.md`,
  `specs/cli/*`, and `docs/guides/cli-design.md`
- Evaluation-history readiness as a CLI/status quality signal

Out of scope:

- Dependencies and release infrastructure not owned by this repository
- Runtime behavior outside the deterministic CLI surface
- Hosted CI and published release artifacts, which are verified separately by
  the release runbook

## Evidence basis

- Current source reads of `QUALITY.md`, `SPECIFICATION.md`, `README.md`,
  `specs/cli.md`, relevant `specs/cli/*`, `docs/guides/cli-design.md`, and
  `internal/cli` / `internal/evaluation` code.
- Behavior checks using `go run ./cmd/qualitymd ...` for lint, status, root
  welcome, and evaluation status.
- Focused regression tests for `internal/cli`, `internal/evaluation`, and
  `internal/status`.
- `qualitymd lint QUALITY.md` for model structural validity.

## Rating approach

Deep rigor covers every in-scope Requirement. A Target can reach `target` when
all binding requirements satisfy their authored assessments and remaining gaps
are bounded or non-blocking. `minimum` is used where a requirement holds the
accepted floor but still leaves a material release-relevant gap. Root aggregate
rating follows the weakest release-relevant child Target because the model's
Risks treat spec ambiguity, misleading README behavior, and non-deterministic
CLI behavior as release constraints.

## Limitations

- The CLI assessment samples representative command paths and focused tests
  rather than exhaustively executing every flag permutation.
- No hosted CI or package-publish verification is part of this evaluation run;
  those belong to the release verification stage.
