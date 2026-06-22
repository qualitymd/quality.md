---
type: Change Case
title: Prospective evaluation plan artifacts
description: Tighten the evaluation workflow so design.md and the initial plan.md are authored before assessment begins, not reconstructed after judgment.
status: Design
tags: [evaluation, skill, records]
timestamp: 2026-06-22T00:00:00Z
---

# Prospective evaluation plan artifacts

A **Change Case** tightening when evaluation planning artifacts are authored.
The detail lives in its
[functional spec](0056-prospective-evaluation-plan-artifacts/spec.md).

## Motivation

The evaluation workflow currently says the CLI seeds `design.md` and `plan.md`
when it creates a run, and the skill then fills them with judgment content.
That is directionally right, but the wording still leaves room for a bad
operational pattern: completing one or both artifacts after assessments,
recommendations, or reports are already written.

When that happens, `design.md` and `plan.md` stop being planning and method
artifacts. They become retrospective provenance, which weakens resume
diagnostics, makes coverage harder to audit, and blurs where rating rationale
belongs. The initial design and plan should exist before assessment starts; any
later changes should be explicit amendments, not an after-the-fact reconstruction
of what the evaluator already did.

## Scope

Covered:

- Clarify the `/quality evaluate` ordering so the skill authors `design.md` and
  the initial `plan.md` immediately after `qualitymd evaluation create` and
  before assessment evidence collection or record writes begin.
- Split prospective planning content from retrospective evidence: intended
  evidence basis and coverage belong in the initial plan; actual evidence,
  findings, and rating rationale belong in assessment/analysis records and
  reports.
- Allow plan amendments during a run, provided they are marked as updates and do
  not erase the original planned scope.
- Add or clarify durable artifact contracts for `design.md` and `plan.md` so
  future skill changes do not drift back into after-the-fact authoring.

Deferred / non-goals:

- No change to the numbered run folder layout or record JSON schemas.
- No change to rating semantics, recommendation follow-up, or report rendering.
- No requirement for the CLI to generate judgment content; the CLI still owns
  mechanical scaffolding, while the skill owns design and planning judgment.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0056-prospective-evaluation-plan-artifacts/spec.md#durable-spec-changes)
section. The index below is the full skimmable list, reconciled before
In-Review.

Code:

- [ ] No direct code changes are expected. If implementation discovers that
      status/report tooling depends on the clarified artifact contract, list and
      handle those code paths before In-Review.

Specs:

- [ ] [`specs/skills/quality-skill/evaluation.md`](../../specs/skills/quality-skill/evaluation.md)
      — make design and initial plan authoring a required pre-assessment step.
- [ ] [`specs/evaluation-records/plan-md.md`](../../specs/evaluation-records/plan-md.md)
      — clarify prospective plan content, optional coverage timing, and explicit
      amendment handling.
- [ ] `specs/evaluation-records/design-md.md` — add an artifact contract for
      `design.md` if the final spec keeps the design artifact independently
      reviewable.
- [ ] [`specs/evaluation-records/index.md`](../../specs/evaluation-records/index.md),
      [`specs/evaluation-records.md`](../../specs/evaluation-records.md), and
      [`specs/evaluation-records/run-folder.md`](../../specs/evaluation-records/run-folder.md)
      — update listings and cross-links if a `design.md` child spec is added.

Runtime skill and docs:

- [ ] [`skills/quality/modes/evaluate.md`](../../skills/quality/modes/evaluate.md)
      — make the procedure explicit: create run, fill `design.md` and initial
      `plan.md`, add settled coverage when useful, then assess.
- [ ] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — keep the
      artifact contract aligned with the evaluate-mode procedure if needed.
- [ ] [`CHANGELOG.md`](../../CHANGELOG.md) — add the 0056 entry when the change
      lands.

`SPECIFICATION.md` is **not** affected: this case changes the `/quality` skill's
runtime evaluation process and evaluation-record artifacts, not the QUALITY.md
format or rating model.

## Children

- [Functional spec](0056-prospective-evaluation-plan-artifacts/spec.md) — what
  the evaluation planning artifact contract must require.
- [Design doc](0056-prospective-evaluation-plan-artifacts/design.md) — how the
  specs and runtime skill land the prospective-artifact workflow.

## Status

`Design`. The functional spec is settled enough to design the durable spec and
runtime skill changes; implementation has not started.
