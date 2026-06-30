---
type: Change Case
title: Requirement Findings Only
description: Make Requirement Findings the only Evaluation finding layer, remove Area Findings, and require rated results to be finding-backed.
status: Done
tags: [evaluation, findings, reports, skill, cli]
timestamp: 2026-06-27T00:00:00Z
---

# Requirement Findings Only

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0142-requirement-findings-only/spec.md) - what the case must
  do.
- [Design doc](0142-requirement-findings-only/design.md) - how it is built, and
  why.

## Motivation

Evaluation ratings are only useful when their evidence trail is inspectable.
Several runs exposed a failure mode where Requirements could be rated with no
Requirement Findings, making the rating effectively unsupported. At the same
time, Area Findings duplicated the Requirement evidence layer they were meant to
synthesize. Because Requirements can already connect to multiple Factors, and
Factor/Area analysis already has rating drivers, rationale, confidence, limits,
and incomplete inputs, a second finding layer adds noise without carrying unique
judgment.

Evaluation should have one evidence layer: Requirement Findings. Roll-up
judgment should preserve and cite those lower-level drivers without inventing
new findings. Recommendations can later synthesize action across Requirements.

## Scope

Covered:

- Make Requirement Findings the only Evaluation findings.
- Require rated Requirement results to be backed by at least one paired
  Requirement Finding.
- Require rated Requirement, Factor, and Area analysis results to carry
  `ratingDrivers`.
- Remove `AreaAnalysisResult.findings` and Area Finding Factor relationships
  from Evaluation schema version 3.
- Remove Area and Factor report Findings sections.
- Update `/quality evaluate` guidance, durable specs, generated schema, examples,
  tests, release notes, and logs.

Deferred:

- Recommendation generation or advice synthesis.
- Mechanical judgment sufficiency checks that try to decide whether a finding
  semantically proves a configured Rating Level.
- Migration or compatibility support for existing Evaluation schema version 2
  runs.

## Affected artifacts

Derived by sweeping for Area Findings, `AreaAnalysisResult.findings`,
`factorRelationships`, Finding Core, `ratingDrivers`, schema version, and rated
Requirement payloads.

**Code**

- [x] `internal/evaluation/types.go` - bump Evaluation data `SchemaVersion` to
      `3`.
- [x] `internal/evaluation/data_contract.go` - remove Area Finding contract and
      validation; add cross-payload effective-run validation for rated results,
      finding-backed Requirement ratings, and rating drivers.
- [x] `internal/evaluation/data.go` - remove Area Finding examples and emit v3
      examples.
- [x] `internal/evaluation/report_tree.go` - remove Area/Factor Finding
      rendering and helper code.
- [x] `internal/evaluation/evaluation-data.schema.json` - regenerate for schema
      version 3.
- [x] `internal/evaluation/evaluation_test.go` - replace Area Finding tests,
      add finding-backed rating validation coverage, and remove rated-empty
      Requirement fixtures.

**Durable specs** (substance in the [functional spec](0142-requirement-findings-only/spec.md))

- [x] `SPECIFICATION.md` - remove Area Findings, require rated Requirement
      results to trace to Requirement Findings, and move roll-up explanation to
      drivers/rationale.
- [x] `specs/evaluation/evaluation.md` - preserve the shared invariant that
      reports project persisted data without new findings.
- [x] `specs/evaluation/protocol.md` - keep protocol order while making Area
      analysis driver-only.
- [x] `specs/evaluation/routines/routine-contracts.md` - make Requirement
      assessment/rating and Factor/Area analysis contracts finding-backed and
      driver-backed.
- [x] `specs/evaluation/records/payload-kinds.md` - make Finding Core
      Requirement-only and delete Area Finding payload fields.
- [x] `specs/evaluation/reports/report-tree.md` - remove Area/Factor Findings
      sections and rendering rules.
- [x] `specs/cli/evaluation-data.md` - require `data set` validation for schema
      v3, rated result drivers, and paired Requirement findings.
- [x] `specs/skills/quality-skill/evaluation.md`,
      `specs/skills/quality-skill/workflows/evaluate.md`, and
      `specs/skills/quality-skill/reporting.md` - align `/quality evaluate` and
      reporting behavior.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` - remove Area Finding instructions and require
      rated Requirements to be backed by Requirement Findings.
- [x] `skills/quality/workflows/evaluate.md` - remove Area Finding authoring and
      make roll-up driver synthesis explicit.
- [x] `skills/quality/resources/SPECIFICATION.md` - refresh bundled
      specification resource.
- [x] `CHANGELOG.md` - note the breaking Evaluation schema/report/skill change.
- [x] `specs/log.md`, `specs/evaluation/log.md`,
      `specs/skills/quality-skill/workflows/log.md`,
      `skills/quality/log.md`, `skills/quality/workflows/log.md`, and
      `skills/quality/resources/log.md` - record durable/runtime updates.

No planned impact: lint rules, model parsing, setup authoring, install docs, or
historical archived Change Cases.

## Status

`Done`. Implemented and archived after Evaluation data validation, report
rendering, durable specs, runtime skill guidance, tests, generated schema,
release notes, and update logs were aligned. `go test ./...`,
`mise run fmt-md-check`, and `mise run check` pass.
