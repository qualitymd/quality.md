---
type: Change Case
title: Required display titles
description: Require human-facing titles on model elements so reports, status output, and skill guidance can consistently render readable labels instead of internal identifiers.
status: Done
tags: [specification, schema, lint, report, skill]
timestamp: 2026-06-19T00:00:00Z
---

# Required display titles

A **Change Case** capturing the *why* and *status* for making display labels a
required part of the `QUALITY.md` model. The detail lives in its
[functional spec](0033-required-display-titles/spec.md).

> **Done.** Implementation and durable artifact synchronization are complete; the change is archived.

## Motivation

Evaluation reports, status snapshots, and `/quality` skill output are
human-facing. Today several model elements rely on map keys or rating identifiers
as their only stable display text, which makes output terse and sometimes
cryptic. Keys and `level` names are useful machine identifiers, but they are not
always the best labels for readers.

Requiring `title` on the model root, targets, factors, and rating levels gives
renderers and skills one consistent display property to use. The cost to authors
is low, and the resulting distinction is clearer: identifiers remain stable for
paths, references, and machine output; titles carry the reader-facing label.

## Scope

Covered: requiring non-empty scalar `title` values on the Model, every Target,
every Factor, and every Rating Level; adding `Factor.title`; updating lint so
missing titles are errors; updating scaffolds, examples, README, the authoring
guide, and skill/report/status consumers to use titles for display.

Deferred / non-goals: no `title` property on Requirements because the
requirement statement is already the display text; no title uniqueness rule in
this change; no change to target paths, factor references, rating-level ids,
evaluation roll-up, or report JSON identity fields.

Implementation is complete and archived.

## Affected specs & docs

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) - require `title` on Model,
      Target, Factor, and Rating Level; add `Factor.title`; define the
      identifier/title split; and update report semantics so human renderers use
      titles for display.
- [x] [`specs/cli/lint.md`](../../specs/cli/lint.md) - make `missing-title` an
      error rule covering Model, Target, Factor, and Rating Level titles, with
      context-specific finding messages.
- [x] [`specs/cli/init.md`](../../specs/cli/init.md) - require the scaffold to seed
      valid placeholder titles wherever the format now requires them.
- [x] [`specs/cli/evaluation-build-report.md`](../../specs/cli/evaluation-build-report.md)
      - specify title-based display for rendered target, factor, and rating
      labels while preserving stable ids/paths where needed.
- [x] [`specs/cli/status.md`](../../specs/cli/status.md) - specify that human labels
      in status/source coverage use required titles while JSON paths remain
      identifier-based.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      and
      [`specs/skills/quality-skill/guides/authoring.md`](../../specs/skills/quality-skill/guides/authoring.md)
      - update the skill contract and authoring-guide spec for required display
      titles and title-first human output.
- [x] [`README.md`](../../README.md) - update examples and schema reference.
- [x] [`skills/quality/guides/authoring.md`](../../skills/quality/guides/authoring.md)
      - update runtime authoring guidance.
- [x] [`internal/scaffold/skeleton.md`](../../internal/scaffold/skeleton.md) -
      include required placeholder titles in generated scaffolds.
- [x] [`QUALITY.md`](../../QUALITY.md) and example model fixtures under
      [`specs/`](../../specs/) / [`quality/`](../../quality/) as needed - update live
      and example models so they remain lint-valid under the new schema.

## Children

- [Functional spec](0033-required-display-titles/spec.md) - what the required
  display-title change must provide.

## Status

`Done`. The change is archived.
