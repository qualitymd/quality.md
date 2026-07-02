---
type: Requirement Evaluation Report
title: "Requirement: ledger entry signs match caller intent"
---

# Requirement: ledger entry signs match caller intent

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [correctness](../../factors/correctness/correctness-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🔵 Target | ✅ Assessed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

Contract tests cover debit, credit, refund, and reversal sign semantics, and all sampled records match caller intent.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-ledger-entry-signs-match-intent` | Contract tests cover debit, credit, refund, and reversal sign semantics, and all sampled records match caller intent. | 💪 Strength | — | 🟢 High | The evidence supports the target rating. | ✅ Verified: The contract tests and sampled ledger records were compared directly. |

## Finding details

<a id="finding-auto-ledger-entry-signs-match-intent"></a>

### auto-ledger-entry-signs-match-intent Contract tests cover debit, credit, refund, and reversal sign semantics, and all sampled records match caller intent.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 22 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Do contract tests show ledger signs matching caller intent across operation kinds?

#### Criteria

- `requirement:api::ledger-entry-signs-match-intent / rating:target`: The recorded evidence satisfies the target criterion for this requirement.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The contract tests and sampled ledger records were compared directly.

##### Basis evidence

(none recorded)

#### Effect

The evidence supports the target rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:api/contract-tests`: Contract tests cover debit, credit, refund, and reversal sign semantics, and all sampled records match caller intent.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/ledger-entry-signs-match-intent/requirement-assessment-result.json](../../../../data/areas/api/requirements/ledger-entry-signs-match-intent/requirement-assessment-result.json)
- [data/areas/api/requirements/ledger-entry-signs-match-intent/requirement-rating-result.json](../../../../data/areas/api/requirements/ledger-entry-signs-match-intent/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
