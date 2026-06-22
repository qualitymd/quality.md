---
type: Design Doc
title: Model reference identifiers - design doc
description: How strict model names, canonical references, lint checks, and Area Breakdown rendering should be implemented.
tags: [format, references, reports, cli, lint]
timestamp: 2026-06-22T00:00:00Z
---

# Model reference identifiers - design doc

## Context

This design answers the
[Model reference identifiers functional spec](spec.md). The change makes the
Model's key vocabulary explicit, rejects ambiguous local names, introduces
canonical typed references at human/tool boundaries, and changes the shared Area
Breakdown table to show both human titles and stable Area references.

The important boundary is that durable machine artifacts already store
structured `areaPath` and `factorPath` arrays. This change should add canonical
string references for human-facing surfaces and helpers, not replace the
structured source of truth.

## Approach

Add model-reference helpers to `internal/evaluation`, next to the existing
`AreaPath` and `FactorPath` types:

```go
func (p AreaPath) Reference() string
func FactorReference(areaPath AreaPath, factorPath FactorPath) string
func RatingReference(level string) string
```

`AreaPath.Display()` and `FactorPath.Display()` can keep their current compact
debug labels. Report rendering should use the new reference helpers wherever the
spec asks for stable model references. Parsing support should be limited to
small helpers for canonical `area:`, `factor:`, and `rating:` forms so future
selector surfaces do not need to copy the grammar.

Lint should own strict-name enforcement because the JSON Schema can express only
part of the contract and cannot validate semantic references. Add a shared name
validator for Area names, Factor names, and Rating Level IDs, then emit one
error rule per concept:

- `invalid-area-name`
- `invalid-factor-name`
- `invalid-rating-level-id`

These rules should run while walking the model, before later reference checks
depend on the names. They should not apply to Requirement statement keys.
Existing `unknown-factor` and `unknown-rating-key` checks remain separate:
invalid declared names are definition errors; references to missing valid names
are resolution errors.

The structural schema package can carry optional key-pattern metadata on map
properties whose keys are model names. JSON Schema generation should use
`propertyNames.pattern` for `areas` and `factors`, and a string `pattern` for
`ratingScale[].level`. It should not apply a pattern to `requirements` or
`ratings` override keys.

The report change is narrow. Update the shared Area Breakdown renderer to write:

```md
| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |
```

The Area cell should use the resolved title. The Path cell should use
`area:<path>` with `area:root` for the root row. The two rating cells should
preserve the existing structural, not-assessed, and Rating Level display logic.
The Factors cell should preserve the existing compact factor title/rating list
and explicit `(no factor ratings)` empty state.

Specs, runtime skill files, and examples should be updated in the same pass so
the visible contract, generated reports, and agent instructions use the same
reference vocabulary.

## Alternatives

**Replace `areaPath` and `factorPath` arrays with strings.** Rejected because
arrays are unambiguous machine data and are already used by evaluation records,
report JSON, and validation. Canonical references are a boundary format, not the
durable storage model.

**Make JSON Schema the only strict-name enforcer.** Rejected because lint owns
the user-facing diagnostics and still needs to validate parsed YAML with stable
rule IDs. The schema should help editors catch structural cases, but lint
remains the authoritative tool check.

**Use title-bearing display paths in the report `Path` column.** Rejected
because the old table made the stable handle ambiguous. Titles belong in the
Area column; references belong in the Path column.

**Accept shorthand throughout the codebase.** Rejected because shorthand is safe
only when the expected reference type is fixed at the input edge. Durable mixed
surfaces and generated artifacts should render canonical typed references.

## Trade-offs & Risks

Adding key-pattern metadata to the structural schema introduces another schema
property to keep in sync with lint. The risk is bounded by deriving JSON Schema
from the same metadata and by adding tests for the generated patterns and lint
rules.

Strict names are a compatibility tightening. Existing models with spaces, dots,
slashes, or leading/trailing separators in Area names, Factor names, or Rating
Level IDs will fail lint. The report and evaluation JSON paths remain structured
arrays, which reduces migration risk for existing evaluation records that were
already written with valid keys.

Model-reference parsing can grow into a selector language if overbuilt. Keep
this change to canonical typed parse/render helpers and the currently specified
edge shorthand rules.

## Open Questions

- Should future CLI commands expose a single mixed `--scope` model-reference
  flag, or continue with type-specific commands where shorthand is available?
- Should report JSON add derived `areaRef` or `factorRef` fields in a later
  change, or are structured arrays enough for machine consumers?
