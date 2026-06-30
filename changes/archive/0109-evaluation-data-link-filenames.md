---
type: Change Case
title: Filename text for evaluation data links
description: Render the Evaluation v2 report Data-column links with their payload filename text instead of the generic words analysis, assessment, and rating.
status: Done
tags: [cli, evaluation, reports, navigation]
timestamp: 2026-06-26T00:00:00Z
---

# Filename text for evaluation data links

This parent concept captures the why and status; the detail lives in its
child:

- [Functional spec](0109-evaluation-data-link-filenames/spec.md) - what the
  change must do.

No design doc: the change is a single, localized link-text swap with no design
choices to settle.

## Motivation

Evaluation v2 reports expose each subject's machine-readable payloads through a
`Data` column. The link targets already point at the structured
`*-result.json` files, but the link _text_ is a generic word — `analysis` on
Area and Factor reports, and `assessment`, `rating` on Requirement reports. A
reader cannot tell from the link which file they are about to open, and the
generic word restates information the surrounding report already conveys.

Following the same filename-clarity theme as
[0108](0108-short-evaluation-report-filenames.md), the link text should
name the payload file itself, so a reader knows exactly which `*-result.json`
they are opening and the rendered text matches the file the link resolves to.

## Scope

Covered:

- generated Evaluation v2 `Data`-column link text on Area, Factor, and
  Requirement reports; and
- the durable report-tree contract that governs those data links.

Deferred / non-goals:

- no change to data-link _targets_ or to structured routine data paths under
  `data/`;
- no change to the `Data` column header, navigation trails, or subject links;
  and
- no change to report content, ratings, findings, analysis, evidence, or limits.

## Affected artifacts

### Code

- [x] `internal/evaluation/report_v2.go` - render the three `Data`-column links
      with their payload base filename as link text via a `reportDataLink`
      helper.
- [x] `internal/evaluation/evaluation_test.go` - assert the new Area link text
      and add coverage for the Requirement two-link cell.

### Format spec

- [x] None - `SPECIFICATION.md` does not govern Evaluation v2 generated report
      link text. (Deliberate.)

### Durable specs

- [x] `specs/evaluation-v2/reports/report-tree.md` - fix the data-link text as
      the payload filename so the contract is intentional.

### Durable docs / bundled skill

- [x] None - no durable doc or bundled-skill content names the current
      `analysis`/`assessment`/`rating` link text. (Deliberate.)

### Suggested new durable specs

- None.

## Status

`Done`. Implemented, verified, and archived. See the
[status lifecycle](../index.md#status-lifecycle).
