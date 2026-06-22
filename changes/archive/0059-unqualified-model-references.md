---
type: Change Case
title: Unqualified model references
description: Define bounded unqualified Area, Factor, and Rating references and use them for Area-only report summary paths.
status: Done
tags: [format, references, reports, cli]
timestamp: 2026-06-22T00:00:00Z
---

# Unqualified model references

A **Change Case** defining unqualified model references as the context-bounded
counterpart to 0058's canonical typed references. The detail lives in its
[functional spec](0059-unqualified-model-references/spec.md) and
[design doc](0059-unqualified-model-references/design.md).

## Motivation

0058 introduced canonical typed model references such as
`area:operations/incident-response`. That form is correct when a reference must
stand alone or share a surface with other reference types. It is noisy, however,
inside an Area-only table column such as the report summary's Area Breakdown
`Path` column, where the column itself already fixes the type.

The format needs a formal, intentional allowance for omitting the typed prefix in
fixed-type contexts, without reopening ambiguity on mixed-reference surfaces or
durable machine artifacts. The implementation should also establish consistent
named helper functions for Area, Factor, and Rating references rather than
introducing boolean flags that obscure call-site intent.

## Scope

Covered:

- Define **unqualified reference** as a prefixless reference form allowed only
  when surrounding context fixes the reference type.
- Define unqualified Area, Factor, and Rating Level reference rendering.
- Allow tools to render or accept unqualified references at fixed-type human
  input/output edges.
- Forbid unqualified references on mixed-reference surfaces and anywhere the
  type must be recoverable from the value alone.
- Forbid persisting unqualified references in durable machine-readable
  artifacts.
- Add named render helpers for Area, Factor, and Rating references.
- Update the Area Breakdown `Path` column to render unqualified Area references.

Deferred / non-goals:

- No generic model-reference abstraction.
- No boolean `prefix` / `qualified` argument on reference rendering helpers.
- No new durable machine-readable string fields.
- No Area-local Factor shorthand that drops the declaring Area side; that can be
  designed later when a concrete Area-fixed Factor surface needs it.
- No broad CLI selector redesign beyond fixed-type reference parsing support.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0059-unqualified-model-references/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [x] `internal/evaluation/types.go` and model-reference helpers - add
      `UnqualifiedReference` methods/functions for Area, Factor, and Rating
      references while preserving qualified forms.
- [x] `internal/evaluation/model_reference.go` and tests - accept unqualified
      references only through type-specific parsing paths where context fixes
      the reference type.
- [x] `internal/evaluation/report.go` and report tests - render unqualified
      Area references in the shared Area Breakdown `Path` column.
- [x] Generated report fixtures under `specs/skills/quality-skill/examples/` -
      update report-summary and report Markdown outputs after renderer changes.

Specs:

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) - define unqualified references,
      allowed render/accept contexts, Area/Factor/Rating forms, and persistence
      limits.
- [x] [`specs/reports/report-summary-md.md`](../../specs/reports/report-summary-md.md)
      - allow the Area Breakdown `Path` column to use unqualified Area
      references.
- [x] [`specs/reports/report-md.md`](../../specs/reports/report-md.md) - keep the
      full report's shared compact Area Breakdown aligned with summary behavior.
- [x] [`specs/reports/report-json.md`](../../specs/reports/report-json.md) and
      [`specs/evaluation-records/report-outputs.md`](../../specs/evaluation-records/report-outputs.md)
      - preserve structured path arrays and forbid durable machine reliance on
      unqualified strings.
- [x] [`specs/skills/quality-skill/reporting.md`](../../specs/skills/quality-skill/reporting.md)
      - align skill-facing report guidance with unqualified Area Breakdown paths
      in human reports and structured identifiers in `report.json`.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      and [`specs/skills/quality-skill/evaluation.md`](../../specs/skills/quality-skill/evaluation.md)
      - align fixed-type input guidance with unqualified references.
- [x] OKF logs under [`specs/`](../../specs/log.md) - record durable spec updates.

Runtime skill and docs:

- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) and
      [`skills/quality/modes/evaluate.md`](../../skills/quality/modes/evaluate.md)
      - distinguish canonical typed references from unqualified fixed-type
      references.
- [x] [`skills/quality/resources/SPECIFICATION.md`](../../skills/quality/resources/SPECIFICATION.md)
      - update the bundled specification copy when `SPECIFICATION.md` changes.
- [x] [`CHANGELOG.md`](../../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

## Children

- [Functional spec](0059-unqualified-model-references/spec.md) - what
  unqualified references must allow, forbid, and render.
- [Design doc](0059-unqualified-model-references/design.md) - how the reference
  helpers, parsing boundary, and report rendering should be implemented.

## Status

`Done`. Landed and archived before the v0.9.0 release.
