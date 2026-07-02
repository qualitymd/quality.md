---
type: Requirement Evaluation Report
title: "Requirement: ledger mutations preserve balance invariants"
---

# Requirement: ledger mutations preserve balance invariants

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger persistence](../../persistence-area.md)

Factors: [integrity](../../factors/integrity/integrity-factor.md)

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

The property-based suite passes across all mutation paths and the nightly reconciliation sensor shows zero unexplained drift over the window.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-007` | Balance invariants hold across every tested mutation path with zero reconciliation drift. | 💪 Strength | — | 🟢 High | The failure the model most exists to prevent shows no signal on either sensor, supporting an outstanding integrity rating. | ✅ Verified: Suite results and the reconciliation drift report were both read for the window. |

## Finding details

<a id="finding-strength-007"></a>

### strength-007 Balance invariants hold across every tested mutation path with zero reconciliation drift.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 13 / 21 | ⚪ P4 Low | Outstanding balance integrity is the evaluation's anchor strength. |

#### Condition

The property-based suite passes for debit, credit, failed-write, and interleaved concurrent reversal paths, and the nightly reconciliation report shows zero unexplained drift for the four-week window.

#### Criteria

- `requirement:persistence::balance-invariants / rating:outstanding`: Invariant and reconciliation sensors pass over the window with coverage beyond the required mutation paths.
  Rationale: This is the model's veto requirement; the sharpened unacceptable band is why two independent sensors are consulted.

#### Basis

Status: ✅ Verified

Suite results and the reconciliation drift report were both read for the window.

##### Basis evidence

(none recorded)

#### Effect

The failure the model most exists to prevent shows no signal on either sensor, supporting an outstanding integrity rating.

Rating effect: supports outstanding

#### Evidence

- `synthetic-source:persistence/reconciliation`: The drift report reads zero unexplained cents across all accounts for the window; the suite's concurrent-reversal properties passed 10,000 generated cases.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json)
- [data/areas/persistence/requirements/balance-invariants/requirement-rating-result.json](../../../../data/areas/persistence/requirements/balance-invariants/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

