---
type: Functional Specification
title: qualitymd evaluation data
description: Persist, inspect, and discover Evaluation structured JSON payloads.
tags: [cli, command, evaluation, data]
timestamp: 2026-06-25T00:00:00Z
---

# qualitymd evaluation data

`qualitymd evaluation data` is the Evaluation structured-data resource for a
run. It is the CLI surface agents use to persist routine outputs under `data/`.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Commands

```text
qualitymd evaluation data set <run> < payloads.json
qualitymd evaluation data list <run>
qualitymd evaluation data get <run>
qualitymd evaluation data kinds
qualitymd evaluation data example <kind>
qualitymd evaluation data schema [<kind>]
qualitymd evaluation data verify <run>
```

Commands that accept `<run>` **MUST** also accept `--latest` and `--model`
through the shared Evaluation run selector. `--model <model>` selects the
`QUALITY.md` file whose model-relative workspace supplies `--latest` history.
When `--model` is supplied with a relative positional `<run>` path, the command
**MUST** resolve that path relative to the selected model's workspace root. When
`--model` is absent and a positional `<run>` path is supplied, the command
**MAY** preserve ordinary filesystem-path behavior.

## Background / Motivation

`data set` is driven primarily by the `/quality` evaluate workflow, which
authors many routine payloads for one run. A JSON-array transport lets the skill
validate and persist a scope's payloads in one invocation instead of looping one
CLI call per Requirement, Factor, or Area. The batch is validated before writing,
so a failed submission is a no-op the agent can correct and retry.

`data schema <kind>` is the legible constraint source for authoring one payload
kind. Required fields and closed enum value sets need to be visible in that
single-kind schema itself; otherwise agents learn constraints by trial-and-error
validation failures instead of reading the contract.

## Requirements

`data set` **MUST** read a single non-empty JSON array from stdin whose elements
are structured JSON payload objects. It **MUST** reject stdin that is not a JSON
array, including a bare JSON object, and **MUST** reject an empty array as a
usage error.

For each payload element, `data set` **MUST** validate it, route it by `kind`,
derive its canonical `data/**` path from structured IDs, and write canonical
JSON. A one-payload write is represented as a one-element array.

Validation **MUST** reject unknown or misspelled fields, wrong field types,
out-of-range enum values, and missing required fields for the payload kind. The
diagnostic **MUST** name the offending field when one is known.

Fixed enum validation **MUST** use the same canonical value sets as
`data schema`. `data set` and `data verify` **MUST NOT** accept display labels,
emoji or shape markers, aliases, case variants, or legacy values in place of
canonical enum values.

For Evaluation Findings, validation **MUST** reject `severity` values outside
`critical`, `high`, `medium`, and `low`; `info` is not a severity value.

Validation **MUST** resolve every Area ID, Factor ID, Requirement ID, and Rating
Level ID in the payload against the run's `model-snapshot.md`. If the snapshot is
missing or cannot be parsed, `data set` **MUST** fail closed instead of accepting
unbound model references.

Validation **MUST** reject `AreaAnalysisResult.findings`,
`factorRelationships`, and any analysis-level finding object shape as unknown
fields.

Before validation and path derivation, `data set` **MUST** assign missing
CLI-owned `RecommendationResult.number` values as positive integers unique
within the run. Assignment **MUST** be deterministic by input order and by the
owning run's existing recommendation numbers. `data set` **MUST NOT** assign
artifact IDs to `FindingRankingResult.orderedFindings[]`.

For cross-payload validation, `data set` **MUST** validate against the effective
run data: existing persisted payloads overlaid with every normalized candidate
payload in the batch. After recommendation-number assignment, candidate batch
order **MUST NOT** affect validation.

`data set` **MUST** reject a `RequirementRatingResult` with `status: rated`
unless the effective run data includes a paired `RequirementAssessmentResult`
for the same Requirement with `status: assessed` or `status: partially_assessed`
and at least one Requirement Finding.

`data set` **MUST** reject a rated Requirement, Factor analysis scope, or Area
analysis scope with empty or absent `ratingDrivers`. For Factor and Area
analysis, this applies when the scope has `status: analyzed` and a
`ratingLevelId`.

`data set` **MUST** reject `ratingDrivers[].inputRefs[]` entries that do not
resolve to an existing routine output in the effective run data.

`data set` **MUST** reject Advice references that do not resolve in the
effective run data. Finding ranking entries and finding coverage entries
**MUST** select existing Requirement Findings. Recommendation ranking entries
and coverage entries **MUST** reference existing `RecommendationResult` numbers.
Finding ranking and Finding coverage **MUST** account for every effective
Requirement Finding exactly once.

`data set` **MUST** reject `RecommendationResult` payloads that use required
effort, ROI, quick-win, backlog-priority, priority, or numeric score fields.

`data set` **MUST** validate every payload in the batch before writing any
payload. If any element fails validation, `data set` **MUST** write nothing and
**MUST** report invalid elements with their array indexes. It **SHOULD** report
all invalid elements in one invocation rather than stopping at the first.

`data set` **MUST NOT** accept a `--file` flag.

`data set` **MUST** overwrite the derived path by default. This makes repeated
writes of the same routine output idempotent. If two elements in one batch derive
the same `data/**` path, `data set` **MUST** reject the batch as a usage error
instead of letting one element overwrite another.

`data set` **MUST** support `--dry-run` and `--json`. Under `--json`, it
**MUST** emit a batch write receipt, not the stored JSON artifacts. The receipt
**MUST** include a count and the canonical path of each write in input order.

Under `--dry-run`, `data set` **MUST** perform the same structural, typed,
required-field, enum, and model-binding validation as a real write, and **MUST
NOT** persist any payload. A valid dry run **MUST** report the same would-write
set that a real write would confirm.

`data set` **MUST NOT** accept `EvaluationOutputResult`; that payload is
CLI-owned and generated by `evaluation report build`.

`data list` **MUST** list stored Evaluation JSON artifacts and **MAY** filter
by `--kind`.

`data get` **MUST** emit the stored JSON artifact directly on stdout. It **MUST
NOT** provide a second JSON result-wrapper mode.

`data kinds` **MUST** include every kind accepted by `data set`, and **MUST**
distinguish CLI-owned kinds from agent-writable kinds.

`data example <kind>` **MUST** emit one complete valid representative JSON
artifact for the requested kind. Repeated nested object fields such as findings,
rating drivers, unknowns, missing evidence, analysis input refs, stop
conditions, and limits **MUST** include at least one representative element
where those fields appear in the requested kind. Examples **MUST** demonstrate
canonical Area, Factor, Requirement, and Rating Level reference strings in
relevant subject, input, rating, finding, rating-driver, and report-reference
fields. They **MUST NOT** be treated as an exhaustive corpus of every valid enum
value, status value, or error case; `data schema <kind>` remains the
machine-readable constraint source for those details. `data example <kind>`
**MUST NOT** provide a second JSON result-wrapper mode.

Commands that take a `<kind>` argument **MUST** accept the canonical kind name
and **SHOULD** accept the kebab-case form.

`data schema` **MUST** emit the JSON Schema generated from the same accepted-kind
contract used by `data set` and `data example`, including the same fixed enum
value sets. With no argument it **MUST** emit the schema for the full data
surface; with `<kind>` it **MUST** emit a self-contained schema for that kind so
the required fields and allowed enum values are legible from the emitted
document without dereferencing a top-level `$ref` into a separate `$defs` map.
The no-argument full-surface schema **MAY** use `$defs` and `$ref`. Finding
severity enums in the schema **MUST** use the same reduced value set:
`critical`, `high`, `medium`, and `low`. `data schema` **MUST NOT** provide a
second JSON result-wrapper mode.

When `data schema` output must be plain, including when stdout is not a terminal
or `NO_COLOR` is set, the command **MUST** write the schema as verbatim JSON and
nothing else. On a terminal, it **MAY** syntax-highlight and page the JSON for
readability.

`data verify <run>` **MUST** re-validate every persisted Evaluation JSON payload
under the run's `data/` directory against the same structural, model-bound, and
cross-payload contract used by `data set`. It **MUST NOT** modify payloads,
**MUST** identify each invalid payload and reason, and **MUST** exit non-zero
when any payload fails.

Evaluation data schema version is `3`. Current `data set`, `data verify`, and
report-building surfaces **MUST NOT** migrate, transform, accept, or render
schema version 2 Evaluation payloads as current data.

Commands whose stdout is already a JSON artifact **MAY** recognize `--json` only
to fail with a targeted usage error that says stdout is already JSON and the
command should be rerun without `--json`.
