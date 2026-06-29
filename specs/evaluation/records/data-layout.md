---
type: Functional Specification
title: Evaluation data layout
description: Run-folder data and generated report layout for Evaluation.
tags: [evaluation, records, layout]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation data layout

Evaluation stores structured routine data under `data/` and human-readable
reports outside `data/`.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Data Tree

An Evaluation run **MUST** store structured data under `data/`.

The run **MUST** store the CLI-generated run manifest at:

```text
data/run-manifest.json
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

Area data **MUST** live under `data/areas/` and mirror the evaluated Area tree.

Each Area data folder **MUST** own:

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

Nested Factors **MUST** recurse through nested `factors/` folders.

## Report Tree

The run-level Evaluation report **MUST** be written to:

```text
report.md
```

The run-level ranked findings index **MUST** be written to:

```text
findings.md
```

The root Area report **MUST** be written to the run root when generated:

```text
root-area.md
```

Non-root Area reports **MUST** be written under:

```text
areas/<area>/<area>-area.md
```

Requirement reports **MUST** be written under the owning Area report folder:

```text
requirements/<requirement>/<requirement>-requirement.md
```

Factor reports **MUST** be written under the owning Area report folder:

```text
factors/<factor>/<factor>-factor.md
```

Nested Factor reports **MUST** recurse through nested `factors/` folders and use
the nested Factor's local structural ID segment in the filename:

```text
factors/<factor>/factors/<sub-factor>/<sub-factor>-factor.md
```

The run-level recommendation index **MUST** be written to:

```text
recommendations.md
```

Recommendation detail reports **MUST** be written to:

```text
recommendations/<NNN>-<slug>.md
```

## Path Derivation

The CLI **MUST** derive data paths from payload `kind` and structural model IDs.
For `RecommendationResult`, the CLI **MUST** derive the data path from the
CLI-assigned recommendation `id`.

The CLI **MUST** derive report paths from report kind and structural model IDs.
For recommendation reports, the CLI **MUST** derive the report path from
recommendation ranking order and the recommendation title slug, falling back to
the CLI-assigned recommendation `id` only when needed.

The CLI **MUST NOT** derive persisted paths from display titles, natural labels,
or rendered human labels.
