---
type: Functional Specification
title: qualitymd evaluation status
description: Inspect whether an evaluation run can be rendered into reports.
tags: [cli, command, evaluation]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd evaluation status

`qualitymd evaluation status <run>` reads an evaluation run and reports whether
it is complete enough for `qualitymd evaluation report build`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Requirements

The command **MUST** accept either a positional run path or `--latest`, and
**MUST** error when both or neither are supplied. `--latest` resolves to the most
recent recognized run in the resolved evaluation directory.

The command **MUST NOT** write files. It exits `0` when the run can be inspected,
even when it is not yet reportable. Missing or dangling records are payload
gaps, not command failures. A missing run folder, unreadable record, or malformed
record that prevents inspection fails with the internal-error category.

Human output **MUST** include the run path, record counts, reportability, and any
gaps. Under `--json`, stdout **MUST** include `schemaVersion`, `path`,
`reportable`, counts, gaps, and `nextActions`.

A run is reportable only when exactly one analysis record represents the
in-scope root target. The root analysis record is identified by an empty
`targetPath`. If no such record exists, `status` **MUST** return
`reportable: false` with a `missing-root-analysis` gap.

A run is not reportable when two or more assessment result records cover the same
ordered `targetPath` and `requirement`, unless all but one are superseded by an
active correction record. `status` **MUST** return `reportable: false` with a
`duplicate-assessment-result` gap that references each later active duplicate
record.

Assessment and recommendation superseding references **MUST** resolve to records
in the same run. Dangling or invalid assessment superseding, stale analysis
references to superseded assessment results, and dangling recommendation
superseding **MUST** produce gaps and make the run non-reportable.

When `plan.md` contains `coverage:` frontmatter, `status` **MUST** validate it at
read time and compare planned assessment and analysis identities to written
records. Missing planned records **MUST** produce
`missing-planned-assessment-result` or `missing-planned-analysis`; written
records outside the plan **MUST** produce `unexpected-assessment-result` or
`unexpected-analysis`.

Malformed `coverage:` frontmatter **MUST** produce an `invalid-plan-coverage`
gap. A body-only `plan.md`, or frontmatter without `coverage:`, keeps the same
behavior as a run with no planned coverage metadata.
