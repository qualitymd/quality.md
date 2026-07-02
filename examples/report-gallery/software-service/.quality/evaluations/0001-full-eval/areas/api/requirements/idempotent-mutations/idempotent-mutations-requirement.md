---
type: Requirement Evaluation Report
title: "Requirement: mutation endpoints are idempotent under retry"
---

# Requirement: mutation endpoints are idempotent under retry

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [correctness](../../factors/correctness/correctness-factor.md)

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

The idempotency contract is present but incomplete for retry recovery.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `gap-001` | Mutation retry behavior is not fully specified for duplicate idempotency keys. | 🚩 Gap | 🔴 High | 🔵 Medium | The API reaches the minimum bar but does not meet the target correctness criterion for retry semantics. | ✅ Verified: The synthetic contract excerpt names idempotency keys but omits partial-write replay behavior. |

## Finding details

<a id="finding-gap-001"></a>

### gap-001 Mutation retry behavior is not fully specified for duplicate idempotency keys.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 1 / 7 | 🔴 P1 Highest | Ranked by expected impact on the service quality bar and report-gallery usefulness. |

#### Condition

The synthetic API contract describes idempotency keys, but replayed requests do not have a documented response contract for partial-write recovery.

#### Criteria

- `requirement:api::idempotent-mutations / rating:target`: Mutation endpoints should meet the target correctness bar with retry behavior a maintainer can verify.
  Rationale: The gallery records one finding per requirement so report tables stay easy to inspect.

#### Basis

Status: ✅ Verified

The synthetic contract excerpt names idempotency keys but omits partial-write replay behavior.

##### Basis evidence

(none recorded)

#### Effect

The API reaches the minimum bar but does not meet the target correctness criterion for retry semantics.

Rating effect: constrains target

#### Evidence

- `synthetic-source:api/idempotency-contract`: The synthetic contract covers idempotency-key presence and duplicate detection, but not partial-write replay outcomes.
  Rationale: Synthetic source reference retained to demonstrate evidence rendering.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json](../../../../data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json](../../../../data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

