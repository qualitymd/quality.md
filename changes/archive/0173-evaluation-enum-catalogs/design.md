---
type: Design Doc
title: Evaluation Enum Catalogs - design doc
description: Design for typed Evaluation enum catalogs and shared report display metadata.
tags: [evaluation, schema, reports, enums]
timestamp: 2026-06-29T00:00:00Z
---

# Evaluation Enum Catalogs - design doc

Design for
[Evaluation Enum Catalogs](../0173-evaluation-enum-catalogs.md) and its
[functional spec](spec.md).

## Context

`internal/evaluation/data_contract.go` already owns structural payload
validation and JSON Schema generation, but many enum values are literal
`enum("...")` lists. `internal/evaluation/display.go` owns report labels for
some of the same vocabularies, and `report_tree.go` owns severity and impact
ranking maps. Some display maps also include stale values not accepted by the
payload contract.

## Approach

Keep the catalogs inside `internal/evaluation` because they model Evaluation
payload and report concepts. Add a small generic value type that carries the
canonical persisted value plus display metadata:

```go
type enumValue[T ~string] struct {
    Value  T
    Label  string
    Marker string
    Rank   int
}
```

Use typed string constants for fixed vocabularies that currently lack them:
finding type, finding severity, finding basis status, recommendation impact,
finding ranking tier, and finding coverage disposition. Keep existing typed
vocabularies such as `DataKind`, `ReportKind`, status types, confidence,
`RunGapKind`, and `RatingResultKind`.

For each fixed vocabulary, declare an ordered slice of `enumValue[T]`. Helpers
derive:

- `enumStrings(values)` for `data_contract.go` enum validation and schema
  generation;
- `enumDisplay(values, raw)` for marker-plus-label report display;
- `enumRank(values, raw)` for severity and impact ordering.

The renderer continues to receive decoded JSON as `map[string]any`, so report
helpers still accept raw strings at the call sites. They convert through the
catalog helpers and fall back to `humanizeEnum` only for defensive rendering of
unknown values that bypassed validation. Generated reportability still depends
on `SetData`/`VerifyData` validation before normal report builds.

`data_contract.go` replaces raw literal enum lists with typed helper calls:
assessment status values, rating status values, analysis status values,
confidence values, finding type values, finding severity values, basis status
values, impact values, ranking tier values, coverage disposition values, data
kind values, and report kind values.

Report-specific artifact `type` frontmatter strings stay plain literals. They
are not persisted Evaluation data enum values and do not need emoji markers.

## Spec response

- The same typed catalogs feed validation, generated schemas, examples, display
  titles, and report sorting.
- Invalid fixed values remain rejected at the data set and verify boundaries.
- Existing persisted values do not change.
- Model-defined Rating Level IDs stay outside the fixed enum catalogs and keep
  using model snapshot titles.
- Report rendering gains explicit labels/markers for basis status, ranking tier,
  coverage disposition, and missing report kinds.
- Stale finding-type display values are removed from the active display catalog.

## Alternatives

- **Keep maps and add tests only.** Rejected because tests would detect drift
  after it happens but would not remove the duplicated sources that cause drift.
- **Use one untyped global registry keyed by string names.** Rejected because it
  would trade raw literals for string registry keys and lose Go type checking at
  catalog construction sites.
- **Generate Go from schema.** Rejected because the current source of truth is
  the Go contract that also emits schema; adding code generation would broaden
  the build pipeline for a small internal factoring change.
- **Remove all unknown-value fallbacks from report helpers.** Rejected for this
  case because historical or hand-edited data may still be inspected in
  defensive paths. Strict write/verify behavior remains the authoritative
  current-data boundary.

## Trade-offs & risks

The generic helper keeps the implementation small, but it still introduces a
new local abstraction. The risk is over-generalizing all report labels into the
catalog even when a value is not an enum. The design limits catalogs to fixed
Evaluation vocabularies and leaves model-defined Rating Levels and report
frontmatter taxonomy out.

Changing labels or markers can alter generated Markdown fixtures and report
gallery output. Focused tests should pin intentional display changes.

## Open questions

None.
