---
type: Functional Specification
title: Render interactions through native affordances — functional spec
description: Requirements for making progressive enhancement the default contract for agent-mediated interaction in the /quality skill, with the prose renderings as the text fallback.
status: Draft
timestamp: 2026-06-26T00:00:00Z
---

# Render interactions through native affordances — functional spec

The delta contract for the
[0123 change case](../0123-native-interaction-affordances.md). It governs the
shared `/quality` interaction contract and the workflow specs that inherit it.

**Normative references:**

- [agent-mediated UX guide](../../../docs/guides/agent-mediated-ux.md) — the
  durable design guidance this contract conforms the skill to; its "Channels and
  progressive enhancement" section defines the affordance taxonomy and the
  not-fit-for-purpose tests.
- [`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
  — the durable shared interaction contract these requirements modify.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The `/quality` skill is rendered entirely by an agent, and many agent harnesses
expose native interaction affordances — option pickers, multi-select checklists,
confirm/approve gates, plan-or-diff review — that render a choice more legibly
than Markdown prose. The shared interaction contract currently mandates a prose
rendering (numbered options with an `Answer` line; a `[y]`/`[n]` decision brief),
so every workflow inherits a prose-only default and ignores whatever richer
affordance the runtime offers.

One workflow already solved this locally. `setup`'s "How to present them" block
tells the agent to _"choose the presentation form from your own interaction
capabilities"_ and lays out three tiers — Markdown, a structured question tool,
and a no-affordance fallback — without naming any specific UI. That pattern is
correct and agent-agnostic; it is simply trapped in one workflow. This change
lifts it into the shared contract and reframes the prose renderings as the text
_fallback_ rather than the sole required form, so the skill leverages native
affordances when present while still degrading cleanly when absent.

The contract changes presentation only. No mutation boundary or confirmation
requirement moves; the skill still mutates only after explicit confirmation.

## Scope

Covered: the shared interaction contract in `SKILL.md` and its `quality-skill.md`
spec mirror; the workflow and guide specs that inherit it (`setup`, `update`,
`evaluate`, recommendation follow-up). Presentation contract only.

Deferred / non-goals:

- No QUALITY.md format, model, rating, evaluation, or CLI behavior change. The
  CLI is non-interactive by contract and has no prompts.
- No change to any mutation boundary or confirmation requirement.
- No shared renderer, template library, or Go code change — rendering is the
  agent's runtime job.
- No naming of a specific harness or question tool anywhere in the contract.
- The human-context checkpoint is not converted to a picker (see R8).

## Assumptions & dependencies

- The agent runtime may or may not expose native interaction affordances, and
  the skill cannot detect a specific named tool. Requirements are therefore
  written as capability-conditional behavior, not as calls to a known UI.
- The `setup` "How to present them" three-tier pattern
  ([`skills/quality/workflows/setup.md`](../../../skills/quality/workflows/setup.md))
  is the source pattern being generalized; if it changes, R1 and R9 depend on it.

## Requirements

### Shared contract: progressive enhancement

**R1 — Interactions are intents rendered through a fit-for-purpose affordance.**
The shared interaction contract **MUST** state that each user interaction is an
_intent_ (for example, a single-select closed choice with a recommended default)
rendered through the richest fit-for-purpose native affordance the runtime
exposes, and that a complete text rendering is always authored as the fallback
when no such affordance is present.

> Rationale: the prose-only default ignored native pickers/confirm gates that
> render a choice more legibly; making the intent the unit of design lets the
> skill use them when present and degrade cleanly when absent. — 0123

**R2 — Teaching lives in the message, not the widget.** The contract **MUST**
require that the question, why-it-matters, recommendation, evidence, and shortest
acceptable response are carried in the surrounding message, and **MUST NOT**
require or assume that design-critical rationale is legible inside widget option
labels.

> Rationale: native option labels are small, vary in how much they display, and
> are sometimes truncated or stripped; teaching compressed into a label is lost
> on those surfaces. — 0123

**R3 — Agent-agnostic phrasing.** The contract **MUST** express interactions as
an intent plus an affordance _category_ and **MUST NOT** name a specific question
UI, tool, or harness. Capability-conditional phrasing ("if the runtime exposes a
structured single-select affordance, render through it; otherwise emit the text
fallback") is the required form.

**R4 — Prose renderings are the text fallback, not the mandated form.** The
numbered-option closed-choice rendering and the `[y]`/`[n]` decision-brief
template **MUST** be specified as the text rendering used when no fit-for-purpose
native affordance is present, not as the sole required form. Within that
fallback, the existing rules still hold: closed choices put the recommended
option first and accept `1` as the shortest confirmation, and binary
confirmations keep `y` and `n` as the visible shortest responses.

> Rationale: the renderings are good and remain the floor; the failure was
> mandating them as the ceiling, which foreclosed richer affordances. — 0123

**R5 — Affordance taxonomy.** The contract **SHOULD** enumerate, or reference the
guide's enumeration of, the affordance categories it covers: single-select,
multi-select, binary confirm/approve, plan-or-diff review, harness authorization
prompt, free-text, and progress indicator. Naming the categories is what lets a
workflow map each interaction to one.

### Gates and authorization

**R6 — Do not double-gate a harness-authorized mutation.** Where the harness
provides its own authorization prompt for a mutation (a tool-permission or
approval prompt), the skill **SHOULD** render the confirmation through that native
gate rather than stacking an additional prose decision brief for the same
mutation, keeping the teaching in the preceding message. It **MUST** still not
mutate without confirmation; the rule removes a redundant second gate, never the
confirmation itself.

> Rationale: a hand-rolled `[y]`/`[n]` brief on top of a native authorization
> prompt is redundant friction; the binary-confirm intent is better rendered
> through the native gate — but dropping confirmation entirely is out of bounds.
> — 0123

### Inherited applications

**R7 — Named multi-value choices are single-select intents.** The `setup`
lifecycle, risk-tolerance, and rating-scale discovery questions and the
recommendation-follow-up apply-vs-hand-off outcome **MUST** be specified as
single-select closed-choice intents subject to R1–R4, so the inventory of
picker-eligible choices is explicit rather than left implicit in prose examples.

**R8 — Free-text checkpoints are preserved.** The human-context checkpoint
**MUST** remain a confirm-or-correct free-text interaction and **MUST NOT** be
forced into fixed selectable options. Its cardinality is open, so a picker is not
fit-for-purpose; R1 selects free-text here.

> Rationale: the reframe must not over-correct into widget-everything; an
> open-ended correction surface is the case where the text rendering is the right
> one. — 0123

**R9 — Say the pattern once.** Once the shared contract carries the
progressive-enhancement pattern, `setup`'s local "How to present them" block
**MUST** defer to the shared contract rather than restate the three tiers, so the
contract has a single home and cannot drift.

## Durable spec changes

### To add

None — the contract lives in the existing
[`quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md) shared
interaction section; no new durable spec is warranted.

### To modify

- `specs/skills/quality-skill/quality-skill.md` — add the progressive-enhancement
  contract, intent-over-rendering framing, teaching-in-the-message rule,
  agent-agnostic phrasing, and the affordance taxonomy reference (per R1–R5);
  reframe the numbered-option and decision-brief renderings as the text fallback
  (per R4); add the no-double-gate rule (per R6).
- `specs/skills/quality-skill/workflows/setup.md` — deduplicate "How to present
  them" to defer to the shared contract (per R9); tag the lifecycle,
  risk-tolerance, and rating-scale questions as single-select intents (per R7);
  lock the human-context checkpoint as free-text (per R8).
- `specs/skills/quality-skill/recommendation-follow-up.md` and
  `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` — tag the
  apply-vs-hand-off outcome as a single-select intent (per R7) and apply the
  no-double-gate rule to the issue-creation gate (per R6).
- `specs/skills/quality-skill/workflows/update.md` — the tooling-mutation gate
  inherits the reshaped shared contract and the no-double-gate rule (per R6).
- `specs/skills/quality-skill/evaluation.md` — the Factor/Area disambiguation
  prompt inherits the shared single-select fallback contract (per R1, R4).

### To rename

None.

### To delete

None.

## Open questions

- Whether the affordance taxonomy (R5) should be enumerated inline in
  `quality-skill.md` or referenced from the guide to keep one source of truth.
  Resolved direction: reference the guide and name only the categories the
  contract acts on; settle during implementation.
