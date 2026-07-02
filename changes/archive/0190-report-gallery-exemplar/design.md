---
type: Design Doc
title: Report gallery exemplar design
description: How the report-gallery generator is restructured to carry the exemplar model, payloads, and changelog.
tags: [examples, report-gallery]
timestamp: 2026-07-02T00:00:00Z
---

# Report gallery exemplar design

The generator keeps its shape — write static files, create the run through
`evaluation.CreateRun`, pin identity, validate and write payloads through
`evaluation.SetData`, build reports with `evaluation.BuildReport` — with three
structural changes.

## Embedded content files instead of Go string constants

`QUALITY.md`, `README.md`, and the quality changelog entries move out of Go
constants into `scripts/report-gallery/content/` and are embedded with
`go:embed`. Markdown authoring stops requiring Go edits, and `prettier` formats
the source files directly, which keeps the emitted files `fmt-md-check`-stable
by construction. The changelog directory is removed and rewritten on each run,
the same as `.quality/evaluations/`, so the whole example stays
generator-owned.

Considered and rejected: keeping Go constants (multi-hundred-line raw strings
with escaping hazards and no prettier coverage), and hand-authoring the example
without the generator (loses the report-format fidelity and the
`report-gallery-check` drift gate — the reason the gallery is generated at
all).

## Generalized payload tables

The one-factor / one-finding `requirementCase` struct becomes:

- `requirementCase` with `Factors []string`, a `Findings []findingCase` slice,
  assessment/rating statuses (for the `not_assessed` case), optional
  `ratingDrivers` per finding, and per-requirement applied criteria (so the
  measured override renders its own bands).
- an explicit factor table (area → factor tree) driving factor frames and
  analyses from the requirement→factor mapping, many-to-many; the
  `agent-harnessability` umbrella analysis uses `childFactorAnalysisRefs` over
  the seven sub-factor analyses.
- area analysis drivers derived from factor analyses as today, plus root local
  analysis over the root's own sub-factor requirements (the root is no longer
  requirement-less).

Ranking, recommendation, and coverage payloads move from index-based heuristics
to explicit tables keyed by finding ID, since coverage now includes
multi-finding recommendations and a restore-assessability recommendation.

## Determinism

No new sources of nondeterminism: run identity stays pinned, tables are
authored literals, and map iteration is avoided by keeping ordered slices.
`mise run report-gallery && git diff --exit-code examples/report-gallery`
remains the byte-stability check.
