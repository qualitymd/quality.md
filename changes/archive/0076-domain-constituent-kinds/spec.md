---
type: Functional Specification
title: Domain constituent kinds and stewardship concerns — functional spec
description: Required authoring, setup, and Top 10 guidance for enumerating a composite root's constituents by stewardship concerns and an audience×purpose axis.
tags: [skill, guide, authoring, areas, constituents]
timestamp: 2026-06-24T00:00:00Z
---

# Domain constituent kinds and stewardship concerns — functional spec

Companion to the
[Domain constituent kinds and stewardship concerns](../0076-domain-constituent-kinds.md)
change case. This spec states *what* the guidance must say; the
[design doc](design.md) covers the taxonomy and *why*.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 (RFC 2119, RFC 8174) when, and only when, they
appear in all capitals.

## Scope

Covers the content of the bundled guide
[`skills/quality/guides/authoring.md`](../../../skills/quality/guides/authoring.md),
the bundled guide
[`skills/quality/guides/top-10-quality-md-checks.md`](../../../skills/quality/guides/top-10-quality-md-checks.md),
the setup workflow
[`skills/quality/workflows/setup.md`](../../../skills/quality/workflows/setup.md),
and the alignment of the two governing guide specs. **Deferred:** any change to
`SPECIFICATION.md`, the format, the CLI, or evaluation/reporting behavior; and
any heavy treatment in `getting-started.md` beyond an optional one-line pointer.

## Background / Motivation

[0074](../0074-composite-root-areas.md) named the composite root shape
and the recurring *use-context* constituents but gave *domain* constituents no
equally prescriptive treatment, so a setup-authored model enumerates
constituents by walking the repository's folders and never asks what *kinds* of
constituent the domain implies. Thin, scattered, or missing constituents (tests,
specs, documentation, a threat model) leave no folder to trip over and so vanish
instead of being recorded as findings. This spec adds a domain-agnostic generator
for constituent kinds so that absence becomes visible.

## Requirements

### Enumerating constituent kinds

- R1. The authoring guide **MUST** teach that, once the root is composite, the
  author enumerates its constituents by **constituent kind** inferred from the
  entity's quality domain, and **MUST** state that enumerating by repository
  folder structure alone is insufficient because it omits constituents that are
  thin, scattered, or absent.
- R2. The guide **MUST** present two generators of constituent kinds: a
  **stewardship-concern** axis (R3–R4) and an **audience × purpose** axis (R5).
- R3. The guide **MUST** name a **lifecycle** band of stewardship concerns —
  discover, define, realize, verify, enable, operate, and maintain — and
  **MUST** describe each as a concern that *tends to leave an authored,
  inspectable artifact* that is a candidate constituent, naming the concern (its
  function) rather than a domain-specific artifact (so "enable its audiences",
  not "docs").
- R4. The guide **MUST** name a cross-cutting **protective** pair distinct from
  the lifecycle band: **secure** — guard the entity from harm by the world — and
  **safeguard** — guard stakeholders and the environment (internal and external)
  from harm by the entity. It **MUST** state that the two are orthogonal (an
  entity can be hardened against attackers yet still harm its own users) and
  **MUST NOT** present `safeguard` as merely a synonym for `secure`; the
  direction of harm distinguishes them.

### Audience × purpose

- R5. The guide **MUST** teach that a single stewardship concern can yield
  several constituents when they serve different **audiences** or **purposes**,
  and **MUST** cite [Diátaxis](https://diataxis.fr/) once as this lens applied to
  the *enable* concern (tutorial, how-to, reference, and explanation as
  constituents with different factor families, not one "documentation" area).
- R6. The guide **MUST** teach that the audience side is derivable from the
  body's Needs (which already name the stakeholders), so that each audience a
  Need names should have an enabling — and verifying — constituent that is
  modeled or consciously accounted for.

### The three-projections rule

- R7. The guide **MUST** teach that a stewardship concern projects into the model
  in up to three ways — as a **factor** (the quality lens), a **constituent**
  (the artifact that pursues it), and an **audience** (who it serves or
  protects) — and that shared names (e.g. `secure`) reflect a shared root
  concern, not duplication.
- R8. The guide **MUST** instruct the author to identify which projection is
  being modeled and model it once, so the same concern is not double-counted —
  the security *of* an area is a factor on that area, while a security policy is
  its own constituent area.

### Accounting and guardrails

- R9. The guide **MUST** instruct the author to account for each constituent kind
  the domain implies by one of: modeling it as an area, deferring it in Scope,
  marking it out of Scope, or recording it as an unknown in the relevant section
  — and **MUST** state that silence is a coverage gap, not neutrality.
- R10. The guide **MUST** teach that a germane, high-leverage constituent kind is
  carried as an area even when its artifact is thin or missing, with the gap
  recorded as a finding within it, and **SHOULD** connect this to "Ground
  high-leverage concerns in normative artifacts" as that rule applied to whole
  constituents.
- R11. The guide **MUST** keep the inclusion test earned, not a roster: a
  constituent kind earns an area only when it leaves an owned, inspectable
  artifact, its factor family diverges from its siblings, and it traces to a Need
  or Risk. It **MUST NOT** present the constituent kinds as a roster every model
  must carry, and **MUST** state that a throwaway or narrowly scoped entity earns
  few or none.
- R12. The guide **SHOULD** connect the audience×purpose split to the existing
  decomposition shapes: a constituent that fans out by audience or purpose is
  itself a composite or collection node.
- R13. The guide **SHOULD** include one illustrative table mapping each
  stewardship concern to instances and conventional factor families, scoped to a
  single named domain (software product quality), and **MUST** mark it as a
  prompt for the entity's own domain rather than a checklist.

### Setup workflow

- R14. The setup workflow **MUST** instruct the agent, when the root is
  composite, to enumerate the domain's constituent kinds — walking the
  stewardship concerns and the audiences the Needs name — rather than only the
  components the repository already has folders for, and to account for each kind
  (model, defer in Scope, mark out of Scope, or record as an unknown).
- R15. The setup workflow **MUST** instruct the agent to carry a germane,
  high-leverage constituent kind (e.g. tests, specs, a threat model) as an Area
  even when its artifact is thin or missing, recording the gap as a finding
  rather than omitting the Area, and **MUST** preserve the earn-it guardrail
  (prompt, not quota). It **SHOULD** reference the authoring guide's "Cover the
  domain's constituent kinds" subsection.

### Top 10 checks

- R16. Check 8 (Area and factor shape) **MUST** add a routing finding for a
  domain that implies constituent kinds the body evidences or implies — for
  software product quality, tests, documentation modes, specs/requirements,
  operations, or a security/safety artifact — that the model neither models as
  areas nor accounts for (deferred, out of scope, or unknown), routed to
  authoring.
- R17. The Check 8 addition **MUST NOT** read as a roster mandate: it **MUST**
  keep presence earned and **MUST NOT** flag a kind a throwaway or narrowly
  scoped entity would not carry.
- R18. The condensed close checklist **MUST** gain a line covering, for a
  composite root, whether the domain's constituent kinds are each modeled or
  consciously accounted for rather than silently missing.

### Consistency

- R19. All added content **MUST** preserve domain agnosticism per
  [`AGENTS.md`](../../../AGENTS.md): examples stay illustrative and
  non-exhaustive, agentic language is scoped to context of use, and software
  product quality is not implied as the default modeled domain. The protective
  pair and the lifecycle band **MUST** read as stewardship concerns of caring for
  *any* entity, not as software-specific phases.

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` — add a best-practice
  coverage requirement for enumerating domain constituent kinds by the two axes,
  the three-projections rule, and the secure/safeguard protective pair (per
  R1–R13 above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — extend the
  area-and-factor-shape check contract with the missing-domain-constituent
  finding, earned-not-roster (per R16–R18 above).

### To rename

None.

### To delete

None.
