---
type: Requirement Evaluation Report
title: "Requirement: examples explain retry and error semantics for integrators"
---

# Requirement: examples explain retry and error semantics for integrators

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Service contract](../../service-contract-area.md)

Factors: [understandability](../../factors/understandability/understandability-factor.md); [completeness](../../factors/completeness/completeness-factor.md)

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

Contract examples cover duplicate retries, authorization failures, and validation errors; interrupted-write examples are absent with the same gap as replay semantics.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-examples-explain-retry-and-error-semantics` | Contract examples cover duplicate retries, authorization failures, and validation errors; interrupted-write examples are absent with the same gap as replay semantics. | 🚩 Gap | 🟡 Medium | 🔵 Medium | The evidence leaves a visible gap, constraining the result below target while remaining acceptable. | ✅ Verified: The example inventory was reviewed against the endpoint index. |

## Finding details

<a id="finding-auto-examples-explain-retry-and-error-semantics"></a>

### auto-examples-explain-retry-and-error-semantics Contract examples cover duplicate retries, authorization failures, and validation errors; interrupted-write examples are absent with the same gap as replay semantics.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 29 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Do contract examples teach retry and error decisions callers need to make?

#### Criteria

- `requirement:service-contract::examples-explain-retry-and-error-semantics / rating:target`: Falls short of target with visible gaps or limited evidence, while remaining acceptable.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The example inventory was reviewed against the endpoint index.

##### Basis evidence

(none recorded)

#### Effect

The evidence leaves a visible gap, constraining the result below target while remaining acceptable.

Rating effect: constrains target

#### Evidence

- `synthetic-source:service-contract`: Contract examples cover duplicate retries, authorization failures, and validation errors; interrupted-write examples are absent with the same gap as replay semantics.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/service-contract/requirements/examples-explain-retry-and-error-semantics/requirement-assessment-result.json](../../../../data/areas/service-contract/requirements/examples-explain-retry-and-error-semantics/requirement-assessment-result.json)
- [data/areas/service-contract/requirements/examples-explain-retry-and-error-semantics/requirement-rating-result.json](../../../../data/areas/service-contract/requirements/examples-explain-retry-and-error-semantics/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
