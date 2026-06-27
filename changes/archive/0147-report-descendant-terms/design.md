---
type: Design Doc
title: Report descendant terms design
description: Implementation approach for generated report Child Areas and Sub-Factors labels.
tags: [evaluation, reports, terminology]
timestamp: 2026-06-27T00:00:00Z
---

# Report Descendant Terms Design

## Context

This design answers the [Report descendant terms functional spec](spec.md).
Generated Markdown reports already compute the correct immediate descendants and
roll-up ratings. The change is to render those relationships with the intended
human-facing vocabulary.

## Approach

Keep the data flow and report tree unchanged. Edit only the literal Markdown
labels emitted by `internal/evaluation/report_tree.go`:

- Area reports render the descendant-Area section as `## Child Areas`.
- Area report descendant-Area tables use `+ Child Areas Rating`.
- Factor reports render descendant Factors as `## Sub-Factors`.
- Empty-state rows use `(no Child Areas)` and `(no Sub-Factors)`.

The existing `evaluationSubRatingCell` helper remains appropriate because it
models a descendant-inclusive rating cell for both relationship types. Its
comment should name the new rendered column label.

Durable specs absorb the same vocabulary so future report changes use the
contract rather than rediscovering the renderer strings.

## Spec Response

The renderer string edits satisfy the generated-label requirements without
touching structured data or traversal helpers. The existing navigation report
test already covers the report tree with a root Area, a Child Area, a Factor, and
a Sub-Factor; expanding its assertions is enough to prove the new labels and
guard against the old labels returning.

## Alternatives

- Rename internal helpers such as `childFactors`.
  Rejected because they are implementation vocabulary and not rendered report
  labels. Renaming them would create mechanical churn without improving the
  user-facing report.
- Introduce a label helper for every relationship term.
  Rejected because the current report labels are simple literals and the change
  touches only a few strings. A helper would add indirection without reducing
  real duplication.

## Trade-offs & Risks

The change is intentionally narrow: historical generated reports in existing run
folders will keep their old wording until rebuilt. That is acceptable for a
report-rendering wording fix.

## Open Questions

None.
