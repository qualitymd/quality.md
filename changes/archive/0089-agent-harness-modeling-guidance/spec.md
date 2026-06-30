---
type: Functional Specification
title: Agent-harness modeling guidance - functional spec
description: Requirements for refining the root Agent Harnessability factor with continuity and stronger verification/observability guidance, giving the agent-harness area a fuller domain-agnostic steering-materials factor and requirement template, extending Top 10 check 8, and adding a use-context-constituent explicitness principle with a served-domain guardrail to the doctrine guide.
tags: [docs, doctrine, skill, agent-harness]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-harness modeling guidance - functional spec

Companion to
[Agent-harness modeling guidance](../0089-agent-harness-modeling-guidance.md).
This spec states what the guidance must say. The doctrine it extends is settled in
[Modeling quality across domains](../../../docs/guides/model-quality-across-domains.md);
the format itself is governed by
[`SPECIFICATION.md`](../../../SPECIFICATION.md). This spec changes guidance prose only —
it adds no normative format rule and no schema default. The
[design doc](design.md) settles _how_ the root Agent Harnessability refinement, the
harness-area template, and the doctrine principle are shaped.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The authoring guide specifies the model-wide Agent Harnessability factor richly
(six sub-factors, each with a boundary and example requirements) but gives the
agent-harness _area_ — the steering-materials constituent it tells authors to model
by default — no factor family or requirement guidance, even though it gives the
sibling QUALITY.md self-check constituent a concrete shape. A model generated from
that guidance carried a harness area with one or two thin factors, below the depth a
germane primary-subject constituent warrants.

Harness-engineering practice (see the [design doc](design.md) for sources) confirms
the projection split: Agent Harnessability is the model-wide quality of being
amenable to useful, bounded agent work, while the agent-harness area is the
checked-in steering-materials artifact. Fresh long-running-agent, context, trace,
and skills guidance also shows that the current six-sub-factor decomposition leaves
one fit-for-purpose concern implicit: continuity across compaction, handoff,
interruption, and long tasks. The work also surfaces a doctrine gap: the harness is
the canonical case where _domain agnostic but not context neutral_ must be resolved —
it recurs from the agentic use context (explicit guidance licensed) yet serves
whatever domain the project models (its requirements must not assume software).

## Scope

Covered: a root Agent Harnessability refinement in the authoring guide; an
agent-harness-area modeling template; a served-domain guardrail for harness factors
and requirements; a self-verifiability sensor-quality and observability sharpening;
a use-context-constituent explicitness principle in the doctrine guide; Top 10
check-8 additions; and the durable mirrors and logs for the above.

Deferred / non-goals: no QUALITY.md schema, rating, roll-up, or evaluation-semantic
change; no CLI or Go change; no collapse of Agent Harnessability and the
agent-harness area into one factor family; no re-opening of the 0087
projection-boundary rule; no README change (0088's scope); no new worked harness
example in the example corpus.

## Requirements

### Refine the root Agent Harnessability factor

The authoring guide **MUST** keep Agent Harnessability as the model-wide factor
projection of the agent-collaboration concern and **MUST** keep it distinct from the
agent-harness area. Agent Harnessability rates whether the project as a whole lets an
agent understand the work, receive a bounded assignment, operate the environment,
preserve and resume state, verify output, respect standards, and stay within
permitted authority; the agent-harness area rates the steering materials' own
quality.

The guide **MUST** revise the Agent Harnessability decomposition to include
`continuity` as a first-class sub-factor alongside the existing concerns:
`agent-accessibility`, `task-specifiability`, `agent-operability`, `continuity`,
`self-verifiability`, `enforcement-of-standards`, and `containment-of-action`.
Each sub-factor **MUST** remain independently assessable and carry boundary guidance
that prevents double-counting with sibling sub-factors and with the agent-harness
area.

The example umbrella factor's `description` **MUST** be updated to name the
state-preservation and resumption capability that `continuity` adds, so the umbrella
definition reflects all seven sub-factors rather than enumerating only the original
six capabilities.

The `continuity` sub-factor **MUST** cover the degree to which an agent can preserve
state and resume useful work across long-running tasks, compaction, interruption,
handoff, and fresh sessions. Example requirements **SHOULD** include that handoff or
progress artifacts capture current state, decisions made, remaining work,
verification status, blockers, and next steps; that resumptions do not depend on
unrecoverable chat history; and that progress records reduce false completion or
context-anxiety failure modes. The `continuity` boundary guidance **MUST** distinguish
it from its two nearest sub-factors so it does not double-count: `agent-operability`
covers a fresh session reaching a ready-to-work environment, while `continuity`
covers resuming with prior state, decisions, and progress; and `agent-accessibility`
covers durable, reachable knowledge, while `continuity` covers the handoff and
progress record of an in-flight task.

The guide **MUST** strengthen the existing Agent Harnessability sub-factor
requirements without changing the projection boundary:

- `agent-accessibility` **SHOULD** cover progressive disclosure and context
  selectivity, not only raw availability of information.
- `task-specifiability` **SHOULD** cover decomposition, explicit done criteria, and
  completion discipline that compares the result with the original task before
  declaring success.
- `agent-operability` **SHOULD** cover tool affordance quality, fresh-session
  environment readiness, act/observe loops, and agent-useful tool output.
- `self-verifiability` **SHOULD** cover deterministic checks, inferential evals,
  trace/run evidence, and actionable machine-readable feedback.
- `enforcement-of-standards` **SHOULD** cover input, output, and tool guardrails or
  equivalent domain-neutral controls, not only advisory prose.
- `containment-of-action` **SHOULD** cover least privilege, approval gates,
  sandboxing, sensitive-resource boundaries, and auditability.

The guide **MUST NOT** add "continuous improvement" or "harness evolution" as a
sibling Agent Harnessability sub-factor. Improvement remains part of the quality
loop; observability and trace evidence are requirements that make improvement
possible under `self-verifiability`, `enforcement-of-standards`, and `continuity`.

> Rationale: the original six sub-factors covered the agent's working loop well, but
> long-running-agent practice makes state preservation and handoff too central to
> leave implicit under accessibility or operability. A project that cannot resume
> with the right state, decisions, remaining work, and verification status is not
> fully harnessable. - 0089

The guide **MUST** reconcile every existing reference that hard-codes the old
sub-factor count so the seven-sub-factor decomposition reads consistently. The
"improve the harness over time" guidance is now the _eighth_ candidate beside the
seven, not the seventh beside six; and the legacy-recognition guidance — which today
accepts an existing `harnessability` factor with the "six-sub-factor shape" as
semantic coverage — **MUST** keep reading a pre-existing six-sub-factor model as prior
coverage while directing active authoring or revision to add `continuity`.

### Give the agent-harness area a modeling template

The authoring guide **MUST** give the agent-harness area a concrete modeling
template at parity with the QUALITY.md self-check template beside it, where today it
gives only a bare "model it by default" instruction. The template **MUST** identify
the harness area as an _enable_ and partly _normative_ constituent for the project's
checked-in steering materials — agent entry points, agent guidance files, skills,
prompts, and related instructions that orient and govern an agent — and **MUST**
offer a candidate factor family richer than one or two factors.

> Rationale: the guide specified the umbrella _factor_ and the _self-check_
> constituent concretely but left the _harness_ constituent a bare instruction; an
> author following the guide literally had a template for one use-context
> constituent and a placeholder for the other, and generated models came out thinly
> factored as a result. - 0089

The candidate factor family **MUST** be presented as an illustrative, non-exhaustive
prompt that an author earns for the specific entity — consistent with the guide's
existing "prompt, not a roster" and "factors are earned per Model" rules — and
**MUST NOT** be presented as a required or default factor set.

The candidate family **SHOULD** include, or explicitly cover through nearby names,
these fit-for-purpose qualities of steering materials:

- `completeness` - the steering surface covers entry points, workflow map,
  verification, permissions, escalation, known pitfalls, handoff expectations, and
  update rules that are germane for this project.
- `accuracy` - guidance matches the actual project, tools, workflows, and authority
  boundaries.
- `currentness` - guidance is updated when structure, tools, workflows, standards,
  or use-context assumptions change.
- `understandability` - instructions are clear, scoped, operational, and written at
  the right altitude for agent use.
- `coherence` - steering files, skills, prompts, durable docs, specs, and checks do
  not contradict each other.
- `selectivity` - the guidance is context-efficient and layered; the first entry
  point is a map, not an encyclopedia, and deeper material is loaded only when
  relevant.
- `discoverability` or `triggerability` - an agent can find the right guidance,
  skill, command, or reference when the task calls for it, including through useful
  names and descriptions.
- `maintainability` - the steering materials have clear ownership, update triggers,
  low duplication, and low-friction edit paths.
- `trustworthiness` - when the harness includes skills, executable helpers,
  third-party prompts, or external instructions, their provenance and safety are
  reviewable and untrusted instructions are not silently adopted.
- `assessability` - the project has evidence that the steering materials help
  agents in representative work, such as traces, feedback logs, review findings, or
  task examples.

The template **MUST** keep the agent-harness area distinct from the model-wide Agent
Harnessability factor per the existing projection-boundary rule: Agent
Harnessability rates how each constituent equips an agent, while the agent-harness
area rates the steering materials' own quality (whether the map is accurate,
current, and a map rather than a manual). It **MUST NOT** reintroduce double-counting
between the two projections.

> Rationale: 0087 encoded that boundary; a richer harness-area family must deepen
> the area's own concerns without re-assessing Agent Harnessability evidence under
> the area. - 0089

### Keep harness factors and requirements agnostic to the served domain

The authoring guide **MUST** direct that the harness area's factors and requirements
be authored agnostic to the domain the project models — phrased around what steering
materials do (orient an agent, point to how work is verified, bound what actions are
permitted) rather than around any one domain's toolchain.

The guide **MUST NOT** present software-specific mechanisms (for example lint,
type-check, test, or CI/deploy gates) as the harness's requirements; where it needs
a concrete illustration it **MUST** frame such mechanisms as one domain's instance
of a domain-neutral expectation (how _this_ project verifies work, enforces
standards, or bounds action), which resolves differently for a documentation,
data, or service project.

To keep the guardrail stated once, the authoring guide **SHOULD** carry it by
reference to the doctrine guide's served-domain principle (_Add the
use-context-constituent explicitness principle_ below) rather than restating it in
full.

> Rationale: the agent harness of a non-software project is as real as a codebase's;
> a harness requirement that assumes a software toolchain leaks the use context's
> familiar domain into the served domain, which is the exact agnosticism failure the
> doctrine guide exists to prevent. - 0089

### Sharpen self-verifiability with good-sensor and trace properties

The authoring guide's `self-verifiability` sub-factor **MUST** name the properties of
a good verification signal — that it is fast; actionable (its output carries
remediation guidance and the rationale for the rule, not just a pass/fail);
grounded in concrete evidence; context-aware (rules fit where they apply); and that
its escapes are suppressible through visible, reviewable exceptions rather than
binary — as the bar its example requirements test, rather than only naming that some
pass/fail signal exists. It **MUST** distinguish deterministic (computational)
signals from inferential (LLM-judge / eval) signals and note that behavioral or
non-deterministic outcomes, which lack a deterministic oracle, are covered by
inferential signals. It **SHOULD** also treat traces, run logs, and evaluation
records as verification evidence when they let an agent or reviewer inspect what
happened, diagnose the failure mode, and compare the result against the original
task. Because the "suppressible through visible, reviewable exceptions" property
meets `enforcement-of-standards`'s "suppression escapes are constrained" from the
other side, the guide **MUST** keep that boundary explicit — the sensor surfaces the
escape, the gate constrains it — so the shared suppression evidence is not
double-counted.

> Rationale: harness-engineering practice treats feedback signals (sensors) as the
> central reliability lever and characterizes a _good_ sensor by these properties;
> the sub-factor already gestured at them but did not make them the assessable bar. -
> 0089

### Add the use-context-constituent explicitness principle

[Modeling quality across domains](../../../docs/guides/model-quality-across-domains.md)
**MUST** state a principle distinguishing two cases for how explicit project
guidance may be about a quality context:

- For a **use-context constituent** — the agent harness and the QUALITY.md
  self-check, which recur from the assumed agentic use context rather than from any
  modeled domain — the project **MAY** give explicit, opinionated modeling guidance
  (canned factor families and requirement shapes), on the same license under which
  the project already prescribes the self-check's shape.
- For a **modeled domain** — what a given `QUALITY.md` evaluates — the project
  **MUST NOT** privilege one domain with a default factor roster; factors stay
  earned per Model.

> Rationale: without this distinction, harness guidance is caught between two errors —
> withheld (as if the harness were a modeled domain we must stay silent on) or
> over-generalized into a domain default. Naming the harness and self-check as
> use-context constituents licenses explicit guidance for exactly those two without
> weakening domain agnosticism. - 0089

The principle **MUST** carry the served-domain guardrail: even for a use-context
constituent, its factors and requirements are authored agnostic to the served
domain, so explicit harness or self-check guidance never assumes software. The guide
**MUST** define the _served domain_ on first use as the domain the project's model is
about — the domain a use-context constituent serves — distinct from the agentic use
context itself. The guide **SHOULD** include a good/avoid example pair illustrating a
domain-neutral harness requirement versus a software-leaking one.

The principle **MUST** preserve the guide's existing register discipline (modeled
domain agnostic; agentic use context opinionated) and the existing decision test;
it refines them for the use-context-constituent case rather than replacing them.

### Extend Top 10 check 8

The Top 10 quality-md checks **MUST**, in the area-and-factor-shape check (check 8),
add a finding for an agent-collaborated composite root whose germane agent-harness
area is carried with only one or two thin factors — below the factor-depth aim for a
primary-subject constituent — routing to authoring. A thin harness area **MUST** be
treated as a coverage gap, not as evidence the harness is unimportant.

The check **MUST** also add a finding for harness-area requirements that assume a
software toolchain when the project's served domain is not software, routing to
authoring.

The check's existing legacy-recognition note **MUST** be reconciled with the
seven-sub-factor decomposition: it currently treats an existing `harnessability`
factor with the "expected six sub-factors" as semantic coverage, and **MUST**
continue to accept a pre-existing six-sub-factor model as semantic coverage without
implying six is the current target shape.

> Rationale: the checklist already flags a missing or boundary-less harness but not a
> harness that is present yet thinly factored, or one whose requirements leak the
> software domain — the two failure modes this case's guidance is meant to prevent. -
> 0089

### Record the work

The relevant `specs/` and `docs/` logs and the `CHANGELOG.md` **MUST** record the
guidance and doctrine edits before the case reaches `In-Review`.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` - mirror the root Agent
  Harnessability continuity refinement (including the reconciled sub-factor count
  references and the updated umbrella `description`), the agent-harness-area modeling
  template, the served-domain guardrail, and the self-verifiability
  sensor/trace-quality sharpening (per _Refine the root Agent Harnessability factor_,
  _Give the agent-harness area a modeling template_, _Keep harness factors and
  requirements agnostic to the served domain_, and _Sharpen self-verifiability with
  good-sensor and trace properties_ above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - mirror the
  check-8 additions and the reconciled legacy six-sub-factor recognition note (per
  _Extend Top 10 check 8_ above).

### To rename

None

### To delete

None
