---
type: Design Doc
title: Area / Factor Breakdown Reports - design doc
description: Design for rendering shared Area / Factor breakdown tables in generated Evaluation reports.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Area / Factor Breakdown Reports - design doc

Design for
[Area / Factor Breakdown Reports](../0161-area-factor-breakdown.md) and its
[functional spec](spec.md).

## Context

The run report currently renders `Subject Reports` from
`[]evaluationRenderedReport`, which is an output-artifact list. That is the
wrong data source for a model overview because it includes findings indexes,
recommendation indexes, and recommendation detail reports alongside Areas,
Factors, and Requirements.

Area reports separately render local `Factors` and `Child Areas` tables from
`evaluationArtifacts`. Those tables are closer to the model shape, but they are
split by kind, include path columns, and do not show findings or
recommendations.

## Approach

Add one shared body helper that renders an `Area / Factor Breakdown` section
from Evaluation artifacts instead of report-output artifacts. The helper takes a
root Area ID, the current report path, and a flag for emphasizing the root row.
It writes rows in deterministic pre-order:

1. the root Area;
2. that Area's local root Factors, including nested sub-factors;
3. each direct Child Area, recursively using the same pattern.

Each Area row uses the Area report link and the Area Analysis Result's
`localAndDescendantAnalysis` and `localAnalysis` values. Each Factor row uses
the Factor report link, the inline Factor marker, and the Factor Analysis
Result's `localAndDescendantAnalysis` and `localAnalysis` values. The helper
does not emit Requirements.

Finding counts are derived from `FindingRankingResult.orderedFindings`. Area
rows count ranked findings whose Requirement belongs to that Area or a
descendant Area. Factor rows count ranked findings whose Requirement is attached
to the Factor or one of its sub-factors.

Recommendation counts are derived from ranked Recommendation outputs. Area rows
count ranked recommendations whose trace context resolves to that Area or a
descendant Area. Factor rows count ranked recommendations whose trace context
resolves to the Factor or one of its sub-factors. A recommendation counts at
most once per table row even when multiple trace refs match.

Use the helper in:

- `renderEvaluationRunReport`, replacing `Subject Reports`;
- `renderEvaluationAreaReport`, replacing the separate `Factors` and
  `Child Areas` sections.

Keep `evaluationOutputResult`, `reportOutputs`, report paths, report
frontmatter source-data helpers, Requirement sections, Factor report
`Sub-Factors`, and recommendation reports unchanged unless tests expose a direct
contract mismatch.

## Spec response

- The run report and Area reports both render the same section and table
  columns because the shared helper owns the heading and header row.
- The old run-level manifest table disappears because the run renderer no
  longer calls `writeEvaluationRunReportsTable`.
- The old split Area structure tables disappear because the Area renderer
  replaces both section writes with the shared helper.
- The `Area / Factor` cell owns navigation, so no `Report` column is needed.
- Counts remain deterministic because they are computed from sorted report data
  already loaded into `evaluationArtifacts` and ranked Advice payload order.
- Machine report manifests remain unchanged because `evaluationOutputResult`
  continues to consume `[]evaluationRenderedReport`.

## Alternatives

- **Keep `Subject Reports` and add the new table above it.** Rejected. That
  preserves the file-manifest framing and makes `report.md` longer without
  adding a distinct reader task.
- **Make the table a shallow one-level Area report table.** Rejected for the run
  report because it would still force readers to open child pages to see the
  evaluated Area / Factor shape. A recursive table is more useful and uses the
  same helper for run and Area reports.
- **Add Requirements to the breakdown.** Deferred. Requirements remain linked
  from Area, Factor, and Finding detail surfaces. Adding them now would widen
  and lengthen the overview before the Area / Factor shape is proven.
- **Add confidence or status columns.** Deferred. Those details are still in
  report headers and detail pages. The first pass keeps the overview narrow.
- **Drive the table from `reportOutputs`.** Rejected. `reportOutputs` is the
  generated artifact index, not the model/evaluation structure.

## Trade-offs & risks

- Recursive Area reports can be longer than the old shallow `Factors` and
  `Child Areas` sections. The table is still narrow, and it replaces two
  sections rather than adding a third.
- Recommendation trace refs can point at Areas, Factors, or Requirements. The
  count helpers need to resolve those refs consistently and avoid double-counts
  per row.
- Finding and recommendation counts are intentionally only counts. They do not
  communicate severity, type, impact, or confidence mix in this change.

## Open questions

None.
