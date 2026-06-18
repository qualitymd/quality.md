---
type: Functional Specification
title: Assessment superseding
description: Allow corrected assessment records to identify stale assessments they supersede.
tags: [evaluation, assessments, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Assessment superseding

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

This change covers optional superseding metadata for evaluation assessment
records and the status/report behavior that uses it.

It does not add deletion, in-place replacement, analysis superseding, or
automatic duplicate resolution without explicit superseding metadata.

## Requirements

An assessment payload **MAY** include `supersedes`, a list of assessment record
IDs or assessment record paths that the new assessment replaces as the active
assessment for the same ordered `targetPath` and `requirement`.

An assessment record **MAY** store `supersedes`. When absent, the record does
not supersede other assessments.

`qualitymd evaluation add-record assessment` **MUST** accept `supersedes` for
assessment payloads and stamp it into the assessment record.

`qualitymd evaluation show-status` **MUST** report
`missing-superseded-assessment` when an assessment's `supersedes` entry does not
match an existing assessment ID or path in the same run.

`qualitymd evaluation show-status` **MUST** report
`invalid-assessment-supersedes` when an assessment supersedes a record with a
different ordered `targetPath` plus `requirement`.

`qualitymd evaluation show-status` **MUST** report
`superseded-assessment-reference` when an analysis references a superseded
assessment record.

Superseding-related assessment gaps **MUST** make the run non-reportable.

Duplicate assessment detection **MUST** consider active assessments only. A
stale superseded assessment and its active replacement **MUST NOT** produce a
`duplicate-assessment` gap.

Planned coverage comparisons **MUST** consider active assessments only.

`report.json` **MUST** include all assessment summaries and indicate whether
each assessment is active or superseded.

`report.md` **MUST** preserve superseded assessments in the Requirements
section, but **MUST** visibly mark them as superseded.
