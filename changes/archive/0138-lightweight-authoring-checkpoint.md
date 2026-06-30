---
type: Change Case
title: Lightweight Authoring Checkpoint
description: Add a lightweight direct-authoring path for QUALITY.md edits that infers intent, asks only material follow-up questions, and confirms the intended edit conversationally.
status: Done
tags: [skill, authoring, ux]
timestamp: 2026-06-27T00:00:00Z
---

# Lightweight Authoring Checkpoint

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0138-lightweight-authoring-checkpoint/spec.md) - what the
  case must do.
- [Design doc](0138-lightweight-authoring-checkpoint/design.md) - how it's
  built, and why.

## Motivation

Users often ask an agent to update `QUALITY.md` in ordinary language: add a
concern, tighten a requirement, rename a factor, capture a newly learned risk, or
adjust body context. The current skill has strong setup discovery and
recommendation-follow-up confirmation, but direct model-authoring requests sit
between those paths. A heavy intake form would slow down clear edits, while
blindly applying inferred intent can miss the user's underlying concerns, goals,
needs, or worries.

The skill needs a lightweight direct-authoring checkpoint: infer what the user
means from the request and current model, ask follow-up only when the answer
would materially change the edit, then state the intended edit and invite
adjustments before mutation. The interaction should make `looks good` a valid
confirmation when the checkpoint clearly names the change.

## Scope

Covered:

- Route direct `QUALITY.md` edit/improvement requests to a direct model-authoring
  path distinct from tooling `update`, setup, evaluate, and recommendation
  follow-up.
- Require agents to infer intent from the request, existing model, and relevant
  authoring guides before asking follow-up.
- Ask follow-up only for missing information that materially affects model
  meaning or mutation surface.
- Present a conversational intent checkpoint before mutation, naming the intended
  edit, model/body target, important boundaries, and quality-log expectation.
- Accept `looks good` and equivalent clear approval as explicit confirmation when
  the checkpoint clearly names the mutation.
- Escalate to the existing decision-brief shape for high-impact or risky model
  changes.
- Align durable skill specs, authoring guide contracts, runtime skill guidance,
  and runtime authoring guides.

Deferred:

- New CLI support for authoring `QUALITY.md` edits.
- A separate public `/quality author` workflow name.
- Persisted draft plans for direct authoring.
- Structured forms or a fixed questionnaire for all edits.

## Affected artifacts

Derived by sweeping for direct authoring, `QUALITY.md` authoring/improvement,
decision briefs, quality-log writes, authoring guide contracts, and runtime skill
dispatch.

**Code**

- [x] No Go or TypeScript code impact; this is a skill/spec/docs behavior change.

**Durable specs** (substance in the [functional spec](0138-lightweight-authoring-checkpoint/spec.md))

- [x] `specs/skills/quality-skill/quality-skill.md` - define direct model
      authoring dispatch and the lightweight intent-checkpoint contract.
- [x] `specs/skills/quality-skill/quality-log.md` - align quality-log write
      responsibility with confirmed direct model-authoring changes.
- [x] `specs/skills/quality-skill/guides/authoring.md` - require authoring
      guidance to preserve intent, target, rationale, judgment effect, unknowns,
      and quality-log routing for direct edits.
- [x] `specs/skills/quality-skill/guides/authoring/quality-log.md` - align
      meaningful-change guidance with direct model-authoring confirmations.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` - add direct QUALITY.md authoring routing and
      checkpoint guidance.
- [x] `skills/quality/guides/authoring.md` - add lightweight direct-edit
      authoring considerations.
- [x] `skills/quality/guides/authoring/quality-log.md` - align quality-log
      wording for confirmed direct model-authoring changes.
- [x] `skills/quality/log.md` - append runtime-history entry.
- [x] `CHANGELOG.md` - note the skill UX change.

No planned impact: `SPECIFICATION.md`, CLI command specs, install docs, setup
workflow mechanics, evaluation record/report artifacts, or generated schemas.

## Status

`Done`. Implemented and archived after durable skill specs, runtime skill
guidance, authoring guides, runtime logs, release notes, and Markdown formatting
checks were updated.
