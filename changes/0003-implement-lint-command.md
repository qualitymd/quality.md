---
type: Change
title: Implement the lint command
description: Build `qualitymd lint` to validate QUALITY.md files according to the completed lint sub-spec.
status: Design
tags: [cli, lint]
timestamp: 2026-06-17T00:00:00Z
---

# Implement the lint command

[`qualitymd lint`](../specs/cli/lint.md) now has a complete durable functional
spec for its command-specific behavior. This change tracks the implementation
work needed to make that command real. The detail lives in its children:

- [Functional spec](0003-implement-lint-command/spec.md) - the implementation
  delta, deferring command behavior to the durable lint sub-spec.
- [Design doc](0003-implement-lint-command/design.md) - how `lint` is built on a
  shared parsed model and traversal layer.

## Motivation

`lint` is the mechanical gate for every `QUALITY.md` workflow: it tells authors,
CI, and evaluation skills whether a file is structurally valid before any
judgment-based work begins. The durable lint sub-spec now defines the rule set,
finding schema, location contract, output expectations, ordering, and deferred
features. Implementing it unblocks CI use, scaffold verification, and the
judgment skills that must stop on structurally invalid models.

## Scope

Covered: implement `qualitymd lint` according to
[`specs/cli/lint.md`](../specs/cli/lint.md), including the initial error and
warning rules, deterministic findings, human output, `--json` output, `--fix`
in-place repairs, and tests that pin rule behavior, repair behavior, and output
shape.

Deferred: suppression directives, rule selection, severity overrides,
patch/full-file repair output modes, and a lint-emitted rule catalog remain
deferred by the durable lint sub-spec. Cross-cutting CLI behavior remains owned
by the [CLI spec](../specs/cli.md).

## Affected specs & docs

Updated before this change reaches **Done**:

- [x] [`specs/cli/lint.md`](../specs/cli/lint.md) - scope `--fix` into the
      durable lint contract, including in-place repair behavior and repair
      reporting.
- [ ] [`README.md`](../README.md) - drop the *(planned)* marker on `lint` and
      update the implementation-status note once the command ships.

If implementation reveals another functional gap, update the durable sub-spec
before marking this change done.

## Status

`Design`. See the [status lifecycle](index.md#status-lifecycle). The design
settles the implementation approach and uses `lint [path]`, defaulting to
`QUALITY.md`, as the minimum invocation shape for this change while the parent
CLI spec keeps the broader file/stdin convention open.
