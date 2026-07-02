---
type: Requirement Evaluation Report
title: "Requirement: contract tests cover every public endpoint"
---

# Requirement: contract tests cover every public endpoint

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [testability](../../factors/testability/testability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🔵 Target | ✅ Assessed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

The contract-test manifest maps every public endpoint to at least one success and one failure case.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-contract-tests-cover-public-endpoints` | The contract-test manifest maps every public endpoint to at least one success and one failure case. | 💪 Strength | — | 🟢 High | The evidence supports the target rating. | ✅ Verified: The endpoint index and contract-test manifest were diffed. |

## Finding details

<a id="finding-auto-contract-tests-cover-public-endpoints"></a>

### auto-contract-tests-cover-public-endpoints The contract-test manifest maps every public endpoint to at least one success and one failure case.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 27 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Does the contract-test manifest cover every endpoint in the index?

#### Criteria

- `requirement:api::contract-tests-cover-public-endpoints / rating:target`: The recorded evidence satisfies the target criterion for this requirement.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The endpoint index and contract-test manifest were diffed.

##### Basis evidence

(none recorded)

#### Effect

The evidence supports the target rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:api/contract-tests`: The contract-test manifest maps every public endpoint to at least one success and one failure case.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/contract-tests-cover-public-endpoints/requirement-assessment-result.json](../../../../data/areas/api/requirements/contract-tests-cover-public-endpoints/requirement-assessment-result.json)
- [data/areas/api/requirements/contract-tests-cover-public-endpoints/requirement-rating-result.json](../../../../data/areas/api/requirements/contract-tests-cover-public-endpoints/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
