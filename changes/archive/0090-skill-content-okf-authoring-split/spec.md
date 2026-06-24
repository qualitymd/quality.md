---
type: Functional Specification
title: Skill-content OKF authoring split - functional spec
description: Requirements for splitting the /quality authoring guide into routed OKF sub-guides, adding mirrored durable sub-specs, and structuring runtime skill content as an OKF-shaped bundle without breaking skill loading.
tags: [skill, docs, specs, okf, authoring]
timestamp: 2026-06-24T00:00:00Z
---

# Skill-content OKF authoring split - functional spec

Companion to
[Skill-content OKF authoring split](../0090-skill-content-okf-authoring-split.md).
This spec states what the restructuring must accomplish. The project OKF editing
contract is governed by [Working with OKF](../../../docs/guides/work-with-okf.md).
The runtime skill contract is governed by the durable
[`/quality` skill spec](../../../specs/skills/quality-skill/quality-skill.md).

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The `/quality` authoring guide is intentionally canonical, but it now carries
several independently reviewable concerns in one long file. That hurts two
workflows: agents pay the context cost of the whole guide even when they are
changing one model element, and maintainers cannot easily solicit focused review
on concerns such as Agent Harnessability or agent-harness area modeling. A split
only works if the reading contract is explicit; otherwise agents will keep reading
the entry guide and miss the deeper rules.

The repo already treats `changes/`, `docs/`, and `specs/` as OKF bundles. The
runtime `skills/quality/` content should start following the same pattern:
progressive-disclosure indexes, concept frontmatter where safe, schema/log files,
and a durable spec tree that mirrors the runtime guide tree. The mirror matters:
if runtime guidance is split by authoring concern, durable specs should split by
the same concern so each review surface has one obvious contract.

## Scope

Covered: the runtime `/quality` skill content, especially authoring guidance; the
durable `/quality` skill guide specs; indexes, logs, and schemas needed for OKF
shape; and reading-contract updates in the skill prompt and workflows.

Deferred / non-goals: no CLI behavior, Go code, `SPECIFICATION.md`, rating,
roll-up, evaluation artifact, or scaffold change; no semantic rewrite of the
authoring doctrine; no forced OKF frontmatter on `SKILL.md` if that would conflict
with Agent Skills loader metadata; no migration of frozen archived Change Case
links.

## Requirements

### Keep an authoring entry point and route deeper reads

`skills/quality/guides/authoring.md` **MUST** remain the mandatory entry point for
creating, populating, reviewing, or improving a QUALITY.md file.

`authoring.md` **MUST** become a compact entry guide and router. It **MUST**
retain the core model orientation: QUALITY.md file shape, concept map, model-first
vocabulary, recommended authoring order, and the rule that the format spec governs
on conflict. It **MUST NOT** retain long-form doctrine for every sub-concern once
that doctrine has moved to a routed sub-guide.

`authoring.md` **MUST** include a reading path that maps authoring concerns to
sub-guides. At minimum it **MUST** route:

- Markdown body, unknowns/open questions, review provenance, stakeholder grounding,
  and "sense of good" to `guides/authoring/body.md`.
- Areas, sources, decomposition shapes, traceability graph, normative artifacts,
  constituent kinds, and recurring use-context constituents to
  `guides/authoring/model-structure.md`.
- Factor naming, factor coverage, descriptions, stable-stakes factors, and
  sub-factors to `guides/authoring/factors.md`.
- Requirement statements, assessments, factor connections, splitting/combining
  claims, and rating overrides to `guides/authoring/requirements.md`.
- Rating scale design, rating criteria, required margin, roll-up, veto
  requirements, and `not assessed` handling to `guides/authoring/rating-scale.md`.
- Agent Harnessability as a model-wide factor to
  `guides/authoring/agent-harnessability.md`.
- The agent harness as a modeled area/constituent to
  `guides/authoring/agent-harness.md`.
- Meaningful model-change judgment and quality-log routing to
  `guides/authoring/quality-log.md`.

`SKILL.md` **MUST** state that agents read `guides/authoring.md` first and then
read every routed sub-guide relevant to the model elements they will create,
review, mutate, evaluate, or recommend changing.

Known broad workflows **SHOULD** pin the sub-guide bundles they predictably need.
Setup **SHOULD** read the full authoring bundle needed for first-model creation.
Narrow recommendation follow-up **SHOULD** read only the entry guide plus the
sub-guides matching the confirmed change surface.

> Rationale: the split only saves context if narrow workflows can load narrow
> guidance, but it only stays correct if the entry guide and workflow contracts
> make the required follow-on reads explicit. - 0090

### Split the runtime authoring guide by independently reviewable concerns

The runtime skill **MUST** place authoring sub-guides under
`skills/quality/guides/authoring/` with an `index.md` that lists the sub-guides
and a `log.md` that records updates to that authoring-guide family.

Each authoring sub-guide **MUST** be independently reviewable: it should define
its concern, state when agents read it, state any prerequisite authoring guides,
and carry the doctrine for that concern without requiring a reader to open
unrelated sub-guides.

Each sub-guide **SHOULD** start with a short routing block:

```markdown
# <Guide Title>

Read this when:

- ...

Depends on:

- `../authoring.md`
- ...
```

The split **MUST** preserve existing authoring doctrine unless this Change Case or
a later Change Case explicitly changes that doctrine. Moving text to a sub-guide
**MUST NOT** silently weaken rules around domain agnosticism, model-by-default
constituents, Agent Harnessability, agent-harness area modeling, projection
boundaries, assessable requirements, or quality-log judgment.

Long examples **SHOULD** move out of the entry guide unless they are needed to
explain the entry route. Sub-guides **MAY** retain focused examples for their own
concern.

### Structure runtime skill content as OKF-shaped content

`skills/quality/` **MUST** start carrying OKF bundle affordances:

- `index.md` - progressive-disclosure map of runtime skill content.
- `schema.md` - concept-type registry for runtime skill content.
- `log.md` - update log for runtime skill-content changes.

`skills/quality/guides/` **MUST** carry `index.md` and `log.md`. Runtime workflow
folders **SHOULD** gain indexes/logs when useful for routing, but this case **MAY**
defer deeper workflow OKF conversion if it would distract from the authoring split.

Runtime guides and workflows **SHOULD** carry OKF concept frontmatter with at least
`type`, `title`, and `description` where doing so is compatible with their
consumers. `SKILL.md` **MUST NOT** receive OKF frontmatter changes that risk
breaking Agent Skills loading; if it already has loader-required frontmatter, that
frontmatter remains authoritative for loader metadata.

The skill runtime schema **SHOULD** register concept types such as `Runtime Guide`,
`Runtime Workflow`, `Runtime Resource`, `Runtime Index`, and `Runtime Schema` if
those types are used.

### Mirror runtime guide structure in durable specs

The durable `/quality` skill guide specs **MUST** mirror the runtime guide tree for
the authoring guide family. If the runtime tree is:

```text
skills/quality/guides/authoring.md
skills/quality/guides/authoring/
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

then the durable spec tree **MUST** use the corresponding parent-concept-plus-folder
shape:

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

The parent durable spec `guides/authoring.md` **MUST** govern the entry guide and
shared authoring-guide-family invariants. Each child spec **MUST** govern the
matching runtime sub-guide.

The durable spec tree **MUST** be structured for OKF review: every child spec has
frontmatter with `type: Functional Specification`, the enclosing `index.md` lists
the children, and the relevant `log.md` records the split.

The existing monolithic durable guide spec
`specs/skills/quality-skill/guides/authoring-md.md` **MUST** be renamed or replaced
so it does not remain a second source of truth for the authoring guide family.
If renamed to `guides/authoring.md`, the implementation **MUST** update live links
and leave archived Change Cases untouched.

> Rationale: retaining a monolithic spec mirror after splitting the runtime guide
> would preserve the review bottleneck and make the sub-guides drift. Mirroring the
> tree makes the review unit and the contract unit the same concern. - 0090

### Preserve link and packaging integrity

All live links affected by the runtime and durable spec moves **MUST** be updated.
Append-only logs and archived Change Cases **MUST NOT** be rewritten solely to
chase renamed guide/spec paths.

The npm/package skill-bundle check **MUST** still pass, including relative-link
resolution for `skills/quality`.

The markdown formatting check **MUST** pass after the split.

## Durable spec changes

### To add

- `specs/skills/quality-skill/guides/authoring/` child specs - add sub-specs for
  body authoring, model structure, factors, requirements, rating scale,
  Agent Harnessability, agent-harness modeling, and quality-log judgment (per
  [Mirror runtime guide structure in durable specs](#mirror-runtime-guide-structure-in-durable-specs)).
- `specs/skills/quality-skill/guides/authoring/index.md` - list the authoring
  sub-specs (per [Mirror runtime guide structure in durable specs](#mirror-runtime-guide-structure-in-durable-specs)).

### To modify

- `specs/schema.md` - add any new concept type needed for guide-family specs, or
  leave unchanged if `Functional Specification` remains sufficient (per
  [Mirror runtime guide structure in durable specs](#mirror-runtime-guide-structure-in-durable-specs)).
- `specs/skills/quality-skill/index.md` - update links/descriptions for the
  authoring guide spec family (per
  [Mirror runtime guide structure in durable specs](#mirror-runtime-guide-structure-in-durable-specs)).
- `specs/skills/quality-skill/guides/index.md` - update the authoring entry to
  point at the new parent spec/folder (per
  [Mirror runtime guide structure in durable specs](#mirror-runtime-guide-structure-in-durable-specs)).
- `specs/skills/quality-skill/guides/log.md` - record the durable guide-spec split
  (per [Mirror runtime guide structure in durable specs](#mirror-runtime-guide-structure-in-durable-specs)).
- `specs/skills/quality-skill/quality-skill.md` - update the authoring reading
  contract if it is specified there (per
  [Keep an authoring entry point and route deeper reads](#keep-an-authoring-entry-point-and-route-deeper-reads)).
- `specs/skills/quality-skill/workflows/setup.md` and
  `specs/skills/quality-skill/workflows/evaluate.md` - update only if they specify
  authoring-guide reads that must pin routed sub-guides (per
  [Keep an authoring entry point and route deeper reads](#keep-an-authoring-entry-point-and-route-deeper-reads)).

### To rename

- `specs/skills/quality-skill/guides/authoring-md.md` →
  `specs/skills/quality-skill/guides/authoring.md` - convert the current
  monolithic artifact spec into the parent contract for the authoring guide family
  (per [Mirror runtime guide structure in durable specs](#mirror-runtime-guide-structure-in-durable-specs)).

### To delete

None
