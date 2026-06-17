---
type: Functional Specification
title: target title and description — functional spec
description: The target title and description fields, the narrowing of the root-only constraint to ratingScale, and the schema and lint conformance the change must meet.
tags: [specification, schema, lint]
timestamp: 2026-06-17T00:00:00Z
---

# target title and description — functional spec

Companion to the [Describe targets with title and description](../0008-target-display-fields.md)
change. This spec states the delta the durable
[`SPECIFICATION.md`](../../../SPECIFICATION.md) and the structural schema must
absorb so a target can carry a human-readable identity, and the lint conformance
that follows.

The key words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, and **MAY** are
to be interpreted as described in IETF RFC 2119.

## Scope

Covered: the `title` and `description` keys on a target, their presence and
shape, the narrowing of the root-only constraint to `ratingScale` alone, and the
structural-schema and `misplaced-root-key` conformance (with tests).

Deferred: any consumption of the fields by `init`, `spec`, or report rendering;
and any change to evaluation, roll-up, or source semantics.

## Requirements — the specification

### Target fields

- A target **MAY** declare a `title`: a human-readable display name for the
  entity it evaluates (e.g. "ACME Products API"). `title` is **OPTIONAL** on a
  target, because the target's map key already names it; a declared `title`
  overrides the map key for display.
- A target **MAY** declare a `description`: a concise statement of what the
  target *is* — the entity or scope it covers (e.g. "Functional specifications
  for the ACME Products API"). A target **SHOULD** declare a `description`, except
  where a pure grouping target or a self-evident map key makes one redundant.
- `title` and `description` **MUST** each be a single scalar string when present.
- A `description` **SHOULD NOT** restate or stand in for the target's factors or
  requirements; it fixes *what is evaluated*, not the expectations, which remain
  in `factors` and `requirements`.

### The model root (apex target)

- The model root **MUST** also accept `description`, with the same meaning and
  shape, since the root is the apex target. The root's existing `title`
  (RECOMMENDED) is unchanged.

### The root-only constraint

- The only key a non-root target **MUST NOT** declare is `ratingScale`. The
  specification text that today reserves both `title` and `ratingScale` to the
  model root **MUST** be narrowed so that `ratingScale` is the sole root-only
  key, and **MUST** continue to state that a non-root target **MUST NOT** declare
  `ratingScale`.
- The Target schema in `SPECIFICATION.md` **MUST** show `title` (Optional) and
  `description` (Recommended) alongside `factors`, `requirements`, `targets`, and
  `source`, and the Model schema **MUST** show `description` (Recommended).
- The change **MUST** add Target **Title** and **Description** prose mirroring the
  Factor description guidance, and **MUST NOT** contradict the existing
  recursive-node, source-inheritance, or roll-up text.

## Requirements — the conformance

- The structural schema (`internal/schema`) **MUST** add `description` to the
  model node and `title` and `description` to the target node, with the presence
  and scalar shape above, so the linter derives validation from the single
  source of truth per the
  [schema-source change](../../archive/0005-schema-source-of-truth.md).
- The spec↔schema consistency test **MUST** be updated to reflect the new
  properties and continue to pass.
- The `misplaced-root-key` rule **MUST** flag only `ratingScale` on a non-root
  target; a non-root target that declares `title` **MUST NOT** be flagged, and a
  target declaring `description` **MUST** be accepted. The rule's catalog
  description and the [`lint` sub-spec](../../../specs/cli/lint.md) row **MUST**
  be updated from "`title` or `ratingScale`" to `ratingScale`.
- Tests **MUST** pin the new behavior: `title`/`description` on a nested target
  are valid, and `ratingScale` on a nested target still produces
  `misplaced-root-key`. No other rule, output, or feature behavior changes.

## Done criteria

- [`SPECIFICATION.md`](../../../SPECIFICATION.md) shows `title` and `description`
  on the Target schema and `description` on the Model schema, narrows the
  root-only constraint to `ratingScale`, and carries the Target Title/Description
  prose.
- [`specs/cli/lint.md`](../../../specs/cli/lint.md)'s `misplaced-root-key` row
  names only `ratingScale`.
- `internal/schema` declares the new properties and the consistency test passes.
- `internal/lint` flags only `ratingScale` on a non-root target, with tests for
  the accepted `title`/`description` cases and the still-rejected `ratingScale`
  case.
- The change is moved through the lifecycle and archived per the
  [changes process](../../index.md#status-lifecycle).
