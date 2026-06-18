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

The command **MUST NOT** write files. It exits `0` when the run can be inspected,
even when it is not yet reportable. Missing or dangling records are payload
gaps, not command failures. A missing run folder, unreadable record, or malformed
record that prevents inspection fails with the internal-error category.

Human output **MUST** include the run path, record counts, reportability, and any
gaps. Under `--json`, stdout **MUST** include `schemaVersion`, `path`,
`reportable`, counts, gaps, and `nextActions`.
