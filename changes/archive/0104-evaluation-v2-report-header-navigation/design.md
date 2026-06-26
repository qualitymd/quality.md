---
type: Design Doc
title: Evaluation v2 report header navigation - design
description: Design for labeled Area and Factor trails plus compact Markdown report headers.
tags: [cli, evaluation, reports, navigation]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report header navigation - design

## Context

This design answers the [functional spec](spec.md) for simplifying Evaluation v2
Markdown report headers. The current renderer already centralizes report
navigation in `writeV2Breadcrumbs` and writes each report's identity table near
the top of `renderV2AreaReport`, `renderV2FactorReport`, and
`renderV2RequirementReport`. That keeps the change local to report projection:
loaded routine data and generated report paths do not need to change.

## Approach

Replace `writeV2Breadcrumbs` with narrower helpers:

1. `writeV2AreaTrail` writes `Area: ...` links from root through the supplied
   Area ID.
2. `writeV2FactorTrail` writes `Factor: ...` links from the local root Factor to
   the current Factor.
3. Area reports call only `writeV2AreaTrail` with the current Area.
4. Factor reports call `writeV2AreaTrail` with the declaring Area, then
   `writeV2FactorTrail` with the current Factor.
5. Requirement reports call only `writeV2AreaTrail` with the declaring Area.

The old parent-link writes are removed. Existing `areaReportPath`,
`factorReportPath`, and `reportLink` helpers still own relative-link generation,
so the navigation remains deterministic and tied to the current report path.

Each report renderer then replaces its generic `Field | Value` table with a
small identity line and one report-specific summary table:

- Area: `Path:` followed by `Overall | Local | Confidence | Data`.
- Factor: `Path:` followed by `Overall | Local | Status | Confidence | Data`.
- Requirement: `Name:` followed by
  `Rating | Assessment | Factors | Confidence | Data`.

The helper values already exist: `v2RatingLabel`, `v2RequirementRatingLabel`,
`v2AnalysisStatusPair`, `v2ConfidencePair`, `requirementFactorLinks`, and the
data-path helpers can be reused. The only label change with semantic weight is
rendering Factor `localAndDescendantAnalysis` as `Overall`; the durable spec
should define that mapping.

## Spec response

The Area trail satisfies the required owning Area navigation for all report
kinds. The Factor trail satisfies parent/ancestor Factor navigation for Factor
reports without adding a false Factor parent to Requirement reports.

Removing `Breadcrumb:` and `Parent*:` is safe because every destination they
provided remains reachable through `Area:` or `Factor:` trails. Requirement
Factor destinations remain reachable through the summary table's linked
`Factors` cell.

The compact summary tables preserve the data the report-tree spec already
requires while removing repeated subject labels from the highest-attention
header area. Existing detail sections, child tables, limits, findings, and data
links after the header are unchanged.

## Verification

Add or update focused report-rendering tests that assert:

- an Area report starts with `Area:` and has no `Breadcrumb:` or `Parent Area:`;
- a nested Area report's `Area:` trail links root, ancestors, and current Area;
- a Factor report has both `Area:` and `Factor:` trails and no `Parent Factor:`;
- a Requirement report has an `Area:` trail and linked `Factors` summary cell,
  but no `Factor:` trail; and
- generated machine outputs still contain the same report refs and raw routine
  data values as before.

## Alternatives

### Keep parent links for explicit upward navigation

Rejected. Parent links are explicit, but they duplicate destinations already
present in the trails and make the header longer. The labeled trails show both
the immediate parent and the wider context.

### Add a Factor breadcrumb to Requirement reports

Rejected. Requirements can be attached to multiple Factors. A Factor breadcrumb
would either pick one arbitrarily or require multiple trail lines, both of which
make ownership less clear than keeping linked Factors in the summary.

### Keep the generic `Field | Value` table

Rejected. It is flexible, but it repeats the subject kind directly below a title
that already names the subject. Report-specific summary columns make the state
easier to scan and leave stable paths/names as simple identity lines.

### Add an "In this report" table of contents

Deferred. It could help for long reports, but it is a separate navigation layer
with no current evidence that generated reports need it. The present change
focuses on the always-visible header.

## Trade-offs & risks

The compact summary table is less uniform across report kinds than `Field |
Value`. That is intentional: the report kinds expose different high-signal
state. Tests should cover each kind so future edits do not accidentally drop a
required header value.

The root Area and root Factor cases need attention because their trails contain
only one element. Those should still render as labeled linked trails, not
special-case into plain text, so navigation remains visually consistent.

## Open questions

None.
