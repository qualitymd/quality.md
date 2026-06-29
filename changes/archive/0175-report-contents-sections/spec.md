---
type: Functional Specification
title: Report Contents Sections - functional spec
description: Requirements for generated report Contents sections and removal of compact Jump to lines.
tags: [evaluation, reports, markdown, navigation]
timestamp: 2026-06-29T00:00:00Z
---

# Report Contents Sections - functional spec

Companion to
[Report Contents Sections](../0175-report-contents-sections.md). This spec
states the delta for generated Evaluation Markdown report navigation. The
durable source of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md),
[`qualitymd evaluation report`](../../../specs/cli/evaluation-report.md), and
[`/quality reporting`](../../../specs/skills/quality-skill/reporting.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Generated reports are human review artifacts with stable sections, local keys,
tables, findings, recommendations, and source-data links. A standard Contents
section is the ordinary Markdown shape for navigating those documents, and it is
easier for readers and agents to recognize than a special inline `Jump to:`
idiom. The important boundary is artifact purpose: generated report artifacts
benefit from Contents when they have multiple substantive sections, while OKF
`index.md` files and other listing/index artifacts already are navigation
surfaces and should not add a second local contents list.

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
- generated report filenames and frontmatter fields;
- generated report `type` values;
- local notation keys and marker/label display catalogs;
- non-Markdown or interactive report surfaces.

## Requirements

1. Generated Evaluation Markdown reports **MUST** render a `## Contents` section
   when they contain at least two substantive top-level body sections.

   > Rationale: Multi-section reports need predictable local navigation, and a
   > standard Markdown section is easier to scan and parse than a special inline
   > navigation sentence. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`,
   > `specs/cli/evaluation-report.md`, and
   > `specs/skills/quality-skill/reporting.md`.

2. Generated `## Contents` sections **MUST** list only visible `##` body
   sections that appear later in the same generated report.

   > Rationale: Contents should describe the report's actual top-level shape
   > without becoming a deep outline or linking to itself. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

3. Generated `## Contents` sections **MUST** render as a simple bullet list of
   Markdown links using each target section's visible heading text and anchor.

   > Rationale: The section should remain lightweight, deterministic, and
   > ordinary Markdown. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

4. Generated reports **MUST NOT** render compact `Jump to:` lines.

   > Rationale: Keeping both idioms creates duplicate local navigation and makes
   > report shape less consistent. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`,
   > `specs/cli/evaluation-report.md`, and
   > `specs/skills/quality-skill/reporting.md`.

5. Generated report Contents eligibility **MUST NOT** depend on report length,
   row count, or subjective usefulness once a report has at least two
   substantive top-level body sections.

   > Rationale: Report generation should be deterministic and easy to verify;
   > "materially improves scanning" moved too much judgment into the renderer.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

6. Generated reports **MUST NOT** render a `## Contents` section when the
   artifact is an OKF `index.md`, another listing/index artifact whose primary
   purpose is navigation, or a generated report with fewer than two substantive
   top-level body sections.

   > Rationale: Index artifacts already are navigation surfaces, and a Contents
   > section on a single-section artifact adds noise without navigation value.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

7. The `## Primary Source Data` section **MUST** be eligible for generated
   Contents when it is one of multiple substantive top-level sections in a
   generated report.

   > Rationale: Source data is a report section readers and agents often jump to;
   > excluding it makes Contents less useful. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

8. Generated report frontmatter, H1 titles, local notation keys, and
   `## Primary Source Data` content **MUST NOT** change except as required to
   add Contents and remove `Jump to:`.

   > Rationale: This case corrects local navigation shape without reopening the
   > report identity, display catalog, or provenance contracts. Durable spec:
   > none.

## Acceptance criteria

- Multi-section generated Evaluation reports contain `## Contents`.
- Generated Contents lists visible later `##` sections and does not include
  itself.
- Generated reports do not contain `Jump to:`.
- Generated report frontmatter and H1 title text are unchanged.
- Local notation keys and Primary Source Data lists remain present.
- OKF `index.md` files are not given Contents sections by this change.
- Durable specs, report design guidance, tests, and report gallery output are
  aligned.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define generated Contents
  eligibility, contents list shape, index/listing exceptions, and removal of
  `Jump to:` lines (requirements 1-8).
- `specs/cli/evaluation-report.md` - mirror generated report navigation
  behavior for `qualitymd evaluation report build` (requirements 1, 4, and 6).
- `specs/skills/quality-skill/reporting.md` - mirror generated report
  navigation expectations for `/quality` report workflows (requirements 1, 4,
  and 6).

### To rename

None

### To delete

None
