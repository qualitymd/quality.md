---
type: Design Doc
title: Typed report model - design
description: Implementation approach for explicit typed evaluation-report states.
tags: [evaluation, report, records, types]
timestamp: 2026-06-20T00:00:00Z
---

# Typed report model - design

Companion to [Typed report model - functional spec](spec.md).

## Context

The evaluation package already owns record loading, record writing, report
assembly, and rendering. The safest implementation path is to keep external
record input names stable while tightening internal Go types and report JSON
state objects.

## Approach

Add small string-backed types beside the existing record structs:
`RatingResultKind`, `LocalRatingKind`, `ReportNextStepKind`,
`EvaluationRunGapKind`, `RecordLifecycleState`, `EvaluationRigor`,
`EvaluationLevel`, `MissingMetadataKind`, `TargetPath`, and `FactorPath`.

Keep assessment and analysis input JSON compatible by letting the typed values
marshal as their stable string ids. Where report JSON needs to distinguish
states that were previously implicit, add explicit state objects rather than
overloading nulls and booleans. Markdown rendering uses display helpers on the
typed values.

## Alternatives

Leaving strings in place and adding more validation was rejected because it keeps
policy scattered across report building and status routing. Rewriting record
schemas wholesale was rejected because most input shapes are already adequate;
the problem is the internal/report model, not the CLI payload names.

## Trade-offs & risks

Report JSON changes will affect consumers that read `localRatingResult`,
`active`, or `nextAction.kind`. The fixture update and durable spec update make
that intentional. The implementation should avoid over-typing extensible fields:
`Evidence.Kind` stays open unless renderer behavior starts depending on its
specific values.
