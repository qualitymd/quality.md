---
type: Functional Specification
title: qualitymd evaluation report
description: Build and gate evaluation reports.
tags: [cli, command, evaluation]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd evaluation report

`qualitymd evaluation report` is the report resource for a run.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Commands

```text
qualitymd evaluation report build <run>
qualitymd evaluation report gate <run> --at-or-below <level>
```

Both verbs **MUST** accept either a positional run path or `--latest`, and
**MUST** error when both or neither are supplied.

`build` derives `report-summary.md`, `report.md`, and `report.json` from the
run's assessment, analysis, and recommendation records. It renders recorded
judgment; it **MUST NOT** infer or recompute ratings. It **MUST** fail without
writing a partial report when the run is not renderable. It **MUST** be
deterministic and idempotent: unchanged records produce byte-identical report
files.

`build` **MUST NOT** accept a gate flag. Gating is a separate operation.

`gate` **MUST** read the already-rendered `report.json`, compare the overall
rating to `--at-or-below <level>`, and exit `1` when the in-scope root aggregate
rating is equal to or worse than `<level>`. It exits `0` when better and exits
`2` when `<level>` is not in the run's rating scale. An overall not-assessed
result fails the gate. `gate` **MUST NOT** write or modify any run file and
**MUST** fail when no rendered `report.json` exists.

The report content contract is defined by
[Evaluation records](../evaluation-records.md#reportjson) and
[Evaluation records](../evaluation-records.md#report-summarymd).
