---
type: Design Doc
title: Agent Harnessability naming — design
description: Why the harnessability factor should use the Agent Harnessability title, `agent-harnessability` key, legacy handling, and accountability-preserving definition.
tags: [skill, authoring, setup, factors, harnessability, agentic]
timestamp: 2026-06-24T00:00:00Z
---

# Agent Harnessability naming — design

## Context

Answers the [functional spec](spec.md) for the
[0085 change case](../0085-agent-harnessability-naming.md). The work is a naming
and guidance refinement over the 0081 factor. It does not redesign the factor,
change the six sub-factors, or alter the QUALITY.md schema.

The design question is how to make the factor intelligible in generated models
without blurring the boundary between the model-wide factor and the agent harness
constituent, and without implying that agent work removes human responsibility.

## Approach

### 1. Use Agent Harnessability as the display title

The display title should teach the factor at the point of use. **Agent
Harnessability** keeps the harness-engineering term of art while making the
audience explicit. It reads as "harnessability for agent work," not as a generic
property that a reader has to decode from surrounding prose.

The title deliberately has no hyphen. Human-facing factor titles appear in setup
summaries, status output, and reports; they should read like labels, not stable
identifiers. The stable key carries the hyphenation: `agent-harnessability`.

### 2. Recommend `agent-harnessability` for new and revised models

Factor keys are model-local identifiers, so this is guidance rather than a schema
change. New setup output should converge on `agent-harnessability`, and
model-authoring work should recommend renaming old `harnessability` keys when it
is already editing the model. Existing models are not broken by the old key.

This separates two responsibilities:

- recognition: old `harnessability` with the 0081 sub-factor shape still counts as
  semantic coverage of the concern;
- authoring: new or revised models use `agent-harnessability` so the model teaches
  the concern more clearly.

### 3. Replace attention-scarcity shorthand with accountability wording

The old definition used phrases such as "limited human attention," "largely
unsupervised," and "synchronous supervision." Those phrases point at a real design
goal: a good harness reduces avoidable re-specification and review toil. But they
can also sound like the goal is less human responsibility.

The revised definition makes the boundary explicit:

> Agent Harnessability is the degree to which the project's checked-in materials,
> tools, workflows, feedback signals, standards, and action limits equip an AI
> agent to understand the project, take scoped work, operate the environment,
> verify its output, and stay safely bounded while preserving clear human
> direction, review, and accountability.

The design keeps the operational loop from 0081 but changes the supervision
register. The factor now describes governed agent work, not autonomous
responsibility transfer.

### 4. Preserve the factor/constituent/audience boundary

The rename should not collapse Agent Harnessability into the **agent harness**
area. The three projections remain:

| Projection  | Model term           | Meaning                                     |
| ----------- | -------------------- | ------------------------------------------- |
| Factor      | Agent Harnessability | how each constituent equips agent work      |
| Constituent | agent harness        | the steering artifact's own quality         |
| Audience    | agent                | the actor whose work the project must equip |

This is why the title uses "Agent Harnessability" rather than "Agent Harness
Quality." The former names a project-wide quality; the latter would pull readers
back toward the constituent.

### 5. Leave history historical

The archive should keep 0081's original language. It records the decision as made
at the time. The current guidance, spec mirrors, README, and CHANGELOG should move
forward; historical Change Cases and append-only logs should not be mass-rewritten.

## Alternatives

- **Keep `harnessability` and only improve the definition.** Rejected. The better
  definition helps, but the title still fails the first-read test in generated
  models and still sounds like it may name the harness artifact's quality.
- **Use `Agent-Harnessability` as the display title.** Rejected. The hyphen makes
  sense for the stable key but is awkward in prose and reports. A report should
  read "Agent Harnessability"; YAML should read `agent-harnessability`.
- **Use `Agent Harness Quality`.** Rejected. It implies the factor assesses the
  agent harness constituent, which is exactly the boundary the guidance must
  preserve.
- **Keep "limited human attention" in the definition.** Rejected as the lead
  shorthand. It names a real constraint, but it can imply reduced human
  responsibility. The current wording can discuss reduced re-specification and
  review toil only after human direction, review, and accountability are explicit.
- **Treat old `harnessability` as missing.** Rejected. It would turn a naming
  improvement into unnecessary churn for existing models that already carry the
  concern semantically.

## Trade-offs & risks

- **Temporary dual naming.** Existing models and historical docs may say
  `harnessability` while current guidance says `agent-harnessability`. The
  mitigation is explicit legacy handling: recognize the old key as coverage,
  recommend the new key during model-authoring work.
- **Longer factor title.** Agent Harnessability is less compact than
  harnessability. The added word pays for itself by making the factor intelligible
  without nearby explanatory prose.
- **Responsibility wording could obscure the efficiency goal.** The revised
  definition should still mention the operational loop and can still talk about
  reducing avoidable supervision, but only with accountability visible.

## Open questions

- Whether this repository's own dogfooded `QUALITY.md` should immediately rename
  any root factor to `agent-harnessability` is left to the implementation phase's
  affected-artifact reconciliation.
