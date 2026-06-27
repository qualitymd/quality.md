# Factor: Operability

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factor: [Operability](operability-factor.md)

| Overall Rating | Local Rating | Status | Confidence | Data |
| --- | --- | --- | --- | --- |
| 🔵 Target | 🔵 Target | ✅ Analyzed / ✅ Analyzed | 🟢 High / 🟢 High | [factor-analysis-result.json](../../../../data/areas/api/factors/operability/factor-analysis-result.json) |

Summary:

Operability follows its direct requirement signal.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| Common caller error cases share a documented response shape. | supports target | [{"kind":"RequirementRatingResult","subject":{"requirementId":"requirement:api::predictable-error-contracts"}}] |

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [error responses are predictable for callers](../../requirements/predictable-error-contracts/predictable-error-contracts-requirement.md) | 🔵 Target | ✅ Assessed |

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
