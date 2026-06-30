---
type: Functional Specification
title: Concern Finding Severity
description: Requirements for requiring severity on gap and risk Findings while forbidding it on strength and note Findings.
tags: [evaluation, findings, schema, reports]
timestamp: 2026-06-30T00:00:00Z
---

# Concern Finding Severity

This change makes Finding `severity` concern-only. It does not rename Finding
types, add new severity values, or add legacy compatibility behavior.

The key words "MUST", "MUST NOT", and "MAY" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

The current Finding core requires `severity` for every Finding type. That
overloads severity: it is meaningful for `gap` and `risk`, but artificial for
`strength` and `note`. Reports already use severity as a triage cue, so requiring
it on non-concern Findings either misleads readers or forces renderers to hide a
required value. The data contract should represent the judgment directly:
severity belongs to concerns only.

## Requirements

1. Evaluation Finding `severity` **MUST** be required when `type` is `gap` or
   `risk`.

   > Rationale: Gaps and risks are quality concerns, and severity is the
   > concern-triage signal reports and advice can use.
   >
   > Durable spec: modify `SPECIFICATION.md`, `specs/cli/evaluation-data.md`,
   > `specs/evaluation/records/payload-kinds.md`,
   > `specs/evaluation/routines/routine-contracts.md`, and
   > `specs/skills/quality-skill/evaluation.md`.

2. Evaluation Finding `severity` **MUST NOT** be present when `type` is
   `strength` or `note`.

   > Rationale: Positive evidence and non-driving context are not concerns;
   > assigning them concern severity is semantically false.
   >
   > Durable spec: modify `SPECIFICATION.md`, `specs/cli/evaluation-data.md`,
   > `specs/evaluation/records/payload-kinds.md`,
   > `specs/evaluation/routines/routine-contracts.md`, and
   > `specs/skills/quality-skill/evaluation.md`.

3. When `severity` is present, validation **MUST** accept only `critical`,
   `high`, `medium`, or `low`; validation **MUST** continue to reject `info`.

   > Durable spec: modify `SPECIFICATION.md` and
   > `specs/cli/evaluation-data.md`.

4. `qualitymd evaluation data schema RequirementAssessmentResult` and the
   full-surface schema **MUST** express the same type-conditional `severity`
   requirement and prohibition as validation.

   > Rationale: Schema output is a discovery surface for agents and automation;
   > it must not advertise a looser or stricter shape than `data set`.
   >
   > Durable spec: modify `specs/cli/evaluation-data.md`.

5. `qualitymd evaluation data example RequirementAssessmentResult` **MUST** emit
   a valid representative example under the concern-only severity contract.

   > Durable spec: modify `specs/cli/evaluation-data.md`.

6. Generated ranked Finding tables and Requirement Finding summary tables
   **MUST** render concern severity for `gap` and `risk` Findings and `—` for
   `strength` and `note` Findings.

   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

7. Generated report severity summaries and highest-severity calculations
   **MUST** consider only `gap` and `risk` Findings.

   > Rationale: A report's severity summary is a concern summary. Non-concern
   > Findings should not affect the highest concern severity.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

8. The run report's full Findings link summary **MUST** include an inline
   complete count summary when ranked Findings exist, ordered as gaps, risks,
   strengths, then notes; gap and risk segments **MUST** include observed
   severity counts ordered `critical`, `high`, `medium`, `low`; zero-count
   Finding types and zero-count severity values **MUST NOT** render.

   Example:

   ```markdown
   **Full findings report:** [findings.md](findings.md) (7 total: 4 Gaps - 1 Critical, 2 High, 1 Medium; 1 Risk - 1 High; 2 Strengths)
   ```

   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

9. Runtime `/quality` evaluation guidance **MUST** tell evaluators to write
   `severity` only on `gap` and `risk` Findings and to omit it from `strength`
   and `note` Findings.

   > Durable spec: modify `specs/skills/quality-skill/evaluation.md`.

10. This change **MUST NOT** add backward-compatibility shims, fallback readers,
    legacy aliases, migration commands, or dual writers for the old all-Finding
    severity shape.

    > Rationale: QUALITY.md is early alpha, and a clean break keeps the current
    > contract easier for agents and maintainers to reason about.
    >
    > Durable spec: none.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` - define `severity` as required for `gap`/`risk`, forbidden
  for `strength`/`note`, and limited to `critical`, `high`, `medium`, and `low`
  when present.
- `specs/cli/evaluation-data.md` - define validation, schema, and example
  behavior for concern-only severity.
- `specs/evaluation/records/payload-kinds.md` - define the type-conditional
  Finding Core `severity` contract.
- `specs/evaluation/routines/routine-contracts.md` - define the assessment
  routine's type-conditional severity output.
- `specs/evaluation/reports/report-tree.md` - define concern-only severity
  rendering, highest severity calculations, and the inline full Findings summary.
- `specs/skills/quality-skill/evaluation.md` - align skill-facing evaluation
  rules with the new Finding shape.

### To rename

None

### To delete

None
