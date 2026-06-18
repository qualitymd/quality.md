---
type: Design Doc
title: Recommendation superseding design
description: How recommendation superseding is stored, validated, and rendered.
tags: [evaluation, recommendations, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Recommendation superseding design

## Context

The [Recommendation superseding spec](spec.md) answers the correction workflow
gap from E15: a corrected recommendation can be appended, but reports still
choose the first recommendation as the primary Next Action.

## Approach

Extend recommendation payloads and records with an optional `supersedes`
string list. The list accepts the same two forms already used for recommendation
references:

- a recommendation ID such as `001-fix-the-gap`;
- a recommendation path such as `recommendations/001-fix-the-gap.md`.

Loading remains append-only. Status builds a recommendation lookup from both
forms and emits `missing-superseded-recommendation` for any superseding entry
that has no match in the same run.

Report assembly computes a superseded set from valid `supersedes` entries. A
recommendation is active when its file and ID are not in that set. The
recommendation summaries include a state field, and Next Action uses the first
active recommendation in load order. Advice still lists superseded
recommendations, marked as superseded, so the audit trail remains visible.

## Alternatives

**Delete or replace recommendation files.** Rejected for this change. It would
make correction history harder to audit and would require a broader mutation
command.

**Silently choose the latest recommendation.** Rejected. It would make report
behavior depend on numbering rather than explicit evaluator intent.

**Detect duplicate recommendations automatically.** Deferred. Recommendation
records can legitimately describe different remediation paths for the same
assessment, so duplicate detection needs a sharper identity rule than this
experiment has proven.

## Trade-offs and risks

This adds a small amount of schema surface to recommendation records. The gain
is deterministic active-action selection without deleting history.

The main risk is dangling superseding references. Treating them as status gaps
keeps reports from silently hiding the wrong recommendation.

## Open questions

- Whether future replacement semantics should also support assessment records.
- Whether recommendation IDs should eventually be first-class stable fields
  rather than derived from filenames.
