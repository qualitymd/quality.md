---
type: Change Case
title: Setup review gate and discovery trim
description: Make /quality setup stop for an explicit final recap response before authoring, trim per-question pedagogy to purpose/context only, and revise the setup question set.
status: Done
tags: [skill, setup, ux, pedagogy]
timestamp: 2026-06-23T00:00:00Z
---

# Setup review gate and discovery trim

A **Change Case** to make `/quality setup`'s final recap a hard review gate and
to reduce discovery-question teaching copy to the context and purpose of each
question. Setup should still teach why each dimension matters, but it should not
repeat "how to change it later" guidance on every question. It also removes two
discovery questions that are better handled by agent inference or post-setup
next steps: modeling rigor and review posture. In their place, setup gains a
rating-scale confirmation question that teaches that Rating Levels are
configurable while recommending the standard four-level scale.

Detail lives in:

- [Functional spec](0069-setup-review-gate-and-pedagogy-trim/spec.md) - what
  the change must do.
- [Design doc](0069-setup-review-gate-and-pedagogy-trim/design.md) - how it is
  shaped, and why.

## Motivation

The teaching-first setup flow introduced in
[0067 - Setup discovery pedagogy](0067-setup-discovery-pedagogy.md)
made discovery more explicit, but two follow-up issues remain.

First, the final recap currently says to "invite" a last comment or correction
before authoring. In practice, that can read as optional ceremony, especially
after a structured question tool completes the last discovery question. Setup
needs a distinct review gate: discovery completion is not permission to write
`QUALITY.md`; the user must respond to the recap with either confirmation or
corrections.

Second, the per-question "How to change it later" copy adds repetitive overhead
to a flow that is already intentionally interactive. The useful teaching moment
is the purpose of each question: why the dimension matters and what it shapes in
`QUALITY.md`. The living-document message can be stated once instead of repeated
ten times.

Third, two discovery questions ask for workflow-shaping decisions too early.
Modeling rigor is an internal setup judgment the agent can infer from lifecycle,
risk tolerance, repository context, and available evidence. Review posture asks
about recurrence or loop expectations before the first useful model exists; that
belongs in setup closeout as a next-step option, not as discovery input.

At the same time, setup should make the Rating Scale visible. Rating Levels are
configurable model vocabulary, not baked into the QUALITY.md format, and users
should learn that early. The recommended `outstanding`, `target`, `minimum`,
`unacceptable` scale is still the right default for most first models:
`outstanding` names a stretch band where further ROI should be questioned,
`target` encourages an acceptable good-enough bar without demanding perfection,
`minimum` defines the floor that can be relied on but still needs improvement,
and `unacceptable` names below-floor quality.

## Scope

Covered:

- Require `/quality setup` to stop after discovery, present the final recap, and
  wait for a user response before authoring.
- Clarify that completing a structured question-tool page or final discovery
  answer does not satisfy the final review gate.
- Accept either explicit confirmation or comments/corrections as the review-gate
  response.
- Require corrections and cross-cutting comments from the review gate to be
  incorporated before authoring.
- Remove per-question "How to change it later" setup pedagogy while preserving
  authored purpose/context copy for every discovery question.
- Permit one general living-document note outside the individual question copy.
- Remove the modeling-rigor discovery question from the user-facing setup
  questions.
- Remove the review-posture discovery question from the user-facing setup
  questions.
- Add a rating-scale discovery question that explains Rating Levels are
  configurable and recommends the standard four-level scale.
- Preserve modeling rigor as an internal setup-brief inference when useful.
- Preserve review/loop guidance as setup closeout next-step routing, not
  discovery.
- Update the runtime skill and durable setup spec to match when implemented.

Deferred / non-goals:

- No QUALITY.md format change.
- No CLI or Go code change; setup remains skill-driven.
- No change to the remaining discovery questions' option sets, recommended
  defaults, confidence labels, or evidence notes.
- No restoration of accept-all-and-skip behavior; setup still presents every
  remaining discovery question every run.
- No requirement that the user provide a substantive final comment; explicit
  confirmation remains enough to proceed.

## Affected artifacts

### Code

- [ ] **None expected.** This is a skill/spec workflow change with no planned
      `qualitymd` CLI, Go package, command, or config change.

### Durable specs

- [ ] `specs/skills/quality-skill/workflows/setup.md` - make the final recap a
      hard review gate, revise discovery pedagogy requirements from
      purpose-plus-how-to-change-later to purpose/context only, and remove the
      modeling-rigor and review-posture discovery questions while adding a
      rating-scale confirmation question.
- [ ] `specs/log.md` and
      `specs/skills/quality-skill/workflows/log.md` - record durable setup-spec
      updates when implementation lands.

### Durable docs

- [ ] `docs/log.md` - record durable documentation updates only if public docs
      are changed.

### Bundled skill

- [ ] `skills/quality/workflows/setup.md` - remove per-question
      "How to change it later" copy, preserve question purpose/context, remove
      the modeling-rigor and review-posture discovery questions, add the
      rating-scale confirmation question, and add the explicit review gate before
      `Write QUALITY.md`.

### Release

- [ ] `CHANGELOG.md` - update the unreleased `/quality Skill` note so it
      describes the trimmed pedagogy and hard review gate.

## Children

- [Functional spec](0069-setup-review-gate-and-pedagogy-trim/spec.md) -
  required behavior for the setup review gate and discovery-pedagogy trim.
- [Design doc](0069-setup-review-gate-and-pedagogy-trim/design.md) - how the
  revised question set and review gate are shaped.

## Status

`Done`. Implemented skill-only with no CLI/Go change: setup discovery now asks
nine questions, removes modeling rigor and review posture as user-facing
questions, adds Rating Scale confirmation, trims per-question pedagogy to
purpose/context only, and treats the final recap as a hard review gate before
authoring. Runtime skill, durable setup spec, changelog, and spec logs are
updated. Archived after implementation.
