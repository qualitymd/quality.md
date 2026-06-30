---
type: Change Case
title: Public Review and Improve Workflows
description: Make review and improve first-class public /quality workflows, routed by focus, and simplify README framing around evaluate, review, and improve.
status: Done
tags: [skill, workflows, docs, ux]
timestamp: 2026-06-27T00:00:00Z
---

# Public Review and Improve Workflows

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0143-public-review-improve-workflows/spec.md) - what the
  case must do.
- [Design doc](0143-public-review-improve-workflows/design.md) - how it is
  built, and why.

## Motivation

QUALITY.md is easiest to explain as three things users do with a shared quality
model: evaluate, review, and improve. The existing skill has a well-defined
`evaluate` workflow, setup/update support workflows, direct model authoring, and
recommendation follow-up, but the public surface does not yet reflect the simpler
value proposition.

`review` and `improve` should become immediate public workflows even while their
deeper designs remain intentionally shallow. A new or returning user should be
able to invoke `/quality review` or `/quality improve`, have the skill quickly
infer the focus, confirm the intended path, and then either route into an
existing safe path or stop at a clear stub boundary. The README should teach that
same shape first, with setup as preparation and update as maintenance.

## Scope

Covered:

- Add public `/quality review` and `/quality improve` workflows.
- Introduce **focus** as the user-facing routing concept for review and improve.
- Provide basic runtime and durable-spec stubs for both workflows.
- Route existing model-review, model-improvement, recommendation-follow-up, and
  ad hoc quality concern language through the new public workflow names where
  applicable.
- Keep `evaluate` unchanged as the current bounded Evaluation workflow.
- Simplify README framing around evaluate, review, and improve as the primary
  value proposition.

Deferred:

- Deep review workflow design.
- Deep improve workflow design.
- New CLI commands for review or improve.
- New persisted review/improve artifacts beyond existing quality-log and issue
  handoff behavior.
- Recurring/cadenced review automation.

## Affected artifacts

Derived by sweeping for public workflow dispatch, `review`, `improve`,
recommendation follow-up, direct model authoring, workflow indexes, and README
usage/value-prop framing.

**Code**

- [x] No planned Go or TypeScript code impact; this is a skill/spec/docs
      behavior change.

**Durable specs** (substance in the [functional spec](0143-public-review-improve-workflows/spec.md))

- [x] `specs/skills/quality-skill/quality-skill.md` - add public review/improve
      invocation, focus routing, and dispatch rules.
- [x] `specs/skills/quality-skill/workflows/index.md` - list review and improve.
- [x] `specs/skills/quality-skill/workflows/review.md` - add the review workflow
      behavioral stub.
- [x] `specs/skills/quality-skill/workflows/improve.md` - add the improve
      workflow behavioral stub.
- [x] `specs/skills/quality-skill/recommendation-follow-up.md` - reconcile
      wording so recommendation follow-up is the existing implementation route
      used by `/quality improve` when compatible recommendation artifacts exist.

**Durable docs / bundled skill runtime**

- [x] `README.md` - foreground evaluate, review, and improve; demote setup to
      preparation, update to maintenance, and CLI to support tooling.
- [x] `CHANGELOG.md` - note the public `/quality review` and `/quality improve`
      workflow stubs.
- [x] `skills/quality/SKILL.md` - route public review/improve invocations and
      define focus confirmation behavior.
- [x] `skills/quality/workflows/index.md` - list review and improve.
- [x] `skills/quality/workflows/review.md` - add runtime workflow stub.
- [x] `skills/quality/workflows/improve.md` - add runtime workflow stub.
- [x] Runtime/spec logs for changed skill, workflow, and docs surfaces.

No planned impact: `SPECIFICATION.md`, Evaluation data schema, lint rules,
model parsing, scaffold output, install mechanics, or archived Change Cases.

## Status

`Done`. Implemented, verified, and archived. Runtime skill guidance, durable
skill specs, README framing, release notes, and logs are aligned.
