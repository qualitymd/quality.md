---
type: Design Doc
title: String model-identity fields in evaluation data
description: How the Evaluation data contract moves model identity fields to canonical qualified-reference strings while reserving `root` as an Area name.
status: Draft
tags: [evaluation, records, references, schema]
timestamp: 2026-06-26T00:00:00Z
---

# String model-identity fields in evaluation data

## Context

Implements the [0120 functional spec](spec.md): persisted Evaluation routine and
report JSON encodes model identities as canonical qualified references, keeps the
existing `*Id` field names, bumps the data schema version, rejects the old
structured shapes, and reserves `root` as an Area name so `area:root` remains
lossless.

## Approach

The implementation keeps the change at the existing contract chokepoints.

1. **Reference parsing and encoding.** `data.go` parses `areaId`, `factorId`,
   `requirementId`, and `ratingLevelId` from typed strings using the existing
   qualified reference grammar. The report-tree output writers use the same
   encoders so CLI-owned `EvaluationOutputResult` matches agent-authored
   routine payloads.
2. **Contract and schema.** `data_contract.go` keeps the domain-specific field
   types but changes their JSON Schema and validation to strings. Repeated
   identity fields become arrays of those same field types. `SchemaVersion`
   increments, making old structured payloads fail before write or verify.
3. **Model binding.** The recursive model-reference walker binds the new string
   fields against the run's `model-snapshot.md`, including repeated fields and
   routine/report refs. Rating references resolve through the rating scale.
4. **Reserved Area name.** The linter adds a dedicated reserved-name diagnostic
   for Area keys named `root`. The companion structural schema narrows Area map
   property names enough to catch the same literal key where JSON Schema can
   express it.

## Spec response

- **Persisted identity fields** are string-encoded at parse, validation, schema,
  examples, and report-output sites.
- **No old structured shapes** is enforced by the version bump and by closed
  payload object schemas.
- **Model binding** stays centralized in `validatePayloadModelBindings`, which
  now walks and resolves the string references.
- **`root` reservation** lands in lint, companion schema, durable CLI specs, and
  `SPECIFICATION.md`.

## Alternatives

- **Keep `ratingLevelId` prefixless** was rejected by the functional spec. It
  would leave one identity field outside the common qualified-reference rule.
- **Dual-shape reader for old runs** was rejected. The data schema version is a
  shape marker; existing runs are regenerated rather than migrated.
- **Rename fields to `*Ref`** was rejected because `*Ref` already identifies a
  routine record link, not the model element identity itself.
