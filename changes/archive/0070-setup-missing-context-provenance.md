---
type: Change Case
title: Setup missing-context provenance
description: Prevent /quality setup from treating tacit or low-evidence project context as understood.
status: Done
tags: [quality-skill, setup, agent-accessibility]
timestamp: 2026-06-23T00:00:00Z
---

# Setup missing-context provenance

This case was created at `Draft`, advanced through `Design`, `In-Progress`, and
`In-Review`, then landed as `Done` and archived with its child folder.

- [Functional spec](0070-setup-missing-context-provenance/spec.md) — what the
  case must do.
- [Design doc](0070-setup-missing-context-provenance/design.md) — how the
  skill/spec wording carries the change.

## Motivation

A real setup run produced a missing-context option that allowed the user to
"assume product purpose and ops context are sufficiently understood" even though
the repository did not make those facts agent-accessible. That weakens the
setup contract: tacit maintainer/operator knowledge may exist, but it is not
available evidence unless the user provides it or points setup to an accessible
source.

## Scope

This case tightens `/quality setup` missing-context discovery so material
low-evidence or unavailable context is either recorded as unknown/not
agent-accessible, explicitly provided during setup, or tied to evidence setup
missed. It is skill-only; the CLI does not generate discovery prompts.

Deferred: richer structured-question UI behavior, automated transcript tests,
and any attempt to enumerate every possible missing-context category.

## Affected artifacts

### Code

- None — no `qualitymd` CLI/Go behavior changes.

### Durable specs

- [specs/skills/quality-skill/workflows/setup.md](../../specs/skills/quality-skill/workflows/setup.md) —
  specify missing-context provenance and prohibit assumed-understood options for
  material low/no-evidence gaps.
- [specs/skills/quality-skill/workflows/log.md](../../specs/skills/quality-skill/workflows/log.md) —
  record the setup workflow spec revision.
- [specs/log.md](../../specs/log.md) — record the durable spec revision.

### Format spec

- None — QUALITY.md schema and format semantics do not change.

### Durable docs

- [CHANGELOG.md](../../CHANGELOG.md) — note the unreleased `/quality setup`
  behavior change.

### Bundled skill

- [skills/quality/workflows/setup.md](../../skills/quality/workflows/setup.md) —
  update the runtime setup playbook so missing-context choices preserve
  provenance and do not convert absent evidence into tacit assumptions.

### Install/scaffold files

- None.

## Status

`Done`. Implemented skill-only with durable spec and changelog sync, verified
with markdown formatting and the repo check gate, then archived.
