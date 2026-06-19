---
type: Change Case
title: /quality skill interaction UX
description: Add a durable interaction contract for the /quality skill so runs are easier to understand, decisions are explicit, low-confidence work stops cleanly, improvements close the loop, and evaluation history informs next actions.
status: Done
tags: [skill, quality, evaluation, ux]
timestamp: 2026-06-19T00:00:00Z
---

# /quality skill interaction UX

A **Change Case** for tightening the `/quality` skill's user-facing interaction
contract. The detail lives in its
[functional spec](0038-quality-skill-interaction-ux/spec.md).

## Motivation

The `/quality` skill already has strong boundaries: the CLI owns mechanical
artifact creation, the skill owns judgment, evaluated source is treated as data,
and mutating modes require confirmation. The next usability gap is not a new
quality model or CLI command. It is the interaction layer around those contracts.

When a user invokes `/quality`, they should be able to see what mode and scope
the agent inferred, when a run is read-only or mutating, why a decision is being
asked, when evidence is too weak to rate, what changed after an improvement, and
how prior evaluations affect the next step. This case records those expectations
as a durable skill contract so the runtime prompt and mode files can provide a
more predictable experience without importing an unrelated skill runtime system.

## Scope

Covered:

- Add a user interaction contract to the durable `/quality` skill spec covering
  run frames, decision briefs, stop rules, history-aware operation, improvement
  delta reports, and status-first voice/posture.
- Reflect the contract in the installable `/quality` skill prompt and relevant
  mode files.
- Keep the existing skill/CLI boundary intact: the interaction contract changes
  how the skill frames and sequences judgment, not which mechanical operations
  the skill may perform.

Deferred / non-goals:

- No new CLI command or evaluation record schema is proposed by this case.
- No telemetry, persistent user preference system, project routing injection, or
  model-specific prompt platform is added.
- No change to `QUALITY.md` format semantics, rating vocabulary, or evaluation
  artifact layout.
- No general Agent Skills specification change.

## Affected specs & docs

- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      - add the durable user interaction contract.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - add the shared
      runtime rules for run frames, decision briefs, stop rules, history context,
      and improvement delta output.
- [x] [`skills/quality/modes/wizard.md`](../../skills/quality/modes/wizard.md) -
      align readiness and next-step output with the interaction contract.
- [x] [`skills/quality/modes/evaluate.md`](../../skills/quality/modes/evaluate.md)
      - add evaluate-specific run framing, stop/reroute behavior, and history
      context.
- [x] [`skills/quality/modes/improve.md`](../../skills/quality/modes/improve.md)
      - add mutating decision briefs and before/after improvement reporting.
- [x] [`skills/quality/modes/setup.md`](../../skills/quality/modes/setup.md) -
      apply the mutating decision-brief shape where setup changes an existing
      model.
- [x] [`skills/quality/modes/upgrade.md`](../../skills/quality/modes/upgrade.md) -
      apply the mutating decision-brief shape to skill/CLI maintenance actions.

No `SPECIFICATION.md` update is expected: this case changes the `/quality`
skill's interaction behavior, not the `QUALITY.md` format.

## Children

- [Functional spec](0038-quality-skill-interaction-ux/spec.md) - what the
  interaction contract must require.
- [Design doc](0038-quality-skill-interaction-ux/design.md) - how the contract
  is added to the durable spec and runtime skill prompts.

## Status

`Done`. Implemented and archived.
