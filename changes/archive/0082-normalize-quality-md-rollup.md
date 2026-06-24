---
type: Change Case
title: Normalize QUALITY.md self-check roll-up
description: Remove the special out-of-roll-up treatment for the QUALITY.md self-check area so `quality-md` evaluates and rolls up like any other modeled area, while preserving quality-log handling for model changes.
status: Done
tags: [skill, authoring, evaluation, roll-up, quality-md]
timestamp: 2026-06-24T00:00:00Z
---

# Normalize QUALITY.md self-check roll-up

A **Change Case** to make the `quality-md` self-check area ordinary for
evaluation and roll-up. The current authoring guidance models the repository's
`QUALITY.md` artifact as a recurring use-context constituent, then carves it out
of the root area's roll-up on a separate "learn loop" axis. That exception is not
represented in the format, CLI, record schema, or report renderer, and it makes
full evaluation ambiguous.

This case removes the exception: when `quality-md` is modeled and in scope, it is
assessed, analyzed, reported, and rolled up like any other area. The remaining
distinction is follow-up behavior: applying a recommendation that changes
`QUALITY.md` may be a model change, so the existing quality-log rule still
applies.

Detail lives in:

- [Functional spec](0082-normalize-quality-md-rollup/spec.md) - what the guidance
  and evaluation contract must say.
- [Design doc](0082-normalize-quality-md-rollup/design.md) - why normal area
  semantics are the simpler design and why mutation logging is the right home for
  the remaining distinction.

## Motivation

The self-check exception creates a split-brain model: `quality-md` is a normal
area structurally, with `source: ./QUALITY.md`, factors, and requirements, but
guidance says it should not affect the root result. That special case hides a
first-class project artifact from the aggregate quality signal and forces agents
to decide whether "full evaluation" means the whole model tree or the model tree
minus one named area.

Treating `quality-md` normally aligns the guidance with the existing CLI and
report machinery. A weak, stale, or unassessable `QUALITY.md` is real project
quality debt and should be consequential in the aggregate result when the area is
in scope.

## Scope

Covered: authoring guidance and its durable spec mirror; evaluation workflow
guidance that full evaluation covers every in-scope modeled area uniformly; and
the repository's own `QUALITY.md` wording where it describes the role of the
`quality-md` area.

Deferred / non-goals: no change to the QUALITY.md format or schema, no new
frontmatter field for roll-up exclusion, no CLI/report schema or Go
implementation change, and no migration command. Re-evaluating this repository's
own full model after the guidance changes is a follow-up.

## Affected artifacts

Found by sweeping for `quality-md`, `QUALITY.md self-check`, `learn loop`, and
`out of roll-up` across live skill guidance, spec mirrors, report specs, CLI
code, and the dogfooded model. Grouped by kind; an empty kind is a deliberate "no
impact," not an oversight.

### Code

None - current CLI/status/report code treats areas generically and has no
`quality-md` special case.

### Format spec (`SPECIFICATION.md`)

None - this changes skill authoring and evaluation guidance, not the document
format.

### Durable specs (`specs/`)

The functional spec's [Durable spec changes](0082-normalize-quality-md-rollup/spec.md)
section is the authoritative breakdown. In summary:

- `specs/skills/quality-skill/guides/authoring-md.md` - remove the self-check
  out-of-roll-up rule and require ordinary area semantics.
- `specs/skills/quality-skill/evaluation.md` - clarify that full evaluation
  covers all in-scope modeled areas uniformly, including `quality-md`.

### Durable docs

- `QUALITY.md` - update the dogfooded model body if needed so the `quality-md`
  area is described as a normal first-class area, not an out-of-band learn-loop
  axis.

### Bundled skill (`skills/quality/`)

- `skills/quality/guides/authoring.md` - runtime counterpart of the authoring-md
  spec changes.
- `skills/quality/workflows/evaluate.md` - runtime counterpart of the full-scope
  evaluation clarification, if the implementation pass decides the existing
  generic wording is not explicit enough.

### Install / scaffold

None.

### Changelog

- `CHANGELOG.md` - note the guidance simplification when implemented.

## Children

- [Functional spec](0082-normalize-quality-md-rollup/spec.md) - what the guidance
  and evaluation contract must say.
- [Design doc](0082-normalize-quality-md-rollup/design.md) - normal area
  semantics, rejected alternatives, and residual risks.

## Status

`Done`. Implemented, verified, and archived. The runtime authoring/evaluation
guidance and durable spec mirrors now treat `quality-md` as an ordinary in-scope
area for evaluation and roll-up; quality-log behavior remains limited to
meaningful confirmed model changes. No `SPECIFICATION.md`, CLI schema, or Go code
change was needed.
