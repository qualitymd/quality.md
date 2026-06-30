---
type: Functional Specification
title: String model-identity fields in evaluation data
description: Persist Evaluation routine and report identity fields as single canonical qualified model-reference strings, keep the `*Id` names, and reserve `root` as an Area name so the string form is lossless.
status: Draft
tags: [evaluation, records, references, schema]
timestamp: 2026-06-26T00:00:00Z
---

# String model-identity fields in evaluation data

This spec governs how persisted Evaluation routine and report JSON encodes the
identity of an Area, Factor, Requirement, or Rating Level. It is the delta for
[0120](../0120-string-model-identity-fields.md). The durable contracts it lands
in are
[`specs/evaluation/records/json-conventions.md`](../../specs/evaluation/records/json-conventions.md)
(normative — the Identity And References rule it rewrites) and
[`SPECIFICATION.md`](../../SPECIFICATION.md) (normative — it defines the
qualified model references and the strict name grammar this change persists, and
gains the `root` Area-name reservation).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The big-picture _why_ is recorded in the
[parent concept's Motivation](../0120-string-model-identity-fields.md#motivation):
one identity concept is persisted in three physical shapes (a string
`ratingLevelId`, an array `areaId`, and object `factorId` / `requirementId`),
and [`SPECIFICATION.md`](../../SPECIFICATION.md) already defines a single
canonical qualified model-reference string for each kind whose grammar makes it a
lossless encoding of the composite identity. This change collapses the persisted
identity fields onto that one string shape.

Two load-bearing facts shape the requirements below:

- **Losslessness depends on an unambiguous rendering.** The only ambiguity in the
  qualified-reference grammar is the root-Area token: `area:root` renders the
  empty root path, but the name grammar also currently admits a child Area named
  `root`, which would render to the same string. Reserving `root` as an Area name
  closes that gap; `/` and `:` are already excluded from names.
- **This reverses a prior decision.** `json-conventions.md` currently mandates
  the structured identity shapes and forbids the rendered string forms in
  persisted JSON (the rule cases
  [0058](../archive/0058-model-reference-identifiers.md) and
  [0059](../archive/0059-unqualified-model-references.md) established). This spec
  rewrites it, so the durable rationale below must be promoted into the durable
  spec, not left to die with this case.

## Scope

Covered: the shape of persisted identity fields in Evaluation routine and report
JSON, the `root` Area-name reservation, and the data `SchemaVersion` bump.
Deferred and non-goals are recorded in the
[parent concept](../0120-string-model-identity-fields.md#scope) — in short, no
change to human-facing rendered reports, no migration of existing runs, no new
selector surface, and no change to the `QUALITY.md` authoring format.

## Assumptions & dependencies

- [`SPECIFICATION.md`](../../SPECIFICATION.md)'s qualified model references
  (`area:<path>`, `factor:<area>::<path>`, `requirement:<area>::<name>`,
  `rating:<id>`) and strict name grammar are the encoding. A change to either
  invalidates the losslessness this spec relies on.
- The CLI already parses and resolves these qualified references against a model
  (the existing reference parsers and the `model-snapshot.md` resolution path);
  this change reuses that machinery rather than adding a new one.

## Requirements

### Persisted identity fields

Every persisted Evaluation routine or report JSON field that identifies a single
Area, Factor, Requirement, or Rating Level — `areaId`, `factorId`,
`requirementId`, and `ratingLevelId`, including where they are nested in a
`subject`, a routine ref, or a report ref — **MUST** be a single canonical
qualified model-reference string: `area:<area-path>`,
`factor:<declaring-area-path>::<factor-path>`,
`requirement:<declaring-area-path>::<requirement-name>`, and `rating:<level-id>`
respectively.

> Rationale: the qualified reference is a lossless string encoding of the
> composite identity (the strict name grammar excludes the separators), so the
> structured object adds bytes, not information. `ratingLevelId` is included so
> all four identity fields share one qualified-reference shape rather than
> leaving a bare value that re-introduces a smaller version of the same
> inconsistency. — 0120

A persisted identity field **MUST NOT** carry the structured `declaringAreaId`,
`factorPath`, or `requirementName` sub-fields, and **MUST NOT** carry an
unqualified (prefixless) reference.

> Rationale: keeping the qualified prefix matches
> [`SPECIFICATION.md`](../../SPECIFICATION.md)'s rule that durable
> machine-readable artifacts use qualified, never unqualified, references — the
> field name fixing the type does not license dropping the prefix in persisted
> data. — 0120

A repeated identity field — `areaIds`, `factorIds`, `localRequirementIds`,
`rootFactorIds`, `childAreaIds`, `ratingLevelIds`, and the secondary Factor lists
(`factors`, subject `factorIds`) — **MUST** be an array of those qualified
model-reference strings, defaulting to `[]` when empty.

> Rationale: the secondary Factor lists are already strings but area-relative and
> unqualified; rendering them as qualified `factor:` references brings every
> identity field, scalar or repeated, to one rule. — 0120

The identity field names **MUST** keep the `*Id` suffix; this change **MUST NOT**
rename them to `*Ref`.

> Rationale: `*Id` still names the model element's identity, now string-encoded;
> `*Ref` already denotes a link to another routine record (e.g.
> `factorAnalysisRef`), and overloading it would blur that distinction. — 0120

### Reserved `root` Area name

`root` **MUST NOT** be accepted as an Area name. `qualitymd lint` **MUST** report
a model that uses `root` as an Area name as an error, and the rule **MUST** be
documented in the lint rules.

> Rationale: `root` is the reserved root-path token in Area, Factor, and
> Requirement references; an Area named `root` would render `area:root`, colliding
> with the empty root path and breaking the round-trip that makes the string form
> lossless. — 0120

### Validation, schema, and versioning

The CLI **MUST** resolve each persisted identity string against the run's
`model-snapshot.md` before accepting a write, and **MUST** reject a payload whose
identity string names an absent Area, Factor, Requirement, or Rating Level —
preserving the existing model-binding guarantee on the new shape.

The generated `evaluation-data.schema.json` **MUST** type the identity fields as
strings (not arrays or objects), and the committed schema **MUST** be regenerated
to match.

The Evaluation data `SchemaVersion` **MUST** be incremented for this
payload-shape change. The CLI **MUST** reject payloads carrying the prior
structured identity shapes rather than silently migrating them, consistent with
the `schemaVersion`-is-a-marker rule in
[`json-conventions.md`](../../specs/evaluation/records/json-conventions.md).

> Rationale: identity shape is part of the payload contract; bumping the version
> and rejecting old shapes keeps write/verify parity and avoids a silent
> dual-shape reader. Existing runs are regenerated, not migrated. — 0120

## Durable spec changes

### To add

None. (The parent concept suggests a future `evaluation-data-schema-json.md`
artifact spec, but creating it is not required for this case to land.)

### To modify

- `specs/evaluation/records/json-conventions.md` — rewrite "Identity And
  References" so persisted identity fields are canonical qualified
  model-reference strings, dropping the structured-shape mandate and the ban on
  string forms (per the _Persisted identity fields_ requirements above); promote
  the losslessness and reversal rationale into the spec.
- `specs/evaluation/records/payload-kinds.md` — align the per-kind identity field
  descriptions with the string shape (per the _Persisted identity fields_
  requirements above).
- `specs/evaluation/reports/report-tree.md` — align the ordering rules that name
  "declaring Area ID and structural Factor path / Requirement name" with the
  string shape (per the _Persisted identity fields_ requirements above).
- `SPECIFICATION.md` — reserve `root` as a forbidden Area name in the name-grammar
  section (per the _Reserved `root` Area name_ requirement above).
- `specs/cli/lint-rules.md` and `specs/cli/lint.md` — add the reserved-`root`
  Area-name diagnostic (per the _Reserved `root` Area name_ requirement above).

### To rename

None.

### To delete

None.
