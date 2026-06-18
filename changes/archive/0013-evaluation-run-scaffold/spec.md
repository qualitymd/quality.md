---
type: Functional Specification
title: Evaluation run scaffold — functional spec
description: The behavior of qualitymd evaluation create-run, which scaffolds an evaluation run folder.
tags: [evaluation, specs, cli]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation run scaffold — functional spec

Companion to [Evaluation run scaffold](../0013-evaluation-run-scaffold.md). This
spec states *what* `qualitymd evaluation create-run` does; it does not specify the
implementation.

`evaluation create-run` inherits the cross-cutting CLI contract from the
[CLI spec](../../../specs/cli.md) — non-interactive operation, stdout/stderr
split, determinism, exit-code categories, and `--json` receipts. The run-folder
layout, naming, and numbering it produces are governed by the
[evaluation-record contract](../0012-evaluation-record-format/spec.md); this spec
defers to that contract and does not restate it.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: resolving the evaluation directory, computing the next run number,
creating the run folder with its required subdirectories and seed files, what the
command prints, and its flags and error cases.

Deferred: writing assessment, analysis, or recommendation records (change 0014,
`qualitymd evaluation add-record`); rendering `report.md`/`report.json` and gating (change
0015, `qualitymd evaluation build-report`); any narrowing semantics beyond carrying the
slug into the folder name.

## Invocation

The command is `qualitymd evaluation create-run`. It takes no positional path for the
subject; the model snapshot it seeds is sourced as described under
[Seed files](#seed-files).

### Flags

- `--altitude <subject|model>` (required) — the evaluation altitude. Any other
  value **MUST** fail with a usage error.
- `--narrowing <slug>` (optional) — a scope slug carried into the run-folder
  name. When given it **MUST** be a single path-safe slug; a value that is not
  **MUST** fail with a usage error.
- `--subject <path>` (optional) — the `QUALITY.md` to snapshot into the run's
  `model.md`. Defaults to `QUALITY.md` in the current working directory. For a
  `model` altitude this is the user's file (the subject); the active model
  snapshotted is the meta-model (see [Seed files](#seed-files)).
- `--evaluation-dir <path>` (optional) — override the evaluation directory for
  this invocation. It takes precedence over `.quality/config.yaml` and follows
  the same repository-relative normalization and escape checks.

## Resolving the evaluation directory

- The command **MUST** resolve the evaluation directory using this precedence:
  `--evaluation-dir <path>` when present; otherwise `evaluationDir` from
  `.quality/config.yaml` at the repository root when present; otherwise
  `quality/evaluations/`.
- A resolved `evaluationDir` **MUST** be a repository-relative normalized path.
  An absolute path, or a path that escapes the repository, **MUST** fail with the
  internal-error exit category and a diagnostic.
- The evaluation directory **MUST** be created if it does not yet exist. A
  missing evaluation directory is not an error — `evaluation create-run` is the
  command that populates it.

## Run number

- The command **MUST** compute `NNNN` deterministically as one past the highest
  existing run in the evaluation directory, across **all** altitudes, per the
  [evaluation-record contract](../0012-evaluation-record-format/spec.md#run-folder).
  An empty (or freshly created) evaluation directory yields `0001`.
- A `subject` run **MUST NOT** reuse a number already held by a `model` run, or
  vice versa: the sequence is shared.
- Existing entries in the evaluation directory that do not match the run-folder
  naming pattern **MUST** be ignored when computing the next number.

## Run folder

- The command **MUST** create the run folder named
  `NNNN-<altitude>[-<narrowing>]-quality-eval`, with `<narrowing>` present only
  when `--narrowing` was given.
- It **MUST** create the required subdirectories `assessments/`, `analysis/`, and
  `recommendations/`, so the writers in change 0014 have somewhere to write.
- It **MUST NOT** create the run folder if one with the computed name already
  exists; it **MUST** instead fail with the internal-error category, leaving the
  existing folder untouched. (Determinism plus the shared sequence means a
  collision signals concurrent or corrupt state, not something to overwrite.)
- The run folder and its contents are raw runtime outputs, not OKF concepts: the
  command **MUST NOT** seed any OKF `index.md`/`log.md`/`schema.md`, per the
  [contract](../0012-evaluation-record-format/spec.md#artifacts-are-not-okf).

## Seed files

The command **MUST** seed three files the skill then fills in:

- `model.md` — the active model the run evaluates against, **snapshotted by this
  command**, not by the skill. For `--altitude subject` it is the resolved
  `--subject` `QUALITY.md`; for `--altitude model` it is the bundled
  quality meta-model (the same model the skill obtains via
  `qualitymd models view quality-meta-model`). Snapshotting the model in the
  deterministic surface is deliberate: the model is mechanically resolvable
  content, and capturing it here keeps the run's record of *what was evaluated*
  off the judging skill, consistent with the contract's CLI-writes /
  skill-judges division. The skill still authors `design.md` and `plan.md`,
  which are judgment.
- `design.md` — a stub for how the evaluation was scoped and run.
- `plan.md` — a stub for the assessment plan.

The `design.md` and `plan.md` stubs **MUST** be created empty or with a minimal
heading; the command **MUST NOT** invent their judgment content.

## Reporting

- On success the command **MUST** report the created run-folder path on standard
  error, so standard output is not polluted, per the
  [CLI baseline](../../../specs/cli.md#baseline).
- Under `--json` the command **MUST** emit a single result receipt on stdout
  instead of the human confirmation, carrying at least `schemaVersion`, the
  created run-folder `path`, the resolved `altitude`, and any `nextActions`.
- It **SHOULD** close with [next actions](../../../specs/cli.md#conventions)
  pointing at the next step in the run — recording results (change 0014).

## Exit codes

Per the [CLI exit-code categories](../../../specs/cli.md#exit-codes):

- `0` — the run folder was scaffolded.
- `2` — usage error: missing or invalid `--altitude`, an invalid `--narrowing`
  slug, or an unknown flag.
- `70` — internal error: a configured `evaluationDir` that is absolute or escapes
  the repository, a run-folder name collision, the `--subject` file missing for a
  `subject` altitude, or an I/O failure.
