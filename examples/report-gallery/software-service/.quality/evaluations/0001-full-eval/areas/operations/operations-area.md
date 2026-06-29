---
type: Area Evaluation Report
title: "Area: Operations"
---

# Area: Operations

Run: [Run 0001](../../report.md) - Run ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../report.md) - [Findings](../../findings.md) - [Recommendations](../../recommendations.md)

Area: [LedgerLite Service](../../root-area.md) / [Operations](operations-area.md)

## Key Details

| Overall Rating | Local Rating | Confidence |
| --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🟡 Low / 🟡 Low |

Ratings: 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.
Empty: `—`.

## Summary

Customer-impact telemetry is useful, but recovery drill ownership is ambiguous.

## Area / Factor Breakdown

| Area / Factor | Overall Rating | Local Rating | Findings | Recommendations |
| --- | --- | --- | --- | --- |
| **[▦ Operations](operations-area.md)** | 🟡 Minimum | 🟡 Minimum | 2 | 1 |
| ↳ [□ Observability](factors/observability/observability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ [□ Recoverability](factors/recoverability/recoverability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |

Rows: `▦` Area, `□` Factor.

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [health signals explain customer impact](requirements/customer-impact-telemetry/customer-impact-telemetry-requirement.md) | 🔵 Target | ✅ Assessed | [observability](factors/observability/observability-factor.md) |
| [recovery drills have current owners](requirements/recovery-drill-ownership/recovery-drill-ownership-requirement.md) | 🟡 Minimum | ✅ Assessed | [recoverability](factors/recoverability/recoverability-factor.md) |

Assessment: ✅ Assessed, 🟡 Partially Assessed, ⚪ Not Assessed, ⛔ Blocked.
Empty: `—`.

## Limits & Incomplete Inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |

## Primary Source Data

- [data/run-manifest.json](../../data/run-manifest.json)
- [data/areas/operations/area-analysis-result.json](../../data/areas/operations/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](../../data/advice/recommendation-ranking-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json](../../data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json](../../data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json](../../data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](../../data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)

