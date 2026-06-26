---
type: Functional Specification
title: Evaluation v2 report tree
description: Deterministic Area, Factor, and Requirement Markdown reports for Evaluation v2.
tags: [evaluation, reports, markdown]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation v2 report tree

Evaluation v2 reports are deterministic Markdown projections over completed
structured routine outputs.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Source Of Truth

Report generation **MUST** consume `EvaluationOutputResult` and referenced
structured routine outputs.

Report generation **MUST NOT** inspect new source evidence.

Report generation **MUST NOT** introduce new findings, ratings, evidence, limits,
analysis, or recommendations.

## Navigation

Every report **MUST** start with linked breadcrumbs from the root Area to the
current report subject.

Every non-root report **MUST** include a parent link.

Area reports **MUST** link to local root Factor reports, local Requirement
reports, and direct child Area reports.

Factor reports **MUST** link to their owning Area report, parent Factor report
when present, child Factor reports, and direct Requirement reports.

Requirement reports **MUST** link to their owning Area report and every attached
Factor report.

## Area Reports

Area reports **MUST** include:

- Area title and path;
- overall and local ratings;
- overall and local confidence;
- data links;
- summary;
- rating drivers;
- local root Factors;
- direct child Areas;
- local Requirements; and
- limits and incomplete inputs.

## Factor Reports

Factor reports **MUST** include:

- owning Area link;
- Factor title and path;
- local and local-and-descendant ratings;
- local and local-and-descendant statuses;
- confidence;
- data links;
- summary;
- rating drivers;
- direct Requirements;
- direct child Factors; and
- limits and incomplete inputs.

## Requirement Reports

Requirement reports **MUST** include:

- owning Area link;
- Requirement title and name;
- Requirement rating status and selected rating when present;
- assessment status;
- attached Factor links;
- confidence;
- data links;
- summary;
- findings summary;
- finding detail sections; and
- unknowns and missing evidence.

## Rendering Rules

Reports **MUST** render empty tables with explicit empty-state rows.

Reports **MUST** render selected Rating Levels with the Rating Level `title`
resolved from the run's `model.md` snapshot, falling back to the stable Rating
Level ID only when a title is unavailable.

> Rationale: Markdown reports are the human review surface, and the model
> snapshot is the historical source for display vocabulary. Structured routine
> data and machine receipts keep stable Rating Level IDs. — 0102

Reports **MUST** render `not_assessed`, `not_rated`, `empty`, `not_analyzed`,
and `blocked` distinctly from Rating Level labels.

Reports **MUST** render CLI-owned enum-like report values, including statuses,
confidence levels, boolean values, report kinds, limits/incomplete-input types,
unknown/missing-evidence types, and known finding classifications, with
human-readable display titles in Markdown while preserving the raw values in
routine JSON, `EvaluationOutputResult`, and report-build receipts.

> Rationale: Markdown reports are optimized for human review and scanning, but
> agents and tools need stable values in the structured data. Unknown or
> free-form values should remain readable through fallback title-casing rather
> than turning presentation decoration into schema validation. — 0103

Reports **MUST** omit Rating Level values when the source result status says the
rating or scoped analysis was not produced.

Reports **MUST** preserve secret-handling boundaries. They may name the locator
and credential type but **MUST NOT** reproduce secret values or unsafe raw
content.

Ordering **MUST** be deterministic:

- Areas by structural Area ID, with the root Area first;
- Factors by declaring Area ID and structural Factor path;
- Requirements by declaring Area ID and Requirement name;
- rating drivers in source result order;
- findings in Requirement Assessment Result order; and
- evidence in recorded order.
