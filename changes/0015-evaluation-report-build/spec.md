---
type: Functional Specification
title: Evaluation status and report build — functional spec
description: qualitymd evaluation show-status inspects a run and qualitymd evaluation build-report derives report.md and report.json from its records.
tags: [evaluation, cli, specs, skill]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation status and report build — functional spec

Companion to [Evaluation status and report build](../0015-evaluation-report-build.md). This
spec states *what* `qualitymd evaluation show-status` and
`qualitymd evaluation build-report` must do; it specifies behavior, not
implementation.

Both commands inherit the cross-cutting CLI contract — non-interactive
operation, stdout/stderr split, exit-code categories, and `--json` conventions —
from the [CLI spec](../../specs/cli.md). This file specifies only what is
particular to these commands.

Two source-of-truth boundaries bind this spec:

- The **Evaluation Report** — what a report must present — is defined in
  [`SPECIFICATION.md` → Report](../../SPECIFICATION.md#report). `build-report`
  *renders* that report; it does not redefine it.
- The **record contract** — the run-folder layout and the schema of
  `assessments/*.json`, `analysis/*.json`, `recommendations/*.md`, and
  `report.json` — is defined in the
  [evaluation-record format spec](../0012-evaluation-record-format/spec.md).
  These commands read and write those records; they do not redefine them.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: showing a run's renderability status, deriving `report.md` and
`report.json` from a run's records, idempotency, the `--fail-at-or-below` CI gate
and its exit codes, and behavior when records are missing or incomplete.

Deferred: scaffolding a run ([0013](../0013-evaluation-run-scaffold.md),
`qualitymd evaluation create-run`) and writing records into it
([0014](../0014-evaluation-record-write.md), `qualitymd evaluation add-record`).
These commands neither create run folders nor write assessment, analysis, or
recommendation records.

## Rendering, not judgment

The roll-up **ratings** at every level — requirement, factor, local, aggregate,
and the in-scope root — are skill judgment, already recorded in the run's
`assessments/*.json` and `analysis/*.json` records (written by
[0014](../0014-evaluation-record-write.md)). `build-report` **renders** those
recorded ratings into the report; it **MUST NOT** infer, recompute, or alter any
rating, rationale, or *not assessed* outcome. Where the records and a derived
report disagree, the records are authoritative.

## Status command

`qualitymd evaluation show-status <run>` takes one positional argument, `<run>` —
the run folder to inspect. The command **MUST** read the run's
`assessments/*.json`, `analysis/*.json`, and `recommendations/*.md` records and
report whether the run is complete enough for `build-report`.

The human output **MUST** include the run path, counts of assessment, analysis,
and recommendation records, whether the run is reportable, and any missing or
inconsistent records that would prevent report rendering. Under `--json`, the
command **MUST** emit a stable JSON document carrying at least `schemaVersion`,
`path`, `reportable`, record counts, missing or inconsistent record references,
and `nextActions`.

`show-status` **MUST NOT** write any files. It exits `0` when the run can be
inspected, even when `reportable` is `false`; missing, dangling, or incomplete
record references are the status payload, not a command failure. For example, an
assessment that references a missing `recommendations/*.md` file, or one whose
frontmatter is not parseable enough to identify it as a recommendation record,
is reported as a renderability gap. A missing or non-run target, an unreadable
record file, or a malformed record that prevents the run from being inspected
**MUST** fail with the internal-error category.

## Build-report command

`qualitymd evaluation build-report <run>` takes one positional argument, `<run>`
— the run folder to render. The command:

- **MUST** read the run's `assessments/*.json`, `analysis/*.json`, and
  `recommendations/*.md` records.
- **MUST** derive `report.md` (the human-readable Evaluation Report) and
  `report.json` (the machine-readable rendering) and write both into the run
  folder, replacing any existing report files.
- **MUST NOT** require the skill to hand-author either report file; deriving them
  is this command's job, and the skill **MUST NOT** write `report.md` or
  `report.json` itself.

## report.md

`report.md` **MUST** present every element the
[Report phase](../../SPECIFICATION.md#report) requires, drawn from the records:

- the in-scope **root aggregate rating** and its **rationale**;
- the **scope** the rating was produced under;
- for each in-scope target (root first, recursively): each requirement's findings
  summary, rating, and rationale; each factor's rating and rationale, including
  every sub-factor at every depth; and the target's **local** and **aggregate**
  ratings, each with rationale; and
- the **advice** — key gaps, options, and recommendations.

*Not assessed* outcomes **MUST** be shown wherever they occur, distinct from
rated outcomes, at every level of the report.
[Appendix A](../../SPECIFICATION.md#appendix-a-sample-evaluation-report) of the
format spec is the non-normative reference rendering.

## report.json

`report.json` **MUST** be the machine rendering of the same Evaluation Report as
`report.md`, conforming to the
[record contract](../0012-evaluation-record-format/spec.md#reportjson): it
presents the in-scope root rating and rationale, the scope, and the per-target
results; it **MUST** reference findings by record rather than restating them, and
**SHOULD** carry only minimal finding summaries inline. The full finding detail
stays in `assessments/*.json`.

Recommendation advice **MUST** be rendered in deterministic filename order by
`recommendations/NNN-<slug>.md`, unless a later contract adds an explicit
ordering field.

## Idempotency

The command **MUST** be deterministic and idempotent: rendering a run whose
records are unchanged **MUST** reproduce byte-for-byte identical `report.md` and
`report.json`, and re-running over its own output **MUST NOT** change either file.
Output carries no timestamps, ordering jitter, or other run-to-run variation.

## Missing or incomplete records

A run is renderable only when its records are complete enough to present the
report the [Report phase](../../SPECIFICATION.md#report) requires:

- When the run folder does not exist, or its record graph is not renderable — for
  example the `analysis/` roll-up is absent or incomplete, no analysis record
  has an empty `targetPath` for the in-scope root, a target referenced by another
  record has no `analysis/<target>.json`, an analysis record references an
  `assessments/*.json` that is not present, or an assessment references a
  `recommendations/*.md` record that is not present or lacks parseable runtime
  frontmatter — the command **MUST** fail with an internal-error exit (`70`),
  name the missing or inconsistent record on stderr, and **MUST NOT** write a
  partial report. The same missing or inconsistent records are reported by
  `show-status` with exit `0` and `reportable: false`.
- A *not assessed* outcome recorded in the records is **not** a missing record:
  it is valid, expected content and **MUST** be rendered as such (see
  [report.md](#reportmd)), never treated as an error.

## The `--fail-at-or-below` gate

`--fail-at-or-below <level>` turns `build-report` into a CI gate over the
rendered result. `<level>` names a rating level of the evaluated model's rating
scale (see the [format spec](../../SPECIFICATION.md#rating-scale)); the scale is
ordered best to worst, and the gate compares against that order.

- The gate is driven by the **in-scope root aggregate rating** — the same rating
  `report.md` presents as the Rating.
- The command **MUST** exit non-zero with the *ran-but-found-problems* code (`1`)
  when that rating is **at or below** `<level>` in the scale's order (i.e. equal
  to `<level>` or worse). It **MUST** exit `0` when the rating is strictly better
  than `<level>`.
- When the in-scope root aggregate rating is *not assessed*, the gate **MUST**
  treat the run as failing and exit `1`: an unrated root cannot clear a quality
  bar.
- When `<level>` does not name a level of the model's rating scale, the command
  **MUST** fail with the usage-error exit (`2`).

The gate changes only the exit code; with or without `--fail-at-or-below`, the
command still derives and writes both report files. Without
`--fail-at-or-below` the command exits `0` on a successful render regardless of
the rating. On stderr the gate **SHOULD** state which rating was compared against
which `<level>` and whether it passed, so a CI log explains the outcome.
