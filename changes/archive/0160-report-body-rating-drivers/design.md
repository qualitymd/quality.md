---
type: Design Doc
title: Report Body Rating Drivers - design doc
description: Design for removing standalone rating-driver sections from generated Markdown reports.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Report Body Rating Drivers - design doc

Design for
[Report Body Rating Drivers](../0160-report-body-rating-drivers.md) and its
[functional spec](spec.md).

## Context

The report renderer currently writes `## Rating Drivers` sections directly in
three report body paths: run report, Area report, and Factor report. A shared
helper renders the table as `Driver | Effect | Inputs`, with `Inputs` serialized
from `ratingDrivers[].inputRefs` as compact JSON.

Those same reports now carry source-data frontmatter pointing at the structured
payloads used to render the artifact. The underlying Area and Factor Analysis
Result payloads still contain `ratingDrivers`, so the machine trace remains
available without projecting the raw driver table into the human report body.

## Approach

Remove the body section writes from the run, Area, and Factor renderers:

- delete the `## Rating Drivers` heading and table call after each summary;
- let the next useful section follow the summary directly (`Top Findings`,
  `Factors`, or `Requirements`);
- delete the shared driver-table rendering helper if no remaining report body
  uses it.

Keep all structured Evaluation data behavior unchanged. The data setter,
effective-run validation, schemas, and routine payload contracts continue to
require and validate `ratingDrivers` for rated outputs. Report frontmatter keeps
listing the relevant analysis payloads through existing source-data helpers, so
readers and agents can open the JSON when they need the trace.

Update tests by converting positive `## Rating Drivers` assertions into absence
assertions, while keeping neighboring-section assertions to prove the report did
not lose the human-facing content around the removed table. Regenerate the
report gallery through the existing task so checked-in examples match the new
body contract.

## Spec response

- Run, Area, and Factor report bodies stop rendering standalone rating-driver
  sections because the only calls that emit those sections are removed.
- Structured rating-driver preservation is unaffected because no data contracts,
  validation, payload schema, or report source-data helper is weakened.
- Human report content remains because the summary and adjacent sections stay in
  place and tests continue to assert them.
- Durable report specs and the `/quality` reporting spec move the explanation
  boundary from visible driver tables to summary, ratings, findings,
  recommendations, limits, incomplete inputs, and source payloads.

## Alternatives

- **Render prettier driver links instead of removing the section.** Deferred.
  Clickable driver links would still create another explanation surface that
  competes with findings and recommendations, and it is not needed to satisfy
  the current report-body simplification.
- **Keep the section only on detail reports.** Rejected. Area and Factor detail
  reports are still human review surfaces, and their source payloads already
  provide the trace for agents or readers who need it.
- **Remove `ratingDrivers` from structured data.** Rejected. That would weaken
  the rating audit trail and conflict with the Evaluation data validation
  contract.

## Trade-offs & risks

- Readers lose a visible raw trace table in Markdown. The source-data
  frontmatter and JSON payloads preserve that trace, but the report body now
  assumes summaries, findings, recommendations, and limits carry the useful
  review narrative.
- Poorly written summaries may become more noticeable because the driver table
  no longer compensates for them. That is an evaluation-quality issue, not a
  report-rendering reason to expose raw payload references.

## Open questions

None.
