---
type: Change Case
title: Report Contents Sections
description: Restore generated report Contents sections and remove compact Jump to lines.
status: Done
tags: [evaluation, reports, markdown, navigation]
timestamp: 2026-06-29T00:00:00Z
---

# Report Contents Sections

Generated Evaluation reports should use a consistent `## Contents` section when
they contain multiple substantive sections. The previous compact `Jump to:` line
kept navigation light, but it made report navigation a special inline idiom
instead of a normal Markdown section and made the report-design guidance push
against the value of stable local document structure.

- [Functional spec](0175-report-contents-sections/spec.md) - what generated
  reports, durable specs, and guidance must change.
- [Design doc](0175-report-contents-sections/design.md) - how the renderer will
  add Contents sections and remove compact Jump to lines.

## Motivation

Report artifacts are reader-facing Markdown documents, not just terminal output.
When a generated report has multiple top-level sections, a standard Contents
section gives people and agents a predictable place to scan the document shape.
Index files are the counterexample: their whole job is already listing and
navigation, so adding a Contents section there duplicates the artifact's
purpose.

## Scope

Covered:

- generated Evaluation Markdown report artifacts;
- deterministic `## Contents` sections for multi-section generated reports;
- removal of generated `Jump to:` lines;
- the exception for OKF `index.md` and other listing/index artifacts;
- report renderer code and tests;
- durable report specs, CLI report specs, skill reporting specs, and report
  design guidance;
- generated report gallery output.

Deferred:

- generated report filenames, frontmatter fields, and report kind taxonomy;
- structured Evaluation JSON schema changes;
- local notation keys and marker/label display catalogs;
- non-Markdown or interactive report surfaces.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - render Contents sections instead of
    compact `Jump to:` lines.
  - `internal/evaluation/evaluation_test.go` - update report navigation
    assertions and negative coverage for `Jump to:`.
- Durable specs:
  - `specs/evaluation/reports/report-tree.md` - define Contents placement,
    eligibility, index exceptions, and removal of `Jump to:`.
  - `specs/cli/evaluation-report.md` - mirror generated report navigation
    behavior for `qualitymd evaluation report build`.
  - `specs/skills/quality-skill/reporting.md` - mirror report navigation and
    report consumption expectations.
- Durable docs:
  - `docs/guides/reporting-design.md` - update report-design guidance,
    examples, and checklist.
  - `CHANGELOG.md` - record the user-visible generated report output change.
- Bundled skill:
  - `skills/quality/SKILL.md` - align runtime skill guidance if report-reading
    instructions mention generated Contents or Jump to behavior.
- Generated examples:
  - `examples/report-gallery/` - regenerated report output.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/log.md`, `docs/log.md`, and `skills/quality/log.md` as needed for
    durable spec, doc, and skill updates.

## Status

`Done`. Implemented and archived. Generated Evaluation reports now use standard
Contents sections for multi-section report artifacts and no longer render
compact `Jump to:` local navigation lines.
