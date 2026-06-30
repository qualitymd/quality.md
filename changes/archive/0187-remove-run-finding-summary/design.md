---
type: Design Doc
title: Remove Run Finding Summary - design
description: Renderer approach for removing the redundant run-level Finding Summary table.
tags: [evaluation, reports, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Remove Run Finding Summary - design

Companion to
[Remove Run Finding Summary](../0187-remove-run-finding-summary.md) and its
[functional spec](spec.md).

## Context

The run report renderer writes the opening in `renderEvaluationRunReport`: H1,
`Evaluation links:`, `Summary`, `Key Details`, the standalone `Finding Summary`
table, `Contents`, `Model Evaluation`, then Top Findings and Top Recommendations.
The full Findings report link under `## Top Findings` now carries the complete
ranked Finding count and inline type/severity summary, which makes the earlier
standalone table duplicative.

## Approach

Remove the `writeRunReportFindingSummary` call from the run report opening and
delete the now-unused helper types and functions that only serve that table.
Leave `writeRunReportKeyDetails` unchanged so the opening still carries the
total Finding count. Leave `writeFullFindingsReportLink` and
`findingInlineCountSummary` unchanged so the richer linked count breakdown stays
beside the capped Top Findings table.

Requirement report rendering uses a separate `## Findings Summary` path and
should not change.

## Spec response

The approach satisfies the run-report removal requirement by eliminating the
only renderer call that writes the standalone table. It satisfies the preserved
count requirements by leaving Key Details and full Findings link rendering in
place. It satisfies the Requirement report non-goal by not touching the
Requirement report summary renderer.

## Alternatives

- **Keep the table but move it under Top Findings.** Rejected because the full
  Findings report link already carries the same summary there; two adjacent
  count summaries would be worse than the current split.
- **Remove all Finding count summaries from `report.md`.** Rejected because the
  Key Details total and full-list link summary serve distinct scan tasks: quick
  opening metadata and complete-list context.
- **Hide zero-count rows only.** Rejected because the redundancy is the whole
  standalone table, not just its zero rows.

## Trade-offs and risks

Readers lose one explicit zero-count taxonomy table in the run report opening.
That is intentional: zero-count absence is still inferable from the inline full
Findings summary omitting absent groups, while the full Findings report remains
the complete ranked list.

Tests and gallery output need exact Markdown updates because the opening stack
and surrounding section adjacency change.

## Open questions

None.
