---
type: Functional Specification
title: Durable spec rationale - functional spec
description: What the three contributor guides must require and teach so durable rationale lives in the spec, both big-picture and per-requirement.
tags: [specs, docs, contributing]
timestamp: 2026-06-18T00:00:00Z
---

# Durable spec rationale - functional spec

Companion to [Durable spec rationale](../0025-durable-spec-rationale.md). This
spec states the _delta_ to three durable contributor guides; the
[design doc](design.md) covers _why this shape_.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background

A [Change Case](../../../docs/guides/work-with-change-cases.md) is archived once
it lands, and the enduring [`specs/`](../../../specs/index.md) bundle carries
the result forward. That hand-off has been lossy: the durable spec inherits the
requirement but not the _why_ — the case's motivation and the design doc's
rationale stay in [`archive/`](../) with the case. Editors then meet rules
stripped of
the failure-modes that produced them, and re-introduce fixed bugs by
"simplifying" rules that only looked arbitrary. The guides below are revised so
the spec itself holds durable rationale at two grains — one big-picture, one
per-requirement — and so a landing case deposits its _why_ there instead of in
the archive.

## Scope

Covered: the requirements the three contributor guides must state or teach — the
two-layer rationale convention, the annotation form, the litmus for when to
annotate, and the absorb-on-landing step. This spec governs the guides' content,
not the `specs/` bundle: existing specs adopt the convention as they are next
edited.

Deferred: any sweep that retrofits rationale into existing specs, any tooling
that detects missing annotations, and any `schema.md` change.

## Two-layer rationale convention

The [functional-spec guide](../../../docs/guides/write-functional-specs.md) **MUST**
establish two distinct, co-located homes for durable rationale in a spec:

1. **Background / Motivation** (spec-level). The guide **MUST** add this to a
   spec's shape: a short prose section near the top stating _why the capability
   exists_ — the problem or failure-mode it addresses and any spec-scale lessons.
   The guide **MUST** distinguish it from **Scope** (what is covered or deferred)
   and from the companion note (what the spec governs, plus the source-of-truth
   link), so the three are not confused for one another.

   > Rationale: a change's big-picture motivation had no durable home in the
   > spec and died in the archive; Scope and the companion note answer _what_,
   > not _why_, so the _why_ needs its own named section.

2. **Per-requirement annotation** (requirement-level). The guide **MUST**
   establish that an individual requirement **MAY** carry a subordinate rationale
   annotation directly beneath it.

The two grains **MUST NOT** restate each other: Background carries the
spec-scale _why_; an annotation carries one requirement's _why_. The guide
**MUST** direct authors to supersede stale rationale rather than let it accrete,
so requirements stay skimmable.

> Rationale: two homes invite saying the same thing twice; without a
> say-it-once rule the duplication is the predictable failure-mode, and stale
> rationale left to pile up is what buried the _what_ under the _why_ before.

## Annotation form

The guide **MUST** specify the annotation's form precisely enough to apply
consistently:

- It **MUST** be a blockquote led by `Rationale:` (the terser `Why:` is
  acceptable).
- It **SHOULD** be one or two sentences.
- It **MAY** cite the originating change id (for example, `— 0012`) as
  provenance.
- The requirement **MUST** remain the lead, testable sentence; the annotation is
  subordinate and **MUST NOT** wrap around or bury it.

## Litmus for annotating

The guide **MUST** give the litmus for _when_ a requirement earns an annotation:
annotate when a future editor would otherwise repeat a mistake or be misled.
Dead-end alternatives and full decision records **MUST** stay in the (archived)
design doc; only durable intent and lessons are promoted into the spec.

> Rationale: without a bar, annotations become a running commentary on every
> rule; the design doc already holds the full decision record, so the spec
> should promote only what a future editor needs to avoid a regression.

## Absorb the why on landing

The [change-cases guide](../../../docs/guides/work-with-change-cases.md) **MUST**
extend the existing "the enduring artifacts absorb [the delta]" account so that
absorbing a case includes its **enduring _why_**, not only its functional
delta. When a case updates a durable spec, its motivation and its design doc's
durable rationale **SHOULD** be promoted into that spec's **Background /
Motivation** and per-requirement annotations. Durable specs **MAY** be edited at
any time (with or without a change), so this promotion is encouraged whenever a
spec is touched rather than gated on the **Before setting In-Review** step.

The [design-doc guide](../../../docs/guides/write-design-docs.md) **MUST** record
the other side of that promotion: the design doc's enduring rationale is lifted
into the functional spec when the Change Case lands, while the design doc remains the
fuller, archived record of alternatives and trade-offs.

## Dogfooding

This spec **MUST** itself follow the convention it specifies: it carries a
**Background** section above, and at least one requirement above carries a
`Rationale:` annotation. A reviewer **MUST** be able to read this file as a
worked example of the shape the guides now require.
