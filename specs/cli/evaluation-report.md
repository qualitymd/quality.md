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

`build` derives one assembled report model from the run's model snapshot, plan,
assessment records, analysis records, recommendation records, and run metadata,
then renders `report-summary.md`, `report.md`, and `report.json` from that model.
It renders recorded judgment; it **MUST NOT** reread subject source, infer or
recompute ratings, invent findings, or choose new recommendations by evaluator
judgment. It **MUST** fail without writing a partial report when the run is not
renderable, including when any run record is malformed, schema-incompatible, or
structurally incomplete under the current evaluation-record contract. The
failure **MUST** identify the blocking gap and point the caller to
`qualitymd evaluation status <run>` for the complete gap list. It **MUST** be
deterministic and idempotent: unchanged records produce byte-identical report
files.

The assembled report model **MUST** preserve typed states from the record
contract in `report.json`: rating-result kind, local-rating state, record
lifecycle state, next-step kind, missing metadata fields, finding severity, and
stable target/factor paths. Human Markdown can render those states with display
labels, but it **MUST NOT** collapse structural local ratings into not-assessed
ratings or superseded records into active advice.

`build` **MUST NOT** accept a gate flag. Gating is a separate operation.

`gate` **MUST** read the already-rendered `report.json`, compare the in-scope
root aggregate verdict to `--at-or-below <level>`, and exit `1` when that
verdict is equal to or worse than `<level>`. It exits `0` when better and exits
`2` when `<level>` is not in the run's rating scale. A not-assessed verdict
fails the gate. `gate` **MUST NOT** write or modify any run file and **MUST**
fail when no rendered `report.json` exists. Before trusting an existing
`report.json`, `gate` **MUST** inspect the selected run and fail with the same
non-reportable-run diagnostic as `build` when current records are incompatible.

The report content contract is defined by
[Evaluation records](../evaluation-records.md#reportjson) and
[Evaluation records](../evaluation-records.md#report-summarymd).
