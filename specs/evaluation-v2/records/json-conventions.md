---
type: Functional Specification
title: Evaluation v2 JSON conventions
description: Shared JSON conventions for Evaluation v2 routine payloads.
tags: [evaluation, json, records]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation v2 JSON conventions

Evaluation v2 stores direct routine payloads as JSON files under `data/`.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Common Fields

Every Evaluation v2 JSON payload **MUST** include `schemaVersion`.

Every Evaluation v2 JSON payload **MUST** include `kind`.

`kind` **MUST** name the payload type.

`schemaVersion` **MUST** be a payload-shape marker only. The CLI **MUST NOT**
silently migrate, upgrade, downgrade, or transform payloads across schema
versions.

## Identity And References

Persisted routine JSON **MUST** use structural model identity fields:

- `AreaId`
- `FactorId`
- `RequirementId`
- `RatingLevelId`

Rendered model refs such as `area:api`, `factor:api::reliability`,
`requirement:api::retry-window`, and `rating:target` **MUST NOT** replace
structural IDs inside persisted routine JSON.

Generated routine outputs, protocol guidance, report artifacts, and payload-local
artifacts **MUST** use `*Ref` names.

Payload-local IDs are local to the containing payload unless the payload kind
defines a wider owner.

## Optional And Repeated Fields

Optional fields **SHOULD** be omitted when absent.

Repeated fields **SHOULD** be present as arrays and default to `[]`.

Fields **SHOULD NOT** use `null` unless the payload kind explicitly defines
`null` as meaningful.

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
