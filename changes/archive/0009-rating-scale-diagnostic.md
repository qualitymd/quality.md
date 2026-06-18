---
type: Change
title: Diagnose rating-scale soundness in the meta-model
description: Add a meta-model Functionality requirement that assesses a subject model's rating scale and per-requirement criterion overrides for semantic soundness, not just structural shape.
status: Done
tags: [diagnostics, meta-model]
timestamp: 2026-06-17T00:00:00Z
---

# Diagnose rating-scale soundness in the meta-model

A **Change** to the CLI's built-in
[quality meta-model](../../internal/models/quality-meta-model.md)
— the model it uses to judge whether a project's `QUALITY.md` is a good quality
model. It adds one Functionality requirement so the rating scale, the instrument
that turns assessments into verdicts, is assessed for *meaning* and not only
structure. Detail lives in the child:

- [Functional spec](0009-rating-scale-diagnostic/spec.md) — what the change must do.

A design doc is omitted: the change adds one requirement to a content artifact
and needs no technical-design discussion.

## Motivation

The meta-model lavishes attention on a model's prose — six separate requirements
cover the Overview, Scope, Needs, Risks, Factors, and Known gaps body sections —
yet the rating scale, which is what actually turns a requirement's findings into
a verdict, is assessed almost in passing. Today only two requirements touch it:
*the model passes structural lint* checks that the scale **parses** ("the rating
scale is well-shaped"), and *the model correctly applies the QUALITY.md format
spec* mentions "rating criteria are used for their intended purposes" inside a
longer conformance list. Neither asks whether the scale's levels carry distinct,
coherent meaning for the subject, whether the acceptable floor sits where the
model's needs and risks put it, or whether a requirement's own `ratings`
overrides are sound.

That gap matters most where the format gives an author the most rope: a
per-requirement `ratings` override (see
[`SPECIFICATION.md`](../../SPECIFICATION.md) § Requirement) replaces a level's
criterion for one requirement — e.g. a measured `p99 ≤ 150 ms` threshold. A
miscalibrated threshold, or an override that quietly redefines what a level
*means* — which the spec forbids, since overrides change the criterion alone —
produces wrong verdicts while passing every check the meta-model makes today. The
rating scale is the heart of the gradient; the meta-model should be able to judge
it.

## Scope

Covered:

- One new Functionality requirement, *the rating scale and any overrides are
  well-formed and meaningful*, assessing level ordering, each level's description
  fixing a fixed model-wide standing, band separability for the subject, floor
  placement against the model's needs and risks, and `ratings` overrides that
  change only a level's criterion and only where the shared criterion cannot
  express the needed gradient.
- The matching body sync: the Functionality summary and the diagnostic coverage
  checklist name the new concern.

Deferred:

- Splitting the requirement into separate scale-shaping and override-calibration
  checks — left as a single requirement, to split only if it proves overloaded.
- Any change to [`SPECIFICATION.md`](../../SPECIFICATION.md) or the scaffold
  [skeleton](../../internal/scaffold/skeleton.md): both already define the
  rating-scale semantics this requirement traces to, and neither changes.
- Any new Usability requirement, and any wiring of the meta-model into code or
  tests.

## Affected specs & docs

None — a deliberate "no durable changes." The meta-model is the artifact this
change edits, not a durable spec or doc that describes it, and no `specs/`,
[`SPECIFICATION.md`](../../SPECIFICATION.md), `README.md`, or `docs/` content
enumerates the meta-model's requirement set, so none drifts. The new requirement
traces to the rating-scale semantics already normative in
[`SPECIFICATION.md`](../../SPECIFICATION.md) (the **Rating Scale** and
**Requirement** sections), which this change leaves unchanged.

## Status

`Done`. Implemented and archived after adding the *rating scale and any overrides
are well-formed and meaningful* requirement to the meta-model's Functionality
factor, syncing the Functionality summary and the diagnostic coverage checklist,
and confirming `qualitymd lint` still reports the model valid.
