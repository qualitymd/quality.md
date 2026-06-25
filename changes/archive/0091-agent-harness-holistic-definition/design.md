---
type: Design Doc
title: Agent-harness holistic definition - design
description: How the agent-harness guidance is re-thought around a holistic harness definition - the concept-first framing, the three-projection model, the area-scoping decision rule, the redefine-not-rename decision, the feedforward/feedback vocabulary choice, and where each edit lands.
tags: [docs, doctrine, skill, agent-harness]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-harness holistic definition - design

Answers the [functional spec](spec.md). Settles *how* the harness definition is
re-thought and where the edits land. Builds on
[0089](../../archive/0089-agent-harness-modeling-guidance.md)'s projection split
and [0087](../../archive/0087-encode-projection-boundaries.md)'s boundary rule;
neither is reopened.

## The grounding

The definition follows the now-common account of an agent harness:

- **Agent = Model + Harness.** The harness is *everything that isn't the model
  itself* — code, configuration, and execution logic.
- The harness has two halves: **feedforward controls** (guides) that steer the
  agent before it acts, and **feedback controls** (sensors) that observe and
  correct after. Tools, runtime/sandbox, memory/state, orchestration, and
  verification are all harness; instructions are one feedforward component.

0089 already cited this body of practice to justify the projection split and to
add `continuity` and the good-sensor properties. This case applies the same
grounding to the one place 0089 left narrow — the *definition* of the harness in
the constituent guide.

## Concept-first framing

The current guide commits a definitional inversion: it defines the *concept*
("the agent harness is…") by the *area's* contents ("…the instructions that steer
the agent"). The fix inverts that:

1. Define the harness at full breadth (the concept).
2. Derive the area as one honest projection of it (the scope).

So `agent-harness.md` opens with "What the agent harness is" (the system), then
"The harness is projected, not rated as one object" (three projections), then
"Scoping the agent-harness area" (the project-owned slice). The area becomes a
*stated projection with named boundaries* instead of a silent redefinition.

## The three-projection model

| Harness facet                                                                                    | Rated by                           |
| ------------------------------------------------------------------------------------------------ | ---------------------------------- |
| Capability — how well the whole harness equips an agent                                          | Agent Harnessability **factor**    |
| Checked-in, project-owned governing artifacts (guides + owned tools/hooks/sandbox/orchestration) | agent-harness **area**             |
| Verification corpus (sensors as a product asset)                                                 | tests area (cross-referenced)      |
| Runtime environment                                                                              | operations area (cross-referenced) |

This is what keeps the holistic definition from causing double-counting: the
harness *is* broad, but each part has exactly one home. The area owns the
governing artifacts the project authors; the rest is cross-referenced.

## The area-scoping decision rule

For an artifact that could sit in more than one projection:

> Rate it in the agent-harness area when its primary job is to **govern or equip
> the agent's work** and the project owns it. Cede it to a domain constituent when
> it is primarily a **product artifact the agent merely also uses** (product test
> suite → tests; deploy runtime → operations). When one artifact does both — an
> agent-run verify hook, a permission gate wired into the loop — rate its
> agent-governing quality here and cross-reference, under the no-double-count rule.

This makes project-owned runtime machinery (hooks, sandbox/permission policy,
orchestration, tool/MCP definitions) explicitly the area's business — closing the
gap where 0089's "steering materials" framing let such machinery be folded into
prose or dropped.

## Decisions and alternatives

### Redefine, not rename (chosen)

The triggering confusion — an "Agent Harness" area scoped to instructions — could
be solved two ways: rename the area to "Agent Steering Materials," or widen the
definition. **Rejected the rename.** "Agent harness" is the field's correct term
for the system; renaming the area would cede the right word to the wrong (narrow)
meaning and leave the field's term unused. Widening the definition keeps the right
word and fixes the actual error (the part-for-whole definition). The area stays
`agent-harness` in key, title, and concept.

### Broaden the area to the whole harness (rejected)

Could make the area literally "everything but the model." Rejected: it would
swallow the tests and operations constituents and overlap the Agent Harnessability
factor, reintroducing the double-counting 0087/0089 removed. The three-projection
model is the right shape — broad *concept*, scoped *area*.

### Vocabulary: feedforward/feedback vs. proper nouns (chosen: plain terms)

The guide uses "feedforward controls / feedback controls" rather than a specific
author's "Guides / Sensors" proper nouns, keeping the doctrine source-agnostic and
consistent with the guide's voice. A short provenance note at the foot records the
"Agent = Model + Harness" lineage without pinning the doctrine to one article.
0089's runtime guides already use "guides/sensors" descriptively; this case keeps
that usage but leads with the plain feedforward/feedback framing.

## Where the edits land

- **`skills/quality/guides/authoring/agent-harness.md`** — the primary rewrite:
  concept-first definition, three-projection section, scoping decision rule,
  owned-machinery scope, expanded requirement shapes, doctrine `Do`.
- **`skills/quality/guides/authoring/model-structure.md`** — the "Carry the
  recurring use-context constituents" definition fix (one sentence, plus keep the
  routing).
- **`skills/quality/workflows/setup.md`** — generated area `description`,
  projection-boundary comment, recap line, and an active owned-machinery check.
- **`skills/quality/guides/top-10-quality-md-checks.md`** — check-8 additions.
- **Durable mirrors** under `specs/skills/quality-skill/guides/authoring/`,
  `.../workflows/setup.md`, and `.../guides/top-10-quality-md-checks-md.md`.
- **Logs / changelog** — `skills/quality/guides/log.md`,
  `specs/skills/quality-skill/guides/log.md`, and `CHANGELOG.md`; `docs/` only if
  the doctrine guide's "Agentic use context" section restates the narrow
  definition.

The Agent Harnessability factor guide (`agent-harnessability.md`) needs at most a
cross-reference; its definition is already holistic and stays as-is.
