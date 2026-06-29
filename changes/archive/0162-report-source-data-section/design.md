---
type: Design Doc
title: Report Source Data Section — design doc
description: Design for identity-only report frontmatter and bottom source-data sections.
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report Source Data Section — design doc

Design for [Report Source Data Section](../0162-report-source-data-section.md)
and its [functional spec](spec.md).

## Context

The Evaluation report renderer already computes report-local source payload
paths and passes them through the shared report header path as frontmatter
`data`. The change keeps that source-path computation, removes only its YAML
projection, and renders the same paths at the end of each report body.

## Approach

Keep the existing source-data path helpers and report-specific source path
assembly. Change the shared report header renderer so YAML frontmatter emits
only `type` and `title`.

Add one shared report-body helper that appends:

```markdown
## Source Data

- [data/run-manifest.json](data/run-manifest.json)
```

The helper receives the current report path and its run-root-relative
source-data paths. It renders labels as the original run-root-relative paths and
targets through the existing report-relative link helper so nested reports point
back to `data/` correctly.

Append the helper at the end of every generated report renderer after all
human-facing content, including the legend when present. Do not add source-data
links to header summary tables.

## Spec response

- Identity-only frontmatter is centralized in the shared header renderer.
- The source-data section is centralized in one body helper and called by every
  report renderer.
- Report-local source semantics are preserved by reusing the same source path
  inputs the current frontmatter `data` list uses.
- Nested links are handled by the existing relative report-link behavior instead
  of bespoke path math.
- Generated reports remain output-only; the renderer still builds from
  structured Evaluation data and never reads generated Markdown.

## Alternatives

- **Keep source data in frontmatter and add collapsible body links.** Rejected
  because it leaves the noisy first-screen YAML list in place.
- **Move source data into top report headers.** Rejected because it competes with
  rating, confidence, navigation, and other scan-critical report state.
- **Render human-friendly labels instead of paths.** Rejected because path labels
  are more useful to agents, avoid another naming surface, and make the link
  target obvious in plain text.
- **List all run data on every report.** Rejected because it would turn a
  report-local source manifest into a noisy run index.

## Trade-offs & risks

The source-data section is now visible body content, so long source lists add
weight at the bottom of reports. That is acceptable because the section is after
the decision content and remains easier to ignore than frontmatter.

Agents that only parse YAML frontmatter for source data will need to read the
stable `## Source Data` section instead. The path labels and heading make this a
simple parser adjustment.

## Open questions

None.
