# Factor: Correctness

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factor: [Correctness](correctness-factor.md)

| Overall Rating | Local Rating | Status | Confidence | Data |
| --- | --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | ✅ Analyzed / ✅ Analyzed | 🔵 Medium / 🔵 Medium | [factor-analysis-result.json](../../../../data/areas/api/factors/correctness/factor-analysis-result.json) |

Summary:

Correctness follows its direct requirement signal.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| Mutation retry behavior is not fully specified for duplicate idempotency keys. | constrains target | [{"kind":"RequirementRatingResult","subject":{"requirementId":"requirement:api::idempotent-mutations"}}] |

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [mutation endpoints are idempotent under retry](../../requirements/idempotent-mutations/idempotent-mutations-requirement.md) | 🟡 Minimum | ✅ Assessed |

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
