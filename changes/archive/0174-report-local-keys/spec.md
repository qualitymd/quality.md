---
type: Functional Specification
title: Report Local Keys and Navigation - functional spec
description: Requirements for generated report local keys, compact navigation, and footer simplification.
tags: [evaluation, reports, markdown, accessibility]
timestamp: 2026-06-29T00:00:00Z
---

# Report Local Keys and Navigation - functional spec

Companion to
[Report Local Keys and Navigation](../0174-report-local-keys.md). This spec
states the delta for generated Evaluation Markdown reports. The durable source
of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md),
[`qualitymd evaluation report`](../../../specs/cli/evaluation-report.md), and
[`/quality reporting`](../../../specs/skills/quality-skill/reporting.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Generated Markdown reports are human review surfaces and agent-readable
artifacts. Their opening structure should get readers to the result and the
relevant navigation quickly, while report tables should be understandable
without making readers jump to a bottom legend. At the same time, accessibility
requires report values to carry text labels directly in the table cells: local
keys can help orient readers, but they must not make marker-only report content
acceptable.

This change removes heavy generated contents and legend sections, introduces
notation-only local keys immediately after first use, and aligns detail-report
openings with the newer run-report section roles.

## Scope

This change covers generated Evaluation Markdown reports:

- run-level `report.md`;
- root and non-root Area reports;
- Factor reports;
- Requirement reports;
- `findings.md`;
- `recommendations.md`;
- recommendation detail reports.

Deferred:

- structured Evaluation JSON schema changes;
- Rating Scale semantics;
- display catalog value additions or removals;
- report filenames and frontmatter fields;
- non-Markdown or interactive report surfaces.

## Requirements

1. Generated Markdown reports **MUST NOT** render a `## Contents` section.

   > Rationale: A whole contents section is too heavy for generated reports
   > whose opening should stay close to the result. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`,
   > `specs/cli/evaluation-report.md`, and
   > `specs/skills/quality-skill/reporting.md`.

2. Generated Markdown reports **MAY** render a compact `Jump to:` line when
   local section navigation materially improves scanning.

   > Rationale: Inline jump links preserve local navigation without creating a
   > full section. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`,
   > `specs/cli/evaluation-report.md`, and
   > `specs/skills/quality-skill/reporting.md`.

3. The run-level `report.md` **MUST** render H1, `## Summary`,
   `## Key Details`, optional `Jump to:`, and `## Model Evaluation` before Top
   Findings and Top Recommendations.

   > Rationale: The run report remains the decision-ready entrypoint and should
   > keep the result before long ranked tables. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

4. Finding and recommendation list reports **MUST** keep visible H1 titles as
   `# Findings` and `# Recommendations`.

   > Rationale: The frontmatter `type` carries the formal report taxonomy; the
   > visible title should stay reader-simple. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

5. Area, Factor, and Requirement reports **MUST** render opening state tables
   under `## Key Details`.

   > Rationale: The table is scan-critical state, not unlabeled chrome. Durable
   > spec: modify `specs/evaluation/reports/report-tree.md`.

6. Area, Factor, and Requirement reports **MUST** render summary prose under
   `## Summary` instead of a bare `Summary:` label.

   > Rationale: Summary is a report section with the same role across report
   > families. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

7. Generated Markdown reports **MUST NOT** render a bottom `## Legend` section.

   > Rationale: A bottom legend explains notation after first use and competes
   > with primary source data. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`,
   > `specs/cli/evaluation-report.md`, and
   > `specs/skills/quality-skill/reporting.md`.

8. When a generated report first uses an indicator family, it **MUST** render a
   local key immediately after the table or section that introduced that family.

   > Rationale: Keys should appear where readers first need them. Durable spec:
   > modify `specs/evaluation/reports/report-tree.md`.

9. Local keys **MUST** be notation-only lines of the form
   `<Family>: <marker label>, <marker label>.`

   > Rationale: The key orients readers to the marker set without becoming a
   > glossary or duplicating report semantics. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

10. When a table first uses multiple indicator families, generated reports
    **MUST** render each family key on a distinct line.

    > Rationale: One line per family keeps keys scannable and avoids dense
    > mini-paragraphs. Durable spec: modify
    > `specs/evaluation/reports/report-tree.md`.

11. Local keys **MUST NOT** include term definitions, rationale, or explanatory
    prose beyond the marker-label set.

    > Rationale: Definitions belong in the model, structured payloads, and
    > durable docs; generated keys should stay compact. Durable spec: modify
    > `specs/evaluation/reports/report-tree.md`.

12. Local keys for Rating Levels **MUST** render the configured Rating Scale
    labels from the run's model snapshot.

    > Rationale: Rating Levels are model-defined and must not be hardcoded.
    > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

13. Local keys for CLI-owned display catalogs **MUST** render the known catalog
    labels and markers for that family.

    > Rationale: Confidence, finding type, severity, recommendation impact, and
    > advice tier labels are CLI-owned display sets and should remain consistent
    > with table cells. Durable spec: modify
    > `specs/evaluation/reports/report-tree.md`.

14. Generated report table cells **MUST** continue to render text labels for
    ratings, statuses, confidence, finding type, severity, recommendation
    impact, and priority-like values, optionally preceded by markers or icons.

    > Rationale: The local key is supplemental; accessible meaning must remain
    > in the content itself. Durable spec: modify
    > `specs/evaluation/reports/report-tree.md`,
    > `specs/cli/evaluation-report.md`, and
    > `specs/skills/quality-skill/reporting.md`.

15. Generated report table cells **MUST NOT** rely on color or marker shape
    alone to convey semantic values.

    > Rationale: Reports are read in terminals, editors, plain text, and by
    > agents; color and icon shape are not reliable carriers of meaning.
    > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

16. The `## Primary Source Data` section **MUST** remain the only standard
    footer utility section.

    > Rationale: Source-data links are the stable bridge to structured inputs;
    > notation belongs near first use. Durable spec: modify
    > `specs/evaluation/reports/report-tree.md` and
    > `specs/skills/quality-skill/reporting.md`.

17. Generated report frontmatter **MUST NOT** change as part of this case.

    > Rationale: Frontmatter already carries non-judgmental report identity and
    > routing metadata. Durable spec: none.

## Acceptance criteria

- No generated report contains `## Contents`.
- No generated report contains `## Legend`.
- `report.md` opens with H1, `## Summary`, `## Key Details`, optional
  `Jump to:`, and `## Model Evaluation` before ranked lists.
- `findings.md` renders `# Findings`; `recommendations.md` renders
  `# Recommendations`.
- Area, Factor, and Requirement reports render `## Key Details` for their
  opening state table and `## Summary` for summary prose.
- Local keys appear immediately after first use of each indicator family.
- Tables that first use multiple indicator families render one local-key line
  per family.
- Local keys are marker-label sets only, with no term definitions or rationale.
- Rating keys are derived from the run model snapshot.
- CLI-owned display keys are derived from current display catalogs.
- Report table cells still include text labels beside markers/icons.
- `## Primary Source Data` remains the only standard footer utility.
- Durable specs, report design guidance, tests, and report gallery output are
  aligned.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define compact navigation, local
  keys, key formatting, accessible label requirements, no `Contents`/`Legend`,
  and detail-report `Key Details`/`Summary` sections (requirements 1-17).
- `specs/cli/evaluation-report.md` - mirror generated report output and footer
  behavior for `qualitymd evaluation report build` (requirements 1-3, 7, 14,
  and 16).
- `specs/skills/quality-skill/reporting.md` - mirror generated report
  navigation, key, footer, and label expectations for `/quality` report
  workflows (requirements 1-3, 7, 14, and 16).

### To rename

None

### To delete

None
