---
type: Functional Specification
title: Encode projection boundaries in the model — functional spec
description: Requirements for teaching the authoring guidance to encode the boundary between two same-rooted projections of one concern into the emitted model.
tags: [skill, authoring, factors, projections, harnessability]
timestamp: 2026-06-24T00:00:00Z
---

# Encode projection boundaries in the model — functional spec

Companion to the [Encode projection boundaries in the
model](../0087-encode-projection-boundaries.md) Change Case. This spec states
_what_ the guidance must require; the [design doc](design.md) covers _how_ and
_why_.

The key words **MUST**, **SHOULD**, and **MAY** are to be interpreted as described
in IETF RFC 2119 (BCP 14) when, and only when, they appear in all capitals.

## Background / Motivation

The authoring guide's three-projections rule teaches the author to _name the
projection meant and model it once_, preventing double-counting during modeling.
It is addressed to the author's reasoning, not to the model's reader: a concern
that legitimately appears as both a factor and an area (the Agent Harnessability
factor and the agent-harness area) leaves no trace in the emitted model of _why_
both exist or _how_ they differ. A field setup run produced exactly that model —
correct, and even carrying clarifying YAML comments — yet a reader still asked
whether the factor was meant to replace the area. The clarifying comments were the
agent's own judgment, not a guidance requirement, so the legibility is not
reliable. The reason a requirement exists must survive into the model the same way
this project's own guides require rationale to survive into durable specs.

## Scope

Covered: a general authoring requirement to encode a concern's projection boundary
in the emitted model when more than one projection is modeled; the encoding
mechanisms (YAML comment, and a description clause where both projections are rated
nodes that surface in a report); the application of that rule to the Agent
Harnessability / agent-harness boundary; and a Top 10 readiness check for the
missing boundary note.

Deferred / non-goals:

- **No format or schema change.** YAML comments and `description` prose are already
  valid model content; this governs `/quality` authoring judgment, not
  `SPECIFICATION.md`.
- **No deterministic CLI lint rule.** Boundary legibility is a judgment carried by
  the skill and the Top 10 check, not a machine-checkable contract.
- **No retrofit mandate.** Existing models that already carry the boundary are
  conformant; others may be updated opportunistically during authoring work.

## Requirements

- The authoring guide **MUST** teach that when a model carries two or more
  projections of one concern (factor, constituent/area, audience), the author
  encodes the projection boundary in the emitted model rather than only reasoning
  about it during modeling.

  > Rationale: the three-projections rule prevents double-counting for the author,
  > but leaves the boundary invisible to a reader of the generated model, who then
  > re-litigates whether the projections are redundant. — 0087

- The guide **MUST** require, as the primary mechanism, a YAML comment on each
  modeled projection's node that names the sibling projection and states the
  one-line distinction between them.

- The guide **MUST** require, **when both projections are rated nodes that surface
  in an evaluation report**, a short disambiguating clause in each node's
  `description` in addition to the comment.

  > Rationale: YAML comments do not survive rendering into an evaluation report, so
  > a reader of the report would see two same-rooted nodes with no way to tell them
  > apart; the description is the only carrier that survives. — 0087

- The guide **SHOULD** keep the disambiguating description clause to the
  distinction from the sibling projection, consistent with the existing
  "distinguishes, not enumerates" rule for descriptions, so the clause does not
  restate the node's factors or requirements.

- The guide **MUST** apply this rule to the Agent Harnessability / agent-harness
  boundary as the canonical instance: the model **SHOULD** carry the boundary on
  both the `agent-harnessability` factor and the agent-harness area (a YAML comment
  on each, plus a disambiguating clause in each description).

- The Top 10 QUALITY.md checks **MUST** include a finding that routes to authoring
  when a model carries two same-rooted projections of one concern (e.g. an Agent
  Harnessability factor and an agent-harness area) with no boundary note —
  neither a YAML comment nor a disambiguating description clause distinguishing
  them.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` — add the requirement that
  the guide teach encoding the projection boundary in the emitted model (the
  general rule and its two mechanisms), and extend the Agent Harnessability
  factor/area boundary requirement to require that boundary be encoded in the model
  (per the projection-boundary requirements above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — add the
  missing-boundary-note check (per the Top 10 requirement above).

### To rename

None.

### To delete

None.
