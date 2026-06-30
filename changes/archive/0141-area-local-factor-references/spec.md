---
type: Functional Specification
title: Area-local Factor References - functional spec
description: What the format, lint command, and /quality guidance must do to make Requirement factor references Area-local.
tags: [format, lint, skill, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Area-local Factor References - functional spec

Companion to the
[Area-local Factor References](../0141-area-local-factor-references.md) change
case. This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Ancestor Factor references make a `QUALITY.md` model partly graph-shaped. A
Requirement declared in one Area can directly affect a Factor owned by an
ancestor Area, so the ancestor Factor's contributors are not discoverable from
that Area's own Factor tree. That creates hidden roll-up dependencies and nudges
authors toward using ancestor Factors as wildcard roll-up targets. Area-local
references keep ownership and evaluation easier to reason about: Requirements
belong to their declaring Area, and cross-cutting root judgment is modeled by
root Requirements or root sub-Factors.

The remaining ambiguity is inside one Area's own Factor tree. `Requirement.factors`
entries are scalar Factor names, not paths, so two Factors with the same name
inside the same Area's recursive Factor tree are confusing even though their
canonical Factor IDs are distinct. The linter should warn authors about that
ambiguity without rejecting valid path-addressable Factors.

## Scope

Covered: format semantics, lint validation, lint documentation, focused tests,
runtime `/quality` authoring guidance, durable skill guide specs, shared UX
example guidance, scaffold comments, logs, and release notes.

Deferred:

- path-qualified `Requirement.factors` entries;
- automatic migration or rewrite of existing models;
- changing canonical qualified Factor references;
- evaluation JSON schema changes; and
- numeric roll-up semantics.

## Requirements

### Format semantics

- A Requirement's explicit `factors` entries **MUST** resolve only to Factors
  declared on the Area where that Requirement is declared.

  > Rationale: Area-local resolution keeps Factor ownership local and prevents
  > descendant Requirements from silently contributing to ancestor Factor ratings.
  >
  > Durable spec: modify `SPECIFICATION.md` - replace ancestor-inclusive Factor
  > reference scope with declaring-Area-only scope.

- A Requirement's explicit `factors` entries **MUST NOT** resolve to Factors
  declared on ancestor Areas, sibling Areas, descendant Areas, or unrelated
  Areas.

  > Durable spec: modify `SPECIFICATION.md` - state the out-of-scope locations
  > explicitly so agents do not infer cross-Area roll-ups.

- Factors with the same name on different Areas **MUST** remain distinct valid
  Factors.

  > Durable spec: modify `SPECIFICATION.md` - preserve local Factor identity
  > across Areas.

- The format specification **SHOULD** advise authors to avoid repeating the same
  Factor name within one Area's recursive Factor tree unless they accept that
  scalar `factors` references may be ambiguous to readers.

  > Durable spec: modify `SPECIFICATION.md` - add advisory same-Area duplicate
  > Factor-name guidance.

- Factor rating semantics **MUST NOT** describe descendant-Area Requirements as
  direct contributors to an ancestor Factor merely because they name that
  ancestor Factor.

  > Rationale: after Factor references are Area-local, ancestor Factor ratings
  > are driven by Requirements declared under or explicitly referencing Factors in
  > that same Area's Factor tree, not by descendant Area opt-ins.
  >
  > Durable spec: modify `SPECIFICATION.md` - remove the descendant-Area
  > explicit-reference report clause.

### Lint behavior

- `qualitymd lint` **MUST** report `unknown-factor` as an error when a
  Requirement's non-empty scalar `factors` entry does not resolve to any Factor
  name in the Requirement's declaring Area.

  > Durable spec: modify `specs/cli/lint-rules.md` - update the `unknown-factor`
  > rule's enforced contract, description, and repair guidance.

- `qualitymd lint` **MUST** report `unknown-factor` as an error when a
  descendant Area Requirement references a Factor name that exists only on an
  ancestor Area.

  > Durable spec: modify `specs/cli/lint-rules.md` - add ancestor-only references
  > to the invalid cases covered by `unknown-factor`.

- `qualitymd lint` **MUST NOT** report `unknown-factor` when a Requirement
  references a Factor name that exists in its declaring Area, even when an
  ancestor Area also has a Factor with the same name.

  > Durable spec: modify `specs/cli/lint-rules.md` - preserve valid local
  > shadowing across Areas.

- `qualitymd lint` **SHOULD** report a warning when the same Factor name appears
  more than once inside one Area's recursive Factor tree.

  > Rationale: canonical Factor IDs remain path-based, so this is not structural
  > invalidity; the warning protects authors from ambiguous scalar Factor
  > references and unclear prose.
  >
  > Durable spec: modify `specs/cli/lint-rules.md` - add a warning rule for
  > same-Area duplicate Factor names.

- The same-Area duplicate Factor-name warning **MUST NOT** fire for same-named
  Factors declared on different Areas.

  > Durable spec: modify `specs/cli/lint-rules.md` - scope the warning to one
  > declaring Area's Factor tree.

- `empty-factor` warning analysis **MUST** use the same Area-local explicit
  Factor-reference resolver as `unknown-factor`.

  > Durable spec: modify `specs/cli/lint-rules.md` - keep empty-Factor analysis
  > aligned with the Factor-reference contract.

### Skill and authoring guidance

- Runtime `/quality` Requirement authoring guidance **MUST** tell agents that
  `factors` names only Factors declared in the same Area as the Requirement.

  > Durable spec: modify
  > `specs/skills/quality-skill/guides/authoring/requirements.md` - require
  > Area-local Factor-reference guidance.

- Runtime `/quality` Requirement authoring guidance **MUST NOT** tell agents they
  can reach up to ancestor Factors.

  > Durable spec: modify
  > `specs/skills/quality-skill/guides/authoring/requirements.md` - remove
  > ancestor reach-up doctrine.

- Runtime `/quality` Factor authoring guidance **MUST** keep same-named Factors
  in different Areas valid while advising authors to avoid repeating the same
  Factor name within one Area's recursive Factor tree.

  > Durable spec: modify
  > `specs/skills/quality-skill/guides/authoring/factors.md` - add local
  > duplicate-name warning doctrine.

- Shared agent-mediated UX examples **MUST NOT** show descendant Requirements
  being connected to a model-wide Factor for roll-up.

  > Durable spec: none - this is durable docs guidance, not a `specs/` bundle
  > contract.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` - replace ancestor-inclusive Factor-reference scope with
  declaring-Area-only resolution, preserve cross-Area Factor identity, add
  same-Area duplicate-name guidance, and remove descendant-Area explicit
  contribution report semantics.
- `specs/cli/lint-rules.md` - update `unknown-factor`, add a same-Area duplicate
  Factor-name warning, and align `empty-factor` terminology.
- `specs/skills/quality-skill/quality-skill.md` - keep scoped Factor evaluation
  wording local to the Factor's declaring Area.
- `specs/skills/quality-skill/guides/authoring/factors.md` - require guidance
  for same-named Factors across Areas and duplicate names within one Area.
- `specs/skills/quality-skill/guides/authoring/requirements.md` - require
  Area-local Factor-reference guidance.

### To rename

None

### To delete

None

## Verification

- `go test ./internal/lint` **MUST** pass.
- `go test ./internal/cli ./internal/status` **SHOULD** pass to catch lint/status
  fixture fallout.
- `go test ./...` **SHOULD** pass before archiving.
- `mise run fmt-md-check` **SHOULD** pass after Markdown updates.
- Source inspection **MUST** show no live authoring guidance that says
  Requirement factor references can reach ancestor Factors.
