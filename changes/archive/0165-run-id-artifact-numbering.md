---
type: Change Case
title: Run IDs and Artifact Numbering
description: Replace locally-scoped QEVAL/QREC/QFIND artifact IDs with a globally-unique run ID plus per-run recommendation numbers, and drop the finding artifact ID.
status: Done
tags: [evaluation, reports, advice, identifiers]
timestamp: 2026-06-29T00:00:00Z
---

# Run IDs and Artifact Numbering

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0165-run-id-artifact-numbering/spec.md) - what the case must do.
- [Design doc](0165-run-id-artifact-numbering/design.md) - how it is built, and why.

## Motivation

[0163 - Report Artifact IDs](0163-report-artifact-ids.md) gave runs,
recommendations, and ranked findings citable `QEVAL`/`QREC`/`QFIND` IDs. Every
one of those tokens is scoped to a single field: `RunManifest.number`, which the
CLI derives by scanning the local `.quality/evaluations/` directory and taking
`max + 1`. That makes the IDs unique only inside one persisted run directory and
nowhere above it, which produces two collisions once an artifact leaves the repo:

- **Cross-project (loud).** Project A's `QREC-0001-002` equals project B's. A
  human notices because the two reports are visibly about different entities.
- **Ad-hoc reproduction (silent, dangerous).** Because the run counter is
  _derived from the directory_, it is not durable unless the directory is. An
  uncommitted or ephemeral run — a fresh clone, a CI scratch workspace, a
  teammate re-running locally — restarts at `0001`. Two _different_ runs then
  mint the same `QEVAL-0001` / `QREC-0001-002` with conflicting content and
  nothing to tell them apart.

Evaluation outputs are handoff artifacts: a recommendation is assigned to a
person or tracker, a run anchors provenance. So identity needs a uniqueness
guarantee that survives leaving the repo, _without_ sacrificing the short,
sayable handle that motivated 0163.

This case revises 0163's identity model around one rule: **`id` names only the
genuinely-unique thing (the run); everything scoped to a run is a `number`;
findings are addressed in context and carry no synthetic artifact ID.** The run
gets a highly-improbable, globally-unique `id`; recommendations keep a short
per-run `number`; the finding artifact ID is removed because a finding already
has a good address — its requirement plus payload-local selector.

## Scope

Covered:

- Add a globally-unique `RunManifest.id` of the form `<timestamp>-<nanoid>`,
  assigned once at run creation, distinct from the friendly `number`.
- Keep `RunManifest.number` as the friendly, directory-naming, report-header
  identity; do not replace it with the run ID.
- Replace `RecommendationResult.id` (`QREC-<NNNN>-<NNN>`) with a per-run
  `number`; derive the recommendation data path from the zero-padded number.
- Reference recommendations by `number` in `RecommendationRankingResult` and
  finding coverage.
- Remove `FindingRankingResult.orderedFindings[].id` (`QFIND-<NNNN>-<NNN>`)
  entirely; address findings only by `findingRef` plus payload-local selector.
- Render run identity as `Run <NNNN>` with the run ID surfaced as a copyable
  field; render a recommendation `#` column and a typed reference
  `evaluation:<run-id>/recommendation/<number>`; drop the Top Findings ID column.
- Retire the `QEVAL`/`QREC`/`QFIND` grammar as a clean break and regenerate the
  report gallery.

Deferred:

- Surfacing the run ID in `qualitymd evaluation list`/`show` output.
- Repository- or entity-qualified namespaces above the run ID (the improbable
  run ID makes them unnecessary for uniqueness).
- Public IDs for candidate actions.

Non-goals:

- Migrating pre-change Evaluation runs on disk, or any backward-compatibility
  shim, fallback reader, or dual writer for the old `QEVAL`/`QREC`/`QFIND`
  contract.
- Changing Area, Factor, Requirement, or Rating Level Model identity.
- Changing Requirement Finding payload-local IDs, selectors, or report anchors.

## Affected artifacts

Derived by sweeping for `QEVAL`, `QREC`, `QFIND`, `RecommendationResult`,
`orderedFindings`, `recommendationRef`, `findingCoverage`, `recommendations/`,
`RunManifest`, `evaluationArtifactID`, `runArtifactID`, `parseArtifactID`, and
report run-line / Top Findings rendering.

**Code**

- [x] `internal/evaluation/types.go` - add `RunManifest.id`; shared
      create-time clock for `createdAt` and the run ID timestamp.
- [x] `internal/evaluation/create.go` - generate the run ID at creation.
- [x] `internal/evaluation/data.go` - drop finding artifact-ID assignment and
      the stability-by-`findingRef` machinery; rename recommendation `id`→`number`
      assignment; derive recommendation paths/queries from the zero-padded
      number; collapse the `QEVAL`/`QREC`/`QFIND` parse/format helpers.
- [x] `internal/evaluation/data_contract.go` - recommendation `number` contract;
      remove finding `id` contract; recommendation-ref-by-number validation.
- [x] `internal/evaluation/report_tree.go` - render `Run <NNNN>` + copyable run
      ID, recommendation `#` column and typed reference; remove the finding
      artifact-ID column.
- [x] `internal/evaluation/evaluation-data.schema.json` - regenerate from the
      typed contract (run ID added, finding ID removed, recommendation
      `number`).
- [x] `internal/evaluation/evaluation_test.go`,
      `internal/cli/evaluation_test.go` - update fixtures and assertions.
- [x] `scripts/report-gallery/main.go` - pin a deterministic run ID alongside
      the pinned `createdAt`.

**Durable specs**

- [x] `specs/evaluation/records/json-conventions.md` - define the run ID, the
      `id`/`number` rule, recommendation references by number, and the typed
      reference; remove the `QEVAL`/`QREC`/`QFIND` definitions.
- [x] `specs/evaluation/records/payload-kinds.md` - `RunManifest.id`;
      `RecommendationResult.number`; remove finding artifact ID; reference
      contracts.
- [x] `specs/evaluation/records/data-layout.md` - recommendation data path from
      the zero-padded number.
- [x] `specs/evaluation/reports/report-tree.md` - run identity, recommendation
      number column and typed reference, removed finding ID column.
- [x] `specs/cli/evaluation-create.md` - manifest gains `id`.
- [x] `specs/cli/evaluation-data.md` - assignment covers recommendation numbers
      only; selector is the number.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` - Advice authoring
      against numbers; findings referenced by `findingRef` only.

**Durable docs / runtime guidance**

- [x] `skills/quality/SKILL.md` - align runtime Advice authoring guidance.
- [x] `skills/quality/workflows/evaluate.md` - align workflow runtime guidance.
- [x] `docs/guides/reporting-design.md` - update report identity examples.
- [x] `CHANGELOG.md` - clean-break release note.
- [x] `changes/index.md` and `changes/log.md` - Change Case lifecycle.

**Generated fixtures**

- [x] `examples/report-gallery/software-service/.quality/evaluations/**` -
      regenerated via `scripts/report-gallery`.

## Status

`Done`. Implemented and archived. Evaluation runs now persist a
globally-unique `RunManifest.id`; recommendations carry per-run `number`
values; finding ranking entries carry no artifact ID; reports, schema, durable
specs, runtime guidance, and the report gallery are regenerated to the new
contract. No backward-compatibility path is planned — pre-change runs are not
migrated.
