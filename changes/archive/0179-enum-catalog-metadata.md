---
type: Change Case
title: Enum Catalog Metadata
description: Add type-level labels and descriptions plus value descriptions to Evaluation enum catalogs.
status: Done
tags: [evaluation, reports, enums, glossary]
timestamp: 2026-06-29T00:00:00Z
---

# Enum Catalog Metadata

Evaluation enum catalogs now centralize persisted values, labels, markers, and
ordering, but the display name for a vocabulary still lives at call sites such
as report local-key rendering. The catalogs also cannot yet power future
glossary or help surfaces because they have no type-level descriptions or
per-value descriptions.

- [Functional spec](0179-enum-catalog-metadata/spec.md) - what catalog metadata,
  report key labels, and durable documentation must change.
- [Design doc](0179-enum-catalog-metadata/design.md) - how the Evaluation
  package extends enum catalog metadata while preserving strict persisted data.

## Motivation

Generated report keys should be unambiguous without making dense table columns
wide. `Finding type` and `Finding severity` are clearer local-key labels than
generic `Type` and `Severity`, but hard-coding those labels in renderer call
sites repeats metadata that belongs with the enum catalog. Adding catalog and
value descriptions now also prepares the existing typed catalogs for future
glossary, help, or schema-description surfaces.

## Scope

Covered:

- Evaluation enum catalog type-level labels and descriptions;
- Evaluation enum value descriptions;
- generated Markdown report local keys that render fixed Evaluation enum catalog
  labels;
- durable Evaluation report specs and report design guidance for catalog-owned
  key labels;
- checked-in report-gallery Markdown affected by generated report output;
- tests that require complete catalog metadata.

Deferred:

- changing persisted enum values;
- accepting display labels, descriptions, markers, aliases, or case variants as
  structured data;
- adding a public glossary, help command, or schema-description output;
- changing dense table column labels such as `Type`, `Severity`, `Basis`, or
  `Impact`;
- changing model-defined Rating Level titles.

## Affected artifacts

- Code:
  - `internal/evaluation/display.go` - add enum catalog metadata and value
    descriptions.
  - `internal/evaluation/report_tree.go` - render local keys from catalog
    labels.
  - `internal/evaluation/data_contract.go` - keep validation/schema enum lists
    derived from catalog values.
  - `internal/evaluation/evaluation_test.go` - cover catalog metadata and
    updated report keys.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - define catalog-owned fixed enum
    key labels and keep table-column labels compact.
- Durable docs:
  - `docs/guides/reporting-design.md` - update local-key guidance and examples.
- Generated examples:
  - `examples/report-gallery/` - regenerate generated report Markdown.
- Release notes:
  - `CHANGELOG.md` - note generated report key label and enum catalog metadata
    changes.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/evaluation/log.md` and `docs/log.md`.

## Status

`Done`. Implemented and archived. Fixed Evaluation enum catalogs now carry
type-level labels and descriptions plus value descriptions, and generated report
local keys render catalog labels such as `Finding type`, `Finding severity`, and
`Recommendation impact`.
