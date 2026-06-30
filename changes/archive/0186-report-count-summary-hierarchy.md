---
type: Change Case
title: Report Count Summary Hierarchy
description: Make run-report full-list links and count summaries more scannable with consistent placement, semantic markers, and clearer grouping grammar.
status: Done
tags: [evaluation, reports, recommendations, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Report Count Summary Hierarchy

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0186-report-count-summary-hierarchy/spec.md) - what the case must do.
- [Design doc](0186-report-count-summary-hierarchy/design.md) - how it is built, and why.

## Motivation

The run report already links from capped Top Findings and Top Recommendations
sections to complete `findings.md` and `recommendations.md` reports, but the two
sections do not present those links at the same visual altitude: the Findings
link appears before its capped table, while the Recommendations link appears
after its capped table. Readers should see the complete-list path before reading
either capped preview.

The count summaries beside those links also need stronger hierarchy. Findings
counts should expose their Finding type and severity groups with the same
semantic markers used elsewhere in reports. Recommendation counts need an
explicit impact axis so labels such as `High` and `Medium` are not mistaken for
severity, confidence, or priority. The Strength finding marker should change
from a status-like check mark to the selected muscle-arm marker.

## Scope

Covered:

- Move the full Recommendations report link to the top of the Top
  Recommendations section, matching the full Findings report link placement.
- Add semantic markers to inline full Findings count summaries.
- Replace the hyphen between Finding type counts and gap/risk severity summaries
  with a colon.
- Add semantic markers and an explicit `impact:` axis to inline full
  Recommendations count summaries.
- Change the fixed Finding type display marker for `strength` from `✅` to `💪`.

Non-goals:

- Changing report filenames, generated report navigation, or report frontmatter.
- Changing ranked Finding or Recommendation ordering.
- Changing structured Evaluation enum values or JSON payloads.
- Adding pagination, filtering, or expanded row controls.

## Affected artifacts

**Code**

- [x] `internal/evaluation/display.go` - change the fixed `strength` Finding
      type marker.
- [x] `internal/evaluation/report_tree.go` - move the Recommendations full-list
      link and render marked count summaries.
- [x] `internal/evaluation/evaluation_test.go` - update exact Markdown and enum
      display assertions.

**Durable specs**

- [x] `specs/evaluation/reports/report-tree.md` - specify link placement, inline
      count summary grammar, recommendation impact summaries, and the
      `💪 Strength` display label.
- [x] `specs/log.md` - record the durable spec update.

**Durable docs / generated examples**

- [x] `docs/guides/reporting-design.md` - update report examples and marker
      guidance.
- [x] `docs/log.md` - record the durable docs update.
- [x] `glossary.md` - reflect the fixed `💪 Strength` Finding type label.
- [x] `CHANGELOG.md` - note the report display polish.
- [x] `examples/report-gallery/software-service/.quality/evaluations/**` -
      regenerate report gallery output if generated reports change.

**OKF lifecycle**

- [x] `changes/archive/0186-report-count-summary-hierarchy.md` - archived parent
      Change Case.
- [x] `changes/archive/0186-report-count-summary-hierarchy/` - archived child
      spec and design folder.
- [x] `changes/index.md`, `changes/log.md`, and `changes/archive/index.md` -
      Change Case lifecycle.

## Status

`Done`. Implemented and archived. Generated run reports now place both full-list
links above capped preview tables, use marked inline count summaries, label
recommendation count groups by impact, and render Strength Findings as `💪
Strength`.
