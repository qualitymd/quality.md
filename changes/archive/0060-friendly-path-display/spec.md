---
type: Functional Specification
title: Friendly path display
description: Separate display values from model references and render root Area report paths as /.
tags: [format, reports, references, display]
timestamp: 2026-06-22T00:00:00Z
---

# Friendly path display

This Change Case spec defines the delta for QUALITY.md report rendering,
qualitymd evaluation helpers, and `/quality` guidance: display values are
human-facing fallback/rendering values distinct from qualified and unqualified
model references, and human report paths render the root Area as `/`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

Qualified and unqualified references solve addressability. They must stay
stable enough for parsing, mixed-reference surfaces, and future tool input:
`area:root` and `root` are reference grammar. Human reports have a different
job. In a report `Path` column or Area Detail path field, `/` is a clearer
display value for the root Area than `root`, and it matches a familiar visual
convention for "the top of this tree."

Separating display from reference grammar keeps this report improvement from
turning `/` into a new reference spelling or changing machine artifacts. It also
keeps future code reviews simple: reference helpers return references; display
helpers return human-facing fallback text.

## Scope

Covered:

- Display value terminology.
- Area display behavior, including `/` for the root Area.
- Preservation of qualified and unqualified reference grammar.
- Consistent display/reference separation for Factor paths and Rating Level IDs.
- Human report `Path` rendering for `report-summary.md` and `report.md`.
- Preservation of structured machine artifacts.

Deferred / non-goals:

- No acceptance of `/` as an Area reference.
- No change to `area:root`, `root`, `factor:root::security`, or
  `rating:target` reference forms.
- No title-aware display inside path or rating helper functions.
- No new report JSON fields.
- No replacement of structured `areaPath`, `factorPath`, or rating `level`
  fields.

## Terminology

A **qualified reference** is the self-describing model-reference form with a type
prefix, such as `area:root`, `factor:root::security`, or `rating:target`.

An **unqualified reference** is the context-bounded model-reference form without
a type prefix, such as `root`, `root::security`, or `target`.

A **display value** is human-facing fallback or rendering text. A display value
is not necessarily parseable as a model reference.

## Requirements

Area display values **MUST** render the root Area as:

```text
/
```

Nested Area display values **MUST** join Area names with `/`:

```text
webhooks/delivery
```

Area qualified references **MUST** keep the existing grammar:

```text
area:root
area:webhooks/delivery
```

Area unqualified references **MUST** keep the existing grammar:

```text
root
webhooks/delivery
```

Tools **MUST NOT** accept `/` as a qualified or unqualified Area reference in
this change.

> Rationale: `/` is a friendly report display value, not a new spelling in the
> reference grammar. Keeping that distinction avoids making report polish change
> input parsing or mixed-reference behavior. — 0060

Factor path display values **SHOULD** follow the same display/reference
separation. Non-empty Factor path display values may continue to render as
slash-joined Factor names:

```text
security/secrets
```

Factor qualified and unqualified references **MUST** keep the existing grammar:

```text
factor:root::security
root::security
```

Rating Level display values **SHOULD** follow the same display/reference
separation. A Rating Level ID display value may continue to render as the Rating
Level ID:

```text
target
```

Rating qualified and unqualified references **MUST** keep the existing grammar:

```text
rating:target
target
```

Human Markdown reports **SHOULD** use display values for presentation-oriented
path fields, including the Area Breakdown `Path` column and full report Area
Detail `Path` field.

Human Markdown reports **MUST NOT** use display values where a qualified or
unqualified model reference is explicitly required.

`report.json` and evaluation records **MUST** remain structured and **MUST NOT**
persist display path strings.

Title-aware display remains owned by report label resolution. Path and rating
display helpers **MUST NOT** resolve titles by themselves.

## Acceptance Criteria

- `AreaPath{}.Display()` returns `/`.
- `AreaPath{}.Reference()` still returns `area:root`.
- `AreaPath{}.UnqualifiedReference()` still returns `root`.
- `AreaPath{"webhooks", "delivery"}.Display()` returns `webhooks/delivery`.
- `AreaPath{"webhooks", "delivery"}.Reference()` still returns
  `area:webhooks/delivery`.
- `AreaPath{"webhooks", "delivery"}.UnqualifiedReference()` still returns
  `webhooks/delivery`.
- `FactorReference(AreaPath{}, FactorPath{"security"})` still returns
  `factor:root::security`.
- `UnqualifiedFactorReference(AreaPath{}, FactorPath{"security"})` still returns
  `root::security`.
- `FactorPath{"security", "secrets"}.Display()` returns `security/secrets`.
- `RatingReference("target")` still returns `rating:target`.
- `UnqualifiedRatingReference("target")` still returns `target`.
- A Rating Level display helper, if introduced, returns `target` for `target`.
- Generated `report-summary.md` Area Breakdown root rows render the `Path`
  column as `` `/` ``.
- Generated `report.md` Area Breakdown root rows render the `Path` column as
  `` `/` ``.
- Generated `report.md` Area Detail root `Path` renders `/`.
- Nested Area paths in human reports still render without `area:` prefixes, for
  example `webhooks/delivery`.
- `report.json` continues to expose structured `areaPath` and `factorPath`
  arrays and does not add display path strings.
- Existing qualified and unqualified parsing tests reject `/` as an Area
  reference.
- `go test ./...` and `mise run check` pass.

## Durable spec changes

### To add

None.

### To modify

- `SPECIFICATION.md` - distinguish display values from qualified and
  unqualified model references, keep reference grammar unchanged, and define
  root Area display as `/` according to the display and reference requirements
  above.
- `specs/reports/report-summary-md.md` - require friendly Area display values in
  the Area Breakdown `Path` column according to the human report requirements
  above.
- `specs/reports/report-md.md` - keep the full report's shared Area Breakdown
  and Area Detail path fields aligned with the summary behavior above.
- `specs/reports/report-json.md` - clarify that report JSON preserves
  structured paths and does not persist display path strings according to the
  machine artifact requirement above.
- `specs/evaluation-records/report-outputs.md` - align shared report-output
  terminology with display values versus references while preserving structured
  machine identifiers according to the requirements above.
- `specs/skills/quality-skill/reporting.md` and
  `specs/skills/quality-skill/quality-skill.md` - align durable skill guidance
  with display values versus references according to the human report
  requirements above.

### To rename

None.

### To delete

None.
