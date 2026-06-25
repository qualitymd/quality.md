---
type: Change Case
title: Agent-harness holistic definition
description: Re-think the agent-harness authoring guidance so it defines the agent harness as the whole engineered system around the model (feedforward guides plus feedback sensors, everything-but-the-model) and derives the agent-harness area as one honest projection of it, replacing the part-for-whole "the harness is the steering instructions" definition.
status: Done
tags: [docs, doctrine, skill, agent-harness]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-harness holistic definition

A **Change Case** to correct a part-for-whole error in the `/quality` authoring
guidance. The guidance defines the **agent harness** as "the instructions that
steer the agent" — but the agent harness is the whole engineered system around
the model. This case redefines the harness holistically and re-derives the
agent-harness area as one honest projection of it, without renaming the
constituent or re-opening the projection boundary that
[0089](0089-agent-harness-modeling-guidance.md) and
[0087](0087-encode-projection-boundaries.md) settled.

Detail lives in:

- [Functional spec](0091-agent-harness-holistic-definition/spec.md) — what the
  guidance must say.
- [Design doc](0091-agent-harness-holistic-definition/design.md) — the
  holistic definition, the three-projection model, the scoping decision rule, the
  redefine-not-rename decision, and the alternatives weighed.

## Motivation

[0089](0089-agent-harness-modeling-guidance.md) grounded the *projection
split* correctly — the model-wide Agent Harnessability factor rates how the
project equips an agent, while the agent-harness area rates a checked-in artifact
— but it left the constituent's own **definition** narrow. The authoring guide
[`agent-harness.md`](../../skills/quality/guides/authoring/agent-harness.md) opens:
"The agent harness is the instructions that steer the agent working the project."
The same narrow definition appears in
[`model-structure.md`](../../skills/quality/guides/authoring/model-structure.md)
("Carry the recurring use-context constituents") and flows into `setup.md`'s
generated area description and projection-boundary comment.

That definition is a part-for-whole error. The field's framing — Agent = Model +
Harness — defines the harness as **everything that isn't the model itself**: the
code, configuration, and execution logic that turns a model into an agent, split
into **feedforward controls** that steer before action (entry points, guides,
skills, prompts, tool/MCP definitions, sandbox, orchestration) and **feedback
controls** that catch and correct after (verification, tests, logs, evals,
review). Instructions are *one feedforward component*, not the whole harness. An
instructions-only definition silently drops the entire feedback half and the
project-owned runtime controls (hooks, sandbox/permission policy, orchestration
config).

The error has a visible downstream cost. A reader who opens a generated model,
sees an "Agent Harness" area scoped to CLAUDE.md-style steering files, and asks
"why is the harness only instructions?" is reading the guidance correctly — the
name claims the whole system while the definition delivers one slice. The fix is
not to narrow the *name* (rename considered and rejected, see design) but to widen
the *definition*: define the harness at full breadth, then scope the area as the
checked-in, project-owned **governing-artifacts** projection, with explicit
boundaries against the factor (capability), the tests area (the feedback corpus),
and the operations area (the runtime environment) so the breadth does not
double-count.

This builds directly on 0089's projection boundary; it does not reopen it. It
corrects the one thing 0089 left narrow: the definition of the harness itself.

## Scope

Covered:

- Re-think [`agent-harness.md`](../../skills/quality/guides/authoring/agent-harness.md)
  so it **defines the harness holistically first** (feedforward + feedback /
  everything-but-the-model), then derives the area as one projection of it.
- Add a **scoping decision rule**: rate an artifact in the agent-harness area
  when its primary job is to govern or equip the agent's work and the project owns
  it; cede it to a domain constituent when it is primarily a product artifact the
  agent merely also uses (product test suite → tests; deploy runtime →
  operations); rate the agent-governing quality here and cross-reference when one
  artifact does both.
- Make **project-owned runtime harness machinery** (hooks, tool/MCP definitions,
  sandbox/permission policy, orchestration config) explicitly in-scope for the
  area, surfaced or given its own area — never silently folded into prose
  instructions and never silently dropped.
- **Expand requirement shapes** to span both halves of the harness and the owned
  controls, not only prose instructions, while keeping them agnostic to the
  served domain.
- Add doctrine: **never define the harness as "the instructions"** — define it as
  the whole system, then scope the area.
- Fix the same narrow definition in
  [`model-structure.md`](../../skills/quality/guides/authoring/model-structure.md)
  "Carry the recurring use-context constituents", and align `setup.md`'s generated
  area `description`, projection-boundary comment, and user-facing recap.
- Extend Top 10 check 8 to flag an instructions-only harness definition and
  unmodeled project-owned runtime harness machinery.
- Mirror the guide and Top 10 changes into their durable spec mirrors, and record
  the work in the relevant logs and the changelog.

Deferred / non-goals:

- **No rename** of the agent-harness area, its `agent-harness` key, or title.
  "Agent harness" is the correct term for the system; this case fixes the
  *definition*, not the name. (Rename considered and rejected — see design.)
- **No collapse** of the Agent Harnessability factor and the agent-harness area;
  they remain distinct projections with the boundary
  [0087](0087-encode-projection-boundaries.md)/0089 set intact.
- No change to the **Agent Harnessability factor's own definition**, which is
  already holistic; touch it only for a cross-reference if needed.
- No QUALITY.md schema, rating, roll-up, or evaluation-semantic change; no CLI or
  Go behavior change; no `SPECIFICATION.md` normative change.
- No new worked harness example in the example corpus.

## Affected artifacts

Derived by analysis: a sweep for where the harness definition lives and is
mirrored — `grep` for "instructions that steer" / "steering materials" and the
agent-harness area concept across `skills/`, `specs/`, `docs/`, and `AGENTS.md`.
Grouped by kind; empty kinds are deliberate. Checkboxes are reconciled before
In-Review.

### Code

None — documentation, doctrine, and bundled-skill guidance content only.

### Format spec (`SPECIFICATION.md`)

None — no normative format rule changes; the harness factor family and the area
projection stay illustrative, never a schema default.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0091-agent-harness-holistic-definition/spec.md#durable-spec-changes)
section is authoritative. In summary:

- [x] `specs/skills/quality-skill/guides/authoring/agent-harness.md` — mirror the
      holistic-definition rewrite, the scoping decision rule, the owned-runtime
      machinery scope, the expanded requirement shapes, and the
      never-define-as-instructions doctrine.
- [x] `specs/skills/quality-skill/guides/authoring/agent-harnessability.md` —
      mirror the updated boundary cross-reference to the checked-in,
      project-owned governing-artifacts projection.
- [x] `specs/skills/quality-skill/guides/authoring/model-structure.md` — mirror
      the "Carry the recurring use-context constituents" definition fix.
- [x] `specs/skills/quality-skill/workflows/setup.md` — mirror the generated area
      `description`, projection-boundary comment, recap line, and owned-machinery
      check.
- [x] `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — mirror
      the check-8 additions.
- [x] `specs/skills/quality-skill/guides/log.md` — record the mirror updates.
- [x] `specs/skills/quality-skill/workflows/log.md` — record the setup workflow
      mirror update.

### Durable docs

- [x] `docs/guides/model-quality-across-domains.md` — reflect the holistic
      framing in the "Agentic use context" section only if its harness wording
      restates the narrow definition; light touch or `None` after review.
- [x] `AGENTS.md` — reviewed; no edit needed because the existing summary already
      delegates the authoritative use-context rules to the domain-agnostic guide.
- [x] `docs/log.md` — record the doctrine-guide edit if `docs/` is touched.
- [x] `CHANGELOG.md` — a documentation/guidance release note.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/guides/authoring/agent-harness.md` — the primary holistic
      rewrite (concept-first definition, three-projection model, scoping rule,
      owned-machinery scope, expanded requirement shapes, doctrine).
- [x] `skills/quality/guides/authoring/model-structure.md` — the recurring
      use-context-constituents definition fix.
- [x] `skills/quality/workflows/setup.md` — generated area `description`,
      projection-boundary comment, user-facing recap line, and an active check for
      project-owned runtime harness machinery.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` — the check-8 additions.
- [x] `skills/quality/guides/authoring/agent-harnessability.md` — cross-reference
      touch only if the holistic concept needs a pointer; `None` otherwise.
- [x] `skills/quality/guides/log.md` — record the guide edits.
- [x] `skills/quality/workflows/log.md` — record the setup workflow edit.

### Install / scaffold

None — no scaffolded QUALITY.md content file changes; the generated harness-area
copy lives in `setup.md`, listed under bundled skill above.

## Children

- [Functional spec](0091-agent-harness-holistic-definition/spec.md) — what the
  guidance must say.
- [Design doc](0091-agent-harness-holistic-definition/design.md) — the holistic
  definition, the three-projection model, the scoping decision rule, the
  redefine-not-rename decision, and alternatives.

## Status

`Done`. Implemented and archived. The bundled skill guidance, setup workflow, Top
10 checks, durable skill spec mirrors, docs/logs, and changelog now carry the
holistic harness definition and checked-in governing-artifacts area projection.
Verified with `mise run check`. No CLI, Go, format-schema, rating, roll-up, or
evaluation behavior changed.
