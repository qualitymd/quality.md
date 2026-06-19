---
type: Functional Specification
title: Sharpen assessment references and traceability — functional spec
description: Requirements for describing an assessment as inline or a reference, reserving "source" for Target.source, and making the model's traceability graph an authoring concern.
tags: [specification, terminology, guide]
timestamp: 2026-06-19T00:00:00Z
---

# Sharpen assessment references and traceability — functional spec

Companion to the
[Sharpen assessment references and traceability](../0029-sharpen-assessment-references.md)
change case. This spec states *what* the change must do. The format it sharpens
is defined by [`SPECIFICATION.md`](../../../SPECIFICATION.md), the source of truth.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 when, and only when, they appear in all
capitals.

## Background / Motivation

A requirement's `assessment` is the means of judging a target's source against
that requirement. Those means can be written in place or can live in another
artifact the assessment points at. The durable artifacts never name this
inline-vs-reference distinction, so two failure modes recur: criteria get copied
into requirements where they drift from their origin, and a single assessable
claim gets split into several per-factor requirements that each re-cite the same
artifact.

An entity an assessment references is usually itself a target — a spec is judged
for quality *and* is the criteria its implementation is judged against. The edges
between these targets form a traceability graph, but nothing helps an author make
them visible. Compounding this, "source" is already bound to `Target.source`, so
naming the referenced criteria a second "source" overloads the term.

The fix reserves "source" for `Target.source`, frames an assessment as inline or
a reference, and uses "reference" for the requirement→entity edge — the same
terminology backbone [0028](../0028-require-characterized-requirements.md)
established for requirement→factor references.

## Scope

Covered: how the format spec, the authoring guide, and the scaffold describe an
`assessment` and the relationships it creates.

Deferred: a first-class machine-readable assessment→target link type — the edge
rides on shared canonical references for now.

Non-goals: no change to the `assessment` **schema** (it stays a single non-empty
scalar) and no new lint rule — the inline-vs-reference distinction is not
mechanically enforced. Re-modeling this repo's own `QUALITY.md` to capture its
doc/spec/code graph is out of scope.

## Requirements

`SPECIFICATION.md` **MUST** describe an `assessment` as either stating the means
of assessing inline or referencing an entity that defines those means, without
changing the rule that an `assessment` is a single non-empty scalar.

> Rationale: Authors copy criteria into requirements because nothing names
> referencing as the alternative; duplicated criteria drift from their origin. —
> 0029

`SPECIFICATION.md` **SHOULD** note, non-normatively, that a referenced entity may
itself be a Target in the Model and that referencing it by the same selector used
as that Target's `source` makes the dependency traceable without a distinct link
type.

The authoring guide **MUST** reserve "source" for `Target.source` and **MUST NOT**
use "source" as the term for the entity an assessment references.

> Rationale: The model already binds "source" to the entities a target
> evaluates; a second "source" for the referenced criteria reads as scope, the
> opposite role. — 0029

The authoring guide **MUST** describe an `assessment` as inline or a reference to
an entity, and **MUST** use "reference" for the requirement→entity edge,
consistent with the requirement→factor reference terminology in
[0028](../0028-require-characterized-requirements.md).

The authoring guide **MUST** guide authors to treat an entity as both a target
and an assessment reference by role, and **MUST** present the model's traceability
graph — targets plus the assessment references between them — as an authoring
concern.

The authoring guide **MUST** guide authors to size a requirement to one
assessable claim, connecting a cross-lens claim to multiple factors rather than
duplicating one assessment across per-factor requirements, while keeping
genuinely independent claims that share a reference separate.

> Rationale: A claim sliced one-per-factor re-assesses the same entity
> repeatedly and fragments a single judgment; but many independent claims may
> legitimately draw on one rich entity, so the split is by claim, not by source. —
> 0029

The scaffold guidance **SHOULD** describe an inline-or-reference assessment using
"reference" rather than "defer to", and **SHOULD** note that a referenced entity
can be a target in its own right.

## Durable spec changes

### To add

None

### To modify

- [`SPECIFICATION.md`](../../../SPECIFICATION.md) — extend the **Assessment**
  terminology and the Requirement section to describe an `assessment` as inline
  or a reference, and add the non-normative traceability note (per the
  `SPECIFICATION.md` requirements above).

### To delete

None
