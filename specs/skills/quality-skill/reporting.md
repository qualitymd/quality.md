---
type: Functional Specification
title: /quality reporting
description: Component spec for /quality Evaluation v2 reports and run artifacts.
tags: [skill, quality, evaluation, reporting]
timestamp: 2026-06-25T00:00:00Z
---

# /quality reporting

This spec owns the `/quality` skill's Evaluation v2 reporting and run-artifact
contract. It composes the shared contracts in the parent
[/quality skill](quality-skill.md), the judgment workflow in
[/quality evaluation workflow](evaluation.md), and the durable
[Evaluation v2 report tree](../../evaluation-v2/reports/report-tree.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Runtime Artifacts

The CLI creates a numbered evaluation folder per run. The default parent
directory is `.quality/evaluations/` under the workspace quality data directory.
A repository may choose the workspace config file with root `config` frontmatter
on `QUALITY.md`; without that pointer the config file defaults to
`.quality/config.yaml`.

The run folder **MUST** include `model.md`, a snapshot of the `QUALITY.md` as
evaluated.

Evaluation v2 structured data **MUST** live under `data/`. The skill **MUST**
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

Report build **MUST** assemble `data/evaluation-output-result.json` before
rendering Markdown reports.

Reports **MUST** be deterministic projections over completed structured data.
The skill **MUST NOT** add report-only findings, ratings, evidence, limits,
analysis, or recommendations.

Reports **MUST** preserve secret-handling boundaries: cite locator and
credential type only, never secret values or unsafe raw content.

## Report Tree

The root Area report **MUST** be `report.md` at the run root.

Non-root Area reports **MUST** be written under `areas/**/report.md`.

Factor reports **MUST** be written under the owning Area report folder.

Requirement reports **MUST** be written under the owning Area report folder.

Every report **MUST** include linked breadcrumbs. Every non-root report **MUST**
include a parent link.

Area reports **MUST** link to local root Factor reports, local Requirement
reports, and direct child Area reports.

Factor reports **MUST** link to their owning Area report, parent Factor report
when present, child Factor reports, and direct Requirement reports.

Requirement reports **MUST** link to their owning Area report and every attached
Factor report.

## User-Facing Closeout

The agent's user-facing evaluation closeout is governed by the shared
[user interaction contract](quality-skill.md#user-interaction-contract).

The closeout **MUST** state the rating, scope, evidence basis, known limits or
incomplete inputs, changed artifacts, what was not done, and the recommended
next action.

The closeout **MUST** distinguish evaluated-source quality from model weakness,
missing evidence, unknowns, and evaluation limits.

Evaluation v2 v0 **MUST NOT** present generated recommendations as part of the
reporting closeout. Recommendation follow-up remains a separate workflow for
historical recommendation artifacts.
