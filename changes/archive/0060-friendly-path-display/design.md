---
type: Design Doc
title: Friendly path display - design doc
description: How display helpers stay separate from model-reference helpers while reports render / for the root Area path.
tags: [format, reports, references, display]
timestamp: 2026-06-22T00:00:00Z
---

# Friendly path display - design doc

## Context

This design answers the
[Friendly path display functional spec](spec.md). The change keeps qualified and
unqualified model references stable while letting human report path fields render
the root Area as `/`.

The important boundary is that display values are not references. The code
should make that boundary visible at call sites so a future report polish change
does not accidentally alter reference parsing or machine artifact identities.

## Approach

Separate Area path rendering into two internal concepts:

```go
AreaPath.Reference()            // area:root, area:webhooks/delivery
AreaPath.UnqualifiedReference() // root, webhooks/delivery
AreaPath.Display()              // /, webhooks/delivery
```

`Reference()` and `UnqualifiedReference()` should stop depending on
`Display()`. Add a small private helper for reference grammar:

```go
func (p AreaPath) referencePath() string {
    if len(p) == 0 {
        return "root"
    }
    return strings.Join(p, "/")
}
```

Then implement:

```go
func (p AreaPath) Reference() string {
    return "area:" + p.referencePath()
}

func (p AreaPath) UnqualifiedReference() string {
    return p.referencePath()
}

func (p AreaPath) Display() string {
    if len(p) == 0 {
        return "/"
    }
    return strings.Join(p, "/")
}
```

Apply the same shape to Factor paths, even though the visible output is not
expected to change for non-empty Factor paths:

```go
FactorPath.Display()      // security/secrets
FactorPath.referencePath() // root for an empty path if reference grammar needs it
```

`FactorReference` and `UnqualifiedFactorReference` should compose the declaring
Area's reference grammar with the Factor path's reference grammar, not either
type's display value. That preserves `factor:root::security` and
`root::security` if Area display changes to `/`.

Rating Level IDs do not currently have a dedicated type. Keep the existing
`RatingReference(level)` and `UnqualifiedRatingReference(level)` functions, and
add a small `RatingDisplay(level)` function only if a call site benefits from the
symmetry. If no call site uses it, tests can still assert the existing reference
helpers and leave rating display documented but not separately implemented.

Update report rendering to use display values in presentation-oriented path
fields:

- the shared Area Breakdown `Path` column should call `AreaPath.Display()`;
- the full report Area Detail `Path` field should call `AreaPath.Display()`;
- title-aware headings and labels should continue using `reportDisplayLabels`;
- `report.json` should remain unchanged.

Parsing stays unchanged. `ParseAreaReference` continues to require `area:`.
`ParseUnqualifiedAreaReference` continues to accept `root` for the root Area and
must reject `/`.

## Alternatives

**Change `UnqualifiedReference()` to return `/` for root.** Rejected because `/`
would become part of fixed-type reference grammar, not just report display. That
would force parser changes or make render/parse behavior asymmetric.

**Add report-only helper functions and leave `Display()` as `root`.** Rejected
because the code already has `Display()` as the natural non-reference method on
path types. Keeping it as `root` while adding another display helper would leave
the naming confusing.

**Change `AreaPath.Display()` and let references keep calling it.** Rejected
because it would turn `Reference()` into `area:/` and
`UnqualifiedReference()` into `/` unless every reference call site was audited.
The safer design makes reference grammar independent of display.

**Make display title-aware inside path types.** Rejected because title-aware
rendering belongs to `reportDisplayLabels`, which has access to the model
snapshot. Path types should remain structural fallback helpers.

## Trade-offs & Risks

Changing `AreaPath.Display()` from `root` to `/` may affect fallback text outside
the report rows covered by this change. The implementation should search all
`Display()` call sites and update tests where the friendlier display is intended.
If a call site needs reference grammar, it should switch to `Reference()` or
`UnqualifiedReference()`.

`/` is familiar as a root marker, but it can also look like a filesystem path.
That is acceptable in human report display because the surrounding labels say
Area path and machine artifacts keep structured identifiers. The spec keeps `/`
out of reference parsing to avoid ambiguity.

Keeping Factor and Rating display mostly identical to their unqualified
reference forms may feel like extra ceremony. The value is consistency: helpers
document intent at call sites even when current strings happen to match.

## Open Questions

- Should `RatingDisplay(level)` be added immediately for symmetry, or deferred
  until a real call site needs it?
- Should the Scope summary's reconstructed root Area text use `/`, stay
  title-oriented, or continue with its current prose?
