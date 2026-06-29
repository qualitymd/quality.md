---
type: Functional Specification
title: Report Source Data Section — functional spec
description: Requirements for moving generated report source-data pointers from frontmatter to a visible section.
tags: [evaluation, reports, markdown, okf]
timestamp: 2026-06-29T00:00:00Z
---

# Report Source Data Section — functional spec

Companion to [Report Source Data Section](../0162-report-source-data-section.md).
This spec states the delta for generated Evaluation Markdown reports. The
durable source of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md),
[`/quality reporting`](../../../specs/skills/quality-skill/reporting.md), and
[`qualitymd evaluation report`](../../../specs/cli/evaluation-report.md).

The key words **MUST** and **MUST NOT** are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Generated report source-data pointers are useful traceability, but the current
YAML `data` list makes report frontmatter noisy and pushes human report content
down. A stable visible section at the bottom of each report gives humans and
agents the same discoverable source links without turning the first screen into
a machine manifest.

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
- treating generated report frontmatter or Markdown bodies as renderer input;
- listing `model-snapshot.md` as source data; the section is limited to
  structured Evaluation payloads under `data/`.

## Requirements

1. Generated report frontmatter **MUST** contain report identity fields only:
   `type` and `title`.

   > Rationale: Report frontmatter should identify the generated report artifact
   > without forcing readers through report-local machine pointers before the H1.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md` to remove generated report frontmatter
   > `data`.

2. Every generated Markdown report **MUST** end with a `## Source Data` section.

   > Rationale: A stable bottom section is easy for humans to scan past and easy
   > for agents to find by heading.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`,
   > `specs/skills/quality-skill/reporting.md`, and
   > `specs/cli/evaluation-report.md` to define the source-data section as the
   > report-local source pointer surface.

3. A report's `Source Data` section **MUST** list the run-root-relative
   structured Evaluation payload paths used as source data for that specific
   Markdown report artifact.

   > Rationale: The section replaces the current report-local manifest; it
   > should not become a run-global index.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md` to preserve report-local source-data
   > semantics outside frontmatter.

4. Source-data list items **MUST** render as Markdown links whose label is the
   run-root-relative payload path and whose target is relative to the report
   file that contains the section.

   > Rationale: Path labels preserve agent discoverability; relative targets
   > keep nested reports useful in GitHub, editors, and local Markdown viewers.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` to define
   > source-data link rendering.

5. A report's `Source Data` section **MUST** include `data/run-manifest.json`
   when the report body renders run number, creation time, or requested scope
   from the run manifest.

   > Rationale: The visible run context is sourced from the manifest and should
   > remain traceable after the pointer surface moves out of frontmatter.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` to preserve
   > manifest source-data listing when manifest fields are rendered.

6. A report's `Source Data` section **MUST NOT** include
   `data/evaluation-output-result.json` solely because that generated output
   index exists.

   > Rationale: `EvaluationOutputResult` indexes generated outputs after report
   > build; it is not source data for a report unless a report is directly
   > rendered from it.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md` to preserve the existing output-index
   > boundary after the pointer move.

7. Generated report bodies **MUST NOT** duplicate report-level source-data links
   in header summary `Data` columns or equivalent header source-data lines.

   > Rationale: The new bottom section is the visible source-data pointer
   > surface; duplicating it in headers would reintroduce the noise this change
   > removes.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md` to keep report headers focused on
   > report state.

8. Report generation **MUST NOT** read generated report frontmatter or Markdown
   body content as report-generation input.

   > Rationale: Moving source-data pointers into the body must not turn generated
   > reports into an input channel or cached result format.
   > Durable spec: modify none; existing durable specs already require this, and
   > this case preserves it.

## Acceptance criteria

- Every generated report frontmatter contains `type` and `title`, and no `data`
  key.
- Every generated report ends with `## Source Data`.
- Each `Source Data` section lists report-local structured payloads under
  `data/` as path-labeled Markdown links.
- Nested reports link source-data targets with correct relative paths.
- Generated reports that render run context link `data/run-manifest.json`.
- Generated reports do not link `data/evaluation-output-result.json` by default.
- Generated report header summary tables still omit `Data` columns.
- Focused Go tests cover frontmatter without `data`, bottom source-data
  sections, missing `evaluation-output-result.json`, and nested relative links.
- Report gallery generation is deterministic and checked in.
- `mise run check` passes.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define identity-only generated
  report frontmatter and the bottom report-local `Source Data` section
  (requirements 1-7).
- `specs/skills/quality-skill/reporting.md` - align skill reporting
  expectations with visible source-data sections and headers without `Data`
  columns (requirements 2, 7).
- `specs/cli/evaluation-report.md` - align CLI report build output requirements
  with identity frontmatter and visible source-data sections (requirements 1-6).

### To rename

None

### To delete

None
