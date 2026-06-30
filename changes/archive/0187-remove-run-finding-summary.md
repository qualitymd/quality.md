---
type: Change Case
title: Remove Run Finding Summary
description: Remove the redundant run-level Finding Summary table from generated report.md now that the full findings link carries the count breakdown.
status: Done
tags: [evaluation, reports, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Remove Run Finding Summary

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0187-remove-run-finding-summary/spec.md) - what the case must do.
- [Design doc](0187-remove-run-finding-summary/design.md) - how it is built, and why.

## Motivation

Generated run-level `report.md` now shows Finding and Recommendation totals in
`## Key Details`, then shows a full Findings report link with the complete ranked
Finding count and inline type/severity breakdown immediately under
`## Top Findings`. The separate `Finding Summary` table between those two
surfaces repeats the same count hierarchy without linking to the full report,
making the opening stack longer and harder to scan.

The run report should keep one compact count in `Key Details`, keep the richer
breakdown beside the full `findings.md` link, and remove the duplicate summary
table from the opening.

## Scope

Covered:

- Remove the run-level `Finding Summary` table from generated `report.md`.
- Keep the `Findings` total in `## Key Details`.
- Keep the full `findings.md` link and its inline complete count summary under
  `## Top Findings`.
- Keep Requirement report `## Findings Summary` sections unchanged.

Non-goals:

- Changing ranked Finding ordering, filtering, table columns, or row caps.
- Changing `findings.md`, Requirement reports, or structured Evaluation data.
- Changing generated report filenames or `Evaluation links:` navigation.

## Affected artifacts

**Code**

- [x] `internal/evaluation/report_tree.go` - remove run-report Finding Summary
      rendering and unused helper code.
- [x] `internal/evaluation/evaluation_test.go` - update generated report
      assertions.

**Durable specs**

- [x] `specs/evaluation/reports/report-tree.md` - remove the run-report Finding
      Summary table contract while preserving Key Details and full-list link
      count contracts.
- [x] `specs/cli/evaluation-report.md` - update the `report.md` opening stack.
- [x] `specs/evaluation/log.md` and `specs/log.md` - record the durable spec
      updates.

**Durable docs / generated examples**

- [x] `docs/guides/reporting-design.md` - update run-report opening guidance and
      examples.
- [x] `docs/log.md` - record the durable docs update.
- [x] `CHANGELOG.md` - note the generated report cleanup.
- [x] `examples/report-gallery/software-service/.quality/evaluations/**` -
      regenerate report gallery output.

**OKF lifecycle**

- [x] `changes/archive/0187-remove-run-finding-summary.md` - archived parent
      Change Case.
- [x] `changes/archive/0187-remove-run-finding-summary/` - archived child spec
      and design
      folder.
- [x] `changes/index.md`, `changes/log.md`, and `changes/archive/index.md` -
      Change Case lifecycle.

## Status

`Done`. Implemented and archived. Generated run reports no longer render the
standalone `Finding Summary` table near `## Key Details`; the total count remains
in Key Details and the complete type/severity breakdown remains beside the full
`findings.md` link.
