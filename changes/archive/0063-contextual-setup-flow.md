---
type: Change Case
title: Contextual setup flow
description: Rework /quality setup into a short context-informed discovery flow that writes only QUALITY.md and routes next steps.
status: Done
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Contextual setup flow

A **Change Case** reworking `/quality setup` from skeleton-plus-population into
a short context-informed discovery flow that teaches through decisions, writes a
useful first `QUALITY.md`, validates it, and offers explicit next-step choices.
The detail lives in its [functional spec](0063-contextual-setup-flow/spec.md)
and [design doc](0063-contextual-setup-flow/design.md).

## Motivation

Setup is the user's first encounter with QUALITY.md as a practical quality bar.
The current flow correctly creates a valid skeleton and begins guided
population, but it does not make the most important setup judgments explicit:
project lifecycle, risk tolerance, modeling rigor, agent-heavy collaboration,
stakeholder needs, missing context, work-management posture, and recurring
quality-review intent.

Those judgments should be elicited quickly, with context-informed defaults, then
preserved in `QUALITY.md` so future agents and maintainers can understand why
the first model has its shape. Setup should remain scoped: it creates or updates
the model only. Evaluation, automation, issue creation, and recommendation
handoff are follow-on workflows, not setup side effects.

## Scope

Covered:

- Analyze available repository context before asking setup questions.
- Ask a small set of discovery questions with recommended defaults and
  confidence signals.
- Treat current directory as the default root area convention.
- Assume agent-heavy development and ask which human collaborators or
  stakeholders also need alignment.
- Distinguish lifecycle, risk tolerance, and modeling rigor.
- Seed missing-context prompts from the agent's analysis.
- Capture work-management and recurring quality-review posture without creating
  external issues, automations, CI jobs, release gates, or scheduler config.
- Write or update only `QUALITY.md`, then validate and inspect model readiness.
- End with explicit next-step options rather than running evaluation or
  configuring integrations.

Deferred / non-goals:

- No change to the QUALITY.md format specification.
- No `qualitymd` CLI setup wizard or interactive CLI workflow.
- No automatic evaluation during setup.
- No GitHub Issues, Linear, Jira, or other external issue creation during setup.
- No CI, release-gating, scheduled automation, Codex automation, or Claude Code
  routine creation during setup.
- No recommendation to make CI or release gating the default quality loop.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0063-contextual-setup-flow/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [x] No planned `qualitymd` CLI or Go code changes; implementation analysis
      should confirm setup remains skill-driven and CLI-owned only for
      scaffolding, linting, status, and format grounding.

Specs:

- [x] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      - update the shared setup summary, mutation surfaces, run-frame
      expectations, and quality-log writer contract.
- [x] [`specs/skills/quality-skill/modes/setup.md`](../specs/skills/quality-skill/modes/setup.md)
      - replace current guided-population procedure with contextual discovery,
      model write, lint/readiness inspection, and next-step routing.
- [x] [`specs/skills/quality-skill/quality-log.md`](../specs/skills/quality-skill/quality-log.md)
      - remove setup as a quality-log writer because setup should mutate only
      `QUALITY.md`.
- [x] [`specs/skills/quality-skill/guides/getting-started-md.md`](../specs/skills/quality-skill/guides/getting-started-md.md)
      - align first-run guidance with the new setup flow and post-setup
      iteration path.
- [x] [`specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md`](../specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md)
      - confirm whether the already-expanded checklist needs further setup
      readiness wording.
- [x] OKF logs under [`specs/`](../specs/log.md) - record durable spec updates
      when they land.

Runtime skill and docs:

- [x] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) - update setup
      mutation surfaces, quality-log hard rule, and setup interaction contract.
- [x] [`skills/quality/modes/setup.md`](../skills/quality/modes/setup.md) -
      implement the new runtime setup procedure.
- [x] [`skills/quality/guides/getting-started.md`](../skills/quality/guides/getting-started.md)
      - align the first-run playbook with the setup handoff.
- [x] [`skills/quality/guides/top-10-quality-md-checks.md`](../skills/quality/guides/top-10-quality-md-checks.md)
      - confirm no further runtime checklist update is needed beyond the
      durable setup-readiness concerns already added.
- [x] [`README.md`](../README.md), [`install.md`](../install.md), and
      [`docs/guides/use-quality-skill.md`](../docs/guides/use-quality-skill.md)
      - update public setup guidance and quality-loop wording.
- [x] [`CHANGELOG.md`](../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

## Children

- [Functional spec](0063-contextual-setup-flow/spec.md) - what contextual setup
  must ask, write, avoid, validate, and offer next.
- [Design doc](0063-contextual-setup-flow/design.md) - how setup analyzes
  context, elicits decisions, writes `QUALITY.md`, validates readiness, and
  routes next steps.

## Status

`Done`. Implementation is complete across runtime skill guidance, durable
skill specs, public docs, logs, and changelog. Verified `mise run check`.
Landed and archived.
