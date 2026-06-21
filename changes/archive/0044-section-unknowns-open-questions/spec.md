---
type: Functional Specification
title: Section unknowns and open questions - functional spec
description: Requirements for replacing the standalone Known gaps body section with per-section unknowns, open questions, and a review-provenance state line.
tags: [authoring, body, skill, scaffold, spec]
timestamp: 2026-06-20T00:00:00Z
---

# Section unknowns and open questions - functional spec

Companion to
[Section unknowns and open questions](../0044-section-unknowns-open-questions.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

The Markdown body lacked a consistent per-section structure, and a single
catch-all **Known gaps** section held uncertainty far from the content it
qualified. Because a QUALITY.md body is largely agent-authored, readers also had
no signal for what a human had actually vetted. This change makes the body shape
consistent, scopes uncertainty to the section it concerns, and records review
provenance so unreviewed agent edits are visible.

## Scope

Covered: the recommended body sections, the section shape, per-section unknowns
and open questions, the review state line, and the propagation of those across
the format spec, the authoring guidance, the scaffold, and the skill checks.

Deferred / non-goals: no frontmatter schema change; no lint rule or CLI behavior
that parses or validates the state line or the unknowns/open-questions lines —
they remain freeform body prose.

## Requirements

### Recommended body sections

The recommended body sections are **Overview**, **Scope**, **Needs**, and
**Risks**. Guidance and scaffolding **MUST NOT** present **Known gaps** as a
recommended standalone section.

> Rationale: a separate catch-all section sat far from the content it qualified
> and was skimmed past; uncertainty belongs with the section it concerns. — 0044

### Section shape

Authoring guidance **MUST** describe a common shape for every body section:
purpose, contents, the section's unknowns and open questions, and a state line.

### Unknowns and open questions

Each section's unknowns and open questions are scoped to that section's own topic
and are context that feeds the model, not commentary on the model. Guidance
**MUST** keep the two distinct: an **unknown** is a broad area of uncertainty
about the section's topic that may not resolve to a single answer; an **open
question** is a specific question with one particular answer, still unresolved.

Guidance **MUST** instruct authors to record "none known" rather than omit the
unknowns or open questions, so an absence reads as considered.

> Rationale: on a high-leverage file an explicit "none" reads as considered; a
> blank reads as skipped. — 0044

### Review state line

Authoring guidance **MUST** describe a per-section state line that records the
last human review (citing a named person) distinctly from the last agent review
(naming the agent). The human review **MUST** advance only when a person reads
and endorses the section, never for an agent or mechanical edit. A section whose
state line carries no human review is unreviewed; guidance **MUST NOT** instruct
authors to backfill a human reviewer who did not endorse the section.

> Rationale: the body is largely agent-authored, so the only freshness signal
> worth trusting is when a person last stood behind the section. — 0044

### Author-declared context vs not assessed

The format spec **MUST** keep documented unknowns and open questions as
author-declared model context, distinct from a `not assessed` Rating Result
produced during evaluation when evidence is absent or insufficient.

### Skill alignment

The skill's first-population guidance and body-context check **MUST** reference
the recommended sections (Overview, Scope, Needs, Risks) and the per-section
unknowns and open questions, not a standalone Known gaps section.

### Scaffold

The `qualitymd init` scaffold **MUST** seed the recommended body sections without
a standalone Known gaps section, and **MUST** show each section's unknowns, open
questions, and state line so the starter file teaches the convention.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` — body context names unknowns and open questions instead of
  known gaps; keep them distinct from `not assessed` (per the recommended-sections
  and author-declared-context requirements above).
- `specs/cli/init.md` — recommended body sections are Overview, Scope, Needs,
  Risks (per the scaffold and recommended-sections requirements above).
- `specs/skills/quality-skill/guides/authoring.md` — recommended-section MUST and
  body best-practice coverage adopt per-section unknowns/open questions and the
  state line (per the section-shape, unknowns/open-questions, and state-line
  requirements above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks.md` — body-context
  check inspects the recommended sections and their unknowns/open questions (per
  the skill-alignment requirement above).
- `specs/skills/quality-skill/quality-skill.md` — setup first-population body list
  drops Known gaps (per the skill-alignment requirement above).

### To delete

None
