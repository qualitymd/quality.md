---
type: Change Case
title: Provider-affine SDK evaluator selection
description: Prefer a ready SDK evaluator matching the invoking agent, fall back through CLI auto discovery before harness judgment, and report the selected transport without asking the user to choose.
status: Done
tags: [evaluation, evaluator, selection, skill, sdk, harness]
timestamp: 2026-07-15T00:00:00Z
---

# Provider-affine SDK evaluator selection

Status note: case is **Done** and archived. Runtime guidance, durable skill
contracts, logs, release notes, R1-R5 review evidence, targeted evaluator
tests, and the complete repository gate are in sync.

## Motivation

`/quality evaluate` currently prefers in-session `harness` judgment whenever
the invoking agent can service checkpoints, even when authenticated SDK-backed
evaluators are available. That makes the ordinary path depend on the current
conversation and harness checkpoint loop instead of the runner's fresh,
isolated provider sessions. It also turns a provider-named request such as
"have Claude evaluate this" into a transport-choice question.

The default should use an independent SDK evaluator when one is ready, prefer
the provider already mediating the user's session, and reserve `harness` for
the no-SDK fallback. The workflow should determine that path without another
user choice and make the selected method and reason clear before writing
evaluation artifacts.

## Scope

Covered:

- move CLI automatic SDK discovery ahead of the harness fallback in the
  `/quality evaluate` selection policy;
- prefer the usable built-in SDK evaluator matching the invoking Codex or
  Claude harness, while preserving the CLI's deterministic winner when no
  usable provider match exists;
- resolve bare provider-named requests to that provider's SDK evaluator without
  a harness-versus-SDK question;
- explain the selected evaluator, transport, and selection basis before the
  first evaluation mutation;
- align durable skill specs, bundled runtime guidance, logs, release notes, and
  patch-release metadata; and
- cut and verify the resulting patch release.

Non-goals:

- changing standalone CLI `auto` ordering, evaluator candidate probing,
  readiness rules, receipts, or persisted evaluation artifacts;
- inferring a parent harness from undocumented environment variables or adding
  an invoker flag to the CLI;
- changing explicit evaluator or non-`auto` workspace configuration precedence;
- silently switching evaluators after a run has selected and recorded one; and
- changing the QUALITY.md format or specification version.

## Affected artifacts

Derived from a repository sweep for evaluator precedence, provider-named
disambiguation, `evaluation.evaluator`, CLI automatic discovery, harness
fallback, selection explanations, and release compatibility.

- **Change record:** this parent, `spec.md`, `design.md`, and `review.md` under
  `changes/0208-provider-affine-sdk-selection/`; `changes/index.md`,
  `changes/log.md`, and the archive index at completion.
- **Durable skill specs:**
  `specs/skills/quality-skill/quality-skill.md` updates the evaluate summary;
  `specs/skills/quality-skill/evaluation.md` owns the provider-affine,
  SDK-before-harness selection contract; and
  `specs/skills/quality-skill/workflows/evaluate.md` owns its ordered workflow.
  `specs/log.md` and `specs/skills/quality-skill/workflows/log.md` record the
  revisions.
- **Bundled skill runtime:** `skills/quality/SKILL.md` and
  `skills/quality/workflows/evaluate.md` implement the selection and explanation
  policy; `skills/quality/log.md` and `skills/quality/workflows/log.md` record
  the revisions.
- **Release surfaces:** `CHANGELOG.md`, `package.json`,
  `npm/quality.md/package.json`, and `skills/quality/SKILL.md` release metadata
  advance for the patch release.
- **CLI code, tests, and durable CLI/evaluation specs:** no behavior change is
  planned. Existing `auto` candidate receipts provide the readiness evidence
  consumed by the skill. Verification uses contract review, targeted wording
  sweeps, existing evaluator-selection tests, and the full release gate.
- **Format specification, project model, durable docs, scaffold, generated
  reports, installer, and dependencies:** no impact is planned.

## Children

- [Functional spec](0208-provider-affine-sdk-selection/spec.md) — selection,
  affinity, provider-request, explanation, and fallback requirements.
- [Design doc](0208-provider-affine-sdk-selection/design.md) — automatic
  discovery preview, provider-affinity decision table, explicit pinning, and
  no-SDK harness fallback.
- [Review ledger](0208-provider-affine-sdk-selection/review.md) — R1-R5
  durable/runtime traceability, stale-policy sweep, evaluator-test evidence, and
  complete repository gate.

## Status

`Done`. Every R1-R5 requirement passed the review ledger, and the parent and
child bundle are archived together.
