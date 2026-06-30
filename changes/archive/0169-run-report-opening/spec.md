---
type: Functional Specification
title: Run Report Opening - functional spec
description: Requirements for generated report.md opening structure, contents, and non-judgmental metadata.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Run Report Opening - functional spec

Companion to [Run Report Opening](../0169-run-report-opening.md). This spec
states the delta for generated Evaluation run reports. The durable source of
truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md) and
[`qualitymd evaluation report`](../../../specs/cli/evaluation-report.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The run-level `report.md` is the primary human entrypoint for an Evaluation
run. Its opening should lead with judgment, decision-critical details, and
navigation. Run IDs, creation timestamps, and stable subject references are
still valuable for agents and traceability, but they should not crowd the first
visible lines when frontmatter can carry them and a lower details section can
preserve visible traceability.

## Scope

This change covers generated run-level `report.md` only:

- frontmatter metadata;
- visible opening order;
- `Summary`, `Key Details`, `Contents`, and `Report Details` sections;
- omission of visible `Limits & Incomplete Inputs`.

Deferred:

- generated detail report openings;
- Evaluation structured data schema changes;
- report output paths;
- rating, roll-up, Advice, coverage, and source-data semantics;
- future visible limits treatment.

## Requirements

1. Generated run-level `report.md` frontmatter **MUST** include `type`,
   `title`, `run`, `runId`, `created`, `scope`, and `subject` fields.

   > Rationale: These fields are non-judgmental routing and indexing metadata
   > that agents and static tooling can consume without forcing low-value run
   > facts into the visible opening. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

2. Generated run-level `report.md` frontmatter **MUST NOT** include ratings,
   confidence, summaries, findings, recommendations, limits text, evidence, or
   source-data manifests.

   > Rationale: Evaluation judgment stays in structured payloads and visible
   > report sections; frontmatter must not become a second result format.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

3. Generated run-level `report.md` **MUST NOT** render the previous visible
   top `Run:` line.

   > Rationale: Run number, run ID, creation time, scope, and subject are now
   > either frontmatter metadata, key decision details, or lower trace details.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

4. Generated run-level `report.md` **MUST** render the visible opening in this
   order after the H1: report navigation line, Area context line, `## Summary`,
   `## Key Details`, and `## Contents`.

   > Rationale: The primary report should lead with the answer and scan-critical
   > facts before detailed finding, recommendation, and breakdown tables.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

5. The `## Summary` section **MUST** render the scoped Area summary and
   **SHOULD** render a deterministic recommended-next-action sentence when a
   ranked recommendation is available.

   > Rationale: A run report summary should state the bottom line and the best
   > next move without requiring the reader to scan the recommendation table.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

6. The `## Key Details` section **MUST** render a compact table containing
   Overall Rating, Confidence, Scope, Findings, and Recommendations.

   > Rationale: These facts are decision-critical enough for the opening but are
   > not prose summary. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

7. The `## Key Details` section **MUST NOT** include limits or incomplete-input
   counts.

   > Rationale: Limits display is being removed from the visible run report for
   > now; keeping a limits count in the opening would preserve the same concern
   > in summary form. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

8. The `## Contents` section **MUST** replace the previous `Jump to:` line and
   **MUST** link only to visible sections that the run report renders.

   > Rationale: The run report is a primary artifact, so contents should be a
   > first-class section, while hidden or omitted sections should not be linked.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

9. Generated run-level `report.md` **MUST NOT** render `## Limits & Incomplete
Inputs`.

   > Rationale: The visible limits section is deliberately removed from the run
   > report for this pass, without changing structured Area Analysis data.
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

10. Generated run-level `report.md` **MUST** render a lower `## Report Details`
    section before `## Legend` and `## Source Data`, with Run, Run ID, Created,
    and Subject rows.

    > Rationale: Visible traceability stays available without occupying the
    > first screen. Durable spec: modify
    > `specs/evaluation/reports/report-tree.md`.

11. Generated non-run reports **MUST NOT** change their visible run-context
    headers, frontmatter fields, or local navigation as part of this case.

    > Rationale: Deep-linked detail reports still benefit from immediate run
    > context, and this change is scoped to the primary run report. Durable
    > spec: none; existing durable specs already govern detail reports.

## Acceptance criteria

- `report.md` frontmatter contains `type`, `title`, `run`, `runId`, `created`,
  `scope`, and `subject`.
- `report.md` frontmatter does not contain ratings, findings,
  recommendations, limits text, or source-data paths.
- `report.md` does not render the top `Run:` line.
- `report.md` renders `## Summary`, `## Key Details`, and `## Contents` before
  `## Top Findings`.
- `## Summary` includes the scoped Area summary and top recommendation title
  when a ranked recommendation exists.
- `## Key Details` includes rating, confidence, scope, ranked finding count, and
  ranked recommendation count.
- `## Contents` links only to visible rendered sections.
- `report.md` does not render `## Limits & Incomplete Inputs`.
- `## Report Details` renders before `## Legend` and `## Source Data`.
- Structured source data and Evaluation payloads remain unchanged.
- Generated report gallery output is regenerated.
- Focused Go tests pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define run report frontmatter,
  opening order, contents, report details, and visible limits removal
  (requirements 1-10).
- `specs/cli/evaluation-report.md` - align `qualitymd evaluation report build`
  with the new run report shape and visible limits removal (requirements 1, 4,
  and 9).

### To rename

None

### To delete

None
