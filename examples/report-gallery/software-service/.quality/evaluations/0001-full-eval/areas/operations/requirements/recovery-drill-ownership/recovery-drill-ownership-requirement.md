---
type: Requirement Evaluation Report
title: "Requirement: recovery drills have current owners"
---

# Requirement: recovery drills have current owners

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Operations](../../operations-area.md)

Factors: [recoverability](../../factors/recoverability/recoverability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🟡 Low / 🟡 Low |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

Recovery ownership is visible but ambiguous across the sampled records.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `gap-002` | The current owner for ledger recovery drills is ambiguous. | 🚩 Gap | 🟡 Medium | 🟡 Low | The finding limits confidence and constrains the requirement to minimum. | ⚪ Not Assessed: The available synthetic records conflict, so current ownership cannot be verified. |

## Finding details

<a id="finding-gap-002"></a>

### gap-002 The current owner for ledger recovery drills is ambiguous.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 6 / 7 | 🔴 P1 Highest | Ranked by expected impact on the service quality bar and report-gallery usefulness. |

#### Condition

The synthetic recovery calendar names a team, while the incident playbook names a former individual owner.

#### Criteria

- `requirement:operations::recovery-drill-ownership / rating:target`: Recovery drills should meet the target recoverability bar with current named ownership.
  Rationale: The gallery records one finding per requirement so report tables stay easy to inspect.

#### Basis

Status: ⚪ Not Assessed

The available synthetic records conflict, so current ownership cannot be verified.

##### Basis evidence

(none recorded)

#### Effect

The finding limits confidence and constrains the requirement to minimum.

Rating effect: constrains target

#### Evidence

- `synthetic-source:operations/recovery-calendar`: The synthetic calendar and playbook disagree about who owns the next recovery drill.
  Rationale: Synthetic source reference retained to demonstrate evidence rendering.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](../../../../data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json](../../../../data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

