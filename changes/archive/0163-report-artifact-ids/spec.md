---
type: Functional Specification
title: Report Artifact IDs — functional spec
description: Requirements for citable Evaluation report artifact IDs and typed references.
tags: [evaluation, reports, advice, identifiers]
timestamp: 2026-06-29T00:00:00Z
---

# Report Artifact IDs — functional spec

Companion to [Report Artifact IDs](../0163-report-artifact-ids.md). This spec
states the delta for Evaluation structured payloads and generated reports. The
durable source of truth is absorbed into
[`Evaluation JSON conventions`](../../../specs/evaluation/records/json-conventions.md),
[`Evaluation payload kinds`](../../../specs/evaluation/records/payload-kinds.md),
[`Evaluation data layout`](../../../specs/evaluation/records/data-layout.md),
and [`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md).
The CLI data-write behavior is absorbed into
[`qualitymd evaluation data`](../../../specs/cli/evaluation-data.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Evaluation outputs are not only files; they are handoff artifacts. A
recommendation or top finding may be cited in chat, copied into an issue
tracker, assigned to another agent, or discussed aloud. Artifact identity
therefore needs to be short and stable enough for people while retaining typed,
scoped references for tools. This mirrors the Model identity doctrine:
local IDs are scoped, references are typed, titles are labels, and paths are
derived.

## Scope

This change covers report/evaluation artifacts created by an Evaluation run:

- the run identity rendered from `RunManifest`;
- CLI-assigned `RecommendationResult` IDs and recommendation references;
- CLI-assigned ranked finding public IDs in `FindingRankingResult`; and
- generated Markdown report tables/detail pages for findings and
  recommendations.

Non-goals:

- changing Area, Factor, Requirement, or Rating Level Model references;
- changing Requirement Finding payload-local IDs;
- making artifact IDs globally unique without run context;
- adding UUIDs or random IDs;
- migrating historical runs.

## Requirements

1. Generated Evaluation reports **MUST** render run identity as
   `QEVAL-<NNNN>`, where `<NNNN>` is the zero-padded `RunManifest.number`.

   > Rationale: The run is the provenance owner for recommendation and finding
   > artifact IDs, and the run header is the first place readers see that
   > ownership.
   > Durable spec: modify `specs/evaluation/records/json-conventions.md` and
   > `specs/evaluation/reports/report-tree.md` to define the run artifact ID.

2. `qualitymd evaluation data set` **MUST** assign a user-citable
   `RecommendationResult.id` matching `QREC-<NNNN>-<NNN>` when a
   `RecommendationResult` input payload omits `id`, and `<NNNN>` **MUST** match
   the owning `RunManifest.number`.

   > Rationale: Users naturally ask for the recommendation "ID" when handing it
   > off. The ID should be the stable token they copy or say, and allocation is
   > mechanical enough for the CLI to own.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
   > `specs/evaluation/records/json-conventions.md`, and
   > `specs/evaluation/records/data-layout.md`.

3. `qualitymd evaluation data set` **MUST** persist `RecommendationResult`
   payloads with `id` present, whether the ID was supplied for an existing
   update or assigned for a new recommendation.

   > Rationale: Stored Evaluation data should remain self-contained and should
   > not require readers to rerun assignment logic.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md`,
   > `specs/evaluation/records/json-conventions.md`, and
   > `specs/evaluation/records/data-layout.md`.

4. `qualitymd evaluation data set` **MUST** assign recommendation IDs
   sequentially from the owning run's existing `QREC` IDs, preserving a supplied
   valid `QREC` ID only when the payload is intentionally replacing or rewriting
   that recommendation.

   > Rationale: The CLI should prevent duplicate public IDs while still allowing
   > an already-assigned recommendation payload to be corrected.
   > Durable spec: modify `specs/evaluation/records/json-conventions.md` and
   > `specs/evaluation/records/payload-kinds.md`.

5. `RecommendationRankingResult` recommendation references and finding coverage
   recommendation references **MUST** use persisted `RecommendationResult.id`
   values.

   > Rationale: Ranking and coverage should point at assigned recommendation
   > artifacts, not title text, slugs, or temporary authoring labels.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md`.

6. `qualitymd evaluation data set` **MUST** assign
   `FindingRankingResult.orderedFindings[].id` values matching
   `QFIND-<NNNN>-<NNN>`, and `<NNNN>` **MUST** match the owning
   `RunManifest.number`.

   > Rationale: Requirement Finding IDs are intentionally local to their
   > assessment payload. The ranked finding entry is the first run-level
   > artifact surface where a citable public finding ID can be assigned without
   > weakening payload-local selector semantics.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
   > `specs/evaluation/records/json-conventions.md`.

7. `qualitymd evaluation data set` **MUST** preserve existing ranked finding IDs
   by `findingRef` when rewriting `FindingRankingResult`, assign new IDs only to
   newly ranked findings, and keep `FindingRankingResult.orderedFindings[].id`
   unique within one run.

   > Rationale: Re-ranking should not churn citable finding IDs for the same
   > underlying Requirement Finding.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
   > `specs/evaluation/records/json-conventions.md`.

8. `FindingRankingResult.orderedFindings[].id` **MUST NOT** replace
   `findingRef` as the exact structured reference to the underlying Requirement
   Finding.

   > Rationale: The citable ID and the data reference solve different jobs:
   > handoff versus exact traceability.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
   > `specs/evaluation/records/json-conventions.md`.

9. Generated finding and recommendation index tables **MUST** render an `ID`
   column for the citable artifact ID.

   > Rationale: Top-level report tables are the scanning and handoff surfaces;
   > IDs must be visible there rather than hidden in source data.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

10. Generated recommendation detail reports **MUST** render the recommendation
    artifact ID and a typed recommendation reference of the form
    `evaluation:QEVAL-<NNNN>/recommendation:QREC-<NNNN>-<NNN>`.

> Rationale: The detail report should provide both the short ID and the
> context-complete reference for external systems.
> Durable spec: modify `specs/evaluation/reports/report-tree.md`.

11. Generated finding detail sections **MUST** render the ranked finding artifact
    ID when a finding has a ranking entry, and **MUST** keep anchors derived from
    the payload-local Requirement Finding ID.

> Rationale: Existing links should remain stable to the exact finding
> selector, while readers get the citable report artifact ID near the detail.
> Durable spec: modify `specs/evaluation/reports/report-tree.md`.

12. Generated recommendation report filenames **SHOULD** continue to use ranking
    order plus a title-derived slug, and generated data paths **MUST** derive
    from `RecommendationResult.id`.

> Rationale: Friendly report filenames and durable data identity are separate
> concerns. Rank-prefixed report filenames aid browsing but are not identity.
> Durable spec: modify `specs/evaluation/records/data-layout.md` and
> `specs/evaluation/reports/report-tree.md`.

13. Generated reports **MUST NOT** label Model references as artifact IDs or
    replace Model references with `QEVAL`, `QREC`, or `QFIND` IDs.

    > Rationale: Model identity and Evaluation artifact identity are distinct
    > contracts that should share style, not syntax.
    > Durable spec: modify `specs/evaluation/records/json-conventions.md` and
    > `specs/evaluation/reports/report-tree.md`.

## Acceptance criteria

- `qualitymd evaluation data set` assigns missing `RecommendationResult.id`
  values and persists `QREC-<NNNN>-<NNN>` IDs.
- `qualitymd evaluation data set` assigns missing
  `FindingRankingResult.orderedFindings[].id` values and persists
  `QFIND-<NNNN>-<NNN>` IDs.
- `qualitymd evaluation data set` preserves ranked finding IDs by `findingRef`
  when rewriting a ranking payload.
- `qualitymd evaluation data set` rejects supplied malformed `QREC` and `QFIND`
  IDs.
- `qualitymd evaluation data set` rejects supplied `QREC` and `QFIND` IDs whose
  run segment does not match `RunManifest.number`.
- `RecommendationResult` data paths use `QREC-<NNNN>-<NNN>`.
- Persisted `FindingRankingResult` carries citable finding IDs while Requirement
  Findings keep their payload-local IDs and selectors.
- `report.md`, `findings.md`, `recommendations.md`, recommendation detail
  reports, and Requirement finding detail sections render citable IDs.
- Recommendation detail reports render typed recommendation references.
- Report anchors for Requirement Finding details remain payload-local.
- Durable specs and runtime skill guidance describe the new ID contract.
- The generated Evaluation data JSON Schema is current.
- Focused Go tests and `mise run check` pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/records/json-conventions.md` - define report artifact IDs,
  CLI-owned `QREC`/`QFIND` assignment, typed artifact references, and the
  boundary with Model references (requirements 1-4, 6-8, 13).
- `specs/evaluation/records/payload-kinds.md` - update
  `FindingRankingResult`, `RecommendationResult`, and
  `RecommendationRankingResult` ID/reference contracts (requirements 2-8).
- `specs/evaluation/records/data-layout.md` - define recommendation data paths
  from `QREC` IDs and keep recommendation report filenames slug/rank-derived
  (requirements 2, 3, 12).
- `specs/evaluation/reports/report-tree.md` - render run, finding, and
  recommendation artifact IDs and typed recommendation references without
  relabeling Model references (requirements 1, 9-13).
- `specs/cli/evaluation-data.md` - define `data set` assignment before
  validation and path derivation (requirements 2, 3, 6, 7).
- `specs/skills/quality-skill/workflows/evaluate.md` - align Advice routine
  authoring guidance with CLI-assigned `QREC` and `QFIND` IDs (requirements
  2-8).

### To rename

None

### To delete

None
