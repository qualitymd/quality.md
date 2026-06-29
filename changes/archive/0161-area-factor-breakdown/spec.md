---
type: Functional Specification
title: Area / Factor Breakdown Reports - functional spec
description: Requirements for rendering compact Area / Factor breakdown tables in generated Evaluation reports.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Area / Factor Breakdown Reports - functional spec

Companion to
[Area / Factor Breakdown Reports](../0161-area-factor-breakdown.md). This spec
states the delta for generated Evaluation Markdown reports. The durable source
of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md) and
[`/quality reporting`](../../../specs/skills/quality-skill/reporting.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Generated run reports need an overview that communicates the evaluated model
shape and the quality signals attached to it. A flat `Subject Reports` table
does not do that because it presents generated report files as the organizing
concept and mixes model subjects with index and recommendation artifacts. Area
reports have a related problem in smaller form: local Factors and Child Areas
are split into separate tables, making the local model shape less scannable than
it needs to be.

An Area / Factor breakdown gives readers one structural scan surface: the local
root Area, descendant Areas, Factors, ratings, ranked finding counts, and ranked
recommendation counts. The machine-oriented report manifest remains available
through `EvaluationOutputResult.reportOutputs`.

## Scope

This change covers generated Markdown bodies for:

- `report.md`;
- root and non-root Area reports.

Deferred:

- Factor report `Sub-Factors` rendering;
- Requirement rows in the Area / Factor breakdown;
- confidence, status, finding-type mix, recommendation-impact mix, and
  Requirement count columns;
- generated report paths and `EvaluationOutputResult` shape.

## Requirements

1. The run-level `report.md` **MUST** render an `Area / Factor Breakdown`
   section that lists the scoped Area as the first row, followed by in-scope
   descendant Areas and Factors in deterministic model order.

   > Rationale: The run report needs a structural quality overview, not a flat
   > generated-file manifest.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md`.

2. Area reports **MUST** render an `Area / Factor Breakdown` section that lists
   the reported Area as the first row, followed by its evaluated descendant
   Areas and Factors in deterministic model order.

   > Rationale: Area reports should use the same structure-summary surface as
   > the run report instead of splitting local Factors and Child Areas into
   > separate tables.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md`.

3. The Area / Factor breakdown table **MUST** use the columns `Area / Factor`,
   `Overall Rating`, `Local Rating`, `Findings`, and `Recommendations`, in that
   order.

   > Rationale: The table should stay narrow enough for generated Markdown while
   > showing the current roll-up rating and action/evidence density.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

4. The Area / Factor breakdown table **MUST** hyperlink the `Area / Factor`
   value to the generated Area or Factor report when that report exists and
   **MUST NOT** render a separate `Report` column.

   > Rationale: The row subject is the natural navigation target. A separate
   > filename column widens the table without adding reader value.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

5. The Area / Factor breakdown table **MUST** emphasize only the table's root
   Area row in the `Area / Factor` cell and **MUST** mark Factor rows inline
   with the existing Factor report-kind marker.

   > Rationale: The bold root row anchors the table to the current report
   > subject, while an inline Factor marker avoids a separate Kind column.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

6. The Area / Factor breakdown table **MUST** count ranked findings and ranked
   recommendations that resolve to each row's Area or Factor.

   > Rationale: Counts should reflect the Advice outputs readers can inspect in
   > the generated findings and recommendations reports rather than raw or
   > unranked intermediate data.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

7. Generated run reports **MUST NOT** render the old `Subject Reports` section
   in the Markdown body, and Area reports **MUST NOT** render the old separate
   `Factors` or `Child Areas` sections in the Markdown body.

   > Rationale: Keeping the old tables alongside the new breakdown would make
   > reports wider and more repetitive while preserving the file-manifest
   > framing this change removes.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `docs/guides/reporting-design.md`.

8. Report generation **MUST** preserve `EvaluationOutputResult.reportOutputs`
   and report frontmatter source-data behavior.

   > Rationale: This change affects human Markdown body presentation only. The
   > machine-readable generated-output manifest and source-data pointers remain
   > the report artifact index.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

## Acceptance criteria

- Generated `report.md` contains `## Area / Factor Breakdown` and does not
  contain `## Subject Reports`.
- Generated Area reports contain `## Area / Factor Breakdown` and do not contain
  `## Factors` or `## Child Areas`.
- Breakdown tables use the exact required columns and no separate report-link
  column.
- Breakdown table root Area rows are emphasized, and Factor rows carry the
  inline Factor marker.
- Breakdown table finding and recommendation counts are deterministic and based
  on ranked findings and recommendations.
- `data/evaluation-output-result.json` still contains `reportOutputs`.
- Focused Go tests cover run and Area report rendering.
- Report gallery generation is deterministic and checked in.
- `mise run check` passes.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - replace the old run-level subject
  reports and split Area structure-table contracts with the Area / Factor
  breakdown contract (requirements 1-8).
- `specs/skills/quality-skill/reporting.md` - align generated report
  expectations with the Area / Factor breakdown (requirements 1, 2, 7).

### To rename

None

### To delete

None
