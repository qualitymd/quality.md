---
type: Functional Specification
title: Assessment result record
description: Runtime JSON contract for assessment result records.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-22T00:00:00Z
---

# Assessment result record

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It owns the independently reviewable runtime contract below; shared
responsibility, runtime-not-OKF status, schema-version, and historical-record
compatibility rules live in the parent spec.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Assessment Result Record

An assessment result record is one JSON file per evaluated requirement. It is
the result of carrying out the Requirement's authored `assessment` instruction
against run evidence. Required fields:

- `schemaVersion`
- `areaPath`
- `requirement`
- `factorPaths`
- `ratingResult`
- `criterionSource`
- `findings`
- `recommendations`
- `supersedes`, optional

`areaPath`, `factorPaths`, and `ratingResult.level` values are stable model
identifiers: ordered Area paths, ordered Factor paths, and rating `level` ids.
They are not human display titles.

`ratingResult` **MUST** be an object with:

- `kind`, either `rated` or `not-assessed`
- `level`, required when `kind` is `rated` and omitted when `kind` is
  `not-assessed`
- `rationale`

The `kind` value is a typed rating-result state. A rated result without a
`level`, a not-assessed result with a `level`, an unknown `kind`, or an empty
`rationale` makes the record invalid for reporting.

Each finding **MUST** carry `locator`, `observation`, `category`, and
`severity`; it **MAY** carry `evidence` and `attributes`.

`severity` **MUST** be one of the canonical severity levels below. The `level`
is the stable record value; the `title` is the human display label used by
reports.

| level      | title    |
| ---------- | -------- |
| `critical` | Critical |
| `high`     | High     |
| `medium`   | Medium   |
| `low`      | Low      |
| `info`     | Info     |

Findings with `severity: "info"` are neutral evidence or supporting
observations. Findings with `critical`, `high`, `medium`, or `low` are risk
findings eligible for selected-finding summaries.

Evidence verification and locator rigor ride on these existing fields
deliberately, with no new schema field; this keeps `schemaVersion` stable and
the record mechanically gate-able. Add a dedicated field only when repeated
real-repo use shows the existing fields insufficient.

`evidence[].kind` is an intentionally open classification string. Report
renderers can display or group it, but they must not assign special semantics to
undocumented kind values.

A run **MUST NOT** contain more than one active assessment result record for the
same ordered `areaPath` and `requirement`. Duplicate active assessment result
records make the run non-reportable. This uniqueness rule exists because a
correction or resume workflow that re-adds an assessment result would otherwise
append a conflicting second record while the run still reported as renderable,
producing a report whose requirement entry and roll-up disagree. A corrected
assessment result **MAY** include `supersedes`, a list of earlier assessment
result IDs or paths for the same ordered `areaPath` and `requirement`;
superseded assessment result records remain part of the run but are not active.
Analysis records **MUST** reference active assessment result records. This is
stricter than recommendation superseding (below): because analysis ratings bind
to assessment result references, a corrected assessment result **MUST** be paired
with an updated analysis, or roll-ups would silently inherit stale judgment.
