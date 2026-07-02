---
type: Requirement Evaluation Report
title: "Requirement: migrations have rollback paths rehearsed against the current schema"
---

# Requirement: migrations have rollback paths rehearsed against the current schema

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger persistence](../../persistence-area.md)

Factors: [recoverability](../../factors/recoverability/recoverability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

The runbook has rollback steps, but the last rehearsal predates the two most recent schema migrations.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `risk-002` | Rollback steps have not been rehearsed since the two most recent schema migrations landed. | ⚠️ Risk | 🟡 Medium | 🔵 Medium | A failed release over the partitioned tables could not be confidently rolled back, constraining recoverability to minimum. | 🟡 Plausible: The staleness is verified; that the steps would fail is projected from the incidents that motivated the tightened criterion. |

## Finding details

<a id="finding-risk-002"></a>

### risk-002 Rollback steps have not been rehearsed since the two most recent schema migrations landed.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 2 / 42 | 🔴 P1 Highest | An unrehearsed rollback over freshly partitioned ledger tables is the failure mode two past incidents already demonstrated. |

#### Condition

The runbook's rehearsal record is dated before migrations 041 (ledger partitioning) and 042 (currency precision), both of which touch tables the rollback steps reorganize.

#### Criteria

- `requirement:persistence::migration-rollback / rating:target`: The most recent rollback rehearsal is newer than the latest schema change.
  Rationale: The criterion was tightened to rehearsal recency after two incidents where documented steps failed against drifted schemas.

#### Basis

Status: 🟡 Plausible

The staleness is verified; that the steps would fail is projected from the incidents that motivated the tightened criterion.

##### Basis evidence

(none recorded)

#### Effect

A failed release over the partitioned tables could not be confidently rolled back, constraining recoverability to minimum.

Rating effect: constrains target

#### Evidence

- `synthetic-source:persistence/migration-runbook`: The runbook rehearsal log's latest entry predates the migration history's entries for 041 and 042.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-rating-result.json](../../../../data/areas/persistence/requirements/migration-rollback/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
