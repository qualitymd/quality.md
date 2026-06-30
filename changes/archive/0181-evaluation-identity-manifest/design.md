---
type: Design Doc
title: Evaluation identity manifest design
description: Design for renaming the Evaluation manifest, nesting run metadata, and simplifying report frontmatter.
tags: [evaluation, reports, identity]
timestamp: 2026-06-29T00:00:00Z
---

# Evaluation identity manifest design

## Context

This design answers the [Evaluation identity manifest](spec.md) functional spec.
The change separates durable Evaluation identity from local run numbering while
leaving the current scope model intact.

## Approach

The implementation should keep the existing Evaluation package flow and update
the data shape at the edges:

- rename `RunManifest` to `EvaluationManifest`;
- rename `id` to `evaluationId`;
- replace top-level `number` with nested `run.number` and `run.label`;
- change the manifest path constant to `data/evaluation-manifest.json`;
- keep `requestedScope` and `plannedScope` unchanged;
- update report frontmatter generation to emit `evaluationId`, `created`,
  `model`, and `run`, removing scope duplication.

Run folder naming still uses the same next-number and slug logic. After the
folder name is known, `evaluation create` can populate both `run.number` and
`run.label` before writing the manifest. Status and list readers use
`manifest.Run.Number` for ordering and `manifest.Run.Label` for display or
frontmatter.

The report renderer already routes manifest-derived facts through a small set of
helpers. Updating those helpers and source-data path functions should propagate
the renamed source-data link across run, area, factor, requirement, findings,
and recommendation reports.

Generated schemas should continue to come from the existing Evaluation data
contract declarations. The contract entry changes kind/path/field names while
preserving the existing CLI-owned write rejection.

## Spec response

- Manifest creation writes the renamed payload and path without changing the
  folder creation lifecycle.
- Readers, status, list, validation, and report build continue to require the
  CLI-owned manifest, now as `EvaluationManifest`.
- Report frontmatter becomes compact and avoids duplicating scope fields that
  remain in the manifest and visible report body.
- Typed recommendation artifact references use `manifest.EvaluationID`, keeping
  external references anchored on Evaluation identity.

## Alternatives

- **Keep `RunManifest` and only rename `runId` in report frontmatter.** Rejected
  because the structured source of truth would still teach integrations that
  the run is primary.
- **Use `localRun` instead of `run`.** Rejected because the nested object already
  makes the run metadata secondary, and `run` is clearer in user-facing
  frontmatter.
- **Remove scope from the manifest.** Rejected because `requestedScope` and
  `plannedScope` remain useful provenance and deterministic report-build input.
  The problem is duplicated scope in report frontmatter, not scope in structured
  data.
- **Rename run folders to include the Evaluation ID.** Deferred because folder
  numbers are useful local navigation, and changing folder names would add churn
  without improving durable references.

## Trade-offs and risks

This is a breaking Evaluation data-contract change. That is acceptable under the
project's early-alpha compatibility policy and avoids dual readers or fallback
writers.

The word "run" remains in CLI path terminology and folder names. That is
deliberate: run remains the local artifact package, while Evaluation identity is
primary in structured data and generated reports.

## Open questions

None.
