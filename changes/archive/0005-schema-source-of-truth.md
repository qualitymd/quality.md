---
type: Change Case
title: Single source of truth for the structural schema
description: Extract the QUALITY.md structural schema into one authoritative definition the linter derives from, so the format spec and the linter can't drift.
status: Done
tags: [lint, spec, schema]
timestamp: 2026-06-17T00:00:00Z
---

# Single source of truth for the structural schema

This change extracts the QUALITY.md _structural schema_ — the valid frontmatter
keys at each node, which are required, and the rating-scale shape — into a single
authoritative definition the linter derives from, rather than re-encoding it in
rule code. The detail lives in its children:

- [Functional spec](0005-schema-source-of-truth/spec.md) — what the change must do.
- [Design doc](0005-schema-source-of-truth/design.md) — how the linter derives
  structural checks from one typed schema declaration.

## Motivation

The structural schema is currently expressed twice: implicitly in the lint rule
code (`internal/lint/rules.go` hard-codes valid keys, required properties, and the
rating-scale shape) and again in prose in
[`SPECIFICATION.md`](../../SPECIFICATION.md) and the [`specs/`](../../specs/index.md)
bundle. These can drift — a key added to the documented format may go unrecognized
by the linter, or a rule may reject something the spec allows. design.md's linter
(reviewed alongside ours) avoids this by deriving its structural rules from one
schema artifact. A single authoritative definition the linter consumes keeps the
two in lockstep and makes "what is valid" answerable in one place.

## Scope

**Covered.** One authoritative definition of the structural schema, the linter's
structural validation derived from it (no second hand-maintained valid-key list),
and the durable format spec reconciled to stay consistent with it. The refactor is
behavior-preserving for documents that already agree with the spec.

**Deferred.** Generating prose specification or user docs from the source; a
runtime configuration surface (rule selection, severity overrides, suppression —
already deferred by [`specs/cli/lint.md`](../../specs/cli/lint.md)); "did you mean?"
typo suggestions on unknown keys (a separate candidate change); and any change to
the rule catalog's severities or membership.

## Affected specs & docs

Decided up front; reconciled before this change reaches `Done`. Direction (to be
detailed in the design doc): the schema source is a **typed Go declaration** the
linter derives from directly — not an embedded data file or a `specs/` concept —
and spec/linter consistency is enforced by a test that checks
[`SPECIFICATION.md`](../../SPECIFICATION.md) against it, rather than by generating
docs. This may add an internal package to the implementation; it does not add to
the durable specs/docs below.

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) — point the structural-schema
      prose at the authoritative definition, reconciling any drift in favor of the
      documented format.
- [x] [`specs/cli/lint.md`](../../specs/cli/lint.md) — note that structural
      validation derives from the single schema definition.

## Status

`Done`. Implemented and archived after adding the typed structural schema
declaration, deriving lint structural checks from it, reconciling the public
format prose, adding the lint sub-spec note, and pinning spec/schema consistency
with tests.
