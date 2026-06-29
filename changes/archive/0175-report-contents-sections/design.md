---
type: Design Doc
title: Report Contents Sections - design doc
description: Design for generated report Contents sections and removal of compact Jump to lines.
tags: [evaluation, reports, markdown, navigation]
timestamp: 2026-06-29T00:00:00Z
---

# Report Contents Sections - design doc

Design for [Report Contents Sections](../0175-report-contents-sections.md) and
its [functional spec](spec.md).

## Context

Generated reports are assembled in `internal/evaluation/report_tree.go`.
`report.md` has a custom renderer, while findings, recommendations, Area,
Factor, Requirement, and recommendation detail reports share the `reportHeader`
path for frontmatter, H1, run context, report navigation, context lines, opening
tables, local notation keys, and the old optional `Jump to:` line.

The current renderer writes report bodies directly to a `strings.Builder`. It
already knows each report's section order at the call site, and section anchors
are ordinary GitHub-style heading anchors already used by the old jump links.

## Approach

Replace the `reportJumpLink`/`JumpLinks` abstraction with a
`reportContentLink`/`Contents` abstraction that renders:

```markdown
## Contents

- [Summary](#summary)
- [Key Details](#key-details)
- [Primary Source Data](#primary-source-data)
```

The shared `renderReportHeader` will render `## Contents` after the opening
navigation/context/local-key block when the caller provides at least two
content links. The run report will use the same helper after `## Key Details`
and its local keys, before `## Model Evaluation`, so the report opening remains
H1, Summary, Key Details, Contents, then substantive body sections.

Each report renderer will pass the visible `##` sections it is about to render,
in order, excluding `Contents` itself. The link list is static per renderer
because the rendered section sequence is static per report kind. Primary Source
Data stays in the list because it is a visible top-level section and a useful
agent target. The renderer does not scan generated Markdown after the fact; it
uses the same known section names that it writes.

Remove `writeJumpLinksLine`, `reportJumpLinks`, and `JumpLinks` fields. This
prevents the old idiom from surviving in any generated report path. Tests will
assert positive Contents output and negative `Jump to:` output across the run
report, detail reports, list reports, and recommendation detail reports.

Durable docs and specs will be updated to state the new rule, remove the old
"materially improves scanning" discretionary rule, and preserve the exception
for OKF `index.md` and other listing/index artifacts.

## Spec response

- Rendering Contents from explicit per-report section lists satisfies the
  deterministic multi-section navigation requirements.
- Gating on at least two supplied links satisfies the no-single-section-noise
  requirement without a subjective length heuristic.
- Removing the old jump-link helpers satisfies the `Jump to:` prohibition.
- Keeping links to visible `##` sections preserves shallow Contents and includes
  Primary Source Data.
- Leaving frontmatter, H1s, local keys, and source-data list writers unchanged
  confines the change to local navigation.

## Alternatives

- **Post-process headings from generated Markdown.** Rejected because report
  renderers already know their section sequence, and parsing the just-rendered
  Markdown would add a second source of truth.
- **Keep `Jump to:` for short reports.** Rejected because the goal is one
  consistent report navigation idiom.
- **Always render Contents even for one-section artifacts.** Rejected because it
  adds noise and violates the user's stated counterexample boundary.
- **Exclude Primary Source Data.** Rejected because it is a visible report
  section and a common reader/agent destination.

## Trade-offs & risks

Contents adds a full section to generated reports, including short detail
reports. The rule remains bounded by requiring at least two substantive
sections, and all current generated report artifacts have enough section
structure for the navigation to be useful.

Static per-renderer Contents lists can drift if a future renderer adds, removes,
or conditionally omits a top-level section without updating the list. Tests over
representative reports should check that the lists contain the visible section
set expected for each report family.

## Open questions

None.
