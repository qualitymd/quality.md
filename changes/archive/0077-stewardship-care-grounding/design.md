---
type: Design Doc
title: Care-grounded stewardship concerns — design doc
description: The five refinements the phenomenology of care makes to 0076's stewardship-concern generator, the care concept behind each, where each lands in the authoring guide, and the alternatives weighed.
tags: [skill, guide, authoring, areas, constituents, stewardship]
timestamp: 2026-06-24T00:00:00Z
---

# Care-grounded stewardship concerns — design doc

Design behind the
[Care-grounded stewardship concerns](../0077-stewardship-care-grounding.md)
change case and its [functional spec](spec.md).

## Context

[0076](../0076-domain-constituent-kinds.md) introduced a domain-agnostic
generator for a composite root's constituent kinds, with stewardship concerns as
one of its two axes. Its framing sentence — a concern "tends to leave an
authored, inspectable artifact that is a candidate constituent" — is
artifact-first: it reads as _enumerate the residue_, and it lets a reader treat a
folder's presence as a concern met. That is serviceable but slightly off from
what the generator is for.

Two essays by B. Scot Rousse and David Spivak —
[_Notes on care_](https://withoutwhy.substack.com/p/notes-on-care-with-david-spivak-a)
and [_Clarifications of care_](https://withoutwhy.substack.com/p/clarifications-of-care)
— supply a sharper footing. Their load-bearing claims, for our purposes:

- Care is something one _does and is_, "tending and attending to what matters",
  not a possession or a list of values. Its artifacts are its trace.
- "Doing what is required, _well_" — the _well_ is attunement, not mere presence
  or efficiency.
- The _required_ is what the situation addresses to the carer: "to experience
  something as required is to find oneself addressed; the situation makes a
  claim."
- Caregiving is the intensified, _asymmetric_ mode of care — responsibility
  toward a vulnerable other who cannot fully protect itself.
- Care is open-ended: "always more to do."

A caveat that shapes scope: despite Spivak's billing, the published essays are
phenomenological, not categorical — they contain no functors, lenses, or
adjunctions. So this change takes the _framing_, not a formal apparatus (see
[Alternatives](#alternatives)).

## Approach

Five refinements, each re-grounding one phrase the generator already uses. None
adds a concern, a factor family, or an area; none touches the format. All land in
the single authoring-guide subsection "Cover the domain's constituent kinds" and
its durable spec.

| #   | Care concept                                     | Refinement                                                                                               | Where it lands                           |
| --- | ------------------------------------------------ | -------------------------------------------------------------------------------------------------------- | ---------------------------------------- |
| 1   | Care is tending; artifacts are its trace         | Present each concern as an activity whose artifact is _evidence_, not the care itself; presence ≠ proof  | Subsection opening (R2)                  |
| 2   | The required is what the situation _claims_      | The constituent kinds are a prompt; the claim that earns an area comes from a Need or Risk, not the list | Subsection opening + earn-it bullet (R3) |
| 3   | Caregiving = asymmetric care under vulnerability | Ground `secure`/`safeguard` in vulnerability; the band cross-cuts because it tracks who is exposed       | Protective-band bullet (R5)              |
| 4   | Steward = caring for what is held in trust       | One-sentence definition of _stewardship_                                                                 | Subsection opening (R1)                  |
| 5   | Care is open-ended                               | Lifecycle band is recurring, not a once-through pipeline; `operate`/`maintain` never conclude            | Lifecycle-band bullet (R6)               |

Plus one structural consequence of #1 and the _well_: a new "Do" bullet
separating the **constituent question** (does a tending leave an owned artifact?
— an _area_) from the **stewardship-quality question** (is it done _well_? — a
_factor_). This is where the care framing pays off operationally: it routes
"present but poor" to the factor layer and keeps "absent but high-leverage" as an
empty area with a missing-anchor finding — both already latent in 0076, now
named (R4).

### Why these are framing, not mechanism

The generator's machinery is unchanged. The nine concerns, the lifecycle ×
protective banding, the audience × purpose axis, the three-projections rule, and
the earn-it guardrail all stand verbatim. What changes is the _reading_: from
"collect the artifacts a domain leaves" to "hear what this entity asks to be
cared for, and model where that asking has — or should have — left an inspectable
trace." That reading makes the existing guardrails (prompt-not-quota,
traced-to-a-Need-or-Risk, carry-thin-kinds-as-findings) follow from a principle
instead of standing as bare cautions.

### The three-projections rule, and the formalism we leave out

0076's "a concern projects as factor, constituent, and audience" is the most
lens-like structure in the guide — one root concern, three coherent views. It is
tempting to formalize it with Spivak's categorical vocabulary. We deliberately do
not: the source essays do not develop it, and importing functor/lens language
would dress a phenomenological point in a formalism it cannot cash. The analogy
is noted here for a future editor; it stays out of the guide.

## Alternatives

- **Promote the care vocabulary into `SPECIFICATION.md`.** Rejected — like
  0076's taxonomy, this is an authoring heuristic, not format semantics. Keeping
  it guide-only avoids over-constraining the format.
- **Adopt formal category-theory framing (functors, lenses, adjunctions).**
  Rejected — not present in the source essays, which are phenomenological;
  formalism here would over-claim and mislead. The lens analogy for the
  three-projections rule is recorded in this doc only.
- **Add a new `care` concern, or a "stewardship-quality" factor family.**
  Rejected — care is not a tenth concern; it is what the whole generator is
  _about_. The presence-vs-quality distinction is already carried by the existing
  area/factor split; R4 names it rather than inventing a new construct.
- **Replace "stewardship" with "care" throughout.** Rejected — the two do
  distinct work. _Care_ is the activity (tending); _steward_ is the role (holding
  in trust, answerable to others and the future). Keeping both lets R1 define the
  role and R2 describe the activity.
- **Drop the secure/safeguard direction-of-harm split in favour of a single
  "vulnerability" framing.** Rejected — direction of harm is the operational
  distinction (who is protected from whom) and stays load-bearing, especially for
  agentic projects. Vulnerability grounds _why_ the band exists; it does not
  replace the split. R5 keeps both.
- **Redraw the lifecycle as an explicit cycle diagram.** Rejected as overkill —
  one clause ("recurring, not a once-through pipeline; `operate` and `maintain`
  never conclude") carries the open-endedness without a diagram, and keeps the
  arrows that still convey a typical first pass.
- **Treat this as a routine guide edit, no change case.** Reasonable under the
  routine-changes rule, but the refinements are subtle: a later editor could
  "simplify" artifact-as-evidence back into artifact-first, or read the care
  framing as relaxing the earn-it test. A change case promotes the rationale onto
  the durable spec so the distinction survives. (The user also asked for the full
  lifecycle.)

## Trade-offs & risks

- **Density.** The subsection is already dense; this adds philosophical framing.
  Mitigation: the refinements are woven into existing sentences plus one bullet —
  no new section — and a primary-subject author skips the whole subsection.
- **Softness / jargon.** "Care", "vulnerability", "held in trust" risk reading as
  soft. Mitigation: each cashes out into an operational distinction (area vs.
  factor; direction of harm; Need/Risk as the claim source). No phenomenology
  jargon enters the guide.
- **Misread as relaxing the guardrail.** Framing concerns as "tending" could be
  taken to make artifacts optional. Mitigation: R3 and R4 explicitly preserve the
  earn-it test and the missing-anchor finding; the claim still must trace to a
  Need or Risk.
- **Source fidelity.** Over-claiming a Spivak formalism that the essays do not
  contain. Mitigation: no categorical apparatus in the guide; the analogy is
  fenced to this design doc.

## Open questions

- 0076 asked whether `discover` is reliably an inspectable artifact across
  domains. The artifact-as-evidence reframing answers it: `discover` is a tending
  that often leaves only a thin trace, legitimately carried as a finding rather
  than a populated area — the presence-vs-quality split (R4) covers the case
  without special-casing the concern.
- 0076 asked whether the protective pair should graduate into named factors. The
  vulnerability grounding plus the three-projections rule reinforce keeping it
  cross-cutting and constituent-side by default; left unchanged.
- Should the _well_ / attunement reading earn explicit treatment in the
  stable-stakes factor guidance (not just at the area/factor boundary)? Deferred;
  R4 names it where the boundary is decided, and the factor guidance is out of
  scope here.
