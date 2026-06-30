---
type: Functional Specification
title: Agent Harnessability naming — functional spec
description: What the /quality skill must teach when renaming harnessability to Agent Harnessability with an accountability-preserving definition.
tags: [skill, authoring, setup, factors, harnessability, agentic]
timestamp: 2026-06-24T00:00:00Z
---

# Agent Harnessability naming — functional spec

Companion to the
[Agent Harnessability naming](../0085-agent-harnessability-naming.md) change
case. This spec states _what_ the `/quality` skill guidance must say when it
renames the 0081 harnessability factor to **Agent Harnessability**. The
[design doc](design.md) covers why the display name, recommended key, and
accountability wording were chosen.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

0081 added harnessability as the model-wide factor projection of the
agent-collaboration concern for agent-collaborated composite roots. The modeling
shape is still right: the factor names how the project equips an agent, the agent
harness remains the constituent, and the agent remains the audience. The name and
definition need tightening.

Plain "harnessability" carries useful harness-engineering lineage, but it is not
immediately intelligible inside a quality model and can be mistaken for the quality
of the agent harness artifact. **Agent Harnessability** makes the audience and use
context visible in the factor title while preserving the term of art.

The existing definition also says the project equips an agent to do good work with
"limited human attention" or "largely unsupervised." That correctly points to
avoidable re-specification and review toil, but it can suggest that human
responsibility disappears. The current definition should instead say that the
project equips agent work while preserving clear human direction, review, and
accountability.

## Scope

Covered: human-facing factor naming, recommended model key, definition wording,
setup guidance, authoring guidance, Top 10 coverage checks, README wording,
CHANGELOG wording, and durable spec mirrors for the `/quality` skill.

Deferred / non-goals: no QUALITY.md schema change, no CLI behavior change, no
automatic migration, no rewrite of archived Change Cases or append-only historical
logs, and no change to the 0081 six-sub-factor decomposition.

## Requirements

### Human-facing name and model key

The authoring guide and setup workflow **MUST** use **Agent Harnessability** as
the human-facing title for the factor and **MUST** recommend
`agent-harnessability` as the stable YAML factor key for new or revised models.

```yaml
agent-harnessability:
  title: Agent Harnessability
```

The guidance **MUST NOT** use `Agent-Harnessability` as the human-facing title.
Hyphenation belongs in the stable key, not in the display label.

> Rationale: the title should read naturally in reports and setup summaries, while
> the key follows existing kebab-case model identifier conventions. — 0085

### Accountability-preserving definition

The authoring guide, setup workflow, and README **MUST** define Agent
Harnessability as the degree to which the project's own materials and tooling
equip an AI agent to understand the project, take scoped work, operate the
environment, verify its output, and stay safely bounded while preserving clear
human direction, review, and accountability.

Where concise model YAML is needed, the guidance **SHOULD** use this definition:

```yaml
description: >
  The degree to which the project's checked-in materials, tools, workflows,
  feedback signals, standards, and action limits equip an AI agent to
  understand the project, take scoped work, operate the environment, verify its
  output, and stay safely bounded while preserving clear human direction,
  review, and accountability.
```

The guidance **MUST NOT** describe the factor's purpose as eliminating human
responsibility, and **SHOULD NOT** rely on "limited human attention," "largely
unsupervised," or "synchronous supervision" as the main shorthand for the factor.
It **MAY** discuss reducing avoidable re-specification, review toil, or
unnecessary supervision when human accountability remains explicit.

> Rationale: Agent Harnessability should make agent work more governable, not
> imply that responsibility has moved from the human/project owner to the agent.
> — 0085

### Preserve the 0081 modeling shape

The guidance **MUST** keep Agent Harnessability as the model-wide factor
projection of the agent-collaboration concern for an agent-collaborated composite
root. It **MUST** keep the **agent harness** as the constituent projection and the
**agent** as the audience projection.

The guidance **MUST** keep the existing six sub-factors unless a separate change
case changes that decomposition:

- `agent-accessibility`
- `task-specifiability`
- `agent-operability`
- `self-verifiability`
- `enforcement-of-standards`
- `containment-of-action`

The guidance **MUST** keep the existing double-count boundary: Agent
Harnessability rates how each constituent equips an agent; the agent-harness area
rates the steering artifact's own quality; Agent Harnessability is not assessed on
the agent-harness area as a recursion of the same evidence.

> Rationale: this change corrects naming and pedagogy, not the factor/constituent
> model introduced by 0081. — 0085

### Setup and maturity routing

Setup **MUST** propose `agent-harnessability` / Agent Harnessability by default
for an agent-collaborated composite root. It **MUST NOT** omit the factor because
the project's harness is thin or absent; thinness or absence remains rating
evidence, not an omission reason.

The Top 10 check and setup maturity guidance **MUST** refer to Agent
Harnessability when checking whether an agent-collaborated composite root carries
the factor or its sub-factors as model-wide factors.

### Legacy `harnessability` handling

The skill guidance **SHOULD** treat an existing `harnessability` factor with the
0081 sub-factor shape as semantic coverage for the model-wide concern, so a model
is not misclassified as missing the concern solely because it uses the previous
key or title.

When the skill is already authoring or revising the model, it **SHOULD** recommend
renaming that factor to `agent-harnessability` / Agent Harnessability unless the
project has an explicit reason to preserve the older key.

> Rationale: factor keys are project-authored identifiers, not schema fields.
> Existing models should remain usable, while new guidance converges on the more
> intelligible name. — 0085

### Public docs and changelog

The README's "Evaluate and Improve Agent Harnessability" section **MUST** name
Agent Harnessability in the body, not plain harnessability on first use, and
**MUST** use accountability-preserving wording.

The unreleased CHANGELOG entry for the 0081 harnessability work **MUST** be
revised to describe Agent Harnessability and the new definition.

Historical Change Cases, archived design/spec files, and append-only historical
logs **MUST NOT** be mass-rewritten. A current implementation log entry may record
the rename when the durable spec mirrors are updated.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` - rename the factor guidance
  to Agent Harnessability, require `agent-harnessability` as the recommended key,
  preserve the 0081 sub-factor shape, and add the accountability-preserving
  definition (per the human-facing name, definition, and modeling-shape
  requirements above).
- `specs/skills/quality-skill/workflows/setup.md` - require setup to propose
  `agent-harnessability` / Agent Harnessability and use the new definition (per
  the setup and maturity routing requirement above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - update the
  missing-factor check to name Agent Harnessability and recognize legacy
  `harnessability` coverage as stale naming, not absence (per the legacy handling
  requirement above).
- `specs/skills/quality-skill/guides/log.md` and
  `specs/skills/quality-skill/workflows/log.md` - record the durable guide and
  workflow revisions when implemented (per the public docs and changelog
  requirement above).

### To rename

None

### To delete

None
