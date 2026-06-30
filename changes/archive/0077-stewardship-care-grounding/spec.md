---
type: Functional Specification
title: Care-grounded stewardship concerns — functional spec
description: Required refinements to the authoring guide's stewardship-concern generator — steward as holding-in-trust, the claim from a Need or Risk, the artifact as evidence, the constituent vs. stewardship-quality split, the protective pair under vulnerability, and a recurring lifecycle band.
tags: [skill, guide, authoring, areas, constituents, stewardship]
timestamp: 2026-06-24T00:00:00Z
---

# Care-grounded stewardship concerns — functional spec

Companion to the
[Care-grounded stewardship concerns](../0077-stewardship-care-grounding.md)
change case. This spec states _what_ the guidance must say; the
[design doc](design.md) covers the care concepts behind each refinement and _why_.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 (RFC 2119, RFC 8174) when, and only when, they
appear in all capitals.

## Scope

Covers the "Cover the domain's constituent kinds" subsection of the bundled guide
[`skills/quality/guides/authoring.md`](../../../skills/quality/guides/authoring.md)
and the alignment of its governing durable spec
[`specs/skills/quality-skill/guides/authoring-md.md`](../../../specs/skills/quality-skill/guides/authoring-md.md).

**Deferred / non-goals:** any change to `SPECIFICATION.md`, the format, the CLI,
or evaluation/reporting behavior; any new stewardship concern; any relaxation of
the [0076](../0076-domain-constituent-kinds.md) earn-it guardrail; any change to
the setup workflow or the Top 10 checks (their use of the concerns is about
constituent presence and accounting, which these refinements do not alter); and
any formal category-theory apparatus in the guide.

## Background / Motivation

[0076](../0076-domain-constituent-kinds.md) shipped a domain-agnostic generator
for constituent kinds, framed as "caring for any entity carries a recurring set
of concerns, and each tends to leave an authored, inspectable artifact that is a
candidate constituent." Read literally that framing is artifact-first: it invites
collecting residue and reading a folder's presence as a concern met. The
phenomenology of care in the two source essays (see the
[change case](../0077-stewardship-care-grounding.md)) names care as an activity —
"tending and attending to what matters", "doing what is required, well" — whose
artifacts are its _trace_, and names the _required_ as what the situation, by
making a claim, addresses to the carer. This spec re-grounds the generator on
that footing. It changes framing, not the generator's shape: the nine concerns,
the two axes, the three-projections rule, and the earn-it guardrail all stand.

## Requirements

### Stewardship as care

- R1. Where the guide introduces the stewardship-concern generator, it **MUST**
  define _stewardship_ as caring for an entity held in trust — answerable to its
  stakeholders and its future — rather than as private ownership.
- R2. The guide **MUST** present each stewardship concern as an ongoing activity
  of _tending_, whose authored, inspectable artifact is the concern's **trace**
  and the candidate constituent. It **MUST** state that the artifact is read as
  _evidence_ of tending, not as the tending itself: a thin or missing artifact
  does not by itself prove the concern unmet, and a populated artifact does not
  prove it met well.

### The claim that earns a constituent

- R3. The guide **MUST** frame the constituent kinds as a prompt for what _this_
  entity asks to be cared for, and **MUST** state that the claim which earns a
  constituent an area originates in a **Need or Risk** the entity presents, not
  in the generator list itself. This **MUST** be consistent with, and **MUST
  NOT** weaken, the 0076 earn-it guardrail (owned, inspectable artifact;
  divergent factor family; traced to a Need or Risk; prompt, not quota).

### Constituent presence vs. stewardship quality

- R4. The guide **MUST** distinguish the **constituent question** — whether a
  tending leaves an owned, inspectable artifact, modeled as an _area_ — from the
  **stewardship-quality question** — whether it is done _well_, attuned to the
  entity's situation rather than merely present or complete, carried by
  _factors_. It **MUST** preserve the existing rule that a germane, high-leverage
  kind is carried as an area even when its artifact is thin or missing, recording
  the gap as a finding within it.

### The protective band under vulnerability

- R5. The guide **MUST** present the protective pair (`secure`, `safeguard`) as
  the mode of stewardship that answers to **vulnerability** — a party that can be
  harmed and cannot fully protect itself — and **MUST** state that the band
  cross-cuts the lifecycle because it tracks _who is exposed_ rather than which
  phase the work is in. It **MUST** preserve the 0076 direction-of-harm
  distinction (`secure` guards the entity from the world; `safeguard` guards
  stakeholders and the environment from the entity) and **MUST NOT** collapse the
  pair into a lifecycle stage.

### The lifecycle band is recurring

- R6. The guide **MUST** present the **lifecycle** band as recurring rather than a
  once-through pipeline: the order is a typical first pass, the concerns are
  standing — revisited for as long as the entity is tended — and `operate` and
  `maintain` do not conclude. It **MUST NOT** imply the lifecycle completes.

### Consistency

- R7. All refined content **MUST** preserve domain agnosticism per
  [`AGENTS.md`](../../../AGENTS.md) (illustrative, non-exhaustive examples;
  agentic language scoped to context of use; software product quality not implied
  as the default modeled domain) and **MUST NOT** introduce formal
  category-theory apparatus; the source framing is phenomenological. The
  refinements **MUST** leave the nine concerns, the two axes, and the
  three-projections rule intact.

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` — extend the
  domain-constituent-kinds contract with the care-grounding clauses: the claim
  originating in a Need or Risk rather than the generator list (per R3), the
  constituent vs. stewardship-quality split (per R4), and the protective pair as
  cross-cutting stewardship under vulnerability (per R5). The 0076
  direction-of-harm and earn-it requirements are preserved, not rewritten.

### To rename

None.

### To delete

None.
