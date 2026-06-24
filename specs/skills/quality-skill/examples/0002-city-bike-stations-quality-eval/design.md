# Evaluation design

The resolved inputs this run is bound to. The [report](report.md) states the
**Scope** for the reader; the full parameterization lives here, stated once.

## Parameters

| Parameter      | Value                                         |
| -------------- | --------------------------------------------- |
| **Mode**       | `evaluate`                                    |
| **Model file** | `./QUALITY.md`                                |
| **Scope**      | full evaluation - no area or factor narrowing |
| **Rigor**      | `standard`                                    |

## Model snapshot

Bound to the model captured in [`model.md`](model.md) - the `./QUALITY.md` in
force at evaluation time, reproduced verbatim so the report's findings trace to
the exact requirements and `source` selectors they were rated against.

## Root area

- **Source:** `city-bike-stations-data` at commit `2ac94ef`, with the
  `./QUALITY.md` model at the same revision.
- **Source resolved from:** `./data/stations.csv`,
  `./schema/stations.schema.yaml`, `./metadata/manifest.yaml`, and
  `./metadata/lineage.md` - the source selectors the in-scope areas declare.
- **Evaluator:** the `/quality` skill - the deterministic `qualitymd` surface for
  structure and source resolution, agent assessment for the findings - on
  2026-06-24.
