---
type: Design Doc
title: Quality Evaluation Title - design doc
description: Design for deriving generated run report titles from scoped Area and factor filter labels.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Quality Evaluation Title - design doc

Design for [Quality Evaluation Title](../0168-quality-evaluation-title.md) and
its [functional spec](spec.md).

## Context

The Evaluation report tree already resolves the planned Area and factor filter
before rendering `report.md`. The current run report renderer calls
`evaluationScopedAreaLabel`, then prepends `Evaluation Report: Area:` to build
the visible title. The same heading feeds frontmatter `title`, so changing the
heading source changes both frontmatter and H1 together.

## Approach

Replace the run report heading construction with a helper that returns the full
document title:

- `Quality Evaluation - <Area title>` for Area-only runs;
- `Quality Evaluation - <Area title> (<Factor titles>)` for factor-filtered
  runs.

The helper uses the same Area and Factor title resolvers already used in scope
tables and report links. Factor titles are gathered from
`plan.FactorFilter`, which is resolved from
`RunManifest.plannedScope.factorFilter`, so ordering follows the manifest
without a separate sorting rule.

The existing `renderReportHeader` path still writes frontmatter `title` from
the H1 heading. The Scope section remains unchanged and continues to show stable
Area and Factor references.

## Spec response

- Area-only titles come from the scoped Area title with the fixed
  `Quality Evaluation -` prefix.
- Factor-filtered titles append a parenthesized comma-separated list of all
  resolved Factor titles.
- Stable references stay in the Scope table rather than the H1.
- Detail report renderers are untouched.

## Alternatives

- **Use `Evaluation - <Area>`.** Rejected because it is shorter but too generic
  in generated artifacts and browser/editor titles.
- **Use `Quality Evaluation - <Factor> in <Area>`.** Rejected because it scales
  poorly to multiple factors and makes the Area less visually stable.
- **Truncate to `<N> factors`.** Rejected because likely factor-filter counts
  are small and truncation hides user-requested scope.
- **Keep `Evaluation Report` and only remove `Area:`.** Rejected because
  `Report` is still redundant in the report artifact's visible title.

## Trade-offs & risks

Long Factor titles can produce a longer H1 for multiple factor filters. That is
acceptable because the title remains exact and typical factor-filter runs are
expected to include only a few Factors.

## Open questions

None.
