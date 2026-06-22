---
type: Functional Specification
title: Area factor report breakdown - functional spec
description: Expose a compact Area-by-Factor breakdown in generated evaluation reports from a first-class report model.
tags: [evaluation, report, cli, model, ux]
timestamp: 2026-06-22T00:00:00Z
---

# Area factor report breakdown - functional spec

Companion to
[Area factor report breakdown](../0048-area-factor-report-breakdown.md). This
spec states the report-output and report-model delta for showing evaluated
Areas broken down by Factors.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Background / Motivation

`report-summary.md` should let a reader understand the shape of an evaluation
without opening the full audit trail. The current summary has Area ratings, but
not the Factor results that explain them. That makes nested models hard to scan:
a weak `services/payments/webhooks` Area and a weak
`operations/incident-response` Area can look equally bad even when one is held
down by Observability and the other is simply not assessed.

The report model already records Factor rating results on Area details. This
change makes the compact Area breakdown a first-class projection of the
assembled report model, so the concise summary, full Markdown report, and JSON
report can share the same data instead of drifting through presentation-specific
logic.

## Scenario

An evaluator runs a full evaluation of a non-trivial quality model with nested
Areas and nested Factors, such as:

```text
Atlas Commerce
├── Services
│   ├── Payments
│   │   ├── Checkout API
│   │   └── Webhooks
│   └── Fulfillment
│       └── Inventory
│           └── Reservations
└── Operations
    └── Incident Response
```

The reader opens `report-summary.md` and can immediately see the aggregate
rating for each in-scope Area and the Factor ratings that explain it, including
paths such as `Services / Fulfillment / Inventory / Reservations` and Factors
such as `Operability / Observability`.

## Scope

Covered: generated report data and rendering for Area-by-Factor breakdowns,
nested Area paths, nested Factor paths, area-group Area-only rating states,
not-assessed states, title-based display labels, empty Factor sets,
machine-readable stable paths, and tests/examples proving those cases.

Deferred / non-goals: no changed rating semantics, no changed roll-up
semantics, no new evaluation record kind, no new interactive output format, no
automatic inference of missing Factor ratings, no compatibility support for
legacy `targetSummary` or `targets` report shapes, and no display-path escaping
for Area or Factor titles that themselves contain the path separator — the
canonical element arrays in `report.json` stay unambiguous, while the human
display path assumes separator-free titles.

## Requirements

The assembled report model **MUST** expose a compact Area breakdown that includes
each in-scope Area's stable `areaPath`, Area-only rating state,
Area-with-descendants rating result, covered requirement count,
structural/not-assessed typed state when relevant, and Factor rating results.

> Rationale: The concise report should project structured report data, not
> reconstruct missing structure from detailed Markdown sections. This keeps
> `report-summary.md`, `report.md`, and `report.json` aligned.

The report model **MUST** name the two Area ratings with self-documenting wire
fields: `areaRatingResult` for the Area-only rating (the requirements declared
directly on that Area) and `areaWithDescendantsRatingResult` for the roll-up
including child Areas. It **MUST** carry the typed Area-only state — rated,
not-assessed, or structural/area group — in a single `areaRatingState` field, and
**MUST NOT** keep a separate boolean or derived note that re-encodes that
structural state.

> Rationale: "local" and "aggregate" are precise to implementers but opaque to
> readers, and a parallel `structural` bool plus a derived "structural grouping
> area" note give the Area-only state two more homes that drift out of sync. One
> typed state field, and ratings named for what they cover, keep the durable model
> honest. This change already expands the `areaSummary` JSON schema and
> regenerates every fixture, so the rename rides along as one churn event instead
> of forcing a second one later.

The compact Area breakdown **MUST** preserve stable machine identifiers in
`report.json` as ordered element arrays — `areaPath` entries and `factorPath`
entries — plus rating `level` ids. These ordered element arrays are the
canonical, unambiguous identifiers; renderers **MUST** build display paths from
the elements, never by parsing a joined string. Human Markdown **MUST** render
Model, Area, Factor, and Rating Level titles from the run's `model.md` snapshot,
falling back to stable identifiers only when a title is unavailable.

> Rationale: the joined display path (`/Services/Payments/Webhooks`) is a derived,
> best-effort label, not an identifier; a title containing the separator would
> make it ambiguous. Identity lives in the element array, so nothing downstream
> reconstructs a path by splitting the display string.

Human Markdown **MUST** render nested Area and Factor paths as path labels, not
only terminal-node labels. Area paths **MUST** render as absolute display paths:
the root Area renders as `/ (<root title>)`, and descendants render with a
leading `/`, such as `/Services/Payments/Webhooks`. For example, a row should
distinguish `/Services/Payments/Webhooks` from a sibling
`/Operations/Webhooks`, and a Factor should distinguish
`Operability / Observability` from another `Observability` sub-factor.

`report-summary.md` **MUST** include an at-a-glance Area breakdown showing each
in-scope Area's display path, Area-with-descendants rating, and Factor
breakdown.

Human Markdown **SHOULD** label the display-path column as `Path` and the two
rating columns as `Area` and `+ Sub-Areas` when both ratings are shown.

> Rationale: `Area` communicates the rating for this Area's own requirements;
> `+ Sub-Areas` communicates the roll-up including descendants without exposing
> the implementation word "aggregate".

The Area breakdown **MUST** support nested absolute Area display paths at depth
three or greater, such as `/Services/Fulfillment/Inventory/Reservations`.

The Factor breakdown **MUST** support nested Factor display paths, such as
`Operability / Observability`.

The Factor breakdown **MUST** distinguish rated results from `not assessed`
results. It **MUST NOT** render `not assessed` as a rating level.

Areas with child Areas but no direct requirements **MUST** render their
Area-only rating as `(area group)` in human Markdown, while preserving their
Area-with-descendants rating and any recorded Factor ratings. Machine-readable
report data **MUST** still preserve the typed structural/grouping state.

Areas with no recorded Factor ratings **MUST** render an explicit empty-state in
human Markdown, rather than implying that Factor ratings were omitted by
accident.

`report.md` **MUST** reuse the same compact Area breakdown data or rendering
helper used by `report-summary.md` for its summary-first Area overview. The full
report **MUST** still preserve detailed Area, Factor, Requirement, Finding, and
Advice sections.

The concise Area breakdown **MUST NOT** replace detailed rating rationale in the
full report. Area-with-descendants and Factor rationales remain part of the full
Area detail and machine-readable rating results.

`report.json` **MUST** expose the compact Area breakdown as a first-class
machine-readable shape. It **MUST NOT** require consumers to join separate
detailed Area arrays only to recover the summary-level Factor breakdown.

The report renderer **MUST NOT** reread evaluated source, recompute ratings,
invent Factor ratings, infer roll-ups, or choose new recommendations while
building the breakdown. It may sort, format, label, and summarize recorded
report data deterministically.

The implementation **SHOULD** remove or simplify stale report-summary model
fields when the new compact breakdown makes them redundant.

> Rationale: The previous summary shape left a split between `AreaSummary` and
> detailed Area data. This change should not preserve duplicate structures that
> make future report output harder to reason about.

Analysis record validation **SHOULD** reject duplicate Factor rating paths within
one analysis record.

Analysis record validation **SHOULD** reject Factor rating paths that cannot be
resolved against the Area's declared or inherited Factor vocabulary.

> Rationale: A stronger report display is only trustworthy if the underlying
> Factor paths are valid enough for title resolution and tool consumption.

Tests **MUST** cover at least:

- a structural root Area with child Areas;
- a nested Area path at depth three or greater;
- a nested Factor path;
- mixed rated outcomes across Areas and Factors;
- `not-assessed` Area and Factor outcomes;
- an Area with no Factor ratings;
- title rendering with identifier fallback;
- stable `report.json` identifiers;
- idempotent report generation.

## Example shape

```md
## Area Breakdown

| Path                                         | Area         | + Sub-Areas  | Factors                                                                                                 |
| -------------------------------------------- | ------------ | ------------ | ------------------------------------------------------------------------------------------------------- |
| / (Atlas Commerce)                           | (area group) | 🟡 Minimum   | Trust: 🔵 Target; Operability: 🟡 Minimum; Changeability: 🔵 Target                                     |
| /Services/Payments/Webhooks                  | 🟡 Minimum   | 🟡 Minimum   | Secret handling: 🔵 Target; Availability: 🟡 Minimum; Observability: 🟡 Minimum; Testability: 🔵 Target |
| /Services/Fulfillment/Inventory/Reservations | 🟡 Minimum   | 🟡 Minimum   | Data isolation: 🔵 Target; Availability: 🟡 Minimum; Maintainability: 🔵 Target                         |
| /Operations/Incident Response                | not assessed | not assessed | Availability: not assessed; Observability: not assessed                                                 |
```

## Durable spec changes

### To add

- `specs/reports/report-summary-md.md` - add a 1:1 durable spec for the concise
  human triage artifact and its Area-by-Factor breakdown requirements.
- `specs/reports/report-md.md` - add a 1:1 durable spec for the complete human
  Evaluation Report and its shared Area breakdown.
- `specs/reports/report-json.md` - add a 1:1 durable spec for the
  machine-readable Evaluation Report and compact Area summary layer.

### To modify

- `specs/evaluation-records.md` - update the `report.json`,
  `report-summary.md`, and `report.md` contracts for the compact
  Area-by-Factor breakdown, including the renamed `areaRatingResult` /
  `areaWithDescendantsRatingResult` / `areaRatingState` fields, element-array
  identifiers, and human title rendering; delegate artifact-specific shape to the
  new `specs/reports/` contracts (per the report-model and rendering requirements
  above).
- `specs/cli/evaluation-report.md` - link the report command contract to the
  shared report-output contract and new artifact-specific specs (per the shared
  report-model requirement above).
- `specs/skills/quality-skill/quality-skill.md` - align the skill reporting
  contract with the compact Area-by-Factor breakdown and typed report states
  (per the report consumption requirements above).
- `SPECIFICATION.md` - update only if minimum Evaluation Report semantics should
  require Area-by-Factor summary visibility in conforming report renderers (per
  the Area breakdown requirement above).

### To rename

- `specs/skills/quality-skill/guides/authoring.md` →
  `authoring-md.md`,
  `specs/skills/quality-skill/guides/getting-started.md` →
  `getting-started-md.md`, and
  `specs/skills/quality-skill/guides/top-10-quality-md-checks.md` →
  `top-10-quality-md-checks-md.md` - apply the artifact-spec filename convention
  the new report specs introduce. A rename is a delete-plus-add, so the old paths
  are removed and inbound links updated; no runtime guide artifact contract
  changes are required (per the artifact-spec convention referenced above).

### To delete

None
