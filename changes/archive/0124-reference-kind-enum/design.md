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
  and `reportRefContract()` (same) declare the reference shapes. Their `kind`
  fields name **different** closed vocabularies: a routine-reference `kind` is a
  payload kind (`AreaAnalysisResult`, …); a report-reference `kind` is a report
  kind (`area`, `factor`, `requirement`), emitted by `evaluationReportRef` from
  the `ReportKind` constants. `reportRefContract` is used only inside the
  CLI-owned `EvaluationOutputResult`.
- `data_contract.go` — `enum(values...)` sets a field's closed set;
  `validateDataValue` already rejects an out-of-set string with
  `"%s = %q, want one of %s"`; `schemaForField` already emits `enum` into the
  JSON schema. None of these change.
- `data.go` — `acceptedDataKinds` is the **agent-writable** subset (nine kinds,
  excluding `EvaluationOutputResult`) that gates `evaluation data set <kind>`.
- `display.go` — the `ReportKind` constants (`area`, `factor`, `requirement`).
- `evaluation-data.schema.json` — committed generated schema, guarded by a
  staleness test in `evaluation_test.go`.

## Approach

1. **One typed source per reference vocabulary.** Add a single slice of all
   supported payload kinds (all ten `DataKind` constants) in `data.go` — call it
   `supportedDataKinds` — and define the existing `acceptedDataKinds` as that set
   minus the CLI-owned `EvaluationOutputResult`, so "agent-writable" stays a
   derived view of one list rather than a second hand-maintained list. Add a
   `reportKinds` slice (the `ReportKind` constants) in `display.go` as the source
   for the report-reference vocabulary. Render either typed slice as strings with
   one generic helper, `kindStrings[T ~string]`, rather than a per-type helper.

2. **Constrain the two reference `kind` fields, each to its own set.** Change
   `routineRefContract`'s declaration to
   `field("kind", dataString, true, enum(kindStrings(supportedDataKinds)...))`
   and `reportRefContract`'s to
   `field("kind", dataString, true, enum(kindStrings(reportKinds)...))`.
   Validation and schema emission then cover each reference `kind` with no further
   code: an out-of-set value is rejected at `evaluation data set` time with the
   field path and value, and `evaluation data schema` shows the enum.

3. **Regenerate the committed schema.** Run `qualitymd evaluation data schema`
   to refresh `evaluation-data.schema.json`; the staleness guard fails the build
   otherwise.

4. **Add a rejection test** asserting that a payload with an out-of-vocabulary
   routine reference `kind` is rejected, alongside the existing enum-rejection
   cases. Report references are not reachable through `evaluation data set` (they
   appear only in the CLI-generated `EvaluationOutputResult`), so their enum is
   covered by the regenerated golden schema rather than a `set` rejection test.

## Spec response

- _Constrain reference `kind`_ and _single typed source per vocabulary_: steps
  1–2 — `routineRefContract` reads `supportedDataKinds` (the same list that
  defines what the CLI persists) and `reportRefContract` reads `reportKinds` (the
  `ReportKind` constants); neither is a hand-maintained second list.
- _Full payload vocabulary including `EvaluationOutputResult`_: step 1 uses the
  ten-kind superset for routine references, not `acceptedDataKinds`.
- _Reject out-of-vocabulary at write time, naming field and value_: inherited
  from `validateDataValue`'s existing enum check (step 2), no new code.
- _Each constraint appears as an `enum` in `data schema`_: inherited from
  `schemaForField` (step 2), verified by the regenerated golden (step 3).

## Alternatives

- **Constrain the report reference `kind` to payload kinds too** (this design's
  original premise). Rejected once implementation looked at the field: a report
  reference `kind` holds a report kind (`area`, `factor`, `requirement`), not a
  payload kind, so a payload-kind enum would reject every valid report reference
  and break the generated `EvaluationOutputResult` and its committed example. The
  correct move is one enum per shape, each over the vocabulary that shape's `kind`
  actually names. A narrower alternative — leave the report reference `kind` free
  — was also rejected: it is a closed vocabulary with a typed source already
  (`ReportKind`), so constraining it costs almost nothing and removes the last
  free-string `kind` in the reference shape. The report-reference enum is a
  contract self-check (the value is CLI-generated) rather than typo prevention,
  but it is still worth pinning.
- **Restrict the routine enum to the agent-writable nine.** Rejected: a reference
  _names_ a payload type; whether an agent may _write_ that type through `set` is a
  a payload type; whether an agent may _write_ that type through `set` is a
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
- **Each reference contract tracks a typed source.** The routine contract reads
  `supportedDataKinds` and the report contract reads `reportKinds`, so neither
  enum can drift from the CLI's real kind sets, and the two cannot silently use
  the wrong vocabulary for each other.
- **Low blast radius.** No new flags, no shape change, no migration; the change
  is two typed slices plus two `enum(...)` arguments plus a regenerated golden.

## Open questions

- Should a later case cross-validate reference `kind` against `subject` identity
  type? Deferred; the enum is independently valuable and unblocks that.
- The reference `selector` field is the next bounded-string candidate; its full
  vocabulary needs confirming before it gets the same treatment. Tracked outside
  this case.
