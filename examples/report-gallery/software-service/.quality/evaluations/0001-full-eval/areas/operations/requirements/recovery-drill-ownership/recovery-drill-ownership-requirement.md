---
type: Requirement Evaluation Report
title: recovery drills have current owners
---

# Requirement: recovery drills have current owners

Run: [Run 0001](../../../../report.md) - Run ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../../../report.md) - [Findings](../../../../findings.md) - [Recommendations](../../../../recommendations.md)

Area: [LedgerLite Service](../../../../root-area.md) / [Operations](../../operations-area.md)

Factors: [recoverability](../../factors/recoverability/recoverability-factor.md)

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🟡 Low / 🟡 Low |

Jump to: [Findings Summary](#findings-summary) - [Finding Details](#finding-details) - [Unknowns & Missing Evidence](#unknowns--missing-evidence)

Summary:

Recovery ownership is visible but ambiguous across the sampled records.

## Findings Summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `unknown-001` | The current owner for ledger recovery drills is ambiguous. | ❓ Unknown | 🟡 Medium | 🟡 Low | The finding limits confidence and constrains the requirement to minimum. | Not Assessed: The available synthetic records conflict, so current ownership cannot be verified. |

## Finding Details

<a id="finding-unknown-001"></a>

### unknown-001 The current owner for ledger recovery drills is ambiguous.

| Advice Rank | Tier | Ranking Rationale |
| --- | --- | --- |
| 6 / 7 | P2 | Ranked by expected impact on the service quality bar and report-gallery usefulness. |

#### Condition

The synthetic recovery calendar names a team, while the incident playbook names a former individual owner.

#### Criteria

- `requirement:operations::recovery-drill-ownership / rating:target`: Recovery drills should meet the target recoverability bar with current named ownership.
  Rationale: The gallery records one finding per requirement so report tables stay easy to inspect.

#### Basis

Status: Not Assessed

The available synthetic records conflict, so current ownership cannot be verified.

##### Basis Evidence

(none recorded)

#### Effect

The finding limits confidence and constrains the requirement to minimum.

Rating effect: constrains target

#### Evidence

- `synthetic-source:operations/recovery-calendar`: The synthetic calendar and playbook disagree about who owns the next recovery drill.
  Rationale: Synthetic source reference retained to demonstrate evidence rendering.

## Unknowns & Missing Evidence

| Type | Detail |
| --- | --- |
| (none recorded) |  |

## Legend

- `—` - not applicable or not recorded.

## Source Data

- [data/run-manifest.json](../../../../data/run-manifest.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](../../../../data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json](../../../../data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

