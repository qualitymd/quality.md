---
type: Change Case
title: Model query commands
description: Add a read-only `qualitymd model` command group — `tree`, `list`, `get` — that projects a quality model's structure and canonical reference IDs, and wire the `/quality` evaluate workflow to query those IDs instead of hand-deriving them; bounded so `model` does not overlap the `status` (state/readiness) or `evaluation` (runs/results) surfaces.
status: Done
tags: [cli, model, query]
timestamp: 2026-06-26T00:00:00Z
---

# Model query commands

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0125-model-query-commands/spec.md) - what the change must do.
- [Design doc](0125-model-query-commands/design.md) - how the code delivers it
  (the projection's home and the reference-grammar move).

## Motivation

Authoring evaluation payloads requires the canonical reference IDs of a model's
elements — `area:<path>`, `factor:<area>::<path>`,
`requirement:<area>::<name>`. No command emits them today: `status --json`
reports `sourceCoverage` by area _path_ and label, not canonical IDs, and does
not enumerate factors or requirements. So an agent hand-derives tens of IDs from
`QUALITY.md` — slow, and a silent-typo source that the 0124 reference-`kind`
enum and the existing identity-resolution checks can only catch _after_ the fact.

The acquire-roi-next evaluation made this concrete: the orchestrator manually
derived every area/factor/requirement ID it needed for ~115 payloads. A
read-only projection of the model — its elements, their IDs, and their
containment — turns that into a query, for an agent (`--json`) and a person
(tree) alike.

## Scope

Covered: a read-only `model` noun with three verbs — `tree` (hierarchical view),
`list` (flat, filterable enumeration), `get <id>` (single-element detail) — over
a model file, each with a human default form and `--json`, all emitting canonical
reference IDs. Also covered: wiring the `/quality` evaluate workflow to consume
those IDs (querying the run's `model-snapshot.md`) instead of hand-deriving them
— the originating pain, and a required outcome of this case.

This case draws a hard boundary so `model` does not overlap existing surfaces:

- **`model`** owns _logical structure and identity_: elements, canonical IDs,
  labels, and containment relationships.
- **`status`** keeps _state & readiness_: lint summary, evaluation history, next
  actions, source coverage/provenance, and the aggregate `shape` counts.
- **`evaluation`** keeps _runs, payloads, reports, and the `model-snapshot.md`
  artifact_.

Deferred / non-goals:

- No selector/query expression language (`model query …`) — deferred until a real
  need.
- No mutation — `model` never writes; authoring stays in editing `QUALITY.md`.
- No model validation — diagnostics stay in `lint`; `model` assumes a parseable
  model.
- No source/provenance, coverage, or aggregate count summaries — those stay in
  `status`.
- No evaluation-run awareness — no `--run` flag, no run-folder resolution; `model`
  reads a model _file_ (a `model-snapshot.md` is reachable only as a file path).
- No ratings or evaluation results.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0125-model-query-commands/spec.md#durable-spec-changes)
section. The index below is the full skimmable list to reconcile before
In-Review.

### Code

- [x] `internal/cli/model.go` (new) - the `model` parent command and its `tree` /
      `list` / `get` subcommands; registered in `internal/cli/root.go` under the
      common command group (alongside `lint`/`status`/`spec`/`schema`).
- [x] `internal/model/` - a read-only projection (`internal/model/projection.go`)
      that walks a parsed model and emits each element's canonical reference ID,
      label, kind, and containment. The duplicate model-tree walk in
      `internal/status` (`modelAccumulator.walkAreas`/`walkFactors`) is retired:
      `status` now derives its shape counts from the shared projection and keeps
      only a source-aware coverage walk.
- [x] Canonical reference encoding - the path types and builders
      (`AreaPath`/`FactorPath`, `AreaPath.Reference()` / `FactorReference()` /
      `RequirementReference()`) plus the parsers and existence helpers moved from
      `internal/evaluation` into a neutral `internal/model/reference.go` shared by
      both. Rating references stay in `evaluation` (they address the rating scale,
      not model structure). `model` depends on neither `status` nor `evaluation`.
- [x] Tests - `internal/cli/model_test.go` (new) for the three verbs, filters,
      `--json` shapes, canonical-ID round-trips, and the unknown-id usage error;
      `internal/model/reference_test.go` and `internal/model/projection_test.go`.
      No root-command help golden enumerates subcommands, so none needed updating.

### Format spec

- [x] [`SPECIFICATION.md`](../../SPECIFICATION.md) - reviewed; no change. It already
      defines the canonical reference grammar `model` emits.

### Durable specs

- [x] `specs/cli/model.md` (new) - the durable command spec for the `model` group.
- [x] `specs/cli.md` - add `model` to the command table and the `Commands` list.
- [x] `specs/cli/index.md` - register `model.md`.

### Durable docs / bundled skill

- [x] `skills/quality/workflows/evaluate.md` - wired the evaluate workflow to
      fetch in-scope canonical IDs from `model list --json` (against the run's
      `model-snapshot.md`) instead of hand-deriving them from `QUALITY.md`
      (steps 8–9). This is the originating pain (see Motivation). See the spec's
      [Skill integration](0125-model-query-commands/spec.md#skill-integration-evaluate-workflow)
      requirements (SK1–SK4).
- [x] `skills/quality/resources/cli-quick-reference.md` - added the `model` verbs
      to the quick-reference table, the introspection list, and the evaluating
      decision tree.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` - added the normative
      "Canonical references" section requiring payload references come from
      `model`, not hand-derived IDs.
- [x] `README.md` - enumerates commands, so added a `qualitymd model tree` row.

### Suggested new durable specs

- The model projection (elements + canonical IDs + containment) has no durable
  data-contract spec of its own; if the `--json` shape grows, a 1:1 artifact spec
  for it may be worth lifting out. Suggesting only.

## Status

`Done`. The `model` group (`tree`/`list`/`get`) ships with the reference grammar
and projection moved into `internal/model`, `status` folded onto the shared
projection, and the evaluate workflow wired to query canonical IDs from the run's
`model-snapshot.md`. Durable specs (`specs/cli/model.md`, `specs/cli.md`,
`specs/cli/index.md`, the evaluate skill spec) and docs (the bundled skill,
README) are in sync. Full build and test suite green.
