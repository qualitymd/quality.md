---
type: Change Case
title: Report Full List Links
description: Make run-report links to full findings and recommendations lists more scannable and show total ranked counts.
status: Done
tags: [evaluation, reports, recommendations, findings]
timestamp: 2026-06-29T00:00:00Z
---

# Report Full List Links

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0177-report-full-list-links/spec.md) - what the case must do.
- [Design doc](0177-report-full-list-links/design.md) - how it is built, and why.

## Motivation

The run report shows capped `Top Findings` and `Top Recommendations` tables,
then links to the complete list reports. The current follow-up lines read as
plain sentence text, use title-style casing inside a sentence, and do not say
whether the full list contains more rows than the capped table. Readers scanning
`report.md` should be able to notice the complete-list link and know the total
ranked count without opening another file first.

## Scope

Covered:

- Emphasize the full Findings and Recommendations report labels in `report.md`.
- Use sentence-case labels for those lines.
- Show the total ranked finding and recommendation counts next to the full-list
  links.
- Keep the linked artifact label as the report filename.

Deferred:

- Changing the Top Findings or Top Recommendations table columns, rank cap, or
  row ordering.
- Changing `findings.md` or `recommendations.md` content beyond any generated
  gallery refresh caused by `report.md`.

Non-goals:

- Adding interactive pagination, sorting, filtering, or hidden row expansion.
- Changing report filenames, report navigation, or report frontmatter.

## Affected artifacts

Derived by sweeping for `Full Recommendations report`, `Full Findings report`,
`recommendations.md`, run report ranked-list output, and report gallery output.

**Code**

- [x] `internal/evaluation/report_tree.go` - render emphasized sentence-case
      full-list labels with total ranked counts.
- [x] `internal/evaluation/evaluation_test.go` - update run-report assertions.
- [x] `scripts/report-gallery/main.go` and generated gallery output if the
      checked-in report gallery changes.

**Durable specs**

- [x] `specs/evaluation/reports/report-tree.md` - specify the emphasized
      full-list links and total counts in `report.md`.

**Durable docs / runtime guidance**

- [x] `docs/guides/reporting-design.md` - align the run-report header pattern
      with emphasized complete-list links and counts.
- [x] `CHANGELOG.md` - add a current-release note if release notes track report
      output polish.

## Status

`Done`. Implemented and archived. Generated run reports now emphasize the full
Findings and Recommendations report links below capped top-list tables and show
the complete ranked count for each linked report.
