---
type: Functional Specification
title: Filename text for evaluation data links - functional spec
description: Requirements for rendering Evaluation v2 Data-column links with their payload base filename as link text while leaving link targets and data paths unchanged.
tags: [cli, evaluation, reports]
timestamp: 2026-06-26T00:00:00Z
---

# Filename text for evaluation data links - functional spec

Companion to the
[Filename text for evaluation data links](../0109-evaluation-data-link-filenames.md)
change case. This spec states what the link-text change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

Evaluation v2 reports link each subject's machine-readable payloads from a
`Data` column. The targets already point at the structured `*-result.json`
files, but the link text is a generic word (`analysis`, `assessment`, `rating`)
that does not identify the file and restates context the report already carries.
Naming the payload file in the link text lets a reader see exactly which
`*-result.json` a link opens, and keeps the rendered text consistent with the
file the link resolves to.

## Scope

Covered: generated Evaluation v2 `Data`-column link text on Area, Factor, and
Requirement reports, and the durable contract that governs it.

Not covered: data-link targets, structured routine data paths under `data/`, the
`Data` column header, navigation trails, subject links, and report content
semantics.

## Requirements

### Data-link text

- Each generated Evaluation v2 `Data`-column link **MUST** use the base filename
  of its linked payload as the link text — for example
  `area-analysis-result.json` rather than `analysis`.

- The Area report `Data` link **MUST** read `area-analysis-result.json`; the
  Factor report `Data` link **MUST** read `factor-analysis-result.json`; and the
  Requirement report `Data` cell **MUST** list
  `requirement-assessment-result.json` and `requirement-rating-result.json`.

- Generated data-link *targets* **MUST NOT** change; only the rendered link text
  changes.

  > Rationale: the filename is the one piece of information the generic word
  > omitted, and it matches the tab a reader sees after following the link.
  >
  > - 0109

### Boundaries

- The `Data` column header, navigation trails, and subject links **MUST NOT**
  change as part of this work.

- Structured routine data paths under `data/` **MUST NOT** change.

## Acceptance Criteria

- An Area report's `Data` cell renders a link whose text is
  `area-analysis-result.json` and whose target is the existing area analysis
  payload path.
- A Factor report's `Data` cell renders a link whose text is
  `factor-analysis-result.json`.
- A Requirement report's `Data` cell renders two links whose texts are
  `requirement-assessment-result.json` and `requirement-rating-result.json`.
- Every data-link target resolves to the same payload file it did before this
  change.
- Tests assert the new Area data-link text and the Requirement two-link cell.

## Durable spec changes

### To add

None.

### To modify

- [`specs/evaluation-v2/reports/report-tree.md`](../../../specs/evaluation-v2/reports/report-tree.md)
  - require generated `Data`-column links to use the payload base filename as
    link text. Driven by [Data-link text](#data-link-text).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, each Evaluation v2 `Data`-column link
names the payload file it opens, while link targets, data paths, and all other
report content stay exactly as they are. That achieves the motivation without
widening into target or layout changes.
