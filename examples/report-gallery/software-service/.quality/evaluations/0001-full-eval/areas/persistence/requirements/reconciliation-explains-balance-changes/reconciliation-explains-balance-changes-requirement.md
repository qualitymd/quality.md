---
type: Requirement Evaluation Report
title: "Requirement: reconciliation explains every balance change"
---

# Requirement: reconciliation explains every balance change

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger persistence](../../persistence-area.md)

Factors: [integrity](../../factors/integrity/integrity-factor.md); [auditability](../../factors/auditability/auditability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟢 Outstanding | ✅ Assessed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

The reconciliation job traces every sampled balance delta to one ordered ledger event.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-reconciliation-explains-balance-changes` | The reconciliation job traces every sampled balance delta to one ordered ledger event. | 💪 Strength | — | 🟢 High | The evidence supports the outstanding rating. | ✅ Verified: The reconciliation job output and audit-event stream were compared for the four-week window. |

## Finding details

<a id="finding-auto-reconciliation-explains-balance-changes"></a>

### auto-reconciliation-explains-balance-changes The reconciliation job traces every sampled balance delta to one ordered ledger event.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 30 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Can reconciliation explain every sampled balance change from ordered ledger events?

#### Criteria

- `requirement:persistence::reconciliation-explains-balance-changes / rating:outstanding`: The recorded evidence exceeds the requirement with verified margin.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The reconciliation job output and audit-event stream were compared for the four-week window.

##### Basis evidence

(none recorded)

#### Effect

The evidence supports the outstanding rating.

Rating effect: supports outstanding

#### Evidence

- `synthetic-source:persistence/reconciliation-job`: The reconciliation job traces every sampled balance delta to one ordered ledger event.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-assessment-result.json)
- [data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-rating-result.json](../../../../data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
