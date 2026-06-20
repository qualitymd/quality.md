# Evaluation design

The resolved inputs this run is bound to. The [report](report.md) states the
**Scope** for the reader; the full parameterization lives here, stated once.

## Parameters

| Parameter       | Value                                                           |
| --------------- | --------------------------------------------------------------- |
| **Mode**        | `evaluate`                                                      |
| **Altitude**    | subject — the entities the model measures, not the model itself |
| **Target file** | `./QUALITY.md`                                                  |
| **Scope**       | full evaluation — no target or factor narrowing                 |
| **Rigor**       | `standard`                                                      |

## Model snapshot

Bound to the model captured in [`model.md`](model.md) — the `./QUALITY.md` in
force at evaluation time, reproduced verbatim so the report's findings trace to
the exact requirements and `source` selectors they were rated against.

## Subject

- **Source:** `sparrow-payments` at commit `9f2c1ab`, with the `./QUALITY.md`
  model at the same revision.
- **Source resolved from:** `./` (root), `./ledger` (Ledger), `./webhooks`
  (Webhooks), and `./webhooks/delivery` (Delivery) — the `source` selectors the
  in-scope targets declare.
- **Evaluator:** the `/quality` skill — the deterministic `qualitymd` surface for
  structure and source resolution, agent assessment for the findings — on
  2026-06-17.
