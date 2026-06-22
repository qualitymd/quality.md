---
type: Functional Specification
title: qualitymd lint
description: Validate a QUALITY.md file's structure against the format spec.
tags: [cli, command, lint]
timestamp: 2026-06-17T00:00:00Z
---

# qualitymd lint

`lint` inherits the cross-cutting CLI contract — invocation, global flags, output
formats, exit codes, and agent accessibility — from the [CLI spec](../cli.md).
This file specifies only what is particular to `lint`.

`qualitymd lint` validates a QUALITY.md file's structure against the
[format specification](../../SPECIFICATION.md), fast and deterministically,
exiting non-zero on errors so it drops into CI.

**Boundary.** `lint` checks *format conformance* — whether the file is a valid
QUALITY.md file — only. It does not assess whether the model is a *good* quality
model; that judgment lives in the evaluation skills, not the deterministic CLI.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", "RECOMMENDED", and "MAY" are to be
interpreted as described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Scope

Covered: the rule boundary for mechanical format validation, rule metadata,
finding messages, finding locations, human and JSON output, deterministic
ordering, in-place repair of fixable findings, and the initial rule set.

Deferred: suppression directives, rule selection or severity overrides, emitting
a rule catalog from `lint` itself, and repair output modes other than in-place
writes (for example patch output or emitting a full rewritten file to stdout).

## Flags

`lint` inherits the cross-cutting flags from the [CLI spec](../cli.md), including
`--json`, and adds one command-specific flag:

- `--fix` — apply every fixable finding that can be repaired deterministically,
  writing the repaired `QUALITY.md` back to the same path.

`--fix` **MAY** be combined with `--json`; the JSON output then reports the
post-repair findings and the repairs that were applied.

## Rule System

The lint rule-system, rule-authoring guidance, and rule catalog live in
[`qualitymd lint rules`](lint-rules.md). `lint` rules remain mechanically
grounded in the format specification, deterministic, format-conformance focused,
and self-contained within the linted file and declared cross-references.

### Repair behavior

When `--fix` is passed, `lint` **MUST** apply every emitted finding whose
`fixable` value is `true`, then lint the repaired file again and report the
post-repair findings. Non-fixable findings remain findings; they are never
silently changed or suppressed.

Repairs are **transactional per file**:

- `lint` **MUST** compute repairs from the original parsed document.
- If two repairs cannot be applied together, `lint` **MUST** leave the file
  unchanged and report a repair failure rather than applying a partial set.
- If a write fails, `lint` **MUST** report the failure and **MUST NOT** report the
  file as repaired.
- Applying `--fix` twice **MUST** have the same effect as applying it once.

When repair fails before a post-repair lint result exists, `lint` **MUST** exit
non-zero and report the failure through the CLI's error-reporting path. It
**MUST NOT** emit a successful lint result for a file it did not repair.

In-place writes **MUST** preserve authored content outside the repaired
frontmatter nodes:

- Markdown body content **MUST** be preserved byte-for-byte.
- YAML map order, comments, scalar style, and whitespace outside repaired nodes
  should be preserved where the parser and emitter make that possible.
- Unrelated YAML keys **MUST NOT** be reordered or rewritten unless the rewrite is
  necessary to apply the repair deterministically.
- The write should be atomic from the caller's perspective: write a complete
  replacement and then replace the area path, rather than truncating the
  original before the replacement is ready.
- To avoid ambiguous replacement behavior, `lint --fix` should refuse to
  repair a linted path that is a symbolic link until symlink write semantics are
  specified.

This phase defines only in-place repair through `--fix`. `lint` **MUST NOT** emit
patches or full rewritten files as alternate repair output modes in this phase.

## Findings and output

The lint finding schema, JSON result shape, repair-result objects, locations,
and human-output requirements live in
[`qualitymd lint output`](lint-output.md). Human-readable output and JSON output
report the same findings for the same input; `--json` changes only the format.

### Exit status

Without `--fix`, `lint` exits non-zero when it emits one or more `error`
findings. Warnings and info findings do not affect the exit code.

With `--fix`, `lint` exits non-zero when repair fails or when the post-repair
lint result still contains one or more `error` findings. A run that fixes all
errors and leaves only warnings exits zero.

### Ordering and blocking

Findings **MUST** be emitted in deterministic order:

1. Earlier source position first, when both findings have source positions.
2. Otherwise, shallower `modelPath` first, then lexicographic comparison of path
   segments with numeric indexes ordered numerically.
3. For the same location, `error` before `warning` before `info`.
4. For the same location and severity, lexicographic `ruleId`.

A malformed or structurally invalid parent can block rules that depend on
that parent's parsed shape. For example, absent frontmatter or invalid YAML emits
`invalid-frontmatter` and prevents model-level rules from running; a malformed
`ratingScale` may prevent per-level checks from running. A blocked downstream
rule **MUST NOT** emit a speculative finding.
