---
type: Functional Specification
title: Quality Evaluation Title - functional spec
description: Requirements for generated run report titles that name quality evaluation scope.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Quality Evaluation Title - functional spec

Companion to
[Quality Evaluation Title](../0168-quality-evaluation-title.md). This spec
states the delta for generated Evaluation run reports. The durable source of
truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md) and
[`qualitymd evaluation report`](../../../specs/cli/evaluation-report.md).

The key words **MUST** and **MUST NOT** are to be interpreted as described in
BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The run report title should read as the user's quality-evaluation result, not as
a stack of internal report and subject labels. `report.md` and the frontmatter
`type` already identify the artifact as a report; the H1 should name the
quality evaluation scope. When a user asks for specific Factors, the title
should preserve the Area context and include those Factor labels because they
are the primary requested scope.

## Scope

This change covers generated `report.md` only:

- its frontmatter `title`;
- its visible H1;
- Area-only titles;
- factor-filtered titles with one or more Factors.

Deferred:

- detail report titles for Area, Factor, Requirement, and Recommendation
  reports;
- index report titles;
- generated report `type` values;
- report file paths, slugs, and navigation labels.

## Requirements

1. Generated run report frontmatter `title` and visible H1 **MUST** use
   `Quality Evaluation - <Area title>` when the planned scope has no factor
   filter.

   > Rationale: The run report is already identified as a report artifact by
   > path and `type`; the visible title should name the quality-evaluation
   > subject. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

2. Generated run report frontmatter `title` and visible H1 **MUST** use
   `Quality Evaluation - <Area title> (<Factor title list>)` when the planned
   scope includes one or more factor filters.

   > Rationale: Factor-scoped evaluations are user-requested factor concerns in
   > an Area context, so the title should keep both the Area and the exact
   > Factor scope. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

3. The `<Factor title list>` **MUST** render all planned factor filters as
   comma-separated Factor titles in the order stored in
   `RunManifest.plannedScope.factorFilter`.

   > Rationale: Truncation hides requested scope, and the manifest preserves the
   > stable planned order the user and skill already resolved. Durable spec:
   > modify `specs/evaluation/reports/report-tree.md`.

4. Generated run report titles **MUST NOT** include `Evaluation Report`, `Area:`,
   raw Area references, or raw Factor references in the H1 or frontmatter
   `title`.

   > Rationale: Those values either duplicate report metadata or belong in the
   > visible scope table for traceability. Durable spec: modify
   > `specs/evaluation/reports/report-tree.md`.

5. Generated run report scope details **MUST** continue to render the planned
   Area reference and factor-filter references in the Scope section.

   > Rationale: Human-friendly titles should not remove stable references needed
   > for traceability and agent follow-up. Durable spec: none; this preserves
   > existing durable behavior.

6. Generated Area, Factor, Requirement, index, and recommendation detail report
   titles **MUST NOT** change as part of this case.

   > Rationale: The title problem is specific to the run entrypoint; detail
   > reports still benefit from kind-prefixed titles as disambiguators. Durable
   > spec: none; existing durable specs already define those titles.

## Acceptance criteria

- Area-only `report.md` emits frontmatter `title` content
  `Quality Evaluation - <Area>` and H1 `# Quality Evaluation - <Area>`.
- Single-Factor `report.md` emits
  `Quality Evaluation - <Area> (<Factor>)`.
- Multiple-Factor `report.md` emits every planned Factor title in manifest
  order, comma-separated, inside the parentheses.
- The run report Scope section still shows stable planned Area and factor-filter
  references.
- Detail report title expectations remain unchanged.
- Report gallery output is regenerated.
- Focused Go tests pass.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - define run report title wording
  and factor-filter title rendering (requirements 1-4).
- `specs/cli/evaluation-report.md` - align `qualitymd evaluation report build`
  with the `Quality Evaluation - <Area> (<Factors>)` run report title contract
  (requirements 1-2).

### To rename

None

### To delete

None
