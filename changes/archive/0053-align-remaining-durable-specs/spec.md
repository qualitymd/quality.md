---
type: Functional Specification
title: Align remaining durable specs - functional spec
description: Requirements for splitting remaining large durable specs into parent and component/artifact contracts.
tags: [specs, cli]
timestamp: 2026-06-22T00:00:00Z
---

# Align remaining durable specs - functional spec

Companion to the
[Align remaining durable specs](../0053-align-remaining-durable-specs.md) change
case. This spec states what the alignment must do; no design doc is required
unless implementation discovers a non-mechanical restructuring decision.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Scope

This change aligns durable internal functional specs with the parent/component
and artifact-spec guidance. It changes spec structure and links only.

Non-goals:

- changing runtime behavior, CLI command shapes, generated record schemas,
  report formats, lint rules, or update behavior;
- splitting `SPECIFICATION.md`, which remains the single primary format
  deliverable;
- changing bundled skill runtime files or CLI implementation.

## Requirements

### Evaluation record specs

`specs/evaluation-records.md` **MUST** remain the shared parent contract for
responsibility, runtime-not-OKF status, schema-version rules, historical-record
compatibility, and links to child record/artifact specs.

The evaluation-record spec set **MUST** include child specs for independently
reviewable runtime contracts: run folder, assessment result record, analysis
record, plan/planned coverage, debug log, recommendation record, and shared
report-output invariants.

### Lint specs

`specs/cli/lint.md` **MUST** remain the command contract for lint scope, flags,
repair behavior, command-level output summary, ordering/blocking, and exit
behavior.

The lint spec set **MUST** include child or sibling specs for the lint rule
system and lint output contract, so rule authoring/catalog requirements and JSON
finding/output schemas are independently reviewable.

### Update notice spec

`specs/cli/update.md` **MUST** remain the explicit update command contract.

The cross-command ambient update notice **MUST** move to a separate durable spec
because it applies to ordinary commands as a shared CLI behavior, not only to
the `qualitymd update` command.

### Bundle navigation

The alignment **MUST** update relevant OKF indexes so a reader can discover the
parent specs and new component/artifact specs.

The alignment **MUST** add `changes/log.md` and `specs/log.md` entries
summarizing the durable spec changes.

## Durable spec changes

### To add

- `specs/evaluation-records/index.md` - index for evaluation-record child specs.
- `specs/evaluation-records/run-folder.md` - runtime run-folder contract.
- `specs/evaluation-records/assessment-result-record.md` - assessment result
  record contract.
- `specs/evaluation-records/analysis-record.md` - analysis record contract.
- `specs/evaluation-records/plan-md.md` - `plan.md` and planned coverage
  contract.
- `specs/evaluation-records/debug-log-md.md` - `debug-log.md` contract.
- `specs/evaluation-records/recommendation-record-md.md` - recommendation record
  contract.
- `specs/evaluation-records/report-outputs.md` - shared generated report-output
  invariants.
- `specs/cli/lint-rules.md` - lint rule-system and rule-catalog contract.
- `specs/cli/lint-output.md` - lint findings, JSON, and human-output contract.
- `specs/cli/update-notice.md` - ambient update notice contract.

### To modify

- `specs/evaluation-records.md` - keep shared parent invariants and link to child
  specs.
- `specs/cli/lint.md` - keep command contract and link to lint rule/output
  specs.
- `specs/cli/update.md` - keep explicit update command contract and link to
  update notice spec.
- `specs/index.md` - list the new evaluation-record child folder.
- `specs/cli/index.md` - list new CLI component specs.
- `specs/log.md` - record the durable spec split.

### To rename

None.

### To delete

None.
