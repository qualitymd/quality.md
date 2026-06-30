---
type: Change Case
title: Glossary and Report Links
description: Add a shared glossary, replace generated report legends with glossary-backed Evaluation links, and seed the glossary with core terms and vocabularies.
status: Done
tags: [evaluation, reports, glossary]
timestamp: 2026-06-30T00:00:00Z
---

# Glossary and Report Links

Generated Evaluation reports currently repeat local `Legend` blocks near
tables. Those keys keep notation close to first use, but they also duplicate
fixed vocabulary across artifacts and do not scale into durable definitions for
report readers. The project needs one glossary for shared report terms and
vocabularies, plus a compact cross-report navigation line that makes the
glossary and the key generated artifacts reachable from every report.

- [Functional spec](0183-glossary-report-links/spec.md) - what the glossary and
  report-link changes must provide.
- [Design doc](0183-glossary-report-links/design.md) - how the glossary,
  durable spec, report renderer, and generated examples absorb the change.

## Motivation

Report readers should be able to move between the overview, Findings,
Recommendations, and definitions without each report carrying its own partial
legend. A shared glossary can define fixed Evaluation vocabularies, model-defined
quality ratings, and core QUALITY.md terms once, while generated reports stay
compact and continue to render text labels directly in tables.

## Scope

Covered:

- a workspace-root `glossary.md` with flat alphabetical entries;
- core glossary terms for Area, Factor, Finding, Recommendation, and
  Requirement;
- a `Quality rating` glossary entry derived from this repository's configured
  `QUALITY.md` Rating Scale;
- fixed Evaluation enum catalog entries rendered as glossary tables;
- generated report `Evaluation links:` navigation to `report.md`,
  `findings.md`, `recommendations.md`, and `glossary.md`;
- removal of generated local `Legend` blocks from individual report artifacts;
- durable report specs, report design guidance, tests, generated examples, and
  related logs.

Deferred:

- generating `glossary.md` from Go enum catalogs;
- adding a CLI command to inspect glossary or enum catalog content;
- changing canonical Evaluation enum values;
- changing configured Rating Level IDs or Rating Scale semantics;
- adding every possible QUALITY.md concept beyond the initial glossary seed.

## Affected artifacts

- Code:
  - `internal/evaluation/report_tree.go` - remove local `Legend` blocks and
    render the standard `Evaluation links:` line with artifact-relative links.
  - `internal/evaluation/display.go` - continue to supply enum labels, markers,
    descriptions, and ordering used by the glossary source material.
  - `internal/evaluation/evaluation_test.go` - update report expectations for
    removed legends and added Evaluation links.
- Durable specs:
  - `specs/glossary-md.md` - new artifact spec for the workspace-root `glossary.md`
    structure, source material, and reader contract.
  - `specs/evaluation/reports/report-tree.md` - define the glossary-backed
    report navigation contract and remove local Legend requirements.
  - `specs/cli/evaluation-report.md` - keep report command expectations aligned
    with generated artifact navigation.
  - `specs/evaluation/records/payload-kinds.md` - remove the stale `unknown`
    Finding type from the valid Finding type list.
- Durable docs:
  - `glossary.md` - new shared glossary for report readers.
  - `docs/guides/reporting-design.md` - update report navigation and local-key
    guidance.
  - `README.md` - link the glossary if the docs entry point needs a user-facing
    reference.
  - `CHANGELOG.md` - add the user-facing report/glossary change.
- Generated examples:
  - `examples/report-gallery/` - regenerate generated Markdown reports.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, and this case.
  - `specs/index.md`, `specs/log.md`, `specs/evaluation/log.md`, and
    `docs/log.md`.

## Status

`Done`. Implemented and archived with the glossary artifact, durable specs,
generated report rendering, docs, tests, and generated examples aligned.
