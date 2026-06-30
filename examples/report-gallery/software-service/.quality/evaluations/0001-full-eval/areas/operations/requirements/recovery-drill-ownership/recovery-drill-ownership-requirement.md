---
type: Requirement Evaluation Report
title: "Requirement: recovery drills have current owners"
---

# Requirement: recovery drills have current owners

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../../../report.md) - [Findings](../../../../findings.md) - [Recommendations](../../../../recommendations.md)

Area: [LedgerLite Service](../../../../root-area.md) / [Operations](../../operations-area.md)

Factors: [recoverability](../../factors/recoverability/recoverability-factor.md)

## Key Details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🟡 Low / 🟡 Low |

Legend

- *Quality rating:* 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable
- *Assessment status:* ✅ Assessed, 🟡 Partially Assessed, ⚪ Not Assessed, ⛔ Blocked
- *Confidence:* 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None
- *Empty:* `—`

## Contents

- [Summary](#summary)
- [Findings Summary](#findings-summary)
- [Finding Details](#finding-details)
- [Unknowns & Missing Evidence](#unknowns--missing-evidence)
- [Primary Source Data](#primary-source-data)

## Summary

Recovery ownership is visible but ambiguous across the sampled records.

## Findings Summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `gap-002` | The current owner for ledger recovery drills is ambiguous. | 🚩 Gap | 🟡 Medium | 🟡 Low | The finding limits confidence and constrains the requirement to minimum. | ⚪ Not Assessed: The available synthetic records conflict, so current ownership cannot be verified. |

Legend

- *Finding type:* 🚩 Gap, ⚠️ Risk, ✅ Strength, ℹ️ Note
- *Finding severity:* 🔴 Critical, 🔴 High, 🟡 Medium, 🔵 Low
- *Confidence:* 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None
- *Finding basis:* ✅ Verified, 🟡 Plausible, ⚪ Not Assessed, ⬜ Not Applicable
- *Empty:* `—`

## Finding Details

<a id="finding-gap-002"></a>

### gap-002 The current owner for ledger recovery drills is ambiguous.

| Advice Rank | Tier | Ranking Rationale |
| --- | --- | --- |
| 6 / 7 | 🔴 P1 Highest | Ranked by expected impact on the service quality bar and report-gallery usefulness. |

Legend

- *Finding rank:* 🔴 P1 Highest, 🟠 P2 High, 🟡 P3 Medium, ⚪ P4 Low
- *Empty:* `—`

#### Condition

The synthetic recovery calendar names a team, while the incident playbook names a former individual owner.

#### Criteria

- `requirement:operations::recovery-drill-ownership / rating:target`: Recovery drills should meet the target recoverability bar with current named ownership.
  Rationale: The gallery records one finding per requirement so report tables stay easy to inspect.

#### Basis

Status: ⚪ Not Assessed

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
| (none recorded) | — |

## Primary Source Data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](../../../../data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json](../../../../data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

