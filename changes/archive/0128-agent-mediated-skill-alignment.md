---
type: Change Case
title: Agent-mediated skill alignment
description: Close the remaining /quality skill gaps found in the agent-mediated UX audit: setup first-output ordering, setup feedback-log disclosure, recommendation follow-up framing, and read-only orientation output shape.
status: Done
tags: [skill, ux, agents, workflows]
timestamp: 2026-06-26T00:00:00Z
---

# Agent-mediated skill alignment

A **Change Case** to align the remaining `/quality` skill interaction contracts
with the agent-mediated UX guide after the audit of the bundled skill and durable
skill specs.

Detail lives in:

- [Functional spec](0128-agent-mediated-skill-alignment/spec.md) — the delta
  requirements and durable spec changes.
- [Design doc](0128-agent-mediated-skill-alignment/design.md) — the implementation
  approach and trade-offs.

## Motivation

The shared `/quality` interaction contract already carries the major
agent-mediated UX rules: status-first output, progressive enhancement, run
frames, decision briefs, explicit answer paths, progress beats, and closeouts.
The audit found a few remaining places where the durable skill and runtime
instructions still leave user-facing state ambiguous:

- setup's runtime template says the opening is first output, but the example
  places welcome prose before the run frame;
- setup promises that nothing is written until confirmation, while setup may
  write a workflow feedback log before `QUALITY.md` authoring;
- setup feedback-log timing differs between the workflow and CLI workflow
  conventions resource;
- recommendation follow-up is a user-visible post-evaluation workflow but has no
  opening frame;
- read-only orientation has routing rules but no standard status-first output
  shape.

These are interaction-contract gaps, not CLI or format changes. Closing them
keeps the user's current state and next action obvious across every `/quality`
entry path.

## Scope

Covered:

- Put setup's run frame first in the runtime opening example and keep the welcome
  and roadmap in the same first-output block after it.
- Replace setup's broad "nothing is written" line with an explicit distinction
  between model changes and the local workflow feedback log.
- Align setup feedback-log timing across runtime workflow, workflow spec, and
  CLI workflow conventions.
- Add a recommendation-follow-up opening frame before recommendation inspection,
  outcome selection, issue creation, local apply, or quality-log writes.
- Add a read-only orientation output shape that reports state, evidence limits,
  one recommended next action, and alternatives without mutating anything.
- Bring the bundled runtime skill files and durable skill specs/logs into sync.

Deferred / non-goals:

- No QUALITY.md format, model schema, rating, evaluation, or CLI command change.
- No Go code change; the CLI remains non-interactive.
- No change to mutation boundaries or confirmation requirements.
- No rewrite of historical Change Cases or append-only archived records.

## Affected artifacts

Derived by searching for run-frame, orientation, recommendation-follow-up,
feedback-log, and setup-opening guidance across `skills/quality/`,
`specs/skills/quality-skill/`, and `docs/guides/agent-mediated-ux.md`.

### Code

- [x] None — no `cmd/` or `internal/` behavior changes.

### Durable specs

- [x] `specs/skills/quality-skill/quality-skill.md` — add standard shapes for
      read-only orientation and non-public recommendation-follow-up framing.
- [x] `specs/skills/quality-skill/workflows/setup.md` — require the setup run
      frame first in the first-output block and clarify the feedback-log
      disclosure.
- [x] `specs/skills/quality-skill/recommendation-follow-up.md` and
      `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` — add the
      recommendation-follow-up opening frame.
- [x] `specs/skills/quality-skill/workflows/log.md`,
      `specs/skills/quality-skill/guides/log.md`, and `specs/log.md` — record the
      durable spec updates.

### Format spec

- [x] None — no `SPECIFICATION.md` impact.

### Durable docs and bundled skill

- [x] `skills/quality/SKILL.md` — add runtime orientation and
      recommendation-follow-up frame guidance.
- [x] `skills/quality/workflows/setup.md` — put the run frame first, clarify the
      setup write boundary, and keep feedback-log timing consistent.
- [x] `skills/quality/guides/recommendation-follow-up.md` — add the opening
      frame for follow-up.
- [x] `skills/quality/resources/cli-workflow-conventions.md` — align setup
      feedback-log timing with the setup workflow.
- [x] `skills/quality/log.md`, `skills/quality/workflows/log.md`,
      `skills/quality/guides/log.md`, and `skills/quality/resources/log.md` —
      record runtime skill updates.
- [x] `CHANGELOG.md` — add an Unreleased `/quality Skill` note when the repo's
      changelog has a current unreleased section.

## Status

`Done`. Functional spec and design doc are written; implementation is complete
across the durable `/quality` skill specs and bundled runtime skill files.
Verified with `mise run check`; archived with the parent and child folder under
`changes/archive/`.
