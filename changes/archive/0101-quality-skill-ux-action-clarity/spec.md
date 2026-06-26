---
type: Functional Specification
title: Quality skill UX action clarity — functional spec
description: Requirements for aligning /quality skill prompts, checkpoints, gates, and closeouts with the current agent-mediated UX guide.
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Quality skill UX action clarity — functional spec

Companion to the
[Quality skill UX action clarity](../0101-quality-skill-ux-action-clarity.md)
change case. This spec states *what* the change must do; it does not choose the
implementation sequence.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

The current [agent-mediated UX guide](../../../docs/guides/agent-mediated-ux.md)
is the binding UX source for live `/quality` workflow output. It requires
status-first interaction blocks, visually dominant primary questions or calls to
action, adjacent recommendation/evidence/answer fields, numbered closed-choice
prompts where `1` is the shortest accept path, explicit decision briefs before
mutation, and closeouts that name outcome, validation, gaps, and next action.

A review of the runtime `/quality` skill found several conformance gaps. The
largest failure mode is not missing prose; it is weak hierarchy. Some prompts
contain the right facts but leave the user to infer what to type, which option is
recommended, whether confirmation mutates an artifact, or what happens next.
This change makes the prompt shapes themselves carry the UX contract.

## Scope

Covered: durable `/quality` skill specs and bundled runtime skill guidance for
setup, evaluate, update, and recommendation follow-up interaction shapes.

Not covered: CLI behavior, Go implementation, format schema, evaluation record
shape, rating semantics, report rendering, issue-tracker integration mechanics,
or model-authoring/evaluation judgment rules beyond how the skill presents user
interactions.

## Assumptions & dependencies

- The current
  [agent-mediated UX guide](../../../docs/guides/agent-mediated-ux.md) is the
  normative UX source; this case applies it rather than replacing it.
- The current
  [functional-spec guide](../../../docs/guides/write-functional-specs.md) and
  [change-case guide](../../../docs/guides/work-with-change-cases.md) govern this
  case's spec quality and Draft-to-Design validation.
- The durable `/quality` skill specs under `specs/skills/quality-skill/` mirror
  runtime skill behavior closely enough that implementation must update both
  durable specs and bundled runtime guidance where a prompt contract changes.

## Requirements

### Shared interaction contract

- The durable `/quality` skill specs and bundled runtime skill guidance **MUST**
  preserve the current agent-mediated UX guide as the normative source for
  status-first output, visual hierarchy, numbered choices, decision briefs,
  progress updates, closeouts, and semantic-only emoji.

- Where this change names a prompt, checkpoint, gate, or closeout, the
  corresponding durable spec and runtime guidance **MUST** make the shortest
  acceptable user response explicit.
  > Rationale: the reviewed failure mode was often not a missing instruction but
  > a weak affordance. The guide's checklist now names "the shortest acceptable
  > user response is clear" as a shipping check, so each affected prompt needs
  > that response path in the local contract. — 0101

- User-facing examples introduced or revised by this change **SHOULD** wrap
  concrete files, commands, model references, fields, and literal user replies in
  code spans.
  > Rationale: the UX guide requires code formatting for exact files, commands,
  > fields, model references, IDs, and literal values. The reviewed skill output
  > sometimes named concrete artifacts such as QUALITY.md or `/quality evaluate`
  > as plain prose, which weakens scanning and precision. — 0101

### Setup workflow

- The setup runtime and durable workflow specs **MUST** make the first two
  discovery questions, root area and domain, explicit agent-mediated questions:
  visually emphasized primary question, adjacent `Why it matters`,
  `Recommended`, `Confidence`, and `Answer` labels when Markdown is available,
  and a shortest accept path such as replying `yes` or naming a correction.
  > Rationale: closed-choice setup prompts were already strengthened in 0099, but
  > the open-ended setup questions still lacked the same clear answer affordance.
  > Users should not infer how to accept a recommended default. — 0101

- The setup human context checkpoint **MUST** present its primary correction
  action as the visually strongest element and **MUST** include an explicit
  `Answer` line that says how to accept the draft or provide terse corrections.
  When reviewing several inferred values together, it **SHOULD** render a compact
  table or equivalent structure with item, draft, confidence, and requested user
  action.
  > Rationale: the guide says checkpoints are for correcting inferred context and
  > that the primary call to action must still be explicit. The current runtime
  > template can bury "correct this draft" in prose while repeating "Why it
  > matters" blocks, making the CTA weaker than the supporting context. — 0101

- The setup human context checkpoint **MUST** place unresolved-context caveats
  after the primary action and before or near the table, without making the caveat
  the first or strongest user-facing element.
  > Rationale: warning that unresolved items become Unknowns is important, but it
  > supports the action. Leading with consequence makes the user reconstruct the
  > actual task from surrounding explanation. — 0101

- The setup final review gate **MUST** use a decision-brief shape before setup
  writes or edits `QUALITY.md`, whether the file is newly scaffolded or already
  exists. The brief **MUST** name the proposed change, evidence or reason,
  recommended option, non-mutating alternatives, and done criterion or
  verification.
  > Rationale: the final review gate is the mutation boundary for creating or
  > editing `QUALITY.md`. The current existing-file path has a decision brief,
  > but the new-file path can advance from a conversational "looks good" prompt
  > without the full mutation brief the UX guide requires. — 0101

- The setup final review prompt **MUST** keep the confirmation or correction
  action visually primary. Broader last-call context such as priorities, worries,
  edge cases, or repo-invisible facts **MAY** remain, but it **MUST NOT** be the
  final or visually dominant call to action.
  > Rationale: open-ended final context is useful, but the current prompt can end
  > on a broad catch-all. The guide warns that broad catch-alls make specific
  > dimensions feel optional and can obscure the primary action. — 0101

### Evaluate workflow

- When evaluation scope resolution yields several runnable Area or Factor
  choices, the skill **MUST** present a numbered choice list with human-readable
  labels first, qualified model references as secondary context where useful, and
  an `Answer` line that accepts a number. If there is no evidence-backed
  recommendation, the prompt **MUST NOT** invent one.
  > Rationale: ambiguity prompts are closed choices even when no option is
  > recommended. Numbering gives the user a terse response path without forcing
  > the skill to pretend one Area is preferred. — 0101

- Evaluate stop responses **MUST** include a clear answer path after listing
  concrete runnable options, such as replying with a number to start a next
  workflow or saying `stop`.
  > Rationale: the current stop shape names options but not what the user should
  > type. The UX guide requires every user-facing step to make the next action
  > obvious. — 0101

### Update workflow

- The update runtime and durable workflow specs **MUST** require `/quality update`
  to emit a public run frame before tool inspection, naming the mode, mutation
  surface, artifacts or lack of local project artifacts, and next visible gate.
  > Rationale: the shared skill contract says public modes emit a run frame, but
  > the update workflow begins with inspection steps. Update can mutate installed
  > tooling, so the user needs orientation before those checks run. — 0101

- The update workflow **SHOULD** show concise progress or status after version
  inspection and before any mutation gate, without turning the output into an
  internal transcript.
  > Rationale: the UX guide calls for progress in multi-step workflows especially
  > after tool-dependent phases and before mutation gates. Update is short, but
  > its tool inspection can materially affect the user's trust in the plan. — 0101

### Recommendation follow-up

- When recommendation follow-up has not yet received an explicit apply-vs-handoff
  choice, it **MUST** present the two productive outcomes as numbered options and
  include an `Answer` line. When the available recommendation evidence supports a
  recommended path, that path **SHOULD** be option `1`; otherwise the prompt
  **MUST** avoid inventing a recommendation.
  > Rationale: the workflow allows only apply or handoff as productive outcomes.
  > Presenting them as numbered options makes that boundary operational instead
  > of just descriptive. — 0101

- Recommendation follow-up issue creation **MUST** use a decision brief before
  creating any external issue. The brief **MUST** name what external artifact will
  be created, what local artifacts will not change, the evidence or reason, the
  recommended option, alternatives, and verification.
  > Rationale: creating an external issue is a mutation. The current guide
  > requires explicit confirmation, and the shared UX guide requires mutation
  > gates to state the change, reason, alternatives, and done criterion. — 0101

- Recommendation result closeouts **MUST** include a `Next` field and **SHOULD**
  include `Not done` when the boundary matters, such as no evaluation rerun, no
  issue creation, no model change, or no quality-log entry.
  > Rationale: closeout should prevent the user from reconstructing success and
  > scope from logs. The current result template names what happened but can omit
  > the next action and explicit non-actions. — 0101

## Durable spec changes

### To add

None.

### To modify

- [`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
  — strengthen the shared interaction contract for explicit shortest responses,
  public mode run frames, decision-brief coverage, and code-span precision.
  Driven by [Shared interaction contract](#shared-interaction-contract).
- [`specs/skills/quality-skill/workflows/setup.md`](../../../specs/skills/quality-skill/workflows/setup.md)
  — update setup discovery, checkpoint, and final review-gate contracts. Driven
  by [Setup workflow](#setup-workflow).
- [`specs/skills/quality-skill/evaluation.md`](../../../specs/skills/quality-skill/evaluation.md)
  and
  [`specs/skills/quality-skill/workflows/evaluate.md`](../../../specs/skills/quality-skill/workflows/evaluate.md)
  — update scope ambiguity and stop-response interaction contracts. Driven by
  [Evaluate workflow](#evaluate-workflow).
- [`specs/skills/quality-skill/workflows/update.md`](../../../specs/skills/quality-skill/workflows/update.md)
  — update run-frame and progress sequencing contracts. Driven by
  [Update workflow](#update-workflow).
- [`specs/skills/quality-skill/recommendation-follow-up.md`](../../../specs/skills/quality-skill/recommendation-follow-up.md)
  and
  [`specs/skills/quality-skill/guides/recommendation-follow-up-md.md`](../../../specs/skills/quality-skill/guides/recommendation-follow-up-md.md)
  — update outcome selection, issue-creation gate, and closeout contracts. Driven
  by [Recommendation follow-up](#recommendation-follow-up).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, the `/quality` skill will no longer
carry the identified agent-mediated UX gaps: the reviewed prompts will expose
the primary action, the shortest response path, mutation consequence, and next
step directly at the point of interaction. That matches the motivation without
changing CLI behavior, evaluation semantics, model authoring semantics, or the
QUALITY.md format.
