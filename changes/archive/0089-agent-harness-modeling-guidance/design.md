---
type: Design Doc
title: Agent-harness modeling guidance - design
description: How the root Agent Harnessability refinement, the agent-harness-area modeling template, the served-domain guardrail, the self-verifiability sensor/observability sharpening, and the use-context-constituent explicitness principle are shaped — the harness-engineering grounding, the candidate factor families, placement decisions, and alternatives weighed.
tags: [docs, doctrine, skill, agent-harness]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-harness modeling guidance - design

## Context

Answers the [functional spec](spec.md) for
[Agent-harness modeling guidance](../0089-agent-harness-modeling-guidance.md). The
change is guidance prose across three files — the bundled authoring guide, the Top
10 checks, and the doctrine guide — plus their durable mirrors and logs. No code,
no schema. This doc records the harness-engineering grounding, the root
Agent Harnessability refinement, the harness-area factor family chosen, where each
edit lands, and the alternatives rejected.

## Harness-engineering grounding

The candidate factor families and the self-verifiability sharpening are grounded in
current agent-harness, context-engineering, and trace-driven improvement sources
(fetched and verified 2026-06-24):

- Martin Fowler / Birgitta Böckeler, _Harness Engineering_ — **Agent = Model +
  Harness**; the harness is everything except the model. It splits into **guides**
  (feedforward controls that steer _before_ action: conventions, skills, how-tos,
  AGENTS.md, bootstrap scripts, language servers) and **sensors** (feedback controls
  that observe _after_ action: tests, linters, type checkers, review agents).
  Controls are **computational** (deterministic, fast) or **inferential** (AI,
  semantic, slower). The article names the system property **"harnessability"** —
  how amenable a system is to control — and stresses the **steering loop** (humans
  iterate the harness when issues recur) and **coherence** (guides and sensors must
  not contradict).
- Martin Fowler / Birgitta Böckeler, _Sensors for Coding Agents_ — a _good_ sensor
  is grounded in concrete evidence, actionable (custom messages carrying
  self-correction guidance and the rule's rationale), context-aware, **suppressible**
  with visible reviewable exceptions, has **flexible thresholds** (not binary), and
  gives fast feedback; computational and inferential sensors complement each other.
- LangChain, _The Anatomy of an Agent Harness_ — the harness is state/storage,
  execution environment (sandboxes, allow-listing, network isolation), the
  tool/capability layer, control/orchestration, **context management** (compaction,
  offloading, progressive disclosure against "context rot"), verification/feedback,
  knowledge/learning, and middleware/hooks.
- Anthropic, _Effective harnesses for long-running agents_ and _Harness design for
  long-running application development_ — long-running work fails when agents try to
  do too much at once, lose coherent state across context windows, or declare victory
  too early; effective harnesses decompose work, leave structured progress and
  handoff artifacts, preserve enough state for the next session, and use evaluator
  loops for closure.
- Anthropic, _Effective context engineering for AI agents_ and _Equipping agents for
  the real world with Agent Skills_ — context is a scarce, task-shaped resource;
  progressive disclosure, just-in-time retrieval, skill metadata, split reference
  files, and observed agent use are key to making steering materials scalable.
- OpenAI, _Harness engineering: leveraging Codex in an agent-first world_ — a short
  `AGENTS.md` acts as a table of contents to deeper sources of truth; agent struggles
  become signals for missing tools, guardrails, or documentation; and agents produce
  code, tests, CI, docs, evaluation harnesses, review responses, and scripts when the
  repository encodes the loop.
- OpenAI Agents SDK guardrails documentation — guardrails attach at input, output,
  and tool boundaries, supporting a domain-neutral refinement of
  `enforcement-of-standards`.
- LangChain, _Improving Deep Agents with harness engineering_, and the AHE paper
  (_Agentic Harness Engineering_) — traces and component-level observability make
  harness improvement evidence-driven; the editable harness components include
  prompts, tools, memory, and middleware, and component interactions can create
  regressions if optimized blindly.

The load-bearing finding: the projection boundary still holds, but the original
"six sub-factors unchanged" assumption should not. The agent-harness **area** is
the checked-in steering-materials artifact (the map, skills, prompts, and guidance
that an agent reads or triggers). The root **Agent Harnessability** factor is the
project-wide capability of letting an agent do useful work from prompt to verified
outcome, including long-running work. The sources map as follows:

| Harness-engineering concept                          | Projection in the model                         |
| ---------------------------------------------------- | ----------------------------------------------- |
| guides / context — what the agent sees (feedforward) | `agent-accessibility` + agent-harness area      |
| execution environment, fresh-env setup               | `agent-operability`                             |
| task framing                                         | `task-specifiability`                           |
| handoff, memory, compaction, progress state          | `continuity`                                    |
| sensors / feedback loops                             | `self-verifiability`                            |
| traces, eval records, run evidence                   | `self-verifiability` + `continuity`             |
| guardrails (prevention; input/output/tool checks)    | `enforcement-of-standards`                      |
| guardrails (containment, sandboxing)                 | `containment-of-action`                         |
| coherence / synchronization                          | harness-area `coherence`                        |
| context efficiency / progressive disclosure          | harness-area `selectivity`                      |
| skill metadata and entry-point routing               | harness-area `discoverability`/`triggerability` |
| provenance and executable/untrusted instructions     | harness-area `trustworthiness`                  |
| representative traces and feedback logs              | harness-area `assessability`                    |
| steering loop                                        | model-wide learn loop (not another sub-factor)  |

So the factor side is validated but needs one fit-for-purpose refinement:
`continuity`. This is not the model-wide learn loop; it is the assessable quality
that lets an agent preserve state and resume without depending on unrecoverable chat
history. The articles stay in this design record as grounding, not as citations in
the durable guide, and the project adopts none of their software-specific rosters
(the three "regulation categories" — maintainability/architecture/behaviour — and
tool names like ESLint, dependency-cruiser, Stryker) as defaults; that restraint is
exactly the served-domain guardrail this case adds.

## Approach

### 1. Root Agent Harnessability refinement (authoring guide)

Keep Agent Harnessability as the root model-wide factor, but revise the
decomposition to seven sub-factors:

- **agent-accessibility** — the agent can reach the right knowledge at the right
  time. Refine requirements toward progressive disclosure, context selectivity,
  durable decision records, and entry points that point rather than dump.
- **task-specifiability** — work can be handed to an agent as a bounded assignment.
  Refine requirements toward goals, non-goals, success criteria, starting points,
  done criteria, decomposition, and completion discipline that compares the outcome
  against the original task.
- **agent-operability** — the agent can operate the working environment. Refine
  requirements toward fresh-session setup, tool affordance quality, required
  tools/data/accounts, act/observe loops, and agent-useful tool output.
- **continuity** — the agent can preserve state and resume useful work across
  long-running tasks, compaction, interruption, handoff, and fresh sessions. Example
  requirements: progress or handoff artifacts capture current state, decisions,
  remaining work, verification status, blockers, and next steps; resumptions do not
  depend on unrecoverable chat history; progress records reduce false completion and
  context-anxiety failure modes. Boundary: `agent-operability` covers a fresh session
  reaching a ready-to-work environment, while `continuity` covers resuming with prior
  state and decisions; `agent-accessibility` covers durable, reachable knowledge,
  while `continuity` covers the progress and handoff record of an in-flight task.
- **self-verifiability** — the agent can check whether work is correct. Refine
  requirements toward deterministic checks, inferential evals where judgment is
  required, traces/run evidence, actionable machine-readable feedback, and comparison
  against the original task rather than only the produced artifact.
- **enforcement-of-standards** — important standards hold even when the agent forgets
  or misreads prose. Refine requirements toward deterministic gates, schemas, hooks,
  review gates, input/output/tool guardrails, constrained suppressions, and
  high-signal failures.
- **containment-of-action** — the agent cannot exceed intended authority. Refine
  requirements toward least privilege, approvals, sandboxing, sensitive-resource
  boundaries, consequential external actions requiring human approval, and
  auditability.

Do not add "continuous improvement" as a sub-factor. The quality loop improves the
harness; `continuity` is what lets a work run carry state across time. Observability
also stays as requirement evidence under `self-verifiability`, `continuity`, and
`enforcement-of-standards`, not as a standalone root sub-factor for now, because
otherwise it risks double-counting the evidence surface for those concerns.

Adding `continuity` ripples through the guide's existing count language, which the
implementation reconciles in the same pass: the example umbrella `description` gains
the state-preservation/resume capability; the "improve the harness over time"
avoid-note becomes the _eighth_ candidate beside the seven; and the legacy-recognition
notes (authoring guide and Top 10 check 8) keep reading a pre-existing six-sub-factor
`harnessability` as prior coverage while active authoring adds `continuity`.

### 2. Agent-harness-area template (authoring guide)

Extend [Carry the recurring use-context constituents] so the harness bullet reaches
the same concreteness as the self-check bullet beside it. Add:

- A one-line framing: the harness area is an _enable_ + partly _normative_
  steering-materials constituent — the checked-in guide layer of the broader
  harness — judged as an artifact (is the map accurate, current, coherent, selective,
  discoverable, and maintainable). Runtime harness code, sandboxes, tools, hooks, or
  orchestration that the project owns may deserve their own areas or requirements;
  they are not silently folded into the steering-materials area.
- A candidate, illustrative factor family for the steering materials as an artifact:
  `completeness`, `accuracy`, `currentness`, `understandability`, `coherence`,
  `selectivity`, `discoverability`/`triggerability`, `maintainability`,
  `trustworthiness`, and `assessability`. Presented as a prompt earned per entity,
  never a roster — the same hedge the guide already applies to every factor list.
- Domain-agnostic requirement shapes: a stable minimal entry point orients an agent
  and links deeper without exhausting context; recorded conventions match actual
  practice; steering materials point to _how work is verified_ and _what is
  off-limits_; no two steering documents contradict each other; skill names and
  descriptions trigger the right guidance; executable or third-party guidance is
  trusted only through reviewable provenance; representative traces or feedback logs
  show whether the guidance helps in real work; steering materials are updated when
  structure or workflows change.

The candidate family is deliberately closer to the guide's ~ten-factor aim for a
primary-subject constituent, while still illustrative rather than prescriptive.

### 3. Served-domain guardrail (authoring guide + doctrine guide)

State once in the doctrine guide and reference it from the authoring template
(say-it-once). The guardrail: harness factors/requirements are authored agnostic to
the _served_ domain; where a concrete mechanism is shown (lint/type-check/CI), frame
it as one domain's instance of "how this project verifies work / enforces standards
/ bounds action," which resolves to a link-checker or build (docs), a validation
suite (data), or a review checklist (a judgment domain).

### 4. Self-verifiability sharpening (authoring guide)

In the `self-verifiability` sub-factor's example requirements, name the good-sensor
bar from _Sensors for Coding Agents_: fast; actionable (output carries remediation
guidance and the rule's rationale); grounded in concrete evidence; context-aware;
and suppressible through visible, reviewable exceptions rather than binary — with
deterministic (computational) signals distinguished from inferential (LLM-judge)
ones, behavioral/non-deterministic outcomes covered by inferential signals, and
traces/run logs/evaluation records treated as verification evidence where they expose
what happened and why. Note the cross-link: "suppressible with visible exceptions" meets
`enforcement-of-standards`'s "suppression escapes are constrained" from the other
side — the sensor surfaces the escape, the gate constrains it.

### 5. Use-context-constituent explicitness principle (doctrine guide)

Add a short subsection to
[Modeling quality across domains](../../../docs/guides/model-quality-across-domains.md),
adjacent to _Agentic use context_, that refines the existing "domain agnostic is not
context neutral" line into an actionable rule: explicit, opinionated modeling
guidance is licensed for **use-context constituents** (the agent harness and the
QUALITY.md self-check), never for a **modeled domain**; and even for a use-context
constituent, its factors and requirements stay agnostic to the served domain. Close
with one good/avoid pair:

- Good: "Steering materials point an agent to how the project verifies work."
- Avoid: "Steering materials document the lint, type-check, and test commands."
  (assumes software)

### 6. Top 10 check 8

Two added findings in the area-and-factor-shape check: a germane harness area with
only one or two thin factors is a coverage gap (route to authoring); harness
requirements that assume a software toolchain when the served domain is not software
(route to authoring).

## Alternatives

- **Treat the agent harness as a modeled domain with a catalog entry** in the
  doctrine guide's quality-context catalog (beside documentation, data, service).
  Rejected: the catalog is for _modeled_ domains, and giving the harness a catalog
  slot would imply the harness is a domain QUALITY.md is _about_, privileging the use
  context as a modeled domain — the precise anti-pattern the guide forbids. The
  harness belongs in the use-context register, so it is licensed as a _constituent_,
  not catalogued as a domain.
- **Keep the original six sub-factors unchanged** and only sharpen
  `self-verifiability`. Rejected after fresh research: long-running-agent work makes
  state preservation and handoff a distinct fit-for-purpose concern, not merely a
  detail of accessibility or operability.
- **Add observability as an eighth Agent Harnessability sub-factor.** Rejected for
  now: traces and run evidence are essential, but they are evidence surfaces for
  verification, continuity, and standards enforcement. A standalone observability
  factor would likely double-count unless a project owns a distinct observability
  subsystem as an area.
- **Rename the agent-harness area to agent steering materials.** Rejected for now:
  the existing `agent-harness` term is already established in the guide. The design
  narrows its definition instead — the default area covers checked-in steering
  materials, while broader runtime harness infrastructure may be modeled separately
  when the project owns it.
- **Mandate a fixed harness factor set.** Rejected: contradicts "factors are earned
  per Model" and "prompt, not a roster." The family is illustrative.
- **Put the explicitness principle only in the authoring guide.** Rejected: it is
  doctrine about the agnostic/use-context registers, whose home is the doctrine guide
  (with the `AGENTS.md` summary pointing to it); the authoring guide references it.
- **Also re-scope the README harness material and add a worked harness example.**
  Rejected as overlap: README modeled-domain/harness re-scoping was
  [0088](../0088-domain-agnostic-corpus-alignment.md)'s scope (since landed),
  and a worked harness example is example-corpus work better sequenced with that
  effort. This case stays on guidance and doctrine.

## Trade-offs & risks

- **Guidance growth.** The authoring guide is already long; adding a root
  Agent Harnessability refinement and a fuller harness-area template grows it
  further. Mitigated by keeping both families illustrative, using terse requirement
  shapes, and stating the served-domain guardrail once in the doctrine guide and
  referencing it.
- **Double-counting regression.** A richer harness-area family risks re-absorbing
  Agent Harnessability evidence. Mitigated by restating the 0087 boundary in the
  template and keeping the area's factors about the steering artifact's own quality.
- **Compatibility churn.** Adding `continuity` changes the expected
  Agent Harnessability shape for newly authored or revised models. Mitigated by
  treating existing six-factor `agent-harnessability` models as useful prior
  coverage during review, while recommending the continuity addition during active
  model authoring.
- **Over-fitting to current tools.** The harness-engineering sources are
  software-coding-agent-centric; the guardrail exists precisely to keep their
  vocabulary from re-importing software assumptions into the served domain.

## Open questions

- Whether the `AGENTS.md` summary should carry a one-line version of the
  explicitness principle, or whether the existing "agentic use context" summary plus
  the guide link is enough. Leaning toward a single line, dropped if it bloats.
- Whether to name the routing factor `discoverability`, `triggerability`, or both;
  settle during authoring against the guide's preference for conventional quality
  names and the skill-specific need for metadata to trigger at the right time.
