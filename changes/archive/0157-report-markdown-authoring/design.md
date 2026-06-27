---
type: Design Doc
title: Report Markdown Authoring
description: Design for centralizing generated Evaluation report Markdown authoring.
tags: [evaluation, reports, markdown, refactor]
timestamp: 2026-06-27T00:00:00Z
---

# Report Markdown Authoring

Design behind the [Report Markdown Authoring](../0157-report-markdown-authoring.md)
Change Case and its [functional spec](spec.md).

## Context

Evaluation report generation already has the right data boundary: collect
completed structured payloads, resolve the report plan, then render deterministic
Markdown files. The weak point is authoring mechanics. Report sections now
repeat pipe tables, relative links, data links, empty-cell markers, and partial
escaping inline across `internal/evaluation/report_tree.go`.

The design keeps the data boundary and moves only Markdown mechanics into a
small internal writer.

## Approach

Add an `internal/markdown` package that owns Markdown primitives, not Evaluation
semantics. Its core type is a `Writer` backed by `strings.Builder` with methods
for the shapes generated reports use:

- headings and raw blocks;
- pipe tables with centralized cell escaping;
- unordered lists;
- relative links and data links;
- code spans;
- the shared empty marker.

The package exposes pure functions for reusable escaping and link behavior, such
as `Cell`, `Link`, `RelLink`, `Code`, and `DataLink`. The Evaluation renderer
uses those helpers while still deciding which reports, sections, rows, and data
fields exist.

Table cells normalize `|`, carriage returns, and newlines so persisted
Evaluation text cannot break table shape. Empty strings render as `—`, matching
the existing report legend. Links continue to use repository-relative slash
paths derived from the current report path.

Migration is deliberately incremental but complete for the risky constructs:
report tables, report links, data links, code spans, and the empty marker move to
the writer/helper functions. Larger prose sections may still call `WriteString`
or use raw blocks when they are already Markdown text from persisted structured
data.

## Spec Response

The writer satisfies the authoring-model requirements by making Markdown
primitives explicit while leaving Evaluation selection, ranking, and projection
logic in `report_tree.go`.

The hygiene requirements are handled centrally: table cells escape pipe
separators and normalize multiline text; link helpers share one relative-path
calculation; and empty values go through one marker.

The dependency requirements are satisfied by using only the standard library.
goldmark and Glamour remain available elsewhere for parsing/display, but they do
not construct generated report files.

Verification uses focused unit tests for the writer plus existing Evaluation
report tests and the report gallery check. The report gallery is important
because it exercises a full generated tree outside synthetic unit expectations.

## Alternatives

**Use goldmark or gomarkdown to build an AST.** Rejected because these packages
are strongest for parsing/rendering existing Markdown. They do not remove the
project-specific decisions about report rows, links, data paths, or empty-state
semantics.

**Use a table-rendering dependency.** Rejected because tables are only one
authoring hazard, and report cells intentionally contain Markdown links and code
spans that need project-owned escaping rules.

**Use `text/template` for report files.** Rejected because the reports are
branchy projections over variable report trees. Templates would hide control
flow and still require helper functions for escaping, relative links, and empty
cells.

**Keep helpers inside `internal/evaluation`.** Rejected because the primitives
are Markdown concepts rather than Evaluation concepts, and a small package gives
them focused tests without pulling in Evaluation fixtures.

## Trade-offs & Risks

The package name `markdown` is broad, so it must stay small and primitive. It is
not a place for report policy, domain vocabulary, or renderer orchestration.

Normalizing multiline table content is a deliberate hygiene fix that can change
bytes for malformed existing cells. That is acceptable because preserving broken
tables is not a useful compatibility target, but tests should prove ordinary
report output is otherwise stable.

The migration may leave some raw Markdown writes in the Evaluation renderer.
That is acceptable for persisted Markdown prose or fixed section structure; the
writer boundary matters most where syntax and data intersect.
