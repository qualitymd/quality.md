---
type: Area Evaluation Report
title: "Area: Ledger persistence"
---

# Area: Ledger persistence

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md) / [Ledger persistence](persistence-area.md)

## Key details

| Overall rating | Local rating | Confidence |
| --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Area / factor breakdown](#area--factor-breakdown)
- [Requirements](#requirements)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

Balance integrity is outstanding on two independent sensors; unrehearsed rollback against the current schema caps the area at minimum.

## Area / factor breakdown

| ▦ Area / □ Factor | Overall rating | Local rating | Findings | Recommendations |
| --- | --- | --- | --- | --- |
| **[▦ Ledger persistence](persistence-area.md)** | 🟡 Minimum | 🟡 Minimum | 6 | 1 |
| ↳ [□ Auditability](factors/auditability/auditability-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ [□ Durability](factors/durability/durability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 0 |
| ↳ [□ Integrity](factors/integrity/integrity-factor.md) | 🟢 Outstanding | 🟢 Outstanding | 2 | 0 |
| ↳ [□ Recoverability](factors/recoverability/recoverability-factor.md) | 🟡 Minimum | 🟡 Minimum | 2 | 1 |
| ↳ [□ Security](factors/security/security-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [audit events are ordered and tamper-evident](requirements/audit-events-are-ordered-and-tamper-evident/audit-events-are-ordered-and-tamper-evident-requirement.md) | 🔵 Target | ✅ Assessed | [auditability](factors/auditability/auditability-factor.md) |
| [ledger mutations preserve balance invariants](requirements/balance-invariants/balance-invariants-requirement.md) | 🟢 Outstanding | ✅ Assessed | [integrity](factors/integrity/integrity-factor.md) |
| [migrations have rollback paths rehearsed against the current schema](requirements/migration-rollback/migration-rollback-requirement.md) | 🟡 Minimum | ✅ Assessed | [recoverability](factors/recoverability/recoverability-factor.md) |
| [persistence access is least-privilege](requirements/persistence-access-is-least-privilege/persistence-access-is-least-privilege-requirement.md) | 🔵 Target | ✅ Assessed | [security](factors/security/security-factor.md) |
| [reconciliation explains every balance change](requirements/reconciliation-explains-balance-changes/reconciliation-explains-balance-changes-requirement.md) | 🟢 Outstanding | ✅ Assessed | [integrity](factors/integrity/integrity-factor.md); [auditability](factors/auditability/auditability-factor.md) |
| [restore drills replay current backups without ledger loss](requirements/restore-drills-replay-current-backups/restore-drills-replay-current-backups-requirement.md) | 🔵 Target | ✅ Assessed | [durability](factors/durability/durability-factor.md); [recoverability](factors/recoverability/recoverability-factor.md) |

## Limits and incomplete inputs

| Type | Scope | Impact |
| --- | --- | --- |
| (no limits or incomplete inputs) | — | — |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/persistence/area-analysis-result.json](../../data/areas/persistence/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](../../data/advice/recommendation-ranking-result.json)
- [data/areas/persistence/requirements/audit-events-are-ordered-and-tamper-evident/requirement-rating-result.json](../../data/areas/persistence/requirements/audit-events-are-ordered-and-tamper-evident/requirement-rating-result.json)
- [data/areas/persistence/requirements/audit-events-are-ordered-and-tamper-evident/requirement-assessment-result.json](../../data/areas/persistence/requirements/audit-events-are-ordered-and-tamper-evident/requirement-assessment-result.json)
- [data/areas/persistence/requirements/balance-invariants/requirement-rating-result.json](../../data/areas/persistence/requirements/balance-invariants/requirement-rating-result.json)
- [data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json](../../data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-rating-result.json](../../data/areas/persistence/requirements/migration-rollback/requirement-rating-result.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json](../../data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json)
- [data/areas/persistence/requirements/persistence-access-is-least-privilege/requirement-rating-result.json](../../data/areas/persistence/requirements/persistence-access-is-least-privilege/requirement-rating-result.json)
- [data/areas/persistence/requirements/persistence-access-is-least-privilege/requirement-assessment-result.json](../../data/areas/persistence/requirements/persistence-access-is-least-privilege/requirement-assessment-result.json)
- [data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-rating-result.json](../../data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-rating-result.json)
- [data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-assessment-result.json](../../data/areas/persistence/requirements/reconciliation-explains-balance-changes/requirement-assessment-result.json)
- [data/areas/persistence/requirements/restore-drills-replay-current-backups/requirement-rating-result.json](../../data/areas/persistence/requirements/restore-drills-replay-current-backups/requirement-rating-result.json)
- [data/areas/persistence/requirements/restore-drills-replay-current-backups/requirement-assessment-result.json](../../data/areas/persistence/requirements/restore-drills-replay-current-backups/requirement-assessment-result.json)
