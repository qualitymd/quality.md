---
type: Change Case
title: Render interactions through native affordances
description: Lift the setup workflow's "render through your own interaction capabilities" pattern into the shared /quality interaction contract so every workflow renders questions, closed choices, and confirmation gates through a fit-for-purpose native affordance when one is present — reframing the numbered-option and y/n prose blocks as the text fallback, keeping the teaching in the message not the widget, staying agent-agnostic, and not stacking a prose gate on a mutation the harness already authorizes.
status: Done
tags: [skill, ux, agents, docs, workflows]
timestamp: 2026-06-26T00:00:00Z
---

# Render interactions through native affordances

A **Change Case** to make progressive enhancement the default contract for
agent-mediated interaction in the `/quality` skill. Today the shared interaction
contract mandates a prose rendering — numbered options with an `Answer` line for
closed choices, a `[y]`/`[n]` decision brief for gates — and only one workflow
(`setup`) has locally climbed out of that default to say "render through your own
interaction capabilities." This case lifts that pattern into the shared contract
so every workflow inherits it, and reframes the prose blocks as the *text
fallback* rather than the sole required form.

Detail lives in:

- [Functional spec](0123-native-interaction-affordances/spec.md) — what the
  contract must require and which durable specs change.
- [Design doc](0123-native-interaction-affordances/design.md) — how the pattern
  is lifted to shared policy, alternatives, and risks.

## Motivation

The agent is the user's interface, and many agent harnesses expose native
interaction affordances — a single-select option picker, a multi-select
checklist, a confirm/approve gate, a plan-or-diff review — that render a choice
more legibly than Markdown prose can. The `/quality` skill mostly ignores them.

The shared interaction contract
([`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) and its spec mirror
[`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md))
mandates a prose rendering: number the options, put the recommended one first,
make `1` the shortest accept; for a binary confirmation, render a `[y]`/`[n]`
decision brief. Every workflow inherits that prose-only default, so closed
choices that are textbook option-picker fits — setup's lifecycle, risk-tolerance,
and rating-scale questions; recommendation follow-up's apply-vs-hand-off
outcome — are emitted as numbered prose even on a harness that could render a
real picker.

Only the setup workflow escaped this. Its "How to present them" block
([`skills/quality/workflows/setup.md`](../../skills/quality/workflows/setup.md))
already says *"Choose the presentation form from your own interaction
capabilities. Do not assume or name a specific question UI,"* with three tiers:
Markdown rendering, a structured question tool (page questions through it, keep
the teaching in the message), and a no-affordance fallback. That pattern is
correct and agent-agnostic — it is simply trapped in one workflow.

A second, related gap: some confirmation gates fire for mutations the harness
itself will prompt to authorize (external issue creation, installed-tooling
mutation). Stacking a hand-rolled prose `[y]`/`[n]` gate on top of a native
authorization prompt is redundant friction — the binary-confirm intent should be
rendered *through* the native gate, with the teaching in the preceding message.

The fix changes presentation contract only. It does not change any mutation
boundary, confirmation requirement, or quality, evaluation, or CLI semantics: a
gate stays a gate, and nothing mutates without confirmation.

## Scope

Covered:

- The [agent-mediated UX guide](../../docs/guides/agent-mediated-ux.md) gains a
  "Channels and progressive enhancement" section (intent over rendering, native
  affordance taxonomy, not-fit-for-purpose tests), reframes closed-choice and
  decision gates as renderings of an intent, adds a "don't double-gate a harness
  authorization prompt" rule, and updates the checklist. *(Already applied; this
  case accounts for it.)*
- Lift the three-tier "render through your own interaction capabilities" pattern
  out of `setup` and into the shared `/quality` interaction contract
  (`SKILL.md` and the `quality-skill.md` spec) so all workflows inherit it.
- Reframe the numbered-option closed-choice rendering and the `[y]`/`[n]`
  decision-brief template as the **text fallback** used when no fit-for-purpose
  native affordance is present — not the sole required form. The
  recommended-first and `1`/`y`/`n` shortest-response rules continue to govern
  that fallback.
- Require that the teaching (question, why-it-matters, recommendation, evidence,
  shortest response) lives in the surrounding message, not solely in widget
  option labels, and that interactions are expressed as intent plus affordance
  category, never as a named question UI.
- Add the "don't stack a prose gate on a mutation the harness already
  authorizes" rule to the shared gate policy; apply it to the issue-creation gate
  and the tooling-mutation gate.
- Tag the named multi-value discovery questions (`setup` lifecycle,
  risk-tolerance, rating-scale) and the recommendation-follow-up apply-vs-hand-off
  outcome as single-select closed-choice intents subject to the shared contract.
- Deduplicate `setup`'s local "How to present them" block to reference the shared
  contract once it carries the pattern.
- Align the affected durable `/quality` skill specs with the runtime changes.

Deferred / non-goals:

- No change to the QUALITY.md format, model schema, rating semantics, evaluation
  semantics, or CLI command behavior. The CLI is non-interactive by contract and
  has no prompts to convert.
- No change to any mutation boundary or confirmation requirement; gates stay
  gates and the skill still mutates only after explicit confirmation.
- The human-context checkpoint stays a confirm-or-correct free-text interaction
  and is **not** forced into fixed options — open cardinality makes a picker the
  wrong affordance.
- No shared renderer, template library, or code change: rendering is the agent's
  job at runtime, not the deterministic CLI's.
- No naming of a specific harness or question tool anywhere in the contract; the
  agent-agnostic posture is preserved and strengthened.
- No rewrite of historical Change Cases, archived specs, or append-only logs.

## Affected artifacts

Derived by sweeping the bundled skill and durable skill specs for discovery
questions, closed-choice prompts, decision/confirmation gates, and the existing
"interaction capabilities" pattern (the inventory in the
[functional spec](0123-native-interaction-affordances/spec.md)). Empty or
verification-only kinds are deliberate.

### Code

- [x] `cmd/` and `internal/` — swept for interactive prompts; none. The CLI is
      non-interactive by contract (`specs/cli.md`), so there is nothing to render
      through an affordance. No code impact.

### Durable specs

- [x] `specs/skills/quality-skill/quality-skill.md` — add the
      progressive-enhancement contract to the shared user interaction section;
      reframe the numbered-option and `[y]`/`[n]` renderings as the text fallback;
      require teaching in the message and intent-plus-affordance phrasing; add the
      no-double-gate rule.
- [x] `specs/skills/quality-skill/workflows/setup.md` — deduplicate "How to
      present them" to reference the shared contract; tag the lifecycle,
      risk-tolerance, and rating-scale questions as single-select intents; lock
      the human-context checkpoint as free-text confirm-or-correct.
- [x] `specs/skills/quality-skill/recommendation-follow-up.md` and
      `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` — tag the
      apply-vs-hand-off outcome as a single-select intent; apply the no-double-gate
      rule to the issue-creation gate.
- [x] `specs/skills/quality-skill/workflows/update.md` — the tooling-mutation gate
      inherits the reshaped shared contract, including the no-double-gate rule.
- [x] `specs/skills/quality-skill/evaluation.md` — the Factor/Area disambiguation
      prompt inherits the shared single-select fallback contract.
- [x] `specs/log.md`, `specs/skills/quality-skill/workflows/log.md`, and
      `specs/skills/quality-skill/guides/log.md` — record the durable spec
      updates.

### Format spec

- [x] None — this is a skill interaction-contract change, not a QUALITY.md format
      change.

### Durable docs (guides, AGENTS, README, bundled skill)

- [x] `docs/guides/agent-mediated-ux.md` — "Channels and progressive enhancement"
      section added, closed-choice and decision-gate sections reframed as
      renderings of an intent, no-double-gate rule added, checklist updated.
      *(Applied.)*
- [x] `skills/quality/SKILL.md` — add the progressive-enhancement contract to the
      User Interaction Contract; reframe the numbered-option and decision-brief
      templates as the text fallback; add the no-double-gate rule.
- [x] `skills/quality/workflows/setup.md` — deduplicate "How to present them" to
      defer to the shared contract; tag the multi-value discovery questions as
      single-select intents; keep the checkpoint free-text.
- [x] `skills/quality/workflows/update.md` — apply-plan gate inherits the
      reshaped contract and the no-double-gate rule.
- [x] `skills/quality/guides/recommendation-follow-up.md` — outcome selection as a
      single-select intent; issue-creation gate under the no-double-gate rule.
- [x] `skills/quality/workflows/evaluate.md` — disambiguation prompt inherits the
      shared single-select fallback contract.
- [x] `docs/guides/index.md` and `docs/log.md` — index already registers the
      guide; add a `docs/log.md` revision entry for the guide update.
- [x] `CHANGELOG.md` — add an Unreleased `/quality Skill` note for the
      progressive-enhancement interaction contract.

## Status

`Done`. Implemented, verified with `mise run check` (Go tests, vet, lint,
markdown formatting, npm-pack, tidy all pass), and archived. The
[agent-mediated UX guide](../../docs/guides/agent-mediated-ux.md), the shared
`/quality` interaction contract ([`SKILL.md`](../../skills/quality/SKILL.md) and
[`quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)), the
bundled skill workflows (`setup`, `update`, `evaluate`, recommendation
follow-up), and the durable skill specs all carry the progressive-enhancement
contract: interactions are intents rendered through a fit-for-purpose native
affordance when present, the numbered-option and `y`/`n` forms are the text
fallback, and a harness-authorized mutation is not double-gated.
