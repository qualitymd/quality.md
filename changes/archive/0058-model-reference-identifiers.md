---
type: Change Case
title: Model reference identifiers
description: Define strict model names, canonical Area/Factor/Rating references, edge-only shorthand, and clearer report summary Area breakdown columns.
status: Done
tags: [format, references, reports, cli, lint]
timestamp: 2026-06-22T00:00:00Z
---

# Model reference identifiers

A **Change Case** defining how QUALITY.md names become stable model references
for Areas, Factors, and Rating Levels, and how those references appear in report
summary output. The detail lives in its
[functional spec](0058-model-reference-identifiers/spec.md).

## Motivation

QUALITY.md already distinguishes required human display titles from stable keys,
but the formal vocabulary is not sharp enough for commands, reports, and
evaluation records. Area map keys read like local names, while "Area ID" should
mean the path that uniquely identifies an Area within a Model. Factors have the
same issue with an extra wrinkle: a Factor path is only globally identifying
when paired with the declaring Area path.

The report summary exposes the same ambiguity in user-facing form. Its Area
Breakdown table currently uses a `Path` column that mixes title-bearing display
paths and a column named `Area` that actually contains the Area-only rating. The
summary should instead separate the Area title, stable Area reference, Area-only
rating, aggregate rating, and compact Factor ratings.

## Scope

Covered:

- Define strict Area names, Factor names, and Rating Level IDs.
- Define Area IDs, Factor IDs, and Rating Level IDs within one Model.
- Define canonical typed model-reference syntax for Areas, Factors, and Rating
  Levels.
- Allow tools to support shorthand model references only at human/input edges
  where the expected reference type is already fixed.
- Keep durable machine artifacts structured or canonical, never shorthand.
- Update `report-summary.md` and the shared Area Breakdown renderer to separate
  Area title from stable Area reference/path.

Deferred / non-goals:

- No relaxed identifier grammar yet.
- No new global namespace for Area or Factor names.
- No change to Requirement statement keys, which remain natural-language map
  keys.
- No change to `report.json` replacing structured `areaPath` or `factorPath`
  arrays with string references.
- No new query language beyond canonical model references.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0058-model-reference-identifiers/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [x] `internal/schema/` and companion JSON Schema generation - add strict name
      patterns where the structural schema can express them.
- [x] `internal/lint/` - enforce Area names, Factor names, Rating Level IDs, and
      related diagnostics.
- [x] `internal/evaluation/` - add canonical model-reference rendering/parsing
      helpers as needed while preserving structured `areaPath` and `factorPath`
      machine fields.
- [x] `internal/evaluation/report.go` and report tests - update the shared Area
      Breakdown table columns and root/nested Area rows.
- [x] `internal/cli/` - align help, scope parsing, and any typed selector
      surfaces with model-reference terminology and edge-only shorthand. No
      direct CLI selector surface exists in this phase; runtime and durable
      `/quality` scope guidance carry the user-facing selector behavior.
- [x] `internal/scaffold/` - update starter `QUALITY.md` placeholders so emitted
      Factor names satisfy the strict name grammar.

Specs:

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) - define name grammar, Area ID,
      Factor ID, Rating Level ID, canonical model references, and shorthand
      rules.
- [x] [`specs/reports/report-summary-md.md`](../../specs/reports/report-summary-md.md) - update the concise Area Breakdown table contract.
- [x] [`specs/reports/report-md.md`](../../specs/reports/report-md.md) - keep the
      shared compact Area Breakdown table aligned with the summary.
- [x] [`specs/reports/report-json.md`](../../specs/reports/report-json.md) -
      clarify that JSON preserves structured IDs and may add canonical
      references only without replacing arrays.
- [x] [`specs/evaluation-records/report-outputs.md`](../../specs/evaluation-records/report-outputs.md) - align shared report-model terminology with canonical model references.
- [x] [`specs/cli/evaluation-create.md`](../../specs/cli/evaluation-create.md),
      [`specs/cli/evaluation-report.md`](../../specs/cli/evaluation-report.md),
      and other CLI selector specs as discovered - align user-facing selectors
      with model-reference and shorthand rules. No direct CLI selector specs
      needed changes beyond lint/report contracts in this phase.
- [x] [`specs/cli/lint.md`](../../specs/cli/lint.md) - state that lint enforces
      the strict model-name grammar as part of model validation.
- [x] [`specs/cli/lint-rules.md`](../../specs/cli/lint-rules.md) - add or update
      diagnostics for invalid Area names, Factor names, and Rating Level IDs,
      including concrete rule IDs.
- [x] [`specs/quality-schema-json.md`](../../specs/quality-schema-json.md) -
      include the structural expression of the strict name grammar if the
      generated schema supports it.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md),
      [`specs/skills/quality-skill/evaluation.md`](../../specs/skills/quality-skill/evaluation.md),
      and [`specs/skills/quality-skill/guides/authoring-md.md`](../../specs/skills/quality-skill/guides/authoring-md.md) - align durable skill contracts and authoring-guide contract with model
      references and strict names.
- [x] OKF logs under [`specs/`](../../specs/log.md) - record durable spec updates.

Runtime skill and docs:

- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - align scope
      resolution and stable identifier language with model references.
- [x] [`skills/quality/resources/SPECIFICATION.md`](../../skills/quality/resources/SPECIFICATION.md) - update the bundled specification copy when `SPECIFICATION.md` changes.
- [x] [`skills/quality/resources/cli-quick-reference.md`](../../skills/quality/resources/cli-quick-reference.md) - update CLI selector examples. No selector examples needed changes in
      this phase; record payload examples already use structured arrays.
- [x] [`skills/quality/modes/evaluate.md`](../../skills/quality/modes/evaluate.md)
      and related mode docs - use model-reference terminology for scoped
      evaluations.
- [x] [`skills/quality/guides/authoring.md`](../../skills/quality/guides/authoring.md) - teach strict names, IDs, and canonical model references in the
      self-contained authoring guide.
- [x] [`README.md`](../../README.md) - update examples if command examples adopt
      canonical model references.
- [x] [`CHANGELOG.md`](../../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

## Children

- [Functional spec](0058-model-reference-identifiers/spec.md) - what strict
  names, canonical model references, edge-only shorthand, and report summary
  columns must require.
- [Design doc](0058-model-reference-identifiers/design.md) - how strict names,
  canonical references, lint checks, schema support, and report rendering should
  be implemented.

## Status

`Done`. Implemented and archived. The change defines strict model-name grammar,
canonical Area/Factor/Rating references, edge-only shorthand boundaries, named
lint diagnostics, generated JSON Schema patterns, Area Breakdown report columns,
runtime and durable `/quality` guidance, scaffold updates, docs, and changelog.
Verified with `mise run check` and
`go run ./cmd/qualitymd lint --json QUALITY.md`.
