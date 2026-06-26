---
type: Functional Specification
title: Binary confirmation UX — functional spec
description: Requirements for aligning /quality binary mutation confirmations with y/n response guidance.
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Binary confirmation UX — functional spec

Companion to the [Binary confirmation UX](../0106-binary-confirmation-ux.md)
change case. This spec states what the change must do; the
[design doc](design.md) covers how.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

Agent-mediated confirmations should minimize translation at the exact point
where a user authorizes mutation. The general closed-choice rule from 0099 made
`1` the shortest accept path for small fixed option sets, which remains right for
multi-option choices. A true binary confirmation is different: the user is
answering yes or no to a named action, so `y`/`n` is more idiomatic and less
ambiguous than mixing `1` with words such as `skip`.

The shared UX guide now documents that distinction. This change applies it to
the durable `/quality` skill specs and bundled runtime guidance.

## Scope

Covered: `/quality` shared interaction guidance plus binary mutation gates in
the update workflow, recommendation follow-up, and setup's fallback
existing-file edit gate.

Not covered: setup discovery question content, evaluation scope-routing choices,
recommendation apply-vs-handoff outcome selection, update mechanics,
issue-tracker mechanics, CLI behavior, Go implementation, format schema,
evaluation records, reports, or rating semantics.

## Assumptions & dependencies

- The current
  [agent-mediated UX guide](../../../docs/guides/agent-mediated-ux.md) is the
  normative UX source. This case applies the guide's binary-confirmation
  exception rather than redefining it.
- The durable `/quality` skill specs under `specs/skills/quality-skill/` mirror
  runtime skill behavior closely enough that any interaction contract change must
  update both durable specs and bundled runtime guidance.

## Requirements

### Shared interaction contract

- The durable `/quality` skill specs and bundled runtime skill guidance **MUST**
  distinguish true binary confirmations from other small closed-choice prompts:
  non-binary closed-choice prompts keep numbered options with the recommended
  option first, while true binary confirmations use visible `y`/`n` as the
  shortest responses.
  > Rationale: `1` remains correct when the user chooses among several options.
  > For a yes/no mutation gate, `y`/`n` matches the question's semantic shape and
  > avoids mixed answer vocabulary such as `1` versus `skip`. — 0106

- When a decision brief asks a true binary confirmation, the durable and runtime
  guidance **MUST** include an explicit answer path equivalent to `Reply y to
  <perform the action>, or n to <decline/skip/stop>`.
  > Rationale: a decision brief names the action, evidence, alternatives, and
  > verification, but it still needs to tell the user what to type. The observed
  > update prompt had the right plan fields but the wrong answer shape. — 0106

- Runtime guidance **SHOULD** accept obvious aliases such as `yes`, `no`, `1`,
  action words, or skip/stop words when they unambiguously match the displayed
  options, but those aliases **MUST NOT** replace visible `y`/`n` as the
  shortest responses for binary confirmations.

### Workflow gates

- The update workflow's update-plan confirmation brief **MUST** show `y`/`n` as
  the visible shortest answer path before running any owner command that mutates
  installed tooling.

- Recommendation follow-up's local apply decision brief **MUST** show `y`/`n` as
  the visible shortest answer path after the recommendation option and mutation
  surface have been selected.

- Recommendation follow-up's external issue creation decision brief **MUST** show
  `y`/`n` as the visible shortest answer path before creating an external issue.

- Setup's fallback decision brief for updating an existing `QUALITY.md` **MUST**
  show `y`/`n` as the visible shortest answer path when it is presented as a true
  binary update-or-stop confirmation.

### Preserved numbered prompts

- Setup discovery calibration questions, evaluation ambiguity prompts, evaluate
  stop-response routing options, and recommendation follow-up's initial
  apply-vs-handoff outcome selection **MUST NOT** be converted to binary `y`/`n`
  prompts solely because this change exists.
  > Rationale: those prompts ask the user to choose among multiple options or
  > provide a correction/routing decision. Converting them to yes/no would erase
  > information and weaken the interaction. — 0106

- Setup's final review gate **MAY** continue to use a literal confirmation such
  as `looks good` when the prompt is framed as a review-and-correction gate
  rather than a plain yes/no confirmation.
  > Rationale: setup's final review is intentionally not a bare binary question:
  > it invites corrections and broader last-call context while preserving an
  > explicit fast path. — 0106

## Durable spec changes

### To add

None.

### To modify

- [`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
  - add the binary-confirmation exception to the shared closed-choice and
    decision-brief contracts. Driven by
    [Shared interaction contract](#shared-interaction-contract).
- [`specs/skills/quality-skill/workflows/update.md`](../../../specs/skills/quality-skill/workflows/update.md)
  - require `y`/`n` on the update-plan confirmation brief. Driven by
    [Workflow gates](#workflow-gates).
- [`specs/skills/quality-skill/recommendation-follow-up.md`](../../../specs/skills/quality-skill/recommendation-follow-up.md)
  - require `y`/`n` on local apply and external issue creation decision briefs.
    Driven by [Workflow gates](#workflow-gates).
- [`specs/skills/quality-skill/workflows/setup.md`](../../../specs/skills/quality-skill/workflows/setup.md)
  - require `y`/`n` on the fallback existing-file update brief while preserving
    the final review gate's current fast path. Driven by
    [Workflow gates](#workflow-gates) and
    [Preserved numbered prompts](#preserved-numbered-prompts).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, `/quality` will consistently present
binary mutation confirmations with `y`/`n`, while preserving numbered responses
for non-binary choices. That achieves the motivation without changing CLI
mechanics, model authoring, evaluation semantics, or the QUALITY.md format.
