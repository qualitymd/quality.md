---
type: Change Case
title: Named Requirement identity
description: Give Requirements stable id-like names, display titles, and qualified model references.
status: Done
tags: [format, requirements, model-references]
timestamp: 2026-06-25T00:00:00Z
---

# Named Requirement identity

This change case updates the QUALITY.md model so Requirements have stable,
id-like names instead of using the natural-language Requirement statement as the
map key. It also introduces qualified Requirement model references for durable
evaluation records and future evaluation v2 data.

Children:

- [Functional spec](0093-requirement-identity/spec.md) - what the change must do.
- [Design doc](0093-requirement-identity/design.md) - how the change is built.

## Motivation

Evaluation v2 needs durable, machine-readable references to Requirements and to
Requirement-derived records. The current format identifies a Requirement by its
natural-language map key. That is friendly for authoring, but weak for durable
records: wording edits change identity, long statements make awkward paths and
references, and requirement-level findings, unknowns, and ratings need stable
anchors.

Named Requirements would align Requirements with Areas and Factors: the map key
is a stable id, while `title`, optional `description`, and `assessment` carry the
human-facing content.

## Scope

This case covers the QUALITY.md model shape for Requirements, the formal
Requirement model-reference grammar, validation rules, scaffold/examples, and
the agent guidance that authors or evaluates Requirements.

This case makes a breaking format change: current statement-key Requirements are
not accepted as a compatibility shape. Authors and agents migrate them by giving
each Requirement a stable name, a required `title`, and the retained
`assessment`.

## Affected artifacts

### Format spec

- [x] `SPECIFICATION.md` - define named Requirements, `title`, optional
      `description`, retained `assessment`, and qualified Requirement references.

### Durable specs

- [x] `specs/quality-schema-json.md` - update the JSON Schema contract for the
      named Requirement shape.
- [x] `specs/cli/lint-rules.md` - update Requirement validation, missing title,
      name grammar, uniqueness, and compatibility diagnostics.
- [x] `specs/cli/init.md` - update scaffolded examples if the starter model
      includes Requirements.
- [x] `specs/evaluation-v2/` - no durable evaluation-v2 spec exists yet; the
      active sketch and future specs use the formal Requirement reference grammar.

### Code

- [x] `internal/model/` - update the typed Requirement model.
- [x] `internal/lint/` - validate named Requirements, titles, factor
      connections, ratings, and compatibility cases.
- [x] `internal/schema/` - update generated schema support.
- [x] `internal/scaffold/` - update starter QUALITY.md content and comments.
- [x] `internal/status/` and evaluation/reporting code - adjust model-shape
      traversal and references where Requirements are surfaced.

### Bundled skill and durable skill specs

- [x] `skills/quality/guides/authoring/requirements.md` - update Requirement
      authoring guidance.
- [x] `skills/quality/guides/authoring/model-structure.md` - reviewed; no
      model-structure-specific stale Requirement-shape guidance found.
- [x] `skills/quality/workflows/setup.md` - reviewed; no setup-specific stale
      Requirement-shape examples found.
- [x] `specs/skills/quality-skill/guides/authoring/requirements.md` - mirror the
      durable skill guidance.
- [x] `specs/skills/quality-skill/workflows/setup.md` - reviewed; no
      setup-specific durable mirror change needed.

### Public docs and examples

- [x] `README.md` - update the example QUALITY.md and any prose that describes
      Requirement identity.
- [x] `docs/guides/model-quality-across-domains.md` - update example model
      content if Requirements appear.
- [x] Other docs/examples found by searching for `requirements:` and
      `assessment:` once the final shape is settled.

## Status

`Done`. Implemented, verified with `mise run check`, and archived.
