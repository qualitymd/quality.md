---
type: Functional Specification
title: Report Body Rating Drivers - functional spec
description: Requirements for removing standalone rating-driver sections from generated Markdown report bodies.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-29T00:00:00Z
---

# Report Body Rating Drivers - functional spec

Companion to
[Report Body Rating Drivers](../0160-report-body-rating-drivers.md). This spec
states the delta for generated Evaluation Markdown reports. The durable source
of truth is absorbed into
[`Evaluation report tree`](../../../specs/evaluation/reports/report-tree.md) and
[`/quality reporting`](../../../specs/skills/quality-skill/reporting.md).

The key words **MUST**, **MUST NOT**, and **SHOULD** are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Rating drivers are a structured audit trail for ratings. They are still needed
in Evaluation data because the CLI can validate that rated outputs cite persisted
routine outputs, and agents can inspect those payloads when they need the full
trace. Generated Markdown reports have a different job: help a reader scan the
result, evidence, recommendations, coverage, and limits quickly.

Standalone `## Rating Drivers` sections make the report body more mechanical
without adding much reader value. Their compact JSON `inputRefs` expose payload
shape rather than report insight. The report summary and subject tables should
carry the human-facing roll-up explanation, while source-data frontmatter and
structured payload links keep traceability available.

## Scope

This change covers generated Markdown report bodies for:

- `report.md`;
- root and non-root Area reports;
- Factor reports.

Deferred:

- Evaluation data schema or validation changes;
- removing or weakening structured `ratingDrivers` requirements;
- replacing driver tables with clickable links, collapsible sections, or another
  visible body representation;
- Requirement report Finding sections, findings index, recommendations index,
  and recommendation detail reports.

## Requirements

1. The run-level `report.md` **MUST NOT** render a standalone
   `## Rating Drivers` section.

   > Rationale: The run report is the primary decision surface; summary, Top
   > Findings, Top Recommendations, scope, subject reports, coverage, and limits
   > should be the body sections readers scan first.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` to remove
   > rating drivers from the run report body contract.

2. Area reports **MUST NOT** render standalone `## Rating Drivers` sections.

   > Rationale: Area reports already show the roll-up summary, ratings, child
   > Areas, Factors, local Requirements, limits, and incomplete inputs. The
   > structured Area Analysis Result remains available through report
   > frontmatter `data`.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md`.

3. Factor reports **MUST NOT** render standalone `## Rating Drivers` sections.

   > Rationale: Factor reports should foreground the Factor summary,
   > Requirement ratings, Sub-Factor ratings, limits, and incomplete inputs.
   > The structured Factor Analysis Result remains available through report
   > frontmatter `data`.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md`.

4. Report generation **MUST** preserve structured `ratingDrivers` in persisted
   Evaluation data and source-data frontmatter links.

   > Rationale: This change only removes a noisy Markdown projection. It must not
   > weaken the rating audit trail or the CLI's structured validation contract.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` to state that
   > rating drivers remain available through source payloads rather than
   > standalone body sections.

5. Generated report bodies **MUST** keep the human-facing explanation and action
   surfaces that remain useful without driver tables: summary, ratings,
   confidence/status, findings, recommendations, subject breakdown tables, scope,
   coverage, limits, and incomplete inputs where the report type already renders
   them.

   > Rationale: Removing driver tables should improve scannability without
   > reducing the decision-ready report content.
   >
   > Durable spec: modify `specs/evaluation/reports/report-tree.md` and
   > `specs/skills/quality-skill/reporting.md`.

## Acceptance criteria

- Generated run, Area, and Factor Markdown reports contain no
  `## Rating Drivers` heading and no `Driver | Effect | Inputs` table.
- Structured `ratingDrivers` payloads remain accepted, validated, and listed in
  the relevant source-data payloads.
- Run reports still render summary, Top Findings, Top Recommendations, scope,
  subject reports, coverage, limits, and incomplete inputs.
- Area reports still render summary, ratings, Factor tables, Child Area tables,
  Requirement tables, limits, and incomplete inputs.
- Factor reports still render summary, ratings, Requirement tables, Sub-Factor
  tables, limits, and incomplete inputs.
- Focused Go tests cover the absence of rating-driver body sections and the
  continued presence of neighboring report sections.
- Report gallery generation is deterministic and checked in.
- `mise run check` passes.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/reports/report-tree.md` - remove standalone rating-driver
  section requirements from generated Markdown body contracts while preserving
  structured rating-driver data availability (requirements 1-5).
- `specs/skills/quality-skill/reporting.md` - align skill reporting expectations
  with no standalone rating-driver body sections (requirements 2, 3, 5).

### To rename

None

### To delete

None
