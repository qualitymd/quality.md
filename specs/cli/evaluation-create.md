---
type: Functional Specification
title: qualitymd evaluation create
description: Create a numbered evaluation run folder.
tags: [cli, command, evaluation]
timestamp: 2026-06-21T00:00:00Z
---

# qualitymd evaluation create

`qualitymd evaluation create` creates a numbered evaluation run folder. It
inherits the cross-cutting CLI contract from [qualitymd CLI](../cli.md) and
produces the layout defined by the
[Evaluation](../evaluation/evaluation.md) contract.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Arguments and flags

- `[model]` — selected `QUALITY.md` file to snapshot; defaults to `QUALITY.md`.
- `--narrowing <slug>` — optional path-safe full structural scope path slug.
- `--model <path>` — selected `QUALITY.md` file to snapshot.
- `--evaluation-dir <path>` — override the evaluation directory.
- `--json` — emit a receipt on stdout.

## Requirements

The command **MUST** resolve a QUALITY.md workspace from the selected model file.
The workspace includes the selected model path, the repository root found from
that model path, the config file, the `.quality/` quality data directory, the
evaluation directory, and the quality log directory.

The selected model's root frontmatter key `config` **MAY** point to the
workspace config file. When present, `config` **MUST** be a non-empty scalar
repository-relative path, **MUST NOT** be absolute, and **MUST NOT** escape the
repository after path normalization. When absent, the config file defaults to
`.quality/config.yaml`. If the resolved config file is absent, the command
**MUST** use built-in defaults.

The command **MUST** resolve the evaluation directory using this precedence:
`--evaluation-dir`, then `evaluationDir` in the resolved config file, then
`.quality/evaluations/`. The path **MUST** be repository-relative and **MUST NOT**
escape the repository.

The command **MUST** validate the model path before creating the evaluation
directory or run folder. The model path **MUST** resolve to a file, not a
directory. Invalid model paths **MUST** fail without creating a numbered run
folder.

The command **MUST** compute the next run number as one past the highest
recognized evaluation run folder, create the run directory, create `data/`, and
snapshot `model-snapshot.md`. Recognized run folders use the
`NNNN-<scope>-eval` grammar.

New run folders **MUST** be named `NNNN-<scope>-eval`. When `--narrowing` is
absent, `<scope>` **MUST** be `full`, producing `NNNN-full-eval`. When
`--narrowing` is present, `<scope>` **MUST** be the validated narrowing slug,
producing `NNNN-<narrowing>-eval`.

`--narrowing` **MUST** remain a path-safe slug and **MUST NOT** include `quality`
as a slug segment. Callers that narrow by Area or Factor **SHOULD** use the
scope's full structural path, joining the Area path from the root and, for
Factor scope, the Factor path with single hyphens and no kind marker or boundary
separator. The run number remains the run identity; the scope slug is a human
mnemonic, and `model-snapshot.md` remains the structural source of truth.

The command **MUST NOT** create previous-runtime record folders or planning
coverage files such as `assessments/`, `analysis/`, `recommendations/`,
`design.md`, or `plan.md`.

`model-snapshot.md` is a frozen copy of the resolved model file. New run names
**MUST NOT** include an altitude segment. The command **MUST NOT** expose an
altitude flag, option, or JSON receipt field.

On success, human output **MUST** report the created path on stderr. Under
`--json`, stdout **MUST** contain a receipt with `schemaVersion`, `path`, and
`nextActions`.

The next action **MUST** show stdin data persistence with
`qualitymd evaluation data set <run> < payload.json`.
