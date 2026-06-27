---
type: Functional Specification
title: Finding Basis
description: Requirements for renaming finding-local cause posture to basis.
tags: [evaluation, records, reports, terminology]
timestamp: 2026-06-27T00:00:00Z
---

# Finding Basis

This change case renames the Finding Core field that records an observation's
supported explanatory or grounding posture from `cause` to `basis`.

The key words "MUST", "MUST NOT", and "MAY" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background

The structured Finding Core is intentionally shared by strengths, gaps, risks,
unknowns, and notes. `cause` makes that shared core read like failure diagnosis:
natural for a gap or risk, but strained for a strength and often noisy for an
unknown or note. `basis` keeps the evidence-support posture while making the
field neutral across finding types.

## Requirements

Evaluation data contracts **MUST** require `basis` on every Finding Core object
and **MUST NOT** accept `cause` as a Finding Core field.

> Rationale: QUALITY.md is early alpha, so the clean contract is clearer than a
> dual-field compatibility period that leaves agents unsure which field to
> write.
>
> Durable spec: modify `SPECIFICATION.md` and
> `specs/evaluation/records/payload-kinds.md` - rename the Finding Core field
> and preserve required-field status.

`basis` **MUST** be an object with required `status` and `statement`, optional
`rationale`, and optional `evidence`.

> Durable spec: modify `specs/evaluation/records/payload-kinds.md` - rename the
> object contract while preserving its nested shape.

`basis.status` **MUST** keep the existing status values `verified`,
`plausible`, `not_assessed`, and `not_applicable`.

> Rationale: the vocabulary problem is the field name, not the status model.
> Keeping the enum stable preserves the judgment distinction already used by
> evaluators and reports.
>
> Durable spec: modify `SPECIFICATION.md` and
> `specs/evaluation/records/payload-kinds.md` - rename the status owner without
> changing allowed values.

The evaluation semantics **MUST** define `basis` as the finding-local
explanation or support posture for the observed condition, distinct from the
finding's checkable `evidence` and from report-level "evidence basis" summaries.

> Rationale: `basis` is already used in prose phrases such as "evidence basis";
> the durable semantics need to prevent readers from treating the new field as a
> duplicate of evidence provenance.
>
> Durable spec: modify `SPECIFICATION.md` and
> `specs/evaluation/routines/routine-contracts.md` - add the distinction to the
> Finding and Requirement assessment semantics.

Evaluation guidance **MUST NOT** claim `basis.status: verified` unless the
finding evidence directly supports the basis statement. When a `gap` or `risk`
has enough evidence for condition and effect but not for basis, evaluation
guidance **MUST** require `basis.status: not_assessed` rather than an invented
basis.

> Rationale: the rename should not weaken the no-overclaim rule that made the
> original cause posture useful.
>
> Durable spec: modify `SPECIFICATION.md`,
> `specs/evaluation/routines/routine-contracts.md`, and
> `specs/skills/quality-skill/evaluation.md` - preserve the verified-basis
> evidence rule and gap/risk fallback.

Evaluation guidance for `strength` findings **MUST** allow
`basis.status: verified` when the positive condition's basis is directly
supported by cited evidence, and **MAY** use `basis.status: not_applicable` when
no separate basis beyond the cited evidence is claimed.

> Rationale: the rename was motivated by strengths; the guidance should make the
> intended positive-use case explicit rather than turning `not_applicable` into
> a reflex.
>
> Durable spec: modify `specs/evaluation/routines/routine-contracts.md` and
> `specs/skills/quality-skill/evaluation.md` - add strength-specific basis
> guidance.

Generated reports **MUST** render the finding summary column, finding detail
section, and nested basis evidence section as `Basis`, `Basis`, and
`Basis Evidence`.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` - replace
> `Cause` report labels with `Basis` labels.

The runtime `/quality` skill **MUST** instruct agents to write `basis`,
`basis.status`, and `basis.rationale` in Requirement Findings, and **MUST NOT**
instruct them to write `cause` for the current Finding Core.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md` - align skill
> behavior with the renamed field.

Generated examples, schema output, tests, and scaffold comments **MUST** use
`basis` for current Finding Core examples and fixtures.

> Durable spec: none.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` - rename Finding Core `cause` semantics to `basis` and
  preserve support-status rules (per the data contract, status, basis semantics,
  and no-overclaim requirements).
- `specs/evaluation/records/payload-kinds.md` - rename the Finding Core object
  contract from `cause` to `basis` (per the data contract, nested shape, and
  status requirements).
- `specs/evaluation/routines/routine-contracts.md` - rename Requirement
  assessment guidance and preserve basis support rules (per the basis semantics,
  no-overclaim, and strength guidance requirements).
- `specs/evaluation/reports/report-tree.md` - rename Requirement Finding report
  labels (per the report rendering requirement).
- `specs/skills/quality-skill/evaluation.md` - align skill evaluation behavior
  with `basis` (per the no-overclaim, strength guidance, and runtime skill
  requirements).

### To rename

None

### To delete

None
