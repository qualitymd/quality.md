---
type: Change Case
title: Concern Finding Severity
description: Make Finding severity concern-only, requiring it for gaps and risks while forbidding it for strengths and notes.
status: Done
tags: [evaluation, findings, schema, reports]
timestamp: 2026-06-30T00:00:00Z
---

# Concern Finding Severity

Evaluation Findings currently require `severity` on every Finding, including
`strength` and `note`. That makes positive evidence and neutral context look
like quality concerns, then forces reports to explain or suppress meaningless
severity values. Severity should remain a triage signal for concerns: required
for `gap` and `risk`, forbidden for `strength` and `note`.

- [Functional spec](0185-concern-finding-severity/spec.md) - what the Finding
  severity contract must change.
- [Design doc](0185-concern-finding-severity/design.md) - how validation,
  schema generation, reports, examples, and guidance enforce the new contract.

## Motivation

Severity is useful when a Finding describes a present shortfall or future loss
path. It is misleading when attached to a strength or non-driving note, because
the value suggests urgency where the Finding is positive or contextual. Making
severity concern-only keeps structured data cleaner and lets report summaries
describe gaps and risks without pretending every Finding belongs on the same
severity scale.

## Scope

Covered:

- Requirement Finding data validation and generated JSON Schema;
- `qualitymd evaluation data example` payloads;
- generated Evaluation report severity display and summaries;
- report-gallery synthetic data and generated reports;
- runtime `/quality` skill guidance and durable Evaluation specs;
- tests, changelog, and OKF logs.

Deferred:

- renaming `note` to `info`;
- changing the severity value set for concerns;
- changing Finding type values or their semantics;
- adding compatibility readers, migration commands, or legacy aliases.

## Affected artifacts

- Code:
  - `internal/evaluation/data_contract.go` - make `severity` type-conditional
    in validation and generated schemas.
  - `internal/evaluation/data.go` - update example payloads.
  - `internal/evaluation/report_tree.go` - render severity only for concern
    Findings and summarize concern counts inline.
  - `internal/evaluation/evaluation_test.go` - cover validation, schema,
    examples, report rendering, and summary text.
  - `internal/cli/evaluation_test.go` - keep CLI data fixtures valid.
  - `scripts/report-gallery/main.go` - stop emitting severity for strengths.
- Durable specs:
  - `SPECIFICATION.md` - define severity as a concern-only Finding field.
  - `specs/cli/evaluation-data.md` - define validation and schema behavior.
  - `specs/evaluation/records/payload-kinds.md` - define the Finding Core field
    requirement.
  - `specs/evaluation/routines/routine-contracts.md` - define assessment
    output expectations.
  - `specs/evaluation/reports/report-tree.md` - define severity display and
    summary behavior.
  - `specs/skills/quality-skill/evaluation.md` - align skill-facing evaluation
    contract.
- Runtime skill:
  - `skills/quality/SKILL.md` - tell evaluators when to write or omit
    `severity`.
- Durable docs:
  - `docs/guides/reporting-design.md` - align report key-detail wording.
  - `CHANGELOG.md` - note the breaking Finding data contract and report change.
- Generated artifacts:
  - `internal/evaluation/evaluation-data.schema.json` - regenerate schema.
  - `examples/report-gallery/` - regenerate gallery outputs.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, and this case.
  - `specs/log.md`, `specs/evaluation/log.md`, `skills/quality/log.md`, and
    `docs/log.md`.

## Status

`Done`. Implemented and archived with validation, schema, reports, durable
specs, runtime guidance, tests, and generated examples aligned.
