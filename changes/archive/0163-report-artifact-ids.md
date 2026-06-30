---
type: Change Case
title: Report Artifact IDs
description: Give handoff-ready Evaluation report artifacts stable user-citable IDs while keeping Model references distinct.
status: Done
tags: [evaluation, reports, advice, identifiers]
timestamp: 2026-06-29T00:00:00Z
---

# Report Artifact IDs

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0163-report-artifact-ids/spec.md) - what the case must do.
- [Design doc](0163-report-artifact-ids/design.md) - how it is built, and why.

## Motivation

Evaluation reports produce handoff artifacts: recommendations are assigned to
people or external trackers, findings are cited in chats and reviews, and a run
itself anchors provenance. Sluggy recommendation IDs such as
`001-refresh-pinned-installer-examples` are readable file names, but they mix
identity, title wording, and path concerns. Finding IDs remain payload-local,
which is correct for structured assessment records but weak for report-level
handoff.

QUALITY.md already has a coherent Model identity style: IDs are scoped, typed
references cross boundaries, display labels are not identity, and paths are
derived. Evaluation report artifacts need the same discipline without confusing
report artifact IDs with Model references.

## Scope

Covered:

- Define user-citable Evaluation artifact IDs for runs, recommendations, and
  ranked findings.
- Use `QEVAL-<NNNN>` for run identity rendered from `RunManifest.number`.
- Make the CLI assign `RecommendationResult.id` values as
  `QREC-<NNNN>-<NNN>` when recommendation payloads are written.
- Make the CLI assign `FindingRankingResult.orderedFindings[].id` values as
  `QFIND-<NNNN>-<NNN>` when finding ranking payloads are written.
- Keep Requirement Finding `id` values payload-local and keep cross-payload
  finding references based on `findingRef` plus selector.
- Add visible `ID` columns or summary fields to generated finding and
  recommendation report surfaces.
- Keep Model references (`area:`, `factor:`, `requirement:`, `rating:`)
  unchanged and label typed artifact references as references, not IDs.

Deferred:

- Adding globally unique UUIDs, URNs, or repository-qualified external IDs.
- Migrating historical Evaluation runs.
- Adding public IDs for candidate actions.
- Adding interactive lookup commands for artifact IDs.
- Changing Area, Factor, Requirement, or Rating Level identity.

## Affected artifacts

Derived by sweeping for `RecommendationResult.id`, `FindingRankingResult`,
`orderedFindings`, `recommendationRef`, `findingCoverage`, `recommendations/`,
Top Findings, Top Recommendations, and report-output references.

**Code**

- [x] `internal/evaluation/data_contract.go` - add artifact ID fields,
      validation, and assignment-aware input contracts.
- [x] `internal/evaluation/data.go` - update data path/query semantics,
      assignment normalization, effective validation, and examples.
- [x] `internal/evaluation/report_tree.go` - render artifact IDs and typed
      references in report outputs.
- [x] `internal/evaluation/evaluation_test.go` - update fixtures and assertions.
- [x] `internal/evaluation/evaluation-data.schema.json` - regenerate the
      published Evaluation data schema.

**Durable specs**

- [x] `specs/evaluation/records/json-conventions.md` - define Evaluation
      artifact ID and reference conventions.
- [x] `specs/evaluation/records/payload-kinds.md` - update
      `FindingRankingResult`, `RecommendationResult`, and
      `RecommendationRankingResult` contracts.
- [x] `specs/evaluation/records/data-layout.md` - align recommendation data
      paths with citable IDs and slug-only report filenames.
- [x] `specs/evaluation/reports/report-tree.md` - render citable IDs and
      references across report tables/detail pages.
- [x] `specs/cli/evaluation-data.md` - define `data set` assignment before
      validation/path derivation.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` - align runtime
      recommendation and finding-ranking authoring guidance.

**Durable docs / runtime guidance**

- [x] `skills/quality/SKILL.md` - align runtime Advice authoring guidance.
- [x] `skills/quality/workflows/evaluate.md` - align workflow-specific runtime
      guidance if present.
- [x] `docs/guides/reporting-design.md` - update report header examples for
      `QEVAL`/`QREC` ID rendering.
- [x] `CHANGELOG.md` - release-note entry for the clean-break artifact ID
      contract.
- [x] `changes/index.md` and `changes/log.md` - Change Case lifecycle.

## Status

`Done`. Implemented and archived. `qualitymd evaluation data set` now assigns
handoff-ready `QREC` and `QFIND` artifact IDs, reports render `QEVAL`, `QREC`,
and `QFIND` references, durable specs and runtime guidance are aligned, and the
generated report gallery reflects the new contract.
