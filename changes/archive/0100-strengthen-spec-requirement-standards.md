---
type: Change Case
title: Strengthen spec requirement standards (29148 + EARS)
description: Patch the functional-spec and change-case guides to close requirement-quality gaps found against ISO/IEC/IEEE 29148 and to add an optional EARS statement template.
status: Done
tags: [process, specs, change-cases, requirements]
timestamp: 2026-06-26T00:00:00Z
---

# Strengthen spec requirement standards (29148 + EARS)

This parent concept captures the _why_ and _status_; the detail lives in its
child:

- [Functional spec](0100-strengthen-spec-requirement-standards/spec.md) — what
  the change must do.

No separate design doc: the change is bounded wording added to two existing
guides, and the per-requirement rationale in the spec carries the _how_ and
_why_.

## Motivation

We assessed our spec-authoring standard against ISO/IEC/IEEE 29148:2018
(requirements engineering) and the Easy Approach to Requirements Syntax (EARS).
The headline finding is reassuring: our per-requirement
[quality bar](../docs/guides/write-functional-specs.md#requirement-quality-bar)
already re-derives 29148's nine individual-requirement characteristics (§5.2.5)
almost one-to-one. But the assessment surfaced concrete gaps the standard
addresses and we do not:

- **No set-level review.** Our bar checks each requirement alone; 29148 §5.2.6
  asks the requirement _set_ to be consistent (no conflicts/overlaps, consistent
  terminology), complete (no unresolved TBD/TBR), and able to be validated.
- **No validation-vs-verification distinction.** We check "did we build it
  right?" (a per-requirement verification path) but never "did we build the right
  thing?" — whether satisfying the whole set achieves the change's motivation.
- **No Assumptions & dependencies element** (29148 §9.6.8): a named home for
  external facts that, if they change, invalidate requirements — distinct from
  Scope and Background.
- **References not classified** (29148 §9.2.4): we name a source of truth but do
  not separate binding/normative references from informational ones.
- **"Unambiguous" is not a named bar item**; it lives only inside an example.
- **No standard statement template** — the one 29148 characteristic
  ("Conforming," §5.2.5) we lack. EARS supplies a lightweight, optional template
  (`While`/`When`/`Where`/`If-Then` + a response) that fills this gap and gives
  explicit syntax for the trigger and edge/error cases our bar already demands.

Closing these keeps our prose-first, palette-not-checklist, BCP 14 conventions
intact while making the standard match the best available guidance — and lets the
guides cite their provenance so the reasoning survives.

## Scope

Covered:

- Strengthen [Writing functional specs](../docs/guides/write-functional-specs.md)
  with: a set-level requirement check; an optional **Assumptions &
  dependencies** Shape element; reference classification (normative vs
  informational); an **Unambiguous** bar item plus active-voice/subject-explicit
  language guidance and a note recording our BCP 14 vs `shall` choice as
  deliberate; and an optional EARS statement template.
- Strengthen [Working with change cases](../docs/guides/work-with-change-cases.md)
  with a Draft→Design **validation gate**: confirm that satisfying the full
  requirement set achieves the case's motivation.
- Carry brief 29148/EARS provenance where each rule is added, per the guides' own
  "two whys" rule.

Deferred / non-goals:

- `SPECIFICATION.md` and the QUALITY.md format's own guidance on authoring
  `requirements` are **out of scope** — a separate pass if wanted.
- No heavyweight 29148 machinery: no unique immutable requirement IDs, owner,
  version, priority, numeric risk, or difficulty attributes; no BRS/StRS/SyRS/SRS
  document hierarchy; no formal Requirements Traceability Matrix.
- EARS is **optional**, never mandatory.
- No code, CLI, bundled-skill, rating, evaluation, report, or format-schema
  change.

## Affected artifacts

### Code

- [x] None — docs-only change. (Deliberate.)

### Format spec

- [x] None — `SPECIFICATION.md` unaffected. (Deliberate scope: internal process
      only.)

### Durable specs

- [x] None — no `specs/` concept governs spec authoring; the standard lives in
      `docs/guides/`.

### Durable docs / bundled skill

- [x] [`docs/guides/write-functional-specs.md`](../docs/guides/write-functional-specs.md)
      — added the set-level **The requirement set** check, the **Assumptions &
      dependencies** Shape element, normative/informational reference
      classification, the **Unambiguous** bar item plus active-voice/subject
      language guidance and the BCP 14-vs-`shall` note, and the optional EARS
      **statement template**, each with brief 29148/EARS provenance.
- [x] [`docs/guides/work-with-change-cases.md`](../docs/guides/work-with-change-cases.md)
      — added the Draft→Design validation check and cross-reference to the
      set-level check.

### Suggested new durable specs

- None.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Spec settled
(no design doc needed), both guides updated, and archived. No code clock — this
is a durable-docs change.
