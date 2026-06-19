---
type: Functional Specification
title: Required display titles - functional spec
description: Require human-facing titles on model elements and make lint enforce them.
tags: [specification, schema, lint, report, skill]
timestamp: 2026-06-19T00:00:00Z
---

# Required display titles - functional spec

Companion to
[Required display titles](../0033-required-display-titles.md). This spec states
the format, lint, scaffold, and display-output delta for required model titles.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background / Motivation

`QUALITY.md` distinguishes stable model identity from human-readable
explanation, but that distinction is incomplete for display labels. A target key,
factor key, or rating `level` can be concise, slug-like, or scoped for
references, while reports and skill output need labels that read well in prose.

Required titles make the contract explicit: structural identifiers are for
paths, references, comparison, and JSON identity; `title` is for human display.
This lets renderers and the `/quality` skill consistently show readable labels
without guessing or falling back to cryptic identifiers in valid documents.

## Scope

Covered: `title` requiredness for the Model, Target, Factor, and Rating Level;
`Factor.title` as a new schema property; lint enforcement; scaffold validity;
and title-first display in human-facing report, status, and skill output.

Deferred / non-goals: no `title` property on Requirement, no title uniqueness
rule, no change to model paths or identity, no change to factor-reference syntax,
no change to rating override keys, and no change to evaluation roll-up.

## Requirements

The Model **MUST** declare a non-empty scalar `title`.

Every Target **MUST** declare a non-empty scalar `title`.

Every Factor **MUST** declare a non-empty scalar `title`.

Every Rating Level **MUST** declare a non-empty scalar `title`.

`Factor.title` **MUST** be added as a Factor property. It is the human display
label for the Factor. The Factor map key remains the stable identifier local to
the Target where the Factor is declared.

`Target.title` **MUST** be the human display label for the Target. The Target map
key remains the stable identifier within its parent `targets` map and remains
part of the Target path.

`RatingLevel.title` **MUST** be the human display label for the Rating Level. The
`level` value remains the stable identifier for rating results, gates, and
requirement-level `ratings` override keys.

The format **MUST NOT** add `Requirement.title`. The requirement statement
remains the Requirement's display text and map key.

> Rationale: adding a separate title would create two authored summaries for the
> same assessable expectation and invite drift. Requirement statements already
> carry the text rendered to readers.

Human-facing renderers and skill output **SHOULD** use `title` as the primary
display label for Models, Targets, Factors, and Rating Levels. They **MAY**
include the stable key, path, or `level` id as secondary context when needed for
disambiguation, links, gates, or machine traceability.

Machine-readable output **MUST** preserve stable identifiers and paths. Requiring
titles **MUST NOT** replace target paths, factor-reference names, `level` ids, or
rating override keys in JSON contracts that use them as identity.

`qualitymd lint` **MUST** report a `missing-title` error when a required `title`
is absent, null, or empty on the Model, a Target, a Factor, or a Rating Level.

`missing-title` finding messages **MUST** identify the specific element that
lacks the title, using context-specific wording for the model root, a target, a
factor, or a rating level.

`missing-title` **MUST NOT** be fixable. A placeholder title is scaffold content,
not a deterministic repair for an authored model.

`qualitymd init` **MUST** generate a scaffold that satisfies the new required
title contract, including placeholder titles for the model root, the seeded
rating levels, and the placeholder factor.

Examples, runtime authoring guidance, and skill-facing specs **MUST** show and
describe the required titles so new models follow the display-label contract by
default.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` - require titles on Model, Target, Factor, and Rating Level;
  add the `Factor.title` property; describe the stable identifier versus display
  label split; and update report semantics for title-first display (per the
  schema and display requirements above).
- `specs/cli/lint.md` - move `missing-title` to the error rule table and define
  it as a required-title rule for Model, Target, Factor, and Rating Level (per
  the lint requirements above).
- `specs/cli/init.md` - require scaffolded placeholder titles wherever the
  format requires them (per the init requirement above).
- `specs/cli/evaluation-build-report.md` - specify title-first human display
  while preserving stable identifiers for traceability (per the display
  requirements above).
- `specs/cli/status.md` - specify title-based human labels and identifier-based
  JSON paths in status output (per the display and machine-output requirements
  above).
- `specs/skills/quality-skill/quality-skill.md` - update the skill contract so
  setup, wizard, evaluation, and reporting guidance assume required titles and
  use them for human-facing labels (per the display requirements above).
- `specs/skills/quality-skill/authoring-guide.md` - update the guide contract so
  the runtime authoring guide teaches required display titles (per the examples
  and guidance requirement above).

### To delete

None
