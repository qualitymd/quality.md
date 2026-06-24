---
type: Runtime Guide
title: Authoring the Agent Harness Area
description: Guidance for modeling the agent harness as a QUALITY.md area/constituent.
tags: [quality, authoring, agent-harness]
---

# Authoring the Agent Harness Area

Read this when:

- creating, revising, reviewing, or evaluating the agent-harness Area;
- deciding whether agent steering materials are germane to a composite root;
- distinguishing the agent-harness Area from the Agent Harnessability factor.

Depends on:

- `../authoring.md`
- `model-structure.md`
- `factors.md`
- `requirements.md`

---

## Agent harness as a recurring use-context constituent

The **agent harness** is the instructions that steer the agent working the project:
agent entry points, guidance files, skills, prompts, and related instructions that
orient and govern agent work. It recurs from QUALITY.md's agent/AI-assistant use
context, not from the modeled domain.

- **Do** model the agent harness by default in a composite root when it is germane.
  *A harness-less or throwaway project hits not-germane and carries no harness
  area; a germane but too-thin harness surfaces its gap rather than being silently
  dropped.*
- **Do** treat the agent harness as partly **normative** — it governs agent
  behavior, so it plays the dual area/assessment-reference role (see
  [An entity can be both an area and an assessment reference](model-structure.md#an-entity-can-be-both-an-area-and-an-assessment-reference)).
  *Watch for double-counting if its influence is also assessed inside a domain
  constituent.*
- **Do** model the agent-harness area as the project's checked-in steering
  materials: agent entry points, guidance files, skills, prompts, and related
  instructions that orient and govern agent work. *This is an **enable**
  constituent and partly **normative** artifact; runtime harness code,
  sandboxes, tools, hooks, or orchestration the project owns may deserve their
  own areas or requirements rather than being silently folded into the steering
  materials.*
- **Do** give a germane agent-harness area a real factor family, not one or two
  placeholder factors. *Illustrative candidates include `completeness`,
  `accuracy`, `currentness`, `understandability`, `coherence`, `selectivity`,
  `discoverability` or `triggerability`, `maintainability`, `trustworthiness`,
  and `assessability`. Earn the actual factors from this harness's risks and
  needs; the list is a prompt, not a required roster.*
- **Do** phrase harness-area requirements around what steering materials do,
  agnostic to the domain the project models. *Good shapes: a stable minimal
  entry point orients an agent and links deeper without exhausting context;
  recorded conventions match actual practice; steering materials point to how
  work is verified and what actions are off-limits; steering documents do not
  contradict each other; skill names and descriptions trigger the right guidance;
  executable or third-party guidance has reviewable provenance; representative
  traces or feedback logs show whether the guidance helps in real work; and the
  materials are updated when structure or workflows change. Do not make lint,
  type-check, test, or CI commands the requirement unless the served domain is
  actually software; those are one domain's instance of how work is verified or
  bounded. See the repository doctrine guide
  `../../../../docs/guides/model-quality-across-domains.md`, section "Agentic use
  context".*
- **Do** distinguish the three projections of the agent-collaboration concern:
  **Agent Harnessability** is the model-wide factor, the **agent harness** is the
  constituent, and the **agent** is the audience. *Keep the harness as an area
  when it is germane; Agent Harnessability rates how each constituent equips an
  agent, while the harness area rates the steering artifact's own quality. Encode
  that boundary in the model per the projection-boundary rule under
  [Cover the domain's constituent kinds](model-structure.md#cover-the-domains-constituent-kinds).*
