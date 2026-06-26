---
type: Design Doc
title: Remove info finding severity — design
description: How `info` is removed from the Evaluation finding severity vocabulary.
tags: [evaluation, schema, reports, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Remove info finding severity — design

## Context

Answers the [functional spec](spec.md) for change case
[0132](../0132-remove-info-finding-severity.md). The spec removes `info` from the
Evaluation finding severity vocabulary while keeping `note` as the finding type
for informational observations.

The current implementation has one map-backed typed contract that drives
validation, examples, and generated schemas in
`internal/evaluation/data_contract.go`; report sorting/display use small helper
maps in `report_tree.go` and `display.go`.

## Approach

### Remove `info` at the typed contract source

Update both finding contract helpers:

- `findingContract()` for Requirement Assessment findings.
- `areaFindingContract()` for Area Findings.

Their `severity` enum becomes:

```text
critical, high, medium, low
```

Because `data set`, `data verify`, `data example`, and `data schema` are all
derived from the same contract, this is the single source of truth for the
runtime behavior. Regenerating `internal/evaluation/evaluation-data.schema.json`
then updates the checked-in full-surface schema artifact.

### Keep report fallbacks defensive but remove `info` as known

Remove `info` from:

- `findingSeverityTitles` in `display.go`.
- `findingSeverityRank` in `report_tree.go`.

Unknown values already fall through to humanized display and sort after known
values. That defensive behavior is useful for hand-edited or historical artifacts
but does not make `info` a valid current value.

### Tests lock both schema and validation

Update schema tests so Area Finding `severity` exposes only the four allowed
values. Add validation cases proving `severity: "info"` is rejected for:

- Requirement Findings; and
- Area Findings.

Existing report rendering fixtures that used `type: note` with `severity: info`
should switch to a valid severity-bearing type or to a note without relying on
`info`, depending on what the fixture is trying to show. Since conditional
severity applicability is deferred, this case does not remove `severity` from
non-adverse types unless a fixture needs to stop demonstrating the bad pattern.

### Specs and skill guidance

Durable specs absorb the current rule:

- the format spec states the reduced severity vocabulary and the `note` routing
  rationale;
- payload and CLI specs state schema/validation behavior;
- report specs update sort order;
- skill specs and runtime guidance tell agents to use `type: note` for
  informational observations.

The bundled runtime `skills/quality/resources/SPECIFICATION.md` mirrors the root
specification so a released skill reads the same semantics.

## Spec response

- **Severity vocabulary** — satisfied by the typed contract enum change and
  regenerated schema artifact.
- **Validation** — satisfied by existing enum validation with new rejection
  tests.
- **Reporting and skill guidance** — satisfied by removing `info` from known
  report display/sort helpers and updating specs/runtime guidance.

## Alternatives

- **Keep `info` until conditional severity applicability lands.** Rejected
  because the value is already misleading and can be removed independently:
  `low` remains the least-severe adverse value, and `note` covers informational
  observations.
- **Alias `info` to `low`.** Rejected because it silently changes meaning and
  preserves stale payloads in an early-alpha contract where clean breaks are
  preferred.
- **Add a separate `significance` field now.** Rejected as premature. If positive
  findings later need ordering or impact, that should be designed as its own
  contract rather than overloading `severity`.

## Trade-offs & risks

- Existing runs containing `severity: "info"` become invalid under current
  validation. This is acceptable under the early-alpha compatibility policy and
  should be called out in release notes.
- Report generation remains defensive for historical data, but active data
  validation rejects the stale value.
- Conditional severity-by-type is still unresolved after this case. The specs
  explicitly defer that rule so this cleanup does not overclaim.

## Open questions

None for this design.
