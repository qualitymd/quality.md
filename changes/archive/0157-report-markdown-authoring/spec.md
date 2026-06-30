---
type: Functional Specification
title: Report Markdown Authoring
description: Requirements for centralizing generated Evaluation report Markdown authoring.
tags: [evaluation, reports, markdown, refactor]
timestamp: 2026-06-27T00:00:00Z
---

# Report Markdown Authoring

Companion to the [Report Markdown Authoring](../0157-report-markdown-authoring.md)
Change Case. This spec states _what_ the change must do. It defers to
[Evaluation report tree](../../specs/evaluation/reports/report-tree.md) and
[qualitymd evaluation report](../../specs/cli/evaluation-report.md) as normative
context for generated Evaluation report behavior.

The key words "MUST", "MUST NOT", and "MAY" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Generated Evaluation reports are deterministic Markdown projections over
structured data. The current direct string-authoring approach is easy to follow
for small report sections, but the report tree now repeats Markdown table syntax,
relative link calculation, empty-cell markers, and escaping across many
functions. That repetition increases the chance that a future report edit skips
escaping, creates malformed Markdown, or diverges from the report tree contract.

Research into current Go alternatives did not justify a wholesale dependency
switch. Parser/renderers such as goldmark and gomarkdown are useful for parsing,
validation, or conversion, but they do not remove the domain-specific work of
choosing report sections, rows, links, and deterministic paths. Markdown builder
and table packages reduce some syntax typing, but still need project-specific
rules for relative links, table-cell escaping, empty cells, and byte-stable
output. The change should therefore centralize this project's report primitives
instead of replacing the renderer with a general-purpose library.

## Scope

Covered: generated Evaluation report Markdown authoring, report writer helpers,
escaping and empty-cell behavior, relative links, code spans, table rendering,
and focused tests that prove preserved report behavior.

Deferred / non-goals: no report content redesign, no structured Evaluation data
contract changes, no new report format, no terminal styling in generated report
files, no template-first report authoring model, and no replacement for the
existing Glamour-based CLI display path for `qualitymd spec`.

## Requirements

### Authoring model

R1. Generated Evaluation report authoring **MUST** remain a deterministic
projection from completed structured Evaluation data and the model snapshot.

> Rationale: report authoring mechanics must not weaken the existing projection
> boundary or let a formatter dependency become a second source of report
> judgment.
>
> Durable spec: none.

R2. The implementation **MUST** provide a focused internal Markdown writer for
the primitives used by generated Evaluation reports: headings, paragraphs, raw
Markdown blocks, lists, tables, relative links, code spans, empty cells, and
source data links.

> Rationale: these are the repeated report-authoring operations that currently
> spread Markdown syntax and escaping decisions across the renderer.
>
> Durable spec: none.

R3. Generated Evaluation report rendering **MUST** use the internal writer for
new or migrated headings, tables, report links, data links, and empty-cell
markers instead of hand-assembling those Markdown constructs inline.

> Durable spec: none.

R4. The internal writer **MUST NOT** own Evaluation semantics such as which
reports exist, which rows are included, how findings or recommendations are
ranked, or which structured data fields are projected.

> Rationale: the writer should make Markdown authoring safer without hiding the
> report contract or moving domain behavior into formatting utilities.
>
> Durable spec: none.

### Markdown hygiene

R5. Report table rendering **MUST** escape Markdown table cell separators and
normalize multiline cell content so structured text cannot break table shape.

> Rationale: the current helper escapes `|`, but newline handling is not
> centralized; both characters can corrupt generated tables when they appear in
> persisted Evaluation text.
>
> Durable spec: modify `specs/evaluation/reports/report-tree.md` if the durable
> report contract does not already state the required table-cell behavior.

R6. Report link rendering **MUST** centralize relative-path calculation and
label escaping for generated report links and source data links.

> Durable spec: none.

R7. Empty optional report values **MUST** continue to render as `—` wherever
the existing report contract treats the value as not applicable or not recorded.

> Durable spec: modify `specs/evaluation/reports/report-tree.md` only if the
> implementation discovers an undocumented empty-cell case that should become
> durable behavior.

R8. The refactor **MUST** preserve existing generated report bytes unless a
change is required to fix malformed Markdown, centralize an existing documented
empty-state rule, or satisfy an explicit requirement in this spec.

> Rationale: this is a maintainability change; unrelated report-content churn
> would make review harder and risk changing the human Evaluation contract.
>
> Durable spec: none.

### Dependency boundary

R9. Core generated report authoring **MUST NOT** depend on a general-purpose
Markdown parser, renderer, table package, or template engine unless the design
doc demonstrates a concrete benefit that the focused internal writer cannot
provide with comparable simplicity.

> Rationale: research found the main pain is project-specific report primitive
> consistency, not CommonMark parsing or generic Markdown conversion.
>
> Durable spec: none.

R10. Parser/renderers such as goldmark or Glamour **MAY** be used for optional
validation or display tests, but they **MUST NOT** become required to construct
the generated Markdown report files.

> Durable spec: none.

### Verification

R11. Tests **MUST** cover table cells containing `|`, multiline text, empty
strings, links, and code spans through the internal writer.

> Durable spec: none.

R12. Evaluation report tests **MUST** prove representative generated reports
still include the existing required sections, links, and table shapes after the
renderer migration.

> Durable spec: none.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - clarify generated Markdown table
  hygiene and empty-cell behavior only if implementation changes durable report
  behavior or reveals an undocumented rule that should be carried forward (R5,
  R7).

### To rename

None

### To delete

None

## Open questions

- Should Markdown parse validation be added as a test-only guard after the
  internal writer is introduced, or are golden report and focused writer tests
  enough?
- Should the writer live as a general `internal/markdown` package, or as an
  evaluation-local helper until another generated Markdown surface needs it?
