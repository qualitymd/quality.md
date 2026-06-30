---
type: Change Case
title: Quality skill UX action clarity
description: Align /quality skill prompts, checkpoints, gates, and closeouts with the current agent-mediated UX guide.
status: Done
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Quality skill UX action clarity

This parent concept captures the _why_ and _status_; the detail lives in its
child:

- [Functional spec](0101-quality-skill-ux-action-clarity/spec.md) — what the
  change must do.
- [Design doc](0101-quality-skill-ux-action-clarity/design.md) — how the change
  is implemented.

The case is implemented and archived.

## Motivation

The current [agent-mediated UX guide](../../docs/guides/agent-mediated-ux.md)
requires each user-facing interaction to make the current state and next action
obvious, with the primary question or call to action visually dominant, the
shortest acceptable response clear, mutation gates expressed as decision briefs,
and closeouts reporting outcome, validation, gaps, and next action.

A review of the bundled `/quality` skill found places where the runtime skill and
durable skill specs still let the agent bury or soften the action the user must
take:

- the setup human context checkpoint can render as prose instead of a compact,
  explicit correction affordance;
- early setup discovery questions do not consistently name the shortest accept
  path;
- the final setup review gate can authorize writing without a full decision
  brief;
- evaluation ambiguity and stop prompts can offer options without a numbered
  answer path;
- update can inspect tooling before orienting the user with a run frame and
  progress state; and
- recommendation follow-up can leave outcome selection, issue creation, and
  result closeout less explicit than the shared UX contract requires.

This case converts those findings into implementation requirements while taking
the latest guide updates into account: the functional spec must pass the current
requirement quality bar, including set-level validation from
[Working with change cases](../../docs/guides/work-with-change-cases.md) and
[Writing functional specs](../../docs/guides/write-functional-specs.md).

## Scope

Covered:

- tighten the durable `/quality` skill specs so the agent-mediated UX contract is
  concrete at the prompt shapes that previously failed;
- tighten the bundled runtime skill files that agents actually execute;
- preserve the current UX guide as the normative source of truth unless
  implementation uncovers a guide wording gap that must be clarified; and
- update append-only skill/spec logs as needed for the touched durable specs and
  bundled skill guidance.

Deferred / non-goals:

- no CLI, Go, format-schema, rating, roll-up, evaluation-record, report-rendering,
  or `QUALITY.md` format change;
- no change to setup's model-authoring semantics, evaluation judgment semantics,
  recommendation data shape, or update installation mechanics;
- no new public `/quality` mode or command; and
- no redesign of the agent-mediated UX guide beyond narrow clarification if the
  implementation reveals an ambiguity.

## Affected artifacts

### Code

- [x] None — no Go, CLI, or generated report implementation change expected.

### Format spec

- [x] None — `SPECIFICATION.md` and the QUALITY.md format are unaffected.

### Durable specs

- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — strengthen the shared interaction contract where public modes, run
      frames, decision briefs, numbered choices, and artifact code spans need
      mode-level conformance.
- [x] [`specs/skills/quality-skill/workflows/setup.md`](../../specs/skills/quality-skill/workflows/setup.md)
      — require explicit setup checkpoint, discovery, and final review-gate
      affordances.
- [x] [`specs/skills/quality-skill/evaluation.md`](../../specs/skills/quality-skill/evaluation.md)
      and
      [`specs/skills/quality-skill/workflows/evaluate.md`](../../specs/skills/quality-skill/workflows/evaluate.md)
      — require numbered ambiguity prompts and stop-response answer paths.
- [x] [`specs/skills/quality-skill/workflows/update.md`](../../specs/skills/quality-skill/workflows/update.md)
      — require update's run frame and progress/status sequencing.
- [x] [`specs/skills/quality-skill/recommendation-follow-up.md`](../../specs/skills/quality-skill/recommendation-follow-up.md)
      and
      [`specs/skills/quality-skill/guides/recommendation-follow-up-md.md`](../../specs/skills/quality-skill/guides/recommendation-follow-up-md.md)
      — require outcome selection, issue creation, and result closeout shapes.

### Durable docs / bundled skill

- [x] [`docs/guides/agent-mediated-ux.md`](../../docs/guides/agent-mediated-ux.md)
      — source-of-truth review target; update only if implementation finds a
      guide ambiguity rather than a skill conformance gap.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — align the shared
      runtime interaction contract and ambiguous evaluation scope prompts.
- [x] [`skills/quality/workflows/setup.md`](../../skills/quality/workflows/setup.md)
      — align discovery questions, human context checkpoint, and final review
      gate.
- [x] [`skills/quality/workflows/evaluate.md`](../../skills/quality/workflows/evaluate.md)
      — align ambiguity and stop-response prompt shapes.
- [x] [`skills/quality/workflows/update.md`](../../skills/quality/workflows/update.md)
      — align run frame, progress, and update closeout guidance.
- [x] [`skills/quality/guides/recommendation-follow-up.md`](../../skills/quality/guides/recommendation-follow-up.md)
      — align outcome selection, issue creation, and result closeout guidance.

### Suggested new durable specs

- None. The current `/quality` skill specs already have durable homes for each
  affected runtime surface.

## Status

`Done`. Implemented across durable skill specs and bundled runtime skill
guidance, verified with `mise run fmt-md-check`, and archived.
