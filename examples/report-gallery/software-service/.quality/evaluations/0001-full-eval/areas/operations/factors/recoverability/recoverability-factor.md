# Factor: Recoverability

Area: [LedgerLite Service](../../../../root-area.md) / [Operations](../../operations-area.md)

Factor: [Recoverability](recoverability-factor.md)

| Overall Rating | Local Rating | Status | Confidence | Data |
| --- | --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | ✅ Analyzed / ✅ Analyzed | 🟡 Low / 🟡 Low | [factor-analysis-result.json](../../../../data/areas/operations/factors/recoverability/factor-analysis-result.json) |

Summary:

Recoverability follows its direct requirement signal.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| The current owner for ledger recovery drills is ambiguous. | constrains target | [{"kind":"RequirementRatingResult","subject":{"requirementId":"requirement:operations::recovery-drill-ownership"}}] |

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [recovery drills have current owners](../../requirements/recovery-drill-ownership/recovery-drill-ownership-requirement.md) | 🟡 Minimum | ✅ Assessed |

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
