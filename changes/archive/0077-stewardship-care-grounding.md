---
type: Change Case
title: Care-grounded stewardship concerns
description: Refine the authoring guide's stewardship-concern generator with the phenomenology of care — frame the artifact as evidence of tending rather than the care itself, locate the claim that earns a constituent in a Need or Risk rather than the list, ground the protective pair in vulnerability, define stewardship as holding-in-trust, and present the lifecycle band as recurring rather than a once-through pipeline.
status: Done
tags: [skill, guide, authoring, areas, constituents, stewardship]
timestamp: 2026-06-24T00:00:00Z
---

# Care-grounded stewardship concerns

A **Change Case** that refines the stewardship-concern generator
[0076](0076-domain-constituent-kinds.md) introduced, using the phenomenology of
care set out in two essays by B. Scot Rousse and David Spivak —
[*Notes on care*](https://withoutwhy.substack.com/p/notes-on-care-with-david-spivak-a)
and [*Clarifications of care*](https://withoutwhy.substack.com/p/clarifications-of-care).
0076 named the generator and shipped it, framed as "caring for any entity
carries a recurring set of concerns, and each tends to leave an authored,
inspectable artifact that is a candidate constituent." The care essays sharpen
that framing in five places without changing the generator's shape, the nine
concerns, or the earn-it guardrail.

The refinements are vocabulary and framing in one authoring-guide subsection and
its durable spec. Like 0076, this case changes no format semantics, no schema, no
CLI, and no evaluation behavior.

Detail lives in:

- [Functional spec](0077-stewardship-care-grounding/spec.md) — what the guidance
  must say.
- [Design doc](0077-stewardship-care-grounding/design.md) — the five refinements,
  the care concepts behind each, and the alternatives weighed.

## Motivation

0076's generator works, but it cashes "care" out as the artifact it leaves: a
concern "tends to leave an authored, inspectable artifact that is a candidate
constituent." Read literally that is artifact-first — it invites enumerating
residue (specs, tests, docs) and treating the presence of a folder as the
presence of care. The care essays name care as an activity — "tending and
attending to what matters", "doing what is required, well" — whose artifacts are
its *trace*, not its substance. That distinction, and four others the essays
draw, tighten language 0076 already half-encodes:

- **Artifact as evidence, not as care.** A thin or missing artifact does not
  prove a concern unmet; a thick one does not prove it met *well*. 0076 already
  carries the first half (an empty area with a missing-anchor finding); the
  second half — that presence is the area question and *done-well* is the factor
  question — is unstated.
- **The claim comes from the entity, not the list.** "To experience something as
  required is to find oneself addressed. The situation makes a claim." The
  legitimate source of a constituent is a Need or Risk the entity presents; the
  nine concerns are a device for *hearing* that claim. 0076's "prompt, not a
  quota" guardrail asserts this; the care framing supplies its principle.
- **The protective band is care under vulnerability.** The essays' caregiving
  mode — asymmetric responsibility toward a party that can be harmed and cannot
  fully protect itself — is exactly what `secure`/`safeguard` answer to. This
  grounds *why* the band cross-cuts the lifecycle rather than sitting inside it:
  it tracks who is exposed, not which phase the work is in.
- **"Steward" is the caregiving word.** A steward cares for what is held in
  trust, answerable to stakeholders and the future, not as private property.
  0076 uses "stewardship" without defining it; one sentence makes it load-bearing
  instead of decorative.
- **Care is open-ended.** "Always more to do." 0076 draws the lifecycle band as a
  "roughly sequential" arrow, which can imply completion. The concerns are
  standing and recurring; `operate` and `maintain` never conclude.

None of these add a concern, relax the earn-it test, or touch the format. They
make the generator read as what it is — a way to hear what an entity asks to be
cared for — instead of a checklist of artifacts to collect.

## Scope

Covered:

- Reframe the opening of the authoring-guide subsection "Cover the domain's
  constituent kinds":
  - define *steward* as caring for something held in trust, answerable to its
    stakeholders and its future;
  - state that the constituent kinds are a prompt for what *this* entity asks to
    be cared for, with the claim that earns one originating in a Need or Risk it
    presents, never in the list itself;
  - present each concern as an activity of tending whose artifact is its *trace*,
    read as evidence rather than proof of quality.
- Present the **lifecycle** band as recurring rather than a once-through
  pipeline: the order is a typical first pass, the concerns are standing, and
  `operate` and `maintain` never conclude.
- Ground the **protective** pair (`secure`, `safeguard`) in **vulnerability**,
  and state that it cross-cuts the lifecycle because it tracks who is exposed,
  not which phase the work is in — preserving the 0076 direction-of-harm
  distinction.
- Add a "Do" bullet separating the constituent question (whether a tending leaves
  an owned, inspectable artifact — an *area*) from the stewardship-quality
  question (whether it is done *well* — a *factor*).
- Align the durable authoring guide spec with matching requirements and promote
  the rationale.
- Record the guide-spec update in the guides log and add a CHANGELOG note.

Deferred / non-goals:

- No QUALITY.md format change and no `SPECIFICATION.md` change; the care-grounding
  stays authoring-guide vocabulary, as 0076's generator did.
- No new concern, no new mandatory area, and no relaxation of the earn-it
  guardrail — the care-grounding strengthens its rationale, it does not loosen
  the test.
- No formal category-theory apparatus. The source essays are phenomenological;
  the lens-like shape of the three-projections rule is noted in the design doc,
  not promoted into the guide as formalism.
- No change to the setup workflow or the Top 10 checks: their use of the
  stewardship concerns is about constituent *presence and accounting*, which
  these refinements do not alter (see Affected artifacts).
- No change to evaluation, reporting, or CLI behavior.

## Affected artifacts

### Code

- [x] None — documentation-only change.

### Durable specs

- [x] `specs/skills/quality-skill/guides/authoring-md.md` — extend the
      domain-constituent-kinds contract with the care-grounding clauses (claim
      originating in a Need or Risk rather than the list; the constituent vs.
      stewardship-quality split; the protective pair as cross-cutting stewardship
      under vulnerability), with a 0077 rationale note. The 0076 direction-of-harm
      and earn-it requirements are preserved, not rewritten.
- [x] `specs/skills/quality-skill/guides/log.md` — record the guide-spec update.

### Format spec

- [x] None — no change to `SPECIFICATION.md`.

### Durable docs (bundled skill)

- [x] `skills/quality/guides/authoring.md` — reframe the "Cover the domain's
      constituent kinds" subsection per the spec (steward definition, claim-first
      framing, artifact-as-evidence, recurring lifecycle band, vulnerability-grounded
      protective pair, and the new presence-vs-quality "Do" bullet).
- [x] `skills/quality/workflows/setup.md` — **no change (assessed).** The
      composite-root stewardship walk concerns presence and accounting and already
      states "prompt, not a quota — traced to a Need or Risk" and carries
      thin/missing kinds as Areas with findings; the care refinements add no
      obligation it lacks.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` — **no change
      (assessed).** Check 8 detects constituent kinds that are neither modeled nor
      accounted for (presence); it is orthogonal to the claim-source and
      presence-vs-quality framing.

### Release

- [x] `CHANGELOG.md` — add the `/quality Skill` note under `Unreleased`.

## Children

- [Functional spec](0077-stewardship-care-grounding/spec.md) — required guidance
  content and acceptance criteria.
- [Design doc](0077-stewardship-care-grounding/design.md) — the five refinements,
  the care concepts, and the alternatives.

## Status

`Done`. Landed the care-grounding of the stewardship-concern generator: reframed
the "Cover the domain's constituent kinds" authoring subsection
(steward-as-holding-in-trust, claim-from-a-Need-or-Risk, artifact-as-evidence,
recurring lifecycle band, vulnerability-grounded protective pair, and the
constituent-vs-stewardship-quality "Do" bullet), aligned the durable authoring
guide spec and promoted the rationale, updated the guides log, and added the
CHANGELOG note. The setup workflow and the Top 10 checks were assessed as
no-change (their use of the concerns is about constituent presence and
accounting). Documentation-only (no `SPECIFICATION.md` or code change). Verified
with `mise run check` (markdown format, bundle link resolution, lint, and Go gates
all pass). Archived.
