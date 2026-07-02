---
type: Area Evaluation Report
title: "Area: Operations"
---

# Area: Operations

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md) / [Operations](operations-area.md)

## Key details

| Overall rating | Local rating | Confidence |
| --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🟡 Low / 🟡 Low |

## Contents

- [Summary](#summary)
- [Area / factor breakdown](#area--factor-breakdown)
- [Requirements](#requirements)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

Customer-impact telemetry meets target, but security and capacity sit at minimum and drill ownership could not be assessed.

## Area / factor breakdown

| ▦ Area / □ Factor | Overall rating | Local rating | Findings | Recommendations |
| --- | --- | --- | --- | --- |
| **[▦ Operations](operations-area.md)** | 🟡 Minimum | 🟡 Minimum | 4 | 1 |
| ↳ [□ Capacity](factors/capacity/capacity-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 0 |
| ↳ [□ Observability](factors/observability/observability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ [□ Recoverability](factors/recoverability/recoverability-factor.md) | ⛔ Blocked | ⛔ Blocked | 2 | 1 |
| ↳ [□ Security](factors/security/security-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 0 |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [break-glass access is reviewed after use](requirements/break-glass-access-is-reviewed/break-glass-access-is-reviewed-requirement.md) | 🟡 Minimum | ✅ Assessed | [security](factors/security/security-factor.md); [recoverability](factors/recoverability/recoverability-factor.md) |
| [health signals explain customer impact](requirements/customer-impact-telemetry/customer-impact-telemetry-requirement.md) | 🔵 Target | ✅ Assessed | [observability](factors/observability/observability-factor.md) |
| [holiday-peak capacity is supported by load evidence](requirements/holiday-peak-capacity-is-supported-by-load-evidence/holiday-peak-capacity-is-supported-by-load-evidence-requirement.md) | 🟡 Minimum | ✅ Assessed | [capacity](factors/capacity/capacity-factor.md) |
| [recovery drills have current owners and recent practice records](requirements/recovery-drill-ownership/recovery-drill-ownership-requirement.md) | ⚪ Not Rated | ⚪ Not Assessed | [recoverability](factors/recoverability/recoverability-factor.md) |

## Limits and incomplete inputs

| Type | Scope | Impact |
| --- | --- | --- |
| 🧩 Incomplete Inputs | reconciled-ownership-records | Until then the area's recoverability contributes no rating signal. |
| ⚠️ Evaluation Limits | drill-ownership-unassessed | The area rating reflects observability evidence only; recoverability is missing, not weak. |
| 🧩 Incomplete Inputs | reconciled-ownership-records | Until then the area's recoverability contributes no rating signal. |
| ⚠️ Evaluation Limits | drill-ownership-unassessed | The area rating reflects observability evidence only; recoverability is missing, not weak. |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/operations/area-analysis-result.json](../../data/areas/operations/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](../../data/advice/recommendation-ranking-result.json)
- [data/areas/operations/requirements/break-glass-access-is-reviewed/requirement-rating-result.json](../../data/areas/operations/requirements/break-glass-access-is-reviewed/requirement-rating-result.json)
- [data/areas/operations/requirements/break-glass-access-is-reviewed/requirement-assessment-result.json](../../data/areas/operations/requirements/break-glass-access-is-reviewed/requirement-assessment-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json](../../data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json](../../data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json)
- [data/areas/operations/requirements/holiday-peak-capacity-is-supported-by-load-evidence/requirement-rating-result.json](../../data/areas/operations/requirements/holiday-peak-capacity-is-supported-by-load-evidence/requirement-rating-result.json)
- [data/areas/operations/requirements/holiday-peak-capacity-is-supported-by-load-evidence/requirement-assessment-result.json](../../data/areas/operations/requirements/holiday-peak-capacity-is-supported-by-load-evidence/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json](../../data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](../../data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)
