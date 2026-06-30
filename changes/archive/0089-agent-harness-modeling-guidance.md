---
type: Change Case
title: Agent-harness modeling guidance
description: Strengthen how the /quality guidance models a project's agent harness — refine the root Agent Harnessability factor with continuity and stronger verification/observability requirements, give the agent-harness area a fuller domain-agnostic steering-materials factor family, and add a doctrine principle that use-context constituents may carry explicit guidance while their factors and requirements stay agnostic to the served domain.
status: Done
tags: [docs, doctrine, skill, agent-harness]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-harness modeling guidance

A **Change Case** to fix an asymmetry and a coverage gap in the `/quality`
authoring guidance: the model-wide
[Agent Harnessability](../../skills/quality/guides/authoring.md) factor is specified
as an umbrella over the agent's working loop, while the **agent-harness area** —
the steering-materials constituent the same guidance says to model by default —
gets no factor family or requirement guidance at all. The result, seen in a
freshly generated `QUALITY.md`, is a harness area carried with one or two thin
factors (`completeness`, `currentness`) and a couple of placeholder requirements,
well below the depth the constituent warrants. Fresh harness-engineering research
also shows that long-running work continuity is too central to leave implicit in
the root Agent Harnessability decomposition.

This case closes that gap, grounded in harness-engineering practice, and settles a
doctrinal question the work surfaces: the agent harness is a _use-context_
constituent, so the project may give it explicit, opinionated modeling guidance —
but that guidance must keep the harness's factors and requirements agnostic to the
_served_ domain, so a generated harness area does not assume software engineering.

Detail lives in:

- [Functional spec](0089-agent-harness-modeling-guidance/spec.md) - what the
  guidance must say.
- [Design doc](0089-agent-harness-modeling-guidance/design.md) - the root
  Agent Harnessability refinement, the harness-area factor family, where the new
  doctrine principle lands, the harness-engineering grounding, and the alternatives
  weighed.

## Motivation

The `/quality` skill is meant to produce the fullest, most sufficient model a
project's evidence supports in one setup pass. For an agent-collaborated project it
does this well on the _factor_ side — [Carry Agent Harnessability for
agent-collaborated composite roots](../../skills/quality/guides/authoring.md) gives
the umbrella factor six independently assessable sub-factors with explicit
boundaries. But on the _constituent_ side it is lopsided: [Carry the recurring
use-context constituents](../../skills/quality/guides/authoring.md) tells the author
to model the agent harness _and_ the QUALITY.md self-check by default, then gives
the self-check a concrete shape (key, title, `source`, factor kinds, a governing
requirement) while giving the harness no comparable factor-family or requirement
template — only its normative/dual-role and projection-boundary guidance. An author
following the guide literally has a modeling template for one use-context
constituent and a bare instruction for the other.

That asymmetry shows up downstream. A generated model's `agent-harness` area
carried only `completeness` and `currentness` and two thin requirements — paltry
for a germane primary-subject constituent, which the same guide elsewhere holds to
a much higher factor-depth aim. The harness is not hard to model richly; the
guidance simply does not say how.

Research into harness engineering confirms two things. First, the field's central
framing — Agent = Model + Harness, with the harness split into _guides_
(feedforward, steering before action), _sensors_ (feedback, observing after),
tools, runtime environment, memory/state, orchestration, guardrails, and traces —
maps cleanly onto the model's projection split: the agent-harness area is the
checked-in steering-materials artifact, while the model-wide Agent Harnessability
factor rates whether the project as a whole lets an agent understand, act, resume,
verify, and stay bounded. Second, recent long-running-agent and trace-driven
harness work makes **continuity** a first-class harnessability concern: a project
that cannot preserve state, decisions, remaining work, and verification status
across compaction, handoff, interruption, or long tasks is not fully harnessable,
even if its entry point and checks are good. The "good sensor" properties from
_Sensors for Coding Agents_ still sharpen `self-verifiability`: feedback should be
fast, actionable, grounded in evidence, context-aware, and suppressible through
visible reviewable exceptions, with inferential signals where no deterministic
oracle exists.

Finally, the work surfaces a doctrine gap. [Modeling quality across domains](../../docs/guides/model-quality-across-domains.md) settles that QUALITY.md is
_domain agnostic but not context neutral_ — but it does not say how that resolves
for the harness specifically. The harness is the canonical case where the two
registers meet: it recurs from the agentic **use context** (so explicit guidance is
licensed, exactly as for the self-check), yet it **serves** whatever domain the
project models (so its requirements must not assume software). Without that
principle stated, harness guidance risks either being withheld (treating the
harness as a modeled domain we must stay silent on) or leaking software assumptions
(lint/type-check/CI as if every harnessed project were a codebase).

## Scope

Covered:

- Give the **agent-harness area** a concrete modeling template in the authoring
  guide — a fuller, illustrative-not-mandatory factor family for the
  steering-materials constituent, and requirement guidance — at parity with the
  self-check treatment beside it.
- Refine the root **Agent Harnessability** factor by adding **continuity** as a
  first-class sub-factor and strengthening the existing sub-factor requirements
  around progressive disclosure, tool affordances, verification evidence,
  guardrails, and containment.
- Direct that the harness area's factors and requirements be authored **agnostic to
  the served domain**: phrased around what steering materials do (orient, point to
  verification, bound action), not around a software toolchain.
- Sharpen the **self-verifiability** sub-factor with the "good sensor" properties
  (fast; actionable; grounded in concrete evidence; context-aware; and suppressible
  through visible, reviewable exceptions rather than binary; behavioral outcomes get
  runnable evals) and observability/trace evidence.
- Add to the doctrine guide a principle distinguishing **use-context constituents**
  (the agent harness and the QUALITY.md self-check — explicit guidance is licensed)
  from **modeled domains** (never privileged), plus the **served-domain guardrail**
  and a good/avoid pair.
- Extend Top 10 check 8 to flag a germane harness area carried with only one or two
  thin factors, and harness requirements that assume a software toolchain when the
  served domain is not software.
- Mirror the authoring-guide and Top 10 changes into their durable spec mirrors, and
  record the work in the relevant logs and the changelog.

Deferred / non-goals:

- No QUALITY.md schema, rating, roll-up, or evaluation-semantic change.
- No CLI or Go behavior change.
- No collapse of Agent Harnessability and the agent-harness area into the same
  factor family; they remain distinct projections with an explicit boundary.
- No re-opening of [0087](0087-encode-projection-boundaries.md)'s
  projection-boundary rule (Agent Harnessability factor vs. agent-harness area);
  this builds on it and keeps the boundary intact.
- No README modeled-domain or harness re-scoping — that landed with
  [0088](0088-domain-agnostic-corpus-alignment.md); this case does not touch
  the README.
- No new worked harness example in the example corpus (deferred; the related
  example-corpus work landed with 0088, but a harness example is out of scope here).

## Affected artifacts

Derived by analysis: a sweep for where harness modeling guidance lives and is
mirrored — the bundled `skills/quality/guides/` authoring and Top 10 guides, their
durable mirrors under `specs/skills/quality-skill/guides/`, the doctrine guide in
`docs/guides/`, and the `AGENTS.md` summary that points to it. Grouped by kind;
empty kinds are deliberate.

### Code

None - documentation, doctrine, and bundled-skill guidance content only.

### Format spec (`SPECIFICATION.md`)

None - no normative format rule changes; the harness factor family stays
illustrative, never a schema default.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0089-agent-harness-modeling-guidance/spec.md#durable-spec-changes)
section is authoritative. In summary:

- [x] `specs/skills/quality-skill/guides/authoring-md.md` - mirror the
      Agent Harnessability continuity refinement (including the reconciled
      six-/seventh-sub-factor count references and the updated umbrella
      `description`), the harness-area factor/requirement template, the
      served-domain guardrail, and the self-verifiability sensor/observability
      sharpening.
- [x] `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - mirror the
      check-8 additions and the reconciled legacy six-sub-factor recognition note.
- [x] `specs/skills/quality-skill/guides/log.md` - record the mirror updates.

### Durable docs

- [x] `docs/guides/model-quality-across-domains.md` - add the use-context-constituent
      explicitness principle and the served-domain guardrail.
- [x] `AGENTS.md` - optional one-line summary of the new principle in the
      "Quality-domain agnostic examples and agentic use context" section (the single
      source for the gitignored `CLAUDE.md`/`GEMINI.md` symlinks). Drop if it bloats.
- [x] `docs/log.md` - record the doctrine-guide edit.
- [x] `CHANGELOG.md` - a documentation/guidance release note.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/guides/authoring.md` - the Agent Harnessability continuity
      refinement (reconciling the existing six-/seventh-sub-factor count references
      and updating the umbrella factor `description`), the agent-harness-area
      template, the served-domain guardrail, and the self-verifiability
      sensor/observability sharpening.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` - the check-8 additions and
      the reconciled legacy six-sub-factor recognition note.

### Install / scaffold

None - no scaffolded QUALITY.md content changes.

## Children

- [Functional spec](0089-agent-harness-modeling-guidance/spec.md) - what the
  guidance must say.
- [Design doc](0089-agent-harness-modeling-guidance/design.md) - the root
  Agent Harnessability refinement, the harness-area factor family, placement,
  harness-engineering grounding, and alternatives.

## Status

`Done`. Implemented and archived. The bundled authoring guide now adds
`continuity` to Agent Harnessability, strengthens the surrounding sub-factor
requirements, gives the agent-harness area a domain-agnostic steering-materials
template, and keeps the projection boundary explicit. The Top 10 checks now flag
thinly factored and software-leaking harness areas while preserving legacy
six-sub-factor coverage recognition. The durable spec mirrors, doctrine guide,
`AGENTS.md`, logs, and CHANGELOG are updated. No code, CLI, schema, rating, roll-up,
or `SPECIFICATION.md` normative change was required.
