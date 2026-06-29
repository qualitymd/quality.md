---
type: Design Doc
title: Report Full List Links — design doc
description: Design for emphasized run-report full-list links with total ranked counts.
tags: [evaluation, reports, recommendations, findings]
timestamp: 2026-06-29T00:00:00Z
---

# Report Full List Links — design doc

Design for [Report Full List Links](../0177-report-full-list-links.md) and its
[functional spec](spec.md).

## Context

`renderEvaluationRunReport` already renders the capped Top Findings and Top
Recommendations sections and has access to the full ranked Evaluation artifacts.
The generated full-list line can stay local to the run-report renderer, but the
count must come from the same ranked artifact lists used by the complete list
reports rather than from the capped table writer.

## Approach

Add a small writer for complete-list links. The run-report renderer can pass the
total length of `artifacts.rankedFindings()` or
`artifacts.rankedRecommendations()`, then write the full-list line in one
deterministic Markdown shape:

```markdown
**Full recommendations report:** [recommendations.md](recommendations.md) (23 total)
```

This keeps the visible filename link unchanged while making the label bold and
sentence-case. The same helper handles Findings and Recommendations so their
lines stay symmetrical.

Tests should assert the exact generated Markdown for both lines. Existing
report-gallery generation then refreshes the checked-in example output.

## Spec response

- The run-report renderer emits the emphasized sentence-case labels required by
  the spec.
- Counts come from complete ranked artifact lists, not from the capped
  table writers.
- Link text remains the literal report filename, and the count is rendered
  outside the link text.

## Alternatives

- **Put the count in link text.** Rejected because the report artifact path is
  the useful scan target in plain Markdown.
- **Only show counts when more than 10 rows exist.** Rejected because a stable
  line shape is easier to scan and test, and the count remains useful when it is
  less than or equal to the cap.
- **Add a separate callout block.** Rejected as too heavy for a local link below
  an already visible table.

## Trade-offs & risks

The line adds a small amount of repeated count information because Key Details
already includes ranked counts. That duplication is acceptable here: Key Details
orients the whole report, while the full-list line explains the local path from a
capped table to its complete report.

## Open questions

None.
