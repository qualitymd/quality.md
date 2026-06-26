---
type: Change Case
title: Finding-level candidate actions
description: Capture non-binding, finding-local candidate actions on gap/risk findings at assessment time as typed raw material for a future Advise phase, without crossing the Evaluation v0 "no recommendations" boundary.
status: Done
tags: [evaluation, skill, schema, advise]
timestamp: 2026-06-26T00:00:00Z
---

# Finding-level candidate actions

A **Change Case** to turn the existing, untyped `findings[].actions` stub into a
typed, populated **candidate-action** field: the assessor records, on each
shortcoming finding, one or more non-binding leads describing what fixing it
might look like, captured where the evidence and context are richest. These are
*raw material* for a later recommendation-producing **Advise** phase — not
recommendations themselves, and not surfaced in the Evaluation v0 report.

This is option **A** ("per-finding candidate actions") from the recommendation
infrastructure discussion. Option **B** (the dedicated Advise phase that
synthesizes a final recommendation set) is captured separately as forward-looking
context in this case's
[advise sketch](0122-finding-candidate-actions/advise-sketch.md); it is **not**
in scope here.

Detail lives in:

- [Functional spec](0122-finding-candidate-actions/spec.md) — what this change
  must do.
- [Advise sketch](0122-finding-candidate-actions/advise-sketch.md) — non-binding
  design notes for the future Advise phase (option B) that shape why the harvest
  layer is the way it is.

No design doc is required unless implementation discovers a reusable abstraction
worth extracting.

## Motivation

The latest full evaluation of a real project surfaced 35 findings across 11
root requirements; **every one rendered `Actions | —`**. The `findings[].actions`
field already exists end-to-end — declared in the evaluation data contract
(`internal/evaluation/data_contract.go`), present in the generated schema, and
rendered by the report tree — but it is an untyped `arrayOfAny()` with no shape,
no authoring guidance, and is never populated.

Findings today are purely diagnostic. The remediation context is richest at
assessment time — the assessor has the source open and knows exactly *why* a gap
is a gap and what closing it would take. A future Advise phase that synthesizes a
final recommendation set will only see finding *descriptions*, not the source, so
discarding that context forces it to re-derive what the assessor already knew.

Capturing finding-local candidate actions as **typed raw material** preserves
that context cheaply, fills an existing stub rather than adding a new pipe, and
gives the Advise phase real inputs to cluster, dedup, and prioritize later. It
must do so **without** generating recommendations, which Evaluation v0 forbids:
candidate actions stay finding-local, non-binding, un-prioritized, and out of the
v0 report and closeout.

## Scope

Covered:

- Replace the untyped `findings[].actions` array with a **typed candidate-action
  object** (a required `description`, an optional `rationale`), keeping the field
  optional on a finding.
- Validate candidate-action objects when persisting `RequirementAssessmentResult`
  (reject unknown fields, missing `description`, wrong types), consistent with the
  payload-kinds validation contract.
- Update the bundled assessment prompt so the skill records candidate actions for
  `gap` and `risk` findings and omits them for `strength` findings.
- Include a candidate action in the `RequirementAssessmentResult` example payload
  so agents discover the shape.
- Keep candidate actions out of the Evaluation v0 report and closeout: they live
  in `data/` only, awaiting the Advise phase.
- Bring the affected durable specs and the bundled skill into conformance.

Deferred / non-goals:

- **The Advise phase itself (option B).** No synthesis, prioritization, dedup,
  option modeling, recommendation records, or new payload kinds. Sketched in
  [advise-sketch.md](0122-finding-candidate-actions/advise-sketch.md), built
  later.
- **No closeout or report presentation** of candidate actions, and **no
  cross-finding aggregation** — these would be recommendation generation, which
  Evaluation v0 prohibits.
- No change to rating semantics, the assessment → finding → rating chain, or the
  evidence-verification re-check (candidate actions are speculative leads and are
  not subject to the binding-finding re-check).
- No change to the QUALITY.md model schema or CLI command surface beyond the
  data-contract field shape and its generated schema/example.

## Affected artifacts

Derived by sweeping the evaluation data contract, generated schema, report
renderer, the bundled `/quality` skill assessment guidance, and the durable
evaluation and skill specs for where finding `actions` are defined, authored, or
rendered, and where the Evaluation v0 recommendation boundary is stated. Empty
kinds are deliberate.

### Code

- [x] `internal/evaluation/data_contract.go` — type the finding `actions` field
      (`arrayOfAny()` → array of candidate-action objects with a required
      `description` and optional `rationale`); add a candidate action to the
      `RequirementAssessmentResult` example.
- [x] `internal/evaluation/evaluation-data.schema.json` — regenerated artifact
      reflecting the typed shape (a staleness test enforces this).
- [x] `internal/evaluation/report_tree.go` — stop rendering the per-finding
      `Actions` row so candidate actions are not presented in the v0 report.
- [x] `internal/evaluation/*_test.go`, `internal/cli/evaluation_test.go` — cover
      candidate-action validation (accept valid, reject unknown/missing-field),
      the regenerated schema, and the report no longer emitting the Actions row.

### Durable specs

- [x] `specs/evaluation/routines/routine-contracts.md` — requirement assessment
      may record finding-local candidate actions and MUST NOT synthesize,
      aggregate, or prioritize them.
- [x] `specs/evaluation/reports/report-tree.md` — finding detail does not render
      candidate actions in v0.
- [x] `specs/skills/quality-skill/evaluation.md` — the skill records non-binding
      candidate actions on shortcoming findings as raw material for a later
      Advise phase.
- [x] `specs/skills/quality-skill/reporting.md` — candidate actions are not
      recommendations; the v0 report and closeout exclude them.
- [x] `SPECIFICATION.md` — note that a Finding MAY carry non-binding candidate
      actions, distinct from Advice-phase recommendations.
- [x] `specs/evaluation/log.md` and the top-level `specs/log.md` (the nearest
      enclosing log for the quality-skill component specs) — record the durable
      spec updates.

### Format spec

- [x] Covered above under Durable specs (`SPECIFICATION.md`) — one clarifying
      distinction between finding-local candidate actions and Advice
      recommendations.

### Durable docs (guides, README, bundled skill)

- [x] `skills/quality/workflows/evaluate.md` — assessment guidance: record
      candidate actions for `gap`/`risk` findings, omit for `strength`; ground
      the shape from the example payload.
- [x] `CHANGELOG.md` — add an Unreleased note for the typed candidate-action
      field and assessment guidance.

## Status

`Done`. The typed candidate-action field, its CLI validation, the example
payload, and the v0 report/closeout exclusion shipped; the durable specs,
bundled skill, format spec, changelog, and logs are in conformance; and the full
`mise run check` gate passes. Archived with its
[functional spec](0122-finding-candidate-actions/spec.md) and option-B
[advise sketch](0122-finding-candidate-actions/advise-sketch.md). See the
[status lifecycle](../index.md#status-lifecycle).
