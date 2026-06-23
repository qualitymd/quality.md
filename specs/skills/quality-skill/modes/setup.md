---
type: Functional Specification
title: /quality setup
description: Behavioral component spec for context-informed QUALITY.md setup through the /quality skill.
tags: [skill, quality, mode, setup]
timestamp: 2026-06-23T00:00:00Z
---

# /quality setup

`setup` is the `/quality` skill mode that creates or updates a useful first
`QUALITY.md` through a short, context-informed discovery flow. It implements the
shared contracts in the parent [/quality skill](../quality-skill.md) spec and
owns only the setup-specific behavior below.

The runtime procedure lives at
[`skills/quality/modes/setup.md`](../../../../skills/quality/modes/setup.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`setup` is selected when no model file is present, when the user explicitly asks
to create or initialize a QUALITY.md file, or when read-only orientation routes
to bootstrap or first-population work.

The mode's purpose is to produce or improve a valid, useful project-specific
model while keeping setup's mutation boundary narrow. It is not an evaluation
mode and does not rate evaluated source.

## Mutation surface and artifacts

`setup` may mutate only the target `QUALITY.md` model file.

`setup` **MUST NOT** run evaluation, create evaluation artifacts, write the
quality log, create external issues, configure issue trackers, create CI or
release workflows, create scheduled automations, configure Codex automations, or
configure Claude Code routines.

## Context analysis and discovery

`setup` **MUST** inspect available repository context before asking setup
questions. Relevant context includes README and docs, repository structure,
package metadata, tests, contributor docs, existing agent instructions, and
visible workflow or work-management signals. This inspection **MUST** stay
bounded to setup signals and **MUST NOT** become source-quality evaluation.

Setup **MUST** treat the current directory as the default root area convention
unless the user supplied an explicit model path or repository context strongly
indicates a narrower root.

`setup` **MUST** present inferred setup defaults with confidence signals that
distinguish strongly inferred, weakly inferred, and assumed defaults.

Setup discovery **MUST** cover lifecycle, risk tolerance, modeling rigor,
collaboration context, and a needs map. It **MUST** treat lifecycle, risk
tolerance, and modeling rigor as separate setup concepts. Collaboration context
**MUST** assume agent-heavy development and ask which human collaborators or
stakeholders also need alignment.

The missing-context prompt **MUST** be seeded from context analysis rather than
asked as a blank "anything else?" question.

Optional work-management and recurring-review questions **MUST** be framed as
context capture, not permission to create external artifacts or automations. If
setup asks about recurring quality review, ad hoc `/quality evaluate` **MUST**
be treated as always available rather than a selectable automation option.
Setup **MUST NOT** recommend CI or release gating as the default quality loop.

## Model authoring

`setup` **MUST** drive `qualitymd init` for deterministic scaffolding when the
model file is absent. It **MUST NOT** reimplement scaffolding, validation, CLI
installation tooling, or source-driven authoring judgment.

After discovery and scaffolding when needed, `setup` **MUST** write a model that
follows the authoring guide and active specification. The model **MUST** address
the Markdown body's Overview, Scope, Needs, and Risks, including each section's
unknowns, open questions, and any material support that is not
agent-accessible. The body **MUST** preserve setup assumptions where they shape
the model: lifecycle, risk tolerance, modeling rigor, collaboration context,
stakeholder needs, and important missing or non-agent-accessible context.

Setup-authored factors **MUST** derive from project-specific needs and risks,
not generic quality labels alone. Child areas **SHOULD** be added only when they
represent distinct evaluated entities. Starter requirements **MUST** be
concrete and assessable from agent-accessible evidence or explicitly name
missing evidence or assessment constraints.

Setup **SHOULD** use the standard rating scale unless discovery shows a real
need to customize it.

Setup **SHOULD** include a `quality-md` area that evaluates the `QUALITY.md`
artifact itself against the active authoring guide unless the user declines or
the model file is not in the root area it governs. The area **SHOULD** use the
key `quality-md`, a title of the form `<Root Title> QUALITY.md`, an area
`description`, and an explicit path-based `source` such as `./QUALITY.md`. It
**MUST NOT** use prose aliases such as `(this file)` for `source`.

When setup adds that area, it **SHOULD** include concise YAML comments that
distinguish the area `source` from the requirement `assessment`. It **SHOULD**
use one area-level requirement with `factors` when the active authoring guide
defines one coherent judgment across multiple factors.

## Stop conditions

`setup` **MUST** stop before CLI-dependent work when the `qualitymd` CLI is
missing, outside the released-install SemVer range declared by the skill, or a
local development build lacks required commands.

`setup` **MUST** run `qualitymd lint` after writing `QUALITY.md`. It **MUST**
report lint failures before offering evaluation as a next step.

`setup` **MUST** inspect the resulting model against the bundled Top 10
QUALITY.md checks before reporting completion. This inspection **MUST** remain a
model-readiness inspection and **MUST NOT** evaluate root area source quality.

## Completion criteria

`setup` is complete when the target model exists, lint has run, the model has
received context-informed authoring or a clearly reported user-deferred
authoring step, and setup has reported model readiness. Completion output
**MUST** summarize the `QUALITY.md` change, lint result, readiness
classification, and next-step choices.

Next-step choices **SHOULD** include continuing to iterate on `QUALITY.md`,
running evaluation, setting up a recurring quality review loop, setting up
recommendation handoff, and stopping. `setup` **MUST NOT** automatically take
any next-step action.
