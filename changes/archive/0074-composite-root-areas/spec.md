---
type: Functional Specification
title: Composite root areas and use-context constituents — functional spec
description: Required authoring-guidance content for the composite root shape, recurring use-context constituents, and the corrected root-factor-count heuristic.
tags: [skill, guide, authoring, areas, factors]
timestamp: 2026-06-24T00:00:00Z
---

# Composite root areas and use-context constituents — functional spec

Companion to the
[Composite root areas and use-context constituents](../0074-composite-root-areas.md)
change case. This spec states _what_ the guidance must say; the
[design doc](design.md) covers the taxonomy and _why_.

The key words **MUST**, **SHOULD**, and **MAY** are to be interpreted as
described in IETF RFC 2119.

## Scope

Covers the content of the two bundled guides
[`skills/quality/guides/authoring.md`](../../../skills/quality/guides/authoring.md)
and
[`skills/quality/guides/top-10-quality-md-checks.md`](../../../skills/quality/guides/top-10-quality-md-checks.md),
and the alignment of their governing guide specs. **Deferred:** any change to
`SPECIFICATION.md`, the format, the CLI, or evaluation/reporting behavior.

## Requirements

### Root shapes

- R1. The authoring guide **MUST** name three decomposition shapes —
  **primary-subject**, **collection**, and **composite** — and give the test that
  distinguishes them: whether a node's children are absent/refinements
  (primary-subject), the same kind (collection), or different kinds each with
  their own factor family (composite). It **MUST** present them as patterns that
  apply at the root and, recursively, at any area — not a one-time root
  classification.
- R2. The guide **MUST** define a **composite** root as one whole decomposed into
  unlike, interdependent constituents of different kinds, and **MUST**
  distinguish it from a **collection/grouping** area of homogeneous members. It
  **MUST NOT** use "portfolio" as the term for the composite case.
- R3. The guide **MUST** state that a near-disjoint factor family is a first-class
  signal to split a constituent into its own area.
- R4. The guide **MUST** identify the composite root's emergent concern as
  cross-part coherence (alignment, conformance, currentness _between_ parts), and
  **SHOULD** connect it to the assessment edges between areas already described in
  "Make the traceability graph visible".

### Use-context constituents

- R5. The guide **MUST** distinguish **domain constituents** (intrinsic to the
  modeled entity's domain, so they vary) from **use-context constituents**
  (present because of how QUALITY.md is used, so they recur across domains).
- R6. The guide **MUST** name the **agent harness** and the **QUALITY.md
  self-check** as the recurring use-context constituents, and **MUST** attribute
  their recurrence to QUALITY.md's agentic context of use rather than to any
  modeled quality domain, consistent with the use-context-vs-model-domain
  distinction in [`AGENTS.md`](../../../AGENTS.md).
- R7. The guide **MUST** state that the QUALITY.md self-check constituent stays on
  the learn-loop axis and is kept out of the entity's roll-up — the model's own
  quality is never averaged into the root area's rating — consistent with the
  two-loops principle in "When to update QUALITY.md".
- R8. The guide **MUST** note that the agent harness is partly normative (it
  governs agent behavior) and so plays the dual area/assessment-reference role,
  with a caution against double-counting its influence.
- R9. The guide **MUST** frame use-context constituents as expected defaults that
  must be consciously justified when omitted, **not** as required areas, and
  **MUST** keep the existing inclusion test (high-leverage, germane here, owned
  and inspectable, traced to a Need or Risk). It **MUST NOT** introduce a roster
  of areas every model must carry.

### Factor-count heuristic

- R10. The guide **MUST** scope the "aim for roughly ten root-level factors /
  fewer than eight triggers a coverage review" heuristic to a **primary-subject**
  root.
- R11. The guide **MUST** state that at a composite root the factor-count
  heuristic applies per constituent, and that the composite root itself carries
  only the factors that recur across constituents — typically stewardship lenses
  (currentness, traceability, consistency, maintainability) — each refined per
  child.
- R12. The "Keep the root lean when child areas carry the detail" subsection
  **SHOULD** tie "model-wide factors at the root" to this precise meaning: the
  factors that recur across constituents, not an arbitrary subset.

### Top 10 checks

- R13. Check 8 (Area and factor shape) **MUST** add a routing finding for a
  composite entity flattened into a primary-subject root — distinct constituent
  artifacts of different kinds described in the body, but all factors held at the
  root as one family — routed to authoring.
- R14. Check 8 **MUST** add a routing finding for a missing expected use-context
  constituent (an agent-collaborated project whose body shows an owned,
  high-leverage harness or self-check that is not modeled), routed to authoring,
  phrased as an expected default rather than a requirement.
- R15. The Top 10 additions **MUST NOT** read as a roster mandate; the guide
  **MUST** keep presence earned (germane, owned, traced to a Need or Risk), and
  **MUST NOT** flag a harness-less or throwaway project for the absence.

### Composition and recursion

- R17. The guide **MUST** present the three shapes as recursive and composable,
  not mutually exclusive root choices: a composite area **MAY** contain
  collection children, a collection member **MAY** itself be composite or
  primary-subject, and a constituent **MAY** decompose further to any depth. It
  **MUST** state that the factor-count heuristic (R10–R11) applies at each
  primary-subject node wherever it sits in the tree, not only at the root.
- R18. The guide **MUST** include a worked example of a composite root whose
  children mix kinds — at minimum a use-context constituent (harness or
  QUALITY.md self-check) alongside a **collection** child that holds homogeneous
  members (e.g. an `apps` collection containing `apps/product1` and
  `apps/product2`). The example **MUST** remain illustrative and domain-scoped per
  R16, not implied as a default model layout.

### Consistency

- R16. All added content **MUST** preserve domain agnosticism per
  [`AGENTS.md`](../../../AGENTS.md): examples stay illustrative and
  non-exhaustive, agentic language is scoped to context of use, and software
  product quality is not implied as the default modeled domain.

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` — confirm and, where the
  contract constrains the affected subsections, update it to admit the
  root-shape and use-context-constituent material.
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — reflect the
  new Check 8 routing findings if the contract enumerates per-check findings.

### To rename

None.

### To delete

None.
