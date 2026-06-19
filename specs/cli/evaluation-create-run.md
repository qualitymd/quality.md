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

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Flags

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
evaluation run folder, create the run directory, create `assessments/`,
`analysis/`, and `recommendations/`, and seed `model.md`, `design.md`, and
`plan.md`.

`model.md` is the resolved subject file.

On success, human output **MUST** report the created path on stderr. Under
`--json`, stdout **MUST** contain a receipt with `schemaVersion`, `path`, and
`nextActions`. The receipt **MAY** include `altitude: "subject"` for
compatibility with existing evaluation consumers.
