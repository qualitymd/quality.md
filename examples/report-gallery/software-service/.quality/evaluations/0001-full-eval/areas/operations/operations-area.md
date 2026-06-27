# Area: Operations

Area: [LedgerLite Service](../../root-area.md) / [Operations](operations-area.md)

| Overall Rating | Local Rating | Confidence | Data |
| --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🟡 Low / 🟡 Low | [area-analysis-result.json](../../data/areas/operations/area-analysis-result.json) |

Summary:

Customer-impact telemetry is useful, but recovery drill ownership is ambiguous.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| Observability is driven by health signals explain customer impact. | supports target | [{"kind":"FactorAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"factorId":"factor:operations::observability"}}] |
| Recoverability is driven by recovery drills have current owners. | constrains target | [{"kind":"FactorAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"factorId":"factor:operations::recoverability"}}] |

## Factors

| Factor | Path | Local Rating | + Sub-Factors Rating | Sub-Factors |
| --- | --- | --- | --- | --- |
| [Observability](factors/observability/observability-factor.md) | `operations::observability` | 🔵 Target | — | — |
| [Recoverability](factors/recoverability/recoverability-factor.md) | `operations::recoverability` | 🟡 Minimum | — | — |

## Child Areas

| Area | Path | Local Rating | + Child Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| (no Child Areas) |  |  |  |  |

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
