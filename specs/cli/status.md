---
type: Functional Specification
title: qualitymd status
description: Emit a deterministic project-state snapshot for a QUALITY.md file.
tags: [cli, command, status]
timestamp: 2026-06-19T00:00:00Z
---

# qualitymd status

`qualitymd status [path]` emits a read-only snapshot of the current `QUALITY.md`
project state: model validity, model shape, source coverage, evaluation history,
staleness, readiness, and next actions.

It inherits the cross-cutting CLI contract from the [CLI spec](../cli.md). This
file specifies only what is particular to `status`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

Agents and CI need one cheap, deterministic way to answer "what state is this
project in?" without scraping `QUALITY.md`, evaluation reports, or run folders by
hand. `status` is that routing surface. It reports mechanical state only; it does
not judge quality, recompute ratings, or replace evaluation reports.

## Scope

Covered: `qualitymd status [path]`, `--json`, model validity and shape, per-
Area source coverage, recognized Evaluation run history, stale and incomplete
run counts, readiness, next actions, deterministic ordering, and exit behavior.

Deferred / non-goals: no model-quality judgment, no report rendering, no lint
repair, no new lint rules, no report-body scraping, no schema change to the
QUALITY.md format, and no interactive workflow.

## Invocation

`qualitymd status` **MUST** inspect `QUALITY.md` in the current working directory
by default.

`qualitymd status <path>` **MUST** inspect the file named by `<path>` instead.

`qualitymd status -` **MUST** fail with a usage error. Status needs a filesystem
path to relate the model to evaluation runs and staleness.

`status` **MUST NOT** write, create, repair, or delete files.

## Snapshot

`status` **MUST** run the same mechanical validation as `qualitymd lint` for the
selected model path.

When lint can inspect the file, `status` **MUST** report whether the model is
valid, lint summary counts, and lint findings under `--json`.

When the model is lint-valid, `status` **MUST** report deterministic model-shape
counts for Areas, Factors, Requirements, and rating-scale levels. The Area
count includes the root Model as the root Area.

When the model is lint-valid, `status` **MUST** report source coverage for every
Area, including the root Model. Each source-coverage row **MUST** include the
ordered `areaPath`, a label, `sourceState`, known `source` value when present,
and direct Factor, Requirement, and child Area counts.

`sourceState` reports how the Area's evaluation Source is resolved:

- `declared` — the Area declares its own `source`; the `source` value is
  reported.
- `inherited` — the Area inherits a `source` declared by an ancestor Area; the
  inherited `source` value is reported.
- `default` — no Area in the chain declares a `source`, so the Area resolves to
  the document's default Source: the directory containing the QUALITY.md file.
  This is a deliberate, valid choice — a lean root that only groups child Areas
  takes the directory default and **MUST** report `default`, not a defect-like
  label.
- `missing` is reserved for a source that cannot be resolved. A lint-valid model
  **MUST NOT** produce it, because an undeclared source always resolves to the
  document default.

Human-facing labels **SHOULD** use required Area `title` values. JSON
`areaPath` values **MUST** remain identifier-based and **MUST NOT** be replaced
by titles.

When the model is invalid, `status` **MUST NOT** derive partial model-shape or
source-coverage counts from it.

`status` **MUST** resolve a QUALITY.md workspace from the selected model file,
using the same root `config` pointer, config file default, and evaluation
directory precedence as `qualitymd evaluation create`: command override where
available, then `evaluationDir` in the resolved config file, then
`.quality/evaluations/`. Status has no evaluation-directory override flag, so it
uses the config/default portions of that precedence.

Status output **MUST** report `.quality/evaluations` as the default evaluation
history path when no resolved config file changes it.

If the evaluation directory is absent, `status` **MUST** report zero recognized
runs rather than failing.

`status` **MUST** recognize Evaluation run folders by the current
[Evaluation](../evaluation/evaluation.md) contract and inspect them in
deterministic order by run number, then folder name.

For recognized runs, `status` **MUST** report total run count, latest run,
reportable run count, incomplete run count, and stale run count. A stale run is
one whose `model-snapshot.md` snapshot bytes differ from the selected model file bytes.

Malformed, unreadable, schema-incompatible, missing, or structurally incomplete
Evaluation data **MUST** make its run incomplete with status gaps. Broken
Evaluation skeletons **MUST** be reported as run inspection problems.

For each run summary under `--json`, `status` **MUST** include run path,
reportability, stale state, evaluation data artifact count, gap count, and any
inspection problem.

`status` **MUST NOT** read generated Markdown report bodies to compute the
snapshot.

Malformed run folders that prevent basic run inspection **MUST** appear in the
snapshot with an inspection problem, and `status` **MUST** continue inspecting
later runs. Failure to read the evaluation directory itself remains an internal
error.

## Readiness

`status` **MUST** derive one coarse readiness state from mechanical signals:

- `missing-model` when the selected model path does not exist;
- `invalid-model` when lint can inspect the file and reports error findings;
- `ready-to-evaluate` when the model is valid and no recognized evaluation runs
  exist;
- `needs-evaluation-reconciliation` when the model is valid and at least one run
  is stale, incomplete, malformed, or unsupported; and
- `has-evaluation-history` when the model is valid, one or more recognized runs
  exist, and none require reconciliation.

Readiness **MUST NOT** be treated as a quality rating.

`status` **MUST** provide deterministic `nextActions` using the shared CLI action
shape. Suggested actions should point to the most useful next command for the
readiness state: initialize the model, review or fix lint findings, create an
evaluation run, inspect an incomplete run, build a report, or create a fresh run
for stale state.

## Output

Human output **MUST** be compact and route-oriented. It **MUST** summarize model
presence and validity, model shape when available, evaluation history,
readiness, and the recommended next action. It **MUST NOT** print every lint
finding, source-coverage row, or run detail by default.

Under `--json`, `status` **MUST** emit one JSON document on stdout with
`schemaVersion: 1`, selected model path, readiness, model status, evaluation
history, and `nextActions`.

The JSON document **MUST NOT** include terminal styling, terminal control
sequences, or implementation-only fields.

`status` output **MUST** be deterministic: unchanged model file, configuration,
and evaluation run files produce byte-equivalent plain output and equivalent
JSON.

## Exit status

`status` exits `0` when it successfully emits a snapshot, even when the snapshot
reports a missing model, lint errors, incomplete runs, stale runs, or active
recommendations.

`status` exits `2` for malformed invocation.

`status` exits `70` when it cannot emit a trustworthy snapshot because of an I/O
failure or configuration problem outside an individual malformed run.
