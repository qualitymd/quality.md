# Getting Started with QUALITY.md

Use this guide after `qualitymd init` has created a valid skeleton. The goal is
to replace the placeholders with the first useful quality model: small enough to
finish, specific enough to evaluate, and clear enough to improve later.

Read [Authoring QUALITY.md](authoring.md) first. This guide assumes that
guidance and focuses on the first-run process and the desired outcome of each
step.

## Starting Point

Setup leaves you with a structurally valid file. It does not know the project,
the decisions quality work should support, or the evidence that will make a
requirement assessable.

Before editing, run:

```sh
qualitymd lint QUALITY.md
```

If lint fails, fix validity first. Do not build a model on an invalid skeleton.

## First Pass

### Fill the Body First

Outcome: the Markdown body explains subject, scope, needs, and risks — with each
section's unknowns and open questions — well enough to justify the first quality
model.

Use authoring guidance: [The Markdown body](authoring.md#the-markdown-body).

Check before moving on:

- Overview names the real subject, dependents, and why quality matters.
- Scope names what is included, excluded, and where the model boundary sits.
- Needs names the outcomes and stakeholders the model serves.
- Risks name the important failure modes.
- Each section records its own unknowns and open questions, or "none known".

### Confirm the Rating Scale

Outcome: the rating levels can distinguish `unacceptable`, `minimum`, `target`,
and `outstanding` for this subject.

Use authoring guidance: [Rating Scale](authoring.md#rating-scale).

Check before moving on:

- The default scale still fits, or the body explains why a different scale is
  needed.
- The level criteria are clear enough to rate future findings consistently.

### Name the Subject

Outcome: the root `title`, body, file location, and root `source` describe the
same evaluated subject.

Use authoring guidance: [Quality Model](authoring.md#quality-model) and
[Target](authoring.md#target).

Check before moving on:

- The root `title` names the entity described by the body.
- The root `source` stays implicit unless the model needs to narrow or relocate
  the default scope.

### Pick Two to Five Factors

Outcome: the initial Factors cover the most important Needs and Risks without
overlapping heavily.

Use authoring guidance: [Factor](authoring.md#factor).

Check before moving on:

- Each major body need or risk has a quality lens.
- Each Factor has a required `title`.
- Each Factor description explains the lens in terms of this subject.

### Write Assessable Requirements

Outcome: each initial Requirement can produce evidence, findings, and a rating.

Use authoring guidance: [Requirement](authoring.md#requirement).

Check before moving on:

- a concrete expectation as the map key;
- exactly one `assessment`;
- evidence a future evaluator can actually inspect;
- enough specificity to distinguish at least adjacent rating levels.

Prefer a small first model with assessable Requirements over a broad model full
of aspirations.

## First Validation

After the first pass, run:

```sh
qualitymd lint QUALITY.md
qualitymd status QUALITY.md --json
```

Lint proves structure. Status helps route the next workflow. Neither proves the
model is useful enough to evaluate; read the file once more against these
questions:

- Does every Requirement name evidence or criteria that can be inspected?
- Can an evaluator tell `target` from `minimum` for the important Requirements?
- Does the body explain why these Factors and Requirements matter?
- Are each section's unknowns and open questions captured as context without
  turning them into ratings?

If the answer is no, revise the model before running an evaluation.

## Next Workflow

When the file is valid and the first model is useful enough, run:

```text
/quality wizard
```

Wizard should route to model review, whole-subject evaluation, scoped
evaluation, or another concrete next step based on the current status.
