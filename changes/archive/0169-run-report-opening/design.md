---
type: Design Doc
title: Run Report Opening - design doc
description: Design for rendering generated report.md with Summary, Key Details, Contents, and report metadata separation.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Run Report Opening - design doc

Design for [Run Report Opening](../0169-run-report-opening.md) and its
[functional spec](spec.md).

## Context

The shared `renderReportHeader` helper currently writes frontmatter, H1, the run
context line, report navigation, context trails, a summary table, and jump links
for every report type that opts into those fields. The run report then appends a
bare `Summary:` label, Top Findings, Top Recommendations, Scope, Coverage, and
Limits sections.

The new behavior is run-report specific. Detail reports should keep the shared
header behavior because a deep-linked Area, Factor, Requirement, or
recommendation report still needs immediate run context.

## Approach

Keep `renderReportHeader` unchanged for non-run reports. Add a run-report
specific header renderer or options path that writes:

- frontmatter with run metadata fields;
- H1;
- report navigation;
- Area context line;
- `## Summary`;
- `## Key Details`;
- `## Contents`.

The run frontmatter metadata comes from already-loaded objects:

- `run`: derived from `RunManifest.number` and rendered as the run folder label
  used by humans, such as `0001-full-eval`;
- `runId`: `RunManifest.id`;
- `created`: `RunManifest.createdAt`;
- `scope`: `requestedScopeLabel(RunManifest.requestedScope)`;
- `subject`: the planned Area model reference, such as `area:root`.

The key-details table reuses existing display helpers for rating, confidence,
and requested scope. Ranked finding and recommendation counts come from the
already-collected Advice ranking helpers.

`## Summary` uses the existing scoped Area summary. When the ranked
recommendation list is non-empty, render a deterministic one-sentence next
action from the first ranked recommendation title. If no recommendation exists,
omit the sentence rather than adding an empty-state paragraph.

`## Contents` should be assembled from the sections the run report actually
writes. Since this case removes visible `Limits & Incomplete Inputs`, contents
must not include that anchor. Top Findings and Top Recommendations stay visible
because Advice is required for reportability.

Move run traceability into a `## Report Details` table near the bottom, before
the existing legend and source-data sections. The table uses existing manifest
and plan values and does not introduce new Evaluation judgment.

Remove the run-report visible Limits section and the factor-filter synthetic
limit row. Do not remove Area Analysis limit data from source payloads or
source-data links.

## Spec response

- The run report gets metadata-rich, non-judgmental frontmatter without
  changing other report frontmatter.
- The old visible run line disappears only from `report.md`.
- The run opening order is controlled by the run-specific renderer.
- Summary prose, key-details facts, and contents links are separate sections.
- Visible limits are omitted from `report.md`, while structured data remains
  unchanged.
- Report details preserve human-visible traceability near the bottom.

## Alternatives

- **Change `renderReportHeader` globally.** Rejected because detail reports use
  the run line as useful deep-link context and should not inherit the run report
  opening pattern.
- **Keep run metadata visible at the top and only add frontmatter.** Rejected
  because it preserves the duplication that motivated the change.
- **Keep `Jump to:` and rename only the summary table.** Rejected because
  `report.md` is the primary artifact and deserves first-class contents
  navigation.
- **Move limits into Key Details as a count.** Rejected because the current
  scope removes visible limits/incomplete-input content for now.

## Trade-offs & risks

Removing the visible limits section can hide important evaluation caveats from a
casual reader. This is an intentional scope choice for now, but the structured
payloads and source-data links still preserve the underlying data for agents and
future report designs.

Adding run-report-specific rendering duplicates a little header logic, but it
keeps the shared detail-report header stable and limits the blast radius.

## Open questions

None.
