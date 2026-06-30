---
type: Change Case
title: Skill-content OKF authoring split
description: Restructure the /quality skill's authoring guidance into OKF-shaped runtime sub-guides and mirrored durable sub-specs, keeping authoring.md as the routed entry point and making agent-read obligations explicit.
status: Done
tags: [skill, docs, specs, okf, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# Skill-content OKF authoring split

A **Change Case** to split the long `/quality` authoring guide into independently
reviewable, context-efficient sub-guides, while starting to structure the runtime
skill content itself as an OKF bundle. The durable spec structure mirrors the
runtime guide structure so each guide concern has a matching sub-spec.

Detail lives in:

- [Functional spec](0090-skill-content-okf-authoring-split/spec.md) - what the
  runtime skill bundle and durable specs must do.
- [Design doc](0090-skill-content-okf-authoring-split/design.md) - the OKF-shaped
  tree, routing contract, migration approach, and alternatives.

## Motivation

`skills/quality/guides/authoring.md` has become the canonical source for too many
independent concerns: body authoring, model structure, factors, requirements,
rating scale, Agent Harnessability, agent-harness area modeling, and quality-log
judgment. That makes it expensive for agents to load even for narrow work, and
hard for maintainers to review or solicit feedback on one concern without reading
the whole guide. Agent Harnessability and agent-harness quality modeling are the
clearest current examples: they need focused review surfaces, but today they are
embedded inside one large file.

At the same time, the skill's runtime content is not yet organized as an OKF
bundle even though the repo's durable `specs/`, `docs/`, and `changes/` bundles
are. Starting the skill content as OKF gives agents a progressive-disclosure map,
gives maintainers smaller review units, and lets runtime guides and durable guide
specs stay aligned by structure rather than by memory.

## Scope

Covered:

- Keep `skills/quality/guides/authoring.md` as the mandatory authoring entry point
  and router.
- Split the long authoring guide into routed sub-guides for body authoring,
  model structure, factors, requirements, rating scale, Agent Harnessability,
  agent-harness area modeling, and quality-log judgment.
- Create matching durable sub-specs under the `/quality` skill specs, with the
  spec tree mirroring the runtime guide tree.
- Add an explicit reading contract to `SKILL.md` and relevant workflows: agents
  read `guides/authoring.md` first, then every routed sub-guide relevant to the
  model elements they will create, review, mutate, evaluate, or recommend
  changing.
- Start structuring `skills/quality/` as an OKF-shaped runtime bundle with
  indexes, schema, and logs where they help progressive disclosure.
- Preserve installable skill compatibility, especially the skill manifest
  frontmatter in `SKILL.md`.

Deferred / non-goals:

- No QUALITY.md format-schema, rating, roll-up, or evaluation-semantic change.
- No CLI or Go behavior change.
- No content redesign of Agent Harnessability or agent-harness modeling beyond
  moving the current guidance into independent review surfaces.
- No requirement that every runtime file become a full OKF concept if doing so
  would conflict with Agent Skills loader metadata.
- No migration of historical Change Case links.

## Affected artifacts

Derived by analysis: this changes the runtime skill guide bundle, the durable
skill-guide specs that mirror it, the skill prompt that routes required reading,
and logs/indexes for the affected OKF bundles. Empty kinds are deliberate.

### Code

None - documentation, spec, and bundled skill guidance content only.

### Format spec (`SPECIFICATION.md`)

None - no format rule changes.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0090-skill-content-okf-authoring-split/spec.md#durable-spec-changes)
section is authoritative. In summary:

- [x] Add or update `/quality` skill durable specs so authoring-guide sub-specs
      mirror runtime authoring sub-guides.
- [x] Rename or reshape the existing authoring guide contract as needed so it
      becomes the parent contract for the authoring guide family, not another
      monolithic mirror.
- [x] Update `specs/skills/quality-skill/index.md`,
      `specs/skills/quality-skill/guides/index.md`, relevant guide logs, and
      workflow/follow-up specs that carry authoring read contracts.

### Durable docs

- [x] `docs/guides/work-with-okf.md` - no update required; existing OKF
      conventions cover the runtime skill bundle.
- [x] `CHANGELOG.md` - release note for the skill-guidance/spec restructuring.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/SKILL.md` - require the routed authoring reading contract.
- [x] `skills/quality/index.md`, `skills/quality/schema.md`, `skills/quality/log.md` - root runtime OKF map, concept-type registry, and update log.
- [x] `skills/quality/guides/authoring.md` - entry point and router.
- [x] `skills/quality/guides/authoring/` - sub-guides and local index/log.
- [x] `skills/quality/guides/index.md` and `skills/quality/guides/log.md` -
      runtime guide map and update log.
- [x] `skills/quality/workflows/setup.md`,
      `skills/quality/guides/recommendation-follow-up.md`, and supporting runtime
      workflow/resource frontmatter and indexes - routed read contracts and OKF
      shape.

### Install / scaffold

None - no scaffolded QUALITY.md content changes.

## Children

- [Functional spec](0090-skill-content-okf-authoring-split/spec.md) - what the
  runtime OKF structure, sub-guide split, mirrored sub-specs, and reading contract
  must do.
- [Design doc](0090-skill-content-okf-authoring-split/design.md) - the tree shape,
  routing rules, compatibility guardrails, and migration sequence.

## Status

`Done`. Implemented and archived. The runtime `/quality` skill content now has
OKF-shaped indexes, schema, and logs; `guides/authoring.md` is the authoring entry
point and router; concern-specific authoring sub-guides live under
`skills/quality/guides/authoring/`; durable guide specs mirror that tree under
`specs/skills/quality-skill/guides/authoring/`; and `SKILL.md`, setup,
recommendation follow-up, logs, and CHANGELOG carry the routed-read contract. No
code, CLI, format-schema, rating, roll-up, or evaluation behavior changed.
Verified with `mise run check`.
