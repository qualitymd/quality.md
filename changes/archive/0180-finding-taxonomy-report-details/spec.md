---
type: Functional Specification
title: Finding taxonomy and report details
description: Requirements for removing unknown Finding type and improving generated run report details.
tags: [evaluation, reports, findings, recommendations]
timestamp: 2026-06-29T00:00:00Z
---

# Finding taxonomy and report details

This change updates the current Evaluation Finding taxonomy and generated run
report opening. It does not preserve legacy Finding type aliases because
QUALITY.md is early alpha and the current contract should stay simpler than a
compatibility layer.

The key words "MUST", "MUST NOT", and "MAY" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

`unknown` currently reads like a separate concern type in reports, but missing
evidence already has explicit representation through not-assessed/not-rated
statuses, assessment `unknowns`, and rating `missingEvidence`. When an ambiguous
current state constrains a rating, the reader needs to see a current shortfall,
not a fifth concern class. Reports also need to distinguish observed gaps from
future risks and make the run-level key details useful without duplicating
confidence or exposing internal ranking terms.

## Scope

Covered:

- current Finding type values and report displays;
- guidance for ambiguous evidence, missing evidence, and current-state
  shortfalls;
- generated run report key-details and top recommendation tables;
- report-gallery examples generated from current structured Evaluation data.

Deferred:

- optional/type-specific severity;
- migration commands or compatibility aliases for `unknown`;
- historical archived runs outside the current generated report gallery.

## Requirements

1. Evaluation data validation **MUST** accept `strength`, `gap`, `risk`, and
   `note` as Finding types and **MUST NOT** accept `unknown`.

   > Rationale: Missing evidence belongs in explicit not-assessed/not-rated and
   > missing-evidence fields; ambiguous current shortfalls are gaps. Keeping
   > `unknown` as a Finding type makes reports present a concern that is neither
   > clearly a shortfall nor clearly an evidence limit.
   >
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
   > `specs/evaluation/routines/routine-contracts.md`, and
   > `specs/skills/quality-skill/evaluation.md` to define the four Finding
   > types and non-finding missing-evidence path.

2. Generated Markdown reports **MUST** display Finding types as `✅ Strength`,
   `🚩 Gap`, `⚠️ Risk`, and `ℹ️ Note`.

   > Rationale: `gap` and `risk` need distinct scan markers. A flag marks an
   > observed concern without implying the future/conditional posture that the
   > warning marker carries for `risk`.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` fixed enum
   > display requirements.

3. The `/quality` evaluation workflow **MUST** classify ambiguous current-state
   concerns that constrain a rating as `gap`, and **MUST** record missing or
   insufficient evidence that prevents rating through assessment/rating status,
   `unknowns`, or `missingEvidence` rather than a Finding type.

   > Rationale: The evaluator should produce one interpretation of the same
   > evidence limit across skill guidance, structured records, and generated
   > reports.
   >
   > Durable spec: modify `specs/skills/quality-skill/evaluation.md`; runtime
   > skill guidance changes live in `skills/quality/`.

4. The run report `Key Details` table **MUST** show one Confidence value aligned
   to the visible Overall Rating, a descriptive Scope value, and total Finding
   and Recommendation counts without the word `ranked`.

   > Rationale: The run report shows one Overall Rating, so a paired confidence
   > value is visually unbound. Counts in Key Details are inventory facts; ranking
   > is already communicated by table headings and rank/number columns.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

5. The run report **MUST** render a compact Finding Breakdown by Finding type
   near Key Details. The breakdown **MUST** show counts for each type present and
   **MUST** show severity breakdown only for `gap` and `risk`.

   > Rationale: Severity is meaningful for concerns that can constrain quality,
   > but it is not meaningful for strengths or neutral notes. The breakdown gives
   > users the marker-based count they asked for without implying every Finding
   > type has the same severity semantics.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

6. The run report Top Recommendations table **MUST** include recommendation
   confidence while preserving the recommendation number, linked title, area /
   factor trace, impact, and reason.

   > Rationale: Confidence is the most decision-useful extra signal already
   > present in `RecommendationRankingResult`; it adds detail without widening
   > the table with implementation-only data.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - requirements 2, 4, 5, and 6 update
  fixed enum displays and generated run report tables.
- `specs/evaluation/routines/routine-contracts.md` - requirements 1 and 3 update
  Finding type classification and missing-evidence guidance.
- `specs/evaluation/records/payload-kinds.md` - requirement 1 updates accepted
  Finding type values and guidance.
- `specs/skills/quality-skill/evaluation.md` - requirements 1 and 3 update
  evaluator behavior for Finding classification.
- `specs/cli/evaluation-report.md` - requirement 4 keeps the report command
  contract aligned with the run report shape.

### To rename

None

### To delete

None
