---
type: Functional Specification
title: Agent-harness holistic definition - functional spec
description: Requirements for redefining the agent harness in the /quality authoring guidance as the whole engineered system around the model (feedforward guides plus feedback sensors), deriving the agent-harness area as one honest projection with explicit boundaries, adding a scoping decision rule and never-define-as-instructions doctrine, fixing the same narrow definition in model-structure and setup, extending Top 10 check 8, and mirroring into durable specs.
tags: [docs, doctrine, skill, agent-harness]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-harness holistic definition - functional spec

Companion to
[Agent-harness holistic definition](../0091-agent-harness-holistic-definition.md).
This spec states what the guidance must say. It builds on the projection boundary
settled in [0087](../../archive/0087-encode-projection-boundaries.md) and the
harness-area modeling template from
[0089](../../archive/0089-agent-harness-modeling-guidance.md); it does not reopen
either. The format itself is governed by
[`SPECIFICATION.md`](../../../SPECIFICATION.md). This spec changes guidance prose
only — it adds no normative format rule and no schema default. The
[design doc](design.md) settles _how_ the holistic definition, the three-projection
model, and the scoping rule are shaped, and records the redefine-not-rename
decision.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

0089 grounded the projection split correctly but left the constituent's own
definition narrow. The authoring guide opens "The agent harness is the
instructions that steer the agent working the project," and the same narrow
definition recurs in `model-structure.md` and `setup.md`. That is a part-for-whole
error: the field defines the harness as everything that isn't the model — code,
configuration, and execution logic — split into feedforward controls that steer
before action and feedback controls that catch and correct after. Instructions are
one feedforward component. An instructions-only definition silently drops the
feedback half and the project-owned runtime controls, and produces the visible
confusion of an "Agent Harness" area whose name claims the system but whose
definition delivers one slice.

## Scope

Covered: a holistic harness definition in the authoring guide; a derived
area-scoping decision rule with explicit boundaries against the factor, tests, and
operations; project-owned runtime machinery brought into the area's scope;
expanded requirement shapes spanning both halves; a never-define-as-instructions
doctrine; the matching definition fix in `model-structure.md` and `setup.md`; a
Top 10 check-8 addition; and the durable mirrors and logs for the above.

Deferred / non-goals: no rename of the agent-harness area, key, or title; no
collapse of Agent Harnessability and the agent-harness area; no change to the
Agent Harnessability factor's own (already holistic) definition beyond
cross-references; no QUALITY.md schema, rating, roll-up, or evaluation-semantic
change; no CLI or Go change; no `SPECIFICATION.md` normative change; no new worked
harness example.

## Requirements

### Define the agent harness holistically

The authoring guide
[`agent-harness.md`](../../../skills/quality/guides/authoring/agent-harness.md)
**MUST** open by defining the agent harness as the whole engineered system around
the model — everything that isn't the model itself: the code, configuration, and
execution logic that turns a model into an agent working the project. The
definition **MUST** name both halves of the harness: feedforward controls that
steer the agent before it acts (agent entry points, guidance files, skills,
prompts, tool and MCP definitions, sandbox and filesystem, orchestration) and
feedback controls that catch and correct it after (verification commands, tests,
run logs, evals, review).

The guide **MUST** state that instructions are one feedforward component of the
harness, not the whole of it. The guide **MUST NOT** define the agent harness as
"the instructions that steer the agent" or any equivalent that equates the harness
with its steering-prose slice.

> Rationale: the field's framing — Agent = Model + Harness, the harness being
> everything but the model, split into feedforward guides and feedback sensors —
> is the established definition; an instructions-only definition drops the entire
> feedback half and the project-owned runtime controls. - 0091

### Derive the area as one projection, with explicit boundaries

The guide **MUST** state that a composite root does not carry one "harness" area
that owns the whole system, and that the harness is rated through three
projections, each part of the harness belonging to exactly one:

- the _capability_ of equipping an agent — the model-wide **Agent Harnessability**
  factor;
- the _checked-in, project-owned governing artifacts_ — the **agent-harness**
  area;
- the _verification corpus_ and the _runtime environment_ — the tests and
  operations constituents respectively, cross-referenced rather than absorbed.

The guide **MUST** keep the agent-harness area distinct from the Agent
Harnessability factor per the existing projection-boundary rule and **MUST NOT**
reintroduce double-counting: harness files are evidence for the factor and the
evaluated entity for the area.

### Scope the area to the project-owned governing artifacts

The guide **MUST** scope the agent-harness area to the harness artifacts the
project checks in and owns — agent entry points, guidance files, skills, and
prompts, together with the project-owned hooks, tool/MCP definitions,
sandbox/permission policy, and orchestration config where these exist — rather
than to prose instructions alone.

The guide **MUST** treat project-owned runtime harness machinery (hooks, tool/MCP
definitions, sandbox/permission policy, orchestration config) as in-scope for the
area, surfaced here or given its own area when large enough to warrant distinct
factors. It **MUST NOT** be silently folded into prose instructions and **MUST
NOT** be silently dropped.

The guide **MUST** give a scoping **decision rule** for an artifact that could
belong to more than one projection: rate it in the agent-harness area when its
primary job is to govern or equip the agent's work and the project owns it; cede
it to a domain constituent when it is primarily a product artifact the agent
merely also uses (the product test suite → the tests area; the deploy runtime →
the operations area). When one artifact does both, the guide **MUST** direct that
its agent-governing quality be rated in the area and the domain constituent
cross-referenced, under the no-double-count rule.

> Rationale: the harness genuinely spans verification and environment, but those
> are full constituents in their own right; the decision rule keeps the area
> honest about its breadth without re-rating what a sibling owns. - 0091

### Expand the requirement shapes across both halves

The guide **MUST** phrase harness-area requirement shapes to span both halves of
the harness and the owned controls, not prose instructions alone, while keeping
them agnostic to the served domain (preserving 0089's served-domain guardrail).
The shapes **SHOULD** cover, illustratively and non-exhaustively:

- feedforward: a stable minimal entry point orients an agent and links deeper
  without exhausting context; recorded conventions match actual practice; skill
  and command names and descriptions trigger the right guidance; steering
  documents do not contradict each other or the guides they reference; executable
  or third-party guidance has reviewable provenance;
- feedback: the harness points to how work is verified and routes to signals the
  agent can run or inspect; representative traces or feedback logs show whether the
  guidance helps in real work;
- owned controls: project-owned hooks, sandbox, or permission policy that bound
  consequential action are coherent, current, and inspectable; orchestration or
  subagent config is internally consistent.

The guide **MUST** direct that requirements rate the artifact's own quality, not
the capability it confers — "is the permission policy coherent and current"
belongs to the area, while "does it contain the agent" is the
`containment-of-action` sub-factor — so the projection boundary holds in the
requirement phrasing.

### Add the never-define-as-instructions doctrine

The guide **MUST** carry an explicit doctrine principle: define the agent harness
as the whole engineered system around the model (feedforward and feedback
controls), then scope the area to the checked-in, project-owned slice; never
define the harness as "the instructions," because tools, sandbox, orchestration,
and verification are equally the harness.

### Fix the matching definition in model-structure and setup

[`model-structure.md`](../../../skills/quality/guides/authoring/model-structure.md)'s
"Carry the recurring use-context constituents" section **MUST** stop defining the
agent harness as "the instructions that steer the agent" and **MUST** refer to it
as the harness system (or the project-owned harness artifacts), consistent with
`agent-harness.md`, while keeping its routing to the agent-harness guide intact.

[`setup.md`](../../../skills/quality/workflows/setup.md) **MUST** align the
generated agent-harness area `description` and its projection-boundary comment with
the holistic definition, so a generated model does not describe the area as
instructions-only. Setup **MUST** add an active check during constituent modeling
for project-owned runtime harness machinery (hooks, tool/MCP definitions,
sandbox/permission policy, orchestration config), so such machinery is surfaced in
the area or as its own area rather than silently dropped. Setup **SHOULD** include
in its user-facing recap a one-line statement of how the harness was scoped — that
the agent-harness area is the checked-in steering and owned-control artifacts,
while the broader equipping capability is the Agent Harnessability factor.

> Rationale: the user-visible confusion arose after setup, reading the generated
> model; surfacing the scope in the recap and the generated description pre-empts
> the "why is the harness only instructions?" question. - 0091

### Extend Top 10 check 8

The Top 10 quality-md checks **MUST**, in the area-and-factor-shape check (check
8), add a finding for an agent-harness area defined or scoped as instructions only
— missing the feedback half or project-owned runtime controls the harness
includes — routing to authoring. The check **MUST** also flag project-owned
runtime harness machinery that is present in the repo but unmodeled (neither in
the agent-harness area nor its own area nor explicitly out of scope).

### Keep harness factors and requirements agnostic to the served domain

The guide **MUST** preserve 0089's served-domain guardrail: the expanded
requirement shapes stay agnostic to the domain the project models, and the guide
**MUST NOT** present software-specific mechanisms (lint, type-check, test, CI) as
the harness's requirements except as one domain's instance of a domain-neutral
expectation.

### Record the work

The relevant `specs/` and `docs/` logs and the `CHANGELOG.md` **MUST** record the
guidance edits before the case reaches `In-Review`.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/guides/authoring/agent-harness.md` - mirror the
  holistic definition, the three-projection derivation, the area-scoping decision
  rule, the owned-runtime-machinery scope, the expanded requirement shapes, and the
  never-define-as-instructions doctrine (per _Define the agent harness
  holistically_, _Derive the area as one projection_, _Scope the area to the
  project-owned governing artifacts_, _Expand the requirement shapes across both
  halves_, and _Add the never-define-as-instructions doctrine_ above).
- `specs/skills/quality-skill/guides/authoring/agent-harnessability.md` - mirror
  the updated cross-reference that distinguishes Agent Harnessability from the
  agent-harness Area as checked-in, project-owned governing artifacts, not the
  whole equipping capability or instructions alone (per _Derive the area as one
  projection, with explicit boundaries_ above).
- `specs/skills/quality-skill/guides/authoring/model-structure.md` - mirror the
  "Carry the recurring use-context constituents" definition fix (per _Fix the
  matching definition in model-structure and setup_ above).
- `specs/skills/quality-skill/workflows/setup.md` - mirror the generated area
  `description`, projection-boundary comment, recap line, and owned-machinery check
  (per _Fix the matching definition in model-structure and setup_ above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - mirror the
  check-8 additions (per _Extend Top 10 check 8_ above).

### To rename

None — the `agent-harness` area, key, and title are deliberately unchanged; this
case redefines, it does not rename.

### To delete

None
