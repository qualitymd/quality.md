---
type: Functional Specification
title: Glossary and Report Links
description: Requirements for adding a shared glossary and replacing generated report legends with compact Evaluation links.
tags: [evaluation, reports, glossary]
timestamp: 2026-06-30T00:00:00Z
---

# Glossary and Report Links

This change adds a shared glossary for generated report readers and replaces
per-report local legends with compact links to the key Evaluation artifacts.
It does not change structured Evaluation data values, Rating Scale semantics, or
the requirement that generated reports render text labels for semantic values.

The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

Generated reports repeat local `Legend` blocks to explain markers for ratings,
confidence, Finding types, severities, recommendation impact, and similar
indicator families. That made sense while reports lacked a shared reference
surface, but it scatters vocabulary display sets across every artifact and keeps
definitions out of reach. A workspace-root glossary can define shared terms and
fixed vocabularies once, while reports stay compact and link to the glossary
alongside the overview, Findings, and Recommendations artifacts.

## Scope

Covered:

- workspace-root `glossary.md` structure and initial content;
- core glossary terms for Area, Factor, Finding, Recommendation, and
  Requirement;
- a `Quality rating` glossary entry derived from the configured quality model;
- fixed Evaluation enum catalog entries in glossary table form;
- generated report `Evaluation links:` navigation;
- removal of generated local `Legend` blocks;
- durable report specs, report design guidance, generated examples, and tests.

Deferred:

- generated glossary maintenance;
- a CLI glossary or enum-catalog inspection command;
- linking every table cell or enum value to the glossary;
- changing canonical enum values, aliases, validation, or schema output;
- changing model-defined Rating Level semantics.

## Requirements

1. The QUALITY.md workspace root **MUST** include `glossary.md` as the shared
   human glossary for QUALITY.md report readers.

   > Rationale: Generated reports need a durable definition surface outside any
   > one report artifact.
   >
   > Durable spec: add `specs/glossary-md.md`.

2. `glossary.md` **MUST** organize entries as one flat alphabetical list of
   `##` headings.

   > Rationale: The glossary will mix terms from model vocabulary, report
   > vocabulary, fixed Evaluation enum catalogs, and configured model values, so
   > source-based categories would make lookup harder.
   >
   > Durable spec: add `specs/glossary-md.md`.

3. `glossary.md` **MUST** include initial concept entries for `Area`, `Factor`,
   `Finding`, `Recommendation`, and `Requirement`.

   > Durable spec: add `specs/glossary-md.md`.

4. Fixed vocabulary entries in `glossary.md` **MUST** use a table with columns
   `Label`, `Value`, and `Description`, in that order.

   > Rationale: The table shape distinguishes human report display from
   > canonical persisted values while keeping every vocabulary entry consistent.
   >
   > Durable spec: add `specs/glossary-md.md`.

5. Fixed vocabulary entry tables **MUST** render values in their catalog order,
   not alphabetical order.

   > Rationale: Several vocabularies carry severity, impact, priority, or Rating
   > Scale order; sorting those values alphabetically would erase useful
   > ordering.
   >
   > Durable spec: add `specs/glossary-md.md`.

6. `glossary.md` **MUST** include a `Quality rating` entry whose table uses the
   configured Rating Scale from this repository's `QUALITY.md`.

   > Rationale: Quality ratings render like a vocabulary in reports, but they
   > are model-defined rather than fixed Evaluation enum values; readers need to
   > know the labels and values come from the quality model.
   >
   > Durable spec: add `specs/glossary-md.md`.

7. The `Quality rating` entry **MUST** state that its labels and values come
   from this project's `QUALITY.md` Rating Scale.

   > Rationale: This is the one source note that helps report readers because
   > Rating Levels can vary across models.
   >
   > Durable spec: add `specs/glossary-md.md`.

8. `glossary.md` **MUST** include fixed Evaluation enum catalog entries for
   Analysis status, Assessment status, Confidence, Data kind, Finding basis,
   Finding coverage, Finding rank, Finding severity, Finding type, Rating
   result, Rating status, Recommendation impact, Report kind, and Run gap kind.

   > Durable spec: add `specs/glossary-md.md`.

9. Generated Markdown report artifacts **MUST** render a compact navigation line
   whose label is exactly `Evaluation links:`.

   > Rationale: Reports need cross-artifact navigation after local legends are
   > removed, but that navigation should not become another report section.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/cli/evaluation-report.md`.

10. The generated `Evaluation links:` line **MUST** render inline links in this
    order with filename link text: `report.md`, `findings.md`,
    `recommendations.md`, and `glossary.md`.

    > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

11. Generated `Evaluation links:` targets **MUST** be relative to the current
    report artifact.

    > Rationale: Generated reports are browsed from nested Area, Factor,
    > Requirement, and recommendation detail paths; links must work from the file
    > a reader opened.
    >
    > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

12. Generated reports **MUST** render `Evaluation links:` after the report's
    opening summary, key details, or report-specific orientation, and before
    `## Contents` when a Contents section is present.

    > Rationale: Readers should see the report's local answer first, then
    > cross-artifact navigation, then in-page navigation.
    >
    > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

13. Generated reports **MUST NOT** render local `Legend` blocks.

    > Rationale: The glossary becomes the definition and vocabulary reference
    > surface; repeated local legends would preserve the duplication this change
    > removes.
    >
    > Durable spec: modify `specs/evaluation/reports/report-tree.md`,
    > `specs/cli/evaluation-report.md`, and `docs/guides/reporting-design.md`.

14. Generated report cells **MUST** continue to render text labels for ratings,
    statuses, confidence, Finding type, Finding severity, recommendation impact,
    Finding rank, Finding coverage, and other semantic values.

    > Rationale: Removing legends must not make report tables depend on marker
    > shape, color, or external reference for basic comprehension.
    >
    > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

15. Generated reports **MUST NOT** link every table cell or enum value to
    `glossary.md`.

    > Rationale: One stable navigation link keeps reports readable and avoids
    > noisy Markdown while still making definitions reachable.
    >
    > Durable spec: modify `specs/evaluation/reports/report-tree.md`.

16. The valid Finding type list in
    `specs/evaluation/records/payload-kinds.md` **MUST** omit the stale
    `unknown` value.

    > Rationale: The active typed catalog and validation accept `strength`,
    > `gap`, `risk`, and `note`; the durable payload spec should not contradict
    > that contract.
    >
    > Durable spec: modify `specs/evaluation/records/payload-kinds.md`.

## Durable spec changes

### To add

- `specs/glossary-md.md` - define the workspace-root `glossary.md` artifact contract,
  including purpose, location, flat alphabetical entry structure, fixed
  vocabulary table shape, initial concept terms, quality rating source note, and
  fixed Evaluation enum entries.

### To modify

- `specs/evaluation/reports/report-tree.md` - define the glossary-backed report
  vocabulary reference, `Evaluation links:` navigation line, and removal of
  local `Legend` blocks.
- `specs/cli/evaluation-report.md` - align the CLI report command contract with
  generated report navigation and removed legends.
- `specs/evaluation/records/payload-kinds.md` - remove stale `unknown` from the
  valid Finding type list.

### To rename

None.

### To delete

None.
