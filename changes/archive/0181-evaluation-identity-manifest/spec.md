---
type: Functional Specification
title: Evaluation identity manifest
description: Requirements for making Evaluation identity primary and simplifying generated report frontmatter.
tags: [evaluation, reports, identity]
timestamp: 2026-06-29T00:00:00Z
---

# Evaluation identity manifest

This change updates current Evaluation data and generated report metadata so
Evaluation identity is primary and run numbering is repo-local. It does not
preserve legacy manifest shapes because QUALITY.md is early alpha and current
contracts should stay simpler than compatibility shims.

The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

The current manifest stores a globally unique `id` on `RunManifest`, while the
report frontmatter presents it as `runId`. That naming makes the execution
package feel primary even though handoffs and external systems need a durable
Evaluation ID. Run numbers remain useful as local folder labels, but they are
not durable identifiers. At the same time, report frontmatter duplicates scope
data already owned by the structured manifest.

## Scope

Covered:

- current Evaluation manifest payload kind, path, and fields;
- generated report frontmatter for `report.md`;
- source-data references to the manifest artifact;
- artifact references that include the Evaluation ID;
- current tests, generated schemas, and checked-in report-gallery examples.

Deferred:

- compatibility readers for `RunManifest` or `data/run-manifest.json`;
- CLI command renames;
- run folder naming changes;
- changes to `requestedScope` or `plannedScope` semantics.

## Requirements

1. `qualitymd evaluation create` **MUST** write the CLI-owned manifest to
   `data/evaluation-manifest.json` with `kind: "EvaluationManifest"`.

   > Rationale: The file and payload kind should name the durable object the
   > data describes. The run folder remains an implementation and navigation
   > package around the Evaluation.
   >
   > Durable spec: modify `specs/cli/evaluation-create.md`,
   > `specs/evaluation/records/data-layout.md`, and
   > `specs/evaluation/records/payload-kinds.md`.

2. `EvaluationManifest` **MUST** include `evaluationId`, `createdAt`, `model`,
   `requestedScope`, `plannedScope`, and `run`. `evaluationId` **MUST** be a
   globally unique opaque Evaluation identifier generated from the creation
   instant plus randomness. `run` **MUST** include the repository-local run
   `number` and folder `label`.

   > Rationale: Keeping the run fields nested makes their secondary, local role
   > visible without removing the useful folder number/label.
   >
   > Durable spec: modify `specs/evaluation/records/json-conventions.md`,
   > `specs/cli/evaluation-create.md`, and
   > `specs/evaluation/records/payload-kinds.md`.

3. Evaluation readers and writers **MUST NOT** accept agent-written
   `EvaluationManifest` payloads through `qualitymd evaluation data set`.

   > Rationale: The manifest remains CLI-owned mechanical data; agent judgment
   > belongs in routine outputs.
   >
   > Durable spec: modify `specs/evaluation/records/payload-kinds.md`.

4. Generated `report.md` frontmatter **MUST** contain only `type`, `title`,
   `evaluationId`, `created`, `model`, and `run`. It **MUST NOT** contain
   `runId`, `scope`, or `subject`.

   > Rationale: Scope is already visible in the report body and structured in
   > the manifest. Frontmatter should be compact routing metadata, not a second
   > partial scope representation.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

5. Generated reports that cite manifest source data **MUST** link
   `data/evaluation-manifest.json`, not `data/run-manifest.json`.

   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/evaluation/records/data-layout.md`.

6. Typed artifact references in reports and handoff text **SHOULD** combine the
   Evaluation ID with artifact IDs, using the `evaluation:<evaluationId>/...`
   prefix.

   > Rationale: External references should be anchored on the durable Evaluation
   > identity rather than a repository-local run number.
   >
   > Durable spec: modify `specs/evaluation/records/json-conventions.md`.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/records/json-conventions.md` - requirements 2 and 6 make
  Evaluation ID durable and run numbering local.
- `specs/evaluation/records/data-layout.md` - requirements 1 and 5 rename the
  manifest source-data artifact.
- `specs/evaluation/records/payload-kinds.md` - requirements 1, 2, and 3 update
  the manifest kind, fields, and ownership rule.
- `specs/evaluation/reports/report-tree.md` - requirements 4 and 5 simplify
  frontmatter and source-data references.
- `specs/cli/evaluation-create.md` - requirements 1 and 2 update creation
  output.
- `specs/cli/evaluation-report.md` - requirement 4 keeps report build aligned
  with `report.md` frontmatter.
- `specs/cli/evaluation-list.md`, `specs/evaluation/protocol.md`,
  `specs/evaluation/orchestration.md`, and
  `specs/skills/quality-skill/evaluation.md` - update references to the
  manifest kind and scope source.

### To rename

None

### To delete

None
