---
type: Factor Evaluation Report
title: Observability
data:
  - data/evaluation-output-result.json
  - data/areas/operations/factors/observability/factor-analysis-result.json
---

# Factor: Observability

Run: [#1](../../../../report.md) - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../../../report.md) - [Findings](../../../../findings.md) - [Recommendations](../../../../recommendations.md)

Area: [LedgerLite Service](../../../../root-area.md) / [Operations](../../operations-area.md)

Factor: [Observability](observability-factor.md)

| Overall Rating | Local Rating | Status | Confidence | Data |
| --- | --- | --- | --- | --- |
| 🔵 Target | 🔵 Target | ✅ Analyzed / ✅ Analyzed | 🔵 Medium / 🔵 Medium | [factor-analysis-result.json](../../../../data/areas/operations/factors/observability/factor-analysis-result.json) |

Summary:

Observability follows its direct requirement signal.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| Health dashboards connect service errors to failed customer actions. | supports target | [{"kind":"RequirementRatingResult","subject":{"requirementId":"requirement:operations::customer-impact-telemetry"}}] |

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [health signals explain customer impact](../../requirements/customer-impact-telemetry/customer-impact-telemetry-requirement.md) | 🔵 Target | ✅ Assessed |

## Sub-Factors

| Factor | Path | Local Rating | + Sub-Factors Rating |
| --- | --- | --- | --- |
| (no Sub-Factors) |  |  |  |

## Limits & Incomplete Inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |

## Legend

- `—` - not applicable or not recorded.
