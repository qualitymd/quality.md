---
type: Functional Specification
title: Planned coverage status
description: Compare optional planned coverage metadata against written evaluation records.
tags: [evaluation, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Planned coverage status

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

This change covers an optional planned-coverage runtime artifact for evaluation
runs and status checks that compare the artifact to written assessment and
analysis records.

It does not make every run cover every model requirement. It does not add record
replacement, deletion, or superseding behavior.

## Requirements

An evaluation run **MAY** include planned coverage metadata that lists intended
assessment requirements and analysis targets for the run.

Planned coverage metadata **MUST** be optional. A run with no planned coverage
metadata **MUST** keep the current status behavior.

When planned coverage metadata exists, `qualitymd evaluation show-status`
**MUST** compare planned assessment entries against written assessment records.
The comparison **MUST** use ordered `targetPath` plus `requirement`.

When planned coverage metadata exists, `qualitymd evaluation show-status`
**MUST** compare planned analysis entries against written analysis records. The
comparison **MUST** use ordered `targetPath`.

When planned coverage metadata exists and planned assessment records are
missing, `show-status` **MUST** report a planned-coverage gap that identifies
each missing planned assessment.

When planned coverage metadata exists and planned analysis records are missing,
`show-status` **MUST** report a planned-coverage gap that identifies each
missing planned analysis.

When planned coverage metadata exists and extra assessment or analysis records
are present outside the plan, `show-status` **SHOULD** report them as unexpected
planned-coverage records.

A planned-coverage gap **MUST** make the run non-reportable unless the gap is
explicitly defined as a warning in the durable status spec before
implementation.

The planned coverage artifact **MUST NOT** replace `design.md` or `plan.md`.
Those files remain the human-readable record of scope, effort, and limitations.

## Deferred

- CLI commands for generating planned coverage from a model.
- Automatic full-model coverage enforcement when no planned coverage metadata
  exists.
- Replacement or superseding semantics for correcting written records.
