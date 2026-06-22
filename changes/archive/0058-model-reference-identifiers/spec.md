---
type: Functional Specification
title: Model reference identifiers
description: Define strict names, canonical model references, edge-only shorthand, and clearer Area Breakdown report columns.
tags: [format, references, reports, cli, lint]
timestamp: 2026-06-22T00:00:00Z
---

# Model reference identifiers

This Change Case spec defines the delta for the QUALITY.md format,
qualitymd CLI, report rendering, and `/quality` skill: strict local names for
addressable model elements, canonical typed model references for Areas, Factors,
and Rating Levels, edge-only shorthand rules, and a clearer report summary Area
Breakdown table.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The Model needs two separate naming layers. Humans need required `title` values;
tools need stable references that can address a particular Area, Factor, or
Rating Level without guessing from display text or filesystem paths. Local map
keys are names, not necessarily full IDs. An Area is uniquely identified by its
path through the Area tree, and a Factor is uniquely identified only by the
declaring Area plus its nested Factor path.

This distinction also improves the concise report. The current Area Breakdown
table combines a title-bearing display path in one column and labels the
Area-only rating as `Area`. Separating Area title, model reference, local rating,
aggregate rating, and compact Factor ratings makes the summary easier to scan
and gives follow-up tools a stable handle.

## Scope

Covered:

- Strict name grammar for Area names, Factor names, and Rating Level IDs.
- Formal Area ID, Factor ID, and Rating Level ID definitions.
- Canonical typed model-reference syntax for Areas, Factors, and Rating Levels.
- Edge-only shorthand support where the expected reference type is fixed.
- Report summary Area Breakdown columns and root/nested row rendering.
- Preservation of structured `areaPath` and `factorPath` arrays in machine
  artifacts.

Deferred / non-goals:

- No relaxed or Unicode identifier grammar yet.
- No global uniqueness for Area names or Factor names.
- No strict grammar for Requirement statements.
- No replacement of structured report/evaluation JSON fields with string
  references.
- No general selector expression language, globbing, or multi-target syntax.

## Terminology

An **Area name** is a single map key under `areas`. It is unique among sibling
Areas in the same `areas` map.

An **Area ID** is the ordered Area path from the root Area to an Area. The root
Area ID is the empty path.

An **Area title** is the required human display label stored in `title`.

A **Factor name** is a single map key under `factors`. It is unique among
sibling Factors in the same `factors` map.

A **Factor ID** is the declaring Area ID plus the ordered Factor path from that
Area's `factors` map to the Factor.

A **Rating Level ID** is the `level` value of a Rating Level, unique within the
Model's `ratingScale`.

A **model reference** is the canonical typed text form used at human/tool
boundaries to address an Area, Factor, or Rating Level.

## Name Grammar

Area names, Factor names, and Rating Level IDs **MUST** match:

```regex
^[A-Za-z0-9](?:[A-Za-z0-9_-]*[A-Za-z0-9])?$
```

Requirement statement keys **MUST NOT** be constrained by this grammar. They
remain natural-language Requirement statements.

Tools **MUST** reject Area names, Factor names, and Rating Level IDs that do not
match the grammar.

The grammar intentionally excludes `/`, `:`, spaces, dots, and leading or
trailing separators so canonical model references are unambiguous and avoid
filesystem-path confusion.

## Canonical Model References

Canonical Area references use:

```text
area:<area-path>
```

The root Area reference is:

```text
area:root
```

Nested Area references join Area names with `/`:

```text
area:webhooks
area:webhooks/delivery
area:platform/runtime/scheduler/queues
```

Canonical Factor references use:

```text
factor:<declaring-area-path>::<factor-path>
```

The root declaring Area is written as `root`:

```text
factor:root::security
factor:root::security/secrets
```

Nested declaring Areas and nested Factors use `/` within each side of the
`::` separator:

```text
factor:webhooks::reliability
factor:webhooks/delivery::reliability/retry-behavior
factor:platform/runtime/scheduler/queues::operability/backpressure
```

Canonical Rating references use:

```text
rating:<rating-level-id>
```

For example:

```text
rating:target
rating:minimum
rating:unacceptable
```

Tools that render canonical model references **MUST** use the typed prefixes
`area:`, `factor:`, and `rating:`.

Tools that parse canonical model references **MUST** reject references whose
segments fail the strict name grammar or whose referenced model element does not
exist.

## Shorthand References

Tools **MAY** accept shorthand references that omit the type prefix only at
human/input edges where the expected reference type is fixed by the command,
field, UI control, or API parameter.

Examples:

```sh
qualitymd evaluate area webhooks/delivery
qualitymd evaluate factor webhooks/delivery::reliability/retry-behavior
```

Tools **MUST NOT** persist shorthand references in evaluation records,
`report.json`, generated reports, or other durable machine-readable artifacts.

Tools **MUST NOT** infer a shorthand reference type from filesystem state,
source paths, titles, or other contextual guesses.

Mixed-reference surfaces **MUST** require canonical typed model references.

Examples:

```sh
qualitymd evaluate --scope area:webhooks/delivery
qualitymd evaluate --scope factor:webhooks/delivery::reliability/retry-behavior
```

## Report Summary Area Breakdown

The concise report's Area Breakdown table **MUST** use these columns, in order:

```md
| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |
```

`Area` **MUST** render the Area title.

`Path` **MUST** render a stable Area model reference or an explicitly specified
Area-reference display form that is not confused with the Area title.

`Area Rating` **MUST** render the Area-only rating state. Structural Area groups
with child Areas but no direct requirements continue to render distinctly from
rated Areas.

`Area + Sub-Areas Rating` **MUST** render the single aggregate rating for the
Area including descendants.

`Factors` **MUST** preserve the compact Factor title/rating breakdown and the
explicit empty state when no Factor ratings were recorded.

Root and nested rows should read like:

```md
| Area                 | Path                     | Area Rating  | Area + Sub-Areas Rating | Factors                                     |
| -------------------- | ------------------------ | ------------ | ----------------------- | ------------------------------------------- |
| Sparrow Payments API | `area:root`              | Unacceptable | Unacceptable            | Security: Unacceptable; Reliability: Target |
| Delivery             | `area:webhooks/delivery` | Minimum      | Minimum                 | Reliability: Minimum                        |
```

The full report's shared compact Area Breakdown table **MUST** stay aligned with
the concise report's column semantics.

## Machine Artifacts

Evaluation records and report JSON **MUST** continue to preserve structured
`areaPath` and `factorPath` arrays.

Machine artifacts **MUST NOT** replace structured paths with display titles or
shorthand references.

Machine artifacts **MAY** add canonical model-reference strings as derived
fields if they preserve the structured path fields as the durable source of
truth.

## Acceptance Criteria

- `SPECIFICATION.md` defines Area name, Area ID, Area title, Factor name, Factor
  ID, Rating Level ID, and model reference.
- `SPECIFICATION.md` applies the strict name grammar to Area names, Factor
  names, and Rating Level IDs.
- `SPECIFICATION.md` explicitly excludes Requirement statement keys from the
  strict name grammar.
- `SPECIFICATION.md` defines canonical `area:`, `factor:`, and `rating:`
  references with root, nested Area, nested Factor, and Rating Level examples.
- `SPECIFICATION.md` defines edge-only shorthand support and forbids shorthand
  persistence in durable mixed-reference artifacts.
- Lint reports errors for invalid Area names, Factor names, and Rating Level
  IDs through named rule IDs.
- Report summary Area Breakdown renders the header
  `| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |`.
- Report summary separates Area title from stable Area reference for root and
  nested Areas.
- Structural Area groups and not-assessed states still render distinctly from
  Rating Levels.
- Areas with no recorded Factor ratings still render an explicit empty state.
- `report.json` continues to expose structured `areaPath` and `factorPath`
  arrays.
- CLI help and examples use "model reference" terminology where arguments might
  otherwise be confused with filesystem paths.
- Edge shorthand is accepted only where the command or field fixes the expected
  reference type.

## Durable spec changes

### To add

None.

### To modify

- `SPECIFICATION.md` - define strict names, Area IDs, Factor IDs, Rating Level
  IDs, canonical model references, and edge-only shorthand rules according to
  the requirements above.
- `specs/reports/report-summary-md.md` - update the Area Breakdown table
  contract and examples according to the report summary requirements above.
- `specs/reports/report-md.md` - align the full report's shared compact Area
  Breakdown table with the summary table requirements above.
- `specs/reports/report-json.md` - clarify structured path preservation and
  optional derived canonical references according to the machine artifact
  requirements above.
- `specs/evaluation-records/report-outputs.md` - align shared report-model
  terminology with Area IDs, Factor IDs, and canonical model references.
- `specs/cli/evaluation-create.md`, `specs/cli/evaluation-report.md`, and other
  discovered CLI selector specs - align scope/selector text with canonical
  model references and edge-only shorthand requirements above.
- `specs/cli/lint.md` - state that lint enforces the strict model-name grammar
  as part of model validation according to the name grammar requirements above.
- `specs/cli/lint-rules.md` - add or update diagnostics for strict Area name,
  Factor name, and Rating Level ID validation, including concrete rule IDs,
  according to the name grammar requirements above.
- `specs/quality-schema-json.md` - describe any structural JSON Schema pattern
  support for strict names if added by implementation.

### To rename

None.

### To delete

None.
