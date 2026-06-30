---
type: Design Doc
title: Stewardship vocabulary discipline — design
description: Where the motivation-vs-taxonomy register rule lives, how the two "stewardship lenses" fusions are rephrased, and why the setup guard is operational rather than a spec contract.
tags: [skill, guide, authoring, setup, stewardship, vocabulary, terminology]
timestamp: 2026-06-24T00:00:00Z
---

# Stewardship vocabulary discipline — design

## Context

Implements the
[Stewardship vocabulary discipline](../0079-stewardship-vocabulary-discipline.md)
change case and its [functional spec](spec.md). The spec requires the
stewardship/care vocabulary to stay in its motivation register and never modify
or replace a taxonomy noun (R1–R2), removes the two "stewardship lenses" fusions
in the authoring guide (R3), guards the setup model summary (R4), and preserves
the 0076/0077 grounding (R5). This doc records _where_ each lands and _why_ that
placement.

## Approach

The change is vocabulary discipline, not a code path, so the design is a set of
placement and wording decisions.

**Three homes, one rule, said once (R1, R6).** The register rule is a repo-wide
convention with one canonical statement and two operational echoes:

- `AGENTS.md` → canonical. A new bullet group under **QUALITY.md vocabulary
  capitalization** states the rule for every contributor and agent: the
  stewardship/care vocabulary (stewardship, care, tending, vulnerability,
  concern) is motivation-layer and never modifies or replaces a taxonomy noun
  (factor, area, requirement, constituent, audience); name recurring root factors
  as model-wide / cross-cutting factors. This sits next to the existing
  capitalization conventions because it is the same kind of rule — how to write
  the vocabulary, repo-wide.
- `skills/quality/guides/authoring.md` → operational. The rule is stated where
  the projection boundary it depends on is already defined: the **three-
  projections rule** (~`:662–671`). A new **Avoid** bullet there says not to let
  a motivation-layer word modify a taxonomy noun, naming "stewardship factor /
  stewardship lens" as the failure. This is the spot the agent reads while
  authoring a model, so it is where the discipline bites.
- `specs/skills/quality-skill/guides/authoring-md.md` → durable contract. The
  three-projections requirement (~`:253`) gains a clause that the guide MUST keep
  the motivation-layer vocabulary from modifying or replacing a taxonomy noun and
  MUST name recurring root factors as model-wide factors, with a 0079 rationale
  annotation.

**The two fusions (R3).** Both occurrences keep their meaning and lose the fused
adjective:

- `authoring.md:109` — "often stewardship lenses like currentness or
  traceability" → "often factors that trace to stewardship concerns, like
  currentness or traceability". `factor` is restored as the noun; the stewardship
  link becomes a _trace to_, per R2.
- `authoring.md:878` — "typically stewardship lenses (currentness, traceability,
  consistency, maintainability)" → "typically those tracing to stewardship
  concerns (currentness, traceability, consistency, maintainability)". "those"
  refers back to "the factors that recur across constituents" in the same
  sentence, so the taxonomy noun is already in hand.

The singular gloss "a factor is a quality _lens_" (`:792`, `:102`, `:844`,
`:885`) is untouched (R5) — it _defines_ factor rather than substituting a
philosophical word for it.

**Setup guard is operational, not a new spec contract (R4).** The setup closeout
is status-first and terse and does not enumerate factors; the observed
"stewardship factors" output came from the agent's free-form recap, learned from
the guide's "stewardship lenses" wording. Fixing the guide + AGENTS.md removes
the source. A one-line reminder in `skills/quality/workflows/setup.md` (near the
candidate-model-shape / recap guidance) names factors as factors / model-wide
factors in any recap and cross-references the rule. The durable setup spec
(`specs/skills/quality-skill/workflows/setup.md`) is assessed **no-change**: it
sets no factor-naming contract for the recap, so the guard is operational phrasing
finer than the spec, not a new MUST.

## Alternatives

- **Drop "lens" as a factor gloss entirely.** Rejected: the singular "quality
  lens" gloss is well-established and clear, and defines what a factor _is_. Only
  the fused plural "stewardship lenses" — a philosophical adjective on the
  taxonomy noun — is the problem.
- **Forbid "stewardship" anywhere near "factor."** Rejected: "factors that trace
  to stewardship concerns" is legitimate and useful; it expresses the projection
  the three-projections rule defines. The rule targets _modifying or replacing_
  the noun, not co-occurrence.
- **Put the rule only in AGENTS.md.** Rejected: the leak originated in the
  authoring guide and surfaced in agent output; the operational echo at the
  three-projections rule is where the discipline is applied, so the guide must
  carry it.
- **Add a new templated "what the model captures" summary to setup with fixed
  factor labels.** Rejected as scope creep (YAGNI): no such template exists today,
  and inventing one to fix a wording leak is heavier than the operational guard
  plus the source fix.

## Trade-offs & risks

- Stating the rule in three places risks drift. Mitigated by say-it-once: AGENTS.md
  is canonical and full; the guide and durable spec state the operational form and
  cross-reference rather than re-deriving it.
- The setup guard is guidance, not an enforced contract, so a future free-form
  recap could still fuse the registers. Accepted: the source fix (guide +
  AGENTS.md) addresses the cause; a spec-level factor-naming contract for a terse
  recap would be over-specification.

## Open questions

None.
