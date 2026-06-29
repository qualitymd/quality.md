---
type: Design Doc
title: Recommendation IDs and Numbers — design doc
description: Design for opaque recommendation IDs and ranking-derived recommendation numbers.
tags: [evaluation, advice, recommendations, identifiers]
timestamp: 2026-06-29T00:00:00Z
---

# Recommendation IDs and Numbers — design doc

Design for
[Recommendation IDs and Numbers](../0176-recommendation-ids-and-numbers.md) and
its [functional spec](spec.md).

## Context

The current implementation already has the right sequencing for this change:
`SetData` normalizes recommendation payloads before deriving their data paths,
then later validates `RecommendationRankingResult` against the effective run
data. The issue is the normalized value: it is a positive integer named
`number`, which reports and users can confuse with ranked order.

The design keeps the existing flow and changes only the identity type.
Recommendation results get opaque IDs before ranking; ranking assigns the
human-facing number through its existing `rank` field.

## Approach

### ID generation and validation

Replace recommendation number assignment in `internal/evaluation/data.go` with
recommendation ID assignment. The assignment helper keeps the same lifecycle:
scan existing recommendation result payloads, preserve supplied valid IDs, and
assign only missing IDs before path derivation.

Generated IDs use `qrec_` plus a short lowercase alphanumeric token from
`crypto/rand`. Validation is intentionally simple and matches the spec:
`^qrec_[a-z0-9]+$`. There is no semantic payload in the token and no ordering
relationship between IDs.

The helper checks for uniqueness in the candidate batch and effective data. If a
generated token collides with an existing ID, it generates another. Because IDs
are run-local and random, no global registry or directory-number scan is needed.

### Data paths and query refs

`recommendationDataPath` changes from a zero-padded integer path to:

```text
data/advice/recommendations/<qrec-id>/recommendation-result.json
```

`dataPathForRoutineRef` and any CLI data query selector paths treat
recommendation selectors as IDs. Ranking and coverage validators build their
recommendation lookup maps by `RecommendationResult.id`, and
`orderedRecommendations[].recommendationRef` plus
`findingCoverage[].recommendationRefs` become strings.

### Report rendering

The report model can continue to derive a ranked recommendation list from
`RecommendationRankingResult.orderedRecommendations`, but the list item needs
two fields:

- `number`: the user-facing recommendation number, copied from `rank`;
- `id`: the opaque `RecommendationResult.id`.

Run and recommendation index tables render the existing `#` column from
`number`. They do not render an ID column. Recommendation detail report filenames
continue to use the number prefix and slug. Detail report source-data links point
to the ID-based JSON path, and the detail report's typed artifact reference uses
`evaluation:<run-id>/recommendation/<qrec-id>`.

### Schema and examples

`data_contract.go` is the typed source for the generated Evaluation data schema:
change the recommendation result field, ranking ref field, and coverage ref
array from number to string. Regenerate `evaluation-data.schema.json`.

Tests and examples should stop asserting numeric recommendation refs. Where a
stable fixture is needed, provide explicit `qrec_...` IDs in test payloads rather
than depending on random assignment.

### Skill guidance

Runtime evaluation guidance keeps the two-step write:

1. write `RecommendationResult` payloads without IDs or with explicit IDs;
2. read assigned IDs from persisted paths or payloads;
3. author `RecommendationRankingResult` and finding coverage with those IDs.

Recommendation follow-up wording treats numeric user input as ranked
recommendation-number selection and `qrec_...` as an exact artifact ID.

## Spec response

- Opaque ID assignment satisfies the identity requirements without changing
  Advice sequencing.
- ID-based data paths and refs keep structured links stable before ranking.
- Ranking-derived report numbers make "recommendation #1" mean the first ranked
  recommendation.
- Source-data links and typed references preserve exact artifact traceability.
- Removing the current `number` field without a compatibility path satisfies the
  early-alpha clean-break boundary.

## Alternatives

- **Make `RecommendationResult.number` equal rank.** Rejected because
  recommendation results are written before ranking. Assigning rank as identity
  would require drafting recommendations in memory, ranking them, then writing
  results, increasing workflow complexity.
- **Keep numeric artifact numbers but hide them.** Rejected because agents still
  consume JSON paths and structured refs. Two small integers would remain in the
  system even if one is hidden in reports.
- **Use alphabetic sequence IDs (`A`, `B`, `C`).** Rejected because alphabetic
  labels still imply order and have awkward rollover behavior.
- **Use title slugs as IDs.** Rejected because titles are user-facing prose,
  mutable, and collision-prone.
- **Use full UUIDs.** Rejected because they are visually noisy in paths and
  reports. A prefixed short random token gives enough run-local uniqueness with
  better ergonomics.
- **Use `qmd_rec_` as the prefix.** Rejected as unnecessarily noisy inside the
  QUALITY.md evaluation format. `qrec_` is shorter and still typed.

## Trade-offs & risks

- **Random IDs add nondeterminism.** Tests and generated examples should use
  explicit IDs in payloads where stable output matters; assignment itself is
  covered by focused tests.
- **Typed references no longer expose the visible recommendation number.** This
  is deliberate: durable handoff should point to artifact identity. Reports still
  provide the user-facing number in the ranked list.
- **Pre-change runs break under current strict readers.** This is acceptable
  under early-alpha policy and avoids preserving the ambiguous contract.

## Open questions

None.
