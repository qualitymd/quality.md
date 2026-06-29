---
type: Design Doc
title: Report Frontmatter H1 Titles - design doc
description: Design for deriving generated report frontmatter titles from H1 titles.
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report Frontmatter H1 Titles - design doc

Design for
[Report Frontmatter H1 Titles](../0167-report-frontmatter-h1-titles.md) and its
[functional spec](spec.md).

## Context

Generated Evaluation reports already use one shared header renderer that writes
frontmatter and then writes the H1. Individual report renderers pass both a
subject-oriented `Title` and a document-oriented `Heading`; the mismatch comes
from the shared header writing frontmatter from `Title` while the body writes
the H1 from `Heading`.

## Approach

Remove the separate subject-title input from the shared report header and write
frontmatter `title` from the same `Heading` value used for the visible H1.
Report-specific renderers keep constructing their current H1 headings:

- `Evaluation Report: Area: <title>`;
- `Area: <title>`;
- `Factor: <title>`;
- `Requirement: <title>`;
- `Recommendation: <title>`;
- `Findings`;
- `Recommendations`.

This keeps the title contract centralized and avoids per-report duplication. The
renderer still writes only `type` and `title` frontmatter fields, and report
paths and recommendation title slugs remain based on their existing inputs.

## Spec response

- H1-matching frontmatter titles are enforced by the shared report header.
- Kind-prefixed reports need no special cases because their existing `Heading`
  strings already include the prefix.
- Index reports remain unchanged because their `Heading` values already match
  their prior frontmatter titles.
- Generated reports remain output-only; the renderer still builds from
  structured Evaluation data and never reads generated Markdown.

## Alternatives

- **Keep subject-only `title` and add a second document-title field.** Rejected
  because it expands frontmatter and preserves the ambiguity this case removes.
- **Strip kind prefixes out of H1 titles.** Rejected because previous report UX
  work made those prefixes the human scanning contract.
- **Derive titles by parsing generated Markdown after render.** Rejected because
  generated Markdown is not a report-generation input.

## Trade-offs & risks

Frontmatter titles become longer for run and detail reports. That is acceptable
because they are now unambiguous document titles and remain small scalar values.

Consumers that used frontmatter `title` as a subject-only display label need to
read structured Evaluation payloads or visible body context instead. That is the
intended boundary: frontmatter identifies the report document, while structured
payloads carry source identity.

## Open questions

None.
