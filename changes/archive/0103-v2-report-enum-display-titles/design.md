---
type: Design Doc
title: Evaluation v2 report enum display titles - design
description: Design for typed display catalogs used by Evaluation v2 Markdown report rendering.
tags: [cli, evaluation, reports, display]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report enum display titles - design

## Context

This design answers the
[functional spec](spec.md) for rendering Evaluation v2 report enum-like values
as human titles. The current renderer already has a clean projection boundary:
routine JSON is loaded into maps, then `report_v2.go` writes deterministic
Markdown. That boundary is the right place to decorate values without changing
stored data.

## Approach

Add a small display helper file in `internal/evaluation`:

1. Define a generic `displayCatalog[T ~string]` with a `title` method.
2. Define named string types and constants for CLI-owned status and confidence
   vocabularies used by reports.
3. Add typed catalogs for owned values and string fallback catalogs for
   currently open finding fields.
4. Add focused helper functions such as `analysisStatusTitle`,
   `ratingStatusTitle`, `confidenceTitle`, and `findingSeverityTitle` that
   convert loose JSON strings at the renderer boundary.
5. Keep `ratingTitle(spec, level)` in `report_v2.go` for model-owned Rating
   Level labels.

The fallback helper humanizes unknown strings by splitting common enum forms
(`not_rated`, `missing-evidence`, camel-case values) and title-casing the words.
This keeps future or free-form values readable without pretending they are
closed by the current schema.

## Spec response

The typed catalogs satisfy the type-safety requirement for owned vocabularies
without changing routine payload structs or validators. Renderer call sites
convert string values from JSON into the relevant named type only when producing
Markdown.

Because display titles are applied only inside Markdown writer helpers,
`EvaluationOutputResult`, report refs, build receipts, and routine JSON continue
to carry raw values. Report tests cover both surfaces by checking titled Markdown
and unchanged receipt or output JSON values.

## Alternatives

### Add title fields to JSON

Rejected. The user need is a human report scanning aid, not a new machine data
contract. Adding title fields would create drift risk and force downstream tools
to distinguish identity from presentation.

### Use one untyped `map[string]string`

Rejected for owned vocabularies. It would be simple, but it would also make
payload kinds, statuses, confidence levels, and report kinds interchangeable in
the catalog layer. Generic typed catalogs keep the implementation small while
catching those mistakes at compile time.

### Type every finding field now

Rejected. Finding `type` and `severity` are not currently closed by the v2
contract. The renderer can title known strings, but making them typed constants
would imply a schema constraint this change does not introduce.

### Put emoji in the stored values

Rejected. Stable stored values need to stay terse and machine-friendly. Emoji
belongs in the report display title, paired with a word label.

## Trade-offs & risks

Emoji can make dense tables noisy if applied everywhere. The catalogs focus on
high-signal states, confidence, and classification cells, and preserve word
labels for readability and accessibility.

The main implementation risk is missing a raw status call site in the renderer.
A search for `v2String(..., "status")`, `confidence`, `severity`, `type`, and
the helper tables in `report_v2.go` identifies the relevant Markdown cells.

## Open questions

None.
