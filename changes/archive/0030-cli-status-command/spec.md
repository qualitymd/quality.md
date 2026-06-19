---
type: Functional Specification
title: CLI status snapshot command - functional spec
description: Specify a read-only qualitymd status command that emits a deterministic project-state snapshot.
tags: [cli, command, wizard]
timestamp: 2026-06-19T00:00:00Z
---

# CLI status snapshot command - functional spec

Companion to
[CLI status snapshot command](../0030-cli-status-command.md). This spec states
the behavior delta for a read-only `qualitymd status` command.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The `/quality` wizard needs a fast way to answer "what state is this project
in?" without becoming a second deterministic implementation. Today that state is
spread across lint output, the `QUALITY.md` target tree, evaluation run folders,
recommendation records, and generated reports. When a skill hand-parses those
artifacts, it duplicates CLI-owned mechanics and drifts from the source of truth.

`qualitymd status` makes the project-state probe a first-class CLI command. It
does not judge model quality or evaluation outcomes; it reports the existing
mechanical state so skills, CI, and humans can route to the next useful action.

## Scope

Covered: a read-only `qualitymd status [path]` command, human and `--json`
output, lint validity, model shape counts, source coverage, evaluation run
history, active recommendation counts, staleness signals, deterministic ordering,
exit behavior, and the wizard-facing next-action data.

Deferred / non-goals: no schema or `QUALITY.md` format change, no new lint
rules, no model-quality judgment, no rating recomputation, no report rendering,
no report-body scraping, no cross-repository aggregation, and no interactive
workflow. This change does not add repair or improvement behavior; it only
reports state.

## Requirements

### Invocation

`qualitymd status` **MUST** inspect `QUALITY.md` in the current working
directory by default.

`qualitymd status <path>` **MUST** inspect the file named by `<path>` instead.
The path is a caller-facing file path, matching the convention used by
`qualitymd lint [path]`.

`qualitymd status -` **MUST** fail with a usage error. Status needs a filesystem
path to relate the model to evaluation runs and staleness signals.

`qualitymd status` **MUST NOT** write, create, repair, or delete files.

### Model validity and shape

`status` **MUST** run the same mechanical validation as `qualitymd lint` for the
selected path.

When lint can inspect the file, `status` **MUST** report whether the model is
valid and the lint summary counts for errors, warnings, info findings, and
fixable findings.

When lint reports findings, `status --json` **MUST** include the lint finding
objects using the same public fields and deterministic ordering as
`qualitymd lint --json`.

When the model is lint-valid, `status` **MUST** report deterministic model-shape
counts:

- total Targets, counting the root Model as the root Target;
- total Factors, including nested sub-factors;
- total Requirements, including requirements under Factors and direct
  Target-level requirements; and
- rating-scale level count.

When the model is lint-valid, `status` **MUST** report Target source coverage for
every Target, including the root Model:

- ordered `targetPath`, with an empty array for the root Model;
- target label, using the model title when present for the root and the target
  key for child Targets;
- `sourceState`, one of `declared`, `inherited`, or `missing`;
- `source`, when a declared or inherited Source is known; and
- direct Factor, Requirement, and child Target counts.

If the model is not lint-valid, `status` **MUST** omit or null the model-shape
and source-coverage sections rather than deriving partial counts from an invalid
model.

### Evaluation history

`status` **MUST** resolve the evaluation directory using the same project
configuration as evaluation commands: `.quality/config.yaml` `evaluationDir`
when present, otherwise `quality/evaluations/`.

If the evaluation directory is absent, `status` **MUST** report zero runs rather
than failing.

`status` **MUST** recognize evaluation run folders by the existing run-folder
name contract from the Evaluation records spec.

`status` **MUST** inspect recognized runs in deterministic order by run number,
then by folder name.

`status` **MUST** report:

- total recognized run count;
- latest recognized run, if any;
- count of reportable runs;
- count of incomplete runs, defined as runs whose `evaluation show-status`
  reportability check would return `reportable: false`;
- count of stale runs, defined as runs whose `model.md` snapshot bytes differ
  from the selected model file bytes; and
- active recommendation counts.

For each reported run summary, `status` **MUST** include the run path,
reportability, stale state, record counts, gap counts, and active recommendation
count.

`status` **MUST** compute active recommendations from recommendation records and
their `supersedes` metadata. A recommendation is active when no later valid
recommendation record in the same run supersedes it.

`status` **MUST NOT** read `report.md` bodies to compute its snapshot. It MAY
read machine records such as `report.json` only when the field it reports is
already a generated machine datum and the command remains deterministic without
rerendering reports.

When a malformed run prevents inspection of that run, `status` **MUST** include
the run in the history with an inspection problem and continue inspecting other
runs. A malformed run is project state, not a command crash, unless the command
cannot read the evaluation directory itself.

### Output

Human output **MUST** summarize the selected model path, model validity, model
shape when available, evaluation history, stale or incomplete run counts, active
recommendation count, readiness, and the recommended next action.

Human output **MUST NOT** print full lint finding detail, full source coverage,
or every run by default when that would turn the status snapshot into an audit
report. It should point the caller to `qualitymd lint`, `qualitymd evaluation
show-status`, or `qualitymd evaluation build-report` for detail.

Under `--json`, `status` **MUST** emit one JSON document on stdout with
`schemaVersion: 1`, the selected model path, model validity, model shape when
available, source coverage when available, evaluation history, readiness, and
`nextActions`.

The JSON document **MUST NOT** include terminal styling, terminal control
sequences, or implementation-only fields.

`status` output **MUST** be deterministic: unchanged model file, configuration,
and evaluation run files produce byte-equivalent plain output and equivalent
JSON.

### Readiness and next actions

`status` **MUST** derive a coarse readiness state from mechanical signals:

- `missing-model` when the selected model path does not exist;
- `invalid-model` when lint can inspect the file and reports error findings;
- `ready-to-evaluate` when the model is valid and no recognized evaluation runs
  exist;
- `needs-evaluation-reconciliation` when the model is valid and at least one
  run is stale, incomplete, or has active recommendations; and
- `has-evaluation-history` when the model is valid, one or more recognized runs
  exist, and none of those runs require reconciliation.

`status` **MUST NOT** turn readiness into a quality rating. Ratings come from
evaluation records and reports, not from this command.

`status` **MUST** provide deterministic `nextActions` using the shared CLI action
shape. Suggested actions should point to the most useful next command for the
readiness state: initialize or create the model when missing, run or fix lint
when invalid, create an evaluation run when ready, inspect incomplete runs when
needed, build reports for reportable runs, or review active recommendations when
they exist.

### Exit behavior

`status` exits `0` when it successfully emits a snapshot, even when the snapshot
reports a missing model, lint errors, incomplete runs, stale runs, or active
recommendations.

`status` exits `2` for malformed invocation, including too many positional
arguments or `-` as the model path.

`status` exits `70` when it cannot emit a trustworthy snapshot because of an I/O
failure or configuration problem outside an individual malformed run.

## JSON example

```json
{
  "schemaVersion": 1,
  "path": "QUALITY.md",
  "readiness": "needs-evaluation-reconciliation",
  "model": {
    "present": true,
    "valid": true,
    "lint": {
      "summary": { "errors": 0, "warnings": 0, "info": 0, "fixable": 0 }
    },
    "shape": {
      "targets": 3,
      "factors": 5,
      "requirements": 12,
      "ratingLevels": 4
    },
    "sourceCoverage": [
      {
        "targetPath": [],
        "label": "Root subject",
        "sourceState": "declared",
        "source": ".",
        "factors": 2,
        "requirements": 1,
        "childTargets": 2
      }
    ]
  },
  "evaluations": {
    "path": "quality/evaluations",
    "runs": 2,
    "latest": {
      "path": "quality/evaluations/0002-subject-quality-eval",
      "reportable": false,
      "stale": true,
      "counts": { "assessments": 8, "analyses": 2, "recommendations": 3 },
      "gaps": 1,
      "activeRecommendations": 2
    },
    "summary": {
      "reportable": 1,
      "incomplete": 1,
      "stale": 1,
      "activeRecommendations": 2
    }
  },
  "nextActions": [
    {
      "id": "show-latest-run-status",
      "label": "Inspect the incomplete evaluation run",
      "command": "qualitymd evaluation show-status quality/evaluations/0002-subject-quality-eval"
    }
  ]
}
```

The example is illustrative, not an exhaustive schema. The durable command spec
must define every required JSON field before implementation.

## Durable spec changes

### To add

- `specs/cli/status.md` - specify `qualitymd status [path]`, its read-only
  behavior, snapshot contents, human output, JSON output, readiness states,
  deterministic ordering, and exit behavior (per the requirements above).

### To modify

- `specs/cli.md` - list `status` in the CLI surface overview and commands list
  (per the invocation and output requirements above).
- `specs/cli/index.md` - register the new `qualitymd status` command sub-spec
  (per the durable command spec addition above).

### To delete

None
