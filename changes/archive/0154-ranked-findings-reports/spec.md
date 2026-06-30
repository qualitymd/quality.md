---
type: Functional Specification
title: Ranked Findings Reports
description: Requirements for ranked finding tables, a full findings index, and Advice ranking context in Requirement reports.
tags: [evaluation, reports, advice, findings]
timestamp: 2026-06-27T00:00:00Z
---

# Ranked Findings Reports

Companion to the [Ranked Findings Reports](../0154-ranked-findings-reports.md)
Change Case. This spec states _what_ the change must do. It defers to
[Evaluation report tree](../../specs/evaluation/reports/report-tree.md),
[qualitymd evaluation report](../../specs/cli/evaluation-report.md), and
[Evaluation Advice](../0150-evaluation-advice/spec.md) as normative context for
the existing report and Advice contracts.

The key words "MUST", "MUST NOT", and "MAY" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

`FindingRankingResult` already contains the Advice-phase ordering that tells a
reader which findings most shape next quality-management action. The generated
reports should make that ordering useful without duplicating the full finding
audit trail: `report.md` should stay top-10 and decision-ready, a full
`findings.md` index should carry the complete ordered list, and Requirement
finding details should show the ranking context when a reader lands on a
specific finding.

## Scope

Covered: generated Markdown report content, generated report links, stable
finding-detail anchors, enum display labels in ranked finding tables, and
Requirement finding-detail ranking context.

Deferred / non-goals: no new Advice data fields, no change to
`FindingRankingResult` validation, no finding filtering UI, no recommendation
coverage columns in findings tables, and no confidence column in Top Findings or
the full findings index.

## Requirements

### Ranked findings tables

R1. The run-level `report.md` Top Findings table **MUST** render the columns
`Rank`, `Finding`, `Area`, `Factors`, `Type`, and `Severity`, in that order.

> Rationale: the table should expose the ranked claim and its model context
> without duplicating confidence, basis, effect, criteria, or evidence from the
> detail report.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` and
> `specs/cli/evaluation-report.md` - update the Top Findings table contract.

R2. The Top Findings table **MUST** render at most the first 10
`FindingRankingResult.orderedFindings` entries by rank.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - preserve the
> capped run-report summary behavior.

R3. `qualitymd evaluation report build <run>` **MUST** generate `findings.md` at
the run root as the full ranked findings index.

> Rationale: users need a complete ranked findings surface parallel to the full
> recommendation index, while `report.md` remains a compact entrypoint.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md`,
> `specs/cli/evaluation-report.md`, and
> `specs/evaluation/records/data-layout.md` - add the full findings index
> artifact.

R4. `findings.md` **MUST** render the same columns as the Top Findings table and
**MUST** include every `FindingRankingResult.orderedFindings` entry ordered by
rank.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define the
> full findings index content.

R5. The run-level `report.md` **MUST** always link to `findings.md`, including
when the run has 10 or fewer ranked findings or no ranked findings.

> Rationale: the full findings index is the stable entrypoint for finding review
> and should not disappear based on row count.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` and
> `specs/cli/evaluation-report.md` - require the full findings index link.

### Links and display

R6. In Top Findings and `findings.md`, the `Finding` cell **MUST** use the
finding `statement` as the link text and **MUST** link to the exact finding
detail section in the generated Requirement report.

> Rationale: the statement is the report-table claim, while the linked detail
> owns condition, criteria, basis, effect, and evidence.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - require
> exact finding-detail links.

R7. Requirement report finding detail sections **MUST** provide stable anchors
derived from the finding ID, so ranked findings links do not depend on statement
wording.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define stable
> finding-detail anchors.

R8. In Top Findings and `findings.md`, the `Area` cell **MUST** render the
declaring Area title as link text pointing to that Area report.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define Area
> link rendering for ranked findings rows.

R9. In Top Findings and `findings.md`, the `Factors` cell **MUST** render the
Requirement's attached Factor titles as comma-separated links to their Factor
reports, or `—` when no Factor link can be rendered.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define Factor
> link rendering for ranked findings rows.

R10. In Top Findings and `findings.md`, `Type` and `Severity` **MUST** render
existing report display names, including their emoji, for known finding type and
severity enum values.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - require
> display-label rendering for ranked finding enum cells.

### Requirement finding details

R11. Each Requirement report finding detail section **MUST** render Advice
ranking context when the finding appears in `FindingRankingResult`.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - add
> per-finding Advice ranking context to Requirement reports.

R12. The Advice ranking context **MUST** include Advice rank as `<rank> / <total
ranked findings>`, tier, and ranking rationale.

> Rationale: the detail report should explain why a reader arrived from a ranked
> findings surface without turning the finding itself into Advice.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define the
> ranking context fields.

R13. If a Requirement finding detail has no matching ranking entry,
Requirement reports **MUST** render an explicit not-ranked state instead of
hiding the ranking context.

> Rationale: a missing ranking entry is meaningful report state and should not be
> mistaken for an omission in the renderer.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` - define the
> missing-ranking presentation.

R14. Generated reports **MUST NOT** introduce report-only findings, evidence,
ratings, recommendations, recommendation coverage, or ranking judgments while
rendering the ranked findings surfaces.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - preserve the
> deterministic projection boundary for the new surfaces.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - update the run report Top Findings
  contract, add the `findings.md` artifact, define stable finding-detail anchors,
  and add Requirement finding-detail Advice ranking context (R1-R14).
- `specs/cli/evaluation-report.md` - require `evaluation report build` to render
  `findings.md` and link it from `report.md` (R1, R3, R5).
- `specs/evaluation/records/data-layout.md` - include `findings.md` in the
  generated report tree paths (R3).

### To rename

None

### To delete

None
