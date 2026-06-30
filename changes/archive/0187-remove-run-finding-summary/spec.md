---
type: Functional Specification
title: Remove Run Finding Summary - functional spec
description: Requirements for removing the redundant run-level Finding Summary table from generated report.md.
tags: [evaluation, reports, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Remove Run Finding Summary - functional spec

Companion to
[Remove Run Finding Summary](../0187-remove-run-finding-summary.md). This spec
states the delta for generated run-level `report.md` output. The durable source
of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md) and
the `qualitymd evaluation report` command contract.

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The run report opening should let a reader understand the overall judgment and
then move quickly to the complete ranked lists. After the full Findings report
link gained the complete ranked count and inline type/severity summary, the
separate run-level `Finding Summary` table became redundant: it repeats the
count breakdown without providing the destination link or the capped-row context.
Removing it keeps the opening stack focused while preserving both the total count
and the richer linked breakdown.

## Scope

Covered: generated Markdown for run-level `report.md`, durable report specs,
reporting guidance, tests, changelog, and report-gallery output.

Non-goals:

- changing `findings.md` or ranked Finding table contents;
- changing Requirement report `## Findings Summary` sections;
- changing structured Evaluation data, enum values, validation, or persisted
  JSON payload shapes;
- changing generated report filenames, navigation blockquotes, or primary
  source-data links.

## Requirements

1. Generated run-level `report.md` **MUST NOT** render the standalone
   `Finding Summary` table between `## Key Details` and `## Contents`.

   > Rationale: The same count hierarchy is now available beside the full
   > `findings.md` link, where it explains the capped Top Findings table and
   > points to the complete report.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

2. Generated run-level `report.md` **MUST** keep the `Findings` total in the
   `## Key Details` table.

   > Rationale: The opening still needs a one-second total count before the
   > reader reaches the Top Findings section.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

3. Generated run-level `report.md` **MUST** keep the full Findings report link
   and complete inline Finding count summary immediately below the
   `## Top Findings` heading when ranked Findings exist.

   > Rationale: The linked summary is the remaining count-breakdown surface and
   > must continue to explain that `report.md` is capped while `findings.md` is
   > complete.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

4. Requirement reports **MUST** keep their `## Findings Summary` sections.

   > Rationale: Requirement-level summaries describe local assessment Findings
   > and are not made redundant by the run-level full Findings report link.
   > Durable spec: none.

5. The change **MUST NOT** alter structured Evaluation data enum values,
   validation, or persisted JSON payload shapes.

   > Rationale: This is a generated Markdown simplification only; source data
   > and compatibility surfaces stay unchanged.
   > Durable spec: none.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - remove the run-report Finding
  Summary table contract and preserve the Key Details total plus full Findings
  report count summary contracts (requirements 1, 2, and 3).
- `specs/cli/evaluation-report.md` - remove the compact Finding Summary from
  the `report.md` opening-stack description (requirement 1).

### To rename

None

### To delete

None

## Acceptance criteria

- Generated run-level `report.md` has no standalone `Finding Summary` table.
- Generated run-level `report.md` still shows the total Finding count in
  `## Key Details`.
- Generated run-level `report.md` still shows the full Findings report link with
  a complete inline count summary under `## Top Findings`.
- Requirement report `## Findings Summary` output remains present.
- Report tests, formatting, and report-gallery checks pass.
