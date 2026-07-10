---
type: Functional Specification
title: /quality reporting
description: Component spec for /quality evaluation reports and run artifacts.
tags: [skill, quality, evaluation, reporting]
timestamp: 2026-07-09T00:00:00Z
---

# /quality reporting

This spec owns the `/quality` skill's evaluation reporting and run-artifact
contract. It composes the shared contracts in the parent
[/quality skill](quality-skill.md), the judgment workflow in
[/quality evaluation workflow](evaluation.md), and the durable
[Evaluation report tree](../../evaluation/reports/report-tree.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Runtime artifacts

The CLI creates a numbered evaluation folder per run. The default parent
directory is `.quality/evaluations/` under the workspace quality data directory
beside the selected `QUALITY.md`. A repository may choose the workspace config
file with root `config` frontmatter on `QUALITY.md`; without that pointer the
config file defaults to `.quality/config.yaml` beside the selected model.

The run folder **MUST** include `model-snapshot.md`, a snapshot of the
`QUALITY.md` as evaluated.

For new runner-created runs, structured evaluation data lives in the
authoritative `evaluation.json` at the run root, written atomically by
`qualitymd evaluation run` together with run-local logs under `logs/`.
Historical runs keep their multi-file structured data under `data/`, persisted
through `qualitymd evaluation data set`. The skill **MUST NOT** hand-author
structured data files directly on either path.

For historical runs, the CLI-generated `data/evaluation-output-result.json`
indexes completed structured outputs and generated report paths; runner-created
runs carry that output index inside `evaluation.json`.

Workflow-experience feedback for current runs lives in
`.quality/logs/<timestamp>-evaluate-feedback-log.md`, not in the evaluation run
folder.

Runtime evaluation artifacts are raw outputs in the evaluated repository.
Generated Markdown reports carry identity frontmatter and visible bottom
`Primary source data` sections as defined by the
[Evaluation report tree](../../evaluation/reports/report-tree.md). Report
frontmatter `title` matches the visible H1 document title; `type` carries the
report artifact taxonomy. The `Primary source data` section lists report-local
primary structured evaluation payloads used to render the specific report
artifact. Generated reports use standard `Contents` sections and
`Evaluation links:` navigation instead of compact `Jump to:` lines or local
`Legend` blocks or bottom `Legend` sections. The evaluation run folder is not
yet a full OKF bundle: it does not require
generated `index.md`, `schema.md`, or `log.md` files, and generated reports do
not require registration in `specs/schema.md`.

## Report generation

`qualitymd evaluation run` builds the Markdown reports as part of the run. When
reports must be rebuilt or a historical run rendered, the skill **MUST** build
them through:

```text
qualitymd evaluation report build <run>
```

The skill **MUST** run `qualitymd evaluation status <run>` before report build
when it needs to diagnose missing or invalid structured data.

When the selected model is not the default `QUALITY.md` in the current working
directory and the run path is model-relative, the skill **MUST** pass
`--model <model>` to `qualitymd evaluation status`, `qualitymd evaluation data`,
and `qualitymd evaluation report build` so run resolution uses the selected
model's workspace root.

Reports **MUST** be rendered from the run's authoritative structured data:
`evaluation.json` for runner-created runs, or the assembled
`data/evaluation-output-result.json` for historical multi-file runs.

Reports **MUST** be deterministic projections over completed structured data.
The skill **MUST NOT** add report-only findings, ratings, evidence, limits,
analysis, or recommendations.

Reports **MUST NOT** read generated report frontmatter or Markdown body content
from other generated reports as source data.

Area and factor reports **MUST NOT** render `Findings` sections or standalone
`Rating Drivers` sections. Their human-facing roll-up explanation belongs in
summary, ratings, confidence, limits, incomplete inputs, and breakdown tables,
while structured `ratingDrivers` remain available through the payloads listed in
the report's `Primary source data` section.

Run reports **MUST** render the Model evaluation section required by the
[Evaluation report tree](../../evaluation/reports/report-tree.md). Area reports
**MUST** render the Area / factor breakdown required by that report tree. These
sections are the human-facing area / factor structure and status surface; the
machine-readable generated-report manifest remains `evaluation.json` for
runner-created runs and `data/evaluation-output-result.json` for historical
runs.

Reports **MUST** preserve secret-handling boundaries: cite locator and
credential type only, never secret values or unsafe raw content.

## Report tree

The run-level evaluation report **MUST** be `report.md` at the run root.

The root area report **MUST** be `root-area.md` at the run root when the root
area has an area analysis result in the run.

Non-root area reports **MUST** be written with short subject-aware filenames
under their area folder, such as `areas/<area>/<area>-area.md`.

Factor reports **MUST** be written with short subject-aware filenames under the
owning area report folder, such as `factors/<factor>/<factor>-factor.md`.

Requirement reports **MUST** be written with short subject-aware filenames under
the owning area report folder, such as
`requirements/<requirement>/<requirement>-requirement.md`.

Every report **MUST** include the navigation trails required by the
[Evaluation report tree](../../evaluation/reports/report-tree.md).

Every generated Markdown report **MUST** include the frontmatter, run context,
`Evaluation links:` navigation, and report-specific header summary required for
that report kind by the
[Evaluation report tree](../../evaluation/reports/report-tree.md). Every
generated Markdown report **MUST** include the `Primary source data` section
required by that report tree. Report bodies **MUST NOT** duplicate report-level
source-data links in header `Data` columns; the `Primary source data` section
owns those pointers.

Area reports **MUST** link to local and descendant area and factor reports
through the Area / factor breakdown, and to local requirement reports through
their requirement table.

Factor reports **MUST** link to their owning area report when that area report
was generated, parent factor report when present, sub-factor reports, and
direct requirement reports.

Requirement reports **MUST** link to their owning area report and every attached
factor report.

## User-facing closeout

The agent's user-facing evaluation closeout is governed by the shared
[user interaction contract](quality-skill.md#user-interaction-contract).

The closeout **MUST** state the rating, scope, evidence basis, recommendations,
known limits or incomplete inputs, changed artifacts, what was not done, and the
recommended next action.
The closeout **MUST** use labeled fields for rating, scope, evidence basis,
recommendations, known limitations, changed artifacts, not-done boundary,
report-reading CTA, and next action.

The closeout **MUST** make `report.md` the primary human-report CTA by naming
the completed run's full report path and describing its value as the
decision-ready evaluation result with rating, evidence basis, limits, top
findings, and top recommendations. The closeout **MUST** name the completed
run's full `recommendations.md` path and describe its value as the
action-planning report with ranked recommendations, why they matter, expected
benefit, and how to know each worked. The primary report-reading CTA **MUST NOT**
include `evaluation.json`, `data/evaluation-output-result.json`, or other
machine-oriented run artifacts.

The closeout **MUST** distinguish evaluated-source quality from model weakness,
missing evidence, unknowns, and evaluation limits.

The closeout **MUST** name the top recommendation or direct the user to
`recommendations.md` when no single recommendation should be singled out.
Recommendation follow-up remains a separate workflow for applying or handing off
recommendation artifacts.

Finding-local `candidateActions` are **not** recommendations: reports and the
closeout **MUST NOT** present them as selected next moves. Selected advice
belongs in `RecommendationResult` and generated recommendation reports.
