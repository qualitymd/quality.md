---
type: Design Doc
title: Report header kind prefix and title-first layout
description: How the report renderers deliver a kind-prefixed, title-first header — reordering the existing writers, dropping the identifier line, and the alternatives weighed.
status: Draft
tags: [evaluation, reports, markdown]
timestamp: 2026-06-26T00:00:00Z
---

# Report header kind prefix and title-first layout

## Context

Answers [the 0119 spec](spec.md): give every generated report a kind-prefixed
H1, render it first, drop the `Path:` / `Name:` line, and lock the `Area:`
trail's root element to the Model `title`. All three renderers live in
[`internal/evaluation/report_tree.go`](../../../internal/evaluation/report_tree.go)
(`renderEvaluationAreaReport`, `renderEvaluationFactorReport`,
`renderEvaluationRequirementReport`), and each already builds the header by
calling the shared trail writers and then writing an `#` + title line, a
`Path:` / `Name:` line, and a summary table.

## Approach

The header today is built bottom-irrelevant order: trail writer(s), then title,
then identifier line, then table. The change is almost entirely a **reordering**,
not new machinery:

1. Write the kind-prefixed H1 first. Inline the prefix at each call site —
   `"# Area: "`, `"# Factor: "`, `"# Requirement: "` + the resolved title — rather
   than threading a kind enum through a helper; there are exactly three sites and
   the literal reads clearest where it is.
2. Move the existing trail-writer calls (`writeEvaluationAreaTrail`, and for
   Factor `writeEvaluationFactorTrail`, for Requirement
   `writeEvaluationRequirementFactorsLine`) to run *after* the title line,
   unchanged.
3. Delete the `Path:` (Area, Factor) and `Name:` (Requirement) lines.

The root-element Model `title` already resolves through `areaTitle(spec, nil)`
(it returns `spec.Title`), so step 1's title for the root Area report — and the
root link the trail writer emits — both already carry it. No renderer change is
needed for the trail-root requirement; it is satisfied by the existing helper and
locked by a new test assertion.

Resulting Factor header (sketch):

```
# Factor: Authentication

Area: Root / Security

Factor: Authentication

| Overall Rating | Local Rating | Status | Confidence | Data |
```

### The one real decision: how the trail lines stack

The agreed "stacked" layout shows the `Area:` and `Factor:` lines on consecutive
lines. In Markdown, two adjacent non-blank lines collapse into one soft-wrapped
paragraph, so a literal tight stack would need hard line breaks (trailing spaces
or `\`). The existing trail writers each emit a trailing blank line (`\n\n`), so
today the lines already render as **separate paragraphs** with vertical space
between — proven, and what readers see now.

This design keeps that: reuse the writers verbatim and only reorder. The lines
stack (each on its own line) via paragraph separation, not via new hard-break
handling. The spec requires preserving the trails' "content, separators, and link
targets" — it does not mandate collapsing the inter-line spacing — so reusing the
writers satisfies it at the lowest risk.

## Spec response

- **Kind-prefixed, title-first H1** — step 1 writes `# <Kind>: <title>` as the
  first `WriteString`, before any trail writer, satisfying both the first-line
  and the prefix requirements.
- **Trails follow the title** — step 2 moves the unchanged writers below the
  title; their content, order, and separators are untouched.
- **No separate identifier line** — step 3 deletes the `Path:` / `Name:` writes.
- **Trail root element** — satisfied by the existing `areaTitle(spec, nil)`
  resolution; verified, not re-implemented.

## Alternatives

- **Kicker line under a bare H1** (`# Authentication` / `**Factor** · id`) — keeps
  the H1 equal to the entity title (matching parent-table link text) but adds a
  new header line and a new convention. Rejected in favor of the inline prefix
  when the layout was chosen against concrete sketches; the prefix lives only at
  the H1 render site, so link text (built from `factorTitle` / `requirementTitle`)
  stays bare regardless.
- **Keep trail-first, add the prefix only** — answers "what kind" without the
  reorder, but leaves navigation chrome as the opening line. Rejected: title-first
  was an explicit goal.
- **Tight-stack the trails with hard line breaks** — matches the idealized sketch
  byte-for-byte but introduces trailing-whitespace / `\` break handling for no
  reader-visible gain over paragraph separation. Rejected as needless risk.
- **Relabel the `Area:` trail** (e.g. `In:` / `Location:`) to kill the
  `# Area: X` / `Area: …` echo on Area reports — deferred by the spec; the echo is
  accepted, not fixed, here.
- **A `reportKindTitle(kind, title)` helper** — over-abstracts three one-line call
  sites. Left as an In-Progress micro-decision, not a design commitment.

## Trade-offs & risks

- **The `Area:` echo** on Area reports (`# Area: X` above `Area: … / X`) is a
  knowingly accepted cosmetic cost; the prefix is a kind label and the trail is
  navigation.
- **The root report repeats the Model title** in its H1 and its trail. Minor and
  inherent to title-first + kind prefix on the root Area report.
- **Test churn**: header assertions for all three report kinds change (new H1
  text, trail-after-title order, absent `Path:` / `Name:`, root-element title).
  Mechanical and low-risk; no structured-data or path assertions move.

## Open questions

None. The helper-vs-inline choice and any exact blank-line bytes are routine
In-Progress decisions within the approach above.
