---
type: Design Doc
title: Report empty-cell marker and legend
description: How the report renderer marks empty scalar cells with an em dash and appends one static per-report legend.
status: Draft
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Report empty-cell marker and legend

## Context

Implements the [0118 functional spec](spec.md): every empty scalar cell in a
generated Evaluation report renders as `—`, composite pair cells render `— / —`
when a component is absent, and each report carries one static legend defining
the marker. All rendering lives in `internal/evaluation/report_tree.go`, with the
enum/title helpers in `internal/evaluation/display.go`. This doc decides *where*
the empty-cell policy lives so the marker lands once, at the right altitude,
without a wide blast radius.

## Approach

The renderer already routes nearly every table cell through one function,
`markdownCell` (report_tree.go:1199) — its only job is to escape pipes for a
Markdown table cell. That makes it the natural owner of the empty-cell policy.

**1. Cell chokepoint — `markdownCell`.** Change its empty guard to return the em
dash instead of a blank:

```go
func markdownCell(s string) string {
    if s == "" {
        return "—"
    }
    return strings.ReplaceAll(s, "|", "\\|")
}
```

This covers every plain scalar cell in one edit. Cells that pass through a title
helper first already collapse to `""` on absence and then hit this guard — e.g.
`findingSeverityTitle(firstString(finding, "severity", …))` → `stringTitle("")` →
`humanizeEnum("")` → `""` → `markdownCell` → `—`. Same for finding type,
description, location/evidence (`compactJSON(nil)` → `""`), rationale, actions,
rating drivers, limits, and unknowns. The existing roll-up `—` in
`evaluationSubRatingCell` (report_tree.go:1000) already matches this glyph, so
the two unify rather than collide.

**2. Composite pair cells — fix the builders.** Three cells assemble two scoped
values *before* `markdownCell` sees them, so an all-empty pair reaches the guard
as `" / "` (non-empty) and slips through:

- `evaluationConfidencePair` (report_tree.go:1174) — Area and Factor headers;
- `evaluationAnalysisStatusPair` (report_tree.go:1170) — Factor header `Status`;
- the inline Requirement-header confidence pair (report_tree.go:454).

Introduce one small helper that em-dashes an empty component, and route all three
through it:

```go
func emDashIfEmpty(s string) string {
    if s == "" {
        return "—"
    }
    return s
}

func evaluationConfidencePair(overall, local map[string]any) string {
    return emDashIfEmpty(confidenceTitle(evaluationString(overall, "confidence"))) +
        " / " + emDashIfEmpty(confidenceTitle(evaluationString(local, "confidence")))
}
```

For consistency, extract the inline Requirement confidence pair (line 454) into a
named builder (e.g. `evaluationRequirementConfidencePair(req)`) that uses the
same helper, so no pair cell hand-rolls the `a + " / " + b` join.

`humanizeEnum` stays untouched — it remains a pure humanizer returning `""` for
empty, and the cell layer owns the marker. This keeps the blast radius to "table
cells," not "every enum title in the system."

**3. Per-report legend.** Add a writer and call it last in each renderer:

```go
func writeEvaluationLegend(b *strings.Builder) {
    b.WriteString("## Legend\n\n")
    b.WriteString("- `—` — not applicable or not recorded.\n")
}
```

Call it at the end of `renderEvaluationAreaReport`, `renderEvaluationFactorReport`,
and `renderEvaluationRequirementReport`. The root `report.md` is an Area report
(`renderEvaluationReportTree` builds it via `renderEvaluationAreaReport`), so all
four report kinds are covered. The writer takes no data and emits unconditionally,
satisfying the static-legend requirement.

## Spec response

- **Empty scalar cells / pair cells** — covered by (1) and (2). The shared
  `emDashIfEmpty` helper is the verification anchor: an all-empty confidence or
  status header renders `— / —`, asserted from a fixture with no confidence.
- **Distinctness from outcomes** — untouched by design. `evaluationRatingLabel`,
  `evaluationRequirementRatingLabel`, `assessmentStatusTitle`, etc. resolve a
  missing *rating/assessment outcome* to its own label (`⚪ Not Rated`,
  `⚪ Not Analyzed`, `⛔ Blocked`) *before* the value reaches `markdownCell`, so a
  non-empty status string flows through the guard unchanged and never becomes
  `—`. The em dash only appears where the upstream value is genuinely empty.
- **Section rows unchanged** — the `(no findings)` / `(no rating drivers)` /
  `(none recorded)` rows are written as literal strings, not via `markdownCell`,
  so they are untouched; their trailing padding cells stay literal-blank, which
  reads as part of the worded empty-state row rather than per-cell absence.
- **Legend presence + static** — covered by (3); the writer is unconditional and
  appended to every renderer.

## Alternatives

- **Literal `N/A`** — rejected. "Not applicable" overclaims; the absence is
  usually "not recorded." A neutral dash states absence without asserting why.
- **Per-cell asterisk + footnote** (the original idea) — rejected. A Markdown
  table has no stable "first instance" to anchor the asterisk, and the anchor
  migrates as data changes, churning committed-report diffs on every re-run. The
  static legend is the fix for exactly that.
- **Single chokepoint at `humanizeEnum`** — rejected. It would em-dash every
  empty enum title system-wide (status badges, data-kind titles, run-gap titles),
  a far wider blast radius than "table cells," *and* still wouldn't catch raw
  string cells (finding description, `compactJSON` evidence) that never pass
  through a title helper — so `markdownCell` would need the guard anyway. Two
  chokepoints, the wider one unnecessary.
- **Dedicated `cellOrDash` helper at every cell site** instead of overloading
  `markdownCell` — rejected as churn. `markdownCell` already *means* "render a
  table cell"; its two non-table uses (a link label at report_tree.go:819, the
  `Name:` backtick at :451) only ever receive non-empty input, so overloading it
  is safe and touches one line instead of ~25 sites that could be missed.
- **Conditional legend** (only when an em dash is present) — rejected. It
  reintroduces the appear/disappear diff churn the static legend exists to avoid.
- **En dash vs em dash** to split "not recorded" from "not applicable" — rejected
  in [discussion](spec.md#background--motivation); the two are visually
  indistinguishable and the distinction rarely matters to a reader.

## Trade-offs & risks

- **`markdownCell` now owns empty-cell policy** alongside pipe-escaping. Its two
  non-table callers would also em-dash an empty input, but those inputs (report
  link label, Requirement `Name`) are never empty in practice; if that ever
  changes, surfacing the marker there is acceptable, not wrong.
- **The em dash carries two senses** — roll-up-not-applicable (existing
  `evaluationSubRatingCell`) and empty-cell. The legend gloss "not applicable or
  not recorded" covers both; both mean "nothing to render here," so the overload
  is coherent rather than ambiguous.
- **Legend boilerplate on every report**, including fully-populated ones — the
  accepted cost of diff stability over a data-conditional legend.
- **Pair-builder miss risk** — the inline Requirement confidence pair (line 454)
  is easy to overlook. Mitigation: route all three pairs through `emDashIfEmpty`
  and extract line 454 into a named builder so no pair hand-rolls the join.

## Open questions

- **Legend format** — a `## Legend` section (chosen, for consistency with the
  report's other `##` sections and easy assertion) versus a lighter italic
  footnote after a horizontal rule. Low-stakes; settle at implementation if the
  heading reads too heavy on short reports.
