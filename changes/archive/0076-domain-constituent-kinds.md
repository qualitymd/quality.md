---
type: Change Case
title: Domain constituent kinds and stewardship concerns
description: Teach the authoring guidance to enumerate a composite root's constituents by the domain's stewardship concerns (lifecycle plus the secure/safeguard protective pair) and an audience×purpose axis, so high-leverage constituents like specs, tests, docs, and protective artifacts are modeled or consciously accounted for instead of silently missing.
status: Done
tags: [skill, guide, authoring, areas, constituents]
timestamp: 2026-06-24T00:00:00Z
---

# Domain constituent kinds and stewardship concerns

A **Change Case** to close a coverage gap the [0074 composite-root work](0074-composite-root-areas.md) left open. 0074 named the composite
root shape and the two recurring *use-context* constituents (agent harness,
QUALITY.md self-check), but gave **domain constituents** no equally prescriptive
treatment. The result: a setup-authored model enumerates constituents by walking
the repository's folders — modeling the services it can see — and never asks what
*kinds* of constituent the entity's domain implies. Constituents that are thin,
scattered, or missing (tests, specs/requirements, documentation, a threat model)
produce no folder to trip over, so their absence is invisible rather than
recorded as a finding.

This case adds a domain-agnostic generator for constituent kinds, built on two
axes: a set of **stewardship concerns** (a lifecycle band — discover, define,
realize, verify, enable, operate, maintain — plus a cross-cutting *protective*
pair, **secure** and **safeguard**) and an **audience × purpose** axis (the
generalization of Diátaxis). It adds the rule that a concern projects into the
model in up to three ways — as a factor, a constituent, and an audience — so the
shared names do not double-count. It carries the existing earn-it guardrail
unchanged: kinds are a prompt, not a quota.

Detail lives in:

- [Functional spec](0076-domain-constituent-kinds/spec.md) — what the guidance
  must say.
- [Design doc](0076-domain-constituent-kinds/design.md) — the two-axis taxonomy,
  the secure/safeguard split, the three-projections rule, and why it is shaped
  this way.

## Motivation

A real evaluated entity is usually a composite (0074). Deciding the root is
composite raises the next question — *which* constituents — and the authoring
guide answers it only indirectly. "Ground high-leverage concerns in normative
artifacts" names requirements/intent, data quality, and interface contracts as
software examples and says to carry a germane high-leverage concern even when its
anchor is missing; but it is framed around normative anchors, not as a way to
*enumerate the constituent kinds a domain implies*, and it is hedged repeatedly
against "a universal roster of areas every model must carry."

That hedge is correct across domains and overcorrected within one. Once the
domain is named (software product, data set, reference document, service), there
*is* a recognizable set of constituent kinds, and the guide stops one step short
of letting the author use it. The asymmetry is the gap: factor coverage is
prescribed sharply ("aim for ~ten factors; fewer than eight triggers a coverage
review"), and the use-context constituents are prescribed sharply (a whole 0074
subsection), but domain-constituent coverage is left to "derive from component
boundaries and available evidence" — exactly the improvisation that drops tests
and specs when no folder exists.

The fix has two parts a domain-agnostic generator must supply:

- **A way to enumerate kinds without naming a domain.** Caring for any entity
  carries a recurring set of *stewardship concerns*; each tends to leave an
  authored, inspectable artifact that is a candidate constituent. Naming the
  concerns (not the artifacts) keeps the generator domain-neutral — "enable its
  audiences" rather than "docs", which already presumes a written-document
  domain.
- **A way to see when one concern is several constituents.** The same concern
  serves different audiences and purposes; documentation is not one area but
  (per Diátaxis) up to four with different factor families. The audience side is
  derivable from the body's Needs, which already name the stakeholders.

## Scope

Covered:

- Add an authoring-guide subsection, "Cover the domain's constituent kinds",
  parallel to "Carry the recurring use-context constituents", that:
  - enumerates constituent kinds by two axes — a **stewardship-concern** axis
    (lifecycle band: discover, define, realize, verify, enable, operate,
    maintain; plus a cross-cutting **protective** pair: **secure** — guard the
    entity from harm by the world; **safeguard** — guard stakeholders and the
    environment from harm by the entity) and an **audience × purpose** axis;
  - cites Diátaxis once as the audience×purpose lens applied to the *enable*
    concern;
  - states the **three-projections rule**: a concern projects as a factor, a
    constituent, and an audience; model the projection you mean and do not
    double-count (the security *of* an area is a factor; a security policy is a
    constituent);
  - derives the audience side from the body's Needs;
  - reuses the "carry a germane high-leverage kind even when its artifact is
    thin or missing, recording the gap as a finding" rule and links it to
    "Ground high-leverage concerns in normative artifacts";
  - keeps the earn-it guardrail (owned inspectable artifact, divergent factor
    family, traced to a Need or Risk; prompt not quota);
  - includes one illustrative, domain-scoped table (software product as one
    column) mapping each concern to instances and conventional factor families.
- Add a cross-reference from "Carry the recurring use-context constituents" to
  the new subsection where it already mentions domain constituents.
- Add a setup model-building step: when the root is composite, enumerate the
  domain's constituent kinds (walk the stewardship concerns and the audiences
  the Needs name) and account for each — model, defer in scope, mark out of
  scope, or record as an unknown — carrying germane high-leverage kinds as Areas
  even when thin/missing.
- Add a Top 10 Check 8 routing finding for a domain that implies constituent
  kinds the body evidences but the model neither models nor accounts for, and a
  matching line in the condensed close checklist; both earned-not-roster.
- Align the two governing guide specs and the guides log.

Deferred / non-goals:

- No QUALITY.md format change and no `SPECIFICATION.md` change; "stewardship
  concerns", "constituent kinds", "secure/safeguard", and "audience × purpose"
  stay authoring-guide vocabulary, not format semantics.
- No new mandatory areas; every constituent kind stays earned, not required.
- No change to evaluation, reporting, or CLI behavior.
- No heavy treatment in getting-started.md beyond an optional one-line pointer;
  that guide is deliberately thin for starter models.

## Affected artifacts

### Code

- [x] None — documentation-only change.

### Durable specs

- [x] `specs/skills/quality-skill/guides/authoring-md.md` — add a best-practice
      coverage requirement for enumerating domain constituent kinds by the two
      axes, the three-projections rule, and the secure/safeguard protective
      pair, with a 0076 rationale note.
- [x] `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` — extend
      the area-and-factor-shape check contract with the missing-domain-constituent
      finding, earned-not-roster.
- [x] `specs/skills/quality-skill/guides/log.md` — record the guide-spec updates.

### Format spec

- [x] None — no change to `SPECIFICATION.md`.

### Durable docs (bundled skill)

- [x] `skills/quality/guides/authoring.md` — add the "Cover the domain's
      constituent kinds" subsection and the cross-reference from "Carry the
      recurring use-context constituents".
- [x] `skills/quality/guides/top-10-quality-md-checks.md` — add the Check 8
      missing-domain-constituent finding and the condensed-checklist line.
- [x] `skills/quality/workflows/setup.md` — add the composite-root
      constituent-enumeration model-building step.

### Release

- [x] `CHANGELOG.md` — add the `/quality Skill` note under `Unreleased`.

## Children

- [Functional spec](0076-domain-constituent-kinds/spec.md) — required guidance
  content and acceptance criteria.
- [Design doc](0076-domain-constituent-kinds/design.md) — the taxonomy,
  alternatives, and trade-offs.

## Status

`Done`. Added the "Cover the domain's constituent kinds" authoring subsection and
the use-context cross-reference, the composite-root enumeration step in setup, the
Top 10 Check 8 missing-domain-constituent finding and condensed-checklist line,
aligned both guide specs and the guides log, and added the CHANGELOG note.
Documentation-only (no `SPECIFICATION.md` or code change). Verified with
`mise run check` (markdown format, bundle link resolution, lint, and Go gates all
pass). Archived.
