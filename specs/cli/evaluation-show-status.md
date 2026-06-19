---
type: Functional Specification
title: qualitymd evaluation show-status
description: Inspect whether an evaluation run can be rendered into reports.
tags: [cli, command, evaluation]
timestamp: 2026-06-18T00:00:00Z
---

# qualitymd evaluation show-status

`qualitymd evaluation show-status <run>` reads an evaluation run and reports
whether it is complete enough for `qualitymd evaluation build-report`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

The command **MUST NOT** write files. It exits `0` when the run can be inspected,
even when it is not yet reportable. Missing or dangling records are payload
gaps, not command failures. A missing run folder, unreadable record, or malformed
record that prevents inspection fails with the internal-error category.

Human output **MUST** include the run path, record counts, reportability, and any
gaps. Under `--json`, stdout **MUST** include `schemaVersion`, `path`,
`reportable`, counts, gaps, and `nextActions`.

A run is reportable only when exactly one analysis record represents the
in-scope root target. The root analysis record is identified by an empty
`targetPath`. If no such record exists, `show-status` **MUST** return
`reportable: false` with a `missing-root-analysis` gap.

A run is not reportable when two or more assessment records cover the same
ordered `targetPath` and `requirement`, unless all but one are superseded by an
active correction record. `show-status` **MUST** return `reportable: false` with
a `duplicate-assessment` gap that references each later active duplicate record.

When an assessment record supersedes another assessment, every superseding
reference **MUST** resolve to an existing earlier assessment ID or path in the
same run. A dangling reference **MUST** produce a
`missing-superseded-assessment` gap and make the run non-reportable.

An assessment **MUST NOT** supersede a record for a different ordered
`targetPath` plus `requirement`. Such a mismatch **MUST** produce an
`invalid-assessment-supersedes` gap and make the run non-reportable.

An analysis record **MUST NOT** reference a superseded assessment. Such a
reference **MUST** produce a `superseded-assessment-reference` gap and make the
run non-reportable.

When a run has `planned-coverage.json`, `show-status` **MUST** compare planned
assessment entries to written assessment records by ordered `targetPath` plus
`requirement`. Missing planned assessment records **MUST** produce
`missing-planned-assessment` gaps.

When a run has `planned-coverage.json`, `show-status` **MUST** compare planned
analysis entries to written analysis records by ordered `targetPath`. Missing
planned analysis records **MUST** produce `missing-planned-analysis` gaps.

When a run has `planned-coverage.json`, assessment or analysis records outside
the plan **MUST** produce `unexpected-assessment` or `unexpected-analysis` gaps.

Planned-coverage gaps **MUST** make the run non-reportable. When
`planned-coverage.json` is absent, `show-status` **MUST** keep the same behavior
as a run with no planned coverage metadata.

When a non-reportable run has duplicate-record or unexpected-record gaps,
`nextActions` should direct the user to review the gaps rather than adding
more records.

When a recommendation record supersedes another recommendation, every
superseding reference **MUST** resolve to an existing recommendation ID or path
in the same run. A dangling reference **MUST** produce a
`missing-superseded-recommendation` gap and make the run non-reportable.
