---
type: Design Doc
title: Composite root areas and use-context constituents — design doc
description: The three-root-shape taxonomy, the domain-vs-use-context constituent split, and how it lands in the authoring and Top 10 guides.
tags: [skill, guide, authoring, areas, factors]
timestamp: 2026-06-24T00:00:00Z
---

# Composite root areas and use-context constituents — design doc

Design behind the
[Composite root areas and use-context constituents](../0074-composite-root-areas.md)
change case and its [functional spec](spec.md).

## Context

The authoring guide implicitly treats the root area as a single primary subject
with one factor family — the "aim for ~ten root-level factors" heuristic assumes
it. For most real entities in QUALITY.md's assumed context of use the root is a
**composite**: one whole decomposed into distinct _kinds_ of constituent
artifacts, each with its own largely-disjoint factor family. The guide's existing
split test already implies this but never names the shape, so authors default to
a flat, product-factored root that either drops other high-leverage artifacts or
flattens incompatible factor families together. See the change case
[motivation](../0074-composite-root-areas.md#motivation) for the full problem
statement; the spec states the required content.

## Approach

### Three decomposition shapes (recursive, composable)

The shapes describe how _any_ area decomposes, not a one-time root
classification. They apply at the root and recurse at every node:

| Shape               | A node's children are                       | The node's own concern                                                       | Factor-count heuristic      |
| ------------------- | ------------------------------------------- | ---------------------------------------------------------------------------- | --------------------------- |
| **Primary-subject** | none, or refinements; one factor family     | the subject's own factors                                                    | applies here                |
| **Collection**      | many of the _same_ kind                     | set-level concerns no member has (coverage, balance, non-redundancy)         | applies per member family   |
| **Composite**       | many _different_ kinds, each its own family | cross-part **coherence** (alignment, conformance, currentness between parts) | applies **per constituent** |

The collection case is already covered by the guide's whole-set-concern material
("a test suite is a different area"); the composite case is the gap. The two are
genuinely different: a collection holds _like_ things judged by set-level
concerns; a composite holds _unlike_, interdependent parts whose emergent concern
is coherence _between_ them — and that coherence is exactly the assessment-edge
graph the guide already prizes in "Make the traceability graph visible". A real
project is almost always composite, which is why a product-factored root felt
wrong.

**The shapes compose; they are not mutually exclusive.** A composite area
routinely holds children of different shapes — some primary-subject
constituents, some that are themselves collections or composites — nested to any
depth. A composite root with a use-context harness and self-check alongside an
`apps` collection is the common case:

```
root  (composite)
├── harness        (primary-subject constituent)
├── quality-md     (primary-subject constituent; learn-loop axis, out of roll-up)
└── apps           (collection)
    ├── apps/product1   (primary-subject — or itself composite)
    └── apps/product2   (primary-subject — or itself composite)
```

Because the shapes recurse, the factor-count heuristic lands at each
primary-subject node wherever it sits — `apps/product1` earns its own ~ten-factor
family, the `apps` collection carries set-level concerns, and the composite root
carries only the stewardship factors that recur across its constituents.

### Two sources of composite constituents

- **Domain constituents** — intrinsic to the modeled entity's domain, so they
  _vary_: product + requirements + docs (software); data + schema +
  collection-methodology (dataset); document + terminology-standard + sources
  (reference doc).
- **Use-context constituents** — present because of _how QUALITY.md is used_, so
  they _recur across domains_: the **agent harness** (the instructions steering
  the agent on the project) and the **QUALITY.md self-check** (the model's own
  quality).

This is the `AGENTS.md` use-context-vs-model-domain distinction applied to area
structure: we keep agnosticism in _what the constituents are about_ while letting
the agentic context of use determine _which constituents recur regardless_. A
dataset project and a document project carry different domain constituents but the
same two use-context constituents.

### Two roll-up asymmetries

The use-context constituents are not symmetric with domain children, and the
guide must say so or authors will model them wrong:

1. The **QUALITY.md self-check** stays on the **learn-loop axis, out of the
   entity's roll-up** — the two loops never fold together; the model's own quality
   is never averaged into the root area's rating. It is tracked and reported
   separately. (This promotes the existing "the self-check area already anchors the
   model in this guide" line from a footnote to a named, recurring constituent.)
2. The **agent harness** is partly **normative** (it governs agent behavior), so
   it plays the dual area/assessment-reference role; watch for double-counting if
   its influence is also assessed inside a domain constituent.

### Root-level factors, made precise

At a composite root, declare only factors that genuinely _recur across every
constituent_ — typically the stewardship/internal lenses (currentness,
traceability, consistency, maintainability), each **refined** per child.
Artifact-specific experiential factors (reliability/security for a runtime
product; credibility for docs) stay local to their constituent. This is the
precise meaning "model-wide factors at the root" should carry, and it is why the
~ten-factor heuristic moves to _per primary subject / per constituent_.

### Where it lands

- `authoring.md`, "Area" section: new subsection "Choose the root shape" after
  "Choose areas that are authored, inspectable artifacts" (the three-shape table,
  composite-vs-collection contrast, near-disjoint family as a split signal,
  cross-part coherence as the emergent concern).
- `authoring.md`: new short subsection "Recurring use-context constituents"
  adjacent to "Ground high-leverage concerns in normative artifacts" (harness +
  self-check, the two asymmetries, the expected-default guardrail).
- `authoring.md`, "Cover the domain's stable stakes before specializing": caveat
  the factor-count heuristic to primary-subject / per-constituent.
- `authoring.md`, "Keep the root lean when child areas carry the detail": one line
  tying model-wide factors to "recur across constituents".
- `top-10-quality-md-checks.md`, Check 8: composite-flattening finding +
  missing-expected-constituent finding + earned-not-quota clarifier.

## Alternatives

- **Make harness + self-check mandatory areas.** Rejected — reproduces the
  "universal roster of areas every model must carry" anti-pattern the guide
  explicitly warns against; some projects legitimately lack them.
- **Treat the composite case as the existing collection/grouping area.** Rejected
  — collapses the homogeneous (set-level concerns) and heterogeneous (cross-part
  coherence) cases, which have different root concerns and different factor
  semantics.
- **Roll the QUALITY.md self-check into the root rating as an ordinary child.**
  Rejected — violates the two-loops invariant.
- **Promote "composite/collection/primary-subject" to `SPECIFICATION.md` terms.**
  Deferred — they are authoring heuristics, not format semantics; keeping them
  guide-only avoids over-constraining the format.
- **Rewrite the Factor section's stable-stakes guidance.** Rejected — a caveat is
  enough; a rewrite risks destabilizing guidance that is correct for the
  primary-subject case.

## Trade-offs & risks

- **Added vocabulary.** Three root shapes plus two constituent sources is more to
  learn. Mitigated by additive subsections and a single table; a primary-subject
  author can ignore it.
- **Over-decomposition.** Authors may split trivially. Mitigation: the split is
  _earned by each constituent having a real, separately-trusted factor family_ —
  stated explicitly; a narrow product with thin docs need not decompose.
- **Roll-up confusion** around the self-check's out-of-band axis. Mitigation:
  state it inline and cross-link the two-loops section.

## Open questions

- Does the harness deserve worked guidance on its dual area/assessment-reference
  role, or is a one-line caveat enough?
- Is there a third recurring use-context constituent worth naming (e.g.
  evaluation/recommendation artifacts), or do those stay out of the entity tree by
  design?
- Are "primary-subject" and "composite" truly distinct or two ends of a
  continuum? The guide should frame the split as earned by separately-trusted
  factor families, not taken by default.
