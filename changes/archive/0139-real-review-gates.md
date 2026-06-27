---
type: Change Case
title: Real Review Gates
description: Make feedback invitations real gates in agent-mediated UX and align direct QUALITY.md authoring so it acknowledges quickly and waits for review before mutation.
status: Done
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Real Review Gates

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0139-real-review-gates/spec.md) - what the case must do.
- [Design doc](0139-real-review-gates/design.md) - how it's built, and why.

## Motivation

A direct `QUALITY.md` authoring request exposed two interaction failures. First,
the agent spent tens of seconds reading model and authoring context before
acknowledging the request. Second, the eventual checkpoint asked whether the user
wanted adjustments while also saying the agent was proceeding with the
interpretation. That wording offers agency and removes it in the same breath.

The shared agent-mediated UX guide should make the principle explicit:
acknowledge long inference before doing it, and treat feedback invitations as
real gates. The `/quality` direct-authoring path should then apply that principle
by acknowledging quickly, showing the intended edit, and waiting for `looks good`
or corrections before editing `QUALITY.md`.

## Scope

Covered:

- Add shared UX guide doctrine for immediate acknowledgement before long
  non-workflow reads or inference.
- Add shared UX guide doctrine that review/feedback invitations are real gates,
  not courtesy text.
- Apply the doctrine to `/quality` direct model authoring in runtime skill
  guidance and durable skill specs.
- Align the authoring guide and guide specs so direct edits wait after inviting
  feedback.
- Update docs/spec/runtime logs and release notes.

Deferred:

- New CLI support for direct authoring.
- A new public `/quality author` workflow.
- Changes to setup's existing final review gate, which already waits.
- Harness-specific native affordance behavior beyond the existing
  progressive-enhancement guidance.

## Affected artifacts

Derived by sweeping for direct model authoring, intent checkpoints,
`looks good`, review gates, feedback invitations, and agent-mediated UX guidance.

**Code**

- [x] No Go or TypeScript code impact; this is a docs/spec/skill guidance change.

**Durable specs** (substance in the [functional spec](0139-real-review-gates/spec.md))

- [x] `specs/skills/quality-skill/quality-skill.md` - require direct model
      authoring to acknowledge before long reads and wait after review
      checkpoints.
- [x] `specs/skills/quality-skill/guides/authoring.md` - align direct-edit
      guidance with real review gates.

**Durable docs / bundled skill runtime**

- [x] `docs/guides/agent-mediated-ux.md` - add review-gate doctrine and the
      long-read acknowledgement pattern.
- [x] `docs/log.md` - record the guide update.
- [x] `skills/quality/SKILL.md` - apply the acknowledgement and wait rules to
      direct model authoring.
- [x] `skills/quality/guides/authoring.md` - align direct-edit guidance with real
      review gates.
- [x] `skills/quality/log.md` - append runtime-history entry.
- [x] `skills/quality/guides/log.md` - append runtime-guide history entry.
- [x] `specs/log.md` - append durable-spec history entry.
- [x] `specs/skills/quality-skill/guides/log.md` - append durable-guide history
      entry.
- [x] `CHANGELOG.md` - note the `/quality` skill UX refinement.

No planned impact: `SPECIFICATION.md`, CLI command specs, setup workflow
mechanics, evaluation artifacts, generated schemas, or install docs.

## Status

`Done`. Implemented and archived after shared UX guidance, `/quality` runtime
skill guidance, durable skill specs, authoring guides, logs, release notes, and
Markdown formatting checks were updated.
