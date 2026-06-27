---
type: Runtime Guide
title: Authoring QUALITY.md
description: Entry point and router for authoring, reviewing, and improving QUALITY.md files.
tags: [quality, authoring, guide]
---

# Authoring QUALITY.md

This is the entry point for working with QUALITY.md files. It gives the core
file shape, vocabulary, authoring order, and routing map. Read this first, then
read every routed sub-guide relevant to the model elements you will create,
review, mutate, evaluate, or recommend changing.

This guide conforms to [`SPECIFICATION.md`](../resources/SPECIFICATION.md). The
specification governs on any conflict.

## Reading Path

- Markdown body, unknowns, open questions, review provenance, stakeholder
  grounding, and the governing sense of "good" →
  [`authoring/body.md`](authoring/body.md).
- Areas, sources, decomposition shapes, traceability graph, normative artifacts,
  constituent kinds, and recurring use-context constituents →
  [`authoring/model-structure.md`](authoring/model-structure.md).
- Factor naming, factor coverage, descriptions, stable-stakes factors, and
  sub-factors → [`authoring/factors.md`](authoring/factors.md).
- Requirement names and titles, assessments, factor connections,
  splitting/combining claims, and rating overrides →
  [`authoring/requirements.md`](authoring/requirements.md).
- Rating Scale design, criteria, required margin, roll-up, veto requirements, and
  `not assessed` handling → [`authoring/rating-scale.md`](authoring/rating-scale.md).
- Agent Harnessability as a model-wide factor →
  [`authoring/agent-harnessability.md`](authoring/agent-harnessability.md).
- The agent harness as a modeled area/constituent →
  [`authoring/agent-harness.md`](authoring/agent-harness.md).
- Meaningful model-change judgment and quality-log routing →
  [`authoring/quality-log.md`](authoring/quality-log.md).

For first-model setup, read the body, model-structure, factors, requirements,
rating-scale, Agent Harnessability, agent-harness, and quality-log guides before
authoring. For narrow recommendation follow-up, read only the sub-guides matching
the confirmed mutation surface.

For direct `QUALITY.md` edits, infer the intended change from the user's request
and the current model before asking follow-up. Preserve the user's intent, the
model or body target, the planned change, the value prop, the expected judgment
effect, unresolved unknowns, and whether the change should write the quality log.
Ask follow-up only when one of those would materially change the edit; otherwise
state the intended edit in simple planned-change prose, name the value prop, and
invite concerns, goals, needs, worries, edge cases, naming preferences, scope
boundaries, constraints, or `looks good` as a review gate, then wait for that
response before mutating.

## The QUALITY.md file

A QUALITY.md file is a Markdown file with two parts:

- **YAML frontmatter** — the **quality model**: a structured, declarative
  description of what quality means for the entity being evaluated.
- **Markdown body** — the evaluable judgment context for the model: what the
  root area is, why quality matters, what decisions the model supports, and what
  context is missing or inaccessible.

The whole file represents a single root **area** — the top entity whose quality
is modeled. The file's location carries meaning: a `QUALITY.md` in a directory
makes that directory and everything under it (`**/*`) the default scope of
evaluation, unless an area narrows it with a `source`.

## Core Concepts

| Concept      | Meaning                                                        |
| ------------ | -------------------------------------------------------------- |
| Model        | The root quality model in a QUALITY.md file.                   |
| Area         | The thing being evaluated.                                     |
| Source       | The material assessed for an area, such as a path or selector. |
| Factor       | A quality dimension that matters for an area.                  |
| Requirement  | A specific quality expectation.                                |
| Assessment   | The means of checking a requirement against an area source.    |
| Finding      | An observation produced by an assessment.                      |
| Rating Scale | The ordered model-wide scale used to rate results.             |

## Quality Model

The **quality model** is the frontmatter. The root model is also the apex Area and
carries one model-wide property, `ratingScale`, plus all Area properties.

| Property       | Presence   | What it is                                                      |
| -------------- | ---------- | --------------------------------------------------------------- |
| `title`        | Required   | Human-readable name of the entity whose quality is modeled.     |
| `description`  | Optional   | A concise statement of what the model's area is.                |
| `ratingScale`  | Required   | The model-wide levels every result is rated on.                 |
| `factors`      | Optional\* | Factors — lenses through which quality is described.            |
| `requirements` | Optional\* | Requirements assessed against the root source.                  |
| `areas`        | Optional\* | Child areas, nested to any depth.                               |
| `source`       | Optional   | Scope override; omit at the root to take the directory default. |

\* At least one of `factors`, `requirements`, or `areas` must be present.

Area names, Factor names, Requirement names, and Rating Level IDs use the same
strict name grammar: letters or digits at both ends, with letters, digits, `_`,
or `-` inside. Requirement titles stay natural language. When a tool needs a
stable text handle, it uses canonical model references such as `area:root`,
`area:api`, `factor:api::reliability`,
`requirement:api::retry-window`, or `rating:target`.

## Authoring Order

1. Name the entity with `title` and think through the body's Overview first.
2. Fill the Markdown body before expanding the frontmatter.
3. Confirm the Rating Scale after the body and before writing requirements.
4. Derive factors and requirements from the body context.
5. Trace at least one important concern from body to model before expanding the
   tree: a need names the outcome, a risk names the failure mode, a factor names
   the quality lens, and a requirement names the inspectable expectation.
6. Keep the root lean when child areas carry the detail; declare model-wide
   factors at the root and push narrower factors/requirements down to child areas.
7. Make sure the model reaches requirements somewhere in the area tree before
   treating it as evaluable.

## Maintenance

When authoring changes meaningfully alter what the model *is* or *how it judges*,
read [`authoring/quality-log.md`](authoring/quality-log.md) before changing the
model or writing a quality-log entry.
