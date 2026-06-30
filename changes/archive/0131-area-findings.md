---
type: Change Case
title: Area findings in evaluation reports
description: Add Area Findings to Area analysis results and render them in Area and Factor reports before introducing recommendations.
status: Done
tags: [evaluation, reports, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Area findings in evaluation reports

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0131-area-findings/spec.md) — what the case must do.
- [Design doc](0131-area-findings/design.md) — how it's built, and why.

## Motivation

Evaluation reports currently expose Requirement Findings and roll-up rating
rationales, but they lack a durable analysis-level finding that says what the
Area analysis concluded from the underlying requirements, factor analyses, child
areas, limits, and patterns. That leaves Area and Factor reports without a
stable "findings" layer for readers to understand local quality concerns before
recommendations exist.

Recommendations should eventually address synthesized findings rather than
duplicate Area-level analysis. This case adds that synthesis layer first:
`AreaAnalysisResult.findings`, rendered in Area reports and projected into Factor
reports by explicit Factor relationships.

## Scope

Covered:

- Add **Area Findings** as analysis-phase findings on `AreaAnalysisResult`.
- Keep Area Findings scoped to the containing `AreaAnalysisResult.areaId`.
- Allow Area Findings to relate to zero or more Factors declared in the same
  Area.
- Render Area Findings in Area reports and matching Area Findings in Factor
  reports.
- Sort report findings deterministically by type, severity, local relationship
  where relevant, confidence, and source order.

Deferred:

- Recommendations and recommendation follow-up changes.
- "Top findings across evaluated scope" or any global cross-Area ranking.
- Impact, priority, ROI, Area importance, Factor weight, or other ranking scores.
- New report-only findings; reports remain deterministic projections over
  persisted data.

## Affected artifacts

Derived from the evaluation data/report path and the runtime skill workflow
surfaces that author Area analysis results.

**Code**

- [x] `internal/evaluation/data_contract.go` — extend `AreaAnalysisResult` with
      `findings`; add Area Finding and Factor relationship contracts.
- [x] `internal/evaluation/data.go` — update examples if needed by the data
      contract helpers.
- [x] `internal/evaluation/report_tree.go` — collect, sort, and render Area
      Findings in Area reports and matching Findings in Factor reports.
- [x] `internal/evaluation/display.go` — add or reuse display labels for finding
      type, severity, confidence, and relationship values as needed.
- [x] Tests under `internal/evaluation/` — validation, model binding, examples,
      schema output, Area report rendering, and Factor report filtering/sorting.

**Durable specs** (substance in the [functional spec](0131-area-findings/spec.md))

- [x] `SPECIFICATION.md` — recognize analysis-level Area Findings in evaluation
      and report semantics.
- [x] `specs/evaluation/routines/routine-contracts.md` — require `analyzeArea`
      to produce Area Findings without introducing new evidence.
- [x] `specs/evaluation/records/payload-kinds.md` — add the
      `AreaAnalysisResult.findings` contract.
- [x] `specs/evaluation/reports/report-tree.md` — require Findings sections in
      Area and Factor reports and define sorting.
- [x] `specs/cli/evaluation-data.md` — reflect the accepted data/schema contract
      for Area Findings where payload validation is specified.
- [x] `specs/skills/quality-skill/evaluation.md` — update the skill evaluation
      behavior for Area Finding synthesis and QC.
- [x] `specs/skills/quality-skill/reporting.md` — update report/closeout
      expectations without introducing recommendations.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` — update evaluate
      workflow phase responsibilities where Area analysis outputs are described.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` — update runtime evaluation guidance to author
      Area Findings.
- [x] `skills/quality/workflows/evaluate.md` — update Area analysis and report
      generation steps.
- [x] `skills/quality/resources/SPECIFICATION.md` — mirror public format
      semantics bundled with the runtime skill.
- [x] `skills/quality/log.md` and related runtime logs — append entries when the
      change lands.
- [x] `CHANGELOG.md` — release note when implemented.

No planned impact: `README.md`, install/scaffold files, or `QUALITY.md` model
frontmatter schema.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Implemented
and archived. `AreaAnalysisResult.findings` now validates traceable Area
Findings with same-Area Factor relationships, Area and Factor reports render
those findings deterministically, durable specs and runtime skill guidance are
updated, and `mise run check` passes.
