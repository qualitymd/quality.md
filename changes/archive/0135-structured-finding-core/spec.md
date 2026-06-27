---
type: Functional Specification
title: Structured Finding Core — functional spec
description: What the change must do to align Evaluation findings around statement, condition, criteria, cause, effect, and evidence.
tags: [evaluation, findings, reports, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Structured Finding Core — functional spec

Companion to the
[Structured Finding Core](../0135-structured-finding-core.md) change case. This
spec states *what* the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Evaluation findings are the bridge between evidence and judgment. They need to
be readable in reports, precise enough for agents to QC, and structured enough
for ratings to cite. The current split between Requirement Finding
`description`, Area Finding `summary`, and generic `rationale` weakens that
bridge: the finding's claim, observed condition, criteria, cause posture, and
effect can blur together.

The shared Finding Core makes the audit-equivalent structure explicit while
preserving QUALITY.md vocabulary. A finding states a short claim, records the
condition observed in evidence, names the Model criteria used for judgment,
states cause only to the level supported by evidence, explains the quality or
rating effect, and cites evidence.

## Scope

Covered: Evaluation Requirement Findings, Area Findings, generated examples,
JSON schema output, Markdown report rendering, durable Evaluation specs, and
runtime `/quality` skill guidance.

Deferred:

- older evaluation-run migration or mixed-shape compatibility;
- recommendation/advice rendering;
- global or cross-run finding identity;
- mandatory verified root-cause analysis; and
- new CLI commands for finding query or lookup.

## Requirements

### Finding Core payload shape

- `RequirementAssessmentResult.findings[]` and `AreaAnalysisResult.findings[]`
  **MUST** use one shared Finding Core containing `id`, `type`, `severity`,
  `confidence`, `statement`, `condition`, `criteria`, `cause`, `effect`, and
  `evidence`.

  > Rationale: one shape lets agents, validators, and reports treat Requirement
  > and Area Findings as the same kind of judgment object with different owners.
  >
  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/records/payload-kinds.md` — define the shared Finding
  > Core and apply it to both finding owners.

- A Finding Core `statement` **MUST** be a non-empty scalar short claim and
  **MUST NOT** replace the fuller observed `condition`.

  > Rationale: reports need a compact scan surface, while detailed review still
  > needs the observed condition separately.
  >
  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/reports/report-tree.md` — define `statement` as the report
  > list claim.

- A Finding Core `condition` **MUST** be a non-empty scalar describing the
  observed state or missing-evidence state.

  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/records/payload-kinds.md` — define `condition` as the
  > evidence-backed observation.

- A Finding Core `criteria` field **MUST** be a non-empty array whose entries
  identify the Model bar used for judgment with `requirementId`,
  `ratingLevelId`, `criterion`, and optional `rationale`.

  > Rationale: findings should point at the active QUALITY.md criteria instead
  > of embedding free-floating standards.
  >
  > Durable spec: modify `SPECIFICATION.md`,
  > `specs/evaluation/records/payload-kinds.md`, and
  > `specs/evaluation/routines/routine-contracts.md` — require model-grounded
  > finding criteria.

- A Finding Core `cause` field **MUST** be an object with `status`,
  `statement`, and optional `rationale` and `evidence`; `status` **MUST** be one
  of `verified`, `plausible`, `not_assessed`, or `not_applicable`.

  > Rationale: requiring a cause posture prevents unsupported root-cause claims
  > while still making the cause gap visible.
  >
  > Durable spec: modify `SPECIFICATION.md`,
  > `specs/evaluation/records/payload-kinds.md`, and
  > `specs/skills/quality-skill/evaluation.md` — require explicit cause posture
  > without requiring verified root cause.

- A Finding Core `effect` field **MUST** be an object with `statement` and
  optional `rationale` and `ratingEffect`.

  > Rationale: effect belongs on findings because it explains why the condition
  > matters to quality or rating; it is not a recommendation, priority, effort,
  > or ROI field.
  >
  > Durable spec: modify `SPECIFICATION.md`,
  > `specs/evaluation/records/payload-kinds.md`, and
  > `specs/evaluation/reports/report-tree.md` — distinguish effect from
  > recommendation and ranking fields.

- A Finding Core `evidence` field **MUST** be a non-empty array of objects with
  `sourceRef` and `statement`, plus optional `rationale`.

  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/records/payload-kinds.md` — require evidence entries on
  > findings instead of free-form evidence objects.

- Requirement Findings **MAY** carry finding-local candidate `actions`, and
  Area Findings **MUST NOT** carry `actions`.

  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/records/payload-kinds.md` — preserve candidate actions on
  > Requirement Findings only.

- Area Findings **MUST** continue to require non-empty `inputRefs` and **MAY**
  carry `factorRelationships`.

  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — keep Area
  > Finding provenance and Factor projection behavior.

- Evaluation data validation **MUST** reject legacy finding fields
  `description`, `summary`, and top-level `rationale` on Requirement and Area
  Findings.

  > Rationale: QUALITY.md is early alpha; a clean break avoids ambiguous dual
  > writers and unclear report precedence.
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
  > `specs/evaluation/records/json-conventions.md` — document the clean payload
  > shape and absence of compatibility migration.

### IDs and references

- Finding `id` values **MUST** be non-empty, payload-local IDs and **MUST NOT**
  be treated as stable cross-run or Model IDs.

  > Durable spec: modify `specs/evaluation/records/json-conventions.md` and
  > `specs/evaluation/records/payload-kinds.md` — clarify payload-local finding
  > identity.

- References to findings from rating drivers, analysis drivers, or reports
  **MUST** use routine `inputRefs` with selectors such as
  `findings[gap-001]`, qualified by the containing payload subject.

  > Durable spec: modify `specs/evaluation/records/json-conventions.md` and
  > `specs/evaluation/reports/report-tree.md` — keep finding references
  > contextual instead of naked global IDs.

### Finding analysis behavior

- The `/quality` skill **MUST** classify findings by type using the shared
  semantics: `gap` falls short of criteria, `risk` could plausibly cause future
  quality loss, `strength` supports or exceeds criteria, `unknown` records
  missing or ambiguous evidence, and `note` preserves relevant non-driving
  context.

  > Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
  > `specs/evaluation/routines/routine-contracts.md` — make type-specific
  > finding analysis part of the routine contract.

- The `/quality` skill **MUST NOT** claim `cause.status: verified` unless the
  finding evidence directly supports the cause statement.

  > Rationale: cause is useful only if the agent does not turn every plausible
  > explanation into a root-cause conclusion.
  >
  > Durable spec: modify `specs/skills/quality-skill/evaluation.md` — add the
  > cause-posture rule.

- The `/quality` skill **MUST** record `cause.status: not_assessed` when a
  `gap` or `risk` finding has enough evidence for condition and effect but not
  enough evidence for cause.

  > Durable spec: modify `specs/skills/quality-skill/evaluation.md` — keep
  > ratings moving from evidence without inventing cause.

### Unified report rendering

- Requirement, Area, and Factor reports **MUST** render findings through one
  shared list shape: ID, Statement, Type, Severity, Confidence, Effect, and
  Cause.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` — replace
  > separate Requirement and Area/Factor finding table shapes.

- Finding detail sections **MUST** render fields in this order: condition,
  criteria, cause, effect, evidence, relationships or inputs, and candidate
  actions only where the report contract permits them.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` — define the
  > shared finding-detail order.

- Evaluation reports **MUST NOT** render finding-local candidate actions in
  Evaluation v0.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` — preserve
  > the existing v0 recommendation boundary under the new detail renderer.

### Verification

- Generated JSON Schema for Evaluation data **MUST** expose the new Finding Core
  fields, nested enums, and required fields for both Requirement and Area
  Findings.

  > Durable spec: none.

- Tests **MUST** verify that data validation accepts representative Requirement
  and Area Findings using the new Finding Core and rejects legacy top-level
  finding fields.

  > Durable spec: none.

- Tests **MUST** verify that generated reports render Requirement, Area, and
  Factor Findings with the unified table and detail structure.

  > Durable spec: none.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) — the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `SPECIFICATION.md` — define Finding Core semantics in Evaluation and preserve
  the not-assessed-over-guessing cause posture.
- `specs/evaluation/records/payload-kinds.md` — define Finding Core payload
  fields, nested objects, Area Finding specialization, candidate action
  boundaries, and legacy field rejection.
- `specs/evaluation/reports/report-tree.md` — require unified findings table and
  detail rendering for Requirement, Area, and Factor reports.
- `specs/evaluation/routines/routine-contracts.md` — require agent routines to
  produce structured findings and type-specific analysis.
- `specs/evaluation/records/json-conventions.md` — clarify payload-local finding
  IDs and selector-qualified finding references.
- `specs/skills/quality-skill/evaluation.md` — align `/quality` skill behavior
  with Finding Core analysis and cause posture rules.

### To rename

None.

### To delete

None.
