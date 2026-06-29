---
type: Design Doc
title: Run Report Simplification - design doc
description: Design for simplified generated run-level Evaluation reports.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Run Report Simplification - design doc

Design for
[Run Report Simplification](../0171-run-report-simplification.md) and its
[functional spec](spec.md).

## Context

The report renderer already separates the run report from detail reports:
`renderEvaluationRunReport` builds `report.md`, while `renderReportHeader`
handles Area, Factor, Requirement, index, and recommendation pages. That split
lets the primary run report lose opening chrome without weakening direct-entry
orientation on granular pages.

Source-data rendering is also centralized through `writeSourceDataSection` and
per-report source-data helper functions. Renaming the section and narrowing the
run report source list can therefore happen without changing structured
Evaluation data or the generated report tree.

## Approach

Keep the run report as a custom renderer rather than routing it through the
shared detail-report header. `renderRunReportHeader` will write only frontmatter
and the H1. The run report body will then render:

1. `## Summary`;
2. `## Key Details`;
3. `## Contents`;
4. `## Model Evaluation`;
5. `## Top Findings`;
6. `## Top Recommendations`;
7. `## Legend`;
8. `## Primary Source Data`.

Move the existing Area/Factor table by parameterizing the breakdown-section
writer with a heading. Run reports pass `Model Evaluation`; Area reports keep
`Area / Factor Breakdown`. This avoids duplicating table rendering and keeps the
renaming scoped to `report.md`.

Delete the run-report recommended-next-action append. The ranked recommendation
table remains the next-action surface.

Remove the `Scope`, `Coverage`, and `Report Details` body sections from
`report.md`. Scope remains in frontmatter and `Key Details`. Run identity,
creation time, and subject remain in frontmatter and primary source data. To
preserve scoped-run coverage safety, add a compact paragraph below the Model
Evaluation table when the root Area report was not generated, linking the root
Area report when it exists and otherwise naming that the root Area was outside
the run.

Rename `writeSourceDataSection` output to `## Primary Source Data`. Keep
path-as-label links and report-relative targets. Narrow the run report's source
data helper to primary run/report inputs: run manifest, scoped Area analysis,
finding ranking, and recommendation ranking. Detail reports keep their own
report-local primary inputs, trimmed only when they were listing descendant
payloads solely for a high-level table.

Update tests to assert both absence and order: run report H1 flows to Summary,
Key Details, Contents, Model Evaluation, Top Findings, and Top Recommendations;
removed sections and lines are absent; `Primary Source Data` appears; detail
reports still carry their trails and Area reports keep `Area / Factor Breakdown`.

Regenerate the report gallery after renderer changes so the checked-in example
matches the report contract.

## Spec response

- Custom run header rendering satisfies the removed post-H1 navigation
  requirements without touching detail reports.
- Parameterized breakdown headings allow `report.md` to say
  `Model Evaluation` while Area reports keep `Area / Factor Breakdown`.
- Removing the recommendation sentence leaves Top Recommendations as the single
  ranked recommendation surface.
- Keeping scope in Key Details and frontmatter preserves decision-relevant scope
  while deleting the long Scope section.
- A compact scoped-run note preserves the old Coverage safety signal without
  keeping the Coverage section.
- Centralized primary-source rendering gives every report the new heading and
  keeps source-data links report-local.

## Alternatives

- **Move removed navigation into a smaller top nav line.** Rejected because the
  run report already has Contents and direct links in the ranked tables; another
  opening nav line still delays Summary.
- **Keep `Area / Factor Breakdown` but move it upward.** Rejected because the
  old label is model-internal and table-shaped; `Model Evaluation` better names
  the reader-facing purpose of the section.
- **Rename every breakdown section to `Model Evaluation`.** Rejected because
  Area detail reports are subject-specific pages, and the old label remains
  precise there.
- **Keep `Source Data` and only trim lists.** Rejected because the old heading
  implies exhaustiveness. `Primary Source Data` better matches the new
  report-local contract.
- **Render no source-data links on run reports.** Rejected because the run
  report still needs a short, stable bridge to machine-readable identity, scope,
  and advice ranking payloads.

## Trade-offs & risks

The simplified run report depends more heavily on frontmatter for run ID,
created time, and subject reference. That is acceptable for the primary human
path because those values are not judgment-critical, but tools that scrape
visible `Report Details` will need to read frontmatter or structured data.

The primary-source principle is intentionally non-exhaustive. It improves human
readability, but readers who want full traceability must follow links into
granular reports and payloads rather than expecting `report.md` to list every
input.

## Open questions

None.
