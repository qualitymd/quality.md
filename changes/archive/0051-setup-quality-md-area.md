---
type: Change Case
title: Setup quality-md Area
description: Add a setup-authored quality-md Area that evaluates the QUALITY.md model file itself against the active authoring guide.
status: Done
tags: [skill, setup, quality-model]
timestamp: 2026-06-22T00:00:00Z
---

# Setup quality-md Area

This Change Case makes `/quality setup` normally include a `quality-md` Area
that evaluates the project's `QUALITY.md` artifact itself, so first-run models
make their own authoring quality explicit.

Details:

- [Functional spec](0051-setup-quality-md-area/spec.md) — what setup and the
  authoring guide must do.
- [Design doc](0051-setup-quality-md-area/design.md) — how the skill prompt and
  guide changes deliver it.

## Motivation

`QUALITY.md` is an owned project artifact: it governs what agents and
maintainers evaluate, what evidence they collect, and how future improvements
are judged. If the initial model omits the quality of that artifact, the setup
flow leaves the model's own credibility, assessability, traceability, and
maintainability implicit. First-time users also need a concise explanation of
why an Area points at `./QUALITY.md` while its Requirement assesses against an
authoring guide.

## Scope

This case covers the setup-time `quality-md` Area pattern, the YAML comments that
explain `source` versus `assessment`, and the authoring-guide advice needed to
keep the Area well-modeled: quality-attribute Factor names and one referenced
assessment connected to multiple Factors.

Deferred:

- Changes to `qualitymd init` or `internal/scaffold/skeleton.md`; the CLI
  scaffold stays generic and context-light.
- A dedicated CLI flag or template for `quality-md` Area generation.
- Any new durable format field for "this file"; the Area uses normal
  path-based `source`.

## Affected artifacts

### Code

- None.

### Durable specs

- [specs/skills/quality-skill/quality-skill.md](../specs/skills/quality-skill/quality-skill.md)
  — specify the setup-time `quality-md` Area default and explanatory YAML
  comments.
- [specs/skills/quality-skill/guides/authoring-md.md](../specs/skills/quality-skill/guides/authoring-md.md)
  — specify Factor naming guidance and the one referenced assessment /
  many-Factors pattern.

### Format spec

- None.

### Durable docs

- None.

### Bundled skill

- [skills/quality/modes/setup.md](../skills/quality/modes/setup.md) — instruct
  setup to normally include the `quality-md` Area and comments during guided
  population.
- [skills/quality/guides/authoring.md](../skills/quality/guides/authoring.md)
  — teach quality-attribute Factor names and single referenced assessments across
  multiple Factors.

### CLI scaffold and install artifacts

- None.

## Status

`Done`. Landed and archived after updating the durable skill specs and bundled
skill guidance. No Go code or CLI scaffold work was in scope.
