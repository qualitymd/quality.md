---
type: Functional Specification
title: Setup review gate and discovery trim - functional spec
description: Requirements for making /quality setup wait for a final recap response, trimming per-question pedagogy to purpose/context only, and revising the setup question set.
tags: [skill, setup, ux, pedagogy]
timestamp: 2026-06-23T00:00:00Z
---

# Setup review gate and discovery trim - functional spec

Companion to the
[Setup review gate and discovery trim](../0069-setup-review-gate-and-pedagogy-trim.md).
This spec defines the delta from the current teaching-first setup flow to a
tighter question-pedagogy contract, a shorter discovery question set, and an
explicit final review gate.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background

[0067 - Setup discovery pedagogy](../0067-setup-discovery-pedagogy.md)
made setup teaching-first by adding authored per-question copy, requiring every
question every run, and adding a final recap before authoring. The follow-up
need is to keep the useful part of that teaching while reducing repetition and
closing an authoring loophole.

The useful per-question teaching is the question's purpose: what judgment the
user is making and what the answer shapes in `QUALITY.md`. Repeating
"how to change it later" on every question makes the flow heavier without adding
much decision value. Modeling rigor is also too meta for first-run discovery:
the agent can infer an initial model depth from lifecycle, risk tolerance,
repository context, and available evidence. Review posture is better handled as
post-setup next-step routing than as a question before the first model exists.

Rating Scale choice, however, should be visible as a confirmation/calibration
question. Rating Levels are configurable model vocabulary, not baked into the
QUALITY.md format. Setup should recommend the four-level
`outstanding` > `target` > `minimum` > `unacceptable` scale because it gives a
useful quality gradient: `outstanding` is a stretch band where further ROI should
be questioned, `target` is the expected good-enough bar, `minimum` is the
acceptable floor that still warrants improvement, and `unacceptable` is below
the floor.

Separately, the final recap must be a real gate. A response to the last discovery
question, including through a structured question tool, completes discovery only;
it must not imply permission to write the model.

## Requirements

### Discovery question set

Setup MUST remove the modeling-rigor question from user-facing discovery.

Setup MAY retain modeling rigor as an internal setup-brief inference used to
shape the first model.

Setup MUST remove the review-posture question from user-facing discovery.

Setup MAY keep review cadence, recurrence, and quality-loop options in setup
closeout as next-step routing.

Setup MUST add a rating-scale discovery question that asks whether to use the
recommended four-level scale:
`outstanding`, `target`, `minimum`, `unacceptable`.

The rating-scale question MUST explain that Rating Levels are configurable in
`QUALITY.md` and are not baked into the format.

The rating-scale question MUST recommend the standard four-level scale for most
first models.

The rating-scale question's purpose/context copy MUST explain the role of each
recommended level:

- `outstanding` names a stretch band where further investment may need ROI
  justification.
- `target` names the expected good-enough bar without demanding perfection.
- `minimum` names the acceptable floor that can be relied on but still warrants
  improvement.
- `unacceptable` names quality below the floor.

Setup MUST NOT ask the user to invent custom Rating Level names during discovery.

When the user rejects the recommended scale, setup MAY choose a simple alternate
scale only when project context clearly supports it, such as a pass/fail gate;
otherwise setup SHOULD use the recommended scale and record the scale decision as
an open question or assumption in the model body.

### Discovery pedagogy

Each remaining setup discovery question MUST include authored purpose/context
copy explaining why the dimension matters and what the answer shapes in
`QUALITY.md`.

Setup MUST NOT include per-question "How to change it later" copy or equivalent
per-question lifecycle guidance.

Setup MAY state once outside the individual question copy that `QUALITY.md` is a
living document and that setup answers can be revised later.

Setup MUST continue to present every remaining discovery question on every run,
including questions whose inferred default has high confidence.

Setup MUST continue to carry each question's recommended default, confidence
label, and evidence note. This change does not alter option sets, defaulting
behavior, or confidence vocabulary for discovery questions except for removing
the modeling-rigor and review-posture questions and adding the rating-scale
question.

Setup MUST NOT restore an accept-all-and-skip shortcut. A per-question fast
confirm and a user-requested show-all-at-once presentation remain allowed.

### Final review gate

After all discovery questions are answered, setup MUST stop before authoring and
present a consolidated final recap that lists every asked discovery question with
its final answer.

Setup MUST wait for a user response to the final recap before writing or editing
`QUALITY.md`.

Completing the final discovery question or a structured question-tool page MUST
NOT satisfy the review gate.

The review-gate response MAY be a correction, a cross-cutting comment, or an
explicit confirmation such as "looks good", "continue", "write it", or an
equivalent phrase. Setup MUST NOT require the user to provide a substantive
comment.

When the user provides corrections or cross-cutting comments at the review gate,
setup MUST incorporate them into the working setup answers before authoring.

The final recap MUST supplement discovery, not replace it. Setup MUST NOT
collapse the discovery questions into the review gate alone.

## Acceptance criteria

- The runtime setup workflow no longer contains live per-question
  "How to change it later" lines.
- The user-facing setup discovery list removes the modeling-rigor question.
- The user-facing setup discovery list removes the review-posture question.
- The user-facing setup discovery list adds a rating-scale question that
  recommends `outstanding`, `target`, `minimum`, `unacceptable`.
- The rating-scale question explains that Rating Levels are configurable and not
  baked into QUALITY.md.
- The rating-scale question does not ask the user to invent custom Rating Level
  names during setup.
- Each remaining setup discovery question still contains authored
  purpose/context copy
  explaining what the answer shapes.
- The runtime setup workflow contains an explicit review gate before
  `Write QUALITY.md`.
- The review gate says to wait for a user response before writing or editing
  `QUALITY.md`.
- The workflow explicitly says structured question-tool completion does not
  satisfy the review gate.
- The durable setup spec mirrors the runtime behavior.
- The unreleased changelog entry reflects the trimmed pedagogy and hard review
  gate.
- `rg "How to change it later|how-to-change-later|change the answer later"`
  finds no live setup-pedagogy requirement text except historical archived
  change records or change-case discussion of the removed behavior.
- Markdown formatting/checks pass.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/workflows/setup.md` - remove the modeling-rigor
  and review-posture discovery questions, add the rating-scale question, update
  the discovery pedagogy requirement to purpose/context copy only, remove
  per-question how-to-change-later requirements, and make the final recap an
  explicit review gate that waits for a user response before authoring.

### To rename

None

### To delete

None
