---
type: Design Doc
title: Evaluation report rating labels - design
description: How the Evaluation v2 report builder renders explicit rating column labels and a local-plus-roll-up breakdown.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation report rating labels - design

## Context

Answers the
[Evaluation report rating labels functional spec](spec.md). All work is in the
deterministic Markdown projection in `internal/evaluation/report_v2.go`; no
structured data, schema, or path logic changes. The scoped analysis maps the
builder needs are already in hand at every call site —
`scopedMap(node.Analysis, "localAnalysis")` and
`scopedMap(node.Analysis, "localAndDescendantAnalysis")` — so this is a
relabel-and-resource change, not new plumbing.

## Approach

**Header tables.** In `renderV2AreaReport` and `renderV2FactorReport` the header
rows already render `overall` then `local` via `v2RatingLabel`; only the literal
header strings change: `| Overall | Local | …` becomes
`| Overall Rating | Local Rating | …`. The cells are already correct.

**Breakdown tables.** Three rows are reshaped — the Area report Factors table,
the Area report Sub-Areas table, and the Factor report child Factors table. Each
row currently renders `v2RatingLabel(spec, localAndDescendantAnalysis)` in the
`Rating` column and `v2BoolLabel(len(children) > 0)` in the `+ Sub-X` column. Each
becomes:

- `Local Rating` ← `v2RatingLabel(spec, scopedMap(child.Analysis, "localAnalysis"))`
- `+ Sub-X Rating` ← a new helper that returns the descendant-inclusive rating
  when the node has descendants and an em dash otherwise:

```go
func v2SubRatingCell(spec *model.Spec, aggregate map[string]any, hasDescendants bool) string {
	if !hasDescendants {
		return "—"
	}
	return v2RatingLabel(spec, aggregate)
}
```

The `hasDescendants` bool reuses the exact predicate the boolean column used
(`len(artifacts.childFactors(id)) > 0` / `len(artifacts.childAreas(id)) > 0`), so
the "has children" signal the boolean carried is preserved as the em-dash branch.

The child-Factors table header gains a column (`| Factor | Path | Rating |` →
`| Factor | Path | Local Rating | + Sub-Factors Rating |`), so its empty-state
row gains one trailing cell. The Factors and Sub-Areas tables keep their column
count, so their empty-state rows are unchanged.

**Dead code.** `v2BoolLabel` was used only by the two reshaped rows. After the
change nothing in the report builder renders a boolean, so the wrapper is
removed and the one display-title test that exercised it points at the
underlying `boolTitle`, which stays as the shared boolean display-title helper.

## Spec response

- [Header rating labels](spec.md#header-rating-labels): satisfied by the two
  header-string edits; cells already sourced correctly.
- [Breakdown rating columns](spec.md#breakdown-rating-columns): satisfied by
  re-sourcing the local column from `localAnalysis` and the sub column from
  `localAndDescendantAnalysis` via `v2SubRatingCell`, whose `—` branch covers the
  no-descendants requirement and whose existence removes every boolean rating
  cell (the 0097 carry-over).
- Verification: the navigation report test asserts each new header and a
  representative row. To prove the local and sub columns read from different
  sources, the navigation fixture gains a second rating level and sets one
  Factor's `localAndDescendantAnalysis` to that level so the two columns render
  visibly different labels; a no-descendant row asserts the `—` cell.

## Alternatives

- **Label the sub column `Overall Rating` too**, matching the header's term for
  the same `localAndDescendantAnalysis` quantity. Rejected: in a breakdown row
  `+ Sub-Factors Rating` says exactly what is rolled in, matches the original
  sketch naming, and reads better than "overall" for a child listed under its
  parent. The header keeps `Overall Rating` because there the node is the
  subject.
- **Repeat the local rating in the sub column when there are no descendants.**
  Rejected: it presents a redundant rating and loses the "has children" signal;
  the em dash keeps both honest.
- **Keep `v2BoolLabel`** as latent capability. Rejected: it would be dead
  production code kept alive only by its test; `boolTitle` already preserves the
  boolean display-title capability the spec enumerates.

## Trade-offs & risks

- The em-dash convention introduces a non-rating glyph into rating columns; it is
  documented in the durable spec so it is not mistaken for a missing value
  (distinct from the `not_analyzed` / `blocked` status labels, which still render
  as their own titles when analysis was not produced).
- Low risk overall: deterministic Markdown only, covered by the report tests; no
  structured output or path change.

## Open questions

None.
