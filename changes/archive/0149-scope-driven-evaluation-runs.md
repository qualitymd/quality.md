---
type: Change Case
title: Scope-driven evaluation runs
description: Capture requested/planned scope as structured run data at create time and render the run report as the scoped area report, replacing the narrowing slug and the positional headline rule.
status: Done
tags: [evaluation, cli, skill, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Scope-driven evaluation runs

This parent concept captures the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0149-scope-driven-evaluation-runs/spec.md) — what the case must do.
- [Design doc](0149-scope-driven-evaluation-runs/design.md) — how it's built, and why.

## Motivation

An evaluation run's report selects a single "headline" subject, and that
selection is positional: `report build` reads `EvaluationFrame.inputs.factorIds`
and takes the **first listed Factor** that has an analysis. A full evaluation
populates `factorIds` with every Factor in model order, so the headline always
collapses to whatever Factor sits first in the model — Security, in the observed
case — and reports a 🟡 Minimum factor while the root Area's overall result is
🔴 Unacceptable. The report build is deterministic, but it derives the headline
from an agent-authored artifact whose ordering is incidental, so the headline is
effectively arbitrary.

The root cause is that scope is never captured as authoritative structured data.
`qualitymd evaluation create` records scope only as a `--narrowing` folder-name
slug (lossy, unparseable), and the structured scope that drives the headline is
reconstructed later by the agent into the frame. Naming, traversal, and headline
each read a different, hand-maintained source.

This case makes scope a first-class, CLI-owned run parameter captured at
`create`, and reshapes the run report so it simply _is_ the scoped Area's report.
The "headline" concept is removed: there is no subject to select, so there is no
selection to get wrong.

## Scope

Covered:

- A CLI-owned `RunManifest` record written by `create`, holding `requestedScope`
  (faithful) and `plannedScope` (normalized, root-defaulted).
- `create` flags `--area`/`--factor` replacing `--narrowing`, with snapshot
  validation and single-Area coherence; slug derived from `plannedScope`.
- The run report (`report.md`) rendered as the scoped Area report; the headline
  concept and its references removed.
- Scope removed from `EvaluationFrame`; `evaluation list` and the coverage check
  sourced from the manifest / derived from `plannedScope`.
- The `/quality` skill resolving requested scope to canonical IDs before
  `create` and passing it through the flags.

Deferred / non-goals:

- No back-compat shim for `--narrowing` or for runs created before this case
  lacking a `RunManifest` (early-alpha clean break).
- No change to evaluation judgment routines, rating semantics, or analysis
  payloads beyond removing scope from the frame.
- The natural-language ask is not persisted as structured data; it stays in the
  evaluate feedback log.

## Affected artifacts

Derived by sweeping `--narrowing`, `narrowing`, `headline`/`Headline`,
`factorIds`/`areaIds`, and `requestedScope` across `cmd/`, `internal/`, `specs/`,
`skills/`, `SPECIFICATION.md`, and scaffold. Append-only `log.md` files and
`archive/` are historical and stay frozen.

### Code

- [x] `internal/cli/evaluation.go` — `create` flags: drop `--narrowing`, add
      `--area`/`--factor`; `list` reads scope from manifest.
- [x] `internal/evaluation/create.go` — validate/normalize scope, write
      `RunManifest`, derive slug from `plannedScope`.
- [x] `internal/evaluation/types.go` — `Options` and receipt fields (drop
      `Narrowing`; add area/factor + manifest).
- [x] `internal/evaluation/path.go` — run-name parsing no longer the scope
      source; slug derivation from scope.
- [x] `internal/evaluation/list.go` — read number/scope from manifest.
- [x] `internal/evaluation/report_tree.go` — remove headline resolution; render
      `report.md` as the scoped Area report from `plannedScope`.
- [x] `internal/evaluation/data_contract.go` — register `RunManifest`; remove
      scope fields from `EvaluationFrame`; drop headline refs from
      `EvaluationOutputResult`.
- [x] `internal/evaluation/data.go`, `internal/evaluation/display.go` — register
      and label the `RunManifest` kind.
- [x] `internal/evaluation/report.go` — output-result/headline receipt fields.
- [x] Tests: `internal/evaluation/evaluation_test.go`,
      `internal/cli/evaluation_test.go` (and report/headline coverage).

### Durable specs & format spec

See the spec's [Durable spec changes](0149-scope-driven-evaluation-runs/spec.md#durable-spec-changes)
for the per-requirement breakdown. Index:

- [x] `specs/cli/evaluation-create.md` — flags, manifest, validation, slug.
- [x] `specs/cli/evaluation-list.md` — scope from manifest.
- [x] `specs/cli/evaluation-report.md` — report-as-area-report; no headline.
- [x] `specs/evaluation/protocol.md` — remove "Headline Result"; scope semantics.
- [x] `specs/evaluation/records/payload-kinds.md` — `RunManifest` kind; frame
      loses scope; output-result loses headline refs.
- [x] `specs/evaluation/reports/report-tree.md` — `report.md` shape/title.
- [x] `specs/evaluation/orchestration.md` — scope/headline flow.
- [x] `specs/evaluation/records/data-layout.md` — run layout includes
      `run-manifest.json`.
- [x] `specs/skills/quality-skill/evaluation.md` — resolve-before-create; flags;
      frame sheds scope.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` — same.
- [x] `specs/skills/quality-skill/quality-skill.md` — narrowing references.
- [x] `SPECIFICATION.md` — reframe any headline/scope wording.
- [x] _(suggested new, not added)_ `specs/evaluation/records/run-manifest-json.md`
      — not needed because the contract remains in `payload-kinds.md`.

### Durable docs

- [x] `skills/quality/SKILL.md` — scope guidance (narrowing → area/factor).
- [x] `skills/quality/workflows/evaluate.md` — create step, scope resolution
      ordering, frame authoring.
- [x] `skills/quality/resources/cli-workflow-conventions.md` — replace the
      "Narrowing Slug" section.

### Install / scaffold

- [x] No impact. _(Run scaffolding is produced by `create`, covered under Code.)_

## Status

`Done`. Implemented, verified, and archived.
