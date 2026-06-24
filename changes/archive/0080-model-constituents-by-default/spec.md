---
type: Functional Specification
title: Model constituents by default — functional spec
description: The model-by-default constituent-coverage contract for the /quality skill — model every germane constituent, skip only on two disqualifiers, never silently omit, and bar an under-covered model from evaluation-ready.
tags: [skill, authoring, setup, constituents, coverage]
timestamp: 2026-06-24T00:00:00Z
---

# Model constituents by default — functional spec

Companion to the [Model constituents by default](../0080-model-constituents-by-default.md)
change case. This spec states *what* the skill's constituent-coverage guidance
must say; the [design doc](design.md) covers *how* the reframe is worded and why.
It governs the bundled skill ([`skills/quality/`](../../../skills/quality/)) and
its functional-spec mirror
([`specs/skills/quality-skill/`](../../../specs/skills/quality-skill/)), and defers
the QUALITY.md format itself to
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

The key words "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", and "MAY" are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The skill made an area something a constituent must **earn**: a three-part test
(owned, inspectable artifact; divergent factor family; traced to a Need or Risk)
inside "a prompt, not a quota" / "not a roster" language. That guidance guards
only against *over*-modeling; nothing counters *under*-modeling, and "defer it in
Scope" sits as a cost-free peer of "model it." On a real multi-service monorepo,
setup produced a flat root with five model-wide factors and deferred every
per-constituent area to "the next iteration," then passed the maturity bar because
a deferral note counts as "accounted for."

A first-pass model should be as full and sufficient as the evidence supports.
Setup runs roughly once per project and the user is most engaged then. The fix is
to make modeling the default and deferral the exception, and to make a germane
concern impossible to drop silently: an absent or thin artifact for a real concern
is a high-value early signal, so it MUST surface as a ratable element of the model,
not as prose.

This case touches only the inclusion default and its enforcement. The
stewardship-concern generator, the audience×purpose axis, the three-projections
rule, and the motivation/taxonomy vocabulary discipline are unchanged.

## Scope

Covered: the rule that decides which constituents of a composite root become their
own areas, across the authoring guide, the setup workflow, and the Top 10 checks
(and their spec mirrors) — the model-by-default default, the disqualifier set, the
no-silent-omission rule, the routing between a minimal area and a requirement
elsewhere, the first-pass completeness bar, and the demotion of deferral.

Deferred / non-goals: no change to the QUALITY.md format or schema, to
`SPECIFICATION.md`, or to the CLI; the stewardship/audience generators and
vocabulary discipline are unchanged; re-checking this repo's own `QUALITY.md`
against the new rule is a follow-up.

## Requirements

### Model-by-default

The authoring guide, the setup workflow, and their spec mirrors **MUST** state
that, for a composite root, the author enumerates the constituent kinds the domain
implies and **models each as its own area by default**. Giving a constituent its
own area is the default outcome, **MUST NOT** be framed as something the
constituent must *earn*, and the thinness of a first pass **MUST NOT** be given as
a reason to defer or omit a constituent.

> Rationale: setup runs once per project with the user most engaged; an under-built
> first model teaches that a flat root is normal and coverage is optional. The old
> "earn an area" framing supplied only anti-over-modeling pressure. — 0080

The guidance **MUST** still teach the author to scale coverage to *this* entity:
the enumerated kinds remain a prompt for what the entity asks to be cared for, not
a universal roster, and a constituent is modeled because it is germane here, not
because the generator names it. The model-by-default default and the
"prompt, not a roster" caution **MUST** be expressed as one coherent rule, not as
opposing pressures — coverage is earned by the entity's own Needs and Risks, and a
germane constituent is modeled rather than deferred.

### The disqualifiers

The guidance **MUST** state that a germane constituent is given its own area
unless one of exactly two disqualifiers holds:

1. **No distinct concerns** — its quality is already fully judged by its parent's
   or a sibling's factors, so a separate area would assess nothing new; fold the
   concern into that area.
2. **Not germane / outside the boundary** — the domain does not imply it for this
   entity, or it belongs to another system or owner; mark it out of Scope.

The guidance **MUST NOT** list "no artifact exists yet" (or equivalent absence of
the constituent's artifact) as a reason to omit a germane constituent.

> Rationale: the two disqualifiers are the old earn-test's "divergent factor
> family" and boundary checks, inverted into skip conditions. "Traced to a Need or
> Risk" is dropped as a gate — it drives factor selection, not area inclusion — so
> the test stays short. — 0080

### No silent omission

The guidance **MUST** state that a germane concern is **never** omitted by being
recorded only in prose (a Scope or "deferred" note). When a germane concern's
artifact is absent or thin, its absence **MUST** be surfaced as a *ratable element
of the model* — one that produces a rating — by one of:

- **(a)** modeling it as a minimal area carrying a missing-anchor finding, or
- **(b)** a requirement on its parent or a sibling area that rates poorly because
  the artifact is missing.

The guidance **MUST** make clear that a Scope or deferral note alone does **not**
satisfy this rule; only a ratable area or requirement does.

> Rationale: an absent or thin artifact for a real concern is a high-value early
> signal — the worst case the model exists to catch. Dropping it to prose hides it
> behind a clean Scope section. — 0080

### Routing between a minimal area and a requirement elsewhere

The guidance **MUST** give a criterion for choosing between (a) and (b), tied to
the same distinct-concerns test used everywhere else:

- model the absence as its **own minimal area (a)** when the kind would carry its
  own factor family once it exists — a high-leverage, first-class kind, or a
  whole-constituent gap — so the gap gets its own axis and rating; and
- surface it as a **requirement on an existing area (b)** when the concern belongs
  to an existing area's factors, the gap is partial or a matter of degree, or a
  standalone area would be a single-finding stub.

### Deferral is a narrow exception

The guidance **MUST** present in-scope deferral ("in scope, not yet modeled") as a
narrow exception, not a routine accounting option peer to modeling. Deferral
**MUST** be reserved for a constituent that is genuinely blocked (for example, on
an undecided boundary), and the specific blocker **MUST** be recorded. "Next
iteration" and "the first model is thin" **MUST NOT** be accepted as deferral
reasons.

### First-pass completeness bar

The setup workflow (and its spec mirror) **MUST** treat a composite model that
leaves a non-disqualified germane constituent unmodeled — or recorded only as a
deferral — as an important gap that bars `evaluation-ready` maturity. A bare
deferral note **MUST NOT** be sufficient to reach `evaluation-ready`.

The Top 10 checks guide (and its spec mirror) **MUST** flag a constituent the
domain implies that is neither modeled as an area nor disqualified by one of the
two disqualifiers, and **MUST** state that a bare deferral or Scope note does not
satisfy the check. The check **MUST** still avoid flagging a constituent that a
throwaway or narrowly scoped entity would not carry (i.e. a constituent that hits
the not-germane disqualifier).

> Rationale: the old check treated constituent presence as an "earned default" and
> explicitly did *not* flag a missing one, so a deferral note passed maturity
> silently. The bar is what makes model-by-default enforceable rather than
> advisory. — 0080

### Use-context constituents and first-model wording

The "recurring use-context constituents" guidance (agent harness, QUALITY.md
self-check) **MUST** be consistent with model-by-default: each is modeled unless it
hits a disqualifier (a too-thin harness with no owned artifact hits
not-germane / nothing-to-evaluate and is recorded as a gap, not silently dropped).
The QUALITY.md self-check **MUST** stay on the learn loop and out of the entity's
roll-up, unchanged.

The getting-started first-model guidance **MUST NOT** advise a "small first model"
in a way that reads as endorsing under-coverage; it **MUST** instead guard
*assessability* (prefer assessable requirements over aspirational ones) while
leaving constituent coverage governed by model-by-default.

## Durable spec changes

Durable **specs** this case rewrites — the [`specs/`](../../../specs/index.md)
bundle and [`SPECIFICATION.md`](../../../SPECIFICATION.md). See
[Writing functional specs](../../../docs/guides/write-functional-specs.md#durable-spec-changes).

### To add

None.

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` — restate the
  constituent-coverage contract as model-by-default with the two disqualifiers and
  the no-silent-omission / routing rules; replace the "earn an area" and "roster"
  MUSTs (per Model-by-default, The disqualifiers, No silent omission, Routing, and
  Use-context constituents above).
- `specs/skills/quality-skill/workflows/setup.md` — require the model-by-default
  constituent walk and the first-pass completeness bar in the maturity
  classification and completion criteria, and demote deferral (per First-pass
  completeness bar and Deferral is a narrow exception above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — flip the
  constituent-coverage check to flag a non-disqualified constituent left unmodeled
  or merely deferred (per First-pass completeness bar above).
- `specs/skills/quality-skill/guides/log.md` — record the contract revision.

### To rename

None.

### To delete

None.
