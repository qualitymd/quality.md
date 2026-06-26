---
type: Functional Specification
title: Remove info finding severity — functional spec
description: What the change must do to remove `info` from Evaluation finding severity.
tags: [evaluation, schema, reports, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Remove info finding severity — functional spec

Companion to the
[Remove info finding severity](../0132-remove-info-finding-severity.md) change
case. This spec states *what* the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

Evaluation finding `severity` is an adverse-concern scale. Keeping `info` in
that value set blurs the distinction between the field and finding `type: note`;
it also makes report ordering treat neutral informational observations as a
severity level. The active contract should keep `low` as the least-severe adverse
finding and use `note` for informational observations.

## Scope

Covered: the Evaluation finding severity value set for Requirement Findings and
Area Findings, emitted data schemas, write/verify validation, report sort/display
helpers, durable specs, bundled skill guidance, and release notes.

Deferred: conditional severity applicability by finding type. This case removes
one invalid value from the severity vocabulary; it does not yet make severity
required for `gap`/`risk` and disallowed for `strength`/`unknown`/`note`.

## Requirements

### Severity vocabulary

- Evaluation finding `severity` **MUST** use exactly these values:
  `critical`, `high`, `medium`, and `low`.

  > Rationale: `info` is not a severity of an adverse finding; informational
  > observations should be represented by finding `type: note`.
  >
  > Durable spec: modify `SPECIFICATION.md` and
  > `specs/evaluation/records/payload-kinds.md` — remove `info` from the active
  > Evaluation finding severity vocabulary and carry the rationale forward.

- `qualitymd evaluation data schema [<kind>]` **MUST** expose the same reduced
  finding severity enum wherever Requirement Finding or Area Finding `severity`
  appears.

  > Durable spec: modify `specs/cli/evaluation-data.md` — the data schema is the
  > authoritative payload contract and must expose the reduced severity set.

- `qualitymd evaluation data set` and `qualitymd evaluation data verify` **MUST**
  reject `severity: "info"` for Requirement Findings and Area Findings.

  > Durable spec: modify `specs/cli/evaluation-data.md` — data validation rejects
  > out-of-vocabulary enum values, now including stale `info` severity.

### Reporting and skill guidance

- Area and Factor report sorting **MUST** order known severities as `critical`,
  `high`, `medium`, then `low`; no report sort order **MUST** treat `info` as a
  known severity.

  > Durable spec: modify `specs/evaluation/reports/report-tree.md` — update
  > severity sort order.

- The `/quality` skill **MUST NOT** author `severity: "info"` for Evaluation
  findings. When it needs an informational observation, it **MUST** use finding
  `type: note` and omit any implication that `info` is a severity value.

  > Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
  > `specs/skills/quality-skill/workflows/evaluate.md` — route informational
  > observations to `type: note`.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) — the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `SPECIFICATION.md` — remove `info` from the active finding severity vocabulary
  and state that informational observations use finding `type: note`.
- `specs/evaluation/records/payload-kinds.md` — update Area Finding severity
  values to `critical`, `high`, `medium`, `low`.
- `specs/evaluation/reports/report-tree.md` — update severity report sort order.
- `specs/cli/evaluation-data.md` — record that schema and validation expose and
  enforce the reduced finding severity vocabulary.
- `specs/skills/quality-skill/evaluation.md` and
  `specs/skills/quality-skill/workflows/evaluate.md` — route informational
  findings to `type: note` rather than `severity: info`.

### To rename

None.

### To delete

None.
