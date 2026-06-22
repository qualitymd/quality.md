---
type: Design Doc
title: Area factor report breakdown - design
description: Strengthen the report model so concise and full reports share one Area-by-Factor breakdown.
tags: [evaluation, report, cli, model, ux]
timestamp: 2026-06-22T00:00:00Z
---

# Area factor report breakdown - design

Companion to [Area factor report breakdown](../0048-area-factor-report-breakdown.md)
and its [functional spec](spec.md).

## Context

The generated report already has one assembled in-memory `ReportDocument` that
feeds `report-summary.md`, `report.md`, and `report.json`. That model currently
splits compact Area data into `AreaSummary` while Factor ratings live only on the
detailed `Areas` entries. The summary renderer can therefore show Area ratings
but not the Factor ratings that explain them without reaching around its own
summary data.

## Approach

Keep the single assembled report model and make `AreaSummary` the canonical
compact Area breakdown. Add Factor rating results to `AreaRatingSummary` by
copying them from each `AreaEvaluationDetail` during report assembly. This keeps
the existing JSON `areaSummary` field as the machine-readable summary layer
while strengthening it instead of adding a second parallel summary object.

Name the two Area rating concepts in the wire model, not only in prose:

- **Area-only rating** (`areaRatingResult`) - the rating for requirements
  declared directly on that Area.
- **Area-with-descendants rating** (`areaWithDescendantsRatingResult`) - the
  roll-up for that Area including child Areas.

Rename the opaque `localRatingResult` / `aggregateRatingResult` fields to these
names on both `AreaRatingSummary` and `AreaEvaluationDetail`, and carry the typed
Area-only state (rated / not-assessed / structural) in a single `areaRatingState`
field. Retire the parallel `structural` bool and the derived "structural grouping
area" note, which re-encode state the typed `areaRatingState` already carries. The
rename is in scope because this change already expands the `areaSummary` schema and
regenerates every fixture — doing it now is one churn event, deferring it is two.

The rendering path stays:

```text
analysis records -> ReportDocument -> report-summary.md / report.md / report.json
```

`report-summary.md` renders an `Area Breakdown` table from `report.AreaSummary`.
The full `report.md` renames its summary-first Area overview to the same concept
and uses the same row helper, while its detailed Area sections continue to
render full rationales, analysis records, not-assessed requirements,
requirements, findings, and advice.

The summary table should replace the current "rating basis" prose column with a
Factor breakdown column. The preferred concise table labels are:

```md
| Path | Area | + Sub-Areas | Factors |
```

`Area` is the Area-only rating. `+ Sub-Areas` is the
Area-with-descendants rating. Area-with-descendants rationale stays in the
verdict and detailed Area sections, where it has room to remain readable and
traceable.

Factor display stays deterministic and label-aware:

- `report.json` preserves stable `areaPath`, `factorPath`, and rating `level`
  ids.
- Markdown resolves Area, Factor, and Rating Level titles from the run's
  `model.md` snapshot through the existing display-label resolver.
- Markdown path rendering uses path-aware label helpers, not only the terminal
  Area or Factor title, so nested rows and nested Factors remain distinguishable.
- Area path rendering uses an absolute display form: `/ (<root title>)` for the
  root Area and slash-prefixed descendant paths such as
  `/Services/Payments/Webhooks`.
- A typed structural/grouping Area-only state renders as `(area group)` in human
  Markdown, while `report.json` preserves the typed state for tools.
- Empty Factor sets render as an explicit human empty-state.

Strengthen analysis validation at the record boundary. `qualitymd evaluation
analysis set` can validate duplicate Factor paths and resolve each Factor path
against the run model before writing. Historical run inspection should keep its
current tolerance posture: malformed or invalid current-schema records surface as
non-reportable run gaps rather than being rewritten.

Cleanup should stay scoped. The expected cleanup is to retire summary-only
helpers that exist only to manufacture a rating-basis prose column once the
compact breakdown carries Factor results directly, to rename the Area rating
fields as above, and to remove the redundant `structural` bool and derived note.
It is not a license to rewrite the full report layout or rename stable JSON fields
unrelated to the Area rating model.

## Alternatives

**Render summary Factors by joining `report.Areas` in the Markdown renderer.**
Rejected. It would make the user-facing summary look better but preserve the
model gap and leave `report.json.areaSummary` weaker than the human summary.

**Add a new `areaBreakdown` JSON field beside `areaSummary`.** Rejected for now.
The existing field already names the compact report layer and has consumers in
tests and docs. Strengthening it avoids carrying two similar summary concepts.

**Rename the opaque Area rating fields to `areaRatingResult` /
`areaWithDescendantsRatingResult` and collapse the structural state into one
`areaRatingState`.** Adopted (see [Approach](#approach)). Originally deferred to
avoid a JSON schema churn, but this change already expands the `areaSummary`
schema and regenerates every fixture, so the rename rides along at no extra churn
and removes the duplicate structural encodings before they rot. The fence stays:
JSON fields unrelated to the Area rating model are not renamed here.

**Remove `AreaSummary` and make all summaries read from `Areas`.** Rejected for
this change. `AreaSummary` is useful as the compact machine shape; deleting it
would force lightweight consumers to parse the full detail array. The better
cleanup is to make `AreaSummary` complete enough for its job.

## Trade-offs & risks

Adding Factor results to `areaSummary` is a JSON schema expansion. It is
backward-compatible for normal JSON consumers but makes generated fixtures and
tests change.

Model-aware Factor path validation can expose existing hand-authored analysis
payload mistakes earlier. That is desirable for new writes, but historical reads
should continue reporting gaps instead of blocking status inspection outright.

Wide Markdown tables remain a risk for deeply nested paths or many Factors. The
first pass keeps the table because it is the most scannable shape for the
requested at-a-glance view; the full report still carries detailed prose below
it.

Keeping `areaSummary` and expanding it means the JSON name is less literal than
`areaBreakdown`, but it avoids schema churn and keeps the compact layer in one
place. The durable spec update should explain that `areaSummary` is the compact
Area breakdown.

## Open questions

None.
