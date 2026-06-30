---
type: Functional Specification
title: Lightweight Authoring Checkpoint - functional spec
description: What the /quality skill must do for lightweight direct QUALITY.md authoring.
tags: [skill, authoring, ux]
timestamp: 2026-06-27T00:00:00Z
---

# Lightweight Authoring Checkpoint - functional spec

Companion to the
[Lightweight Authoring Checkpoint](../0138-lightweight-authoring-checkpoint.md)
change case. This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Direct `QUALITY.md` edits are part of the primary agent-mediated experience, but
they should not feel like setup discovery or recommendation follow-up unless the
risk warrants it. The skill needs to protect model meaning and capture user
intent while keeping clear edits conversational: infer first, ask only material
follow-up questions, state the intended edit, and invite quick correction or
approval before mutation.

## Scope

Covered: direct `QUALITY.md` authoring dispatch, lightweight intent capture,
follow-up-question thresholds, mutation confirmation wording, escalation to
decision briefs for high-impact model changes, quality-log routing, durable skill
specs, runtime skill guidance, authoring guides, and release notes.

Deferred:

- a new public `/quality author` workflow;
- CLI-managed model-authoring edits;
- persisted authoring plans or draft artifacts; and
- mandatory structured forms for every direct edit.

## Requirements

### Direct authoring dispatch

- When the user asks the skill to directly edit, revise, improve, add to, remove
  from, or otherwise change an existing `QUALITY.md`, and the request is not
  setup, evaluate, tooling `update`, or recommendation follow-up, the skill
  **MUST** route the request to direct model authoring.

  > Rationale: direct authoring is already a supported purpose of the skill, but
  > without its own routing rule it falls through to evaluation-shaped defaults
  > or heavyweight follow-up gates.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - add
  > direct model authoring to invocation dispatch and workflow summaries.

- Direct model authoring **MUST NOT** introduce a new public workflow invocation
  name. It is a direct authoring path for requests to change `QUALITY.md`, not a
  fourth public workflow alongside `setup`, `evaluate`, and `update`.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - keep the
  > public workflow surface unchanged.

### Intent inference and follow-up

- Before asking the user follow-up questions, the skill **MUST** infer the likely
  authoring intent from the user's request, the current `QUALITY.md`, and the
  routed authoring guides relevant to the likely mutation surface.

  > Rationale: the agent should spend its context on understanding the request,
  > not making the user restate information already visible in the model.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` and
  > `specs/skills/quality-skill/guides/authoring.md` - require inference before
  > follow-up and route relevant authoring guides.

- The skill **SHOULD** ask follow-up only when missing or ambiguous information
  would materially change the model/body target, mutation surface, judgment
  effect, quality-log decision, or safety boundary. The skill **MUST NOT** use a
  fixed full questionnaire for routine direct edits.

  > Rationale: direct authoring needs intent capture, not ceremony; a mandatory
  > questionnaire would make clear edits feel worse than manual editing.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - define
  > material-follow-up thresholds.

- The skill **SHOULD** ask before choosing between body-only context and a
  structured model change when both are plausible; before adding, removing, or
  moving Areas, Factors, Requirements, or Rating Levels when the target is
  unclear; and before changing Rating Scale criteria, weights, required margin,
  scope, or apex unless the user's request already clearly decides the choice.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - define
  > common follow-up triggers.

### Intent checkpoint and confirmation

- Before mutating `QUALITY.md` through direct model authoring, the skill **MUST**
  present a lightweight intent checkpoint that states the inferred intent, the
  intended edit target, important boundaries, and whether a quality-log entry is
  expected.

  > Rationale: the checkpoint lets the user correct goals, needs, worries, and
  > constraints while the edit is still cheap to change.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - add the
  > lightweight intent checkpoint contract.

- The intent checkpoint **MUST** invite adjustments in conversational terms and
  **MUST** make a short approval path explicit. If the checkpoint clearly names
  the mutation, `looks good` or an equivalent clear approval **MUST** count as
  explicit confirmation to proceed.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - permit
  > `looks good` as explicit confirmation for clear direct-authoring
  > checkpoints.

- For high-impact changes that alter rating semantics, remove model coverage,
  shift scope or apex, or otherwise carry substantial judgment risk, the skill
  **MUST** use the existing decision-brief confirmation shape instead of the
  lightweight checkpoint alone.

  > Rationale: lightweight by default should not weaken confirmation for changes
  > that can substantially reshape future evaluations.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - define
  > escalation from checkpoint to decision brief.

### Quality-log routing

- When a confirmed direct model-authoring edit meaningfully alters what the
  model is or how it judges, the skill **MUST** write one quality-log entry for
  the coherent change. The entry's rationale **SHOULD** reuse the intent and
  reason presented in the checkpoint or decision brief.

  > Durable spec: modify `specs/skills/quality-skill/quality-log.md` and
  > `specs/skills/quality-skill/guides/authoring/quality-log.md` - name direct
  > model authoring as a quality-log write source.

- Direct authoring **MUST NOT** write a quality-log entry for wording-only,
  typo, formatting, or body-only clarification edits that do not alter what the
  model is or how it judges.

  > Durable spec: modify `specs/skills/quality-skill/quality-log.md` and
  > `specs/skills/quality-skill/guides/authoring/quality-log.md` - preserve the
  > curated-log boundary for direct authoring.

## Verification

- Source inspection **MUST** show direct model-authoring dispatch in the runtime
  skill prompt and durable skill spec.
- Source inspection **MUST** show the checkpoint guidance in the runtime skill
  prompt with `looks good` as an accepted clear confirmation.
- Source inspection **MUST** show authoring guide guidance for preserving user
  intent, target, rationale, judgment effect, unresolved unknowns, and
  quality-log routing.
- Source inspection **MUST** show quality-log guidance covering meaningful
  confirmed direct model-authoring changes and excluding wording-only edits.
- Markdown formatting checks **SHOULD** pass for touched Markdown files.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) - the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` - add direct model-authoring
  dispatch, material-follow-up thresholds, lightweight intent checkpoints,
  `looks good` confirmation, and decision-brief escalation.
- `specs/skills/quality-skill/quality-log.md` - name confirmed direct
  model-authoring changes as quality-log write sources when meaningful.
- `specs/skills/quality-skill/guides/authoring.md` - require authoring guidance
  to preserve intent, target, rationale, judgment effect, unresolved unknowns,
  and quality-log routing.
- `specs/skills/quality-skill/guides/authoring/quality-log.md` - align
  meaningful-change guidance with confirmed direct model-authoring edits.

### To rename

None.

### To delete

None.
