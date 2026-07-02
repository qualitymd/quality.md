---
type: Factor Evaluation Report
title: "Factor: Integrity"
---

# Factor: Integrity

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger persistence](../../persistence-area.md)

Factor: [Integrity](integrity-factor.md)

## Key details

| Overall rating | Local rating | Status | Confidence |
| --- | --- | --- | --- |
| 🟢 Outstanding | 🟢 Outstanding | ✅ Analyzed / ✅ Analyzed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Requirements](#requirements)
- [Sub-factors](#sub-factors)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

Two independent sensors show balance preservation with margin: property coverage and zero reconciliation drift.

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [ledger mutations preserve balance invariants](../../requirements/balance-invariants/balance-invariants-requirement.md) | 🟢 Outstanding | ✅ Assessed |
| [reconciliation explains every balance change](../../requirements/reconciliation-explains-balance-changes/reconciliation-explains-balance-changes-requirement.md) | 🟢 Outstanding | ✅ Assessed |

## Sub-factors

| Factor | Path | Local rating | + Sub-factors rating |
| --- | --- | --- | --- |
| (no sub-factors) | — | — | — |

## Limits and incomplete inputs

| Type | Scope | Impact |
| --- | --- | --- |
| (no limits or incomplete inputs) | — | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/persistence/factors/integrity/factor-analysis-result.json](../../../../data/areas/persistence/factors/integrity/factor-analysis-result.json)
- [data/areas/persistence/requirements/balance-invariants/requirement-rating-result.json](../../../../data/areas/persistence/requirements/balance-invariants/requirement-rating-result.json)
- [data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json)
- [data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-rating-result.json](../../../../data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-rating-result.json)
- [data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-assessment-result.json)
