---
type: Change Case
title: Finding Taxonomy and Report Details
description: Remove unknown as a Finding type, distinguish gap/risk markers, and make run report key details more useful.
status: Done
tags: [evaluation, reports, findings, recommendations]
timestamp: 2026-06-29T00:00:00Z
---

# Finding Taxonomy and Report Details

Generated Evaluation reports currently blur several reader-facing signals:
`gap` and `risk` share the same marker, `unknown` behaves like a rating concern
even though missing evidence already has explicit non-finding fields, and the
run report key-details table shows two confidence values, terse scope text, and
counts described as "ranked".

- [Functional spec](0180-finding-taxonomy-report-details/spec.md) - what the
  Finding taxonomy and report opening must change.
- [Design doc](0180-finding-taxonomy-report-details/design.md) - how the
  renderer, data contract, skill guidance, and gallery absorb the cleanup.

## Motivation

The generated run report should help a reader identify the current quality
concerns and next moves without decoding implementation terms. `Unknown` is
doing two jobs: missing or ambiguous evidence and an observed ambiguous current
state. Missing evidence is better represented by not-assessed/not-rated results,
`unknowns`, and `missingEvidence`; an ambiguous current state that constrains a
rating is a current shortfall and should be a `gap`. Separating this also lets
`gap` and `risk` carry distinct report markers and lets the key-details area
summarize findings without implying severity applies to strengths or notes.

## Scope

Covered:

- Finding type taxonomy for current Evaluation data and skill-authored
  findings;
- generated report displays, legends, and report-gallery examples for Finding
  types;
- run report `Key Details`, scope text, finding/recommendation counts, and
  compact Finding breakdown;
- run report Top Recommendations detail columns;
- durable Evaluation, CLI, and skill specs plus report design guidance;
- runtime `/quality` skill guidance for classifying evidence gaps and findings;
- tests and checked-in generated report examples.

Deferred:

- historical archived evaluation records that predate the current structured
  Evaluation data contract;
- compatibility aliases or fallback readers for the removed `unknown` Finding
  type;
- a full redesign of `severity` as an optional or type-specific field;
- changes to the `unknowns` and `missingEvidence` payload fields beyond clarifying
  when to use them;
- recommendation detail report shape beyond any field needed by the run report.

## Affected artifacts

- Code:
  - `internal/evaluation/display.go` - remove `unknown` from Finding type
    catalog and change the `gap` marker.
  - `internal/evaluation/data_contract.go` - update Finding type validation and
    generated schema behavior through the catalog.
  - `internal/evaluation/report_tree.go` - render revised key details,
    Finding breakdown, and Top Recommendations columns.
  - `internal/evaluation/evaluation_test.go` - update report, enum, and schema
    expectations.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - update fixed enum displays, run
    report key details, Finding breakdown, and Top Recommendations columns.
  - `specs/evaluation/routines/routine-contracts.md` - update Finding type
    classification and candidate-action wording.
  - `specs/evaluation/records/payload-kinds.md` - update Finding type contract
    and missing-evidence guidance.
  - `specs/skills/quality-skill/evaluation.md` - update skill evaluation
    behavior for Finding classification and non-finding unknowns.
  - `specs/cli/evaluation-report.md` - keep the CLI report command contract in
    sync with the run report shape.
- Runtime skill:
  - `skills/quality/SKILL.md` - update evaluator guidance for Finding types,
    missing evidence, and gap/risk candidate actions.
  - `skills/quality/workflows/evaluate.md` - update authoring workflow guidance
    for gap/risk findings and non-finding missing evidence.
- Durable docs:
  - `docs/guides/reporting-design.md` - update generated report design guidance
    and examples.
- Generated examples:
  - `examples/report-gallery/` - regenerate generated Markdown and update the
    synthetic `unknown` Finding example.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md`, `specs/evaluation/log.md`, `skills/quality/log.md`,
    `skills/quality/workflows/log.md`, and `docs/log.md` where the affected
    bundles record updates.

## Status

`Done`. Implemented and archived. Evaluation Finding types now use `strength`,
`gap`, `risk`, and `note`; generated reports render `🚩 Gap`, use a single
Overall Rating confidence in run Key Details, add a compact Finding Breakdown,
and include confidence in run Top Recommendations.
