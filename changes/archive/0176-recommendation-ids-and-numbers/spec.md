---
type: Functional Specification
title: Recommendation IDs and Numbers — functional spec
description: Requirements for opaque recommendation IDs and ranking-derived recommendation numbers.
tags: [evaluation, advice, recommendations, identifiers]
timestamp: 2026-06-29T00:00:00Z
---

# Recommendation IDs and Numbers — functional spec

Companion to
[Recommendation IDs and Numbers](../0176-recommendation-ids-and-numbers.md).
This spec states the delta for Evaluation recommendation identity, structured
payload references, generated reports, and recommendation follow-up language. The
durable source of truth is absorbed into
[`Evaluation JSON conventions`](../../../specs/evaluation/records/json-conventions.md),
[`Evaluation payload kinds`](../../../specs/evaluation/records/payload-kinds.md),
[`Evaluation data layout`](../../../specs/evaluation/records/data-layout.md),
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md), and
[`qualitymd evaluation data`](../../../specs/cli/evaluation-data.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The pre-change contract assigns `RecommendationResult.number` before
recommendation ranking. Generated reports then also render recommendation rank.
Because both values are small integers attached to the same advice item, people
and agents can naturally read "recommendation #1" as either the first ranked
recommendation or the first persisted recommendation payload. The format should
make that phrase unambiguous.

The fix is to reserve small recommendation numbers for the ranked advice list and
move payload identity to a visibly non-ordinal ID. Recommendation payloads still
need identity before ranking because the evaluation workflow writes
`RecommendationResult` payloads before `RecommendationRankingResult`. An opaque
`qrec_...` ID satisfies that sequencing without competing with the user-facing
number.

## Scope

Covered: recommendation identity, recommendation ranking references, finding
coverage references, recommendation data paths, generated recommendation report
tables and detail metadata, runtime evaluation guidance, and recommendation
follow-up selection language.

Non-goals:

- changing evaluation run identity or run numbering;
- changing Requirement Finding payload-local IDs, selectors, report anchors, or
  finding ranking identity;
- changing Area, Factor, Requirement, or Rating Level Model references;
- adding migration or compatibility readers for pre-change recommendation-number
  data;
- changing the Advice requirement to produce at least one recommendation.

## Requirements

1. `RecommendationResult` **MUST** carry an opaque `id` string instead of
   `number`; valid recommendation IDs **MUST** match `^qrec_[a-z0-9]+$`.

   > Rationale: The `qrec_` prefix makes the value self-describing while the
   > opaque token prevents humans and agents from reading artifact identity as
   > recommendation order.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
   > `specs/evaluation/records/json-conventions.md`.

2. `qualitymd evaluation data set` **MUST** assign a missing
   `RecommendationResult.id` before persisting the payload, and assigned IDs
   **MUST** be unique within the owning run.

   > Rationale: The CLI already owns mechanical artifact assignment. Keeping ID
   > assignment there preserves the current write-before-ranking workflow.
   > Durable spec: modify `specs/cli/evaluation-data.md`,
   > `specs/evaluation/records/payload-kinds.md`, and
   > `specs/evaluation/records/json-conventions.md`.

3. `qualitymd evaluation data set` **MUST** preserve a supplied valid
   `RecommendationResult.id` when the payload intentionally rewrites an existing
   recommendation, and **MUST** reject duplicate or malformed IDs in the
   effective run data.

   > Rationale: Correcting a persisted recommendation should not churn the
   > artifact identity, but malformed or duplicate identities make ranking and
   > coverage ambiguous.
   > Durable spec: modify `specs/cli/evaluation-data.md` and
   > `specs/evaluation/records/payload-kinds.md`.

4. `qualitymd evaluation data set` **MUST** derive the `RecommendationResult`
   data path from the recommendation ID at
   `data/advice/recommendations/<recommendation-id>/recommendation-result.json`.

   > Rationale: Recommendation data paths need stable artifact identity before
   > ranking exists; the opaque ID provides that identity without consuming the
   > user-facing recommendation number.
   > Durable spec: modify `specs/evaluation/records/data-layout.md`.

5. `RecommendationRankingResult.orderedRecommendations[].recommendationRef` and
   `findingCoverage[].recommendationRefs` **MUST** reference recommendations by
   `RecommendationResult.id`, and a write **MUST** be rejected when a referenced
   ID has no corresponding `RecommendationResult` in the run.

   > Rationale: Ranking and coverage are structured links to recommendation
   > artifacts; using IDs preserves write sequencing and avoids ordinal
   > ambiguity.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
   > `specs/evaluation/records/json-conventions.md`, and
   > `specs/cli/evaluation-data.md`.

6. Generated recommendation numbers **MUST** be derived from
   `RecommendationRankingResult.orderedRecommendations[].rank`, and human reports
   **MUST** label the value as `#` or `Number`, not as the recommendation ID.

   > Rationale: This aligns generated artifacts with common parlance:
   > recommendation #1 is the first ranked recommendation.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

7. Generated recommendation detail report filenames **MUST** use the
   recommendation number plus title slug at
   `recommendations/<NNN>-<slug>.md`; `<NNN>` is the zero-padded recommendation
   number derived from ranking.

   > Rationale: Detail reports are human-facing advice artifacts, so their
   > filenames should match the numbered ranked list users navigate.
   > Durable spec: modify `specs/evaluation/records/data-layout.md` and
   > `specs/evaluation/reports/report-tree.md`.

8. Generated recommendation detail reports **SHOULD** include the opaque
   recommendation ID in source-data or metadata context, but generated run and
   recommendation index tables **MUST NOT** present the opaque ID as a competing
   recommendation number.

   > Rationale: The ID remains useful for debugging, structured handoff, and data
   > traceability, but the main scan surface should have one small number.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

9. Typed recommendation artifact references used in reports and external handoff
   text **SHOULD** combine the run ID and recommendation ID, for example
   `evaluation:<run-id>/recommendation/<qrec-id>`.

   > Rationale: The durable handoff reference should point to stable artifact
   > identity, not ranked order that exists only in the generated advice view.
   > Durable spec: modify `specs/evaluation/records/json-conventions.md` and
   > `specs/evaluation/reports/report-tree.md`.

10. Recommendation follow-up guidance **MUST** treat `recommendation 1`,
    `rec #1`, and bare numeric selection as recommendation-number selection from
    ranking order, while accepting `qrec_...` values as explicit ID selection
    when present.

    > Rationale: The skill is the primary interface for acting on
    > recommendations, so it must match the report language and still support
    > exact artifact references.
    > Durable spec: modify `specs/skills/quality-skill/recommendation-follow-up.md`.

11. Evaluation runtime guidance **MUST** instruct agents to write
    `RecommendationResult` payloads before `RecommendationRankingResult`, read
    assigned recommendation IDs, and use those IDs in recommendation ranking and
    finding coverage.

    > Rationale: The change is deliberately preserving sequencing; the skill
    > needs to stop asking agents to read assigned numbers before ranking.
    > Durable spec: modify `specs/skills/quality-skill/workflows/evaluate.md`.

12. The CLI and runtime guidance **MUST NOT** retain
    `RecommendationResult.number` as an accepted current-contract field, and
    **MUST NOT** add a compatibility shim, fallback reader, migration command, or
    dual writer for pre-change recommendation-number data.

    > Rationale: QUALITY.md is early alpha; carrying both fields would recreate
    > the ambiguity this case removes.
    > Durable spec: modify `specs/cli/evaluation-data.md` and
    > `specs/evaluation/records/payload-kinds.md`.

## Acceptance criteria

- `qualitymd evaluation data set` assigns missing `RecommendationResult.id`
  values matching `qrec_<token>`, preserves supplied valid IDs, and rejects
  malformed or duplicate IDs.
- Recommendation JSON data paths use
  `data/advice/recommendations/<qrec-id>/recommendation-result.json`.
- `RecommendationRankingResult` and finding coverage references use `qrec_...`
  IDs and reject unknown IDs.
- `RecommendationResult.number` is absent from the current JSON contract, schema,
  examples, specs, tests, and runtime guidance.
- Generated report tables render user-facing recommendation numbers from ranking
  order and do not show opaque IDs in the main scan table.
- Recommendation detail report filenames remain number-prefixed by ranked order.
- Recommendation detail report source-data or metadata context includes the
  `qrec_...` ID or the data path that contains it.
- Recommendation follow-up guidance resolves numeric user input by
  recommendation number and `qrec_...` input by ID.
- Generated Evaluation data JSON Schema is regenerated and current.
- Focused Go tests, report-gallery checks, and `mise run check` pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/records/json-conventions.md` - define recommendation IDs,
  recommendation numbers, ID-based structured refs, and typed references.
- `specs/evaluation/records/payload-kinds.md` - replace
  `RecommendationResult.number` with `id` and make ranking/coverage refs
  strings.
- `specs/evaluation/records/data-layout.md` - recommendation data paths derive
  from IDs; recommendation report paths derive from recommendation numbers.
- `specs/evaluation/reports/report-tree.md` - report tables and detail reports
  render number and ID in distinct roles.
- `specs/evaluation/protocol.md`,
  `specs/evaluation/orchestration.md`, and
  `specs/evaluation/routines/routine-contracts.md` - align Advice sequencing and
  reference semantics.
- `specs/cli/evaluation-data.md` - data assignment, validation, and examples use
  recommendation IDs.
- `specs/skills/quality-skill/workflows/evaluate.md` - evaluation workflow uses
  assigned IDs before ranking.
- `specs/skills/quality-skill/recommendation-follow-up.md` - numeric follow-up
  selection means recommendation number; `qrec_...` means recommendation ID.

### To rename

None

### To delete

None
