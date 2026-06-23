---
type: Design Doc
title: Setup review gate and discovery trim - design
description: Design rationale for the revised /quality setup discovery question set and final review gate.
tags: [skill, setup, ux, pedagogy]
timestamp: 2026-06-23T00:00:00Z
---

# Setup review gate and discovery trim - design

Companion to the
[functional spec](spec.md) for
[Setup review gate and discovery trim](../0069-setup-review-gate-and-pedagogy-trim.md).

## Context

This change refines the teaching-first setup workflow from
[0067](../0067-setup-discovery-pedagogy.md). The goal is to keep setup
educational while reducing repeated explanation and closing the gap where a
structured question response could be treated as permission to author
`QUALITY.md`.

## Approach

The setup question set stays fixed and authored in the workflow, but changes
shape:

- remove **Modeling rigor** as a user-facing question;
- remove **Review posture** as a user-facing question;
- add **Rating Scale** as a user-facing confirmation/calibration question.

Modeling rigor remains internal setup judgment. The agent can infer first-model
depth from lifecycle, risk tolerance, repository shape, evidence availability,
and the body it authors. Review cadence remains a closeout next-step option, not
pre-authoring discovery.

The Rating Scale question teaches the format boundary: Rating Levels are
configurable model vocabulary, not baked into QUALITY.md. Setup still recommends
the standard four-level scale and asks for confirmation rather than inviting the
user to invent scale names. If the user rejects the recommendation, setup uses a
simple alternate scale only when context clearly supports it, such as a pass/fail
gate; otherwise it keeps the recommendation and records the uncertainty in the
model body.

The final recap becomes a review gate. Discovery completion and review-gate
confirmation are separate states: after the last discovery answer, setup recaps
the answers and waits for a user response before writing or editing
`QUALITY.md`.

## Alternatives

Keep modeling rigor as a question: rejected because it asks users to reason about
the model's implementation detail instead of the evaluated thing's needs and
risks. The agent can make a better initial call from context.

Keep review posture as a question: rejected because recurrence and loop design
are easier to choose after a first model exists. Asking during discovery makes
setup feel like it may create schedules or automation even when it does not.

Ask users to design a custom Rating Scale: rejected because setup should not ask
users to invent Rating Level names cold. The pedagogical point is that the scale
is configurable; the practical default is still the recommended scale.

Treat the last structured question response as confirmation: rejected because it
collapses discovery and authorization. The recap is the only place the user sees
the full answer set together.

## Trade-offs and risks

Removing two questions reduces setup overhead, but it also means the body must
preserve inferred modeling rigor only when it materially shapes the model.

Adding the Rating Scale question keeps the interaction count at nine rather than
eight. The added question earns its place because it teaches a format concept
users otherwise might mistake as hard-coded.

The review gate adds one explicit stop before authoring. That is intentional:
setup writes `QUALITY.md`, so the user should see and confirm the consolidated
answers first.

## Open questions

None.
