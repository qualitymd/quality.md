---
type: Change Case
title: Candidate Actions Payload
description: Rename finding-local actions to candidateActions and give each candidate action a local ID.
status: Done
tags: [evaluation, findings, recommendations, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Candidate Actions Payload

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0136-candidate-actions-payload/spec.md) - what the case must do.
- [Design doc](0136-candidate-actions-payload/design.md) - how it's built, and why.

## Motivation

Evaluation can record finding-local candidate actions as raw material for a
later Advice phase, but the payload field is currently named `actions`. That
generic name is too close to CLI `nextActions`, future recommendation options,
issue handoff actions, and selected remediation work. Before recommendations are
modeled again, the Evaluation payload should name candidate actions exactly and
make them precise enough for future advice to cite without making them
recommendations.

## Scope

Covered:

- Rename Requirement Finding `actions` to `candidateActions`.
- Add a payload-local candidate action `id`.
- Keep the candidate action object lean: `id`, `description`, and optional
  `rationale`.
- Keep candidate actions only on Requirement Findings.
- Keep candidate actions out of Evaluation v0 reports and closeouts.
- Reject the legacy `actions` field without compatibility shims.

Deferred:

- Recommendation and Advice record modeling.
- Candidate actions on Area Findings.
- Done criteria, closure criteria, verification fields, priority, effort, or ROI.
- Durable cross-run candidate action identity.
- Automatic migration of older Evaluation runs.

## Affected artifacts

Derived by sweeping for `candidate action`, `candidate actions`,
`candidateActions`, and finding-local `actions` across code, durable specs,
runtime skill content, generated schemas, tests, and examples.

**Code**

- [x] `internal/evaluation/data_contract.go` - rename the Requirement Finding
      extension to `candidateActions`, add required candidate action `id`, and
      validate duplicate candidate action IDs within a containing Finding.
- [x] `internal/evaluation/data.go` - update generated examples.
- [x] `internal/evaluation/evaluation-data.schema.json` - regenerate the schema
      from the typed data contract.
- [x] `internal/evaluation/evaluation_test.go` - update validation, schema,
      example, and report-omission coverage.

**Durable specs** (substance in the [functional spec](0136-candidate-actions-payload/spec.md))

- [x] `SPECIFICATION.md` - rename the payload field and clarify candidate
      action identity.
- [x] `specs/evaluation/records/payload-kinds.md` - define
      `candidateActions`, candidate action fields, and legacy `actions`
      rejection.
- [x] `specs/evaluation/records/json-conventions.md` - document candidate
      action selector form.
- [x] `specs/evaluation/reports/report-tree.md` - keep candidate actions out of
      Evaluation v0 report rendering under the new field name.
- [x] `specs/evaluation/routines/routine-contracts.md` - update routine
      behavior to write `candidateActions`.
- [x] `specs/skills/quality-skill/evaluation.md` - align `/quality` evaluation
      behavior with the renamed field.
- [x] `specs/skills/quality-skill/reporting.md` - keep candidate actions out of
      evaluation closeouts under the new field name.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` - update runtime Area Finding exclusion wording.
- [x] `skills/quality/workflows/evaluate.md` - update candidate action payload
      guidance for Requirement Findings.

No planned impact: CLI command names, workspace layout, public recommendation
follow-up behavior, install docs, or `qualitymd init` model schema.

## Status

`Done`. Implemented and archived after code, durable specs, generated schema,
runtime skill guidance, and tests were updated.
