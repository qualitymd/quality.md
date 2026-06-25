---
type: Runtime Guide
title: Authoring the Agent Harness Area
description: Guidance for modeling the agent harness as a QUALITY.md area/constituent.
tags: [quality, authoring, agent-harness]
---

# Authoring the Agent Harness Area

Read this when:

- creating, revising, reviewing, or evaluating the agent-harness Area;
- deciding whether agent governing artifacts are germane to a composite root;
- distinguishing the agent-harness Area from the Agent Harnessability factor.

Depends on:

- `../authoring.md`
- `model-structure.md`
- `factors.md`
- `requirements.md`

---

## Agent harness as a recurring use-context constituent

The **agent harness** is the whole engineered system around the model: everything
that is not the model itself, including the code, configuration, and execution
logic that turns a model into an agent working the project. It has two halves:
feedforward controls that steer the agent before it acts (agent entry points,
guidance files, skills, prompts, tool and MCP definitions, sandbox and filesystem
rules, orchestration) and feedback controls that catch and correct it after
(verification commands, tests, run logs, evals, review). Instructions are one
feedforward component of the harness, not the harness as a whole.

The agent harness recurs from QUALITY.md's agent/AI-assistant use context, not
from the modeled domain. Define the harness at full breadth first, then scope the
agent-harness Area as one projection of it.

- **Do** model the agent harness by default in a composite root when it is germane.
  *A harness-less or throwaway project hits not-germane and carries no harness
  area; a germane but too-thin harness surfaces its gap rather than being silently
  dropped.*
- **Do** treat the agent harness as partly **normative** — it governs agent
  behavior, so it plays the dual area/assessment-reference role (see
  [An entity can be both an area and an assessment reference](model-structure.md#an-entity-can-be-both-an-area-and-an-assessment-reference)).
  *Watch for double-counting if its influence is also assessed inside a domain
  constituent.*
- **Do** distinguish the three projections of the agent-collaboration concern:
  **Agent Harnessability** is the model-wide factor, the **agent harness** is the
  constituent, and the **agent** is the audience. *Keep the harness as an Area
  when it is germane; Agent Harnessability rates how each constituent equips an
  agent, while the harness Area rates the project-owned governing artifacts'
  own quality. Harness files are evidence for the factor and the evaluated
  entity for the Area. Encode that boundary in the model per the
  projection-boundary rule under
  [Cover the domain's constituent kinds](model-structure.md#cover-the-domains-constituent-kinds).*
- **Do** model the agent-harness Area as the project's checked-in,
  project-owned governing artifacts: agent entry points, guidance files, skills,
  prompts, and related instructions, plus project-owned hooks, tool/MCP
  definitions, sandbox or permission policy, and orchestration config when those
  exist. *This is an **enable** constituent and partly **normative** artifact.
  Project-owned runtime harness machinery must be surfaced here or, when large
  enough to warrant distinct factors, given its own Area; never silently fold it
  into prose instructions or drop it.*
- **Do** scope artifacts by primary job. *Rate an artifact in the agent-harness
  Area when its primary job is to govern or equip the agent's work and the
  project owns it. Cede it to a domain constituent when it is primarily a
  product artifact the agent merely also uses: the product test suite belongs to
  tests, and the deploy runtime belongs to operations. When one artifact does
  both, rate its agent-governing quality here and cross-reference the domain
  constituent under the no-double-count rule.*
- **Do** give a germane agent-harness area a real factor family, not one or two
  placeholder factors. *Illustrative candidates include `completeness`,
  `accuracy`, `currentness`, `understandability`, `coherence`, `selectivity`,
  `discoverability` or `triggerability`, `maintainability`, `trustworthiness`,
  and `assessability`. Earn the actual factors from this harness's risks and
  needs; the list is a prompt, not a required roster.*
- **Do** phrase harness-area requirements around what harness governing artifacts do,
  agnostic to the domain the project models. *Good shapes: a stable minimal
  entry point orients an agent and links deeper without exhausting context;
  recorded conventions match actual practice; skill names and descriptions
  trigger the right guidance; steering documents do not contradict each other or
  the guides they reference; executable or third-party guidance has reviewable
  provenance; the harness points to how work is verified and routes to signals
  the agent can run or inspect; representative traces or feedback logs show
  whether the guidance helps in real work; project-owned hooks, sandbox, or
  permission policy that bound consequential action are coherent, current, and
  inspectable; and orchestration or subagent config is internally consistent.
  Rate the artifact's own quality, not the capability it confers: "the
  permission policy is coherent and current" belongs here, while "the project
  contains agent action" belongs to `containment-of-action`. Do not make lint,
  type-check, test, or CI commands the requirement unless the served domain is
  actually software; those are one domain's instance of how work is verified or
  bounded. See the repository doctrine guide
  `../../../../docs/guides/model-quality-across-domains.md`, section "Agentic use
  context".*
- **Avoid** defining the agent harness as "the instructions" or equivalent
  steering-prose shorthand. *Define the harness as the whole engineered system
  around the model - feedforward and feedback controls - then scope the Area to
  the checked-in, project-owned slice. Tools, sandbox, orchestration, and
  verification are equally harness.*
