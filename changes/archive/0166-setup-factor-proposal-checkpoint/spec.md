---
type: Functional Specification
title: Setup Factor Proposal Checkpoint — functional spec
description: Requirements for teaching and applying factor desiderata during setup.
tags: [skill, quality, setup, factors]
timestamp: 2026-06-29T00:00:00Z
---

# Setup Factor Proposal Checkpoint — functional spec

Companion to
[Setup Factor Proposal Checkpoint](../0166-setup-factor-proposal-checkpoint.md).
This spec states the delta for runtime `/quality setup` guidance, factor
authoring guidance, durable skill specs, and README framing.

The durable source of truth is absorbed into
[`/quality setup`](../../../specs/skills/quality-skill/workflows/setup.md),
[`QUALITY.md factor authoring guide`](../../../specs/skills/quality-skill/guides/authoring/factors.md),
the parent [`/quality` skill spec](../../../specs/skills/quality-skill/quality-skill.md),
the runtime setup workflow, the runtime factor authoring guide, and `README.md`.

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Initial `QUALITY.md` setup asks good discovery questions, but factor selection is
where the model's judgment becomes concrete. Without an explicit factor proposal
checkpoint, the user can approve a model without seeing whether the proposed
factor set covers what matters, emphasizes what matters most, and remains usable.
The same checkpoint should teach individual-factor quality so feedback is better
than "looks fine": a user can name a missing factor, an over-modeled concern, a
misplaced factor, or a poor name.

This change makes setup use and teach two-layer factor desiderata. The set-level
qualities are comprehensive, proportionate, and sustainable. The
individual-factor qualities are consequential, bounded, operational, traceable,
and neutral.

## Scope

Covered:

- factor desiderata in the runtime and durable factor authoring guide contract;
- a setup factor proposal checkpoint after discovery and before final review;
- setup-brief fields, final recap fields, model-authoring guidance, and close-gap
  checks that consume the proposal;
- public README framing.

Non-goals:

- CLI behavior, schema changes, or generated scaffolds;
- a complete factor-selection algorithm;
- domain-specific default factor families.

## Requirements

1. The factor authoring guide **MUST** teach factor quality at two layers:
   Factor sets are comprehensive, proportionate, and sustainable; individual
   Factors are consequential first, then bounded, operational, traceable, and
   neutral.

   > Rationale: The set-level and individual-factor qualities answer different
   > questions. Combining them into one checklist makes setup likely to drop hard
   > but important concerns or over-model easy ones.
   > Durable spec: modify `specs/skills/quality-skill/guides/authoring/factors.md`.

2. The factor authoring guide **MUST** state that a consequential Factor that is
   weak on the refinement qualities should usually be improved rather than
   dropped.

   > Rationale: Important factors are often initially hard to bound, inspect, or
   > name; setup should refine those aspects instead of biasing toward concerns
   > that are merely easy to assess.
   > Durable spec: modify `specs/skills/quality-skill/guides/authoring/factors.md`.

3. The setup workflow **MUST** build working context for candidate Factor-set
   quality and candidate Factor rationales, alongside the existing setup brief.

   > Rationale: The user-facing proposal should come from deliberate analysis,
   > not a presentational table assembled after the fact.
   > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

4. After discovery questions and the human context checkpoint, setup **MUST**
   present a factor proposal checkpoint before the final review gate.

   > Rationale: Users need to see and correct the proposed factor set before
   > authoring, but after setup has enough context to propose something concrete.
   > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

5. The factor proposal checkpoint **MUST** teach the set-level and
   individual-factor desiderata in compact user-facing copy and **MUST NOT** ask
   the user to design Factors, child Areas, Requirements, or YAML cold.

   > Rationale: The checkpoint exists to improve feedback quality, not to push
   > modeling labor onto the user.
   > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

6. The factor proposal checkpoint **MUST** present candidate Factors by Area with
   a short rationale and visible initial depth, using `light`, `normal`, and
   `deep` as the depth labels.

   > Rationale: The setup output should make both selection and proportionality
   > reviewable. Depth communicates where the first model will spend requirement
   > and assessment rigor.
   > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

7. The factor proposal checkpoint **MUST** ask for targeted corrections about
   concerns that are missing, overemphasized, misplaced, or badly named.

   > Rationale: Those correction categories map directly to factor-set and
   > individual-factor quality without requiring the user to know the model
   > schema.
   > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

8. The final setup review recap **MUST** include the reviewed factor proposal
   before the user authorizes writing `QUALITY.md`.

   > Rationale: The review gate authorizes the actual model, so factor selection
   > must be part of the reviewed state, not an earlier transient prompt.
   > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

9. Setup-authored `QUALITY.md` models **MUST** use the reviewed factor proposal
   and depth labels when shaping starter Factors, Requirements, and assessments.

   > Rationale: The proposal checkpoint must bind authoring behavior; otherwise it
   > teaches users without changing the resulting model.
   > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

10. Setup's important-gap inspection **MUST** treat missing consequential factor
    coverage or unreasonable depth allocation as first-model usefulness gaps.

    > Rationale: The closeout should catch failure to apply the same desiderata
    > setup just taught.
    > Durable spec: modify `specs/skills/quality-skill/workflows/setup.md`.

11. The public README **SHOULD** include a concise explanation that QUALITY.md
    does not import default factor checklists and that factors are selected
    deliberately using set-level and individual-factor qualities.

    > Rationale: README readers should see the principle before reaching runtime
    > skill guidance, but the README should not become the full authoring guide.
    > Durable spec: none.

## Durable spec changes

To add:

- None.

To modify:

- `specs/skills/quality-skill/guides/authoring/factors.md` - require the
  two-layer factor desiderata and refine-not-drop rule.
- `specs/skills/quality-skill/workflows/setup.md` - require the factor proposal
  checkpoint and the setup surfaces that consume it.
- `specs/skills/quality-skill/quality-skill.md` - mention the checkpoint in the
  setup workflow summary.

To rename:

- None.

To delete:

- None.

## Verification

- Inspect runtime and durable guidance for the factor desiderata and setup
  checkpoint.
- Run Markdown formatting checks.
- Confirm the Change Case is archived with logs and indexes updated.
