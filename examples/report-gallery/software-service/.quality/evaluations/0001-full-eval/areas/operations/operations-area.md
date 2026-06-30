---
type: Area Evaluation Report
title: "Area: Operations"
---

# Area: Operations

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md) / [Operations](operations-area.md)

## Key Details

| Overall Rating | Local Rating | Confidence |
| --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🟡 Low / 🟡 Low |

**Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

## Contents

- [Summary](#summary)
- [Area / Factor Breakdown](#area--factor-breakdown)
- [Requirements](#requirements)
- [Limits & Incomplete Inputs](#limits--incomplete-inputs)
- [Primary Source Data](#primary-source-data)

## Summary

Customer-impact telemetry is useful, but recovery drill ownership is ambiguous.

## Area / Factor Breakdown

| Area / Factor | Overall Rating | Local Rating | Findings | Recommendations |
| --- | --- | --- | --- | --- |
| **[▦ Operations](operations-area.md)** | 🟡 Minimum | 🟡 Minimum | 2 | 1 |
| ↳ [□ Observability](factors/observability/observability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ [□ Recoverability](factors/recoverability/recoverability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |

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

## Primary Source Data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/operations/area-analysis-result.json](../../data/areas/operations/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](../../data/advice/recommendation-ranking-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json](../../data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json](../../data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json](../../data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](../../data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)

