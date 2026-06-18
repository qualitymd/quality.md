---
type: Functional Specification
title: Recommendation superseding
description: Allow corrected recommendation records to identify stale recommendations they supersede.
tags: [evaluation, recommendations, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Recommendation superseding

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

This change covers optional superseding metadata for evaluation recommendation
records and the report/status behavior that uses it.

It does not add deletion, in-place replacement, assessment superseding, analysis
superseding, or automatic duplicate recommendation detection.

## Requirements

A recommendation payload **MAY** include `supersedes`, a list of recommendation
record IDs or recommendation record paths that the new recommendation replaces
as the active recommendation for the same gap.

A recommendation record **MAY** store `supersedes`. When absent, the record does
not supersede other recommendations.

`qualitymd evaluation add-record recommendation` **MUST** accept `supersedes`
for recommendation payloads, stamp it into the recommendation record metadata,
and render it in the recommendation body when non-empty.

`qualitymd evaluation show-status` **MUST** report
`missing-superseded-recommendation` when a recommendation's `supersedes` entry
does not match an existing recommendation ID or path in the same run.

A missing superseded recommendation **MUST** make the run non-reportable.

`qualitymd evaluation build-report` **MUST** treat a recommendation as
superseded when any later-loaded recommendation lists its ID or path in
`supersedes`.

`report.json` **MUST** include all recommendation records. For each
recommendation summary, it **MUST** indicate whether the recommendation is
active or superseded.

`report.md` **MUST** preserve superseded recommendations in Advice, but **MUST**
visibly mark them as superseded.

The report Next Action **MUST** choose from active recommendations only. When no
active recommendation exists, it **MUST** render the same no-remediation state
used by runs with no recommendation records.
