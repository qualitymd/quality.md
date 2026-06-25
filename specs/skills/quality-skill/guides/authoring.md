---
type: Functional Specification
title: QUALITY.md authoring guide family
description: Parent contract for the /quality skill authoring entry guide and its routed authoring sub-guides.
tags: [skill, quality, guide, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# QUALITY.md authoring guide family

This spec governs the authoring guide family the [`/quality` skill](../quality-skill.md)
ships at [`skills/quality/guides/authoring.md`](../../../../skills/quality/guides/authoring.md)
and [`skills/quality/guides/authoring/`](../../../../skills/quality/guides/authoring/).
The entry guide is the mandatory runtime entry point and router; child specs in
[`authoring/`](authoring/index.md) govern the routed sub-guides. The format the
guides teach is defined by [`SPECIFICATION.md`](../../../../SPECIFICATION.md), the
source of truth the guides conform to.

This document uses BCP 14 keywords only for testable conformance requirements. The
key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Purpose

The guide family exists to be the canonical reference and best-practices guidance
for understanding and working with QUALITY.md files while keeping independent
authoring concerns reviewable and context-efficient. The entry guide gives the
core file shape, vocabulary, authoring order, and routing map. Child guides carry
concern-specific doctrine.

## Entry guide contract

`skills/quality/guides/authoring.md` **MUST** remain the guide agents read first
when creating, populating, reviewing, or improving a QUALITY.md file.

The entry guide **MUST** state that it conforms to `SPECIFICATION.md` and that the
specification governs on conflict. It **MUST** define the two-layer QUALITY.md file
shape, the core Model/Area/Source/Factor/Requirement/Assessment/Finding/Rating
Scale vocabulary, the root model's key properties, strict name grammar for Area
names, Factor names, and Rating Level IDs, and the recommended authoring order.

The entry guide **MUST** route authoring concerns to these child runtime guides:

- `authoring/body.md` - Markdown body, unknowns, open questions, review provenance,
  stakeholder grounding, and sense of good.
- `authoring/model-structure.md` - Areas, Sources, decomposition, traceability,
  normative artifacts, constituent kinds, and recurring use-context constituents.
- `authoring/factors.md` - factor naming, coverage, descriptions, stable-stakes
  factors, and sub-factors.
- `authoring/requirements.md` - Requirement names and titles, Assessments, factor
  connections, splitting/combining claims, and rating overrides.
- `authoring/rating-scale.md` - Rating Scale design, criteria, required margin,
  roll-up, veto requirements, and `not assessed` handling.
- `authoring/agent-harnessability.md` - Agent Harnessability as a model-wide
  factor.
- `authoring/agent-harness.md` - the agent harness as a modeled Area/constituent.
- `authoring/quality-log.md` - meaningful model-change judgment and quality-log
  routing.

The entry guide **MUST NOT** retain the full long-form doctrine for each concern
once that doctrine lives in a routed sub-guide.

## Child guide contract

Each child guide **MUST** state when agents read it and which authoring guides it
depends on. Each child guide **MUST** be independently reviewable for its concern
and **MUST NOT** require reading unrelated child guides to understand its local
rules.

The guide family **MUST** preserve existing authoring doctrine unless a Change Case
explicitly changes it. Moving doctrine from the former monolithic guide to a child
guide **MUST NOT** weaken domain agnosticism, model-by-default constituents,
Agent Harnessability, agent-harness area modeling, projection boundaries,
assessable requirements, or quality-log judgment.

## Reading contract

The root skill prompt **MUST** require agents to read `guides/authoring.md` first
and then every routed sub-guide relevant to the model elements they will create,
review, mutate, evaluate, or recommend changing.

Broad workflows such as setup **SHOULD** pin the authoring sub-guide bundle they
predictably need. Narrow workflows such as recommendation follow-up **SHOULD** read
only the child guides matching the confirmed mutation surface.
