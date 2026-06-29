---
type: Recommendation Index Report
title: Recommendations
---

# Recommendations

Run: [Run 0001](report.md) - Run ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](report.md) - [Findings](findings.md) - Recommendations

## Key Details

| Recommendations | Highest Impact | Coverage |
| --- | --- | --- |
| 3 ranked recommendations | ▲ High | ✅ Addressed by Recommendation: 3 / ⬜ Not Advice Driving: 4 |

Impact: ◆ Very high, ▲ High, ● Medium, ○ Low.
Coverage: ✅ Addressed by Recommendation, ⬜ Not Advice Driving.
Empty: `—`.

## Contents

- [Ranked Recommendations](#ranked-recommendations)
- [Coverage](#coverage)
- [Primary Source Data](#primary-source-data)

## Ranked Recommendations

| Rank | # | Recommendation | Area / Factors | Impact | Confidence | Reason | Ranking Rationale |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 1 | [Tighten the idempotency replay contract](recommendations/001-tighten-the-idempotency-replay-contract.md) | [Public API](areas/api/api-area.md) / [Correctness](areas/api/factors/correctness/correctness-factor.md) | ▲ High | 🔵 Medium | Callers and agents can verify retry behavior without inferring undocumented recovery semantics. | Recommendation rank follows the synthetic finding priority and expected quality-management value. |
| 2 | 2 | [Rehearse migration rollback after schema changes](recommendations/002-rehearse-migration-rollback-after-schema-changes.md) | [Ledger Persistence](areas/persistence/persistence-area.md) / [Recoverability](areas/persistence/factors/recoverability/recoverability-factor.md) | ▲ High | 🔵 Medium | Release risk drops because rollback instructions are proven against current migrations. | Recommendation rank follows the synthetic finding priority and expected quality-management value. |
| 3 | 3 | [Assign a current recovery drill owner](recommendations/003-assign-a-current-recovery-drill-owner.md) | [Operations](areas/operations/operations-area.md) / [Recoverability](areas/operations/factors/recoverability/recoverability-factor.md) | ● Medium | 🟡 Low | Incident preparation has a clear owner agents and maintainers can route to. | Recommendation rank follows the synthetic finding priority and expected quality-management value. |

Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.

## Coverage

- ✅ Addressed by Recommendation: 3
- ⬜ Not Advice Driving: 4

## Primary Source Data

- [data/run-manifest.json](data/run-manifest.json)
- [data/advice/recommendation-ranking-result.json](data/advice/recommendation-ranking-result.json)
- [data/advice/recommendations/001/recommendation-result.json](data/advice/recommendations/001/recommendation-result.json)
- [data/advice/recommendations/002/recommendation-result.json](data/advice/recommendations/002/recommendation-result.json)
- [data/advice/recommendations/003/recommendation-result.json](data/advice/recommendations/003/recommendation-result.json)

