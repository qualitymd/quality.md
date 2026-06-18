---
type: Functional Specification
title: rating-scale soundness diagnostic — functional spec
description: The new meta-model Functionality requirement that assesses a subject model's rating scale and per-requirement criterion overrides for semantic soundness, and the body sync it carries.
tags: [diagnostics, meta-model]
timestamp: 2026-06-17T00:00:00Z
---

# rating-scale soundness diagnostic — functional spec

Companion to the
[Diagnose rating-scale soundness in the meta-model](../0009-rating-scale-diagnostic.md)
change. This spec states the delta the CLI's built-in
[quality meta-model](../../../internal/models/quality-meta-model.md)
must absorb so its Functionality factor assesses the rating scale's meaning, not
only its structure.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Scope

Covered: one new Functionality requirement on the meta-model, and the body sync
(Functionality summary and diagnostic coverage checklist) that keeps the prose
consistent with it.

Deferred: splitting the requirement into separate scale-shaping and
override-calibration checks; any change to
[`SPECIFICATION.md`](../../../SPECIFICATION.md) or the scaffold
[skeleton](../../../internal/scaffold/skeleton.md), which already hold the
rating-scale semantics the requirement traces to; any new Usability requirement;
and any wiring of the meta-model into code or tests.

## Requirements — the new requirement

- The meta-model **MUST** add one requirement under its **Functionality** factor
  that assesses whether a subject model's rating scale, and any per-requirement
  `ratings` overrides, are well-formed and meaningful — not merely structurally
  valid.
- Its assessment **MUST** cover: that levels are ordered best-to-worst; that each
  level's `description` fixes a distinct, coherent standing held fixed across the
  whole model; that the bands are separable for the subject, so an evaluator can
  tell where findings land; and that the acceptable floor sits where the model's
  needs and risks place it.
- Its assessment **MUST** require that, where a requirement supplies its own
  `ratings`, each override names a real level of the model's scale and replaces
  only that level's criterion — never the level's meaning, order, or title — and
  appears only where the shared criterion cannot express the gradient that
  matters.
- The requirement **MUST** trace to the rating-scale semantics already normative
  in [`SPECIFICATION.md`](../../../SPECIFICATION.md) (the **Rating Scale** and
  **Requirement** sections) and **MUST NOT** restate authoring guidance that is
  not an evaluable criterion.
- The requirement **MUST** be satisfied trivially by a model that inherits the
  suggested rating scale unchanged and declares no overrides, so that it raises
  findings only against real authoring choices rather than against every model.
- The change **SHOULD** express this as a single requirement, and **MAY** later
  split scale-shaping from override-calibration if the combined form proves
  overloaded.

## Requirements — the body sync

- The meta-model's Functionality summary **MUST** name rating-scale soundness
  among the concerns the factor covers.
- The diagnostic coverage checklist **MUST** list the new requirement.

## Done criteria

- The meta-model's **Functionality** factor carries the new *rating scale and any
  overrides are well-formed and meaningful* requirement, with an assessment
  meeting the requirements above.
- The Functionality summary and the diagnostic coverage checklist name the new
  concern.
- `qualitymd lint` reports the meta-model valid (no error-level findings).
- The change is moved through the lifecycle and archived per the
  [changes process](../../index.md#status-lifecycle).
