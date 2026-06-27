---
type: Design Doc
title: Casual Review Gate Wording - design
description: How shared UX guidance and direct QUALITY.md authoring adopt casual planned-change review gates.
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Casual Review Gate Wording - design

## Context

Answers the [functional spec](spec.md) for change case
[0140](../0140-casual-review-gate-wording.md). The change builds on
[0139](../0139-real-review-gates.md), which made feedback invitations
real gates. This case refines the voice and structure of those gates so routine
direct `QUALITY.md` edits feel more conversational.

## Approach

### Prefer planned-change prose with a value prop

Update `docs/guides/agent-mediated-ux.md` so review gates for content or model
edits default to a short plan block:

```text
Here's what I'm planning to do:

<simple common-sense prose of the change>, so that <value prop>.
```

The `so that` clause is not ceremony; it exposes the intended value of the edit.
That lets the user correct not only the mechanics, but also the goal the agent is
optimizing for.

### Keep the gate friendly and real

Keep the 0139 wait rule intact. The feedback prompt should feel like an
invitation to shape the edit, not a legal approval line. The guide and skill use
wording such as "what should I adjust or watch out for" and keep `looks good` as
the shortest approval path.

### Avoid default action lists

Use a prose sentence for ordinary direct edits. A numbered action list is still
allowed when a multi-part edit would be hard to scan in prose, but it should not
be the default review-gate shape.

### Mirror durable specs and guide contracts

Update the durable parent skill spec and authoring guide spec to require the
planned-change and value-prop checkpoint shape. Update the runtime authoring
guide so agents preserve the value prop as one of the considerations behind a
direct edit.

### Logs and release notes

Record the guide, spec, runtime skill, and runtime guide updates in their
existing logs, and update the Unreleased `/quality Skill` changelog entry. No
format, CLI, generated schema, or setup workflow change is needed.

## Spec response

- **Shared UX guidance** - satisfied by updating the Review Gates section and
  checklist language.
- **Direct QUALITY.md authoring** - satisfied by sharpening the runtime skill
  checkpoint template and durable skill spec.
- **Authoring guide alignment** - satisfied by updating the runtime authoring
  guide and durable authoring guide spec.
- **Verification** - satisfied by source inspection and Markdown formatting.

## Alternatives

- **Use numbered planned next actions by default.** Rejected. It is clearer for
  complex edits, but it makes small authoring requests feel procedural.
- **Ask only "What about that?" with no examples.** Rejected. It is friendly, but
  giving examples such as concerns, goals, edge cases, and naming preferences
  helps users know what kind of feedback is welcome.
- **Only update the `/quality` skill.** Rejected. The wording is a general
  agent-mediated UX pattern and belongs in the shared guide too.

## Trade-offs & risks

- Casual wording could become too vague if agents omit the target, boundary, or
  quality-log decision. The template keeps those elements visible when relevant.
- The `so that` clause can sound awkward for tiny edits. The guide allows
  equivalent value-prop wording when `so that` does not fit naturally.

## Open questions

None.
