---
type: Functional Specification
title: Scannable interaction hierarchy — functional spec
description: Requirements for reshaping agent-mediated decision gates, discovery questions, progress, and result blocks so the primary call to action is scannable and survives bold-stripping, across the /quality skill and its durable specs.
tags: [skill, ux, agents, workflows]
timestamp: 2026-06-26T00:00:00Z
---

# Scannable interaction hierarchy — functional spec

Companion to
[Scannable interaction hierarchy](../0121-scannable-interaction-hierarchy.md).
This spec states what the conformance pass must do. No design doc is required
unless implementation discovers a reusable rendering abstraction.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Scope

This change applies the reshaped Decision gates and Emphasis guidance in
[Designing agent-mediated UX](../../../docs/guides/agent-mediated-ux.md) to live
agent-mediated workflow surfaces: the bundled `/quality` skill, its durable
specs, and recommendation follow-up guidance.

Historical records, archived Change Cases, append-only logs, recorded evaluation
examples, and CLI-only human output are out of scope unless they are live
templates for future agent-mediated output. The change adjusts presentation
only; it MUST NOT alter any mutation boundary, confirmation requirement, or
quality, evaluation, or CLI semantics.

## Requirements

### Inventory

The implementation MUST audit live agent-mediated workflow surfaces before
editing them: the shared `/quality` interaction contract, `setup`, `evaluate`,
`update`, recommendation follow-up, and any code or scaffolded text that
hardcodes agent-mediated workflow output. The audit MUST identify every decision
gate, discovery question, progress beat, and result/closeout block.

The implementation MUST leave historical material alone unless it is also a live
template or current source of truth.

### Decision gate shape

Every decision gate that asks the user to confirm a mutation MUST lead with the
decision question as the strongest element of the block, before any supporting
label.

A decision gate MUST render its choices as a visually separated block — one
choice per line — distinct from the supporting context, rather than as a single
prose answer sentence appended below the rationale. The recommended choice MUST
be marked inline.

A decision gate MUST NOT present the question, its supporting labels, and the
call to action at the same visual weight. The user MUST be able to locate the
question and the response path by glance, not by reading the whole block.

A decision gate's hierarchy MUST NOT depend on bold rendering alone: with
emphasis removed, the question and the choice block MUST still be distinguishable
from the supporting context by position and structure.

A decision gate SHOULD carry at most three supporting fields beyond the question
and choices. For a binary gate, the alternatives MUST be folded into the choice
lines rather than listed as a separate field.

A decision gate MUST still require an explicit `y`/`n` response and MUST state
what will not happen when that mutation boundary matters.

> Rationale: a real `Apply update plan?` gate buried its `y`/`n` choice under six
> equally-weighted bolded labels; the user could not tell what was being asked.
> Leading with the question, separating the choices, capping the fields, and not
> relying on bold are what make the ask scannable.

### Discovery questions

Setup discovery questions MUST keep the purpose, recommended answer, and
confidence/evidence adjacent to the question and ordered _before_ the answer
line, so rationale supports the choice rather than trailing it.

Discovery rationale MUST NOT grow into a wall that buries the choice. Where a
question needs extended per-option explanation, the inline rationale SHOULD be
capped and the detail moved to a reference the agent cites only on request,
without removing the teaching the workflow requires.

### Progress

Multi-step workflows MUST re-emit a short, factual progress indication at the
points where the user's mental model would otherwise drift — before a long
context scan, after a tool-dependent phase, and before a mutation gate — not only
in the opening frame. The evaluate workflow MUST surface progress before its run
creation gate and before its long per-requirement loop.

Progress output MUST remain user-facing and factual; it MUST NOT become a
transcript of internal reasoning or tool chatter.

### Result and closeout blocks

Result and closeout blocks MUST keep their required content (outcome, changed
artifacts, validation, gaps or limitations, next action, and any workflow-
specific required content) but MUST present it so the emphasized labels stay
scannable. A block SHOULD lead with a primary outcome line and demote secondary
detail; it MUST NOT stack so many equally-weighted bold labels that the outcome
is no longer the strongest element.

### Shared interaction contract

The shared `/quality` interaction contract MUST define the decision-gate shape,
the supporting-field cap, and the "hierarchy not by bold alone" rule once, so the
individual workflows inherit a single source of truth rather than each carrying
its own gate template.

### Verification

Before moving the Change Case to `In-Review`, the implementation MUST reconcile
the parent's affected-artifact list against the actual edits, and MUST confirm
that for each reshaped gate, stripping Markdown emphasis still leaves the
question and the response path obvious.

Verification MUST include `mise run fmt-md-check`. If the implementation touches
Go code or CLI output tests, verification MUST include the narrowest relevant Go
tests and `mise run check` unless a documented environment blocker prevents it.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` — reshape the shared user
  interaction contract's decision-brief requirement to question → separated
  choices → demoted rationale; add the supporting-field cap and the
  "hierarchy not by bold alone" rule.
- `specs/skills/quality-skill/workflows/setup.md` — require discovery rationale
  to precede the answer line, cap the rating-scale rationale, and apply the new
  gate shape to the write and update gates.
- `specs/skills/quality-skill/workflows/evaluate.md` — require re-emitted factual
  progress at drift points and the new stop/gate shape, preserving evaluation
  content requirements.
- `specs/skills/quality-skill/workflows/update.md` — require the reshaped
  apply-plan gate.
- `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` and
  `specs/skills/quality-skill/recommendation-follow-up.md` — require the reshaped
  apply and issue-creation gates and the de-stacked result block.
- `specs/log.md`, `specs/skills/quality-skill/workflows/log.md`, and
  `specs/skills/quality-skill/guides/log.md` — record the durable spec updates.

### To rename

None.

### To delete

None.
