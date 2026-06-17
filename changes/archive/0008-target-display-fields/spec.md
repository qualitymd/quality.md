---
type: Functional Specification
title: target title and description — functional spec
description: The target title and description fields, the reframing of the root as a Model composed of ratingScale plus Target properties, and the schema and lint conformance the change must meet.
tags: [specification, schema, lint]
timestamp: 2026-06-17T00:00:00Z
---

# target title and description — functional spec

Companion to the [Describe targets with title and description](../0008-target-display-fields.md)
change. This spec states the delta the durable
[`SPECIFICATION.md`](../../../SPECIFICATION.md) and the structural schema must
absorb so a target can carry a human-readable identity, reframes the root node
as a Model composed of `ratingScale` plus the Target properties, and states the
lint conformance that follows.

The key words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, and **MAY** are
to be interpreted as described in IETF RFC 2119.

## Scope

Covered: the `title` and `description` keys on a target, their presence and
shape; the reframing of the Model/Target relationship; and the
structural-schema and `misplaced-root-key` conformance (with tests).

Deferred: any consumption of the fields by `init`, `spec`, or report rendering;
and any change to evaluation, roll-up, or source semantics. The evaluation
prose's "root target" wording (the top of the target tree) is unaffected.

## Requirements — the specification

### Target fields

- A target **MAY** declare a `title`: a human-readable display name for the
  entity it evaluates (e.g. "ACME Products API"). A target **SHOULD** declare a
  `title` — it is **RECOMMENDED** for readable reports — and a declared `title`
  overrides the map key for display; the map key remains the identifier.
- A target **MAY** declare a `description`: a concise statement of what the
  target *is* — the entity or scope it covers (e.g. "Functional specifications
  for the ACME Products API"). `description` is **OPTIONAL**.
- `title` and `description` **MUST** each be a single scalar string when present.
- A `description` **SHOULD NOT** restate or stand in for the target's factors or
  requirements; it fixes *what is evaluated*, not the expectations, which remain
  in `factors` and `requirements`.

### The Model / Target relationship

- The specification **MUST** describe the root node as a **Model**: the
  model-wide `ratingScale` together with all the properties of a
  [Target](../../../SPECIFICATION.md) — `title`, `description`, `factors`,
  `requirements`, `targets`, and `source`. The Model root therefore tops the
  target tree because it carries the target properties.
- `ratingScale` **MUST** be the only property unique to the Model; a **Target**
  **MUST NOT** declare `ratingScale`. The specification **MUST** state this as a
  type distinction — `ratingScale` is a Model property a Target does not have —
  rather than as a list of keys "a non-root target MUST NOT declare."
- The Target schema in `SPECIFICATION.md` **MUST** show `title` (Recommended) and
  `description` (Optional) alongside `factors`, `requirements`, `targets`, and
  `source`; the Model schema **MUST** show `description` (Optional) and retain
  its `ratingScale` and RECOMMENDED `title`.
- The change **MUST** replace the "apex target / a non-root target MUST NOT
  declare either" prose accordingly, add Target **Title** and **Description**
  prose mirroring the Factor description guidance, and **MUST NOT** contradict
  the existing recursive-node, source-inheritance, or roll-up text.

## Requirements — the conformance

- The structural schema (`internal/schema`) **MUST** add `description` to the
  model node and `title` and `description` to the target node, with the presence
  and scalar shape above, so the linter derives validation from the single
  source of truth per the
  [schema-source change](../../archive/0005-schema-source-of-truth.md). It **MAY**
  build the Model node's target properties from the Target node to avoid drift;
  the [design doc](design.md) decides.
- The spec↔schema consistency test **MUST** be updated to reflect the new
  properties and continue to pass.
- The `misplaced-root-key` rule **MUST** flag only `ratingScale` on a target;
  a target that declares `title` **MUST NOT** be flagged, and a target declaring
  `description` **MUST** be accepted. The rule's catalog description and the
  [`lint` sub-spec](../../../specs/cli/lint.md) row **MUST** be updated from
  "`title` or `ratingScale`" to `ratingScale`.
- Tests **MUST** pin the new behavior: `title`/`description` on a nested target
  are valid, and `ratingScale` on a nested target still produces
  `misplaced-root-key`. No other rule, output, or feature behavior changes.

## Done criteria

- [`SPECIFICATION.md`](../../../SPECIFICATION.md) shows `title` and `description`
  on the Target schema and `description` on the Model schema, and frames the root
  as a Model = `ratingScale` + Target properties, with `ratingScale` named the
  sole Model-only property.
- [`specs/cli/lint.md`](../../../specs/cli/lint.md)'s `misplaced-root-key` row
  names only `ratingScale`.
- `internal/schema` declares the new properties and the consistency test passes.
- `internal/lint` flags only `ratingScale` on a target, with tests for the
  accepted `title`/`description` cases and the still-rejected `ratingScale` case.
- The change is moved through the lifecycle and archived per the
  [changes process](../../index.md#status-lifecycle).
