---
type: Functional Specification
title: Agent-mediated UX conformance — functional spec
description: Requirements for bringing live agent-mediated workflow guidance and durable skill specs into conformance with the agent-mediated UX guide.
tags: [skill, ux, agents, workflows]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-mediated UX conformance — functional spec

Companion to
[Agent-mediated UX conformance](../0084-agent-mediated-ux-conformance.md). This
spec states what the conformance pass must do. No design doc is required unless
the implementation phase discovers a reusable rendering abstraction or code path.

The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Scope

This change applies
[Designing agent-mediated UX](../../../docs/guides/agent-mediated-ux.md) to live
agent-mediated workflow surfaces in this repository: the bundled `/quality`
skill, durable `/quality` skill specs, recommendation follow-up guidance, and
contributor/public docs that teach those workflows.

Historical records, archived Change Cases, append-only logs, recorded evaluation
examples, and CLI-only human output are out of scope unless they are live
templates for future agent-mediated output.

## Requirements

### Inventory

The implementation **MUST** audit live agent-mediated workflow surfaces before
editing them. The audit **MUST** include the parent `/quality` skill interaction
contract, `setup`, `evaluate`, `update`, recommendation follow-up, public
setup/evaluation usage snippets, and any code or scaffolded text that appears to
hardcode agent-mediated workflow output.

The implementation **MUST** leave historical material alone unless it is also a
live template or current source of truth.

> Rationale: a repo-wide UX conformance pass can easily rewrite history or
> non-target surfaces. The artifact sweep is part of the work, not a nice-to-have;
> it protects archived records and keeps CLI UX separate from agent-mediated UX.

### Shared interaction contract

The shared `/quality` interaction contract **MUST** require live user-facing
workflow output to follow the agent-mediated UX guide's output hierarchy:
status first, primary user action second, supporting context after that.

For each interaction block where the user must answer, approve, correct, choose,
or act, the primary question or call to action **MUST** be visually emphasized,
preferably with bold Markdown.

Supporting labels such as `Recommended`, `Why it matters`, `Confidence`,
`Changed`, `Validation`, `Important gaps`, and `Next` **SHOULD** be bold when
rendered in Markdown so the left edge of the output is scannable.

The shared contract **MUST NOT** change mutation boundaries, confirmation
requirements, rating semantics, evaluation semantics, or CLI responsibilities.

> Rationale: this case is about how the agent renders workflow state and choices,
> not about what the workflows are allowed to do.

### Discovery and checkpoints

Setup discovery questions **MUST** present the primary question as the strongest
visual element in the interaction block. The supporting material **MUST** keep
the question's purpose, recommended answer, confidence/evidence, and shortest
acceptable response adjacent to the question.

Human context checkpoints **MUST** make the primary correction action explicit
and easy to scan. When several inferred values are reviewed together, the
checkpoint **SHOULD** use a table or similarly compact structure unless the
agent's interaction surface makes that less usable.

Setup **MUST NOT** let emphasis obscure the teaching copy required by the setup
workflow. The formatting should make the teaching easier to scan, not remove it.

### Progress and long workflows

Multi-step workflows **SHOULD** show visible progress at phase boundaries where
the user's mental model would otherwise drift: before a context scan, after a
tool-dependent phase, before a mutation gate, and at closeout.

Progress output **MUST** remain factual and user-facing. It **MUST NOT** become a
transcript of internal reasoning or tool chatter.

### Decision gates

Decision briefs **MUST** visually emphasize the decision question or call to
action. They **MUST** keep the proposed changes, evidence/reason, recommended
option, alternatives, and done criterion in a consistent, scannable shape.

Where a workflow's mutation boundary matters, the decision brief or nearby
context **MUST** state what will not happen as well as what will happen.

### Closeout and stop responses

Workflow closeouts **MUST** report outcome, changed artifacts when any,
validation, important gaps or limitations, and the recommended next action in a
status-first shape.

Stop responses **MUST** make the reason and next useful action clear without
making the user reconstruct the path from logs or command output.

Evaluate closeouts **MUST** preserve required evaluation content: rating, scope,
evidence basis, recommendations or lack of gaps, and known limitations. The
agent-mediated UX pass may improve the presentation, but it **MUST NOT** weaken
the evaluation report contract.

### Emoji and emphasis

Emoji **MAY** be used as semantic markers only. The implementation **MUST NOT**
add decorative emoji to every heading or label.

Existing Rating Level emoji (`🟢`, `🔵`, `🟡`, `🔴`) **MUST** remain scanning aids
for human display titles only; they **MUST NOT** become rating identity,
ordering, or semantics.

Markdown emphasis **MUST** be used as hierarchy, not decoration. Whole paragraphs
and repeated prose **MUST NOT** be bolded simply to make the output feel more
polished.

### Verification

Before moving the Change Case to `In-Review`, the implementation **MUST**
reconcile the parent Change Case's affected-artifact list against the actual
edits.

Verification **MUST** include `mise run fmt-md-check`. If the implementation
touches Go code or CLI output tests, verification **MUST** include the narrowest
relevant Go tests and `mise run check` unless a documented environment blocker
prevents it.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` — add the shared
  agent-mediated UX output hierarchy and primary CTA/question emphasis
  requirements.
- `specs/skills/quality-skill/workflows/setup.md` — require setup discovery,
  checkpoint, review, and closeout presentation to follow the shared UX guide.
- `specs/skills/quality-skill/workflows/evaluate.md` — require evaluation
  progress, stop, summary, limitation, and next-action output to follow the
  shared UX guide while preserving evaluation content requirements.
- `specs/skills/quality-skill/workflows/update.md` — require update planning,
  confirmation, verification, and restart/reload output to follow the shared UX
  guide.
- `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` and
  `specs/skills/quality-skill/recommendation-follow-up.md` — require
  recommendation apply/handoff decisions and closeouts to follow the shared UX
  guide.
- `specs/skills/quality-skill/reporting.md` — verify and update if user-facing
  evaluation summaries need a pointer to the shared interaction contract.
- `specs/log.md` and `specs/skills/quality-skill/workflows/log.md` — record the
  durable spec updates.

### To rename

None.

### To delete

None.
