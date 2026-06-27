# Area: LedgerLite Service

Area: [LedgerLite Service](root-area.md)

| Overall Rating | Local Rating | Confidence | Data |
| --- | --- | --- | --- |
| 🟡 Minimum | ⬜ Empty | 🔵 Medium / ⚪ None | [area-analysis-result.json](data/areas/root/area-analysis-result.json) |

Summary:

LedgerLite is usable in the synthetic evaluation, but API idempotency, rollback rehearsal, and recovery ownership keep the overall service below target.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| api area contributes to the full service roll-up. | contributes to minimum | [{"kind":"AreaAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"areaId":"area:api"}}] |
| persistence area contributes to the full service roll-up. | contributes to minimum | [{"kind":"AreaAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"areaId":"area:persistence"}}] |
| operations area contributes to the full service roll-up. | contributes to minimum | [{"kind":"AreaAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"areaId":"area:operations"}}] |
| agent-harness area contributes to the full service roll-up. | contributes to minimum | [{"kind":"AreaAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"areaId":"area:agent-harness"}}] |

## Factors

| Factor | Path | Local Rating | + Sub-Factors Rating | Sub-Factors |
| --- | --- | --- | --- | --- |
| (no local Factors) |  |  |  |  |

## Child Areas

| Area | Path | Local Rating | + Child Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| [Agent Harness](areas/agent-harness/agent-harness-area.md) | `/agent-harness` | 🔵 Target | — | [Agent Accessibility](areas/agent-harness/factors/agent-accessibility/agent-accessibility-factor.md) 🔵 Target |
| [Public API](areas/api/api-area.md) | `/api` | 🟡 Minimum | — | [Correctness](areas/api/factors/correctness/correctness-factor.md) 🟡 Minimum; [Operability](areas/api/factors/operability/operability-factor.md) 🔵 Target |
| [Operations](areas/operations/operations-area.md) | `/operations` | 🟡 Minimum | — | [Observability](areas/operations/factors/observability/observability-factor.md) 🔵 Target; [Recoverability](areas/operations/factors/recoverability/recoverability-factor.md) 🟡 Minimum |
| [Ledger Persistence](areas/persistence/persistence-area.md) | `/persistence` | 🟡 Minimum | — | [Integrity](areas/persistence/factors/integrity/integrity-factor.md) 🔵 Target; [Recoverability](areas/persistence/factors/recoverability/recoverability-factor.md) 🟡 Minimum |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| (no local Requirements) |  |  |  |

## Limits & Incomplete Inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |

## Legend

- `—` - not applicable or not recorded.
