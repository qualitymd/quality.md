---
type: Design Doc
title: Unqualified model references - design doc
description: How unqualified reference helpers, parsing boundaries, and Area Breakdown rendering should be implemented.
tags: [format, references, reports, cli]
timestamp: 2026-06-22T00:00:00Z
---

# Unqualified model references - design doc

## Context

This design answers the
[Unqualified model references functional spec](spec.md). It builds on 0058's
canonical typed references while adding a shorter form for fixed-type contexts,
starting with the Area Breakdown `Path` column.

The design goal is to make the semantic distinction visible at call sites:
qualified references are self-describing; unqualified references depend on the
surrounding contract for their type.

## Approach

Use named helpers, not boolean flags:

```go
AreaPath.Reference()
AreaPath.UnqualifiedReference()

FactorReference(areaPath, factorPath)
UnqualifiedFactorReference(areaPath, factorPath)

RatingReference(level)
UnqualifiedRatingReference(level)
```

The qualified helpers keep the current 0058 behavior:

```text
area:operations/incident-response
factor:operations/incident-response::operability
rating:target
```

The unqualified helpers omit only the typed prefix:

```text
operations/incident-response
operations/incident-response::operability
target
```

For root paths, reuse the existing `root` token:

```text
root
root::security
```

Do not add an Area-local Factor helper in this change. A form such as
`operability/backpressure` is useful only when a surface fixes both the type and
the declaring Area. No current surface needs that extra contraction, and using
`UnqualifiedFactorReference` for two different shapes would make code harder to
review.

Parsing should mirror the render boundary. Keep existing canonical parsing
strict for qualified references. Add or adjust type-specific parse paths only
where the expected type is explicit, for example an Area-only helper can accept
`root` or `operations/incident-response`, while a mixed `--scope`-style parser
must still require `area:`, `factor:`, or `rating:`.

Update the shared Area Breakdown renderer to call
`AreaPath.UnqualifiedReference()` for the `Path` column. Leave `report.json`
unchanged; it already preserves structured path arrays and should not gain a
string dependency on unqualified references.

## Alternatives

**Use `Reference(prefix bool)`.** Rejected because boolean arguments hide the
semantic choice at call sites. `Reference(false)` is shorter, but a reader must
look up whether `false` means no prefix, no qualification, no validation, or
something else.

**Rename canonical references to qualified everywhere.** Rejected because 0058
already established canonical typed references as the self-describing form.
This change can introduce "qualified" as the contrast term without churning the
existing public vocabulary.

**Implement a generic reference object.** Rejected as premature. Area, Factor,
and Rating references have different structure, and the current code already has
typed paths and small helper functions.

**Render Area Breakdown with qualified references and rely on users to ignore
the prefix.** Rejected because the table is human-first and Area-specific. The
prefix is redundant in that column and makes the report harder to scan.

## Trade-offs & Risks

Adding unqualified parsing creates a risk of accidental guessing. Keep the parse
entrypoints type-specific and avoid adding any mixed parser that accepts
unqualified input.

The word "unqualified" is precise but still technical. Durable specs should
define it once and examples should carry the meaning; user-facing prose can say
"Area path" where the context is already clearly Area-specific.

Factor unqualified references preserve the declaring Area side, so they are not
the shortest possible Factor form. This is deliberate: dropping the declaring
Area requires a stronger context guarantee and should get a separate helper when
needed.

## Open Questions

- Should a future Area-fixed Factor surface add an explicit
  `FactorPath.UnqualifiedReference()` or `LocalFactorReference()` helper for the
  `<factor-path>`-only form?
- Should CLI help eventually expose separate examples for qualified mixed-scope
  selectors and unqualified fixed-type selectors?
