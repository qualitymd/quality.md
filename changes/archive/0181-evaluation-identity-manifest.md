---
type: Change Case
title: Evaluation Identity Manifest
description: Make Evaluation identity primary, keep run numbering secondary, and simplify report frontmatter.
status: Done
tags: [evaluation, reports, identity]
timestamp: 2026-06-29T00:00:00Z
---

# Evaluation Identity Manifest

Generated Evaluation data currently exposes the durable identifier as a run ID
inside `RunManifest`, while reports carry both durable identity and duplicated
scope routing fields. This makes "run" read like the primary object even though
the durable user and integration concept is the Evaluation.

- [Functional spec](0181-evaluation-identity-manifest/spec.md) - what the
  Evaluation manifest and report frontmatter must change.
- [Design doc](0181-evaluation-identity-manifest/design.md) - how the Evaluation
  package, specs, skill guidance, and examples absorb the identity cleanup.

## Motivation

External systems, handoffs, and generated reports should identify an Evaluation
by an `evaluationId`. Run numbering remains useful for repository-local folder
ordering and editor navigation, but it should not be the durable identity. The
manifest should stay the complete structured source for requested and planned
scope, while report frontmatter should carry only compact routing metadata.

## Scope

Covered:

- CLI-created Evaluation manifest kind, durable ID field, and run numbering
  shape;
- generated report frontmatter and primary source-data references;
- structured Evaluation data schema and validation;
- typed recommendation artifact references that include the Evaluation ID;
- durable Evaluation, CLI, and skill specs plus report design guidance;
- runtime `/quality` skill guidance for using CLI-owned Evaluation manifests;
- tests and checked-in generated report examples.

Deferred:

- migration or compatibility readers for existing `RunManifest` data;
- renaming CLI commands or user-facing "run folder" path terminology;
- removing `requestedScope` or `plannedScope` from the manifest;
- changing the generated run folder naming convention.

## Affected artifacts

- Code:
  - `internal/evaluation/types.go` - rename the manifest type and ID field,
    and nest run number/label metadata.
  - `internal/evaluation/create.go`, `manifest.go`, `path.go`, `list.go`,
    `report_tree.go`, `data.go`, `data_contract.go`, and `display.go` - consume
    `EvaluationManifest` from `data/evaluation-manifest.json`.
  - `internal/evaluation/evaluation_test.go`, `internal/cli/evaluation_test.go`,
    and generated schemas - update expectations.
  - `scripts/report-gallery/main.go` - rewrite generated gallery manifest IDs.
- Durable specs:
  - `specs/evaluation/records/json-conventions.md` - make Evaluation ID the
    durable identity and run number local.
  - `specs/evaluation/records/data-layout.md` and
    `specs/evaluation/records/payload-kinds.md` - rename the manifest artifact
    and payload kind.
  - `specs/evaluation/reports/report-tree.md` and
    `specs/cli/evaluation-report.md` - simplify report frontmatter and source
    data references.
  - `specs/cli/evaluation-create.md`, `specs/cli/evaluation-list.md`,
    `specs/evaluation/protocol.md`, `specs/evaluation/orchestration.md`, and
    `specs/skills/quality-skill/evaluation.md` - update manifest references.
- Runtime skill:
  - `skills/quality/SKILL.md`, `skills/quality/workflows/evaluate.md`, and
    `skills/quality/resources/cli-workflow-conventions.md` - use the
    Evaluation manifest terminology.
- Durable docs:
  - `docs/guides/reporting-design.md` - update report frontmatter and source
    data examples.
- Generated examples:
  - `examples/report-gallery/` - regenerate generated manifests and reports.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md`, `specs/evaluation/log.md`, `skills/quality/log.md`,
    `skills/quality/workflows/log.md`, and `docs/log.md` where affected bundles
    record updates.

## Status

`Done`. Implemented and archived. Evaluation manifests now use
`EvaluationManifest`, `evaluationId`, nested local `run` metadata, and
`data/evaluation-manifest.json`; generated run report frontmatter now carries
`evaluationId`, `created`, `model`, and `run` without duplicated scope fields.
