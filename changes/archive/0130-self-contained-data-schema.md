---
type: Change Case
title: Self-contained per-kind data schema
description: Make `evaluation data schema <kind>` self-contained so required fields and enum value sets are directly legible, and align the skill to read it as the constraint source.
status: Done
tags: [cli, evaluation, schema, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Self-contained per-kind data schema

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0130-self-contained-data-schema/spec.md) — what the case must do.
- [Design doc](0130-self-contained-data-schema/design.md) — how it's built, and why.

## Motivation

A full `evaluate` run against an external monorepo (skill/CLI 0.18.0) recorded a
recurring friction in its
[feedback log](../../changes/archive/0073-evaluation-feedback-log.md): the closed
enum value sets and required fields for evaluation payload kinds were discovered
only by triggering `data set --dry-run` validation failures, not by reading the
schema. Two concrete misses — the analysis `status` enum
(`analyzed|empty|not_analyzed|blocked`, of which the example shows only `empty`)
and the `area:<id>` reference form for `areaId` — both surfaced as dry-run
diagnostics rather than from `data schema`.

The root cause is one unreadable surface, not two problems. `qualitymd
evaluation data schema <kind>` already carries `enum` arrays and `required`, but
even for a single requested kind it emits the **full** multi-kind `$defs` map
with a top-level `$ref` indirection. Authors and agents then fall back to
`data example` (which shows one value per enum) as the practical source of truth,
and learn the closed sets by failing validation. The schema should be the legible
constraint source it is already meant to be.

## Scope

Covered:

- `evaluation data schema <kind>` emits a **self-contained** schema scoped to the
  requested kind — its `required` and each field's `enum` are directly legible at
  the document root, without dereferencing `$ref` into a separate `$defs` block.
- The constraint surface stays derived from the one typed contract that backs
  `data set` validation and `data example` (no second source of truth; flattening
  is presentational, not a constraint change).
- The `/quality` skill is sharpened to read `data schema <kind>` as the source for
  required fields and allowed enum values, and to treat `data example <kind>` as a
  concrete instance only.
- Optional consistency win: `evaluation data schema` reuses the same terminal
  rendering (verbatim-on-plain, optional highlight/page on a terminal) that
  `qualitymd schema` already applies.

Deferred / **rejected** (recorded so they are not re-litigated):

- **Not** collapsing the two commands into a unified `qualitymd schema <kind>`
  with the model as a default kind. The QUALITY.md frontmatter (authored input)
  and the evaluation payloads (run output) are different artifact families across
  a pipeline boundary, not peers in one taxonomy — the forced special-case default
  is the tell. They share _machinery_, not a _command surface_ or a `kinds`
  registry.
- **No** new `data describe` command and **no** annotation of `data example`
  output. Both would create a third surface re-deriving constraints the schema
  already encodes (drift risk) or break `data example` as pure pipeable JSON.

## Affected artifacts

Derived by sweeping the repo for `data schema`, `EvaluationDataSchema`,
`evaluationDataSchemaDoc`, `writeSchema`, and the skill's artifact-contract
guidance.

**Code**

- [x] `internal/evaluation/data_contract.go` — `evaluationDataSchemaDoc`
      single-kind branch emits the kind's object schema with the document header
      directly (no `$defs`/`$ref` envelope). The no-arg branch is unchanged.
- [x] `internal/cli/evaluation.go` — `data schema` calls the existing
      `writeSchema` presenter instead of a raw write (optional rendering
      alignment).
- [x] Tests: `internal/evaluation/*_test.go`, `internal/cli/schema_test.go`,
      `internal/cli/evaluation_test.go`.

> The [design](0130-self-contained-data-schema/design.md) found the kind objects
> carry no inter-kind `$ref`, so the generic `internal/schema` flattener and the
> `internal/cli/schema.go` extraction scoped in Draft are unnecessary and dropped
> from this footprint.

**Durable specs** (substance in the [functional spec](0130-self-contained-data-schema/spec.md))

- [x] `specs/cli/evaluation-data.md` — strengthen the `data schema <kind>`
      requirement from "rooted at that kind" to self-contained legibility; note
      the optional terminal-rendering alignment.
- [x] `specs/skills/quality-skill/quality-skill.md` — the skill's _durable spec_
      (distinct from its runtime content below): sharpen the Evaluation
      payload-discovery `MUST` so required fields and enum values are read from
      `data schema <kind>` and `data example <kind>` is a concrete instance only.

**Durable docs**

- [x] `skills/quality/SKILL.md` — Artifact Contract: name `data schema <kind>` as
      the required-fields/enum source; `data example` as a concrete instance only.
- [x] `skills/quality/workflows/evaluate.md` — the shape-discovery step: same
      division of labor.
- [x] `skills/quality/resources/cli-workflow-conventions.md` — note
      `data schema <kind>` is per-kind and self-contained (optional sharpening).
- [x] `skills/quality/log.md` — append a runtime-history entry when landing.
- [x] `CHANGELOG.md` — entry at release time.

No impact: `README.md` (CLI is not foregrounded), `SPECIFICATION.md` (format
model unchanged), install/scaffold files.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Implemented and
archived. `evaluation data schema <kind>` now emits a self-contained kind schema,
`data schema` uses the shared schema presenter, durable specs and runtime skill
guidance name schema as the required-fields/enum source, and focused
`internal/evaluation` plus `internal/cli` tests pass.
