---
type: Design Doc
title: Assessment superseding design
description: How assessment superseding keeps corrections coherent with analysis roll-ups.
tags: [evaluation, assessments, status, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Assessment superseding design

## Context

The [Assessment superseding spec](spec.md) extends the correction workflow from
recommendations to assessment records. Assessment corrections are riskier than
recommendation corrections because analysis records bind target and root
ratings to assessment record references.

## Approach

Extend assessment payloads and records with an optional `supersedes` string
list. References accept assessment IDs such as `001-root-has-tests` and paths
such as `assessments/001-root-has-tests.json`.

Status computes a superseded set from loaded assessments in file order. A
superseding reference is valid only when it resolves to an earlier assessment
with the same assessment identity: ordered `targetPath` plus `requirement`.

Duplicate assessment detection runs over active assessments only. Missing
analysis references still use all assessment files, but analysis references to
superseded assessments produce `superseded-assessment-reference` gaps. This
requires the evaluator to rewrite the relevant analysis record, which is already
replaceable by target slug.

Planned coverage compares against active assessments only, so a superseded stale
assessment does not satisfy coverage and does not appear as an unexpected extra
record when its replacement is active.

Reports include active state on assessment summaries. Markdown keeps
superseded assessments visible, marked as superseded, but target roll-ups remain
driven by analysis records.

## Alternatives

**Allow analyses to reference superseded assessments.** Rejected. It would make
the active assessment state and the roll-up evidence disagree.

**Silently choose the latest assessment.** Rejected. Corrections must be
explicit because duplicate assessment records can be accidental.

**Delete the old assessment.** Rejected for now. Keeping the old record
preserves the audit trail and avoids adding mutation commands.

## Trade-offs and risks

This is stricter than recommendation superseding because corrected assessments
must be paired with updated analysis records. That extra step is intentional:
ratings should not silently inherit stale assessment references.
