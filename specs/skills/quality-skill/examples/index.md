# /quality skill — reference examples

Worked reference artifacts that make the skill's
[Reporting](../reporting.md#reporting) contract concrete. Each is a captured
instance of what the skill writes at runtime.

At runtime the skill writes these into the **evaluated** repository under
`.quality/evaluations/NNNN[-<narrowing>]-quality-eval/` by default (see
[Reporting](../reporting.md#reporting)). These example files intentionally
match that raw runtime shape: no OKF frontmatter, JSON assessment/analysis
records, recommendation Markdown with runtime YAML frontmatter, and generated
`report-summary.md`, `report.md`, and `report.json`.

Shared across this bundle, so the individual files need not repeat it: the
subject ("Sparrow Payments"), the `9f2c1ab` commit revision, and every
`file:line` locator are fictional and illustrative; the model uses the suggested
four-level rating scale
(**Outstanding** > **Area** > **Minimum** > **Unacceptable**); and `model.md`,
`design.md`, `plan.md`, `debug-log.md`, and the `assessments/` and `analysis/`
records are the snapshot, inputs, method, process diagnostics, and write-once
evidence trail the skill writes alongside the report (see
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
  data, a *not assessed* requirement, a cross-target **secondary factor** (the
  Ledger's double-entry invariant also lensing the root Reliability factor), a
  **two-level nested target** (Webhooks → Delivery) whose intermediate
  **aggregate (Minimum) differs from its local rating (Area)** because a child
  subtree pulls it down, a **layered binding constraint** (rotating the
  committed credential lifts the root only to Minimum, where the webhook-delivery
  gap then binds), and standalone recommendation artifacts with done-criteria.

  **Read it in order:** the
  [model evaluated](0001-quality-eval/model.md) →
  the [design](0001-quality-eval/design.md) and
  [plan](0001-quality-eval/plan.md) →
  the [debug log](0001-quality-eval/debug-log.md) →
  the [assessment result records](0001-quality-eval/assessments/) →
  the [analysis records](0001-quality-eval/analysis/) →
  the [summary](0001-quality-eval/report-summary.md),
  [human report](0001-quality-eval/report.md), and
  [JSON report](0001-quality-eval/report.json) → its recommendations
  ([001 — rotate the committed key](0001-quality-eval/recommendations/001-rotate-committed-gateway-key.md),
  [002 — produce reconciliation evidence](0001-quality-eval/recommendations/002-produce-reconciliation-evidence.md),
  [003 — bound the webhook dedup window](0001-quality-eval/recommendations/003-bound-webhook-dedup-window.md)).
