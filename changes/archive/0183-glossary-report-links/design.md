---
type: Design Doc
title: Glossary and Report Links
description: Design for adding a shared glossary and replacing generated report legends with Evaluation links.
tags: [evaluation, reports, glossary]
timestamp: 2026-06-30T00:00:00Z
---

# Glossary and Report Links

## Context

This design answers the [functional spec](spec.md). The change creates a
workspace-root `glossary.md`, gives that artifact its own durable spec, and
changes generated Evaluation reports so they link to the overview, Findings,
Recommendations, and glossary instead of rendering repeated local `Legend`
blocks.

## Approach

Keep the glossary as a checked-in authored Markdown artifact for this slice.
The fixed Evaluation enum source of truth remains the typed catalog in
`internal/evaluation/display.go`; the first glossary version copies those
labels, values, descriptions, and order into `glossary.md`. Generation is
deferred until the artifact proves stable.

Give `glossary.md` a 1:1 durable artifact spec at `specs/glossary-md.md`. That
spec owns the glossary's purpose, location, entry ordering, table shape, initial
terms, configured Quality rating note, and fixed Evaluation vocabulary entries.
The report-tree spec owns how generated reports link to the glossary.

In report rendering, replace every `writeLocalKeys` call with one shared
`Evaluation links:` line near the top of each report. The line uses filename link
text and relative links from the current report artifact:

```markdown
**Evaluation links:** [report.md](...) | [findings.md](...) | [recommendations.md](...) | [glossary.md](...)
```

Reports keep marker-plus-text labels in cells. Only the explanatory local
`Legend` blocks go away.

## Spec response

- Glossary structure and content are satisfied by the checked-in `glossary.md`
  plus `specs/glossary-md.md`.
- Fixed vocabulary tables preserve `Label`, `Value`, `Description` columns and
  keep catalog order inside each table.
- Quality rating uses this repository's configured `QUALITY.md` Rating Scale and
  carries the source note because Rating Levels are model-defined.
- Generated reports gain one `Evaluation links:` line after the opening
  summary/key details/orientation and before `## Contents`.
- Local `Legend` blocks are removed without removing visible labels from report
  cells.
- The stale `unknown` Finding type entry is removed from the durable payload
  spec.

## Alternatives

**Generate `glossary.md` from Go catalogs now.** Rejected for this slice. It
would avoid manual enum drift but adds a generator, fixtures, and release
workflow surface before the glossary contract is proven.

**Keep glossary rules inside `report-tree.md`.** Rejected. `glossary.md` is a
workspace-root durable artifact used by generated reports, so it should have a
spec that owns its own artifact contract.

**Link every report table value to the glossary.** Rejected. It would make
generated Markdown noisy and harder to scan; one stable navigation line is
enough.

**Leave local legends and add glossary links.** Rejected. That preserves the
duplication the glossary is meant to remove.

## Trade-offs & risks

The first glossary is manually maintained, so enum catalog changes can drift
until generation or a dedicated check is added. Tests should cover the report
side of the contract, while review must check the glossary content against the
typed catalog for this slice.

The `Evaluation links:` line appears in every generated report, including the
current artifact. Keeping all links active is mechanically simpler and keeps the
navigation stable, at the cost of one self-link per artifact.

## Open questions

None for this slice.
