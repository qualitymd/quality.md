# /quality skill — reference examples

Worked reference artifacts that make the skill's
[Reporting](../quality-skill.md#reporting) contract concrete. Each is a captured
instance of what the skill writes at runtime, stored here as OKF concepts so the
spec can point at a real example.

At runtime the skill writes these into the **evaluated** repository under
`quality/evaluations/NNNN-<scope>-quality-eval/` (see
[Reporting](../quality-skill.md#reporting)); the only difference here is the OKF
frontmatter wrapper that lets the bundle catalogue them.

Shared across this bundle, so the individual files need not repeat it: the
subject ("Sparrow Payments"), the `9f2c1ab` commit revision, and every
`file:line` locator are fictional and illustrative; the model uses the suggested
four-level rating scale
(**Outstanding** > **Target** > **Minimum** > **Unacceptable**); and `model.md`
reproduces the evaluated `QUALITY.md` itself — not a runtime output — so the
report's findings trace to declared requirements.

# Examples

- **[0001 — Sparrow Payments, whole model](0001-payments-quality-eval/report.md)**
  — a subject evaluation of a small fictional payments service, held at
  **Unacceptable** by a committed live credential. Exercises `file:line`
  evidence, the secret-by-reference rule, a prompt-injection finding treated as
  data, a *not assessed* requirement, a cross-target **secondary factor** (the
  Ledger's double-entry invariant also lensing the root Reliability factor), a
  **two-level nested target** (Webhooks → Delivery) whose intermediate
  **aggregate (Minimum) differs from its local rating (Target)** because a child
  subtree pulls it down, a **layered binding constraint** (rotating the
  committed credential lifts the root only to Minimum, where the webhook-delivery
  gap then binds), and standalone recommendation artifacts with done-criteria.

  **Read it in order:** the
  [model evaluated](0001-payments-quality-eval/model.md) →
  the [report](0001-payments-quality-eval/report.md) → its recommendations
  ([001 — rotate the committed key](0001-payments-quality-eval/recommendations/001-rotate-committed-gateway-key.md),
  [002 — produce reconciliation evidence](0001-payments-quality-eval/recommendations/002-produce-reconciliation-evidence.md),
  [003 — bound the webhook dedup window](0001-payments-quality-eval/recommendations/003-bound-webhook-dedup-window.md)).
