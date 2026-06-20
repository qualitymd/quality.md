---
type: Design Doc
title: Readable report summary - design
description: Render the concise evaluation summary as a decision brief over the existing report model.
tags: [evaluation, report, cli, ux]
timestamp: 2026-06-19T00:00:00Z
---

# Readable report summary - design

Companion to [Readable report summary](../0040-readable-report-summary.md) and
its [functional spec](spec.md).

## Context

`report-summary.md` is already generated from `EvaluationReportDocument`, the
same in-memory model used for `report.md` and `report.json`. The requested change
is presentation-focused: improve the concise Markdown artifact's outline and
reader-facing vocabulary without changing report records, roll-up semantics, or
the JSON schema.

## Approach

Keep the existing report model and update only the concise Markdown renderer.
The render path remains:

```text
run records -> EvaluationReportDocument -> report-summary.md / report.md / report.json
```

`renderReportSummaryMarkdown` becomes the outline owner:

1. key-details table;
2. Summary, including headline prose and target rating table;
3. Top Issues;
4. Recommendations, including visible `Recommendation ID` values;
5. Scope & Limitations.

The renderer formats "Full evaluation" at display time when the run is not
narrowed, while preserving the recorded scope strings in `report.json`. It also
renames the human-facing headline result to "Overall rating" while leaving
`RatingResult` and aggregate-rating mechanics unchanged.

## Alternatives

- **Rename JSON fields.** Rejected because the user-facing issue is Markdown
  readability, and schema churn would break consumers without adding information.
- **Add a new `topIssues` JSON field first.** Rejected for this change because
  existing `findingSummaries`, `assessmentResults`, `recommendations`, and
  `targetSummary` already support the desired summary. A first-class top-issue
  field can be added later if multiple renderers need identical prioritization.
- **Redesign full `report.md` at the same time.** Rejected to keep the change
  focused. The full report remains the audit trail; the concise summary becomes
  the decision brief.

## Trade-offs & risks

The target summary table's "driver" still comes from existing aggregate-rating
rationales, which can be longer than ideal. This preserves traceability and
avoids new summarization judgment in the deterministic renderer, at the cost of
occasional wide Markdown rows.

## Open questions

None.
