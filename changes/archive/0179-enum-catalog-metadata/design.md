---
type: Design Doc
title: Enum Catalog Metadata - design doc
description: Design for adding type-level and value-level descriptions to Evaluation enum catalogs.
tags: [evaluation, reports, enums, glossary]
timestamp: 2026-06-29T00:00:00Z
---

# Enum Catalog Metadata - design doc

Design for
[Enum Catalog Metadata](../0179-enum-catalog-metadata.md) and its
[functional spec](spec.md).

## Context

`internal/evaluation/display.go` currently represents each fixed Evaluation
vocabulary as a slice of enum values with the canonical value, label, marker,
and optional rank. `report_tree.go` separately supplies local-key family labels
when it calls `fixedEnumKeyLine("Type", findingTypeValues)`, which leaves the
vocabulary's name outside the catalog.

## Approach

Wrap each fixed Evaluation value slice in a small generic catalog type:

```go
type enumCatalog[T ~string] struct {
    Label       string
    Description string
    Values      []enumValue[T]
}

type enumValue[T ~string] struct {
    Value       T
    Label       string
    Marker      string
    Description string
    Rank        int
}
```

Keep the existing catalog variable names so call sites continue to read as the
vocabulary they reference. Update helper functions to accept `enumCatalog[T]`
and derive values from `catalog.Values` for validation, schema enum lists,
report display titles, report local keys, and rank lookup.

`fixedEnumKeyLine` should take only the catalog. It reads `catalog.Label` and
the marker-plus-label values, so report key labels are no longer duplicated at
each renderer call site. This directly changes finding local keys from `Type`
and `Severity` to `Finding type` and `Finding severity`.

Descriptions remain metadata only in this change. They are required by tests but
are not rendered in generated reports or emitted in JSON Schema. That keeps the
scope local to catalog completeness and report key labels while leaving future
glossary or schema-description surfaces free to choose their own presentation.

## Spec response

- Catalog labels and descriptions live at the vocabulary level.
- Value descriptions live beside labels and markers.
- Report local keys render catalog labels and marker-plus-label values only.
- Strict persisted enum values and generated schema enum arrays continue to come
  from catalog values.
- Dense table headers remain literal report layout labels, not catalog labels.

## Alternatives

- **Replace only `"Type"` with `"Finding type"`.** Rejected because it fixes the
  immediate report ambiguity while preserving the same drift-prone call-site
  label pattern.
- **Store descriptions in durable specs only.** Rejected because future
  glossary/help surfaces would still need to duplicate those descriptions in
  code.
- **Emit descriptions in JSON Schema immediately.** Deferred because plain enum
  arrays are the current schema contract; richer schema presentation should be a
  separately scoped compatibility decision.
- **Use one untyped registry for all catalogs.** Rejected because the current
  typed generic helpers keep construction sites type-checked and simple.

## Trade-offs & risks

Adding a wrapper type touches many helper call sites, but it keeps the metadata
shape explicit and small. The main risk is treating catalog descriptions as
report prose; tests and specs keep report local keys notation-only.

Changing local-key labels regenerates report-gallery Markdown. That is an
intentional output change and should be verified through the report gallery
check.

## Open questions

None.
