---
type: Evaluation Design
title: Sparrow Payments — Evaluation Design (0001)
description: Reference Evaluation Design — the resolved parameters and inputs the 0001 whole-model subject evaluation was run under.
tags: [skill, quality, evaluation, example]
timestamp: 2026-06-17T00:00:00Z
---

> **Reference instance — non-normative.** A captured example of the **design**
> artifact the [`/quality` skill](../../quality-skill.md) writes alongside its
> [report](report.md) (see [Reporting](../quality-skill.md#reporting)). It records
> *what* was evaluated and *under what inputs*, so the run can be reproduced. At
> runtime the skill writes it as
> `quality/evaluations/0001-payments-quality-eval/design.md` in the **evaluated**
> repository. See the [examples index](../index.md) for the facts shared across
> this bundle.

# Evaluation design

The resolved inputs this run is bound to. The [report](report.md) states the
**Scope** for the reader; the full parameterization lives here, stated once.

## Parameters

| Parameter       | Value                                                           |
| --------------- | --------------------------------------------------------------- |
| **Mode**        | `evaluate`                                                      |
| **Altitude**    | subject — the entities the model measures, not the model itself |
| **Target file** | `./QUALITY.md`                                                  |
| **Scope**       | whole model — no target or factor narrowing                     |
| **Effort**      | `standard`                                                      |

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
