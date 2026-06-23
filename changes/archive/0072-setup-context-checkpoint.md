---
type: Change Case
title: Setup context checkpoint
description: Replace /quality setup's final open-ended discovery questions with a compact human context checkpoint.
status: Done
tags: [quality-skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup context checkpoint

This case was created at `Draft`, advanced through `Design`, `In-Progress`, and
`In-Review`, then landed as `Done` and archived with its child folder.

- [Functional spec](0072-setup-context-checkpoint/spec.md) - what the case must
  do.
- [Design doc](0072-setup-context-checkpoint/design.md) - how the runtime skill
  and durable spec carry the checkpoint change.

## Motivation

A field run showed that `/quality setup`'s final four discovery questions were
important but hard to answer well. They asked for primary users, maintainers,
other stakeholders, and missing context as separate open-ended prompts, then
ended with the broadest question. Users are likely to answer only the final
prompt, leaving the most important human context implicit. Setup needs that
context for Needs, Risks, Unknowns, and model provenance, but it should make the
user's job correction rather than composition.

## Scope

This case changes setup discovery so the human context dimensions are presented
as one draft context checkpoint for confirmation, correction, or terse fill-in.
Setup still collects primary users/outcomes, maintainers/collaborators, other
stakeholders, and missing/not-agent-accessible context, and still records omitted
or low-confidence material context as Unknown instead of guessing.

Deferred: new structured UI affordances, CLI support for setup prompts, and
changes to the QUALITY.md schema.

## Affected artifacts

### Code

- None - no `qualitymd` CLI/Go behavior changes.

### Durable specs

- [specs/skills/quality-skill/workflows/setup.md](../../specs/skills/quality-skill/workflows/setup.md) -
  specify the human context checkpoint and its provenance/Unknown handling.
- [specs/skills/quality-skill/workflows/log.md](../../specs/skills/quality-skill/workflows/log.md) -
  record the setup workflow spec revision.
- [specs/log.md](../../specs/log.md) - record the durable spec revision.

### Format spec

- None - QUALITY.md schema and format semantics do not change.

### Durable docs

- [CHANGELOG.md](../../CHANGELOG.md) - note the unreleased `/quality setup`
  discovery UX change.

### Bundled skill

- [skills/quality/workflows/setup.md](../../skills/quality/workflows/setup.md) -
  update the runtime setup playbook's discovery prompt shape.

### Install/scaffold files

- None.

## Status

`Done`. Implemented skill-only with durable spec and changelog sync, verified
with markdown formatting and the repo check gate, then archived.
