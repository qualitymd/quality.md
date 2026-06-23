---
type: Change Case
title: Remove wizard mode
description: Remove wizard from the /quality skill's public contract while preserving safe read-only orientation for ambiguous requests.
status: In-Review
tags: [skill, ux, contract]
timestamp: 2026-06-23T00:00:00Z
---

# Remove wizard mode

A **Change Case** simplifying the `/quality` skill contract by removing
`wizard` as a public mode and documented invocation. The detail lives in its
[functional spec](0062-remove-wizard-mode/spec.md) and
[design doc](0062-remove-wizard-mode/design.md).

## Motivation

`wizard` currently names a useful read-only lifecycle routing behavior, but it
does not describe a clear user task alongside `setup`, `evaluate`, and `update`.
It also competes with a likely future use of "setup wizard" to describe guided
setup behavior. The public surface should be simpler: users choose setup,
evaluation, update, or recommendation follow-up, while unclear input remains
safe and read-only.

## Scope

Covered:

- Remove `wizard` as a documented `/quality` mode and invocation.
- Keep bare `/quality` valid as a read-only orientation entrypoint that may
  recommend only public workflows.
- Keep ambiguous requests safe: they must not evaluate, mutate files, update
  tooling, or write artifacts without the user choosing a public workflow.
- Remove `status`, `next`, `review model`, and `review history` from the public
  invocation contract rather than promoting them as replacement modes.
- Preserve recommendation follow-up as a non-mode workflow.

Deferred / non-goals:

- No design of a future setup wizard.
- No new public `status`, `next`, `review`, or `orient` mode.
- No change to the QUALITY.md format specification.
- No planned `qualitymd` CLI command or Go code change unless implementation
  analysis finds hidden coupling.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0062-remove-wizard-mode/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [x] No planned `qualitymd` CLI or Go code changes; implementation analysis
      should confirm no tests or generated examples encode `wizard` as public
      contract.

Specs:

- [x] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      - remove `wizard` from the public mode model, invocation examples, mode
      dispatch, and default routing language.
- [x] `specs/skills/quality-skill/modes/wizard.md` - delete the public wizard
      mode component spec.
- [x] [`specs/skills/quality-skill/modes/index.md`](../specs/skills/quality-skill/modes/index.md)
      and [`specs/skills/quality-skill/index.md`](../specs/skills/quality-skill/index.md)
      - remove public wizard listings and broken links.
- [x] [`specs/skills/quality-skill/modes/setup.md`](../specs/skills/quality-skill/modes/setup.md)
      and [`specs/skills/quality-skill/guides/getting-started-md.md`](../specs/skills/quality-skill/guides/getting-started-md.md)
      - stop routing setup completion to wizard.
- [x] [`specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md`](../specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md)
      and [`specs/skills/quality-skill/guides/index.md`](../specs/skills/quality-skill/guides/index.md)
      - describe the checklist as internal orientation/model-review support, not
      wizard-specific behavior.
- [x] [`specs/skills/quality-skill/quality-log.md`](../specs/skills/quality-skill/quality-log.md)
      and [`specs/skills/quality-skill/recommendation-follow-up.md`](../specs/skills/quality-skill/recommendation-follow-up.md)
      - remove wizard-specific routing/reconciliation wording.
- [x] OKF logs under [`specs/`](../specs/log.md) - record durable spec updates
      when they land.

Runtime skill and docs:

- [x] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) - remove `wizard`
      from description, hard rules, argument parsing, invocation variants, mode
      dispatch, and run-frame exceptions.
- [x] `skills/quality/modes/wizard.md` - delete the runtime wizard mode file.
- [x] [`skills/quality/modes/setup.md`](../skills/quality/modes/setup.md),
      [`skills/quality/guides/getting-started.md`](../skills/quality/guides/getting-started.md),
      and [`skills/quality/guides/top-10-quality-md-checks.md`](../skills/quality/guides/top-10-quality-md-checks.md)
      - remove wizard-specific handoff and checklist language.
- [x] [`README.md`](../README.md) - remove `/quality wizard` usage.
- [x] [`docs/guides/use-quality-skill.md`](../docs/guides/use-quality-skill.md)
      and [`docs/guides/index.md`](../docs/guides/index.md) - remove wizard from
      public skill guidance.
- [x] [`install.md`](../install.md) - remove `/quality wizard` from bootstrap
      instructions.
- [x] [`CHANGELOG.md`](../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

## Children

- [Functional spec](0062-remove-wizard-mode/spec.md) - what removing wizard from
  the public contract must preserve, remove, and leave deferred.
- [Design doc](0062-remove-wizard-mode/design.md) - how the skill removes wizard
  while preserving safe read-only orientation.

## Status

`In-Review`. Implementation is complete across runtime skill guidance, durable
specs, public docs, indexes, logs, and changelog. Verified `mise run check`.
