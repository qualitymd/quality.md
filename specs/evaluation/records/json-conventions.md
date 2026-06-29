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

A reference object's `kind` field **MUST** name a value from a closed
vocabulary, and a write whose reference `kind` falls outside that vocabulary
**MUST** be rejected rather than persisted as a free-form string, parallel to
the identity-resolution rule above. A routine reference (`*Ref` / `inputRefs[]`)
`kind` names a supported payload kind (per
[`payload-kinds.md`](payload-kinds.md), the full set the CLI can persist,
including the CLI-owned `EvaluationOutputResult`). A report reference `kind`
names one of the report kinds `run`, `area`, `factor`, `requirement`,
`findings`, `recommendations`, or `recommendation`. The `run` report kind
identifies the run-level `report.md` and does not include Area, Factor, or
Requirement identity fields. Area, Factor, and Requirement report refs include
`areaId`; Factor report refs also include `factorId`; Requirement report refs
also include `requirementId`. `findings` identifies `findings.md`,
`recommendations` identifies `recommendations.md`, and `recommendation`
identifies a recommendation detail report.

> Rationale: `kind` was the one required reference field left as a free string
> while every other closed vocabulary in the contract is enum-validated; pinning
> it moves a misspelled or invented kind to a write-time rejection instead of a
> dangling reference discovered later. — 0124, 0137

Generated routine outputs, protocol guidance, report artifacts, and payload-local
artifacts **MUST** use `*Ref` names.

Payload-local IDs are local to the containing payload unless the payload kind
defines a wider owner. Finding `id` values are payload-local IDs, not Model IDs
or durable cross-run identifiers. Cross-payload references to findings **MUST**
use a routine reference qualified by the owning payload subject plus a selector,
such as `findings[gap-001]`. Candidate action `id` values are local to their
containing Finding, not the payload or run. References to candidate actions
**MUST** qualify the containing Finding selector, for example
`findings[gap-001].candidateActions[action-001]`.

Evaluation run identity has two parts:

- `RunManifest.id` is the globally-unique, opaque run identifier. New run IDs
  are generated as `<timestamp>-<nanoid>`, where `<timestamp>` is the UTC
  ISO-8601 basic creation timestamp (`YYYYMMDDThhmmssZ`) and `<nanoid>` is at
  least 12 lowercase base32 characters from cryptographic-strength randomness.
  Readers **MUST** treat it as opaque and require only that it is non-empty.
- `RunManifest.number` is the local, friendly run number used in folder names
  and report headers. It is not globally unique and **MUST NOT** be presented as
  a durable handoff identifier.

The structured Evaluation identity rule is: `id` names opaque artifact identity;
user-facing `number` values name ranked report positions. `RunManifest.id`
names the run. `RecommendationResult.id` names a recommendation artifact within
the run and **MUST** use the `qrec_<token>` form, where `<token>` contains only
lowercase ASCII letters and digits. Requirement Findings use their
requirement-scoped selector. `qualitymd evaluation data set` **MUST** assign a
missing `RecommendationResult.id` before writing a recommendation and **MUST
NOT** assign an artifact ID to `FindingRankingResult.orderedFindings[]`.

Recommendation ranking and finding coverage **MUST** reference recommendations by
their assigned `id`. Cross-payload references to findings **MUST** use the
Requirement Assessment routine reference plus selector, for example
`findings[gap-001]`.

The user-facing recommendation number is derived from
`RecommendationRankingResult.orderedRecommendations[].rank`; in reports it is
rendered as `#` or `Number`. Typed artifact references used in reports and
external handoff text **SHOULD** combine the run ID and recommendation ID, for
example
`evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_7h4km2p9`. These
artifact references do not replace routine `*Ref` objects or canonical Model
references.

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

Fixed Evaluation enum values **MUST** be written as their canonical persisted
values, not as report display labels, markers, aliases, or case variants.
