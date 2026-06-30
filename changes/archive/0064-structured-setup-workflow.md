---
type: Change Case
title: Structured setup workflow
description: Turn /quality setup guidance into an explicit workflow with concrete discovery questions and prompt framing.
status: Done
tags: [skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Structured setup workflow

A **Change Case** refining `/quality setup` from a spec-like runtime mode into a
concrete setup workflow with explicit discovery questions, prompt framing, and
operator steps. The detail lives in its
[functional spec](0064-structured-setup-workflow/spec.md) and
[design doc](0064-structured-setup-workflow/design.md).

## Motivation

The contextual setup flow established the right concepts: context inspection,
confidence-labeled defaults, lifecycle, risk tolerance, modeling rigor,
collaboration context, needs, missing context, and follow-on posture. The
runtime setup guidance still reads more like a behavioral spec than an
actionable workflow an agent can reliably execute.

Setup should be a structured interaction: inspect the repo, present a brief,
ask a stable set of concrete questions, synthesize the model, validate it, and
route next steps. The workflow should teach through the questions themselves
without asking users to design factors or areas cold. It should also use
"workflow" language for the user-facing setup process while reserving "mode" for
internal dispatch where useful.

## Scope

Covered:

- Rewrite runtime setup guidance as a step-by-step setup workflow.
- Specify the concrete setup questions, option sets, recommended defaults, and
  confidence framing.
- Define how questions may be grouped into one compact prompt or a short
  sequence without losing required content.
- Make root area, domain, lifecycle, risk tolerance, modeling rigor,
  collaboration, stakeholder needs, missing context, and review posture explicit
  setup inputs.
- Clarify terminology: `/quality setup` starts a setup workflow; "mode" remains
  internal dispatch vocabulary unless a later design chooses to rename file
  paths.
- Keep setup's mutation boundary unchanged: setup writes only `QUALITY.md`.

Deferred / non-goals:

- No change to the QUALITY.md format specification.
- No `qualitymd` CLI setup wizard or interactive CLI workflow.
- No automatic evaluation, issue creation, quality-log writing, or automation
  configuration during setup.
- No required rename of existing `modes/` file paths in this draft.
- No requirement that users choose factors, child areas, or Requirements
  directly.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0064-structured-setup-workflow/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [x] No planned `qualitymd` CLI or Go code changes; implementation analysis
      should confirm setup remains skill-driven and that no generated examples
      encode the old abstract setup wording.

Specs:

- [x] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md) - describe `/quality setup` as a setup workflow in public-facing
      summaries while preserving internal dispatch semantics.
- [x] [`specs/skills/quality-skill/modes/setup.md`](../specs/skills/quality-skill/modes/setup.md) - replace abstract setup-mode requirements with a concrete workflow
      contract, fixed discovery question set, prompt templates, and branching
      rules.
- [x] [`specs/skills/quality-skill/guides/getting-started-md.md`](../specs/skills/quality-skill/guides/getting-started-md.md) - align post-setup iteration guidance with the structured setup
      assumptions recorded in `QUALITY.md`.
- [x] OKF logs under [`specs/`](../specs/log.md) - record durable spec updates
      when they land.

Runtime skill and docs:

- [x] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) - use setup
      workflow wording where user-facing, and keep mode terminology only for
      dispatch/routing.
- [x] [`skills/quality/modes/setup.md`](../skills/quality/modes/setup.md) -
      rewrite as an operator playbook with concrete prompts, required question
      framing, synthesis steps, validation, and completion output.
- [x] [`skills/quality/guides/getting-started.md`](../skills/quality/guides/getting-started.md) - ensure first-run follow-on guidance refers to the structured setup
      assumptions consistently.
- [x] [`README.md`](../README.md) and
      [`docs/guides/use-quality-skill.md`](../docs/guides/use-quality-skill.md) - update setup wording if implementation changes public phrasing.
- [x] [`npm/quality.md/README.md`](../npm/quality.md/README.md) - keep the
      packaged README copy aligned with root README setup wording.
- [x] [`CHANGELOG.md`](../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

## Children

- [Functional spec](0064-structured-setup-workflow/spec.md) - what the setup
  workflow must ask, how prompts must be framed, and which terminology changes
  must land.
- [Design doc](0064-structured-setup-workflow/design.md) - how setup guidance
  becomes an operator workflow with a setup brief and concrete prompt.

## Status

`Done`. Implementation is complete across runtime skill guidance, durable
skill specs, public docs, logs, and changelog. Verification is recorded in the
change log. Landed and archived.
