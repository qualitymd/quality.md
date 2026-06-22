---
type: Functional Specification
title: Remove improve mode - functional spec
description: Requirements for removing /quality improve while preserving recommendation follow-up.
tags: [skill, quality]
timestamp: 2026-06-22T00:00:00Z
---

# Remove improve mode - functional spec

Companion to the [Remove improve mode](../0054-remove-improve-mode.md) change
case. This spec states what the skill-surface simplification must do.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Scope

This change removes `improve` as a named `/quality` mode and replaces it with a
recommendation follow-up workflow available after evaluation or history review.
It changes the bundled skill, durable skill specs, user-facing docs, examples,
and bundle navigation.

Non-goals:

- changing the QUALITY.md format;
- changing `qualitymd` CLI commands;
- changing evaluation run, recommendation, report, or quality-log schemas;
- adding issue-tracker-specific integration behavior beyond preparing issue-ready
  content and using available agent connector tools when the user explicitly
  confirms external creation.

## Requirements

### Mode surface

The `/quality` skill **MUST** remove `improve` from mode parsing, invocation
variants, mode dispatch, and user-facing examples.

The supported named modes **MUST** be `wizard`, `setup`, `evaluate`, and
`update`.

Requests to "improve", "apply", "act on", or "handoff" a recommendation
**MUST** route to recommendation follow-up rather than to a separate mode.

### Evaluation remains advisory

`evaluate` **MUST** continue to emit recommendation records when evaluation
finds gaps.

`evaluate` **MUST NOT** edit evaluated source files, edit `QUALITY.md`, write the
quality log, create external issues, or apply recommendations.

### Recommendation follow-up

Recommendation follow-up **MUST** offer only two explicit productive outcomes:

1. apply a confirmed recommendation option now;
2. hand off the recommendation to an issue tracker.

The workflow **MUST NOT** present `defer`, `skip`, or `keep open` as formal
follow-up options.

If the user does not choose apply or issue-tracker handoff, the workflow
**MUST** stop without mutating evaluated source, `QUALITY.md`, the quality log,
or external systems.

### Apply now

Before applying a recommendation, the skill **MUST** present a decision brief
that identifies the recommendation, selected option, affected artifact class,
evidence or reason, risk when relevant, done criterion, and verification path.

The skill **MUST NOT** edit evaluated source files, edit `QUALITY.md`, or write
the quality log until the user explicitly confirms the recommendation option and
mutation surface.

After applying a confirmed option, the skill **SHOULD** verify the done criterion
with the narrowest useful evidence. When the recommendation's done criterion is
rating-bound or otherwise depends on the QUALITY.md model, the skill **SHOULD**
run a scoped re-evaluation and report the before/after delta.

The result report **MUST** state the recommendation, applied option, changed
artifacts, verification performed, rating movement when known, and remaining
limits.

### Issue-tracker handoff

Issue-tracker handoff **MUST** produce issue-ready content that includes the
recommendation ID and title, source evaluation run, affected area/factor or
requirement, current rating when known, target or done criterion, evidence
summary with locators, suggested implementation option, verification path, and
links or paths to the generated report and recommendation artifact.

Creating an external issue **MUST** require explicit user confirmation and
available issue-tracker tooling.

Issue-tracker handoff **MUST NOT** mutate evaluated source, `QUALITY.md`, or the
quality log.

When issue-tracker tooling is unavailable or the user has not confirmed external
creation, the skill **MUST** stop after producing issue-ready text.

### Quality log

The quality log **MUST** remain limited to meaningful `QUALITY.md` model changes.

`setup` **MUST** continue to seed an inaugural quality-log entry for meaningful
model creation or first population.

A confirmed recommendation-apply or model-authoring workflow **MUST** write one
quality-log entry when it makes a meaningful `QUALITY.md` model change.

Issue-tracker handoff alone **MUST NOT** write the quality log.

### Wizard routing

`wizard` **MUST** remain read-only.

When active recommendations exist, `wizard` **SHOULD** route to recommendation
review, apply, or issue-tracker handoff rather than to `improve`.

When model history needs reconciliation, `wizard` **SHOULD** route to confirmed
model-authoring or recommendation-follow-up work rather than to `improve`.

### Documentation and search cleanup

Runtime skill docs, durable specs, and user docs **MUST NOT** advertise
`/quality improve` as an invocation.

Repository search for `/quality improve`, `improve mode`, and
`modes/improve.md` **MUST** return only historical/archive references or
intentional migration notes after implementation.

Conceptual uses of "improve" and "improvement" **MAY** remain when they refer to
the quality-loop phase or ordinary English usage.

## Durable spec changes

### To add

- `specs/skills/quality-skill/recommendation-follow-up.md` - non-mode contract
  for reviewing, applying, or handing off evaluation recommendations.
- `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` - 1:1
  artifact spec for the runtime recommendation follow-up guide.

### To modify

- `specs/skills/quality-skill/quality-skill.md` - remove `improve` as a mode and
  add shared recommendation-follow-up behavior.
- `specs/skills/quality-skill/modes/evaluate.md` - keep recommendations as
  evaluation outputs and point to follow-up behavior.
- `specs/skills/quality-skill/modes/wizard.md` - route active recommendations to
  review, apply, or issue handoff.
- `specs/skills/quality-skill/modes/index.md` - remove `improve` mode navigation
  or replace it with non-mode follow-up navigation if a child spec remains.
- `specs/skills/quality-skill/index.md` - update skill-spec navigation if child
  specs change.
- `specs/skills/quality-skill/guides/index.md` - list the recommendation
  follow-up guide contract.
- `specs/skills/quality-skill/guides/log.md` - record the new guide contract.
- `specs/log.md` and relevant child `log.md` files - record the durable spec
  updates.

### To rename

None.

### To delete

- `specs/skills/quality-skill/modes/improve.md` - remove the improve mode
  contract.
