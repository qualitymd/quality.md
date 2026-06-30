---
type: Runtime Guide
title: Getting Started with QUALITY.md
description: First-run iteration guidance after setup leaves a valid QUALITY.md with important model gaps.
---

# Getting started with QUALITY.md

Use this guide after setup leaves a valid `QUALITY.md` with important model gaps,
or when a user asks how to keep iterating on the first useful model. The goal is
to turn the current file into a useful quality model: small enough to finish,
specific enough to evaluate, and clear enough to improve later.

Read [Authoring QUALITY.md](authoring.md) first. This guide assumes that
guidance and focuses on the first-run process and the desired outcome of each
step.

## Starting point

Setup should leave you with a structurally valid file and the first pass at the
root area, domain, lifecycle, risk tolerance, modeling rigor, collaboration
context, stakeholder needs, risks, missing context, and model shape. Treat those
as starting assumptions to review, not as settled truth.

Before editing, run:

```sh
qualitymd lint QUALITY.md
```

If lint fails, fix validity first. Do not build a model on an invalid skeleton.

## First pass

### Fill the body first

Outcome: the Markdown body explains the root area, scope, needs, and risks — with each
section's unknowns and open questions — well enough to justify and evaluate the
first quality model.

Use authoring guidance: [Markdown body](authoring/body.md).

Check before moving on:

- Overview names the real root area, dependents, and why quality matters.
- Scope names the domain, what is included, excluded, and where the model
  boundary sits.
- Needs names primary user, maintainer/collaborator, and other stakeholder
  outcomes the model serves.
- Risks name the important failure modes.
- The body captures lifecycle, risk tolerance, modeling rigor, collaboration
  context, and any other context that directly shapes the model.
- Each section records its own unknowns and open questions, or "none known".
- Material support is cited when agent-accessible; important support that is not
  agent-accessible is named as an unknown or open question.

### Confirm the rating scale

Outcome: the rating levels can distinguish `unacceptable`, `minimum`, `target`,
and `outstanding` for this root area.

Use authoring guidance: [Rating Scale](authoring/rating-scale.md).

Check before moving on:

- The default scale still fits, including the `🟢 Outstanding`, `🔵 Target`,
  `🟡 Minimum`, and `🔴 Unacceptable` display titles, or the body explains why a
  different scale or plain-title style is needed.
- The level criteria are clear enough to rate future findings consistently.

### Name the root area

Outcome: the root `title`, body, file location, and root `source` describe the
same evaluated root area.

Use authoring guidance: [Authoring QUALITY.md](authoring.md) and
[Model structure](authoring/model-structure.md).

Check before moving on:

- The root `title` names the entity described by the body.
- The root `source` stays implicit unless the model needs to narrow or relocate
  the default scope.

### Pick two to five factors

Outcome: the initial factors cover the most important Needs and Risks without
overlapping heavily.

Use authoring guidance: [Factors](authoring/factors.md).

Check before moving on:

- Each major body need or risk has a quality lens.
- Each factor has a required `title`.
- Each factor description explains the lens in terms of this root area.

### Write assessable requirements

Outcome: each initial requirement can produce evidence, findings, and a rating.

Use authoring guidance: [Requirements](authoring/requirements.md).

Check before moving on:

- a concrete expectation as the map key;
- exactly one `assessment`;
- evidence a future evaluator can actually inspect;
- enough specificity to distinguish at least adjacent rating levels.

Prefer assessable requirements over aspirational ones: a requirement you cannot
inspect adds no signal. This guards requirement quality, not coverage — model the
constituents the entity has (see
[model structure](authoring/model-structure.md#cover-the-domains-constituent-kinds)),
and keep each requirement under them assessable rather than padding with expectations
you can't judge.

## First validation

After the first pass, run:

```sh
qualitymd lint QUALITY.md
qualitymd status QUALITY.md --json
```

Lint proves structure. Status helps route the next workflow. Neither proves the
model is useful enough to evaluate; read the file once more against these
questions:

- Does every requirement name evidence or criteria that can be inspected?
- Can an evaluator tell `target` from `minimum` for the important requirements?
- Does the body explain why these factors and requirements matter?
- Is the body current, grounded, and agent-accessible enough to evaluate the
  model's quality?
- Are each section's unknowns and open questions captured as context without
  turning them into ratings?

If the answer is no, revise the model before running an evaluation.

## Next workflow

When the file is valid and the first model is useful enough, choose the next
workflow:

```text
**Recommended next:** continue iterating on `QUALITY.md`

**Options**
1. Continue iterating on `QUALITY.md`.
2. Run `/quality evaluate`.
3. Stop here.

**Answer:** Reply `1`, `2`, or `3`.
```

If the validation or status output shows model gaps, continue model authoring
before evaluating. If tooling is missing or stale, run `/quality update` first.
This guide does not create external issues, configure integrations, or configure
automation.
