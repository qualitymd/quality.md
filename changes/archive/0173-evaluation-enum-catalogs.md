---
type: Change Case
title: Evaluation Enum Catalogs
description: Centralize Evaluation enum values, labels, markers, and report ordering so validation, schemas, and Markdown rendering stay aligned.
status: Done
tags: [evaluation, schema, reports, enums]
timestamp: 2026-06-29T00:00:00Z
---

# Evaluation Enum Catalogs

Evaluation data already rejects many invalid enum values, but fixed
Evaluation vocabularies are split across raw string lists, typed constants,
display-label maps, and report sorting tables. This case makes those
vocabularies explicit typed catalogs so the CLI data contract, generated JSON
Schema, examples, report rendering, and tests use the same value sets and
display metadata.

- [Functional spec](0173-evaluation-enum-catalogs/spec.md) - what Evaluation
  enum catalogs, validation, schema output, display metadata, and report
  rendering must do.
- [Design doc](0173-evaluation-enum-catalogs/design.md) - how the Evaluation
  package centralizes enum metadata without changing persisted JSON values.

## Motivation

Evaluation reports are agent-authored and CLI-validated. When a vocabulary's
allowed values, labels, markers, and sort order live in different places,
validation can drift from report rendering or a new value can silently lack a
human label. The report model should keep stable machine values strict while
giving every known value one consistent human presentation.

## Scope

Covered:

- fixed Evaluation payload vocabularies used by routine data, report
  references, Advice data, finding data, statuses, confidence, impact, basis
  status, ranking tier, coverage disposition, and run gaps;
- JSON Schema enum generation and `qualitymd evaluation data set` validation;
- generated Markdown report labels, markers, and enum ordering;
- report renderer tests and schema/example tests;
- durable Evaluation data/report specs, skill evaluation guidance, release
  notes, and OKF logs.

Deferred:

- changing persisted enum values;
- accepting aliases, case-insensitive values, or legacy values;
- changing model-defined Rating Level IDs or Rating Level display titles;
- adding new public CLI flags or commands;
- changing historical Change Cases, archived logs, or released changelog
  entries.

## Affected artifacts

- Code:
  - `internal/evaluation/display.go`, `internal/evaluation/data_contract.go`,
    `internal/evaluation/data.go`, `internal/evaluation/report_tree.go`, and
    `internal/evaluation/types.go` - centralize fixed enum values and display
    metadata.
  - `internal/evaluation/evaluation_test.go` - cover catalog completeness,
    validation, schemas, and report rendering.
  - `internal/evaluation/evaluation-data.schema.json` - regenerate from the
    contract.
- Durable specs:
  - `specs/evaluation/records/payload-kinds.md` - name fixed vocabularies and
    their allowed values.
  - `specs/evaluation/records/json-conventions.md` - align report-reference
    kind and confidence vocabulary notes.
  - `specs/evaluation/reports/report-tree.md` - define report display metadata
    and ordering for fixed enum-like values.
  - `specs/cli/evaluation-data.md` - clarify enum-source and rejection behavior
    for data set/schema/example.
  - `specs/skills/quality-skill/evaluation.md` - align skill-authored
    Evaluation enum guidance.
- Bundled skill:
  - `skills/quality/SKILL.md` - align allowed enum references where the skill
    writes Evaluation data.
- Release notes:
  - `CHANGELOG.md` - Unreleased CLI note.
- Generated examples:
  - `examples/report-gallery/software-service/.quality/evaluations/0001-full-eval/`
    - regenerated report Markdown for enum display changes.
- OKF logs and indexes:
  - `changes/log.md`, `changes/index.md`, `changes/archive/index.md`, and this
    case.
  - `specs/evaluation/log.md` and `specs/cli/log.md`.

## Status

`Done`. Implemented and archived. Fixed Evaluation enum values now use typed
catalogs for validation, generated schemas, display labels, markers, and report
ordering; generated reports and durable guidance are aligned.
