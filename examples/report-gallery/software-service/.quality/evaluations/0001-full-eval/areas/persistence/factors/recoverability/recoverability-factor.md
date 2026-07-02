---
type: Factor Evaluation Report
title: "Factor: Recoverability"
---

# Factor: Recoverability

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger persistence](../../persistence-area.md)

Factor: [Recoverability](recoverability-factor.md)

## Key details

| Overall rating | Local rating | Status | Confidence |
| --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | ✅ Analyzed / ✅ Analyzed | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Requirements](#requirements)
- [Sub-factors](#sub-factors)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

Restore drills pass, but rollback guidance is unrehearsed against the two most recent schema migrations.

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [migrations have rollback paths rehearsed against the current schema](../../requirements/migration-rollback/migration-rollback-requirement.md) | 🟡 Minimum | ✅ Assessed |
| [restore drills replay current backups without ledger loss](../../requirements/restore-drills-replay-current-backups/restore-drills-replay-current-backups-requirement.md) | 🔵 Target | ✅ Assessed |

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
- [data/areas/persistence/factors/recoverability/factor-analysis-result.json](../../../../data/areas/persistence/factors/recoverability/factor-analysis-result.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-rating-result.json](../../../../data/areas/persistence/requirements/migration-rollback/requirement-rating-result.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json)
- [data/areas/persistence/requirements/restore-drills-replay-current-backups/requirement-rating-result.json](../../../../data/areas/persistence/requirements/restore-drills-replay-current-backups/requirement-rating-result.json)
- [data/areas/persistence/requirements/restore-drills-replay-current-backups/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/restore-drills-replay-current-backups/requirement-assessment-result.json)
