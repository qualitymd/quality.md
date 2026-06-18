---
type: Functional Specification
title: Evaluation record format — functional spec
description: The deterministic contract for an evaluation run's on-disk records.
tags: [evaluation, specs, cli]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation record format — functional spec

Companion to [Evaluation record format](../0012-evaluation-record-format.md).
This spec states *what* the evaluation-record contract is; it does not specify
the CLI commands that produce the records (deferred — see
[Scope](#scope)).

The evaluation semantics — Define, Assess and Rate, Analyze, Advise, Report — are
the source of truth in [`SPECIFICATION.md`](../../../SPECIFICATION.md#evaluation);
this spec governs only how a run's outputs are laid out and serialized.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be interpreted as
described in IETF RFC 2119.

## Scope

Covered: the on-disk run-folder layout and naming, the schema of each record and
its required fields, the `schemaVersion` convention, and the division of
responsibility between the deterministic CLI and the judging skill.

Deferred: the CLI command surface that writes these records (changes
0013–0015); any record kinds beyond those defined here.

## Responsibility

The deterministic `qualitymd` CLI **MUST** own writing every record: file
creation, serialization, run-folder numbering, `schemaVersion` stamping, and
rendering `report.md`/`report.json`. The skill supplies the **judgment content**
(findings, ratings, rationales, recommendations); it **MUST NOT** itself
serialize, number, or stamp records.

## Artifacts are not OKF

A run's records are raw runtime outputs, **not** OKF concepts. The run folder
**MUST NOT** be treated as an OKF bundle (no `index.md`/`log.md`/`schema.md`
semantics), and runtime Markdown frontmatter **MUST NOT** be interpreted as OKF
concept frontmatter.

## Run folder

Each evaluation produces one run folder named
`NNNN-<altitude>[-<narrowing>]-quality-eval`, where `<altitude>` is `subject` or
`model` and `<narrowing>` is an optional scope slug.

- `NNNN` **MUST** be a single zero-padded sequence shared across all altitudes
  within the evaluation directory: the next number is one past the highest
  existing run regardless of its altitude. A `subject` run **MUST NOT** reuse a
  number already held by a `model` run, or vice versa.

A run folder **MUST** contain:

```text
NNNN-<altitude>[-<narrowing>]-quality-eval/
  model.md          # the model evaluated (Define)
  design.md         # how the evaluation was scoped and run
  plan.md           # the assessment plan
  assessments/
    NNN-<target>-<requirement>.json
  analysis/
    <target>.json
  recommendations/
    NNN-<slug>.md
  report.md         # human-readable Evaluation Report
  report.json       # machine-readable Evaluation Report
```

`assessments/` and `recommendations/` entries are numbered with independent
zero-padded `NNN` sequences local to their respective subdirectories within the
run.

## schemaVersion

Every JSON record (`assessments/*.json`, `analysis/*.json`, `report.json`)
**MUST** carry a top-level `schemaVersion` integer. Its value is `1` for the
contract defined here.

Every machine-written Markdown record (`recommendations/*.md`) **MUST** carry a
YAML frontmatter block with `schemaVersion: 1`. This frontmatter is runtime
metadata for deterministic CLI reads and writes; it does not make the file an OKF
concept.

## Assessment record

One record per assessed requirement, at
`assessments/NNN-<target>-<requirement>.json`. Each record holds these fields;
those marked required **MUST** be present:

- `schemaVersion` (required).
- `target` (required) — the target's name.
- `targetPath` (required) — the target's path from the root, as an array.
- `requirement` (required) — the requirement assessed.
- `factors` (required) — the factors the requirement is tied to (MAY be empty).
- `rating` (required) — the requirement's rating level, or `null`.
- `notAssessed` (required) — boolean; `true` when the requirement could not be
  rated, in which case `rating` **MUST** be `null`.
- `criterionSource` (required) — what criterion the rating was made against
  (e.g. the rating scale or a requirement override).
- `findings` (required) — an array of findings (MAY be empty).
- `rationale` (required) — why the rating (or `notAssessed`) was reached.
- `recommendations` (required) — references to `recommendations/*.md` files
  (MAY be empty).

Each finding **MUST** carry `locator`, `observation`, and `category`, and **MAY**
carry `severity`, `evidence`, and `attributes`. Domain-specific metadata (e.g. a
credential type) **MUST** be placed under `attributes`, not as a new top-level
finding field.

## Analysis record

One record per target, at `analysis/<target>.json`, recording the roll-up for
that target. Each record **MUST** carry:

- `schemaVersion`, `target`, `targetPath`.
- `localRating` — the target's local rating with its `rationale`, or `null` for
  a grouping target with no own requirements.
- `factorRatings` — each factor's rating with its `rationale` (MAY be empty).
- `aggregateRating` — the target's aggregate rating with its `rationale`.
- `assessmentRecords` — references to the contributing `assessments/*.json`.
- `childAnalysisRecords` — references to child targets' `analysis/*.json`.

Each rating **MUST** record `notAssessed` distinctly from a rating level.

## report.json

A machine-readable rendering of the same Evaluation Report as `report.md`.
`report.json` **MUST** present the in-scope root rating and rationale, the scope,
and the per-target results, and **MUST** reference findings by record (the full
finding detail stays in `assessments/*.json`). It **SHOULD** carry only minimal
finding summaries inline.

## Recommendation record

One Markdown file per key gap, at `recommendations/NNN-<slug>.md`. Recommendation
records are human-readable Markdown with YAML frontmatter so a person can review
the advice directly and the CLI can read it mechanically.

The frontmatter **MUST** carry:

- `schemaVersion`.
- `title` — a short human label for the recommendation.
- `gap` — the gap the recommendation addresses.
- `evidenceLocators` — the evidence positions supporting the gap.
- `assessmentRecords` — references to related `assessments/*.json` records
  (MAY be empty).
- `remediationOptions` — the options considered for remediation.
- `recommendedOption` — the chosen remediation option.
- `doneCriterion` — how completion will be recognized.

The Markdown body **MUST** state the gap, the evidence locators, the remediation
options, one recommended option, and the done criterion in stable, human-readable
sections. For a `notAssessed` gap, the done criterion is to become assessable and
reach at least the acceptable floor.
