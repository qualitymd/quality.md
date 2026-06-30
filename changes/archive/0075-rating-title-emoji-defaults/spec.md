---
type: Functional Specification
title: Rating title emoji defaults — functional spec
description: Required behavior for making emoji-prefixed Rating Level titles the default starter and setup convention without changing formal rating semantics.
tags: [scaffold, skill, rating-scale, docs]
timestamp: 2026-06-24T00:00:00Z
---

# Rating title emoji defaults — functional spec

Companion to the
[Rating title emoji defaults](../0075-rating-title-emoji-defaults.md) change
case. This spec states _what_ the change must do; a later
[design doc](design.md) will cover how the implementation lands.

The key words **MUST**, **SHOULD**, and **MAY** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Scope

Covers the default display convention for the standard four-level Rating Scale
seeded by `qualitymd init`, recommended by `/quality setup`, and shown in
starter/user-facing examples. Deferred: any schema change, lint rule, migration
command, rating semantics change, or requirement that all Rating Level titles use
emoji.

## Requirements

### Default titles

- R1. The standard four-level Rating Scale default **MUST** keep stable Rating
  Level IDs as `outstanding`, `target`, `minimum`, and `unacceptable`, in that
  order.
- R2. When a tool or guide authors the standard scale as a default starter, it
  **SHOULD** use the human-facing titles `🟢 Outstanding`, `🔵 Target`,
  `🟡 Minimum`, and `🔴 Unacceptable`.
  > Rationale: The emoji marker is a display aid for repeated human report and
  > model scanning; the stable `level` ID continues to carry machine meaning.
  > — 0075
- R3. Tools and guides **MUST NOT** treat the emoji marker as the source of rating
  semantics, ordering, identity, or machine references.
- R4. Guidance that recommends emoji-prefixed default titles **MUST** still make
  clear that Rating Level titles are human display labels and can be customized
  when a project uses a different scale or a plain-text house style.
- R5. Guidance **SHOULD** avoid emoji-only titles; default titles include both a
  lightweight marker and the word label.

### Surfaces

- R6. `qualitymd init` and `qualitymd init --minimal` **MUST** seed the standard
  four-level Rating Scale with the emoji-prefixed titles in R2.
- R7. `/quality setup` **MUST** recommend the standard scale using the
  emoji-prefixed human titles in R2 while presenting stable Rating Level IDs as
  plain identifiers.
- R8. The authoring guide **MUST** describe the emoji-prefixed titles as the
  recommended default display convention for the standard scale, not as a
  conformance requirement.
- R9. User-facing starter examples that present the recommended default scale
  **SHOULD** show the emoji-prefixed titles unless the surrounding context is
  specifically about formal syntax or machine identifiers.
- R10. Existing example evaluation artifacts that intentionally demonstrate
  stable IDs or historical report output **MAY** remain unchanged unless they are
  directly presented as current starter defaults.

### Format neutrality

- R11. `SPECIFICATION.md` **MUST NOT** require emoji in Rating Level titles.
- R12. If `SPECIFICATION.md` mentions visual markers, it **MUST** do so only as
  non-normative guidance about human-readable labels.

## Durable spec changes

Durable **specs** this change rewrites — the
[`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/skills/quality-skill/workflows/setup.md` — record setup's
  emoji-prefixed standard-scale default (per R2, R4, and R7).
- `specs/skills/quality-skill/guides/authoring-md.md` — record the authoring
  guide's recommended display convention for default Rating Level titles (per
  R2, R4, and R8).
- `SPECIFICATION.md` — only if needed, add non-normative visual-marker guidance
  while preserving format neutrality (per R11 and R12).

### To rename

None.

### To delete

None.
