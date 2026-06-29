---
type: Functional Specification
title: Run IDs and Artifact Numbering — functional spec
description: Requirements for a globally-unique run ID, per-run recommendation numbers, and removal of the finding artifact ID.
tags: [evaluation, reports, advice, identifiers]
timestamp: 2026-06-29T00:00:00Z
---

# Run IDs and Artifact Numbering — functional spec

Companion to [Run IDs and Artifact Numbering](../0165-run-id-artifact-numbering.md).
This spec states the delta for Evaluation structured payloads and generated
reports. The durable source of truth is absorbed into
[`Evaluation JSON conventions`](../../../specs/evaluation/records/json-conventions.md),
[`Evaluation payload kinds`](../../../specs/evaluation/records/payload-kinds.md),
[`Evaluation data layout`](../../../specs/evaluation/records/data-layout.md),
and [`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md).
The CLI behavior is absorbed into
[`qualitymd evaluation create`](../../../specs/cli/evaluation-create.md) and
[`qualitymd evaluation data`](../../../specs/cli/evaluation-data.md).

This case revises [0163 - Report Artifact IDs](../archive/0163-report-artifact-ids.md);
it replaces the `QEVAL`/`QREC`/`QFIND` contract rather than extending it.

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

0163's artifact IDs are scoped only to `RunManifest.number`, a counter the CLI
derives by scanning the local run directory. That counter is not durable when
the directory is not: an uncommitted or ad-hoc run restarts at `0001`, so two
different runs can mint the same `QEVAL-0001` / `QREC-0001-002` with conflicting
content and no way to disambiguate them. Cross-project handoff collides the same
way. The fix is to give the *run* a highly-improbable, globally-unique `id` and
address everything inside the run through it, while keeping the short tokens
people actually cite. The governing rule: `id` names only the run; recommendations
carry a per-run `number`; findings are addressed by their existing
requirement-scoped reference and carry no synthetic artifact ID.

The friendly `RunManifest.number` is retained on purpose. It still names the run
directory and reads naturally in-repo ("run 1"); its non-durability is now
harmless because the run `id`, not the number, carries the uniqueness guarantee.

## Scope

Covered: the run identity in `RunManifest`; the recommendation number replacing
`RecommendationResult.id`; removal of the ranked-finding artifact ID; and the
generated report rendering of run, recommendation, and finding identity.

Non-goals:

- migrating pre-change runs, or any compatibility shim/fallback reader for the
  `QEVAL`/`QREC`/`QFIND` contract;
- repository- or entity-qualified namespaces above the run ID;
- surfacing the run ID in `evaluation list`/`show`;
- changing Area, Factor, Requirement, or Rating Level Model references;
- changing Requirement Finding payload-local IDs, selectors, or report anchors.

## Requirements

1. `qualitymd evaluation create` **MUST** assign `RunManifest.id` as a
   globally-unique run identifier of the form `<timestamp>-<nanoid>`, where
   `<timestamp>` is the run-creation instant formatted as a UTC ISO-8601 basic
   timestamp (`YYYYMMDDThhmmssZ`) and `<nanoid>` is at least 12 characters drawn
   from a lowercase ambiguity-free base32 alphabet using cryptographic-strength
   randomness.

   > Rationale: The timestamp gives a human-decodable, lexically sortable anchor;
   > the random tail provides the uniqueness guarantee the run number cannot,
   > even under total clock collision across machines. The reduced alphabet keeps
   > the token sayable and transcribable for handoff.
   > Durable spec: modify `specs/evaluation/records/json-conventions.md`,
   > `specs/evaluation/records/payload-kinds.md`, and
   > `specs/cli/evaluation-create.md`.

2. `qualitymd evaluation create` **MUST** persist `RunManifest.id` in the run
   manifest and **MUST NOT** derive it from, or reuse, `RunManifest.number`; the
   `id` and `createdAt` **MUST** be computed from the same creation instant.

   > Rationale: The number is local and non-durable by design, so the unique `id`
   > must be independent of it. Sharing one instant keeps the two manifest fields
   > from skewing.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
   > `specs/cli/evaluation-create.md`.

3. Readers of `RunManifest.id` **MUST** treat it as an opaque token, requiring
   only that it be non-empty; the CLI **MUST NOT** parse a run-segment grammar
   out of it.

   > Rationale: Uniqueness lives in generation, not in a parseable structure;
   > keeping the token opaque avoids reintroducing a grammar the clean break
   > removes.
   > Durable spec: modify `specs/evaluation/records/json-conventions.md`.

4. Generated reports **MUST** render run identity as `Run <NNNN>` using the
   zero-padded `RunManifest.number`, **MUST** surface `RunManifest.id` once as a
   copyable field, and **MUST NOT** label the number as a globally-unique
   identifier.

   > Rationale: The number stays the friendly in-repo handle; the run ID is the
   > handoff anchor and must be present but not dominate. Mislabeling the number
   > as unique is the exact confusion this case removes.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

5. `qualitymd evaluation data set` **MUST** assign a `RecommendationResult.number`
   as a positive integer when a `RecommendationResult` payload omits it, assigning
   sequentially from the owning run's existing recommendation numbers and keeping
   each number unique within the run.

   > Rationale: Recommendations are a flat run-level list of deliverables; a short
   > assigned number is the citable handle, and allocation is mechanical CLI work.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
   > `specs/evaluation/records/json-conventions.md`.

6. `qualitymd evaluation data set` **MUST** preserve a supplied valid
   `RecommendationResult.number` when the payload intentionally rewrites an
   existing recommendation, and **MUST** persist every `RecommendationResult`
   with its `number` present.

   > Rationale: An already-assigned recommendation must be correctable without
   > churning its citable number, and stored data should be self-contained.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md`.

7. `qualitymd evaluation data set` **MUST** derive the `RecommendationResult`
   data path from the zero-padded recommendation `number`, at
   `data/advice/recommendations/<NNN>/recommendation-result.json`.

   > Rationale: Data identity derives from the number now that the `QREC` string
   > is gone; zero-padding keeps directory listings ordered.
   > Durable spec: modify `specs/evaluation/records/data-layout.md`.

8. `RecommendationRankingResult.orderedRecommendations[].recommendationRef` and
   `findingCoverage[].recommendationRefs` **MUST** reference recommendations by
   their assigned `number`, and a write **MUST** be rejected when a referenced
   number has no corresponding `RecommendationResult` in the run.

   > Rationale: Ranking and coverage point at assigned recommendation artifacts;
   > within a run the number is the unambiguous in-run handle.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md`.

9. `FindingRankingResult.orderedFindings[]` **MUST NOT** carry an artifact `id`,
   and `qualitymd evaluation data set` **MUST NOT** assign one; entries continue
   to carry `rank`, `findingRef`, `tier`, and `rationale`.

   > Rationale: A finding already has a good address — its requirement plus
   > payload-local selector. A second, evaluation-global finding ID only restates
   > rank (which churns) and forces stability machinery to keep it from churning.
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md` and
   > `specs/evaluation/records/json-conventions.md`.

10. Cross-payload and report references to a Requirement Finding **MUST** use the
    finding's requirement-scoped reference plus payload-local selector (for
    example `findings[gap-001]`), and report anchors for finding detail **MUST**
    remain derived from the payload-local finding ID.

    > Rationale: With the synthetic artifact ID gone, the requirement-scoped
    > selector is the finding's identity for trace, coverage, and in-report
    > links; existing anchors must not break.
    > Durable spec: modify `specs/evaluation/records/json-conventions.md`.

11. Generated finding index tables **MUST NOT** render a finding artifact-ID
    column, and generated recommendation index tables **MUST** render a `#`
    column for the recommendation number.

    > Rationale: The finding ID column showed an unstable global ordinal; the
    > recommendation number is the genuine citable handle and belongs in the scan
    > surface.
    > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

12. Generated recommendation detail reports **MUST** render the recommendation
    number and a typed reference of the form
    `evaluation:<run-id>/recommendation/<number>`.

    > Rationale: The detail report provides both the short in-run handle and the
    > context-complete reference whose uniqueness is anchored once at the run ID.
    > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
    > `specs/evaluation/records/json-conventions.md`.

13. The CLI **MUST NOT** retain the `QEVAL`/`QREC`/`QFIND` grammar, any reader
    for it, or any migration path for pre-change runs, and the generated report
    gallery **MUST** be regenerated to the new contract.

    > Rationale: Early-alpha clean break — the simpler contract replaces the old
    > one outright rather than carrying a compatibility surface.
    > Durable spec: modify `specs/evaluation/records/json-conventions.md` and
    > `specs/cli/evaluation-data.md`.

## Acceptance criteria

- `qualitymd evaluation create` writes a `RunManifest.id` of the specified form,
  independent of `number`, sharing the `createdAt` instant.
- `qualitymd evaluation data set` assigns missing `RecommendationResult.number`
  values, preserves supplied valid numbers, and persists numbers.
- `RecommendationResult` data paths use `data/advice/recommendations/<NNN>/`.
- Ranking and coverage reference recommendations by number; an unknown number is
  rejected.
- `FindingRankingResult.orderedFindings[]` carries no artifact ID and the CLI
  assigns none; findings resolve only via `findingRef` selectors.
- `report.md`, `findings.md`, `recommendations.md`, and recommendation detail
  reports render `Run <NNNN>`, a copyable run ID, recommendation numbers, and the
  typed recommendation reference; the Top Findings table has no ID column.
- Report anchors for Requirement Finding details remain payload-local.
- No `QEVAL`/`QREC`/`QFIND` token or reader remains in code, specs, or generated
  fixtures.
- The generated Evaluation data JSON Schema is regenerated and current.
- The report gallery is regenerated; durable specs and runtime skill guidance
  describe the new contract.
- Focused Go tests and `mise run check` pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/records/json-conventions.md` - define the run ID and the
  `id`/`number` rule, recommendation references by number, the typed reference,
  and the finding-by-selector rule; remove the `QEVAL`/`QREC`/`QFIND` definitions
  (requirements 1, 3, 5, 9, 10, 12, 13).
- `specs/evaluation/records/payload-kinds.md` - add `RunManifest.id`; replace
  `RecommendationResult.id` with `number`; remove the finding artifact ID; update
  reference contracts (requirements 1, 2, 5, 6, 8, 9).
- `specs/evaluation/records/data-layout.md` - derive the recommendation data path
  from the zero-padded number (requirement 7).
- `specs/evaluation/reports/report-tree.md` - render run identity, the
  recommendation number column and typed reference, and remove the finding ID
  column (requirements 4, 11, 12).
- `specs/cli/evaluation-create.md` - the manifest gains `id` from the creation
  instant (requirements 1, 2).
- `specs/cli/evaluation-data.md` - assignment covers recommendation numbers only;
  the recommendation selector is the number; no finding ID assignment
  (requirements 5, 6, 7, 8, 9, 13).
- `specs/skills/quality-skill/workflows/evaluate.md` - align Advice authoring
  guidance with numbers and finding-by-`findingRef` references (requirements 5,
  8, 9, 10).

### To rename

None

### To delete

None
