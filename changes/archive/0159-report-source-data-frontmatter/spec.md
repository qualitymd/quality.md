---
type: Functional Specification
title: Report source-data frontmatter — functional spec
description: Requirements for source-data report frontmatter and body data-link removal.
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report source-data frontmatter — functional spec

Companion to
[Report source-data frontmatter](../0159-report-source-data-frontmatter.md).
This spec states the delta for generated Evaluation Markdown reports. The
durable source of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md) and
[`/quality reporting`](../../../specs/skills/quality-skill/reporting.md).

The key words **MUST** and **MUST NOT** are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Generated report frontmatter should give agents and secondary tools a cheap,
precise way to find the source data behind a report artifact. Once that source
manifest exists, visible report headers do not need a duplicate `Data` column.
Keeping source-data pointers in one place makes report bodies narrower and keeps
the human review surface focused on ratings, status, findings, recommendations,
limits, and navigation.

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
- new report frontmatter fields;
- treating generated report frontmatter as renderer input;
- listing `model-snapshot.md` in `data`; the `data` field is limited to
  structured Evaluation payloads under `data/`.

## Requirements

1. Generated report frontmatter `data` **MUST** list the run-root-relative
   structured Evaluation data files used as source data for that specific
   Markdown report artifact.

   > Rationale: `data` is a source-data manifest, not a general related-files
   > index. A consumer opening a report should be able to tell which persisted
   > Evaluation payloads back that artifact.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to define `data` as the
   > report-local source-data manifest.

2. Generated report frontmatter `data` **MUST NOT** include
   `data/evaluation-output-result.json` solely because that generated index
   exists.

   > Rationale: The report builder assembles `EvaluationOutputResult` as an
   > output index; it is not source data for the Markdown reports unless a report
   > later renders directly from it.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to remove the blanket
   > `evaluation-output-result.json` frontmatter requirement.

3. Generated report frontmatter `data` **MUST** include `data/run-manifest.json`
   when the report body renders run number, creation time, or requested scope
   from the run manifest.

   > Rationale: The visible run context is sourced from the manifest and should
   > be reflected in the report-local source-data manifest.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` to require listing rendered
   > manifest data when used.

4. Generated report Markdown bodies **MUST NOT** render report-level source-data
   links in header summary `Data` columns or equivalent body-only source-data
   link lines.

   > Rationale: Source-data links have one home in report frontmatter; duplicating
   > them in the body makes report headers wider without adding human judgment
   > context.
   > Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md` to remove body-level `Data`
   > columns from generated report headers.

5. Generated report bodies **MUST** keep visible human orientation content:
   report title, run context, scope and coverage, report navigation, hierarchy
   context, ratings, confidence, status, summaries, findings, recommendations,
   limits, and incomplete inputs.

   > Rationale: Frontmatter is not a substitute for the human report body; it
   > only owns source-data pointers.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` to preserve
   > human-first report body obligations while moving source-data links out.

6. Report generation **MUST NOT** read generated report frontmatter as
   report-generation input.

   > Rationale: Frontmatter is an output convenience and source-data manifest,
   > not an input channel or cached result format.
   > Durable spec: modify none; existing durable specs already require this, and
   > this case preserves it.

## Acceptance criteria

- Every generated report frontmatter `data` entry names a structured source
  payload under `data/` used to render that report.
- Generated report frontmatter no longer lists
  `data/evaluation-output-result.json` by default.
- Generated reports that render run context list `data/run-manifest.json`.
- Generated report header summary tables no longer contain `Data` columns.
- Generated report bodies no longer contain report-level links to
  `*-result.json` payloads solely as source-data pointers.
- Visible run context, navigation, hierarchy context, ratings, status,
  confidence, summaries, findings, recommendations, limits, and incomplete
  inputs remain present.
- Focused Go tests cover the frontmatter source-data manifest and missing body
  `Data` columns.
- Report gallery generation is deterministic and checked in.
- `mise run check` passes.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define `data` frontmatter as the
  report-local source-data manifest, remove blanket
  `evaluation-output-result.json` frontmatter and body data-link requirements,
  and remove `Data` from report header summaries (requirements 1-5).
- `specs/skills/quality-skill/reporting.md` - align skill reporting expectations
  with source-data frontmatter and body headers without `Data` columns
  (requirement 4).
- `specs/cli/evaluation-report.md` - clarify report build source-data
  boundaries if the durable CLI spec needs an explicit pointer (requirements
  1-2).

### To rename

None

### To delete

None
