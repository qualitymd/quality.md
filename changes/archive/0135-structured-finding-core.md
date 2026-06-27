---
type: Change Case
title: Structured Finding Core
description: Align Evaluation findings around statement, condition, criteria, cause, effect, and evidence.
status: Done
tags: [evaluation, findings, reports, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Structured Finding Core

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0135-structured-finding-core/spec.md) — what the case must do.
- [Design doc](0135-structured-finding-core/design.md) — how it's built, and why.

## Motivation

Evaluation findings currently carry loose `description`, `summary`, and
`rationale` fields depending on whether the finding is Requirement-local or
Area-level. That makes findings easy to write but harder to audit, compare, and
render consistently. A reader needs the short claim, observed condition, Model
criteria, cause posture, quality effect, and evidence trail to understand why a
finding matters and how it affects a rating.

This case introduces a shared Finding Core for Requirement Findings and Area
Findings. The structure is functionally equivalent to the internal-audit
condition/criteria/cause/effect frame while keeping QUALITY.md terminology and
domain-agnostic evaluation semantics.

## Scope

Covered:

- Replace loose Requirement and Area Finding fields with a shared core:
  `statement`, `condition`, `criteria`, `cause`, `effect`, and `evidence`.
- Keep short, run-local finding IDs and payload-local uniqueness rules.
- Put rationale on the specific nested field it explains instead of on the
  finding as a whole.
- Render Requirement, Area, and Factor Findings through one unified report
  shape.
- Update the `/quality` skill evaluation guidance so each finding type uses the
  new analysis pattern.

Deferred:

- Recommendation generation and advice reports.
- Durable cross-run finding identity.
- Automatic migration of older evaluation runs.
- Exhaustive root-cause analysis beyond what available evidence supports.

## Affected artifacts

Derived by sweeping for `finding`, `findings`, `summary`, `description`,
`rationale`, `condition`, `criteria`, `cause`, `effect`,
`AreaAnalysisResult.findings`, and `RequirementAssessmentResult.findings`
across code, specs, reports, skill runtime content, and examples.

**Code**

- [x] `internal/evaluation/data_contract.go` — define and validate the shared
      Finding Core for Requirement and Area Findings.
- [x] `internal/evaluation/data.go` — update generated examples to the new
      finding shape.
- [x] `internal/evaluation/display.go` — add cause-status display titles.
- [x] `internal/evaluation/report_tree.go` — render Requirement, Area, and
      Factor Findings through the unified report shape.
- [x] `internal/evaluation/evaluation-data.schema.json` — regenerate the
      schema from the typed data contract.
- [x] `internal/evaluation/evaluation_test.go` — update validation, schema, and
      report tests for the new finding fields.

**Durable specs** (substance in the [functional spec](0135-structured-finding-core/spec.md))

- [x] `SPECIFICATION.md` — define Finding Core semantics in the Evaluation
      contract.
- [x] `specs/evaluation/records/payload-kinds.md` — define the shared payload
      shape and Area Finding specialization.
- [x] `specs/evaluation/reports/report-tree.md` — require one findings table
      and detail order across Requirement, Area, and Factor reports.
- [x] `specs/evaluation/routines/routine-contracts.md` — require agent
      routines to write structured findings.
- [x] `specs/evaluation/records/json-conventions.md` — clarify local finding
      ID scoping and selector use.
- [x] `specs/skills/quality-skill/evaluation.md` — align skill behavior with
      the Finding Core and type-specific analysis rules.
- [x] `specs/evaluation/log.md` and `specs/log.md` — record durable spec
      updates.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` — update runtime evaluation guidance and QC
      checks for structured findings.
- [x] `README.md` — inspected; no update needed because it does not expose
      Evaluation finding payload shape.
- [x] `internal/scaffold/skeleton.md` — update comments that mention cause and
      next steps if needed.

No planned impact: CLI command names, workspace layout, install docs,
`qualitymd init` generated model schema, or release notes.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Implemented
and archived. Evaluation findings now use a shared Finding Core, reports render
Requirement, Area, and Factor Findings through one structure, and `/quality`
runtime guidance now performs type-specific structured finding analysis.
