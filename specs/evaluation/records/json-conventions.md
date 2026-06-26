---
type: Functional Specification
title: Evaluation JSON conventions
description: Shared JSON conventions for Evaluation routine payloads.
tags: [evaluation, json, records]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation JSON conventions

Evaluation stores direct routine payloads as JSON files under `data/`.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Common Fields

Every Evaluation JSON payload **MUST** include `schemaVersion`.

Every Evaluation JSON payload **MUST** include `kind`.

`kind` **MUST** name the payload type.

`schemaVersion` **MUST** be a payload-shape marker only. The CLI **MUST NOT**
silently migrate, upgrade, downgrade, or transform payloads across schema
versions.

## Identity And References

Persisted routine JSON **MUST** encode model identity fields as canonical
qualified model-reference strings:

- `areaId`: `area:<area-path>`
- `factorId`: `factor:<declaring-area-path>::<factor-path>`
- `requirementId`: `requirement:<declaring-area-path>::<requirement-name>`
- `ratingLevelId`: `rating:<rating-level-id>`

Repeated identity fields such as `areaIds`, `factorIds`,
`localRequirementIds`, `rootFactorIds`, `childAreaIds`, `ratingLevelIds`, and
the secondary Factor lists **MUST** be arrays of those same qualified reference
strings and default to `[]` when empty.

Persisted routine JSON **MUST NOT** use structured identity sub-fields such as
`declaringAreaId`, `factorPath`, or `requirementName`, and **MUST NOT** use
unqualified references.

> Rationale: the qualified reference grammar is the lossless string encoding of
> the same composite identity; carrying the parsed object in persisted JSON adds
> shape complexity without adding information. The `root` Area-name reservation
> keeps `area:root` unambiguous. — 0120

The CLI **MUST** resolve persisted identity fields against the run's
`model-snapshot.md` before accepting a write. A payload that names an absent
Area, Factor, Requirement, or Rating Level **MUST** be rejected rather than
persisted as a free-form string.

Generated routine outputs, protocol guidance, report artifacts, and payload-local
artifacts **MUST** use `*Ref` names.

Payload-local IDs are local to the containing payload unless the payload kind
defines a wider owner.

## Optional And Repeated Fields

Optional fields **SHOULD** be omitted when absent.

Repeated fields **SHOULD** be present as arrays and default to `[]`.

Fields **SHOULD NOT** use `null` unless the payload kind explicitly defines
`null` as meaningful.

Payload object fields **MUST** match the accepted kind contract. Unknown or
misspelled fields **MUST** be rejected instead of silently persisted.

## Confidence

Routine outputs that carry judgment confidence **MUST** use:

- `high`
- `medium`
- `low`
- `none`

`high` means evidence is direct, current, sufficient for the claim, and
independently checkable.

`medium` means evidence is relevant and plausible but partial, indirect, sampled,
or not fully verified.

`low` means evidence is thin, ambiguous, stale, inferred, or materially
incomplete.

`none` means no confidence judgment was possible because there was no
assessment-quality evidence.
