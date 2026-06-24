---
type: Functional Specification
title: Stewardship vocabulary discipline — functional spec
description: Required guidance that confines the stewardship/care core language to its motivation register so it never modifies or replaces taxonomy nouns, removes the "stewardship lenses" fusions, and guards the setup model summary.
tags: [skill, guide, authoring, setup, stewardship, vocabulary, terminology]
timestamp: 2026-06-24T00:00:00Z
---

# Stewardship vocabulary discipline — functional spec

Companion to the
[Stewardship vocabulary discipline](../0079-stewardship-vocabulary-discipline.md)
change case. This spec states *what* the guidance must say and which fusions it
must remove; the motivation lives in the change case.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 (RFC 2119, RFC 8174) when, and only when, they
appear in all capitals.

## Scope

Covers register discipline for the stewardship/care vocabulary across:

- the bundled authoring guide
  [`skills/quality/guides/authoring.md`](../../../skills/quality/guides/authoring.md)
  and its governing durable spec
  [`specs/skills/quality-skill/guides/authoring-md.md`](../../../specs/skills/quality-skill/guides/authoring-md.md);
- the bundled setup workflow
  [`skills/quality/workflows/setup.md`](../../../skills/quality/workflows/setup.md)
  and its governing durable spec
  [`specs/skills/quality-skill/workflows/setup.md`](../../../specs/skills/quality-skill/workflows/setup.md)
  (where it shapes the model summary);
- the repository vocabulary conventions in
  [`AGENTS.md`](../../../AGENTS.md).

**Deferred / non-goals:** any change to `SPECIFICATION.md`, the format, the
schema, the CLI, or evaluation/reporting behavior; any retraction of the
[0076](../0076-domain-constituent-kinds.md)/[0077](../0077-stewardship-care-grounding.md)
stewardship grounding; any change to the nine concerns, the two axes, or the
three-projections rule; and removal of the legitimate singular gloss "a factor is
a quality *lens*."

## Background / Motivation

0076 and 0077 grounded the authoring guide in the phenomenology of care: a
stewardship concern is an activity of tending that *projects into* the model as a
factor, a constituent, or an audience. The grounding is sound, but its vocabulary
crossed the projection boundary it defines. The guide twice names the root's
recurring factors "stewardship lenses" (`authoring.md:109`, `:878`); because
"lens" is the guide's gloss for *factor*, that phrasing makes *stewardship* a
kind of factor rather than the concern a factor projects from. A live
`/quality setup` run reproduced the fusion as "stewardship factors," demoting an
established term of art to a subcategory of a newly introduced philosophical one.
A fixed taxonomy exists to keep communication consistent; a fresh word silently
displacing it defeats that purpose. This spec confines the core language to its
register without weakening the grounding.

## Requirements

### Two registers, one rule

- R1. The guidance **MUST** state that the stewardship/care vocabulary —
  *stewardship, care, tending, vulnerability, concern* — is **motivation-layer**:
  it describes *why* a concern exists and what it means to tend an entity. It
  **MUST** state that this vocabulary **MUST NOT** modify or stand in for a
  **taxonomy noun** — *factor, area, requirement, constituent, audience* — which
  names a slot in the Model.

  > Rationale: a stewardship concern *projects into* a factor/constituent/audience
  > (the three-projections rule); it is not itself one. "Stewardship factor" or
  > "stewardship lens" collapses the source into the projection and demotes a term
  > of art to a subcategory of the philosophical word. — 0079

- R2. The guidance **MUST** name the root's recurring factors as **model-wide**
  (or cross-cutting) **factors** — the guide's existing terms. It **MAY** note
  that such factors *trace to* stewardship concerns, but **MUST NOT** render that
  link by making "stewardship" (or another motivation-layer word) an adjective on
  the taxonomy noun.

### Remove the existing fusions

- R3. The authoring guide **MUST NOT** call factors "stewardship lenses." The two
  current occurrences (`skills/quality/guides/authoring.md:109` and `:878`)
  **MUST** be rephrased as model-wide / cross-cutting factors that recur across
  the root's constituents, with the stewardship link expressed per R2.

- R4. The setup workflow's model summary **MUST** name factors as *factors* (or
  *model-wide factors*), and **MUST NOT** label them "stewardship factors" or use
  another motivation-layer word as an adjective on a taxonomy noun.

### Preserve the grounding

- R5. The refinements **MUST** preserve the
  [0076](../0076-domain-constituent-kinds.md)/[0077](../0077-stewardship-care-grounding.md)
  grounding intact: the nine concerns, the two axes, the three-projections rule,
  the earn-it guardrail, and the use of *stewardship concern* / *care* / *tending*
  / *vulnerability* as motivation-layer nouns and verbs all stand. The singular
  gloss "a factor is a quality *lens*" (defining what a factor is) is preserved;
  only "stewardship lens(es)" as a substitute name for factors is removed.

### Consistency

- R6. All refined content **MUST** preserve domain agnosticism per
  [`AGENTS.md`](../../../AGENTS.md) (illustrative, non-exhaustive examples;
  agentic language scoped to context of use; software product quality not implied
  as the default modeled domain), and the register rule **MUST** be stated once
  per home and cross-referenced rather than restated in full (say it once).

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` — add a register-discipline
  clause to the three-projections / domain-constituent-kinds contract: the
  motivation-layer vocabulary must not modify or replace a taxonomy noun, and
  recurring root factors are named as model-wide factors (per R1–R3).
- `specs/skills/quality-skill/workflows/setup.md` — **no change (assessed).** The
  setup recap and closeout set no factor-naming contract (the closeout is terse
  and status-first and does not enumerate factors), so R4 is satisfied by an
  operational guard in the bundled workflow plus the source fix in the authoring
  guide and `AGENTS.md`; a spec-level factor-naming contract for a terse recap
  would be over-specification (see the [design doc](design.md)).

### To rename

None.

### To delete

None.
