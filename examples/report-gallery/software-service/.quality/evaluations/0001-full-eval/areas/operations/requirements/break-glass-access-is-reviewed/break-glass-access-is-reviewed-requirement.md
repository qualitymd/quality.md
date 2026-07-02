---
type: Requirement Evaluation Report
title: "Requirement: break-glass access is reviewed after use"
---

# Requirement: break-glass access is reviewed after use

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Operations](../../operations-area.md)

Factors: [security](../../factors/security/security-factor.md); [recoverability](../../factors/recoverability/recoverability-factor.md)

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

Every break-glass use in the latest quarter has an approver and incident link, but one post-use access review was recorded late.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-break-glass-access-is-reviewed` | Every break-glass use in the latest quarter has an approver and incident link, but one post-use access review was recorded late. | 🚩 Gap | 🟡 Medium | 🔵 Medium | The evidence leaves a visible gap, constraining the result below target while remaining acceptable. | ✅ Verified: Break-glass logs were inspected for the latest quarter. |

## Finding details

<a id="finding-auto-break-glass-access-is-reviewed"></a>

### auto-break-glass-access-is-reviewed Every break-glass use in the latest quarter has an approver and incident link, but one post-use access review was recorded late.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 34 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Does every break-glass use have approver, incident link, and timely review?

#### Criteria

- `requirement:operations::break-glass-access-is-reviewed / rating:target`: Falls short of target with visible gaps or limited evidence, while remaining acceptable.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

Break-glass logs were inspected for the latest quarter.

##### Basis evidence

(none recorded)

#### Effect

The evidence leaves a visible gap, constraining the result below target while remaining acceptable.

Rating effect: constrains target

#### Evidence

- `synthetic-source:operations/break-glass-log`: Every break-glass use in the latest quarter has an approver and incident link, but one post-use access review was recorded late.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/operations/requirements/break-glass-access-is-reviewed/requirement-assessment-result.json](../../../../data/areas/operations/requirements/break-glass-access-is-reviewed/requirement-assessment-result.json)
- [data/areas/operations/requirements/break-glass-access-is-reviewed/requirement-rating-result.json](../../../../data/areas/operations/requirements/break-glass-access-is-reviewed/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
