---
type: Functional Specification
title: Report Frontmatter H1 Titles - functional spec
description: Requirements for aligning generated report frontmatter titles with visible H1 titles.
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report Frontmatter H1 Titles - functional spec

Companion to
[Report Frontmatter H1 Titles](../0167-report-frontmatter-h1-titles.md). This
spec states the delta for generated Evaluation Markdown reports. The durable
source of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md),
[`qualitymd evaluation report`](../../../specs/cli/evaluation-report.md), and
[`/quality reporting`](../../../specs/skills/quality-skill/reporting.md).

The key words **MUST** and **MUST NOT** are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Generated Markdown reports should have one document title. The current contract
puts a subject-only title in frontmatter and a kind-prefixed title in the H1,
even though `type` already records the report artifact kind. Aligning
frontmatter `title` with the visible H1 keeps metadata, editor previews, and
human report bodies from naming the same report differently.

## Scope

This change covers generated Markdown reports only:

- `report.md`;
- `findings.md`;
- `recommendations.md`;
- recommendation detail reports;
- root and non-root Area reports;
- Factor reports;
- Requirement reports.

Deferred:

- changing generated report `type` values;
- adding frontmatter subject fields;
- changing report paths, title slugs, or recommendation detail filenames;
- treating generated report frontmatter or Markdown bodies as renderer input.

## Requirements

1. Generated report frontmatter `title` **MUST** equal the plain-text content of
   the report's H1 title line, with the leading Markdown `#` marker removed.

   > Rationale: The frontmatter title should identify the generated Markdown
   > document, not a second subject-only view of the document.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md` to define frontmatter `title` as the H1
   > document title.

2. Generated Area, Factor, Requirement, and recommendation detail reports
   **MUST** keep their report-kind prefix in both the H1 and frontmatter
   `title`.

   > Rationale: The prefix is part of the reader-facing document title and
   > disambiguates identically named Model concepts or recommendations.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md` to make kind-prefixed titles the
   > shared report identity.

3. The run-level `report.md` frontmatter `title` **MUST** match the run-level H1
   that identifies the report as an Evaluation Report.

   > Rationale: The run entrypoint is a report document, not just an Area
   > subject page; its metadata should carry the same entrypoint title as the
   > visible page.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md` to align the run report title contract.

4. Finding and recommendation index reports **MUST** continue to emit
   frontmatter `title` values that match their visible H1 titles.

   > Rationale: Index reports already satisfy the one-title contract and should
   > not gain extra prefixes or alternative metadata names.
   > Durable spec: none.

5. Generated report frontmatter **MUST NOT** add subject-title, display-label,
   path, scope, rating, confidence, source-data, or run-identity fields to
   preserve the previous subject-only title separately.

   > Rationale: `type` and the H1-matching `title` are the identity layer; extra
   > fields would reintroduce summary metadata that belongs in structured
   > payloads or visible report body content.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` to keep the
   > frontmatter identity layer tiny.

6. Report generation **MUST NOT** read generated report frontmatter or Markdown
   body content as report-generation input.

   > Rationale: Aligning metadata with the H1 must not turn generated reports
   > into source data or a cached result format.
   > Durable spec: none; existing durable specs already require this, and this
   > case preserves it.

## Acceptance criteria

- Every generated report frontmatter contains only `type` and `title`.
- Every generated report frontmatter `title` equals the report's H1 text without
  the leading `#` marker.
- Run, Area, Factor, Requirement, and recommendation detail report frontmatter
  titles include the same kind prefix as their H1 titles.
- Finding and recommendation index frontmatter titles remain unchanged because
  they already match their H1 titles.
- Focused Go tests cover the H1-matching frontmatter title contract for the run
  report, root Area report, Factor report, Requirement report, recommendation
  detail report, and index reports.
- Report gallery generation is deterministic and checked in.
- `mise run check` passes.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define generated report
  frontmatter `title` as the plain-text H1 content and keep frontmatter limited
  to `type` and `title` (requirements 1-5).
- `specs/cli/evaluation-report.md` - align `qualitymd evaluation report build`
  with the H1-matching generated report title contract (requirements 1, 3).
- `specs/skills/quality-skill/reporting.md` - align skill reporting
  expectations with generated report titles as report document titles
  (requirement 2).

### To rename

None

### To delete

None
