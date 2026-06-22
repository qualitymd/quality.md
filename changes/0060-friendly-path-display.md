---
type: Change Case
title: Friendly path display
description: Separate friendly display values from model-reference grammar so human reports can render the root Area path as / without changing references.
status: In-Review
tags: [format, reports, references, display]
timestamp: 2026-06-22T00:00:00Z
---

# Friendly path display

A **Change Case** defining the boundary between display values, qualified model
references, and unqualified model references. The detail lives in its
[functional spec](0060-friendly-path-display/spec.md) and
[design doc](0060-friendly-path-display/design.md).

## Motivation

0058 and 0059 introduced qualified and unqualified model references, including
`area:root` and `root` for the root Area. That reference grammar is useful for
tools and fixed-type inputs, but it is not always the friendliest report display
value. In human reports, the root Area path reads more naturally as `/`.

The implementation should keep reference grammar stable while giving report
rendering and fallback display text a clear, intentionally separate path. That
separation should apply consistently across Area paths, Factor paths, and Rating
Level IDs so future code does not use reference helpers as display helpers.

## Scope

Covered:

- Define display values as human-facing fallback/rendering values distinct from
  qualified and unqualified references.
- Keep qualified and unqualified reference grammar unchanged.
- Render the root Area path as `/` in human report display contexts.
- Keep nested Area display paths as slash-joined Area names.
- Apply the same display/reference separation pattern to Factor paths and Rating
  Level IDs, even when non-root display output currently matches the
  unqualified reference form.
- Keep durable machine artifacts structured and unchanged.

Deferred / non-goals:

- No change to qualified or unqualified model-reference parsing.
- No acceptance of `/` as an Area reference.
- No title-aware display inside `AreaPath`, `FactorPath`, or Rating Level helper
  functions; report label resolvers continue to own title resolution.
- No new fields in `report.json` or evaluation records.
- No broad report redesign beyond friendly path display values.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0060-friendly-path-display/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

Code:

- [x] `internal/evaluation/types.go` - separate `Display` behavior from
      qualified and unqualified reference grammar for Area paths, Factor paths,
      and Rating Level helper functions.
- [x] `internal/evaluation/report.go` and report tests - render friendly Area
      display values in human report `Path` fields.
- [x] Generated report fixtures under `specs/skills/quality-skill/examples/` -
      update report-summary and report Markdown outputs after renderer changes.

Specs:

- [x] [`SPECIFICATION.md`](../SPECIFICATION.md) - distinguish display values
      from model references and keep reference grammar unchanged.
- [x] [`specs/reports/report-summary-md.md`](../specs/reports/report-summary-md.md)
      - require or allow friendly Area path display for the Area Breakdown
      `Path` column, including `/` for root.
- [x] [`specs/reports/report-md.md`](../specs/reports/report-md.md) - keep full
      report Area Breakdown and Area Detail path display aligned with summary
      behavior.
- [x] [`specs/reports/report-json.md`](../specs/reports/report-json.md) and
      [`specs/evaluation-records/report-outputs.md`](../specs/evaluation-records/report-outputs.md)
      - preserve structured path arrays and avoid display string persistence in
      machine artifacts.
- [x] [`specs/skills/quality-skill/reporting.md`](../specs/skills/quality-skill/reporting.md)
      and [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
      - align durable skill guidance with display values versus references.
- [x] OKF logs under [`specs/`](../specs/log.md) - record durable spec updates.

Runtime skill and docs:

- [x] [`skills/quality/SKILL.md`](../skills/quality/SKILL.md) and
      [`skills/quality/modes/evaluate.md`](../skills/quality/modes/evaluate.md)
      - distinguish human display values from model references where reports are
      discussed.
- [x] [`skills/quality/resources/SPECIFICATION.md`](../skills/quality/resources/SPECIFICATION.md)
      - update the bundled specification copy when `SPECIFICATION.md` changes.
- [x] [`CHANGELOG.md`](../CHANGELOG.md) - add the unreleased entry when the
      implementation lands.

## Children

- [Functional spec](0060-friendly-path-display/spec.md) - what display/reference
  separation and friendly report path rendering must require.
- [Design doc](0060-friendly-path-display/design.md) - how display helpers stay
  separate from model-reference helpers while reports render `/` for the root
  Area path.

## Status

`In-Review`. Implementation is complete across display/reference helper
separation, human Markdown report path rendering, durable specs, runtime skill
guidance, generated examples, and changelog. Verified with `mise run check`.
