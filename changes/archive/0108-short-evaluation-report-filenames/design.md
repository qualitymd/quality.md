---
type: Design Doc
title: Short evaluation report filenames - design
description: Design for preserving the root Evaluation v2 report.md entrypoint while replacing repeated descendant report.md filenames with short subject-aware generated Markdown filenames.
tags: [cli, evaluation, reports]
timestamp: 2026-06-26T00:00:00Z
---

# Short evaluation report filenames - design

This design answers the
[Short evaluation report filenames functional spec](spec.md).

## Context

Evaluation v2 report generation already centralizes generated Markdown paths in
the report renderer helpers. The current helpers return `report.md` at each
report folder, and all navigation, table links, and report refs derive from
those helper results.

The change should preserve that shape: keep report path derivation centralized,
change only the generated Markdown filenames, and let existing relative-link
logic continue to calculate links between generated reports.

## Approach

Update the report path helper functions in `internal/evaluation/report_v2.go`:

- `areaReportPath(nil)` keeps returning `report.md`.
- `areaReportPath(nonRoot)` appends `<local-area>-area.md` instead of
  `report.md`.
- `requirementReportPath` appends `<requirement>-requirement.md` under the
  existing `requirements/<requirement>/` folder.
- `factorReportPath` appends `<local-factor>-factor.md` under the existing
  Factor folder, including nested Factors.

The helper functions already return repository-slash paths relative to the run
root. Existing writers, breadcrumbs, table subject links, and report refs should
continue consuming those helper functions rather than learning filename rules
inline.

Update `internal/evaluation/data.go` where example or fallback
`EvaluationOutputResult` data still names the root report. Since the root Area
remains `report.md`, this file is likely a review point rather than a broad
rewrite.

Update tests by replacing path and link expectations with the new descendant
filenames. Add at least one assertion that no descendant generated Markdown file
named `report.md` exists, while the root `report.md` still does.

Durable specs and runtime skill guidance get wording-only updates after the code
path is clear, so the public contract and skill artifact summary match the
generated tree.

## Spec Response

The filename contract is satisfied by changing the path helpers at the only
report-path derivation point. Root Area special-casing stays in
`areaReportPath(nil)`, so `report.md` remains the run entrypoint.

The links-and-indexes requirements are satisfied because report rendering
continues to call the same path helpers for write paths, breadcrumb paths, table
subject links, and report refs.

The compatibility boundary is satisfied by not writing any additional files:
the writer emits exactly the paths returned by the report collection, and no
legacy-copy layer is added.

## Alternatives

### Rename every report, including root

Rejected. `root-area.md` would make every tab subject-aware, but the root report
is already unique within a run and `report.md` is a useful conventional
entrypoint. Keeping root stable reduces migration noise while solving the
multi-tab ambiguity for descendant reports.

### Use display titles in filenames

Rejected. Titles are better human labels but are mutable and can contain
characters that need escaping or normalization. Structural IDs are already the
path contract and keep generated paths deterministic.

### Use full structural paths in filenames

Rejected. Full-path filenames reduce ambiguity inside tabs but make filenames
long. The directory tree already carries full identity; local subject-aware
filenames are enough to distinguish report kind and subject in common tab views.

### Add compatibility copies

Rejected. Duplicate Markdown outputs would complicate report refs, make stale
links harder to detect, and weaken the report tree as a deterministic
projection.

## Trade-offs & Risks

Existing generated runs keep their old paths; this change does not migrate them.
Users comparing old and new runs will see both shapes in history.

Any caller with hard-coded descendant `report.md` paths will need to follow
generated report refs or rebuild links from the documented path contract. That
is acceptable because generated report paths are already part of the Evaluation
v2 report output contract.

## Open Questions

None.
