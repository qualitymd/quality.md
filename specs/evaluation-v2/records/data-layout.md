---
type: Functional Specification
title: Evaluation v2 data layout
description: Run-folder data and generated report layout for Evaluation v2.
tags: [evaluation, records, layout]
timestamp: 2026-06-25T00:00:00Z
---

# Evaluation v2 data layout

Evaluation v2 stores structured routine data under `data/` and human-readable
reports outside `data/`.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Data Tree

An Evaluation v2 run **MUST** store structured data under `data/`.

The run **MUST** store the evaluation frame at:

```text
data/frame/evaluation-frame.json
```

The run **MUST** store the CLI-generated evaluation output result at:

```text
data/evaluation-output-result.json
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

The root Area report **MUST** be written to:

```text
report.md
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

## Path Derivation

The CLI **MUST** derive data paths from payload `kind` and structural model IDs.

The CLI **MUST** derive report paths from report kind and structural model IDs.

The CLI **MUST NOT** derive persisted paths from display titles, natural labels,
or rendered human labels.
