---
type: Functional Specification
title: Evaluation v2 report header navigation - functional spec
description: Requirements for labeled navigation trails and compact headers in Evaluation v2 Markdown reports.
tags: [cli, evaluation, reports, navigation]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation v2 report header navigation - functional spec

Companion to the
[Evaluation v2 report header navigation](../0104-evaluation-v2-report-header-navigation.md)
change case. This spec states the delta for Evaluation v2 Markdown report
headers. It defers the cumulative Evaluation v2 report contract to
[`specs/evaluation-v2/reports/report-tree.md`](../../specs/evaluation-v2/reports/report-tree.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Evaluation v2 Markdown reports are generated for readers who need to move
between Area, Factor, and Requirement detail pages while keeping their place in
the Model. The current `Breadcrumb:` and `Parent Area:` / `Parent Factor:`
headers are navigable, but the labels are generic and the parent links repeat
navigation already present in the trail. A clearer header labels the actual
model hierarchy, then uses the report title and a compact summary to show the
current subject's state.

## Scope

Covered:

- Evaluation v2 Markdown report headers for Area, Factor, and Requirement pages.
- Labeled Area and Factor navigation trails.
- Removal of redundant parent-link header lines.
- Compact report-specific summary tables.
- Preservation of existing required identity, rating/status, confidence, data,
  Area link, Factor link, and Requirement attachment information.
- Focused tests and durable report-tree spec sync.

Deferred / non-goals:

- No changes to report tree paths, output refs, routine JSON payloads, rating
  computation, display-title catalogs, or report-build receipt JSON.
- No new visual renderer beyond Markdown.
- No table-of-contents, previous/next navigation, or backlinks beyond the
  existing report tables.

## Header Shape

Area reports should start in this shape:

```md
Area: [Root](../../report.md) / [Payments](../report.md)

# Checkout

Path: `/payments/checkout`

| Overall   | Local      | Confidence    | Data            |
| --------- | ---------- | ------------- | --------------- |
| 🔵 Target | 🟡 Minimum | Medium / High | [analysis](...) |
```

Factor reports should start in this shape:

```md
Area: [Root](../../report.md) / [Payments](../report.md)

Factor: [Reliability](../report.md) / [Latency](./report.md)

# Latency

Path: `payments::reliability/latency`

| Overall    | Local     | Status              | Confidence    | Data            |
| ---------- | --------- | ------------------- | ------------- | --------------- |
| 🟡 Minimum | 🔵 Target | Analyzed / Analyzed | Medium / High | [analysis](...) |
```

Requirement reports should start in this shape:

```md
Area: [Root](../../report.md) / [Payments](../report.md)

# Checkout handles retry failures

Name: `retry-failures`

| Rating     | Assessment | Factors                               | Confidence    | Data                             |
| ---------- | ---------- | ------------------------------------- | ------------- | -------------------------------- |
| 🟡 Minimum | Assessed   | [Reliability](...), [Resilience](...) | Medium / High | [assessment](...), [rating](...) |
```

The example links and Rating Level titles are illustrative. Renderers continue to
resolve links and labels from the run's actual report tree and model snapshot.

## Requirements

- Every Evaluation v2 Markdown report **MUST** start with an `Area:` navigation
  trail whose elements link to the corresponding Area reports from the root Area
  through the current Area report or owning Area report.

  > Rationale: `Area` names the model hierarchy the trail represents, while a
  > generic `Breadcrumb` label describes a UI widget rather than the report's
  > subject context. - 0104

- Factor reports **MUST** include a `Factor:` navigation trail after the `Area:`
  trail. The trail **MUST** link each Factor ancestor and the current Factor to
  its generated Factor report.

  > Rationale: Factors form their own hierarchy inside an Area; showing that
  > hierarchy as a labeled trail is clearer than a single parent link. - 0104

- Requirement reports **MUST NOT** render a `Factor:` breadcrumb or choose one
  attached Factor as a navigation parent.

  > Rationale: a Requirement can attach to multiple Factors, so a singular
  > Factor trail would imply ownership the Model does not guarantee. - 0104

- Generated Evaluation v2 Markdown report headers **MUST NOT** render standalone
  `Breadcrumb:`, `Parent Area:`, `Parent Factor:`, or `Parent:` lines.

  > Rationale: labeled trails already provide upward navigation; repeated parent
  > lines add noise at the highest-attention part of the report. - 0104

- Area report headers **MUST** show the report title, Area display path, overall
  rating, local rating, confidence pair, and Area analysis data link before the
  summary section.

- Factor report headers **MUST** show the report title, Factor display path,
  local-and-descendant rating labeled as `Overall`, local rating, status pair,
  confidence pair, and Factor analysis data link before the summary section.

  > Rationale: `Overall` is the report-reader label for the aggregate value; the
  > durable spec can define that it corresponds to
  > `localAndDescendantAnalysis`. - 0104

- Requirement report headers **MUST** show the report title, stable Requirement
  name, selected rating or unrated status, assessment status, linked attached
  Factors, confidence pair, and assessment/rating data links before the summary
  section.

- Header summary tables **SHOULD** use report-specific columns instead of a
  generic `Field | Value` key-value table.

  > Rationale: the report title and identity line already name the subject; a
  > one-row summary surfaces state without repeating `Area`, `Factor`, or
  > `Requirement` as table fields. - 0104

- Report tests **MUST** cover the new Area, Factor, and Requirement header
  shapes, including linked trails and the absence of old parent-link labels.

## Durable spec changes

### To add

None

### To modify

- [`specs/evaluation-v2/reports/report-tree.md`](../../specs/evaluation-v2/reports/report-tree.md)
  - replace the generic breadcrumb and parent-link navigation contract with
    labeled Area/Factor trails and compact report-header requirements.

### To rename

None

### To delete

None
