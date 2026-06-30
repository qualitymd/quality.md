---
type: Change Case
title: Setup Factor Proposal Checkpoint
description: Teach and apply factor desiderata during setup so users can give better feedback on initial factor selection.
status: Done
tags: [skill, quality, setup, factors]
timestamp: 2026-06-29T00:00:00Z
---

# Setup Factor Proposal Checkpoint

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0166-setup-factor-proposal-checkpoint/spec.md) - what the case must do.
- [Design doc](0166-setup-factor-proposal-checkpoint/design.md) - how it is built, and why.

## Motivation

Setup already inspects repository context and asks discovery questions before
authoring `QUALITY.md`, but it moves from human context directly to final review.
That means the first factor set can be accepted without the user seeing the
criteria used to select it or the trade-off between coverage and maintainability.

The factor desiderata make that judgment teachable. A factor set should be
comprehensive, proportionate, and sustainable. Individual factors should be
consequential first, then refined until they are bounded, operational, traceable,
and neutral. Setup should use those qualities internally and present them in a
draft factor proposal so user feedback can be targeted: missing,
overemphasized, misplaced, or badly named concerns.

## Scope

Covered:

- Canonicalize factor-set and individual-factor desiderata in factor authoring
  guidance and the durable guide spec.
- Add a setup factor proposal checkpoint between discovery and final review.
- Teach the desiderata in setup without asking users to design factors cold.
- Add candidate factor-set quality and factor rationale fields to the setup
  brief.
- Make setup recap and model authoring consume the reviewed factor proposal.
- Capture a concise public README explanation.

Deferred:

- A separate algorithmic factor-selection procedure beyond the checkpoint and
  desiderata.
- CLI support for setup or model authoring.
- New generated templates or scaffolds for `QUALITY.md`.

Non-goals:

- Changing the QUALITY.md format schema.
- Changing evaluation behavior or report generation.
- Adding default factor checklists for any quality domain.

## Affected artifacts

Derived by sweeping for setup workflow contracts, factor authoring guide mirrors,
README factor guidance, and OKF bundle logs.

**Code**

- [x] No code impact.

**Durable specs**

- [x] `specs/skills/quality-skill/guides/authoring/factors.md` - require
      factor-set and individual-factor desiderata in the durable factor guide
      contract.
- [x] `specs/skills/quality-skill/workflows/setup.md` - require the factor
      proposal checkpoint, setup-brief fields, recap inclusion, model-authoring
      use, and important-gap check.
- [x] `specs/skills/quality-skill/quality-skill.md` - mention the setup factor
      proposal checkpoint in the parent workflow summary.

**Format spec**

- [x] No `SPECIFICATION.md` impact; this is authoring/workflow guidance, not
      schema semantics.

**Durable docs**

- [x] `README.md` - add concise public guidance for deliberate factor selection.

**Bundled skill**

- [x] `skills/quality/guides/authoring/factors.md` - add the canonical two-layer
      factor desiderata.
- [x] `skills/quality/workflows/setup.md` - add the factor proposal checkpoint
      and its runtime teaching/review shape.

**Install/scaffold/examples**

- [x] No impact.

**OKF bundle logs and indexes**

- [x] `changes/log.md`, `changes/index.md`, `changes/archive/index.md`.
- [x] `specs/log.md`, `specs/skills/quality-skill/workflows/log.md`, and
      `specs/skills/quality-skill/guides/log.md`.
