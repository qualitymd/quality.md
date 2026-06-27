# Factor: Recoverability

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger Persistence](../../persistence-area.md)

Factor: [Recoverability](recoverability-factor.md)

| Overall Rating | Local Rating | Status | Confidence | Data |
| --- | --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | ✅ Analyzed / ✅ Analyzed | 🔵 Medium / 🔵 Medium | [factor-analysis-result.json](../../../../data/areas/persistence/factors/recoverability/factor-analysis-result.json) |

Summary:

Recoverability follows its direct requirement signal.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| Rollback guidance exists, but rehearsal evidence is stale. | constrains target | [{"kind":"RequirementRatingResult","subject":{"requirementId":"requirement:persistence::migration-rollback"}}] |

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [migrations have rehearsed rollback paths](../../requirements/migration-rollback/migration-rollback-requirement.md) | 🟡 Minimum | ✅ Assessed |

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
