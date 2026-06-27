---
type: Functional Specification
title: /quality reporting
description: Component spec for /quality Evaluation reports and run artifacts.
tags: [skill, quality, evaluation, reporting]
timestamp: 2026-06-25T00:00:00Z
---

# /quality reporting

This spec owns the `/quality` skill's Evaluation reporting and run-artifact
contract. It composes the shared contracts in the parent
[/quality skill](quality-skill.md), the judgment workflow in
[/quality evaluation workflow](evaluation.md), and the durable
[Evaluation report tree](../../evaluation/reports/report-tree.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Runtime Artifacts

The CLI creates a numbered evaluation folder per run. The default parent
directory is `.quality/evaluations/` under the workspace quality data directory
beside the selected `QUALITY.md`. A repository may choose the workspace config
file with root `config` frontmatter on `QUALITY.md`; without that pointer the
config file defaults to `.quality/config.yaml` beside the selected model.

The run folder **MUST** include `model-snapshot.md`, a snapshot of the
`QUALITY.md` as evaluated.

Evaluation structured data **MUST** live under `data/`. The skill **MUST**
persist routine outputs through `qualitymd evaluation data set`; it **MUST NOT**
hand-author structured data files directly.

The CLI-generated `data/evaluation-output-result.json` indexes completed
structured outputs and generated report paths.

Workflow-experience feedback for current runs lives in
`.quality/logs/<timestamp>-evaluate-feedback-log.md`, not in the evaluation run
folder.

Runtime evaluation artifacts are raw outputs in the evaluated repository, not
OKF concepts. They **MUST NOT** carry OKF frontmatter or require registration in
`specs/schema.md`.

## Report Generation

The skill **MUST** build reports through:

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

Report build **MUST** assemble `data/evaluation-output-result.json` before
rendering Markdown reports.

Reports **MUST** be deterministic projections over completed structured data.
The skill **MUST NOT** add report-only findings, ratings, evidence, limits,
analysis, or recommendations.

Area and Factor reports **MUST NOT** render `Findings` sections. Their roll-up
explanation belongs in `Rating Drivers`, rationale, confidence, limits, and
incomplete inputs.

Reports **MUST** preserve secret-handling boundaries: cite locator and
credential type only, never secret values or unsafe raw content.

## Report Tree

The run-level Evaluation report **MUST** be `report.md` at the run root.

The root Area report **MUST** be `root-area.md` at the run root when the root
Area has an Area Analysis Result in the run.

Non-root Area reports **MUST** be written with short subject-aware filenames
under their Area folder, such as `areas/<area>/<area>-area.md`.

Factor reports **MUST** be written with short subject-aware filenames under the
owning Area report folder, such as `factors/<factor>/<factor>-factor.md`.

Requirement reports **MUST** be written with short subject-aware filenames under
the owning Area report folder, such as
`requirements/<requirement>/<requirement>-requirement.md`.

Every report **MUST** include the navigation trails required by the
[Evaluation report tree](../../evaluation/reports/report-tree.md).

Area reports **MUST** link to local root Factor reports, local Requirement
reports, and direct Child Area reports.

Factor reports **MUST** link to their owning Area report when that Area report
was generated, parent Factor report when present, Sub-Factor reports, and
direct Requirement reports.

Requirement reports **MUST** link to their owning Area report and every attached
Factor report.

## User-Facing Closeout

The agent's user-facing evaluation closeout is governed by the shared
[user interaction contract](quality-skill.md#user-interaction-contract).

The closeout **MUST** state the rating, scope, evidence basis, recommendations,
known limits or incomplete inputs, changed artifacts, what was not done, and the
recommended next action.
The closeout **MUST** use labeled fields for rating, scope, evidence basis,
recommendations, known limitations, changed artifacts, not-done boundary,
reports, and next action.

The closeout **MUST** distinguish evaluated-source quality from model weakness,
missing evidence, unknowns, and evaluation limits.

The closeout **MUST** name the top recommendation or the recommendation index
when no single recommendation should be singled out. Recommendation follow-up
remains a separate workflow for applying or handing off recommendation
artifacts.

Finding-local `candidateActions` are **not** recommendations: reports and the
closeout **MUST NOT** present them as selected next moves. Selected advice
belongs in `RecommendationResult` and generated recommendation reports.
