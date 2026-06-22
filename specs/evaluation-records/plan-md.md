---
type: Functional Specification
title: plan.md
description: Runtime plan.md and planned coverage contract for evaluation runs.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-22T00:00:00Z
---

# plan.md

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It owns the independently reviewable runtime contract below; shared
responsibility, runtime-not-OKF status, schema-version, and historical-record
compatibility rules live in the parent spec.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Planned Coverage

`plan.md` is a YAML-frontmatter + Markdown-body artifact. Its Markdown body is
the run's prose plan. Optional frontmatter `coverage:` lists the assessment
result requirements and analysis areas intended for the run. Planned coverage
exists to support deterministic resume diagnostics; it **MUST NOT** replace the
prose plan or the evaluation design.

When `coverage:` is present, it **MUST** contain:

- `assessmentResults`
- `analyses`

Each assessment result entry **MUST** contain ordered `areaPath` and
`requirement`.
Each analysis entry **MUST** contain ordered `areaPath`.

Coverage frontmatter is hand-authored as part of `plan.md`; there is no separate
CLI write command for planned coverage.

New run scaffolds **SHOULD** include a commented or fenced coverage example in
`plan.md` so the expected shape is discoverable without guessing field names.

A planned assessment result key is ordered `areaPath` plus `requirement`. A
planned analysis key is ordered `areaPath`. Duplicate planned assessment
result or analysis keys are invalid.

When `coverage:` is absent, the run keeps the same status and reportability
behavior it would have without planned coverage metadata. Malformed `coverage:`
frontmatter makes the run non-reportable through an `invalid-plan-coverage`
status gap rather than making the run unloadable. The gap detail should
aggregate detectable coverage-shape problems and name fields by YAML/JSON key
path.
