---
type: Change Case
title: Report Markdown Authoring
description: Centralize generated Evaluation report Markdown authoring without replacing deterministic Go rendering.
status: Done
tags: [evaluation, reports, markdown, refactor]
timestamp: 2026-06-27T00:00:00Z
---

# Report Markdown Authoring

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0157-report-markdown-authoring/spec.md) - what the case must
  do.
- [Design doc](0157-report-markdown-authoring/design.md) - how it is built, and
  why.

## Motivation

Generated Evaluation reports are currently authored directly in Go with
`strings.Builder`, `fmt.Fprintf`, and local helpers. That keeps report generation
deterministic and close to the structured Evaluation data, but report syntax is
now scattered across many headings, tables, links, and list blocks. The most
fragile pieces are Markdown table rows, relative links, empty cells, and escaping:
small report changes must remember which helper to call and which values are
already safe.

Online research did not identify a broad Go Markdown authoring library that
would bring enough benefit to justify handing report shape to a dependency.
Parser/renderers such as goldmark are better suited to validation or conversion,
while table or builder packages still need project-specific link, empty-cell,
and deterministic-output rules. The useful change is a small internal Markdown
writer that expresses the report primitives this project already uses.

## Scope

Covered:

- Add a focused internal Markdown/report writer for generated Evaluation reports.
- Centralize heading, paragraph, list, table, relative-link, code-span, empty
  cell, and table-cell escaping behavior.
- Migrate generated Evaluation report rendering to the writer while preserving
  byte-stable report output except for deliberate Markdown hygiene fixes.
- Add focused tests for report escaping and generated table/link behavior.
- Keep parser/renderers such as goldmark or Glamour out of the core report
  authoring path; they may be used only for optional validation or display.

Deferred:

- Changing Evaluation report content, report paths, report sections, ranking
  semantics, or structured Evaluation data contracts.
- Rendering HTML, PDF, terminal-styled reports, or non-Markdown report formats.
- Introducing templates as the primary report authoring model.
- Replacing the existing CLI Markdown display path for `qualitymd spec`.

## Affected artifacts

Derived by sweeping for `renderEvaluation*`, `markdownCell`, `reportLink`, table
rows in generated reports, `evaluation report build`, and durable report tree
contracts.

**Code**

- [ ] `internal/evaluation/report_tree.go` - migrate string-built report blocks
- [x] `internal/evaluation/report_tree.go` - migrate string-built report blocks
      to the internal writer and keep data-to-report behavior deterministic.
- [x] `internal/markdown/` - centralize report
      primitives, escaping, relative links, and empty-cell behavior.
- [x] `internal/evaluation/evaluation_test.go` - cover preserved report output
      and edge cases for table/link escaping.
- [x] Optional focused tests for the new writer - verify headings, paragraphs,
      lists, tables, links, code spans, empty cells, and multiline cells.

**Format spec and durable specs** (substance in the [functional spec](0157-report-markdown-authoring/spec.md))

- [x] `specs/evaluation/reports/report-tree.md` - clarify generated Markdown
      hygiene rules when durable behavior is affected by escaping or empty-cell
      handling.
- [x] `specs/cli/evaluation-report.md` - no content change expected unless
      report determinism or renderability wording needs tightening.
- [x] `specs/log.md` and `specs/evaluation/log.md` - durable spec log entries if
      durable report specs change.
- [x] `SPECIFICATION.md` - no change expected; report authoring internals do not
      alter the format contract.

**Durable docs / bundled skill runtime**

- [x] `CHANGELOG.md` - release-note entry if the refactor lands user-visible
      Markdown hygiene fixes.
- [x] `changes/index.md`, `changes/archive/index.md`, and `changes/log.md` -
      Change Case lifecycle.
- [x] `/quality` runtime skill files - no change expected; report content and
      closeout guidance stay the same.

## Status

`Done`. Implemented, verified, and archived. Generated Evaluation reports now
use shared Markdown helpers for table rows, links, data links, code spans, and
empty cells; durable report specs document table-cell hygiene; focused helper
tests and the full report gallery check pass. `mise run check` passes.
