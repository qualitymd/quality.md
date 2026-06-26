---
type: Functional Specification
title: Evaluation run folder naming - functional spec
description: Requirements for naming Evaluation v2 run folders NNNN-full-eval (full) and NNNN-<scope-path>-eval (narrowed), where the narrowing slug is the scope's full structural path, while still recognizing legacy run folders.
tags: [cli, evaluation, naming, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation run folder naming - functional spec

Companion to the
[Evaluation run folder naming](../0113-evaluation-run-folder-naming.md) change
case. This spec states what the run-folder naming change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

Evaluation v2 names run folders `NNNN-quality-eval`, or
`NNNN-<narrowing>-quality-eval` when narrowed. The `quality-eval` tag duplicates
the `.quality/evaluations/` parent in folder listings and only earns its keep
when the run name travels alone as a handle (e.g. `run: 0003-quality-eval` in a
quality-log entry). The narrowing segment, when present, is the informative part
of the name, but it follows no convention.

This change shortens the constant tag to `eval` and makes a narrowed run's slug
the scope's full structural path. The run number remains the run's identity and
`model.md` remains the structural source of truth; the slug is a human mnemonic
that now carries maximal scope context. Existing runs are not migrated, so both
the old and new grammars must be recognized.

## Scope

Covered: the run-folder name grammar produced by `qualitymd evaluation create`,
the `--narrowing` slug convention, recognition of new and legacy run folders
across listing/status/next-number computation, and the durable CLI and
`/quality` contracts that describe them.

Not covered: migration or rewrite of existing run folders; an Area-vs-Factor
kind marker or boundary separator inside the slug; the run number as identity or
`nextRunNumber` monotonicity; and structured `data/` layout, generated report
filenames, or report content.

## Assumptions & dependencies

- `--narrowing` is caller-provided; `qualitymd evaluation create` validates it as
  a path-safe slug (lowercase `[a-z0-9]` tokens joined by single hyphens) and
  otherwise treats it as opaque. The full-structural-path convention is therefore
  enforced by the caller (the `/quality` evaluation workflow), not by the CLI.
- The run number remains the run's stable identity; the slug never has to be
  unique on its own.

## Requirements

### Run-folder name grammar

- Every newly created run folder **MUST** be named `NNNN-<scope>-eval`, where
  `NNNN` is the zero-padded run number and `<scope>` is a single scope segment
  between the number and the `eval` tag.

  > Rationale: a uniform `NNNN-<scope>-eval` shape keeps full and narrowed runs
  > symmetric and makes a full run self-declare its scope rather than read as a
  > missing segment. — 0113

- When `--narrowing` is absent, `qualitymd evaluation create` **MUST** set
  `<scope>` to the reserved marker `full`, naming the folder `NNNN-full-eval`.

- When `--narrowing` is present, `qualitymd evaluation create` **MUST** set
  `<scope>` to the validated `--narrowing` slug, naming the folder
  `NNNN-<narrowing>-eval`.

- The constant trailing tag in a newly created run folder **MUST** be `eval`,
  not `quality-eval`.

  > Rationale: the `.quality/evaluations/` parent already says "evaluation", so
  > `quality-eval` echoed it in every listing; `eval` keeps a self-identifying
  > noun on the bare run handle without the redundant `quality-`. — 0113

- `full` is reserved as the full-scope marker. A `--narrowing` slug equal to
  `full` is permitted and produces a name shaped like a full-scope run; the run
  number remains the identity that disambiguates.

  > Rationale: a structural scope path that slugs to exactly `full` is vanishingly
  > rare and not worth a new error path; the run number already carries identity,
  > so the cosmetic collision is tolerated rather than validated against. — 0113

### Narrowing scope slug

- When the `/quality` evaluation workflow creates a run narrowed to an Area, it
  **MUST** pass `--narrowing` the Area's full structural path from the root Area,
  with the path's structural segments joined by single hyphens.

- When the `/quality` evaluation workflow creates a run narrowed to a Factor, it
  **MUST** pass `--narrowing` the owning Area's structural path followed by the
  Factor's structural path, with all segments joined by single hyphens.

- The narrowing slug **MUST** remain a path-safe slug.

- The narrowing slug **MUST NOT** insert an Area-vs-Factor kind marker or any
  boundary separator between the Area-path segments and the Factor-path segments.

  > Rationale: the slug is a mnemonic, not an identifier — the run number is the
  > identity and `model.md` carries full structural identity. Path-safe slugs
  > have no boundary character available, so a kind marker would only add another
  > ambiguous token rather than a reliable boundary. — 0113

### Recognition and legacy runs

- `qualitymd` **MUST** recognize newly created run folders that use the `-eval`
  suffix wherever it recognizes run folders, including `evaluation create` (next
  number), `evaluation list`, and `evaluation status`.

- `qualitymd` **MUST** continue to recognize existing run folders that use the
  legacy `-quality-eval` suffix, including legacy `subject`- and `model`-prefixed
  forms, so existing runs remain listed, inspectable, and counted.

- `qualitymd evaluation create` **MUST** compute the next run number from both
  new `-eval` and legacy `-quality-eval` run folders, so numbering stays
  monotonic across the format change.

- `qualitymd evaluation list` **MUST** report no narrowing for a `full`-marked
  run and the narrowing slug for a narrowed run, for both new `-eval` and legacy
  `-quality-eval` run names.

- This change **MUST NOT** migrate, rename, or rewrite any existing run folder.

## Acceptance criteria

- Creating a full-scope run yields a folder named `0007-full-eval`.
- Creating a run narrowed to the top-level Area `security` yields
  `0007-security-eval`.
- Creating a run narrowed to a sub-Area `network` under `security` yields
  `0007-security-network-eval`.
- Creating a run narrowed to a Factor `reliability` under Area `security` yields
  `0007-security-reliability-eval`.
- Creating a run narrowed to a sub-Factor `latency` under Factor `reliability`
  under Area `security` yields `0007-security-reliability-latency-eval`.
- No newly created run folder contains the segment `quality`.
- `qualitymd evaluation list` and `qualitymd evaluation status` list and inspect
  both a new `0007-full-eval` folder and an existing `0006-quality-eval` folder.
- `qualitymd evaluation create` run after a legacy `0006-quality-eval` folder
  produces `0007-full-eval`, not `0001-full-eval`.
- `qualitymd evaluation list` reports no narrowing for `0007-full-eval`,
  narrowing `security-network` for `0007-security-network-eval`, and the stripped
  narrowing for a legacy `0005-subject-quality-eval` folder.
- Existing run folders on disk are byte-for-byte unchanged after the run.

## Durable spec changes

### To add

None.

### To modify

- [`specs/cli/evaluation-create.md`](../../../specs/cli/evaluation-create.md) -
  state the run-folder grammar `NNNN[-<narrowing>]-eval`, document that
  `--narrowing` carries the scope's full structural path, and note continued
  recognition of legacy `-quality-eval` folders. Driven by
  [Run-folder name grammar](#run-folder-name-grammar) and
  [Recognition and legacy runs](#recognition-and-legacy-runs).
- [`specs/skills/quality-skill/evaluation.md`](../../../specs/skills/quality-skill/evaluation.md)
  - in the create-run workflow step, record that the skill passes `--narrowing`
    the scope's full structural path. Driven by
    [Narrowing scope slug](#narrowing-scope-slug).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, new run folders read `NNNN-full-eval`
or `NNNN-<scope-path>-eval` with the scope's full structural path as the mnemonic,
existing legacy runs stay recognized and correctly numbered against, and no run
is migrated. That achieves the motivation — a cleaner, more informative run name
that drops the redundant tag — without widening into run migration, a new slug
boundary syntax, or any change to run identity or data layout.
