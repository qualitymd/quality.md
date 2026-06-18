---
type: Design Doc
title: Duplicate assessment status design
description: Implementation approach for duplicate assessment renderability gaps.
tags: [evaluation, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Duplicate assessment status design

## Context

`Run.Renderable()` already centralizes reportability checks shared by
`show-status` and `build-report`. It verifies missing analysis records,
missing referenced assessments, missing recommendations, and root-analysis
uniqueness.

The experiment `duplicate-record-correction-trial` showed a correction payload
for the same target path and requirement produces a second numbered assessment.
The old renderability check still passed, so the report showed conflicting
requirement entries.

## Approach

Add duplicate assessment detection to `Run.Renderable()` before analysis
reference checks:

- Build a map keyed by `strings.Join(targetPath, "\x00") + "\x00" +
  requirement`.
- Store the first assessment file for each key.
- For each later assessment with the same key, append a
  `duplicate-assessment` gap whose `Ref` is the later file and whose detail
  names the first conflicting file.

Because `build-report` already refuses any non-empty renderability gap list, no
new report-specific gate is needed.

## Alternatives

- Reject duplicates in `add-record`: better long-term ergonomics, but it would
  require defining replacement and intentional superseding behavior. This change
  only prevents ambiguous report rendering.
- Use target title plus requirement as the key: weaker than target path because
  titles are display labels.
- Silently choose the latest assessment: dangerous because it hides earlier
  record references from analysis records and makes reports order-sensitive.

## Trade-offs and Risks

Existing runs with duplicate assessments will become non-reportable until the
duplicates are removed or a later replacement/superseding mechanism exists. That
is acceptable because conflicting duplicate assessments cannot produce an
unambiguous report.
