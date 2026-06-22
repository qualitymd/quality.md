---
type: Functional Specification
title: /quality wizard
description: Behavioral component spec for read-only QUALITY.md state inspection and workflow routing.
tags: [skill, quality, mode, wizard]
timestamp: 2026-06-22T00:00:00Z
---

# /quality wizard

`wizard` is the `/quality` skill's read-only wayfinder. It implements the shared
contracts in the parent [/quality skill](../quality-skill.md) spec and owns only
the wizard-specific behavior below.

The runtime procedure lives at
[`skills/quality/modes/wizard.md`](../../../../skills/quality/modes/wizard.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`wizard` is selected for a bare `/quality`, ambiguous direction, `status`,
`next`, `review model`, `review history`, and other requests that ask what to do
with a QUALITY.md model rather than directly choosing setup, evaluate, update,
or recommendation follow-up.

Its purpose is to inspect current state, classify readiness, recommend one next
workflow, and offer concrete alternatives. Its judgments are routing findings,
not Evaluation ratings.

## Mutation surface and artifacts

`wizard` **MUST** be read-only. It **MUST NOT** modify files, create evaluation
records, build reports, write the quality log, apply recommendations, update
tooling, or rate evaluated source.

`wizard` may read CLI version/status output, the resolved `QUALITY.md`, bounded
model-quality guidance, evaluation-history metadata, and quality-log metadata.
It **MUST NOT** inspect evaluated source files or read full evaluation report
bodies during the model-lifecycle checklist pass.

## Required flow

`wizard` **MUST** follow a probe -> classify -> recommend -> offer flow:

1. Probe CLI readiness, model-file presence, model validity, model source
   coverage, evaluation-history status, and model-history status through the CLI
   and bounded file inspection.
2. Inspect model lifecycle with the Top 10 QUALITY.md checks when the model is
   structurally valid.
3. Classify readiness using the best matching state: no setup, invalid model,
   starter or skeleton model, usable but immature model, ready to evaluate, has
   evaluation history, or mature but needs maintenance or reconciliation.
4. Recommend one next workflow with evidence from the observed state.
5. Offer concrete runnable alternatives.

`wizard` **MUST** judge CLI readiness offline from the visible version against
the skill's prerequisite range. It **MUST NOT** probe the network for newer
versions; release discovery belongs to [`update`](update.md).

When the user asks to review or improve the `QUALITY.md` itself, `wizard`
**MUST** use the authoring guide as the model-quality reference and route to a
confirmed editing workflow rather than treating `QUALITY.md` as the root Area of
an Evaluation report.

## Output contract

`wizard` output **MUST** be status-first and action-oriented. It should include:

- CLI, QUALITY.md, model, evaluation-history, model-history, and readiness
  status;
- top QUALITY.md inspection findings or an explicit non-blocking state;
- one recommended next step; and
- a short numbered menu of concrete workflow options.

## Stop conditions

`wizard` does not stop as a failure when setup, lint, history, or update work is
needed. It reports the observed state and routes to the next workflow. It
**MUST** stop only when it cannot inspect enough local state to make a truthful
readiness statement, and the stop response **MUST** name the missing evidence.

## Completion criteria

`wizard` is complete when it has reported current state, named the best next
workflow, and offered concrete alternatives without mutating the workspace.
