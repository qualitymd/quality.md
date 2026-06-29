---
type: Design Doc
title: Report source-data frontmatter — design doc
description: Design for source-data frontmatter lists and report body header cleanup.
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report source-data frontmatter — design doc

Design for
[Report source-data frontmatter](../0159-report-source-data-frontmatter.md) and
its [functional spec](spec.md).

## Context

The report renderer already has a shared `reportHeader` path that writes
frontmatter and report-specific summary tables. Each report renderer currently
passes both `Data` frontmatter values and visible `Data` summary-table columns.
The same code path can keep frontmatter while removing body duplication.

## Approach

Keep `reportHeader.Data` as the only report-level source-data pointer surface.
Replace each visible summary table's `Data` column with report-specific state
columns only:

- run report: `Overall Rating`, `Scope`, `Confidence`;
- Area report: `Overall Rating`, `Local Rating`, `Confidence`;
- Factor report: `Overall Rating`, `Local Rating`, `Status`, `Confidence`;
- Requirement report: `Rating`, `Assessment`, `Confidence`;
- findings index: `Findings`, `Highest Severity`;
- recommendations index: `Recommendations`, `Highest Impact`, `Coverage`;
- recommendation detail report: `Rank`, `Impact`, `Confidence`.

Add small source-data helper functions in `internal/evaluation/report_tree.go`
that build de-duplicated, deterministic frontmatter lists from the report's
actual structured inputs. Every report with visible run context includes
`data/run-manifest.json`. Subject reports add their own result payloads and the
direct child or referenced payloads whose data appears in body tables. Index and
run reports add the ranking, recommendation, and referenced assessment/result
payloads their rows render from. The generated
`data/evaluation-output-result.json` output is no longer listed unless a future
renderer reads from it.

## Spec response

- Frontmatter `data` remains generated in one shared header renderer, but the
  values passed to it now represent source payloads rather than generic related
  payloads.
- Body `Data` duplication disappears by shrinking report-specific summary table
  header and row slices.
- Run context remains visible because the header still renders from
  `RunManifest`; the manifest is listed in frontmatter as source data.
- Determinism comes from stable data path helpers, existing sorted artifact
  accessors, and de-duplicating lists while preserving first-use order.
- Generated report frontmatter remains output-only; no report renderer reads it.

## Alternatives

- **Keep `Data` columns and reinterpret frontmatter as only metadata.** Rejected
  because the user wants report data info in frontmatter, not duplicated in
  report bodies.
- **List every payload in the run on every report.** Rejected because `data`
  would stop being report-local source information and would become a noisy run
  index.
- **Add another frontmatter key for source data.** Rejected because existing
  `data` is already the OKF-compatible place for payload pointers; the fix is to
  sharpen its meaning.

## Trade-offs & risks

- Some report frontmatter lists become longer than before because they now list
  the actual payloads rendered in body tables. That is acceptable because the
  frontmatter is the source-data manifest, but the lists should stay
  report-local rather than run-global.
- Reports also render display vocabulary from `model-snapshot.md`. This design
  keeps `data` limited to structured Evaluation payloads under `data/`; the
  durable spec should name that boundary so readers do not expect
  `model-snapshot.md` there.

## Open questions

None.
