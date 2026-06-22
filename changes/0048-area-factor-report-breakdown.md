---
type: Change Case
title: Area factor report breakdown
description: Make generated evaluation reports expose an at-a-glance Area breakdown by Factor without weakening the underlying report model.
status: In-Review
tags: [evaluation, report, cli, model, ux]
timestamp: 2026-06-22T00:00:00Z
---

# Area factor report breakdown

A **Change Case** capturing the *why* and *status* for improving generated
evaluation reports so readers can see, at a glance, which Areas were evaluated
and which Factors drove their ratings. The detail lives in its
[functional spec](0048-area-factor-report-breakdown/spec.md).

> **In-Review.** Implementation complete and verified in `internal/evaluation/`
> (`go test ./...` green); the three `0001-quality-eval` golden report fixtures
> are regenerated. **Durable specs:** the new `specs/reports/` artifact specs and
> the guide-spec renames already landed — durable specs track the current source
> of truth and are not lifecycle-gated.

## Motivation

The current concise report summary exposes Area ratings, but it does not show
the Factor breakdown that explains those ratings. A reader can see that an Area
landed at `minimum` or `not assessed`, but must open the full report and scan
detailed Area sections to answer the practical follow-up questions: which
quality Factors drove the result, where are the weak spots, and which nested
Area paths need attention.

The fix should strengthen the report model, not paper over the gap with a
summary-only Markdown table. The same Area-by-Factor breakdown should be
available to `report-summary.md`, reusable in `report.md`, and exposed in
`report.json` for tools.

## Scope

Covered: the assembled report model, generated `report-summary.md`,
summary-first `report.md`, `report.json`, tests and examples for nested Areas
and Factors, durable report contracts, and cleanup of stale summary/report model
cruft when the new model shape makes it redundant.

Deferred / non-goals: no changed rating semantics, no changed roll-up
semantics, no new evaluation judgment, no interactive report viewer, no new
rating aggregation formula, and no compatibility layer for historical `targets`
report shapes.

## Affected artifacts

### Code

- `internal/evaluation/report.go` - assemble and render the Area-by-Factor
  breakdown from the report model for `report-summary.md`, `report.md`, and
  `report.json`; rename the Area rating fields to `areaRatingResult` /
  `areaWithDescendantsRatingResult` / `areaRatingState` and drop the redundant
  `structural` bool and derived note.
- `internal/evaluation/types.go` - reshape `LocalRatingState` toward the
  `areaRatingState` vocabulary as the single carrier of rated / not-assessed /
  structural state; update report and analysis-facing types for the clearer
  breakdown shape.
- `internal/evaluation/write.go` and `internal/evaluation/load.go` - validate
  Factor rating paths more strongly if the design requires model-aware
  validation.
- `internal/evaluation/evaluation_test.go` - add regression coverage for nested
  Area paths, nested Factor paths, mixed ratings, structural Areas,
  not-assessed states, empty Factor breakdowns, title rendering, and idempotent
  generation.

### Durable specs

- `specs/reports/report-summary-md.md` - add the durable artifact-specific
  contract for the concise human report summary and update it for the
  Area-by-Factor breakdown.
- `specs/reports/report-md.md` - add the durable artifact-specific contract for
  the full human evaluation report and its shared Area breakdown.
- `specs/reports/report-json.md` - add the durable artifact-specific contract
  for the machine-readable evaluation report and compact Area summary layer.
- `specs/reports/index.md`, `specs/index.md`, and `specs/log.md` - register and
  log the new report artifact specs.
- `specs/skills/quality-skill/guides/authoring-md.md`,
  `specs/skills/quality-skill/guides/getting-started-md.md`, and
  `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - rename
  existing 1:1 runtime guide contract specs to the artifact-spec filename
  convention, with links updated.
- `specs/evaluation-records.md` - update the `report.json`,
  `report-summary.md`, and `report.md` report contracts for the Area-by-Factor
  breakdown, and delegate artifact-specific contracts to the new report specs.
- `specs/cli/evaluation-report.md` - link the command-level report contract to
  the shared report-output contract and artifact-specific report specs.
- `specs/skills/quality-skill/quality-skill.md` - align the reporting contract
  and skill expectations for consuming generated reports.
- `SPECIFICATION.md` - update only if minimum Evaluation Report semantics should
  explicitly call out Area-by-Factor summary breakdowns.

### Durable docs, skill, and examples

- `specs/skills/quality-skill/examples/0001-quality-eval/report-summary.md` -
  refresh the maintained generated summary example.
- `specs/skills/quality-skill/examples/0001-quality-eval/report.md` and
  `specs/skills/quality-skill/examples/0001-quality-eval/report.json` - refresh
  the maintained generated report examples when the report model changes.
- `specs/skills/quality-skill/examples/index.md` - update only if the example
  description needs to name the new breakdown.
- `docs/guides/write-functional-specs.md`, `docs/guides/work-with-okf.md`, and
  `docs/guides/work-with-change-cases.md` - document the artifact-spec filename
  convention for future durable specs.
- `skills/quality/SKILL.md` and `skills/quality/modes/evaluate.md` - update
  only if the skill's report consumption guidance needs to mention the new
  report shape.
- `README.md`, `docs/`, and scaffold files - no expected changes unless they
  contain report-output examples or guidance that becomes stale.

## Children

- [Functional spec](0048-area-factor-report-breakdown/spec.md) - what the
  report breakdown must provide.
- [Design doc](0048-area-factor-report-breakdown/design.md) - how the report
  model and renderers deliver the breakdown.

## Status

`In-Review`. The functional spec and design doc are settled and the code is
implemented in `internal/evaluation/` (report-model field rename, compact
Area-by-Factor breakdown, shared Area Breakdown renderer, and analysis-write
Factor-path validation), with regression tests and the three regenerated
`0001-quality-eval` golden report fixtures; `go test ./...` is green. The durable
`specs/reports/` artifact specs and the guide-spec renames already landed, which
durable-spec editing permits in any phase.
