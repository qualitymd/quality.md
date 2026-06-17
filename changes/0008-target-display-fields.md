---
type: Change
title: Describe targets with title and description
description: Let every target carry an optional title and a recommended description, lifting the root-only title restriction so the field is the one root-only key.
status: Draft
tags: [specification, schema, lint]
timestamp: 2026-06-17T00:00:00Z
---

# Describe targets with title and description

A **Change** to give a target a human-readable identity. Detail lives in the
children:

- [Functional spec](0008-target-display-fields/spec.md) — what the change must do.

A design doc follows at **Design** for the one real decision (how the
`misplaced-root-key` rule narrows to `ratingScale` only).

## Motivation

A target's only human-facing label today is its **map key**, which doubles as
identifier. There is no place to record what a target *is* — neither a display
name ("ACME Products API") nor a descriptive line ("Functional specifications
for the ACME Products API"). The spec makes the omission sharper than it should
be: it states "the model root is itself the apex target" and that every target
is "another target of the same shape," yet `title` is reserved to the root and
no target — root or not — can carry a `description`. Reports and a skill scoping
an evaluation have nothing but the slug to go on.

The fix keeps the recursive-node model honest: a target describes itself with a
`title` (display name) and a `description` (what it is / its scope), exactly as
the apex target and a factor already do, leaving `ratingScale` as the one
genuinely root-only key.

## Scope

Covered:

- Non-root targets MAY declare `title` (a display name); the root-only
  restriction on `title` is lifted.
- Every target — the root/apex target included — MAY declare a `description`.
- `ratingScale` remains the sole root-only key.
- The structural schema, the `misplaced-root-key` lint rule, and their tests are
  brought into line, with no other behavior change.

Deferred: any use of the new fields by `init` scaffolding, `spec` rendering, or
report formatting — this change defines the fields and their validation only.
Evaluation roll-up semantics are untouched.

## Affected specs & docs

- [ ] [`SPECIFICATION.md`](../SPECIFICATION.md) — Model schema gains
      `description`; Target schema gains `title` and `description`; the "`title` and
      `ratingScale` … declared only on the model root" sentence narrows to
      `ratingScale`; add the Target **Title** and **Description** prose.
- [ ] [`specs/cli/lint.md`](../specs/cli/lint.md) — the `misplaced-root-key` row
      changes from "`title` or `ratingScale`" to `ratingScale` only.
- [ ] `internal/schema/schema.go` — add `description` to `Model`; add `title` and
      `description` to `Target`.
- [ ] `internal/schema/schema_test.go` — update the spec↔schema consistency test
      for the new properties.
- [ ] `internal/lint` (`result.go`, `rules.go`, `rules_test.go`) — narrow the
      `misplaced-root-key` rule and its description to `ratingScale`; the existing
      "nested target title" test case becomes a valid model, not a finding.

## Status

`Draft`. See the [status lifecycle](index.md#status-lifecycle).
