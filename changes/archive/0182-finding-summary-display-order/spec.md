---
type: Functional Specification
title: Finding Summary Display Order
description: Requirements for renaming the run report Finding Breakdown, showing all Finding types, and using concern-first display order.
tags: [evaluation, reports, findings]
timestamp: 2026-06-30T00:00:00Z
---

# Finding Summary Display Order

This change updates the current generated run report opening. It does not change
structured Finding type values, validation, or ranking semantics.

The key words "MUST" and "MUST NOT" are to be interpreted as described in BCP
14 when, and only when, they appear in all capitals.

## Background / motivation

The current `Finding Breakdown` table shows only Finding types that occur in the
ranked Finding set. That is compact, but it makes absent types implicit and
keeps the reader from seeing the full taxonomy where the report is summarizing
Finding counts. The table also names its third column `Detail`, even though the
only detail shown is severity information for concern types.

## Scope

Covered:

- run report Finding Summary heading and table columns;
- Finding type display order in generated reports and legends;
- zero-count rows for absent Finding types in the run report Finding Summary;
- severity count rendering for `gap` and `risk`;
- durable report specs, CLI report docs, reporting guide examples, tests, and
  generated gallery output.

Deferred:

- changing the structured Finding type vocabulary;
- changing Finding ranking order;
- making Finding `severity` optional or type-specific;
- showing zero-count severity values;
- adding severity summaries to non-run Finding tables.

## Requirements

1. Generated Markdown reports **MUST** display Finding types in this order:
   `🚩 Gap`, `⚠️ Risk`, `✅ Strength`, and `ℹ️ Note`.

   > Rationale: The run opening should surface actionable concerns before
   > supporting strengths and neutral notes. Keeping this in the shared display
   > order avoids a report-specific ordering exception.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

2. The run report **MUST** render the run-level Finding count table under the
   heading `Finding Summary`.

   > Rationale: Once absent Finding types appear as zero-count rows, the table is
   > a fixed taxonomy summary rather than only a breakdown of present categories.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

3. The run report `Finding Summary` table **MUST** render columns `Finding Type`,
   `Count`, and `Severity`, in that order.

   > Rationale: The third column shows severity counts for concern types; naming
   > it `Severity` is clearer than generic detail wording.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

4. The run report `Finding Summary` table **MUST** render one row for every
   current Finding type, including types with count `0`.

   > Rationale: Zero-count rows make absence explicit and keep the report
   > opening self-contained.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

5. In the run report `Finding Summary` table, the `Severity` cell for `gap` and
   `risk` **MUST** list only severity values that occur for that Finding type,
   ordered by the Finding severity catalog. It **MUST NOT** list zero-count
   severity values.

   > Rationale: Finding type absence is useful at taxonomy level, but every
   > zero-count severity would add noise without improving scan value.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

6. In the run report `Finding Summary` table, the `Severity` cell for
   `strength`, `note`, and any concern type with no observed severities **MUST**
   render `—`.

   > Rationale: Severity is meaningful for gaps and risks when present, but the
   > report must not imply severity semantics for strengths or neutral notes.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation/reports/report-tree.md` - update Finding type display order
  and the run report Finding Summary contract.
- `specs/cli/evaluation-report.md` - update the report command contract's run
  report opening description.

### To rename

None.

### To delete

None.
