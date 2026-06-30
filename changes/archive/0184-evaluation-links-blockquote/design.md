---
type: Design Doc
title: Evaluation Links Blockquote
description: Design rationale for rendering generated report Evaluation links as an H1-adjacent blockquote and preserving a compact Area/Factor marker key.
tags: [evaluation, reports, navigation]
timestamp: 2026-06-30T00:00:00Z
---

# Evaluation Links Blockquote

## Approach

The report renderer keeps the existing `writeEvaluationLinks` helper and changes
only its Markdown wrapper from a plain paragraph to a blockquote:

```markdown
> **Evaluation links:** [report.md](report.md) | [findings.md](findings.md) | [recommendations.md](recommendations.md) | [glossary.md](../../../glossary.md)
```

Shared non-run report header rendering calls the helper immediately after the
H1. The run report has its own header path, so it calls the same helper directly
after writing the run H1 and before `## Summary`.

## Rationale

A blockquote is the smallest Markdown-native treatment that gives the navigation
cluster a visual boundary in GitHub, editors, and static renderers. It avoids the
extra table syntax and escaped pipes a one-cell table would require, and it
avoids horizontal rules that would make the report opening feel broken into
separate panels before the reader reaches the summary.

The helper continues to own link construction, so relative target behavior and
the glossary path calculation stay unchanged.

The breakdown table header changes from `Area / Factor` to
`▦ Area / □ Factor` in the existing table writer. That keeps the row-marker key
where readers encounter the glyphs and avoids adding back a local legend block.
