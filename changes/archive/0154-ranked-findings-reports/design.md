---
type: Design Doc
title: Ranked Findings Reports Design
description: Implementation approach for ranked findings tables, findings.md, and per-finding Advice ranking context.
tags: [evaluation, reports, advice, findings]
timestamp: 2026-06-27T00:00:00Z
---

# Ranked Findings Reports Design

## Context

This design answers the [Ranked Findings Reports functional spec](spec.md). The
existing report renderer already loads `FindingRankingResult`, dereferences
ranked findings, and renders a compact Top Findings table in `report.md`. The
change reshapes that projection and adds a full findings index while preserving
the existing judgment/mechanics boundary: Advice ranking remains agent-authored
data, and the CLI remains a deterministic renderer.

## Approach

### Shared ranked finding row model

Add an internal helper that turns each
`FindingRankingResult.orderedFindings[]` entry into a render-ready row:

- rank and total ranked count;
- the referenced Requirement ID;
- the referenced finding ID and finding payload;
- linked finding statement targeting the Requirement report's stable
  `finding-<id>` anchor;
- linked declaring Area title;
- linked attached Factor titles;
- finding type and severity display labels; and
- tier and ranking rationale for Requirement detail rendering.

The helper should sort rows by the persisted `rank` field and leave judgment
fields unchanged. It should tolerate missing dereference data by falling back to
the selector or `—` display cells, but reportability validation should continue
to catch invalid references before normal report build.

### Report tree artifacts

Reuse one table writer for both the run-report Top Findings section and the full
index:

```text
writeRankedFindingsTable(limit)
```

`report.md` calls the shared writer with `limit = 10`, then always writes:

```markdown
Full findings index: [findings.md](findings.md)
```

The report builder adds a first-class `findings.md` artifact at the run root.
`findings.md` renders a short heading, a source data link to
`data/advice/finding-ranking-result.json`, and the same table with no limit.
Adding a `ReportKindFindings` value keeps `data/evaluation-output-result.json`
and subject report tables able to name the artifact consistently if the renderer
chooses to include it in report references.

### Stable finding anchors

Requirement finding details currently render a heading that combines the finding
ID and statement. Markdown-generated anchors therefore depend on renderer and
statement text. Emit an explicit HTML anchor before each finding detail:

```html
<a id="finding-gap-001"></a>
```

The detail heading can remain human-readable. Ranked-finding links target
`requirement-report.md#finding-gap-001`.

### Requirement detail ranking context

Build a lookup from artifact finding reference key to ranked row. Pass the
artifacts object into `writeEvaluationFindingDetails` so each detail section can
render a compact table before Condition/Criteria/Basis/Effect/Evidence:

```markdown
| Advice Rank | Tier | Ranking Rationale |
| ----------- | ---- | ----------------- |
| 12 / 55     | P2   | ...               |
```

When no ranking row exists, render:

```markdown
| Advice Rank  | Tier | Ranking Rationale |
| ------------ | ---- | ----------------- |
| (not ranked) | —    | —                 |
```

This keeps the Advice context visibly separate from the finding's evidence
fields.

## Spec Response

- **Ranked findings tables (R1-R5)** - one shared writer renders the capped
  `report.md` table and the full `findings.md` table, both ordered by persisted
  rank; `report.md` always links to the full index.
- **Links and display (R6-R10)** - row construction dereferences the Requirement,
  Area, Factors, and Finding once, using existing path and display helpers for
  links and enum labels.
- **Requirement finding details (R11-R14)** - stable anchors provide exact link
  targets, and a ranking-context table projects only persisted rank, tier, and
  rationale.

## Alternatives

**Keep Top Findings as the only findings index.** Rejected because a top-10
entrypoint cannot support full review when a run has more ranked findings.

**Use `Location` with Area / Factor / Requirement.** Rejected after discussion:
separate `Area` and `Factors` columns preserve model context without making the
row path-like, and the `Finding` link already names the Requirement detail.

**Show confidence in ranked findings tables.** Rejected for this slice. The top
tables are navigational decision surfaces; confidence remains available in the
Requirement finding summary/details.

**Use Markdown-generated heading anchors.** Rejected because statement edits
would break inbound links from `report.md` and `findings.md`.

## Trade-offs & Risks

- The new full index adds another generated root file, but it mirrors the
  existing recommendation index pattern and keeps `report.md` compact.
- Explicit HTML anchors are less visually pure Markdown, but they give stable
  link targets across renderers.
- Requirement detail pages will repeat ranking context for every finding. The
  table is intentionally small so it does not crowd out the evidence sections.

## Open Questions

None.
