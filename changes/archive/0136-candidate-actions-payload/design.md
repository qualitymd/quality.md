---
type: Design Doc
title: Candidate Actions Payload - design
description: How Evaluation payload validation, examples, reports, and skill guidance adopt candidateActions.
tags: [evaluation, findings, recommendations, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Candidate Actions Payload - design

## Context

Answers the [functional spec](spec.md) for change case
[0136](../0136-candidate-actions-payload.md). The active 0135-shaped
Evaluation data contract already has a shared Finding Core, with Requirement
Findings extending it through an optional `actions` array. This change renames
that extension to `candidateActions`, adds local IDs to entries, and keeps the
reporting boundary unchanged.

## Approach

### Rename the Requirement Finding extension

Change `findingContract()` in `internal/evaluation/data_contract.go` from an
optional `actions` array to an optional `candidateActions` array. Keep
`candidateActionContract()` closed, with required `id` and `description` plus
optional `rationale`.

Do not add an alias for `actions`; closed object validation rejects the old
field.

### Validate local candidate action IDs

Extend the existing Evaluation payload usage validation pass to scan each
Requirement Finding's `candidateActions[]` and reject duplicate `id` values
within that containing Finding. Keep the ID local to the Finding, so two
different Findings can both use `action-001` without conflict.

### Update generated artifacts and tests

Update the Requirement Assessment example to use `candidateActions`. Regenerate
the Evaluation data schema from the typed contract. Update tests around invalid
candidate action shapes, valid examples, schema properties, Area Finding
rejection, and report omission.

### Update specs and skill guidance

Rename the field in `SPECIFICATION.md`, durable Evaluation specs, durable
`/quality` skill specs, and runtime evaluate guidance. Keep all wording clear
that candidate actions are raw material for a later Advice phase and are not
rendered or presented in Evaluation v0.

## Spec response

- **Payload shape** - satisfied by the renamed field, closed candidate action
  object contract, and duplicate-ID validation.
- **Evaluation behavior** - satisfied by durable routine and skill spec updates
  plus runtime evaluate guidance.
- **Report and closeout behavior** - satisfied by preserving report omission and
  updating the reporting specs to the new field name.
- **Verification** - satisfied by focused validation, schema/example, Area
  Finding rejection, and report omission tests.

## Alternatives

- **Keep `actions`.** Rejected. It is short, but it is ambiguous next to CLI
  `nextActions` and future recommendation actions/options.
- **Use `statement` instead of `description`.** Rejected. `statement` is the
  evidence-backed claim on a Finding; a candidate action is a possible
  remediation lead, so `description` is the clearer field.
- **Add done or closure criteria now.** Rejected. That information may help
  recommendations later, but adding it now risks turning candidate actions into
  underspecified recommendations before the Advice contract exists.
- **Make candidate action IDs unique across the payload or run.** Rejected.
  Candidate actions are finding-local raw material, so run-wide identity would
  add scope without current value.

## Trade-offs & risks

- This is another early-alpha breaking Evaluation payload change. That matches
  the repo's compatibility policy, but old authored payloads using `actions`
  will fail validation if rewritten.
- Future recommendation work may want richer candidate action metadata. Keeping
  the object lean now avoids pre-deciding the Advice contract, but that later
  work may add fields after the recommendation model is settled.

## Open questions

None.
