---
type: Functional Specification
title: Candidate Actions Payload - functional spec
description: What the change must do to rename finding-local actions to candidateActions and make them citeable.
tags: [evaluation, findings, recommendations, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Candidate Actions Payload - functional spec

Companion to the
[Candidate Actions Payload](../0136-candidate-actions-payload.md) change case.
This spec states *what* the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Candidate actions are intentionally weaker than recommendations: they are
finding-local remediation leads captured where evidence is richest, not selected
Advice, priority, effort, or committed work. The old payload field name
`actions` made that boundary harder to preserve because the repo already has CLI
`nextActions` and will need recommendation options or actions later. The
Evaluation data contract should name this raw material as `candidateActions` and
give each candidate action a local ID so a future Advice phase can cite it
without making candidate actions report-visible recommendations.

## Scope

Covered: Evaluation Requirement Finding payloads, generated examples, JSON
schema output, durable Evaluation specs, report and closeout boundaries, and
runtime `/quality` evaluation guidance.

Deferred:

- recommendation or Advice record shape;
- candidate actions on Area Findings;
- done criteria, closure criteria, verification fields, priority, effort, ROI,
  or selected-option state;
- durable cross-run candidate action identity; and
- migration or mixed-shape compatibility for older Evaluation runs.

## Requirements

### Payload shape

- `RequirementAssessmentResult.findings[]` **MAY** include
  `candidateActions`, a list of finding-local candidate action objects.

  > Rationale: the field name should preserve the boundary between raw
  > finding-local remediation leads and future selected recommendations.
  >
  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/records/payload-kinds.md` - rename the Requirement
  > Finding extension.

- Area Findings **MUST NOT** include `candidateActions`.

  > Rationale: Area Findings are synthesized observations. Letting them carry
  > candidate actions would blur analysis into Advice.
  >
  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/records/payload-kinds.md` - keep candidate actions on
  > Requirement Findings only.

- Each candidate action object **MUST** include `id` and `description`, and
  **MAY** include `rationale`.

  > Rationale: `description` fits a possible remediation lead better than
  > `statement`, which is already the finding's evidence-backed claim.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` - define
  > candidate action object fields.

- Candidate action `id` values **MUST** be unique within the containing Finding
  and **MUST NOT** be treated as stable cross-run IDs.

  > Rationale: a local ID is enough for a later Advice phase to cite candidate
  > actions inside the run without inventing durable action identity.
  >
  > Durable spec: modify `specs/evaluation/records/json-conventions.md` - define
  > candidate action ID scope and selector use.

- Evaluation data validation **MUST** reject legacy finding field `actions`.

  > Rationale: QUALITY.md is early alpha; a clean break avoids dual-writer
  > behavior before recommendations are modeled.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` - document
  > the rejected legacy field.

### Evaluation behavior

- Requirement assessment **MAY** record `candidateActions` on `gap` and `risk`
  Requirement Findings when a local remediation lead is evident.

  > Durable spec: modify `specs/evaluation/routines/routine-contracts.md` and
  > `specs/skills/quality-skill/evaluation.md` - align routine and skill
  > behavior.

- Requirement assessment **MUST NOT** attach `candidateActions` to `strength`
  findings and **MUST NOT** synthesize, aggregate, prioritize, select, or present
  candidate actions as recommendations.

  > Durable spec: modify `specs/evaluation/routines/routine-contracts.md` and
  > `specs/skills/quality-skill/evaluation.md` - preserve the Evaluation v0
  > Advice boundary.

### Report and closeout behavior

- Evaluation reports **MUST NOT** render finding-local `candidateActions` in
  Evaluation v0.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` - update the
  > non-rendering rule to the renamed field.

- `/quality` evaluation closeouts **MUST NOT** present finding-local
  `candidateActions`.

  > Durable spec: modify `specs/skills/quality-skill/reporting.md` - keep
  > closeouts aligned with reports.

## Verification

- Validation tests **MUST** prove valid `candidateActions` are accepted and
  legacy `actions` are rejected.
- Validation tests **MUST** prove candidate action objects reject missing `id`,
  missing `description`, unknown fields, non-string fields, and duplicate IDs
  within the containing Finding.
- Validation tests **MUST** prove Area Findings reject `candidateActions`.
- Generated example and schema tests **MUST** prove `candidateActions` replaces
  `actions`.
- Report tests **MUST** prove candidate action descriptions are not rendered.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) - the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `SPECIFICATION.md` - rename the candidate action payload field and clarify
  local candidate action identity.
- `specs/evaluation/records/payload-kinds.md` - define `candidateActions`,
  candidate action fields, owner boundary, and legacy `actions` rejection.
- `specs/evaluation/records/json-conventions.md` - define candidate action ID
  scope and selector form.
- `specs/evaluation/reports/report-tree.md` - keep candidate actions out of
  report rendering under the new field name.
- `specs/evaluation/routines/routine-contracts.md` - align routine behavior with
  `candidateActions`.
- `specs/skills/quality-skill/evaluation.md` - align `/quality` evaluation
  behavior with `candidateActions`.
- `specs/skills/quality-skill/reporting.md` - keep evaluation closeouts from
  presenting candidate actions under the new field name.

### To rename

None.

### To delete

None.
