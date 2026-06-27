---
type: Change Case
title: Finding Basis
description: Rename finding-local cause posture to basis across evaluation records and reports.
status: Done
tags: [evaluation, records, reports, terminology]
timestamp: 2026-06-27T00:00:00Z
---

# Finding Basis

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0148-finding-basis/spec.md) - what the case must do.
- [Design doc](0148-finding-basis/design.md) - how it's built, and why.

## Motivation

Requirement Findings need a field for the supported explanation or grounding
posture behind an observation. The current `cause` field fits gaps and risks,
but reads awkwardly for strengths, unknowns, and notes. `basis` is shorter and
more neutral while preserving the same evidence discipline and support-status
model.

## Scope

Covered:

- Rename the Finding Core object field from `cause` to `basis`.
- Rename generated report labels from `Cause` to `Basis`.
- Keep the nested `status`, `statement`, `rationale`, and `evidence` shape.
- Keep the status values `verified`, `plausible`, `not_assessed`, and
  `not_applicable`.
- Update runtime skill guidance, durable specs, data contracts, examples,
  report rendering, tests, scaffold comments, and release notes.

Deferred:

- Adding backward-compatibility aliases, fallback readers, dual writers, or
  migration commands for `cause`.
- Renaming unrelated ordinary-English uses of "cause".
- Rewriting archived Change Cases or append-only historical logs.

## Affected artifacts

Derived by sweeping for `cause`, `Cause`, `cause.status`, `findingCause`,
`Cause Evidence`, `supported cause`, and `cause posture`.

**Code**

- [x] `internal/evaluation/data_contract.go` - Finding Core validation.
- [x] `internal/evaluation/data.go` - generated data examples.
- [x] `internal/evaluation/evaluation-data.schema.json` - generated evaluation
      data schema.
- [x] `internal/evaluation/report_tree.go` - generated Markdown report labels
      and finding summaries.
- [x] `internal/evaluation/display.go` - status title helper naming.
- [x] `internal/evaluation/evaluation_test.go` - data-schema and report
      assertions.
- [x] `internal/cli/evaluation_test.go` - CLI evaluation-data fixtures.

**Format spec and durable specs** (substance in the [functional spec](0148-finding-basis/spec.md))

- [x] `SPECIFICATION.md` - Finding terminology and assessment semantics.
- [x] `skills/quality/resources/SPECIFICATION.md` - bundled spec symlink to the
      root specification.
- [x] `specs/evaluation/records/payload-kinds.md` - Finding Core data
      contract.
- [x] `specs/evaluation/routines/routine-contracts.md` - Requirement
      assessment routine contract.
- [x] `specs/evaluation/reports/report-tree.md` - Requirement Finding report
      rendering contract.
- [x] `specs/evaluation/log.md` and `specs/log.md` - durable spec bundle logs.
- [x] `specs/skills/quality-skill/evaluation.md` - skill evaluation behavior.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` - runtime evaluation guidance.
- [x] `internal/scaffold/skeleton.md` - starter model assessment comment.
- [x] `CHANGELOG.md` - release-note entry for the breaking data/report
      vocabulary change.
- [x] `changes/log.md`, `changes/index.md`, and `changes/archive/index.md` -
      Change Case lifecycle.

## Status

`Done`. Implemented and archived after evaluation records, reports, durable
specs, runtime skill guidance, scaffold comments, tests, release notes, and logs
were updated. `go test ./...` and `mise run fmt-md-check` pass.
