---
type: Change
title: Evaluation report UX
description: Make generated evaluation reports summary-first, scoped, and easier to scan at larger target counts.
status: Done
tags: [evaluation, report, cli, skill]
timestamp: 2026-06-18T00:00:00Z
---

# Evaluation report UX

A unit of work that improves `qualitymd evaluation build-report` output after
the first experiment pass found the baseline report too scan-heavy at real and
advanced repository scale. Detail lives in the child:

- [Functional spec](0018-evaluation-report-ux/spec.md) - what the summary-first
  `report.md` and clearer `report.json` outputs must provide.
- [Design doc](0018-evaluation-report-ux/design.md) - how report rendering
  should assemble the summary layer from existing records plus minimal structured
  run metadata.

## Motivation

The experiment program in [`experiment.md`](../../experiment.md) found that the
baseline generated report is workable for small fixtures but weak for larger
evaluations. DataLoader and ESLint reviews showed the same issues: scope and
limitations are buried in rationales, root grouping targets render like evidence
gaps, empty recommendations appear as `null`, and reviewers must scan details to
answer basic questions.

The V1 summary-first prototype improved the ESLint reviewer walkthrough from
2/5 to 5/5 for user experience without changing any underlying judgments.

## Scope

Covered: the generated shape of `report.md`, the corresponding `report.json`
fields needed by tools, grouping-target rendering, explicit empty
recommendation arrays, explicit null rating objects for not-assessed ratings,
and the skill/run metadata needed to make scope and limitations visible.

Deferred: changing the evaluation rating algorithm, changing assessment or
analysis judgment semantics, adding a web UI, and implementing deep fan-out or
other evaluator workflow improvements.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches
`In-Review`:

- [x] [`specs/cli/evaluation-build-report.md`](../../specs/cli/evaluation-build-report.md)
      - describe the summary-first report rendering.
- [x] [`specs/evaluation-records.md`](../../specs/evaluation-records.md) - update
      the `report.json` contract and any required run metadata fields.
- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) - clarify that human renderings
      may front-load summary, scope, limitations, evidence basis, and target
      summary while preserving the complete report.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      - require the skill to record enough plan/design metadata for the report
      summary.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - update the
      evaluation workflow to record scope, exclusions, evidence basis, and
      limitations in reportable form.

## Status

`Done`. Implemented and archived after implementing summary-first report rendering, syncing the durable specs/docs and skill prompt, and verifying the renderer on copied ESLint and DataLoader runs.
