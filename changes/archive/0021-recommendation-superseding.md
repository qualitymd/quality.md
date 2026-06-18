---
type: Change Case
title: Recommendation superseding
description: Let corrected recommendation records supersede stale recommendations.
status: Done
tags: [evaluation, recommendations, cli]
timestamp: 2026-06-18T00:00:00Z
---

# Recommendation superseding

A unit of work that turns the experiment program's recommendation-correction
finding into a deterministic report behavior. Detail lives in the child:

- [Functional spec](0021-recommendation-superseding/spec.md) - what
  recommendation superseding must do.
- [Design doc](0021-recommendation-superseding/design.md) - how superseding is
  stored, validated, and rendered.

## Motivation

The E15 recommendation-correction trial showed that appending a corrected
recommendation leaves the run reportable and renders both recommendation files,
but the primary report Next Action can still point to the stale original
recommendation. The workflow needs a way to preserve the audit trail while
making the active recommendation unambiguous.

## Scope

Covered: allow recommendation records to declare older recommendation IDs or
paths they supersede; make status validate superseding references; make reports
distinguish active recommendations from superseded recommendations and choose
Next Action from active recommendations only.

Deferred: deleting old recommendation records, replacing assessment records,
superseding assessment or analysis records, and automatic duplicate
recommendation detection.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches
`In-Review`:

- [x] [`specs/evaluation-records.md`](../../specs/evaluation-records.md) - define
      optional recommendation superseding metadata.
- [x] [`specs/cli/evaluation-add-record.md`](../../specs/cli/evaluation-add-record.md)
      - document recommendation payload support for superseding.
- [x] [`specs/cli/evaluation-show-status.md`](../../specs/cli/evaluation-show-status.md)
      - document invalid superseding-reference gaps.
- [x] [`specs/cli/evaluation-build-report.md`](../../specs/cli/evaluation-build-report.md)
      - document active recommendation selection for Next Action and reports.
- [x] [quality skill spec](../../specs/skills/quality-skill/quality-skill.md) -
      guide correction workflows to supersede stale recommendations.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - teach the skill
      to use superseding rather than appending ambiguous recommendation
      corrections.

## Status

`Done`. Implemented and archived after implementing recommendation `supersedes` metadata, dangling-reference status gaps, active/superseded report rendering, durable specs, and skill guidance.
