---
type: Design Doc
title: Run Report Entrypoint - design
description: How Evaluation report build separates the run report from root Area report details.
tags: [evaluation, reports, cli, skill]
timestamp: 2026-06-27T00:00:00Z
---

# Run Report Entrypoint - design

## Context

Answers the [functional spec](spec.md) for change case
[0137](../0137-run-report-entrypoint.md). The current renderer treats each
generated Markdown file as an Area, Factor, or Requirement report and special
cases the root Area path to `report.md`. This change introduces one more report
kind - the run report - and gives the root Area its own detail path,
`root-area.md`.

## Approach

### Add a report plan with a headline subject

Keep the existing artifact collector. After collection, derive a small report
plan:

1. load the `EvaluationFrame`;
2. read `inputs.factorIds` first, then `inputs.areaIds`;
3. choose the first listed Factor or Area that has the matching analysis result;
4. fall back to the root Area Analysis Result when no scoped input is recorded;
   and
5. reject report build when no headline result can be resolved.

This keeps scoped reportability tied to recorded run data instead of the run
folder slug, which is explicitly only a human mnemonic.

### Render the run report separately

Prepend a generated report with `kind: run`, path `report.md`, and content from a
new `renderEvaluationRunReport` function. The run report uses only the collected
artifacts and model snapshot. It links to the headline subject report, root Area
report when present, all generated subject reports, and
`data/evaluation-output-result.json`.

The run report summarizes coverage from the Evaluation Frame scope:

- Factor scope when `inputs.factorIds` is populated;
- Area scope when `inputs.areaIds` is populated; and
- root Area/full evaluation fallback otherwise.

It renders an explicit note when the root Area report is absent.

### Rename the root Area report path

Change `areaReportPath(nil)` to return `root-area.md`. All existing navigation
helpers already route root Area links through `areaReportPath(nil)`, so Area,
Factor, and Requirement reports pick up the new root link automatically.

Non-root Area, Factor, and Requirement path functions stay unchanged.

### Extend output contracts

Add `ReportKindRun` to the typed report-kind vocabulary and contract examples.
`EvaluationOutputResult` gains:

- `runReportRef`;
- `headlineResultRef`; and
- `headlineReportRef`.

Keep `rootAreaAnalysisRef` for full/root runs and compatibility with the current
root roll-up concept, but do not require it to be the headline for scoped runs.

The build receipt keeps `reportMd` as the run report path, adds
`headlineReportMd`, and adds `rootAreaReportMd` when the root Area report exists.
The stable `ratingResult` continues to describe the headline result.

### Scope-aware reportability

Replace the fixed root Area reportability check with a minimal graph check:

- `data/frame/evaluation-frame.json` must exist and be valid;
- the headline result implied by the Evaluation Frame must exist and be valid;
- if no headline scope is recorded, the root Area Analysis Result is required.

This is deliberately narrower than full semantic completeness. Existing
structured-data validation and the `/quality` workflow remain responsible for
ensuring all in-scope Requirements, Factors, and Areas are evaluated before
report build.

### Update specs, docs, tests, and schema

Update durable report/data/CLI/skill specs, runtime skill guidance, and the
changelog. Regenerate the Evaluation data schema from the typed data contract.
Add tests for full, scoped sub-Area, and scoped Factor report builds.

## Spec response

- **Report paths and content** - satisfied by the new run report renderer and
  root Area path change.
- **Output indexing and receipts** - satisfied by the extended report kind,
  `EvaluationOutputResult`, and build receipt.
- **Scoped reportability** - satisfied by deriving the headline subject from
  `EvaluationFrame.inputs`.
- **Verification** - satisfied by focused report tree, output JSON, receipt, and
  scoped reportability tests.

## Alternatives

- **Keep root Area at `report.md` and add `run-report.md`.** Rejected. Users
  expect `report.md` to be the entrypoint; the more specific artifact should
  take the more specific filename.
- **Use the run folder narrowing slug to derive scope.** Rejected. The slug is a
  human mnemonic and loses Area-vs-Factor boundaries. Structured inputs are the
  durable source.
- **Require root Area analysis for every scoped run.** Rejected. It makes
  targeted evaluations pay for and report a subject they did not evaluate.
- **Remove `rootAreaAnalysisRef` immediately.** Rejected for this slice. Adding
  headline refs solves scoped reports without forcing every consumer to migrate
  off the root roll-up field in the same change.

## Trade-offs & risks

- This is a breaking path change for consumers that open the root Area report at
  `report.md`.
- Minimal legacy Evaluation Frames that omit scope continue to require root Area
  analysis, which is intentional fallback behavior but means targeted runs must
  record scope in the frame to avoid the root requirement.
- `EvaluationOutputResult` grows rather than replaces `rootAreaAnalysisRef`.
  That avoids a larger consumer break now, but a later cleanup may remove the
  root-specific field.

## Open questions

None.
