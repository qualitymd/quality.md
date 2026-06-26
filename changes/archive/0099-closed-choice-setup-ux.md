---
type: Change Case
title: Closed-choice setup UX
description: Make /quality setup closed-choice discovery prompts use numbered options with the recommended answer first, and align user-facing option labels with the question's axis.
status: Done
tags: [skill, quality, setup, ux]
timestamp: 2026-06-26T00:00:00Z
---

# Closed-choice setup UX

This parent concept captures the *why* and *status*; the detail lives in its
children:

- [Functional spec](0099-closed-choice-setup-ux/spec.md) — what the change must
  do.
- [Design doc](0099-closed-choice-setup-ux/design.md) — how it's built, and why.

## Motivation

A real setup question showed the interaction cost of a weak closed-choice
prompt:

```text
Question 4: How costly is poor quality here?
Options: high tolerance, moderate tolerance, low tolerance
Recommended: low tolerance
Answer: Reply accept to use low tolerance, or choose another option.
```

The prompt asks on a cost axis but asks the user to answer on a tolerance axis.
It also makes the default harder than it needs to be: the recommendation is
separate from the options, and the shortest confirming answer is a word instead
of a number. For an agent-mediated setup flow, that friction matters because the
agent is the interface.

Closed-choice setup prompts should make the recommended answer the default path:
numbered options, recommended first, and `1` as the shortest confirmation.
Internal model vocabulary can still use `risk tolerance`; the visible prompt
should match the question's framing and map the answer behind the scenes.

## Scope

Covered:

- Strengthen the agent-mediated UX guide for small closed-choice prompts.
- Update `/quality setup` discovery guidance so closed-choice questions use
  numbered options, put the recommended option first, and accept `1` as the
  shortest confirmation.
- Align the risk/cost discovery prompt so the visible options answer the visible
  question, while the setup model still records risk tolerance internally.
- Mirror the interaction contract in durable `/quality` skill specs so future
  skill updates preserve the behavior.

Deferred:

- No new setup question dimensions.
- No change to setup's required discovery inputs, human context checkpoint,
  review gate, model authoring, lint, feedback logs, or closeout.
- No CLI, Go, format-schema, rating, roll-up, evaluation-record, or report
  behavior change.

## Affected artifacts

### Code

- [ ] None — no CLI/Go change. (Deliberate: this is skill/spec/docs guidance
      only.)

### Format spec

- [ ] None — `SPECIFICATION.md` unaffected.

### Durable specs

- [x] `specs/skills/quality-skill/quality-skill.md` — add the general
      closed-choice interaction contract to the user interaction section.
- [x] `specs/skills/quality-skill/workflows/setup.md` — update discovery input
      and prompt-form requirements for numbered, recommended-first closed-choice
      setup questions and cost-to-tolerance mapping.

### Durable docs / bundled skill

- [x] `docs/guides/agent-mediated-ux.md` — add closed-choice guidance and
      checklist coverage.
- [x] `skills/quality/SKILL.md` — mirror the general runtime interaction
      contract.
- [x] `skills/quality/workflows/setup.md` — update setup discovery prompt
      templates and presentation rules.

### Suggested new durable specs

- None.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Spec and
design settled, durable docs/specs and runtime skill guidance updated,
formatting verified, and archived. No CLI/Go code clock.
