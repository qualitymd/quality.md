---
type: Design Doc
title: Real Review Gates - design
description: How shared UX guidance and direct QUALITY.md authoring align around real review gates.
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Real Review Gates - design

## Context

Answers the [functional spec](spec.md) for change case
[0139](../0139-real-review-gates.md). The change is documentation and runtime
prompt behavior: no CLI or Go implementation is involved. It builds on
[0138](../archive/0138-lightweight-authoring-checkpoint.md), which introduced
the lightweight direct-authoring checkpoint but did not explicitly say that the
checkpoint stops and waits.

## Approach

### Make the shared UX doctrine explicit

Update `docs/guides/agent-mediated-ux.md` in the most general terms:

- add the principle that feedback invitations are gates;
- add the quick-acknowledgement pattern for long non-workflow reads;
- add a Review Gates section between Checkpoints and Decision gates; and
- add checklist items so future workflow reviews catch the failure mode.

The Review Gates section distinguishes informational previews, review gates, and
decision gates. That vocabulary is the main design move: it lets workflows ask
for conversational feedback without pretending every review is a binary `y`/`n`
decision, while still preserving the hard rule that asking requires waiting.

### Apply the doctrine to direct authoring

Update `skills/quality/SKILL.md` direct model-authoring guidance to:

1. acknowledge likely direct authoring before long reads;
2. inspect `QUALITY.md` and relevant authoring guides;
3. infer intent and ask only material follow-up;
4. render the checkpoint as a review gate;
5. wait for `looks good`, corrections, or equivalent approval; and
6. only then mutate.

Change the sample checkpoint wording from "Anything you want adjusted first?" to
"Anything you want adjusted before I make that edit?" so the gate's wait state is
part of the prompt.

### Mirror durable specs and guide contracts

Update the durable parent skill spec to carry the same direct-authoring
requirements. Update the authoring guide and its durable guide spec to say that
inviting concerns/goals/needs/worries/constraints or `looks good` requires
waiting before mutation.

### Logs and release notes

Record the guide, spec, runtime skill, and runtime guide updates in their
existing logs, and add an Unreleased `/quality Skill` changelog entry. No format,
CLI, generated schema, or setup workflow change is needed.

## Spec response

- **Shared UX guidance** - satisfied by updating the guide's core principle,
  opening section, new Review Gates section, and checklist.
- **Direct QUALITY.md authoring** - satisfied by sharpening the runtime skill
  prompt and durable skill spec.
- **Authoring guide alignment** - satisfied by updating the runtime authoring
  guide and durable authoring guide spec.
- **Verification** - satisfied by source inspection and Markdown formatting.

## Alternatives

- **Always use a binary decision brief for direct authoring.** Rejected. It would
  fix the wait problem but make conversational model edits heavier than needed.
- **Let clear intent proceed without review.** Rejected for judgment-shaping
  `QUALITY.md` model changes; the user still needs a chance to catch a wrong
  inference.
- **Only update the `/quality` skill.** Rejected. The failure is a general
  agent-mediated UX principle and belongs in the shared guide too.

## Trade-offs & risks

- Requiring a review gate for direct model changes adds one user turn. That is
  intentional for edits that reshape future judgment, and the immediate
  acknowledgement reduces perceived stall before the checkpoint.
- Agents still need judgment to distinguish informational previews from review
  gates. The guide's false-affordance example makes the boundary concrete.

## Open questions

None.
