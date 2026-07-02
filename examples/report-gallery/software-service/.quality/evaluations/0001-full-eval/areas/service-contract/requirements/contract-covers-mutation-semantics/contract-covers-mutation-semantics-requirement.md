---
type: Requirement Evaluation Report
title: "Requirement: the contract defines retry, idempotency, and error semantics for every mutation endpoint"
---

# Requirement: the contract defines retry, idempotency, and error semantics for every mutation endpoint

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Service contract](../../service-contract-area.md)

Factors: [completeness](../../factors/completeness/completeness-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

Twelve of fourteen mutation endpoints have complete semantics; two lack replay clauses.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `gap-005` | Two mutation endpoints have no replay semantics in the contract. | 🚩 Gap | 🟡 Medium | 🟢 High | Integrators retrying reversals or adjustments are working from silence, which constrains completeness to minimum. | ✅ Verified: All fourteen indexed endpoints were checked clause by clause. |

## Finding details

<a id="finding-gap-005"></a>

### gap-005 Two mutation endpoints have no replay semantics in the contract.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 4 / 42 | 🟠 P2 High | Two endpoints with silent replay semantics share a root cause and a fix with the P1 replay gap. |

#### Condition

The reversal and adjustment endpoints appear in the contract's endpoint index with request and response shapes but no retry or idempotency clauses.

#### Criteria

- `requirement:service-contract::contract-covers-mutation-semantics / rating:target`: Every endpoint in the contract's index defines retry, idempotency, and error semantics.
  Rationale: The index is the population, so the totality claim is checkable and the absences are exact.

#### Basis

Status: ✅ Verified

All fourteen indexed endpoints were checked clause by clause.

##### Basis evidence

(none recorded)

#### Effect

Integrators retrying reversals or adjustments are working from silence, which constrains completeness to minimum.

Rating effect: constrains target

#### Evidence

- `synthetic-source:service-contract`: The reversal and adjustment entries end at the response shape; the other twelve entries carry retry, idempotency, and error clauses.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/service-contract/requirements/contract-covers-mutation-semantics/requirement-assessment-result.json](../../../../data/areas/service-contract/requirements/contract-covers-mutation-semantics/requirement-assessment-result.json)
- [data/areas/service-contract/requirements/contract-covers-mutation-semantics/requirement-rating-result.json](../../../../data/areas/service-contract/requirements/contract-covers-mutation-semantics/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
