---
type: Change Case
title: Readable report summary
description: Reshape report-summary.md into a clearer triage artifact with reader-facing vocabulary and prominent recommendation identifiers.
status: Done
tags: [evaluation, report, cli, skill, ux]
timestamp: 2026-06-19T00:00:00Z
---

# Readable report summary

A **Change Case** capturing the *why* and *status* for reshaping
`report-summary.md` into a more readable, action-oriented triage artifact. The
detail lives in its [functional spec](0040-readable-report-summary/spec.md).

> **Done.** Implementation and durable artifact synchronization are complete; the
> change is archived.

## Motivation

The current generated summary artifact contains the right ingredients, but its
outline and labels still expose too much implementation vocabulary. It leads
with process-shaped sections such as "Top Risks", "Rating Summary",
"Limitations", and "Next Action", and it labels the headline result as "Root
rating". Those terms are accurate to the model tree, but they are not the most
direct language for a reader deciding what happened and what to ask the agent to
do next.

The concise report should read like a decision brief: key details, a compact
summary table, the top issues, recommendation handles for follow-up prompts, and
only enough scope and limitation context to keep the result honest.

## Scope

Covered: the generated `report-summary.md` outline, human-facing report-summary
labels, default scope wording, overall-rating wording, recommendation table
shape, active recommendation identifier prominence, examples, tests, and the
durable reporting contract that describes the summary artifact.

Deferred / non-goals: no `report.json` schema rename, no changed rating
semantics, no changed roll-up semantics, no new judgment record type, no
interactive report viewer, and no redesign of the full `report.md` audit trail
beyond any terminology needed for consistency.

## Affected artifacts

### Code

- `internal/evaluation/report.go` - render the revised `report-summary.md`
  outline and human-facing labels.
- `internal/evaluation/evaluation_test.go` - update summary-rendering assertions
  and regression coverage for the new section names and recommendation table.
- `internal/cli/evaluation.go` - update human CLI wording only if it exposes
  "root rating" or related report-summary vocabulary.

### Durable specs

- `specs/evaluation-records.md` - update the `report-summary.md` content
  contract and reader-facing terminology.
- `specs/skills/quality-skill/quality-skill.md` - update the reporting contract
  and skill-facing vocabulary for full evaluations and overall ratings.
- `specs/cli/evaluation-report.md` - update only if the report command spec
  needs to name the revised summary contract directly.
- `SPECIFICATION.md` - update only if durable Evaluation Report language should
  prefer "overall rating" while preserving formal Target aggregate semantics.

### Durable docs and examples

- `specs/skills/quality-skill/examples/0001-subject-quality-eval/report-summary.md`
  - update the worked summary example.
- `specs/skills/quality-skill/examples/index.md` - update example description if
  it names the old summary shape.
- `skills/quality/SKILL.md` and `skills/quality/modes/evaluate.md` - align
  runtime skill wording where it uses user-facing "whole model" or "root rating"
  language.
- `README.md` and `docs/` - no expected changes unless terminology appears in
  reader-facing report examples.

## Children

- [Functional spec](0040-readable-report-summary/spec.md) - what the revised
  summary artifact must provide.
- [Design doc](0040-readable-report-summary/design.md) - how the summary
  renderer delivers the revised artifact without report schema churn.

## Status

`Done`. Verified with `go test ./...` and targeted `dprint check` over the
changed Markdown artifacts. The change is archived.
