---
type: Functional Specification
title: Evaluation data layout
description: Run-folder data and generated report layout for evaluation.
tags: [evaluation, records, layout]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation data layout

Evaluation has two run-folder layouts: the runner layout, where one
authoritative `evaluation.json` carries all structured data, and the historical
and manual multi-file layout under `data/`. Human-readable reports live outside
the structured data in both layouts.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Runner run layout

A run created by [`qualitymd evaluation run`](../../cli/evaluation-run.md)
**MUST** keep one authoritative structured run artifact,
[`evaluation.json`](../evaluation-json.md), at the run root, alongside
`model-snapshot.md`, run-local `logs/`, and the generated reports.

New runner-created runs **MUST NOT** write the multi-file `data/` tree as
their authoritative structured data.

> Rationale: the multi-file tree exists primarily to let skill-authored routine
> payloads be validated and persisted incrementally. A CLI-owned runner keeps
> the same structured concepts without exposing each routine as a separate
> file. — 0192

Existing multi-file evaluation runs stay readable as historical runs through
the existing commands; the runner change requires no migrations, dual writers,
or compatibility payload copies for new runs.

## Historical and manual data tree

The layout in this section is the historical and manual layout: runs created by
`qualitymd evaluation create` and populated through
`qualitymd evaluation data set`.

Such an evaluation run **MUST** store structured data under `data/`.

The run **MUST** store the CLI-generated evaluation manifest at:

```text
data/evaluation-manifest.json
```

The run **MUST** store the evaluation frame at:

```text
data/frame/evaluation-frame.json
```

The run **MUST** store the CLI-generated evaluation output result at:

```text
data/evaluation-output-result.json
```

Advice data **MUST** live under:

```text
data/advice/
```

The advice folder **MUST** contain:

- `finding-ranking-result.json`
- `recommendation-ranking-result.json`
- `recommendations/`

Each recommendation payload **MUST** be stored under:

```text
data/advice/recommendations/<recommendation-id>/recommendation-result.json
```

`<recommendation-id>` is the `RecommendationResult.id`, such as
`qrec_7h4km2p9`.

Area data **MUST** live under `data/areas/` and mirror the evaluated area tree.

Each area data folder **MUST** own:

- `area-evaluation-frame.json`
- `area-analysis-frame.json`
- `area-analysis-result.json`
- local `requirements/`
- local `factors/`
- child `areas/`

Requirement data folders **MUST** contain:

- `requirement-evaluation-frame.json`
- `requirement-assessment-result.json`
- `requirement-rating-result.json`

Factor data folders **MUST** contain:

- `factor-analysis-frame.json`
- `factor-analysis-result.json`

Nested factors **MUST** recurse through nested `factors/` folders.

## Report tree

The run-level evaluation report **MUST** be written to:

```text
report.md
```

The run-level findings report **MUST** be written to:

```text
findings.md
```

The root area report **MUST** be written to the run root when generated:

```text
root-area.md
```

Non-root area reports **MUST** be written under:

```text
areas/<area>/<area>-area.md
```

Requirement reports **MUST** be written under the owning area report folder:

```text
requirements/<requirement>/<requirement>-requirement.md
```

Factor reports **MUST** be written under the owning area report folder:

```text
factors/<factor>/<factor>-factor.md
```

Nested factor reports **MUST** recurse through nested `factors/` folders and use
the nested factor's local structural ID segment in the filename:

```text
factors/<factor>/factors/<sub-factor>/<sub-factor>-factor.md
```

The run-level recommendations report **MUST** be written to:

```text
recommendations.md
```

Recommendation detail reports **MUST** be written to:

```text
recommendations/<NNN>-<slug>.md
```

## Path derivation

The CLI **MUST** derive data paths from payload `kind` and structural model IDs.
For `RecommendationResult`, the CLI **MUST** derive the data path from the
CLI-assigned recommendation `id`.

The CLI **MUST** derive report paths from report kind and structural model IDs.
For recommendation reports, the CLI **MUST** derive the report path from
the recommendation number and the recommendation title slug. The recommendation
number is the ranking entry's `rank`, zero-padded in filenames. The slug falls
back to the CLI-assigned recommendation `id` only when needed.

The CLI **MUST NOT** derive persisted paths from display titles, natural labels,
or rendered human labels.
