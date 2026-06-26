---
type: Change Case
title: Remove info finding severity
description: Remove `info` from the Evaluation finding severity vocabulary so severity remains an adverse-finding scale.
status: Done
tags: [evaluation, schema, reports, skill]
timestamp: 2026-06-26T00:00:00Z
---

# Remove info finding severity

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0132-remove-info-finding-severity/spec.md) — what the case must do.
- [Design doc](0132-remove-info-finding-severity/design.md) — how it's built, and why.

## Motivation

After Area Findings landed, the shared Evaluation finding severity vocabulary
still included `info`. That makes it easy to encode a neutral note or positive
strength as a "severity", and it weakens the meaning of the field in reports.
Severity should remain a scale for adverse findings, with `low` as the least
severe retained concern. Neutral informational observations belong under finding
`type: note`, not under `severity: info`.

This case removes `info` from the active Evaluation finding severity vocabulary
across Requirement Findings and Area Findings, the emitted data schema, data-set
validation, report sorting/display, durable specs, and bundled skill guidance.

## Scope

Covered:

- Remove `info` from the Evaluation finding severity enum used by
  Requirement Finding `severity` and Area Finding `severity`.
- Make `qualitymd evaluation data schema [<kind>]` expose the reduced severity
  set: `critical`, `high`, `medium`, `low`.
- Make `qualitymd evaluation data set` and `data verify` reject `severity:
  "info"` wherever an Evaluation finding severity is validated.
- Update report sorting and display code so `info` is no longer a known finding
  severity.
- Update durable specs and bundled skill/runtime resources so agents use
  `type: note` for informational findings rather than `severity: info`.
- Prepare release notes and skill compatibility for a new release.

Deferred:

- Conditional severity applicability (`severity` required for `gap`/`risk` and
  disallowed for other finding types). That is the next semantic rule, but this
  case only removes `info` from the value set.
- Historical migration or compatibility readers. QUALITY.md is early alpha, and
  active data validation should reject stale `info` severities.

## Affected artifacts

Derived by sweeping for `severity`, `info`, `findingSeverity`, and Area Finding
report ordering across code, specs, generated schema, and bundled skill content.

**Code**

- [x] `internal/evaluation/data_contract.go` — remove `info` from finding
      severity enum declarations for Requirement and Area Findings.
- [x] `internal/evaluation/display.go` — remove `info` as a known finding
      severity display title.
- [x] `internal/evaluation/report_tree.go` — remove `info` from known severity
      sort order; unknown values continue to sort after known values defensively.
- [x] `internal/evaluation/evaluation-data.schema.json` — regenerate from the
      typed contract.
- [x] Tests under `internal/evaluation/` — update schema expectations and add
      rejection coverage for `severity: "info"`.

**Durable specs** (substance in the [functional spec](0132-remove-info-finding-severity/spec.md))

- [x] `SPECIFICATION.md` — carry forward the reduced finding severity vocabulary
      and rationale.
- [x] `specs/evaluation/records/payload-kinds.md` — update the Area Finding
      severity value set.
- [x] `specs/evaluation/reports/report-tree.md` — update severity sort order.
- [x] `specs/cli/evaluation-data.md` — note data schema / data set expose and
      enforce the reduced severity value set.
- [x] `specs/skills/quality-skill/evaluation.md` and
      `specs/skills/quality-skill/workflows/evaluate.md` — guide Area Finding
      authoring away from `severity: info`.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` — update runtime Area Finding guidance.
- [x] `skills/quality/workflows/evaluate.md` — update Area Finding authoring
      guidance.
- [x] `skills/quality/resources/SPECIFICATION.md` — mirror the public
      specification resource.
- [x] `skills/quality/log.md` and related runtime logs — append entries when
      landing.
- [x] `CHANGELOG.md` — release note and compatibility entry.

No planned impact: `README.md`, install/scaffold files, or `QUALITY.md` model
frontmatter schema.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Implemented
and archived. Evaluation finding severity now excludes `info` in the typed
contract, schema artifact, validation, report helpers, durable specs, bundled
skill guidance, and release notes; `go test ./internal/evaluation` passes.
