---
type: Change Case
title: Evaluation Report CTA
description: Make /quality evaluate closeouts direct users to the human reports with clear value-oriented CTAs.
status: Done
tags: [quality-skill, evaluation, reports, agent-mediated-ux]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluation Report CTA

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0151-evaluation-report-cta/spec.md) - what the case must do.
- [Design doc](0151-evaluation-report-cta/design.md) - how it is built, and why.

## Motivation

After `/quality evaluate` completes, the user should immediately know which
human report to open and why it is worth opening. The current closeout names
`report.md` and `recommendations.md`, but it does not make either one a strong
call to action or explain the different value each report provides. It also
risks treating report artifacts as a generic list rather than guiding the next
reader action.

The closeout should sell the value of the two human reports without adding
machine-oriented noise. `report.md` is the decision-ready evaluation result:
rating, evidence basis, limits, top findings, and top recommendations.
`recommendations.md` is the action-planning report: ranked recommendations, why
they matter, expected benefit, and how to know each worked. Machine indexes such
as `data/evaluation-output-result.json` remain important implementation
artifacts, but they should not be framed as the user-facing report CTA.

## Scope

Covered:

- Redesign the `/quality evaluate` closeout's report-reading call to action.
- Require full paths to the human report files, not bare filenames that require
  the user to infer the run folder.
- Describe the value of `report.md` and `recommendations.md` in the closeout.
- Keep generated report data indexes out of the primary user-facing closeout
  CTA.
- Preserve the existing evaluation boundary: no recommendations applied, no
  source edits, no `QUALITY.md` edits, no quality changelog, and no external
  issues during evaluate.

Deferred:

- Changing generated report contents or filenames.
- Changing CLI report-build output or JSON receipts.
- Changing recommendation follow-up behavior.
- Adding tests for agent prose templates unless implementation introduces a
  checkable fixture or rendered example.

## Affected artifacts

Derived by sweeping for `Evaluation complete`, `Reports:`, `report.md`,
`recommendations.md`, `recommendation index`, report paths, report closeout, and
agent-mediated UX guidance.

**Code**

- [x] No planned code impact. Revisit only if the implementation adds a rendered
      closeout fixture or testable prose template.

**Format spec and durable specs** (substance in the [functional spec](0151-evaluation-report-cta/spec.md))

- [x] `specs/skills/quality-skill/reporting.md` - require a value-oriented
      human report CTA in evaluation closeouts and keep machine report indexes
      out of the primary CTA.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` - align evaluate
      completion criteria around human report paths and a report-reading CTA.
- [x] `specs/evaluation/reports/report-tree.md` - no change; current report
      content contracts already support the CTA value propositions.
- [x] `specs/cli/evaluation-report.md` - no change; CLI receipts remain
      machine/tooling output, not the agent closeout surface.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/workflows/evaluate.md` - update the runtime closeout
      template and report-building guidance.
- [x] `skills/quality/SKILL.md` - no change; existing shared reporting sentence
      remains accurate and detailed closeout language lives in the evaluate
      workflow.
- [x] `docs/guides/agent-mediated-ux.md` - no change; existing CTA and
      five-second-scan guidance is sufficient.
- [x] `CHANGELOG.md` - release-note entry for the closeout behavior.
- [x] `specs/log.md` and `skills/quality/log.md` - durable spec/runtime log
      entries.
- [x] `changes/index.md` and `changes/log.md` - Change Case lifecycle.

## Status

`Done`. Implementation, durable skill specs, runtime guidance, release notes,
and Change Case lifecycle artifacts are complete. `mise run fmt-md-check`
passes.
