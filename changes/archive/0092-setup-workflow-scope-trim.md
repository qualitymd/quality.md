---
type: Change Case
title: Setup workflow scope trim
description: Trim setup so it stops asking about future recommendation handling, review cadence, recurring automation, and evaluation-readiness/maturity status, keeping the workflow focused on producing a high-quality initial QUALITY.md.
status: Done
tags: [skill, setup, ux]
timestamp: 2026-06-25T00:00:00Z
---

# Setup workflow scope trim

A **Change Case** to tighten `/quality setup` around its primary job: producing a
high-quality initial `QUALITY.md`. Setup should collect the context needed to
author the model, validate it, name important model gaps, and stop. It should not
ask users to choose future recommendation-handling policy, recurring review
posture, automation posture, or an evaluation-readiness classification during the
first-model workflow.

Detail lives in:

- [Functional spec](0092-setup-workflow-scope-trim/spec.md) - what setup must
  stop asking, stop reporting, and continue preserving.
- [Design doc](0092-setup-workflow-scope-trim/design.md) - how the runtime
  workflow and durable mirrors are trimmed while preserving pedagogical discovery
  and feedback logging.

## Motivation

The current setup workflow has accumulated follow-on workflow questions and
closeout labels that do not directly improve the first `QUALITY.md`:

- an optional work-handoff question asks where future evaluation recommendations
  should go, but recommendation follow-up later owns apply-vs-handoff decisions
  and still requires explicit confirmation;
- review cadence and recurring automation are already follow-on work, but setup
  still carries posture and closeout routing language that can make those concerns
  feel like setup inputs;
- the setup closeout classifies model maturity as `starter`, `immature`, or
  `evaluation-ready`, which adds a second judgment layer after lint and can read
  like a promise that evaluation is now the right or valuable next action.

Those concerns are real, but they are not germane to producing the initial model.
Setup should stay focused on root boundary, domain, lifecycle, risk tolerance,
stakeholders, missing context, rating scale calibration, model synthesis, lint,
and important gaps.

## Scope

Covered:

- Remove setup discovery for future recommendation handling, work-handoff
  destination, review cadence, recurring review posture, and automation posture.
- Remove setup closeout routing options that invite setting up recurring review
  or recommendation handoff as part of setup.
- Remove setup's maturity/evaluation-readiness classification from the runtime
  workflow and closeout.
- Keep lint validation and important model gaps in the setup closeout.
- Keep the core discovery questions and checkpoint that shape the first model.
- Keep the workflow feedback log, but update its outcome contract so it no longer
  records `starter`, `immature`, or `evaluation-ready`.
- Align the bundled skill guidance and durable skill spec mirrors.

Deferred / non-goals:

- Do not remove or reduce the core setup questions only to optimize turn count.
  The retained questions serve a pedagogical purpose and help prevent the agent
  from inventing model context.
- Do not default the rating-scale question away in this case. That may be worth
  revisiting later, but it is a separate trade-off.
- Do not remove the workflow feedback log. It is important for improving the
  setup process even though it is not model content.
- Do not change recommendation follow-up behavior beyond making clear it owns
  recommendation apply and handoff decisions when recommendations actually exist.
- Do not change the `QUALITY.md` format, rating semantics, roll-up behavior,
  evaluation records, CLI behavior, or Go code.

## Affected artifacts

Derived by analysis: searched setup and durable skill specs for
`work-handoff`, `recommendation handoff`, `review posture`, `recurring review`,
`maturity`, and `evaluation-ready`, then separated live guidance from historical
archive/log entries. Grouped by kind; empty kinds are deliberate. Checkboxes are
reconciled before In-Review.

### Code

None expected - bundled skill guidance and durable spec mirrors only.

### Format spec (`SPECIFICATION.md`)

None - no QUALITY.md schema, rating-scale, roll-up, or evaluation-semantic
change.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0092-setup-workflow-scope-trim/spec.md#durable-spec-changes)
section is authoritative. In summary:

- [x] `specs/skills/quality-skill/workflows/setup.md` - remove follow-on
      workflow discovery and maturity/evaluation-ready closeout requirements;
      require lint plus important gaps instead.
- [x] `specs/skills/quality-skill/workflows/setup/feedback-log.md` - replace the
      maturity outcome contract with setup-result values that do not imply
      evaluation readiness.
- [x] `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md` -
      remove setup-close dependence on the maturity/evaluation-ready checklist;
      keep any read-only orientation contract only if still useful outside setup.
- [x] `specs/skills/quality-skill/guides/getting-started-md.md` - remove review
      posture, handoff posture, recurring review, and recommendation handoff as
      first-model/setup continuation concerns.
- [x] `specs/skills/quality-skill/guides/index.md` - align guide descriptions
      with the new getting-started and Top 10 framing.
- [x] `specs/skills/quality-skill/quality-skill.md` - align the public setup
      contract if it still says setup classifies model maturity or offers
      follow-on handoff/recurring-review setup.
- [x] `specs/skills/quality-skill/guides/log.md` and
      `specs/skills/quality-skill/workflows/log.md` - record the durable spec
      mirror updates if those files are edited during implementation.

### Durable docs

- [x] `README.md` - review setup/getting-started copy for any claim that setup
      captures recommendation handoff, recurring review posture, or evaluation
      readiness; no edit needed because live README setup copy does not imply
      those concerns are part of setup.
- [x] `docs/guides/agent-mediated-ux.md` - update the setup closeout example if
      it still uses `Maturity: immature`.
- [x] `docs/log.md` - record doc-guide edits if `docs/` is touched.
- [x] `CHANGELOG.md` - add a user-facing note when the change is implemented.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/workflows/setup.md` - primary runtime workflow trim.
- [x] `skills/quality/guides/top-10-quality-md-checks.md` - remove setup-close
      dependence on maturity/evaluation-ready classification, preserving only
      read-only orientation checks that remain useful.
- [x] `skills/quality/guides/getting-started.md` - remove review/handoff posture
      from first-model continuation guidance.
- [x] `skills/quality/guides/index.md` - align guide descriptions with the new
      getting-started and Top 10 framing.
- [x] `skills/quality/SKILL.md` - align setup workflow summary and any references
      to readiness/maturity routing if needed.
- [x] `skills/quality/workflows/log.md` and `skills/quality/guides/log.md` -
      record the bundled skill guide/workflow edits.

### Install / scaffold

None expected - the scaffolded `QUALITY.md` content and install flow are not part
of this scope.

## Children

- [Functional spec](0092-setup-workflow-scope-trim/spec.md) - what setup must
  stop asking, stop reporting, and continue preserving.
- [Design doc](0092-setup-workflow-scope-trim/design.md) - how the runtime
  workflow and durable mirrors are trimmed while preserving pedagogical discovery
  and feedback logging.

## Status

`Done`. Implemented and archived. Setup no longer asks about future
recommendation handling, handoff destination, review cadence, recurring review,
or automation preferences; setup closeout reports validation plus important model
gaps instead of maturity/evaluation-readiness labels; setup feedback logs record
workflow outcomes. Verified with `mise run check`. No CLI, Go, format-schema,
rating, roll-up, or evaluation-record behavior changed.
