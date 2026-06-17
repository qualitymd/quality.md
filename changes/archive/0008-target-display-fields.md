---
type: Change
title: Describe targets with title and description
description: Let every target carry title and description, and reframe the root as a Model — Target properties plus the model-wide ratingScale — so ratingScale is the one root-only key by type, not by prohibition.
status: Done
tags: [specification, schema, lint]
timestamp: 2026-06-17T00:00:00Z
---

# Describe targets with title and description

A **Change** to give a target a human-readable identity, and to reframe the
root node so the difference between it and a nested target is a clean type
distinction rather than a list of prohibited keys. Detail lives in the children:

- [Functional spec](0008-target-display-fields/spec.md) — what the change must do.
- [Design doc](0008-target-display-fields/design.md) — how it lands, and why.

## Motivation

A target's only human-facing label today is its **map key**, which doubles as
identifier. There is no place to record what a target *is* — neither a display
name ("ACME Products API") nor a descriptive line ("Functional specifications
for the ACME Products API"). The spec makes the omission sharper than it should
be: it states "the model root is itself the apex target" and that every target
is "another target of the same shape," yet `title` is reserved to the root and
no target — root or not — can carry a `description`.

Fixing the fields exposes a second, framing problem. The spec explains the root
as an *apex target* that "shares the structure of the model root but for two
keys … which are declared only on the model root: a non-root target MUST NOT
declare either." Once `title` and `description` belong to every target, the only
remaining difference is `ratingScale`, and the "non-root target MUST NOT" framing
reads as an awkward prohibition. It is clearer — and matches how
[`internal/schema`](../../internal/schema) already models it, as two distinct
`Node`s — to say the root node is a **Model**: the model-wide `ratingScale`
together with all the properties of a **Target**. The root then tops the target
tree because it carries the target properties, and `ratingScale` is simply a
Model property a Target does not have. The constraint falls out of the types
instead of being asserted as a rule.

## Scope

Covered:

- Targets gain a `title` (display name) and a `description` (what the target is).
- The root is reframed as a **Model** = the model-wide `ratingScale` + every
  **Target** property; `ratingScale` is the only property unique to the Model.
- The "apex target / non-root target MUST NOT" prose is replaced by that type
  framing, with `ratingScale` named as the one Model-only property.
- The structural schema, the `misplaced-root-key` lint rule, and their tests are
  brought into line, with no other behavior change.

Deferred: any use of the new fields by `init` scaffolding, `spec` rendering, or
report formatting — this change defines the fields and their validation only.
Evaluation roll-up and the "root target" terminology in the evaluation prose
(the top of the target tree) are untouched.

## Affected specs & docs

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) — Target schema gains `title`
      (Recommended) and `description` (Optional) with their prose; the Model schema
      gains `description`; the Target intro and the "shares the structure … but for
      two keys … a non-root target MUST NOT declare either" sentence are reframed as
      *Model = Target properties + `ratingScale`*, naming `ratingScale` the sole
      Model-only property.
- [x] [`specs/cli/lint.md`](../../specs/cli/lint.md) — the `misplaced-root-key` row
      changes from "`title` or `ratingScale`" to `ratingScale` only.
- [x] [`internal/schema/schema.go`](../../internal/schema/schema.go) — add `description` to `Model`; add `title` and
      `description` to `Target` (the design doc decides whether `Model` derives its
      target properties from `Target`).
- [x] [`internal/schema/schema_test.go`](../../internal/schema/schema_test.go) — update the spec↔schema consistency test
      for the new properties.
- [x] [`internal/lint`](../../internal/lint) (`result.go`, `rules.go`, `rules_test.go`) — narrow the
      `misplaced-root-key` rule and its description to `ratingScale`; the existing
      "nested target title" test case becomes a valid model, not a finding.

## Status

`Done`. Implemented and archived after adding target display fields to the
schema, updating the durable format and lint specs, narrowing
`misplaced-root-key` to `ratingScale`, adding focused tests for nested target
`title`/`description`, and passing `go test ./...`.
