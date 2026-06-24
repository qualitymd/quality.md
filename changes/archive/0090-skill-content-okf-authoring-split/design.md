---
type: Design Doc
title: Skill-content OKF authoring split - design
description: Design for splitting the /quality authoring guide into routed sub-guides, mirroring that structure in durable specs, and introducing OKF-shaped runtime skill indexes, schema, and logs.
tags: [skill, docs, specs, okf, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# Skill-content OKF authoring split - design

## Context

Answers the [functional spec](spec.md) for
[Skill-content OKF authoring split](../0090-skill-content-okf-authoring-split.md).
The change is structural documentation/spec work: runtime skill guides become
smaller routed OKF-shaped concepts, and durable specs mirror the runtime guide
tree. No CLI, Go, format-schema, or evaluation behavior changes.

## Approach

### Runtime skill tree

Keep the existing runtime entry path:

```text
skills/quality/guides/authoring.md
```

That file becomes the entry guide: compact orientation, core concept map, and
reading router. It links to child guides:

```text
skills/quality/guides/authoring/
  index.md
  log.md
  body.md
  model-structure.md
  factors.md
  requirements.md
  rating-scale.md
  agent-harnessability.md
  agent-harness.md
  quality-log.md
```

The child files are split by independently reviewable concern, not by the current
heading layout. The split is intentionally asymmetric where the concern warrants
it: Agent Harnessability and agent harness get separate guides even though they
are related, because reviewers need to critique the model-wide factor and the
steering-materials area without conflating the projections.

Each child guide starts with a small routing block:

```markdown
# Agent Harnessability Authoring

Read this when:

- creating, revising, reviewing, or evaluating the Agent Harnessability factor;
- changing its sub-factor decomposition or boundaries.

Depends on:

- `../authoring.md`
- `model-structure.md`
- `factors.md`
```

The body of each child guide is mostly moved text, lightly edited only to remove
now-local cross-repetition. This preserves doctrine while improving reviewability.

### Runtime OKF affordances

Add a root runtime map:

```text
skills/quality/index.md
skills/quality/schema.md
skills/quality/log.md
```

Add a guide map:

```text
skills/quality/guides/index.md
skills/quality/guides/log.md
```

`SKILL.md` keeps its Agent Skills manifest frontmatter unchanged. It can link to
the runtime OKF index, but the loader-facing frontmatter is not treated as generic
OKF frontmatter unless verified safe. Other guide/workflow files may gain OKF
frontmatter where it is mechanically safe and useful for routing.

The runtime `schema.md` should stay lightweight. It registers the concept types
the skill bundle uses; it is not an enforcement mechanism. This mirrors the repo's
OKF convention without turning runtime skill packaging into a schema-validation
project.

### Durable spec mirror

The durable spec structure mirrors the runtime authoring guide family:

```text
specs/skills/quality-skill/guides/authoring.md
specs/skills/quality-skill/guides/authoring/
  index.md
  body.md
  model-structure.md
  factors.md
  requirements.md
  rating-scale.md
  agent-harnessability.md
  agent-harness.md
  quality-log.md
```

This uses the repo's parent-concept-plus-folder pattern. The old
`authoring-md.md` artifact-spec filename is renamed to `authoring.md` because the
artifact is no longer one flat file with a 1:1 mirror; it is a guide family with a
parent router and child contracts. The design intentionally prioritizes structural
mirror over the older artifact-normalized filename convention for this one guide
family.

Each child spec governs one child guide. The parent `authoring.md` spec owns:

- the entry-guide purpose;
- the reading-router contract;
- shared conformance to `SPECIFICATION.md`;
- shared authoring vocabulary;
- the guarantee that moved doctrine is not weakened.

Child specs own their concern-specific coverage. For example:

- `agent-harnessability.md` owns the model-wide factor, sub-factor boundaries, and
  legacy `harnessability` handling.
- `agent-harness.md` owns the agent-harness area as a steering-materials
  constituent, its factor-family prompts, and the served-domain guardrail.
- `requirements.md` owns requirement statements, assessments, factor references,
  splitting/combining, and overrides.

### Reading contract

The runtime reading rule belongs in three places:

1. `SKILL.md` gives the global rule: read `guides/authoring.md`, then routed
   sub-guides for every touched model element.
2. `guides/authoring.md` gives the routing table.
3. Workflows pin predictable bundles where useful:
   - setup reads the entry guide plus the body, model-structure, factors,
     requirements, rating-scale, Agent Harnessability, agent-harness, and
     quality-log guides;
   - recommendation follow-up reads only the sub-guides matching the confirmed
     mutation surface;
   - evaluation reads authoring sub-guides only when it evaluates the
     `quality-md` self-check or recommends model changes.

This prevents the split from becoming a correctness bug while preserving context
efficiency for narrow work.

## Alternatives

### Keep one authoring guide

Rejected. It preserves the current bottleneck: agents load too much context and
reviewers cannot isolate concerns like Agent Harnessability or agent harness area
modeling.

### Split runtime guides but keep one monolithic durable spec

Rejected. That would make runtime review surfaces smaller while leaving the
contract surface large and easy to drift. The user specifically wants sub-specs
and a mirrored spec structure.

### Put all sub-guides beside `authoring.md`

Rejected. Names such as `authoring-body.md` work, but they do not give the guide
family an OKF-local `index.md` and `log.md`, and they make the relationship less
obvious than a child folder.

### Make every runtime skill file fully OKF in one pass

Deferred. The authoring split is the highest-value first move. Workflows and
resources can be shaped further once the runtime indexes and schema establish the
pattern. `SKILL.md` requires extra care because Agent Skills tooling consumes its
frontmatter.

## Trade-offs & risks

- **More files.** The split increases file count and link maintenance. The offset
  is smaller review units and lower agent context cost.
- **Missed reads.** Agents can miss sub-guides if routing is vague. The mitigation
  is a mandatory reading path in `authoring.md`, a global `SKILL.md` rule, and
  workflow-specific bundles for broad workflows.
- **Spec naming exception.** Renaming `authoring-md.md` to `authoring.md` departs
  from the earlier 1:1 artifact-spec naming convention. The reason is that the
  artifact is becoming a guide family with a parent concept and child specs; the
  mirror structure is more important than preserving the old flat artifact name.
- **Packaging links.** Moving guide content can break relative links inside the
  skill bundle. The npm package check's skill-link validation is the main guard.

## Open questions

- Should workflow specs mirror runtime workflow OKF structure in the same pass, or
  should this case limit the mirror requirement to authoring guides?
- Should `skills/quality/resources/` get an OKF index/log now, or wait until a
  resource-focused split creates a concrete need?
