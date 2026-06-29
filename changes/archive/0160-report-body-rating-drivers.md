---
type: Change Case
title: Report Body Rating Drivers
description: Keep rating drivers in structured Evaluation data while removing standalone Markdown driver sections from generated reports.
status: Done
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Report Body Rating Drivers

Generated Evaluation reports currently render standalone `## Rating Drivers`
sections with `inputRefs` dumped as compact JSON. This case keeps
`ratingDrivers` as structured Evaluation data and removes those sections from
human-facing Markdown report bodies.

- [Functional spec](0160-report-body-rating-drivers/spec.md) - what the reports
  must change.
- [Design doc](0160-report-body-rating-drivers/design.md) - how the renderer
  will do it.

## Motivation

Rating drivers are useful in the data layer: they keep Requirement, Factor, and
Area ratings traceable to lower-level routine outputs. In Markdown reports,
however, the current table is too literal for the review surface. It repeats
machine references as compact JSON and competes with the sections readers use to
decide what happened and what to do next: summary, findings, recommendations,
subject tables, limits, and incomplete inputs.

The report body should stay human-first. The structured payloads and source-data
frontmatter already preserve machine traceability.

## Scope

Covered:

- generated Evaluation Markdown reports that currently render standalone
  `## Rating Drivers` sections;
- preservation of structured `ratingDrivers` payload requirements;
- durable report specs, `/quality` reporting spec, report design guidance,
  focused tests, and regenerated report gallery output.

Deferred:

- changing the Evaluation data schema or rating-driver validation rules;
- adding clickable driver-link rendering;
- changing Requirement Finding detail rendering;
- changing Advice, recommendation, or ranking generation.

## Affected artifacts

Derived by sweeping for `## Rating Drivers`, `writeEvaluationDriversTable`,
`rating drivers`, report-gallery output, and report body guidance.

- Code:
  - `internal/evaluation/report_tree.go` - remove standalone rating-driver
    section rendering from Markdown reports.
  - `internal/evaluation/evaluation_test.go` - update report assertions so
    driver tables are absent while summary/findings/recommendations remain.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - distinguish structured
    rating-driver preservation from Markdown body rendering.
  - `specs/skills/quality-skill/reporting.md` - align `/quality` report
    expectations with no standalone rating-driver body sections.
- Durable docs:
  - `docs/guides/reporting-design.md` - clarify that rating drivers are source
    data, not a default visible report-body section.
- Generated examples:
  - `examples/report-gallery/` - regenerated generated reports.
- Change Case lifecycle:
  - `changes/index.md`, `changes/log.md`, and this case; archived on
    completion.

No planned impact: Evaluation schema version, `ratingDrivers` payload shape,
`qualitymd evaluation data set` validation, setup authoring, install docs, or
historical archived Change Cases.

## Status

`Done`. Implemented and archived. Generated run, Area, and Factor reports no
longer render standalone rating-driver body sections, while structured
`ratingDrivers` remain available in source payloads. `go test ./...`,
`mise run report-gallery-check`, and `mise run check` pass.
