---
type: Functional Specification
title: Evaluation records
description: The deterministic on-disk contract for QUALITY.md evaluation run records.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-18T00:00:00Z
---

# Evaluation records

This spec defines the runtime record contract for a `QUALITY.md` evaluation run:
folder names, record names, record schemas, `schemaVersion`, and the division of
responsibility between the deterministic `qualitymd` CLI and the judging skill.
The evaluation semantics are defined by
[`SPECIFICATION.md` → Evaluation](../SPECIFICATION.md#evaluation).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Responsibility

The CLI **MUST** own file creation, serialization, run-folder numbering,
record numbering, `schemaVersion` stamping, and `report.md` / `report.json`
rendering. The skill supplies judgment content: findings, ratings, rationales,
roll-up judgment, and recommendations.

The skill **MUST NOT** hand-author, number, serialize, or stamp evaluation
records when the corresponding CLI command exists.

## Runtime, Not OKF

Evaluation records are raw runtime outputs, not OKF concepts. A run folder
**MUST NOT** be treated as an OKF bundle and **MUST NOT** contain OKF
`index.md`, `log.md`, or `schema.md` semantics. Runtime Markdown frontmatter in
recommendation records is machine metadata, not OKF concept frontmatter.

## Run Folder

Each run folder is named:

```text
NNNN-<altitude>[-<narrowing>]-quality-eval
```

`<altitude>` is `subject` or `model`. `NNNN` is one zero-padded sequence shared
across all altitudes in the evaluation directory.

A run folder contains:

```text
model.md
design.md
plan.md
assessments/
  NNN-<target>-<requirement>.json
analysis/
  <target>.json
recommendations/
  NNN-<slug>.md
report.md
report.json
```

`assessments/` and `recommendations/` each use their own local `NNN` sequence.

## Schema Version

Every JSON record (`assessments/*.json`, `analysis/*.json`, `report.json`)
**MUST** carry top-level `schemaVersion: 1`.

Every CLI-written recommendation Markdown record **MUST** carry runtime YAML
frontmatter with `schemaVersion: 1`.

## Assessment Record

An assessment record is one JSON file per assessed requirement. Required fields:

- `schemaVersion`
- `target`
- `targetPath`
- `requirement`
- `factors`
- `rating`, or `null`
- `notAssessed`
- `criterionSource`
- `findings`
- `rationale`
- `recommendations`

When `notAssessed` is `true`, `rating` **MUST** be `null`. Each finding **MUST**
carry `locator`, `observation`, and `category`; it **MAY** carry `severity`,
`evidence`, and `attributes`.

## Analysis Record

An analysis record is one JSON file per target. Required fields:

- `schemaVersion`, `target`, `targetPath`
- `localRating`, or `null` for a grouping target with no own requirements
- `factorRatings`
- `aggregateRating`
- `assessmentRecords`
- `childAnalysisRecords`

Every rating result **MUST** record `notAssessed` distinctly from a rating level.

## Recommendation Record

A recommendation record is one Markdown file per key gap. Its runtime
frontmatter **MUST** carry:

- `schemaVersion`
- `title`
- `gap`
- `evidenceLocators`
- `assessmentRecords`
- `remediationOptions`
- `recommendedOption`
- `doneCriterion`

The Markdown body **MUST** state the gap, evidence locators, remediation options,
recommended option, and done criterion in stable human-readable sections.

## report.json

`report.json` is the machine rendering of the same Evaluation Report as
`report.md`. It **MUST** present the in-scope root rating and rationale, scope,
per-target results, and advice. It **MUST** reference findings by assessment
record; full finding detail stays in `assessments/*.json`.
