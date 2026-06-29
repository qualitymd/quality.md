---
type: Design Doc
title: Report Local Keys and Navigation - design doc
description: Design for generated report local keys, compact navigation, and footer simplification.
tags: [evaluation, reports, markdown, accessibility]
timestamp: 2026-06-29T00:00:00Z
---

# Report Local Keys and Navigation - design doc

Design for
[Report Local Keys and Navigation](../0174-report-local-keys.md) and its
[functional spec](spec.md).

## Context

Generated reports are rendered by `internal/evaluation/report_tree.go`.
`report.md` uses a custom run-report renderer, while Area, Factor, Requirement,
Findings, Recommendations, and recommendation detail reports share
`renderReportHeader` for frontmatter, H1, run context, report navigation,
context lines, opening tables, and optional local jump links.

Display labels and markers for fixed Evaluation values are already centralized
through typed display catalogs. Rating Levels are different: their labels come
from the run's model snapshot. The local-key design should reuse those existing
sources so keys and table cells cannot drift.

## Approach

Keep the run report custom. Replace `writeRunReportContents` with a
`writeJumpLinksLine` helper that renders one `Jump to:` line when the report
uses local navigation. The run report will render H1, Summary, Key Details, the
jump line, Model Evaluation, ranked findings, ranked recommendations, and
Primary Source Data.

Expand the shared `reportHeader` path with explicit section labels for the
opening table. Area, Factor, Requirement, Findings, Recommendations, and
recommendation detail reports can render their opening tables under
`## Key Details` when they are part of the report opening. Existing run context,
report navigation, Area, Factor, and Factors lines remain near the H1.

Move summary prose in Area, Factor, and Requirement reports from the bare
`Summary:` label to `## Summary`. This is a direct writer change in the three
renderers and does not require new data.

Remove `writeEvaluationLegend` calls. Keep `writePrimarySourceDataSection` as
the only standard footer writer.

Add local key writers near the tables that first introduce each indicator
family:

- Key Details: ratings and confidence, plus status or assessment when present.
- Model Evaluation / Area-Factor Breakdown: Rating Levels, row markers, and
  empty marker.
- Findings tables: finding type and severity.
- Recommendation tables and detail headers: impact, confidence, and rank/tier
  where used.
- Requirement finding detail ranking tables: advice tier and empty marker where
  used.

Each key writer returns one line per family, for example:

```markdown
Ratings: 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.
Rows: `▦` Area, `□` Factor.
Empty: `—`.
```

The writer should use the same marker+label values as table cells. For Rating
Levels, derive the ordered set from the model snapshot. For fixed catalogs,
derive the ordered set from existing display catalog metadata rather than
duplicating literals inside report rendering. Unknown values remain rendered in
table cells through existing fallback label handling but are not invented into
catalog keys.

Tests should assert both the positive and negative contracts: no `## Contents`,
no `## Legend`, local key lines are present in representative reports, and table
cells still include marker plus text label values such as `▲ High`, `✅
Strength`, and `🔵 Target`.

## Spec response

- Removing `writeRunReportContents` and replacing it with an inline jump-link
  writer satisfies compact navigation requirements.
- Removing `writeEvaluationLegend` calls and keeping Primary Source Data last
  satisfies footer simplification.
- Per-report fixed key placement satisfies "first use" deterministically without
  adding dynamic document scanning.
- Reusing the model snapshot and display catalogs keeps local keys aligned with
  table labels.
- Header rendering changes provide `## Key Details` without changing report
  frontmatter.
- Summary writer changes align detail reports with the run report section role.

## Alternatives

- **A bottom legend table.** Rejected because it explains notation after first
  use and creates another footer utility competing with Primary Source Data.
- **No keys.** Rejected because readers benefit from seeing the full marker set
  when a report introduces compact visual indicators.
- **Definition-style explainers.** Rejected because generated reports should not
  become glossaries; table headings and durable docs own term meaning.
- **Dynamic first-use detection.** Rejected because report renderers already
  know their table sequence, and fixed placement is simpler and more stable.
- **Marker-only table cells plus keys.** Rejected because keys do not solve
  accessibility; the table cell must still carry the text label.

## Trade-offs & risks

Local keys add repeated lines to generated reports. The cost is acceptable
because the lines are short, deterministic, and close to the tables they
explain. Keeping keys notation-only reduces report clutter but means term
semantics still live in table headings, model Rating Scale criteria, structured
payloads, and durable docs.

Fixed key placement can over-render a family on a page where a table happens to
contain only one value from that family. That is intentional: keys show the
known display set for orientation and reduce data-dependent churn.

## Open questions

None.
