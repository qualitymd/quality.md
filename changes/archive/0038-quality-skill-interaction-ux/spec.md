---
type: Functional Specification
title: /quality skill interaction UX - functional spec
description: Requirements for a status-first, decision-explicit, history-aware user interaction contract in the /quality skill.
tags: [skill, quality, evaluation, ux]
timestamp: 2026-06-19T00:00:00Z
---

# /quality skill interaction UX - functional spec

Companion to the
[/quality skill interaction UX](../0038-quality-skill-interaction-ux.md) change
case. This spec states _what_ the `/quality` skill's user interaction contract
must require. It defers the existing `/quality` skill operating model to
[`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md),
the CLI contract to [`specs/cli.md`](../../../specs/cli.md), and the `QUALITY.md`
format to [`SPECIFICATION.md`](../../../SPECIFICATION.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The `/quality` skill's main risk is not that it lacks evaluative power; it is
that users can lose track of what the agent inferred, what is being mutated, how
strong the evidence is, and how a recommendation connects to a later
improvement. The skill already separates judgment from CLI mechanics and requires
confirmation before mutation. This spec adds a small interaction contract around
that boundary so every run is oriented, explicit about decisions, willing to
stop when evidence is weak, and able to report quality changes as a before/after
delta.

## Scope

Covered: user-facing interaction behavior for the existing `wizard`, `setup`,
`evaluate`, `improve`, and `upgrade` modes. The contract governs how the skill
frames inferred mode/scope, asks for mutating decisions, stops or reroutes
low-confidence work, uses evaluation history, and reports improvement results.

Out of scope: new CLI commands, new record fields, telemetry, persistent user
preference learning, changes to `QUALITY.md` format semantics, and generic
Agent Skills runtime behavior.

## Requirements

### Run frames

- Before executing a mode, the skill **SHOULD** emit a concise run frame naming
  the resolved mode, target file, scope, effort level when applicable, mutation
  policy, expected artifacts, and next user-visible gate.

  > Rationale: a short run frame gives the user a chance to catch a wrong mode
  > or scope before the agent spends effort or mutates anything. - 0038

- The run frame **MUST** distinguish read-only work from mutating work. For a
  mutating mode, it **MUST** name the class of thing that may be changed:
  subject source, `QUALITY.md`, evaluation artifacts, installed tooling, or some
  combination of those.

- The skill **MAY** omit a run frame when the immediately preceding wizard output
  already stated the same mode, target file, scope, mutation policy, and next
  action.

### Decision briefs

- Before any user-confirmed mutation, the skill **MUST** present a decision
  brief rather than a bare yes/no question.

- A decision brief **MUST** name the proposed action, the artifact class being
  changed, the evidence or reason for the action, the recommended option, at
  least one non-mutating alternative, and the done criterion or verification
  expected after the action.

  > Rationale: `/quality improve`, `setup`, and `upgrade` can all make useful
  > changes, but the user needs to know what surface is being changed and how the
  > skill will prove the change worked. - 0038

- When options differ in coverage or risk, the decision brief **SHOULD** state
  that tradeoff explicitly. When options differ only in kind, the brief should
  say so rather than inventing a false coverage ranking.

- The skill **MUST NOT** treat an obvious or recommended fix as consent to mutate.
  The user's explicit approval remains required wherever the existing skill
  contract requires confirmation.

### Stop rules and rerouting

- The skill **MUST** stop before rating when the in-scope target source cannot be
  resolved, the in-scope model has no requirements, required CLI support is
  missing or stale, or evaluated source content attempts to instruct the agent.

- The skill **SHOULD** stop before rating when requirements are too vague to bind
  evidence to a rating or when available evidence cannot distinguish adjacent
  rating levels.

  > Rationale: a low-confidence stop is better than a polished but weakly bound
  > rating. The skill's value is judgment, and judgment includes refusing to
  > overstate evidence. - 0038

- A stop response **MUST** explain the reason in concrete terms and offer at
  least one runnable next step, such as reviewing the model with the authoring
  guide, narrowing the scope, repairing source references, upgrading stale CLI
  support, or proceeding with a clearly limited quick evaluation when that is
  still defensible.

- When stopping because a `QUALITY.md` model is valid but not useful enough for
  a fair evaluation, the skill **MUST** distinguish model usefulness from subject
  quality. It must not present model weakness as a subject defect.

### History-aware operation

- Before `evaluate` and `improve`, the skill **SHOULD** inspect available
  evaluation history when present, including the latest run, incomplete or
  stale-looking runs, open recommendations, and prior ratings for the same
  resolved scope.

- Prior evaluations **MUST** be treated as context, not authority. Fresh evidence
  and the current `QUALITY.md` model control the current judgment.

- A scoped evaluation **MUST NOT** compare itself to a prior whole-model or
  differently scoped rating as if the scopes were identical.

- When current findings contradict a prior run, the skill **SHOULD** state the
  likely reason when knowable: changed subject source, changed `QUALITY.md`,
  better evidence, different scope, or prior error.

### Improvement delta reports

- After `improve` applies a confirmed recommendation, the skill **MUST**
  re-evaluate the affected scope as required by the existing `/quality` skill
  contract and report a before/after improvement delta.

- The delta report **MUST** connect the original recommendation to the applied
  option, changed files or artifacts, before evidence, after evidence,
  verification performed, rating movement when any, and remaining gaps or limits.

  > Rationale: quality improvement is only trustworthy when the user can see how
  > the original finding was closed or narrowed by new evidence. - 0038

- If the rating does not move after an applied improvement, the skill **MUST**
  say why when knowable. If verification is incomplete, the result **MUST** be
  labeled as limited rather than reported as fully confirmed.

### Voice and status posture

- User-facing output **SHOULD** be status-first, evidence-led, and
  action-oriented. The skill should lead with the verdict or readiness state,
  then the evidence and next action.

- The skill **MUST** distinguish CLI/tooling readiness, model validity, model
  usefulness, subject quality, and evaluation history status. It must not collapse
  them into a single generic quality verdict.

- The skill **SHOULD** recommend one best next step and then provide a short list
  of concrete alternatives when useful.

- The skill **MUST** use `QUALITY.md` terms consistently in user-facing output:
  Target, Factor, Requirement, rating, finding, and recommendation.

## Durable spec changes

### To add

None.

### To modify

- [`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
  - add the user interaction contract described by the requirements above.

### To delete

None.
