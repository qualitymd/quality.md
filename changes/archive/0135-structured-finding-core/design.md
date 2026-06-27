---
type: Design Doc
title: Structured Finding Core — design
description: How Evaluation payloads, reports, and skill guidance adopt the shared Finding Core.
tags: [evaluation, findings, reports, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Structured Finding Core — design

## Context

Answers the [functional spec](spec.md) for change case
[0135](../0135-structured-finding-core.md). Evaluation finding shape is
currently defined in `internal/evaluation/data_contract.go`, examples in
`internal/evaluation/data.go`, generated report rendering in
`internal/evaluation/report_tree.go`, durable contracts under `SPECIFICATION.md`
and `specs/evaluation/`, and runtime agent behavior in `skills/quality/SKILL.md`.

The implementation already validates payload objects through a typed contract
registry, rejects unknown fields, and derives JSON Schema from the same
contract. That makes the clean break straightforward: change the contract once,
update examples and report rendering, regenerate the schema, and update tests.

## Approach

### One typed Finding Core contract

Add one reusable `findingCoreContract()` in `data_contract.go` for the shared
fields:

- scalar fields: `id`, `type`, `severity`, `confidence`, `statement`,
  `condition`;
- `criteria[]`: `requirementId`, `ratingLevelId`, `criterion`, optional
  `rationale`;
- `cause`: `status`, `statement`, optional `rationale`, optional `evidence[]`;
- `effect`: `statement`, optional `rationale`, optional `ratingEffect`;
- `evidence[]`: `sourceRef`, `statement`, optional `rationale`.

Then define:

- Requirement Finding = Finding Core + optional `actions`;
- Area Finding = Finding Core + required non-empty `inputRefs` + optional
  `factorRelationships`.

Keep the existing finding type, severity, confidence, and Factor relationship
enums. Do not add legacy aliases for `description`, `summary`, or top-level
`rationale`; closed object validation rejects them.

### Examples teach the new shape

Update `requirementAssessmentExample` and Area Analysis examples to use the new
fields. Use domain-neutral statements and concise evidence refs. Keep candidate
actions only on a Requirement `gap` example so the existing advice boundary
stays visible.

### One report renderer

Introduce shared helpers that render any Finding Core:

- list table columns: ID, Statement, Type, Severity, Confidence, Effect, Cause;
- detail sections in order: Condition, Criteria, Cause, Effect, Evidence, then
  owner-specific relationships/inputs.

Requirement reports pass Requirement Findings into the shared renderer. Area
reports pass all Area Findings. Factor reports pass filtered Area Findings with
the matched Factor relationship. The Area/Factor wrappers add relationship and
input details; Requirement reports continue to omit candidate actions because
Evaluation v0 must not present recommendation-like content.

### Schema and tests

Regenerate `internal/evaluation/evaluation-data.schema.json` from the data
contract generator after code changes.

Update validation tests to:

- accept representative Requirement and Area Findings in the new shape;
- reject legacy fields on both finding kinds;
- verify Area Finding duplicate IDs and same-Area Factor relationship validation
  still work;
- assert the generated schema contains required Finding Core fields and cause
  status enum values;
- assert report Markdown contains the unified finding list headers and detail
  sections for Requirement, Area, and Factor reports.

## Spec response

- **Finding Core payload shape** — satisfied by the reusable data contract and
  owner-specific extensions.
- **IDs and references** — satisfied by keeping IDs payload-local and using
  existing `inputRefs.selector` strings for cross-payload references.
- **Finding analysis behavior** — satisfied through durable skill spec and
  runtime `SKILL.md` updates; the CLI validates structure but does not judge
  semantics.
- **Unified report rendering** — satisfied by a shared report renderer used from
  Requirement, Area, and Factor report paths.
- **Verification** — satisfied by validation, schema, example, and report tests.

## Alternatives

- **Keep `description` and `summary` as aliases.** Rejected. The repo's
  early-alpha policy favors clean breaks, and aliases would create precedence
  rules and dual-writer behavior.
- **Make cause a plain string.** Rejected. The important distinction is the
  evidence posture: verified, plausible, not assessed, or not applicable.
- **Move effect only to rating drivers.** Rejected. Rating drivers still explain
  rating selection, but findings need their own effect so Area and Factor
  reports can show why a condition matters before a reader opens rating details.
- **Require verified cause for every gap or risk.** Rejected. That would force
  root-cause analysis into ordinary quality evaluation and encourage
  speculation. Cause posture makes the limit explicit.

## Trade-offs & risks

- The new JSON shape is more verbose. That is acceptable because findings are
  the durable evidence-to-judgment object, and reports become easier to review.
- Old evaluation runs using the previous shape will not validate if rewritten
  through the new data contract. Evaluation explicitly does not define
  migrations or mixed-version run support.
- The skill can still write weak prose inside the new fields. The runtime
  guidance and QC rules reduce that risk, but judgment quality remains an agent
  responsibility.

## Open questions

None for this design.
