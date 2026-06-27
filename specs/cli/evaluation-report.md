---
type: Functional Specification
title: qualitymd evaluation report
description: Build evaluation reports.
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
```

`build` **MUST** accept either a positional run path or `--latest`, and **MUST**
error when both or neither are supplied.

`--model <model>` **MUST** select the `QUALITY.md` file whose model-relative
workspace supplies `--latest` history. When `--model` is supplied with a
relative positional `<run>` path, `build` **MUST** resolve that path relative to
the selected model's workspace root. When `--model` is absent and a positional
`<run>` path is supplied, `build` **MAY** preserve ordinary filesystem-path
behavior.

For Evaluation runs, `build` validates the structured payload graph under
`data/`, assembles `data/evaluation-output-result.json`, and renders the
deterministic Markdown report tree from completed structured outputs. It renders
recorded judgment; it **MUST NOT** reread evaluated source, infer or recompute
ratings, invent findings, or choose new recommendations by evaluator judgment.
It **MUST** fail before writing generated report files when the run is not
renderable, including when required Evaluation data is missing, malformed,
schema-incompatible, or structurally incomplete. The failure **MUST** identify
the blocking gap and point the caller to `qualitymd evaluation status <run>` for
the complete gap list. It **MUST** be deterministic and idempotent: unchanged
structured data produces byte-identical report files.

On success, the build receipt's `reportMd` field **MUST** point to the
run-level `report.md`. The receipt **MAY** include `headlineReportMd` and
`rootAreaReportMd` when those subject reports exist. The receipt's
`ratingResult` **MUST** describe the headline result, not necessarily the root
Area result.

`build` **MUST NOT** accept a gate flag. Report gating is not part of Evaluation
v0.

The Evaluation report content contract is defined by
[Evaluation report tree](../evaluation/reports/report-tree.md).
