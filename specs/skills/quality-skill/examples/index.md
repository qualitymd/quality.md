# /quality skill — reference examples

Worked reference artifacts that make the skill's
[Reporting](../quality-skill.md#reporting) contract concrete. Each is a captured
instance of what the skill writes at runtime, stored here as OKF concepts so the
spec can point at a real example.

At runtime the skill writes these into the **evaluated** repository under
`quality/evaluations/NNNN-<scope>-quality-eval/` (see
[Reporting](../quality-skill.md#reporting)); the only difference here is the OKF
frontmatter wrapper that lets the bundle catalogue them.

# Examples

- [0001 — Sparrow Payments, whole model](0001-payments-quality-eval/report.md) -
  a subject evaluation of a small fictional payments service, held at
  **Unacceptable** by a committed live credential. Exercises `file:line`
  evidence, the secret-by-reference rule, a prompt-injection finding treated as
  data, a *not assessed* requirement, and standalone
  [recommendation](0001-payments-quality-eval/recommendations/001-rotate-committed-gateway-key.md)
  artifacts with done-criteria. The
  [model evaluated](0001-payments-quality-eval/model.md) is reproduced alongside
  so the findings trace to declared requirements.
