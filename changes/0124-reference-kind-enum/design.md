---
type: Design Doc
title: Constrain reference kind fields — design doc
description: How the reference-kind enum lands on the Evaluation data contract and its typed source of truth.
tags: [evaluation, records, references, schema]
timestamp: 2026-06-26T00:00:00Z
---

# Constrain reference kind fields — design doc

Design behind the
[Constrain reference kind fields](../0124-reference-kind-enum.md) change case and
its [functional spec](spec.md).

## Context

The contract's validation, enum, and schema-emission machinery already exists and
is exercised by the other enum fields; the only gap is that two reference
contracts type their `kind` as a free `dataString`. This is a small, additive
constraint, not new machinery.

The relevant existing pieces in `internal/evaluation`:

- `data_contract.go` — `routineRefContract()` (`field("kind", dataString, true)`)
  and `reportRefContract()` (same) declare the reference shapes.
- `data_contract.go` — `enum(values...)` sets a field's closed set;
  `validateDataValue` already rejects an out-of-set string with
  `"%s = %q, want one of %s"`; `schemaForField` already emits `enum` into the
  JSON schema. None of these change.
- `data.go` — `acceptedDataKinds` is the **agent-writable** subset (nine kinds,
  excluding `EvaluationOutputResult`) that gates `evaluation data set <kind>`.
- `evaluation-data.schema.json` — committed generated schema, guarded by a
  staleness test in `evaluation_test.go`.

## Approach

1. **One typed source for the full payload-kind set.** Add a single slice of all
   supported payload kinds (all ten `DataKind` constants) in `data.go` — call it
   `supportedDataKinds` — and define the existing `acceptedDataKinds` as that set
   minus the CLI-owned `EvaluationOutputResult`, so "agent-writable" stays a
   derived view of one list rather than a second hand-maintained list. Expose the
   string values via a small helper (e.g. `supportedDataKindStrings()`).

2. **Constrain the two reference `kind` fields.** Change both
   `field("kind", dataString, true)` declarations to
   `field("kind", dataString, true, enum(supportedDataKindStrings()...))`.
   Validation and schema emission then cover reference `kind` with no further
   code: an out-of-set value is rejected at `evaluation data set` time with the
   field path and value, and `evaluation data schema` shows the enum.

3. **Regenerate the committed schema.** Run `qualitymd evaluation data schema`
   to refresh `evaluation-data.schema.json`; the staleness guard fails the build
   otherwise.

4. **Add a rejection test** asserting that a payload with an out-of-vocabulary
   reference `kind` is rejected, alongside the existing enum-rejection cases.

## Spec response

- *Constrain reference `kind`* and *single typed source*: steps 1–2 — both
  reference contracts read their vocabulary from `supportedDataKinds`, the same
  list that defines what the CLI persists.
- *Full vocabulary including `EvaluationOutputResult`*: step 1 uses the
  ten-kind superset, not `acceptedDataKinds`.
- *Reject out-of-vocabulary at write time, naming field and value*: inherited
  from `validateDataValue`'s existing enum check (step 2), no new code.
- *Appear as an `enum` in `data schema`*: inherited from `schemaForField` (step
  2), verified by the regenerated golden (step 3).

## Alternatives

- **Restrict the enum to the agent-writable nine.** Rejected: a reference *names*
  a payload type; whether an agent may *write* that type through `set` is a
  separate concern. Restricting would reject a legitimate reference to the
  CLI-owned `EvaluationOutputResult` and conflate two ideas. Using the superset is
  both simpler (no special-case exclusion) and more correct.
- **Cross-validate `kind` against the reference's `subject` identity field**
  (e.g. a `requirementId` subject implies a `Requirement*` kind). More semantic,
  but beyond the spec's one obligation and easy to get subtly wrong; the enum is
  the minimal correct step. Left as an open question.
- **Leave `kind` free and normalize/validate in the skill.** Rejected: the
  contract is the enforcement point every writer shares; pushing vocabulary rules
  into the skill is exactly the silent-typo path this case closes.

## Trade-offs & risks

- **Tightening rejects previously-accepted data.** A payload with a bogus
  reference `kind` validated before and will not now. That data was invalid in
  intent (a dangling reference), and under early-alpha clean-break policy this is
  the desired behavior, not a regression. No `SchemaVersion` bump: conforming
  payloads are unaffected.
- **Two reference contracts must stay in sync.** Both now point at the same
  `supportedDataKinds` helper, so they cannot drift from each other or from the
  CLI's real kind set.
- **Low blast radius.** No new flags, no shape change, no migration; the change
  is one slice plus two `enum(...)` arguments plus a regenerated golden.

## Open questions

- Should a later case cross-validate reference `kind` against `subject` identity
  type? Deferred; the enum is independently valuable and unblocks that.
- The reference `selector` field is the next bounded-string candidate; its full
  vocabulary needs confirming before it gets the same treatment. Tracked outside
  this case.
