---
type: Design Doc
title: Report Count Summary Hierarchy — design doc
description: Design for consistent run-report full-list link placement and marked count summaries.
tags: [evaluation, reports, recommendations, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Report Count Summary Hierarchy — design doc

Design for
[Report Count Summary Hierarchy](../0186-report-count-summary-hierarchy.md) and
its [functional spec](spec.md).

## Context

`renderEvaluationRunReport` already owns the Top Findings and Top
Recommendations section order. Finding count summaries are produced from
`artifacts.rankedFindings()`, and recommendation totals already come from
`artifacts.rankedRecommendations()`. The implementation can stay in the report
renderer and fixed enum display catalog; no structured Evaluation data change is
needed.

The marker change for `strength` is broader than the inline count line because
Finding type display labels are shared by report tables, summaries, and
`glossary.md`.

## Approach

Keep the section order local to `renderEvaluationRunReport`:

```text
## Top Recommendations

full recommendations report link
top recommendations table
```

That mirrors the existing Top Findings section and keeps the complete-list path
above the capped preview table.

Update the Finding type enum catalog so `strength` uses marker `💪`. All
callers that use `findingTypeTitle` or catalog-derived glossary rows then pick
up `💪 Strength` without per-surface special cases.

Change the inline Finding summary builder to use catalog marker plus label for
Finding types and severity markers plus labels for gap/risk severity summaries.
The builder should continue to skip zero-count types and zero-count severities,
preserve the existing Finding type order, and pluralize only Finding type
labels:

```markdown
🚩 2 Gaps: 🔴 1 High, 🟡 1 Medium; 💪 4 Strengths
```

Add a recommendation impact summary builder beside the existing full-list link
helper. It should count ranked recommendations by `impact`, iterate the
recommendation impact catalog order, skip zero-count impact values, and render
one parenthetical shape:

```markdown
(3 total; impact: ⬥ 1 High, ● 2 Medium)
```

When no non-zero impact groups can be derived, keep the existing total-only
fallback:

```markdown
(3 total)
```

Tests should assert exact Markdown for representative mixed Finding summaries,
recommendation impact summaries, section ordering, and the `💪 Strength` display
label. Regenerate report-gallery output after renderer changes.

## Spec response

- Link placement is satisfied by moving the full Recommendations link before
  the table while leaving the Findings link before its table.
- Count hierarchy is satisfied by deriving marker/text labels from existing enum
  catalogs and using stable punctuation for group boundaries.
- The `impact:` axis label is produced only for recommendation count summaries,
  where `High` and `Medium` would otherwise be ambiguous.
- Structured data remains unchanged because the implementation edits display
  labels and report rendering only.

## Alternatives

- **Italicize severity labels.** Rejected because severity already has marker,
  label, count, and ordering signals; inline styling would add emphasis without
  adding structure.
- **Use `3 total by impact`.** Rejected because it reads awkwardly and can make
  the total sound impact-qualified instead of separating total from grouping
  axis.
- **Render `High impact` and `Medium impact` on every recommendation group.**
  Rejected as repetitive in a compact parenthetical.
- **Keep `✅ Strength`.** Rejected because the check mark overlaps with status
  and completion semantics. The selected `💪` marker is more distinctive for
  Strength.
- **Only add markers to inline summaries.** Rejected for `strength` because
  Finding type labels are shared display vocabulary; changing only one surface
  would make the glossary and tables disagree.

## Trade-offs & risks

The muscle-arm marker is more expressive than the quieter report markers used
elsewhere. That is intentional for this case because it was selected as the
Strength marker, but the implementation should still keep text labels beside it
so the marker is supplemental.

Recommendation summaries may become long if many impact levels are present.
There are only four fixed impact values, so the full inline summary remains
bounded.

## Open questions

None.
