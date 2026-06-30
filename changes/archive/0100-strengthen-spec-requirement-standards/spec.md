---
type: Functional Specification
title: Strengthen spec requirement standards — functional spec
description: Requirements for the wording the functional-spec and change-case guides must gain to close 29148 gaps and add an optional EARS template.
tags: [process, specs, change-cases, requirements]
timestamp: 2026-06-26T00:00:00Z
---

# Strengthen spec requirement standards — functional spec

Companion to the
[Strengthen spec requirement standards (29148 + EARS)](../0100-strengthen-spec-requirement-standards.md)
change case. This spec states _what_ the change must do. There is no design doc;
each requirement carries its own rationale.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174.

## Background / Motivation

Our [requirement quality bar](../../docs/guides/write-functional-specs.md#requirement-quality-bar)
already re-derives the nine individual-requirement characteristics of
ISO/IEC/IEEE 29148:2018 §5.2.5. The remaining gaps are at the _set_ level, in the
_validation_ (vs verification) question, in two missing document elements
(assumptions, reference classification), in one unnamed bar item (unambiguous),
and in the absence of a statement _template_ — the only §5.2.5 characteristic
("Conforming") we lack. EARS supplies that template in a lightweight, optional
form. This change adds the missing guidance to the two authoring guides without
disturbing their prose-first, palette-not-checklist, BCP 14 posture.

This spec governs _what the guides must say_; it does not itself change tool
behavior. Its own conformance target is the very guide it edits — it is written
to pass that guide's bar.

## Scope

Covered: wording added to
[`docs/guides/write-functional-specs.md`](../../docs/guides/write-functional-specs.md)
and
[`docs/guides/work-with-change-cases.md`](../../docs/guides/work-with-change-cases.md).

Not covered: `SPECIFICATION.md`, the QUALITY.md format's `requirements`
guidance, any code/CLI/skill/rating/evaluation/report/format-schema behavior, and
the heavier 29148 apparatus (immutable requirement IDs, owner/version/priority/
risk/difficulty attributes, the BRS/StRS/SyRS/SRS hierarchy, a formal
traceability matrix). EARS stays optional. See the change case for the full
non-goal list.

## Requirements

### Functional-spec guide

- The functional-spec guide **MUST** add set-level requirement guidance distinct
  from the per-requirement [quality bar](../../docs/guides/write-functional-specs.md#requirement-quality-bar):
  a spec's requirements **MUST** be reviewable as a _set_ for consistency (no
  requirement conflicts with or overlaps another; one term means one thing
  throughout), completeness (no unresolved TBD/TBR markers except those captured
  under Open questions or Deferred), and validatability (satisfying the set would
  achieve the stated need).

  > Rationale: a set of individually well-formed requirements can still conflict,
  > duplicate, drift in terminology, or collectively miss the need. 29148 §5.2.6
  > makes the set its own object of review; our bar was per-requirement only. — 0100

- The guide **MUST** add an optional **Assumptions & dependencies** element to
  the [Shape](../../docs/guides/write-functional-specs.md#shape) palette for
  external facts that, if they change, would invalidate requirements, and **MUST**
  distinguish it from Scope (what is in/out) and Background (the _why_).

  > Rationale: 29148 §9.6.8 carves out assumptions as their own item. For an
  > agent-driven tool resting on shifting substrate (run-folder layouts, "the CLI
  > never calls a model"), a named home lets a broken assumption flag the
  > dependent requirements instead of silently invalidating them. — 0100

- The guide **MUST** instruct authors to distinguish **normative** references
  (binding sources of truth a spec defers to, e.g. `SPECIFICATION.md`) from
  **informational** ones (context, related work).

  > Rationale: 29148 §9.2.4 splits references into compliance vs guidance. The
  > distinction sharpens our "say it once / link to the source of truth" rule by
  > marking which links _bind_. — 0100

- The guide's [quality bar](../../docs/guides/write-functional-specs.md#requirement-quality-bar)
  **MUST** include an explicit **Unambiguous** item (the requirement admits one
  interpretation; it avoids vague predicates), and the
  [Requirements](../../docs/guides/write-functional-specs.md#requirements)
  guidance **SHOULD** recommend subject-explicit, active-voice statements.

  > Rationale: 29148 §5.2.5 lists _unambiguous_ as its own characteristic and
  > §5.2.4 calls for active voice and an explicit subject. We had unambiguity only
  > implicitly, inside the weak-verb example. — 0100

- The guide **MUST** record that our BCP 14 keyword convention (`MUST`/`SHOULD`/
  `MAY`) is a deliberate divergence from the `shall`-based convention 29148 §5.2.4
  and EARS use, and that the convention is retained.

  > Rationale: 29148 and EARS center `shall` and advise avoiding "must"; we
  > center BCP 14. Naming the divergence as a choice stops a future reader from
  > "correcting" it back. — 0100

- The guide **MUST** add an **optional** requirement-statement template based on
  EARS, presented as recommended-not-required and consistent with the
  palette-not-checklist posture. It **MUST** cover the EARS patterns — ubiquitous,
  state-driven (`While`), event-driven (`When`), optional-feature (`Where`),
  unwanted-behaviour (`If`/`Then`), and complex combinations — adapted to keep the
  BCP 14 keyword in the response (e.g. `When <trigger>, <surface> MUST
<observable result>`). It **MUST** tie the `When` form to the existing
  **Bounded condition** bar item and the `If`/`Then` form to the **Divergence
  handled** bar item and the "an unspecified case is a decision delegated"
  convention.
  > Rationale: EARS is the lightweight statement template for the one §5.2.5
  > characteristic ("Conforming") we lacked. Its unwanted-behaviour `If/Then` form
  > gives explicit syntax for the edge/error cases our bar already insists be
  > decided. Optional, because EARS fits behavioral/trigger requirements better
  > than meta/structural ones. — 0100

### Change-case guide

- The change-case guide **MUST** add, at the Draft→Design boundary, a
  **validation** check distinct from per-requirement verification: before leaving
  Draft, confirm that satisfying the full requirement set would achieve the
  case's motivation — no more, no less. It **SHOULD** cross-reference the
  functional-spec guide's set-level guidance.
  > Rationale: our lifecycle gates on the per-requirement verification path
  > ("built it right") but never asks the validation question ("built the right
  > thing"). The change case is where the motivation lives, so the Draft→Design
  > boundary is the natural home. 29148 distinguishes verification from
  > validation throughout. — 0100

### Provenance

- Where these rules are added, the guides **MUST** carry brief rationale (a
  sentence or short annotation) naming the 29148/EARS origin, per the guides' own
  [two-whys](../../docs/guides/write-functional-specs.md#two-whys-each-in-its-place)
  rule, so the reasoning survives this case's archival.

## Durable spec changes

This change touches no durable `specs/` concept and not `SPECIFICATION.md`; the
edited standard lives in `docs/guides/`, tracked in the change case's
[Affected artifacts](../0100-strengthen-spec-requirement-standards.md#affected-artifacts)
index.

### To add

None.

### To modify

None.

### To rename

None.

### To delete

None.
