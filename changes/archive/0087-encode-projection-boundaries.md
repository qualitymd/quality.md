---
type: Change Case
title: Encode projection boundaries in the model
description: Teach the authoring guidance to encode the boundary between two same-rooted projections of one concern (e.g. the Agent Harnessability factor and the agent-harness area) into the emitted model — a YAML comment on each node and, where both surface in a report, a disambiguating description clause — so a reader of the generated QUALITY.md can tell them apart.
status: Done
tags: [skill, authoring, factors, projections, harnessability]
timestamp: 2026-06-24T00:00:00Z
---

# Encode projection boundaries in the model

A **Change Case** to make the [three-projections rule](../../skills/quality/guides/authoring.md) produce a _legible_ model. The guide
already teaches the author to reason about the boundary between a concern's
projections — _name the projection meant, model it once_ — and the Agent
Harnessability factor vs. the agent-harness area is its sharpest instance. But
nothing tells the author to **encode that boundary into the model they emit**. The
reasoning archives with the setup run; the reader of the generated `QUALITY.md` is
left to reconstruct why one concern appears as both a factor and an area.

This case adds a guided requirement: when a model carries two or more projections
of one concern, encode the boundary at the point of definition — a YAML comment on
each projection's node naming its sibling and the one-line distinction, and, when
both projections are rated nodes that surface in an evaluation report, a short
disambiguating clause in each `description`. The rule is general (any concern with
multiple projections); Agent Harnessability / agent-harness is the canonical
instance.

Detail lives in:

- [Functional spec](0087-encode-projection-boundaries/spec.md) — what the guidance
  must say.
- [Design doc](0087-encode-projection-boundaries/design.md) — why comment-primary
  with an optional description clause, why it generalizes, and the boundary against
  the "distinguishes, not enumerates" description rule.

## Motivation

A real setup run produced a structurally correct model for a multi-service
monorepo: it carried both the model-wide **Agent Harnessability** factor and a
distinct **agent-harness** area, with hand-written YAML comments on each node
explaining how they differ. A reader of that model still asked why both exist and
whether the factor was meant to replace the area — exactly the confusion the
boundary discipline is supposed to prevent. The clarifying comments were the
agent's own judgment, not something the guidance requires, so the next model may
omit them and re-create the confusion.

The guide's [three-projections rule](../../skills/quality/guides/authoring.md) (and the 0081/0085 Agent
Harnessability guidance built on it) is addressed to the _author's reasoning_: it
prevents double-counting during modeling. It does not ask the author to leave a
trace of that reasoning _in the model_, so the boundary is invisible to whoever
reads the result. The same gap applies to any concern modeled in more than one
projection — the guide's own `secure` example (a security factor, a threat-model
constituent, an auditor audience). The fix belongs at the general rule, with Agent
Harnessability as one instance.

## Scope

Covered: a general authoring `Do` requiring the projection boundary to be encoded
in the emitted model (YAML comment on each node; disambiguating description clause
when both projections are rated nodes that appear in a report); a reinforcing
pointer at the Agent Harnessability / agent-harness guidance so the canonical
instance applies it; a Top 10 readiness check that flags two same-rooted
projections carrying no boundary note; and the matching durable spec-mirror
requirements with promoted rationale.

Deferred / non-goals: no change to the QUALITY.md format or schema, to
`SPECIFICATION.md`, or to the CLI — YAML comments and description prose are already
valid model content; this is authoring judgment, not format semantics. No retrofit
of existing models (a model already carrying the boundary, like the field model
that motivated this, is conformant; others may be updated opportunistically). No
new lint rule in the deterministic CLI — boundary legibility is a judgment the
skill and the Top 10 check carry, not a machine-checkable contract.

## Affected artifacts

Derived by sweeping for the three-projections rule, the Agent
Harnessability/agent-harness boundary guidance, and the projection vocabulary
across the bundled skill, its spec mirror, the Top 10 checks, the changelog, and
the format spec. Grouped by kind; an empty kind is a deliberate "no impact."

### Code

None — authoring judgment and model content, not CLI or Go behavior.

### Format spec (`SPECIFICATION.md`)

None — YAML comments and `description` prose are already valid; the boundary is an
authoring practice, not a schema rule.

### Durable specs (`specs/`)

The functional spec's [Durable spec changes](0087-encode-projection-boundaries/spec.md) section is authoritative. In
summary:

- `specs/skills/quality-skill/guides/authoring-md.md` — require the guide to teach
  encoding the projection boundary in the emitted model (general rule), and require
  the Agent Harnessability / agent-harness boundary to be encoded in the model.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — add the
  missing-boundary-note check.
- `specs/skills/quality-skill/guides/log.md` — record the revision.

### Durable docs

None — README's Agent Harnessability use case stands; this is authoring practice,
not a user-facing concept change.

### Bundled skill (`skills/quality/`)

- `skills/quality/guides/authoring.md` — runtime counterpart: the general
  projection-boundary `Do` and the reinforced Agent Harnessability bullet.
- `skills/quality/guides/top-10-quality-md-checks.md` — runtime counterpart of the
  new check.

### Install / scaffold

None.

### Changelog

- `CHANGELOG.md` — note the projection-boundary authoring guidance.

## Children

- [Functional spec](0087-encode-projection-boundaries/spec.md) — what the guidance
  must say.
- [Design doc](0087-encode-projection-boundaries/design.md) — comment-primary
  rationale, generalization, and the description-rule boundary.

## Status

`Done`. Implemented and archived after updating the bundled authoring guide and
Top 10 checks, their durable spec mirrors, the guides log, and CHANGELOG. The
general projection-boundary rule lands at the three-projections rule with the Agent
Harnessability / agent-harness boundary as its canonical instance, and the Top 10
checks gain a missing-boundary-note finding. Reconciled the Affected artifacts list;
no `SPECIFICATION.md`, CLI schema, or Go code change was needed. Verified with
`mise run check`.
