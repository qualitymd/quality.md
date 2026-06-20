---
type: Functional Specification
title: Evaluation history compatibility - functional spec
description: Requirements for tolerant inspection of historical, malformed, partial, or hand-edited evaluation runs without adding migrations.
tags: [evaluation, records, compatibility, status]
timestamp: 2026-06-20T00:00:00Z
---

# Evaluation history compatibility - functional spec

Companion to
[Evaluation history compatibility](../0043-evaluation-history-compatibility.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

Evaluation runs are intentionally runtime artifacts, but they are also durable
history. During active development, users will have runs that were not produced
by the current skill/CLI pair: older records, partial runs, hand edits, copied
folders, malformed files, or future records with a different writer stamp. If
the current loader treats every mismatch as a fatal process error, ordinary
maintenance breaks even when the user only wants status, history, or a fresh
evaluation.

The compatibility posture is deliberately narrow: keep current writers strict,
avoid migrations, inspect existing history best-effort, and make a specific run
non-reportable when its records cannot be trusted.

## Scope

Covered: existing evaluation run folders and records read by status, list,
latest-run resolution, record-list, report-build, and report-gate workflows.

Deferred / non-goals: no migration framework; no automatic repair; no rewriting
historical runs; no compatibility writer for obsolete shapes; no acceptance of
invalid input to `qualitymd evaluation assessment add`,
`qualitymd evaluation analysis set`, or
`qualitymd evaluation recommendation add`.

## Requirements

### Strict current writers

Evaluation record write commands **MUST** continue to validate current payloads
against the active record contract and **MUST** stamp records with the current
`schemaVersion`.

> Rationale: The project wants compatibility when reading history, not a second
> write path for old or invalid shapes.

### Tolerant run inspection

Commands that inspect evaluation history **SHOULD** load enough run metadata to
describe the run even when one or more record files are malformed,
schema-incompatible, incomplete, or unreadable.

> Rationale: A bad historical record should not prevent status, list, or new-run
> workflows from working.

### Compatibility gaps

When a run contains a record that cannot be trusted under the current contract,
the run **MUST** surface a reportability gap that names the record path and a
human-actionable reason. The gap **MUST** make the run non-reportable.

At minimum, gaps must distinguish malformed JSON or frontmatter, missing
`schemaVersion`, unsupported `schemaVersion`, and current-schema records that
are structurally incomplete.

### Status and list resilience

`qualitymd status` and `qualitymd evaluation list` **MUST NOT** abort solely
because one discovered run contains incompatible records. They **MUST** include
the run as incomplete/problematic, preserve counts that can be determined
without trusting the malformed record body, and keep next actions aimed at
inspection or creation of a fresh run.

### Latest-run resolution

Resolving `--latest` **MUST** remain based on discovered run folder order, not on
strict reportability. If the latest run is incompatible, commands that operate
on it **MUST** report that run's compatibility gaps instead of silently skipping
to an older reportable run.

### Report and gate behavior

`qualitymd evaluation report build` and `qualitymd evaluation report gate`
**MUST** refuse to render or gate an incompatible selected run. The error or
status output **MUST** identify the first blocking compatibility gap and point to
`qualitymd evaluation status <run>` for the full gap list.

### Skill behavior

The `/quality` skill **MUST** treat incompatible historical runs as evaluation
history status, not as evidence about subject quality. It should recommend
inspecting the run or creating a fresh evaluation rather than hand-editing or
manually migrating records.

### Tests and fixtures

The implementation **MUST** include regression coverage for missing
`schemaVersion`, unsupported future `schemaVersion`, malformed JSON assessment
or analysis records, malformed recommendation frontmatter, structurally
incomplete current-schema records, and history containing both good and bad
runs.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation-records.md` — add the strict-writer / tolerant-reader
  compatibility posture and define incompatible historical records as
  non-reportable gaps.
- `specs/cli/evaluation-list.md` — require listing to survive incompatible runs
  and classify them as incomplete/problematic.
- `specs/cli/evaluation-status.md` — describe compatibility gaps and status
  output for malformed, unsupported, or structurally incomplete records.
- `specs/cli/evaluation-report.md` — require report build/gate to refuse
  incompatible selected runs with clear diagnostics.
- `specs/cli/status.md` — align project status readiness and next actions with
  incompatible evaluation history.
- `specs/skills/quality-skill/quality-skill.md` — encode the skill-facing policy
  for stale or incompatible historical runs.

### To delete

None
