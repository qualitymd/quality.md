---
type: Design Doc
title: Finding Summary Display Order Design
description: Design for updating the run report Finding Summary display, enum ordering, and generated examples.
tags: [evaluation, reports, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Finding Summary Display Order Design

## Context

This design answers the [Finding Summary Display Order](spec.md) functional
spec. The change is a generated report display cleanup: it keeps structured
Finding data unchanged while making the run report opening more explicit.

## Approach

The Evaluation package already centralizes Finding type display values in
`internal/evaluation/display.go`. Reorder `findingTypeValues.Values` there to
make concern-first order the shared report display order: `gap`, `risk`,
`strength`, `note`.

Keep the run report renderer deterministic over ranked Findings:

- `runReportFindingBreakdownRows` continues to derive counts and severity counts
  from `artifacts.rankedFindings()`;
- the row builder iterates the full Finding type catalog and no longer skips
  zero-count types;
- `writeRunReportFindingBreakdown` is renamed around the visible `Finding
  Summary` label and writes the `Severity` column;
- `findingBreakdownDetail` continues to list only observed severities for `gap`
  and `risk`, preserving `—` for non-concern types and empty severity sets.

Durable specs and report design guidance absorb the new visible contract.
Regenerating the report gallery updates the checked-in example report without
changing synthetic source data.

## Spec response

- The catalog reorder satisfies the shared Finding type display order and
  updates legends/local keys through existing enum-key rendering.
- The run report helper updates the section title, column title, and row
  inclusion behavior without altering Finding ranking data.
- Existing severity-count logic already lists only observed severity values in
  catalog order; preserving that path satisfies the sparse severity requirement.

## Alternatives

- **Use a report-specific order.** Rejected because it would make the run
  summary disagree with local keys and other report legends, adding a special
  case for little benefit.
- **Keep the `Finding Breakdown` heading.** Rejected because the fixed
  zero-count taxonomy rows make the table a summary, not only a breakdown of
  present observations.
- **Show every severity value with zero counts.** Rejected because it would turn
  a compact opening summary into a dense matrix and obscure the Finding type
  count signal.

## Trade-offs and risks

Changing the shared enum display order affects report legends and local keys
outside the run report opening. That is intentional for consistency, but tests
and gallery output need to update anywhere they assert the old order.

The run report now spends one extra row on absent `note` Findings in common
runs. The row is small, and making absence explicit is the purpose of the
change.

## Open questions

None.
