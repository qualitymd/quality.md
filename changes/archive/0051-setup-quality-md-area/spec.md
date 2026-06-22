---
type: Functional Specification
title: Setup quality-md Area — functional spec
description: Requirements for adding a setup-authored quality-md Area and supporting authoring-guide guidance.
tags: [skill, setup, quality-model]
timestamp: 2026-06-22T00:00:00Z
---

# Setup quality-md Area — functional spec

Companion to [Setup quality-md Area](../0051-setup-quality-md-area.md). This
spec states the delta for `/quality setup` and the authoring guide.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in
[RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Background

The setup flow already creates a valid skeleton through `qualitymd init`, then
uses the skill's authoring judgment to populate a useful first model. The
`QUALITY.md` file itself is part of that root area's maintainability surface:
when it is vague, unsupported, stale, or unassessable, future evaluations and
improvements degrade. The initial model should therefore make the model file's
own quality inspectable without moving that judgment into the deterministic CLI
scaffold.

## Scope

In scope:

- Setup-time creation or proposal of a `quality-md` Area during guided
  population.
- Human-facing YAML comments that explain the self-referential-looking pattern.
- Authoring-guide guidance for quality-attribute Factor names.
- Authoring-guide guidance for one referenced assessment connected to multiple
  Factors.

Deferred:

- Any change to `qualitymd init`, `internal/scaffold/skeleton.md`, or CLI
  scaffold templates.
- A new source alias such as `(this file)`.
- A new schema field for assessment references.

## Requirements

### Setup quality-md Area

After a valid skeleton is available and before routing to wizard, `/quality
setup` **SHOULD** include a `quality-md` Area that evaluates the active
`QUALITY.md` file itself against the active authoring guide, unless the user
declines or the model file is not in the root area it governs.

The `quality-md` Area **SHOULD** use the concrete model artifact as its identity:
the Area key `quality-md`, a title of the form `<Root Title> QUALITY.md`, an
Area `description`, and an explicit path-based `source` pointing at the model
file, such as `./QUALITY.md`.

The `quality-md` Area **MUST NOT** use prose aliases such as `(this file)` for
`source`; `source` remains a normal machine-resolvable selector.

The `quality-md` Area **SHOULD** include concise YAML comments explaining that
the Area `source` is the `QUALITY.md` artifact being evaluated, while the
Requirement `assessment` names the authoring guide used to judge it.

The `quality-md` Area **SHOULD** use one Area-level Requirement that references
the active authoring guide and lists multiple Factors through `factors` when the
guide defines one coherent judgment that bears on all of them.

### Authoring guide

The authoring guide **MUST** teach authors to name Factors as quality
characteristics the Area can exhibit to a degree, not as practices, workflows,
lifecycle phases, authoring techniques, or evaluation tactics.

The authoring guide **SHOULD** give examples that help authors choose narrower
attributes such as `completeness`, `consistency`, `credibility`, `currentness`,
`understandability`, `traceability`, `assessability`, `maintainability`,
`modifiability`, and `testability` when those names fit the Area and concern.

The authoring guide **MUST** teach that when one guide, spec, or checklist
defines a coherent assessment that bears on several Factors, authors should
write one Requirement, connect it to the affected Factors through `factors`, and
reference the governing entity once. It **MUST** teach authors to split such
Requirements only when the referenced entity defines claims whose results could
legitimately diverge.

## Durable spec changes

### To add

None.

### To modify

- [specs/skills/quality-skill/quality-skill.md](../../../specs/skills/quality-skill/quality-skill.md)
  — add the setup-time `quality-md` Area default and YAML-comment guidance.
- [specs/skills/quality-skill/guides/authoring-md.md](../../../specs/skills/quality-skill/guides/authoring-md.md)
  — add the Factor naming and one referenced assessment / many-Factors guide
  requirements.

### To rename

None.

### To delete

None.
