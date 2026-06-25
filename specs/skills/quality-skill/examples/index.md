# /quality skill — reference examples

Worked legacy reference artifacts from the pre-Evaluation v2 workflow. They are
kept to illustrate historical report and recommendation artifacts, not as the
current runtime shape for new `/quality evaluate` runs.

These examples are domain-illustrative, not defaults for QUALITY.md. `0001`
models a software payments service; `0002` models a non-software data product.
Together they demonstrate that the runtime artifact shape is invariant while
the domain-carried content - Factors, Requirements, Assessments, Findings, and
recommendations - changes with the modeled entity's own needs and risks. For the
doctrine behind that pairing, see
[Modeling quality across domains](../../../../docs/guides/model-quality-across-domains.md)
and its worked documentation-set example.

These example files intentionally match the legacy raw runtime shape: no OKF
frontmatter, JSON assessment/analysis records, recommendation Markdown with
runtime YAML frontmatter, and generated `report-summary.md`, `report.md`, and
`report.json`. New Evaluation v2 runs instead persist structured routine outputs
under `data/` and render the deterministic report tree specified by
[Evaluation v2 reports](../../../evaluation-v2/reports/report-tree.md).

Shared across this bundle, so the individual files need not repeat it: the
subjects, commit revisions, source paths, and every locator are fictional and
illustrative; the models use the suggested four-level rating scale
(**Outstanding** > **Target** > **Minimum** > **Unacceptable**); and `model.md`,
`design.md`, `plan.md`, and the `assessments/` and `analysis/`
records are the snapshot, inputs, method, and write-once evidence trail the
skill writes alongside the report (see
[Reporting](../reporting.md#reporting)), so the report's findings trace to
the model, parameters, coverage, and per-requirement and per-area records in
force at evaluation time. The checked-in reports are generated from the adjacent
runtime records; update them by changing the fixture inputs and rebuilding the
report, not by hand-editing the generated report files. This fixture is
intentionally complete and reportable: it records one assessment per in-scope
requirement and one analysis per area node.

# Examples

- **[0001 — Sparrow Payments, full evaluation](0001-quality-eval/report.md)**
  — a quality evaluation of a small fictional payments service, held at
  **Unacceptable** by a committed live credential. Exercises `file:line`
  evidence, the secret-by-reference rule, a prompt-injection finding treated as
  data, a *not assessed* requirement, a cross-area secondary factor (the
  Ledger's double-entry invariant also lensing the root Reliability factor), a
  two-level nested area (Webhooks → Delivery) whose intermediate
  aggregate (Minimum) differs from its local rating (Target) because a child
  subtree pulls it down, a **layered binding constraint** (rotating the
  committed credential lifts the root only to Minimum, where the webhook-delivery
  gap then binds), and standalone recommendation artifacts with done-criteria.

  **Read it in order:** the
  [model evaluated](0001-quality-eval/model.md) →
  the [design](0001-quality-eval/design.md) and
  [plan](0001-quality-eval/plan.md) →
  the [assessment result records](0001-quality-eval/assessments/) →
  the [analysis records](0001-quality-eval/analysis/) →
  the [summary](0001-quality-eval/report-summary.md),
  [human report](0001-quality-eval/report.md), and
  [JSON report](0001-quality-eval/report.json) → its recommendations
  ([001 — rotate the committed key](0001-quality-eval/recommendations/001-rotate-committed-gateway-key.md),
  [002 — produce reconciliation evidence](0001-quality-eval/recommendations/002-produce-reconciliation-evidence.md),
  [003 — bound the webhook dedup window](0001-quality-eval/recommendations/003-bound-webhook-dedup-window.md)).

- **[0002 — City Bike Station Data, full evaluation](0002-city-bike-stations-quality-eval/report.md)**
  — a quality evaluation of a small fictional station data product, held at
  **Unacceptable** by missing source/acquisition metadata in the dataset
  manifest. Exercises the same runtime artifact shape as `0001` without a
  software source or runnable oracle: reviewer judgment profiles a table,
  compares it with a roster and schema, and inspects provenance/freshness
  metadata. The Factors are earned from the data product's own risks
  (**Fitness for use**, **Provenance**, **Structural validity**, **Freshness**),
  not adopted as a default QUALITY.md factor family.

  **Read it in order:** the
  [model evaluated](0002-city-bike-stations-quality-eval/model.md) →
  the [design](0002-city-bike-stations-quality-eval/design.md) and
  [plan](0002-city-bike-stations-quality-eval/plan.md) →
  the [assessment result records](0002-city-bike-stations-quality-eval/assessments/) →
  the [analysis records](0002-city-bike-stations-quality-eval/analysis/) →
  the [summary](0002-city-bike-stations-quality-eval/report-summary.md),
  [human report](0002-city-bike-stations-quality-eval/report.md), and
  [JSON report](0002-city-bike-stations-quality-eval/report.json) → its
  recommendations
  ([001 — backfill the missing active station](0002-city-bike-stations-quality-eval/recommendations/001-backfill-missing-active-station.md),
  [002 — record snapshot provenance](0002-city-bike-stations-quality-eval/recommendations/002-record-snapshot-provenance.md)).
