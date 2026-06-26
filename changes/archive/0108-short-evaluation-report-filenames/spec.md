---
type: Functional Specification
title: Short evaluation report filenames - functional spec
description: Requirements for preserving the root Evaluation v2 report.md entrypoint while replacing repeated descendant report.md files with short subject-aware generated Markdown report filenames.
tags: [cli, evaluation, reports]
timestamp: 2026-06-26T00:00:00Z
---

# Short evaluation report filenames - functional spec

Companion to the
[Short evaluation report filenames](../0108-short-evaluation-report-filenames.md)
change case. This spec states what the filename change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

Evaluation v2 report paths are structurally deterministic, but every human
Markdown report currently ends in `report.md`. That is clear inside the file
tree and weak in editor or browser tabs, where reviewers often see only the
filename. The root Area report can remain the single run entrypoint at
`report.md`; descendant subject-aware filenames improve tab scanning without
moving identity out of the existing directory structure.

## Scope

Covered: generated Evaluation v2 human Markdown report paths, report links,
navigation link targets, report refs, report output indexes, and durable
contracts that describe those paths.

Not covered: structured routine data paths, report directory names, report
content semantics, compatibility copies, display-title-based paths, existing run
migration, and historical archived Change Cases.

## Requirements

### Filename contract

- The Evaluation v2 report builder **MUST** write the root Area report to
  `report.md`.

- For each non-root Area, the Evaluation v2 report builder **MUST** write the
  Area report to `areas/<area>/<area>-area.md`, where `<area>` is the Area's
  local structural Area ID segment.

- For each Requirement, the Evaluation v2 report builder **MUST** write the
  Requirement report under the owning Area report folder at
  `requirements/<requirement>/<requirement>-requirement.md`, where
  `<requirement>` is the structural Requirement name.

- For each Factor, the Evaluation v2 report builder **MUST** write the Factor
  report under the owning Area report folder at
  `factors/<factor>/<factor>-factor.md`, where `<factor>` is the Factor's local
  structural Factor path segment.

- Nested Factor reports **MUST** keep recursing through nested `factors/`
  folders and **MUST** use the nested Factor's local segment for the filename.
  For example:

  ```text
  factors/reliability/factors/latency/latency-factor.md
  ```

  > Rationale: the directory path already carries full structural identity, while
  > the filename only needs to make open tabs distinguishable. Full-path
  > filenames would solve tab ambiguity at the cost of much longer labels. -
  > 0108

- The Evaluation v2 report builder **MUST NOT** generate non-root Area, Factor,
  or Requirement Markdown reports named `report.md`.

- The Evaluation v2 report builder **MUST NOT** derive persisted report
  filenames from display titles, natural labels, or rendered human labels.

### Links and indexes

- Generated Markdown links between reports **MUST** target the new filenames
  from every source report path.

- Report navigation trails **MUST** preserve their current labels and hierarchy;
  only the link targets change.

- Report tables **MUST** continue to link row subjects to generated human report
  targets when the row has exactly one generated human report target, using the
  new filenames.

- Structured report refs and report output indexes that point to generated
  human Markdown reports **MUST** use the new report paths.

- Structured routine data locations under `data/` **MUST NOT** change as part of
  this filename change.

### Compatibility boundaries

- The report builder **MUST NOT** write duplicate compatibility copies using old
  descendant `report.md` filenames.

- Existing completed evaluation runs **MUST NOT** be migrated or rewritten by
  this change.

- `qualitymd evaluation status` **MUST NOT** read generated Markdown report
  bodies to compute status, regardless of the generated report filenames.

## Acceptance Criteria

- Building an Evaluation v2 report creates only one generated human Markdown
  file named `report.md`: the root Area report at the run root.
- The root Area report exists at `report.md`.
- A non-root Area report exists at
  `areas/<area>/<area>-area.md`.
- A root-local Factor report exists at
  `factors/<factor>/<factor>-factor.md`.
- A nested Factor report exists at
  `factors/<factor>/factors/<child-factor>/<child-factor>-factor.md`.
- A Requirement report exists at
  `requirements/<requirement>/<requirement>-requirement.md`.
- Generated Area, Factor, and Requirement report links resolve to existing files
  using the new filenames.
- Generated navigation labels and table subject labels remain unchanged except
  for their link targets.
- `data/evaluation-output-result.json` report refs and report outputs use
  `report.md` for the root Area report and the new generated Markdown paths for
  descendant reports.
- Structured routine data paths under `data/` are unchanged.
- Existing tests that assert `report.md` paths are updated or replaced with
  assertions for the new filenames.

## Durable spec changes

### To add

None.

### To modify

- [`specs/evaluation-v2/records/data-layout.md`](../../../specs/evaluation-v2/records/data-layout.md)
  - update the Report Tree and Path Derivation contract for the new generated
    Markdown filenames. Driven by
    [Filename contract](#filename-contract).
- [`specs/evaluation-v2/reports/report-tree.md`](../../../specs/evaluation-v2/reports/report-tree.md)
  - require generated report navigation and subject links to target the new
    filenames while preserving report labels and hierarchy. Driven by
    [Links and indexes](#links-and-indexes).
- [`specs/skills/quality-skill/reporting.md`](../../../specs/skills/quality-skill/reporting.md)
  - align the `/quality` reporting artifact contract with the new generated
    report filenames. Driven by
    [Filename contract](#filename-contract).
- [`specs/cli/status.md`](../../../specs/cli/status.md) - replace the stale
  `report.md`-specific body-reading rule with a filename-independent generated
  Markdown report rule. Driven by
  [Compatibility boundaries](#compatibility-boundaries).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, generated Evaluation v2 report tabs
will distinguish Area, Factor, and Requirement reports by short local filenames,
while the directory tree continues to carry full structural identity and existing
structured data remains stable. That achieves the motivation without widening
the change into report-content redesign or run migration.
