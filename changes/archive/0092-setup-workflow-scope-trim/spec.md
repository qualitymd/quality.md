---
type: Functional Specification
title: Setup workflow scope trim - functional spec
description: Requirements for trimming /quality setup so it stops asking about future recommendation handling, review cadence, recurring automation, and evaluation-readiness/maturity status while preserving the context needed to author a high-quality initial QUALITY.md.
tags: [skill, setup, ux]
timestamp: 2026-06-25T00:00:00Z
---

# Setup workflow scope trim - functional spec

Companion to [Setup workflow scope trim](../0092-setup-workflow-scope-trim.md).
This spec states what the setup workflow must stop asking, stop reporting, and
continue preserving. The format itself is governed by
[`SPECIFICATION.md`](../../../SPECIFICATION.md). This case changes bundled skill
guidance and durable skill spec mirrors only; it adds no normative format rule
and no CLI behavior.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Setup is the first-model workflow. Its highest-value work is to make the model
boundary, domain, users, maintainers, stakeholders, risks, unknowns, rating scale,
Factors, Areas, and Requirements visible enough that an agent can author a useful
initial `QUALITY.md` without pretending repo-invisible facts are known. Future
recommendation handling, recurring review, automation posture, and
evaluation-readiness labels do not improve that first model enough to justify the
extra setup surface. They are better handled when an evaluation or follow-up task
actually exists.

## Scope

Covered:

- Remove recommendation-handling and work-handoff destination questions from
  setup discovery.
- Remove review cadence, recurring review, and automation posture from setup
  discovery and setup closeout routing.
- Remove setup's maturity/evaluation-readiness classification and closeout label.
- Keep lint validation, important model gaps, and narrow next-step guidance.
- Keep the existing pedagogical setup questions and human context checkpoint
  unless they only serve the removed future workflow preferences.
- Keep setup workflow feedback logging, with a revised outcome contract.

Deferred / non-goals:

- No change to recommendation follow-up's two productive outcomes: apply a
  confirmed recommendation option now, or hand off the recommendation to an issue
  tracker.
- No removal of the core setup questions for turn-count optimization.
- No removal of setup teaching copy or the review gate.
- No removal of workflow feedback logs.
- No `QUALITY.md` schema, rating, roll-up, evaluation record, CLI, or Go change.

## Requirements

### Keep setup focused on first-model context

`/quality setup` **MUST** collect and preserve only context that helps author,
validate, or immediately refine the first useful `QUALITY.md`: root boundary,
domain, lifecycle, risk tolerance, rating scale calibration, primary users and
outcomes, maintainers/collaborators, other stakeholders, missing or
non-agent-accessible context, and the evidence needed to derive Factors, Areas,
Requirements, unknowns, and open questions.

Setup **MUST NOT** ask the user to choose future recommendation-handling policy,
work-handoff destination, issue tracker, review cadence, recurring review
posture, CI or release-gating posture, calendar/scheduler posture, Codex
automation, Claude Code routine, or other automation posture.

> Rationale: these choices do not materially improve the first model. Asking them
> during setup makes follow-on workflow decisions feel like setup obligations and
> creates expectations that later workflows do not currently consume. - 0092

### Remove recommendation handling from setup

Setup **MUST NOT** ask where future evaluation recommendations should usually go
and **MUST NOT** offer setup-time options such as leaving recommendations in the
evaluation report, GitHub Issues, Linear/Jira, or maintainer-decides-each-time.

Recommendation handling **MUST** remain owned by recommendation follow-up, after
an evaluation has produced active recommendation records or the user has asked to
act on a concrete recommendation. Setup documentation **MAY** mention, as a
boundary, that setup creates no issues and configures no integrations.

### Remove recurring review and automation posture from setup

Setup **MUST NOT** ask for review cadence, recurring quality-loop posture, or
automation setup preferences. Setup closeout **MUST NOT** present `set up
recurring review`, `set up recommendation handoff`, or equivalent follow-on
configuration as setup's next-step choices.

Setup closeout **MAY** state what setup did not do, including no evaluation, no
quality-log entry, no external issues, and no automations, when that boundary
prevents user confusion.

> Rationale: recurring review and automation are valid later workflows, but
> surfacing them during setup distracts from authoring the initial model and can
> make setup feel like it is configuring operations. - 0092

### Remove setup maturity and evaluation-readiness status

Setup **MUST NOT** classify the resulting model as `starter`, `immature`, or
`evaluation-ready`, and **MUST NOT** report `Maturity:` or
`Evaluation readiness:` in setup closeout.

Setup **MUST** run `qualitymd lint [path]` and report whether validation passed
or failed. When lint passes, setup **MUST** report important model gaps that
materially affect first-model usefulness, such as missing stakeholder context,
thin Needs/Risks, missing germane constituents, missing Agent Harnessability
coverage for an agent-collaborated composite root, or vague Requirements. This
gap report **MUST** be framed as model-improvement guidance, not as a readiness
or maturity classification.

Setup **SHOULD** recommend one immediate next step from this reduced set:
continue iterating on `QUALITY.md`, run `/quality evaluate`, or stop here. The
recommendation **MUST** be based on lint status and important gaps, not on a named
evaluation-readiness level.

> Rationale: lint already answers validity. Important gaps answer what to improve
> before or after evaluation. A named `evaluation-ready` status adds a second
> judgment layer and can over-promise that evaluation will produce a useful signal.
>
> - 0092

### Preserve pedagogical discovery and workflow feedback

Setup **MUST** keep the core discovery questions and human context checkpoint
that teach and capture the dimensions needed for a useful first model. This case
**MUST NOT** remove questions merely because the agent has a high-confidence
default or because fewer turns would be convenient.

Setup **MUST** keep the workflow feedback log under `.quality/logs/`. Its outcome
metadata **MUST NOT** use `starter`, `immature`, or `evaluation-ready`. The
feedback-log contract **MUST** instead use setup-result values that describe the
workflow outcome without implying evaluation readiness, such as completed,
completed-with-important-gaps, lint-failed, or interrupted.

> Rationale: the retained questions are part of the setup pedagogy and protect
> against invented context. The feedback log improves the workflow itself, so it
> remains valuable even though it is not model content. - 0092

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/workflows/setup.md` - remove follow-on workflow
  discovery, remove recommendation handoff and recurring review next-step routing,
  and replace maturity/evaluation-ready closeout with lint plus important gaps
  (per the setup-focus, recommendation-handling, recurring-review, and
  evaluation-readiness requirements above).
- `specs/skills/quality-skill/workflows/setup/feedback-log.md` - replace
  maturity outcome values with setup-result values that do not imply evaluation
  readiness (per the workflow-feedback requirement above).
- `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` - remove the
  setup-close dependency on the maturity/evaluation-ready checklist; retain or
  reframe read-only orientation checks only where they remain useful outside setup
  (per the evaluation-readiness requirement above).
- `specs/skills/quality-skill/guides/getting-started-md.md` - remove review
  posture, handoff posture, recurring review, and recommendation handoff as
  first-model continuation concerns (per the setup-focus, recommendation-handling,
  and recurring-review requirements above).
- `specs/skills/quality-skill/guides/index.md` - align guide descriptions with
  the new getting-started and Top 10 framing (per the evaluation-readiness and
  workflow-feedback requirements above).
- `specs/skills/quality-skill/quality-skill.md` - align the public setup contract
  if it still says setup classifies maturity or offers follow-on handoff/recurring
  review setup (per the requirements above).
- `specs/skills/quality-skill/guides/log.md` and
  `specs/skills/quality-skill/workflows/log.md` - record the durable spec mirror
  changes.

### To rename

None

### To delete

None
