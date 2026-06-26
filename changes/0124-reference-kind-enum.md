---
type: Change Case
title: Constrain reference kind fields to the payload-kind vocabulary
description: Make the `kind` member of Evaluation reference objects (`*Ref` / `inputRefs[]`) an enum over the supported payload kinds, so a misspelled or invented reference kind is rejected at write time instead of persisting as a free-form string.
status: Design
tags: [evaluation, records, references, schema]
timestamp: 2026-06-26T00:00:00Z
---

# Constrain reference kind fields to the payload-kind vocabulary

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0124-reference-kind-enum/spec.md) - what the change must do.
- [Design doc](0124-reference-kind-enum/design.md) - how the contract enum and
  its typed source land.

## Motivation

Evaluation payloads reference other payloads by a `{kind, subject, selector}`
reference object: an `AreaAnalysisResult` lists the `FactorAnalysisResult` and
child `AreaAnalysisResult` inputs that drove it, the generated
`EvaluationOutputResult` points back at the area outputs, and rating drivers cite
the assessments and ratings behind them. The `kind` member of those references
names the *type of the referenced payload*.

The data contract already constrains the closed vocabularies it cares about —
`finding.type`, `severity`, `confidence`, and `status` are all enum-validated, so
an out-of-vocabulary value is rejected before the payload is written
(`internal/evaluation/data_contract.go`). The reference `kind` is the conspicuous
exception: it is a required field typed as a free `dataString`
(`routineRefContract` and `reportRefContract`), so a misspelled or invented kind —
`AreaAnalysis`, `RequirementResult`, `coverage_gap` — passes validation and
persists. Identity fields in the same reference are already resolved against the
model snapshot and rejected when absent; `kind` is the one part of a reference
the CLI does not check against its own closed set.

The acquire-roi-next evaluation surfaced the cost concretely: an orchestrator
hand-deriving ~115 payloads (each carrying reference objects) has no early signal
when a `kind` is wrong. A typo that the contract would reject is far cheaper to
fix at `evaluation data set` time than as a dangling reference discovered during
`report build`. Constraining `kind` to an enum closes the last unguarded field in
the reference shape and brings it under the same "reject out-of-range enum
values" rule the payload-kinds spec already states.

## Scope

Covered:

- Constrain the `kind` member of every Evaluation reference object — the
  `routineRefContract` used for `inputRefs[]` and the `*Ref` objects, and the
  `reportRefContract` used in report references — to the supported Evaluation
  payload-kind vocabulary, validated from one typed source of truth.
- Reject a write whose reference `kind` is present but outside that vocabulary,
  with an error naming the offending field and value, before persisting.
- Surface the constraint as an `enum` in `qualitymd evaluation data schema`
  output (and regenerate the committed `evaluation-data.schema.json`), so the
  closed set is discoverable from the typed contract.
- Validate against the full set of supported payload kinds — including the
  CLI-owned `EvaluationOutputResult` — not only the agent-writable subset that
  gates `evaluation data set <kind>`, because naming a payload type in a
  reference is distinct from being permitted to write that type.

Deferred / non-goals:

- No change to the top-level payload `kind` field. It is already pinned by the
  `evaluation data set <kind>` argument path and the agent-writable
  `acceptedDataKinds` gate; this change touches only the `kind` *inside reference
  objects*.
- No change to reference `subject` identity resolution or to the `selector`
  field. `selector` is also a bounded string and a candidate for the same
  treatment, but its full vocabulary needs confirming first; it is tracked
  separately, not in this case.
- No `SchemaVersion` bump. Tightening an existing field from "any string" to "a
  closed set that every valid payload already satisfies" rejects only
  previously-invalid data; it does not change the shape of conforming payloads.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0124-reference-kind-enum/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

### Code

- [ ] `internal/evaluation/data_contract.go` - constrain the `kind` field in
      `routineRefContract` and `reportRefContract` with an `enum` over the
      supported payload kinds; the enum/reject plumbing
      (`validateDataValue`, `enum()`) and schema emission already exist and need
      no change.
- [ ] `internal/evaluation/data.go` - source the enum values from one typed list
      of supported payload kinds (the `DataKind` constants), distinct from the
      agent-writable `acceptedDataKinds` subset.
- [ ] Tests - `internal/evaluation/evaluation_test.go` (and
      `internal/cli/evaluation_test.go` if it asserts reference shapes): add a
      rejection case for an out-of-vocabulary reference `kind`.

### Generated schema

- [ ] `internal/evaluation/evaluation-data.schema.json` - committed generated
      artifact; regenerate after the contract change (the staleness guard in
      `evaluation_test.go` enforces this).

### Format spec

- [ ] [`SPECIFICATION.md`](../SPECIFICATION.md) - reviewed; no change. It defines
      model references and the name grammar, not the per-field reference contract.

### Durable specs

- [ ] `specs/evaluation/records/json-conventions.md` - in "Identity And
      References", state that a reference's `kind` names a supported payload kind
      and is rejected when it does not, parallel to the existing identity-field
      resolution rule.

### Durable docs / bundled skill

- [ ] Reviewed; no change. The bundled skill and its specs already use valid
      reference kinds; the evaluate workflow names `AreaAnalysisResult` only in
      descriptive prose.

### Suggested new durable specs

- As noted in 0120, the Evaluation data contract has no 1:1 durable artifact spec
  (`evaluation-data-schema-json.md`); the reference-`kind` enum is one more
  contract detail living only in `json-conventions.md` prose and generated JSON.
  Suggesting only; not a precondition to land.

## Status

`Design`. Functional spec and design doc authored; code not started. Durable
specs not yet edited.
