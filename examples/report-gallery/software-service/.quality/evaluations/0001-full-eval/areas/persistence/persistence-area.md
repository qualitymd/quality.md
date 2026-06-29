---
type: Area Evaluation Report
title: Ledger Persistence
data:
  - data/evaluation-output-result.json
  - data/areas/persistence/area-analysis-result.json
---

# Area: Ledger Persistence

Run: [#1](../../report.md) - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../report.md) - [Findings](../../findings.md) - [Recommendations](../../recommendations.md)

Area: [LedgerLite Service](../../root-area.md) / [Ledger Persistence](persistence-area.md)

| Overall Rating | Local Rating | Confidence | Data |
| --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🔵 Medium / 🔵 Medium | [area-analysis-result.json](../../data/areas/persistence/area-analysis-result.json) |

Summary:

Ledger integrity is well covered, while rollback rehearsal evidence is stale.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| Integrity is driven by ledger mutations preserve balance invariants. | supports target | [{"kind":"FactorAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"factorId":"factor:persistence::integrity"}}] |
| Recoverability is driven by migrations have rehearsed rollback paths. | constrains target | [{"kind":"FactorAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"factorId":"factor:persistence::recoverability"}}] |

## Factors

| Factor | Path | Local Rating | + Sub-Factors Rating | Sub-Factors |
| --- | --- | --- | --- | --- |
| [Integrity](factors/integrity/integrity-factor.md) | `persistence::integrity` | 🔵 Target | — | — |
| [Recoverability](factors/recoverability/recoverability-factor.md) | `persistence::recoverability` | 🟡 Minimum | — | — |

## Child Areas

| Area | Path | Local Rating | + Child Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| (no Child Areas) |  |  |  |  |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [ledger mutations preserve balance invariants](requirements/balance-invariants/balance-invariants-requirement.md) | 🔵 Target | ✅ Assessed | [integrity](factors/integrity/integrity-factor.md) |
| [migrations have rehearsed rollback paths](requirements/migration-rollback/migration-rollback-requirement.md) | 🟡 Minimum | ✅ Assessed | [recoverability](factors/recoverability/recoverability-factor.md) |

## Limits & Incomplete Inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |

## Legend

- `—` - not applicable or not recorded.
