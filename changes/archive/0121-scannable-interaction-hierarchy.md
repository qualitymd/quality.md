---
type: Change Case
title: Scannable interaction hierarchy
description: Fix the "flat wall" failure in agent-mediated decision gates, discovery questions, progress, and result blocks across the /quality skill — leading with the primary question, rendering choices as a visually separated block, capping supporting labels, and carrying hierarchy in position rather than bold alone — and bring the agent-mediated UX guide and durable skill specs into conformance.
status: Done
tags: [skill, ux, agents, docs, workflows]
timestamp: 2026-06-26T00:00:00Z
---

# Scannable interaction hierarchy

A **Change Case** to fix a recurring scan failure in agent-mediated workflow
output: the _flat wall_, where a decision gate (or discovery question, or result
block) stacks the primary question and several supporting labels at equal
visual weight, so the actual call to action — the choice the user must make —
becomes one more bolded line at the bottom and disappears entirely when bold is
not rendered.

Detail lives in:

- [Functional spec](0121-scannable-interaction-hierarchy/spec.md) — what must be
  brought into conformance.

No design doc is required unless implementation discovers a reusable rendering
abstraction worth extracting.

## Motivation

A user reviewing an `Apply update plan?` gate from the `/quality` update workflow
could not tell at a glance what was being asked. The gate followed the
established pattern faithfully — `Changes`, `Evidence/reason`, `Recommended
option`, `Alternatives`, `Done criterion / verification`, then a prose `Answer:`
line — but that pattern _is_ the problem. Six bolded labels compete for
attention, the `y`/`n` choice is the last line after a wall of rationale, and
with bold stripped nothing distinguishes the question from its supporting
context.

This is not a one-off. The shape originates in the shared interaction contract
([`skills/quality/SKILL.md`](../../skills/quality/SKILL.md)) and the
[agent-mediated UX guide](../../docs/guides/agent-mediated-ux.md)'s own Decision
gates example, then is copied into every mutation gate in the skill. The same
flattening recurs in discovery-question rationale ordering, the absence of
re-emitted progress in long workflows, and result/closeout blocks that pile up
nine equally-weighted labels.

Because the agent is the user's interface, a gate that buries its choice costs
the user a careful read on every confirmation and erodes trust in a workflow
that is otherwise mechanically sound. The fix changes presentation only: lead
with the question, render choices as a separated block, cap supporting fields,
and carry hierarchy in position and structure so it survives bold-stripping —
without changing any mutation boundary, confirmation requirement, or quality,
evaluation, or CLI semantics.

## Scope

Covered:

- Update the [agent-mediated UX guide](../../docs/guides/agent-mediated-ux.md) so
  Decision gates teach a _question → visually separated choices → demoted
  rationale_ shape, cap supporting fields, add a "don't rely on bold alone"
  rule, and include an Avoid counter-example. _(Already applied; this case
  accounts for it.)_
- Reshape the shared `/quality` interaction contract's decision-brief template
  and its "make the primary action stand out (with bold)" prose to match.
- Bring every mutation gate in the bundled skill into the new shape
  (`setup`, `update`, recommendation follow-up, and the issue-creation gate).
- Reorder setup discovery questions so rationale precedes the answer line, and
  cap the rating-scale rationale.
- Add re-emitted, factual progress to long workflows (notably `evaluate`) at the
  guide's named drift points.
- De-stack result/closeout blocks so emphasized labels stay scannable.
- Align the affected durable `/quality` skill specs with the runtime changes.

Deferred / non-goals:

- No change to the QUALITY.md format, model schema, rating semantics, evaluation
  semantics, or CLI command behavior.
- No change to any mutation boundary or confirmation requirement; gates stay
  gates and still require explicit `y`/`n`.
- No redesign of CLI human output as CLI UX; that remains governed by
  [Designing CLI interfaces](../../docs/guides/cli-design.md).
- No rewrite of historical Change Cases, archived specs, append-only logs, or
  recorded evaluation examples.
- No shared renderer or template library unless implementation finds repeated
  structure that justifies one.

## Affected artifacts

Derived by sweeping the bundled skill and durable skill specs for decision
briefs, confirmation gates, discovery prompts, progress language, and
result/closeout blocks (the inventory in the
[functional spec](0121-scannable-interaction-hierarchy/spec.md) and the prior
gap analysis). Empty or verification-only kinds are deliberate.

### Code

- [x] `cmd/` and `internal/` — swept for hardcoded decision-brief/gate output;
      none found. CLI-only human output remains governed by the CLI design guide.
      No code impact.

### Durable specs

- [x] `specs/skills/quality-skill/quality-skill.md` — reshaped the shared user
      interaction contract's decision-brief requirement to the question →
      separated choices → demoted-rationale shape; added the field cap and the
      "hierarchy not by bold alone" rule.
- [x] `specs/skills/quality-skill/workflows/setup.md` — discovery rationale must
      precede the answer line, per-option explanation is capped; the write/update
      gates inherit the reshaped shared contract.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` — requires re-emitted
      factual progress before run creation and the per-requirement loop.
- [x] `specs/skills/quality-skill/workflows/update.md` — no change needed; the
      apply-plan gate already defers to the shared contract, now reshaped.
- [x] `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` and
      `specs/skills/quality-skill/recommendation-follow-up.md` — apply/issue gate
      and result-block requirements reshaped to the new shape.
- [x] `specs/log.md`, `specs/skills/quality-skill/workflows/log.md`, and
      `specs/skills/quality-skill/guides/log.md` — recorded the durable spec and
      guide-contract updates.

### Format spec

- [x] None — this is a skill/workflow presentation change, not a QUALITY.md
      format change.

### Durable docs (guides, AGENTS, README, bundled skill)

- [x] `docs/guides/agent-mediated-ux.md` — Decision gates reshaped, Emphasis
      "don't rely on bold alone" rule added, Avoid counter-example added,
      Checklist updated.
- [x] `skills/quality/SKILL.md` — reshaped the decision-brief template (User
      Interaction Contract) and revised the contract prose to lead with
      position/structure rather than bold.
- [x] `skills/quality/workflows/setup.md` — reshaped write/update gates, reordered
      discovery rationale before the answer line, capped the rating-scale glosses,
      and accept `y` on Q1/Q2.
- [x] `skills/quality/workflows/evaluate.md` — added re-emitted progress beats
      before run creation and before the per-requirement loop.
- [x] `skills/quality/workflows/update.md` — reshaped the apply-plan gate. The
      status block was retained as a legitimate progress beat (supports the
      progress requirement) rather than trimmed.
- [x] `skills/quality/guides/recommendation-follow-up.md` — reshaped the apply and
      issue-creation gates; de-stacked the result block with a primary outcome line.
- [x] `docs/guides/index.md` and `docs/log.md` — index already registers the
      guide; added a `docs/log.md` revision entry.
- [x] `CHANGELOG.md` — added an Unreleased `/quality Skill` note for the gate,
      discovery, and progress changes.

## Status

`Done`. Implemented, verified with `mise run fmt-md-check`, and archived. The
[agent-mediated UX guide](../../docs/guides/agent-mediated-ux.md), the bundled
`/quality` skill, and the durable skill specs all carry the reshaped
question-first gate shape.
