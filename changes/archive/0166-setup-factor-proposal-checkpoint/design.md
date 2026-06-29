---
type: Design Doc
title: Setup Factor Proposal Checkpoint — design doc
description: Design for adding factor desiderata teaching and a setup factor proposal checkpoint.
tags: [skill, quality, setup, factors]
timestamp: 2026-06-29T00:00:00Z
---

# Setup Factor Proposal Checkpoint — design doc

Design for
[Setup Factor Proposal Checkpoint](../0166-setup-factor-proposal-checkpoint.md)
and its [functional spec](spec.md).

## Context

This is a runtime guidance and durable-spec alignment change. No Go code or CLI
behavior changes. The main workflow file already has the right large-scale shape:
setup reads context, asks discovery questions, presents a human context
checkpoint, stops for review, writes, lints, and reports important gaps. The
change inserts one focused checkpoint between human context and final review.

## Approach

Add the canonical factor desiderata to the runtime factor authoring guide and its
durable guide spec. Keep the runtime guide prose practical and the durable spec
contractual: the runtime guide explains the terms, while the spec requires that
the guide teach them.

Extend the setup workflow in four places:

- the phase roadmap and setup brief, so factor-set quality and factor rationales
  are working context rather than an afterthought;
- a new `Factor Proposal Checkpoint` section, with the exact user-facing shape;
- final review and model authoring, so the reviewed proposal binds the generated
  model;
- verify/close gap inspection, so setup reports failure to apply the desiderata.

Update the durable setup workflow spec in the same structure, using normative
requirements rather than operator copy. Update the parent skill spec only as a
summary pointer.

Update the README with a short public-facing explanation and link to the detailed
factor guide.

## Spec response

The runtime and durable factor guide changes satisfy the desiderata requirements.
The runtime setup checkpoint satisfies the teaching and user-feedback
requirements without asking the user to design model elements cold. The final
recap, authoring, and close-gap additions ensure the checkpoint affects the
resulting `QUALITY.md` rather than being only educational copy.

## Alternatives

One alternative was to add a sixth discovery question asking "what factors do you
want?" That was rejected because it asks users to design the model cold and would
produce worse feedback than reviewing a concrete proposal.

Another alternative was to put the full desiderata only in README. That was
rejected because setup agents need operational guidance at runtime and the README
should stay introductory.

## Trade-offs & risks

The checkpoint adds one more setup interaction. That is deliberate: setup runs
rarely, and factor selection is important enough to teach and review. To keep the
cost bounded, the checkpoint uses one table and targeted correction categories
instead of a new questionnaire.

The `light` / `normal` / `deep` labels could be mistaken for importance labels.
The guidance explicitly defines them as initial requirement and assessment rigor,
not importance alone.

## Open questions

- A future change can define a fuller factor-selection method using these
  desiderata. This case only installs the qualities and the setup checkpoint.
