---
type: Change Case
title: Composite root areas and use-context constituents
description: Name the composite root shape and the recurring agent-harness and QUALITY.md self-check constituents in the authoring guidance, and correct the root-factor-count heuristic.
status: Done
tags: [skill, guide, authoring, areas, factors]
timestamp: 2026-06-24T00:00:00Z
---

# Composite root areas and use-context constituents

A **Change Case** to correct a structural assumption in the authoring guidance:
that a root area is a single primary subject with one factor family. For most
real entities in QUALITY.md's assumed context of use the root is a **composite**
— one whole decomposed into distinct _kinds_ of constituent artifacts, each with
its own largely-disjoint factor family. This case names the composite root shape
(distinct from a homogeneous collection/grouping area), names the two
constituents that recur regardless of modeled domain because of QUALITY.md's
agentic context of use (the **agent harness** and the **QUALITY.md self-check**),
and fixes the "aim for ~ten root-level factors" heuristic so it applies per
primary subject rather than at a composite root.

Detail lives in:

- [Functional spec](0074-composite-root-areas/spec.md) - what the change must do.
- [Design doc](0074-composite-root-areas/design.md) - the taxonomy and why it is
  shaped this way.

## Motivation

The authoring guide's "Cover the domain's stable stakes before specializing"
subsection tells authors to aim for roughly ten root-level factors. That target
assumes the root is a single primary subject with one conventional factor family.
For a real project that assumption is usually wrong and produces a structural
error: it silently equates the root with one constituent (typically "the
product") and either drops the other high-leverage artifacts — requirements,
docs, governing specs — or jams incompatible factor families into one flat root
list.

A real evaluated entity is usually a **composite**: unlike, interdependent parts
of different kinds, each with its own near-disjoint factor family. The guide's
existing split test ("split off a child only when it has distinct factors or
requirements") already implies this, but the guide never names the composite
shape, so authors default to a flat, product-factored root. A near-disjoint
factor family is in fact the _strongest_ form of the split signal.

Two facts compound the gap:

- **Domain agnosticism cuts across it.** Which _domain_ constituents exist varies
  by what is modeled (software vs. dataset vs. document); the composite shape
  itself is domain-general.
- **Context of use is invariant.** QUALITY.md's assumed context of use is an
  agent/AI-assistant-collaborated project. That makes two constituents recur
  _regardless of modeled domain_: the agent harness and the QUALITY.md
  self-check. They belong to the use context, not the modeled domain — the exact
  distinction [`AGENTS.md`](../../AGENTS.md) already draws (use context vs. model
  domain), now applied to area structure.

## Scope

Covered:

- Name three decomposition shapes in the authoring guide — **primary-subject**,
  **collection**, and **composite** — the test that selects each, and the fact
  that they are recursive and composable (a composite area can hold collection or
  composite children to any depth), not a one-time root classification.
- Define the composite root: unlike interdependent constituents, each with its
  own factor family; its emergent concern is cross-part coherence (the assessment
  edges), distinct from a collection's set-level concerns.
- Name the recurring **use-context constituents** (agent harness, QUALITY.md
  self-check), separating them from domain constituents, and state the two
  roll-up asymmetries: the self-check stays on the learn-loop axis, out of the
  entity roll-up; the harness is partly normative and plays the dual
  area/assessment-reference role.
- State the expected-default-not-quota guardrail for use-context constituents.
- Correct the root-factor-count heuristic to apply per primary subject; at a
  composite root it applies per constituent, with only recurring stewardship
  factors (refined down) at the root.
- Add the matching routing findings to the Top 10 checks (Check 8).

Deferred / non-goals:

- No QUALITY.md format change and no `SPECIFICATION.md` change; "composite",
  "collection", and "primary-subject" stay authoring-guide vocabulary, not format
  semantics.
- No new mandatory areas; the harness and self-check stay earned, not required.
- No change to evaluation, reporting, or CLI behavior.

## Affected artifacts

### Code

- [x] None — documentation-only change.

### Durable specs

- [x] `specs/skills/quality-skill/guides/authoring-md.md` - added best-practice
      coverage requirements for the three decomposition shapes, the recurring
      use-context constituents, and the per-primary-subject factor-coverage aim,
      with a 0074 rationale note.
- [x] `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` -
      extended the area-and-factor-shape check with the composite-flattening and
      missing-constituent findings and the earned-not-roster qualifier.
- [x] `specs/skills/quality-skill/guides/log.md` - recorded the guide-spec
      updates.

### Format spec

- [x] None — no change to `SPECIFICATION.md`.

### Durable docs (bundled skill)

- [x] `skills/quality/guides/authoring.md` - added the "Choose the decomposition
      shape" subsection in "Area", added the "Carry the recurring use-context
      constituents" subsection after "Ground high-leverage concerns in normative
      artifacts", caveated the factor-count heuristic in "Cover the domain's
      stable stakes before specializing", and tied "Keep the root lean" to the
      precise meaning of model-wide factors.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` - extended Check 8 with
      the composite-flattening and missing-expected-constituent findings and the
      earned-not-quota clarifier.

### Release

- [x] `CHANGELOG.md` - added the `/quality Skill` note under `Unreleased`.

## Children

- [Functional spec](0074-composite-root-areas/spec.md) - required guidance
  content and acceptance criteria.
- [Design doc](0074-composite-root-areas/design.md) - the taxonomy, alternatives,
  and trade-offs.

## Status

`Done`. Applied the authoring-guide and Top 10 edits, aligned both guide specs
and the guides log, and added the CHANGELOG note. Documentation-only (no
`SPECIFICATION.md` or code change). Verified with `mise run check` (markdown
format, bundle link resolution, Go gates all pass). Archived.
