---
type: Functional Specification
title: Analysis record
description: Runtime JSON contract for analysis records.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-22T00:00:00Z
---

# Analysis record

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It owns the independently reviewable runtime contract below; shared
responsibility, runtime-not-OKF status, schema-version, and historical-record
compatibility rules live in the parent spec.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Analysis Record

An analysis record is one JSON file per Area. Required fields:

- `schemaVersion`, `areaPath`
- `localRatingResult`, or `null` for a grouping area with no own requirements
- `factorRatingResults`
- `aggregateRatingResult`
- `assessmentResultRecords`
- `childAnalysisRecords`
- `ratingConstraints`, optional

Every rating result **MUST** use the explicit `ratingResult` object shape
defined above. `areaPath`, `factorRatingResults[].factorPath`, and rating
values are stable model identifiers, not human display titles. A
`localRatingResult: null` on an area with child analyses and no local
assessment result records represents a structural grouping Area; report outputs
must render that as a distinct structural local-rating state, not as a missing
not-assessed rating.
When present, each `ratingConstraints` entry **SHOULD** identify the binding
`assessmentResultRecord`, `requirement`, and constrained `level`.

`qualitymd evaluation analysis set` writes the complete set of analysis payloads
provided in one invocation as a single planned replacement set: every payload is
decoded and validated before any analysis file is replaced. Within that set,
`childAnalysisRecords` values resolve by record path to analysis records written
by the same invocation or already present in the same run folder.
