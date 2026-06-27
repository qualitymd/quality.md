---
type: Change Case
title: Report Descendant Terms
description: Align generated report labels around Child Areas and Sub-Factors for immediate descendants.
status: Done
tags: [evaluation, reports, terminology]
timestamp: 2026-06-27T00:00:00Z
---

# Report Descendant Terms

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0147-report-descendant-terms/spec.md) - what the case must
  do.
- [Design doc](0147-report-descendant-terms/design.md) - how it's built, and why.

## Motivation

Generated Evaluation reports currently mix two naming families for immediate
descendants: Area reports render `Sub-Areas`, while Factor reports render
`Child Factors`. The Model vocabulary already uses Child Areas for immediate
descendant Areas and sub-factors for immediate descendant Factors. Reports should
use that same vocabulary so readers do not have to translate between equivalent
relationship names.

## Scope

Covered:

- Rename generated Area report descendant-Area labels from `Sub-Areas` to
  `Child Areas`.
- Rename generated Factor report descendant-Factor labels from `Child Factors`
  to `Sub-Factors`.
- Align report breakdown table labels, empty-state rows, tests, and durable
  report specs.
- Update release notes for the user-facing report wording change.

Deferred:

- Renaming structured data fields such as `childAreaAnalysisRefs` or
  `childFactorAnalysisRefs`.
- Renaming internal helper functions such as `childAreas` or `childFactors`.
- Rewriting historical Change Case archives or append-only historical log text.
- Changing Factor or Area roll-up semantics.

## Affected artifacts

Derived by sweeping for `Sub-Areas`, `Child Factors`, `child Factors`,
`child Areas`, `Sub-Factors Rating`, `Sub-Areas Rating`, and report-tree
references.

**Code**

- [x] `internal/evaluation/report_tree.go` - generated Markdown report labels.
- [x] `internal/evaluation/evaluation_test.go` - report-rendering assertions and
      rendered fixture rationale text.

**Durable specs** (substance in the [functional spec](0147-report-descendant-terms/spec.md))

- [x] `specs/evaluation/reports/report-tree.md` - report tree label contract.
- [x] `specs/skills/quality-skill/reporting.md` - skill report navigation
      terminology.

**Durable docs / bundled skill runtime**

- [x] `CHANGELOG.md` - release-note entry for the user-facing report label
      change.
- [x] `specs/log.md` - durable specs bundle update log.
- [x] `changes/log.md` and `changes/archive/index.md` - Change Case lifecycle.

## Status

`Done`. Implemented and archived after generated report labels, durable specs,
tests, release notes, and logs were updated. `go test ./...` and
`mise run fmt-md-check` pass.
