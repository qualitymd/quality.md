---
type: Change Case
title: Natural scope labels
description: Make natural Area and Factor labels the primary documented scoped-evaluation input for the /quality skill.
status: Done
tags: [skill, evaluation, ux, references]
timestamp: 2026-06-23T00:00:00Z
---

# Natural scope labels

A **Change Case** updating the `/quality` skill contract and user-facing docs so
scoped evaluations lead with natural Area and Factor labels instead of canonical
model references. The detail lives in its
[functional spec](0061-natural-scope-labels/spec.md) and
[design doc](0061-natural-scope-labels/design.md).

## Motivation

Canonical model references such as `factor:root::security` are precise and
necessary for durable artifacts, but they are too technical as the first thing a
person sees when asking the skill to evaluate a familiar quality concern.

The skill already owns judgment at the human edge. It should let users ask for
the thing they mean in natural project vocabulary, resolve that request against
the grounded `QUALITY.md` model, and ask a targeted clarification question only
when the label is ambiguous. Stable IDs remain the internal and durable form.

## Scope

Covered:

- Document natural Area and Factor labels as the primary scoped-evaluation input
  for `/quality evaluate`.
- Define one-label resolution as a normal path for unique Area or Factor labels.
- Define two-label resolution as `<area-label> <factor-label>`, resolving the
  Area first and then the Factor within that Area.
- Require targeted clarification when a Factor label appears in multiple Areas,
  using the question: "What area do you want to evaluate <Factor> for?"
- Preserve qualified model references as supported exact-addressing input.
- Preserve stable model IDs in evaluation records and `report.json`.

Deferred / non-goals:

- No change to QUALITY.md model-reference grammar.
- No change to evaluation record or report JSON schema.
- No broad CLI selector redesign.
- No fuzzy matching requirement beyond exact label/name matching unless a later
  design explicitly adds it.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0061-natural-scope-labels/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [x] No planned `qualitymd` CLI or Go code changes; this case changes the
      skill's human-edge contract and documentation unless implementation
      analysis finds missing runtime support.

Specs:

- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      - make natural labels the primary scoped-evaluation input and move
      qualified references to exact-addressing guidance.
- [x] [`specs/skills/quality-skill/evaluation.md`](../../specs/skills/quality-skill/evaluation.md)
      - align evaluation-mode scope resolution and ambiguity behavior.
- [x] OKF logs under [`specs/`](../../specs/log.md) - record durable spec updates
      when they land.

Runtime skill and docs:

- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - update argument
      parsing and invocation examples so one-label and two-label scoped
      evaluation are the documented happy path.
- [x] [`skills/quality/modes/evaluate.md`](../../skills/quality/modes/evaluate.md)
      - update the scope-resolution decision tree and ambiguity prompts.
- [x] [`README.md`](../../README.md) - update usage examples to foreground natural
      labels and present qualified references as exact/advanced syntax.
- [x] [`CHANGELOG.md`](../../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

## Children

- [Functional spec](0061-natural-scope-labels/spec.md) - what natural scoped
  input must resolve, preserve, and ask.
- [Design doc](0061-natural-scope-labels/design.md) - how the skill resolves
  natural labels while preserving stable model IDs.

## Status

`Done`. Implemented and archived. Natural labels are now the primary documented
scoped-evaluation input for the `/quality` skill, qualified model references
remain available for exact addressing, and durable artifacts continue to use
stable identifiers.
