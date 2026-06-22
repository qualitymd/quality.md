---
type: Functional Specification
title: Unqualified model references
description: Define bounded unqualified Area, Factor, and Rating references for fixed-type contexts.
tags: [format, references, reports, cli]
timestamp: 2026-06-22T00:00:00Z
---

# Unqualified model references

This Change Case spec defines the delta for the QUALITY.md format, generated
reports, qualitymd evaluation helpers, and `/quality` skill guidance: a formal
unqualified-reference allowance for contexts that already fix the reference type.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

Canonical typed references from 0058 solve ambiguity: `area:...`, `factor:...`,
and `rating:...` are self-describing and safe on mixed-reference surfaces. But a
generated report column named `Path` inside `## Area Breakdown` is already
Area-specific. Rendering `area:operations/incident-response` there repeats type
information the surrounding context supplies and makes the table harder to scan.

The fix should be formal, not incidental. Tools need permission to render or
accept a shorter unqualified form when context fixes the type, and they need the
same bright line that keeps mixed-reference and machine surfaces unambiguous.

## Scope

Covered:

- Unqualified reference terminology.
- Unqualified Area, Factor, and Rating Level reference forms.
- Fixed-type render and accept rules.
- Mixed-reference and durable-machine-artifact limits.
- Named code helpers for qualified and unqualified render forms.
- Area Breakdown `Path` column rendering.

Deferred / non-goals:

- No replacement of canonical typed references.
- No `Reference(prefix bool)` or similar boolean switch.
- No generic model-reference object hierarchy.
- No Area-local Factor shorthand that omits the declaring Area side.
- No new machine-readable report JSON string reference fields.

## Terminology

A **qualified model reference** is the canonical typed form defined by 0058:
`area:<area-path>`, `factor:<declaring-area-path>::<factor-path>`, or
`rating:<rating-level-id>`.

An **unqualified reference** is a reference that omits the typed prefix only when
the surrounding context fixes the reference type.

## Requirements

Tools **MAY** render or accept an unqualified reference when the surrounding
context fixes the reference type. An unqualified reference omits the typed prefix
but otherwise preserves the same path structure as the qualified reference.

> Rationale: the shorter form is safe only when the type is supplied by the
> field, column, command, or parameter contract; otherwise the value is
> ambiguous. — 0059

Unqualified Area references **MUST** render the root Area as:

```text
root
```

Nested unqualified Area references **MUST** join Area names with `/`:

```text
operations/incident-response
```

Unqualified Factor references **MUST** render as:

```text
<declaring-area-path>::<factor-path>
```

The root declaring Area **MUST** render as `root`:

```text
root::security
operations/incident-response::operability/backpressure
```

Unqualified Rating Level references **MUST** render as the Rating Level ID:

```text
target
```

Tools **MUST NOT** render or accept unqualified references on mixed-reference
surfaces or anywhere the reference type must be recoverable from the value alone.

Tools **MUST NOT** persist unqualified references in durable machine-readable
artifacts. Durable machine-readable artifacts continue to preserve structured
`areaPath`, `factorPath`, and rating `level` identifiers.

The report summary Area Breakdown `Path` column **MAY** render unqualified Area
references because the column is Area-specific.

The full report's shared compact Area Breakdown table **MUST** stay aligned with
the concise report's Area Breakdown reference form.

Qualitymd code **SHOULD** expose named qualified and unqualified render helpers
instead of boolean flags whose call sites hide the chosen form.

## Acceptance Criteria

- `SPECIFICATION.md` defines unqualified references for Areas, Factors, and
  Rating Levels.
- `SPECIFICATION.md` uses a formal **MAY** for tools rendering or accepting
  unqualified references when context fixes the reference type.
- `SPECIFICATION.md` forbids rendering or accepting unqualified references on
  mixed-reference surfaces.
- `SPECIFICATION.md` forbids persisting unqualified references in durable
  machine-readable artifacts.
- `specs/reports/report-summary-md.md` allows the Area Breakdown `Path` column
  to use unqualified Area references.
- `specs/reports/report-md.md` keeps the full report's shared Area Breakdown
  `Path` column aligned with `report-summary.md`.
- `AreaPath.Reference()` still returns `area:root` or
  `area:operations/incident-response`.
- `AreaPath.UnqualifiedReference()` returns `root` or
  `operations/incident-response`.
- `FactorReference(areaPath, factorPath)` still returns
  `factor:operations/incident-response::operability`.
- `UnqualifiedFactorReference(areaPath, factorPath)` returns
  `operations/incident-response::operability`.
- `RatingReference("target")` still returns `rating:target`.
- `UnqualifiedRatingReference("target")` returns `target`.
- Generated `report-summary.md` and `report.md` Area Breakdown tables render
  unqualified Area references in the `Path` column.
- `report.json` continues to preserve structured `areaPath` and `factorPath`
  arrays.
- Existing canonical parsing remains strict for qualified references.
- Unqualified parsing is accepted only through type-specific parsing paths where
  context fixes the expected reference type.

## Durable spec changes

### To add

None.

### To modify

- `SPECIFICATION.md` - define unqualified references, allowed render/accept
  contexts, Area/Factor/Rating forms, and machine-artifact limits according to
  the requirements above.
- `specs/reports/report-summary-md.md` - allow the Area Breakdown `Path` column
  to render unqualified Area references according to the report requirements
  above.
- `specs/reports/report-md.md` - keep the full report's shared compact Area
  Breakdown aligned with the summary behavior above.
- `specs/reports/report-json.md` - clarify that durable report JSON does not
  persist unqualified references and preserves structured paths according to the
  machine-artifact requirement above.
- `specs/evaluation-records/report-outputs.md` - align shared report-output
  terminology with unqualified references while preserving structured machine
  identifiers according to the requirements above.
- `specs/skills/quality-skill/reporting.md` - align skill-facing report
  guidance with unqualified Area Breakdown paths in human report artifacts and
  structured identifiers in `report.json`.
- `specs/skills/quality-skill/quality-skill.md` and
  `specs/skills/quality-skill/evaluation.md` - align fixed-type input guidance
  with the render/accept rules above.

### To rename

None.

### To delete

None.
