---
type: Change Case
title: Harnessability factor
description: Add harnessability — how well a project equips an agent to work on it — as a model-wide factor (with sub-factors) the skill proposes by default for agent-collaborated composite roots, so the agent-collaboration concern leads with its factor projection instead of an easily-deferred area.
status: Draft
tags: [skill, authoring, setup, factors, harnessability, agentic]
timestamp: 2026-06-24T00:00:00Z
---

# Harnessability factor

A **Change Case** to give the agent-collaboration concern its *factor* projection
in the `/quality` skill. Today the skill carries that concern only as a
*constituent* — the **agent harness** (see [Carry the recurring use-context
constituents](../skills/quality/guides/authoring.md)) — and 0080 made that
constituent model-by-default. But the concern's most useful projection across a
composite root is a **factor**: *how well does each part of this project equip an
agent to work on it well?* That question recurs across constituents, which is the
signature of a model-wide factor, and unlike an area it is never dropped for
thinness — a thin harness just rates low and surfaces a finding everywhere.

This case adds **harnessability** as a deliberate **umbrella** model-wide factor,
decomposed into six sub-factors that trace the agent's working loop, and teaches
the skill to propose it by default for an agent-collaborated composite root. The
agent harness stays a constituent (the *artifact* projection) and the agent stays
an *audience*; harnessability is the *factor* projection of the same concern.

Detail lives in:

- [Functional spec](0081-harnessability-factor/spec.md) — what the guidance must
  say.
- [Design doc](0081-harnessability-factor/design.md) — why factor-over-area, how
  the six sub-factors were derived and verified, and the boundary discipline that
  keeps them from double-counting.

## Motivation

A real setup run on a multi-service monorepo (the same run that motivated 0080)
reported the agent harness as "thin — only the quality skill symlinked — not
earning its own area yet," and deferred it. 0080 fixes the *area-inclusion* side
of that failure: a germane constituent is now modeled by default and a thin one
surfaces as a ratable gap rather than a deferral. But it leaves the concern on the
**constituent** axis, where it is measured by the substance of one steering
artifact.

The harness-engineering literature frames the concern differently. The quality at
stake is **agent legibility / steerability / verifiability of the whole project**
— *"anything an agent can't see doesn't exist"* (OpenAI), *guides and sensors*
around the model (Böckeler) — a property that applies to every part of the
codebase, not just the steering files. That is a quality lens, i.e. a **factor**,
and a *model-wide* one: it recurs across constituents (the server, the schema, the
tests, the docs are each more or less harnessable). Leading with the factor
projection puts the concern where the guide's own rules say it belongs (a concern
projects as factor / constituent / audience; name the projection meant), and makes
it impossible to drop for thinness — the same robustness `testability` already has
("no tests is not a reason to omit `testability`").

## Scope

Covered: adding harnessability as a model-wide umbrella factor with six
sub-factors and operational definitions; teaching the authoring guide and setup
workflow to propose it by default for an agent-collaborated composite root; a Top
10 check for its coverage; the boundary discipline that keeps the sub-factors
non-overlapping with each other, with existing common factors, and with the agent
harness constituent; and a `## Why QUALITY.md` README use case with citations.

Deferred / non-goals: no change to the QUALITY.md format or schema, to
`SPECIFICATION.md`, or to the CLI — recommending a factor is authoring judgment,
not format semantics (the layering 0080 preserved). The continuous-improvement of
this equipping over time stays the model's existing **model-wide learn-loop
concern**, not a seventh sub-factor. Fleet orchestration and the cross-project
human-attention economy sit above a project Area and are out of scope. Re-checking
this repo's own `QUALITY.md` against the new factor is a follow-up.

## Affected artifacts

Found by sweeping for the agent-harness / use-context-constituent guidance and the
model-wide-factor guidance across the skill and its spec mirror. Grouped by kind;
an empty kind is a deliberate "no impact," not an oversight.

### Code

None — skill-guidance, spec-mirror, and README only; no `cmd/` or `internal/`
change.

### Format spec (`SPECIFICATION.md`)

None — factor selection is authoring judgment, not format semantics (per 0080's
layering rationale).

### Durable specs (`specs/`)

The functional spec's [Durable spec changes](0081-harnessability-factor/spec.md)
section is the authoritative breakdown. In summary:

- `specs/skills/quality-skill/guides/authoring-md.md` — require the harnessability
  model-wide factor, its six sub-factors, the factor projection of the
  use-context harness concern, and the double-count boundary.
- `specs/skills/quality-skill/workflows/setup.md` — require setup to propose
  harnessability by default for an agent-collaborated composite root.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — add the
  harnessability-coverage check.
- `specs/skills/quality-skill/guides/log.md` and
  `specs/skills/quality-skill/workflows/log.md` — record the revision.

### Durable docs

- `README.md` — add the `### Evaluate and Improve Agent Harnessability` use case
  under `## Why QUALITY.md`, with citations.

### Bundled skill (`skills/quality/`)

- `skills/quality/guides/authoring.md` — runtime counterpart of the authoring-md
  spec changes.
- `skills/quality/workflows/setup.md` — runtime counterpart of the setup spec
  changes.
- `skills/quality/guides/top-10-quality-md-checks.md` — runtime counterpart of the
  Top 10 check.

### Install / scaffold

None.

### Changelog

- `CHANGELOG.md` — note the new harnessability factor guidance.

## Children

- [Functional spec](0081-harnessability-factor/spec.md) — what the guidance must
  say.
- [Design doc](0081-harnessability-factor/design.md) — factor-over-area, the
  six-sub-factor derivation and verification, and the double-count boundary
  discipline.

## Status

`Draft`. Functional spec authored. Design doc and implementation to follow; no
code involved, so the case advances Draft → Design → In-Progress for the durable
guidance, README, and spec edits, then In-Review. Durable specs and docs are
edited in In-Progress alongside the skill.
