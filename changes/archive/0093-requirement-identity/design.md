---
type: Design Doc
title: Named Requirement identity — design
description: Implementation approach for stable Requirement names, titles, and references.
tags: [format, requirements, model-references]
timestamp: 2026-06-25T00:00:00Z
---

# Named Requirement identity — design

Companion to the [functional spec](spec.md). This design records how the
Requirement model-shape change is implemented.

## Context

The existing parser and linter treat the Requirement map key as the
human-readable Requirement statement. That makes authoring pleasant, but it is a
poor durable identity for generated evaluation records: wording changes mutate
identity, references are awkward, and Requirement-scoped outputs cannot use the
same reference discipline as Areas, Factors, and Rating Levels.

## Approach

Use one breaking format change rather than a dual-shape compatibility layer.
The Requirement map key becomes the stable Requirement name, following the same
strict grammar as Area names, Factor names, and Rating Level IDs. The
human-facing statement moves to a required `title`; the existing `assessment`
field keeps its meaning; optional `description`, `factors`, and `ratings` remain
ordinary Requirement fields.

Implementation work follows the existing single-source schema path:

- `internal/schema` declares Requirement key patterns and the required `title`.
- `internal/lint` walks Requirements by stable name, emits
  `invalid-requirement-name` for bad map keys, reuses `missing-title` for missing
  Requirement titles, and preserves the existing assessment, factor-reference,
  and rating-override checks.
- `internal/model` adds `Title` and `Description` to the typed Requirement
  struct. Existing `map[string]Requirement` containers stay in place, so the map
  key remains the machine identity.
- `internal/evaluation` adds Requirement reference render/parse helpers beside
  the existing Area, Factor, and Rating helpers:
  `requirement:<declaring-area-path>::<requirement-name>`.
- `quality.schema.json` continues to be generated from `internal/schema`, so the
  companion schema follows the same structural contract as lint.

The repo's own model, scaffold templates, public examples, bundled skill
guidance, and durable specs are migrated in the same change so local validation
tests exercise the new shape.

## Alternatives

Supporting both legacy statement-key Requirements and named Requirements was
rejected for v1 of this format change. It would make lint, schema, model
references, and generated evaluation data ambiguous: a natural-language key that
happens to match the strict grammar could be read as either a legacy statement or
a stable name. The project can add explicit migration help later without making
the accepted format itself ambiguous now.

Using a `name` property inside each Requirement instead of the map key was also
rejected. Areas and Factors already use the map key as stable identity, and
switching only Requirements to an internal `name` would make traversal,
reference generation, and authoring examples inconsistent.

Renaming `assessment` was deferred. The current field is already established,
and this change is about identity, not reworking the basis/assessment vocabulary.

## Trade-offs & risks

This is a breaking format change. Existing QUALITY.md files with statement-key
Requirements will fail lint until migrated. That is intentional, but the release
notes and authoring guidance need to be direct about the required edit:
choose a stable Requirement name and move the previous statement into `title`.

The delimiter stays `::` for Requirement references to match Factor references.
That keeps the parser simple and avoids ambiguity between typed-reference
prefixes and the area/name separator.

## Open questions

Automated migration support remains deferred. A future change can add an
authoring workflow or CLI helper that proposes stable names from legacy
Requirement statements.
