---
type: Design Doc
title: Self-contained per-kind data schema — design
description: How `evaluation data schema <kind>` is made self-contained, and why the generic flattener scoped in Draft is unnecessary.
tags: [cli, evaluation, schema, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Self-contained per-kind data schema — design

## Context

Answers the [functional spec](spec.md) for change case
[0130](../0130-self-contained-data-schema.md). The spec requires
`evaluation data schema <kind>` to emit a schema whose required fields and enum
value sets are legible at the document root, derived from the same typed contract
that backs `data set` validation, while keeping `qualitymd schema` and
`evaluation data schema` as separate commands.

## Approach

### The flatten is shedding an envelope, not resolving references

`evaluationDataSchemaDoc(kind)`
([`internal/evaluation/data_contract.go`](../../../internal/evaluation/data_contract.go))
today builds one document for both forms: a `$defs` map holding all ten kind
schemas, plus either a top-level `$ref` to the requested kind (single-kind) or a
`oneOf` over every kind (no-arg).

The load-bearing finding from reading the generator: **the only `$ref`/`$defs` in
the entire surface is that top-level envelope.** `schemaForObject` and
`schemaForField` emit no references — nested objects (findings, rating drivers,
routine refs) are inlined recursively, and each field's `enum` and the object's
`required` sit directly in the object schema. So every `$defs[<kind>]` entry is
_already_ a self-contained schema. The legibility problem was never reference
depth; it was that the single-kind request wrapped a self-contained object in a
ten-kind envelope and pointed a `$ref` at it.

The change is therefore narrow: when `kind != ""`, return
`schemaForObject(contract.Object)` merged with the document header
(`$schema`, `$id`) directly — no `$defs`, no `$ref`. The requested kind's
`required` array, its `properties`, each closed field's `enum` (including the
self-identifying `kind` enum pinned to the requested kind), and
`additionalProperties: false` all surface at the root. The no-arg form is
unchanged: it keeps `$defs` + `oneOf`, the correct shape for a multi-kind
document.

```
single-kind, before                  single-kind, after
{                                     {
  "$schema": ..., "$id": ...,           "$schema": ..., "$id": ".../<kind>",
  "$defs": { …all 10 kinds… },          "type": "object",
  "$ref": "#/$defs/<kind>"              "required": [ … ],
}                                       "properties": { "status": {"enum": […]}, … },
                                        "additionalProperties": false
                                      }
```

The single-kind document takes a kind-qualified `$id` (base id + `/<kind>`) so the
two emissions are distinguishable; `$schema` (the dialect) is unchanged.

### Presentation reuses the existing presenter

`writeSchema`
([`internal/cli/schema.go`](../../../internal/cli/schema.go)) already implements
the verbatim-on-plain / highlight-and-page-on-terminal behavior the spec's
terminal-rendering group wants, and it lives in `package cli` — the same package
as the `data schema` command. The optional alignment is a one-line swap in
`newEvaluationDataSchemaCmd`
([`internal/cli/evaluation.go`](../../../internal/cli/evaluation.go)): call
`writeSchema(cmd.OutOrStdout(), raw)` instead of `cmd.OutOrStdout().Write(raw)`.
No extraction, no new package, and `qualitymd schema` is untouched.

### Source-of-truth invariant holds for free

Because both forms are still produced from `dataContracts` via `schemaForObject`,
the emitted `enum` sets and `required` arrays remain identical to what
`validateField`/`validateObject` enforce. There is no second source to keep in
sync — the spec's "derived from one typed contract" requirement is preserved by
construction, not by added code.

### Skill edits

Runtime content (`skills/quality/SKILL.md`, `workflows/evaluate.md`,
`resources/cli-workflow-conventions.md`) and the durable skill spec
(`specs/skills/quality-skill/quality-skill.md`) are sharpened to name
`data schema <kind>` as the source for required fields and enum values, with
`data example <kind>` a concrete instance only — per the spec's skill-guidance
requirement. No structural skill change; the schema becoming legible is what makes
that guidance followable.

## Spec response

- **Self-contained single-kind schema** — satisfied by the envelope-shedding
  branch; required + enums surface at the root because the kind object is already
  reference-free.
- **Same constraints / one typed contract** — preserved by construction: the same
  `schemaForObject` output, minus the wrapper, is what validation also reads.
- **No-arg full surface** — explicitly the unchanged branch.
- **Command surface** — no command is added, renamed, or merged; only the body of
  `evaluationDataSchemaDoc`'s single-kind branch changes.
- **Terminal rendering** — delegated wholesale to the existing `writeSchema`,
  which already meets the verbatim-on-plain guarantee.

## Alternatives

- **A generic `$ref`-resolving flattener in `internal/schema`** — scoped in the
  Draft footprint on the assumption that single-kind output needed reference
  resolution. Reading the generator showed the kind objects carry no inter-kind
  `$ref`, so a resolver would be unused generality solving a problem this surface
  doesn't have. Rejected; the parent's **Affected artifacts** is corrected to drop
  `internal/schema/` and the `schema.go` extraction.
- **Unified `qualitymd schema <kind>` with the model as default** — rejected in the
  [change case](../0130-self-contained-data-schema.md#scope): the model frontmatter
  and evaluation payloads are different artifact families across a pipeline
  boundary, not peers in one kind taxonomy.
- **`data describe` / annotating `data example`** — rejected: a third surface
  re-deriving constraints the schema already encodes (drift), or breaking
  `data example` as pure pipeable JSON.

## Trade-offs & risks

- **Two emission shapes.** Single-kind (object-at-root) and no-arg
  (`$defs` + `oneOf`) now differ structurally. This is intended — a consumer
  asking for one kind wants the object, not a union — but tests must cover both
  shapes so a future refactor doesn't accidentally re-wrap the single-kind form.
- **`$id` for the single-kind document.** Distinct content under a distinct
  kind-qualified `$id` avoids two documents claiming one identifier. Low-stakes:
  output is read/validated ad hoc, not registered in a shared schema store.
- **Validation equivalence is asserted, not refactored.** The flattened schema
  must validate the same payloads as `data set`. A regression test should validate
  a known-good and a known-bad payload for one kind against the emitted single-kind
  schema to lock the equivalence the spec requires.

## Open questions

- ~~Should the optional terminal-rendering alignment ship in this case or split
  out?~~ **Resolved:** kept in this case — one-line `writeSchema` swap plus a
  test, no benefit to splitting.
