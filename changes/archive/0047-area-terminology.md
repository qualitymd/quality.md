---
type: Change Case
title: Area terminology changeover
description: Replace the formal Target terminology with Area across the QUALITY.md schema, records, reports, CLI, skill, scaffold, examples, and docs.
status: Done
tags: [terminology, schema, evaluation, cli, skill]
timestamp: 2026-06-21T00:00:00Z
---

# Area terminology changeover

This Change Case records the full vocabulary change from **Target** to
**Area**. Detail lives in:

- [Functional spec](0047-area-terminology/spec.md) - what the change must do.
- [Design doc](0047-area-terminology/design.md) - how to implement the
  changeover.

## Motivation

`Target` is overloaded in QUALITY.md: it names the recursive evaluated node, and
`Target` is also the default rating level title. Reports with a `Target` column
and `Target` rating make the model feel harder to read than it is. `Area` is
shorter, friendlier, and works for projects, documents, code, services, test
suites, and other evaluatable parts without implying only software architecture
components.

The new vocabulary should be simple: every evaluatable node is an **Area**, and
the document's top-level Area is the **root area**. Do not introduce a second
canonical user-facing noun for the same concept.

## Scope

Covered: a full, no-backward-compatibility changeover from `Target` to `Area`
for the model schema, evaluation record schema, report model, report Markdown,
CLI command/help output, lint/status behavior, evaluation run naming, bundled
`/quality` skill instructions and resources, scaffold comments, examples,
fixtures, generated npm README, and project docs.

Deferred / non-goals: preserving `targets:` or `targetPath` as accepted legacy
inputs; migrating old evaluation runs; changing the default rating level
`target` / `Target`; renaming unrelated English uses of "target" that do not
refer to the model node concept, such as "target range" or "target level".

## Affected Artifacts

### Code

- `internal/model/` - rename the recursive model type and parsed field from
  Target/`targets` to Area/`areas`.
- `internal/schema/` - rename schema node/property names and generated schema
  text for Areas.
- `internal/lint/` - validate `areas:` and Area semantics; reject legacy
  `targets:` as unknown/non-conforming rather than accepting it.
- `internal/status/` - report model shape, narrowing, and evaluation history
  using Area terminology.
- `internal/evaluation/` - replace TargetPath/`targetPath`, target summaries,
  target details, analysis target fields, and report display labels with
  AreaPath/`areaPath` equivalents.
- `internal/cli/` and `cmd/qualitymd/` - update command help, diagnostics,
  examples, and receipts that name targets, target paths, subjects, or
  subject-altitude runs.
- Tests under `internal/**` - update model fixtures, record fixtures, expected
  report output, and incompatible-record expectations for the new schema.

### Durable Specs

- `SPECIFICATION.md` - make Area the formal recursive model node and root area
  the formal root descriptor.
- `specs/evaluation-records.md` - replace `targetPath`, target analysis, target
  summaries/details, report sections, and generated artifact wording with Area
  equivalents.
- `specs/cli/` and `specs/cli.md` - update CLI contracts that expose target
  terminology in arguments, diagnostics, examples, output, run naming, or the
  model-file selection flag.
- `specs/skills/quality-skill/quality-skill.md` and related
  `specs/skills/quality-skill/` examples/guides - align the skill contract and
  example evaluation bundle with Area terminology.

### Durable Docs

- `README.md` and `npm/quality.md/README.md` - update usage examples, model
  schema snippets, and terminology tables.
- `install.md` - update any model-discovery or evaluation examples that name
  targets.
- `docs/guides/` and `docs/reference/` - update authoring, CLI, release, and
  versioning prose where "target" refers to the model node concept.
- `CHANGELOG.md` - add the user-visible breaking change once implemented.

### Bundled Skill

- `skills/quality/SKILL.md` - route Area/factor narrowing, run-frame wording,
  and required terminology through Area instead of Target.
- `skills/quality/guides/` - update authoring, getting-started, top-10 checks,
  and mode guidance that describes target shape, target file, target/factor
  scope, or report interpretation.
- `skills/quality/resources/` - update CLI quick reference and any record JSON
  examples that use `targetPath`.
- `skills/quality/modes/` - update wizard, evaluate, improve, setup, and update
  guidance where model-node terminology appears.

### Scaffold, Examples, Fixtures

- `internal/scaffold/skeleton.md` - replace commented `targets:` guidance with
  `areas:` and root-area language.
- `QUALITY.md` - update this repo's dogfood model from `targets:` to `areas:`
  and update its Markdown body terminology.
- `quality/evaluations/` - leave historical run folders as historical data unless
  a fresh dogfood evaluation is intentionally recorded with the new schema.
- `specs/skills/quality-skill/examples/0001-quality-eval/` - update the
  maintained example bundle to the new Area record/report vocabulary.

## Status

`Done`. Implementation is complete, archived, and verified with `go test ./...`.

Implemented changes include the `areas:` schema, Area/`areaPath` evaluation
records, Area report JSON/Markdown, `--model` evaluation creation, new
`NNNN[-<narrowing>]-quality-eval` run names, Area lint/status surfaces, updated
durable specs and docs, aligned `/quality` skill guidance, updated scaffold and
dogfood model, and a regenerated maintained Sparrow example bundle.

Historical evaluation runs and append-only logs retain old terminology as
history; current readers do not translate legacy `targets:` or `targetPath`
records into the new model.
