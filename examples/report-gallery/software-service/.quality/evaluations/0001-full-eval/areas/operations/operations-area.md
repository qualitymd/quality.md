---
type: Area Evaluation Report
title: Operations
data:
  - data/run-manifest.json
  - data/areas/operations/area-analysis-result.json
  - data/areas/operations/factors/observability/factor-analysis-result.json
  - data/areas/operations/factors/recoverability/factor-analysis-result.json
  - data/advice/finding-ranking-result.json
  - data/advice/recommendation-ranking-result.json
  - data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json
  - data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json
  - data/advice/recommendations/rec-003/recommendation-result.json
  - data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json
  - data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json
---

# Area: Operations

Run: [#1](../../report.md) - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../report.md) - [Findings](../../findings.md) - [Recommendations](../../recommendations.md)

Area: [LedgerLite Service](../../root-area.md) / [Operations](operations-area.md)

| Overall Rating | Local Rating | Confidence |
| --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🟡 Low / 🟡 Low |

Summary:

Customer-impact telemetry is useful, but recovery drill ownership is ambiguous.

## Area / Factor Breakdown

| Area / Factor | Overall Rating | Local Rating | Findings | Recommendations |
| --- | --- | --- | --- | --- |
| **[Operations](operations-area.md)** | 🟡 Minimum | 🟡 Minimum | 2 | 1 |
| ↳ 🧩 [Observability](factors/observability/observability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ 🧩 [Recoverability](factors/recoverability/recoverability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [health signals explain customer impact](requirements/customer-impact-telemetry/customer-impact-telemetry-requirement.md) | 🔵 Target | ✅ Assessed | [observability](factors/observability/observability-factor.md) |
| [recovery drills have current owners](requirements/recovery-drill-ownership/recovery-drill-ownership-requirement.md) | 🟡 Minimum | ✅ Assessed | [recoverability](factors/recoverability/recoverability-factor.md) |

## Limits & Incomplete Inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |

## Legend

- `—` - not applicable or not recorded.
