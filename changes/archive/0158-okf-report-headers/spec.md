---
type: Functional Specification
title: OKF-compatible report headers — functional spec
description: Requirements for generated Evaluation report frontmatter and header navigation.
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# OKF-compatible report headers — functional spec

Companion to
[OKF-compatible report headers](../0158-okf-report-headers.md). This spec states
the delta for generated Evaluation Markdown reports. The durable source of truth
is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md) and
[`/quality reporting`](../../../specs/skills/quality-skill/reporting.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Evaluation reports are human-first review artifacts, and direct links often land
readers in nested detail pages. A consistent header gives readers the run,
subject, and navigation context without path archaeology. Small OKF-compatible
frontmatter gives agents and editors cheap metadata while preserving the
structured Evaluation JSON as the single source of truth.

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

- generated `index.md`, `schema.md`, or `log.md` files for Evaluation run
  folders;
- renaming report files to OKF reserved filenames;
- structured Evaluation schema changes other than adding `RunManifest.createdAt`
  for visible run timing;
- any report renderer dependency on generated report frontmatter.

## Requirements

1. Every generated Markdown report **MUST** begin with YAML frontmatter
   containing only `type`, `title`, and `data`.

   > Rationale: This is enough for OKF-compatible discovery and agent routing
   > while keeping top-of-file noise low.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to require pointer-only
   > frontmatter on generated reports.

2. Report frontmatter `type` **MUST** use the report-subject taxonomy:
   `Evaluation Overview Report`, `Area Evaluation Report`,
   `Factor Evaluation Report`, `Requirement Evaluation Report`,
   `Finding Index Report`, `Recommendation Index Report`, or
   `Recommendation Report`.

   > Rationale: The type should name the report artifact's subject without
   > colliding with Model concepts such as Area, Factor, or Requirement.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to define the report type
   > taxonomy.

3. Report frontmatter `title` **MUST** name the report subject without repeating
   the report-kind prefix rendered in the visible H1.

   > Rationale: `type` already carries report kind for metadata consumers; the
   > visible H1 keeps prefixes such as `Area:` and `Requirement:` for human
   > scanning.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to distinguish frontmatter title
   > from visible report heading.

4. Report frontmatter `data` **MUST** list concrete canonical JSON payloads that
   are sufficient to understand the report, and **MUST** include
   `data/evaluation-output-result.json` for every report.

   > Rationale: `evaluation-output-result.json` indexes the report tree and
   > connected payloads; report-local payloads let agents jump directly to the
   > source data behind the rendered page.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to define report data pointer
   > expectations.

5. Report frontmatter **MUST NOT** duplicate generated time, run identity, model
   snapshot, subject identity, scope, ratings, confidence, summaries, rating
   drivers, findings, recommendations, limits, evidence, or rendered display
   labels when those are available from linked JSON payloads or the visible
   Markdown header.

   > Rationale: Duplicating Evaluation result facts creates a second result
   > format and risks drift from structured data.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to preserve the structured-data
   > source-of-truth boundary.

6. Report generation **MUST NOT** rename report files such as `findings.md` and
   `recommendations.md` to `index.md` for OKF compatibility.

   > Rationale: `index.md` is an OKF bundle listing file; `findings.md` and
   > `recommendations.md` are report concepts with their own subjects.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to keep report paths stable and
   > defer bundle listing files.

7. With frontmatter present, every generated report **MUST** render the visible
   H1 as the first Markdown content after the frontmatter block.

   > Rationale: Markdown readers still need a predictable visible title at the
   > start of the report body.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to update the navigation/H1 rule.

8. Every generated report **MUST** include a visible header area near the H1 with
   report-level navigation to the overview, findings index, and recommendation
   index when those reports exist, plus run number, run creation timestamp, and
   requested Evaluation scope.

   > Rationale: Readers should not have to reconstruct the report tree from
   > paths when they land on a nested detail page, and the run timestamp answers
   > the basic freshness question without putting timestamps in report
   > frontmatter.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md` to require top report
   > navigation.

9. Subject reports **MUST** preserve the existing Model context lines:
   `Area:`, `Factor:`, and plural `Factors:` as applicable.

   > Rationale: Report-level navigation should supplement hierarchy trails, not
   > replace the Model context that explains the subject.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to keep the existing context-line
   > contract alongside the new report navigation.

10. Report header summary tables **MUST** remain report-specific rather than
    becoming a generic `Field | Value` table.

    > Rationale: Specific summary columns make the report state scannable and
    > avoid repeating metadata already carried by title, type, navigation, or
    > structured data.
    > Durable spec: modify
    > `specs/evaluation/reports/report-tree.md` to extend the existing
    > report-specific header guidance.

11. Long reports **SHOULD** include a compact `Jump to:` line when it materially
    improves navigation.

    > Rationale: The run report, deep Requirement reports, and recommendation
    > reports benefit from local section navigation; short Area and Factor
    > reports may not.
    > Durable spec: modify
    > `specs/evaluation/reports/report-tree.md` to recommend jump links for long
    > reports.

12. Report rendering **MUST** remain a deterministic projection over completed
    structured Evaluation data and **MUST NOT** read generated report
    frontmatter as report input.

    > Rationale: Frontmatter is a convenience output, not an input or result
    > source.
    > Durable spec: modify
    > `specs/evaluation/reports/report-tree.md`,
    > `specs/skills/quality-skill/reporting.md`, and, if needed,
    > `specs/cli/evaluation-report.md` to preserve report-source boundaries.

13. `qualitymd evaluation create` **MUST** write `RunManifest.createdAt` as a UTC
    RFC 3339 timestamp.

    > Rationale: Report headers need a stable run-created time from structured
    > data; using report-build time would make rebuilds non-deterministic.
    > Durable spec: modify `specs/cli/evaluation-create.md` and
    > `specs/evaluation/records/payload-kinds.md` to require
    > `RunManifest.createdAt`.

## Acceptance criteria

- Every generated report starts with YAML frontmatter containing only `type`,
  `title`, and `data`.
- Frontmatter `type` values match the taxonomy in requirement 2.
- Frontmatter `title` values omit visible H1 prefixes such as `Area:` and
  `Requirement:`.
- Frontmatter `data` includes `data/evaluation-output-result.json` plus
  report-local payloads, using paths relative to the run root.
- No report frontmatter includes duplicated run, scope, subject, rating,
  summary, finding, recommendation, limit, evidence, or display-label values.
- `findings.md` and `recommendations.md` remain report files; no generated
  report content is renamed to `index.md`.
- The H1 is the first Markdown body content after frontmatter.
- Reports include top navigation to overview, findings, and recommendations
  where those outputs exist.
- Visible report headers include run number, run creation timestamp, and
  requested scope.
- Existing Area, Factor, and Requirement context lines remain present.
- Header tables remain report-specific.
- Focused report-renderer tests cover the frontmatter and header contract.
- Markdown formatting checks pass.
- Go tests pass.
- Report gallery generation is deterministic; regenerating the gallery leaves
  the examples diff unchanged.

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation/reports/report-tree.md` - add the cumulative generated
  report frontmatter, type taxonomy, header navigation, and jump-link contract.
- `specs/skills/quality-skill/reporting.md` - align the skill reporting
  contract with the new generated report header and frontmatter behavior.
- `specs/cli/evaluation-report.md` - update only if the CLI report-build spec
  needs to name the generated frontmatter/header behavior explicitly.
- `specs/cli/evaluation-create.md` - add `RunManifest.createdAt`.
- `specs/evaluation/records/payload-kinds.md` - add `RunManifest.createdAt`.

### To rename

None.

### To delete

None.
