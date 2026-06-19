---
type: Change Case
title: Evaluation report summary artifact
description: Generate report-summary.md beside full evaluation reports so readers can triage a run without opening the full audit trail.
status: Done
tags: [evaluation, report, cli, skill]
timestamp: 2026-06-19T00:00:00Z
---

# Evaluation report summary artifact

A **Change Case** capturing the *why* and *status* for adding a concise
`report-summary.md` companion artifact to generated evaluation runs. The detail
lives in its [functional spec](0031-report-summary-artifact/spec.md).

> **Done.** Implementation and durable artifact synchronization are complete; the change is archived.

## Motivation

Generated evaluation reports are intentionally complete: they preserve scope,
ratings, rationales, detailed target and requirement results, evidence limits,
and advice. That audit trail is necessary, but it can become too detailed for
quick review surfaces such as PR comments, CI artifacts, release notes, and
stakeholder triage.

The current `report.md` contract is already summary-first, which helps readers
inside the full report. A separate generated `report-summary.md` would serve a
different use case: a short, routable artifact that answers "what happened, why
does it matter, and what should happen next?" while linking to the full report
for the evidence trail. Keeping it generated from the same report model avoids
manual drift and preserves the CLI-owned mechanical rendering boundary.

## Scope

Covered: generating `report-summary.md` from `qualitymd evaluation build-report`;
defining its relationship to `report.md` and `report.json`; specifying the
minimum summary contents; updating the evaluation run layout, CLI docs, skill
reporting contract, README, and example bundles.

Deferred / non-goals: no new evaluation judgment, no separate summary record
type, no `report-summary.json`, no interactive viewer, and no change to rating
semantics or roll-up behavior. The full `report.md` remains the authoritative
human Evaluation Report; `report.json` remains the machine-readable report.

Implementation is complete and archived.

## Affected specs & docs

- [x] [`specs/cli/evaluation-build-report.md`](../../specs/cli/evaluation-build-report.md)
      - require `build-report` to render `report-summary.md` alongside
      `report.md` and `report.json`, with deterministic contents derived from the
      same summary-layer data.
- [x] [`specs/evaluation-records.md`](../../specs/evaluation-records.md) - add
      `report-summary.md` to the run-folder artifact contract and clarify that it
      is generated output, not a judgment record.
- [x] [`specs/cli.md`](../../specs/cli.md) and
      [`specs/cli/index.md`](../../specs/cli/index.md) - update the command overview
      so `evaluation build-report` names all generated report artifacts.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      - update the reporting/folder layout description for the additional
      generated summary artifact.
- [x] [`README.md`](../../README.md) - update the command summary for
      `evaluation build-report`.
- [x] [`specs/skills/quality-skill/examples/`](../../specs/skills/quality-skill/examples/)
      - update example run bundles to include representative `report-summary.md`
      output or explicitly cover it in the example index.

No `SPECIFICATION.md` update is expected: `report-summary.md` is a CLI-generated
convenience artifact, not the complete conforming Evaluation Report described by
the format spec.

## Children

- [Functional spec](0031-report-summary-artifact/spec.md) - what the generated
  summary artifact must provide.
- [Design doc](0031-report-summary-artifact/design.md) - how the CLI renderer
  should generate the summary artifact from the existing report model.

## Status

`Done`. The change is archived.
