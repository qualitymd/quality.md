---
type: Change Case
title: Stewardship vocabulary discipline
description: Keep the stewardship/care core language in its motivation register so it never modifies or replaces taxonomy nouns (factor, area, requirement, constituent, audience) — fix the "stewardship lenses" fusions, guard the setup output, and record the rule in the guide, its durable spec, and AGENTS.md.
status: Done
tags: [skill, guide, authoring, setup, stewardship, vocabulary, terminology]
timestamp: 2026-06-24T00:00:00Z
---

# Stewardship vocabulary discipline

A **Change Case** that draws a clean line between two registers the recent
stewardship/care work introduced and the established model taxonomy. The
philosophical core language — *stewardship, care, tending, vulnerability,
concern* — explains *why* a concern exists and what it means to tend an entity.
The taxonomy — *Factor, Area, Constituent, Requirement, Audience* — names the
slots in the Model. This case keeps the first register from modifying or
standing in for the second, so the core language can stay without eroding the
terms of art it sits beside.

Like [0076](0076-domain-constituent-kinds.md) and
[0077](0077-stewardship-care-grounding.md), this case changes guidance
vocabulary and framing only — no format semantics, no schema, no CLI, no
evaluation behavior, and no removal of the stewardship/care grounding those
cases established.

Detail lives in:

- [Functional spec](0079-stewardship-vocabulary-discipline/spec.md) — what the
  guidance must say and which fusions it must remove.
- [Design doc](0079-stewardship-vocabulary-discipline/design.md) — where the
  register rule lives (AGENTS.md canonical; the guide and durable spec
  operational), the exact rephrasing of the two fusions, and why the setup guard
  is operational rather than a spec contract.

## Motivation

0076 and 0077 gave the authoring guide a strong philosophical grounding: a
stewardship concern is *care* — an activity of tending — that *projects into* the
model as a factor, a constituent, or an audience (the three-projections rule).
That grounding is good and load-bearing. But the new vocabulary has started to
leak across the projection boundary it defines:

- The guide twice calls the root's recurring factors "**stewardship lenses**"
  (`authoring.md:109`, `:878`). "Lens" is the guide's own gloss for *factor*
  (`authoring.md:792`), so "stewardship lenses" reads as "stewardship factors" —
  it makes *stewardship* a kind of factor, when the three-projections rule says a
  stewardship concern is the *source* a factor projects from, not a factor itself.
- A real `/quality setup` run reproduced the fusion verbatim, reporting
  "**stewardship factors** (maintainability, consistency, currentness,
  traceability)." The agent reached for the nearest grouping word and demoted the
  taxonomy term to a subcategory of the new philosophical one.

This is exactly the substitution to avoid: a freshly introduced word quietly
displacing an established term of art weakens shared communication — the reason
the taxonomy is fixed vocabulary in the first place. The fix is not to retract
stewardship/care but to confine it to its register. The three-projections rule
already implies the discipline (a concern is not a factor); this case states it
as a rule and removes the two fusions that violate it, so the core language
clarifies the taxonomy instead of overwriting it.

## Scope

Covered:

- State the register rule once, where authors and the agent will see it: the
  motivation-layer vocabulary (stewardship, care, tending, vulnerability,
  concern) **describes why** and **never modifies or replaces** a taxonomy noun
  (factor, area, requirement, constituent, audience). The root's recurring
  factors are named as **model-wide / cross-cutting factors** (the guide's
  existing terms), which **may** be noted as tracing to stewardship concerns —
  without making "stewardship" an adjective on the taxonomy noun.
- Remove the two "stewardship lenses" fusions in `skills/quality/guides/authoring.md`
  (`:109`, `:878`), rephrasing each as model-wide/cross-cutting factors that
  recur across constituents and trace to stewardship concerns.
- Add an output-shaping guard to the setup workflow so the model summary names
  factors as factors (or model-wide factors), not "stewardship factors."
- Add a vocabulary-discipline rule to `AGENTS.md` alongside the existing
  QUALITY.md vocabulary conventions.
- Align the durable authoring guide spec with a matching requirement and promote
  the rationale; align the durable setup spec if it shapes the summary output.
- Record the guide/spec updates in the guides log and add a CHANGELOG note.

Deferred / non-goals:

- **No retraction of the stewardship/care grounding.** 0076's generator, 0077's
  care-grounding, the nine concerns, the two axes, and the three-projections rule
  all stand unchanged. This case adds a register guardrail; it removes nothing
  substantive.
- The legitimate gloss "a factor is a quality *lens*" (singular, defining what a
  factor is) stays — only "stewardship lens(es)" as a *substitute name* for the
  recurring factors is the target.
- No QUALITY.md format change and no `SPECIFICATION.md` change.
- No new concern, no change to the three-projections rule, the nine concerns, or
  the two axes.
- No change to evaluation, reporting, schema, or CLI behavior.

## Affected artifacts

Derived by sweeping the repo for the fusion phrasing (`stewardship (factor|lens|
lenses|area|requirement)`) and the stewardship/care vocabulary across every
bundle; empty kinds are deliberate.

### Code

- [x] None — documentation-only change.

### Durable specs

- [x] `specs/skills/quality-skill/guides/authoring-md.md` — add a register-
      discipline clause to the domain-constituent-kinds / three-projections
      contract: the motivation-layer vocabulary must not modify or replace a
      taxonomy noun, and recurring root factors are named as model-wide factors.
      Promote the rationale (0079).
- [x] `specs/skills/quality-skill/workflows/setup.md` — **no change (assessed).**
      The recap/closeout set no factor-naming contract; R4 is satisfied by an
      operational guard in the bundled workflow plus the source fix. See the
      [design doc](0079-stewardship-vocabulary-discipline/design.md).
- [x] `specs/skills/quality-skill/guides/log.md` — record the guide-spec update.

### Format spec

- [x] None — no change to `SPECIFICATION.md`.

### Durable docs (AGENTS.md and bundled skill)

- [x] `AGENTS.md` — added the motivation-vs-taxonomy register rule as a new
      subsection ("Keep the motivation and taxonomy registers distinct") under the
      QUALITY.md vocabulary conventions (symlinked from `CLAUDE.md`/`GEMINI.md`;
      edited `AGENTS.md` only).
- [x] `skills/quality/guides/authoring.md` — removed the two "stewardship lenses"
      fusions (`:109`, `:878`) and added an Avoid bullet stating the register rule
      at the three-projections rule.
- [x] `skills/quality/workflows/setup.md` — added the output-shaping guard so any
      recap names factors as factors / model-wide factors.

### Release

- [x] `CHANGELOG.md` — added the `/quality Skill` note under `Unreleased`.

## Children

- [Functional spec](0079-stewardship-vocabulary-discipline/spec.md) — required
  guidance content and acceptance criteria.
- [Design doc](0079-stewardship-vocabulary-discipline/design.md) — placement of
  the register rule, the fusion rephrasings, and the alternatives weighed.

## Status

`Done`. Implementation complete and the Affected artifacts list reconciled.
Landed: the AGENTS.md register-rule subsection; the two authoring-guide
"stewardship lenses" rephrasings plus an Avoid bullet at the three-projections
rule; the durable authoring-guide spec clause with a 0079 rationale; the
operational setup-workflow guard; and the guides-log and CHANGELOG notes. The
durable setup spec was assessed no-change. Documentation-only — no `SPECIFICATION.md`
or code change. Verified with `mise run check` (markdown format, bundle link
resolution, Go vet/lint/test all pass). Archived to `archive/` on landing.
