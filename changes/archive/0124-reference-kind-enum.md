---
type: Change Case
title: Constrain reference kind fields to closed kind vocabularies
description: Make the `kind` member of Evaluation reference objects an enum over the closed vocabulary it names — supported payload kinds for routine references (`*Ref` / `inputRefs[]`), report kinds for report references — so a misspelled or invented reference kind is rejected at write time instead of persisting as a free-form string.
status: Done
tags: [evaluation, records, references, schema]
timestamp: 2026-06-26T00:00:00Z
---

# Constrain reference kind fields to closed kind vocabularies

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
exception: it is a required field typed as a free `dataString` in both reference
shapes (`routineRefContract` and `reportRefContract`), so a misspelled or
invented kind — `AreaAnalysis`, `RequirementResult`, `coverage_gap` — passes
validation and persists. Identity fields in the same reference are already
resolved against the model snapshot and rejected when absent; `kind` is the one
part of a reference the CLI does not check against its own closed set.

Each reference shape's `kind` names a *different* closed vocabulary. A routine
reference (`*Ref` / `inputRefs[]`) `kind` names the *type of the referenced
payload* — an `AreaAnalysisResult` lists the `FactorAnalysisResult` and child
`AreaAnalysisResult` inputs that drove it — so its vocabulary is the supported
payload kinds. A report reference `kind` names a *report kind* (`area`, `factor`,
`requirement`), the scope of the generated report it points at. Both are closed
sets the CLI already owns as typed constants; neither was pinned in the contract.

The acquire-roi-next evaluation surfaced the cost concretely for routine
references, which agents hand-author: an orchestrator hand-deriving ~115 payloads
(each carrying reference objects) has no early signal when a `kind` is wrong. A
typo that the contract would reject is far cheaper to fix at `evaluation data set`
time than as a dangling reference discovered during `report build`. Report
references are CLI-generated (only inside the CLI-owned `EvaluationOutputResult`),
so for them the constraint is a contract self-check and schema self-documentation
rather than typo prevention — but it closes the same gap. Constraining `kind` to
an enum in both shapes closes the last unguarded field in the reference shape and
brings it under the same "reject out-of-range enum values" rule the payload-kinds
spec already states.

## Scope

Covered:

- Constrain the `kind` member of every Evaluation reference object to the closed
  vocabulary it names, validated from one typed source of truth each: the
  `routineRefContract` used for `inputRefs[]` and the `*Ref` objects to the
  supported Evaluation payload-kind set, and the `reportRefContract` used in
  report references to the report-kind set (`area`, `factor`, `requirement`).
- Reject a write whose reference `kind` is present but outside that vocabulary,
  with an error naming the offending field and value, before persisting.
- Surface the constraint as an `enum` in `qualitymd evaluation data schema`
  output (and regenerate the committed `evaluation-data.schema.json`), so the
  closed set is discoverable from the typed contract.
- Validate routine references against the full set of supported payload kinds —
  including the CLI-owned `EvaluationOutputResult` — not only the agent-writable
  subset that gates `evaluation data set <kind>`, because naming a payload type in
  a reference is distinct from being permitted to write that type.

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

- [x] `internal/evaluation/data_contract.go` - constrain the `kind` field in
      `routineRefContract` (payload kinds) and `reportRefContract` (report kinds)
      with an `enum`; the enum/reject plumbing (`validateDataValue`, `enum()`)
      and schema emission already exist and need no change.
- [x] `internal/evaluation/data.go` - add `supportedDataKinds` (all ten
      `DataKind` constants) as the typed source for the routine-reference enum,
      derive the agent-writable `acceptedDataKinds` from it, and add a generic
      `kindStrings` helper used by both reference enums.
- [x] `internal/evaluation/display.go` - add `reportKinds` (the `ReportKind`
      constants) as the typed source for the report-reference enum.
- [x] Tests - `internal/evaluation/evaluation_test.go`: add a rejection case for
      an out-of-vocabulary routine reference `kind` (report references are
      CLI-generated and not reachable through `evaluation data set`, so their
      enum is covered by the regenerated golden schema).

### Generated schema

- [x] `internal/evaluation/evaluation-data.schema.json` - committed generated
      artifact; regenerated after the contract change (the staleness guard in
      `evaluation_test.go` enforces this).

### Format spec

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) - reviewed; no change. It defines
      model references and the name grammar, not the per-field reference contract.

### Durable specs

- [x] `specs/evaluation/records/json-conventions.md` - in "Identity And
      References", state that a reference's `kind` names a value from a closed
      vocabulary — supported payload kinds for routine references, report kinds
      for report references — and is rejected when it does not, parallel to the
      existing identity-field resolution rule.

### Durable docs / bundled skill

- [x] Reviewed; no change. The bundled skill and its specs already use valid
      reference kinds; the evaluate workflow names `AreaAnalysisResult` only in
      descriptive prose.

### Suggested new durable specs

- As noted in 0120, the Evaluation data contract has no 1:1 durable artifact spec
  (`evaluation-data-schema-json.md`); the reference-`kind` enum is one more
  contract detail living only in `json-conventions.md` prose and generated JSON.
  Suggesting only; not a precondition to land.

## Status

`Done`. Both reference contracts now enum-constrain their `kind` from typed
sources (`supportedDataKinds` for routine references, `reportKinds` for report
references); the committed schema is regenerated; a routine-reference rejection
test is in place; `json-conventions.md` records the rule. During implementation
the design's premise that *both* reference `kind` fields name payload kinds was
corrected: a report reference `kind` names a report kind, so it is constrained to
the report-kind set instead. `mise run check` passes.
