---
type: Change Case
title: Remove improve mode
description: Simplify the /quality skill by removing the separate improve mode while preserving recommendation follow-up.
status: Done
tags: [skill, quality]
timestamp: 2026-06-22T00:00:00Z
---

# Remove improve mode

This case removes `/quality improve` as a named skill mode and keeps
recommendation follow-up as a post-evaluation workflow.

- [Functional spec](0054-remove-improve-mode/spec.md) - what the mode removal
  and recommendation follow-up workflow must do.

## Motivation

The README presents the primary user path as setup, evaluation, and then acting
on or refining from evaluation results. It describes "Improve" as a quality-loop
phase, not as a separate command. The current `improve` mode is thin: it first
does an evaluation, then adds only a confirmed apply step and re-evaluation.

Keeping that as a public mode complicates the skill surface without supporting a
distinct core scenario. The useful behavior is recommendation follow-up:
reviewing evaluation recommendations, either applying a confirmed option now or
handing the recommendation to an issue tracker.

## Scope

In scope:

- remove `/quality improve` from the skill's public mode surface;
- preserve recommendations as outputs of `/quality evaluate`;
- define recommendation follow-up with two explicit outcomes: apply now or
  hand off to an issue tracker;
- preserve explicit confirmation before local mutation or external issue
  creation;
- keep quality-log entries for meaningful confirmed `QUALITY.md` model changes;
- update runtime skill instructions, durable skill specs, user docs, examples,
  and bundle navigation to match the simplified surface.

Deferred:

- adding a deterministic `qualitymd` command for recommendation application or
  issue creation;
- changing evaluation record, recommendation record, report, or quality-log
  file schemas;
- integrating with a specific issue tracker beyond the agent's existing
  connector/tooling capabilities;
- migrating historical evaluation runs or archived change cases that mention
  `improve`.

## Affected artifacts

### Code

None expected.

### Format spec

None.

### Durable specs

- `specs/skills/quality-skill/quality-skill.md` - remove `improve` as a mode
  and replace improve-specific shared contracts with recommendation follow-up
  contracts.
- `specs/skills/quality-skill/recommendation-follow-up.md` - add a non-mode
  contract for applying or handing off evaluation recommendations.
- `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` - add the
  1:1 artifact spec for the runtime follow-up guide.
- `specs/skills/quality-skill/modes/improve.md` - delete the improve mode
  contract.
- `specs/skills/quality-skill/modes/wizard.md` - route active recommendations to
  review, apply, or issue handoff instead of `improve`.
- `specs/skills/quality-skill/modes/evaluate.md` - keep recommendations as
  evaluation outputs and point to follow-up behavior.
- `specs/skills/quality-skill/modes/index.md` and
  `specs/skills/quality-skill/index.md` - update navigation.
- `specs/skills/quality-skill/guides/index.md`,
  `specs/skills/quality-skill/guides/log.md`, `specs/log.md`, and any affected
  child `log.md` files - record durable spec changes.

### Durable docs

- `README.md` - keep the public usage path centered on setup, wizard, and
  evaluate, with improve described only as a quality-loop phase when needed.
- `docs/guides/use-quality-skill.md` - remove improvement-mode guidance and
  describe recommendation follow-up and issue-tracker handoff.
- Other docs found by repository search that advertise `/quality improve` or
  "improve mode" as a public workflow.

### Bundled skill

- `skills/quality/SKILL.md` - remove `improve` from argument parsing,
  invocation variants, hard rules, mode dispatch, and log ownership phrasing.
- `skills/quality/guides/recommendation-follow-up.md` - add non-mode runtime
  guidance if the shared prompt becomes too large.
- `skills/quality/modes/improve.md` - delete the improve mode procedure.
- `skills/quality/modes/evaluate.md` - preserve recommendation output behavior
  and point to follow-up.
- `skills/quality/modes/wizard.md` - route active recommendations to review,
  apply, or issue handoff.
- `skills/quality/guides/authoring.md` and
  `skills/quality/guides/top-10-quality-md-checks.md` - update stale routing and
  quality-loop language found by search.

### Install, scaffold, and packaging

None expected.

## Status

`Done`: implemented and archived. No design doc was required because the work
was a mechanical skill/spec/doc surface change.
