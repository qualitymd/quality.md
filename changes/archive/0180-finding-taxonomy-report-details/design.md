---
type: Design Doc
title: Finding taxonomy and report details design
description: Design for updating Finding type validation, report rendering, skill guidance, and generated examples.
tags: [evaluation, reports, findings, recommendations]
timestamp: 2026-06-29T00:00:00Z
---

# Finding taxonomy and report details design

## Context

This design answers the [Finding taxonomy and report details](spec.md)
functional spec. The change is a current-contract cleanup: remove `unknown` from
Finding types, use `gap` for ambiguous current shortfalls, keep missing evidence
on existing non-finding fields, and make the run report opening more useful.

## Approach

The Evaluation package already centralizes fixed enum values in
`internal/evaluation/display.go`; validation, schema generation, report legends,
and tests derive from those catalogs. The taxonomy change should therefore start
there: remove `FindingTypeUnknown`, remove it from `findingTypeValues`, and
change the `gap` marker to `🚩`.

Report rendering stays deterministic over existing structured data:

- `writeRunReportKeyDetails` renders only the scoped Area confidence because
  that is the confidence paired with the visible Overall Rating.
- a small `writeRunReportFindingBreakdown` helper derives counts from
  `artifacts.rankedFindings()` and renders immediately after the Key Details
  legend;
- severity detail is produced only for `gap` and `risk`;
- `writeTopRecommendationsTable` adds `Confidence` from the ranked
  recommendation entry.

The report-gallery synthetic `unknown` Finding becomes a `gap`: the example
records a conflicting current owner, which constrains recoverability today. It
is not an absence that prevents assessment.

Durable specs and skill guidance absorb the semantic rule: missing evidence that
prevents rating is represented as not assessed/not rated plus `unknowns` or
`missingEvidence`; ambiguous current-state concerns that can be evaluated are
`gap`.

## Spec response

- The enum catalog change makes validation reject `unknown` and propagates the
  new marker to report displays and local keys.
- The renderer changes keep the run report a projection over persisted Advice
  and analysis data; no report-only findings, recommendations, or ratings are
  introduced.
- The gallery data update demonstrates the new classification boundary with the
  only current `unknown` example.
- The skill guidance update keeps future agent-written evaluations aligned with
  the CLI contract.

## Alternatives

- **Rename `gap` and `unknown` to `issue`.** Rejected because `issue` is broader
  but less diagnostic. It would hide the useful distinction between current
  shortfalls and future risks.
- **Keep `unknown` but hide it from summary breakdowns.** Rejected because it
  leaves the ambiguous semantic split in the structured data and skill guidance.
- **Make severity optional by Finding type now.** Deferred because it is a
  larger data-contract change. The report can avoid showing severity where it is
  not meaningful without changing the required field yet.
- **Use `⬇️` for `gap`.** Rejected in favor of `🚩` because the flag better
  communicates a concern without implying a directional metric.

## Trade-offs and risks

Removing `unknown` is a breaking data-contract change for current-format runs
that still contain that Finding type. That is acceptable under the repo's
early-alpha compatibility policy, and the current checked-in v3 usage is limited
to one generated gallery example.

Adding a Finding Breakdown gives the run report more information near the top.
Keeping it compact and type-count focused avoids turning Key Details into a
second findings report.

## Open questions

None.
