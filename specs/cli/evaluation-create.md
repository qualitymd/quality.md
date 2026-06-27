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
- `--area <area-id>` — optional canonical Area reference for the run scope.
- `--factor <factor-id>` — optional canonical Factor reference; repeatable.
- `--model <path>` — selected `QUALITY.md` file to snapshot.
- `--evaluation-dir <path>` — override the model-relative evaluation directory.
- `--json` — emit a receipt on stdout.

## Requirements

The command **MUST** resolve a QUALITY.md workspace from the selected model file.
The workspace includes the selected model path, the workspace root directory
containing that model, the repository root found from that model path, the
config file, the `.quality/` quality data directory, the evaluation directory,
the quality changelog directory, and the workflow feedback-log directory. The
workspace root is the default base for relative tooling paths; the repository
root is the containment boundary.

The selected model's root frontmatter key `config` **MAY** point to the
workspace config file. When present, `config` **MUST** be a non-empty
model-relative scalar path, **MUST NOT** be absolute, and **MUST NOT** escape the
repository after path normalization from the workspace root. When absent, the
config file defaults to `.quality/config.yaml` under the workspace root. If the
resolved config file is absent, the command **MUST** use built-in defaults.

The command **MUST** resolve the evaluation directory using this precedence:
`--evaluation-dir`, then `evaluationDir` in the resolved config file, then
`.quality/evaluations/`. The path **MUST** be model-relative and **MUST NOT**
escape the repository after normalization from the workspace root.

The command **MUST** validate the model path before creating the evaluation
directory or run folder. The model path **MUST** resolve to a file, not a
directory. Invalid model paths **MUST** fail without creating a numbered run
folder.

The command **MUST** compute the next run number as one past the highest
`RunManifest.number` in the evaluation directory, create the run directory,
create `data/`, snapshot `model-snapshot.md`, and write
`data/run-manifest.json`.

`data/run-manifest.json` **MUST** be a CLI-owned `RunManifest` payload containing
`schemaVersion`, `kind`, `number`, `model`, `requestedScope`, and
`plannedScope`. `requestedScope` **MUST** faithfully record supplied scope
arguments. `plannedScope` **MUST** normalize the scope by defaulting `areaId` to
`area:root` and `factorFilter` to an array.

The command **MUST** validate `--area` and every `--factor` against the model
snapshot before creating a numbered run folder. If a Factor does not belong to
the planned Area, the command **MUST** fail without creating the run folder.

New run folders **MUST** be named `NNNN-<scope>-eval`. `<scope>` **MUST** be
derived from `plannedScope`: a root Area with an empty `factorFilter` produces
`NNNN-full-eval`; otherwise the slug is built from the planned Area path and any
filtered Factor structural paths. Callers **MUST NOT** supply the slug.

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
`qualitymd evaluation data set <run> < payloads.json`. When the selected model
is not the default `QUALITY.md` in the current working directory, the next
action **MUST** include `--model <model>` so the model-relative run path is
directly reusable.
