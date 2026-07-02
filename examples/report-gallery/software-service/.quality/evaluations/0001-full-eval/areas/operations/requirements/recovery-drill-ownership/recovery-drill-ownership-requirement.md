---
type: Requirement Evaluation Report
title: "Requirement: recovery drills have current owners and recent practice records"
---

# Requirement: recovery drills have current owners and recent practice records

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Operations](../../operations-area.md)

Factors: [recoverability](../../factors/recoverability/recoverability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| ⚪ Not Rated | ⚪ Not Assessed | ⚪ None / ⚪ None |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

Ownership records conflict; the requirement could not be assessed from the available material.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `note-003` | Recovery drill ownership records contradict each other. | ℹ️ Note | — | 🟡 Low | The requirement is recorded as not assessed; restoring assessability is the actionable next step. | ⚪ Not Assessed: The records conflict, so neither can serve as the basis for a judgment. |

## Finding details

<a id="finding-note-003"></a>

### note-003 Recovery drill ownership records contradict each other.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 5 / 42 | 🟠 P2 High | Missing evidence on drill ownership blocks the recoverability judgment entirely; restoring assessability is cheap. |

#### Condition

The calendar's owner field names the platform team; the playbook's escalation section names an engineer who left the rotation in the spring reorg; no drill record exists after the reorg.

#### Criteria

- `requirement:operations::recovery-drill-ownership / rating:target`: The calendar and playbook agree on the current drill owner and a drill record exists within the last quarter.
  Rationale: The contradiction itself is the observation; it blocks assessment rather than failing it.

#### Basis

Status: ⚪ Not Assessed

The records conflict, so neither can serve as the basis for a judgment.

##### Basis evidence

(none recorded)

#### Effect

The requirement is recorded as not assessed; restoring assessability is the actionable next step.

Rating effect: blocks rating

#### Evidence

- `synthetic-source:operations/recovery-calendar`: The calendar and playbook disagree on the owner, and the drill log's latest entry predates the reorg.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| 🔎 Missing Evidence | The recovery calendar names the platform team while the incident playbook names a former individual owner; neither record postdates the reorg. |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](../../../../data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json](../../../../data/areas/operations/requirements/recovery-drill-ownership/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
