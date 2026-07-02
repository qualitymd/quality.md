---
type: Requirement Evaluation Report
title: "Requirement: contract semantics match shipped handler behavior"
---

# Requirement: contract semantics match shipped handler behavior

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Service contract](../../service-contract-area.md)

Factors: [consistency](../../factors/consistency/consistency-factor.md)

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

Where the contract speaks, handlers agree; one deprecated response field still ships that the contract no longer documents.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-006` | Shipped behavior matches every specified contract clause in the latest sensor run. | 💪 Strength | — | 🟢 High | Judgments made against the contract hold for the shipped service, supporting the target consistency rating. | ✅ Verified: The sensor run report was read alongside the clause list; coverage is one test per clause. |
| `note-002` | The deprecated balance_after field still ships but is no longer documented. | ℹ️ Note | — | 🔵 Medium | Silent behavior outside the contract invites accidental breakage when the field is eventually removed. | ✅ Verified: Response samples and the contract revision history were compared. |

## Finding details

<a id="finding-strength-006"></a>

### strength-006 Shipped behavior matches every specified contract clause in the latest sensor run.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 19 / 42 | ⚪ P4 Low | Contract-to-behavior conformance holds where specified. |

#### Condition

The contract-test suite exercises each specified clause against the running service and passed in full on the most recent run.

#### Criteria

- `requirement:service-contract::contract-matches-shipped-behavior / rating:target`: Every specified contract clause matches shipped handler behavior in the latest contract-test run.

#### Basis

Status: ✅ Verified

The sensor run report was read alongside the clause list; coverage is one test per clause.

##### Basis evidence

(none recorded)

#### Effect

Judgments made against the contract hold for the shipped service, supporting the target consistency rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:api/contract-tests`: The latest run reports all clause-mapped cases passing; the suite's clause coverage table shows no unmapped clause.

<a id="finding-note-002"></a>

### note-002 The deprecated balance_after field still ships but is no longer documented.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 12 / 42 | ⚪ P4 Low | The deprecated field is drift-in-waiting worth watching, not acting on yet. |

#### Condition

Mutation responses include balance_after, removed from the contract in the v1.4 revision; two integrators are known to still read it.

#### Criteria

- `requirement:service-contract::contract-matches-shipped-behavior / rating:target`: Every specified contract clause matches shipped handler behavior in the latest contract-test run.
  Rationale: The field contradicts no clause, but undocumented shipped behavior is drift in the making.

#### Basis

Status: ✅ Verified

Response samples and the contract revision history were compared.

##### Basis evidence

(none recorded)

#### Effect

Silent behavior outside the contract invites accidental breakage when the field is eventually removed.

Rating effect: informational

#### Evidence

- `synthetic-source:api/handler-matrix`: Response samples carry balance_after; the contract's changelog shows its documentation removed in v1.4 without a removal plan for the field itself.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/service-contract/requirements/contract-matches-shipped-behavior/requirement-assessment-result.json](../../../../data/areas/service-contract/requirements/contract-matches-shipped-behavior/requirement-assessment-result.json)
- [data/areas/service-contract/requirements/contract-matches-shipped-behavior/requirement-rating-result.json](../../../../data/areas/service-contract/requirements/contract-matches-shipped-behavior/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
