---
type: Change Case
title: Setup open-ended review gate
description: Make /quality setup's final review prompt friendly and open-ended before authoring QUALITY.md.
status: Done
tags: [quality-skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup open-ended review gate

This case was created at `Draft`, advanced through `Design`, `In-Progress`, and
`In-Review`, then landed as `Done` and archived with its child folder.

- [Functional spec](0071-setup-open-ended-review-gate/spec.md) - what the case
  must do.
- [Design doc](0071-setup-open-ended-review-gate/design.md) - how the runtime
  skill and durable spec carry the wording change.

## Motivation

A real `/quality setup` run ended the final recap with `Reply "looks good" (or
note any corrections) and I'll author QUALITY.md.` The behavior was technically
right: setup waited for confirmation and allowed corrections. The wording was
too narrow. It framed the user's useful last-call input as either approval or
correction, when setup benefits from broader context such as priorities, worries,
preferred wording, edge cases, and repo-invisible facts.

## Scope

This case changes the setup final review gate prompt to be more colloquial,
friendly, and open-ended while preserving the explicit confirmation fast path.
It remains a skill-only UX change; setup still waits for a user response before
writing `QUALITY.md`.

Deferred: structured question-tool UI affordances, transcript automation, and
changes to setup discovery questions.

## Affected artifacts

### Code

- None - no `qualitymd` CLI/Go behavior changes.

### Durable specs

- [specs/skills/quality-skill/workflows/setup.md](../../specs/skills/quality-skill/workflows/setup.md) -
  specify the open-ended final review prompt and broad last-call input handling.
- [specs/skills/quality-skill/workflows/log.md](../../specs/skills/quality-skill/workflows/log.md) -
  record the setup workflow spec revision.
- [specs/log.md](../../specs/log.md) - record the durable spec revision.

### Format spec

- None - QUALITY.md schema and format semantics do not change.

### Durable docs

- [CHANGELOG.md](../../CHANGELOG.md) - note the unreleased `/quality setup`
  review-gate wording change.

### Bundled skill

- [skills/quality/workflows/setup.md](../../skills/quality/workflows/setup.md) -
  update the runtime setup playbook's final review gate wording.

### Install/scaffold files

- None.

## Status

`Done`. Implemented skill-only with durable spec and changelog sync, verified
with markdown formatting and the repo check gate, then archived.
