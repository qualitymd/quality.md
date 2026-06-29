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
deterministic Markdown report tree from completed structured outputs, including
the generated report frontmatter, header, and source-data section contract defined by
[Evaluation report tree](../evaluation/reports/report-tree.md). It renders
recorded judgment; it **MUST NOT** reread evaluated source, infer or recompute
ratings, invent findings, choose new recommendations by evaluator judgment, or
read generated report frontmatter or Markdown body content as source data. It
**MUST** fail before writing generated report files when the run is not
renderable, including when required Evaluation data is missing, malformed,
schema-incompatible, or structurally incomplete. The failure **MUST** identify
the blocking gap and point the caller to `qualitymd evaluation status <run>` for
the complete gap list. It **MUST** be deterministic and idempotent: unchanged
structured data produces byte-identical report files.

`build` **MUST** read the run scope from `RunManifest.plannedScope`.
`report.md` **MUST** render as the scoped Area report for
`plannedScope.areaId`, narrowed by `factorFilter` when present. It **MUST NOT**
choose a headline subject from `EvaluationFrame` or any other agent-authored
payload ordering. `report.md` **MUST** title the run entrypoint as
`Quality Evaluation - <Area title>` and append a parenthesized comma-separated
Factor title list when `plannedScope.factorFilter` is present. `report.md`
**MUST** render non-judgmental run metadata in YAML frontmatter, omit the
visible top run-context line, and open with `Summary`, `Key Details`, and
`Contents` sections before Top Findings. `report.md` **MUST NOT** render a
visible `Limits & Incomplete Inputs` section.

`build` **MUST** render persisted Advice outputs into `report.md`,
`findings.md`, `recommendations.md`, and recommendation detail reports under
`recommendations/`. `report.md` **MUST** include Top Findings and Top
Recommendations sections capped at 10 rows each and **MUST** link to
`findings.md` and `recommendations.md`. Recommendation report content **MUST** be
rendered from persisted Advice data and the model snapshot, not from YAML
frontmatter or Markdown body content in other generated reports.

Generated Markdown report frontmatter **MUST** contain only non-judgmental
report metadata allowed by the report-tree contract, and its `title` **MUST**
match the report's visible H1 title text. Every generated Markdown report
**MUST** end with a `Source Data` section listing the structured Evaluation
payload files used as source data for the specific report artifact. The
generated `data/evaluation-output-result.json` index **MUST NOT** be listed as
report source data unless a report is directly rendered from it.

On success, the build receipt's `reportMd` field **MUST** point to `report.md`.
The receipt's `ratingResult` **MUST** describe the scoped Area result rendered by
`report.md`.

`build` **MUST NOT** accept a gate flag.

The Evaluation report content contract is defined by
[Evaluation report tree](../evaluation/reports/report-tree.md).
