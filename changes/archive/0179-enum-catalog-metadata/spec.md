---
type: Functional Specification
title: Enum Catalog Metadata - functional spec
description: Requirements for Evaluation enum catalog labels, descriptions, and generated report keys.
tags: [evaluation, reports, enums, glossary]
timestamp: 2026-06-29T00:00:00Z
---

# Enum Catalog Metadata - functional spec

Companion to
[Enum Catalog Metadata](../0179-enum-catalog-metadata.md). This spec states the
delta for Evaluation enum catalog metadata and generated report local keys. The
durable source of truth is absorbed into the Evaluation report spec and report
design guidance.

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Fixed Evaluation enum catalogs already keep persisted values, labels, markers,
and ordering together. Report local keys still pass vocabulary names such as
`Type` and `Severity` independently from those catalogs, which lets key labels
drift from the vocabulary they describe. The same catalogs also lack
descriptions, so future glossary or help surfaces would have to invent
definitions outside the typed metadata.

The change should keep raw JSON values strict and stable while making the
catalogs rich enough to own vocabulary labels and concise descriptions. Generated
report keys can then use fully qualified catalog labels such as `Finding type`,
while dense table headers remain compact.

## Scope

Covered:

- fixed Evaluation enum catalog type-level labels;
- fixed Evaluation enum catalog type-level descriptions;
- fixed Evaluation enum value descriptions;
- generated Markdown report local keys for fixed Evaluation enum catalogs;
- tests and generated examples affected by local-key output.

Deferred:

- new public glossary, help, or catalog-inspection commands;
- JSON Schema `oneOf`, `title`, or `description` output for enum values;
- parsing display labels, markers, descriptions, aliases, or case variants as
  structured data;
- dense report table column label changes;
- model-defined Rating Level metadata.

## Requirements

1. Each fixed Evaluation enum catalog **MUST** carry a type-level human label and
   description.

   > Rationale: The catalog should own the display name and glossary-ready
   > summary for the vocabulary instead of leaving report keys or future docs to
   > invent their own names. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

2. Each fixed Evaluation enum value **MUST** carry a concise description of the
   value's meaning.

   > Rationale: Labels and markers aid scanning, but future glossary or help
   > surfaces need stable meaning text that is not embedded in generated report
   > prose. Durable spec: none.

3. Generated Markdown report local keys for fixed Evaluation enum catalogs
   **MUST** render the catalog label as the key label.

   > Rationale: Local keys describe an indicator family, so labels such as
   > `Finding type` and `Finding severity` prevent unrelated indicator families
   > from bleeding together. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

4. Generated Markdown report local keys **MUST NOT** render enum catalog
   descriptions or value descriptions.

   > Rationale: Local keys are notation-only and should stay compact near first
   > use; definitions belong in future glossary/help surfaces. Durable spec:
   > modify `specs/evaluation/reports/report-tree.md`.

5. Generated Markdown reports **MUST** keep dense table headers compact when the
   table context already scopes the enum family.

   > Rationale: Fully qualified key labels solve legend ambiguity without
   > widening high-density tables such as Finding lists. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

6. Evaluation data validation, verification, and JSON Schema enum lists
   **MUST** continue to derive allowed values from the fixed enum catalog values
   and **MUST NOT** accept catalog labels, descriptions, markers, aliases, or
   case variants as persisted values.

   > Durable spec: none; existing strict enum behavior remains unchanged.

## Acceptance criteria

- All fixed Evaluation enum catalogs have non-empty type-level labels and
  descriptions.
- All fixed Evaluation enum values have non-empty labels, markers, and
  descriptions.
- Generated report local keys use catalog labels, including `Finding type` and
  `Finding severity`.
- Dense generated report table headers remain unchanged where they currently use
  compact labels such as `Type`, `Severity`, `Basis`, and `Impact`.
- Generated report-gallery Markdown is regenerated for intentional report key
  changes.
- Focused Go tests and repository checks pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - require fixed Evaluation enum
  local keys to use catalog-owned labels, keep descriptions out of report keys,
  and preserve compact table headers where context scopes the enum family
  (requirements 1, 3-5).

### To rename

None

### To delete

None
