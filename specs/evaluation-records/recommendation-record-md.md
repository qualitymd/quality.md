---
type: Functional Specification
title: Recommendation record
description: Runtime Markdown contract for recommendation records.
tags: [evaluation, records, cli, skill]
timestamp: 2026-06-22T00:00:00Z
---

# Recommendation record

This spec is part of the [Evaluation records](../evaluation-records.md) contract.
It owns the independently reviewable runtime contract below; shared
responsibility, runtime-not-OKF status, schema-version, and historical-record
compatibility rules live in the parent spec.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Recommendation Record

A recommendation record is one Markdown file per key gap. Its runtime
frontmatter **MUST** carry:

- `schemaVersion`
- `title`
- `gap`
- `evidenceLocators`
- `assessmentResultRecords`
- `remediationOptions`
- `recommendedOption`
- `doneCriterion`
- `supersedes`, optional

The Markdown body **MUST** state the gap, evidence locators, remediation options,
recommended option, and done criterion in stable human-readable sections.

When a recommendation corrects an earlier recommendation while preserving the
audit trail, it can include `supersedes`, a list of earlier recommendation
IDs or paths. Superseded recommendation records remain part of the run, but
reports treat them as inactive advice. Active selection is driven by explicit
`supersedes` intent, not by record numbering or recency: appending a corrected
recommendation without `supersedes` leaves the run reportable and renders both
files, so the report's Next Action can still point at the stale original — a
silent error. Requiring explicit superseding makes the active advice unambiguous
without making report output depend on which correction happened to be written
last. Route hints that help a reader act (affected package, path, workflow,
maintainer surface, verification command) belong in the existing recommendation
text fields rather than a dedicated schema field, for the same schema-stability
reason as assessment evidence above.
