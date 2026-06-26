---
type: Functional Specification
title: Constrain reference kind fields — functional spec
description: Requirements for constraining the `kind` member of Evaluation reference objects to the supported payload-kind vocabulary.
tags: [evaluation, records, references, schema]
timestamp: 2026-06-26T00:00:00Z
---

# Constrain reference kind fields — functional spec

Companion to the
[Constrain reference kind fields](../0124-reference-kind-enum.md) change case.
This spec states *what* the change must do; the [design doc](design.md) covers
*how*. The supported payload kinds are defined by
[`Evaluation payload kinds`](../../specs/evaluation/records/payload-kinds.md)
(normative); reference encoding is defined by
[`Evaluation JSON conventions`](../../specs/evaluation/records/json-conventions.md)
(normative).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

An Evaluation reference object — `{kind, subject, selector}` for routine inputs,
and the report reference shape — names another payload by its `kind`. The data
contract already enum-validates every other closed vocabulary it carries
(`finding.type`, `severity`, `confidence`, `status`) and already resolves a
reference's identity `subject` against the run's model snapshot, rejecting absent
IDs. The reference `kind` is the one required field left as a free string, so a
misspelled or invented kind persists silently and only fails — if at all — later,
when a consumer tries to resolve the dangling reference. Constraining it to the
payload-kind set closes that gap and moves the failure to write time, where it is
cheap to fix. See the change case
[Motivation](../0124-reference-kind-enum.md#motivation) for the originating
evidence.

## Scope

Covered: the `kind` member of Evaluation **reference objects** (routine
`inputRefs[]` / `*Ref` objects and report references).

Deferred / non-goals:

- The **top-level** payload `kind` field is out of scope — it is already pinned
  by the `evaluation data set <kind>` argument and the agent-writable kind gate.
- The reference **`selector`** field is out of scope; its vocabulary is bounded
  but not yet confirmed, and is tracked separately.
- No `SchemaVersion` change: the constraint rejects only data that was already
  invalid in intent, and does not alter the shape of conforming payloads.

## Requirements

- The Evaluation data contract **MUST** constrain the `kind` member of every
  reference object — both the routine reference shape (used for `inputRefs[]` and
  the `*Ref` fields) and the report reference shape — to the set of supported
  Evaluation payload kinds.

  >> Rationale: `kind` is the only required reference field not validated against
  >> a closed set, while every other closed vocabulary in the contract already is;
  >> a free-string kind is the contract's one silent typo path. — 0124

- The supported payload-kind vocabulary for a reference `kind` **MUST** include
  every payload kind the CLI can persist, including the CLI-owned
  `EvaluationOutputResult`, and **MUST NOT** be narrowed to only the
  agent-writable kinds.

  >> Rationale: naming a payload type in a reference is distinct from being
  >> permitted to *write* that type through `evaluation data set`; coupling the
  >> two would reject a legitimate reference to a CLI-owned payload. — 0124

- When `qualitymd evaluation data set` validates a payload whose reference `kind`
  is present but names a value outside the supported payload-kind vocabulary, it
  **MUST** reject the write before persisting and report the offending field path
  and value.

- The reference `kind` constraint **MUST** appear as an `enum` in
  `qualitymd evaluation data schema` output, so the closed vocabulary is
  discoverable from the typed contract rather than only from prose.

- The reference-`kind` vocabulary **MUST** derive from a single typed source of
  truth shared with the rest of the contract's kind handling, so the enum cannot
  drift from the set of kinds the CLI actually supports.

  >> Rationale: the payload-kinds spec already requires validation to derive from
  >> one typed source for fields, enum values, schema, and examples; a
  >> hand-maintained second kind list would violate that and rot. — 0124

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../SPECIFICATION.md). Each subsection is required.

### To add

None.

### To modify

- `specs/evaluation/records/json-conventions.md` — in "Identity And References",
  add that a reference's `kind` field MUST name a supported payload kind (per
  [`payload-kinds.md`](../../specs/evaluation/records/payload-kinds.md)) and MUST
  be rejected when it does not, parallel to the existing rule that resolves
  identity subjects against the model snapshot (per the constraint and rejection
  requirements above).

### To rename

None.

### To delete

None.
