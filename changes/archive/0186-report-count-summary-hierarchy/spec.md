---
type: Functional Specification
title: Report Count Summary Hierarchy — functional spec
description: Requirements for consistent run-report full-list link placement and semantic count summary markers.
tags: [evaluation, reports, recommendations, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Report Count Summary Hierarchy — functional spec

Companion to
[Report Count Summary Hierarchy](../0186-report-count-summary-hierarchy.md).
This spec states the delta for generated `report.md` Top Findings and Top
Recommendations section links and count summaries. The durable source of truth
is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

`report.md` is a capped overview. Its Top Findings and Top Recommendations
tables help readers triage quickly, while `findings.md` and
`recommendations.md` hold the complete ranked lists. The full-list links should
therefore sit before each capped table so readers understand the table's
relationship to the complete report before scanning rows.

The count summaries beside those links are also miniature taxonomies. They need
enough hierarchy to show which labels are Finding types, which nested labels are
gap/risk severities, and which recommendation labels are impact levels. Markers
help scanning only when they remain paired with text and preserve the underlying
semantic axis.

## Scope

Covered: generated Markdown for run-level `report.md`, fixed Finding type
display labels used by generated reports and `glossary.md`, durable report
specs, reporting guidance, tests, changelog, and report-gallery output.

Non-goals:

- changing structured Evaluation enum values such as `strength`, `gap`,
  `risk`, `note`, or recommendation impact values;
- changing the 10-row cap for Top Findings or Top Recommendations;
- changing ranked-list report table columns, sorting, or filtering;
- changing generated report filenames or navigation blockquotes.

## Requirements

1. Generated run-level `report.md` **MUST** render the full Recommendations
   report link immediately after the `## Top Recommendations` heading and before
   the capped Top Recommendations table.

   > Rationale: The complete-list link explains the capped table before the
   > reader starts scanning preview rows, matching the Top Findings section.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

2. Generated run-level `report.md` **MUST** keep the full Findings report link
   immediately after the `## Top Findings` heading and before the capped Top
   Findings table.

   > Rationale: This preserves the already-good placement while making the
   > section symmetry explicit.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

3. When ranked Findings exist, the full Findings report count summary **MUST**
   render Finding type groups in catalog order, each non-zero group **MUST**
   include the Finding type marker and text label, and gap/risk groups with
   observed severities **MUST** render severity groups after a colon.

   > Rationale: Type, count, and severity are different hierarchy levels; marker
   > plus text and colon nesting make that structure visible without inline
   > styling.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

4. The full Findings report count summary **MUST** separate Finding type groups
   with semicolons and gap/risk severity groups with commas.

   > Rationale: Stable punctuation gives the inline summary two readable levels:
   > type groups outside, severity groups inside.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

5. The fixed Finding type display label for `strength` **MUST** render as
   `💪 Strength` in generated human-facing Evaluation surfaces and the generated
   glossary.

   > Rationale: The check mark marker overlaps with completion and status
   > semantics already used elsewhere; the selected muscle-arm marker gives
   > Strength its own visible identity while leaving the persisted value
   > unchanged.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

6. When ranked Recommendations exist, the full Recommendations report count
   summary **MUST** include an explicit lowercase `impact:` axis label after the
   total count, then render non-zero recommendation impact groups in impact
   catalog order with marker, count, and text label.

   > Rationale: Labels such as `High` and `Medium` are ambiguous without the
   > axis; `impact:` states what is being grouped without repeating "impact" on
   > every item.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

7. The full Recommendations report count summary **MUST** separate the total
   count from the impact summary with a semicolon, and **MUST** separate impact
   groups with commas.

   > Rationale: Recommendations have one grouping level, so commas are enough
   > inside the impact summary while the semicolon separates the total from the
   > axis label.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

8. Generated count summaries **MUST NOT** put total counts, type summaries,
   severity summaries, or impact summaries inside the report filename link text.

   > Rationale: The link text should remain the target artifact filename, while
   > counts describe the linked report.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

9. The change **MUST NOT** alter structured Evaluation data enum values,
   validation, or persisted JSON payload shapes.

   > Rationale: This is a human Markdown display change; existing data contracts
   > and run artifacts should remain compatible with the current early-alpha
   > contract.
   > Durable spec: none.

## Examples

```markdown
**Full findings report:** [findings.md](findings.md) (7 total: 🚩 2 Gaps: 🔴 1 High, 🟡 1 Medium; ⚠️ 1 Risk: 🟡 1 Medium; 💪 4 Strengths)

**Full recommendations report:** [recommendations.md](recommendations.md) (3 total; impact: ⬥ 1 High, ● 2 Medium)
```

## Acceptance criteria

- `report.md` renders each full-list link before its capped overview table.
- Finding inline summaries include type markers, severity markers for gap/risk
  severities, colon nesting, semicolon type separators, and comma severity
  separators.
- Recommendation inline summaries render `(N total; impact: ...)` with impact
  markers and comma-separated non-zero impact groups.
- The shared Finding type display for `strength` is `💪 Strength` in generated
  reports and `glossary.md`.
- Focused Evaluation renderer tests pass.
- Report-gallery output is regenerated if examples change.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - specify full-list link placement,
  marked inline count summary grammar, recommendation impact count summaries,
  and the `💪 Strength` Finding type display label (per requirements 1-8).

### To rename

None

### To delete

None
