---
type: Functional Specification
title: qualitymd evaluation create-run
description: Create a numbered evaluation run folder.
tags: [cli, command, evaluation]
timestamp: 2026-06-18T00:00:00Z
---

# qualitymd evaluation create-run

`qualitymd evaluation create-run` creates a numbered evaluation run folder. It
inherits the cross-cutting CLI contract from [qualitymd CLI](../cli.md) and
produces the layout defined by [Evaluation records](../evaluation-records.md).

## Flags

- `--altitude <subject|model>` — required.
- `--narrowing <slug>` — optional path-safe scope slug.
- `--subject <path>` — repository-relative `QUALITY.md` file to snapshot;
  defaults to `QUALITY.md`.
- `--evaluation-dir <path>` — override the evaluation directory.
- `--json` — emit a receipt on stdout.

## Requirements

The command **MUST** resolve the evaluation directory using this precedence:
`--evaluation-dir`, then `.quality/config.yaml` `evaluationDir`, then
`quality/evaluations/`. The path **MUST** be repository-relative and **MUST NOT**
escape the repository.

The command **MUST** validate the subject path before creating the evaluation
directory or run folder. The subject path **MUST** be repository-relative,
**MUST NOT** escape the repository, and **MUST** resolve to a file, not a
directory. Invalid subject paths **MUST** fail without creating a numbered run
folder.

This prevalidation exists because creating the run folder before validating
`--subject` could fail mid-write — leaving an empty run skeleton on disk and
consuming a run number with no records. Validating first guarantees an invalid
subject produces no on-disk artifacts and no run-number gap. The subject snapshot
bytes are prepared before the evaluation directory is created so a validation
failure touches no disk state; rolling back partially created folders was
deliberately avoided to prevent deleting a directory a concurrent run created.

The command **MUST** compute the next run number as one past the highest matching
run folder across both altitudes, create the run directory, create
`assessments/`, `analysis/`, and `recommendations/`, and seed `model.md`,
`design.md`, and `plan.md`.

For subject altitude, `model.md` is the resolved subject file. For model
altitude, `model.md` is the bundled `quality-meta-model` with its source pointed
at the validated subject file.

On success, human output **MUST** report the created path on stderr. Under
`--json`, stdout **MUST** contain a receipt with `schemaVersion`, `path`,
`altitude`, and `nextActions`.
