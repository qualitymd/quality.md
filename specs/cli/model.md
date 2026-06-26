---
type: Functional Specification
title: qualitymd model
description: Query a quality model's structure and canonical reference IDs read-only, via tree, list, and get.
tags: [cli, command, model]
timestamp: 2026-06-26T00:00:00Z
---

# qualitymd model

`qualitymd model` projects a `QUALITY.md` model as a read-only structure: its
Areas, Factors, and Requirements, each with the canonical reference ID used in
evaluation payloads, and how they contain one another. It has three verbs:
`tree` (hierarchical view), `list` (flat enumeration), and `get <id>`
(single-element detail).

It inherits the cross-cutting CLI contract from the [CLI spec](../cli.md). This
file specifies only what is particular to `model`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

Authoring evaluation payloads requires the canonical reference IDs of a model's
elements — `area:<path>`, `factor:<area>::<path>`, `requirement:<area>::<name>`.
No other command emits them: `status` reports source coverage by Area path and
label, not canonical IDs, and does not enumerate Factors or Requirements. Agents
therefore hand-derive tens of IDs from `QUALITY.md` text — slow, and a
silent-typo source. `model` makes those IDs a query for an agent (`--json`) and
a person (tree) alike.

The canonical reference grammar `model` emits and resolves is defined by
[`SPECIFICATION.md`](../../SPECIFICATION.md).

## Scope

Covered: the `model` noun and its `tree`, `list`, and `get` verbs over a model
file, each with a human default form and `--json`, all emitting canonical
reference IDs.

Deferred / non-goals (each absence is deliberate):

- No `model query` selector/expression language.
- No mutation — `model` is read-only.
- No model validation — `lint` owns diagnostics; `model` assumes a parseable
  model.
- No source/provenance, coverage, or aggregate count summaries — `status` owns
  those.
- No evaluation-run awareness — no `--run` flag and no run-folder resolution.
- No ratings or evaluation results.

## Boundary

The `model` command group **MUST NOT** emit source declaration/provenance,
source coverage, aggregate element-count summaries, evaluation results, or
readiness state. Those remain the responsibility of [`status`](status.md) and
the `evaluation` commands. `model` owns *logical structure and identity*:
elements, canonical IDs, labels, and containment.

## Invocation

Each `model` verb **MUST** read a model file — the default `QUALITY.md` in the
current working directory, or a `[path]` argument when given — and **MUST NOT**
be aware of evaluation runs, run IDs, or run-folder layout. A run's
`model-snapshot.md` is itself a model file, reachable by path.

`model` **MUST NOT** write, create, repair, or delete files.

A `[path]` of `-` **MUST** fail with a usage error; `model` needs a filesystem
path.

## Canonical IDs and ordering

Every element a verb emits **MUST** carry its canonical qualified reference
string, identical to the form persisted in evaluation payloads. The root Area is
`area:root`. A Requirement is addressed by its declaring Area and name
(`requirement:<area>::<name>`) whether it is declared directly under the Area or
under one of the Area's Factors.

For the same model file and content, each verb's output **MUST** be byte-stable
under one documented element ordering: the rooted Area, then its Factors
(sub-Factors nested), then its Requirements, then its child Areas, recursively,
with siblings of each kind in lexicographic key order. No sort or ordering flag
is offered.

## Flags

Beyond the shared `[path]` and `--json`, each verb **MUST** accept only the flags
named for it below and **MUST NOT** offer write, output-redirection, dry-run,
quiet, or verbose flags.

Where a verb accepts `--area`, its value **MUST** be a canonical Area reference
(`area:<path>`), the same form `model` emits, so the value round-trips with
`list`/`get` output. A value that is not a canonical Area reference, or that does
not resolve to an Area in the model, is a usage error. A bare path is not
accepted.

## `model tree`

`model tree [path]` **MUST** render the model as a hierarchy — the rooted Area,
its Factors with sub-Factors nested, its Requirements, then child Areas
recursively. The human form is an indented tree; under `--json` it is a nested
structure where every node carries its canonical `id`, `label`, and element
`kind`, with nested children under `children`.

`model tree` **MUST** accept `--area <area-id>` to root the output at that Area's
subtree instead of the whole model.

`model tree` **MUST** accept `--depth <n>` to limit nesting, where `0` emits only
the rooted node. A negative `--depth` is a usage error.

## `model list`

`model list [path]` **MUST** emit a flat enumeration of model elements, each with
its canonical `id`, `kind` (`area` | `factor` | `requirement`), `label`, and
parent `id` (absent for the root Area). Under `--json` it is a JSON array of
those objects.

`model list` **MUST** accept `--type` restricting output to one or more element
kinds (`area`, `factor`, `requirement`), and `--area <area-id>` restricting
output to one Area's subtree; the two **MAY** be combined. An invalid `--type`
value is a usage error naming the allowed values.

With no filter, `model list` **MUST** enumerate every element of every kind in
the model.

## `model get`

`model get <id> [path]` **MUST** accept a canonical reference id as a required
positional argument and emit that element's `id`, `kind`, `label`, and its
structural relationships — for an Area, its immediate Factor IDs, Requirement
IDs, and child-Area IDs; for a Factor, its sub-Factor IDs and immediate
Requirement IDs — as an object under `--json` and a human-readable form
otherwise. It **MUST NOT** require a `--type` hint, since the id prefix carries
the kind.

If `<id>` does not resolve to an element in the model, `model get` **MUST** exit
`2`, name the unresolved id, and **SHOULD** suggest the closest matching element
ids.

## Output

Each verb **MUST** default to a human-readable form on stdout and **MUST** offer
`--json` for the structured form. Terminal styling **MUST NOT** change the words,
order, or facts of the output, and **MUST NOT** appear in `--json` or
non-terminal output. There is no `--format` flag.

## Exit status

`model` exits `0` on success.

`model` exits `2` on a malformed invocation: an unknown flag, an invalid `--type`
value, or an `--area`/`<id>` that does not resolve. An unresolved `--area`/`<id>`
is a usage error — the model is valid, the request named a non-element — not a
found-problem.

`model` exits `70` when the model file cannot be read or parsed; the error
**SHOULD** point at `lint` for diagnostics.

`model` **MUST NOT** return `1`; it reports no negative result of its own.
