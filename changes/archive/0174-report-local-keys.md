---
type: Change Case
title: Report Local Keys and Navigation
description: Replace generated report contents and legend sections with compact navigation and local indicator keys.
status: Done
tags: [evaluation, reports, markdown, accessibility]
timestamp: 2026-06-29T00:00:00Z
---

# Report Local Keys and Navigation

Generated Evaluation reports should explain their notation where readers first
need it, without adding a heavy contents section or a bottom legend. This case
replaces `## Contents` and `## Legend` with compact local navigation and
notation-only local keys while preserving text labels as the accessible meaning
of ratings, statuses, confidence, finding classifications, recommendation
impact, and priority-like values.

- [Functional spec](0174-report-local-keys/spec.md) - what generated reports,
  report specs, and guidance must change.
- [Design doc](0174-report-local-keys/design.md) - how the renderer will add
  local keys and remove contents/legend sections.

## Motivation

Recent run-report improvements made `report.md` more decision-ready, but the
report family still mixes heavy generated chrome with under-localized notation:
`## Contents` creates a whole section for links, `## Legend` explains markers at
the bottom after readers have already consumed the tables, and detail reports
still use older opening patterns. Local keys make indicator sets visible at the
point of first use without becoming prose explainers, while the report content
itself keeps text labels for accessibility and agent readability.

## Scope

Covered:

- removal of generated `## Contents` sections;
- compact `Jump to:` navigation where local navigation is useful;
- removal of generated bottom `## Legend` sections;
- notation-only local keys after first use of an indicator family;
- one key line per indicator family when a table uses multiple families;
- simple visible `# Findings` and `# Recommendations` H1s;
- `## Summary` and `## Key Details` alignment for detail reports;
- generated report renderer code and tests;
- durable report specs, CLI report specs, skill reporting specs, and report
  design guidance;
- generated report gallery output.

Deferred:

- structured Evaluation JSON schema changes;
- display catalog value changes;
- Rating Scale model semantics;
- generated report filenames and frontmatter fields;
- non-Markdown or interactive report surfaces.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - render compact navigation, local
    keys, `## Summary`, `## Key Details`, and remove bottom legends.
  - `internal/evaluation/evaluation_test.go` - update report rendering
    assertions and marker/text-label coverage.
  - `internal/evaluation/display.go` - use existing display catalogs as local
    key sources if needed.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - define local navigation,
    local keys, footer rules, and summary/key-details shape.
  - `specs/cli/evaluation-report.md` - mirror generated report output behavior
    for `qualitymd evaluation report build`.
  - `specs/skills/quality-skill/reporting.md` - mirror report generation and
    `/quality` report consumption guidance.
- Durable docs:
  - `docs/guides/reporting-design.md` - align design principles, examples, and
    checklist.
- Generated examples:
  - `examples/report-gallery/` - regenerated report output.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md` and `docs/log.md` for durable spec/doc updates.

## Status

`Done`. Implemented and archived. Generated Evaluation reports now use compact
local navigation and notation-only local keys instead of generated Contents and
bottom Legend sections, while table cells keep marker-plus-text display labels.
