---
type: Functional Specification
title: Run Report Simplification - functional spec
description: Requirements for simplified generated run-level Evaluation reports.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Run Report Simplification - functional spec

Companion to
[Run Report Simplification](../0171-run-report-simplification.md). This spec
states the delta for generated Evaluation Markdown reports. The durable source
of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md)
and [`/quality reporting`](../../../specs/skills/quality-skill/reporting.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The run-level `report.md` is the decision-ready Evaluation entrypoint. Recent
report iterations moved metadata into frontmatter and made the opening more
summary-first, but the run report still repeats navigation, subject, scope,
coverage, and report-detail information that either appears in frontmatter,
appears in `Key Details`, or belongs in granular reports. Its source-data list
also behaves like a transitive provenance manifest, duplicating payloads already
listed by Area, Factor, Requirement, finding, and recommendation reports.

This change keeps the run report human-first: the report should present the
overall judgment, the evaluated Model structure, ranked findings, and ranked
recommendations before provenance. Source-data sections remain useful, but they
name the report's primary structured inputs instead of every downstream payload
that influenced linked detail reports.

## Scope

This change covers generated Evaluation Markdown reports:

- run-level `report.md` opening lines and section order;
- run-level `report.md` local contents links;
- the run-level model-structure table heading;
- bottom source-data section heading and source-data selection principle;
- durable report specs, `/quality` reporting specs, and report design guidance.

Deferred:

- Evaluation structured data schema changes;
- report frontmatter fields;
- finding and recommendation ranking behavior;
- detail report navigation trails;
- global renaming of detail report `Area / Factor Breakdown` sections;
- non-Markdown or interactive report surfaces.

## Requirements

1. The run-level `report.md` **MUST NOT** render post-H1 `Report:` or `Area:`
   navigation lines.

   > Rationale: The run report's H1, frontmatter, Key Details, Contents, and
   > linked detail reports already provide orientation; these lines delay the
   > Summary without adding decision-relevant information. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

2. The run-level `report.md` Summary **MUST NOT** render a
   `Recommended next action:` sentence.

   > Rationale: Top Recommendations is the ranked recommendation surface; a
   > second selected-next-action sentence duplicates that priority signal and
   > makes Summary do recommendation-table work. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

3. The run-level `report.md` **MUST NOT** render `## Scope`, `## Coverage`, or
   `## Report Details` sections.

   > Rationale: Requested scope remains visible in Key Details and frontmatter,
   > while run identity and subject remain in frontmatter and primary source
   > data. Coverage detail belongs in the generated report tree and granular
   > reports instead of a separate run-report block. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

4. When the root Area was not evaluated in the run, the run-level `report.md`
   **MUST** preserve that coverage signal in a compact visible location other
   than a `## Coverage` section.

   > Rationale: Scoped runs must not read as whole-model verdicts, even though
   > the old Coverage section is removed. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

5. The run-level `report.md` **MUST** render its model-structure table under
   `## Model Evaluation`, before Top Findings and Top Recommendations.

   > Rationale: Readers need the evaluated Model structure and rating
   > distribution before the ranked lists; `Model Evaluation` is clearer and
   > less implementation-shaped than `Area / Factor Breakdown`. Durable spec:
   > modify `specs/evaluation/reports/report-tree.md`.

6. Area reports **MUST** continue to render their model-structure table under
   `## Area / Factor Breakdown`.

   > Rationale: The naming concern is strongest on the run report; changing
   > detail reports would broaden the scope beyond the primary entrypoint.
   > Durable spec: none; existing durable specs already govern Area reports.

7. The run-level `report.md` Contents **MUST** link only to visible run-report
   sections in their rendered order.

   > Rationale: Contents should be a trustworthy local map after the report
   > shape changes. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

8. Generated Markdown reports **MUST** use `## Primary Source Data` as the
   bottom source-data section heading.

   > Rationale: The heading makes explicit that the section lists the primary
   > structured inputs for that report, not an exhaustive transitive provenance
   > manifest. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md`.

9. A generated report's `Primary Source Data` section **MUST** list the
   report-local primary structured payloads used to establish that report's
   identity, scope, subject result, ranking, or recommendation content.

   > Rationale: Report-local primary sources keep traceability useful without
   > making high-level reports duplicate the source lists already owned by
   > granular reports. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md`.

10. The run-level `report.md` Primary Source Data section **MUST NOT** include
    descendant Area, Factor, Requirement, assessment, rating, finding, or
    recommendation payloads solely because the run report links to, summarizes,
    or counts data from granular reports.

    > Rationale: Detail reports are the traceability surface for their own
    > subjects; `report.md` should list primary run/report inputs and link into
    > granular reports for deeper traceability. Durable spec: modify
    > `specs/evaluation/reports/report-tree.md`.

11. Generated report frontmatter **MUST NOT** change as part of this case.

    > Rationale: Frontmatter already carries non-judgmental routing metadata and
    > remains the right place for run ID, created time, requested scope, and
    > subject reference. Durable spec: none; existing durable specs already
    > govern report frontmatter.

## Acceptance criteria

- Run-level `report.md` opens with H1, Summary, Key Details, Contents, then
  `Model Evaluation`.
- Run-level `report.md` does not render post-H1 `Report:` or `Area:` lines.
- Run-level `report.md` does not render `Recommended next action:`.
- Run-level `report.md` does not render `## Scope`, `## Coverage`, or
  `## Report Details`.
- Scoped runs still visibly state when the root Area was not evaluated.
- Run-level Contents links match visible run-report sections and order.
- Generated reports use `## Primary Source Data`.
- Run-level Primary Source Data is limited to primary run/report inputs and does
  not duplicate granular report source payloads.
- Detail reports keep their hierarchy/navigation trails.
- Area reports keep `## Area / Factor Breakdown`.
- Durable report specs, `/quality` reporting specs, and report design guidance
  are aligned.
- Generated report gallery output is regenerated.
- Focused Go tests pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define the simplified run report
  shape, `Model Evaluation`, compact scoped-coverage note, and report-local
  `Primary Source Data` contract (requirements 1-5 and 7-10).
- `specs/cli/evaluation-report.md` - mirror the run report shape and
  `Primary Source Data` contract for `qualitymd evaluation report build`
  (requirements 1-5 and 8-10).
- `specs/skills/quality-skill/reporting.md` - mirror `Primary Source Data` and
  run/Area model-structure report requirements for `/quality` report generation
  (requirements 8-9).

### To rename

None

### To delete

None
