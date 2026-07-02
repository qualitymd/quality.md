---
type: Requirement Evaluation Report
title: "Requirement: money movement flow is analyzable from entry point to ledger write"
---

# Requirement: money movement flow is analyzable from entry point to ledger write

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Codebase](../../codebase-area.md)

Factors: [maintainability/analyzability](../../factors/maintainability/factors/analyzability/analyzability-factor.md)

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

The handler-to-ledger trace is readable, but complexity check flags the reversal path as hard to follow.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-money-flow-is-analyzable` | The handler-to-ledger trace is readable, but complexity check flags the reversal path as hard to follow. | 🚩 Gap | 🟡 Medium | 🔵 Medium | The evidence leaves a visible gap, constraining the result below target while remaining acceptable. | ✅ Verified: The trace and complexity check output were inspected. |

## Finding details

<a id="finding-auto-money-flow-is-analyzable"></a>

### auto-money-flow-is-analyzable The handler-to-ledger trace is readable, but complexity check flags the reversal path as hard to follow.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 36 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Can a maintainer trace money movement from handler to ledger write?

#### Criteria

- `requirement:codebase::money-flow-is-analyzable / rating:target`: Falls short of target with visible gaps or limited evidence, while remaining acceptable.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The trace and complexity check output were inspected.

##### Basis evidence

(none recorded)

#### Effect

The evidence leaves a visible gap, constraining the result below target while remaining acceptable.

Rating effect: constrains target

#### Evidence

- `synthetic-source:codebase/handler-ledger-trace`: The handler-to-ledger trace is readable, but complexity check flags the reversal path as hard to follow.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/codebase/requirements/money-flow-is-analyzable/requirement-assessment-result.json](../../../../data/areas/codebase/requirements/money-flow-is-analyzable/requirement-assessment-result.json)
- [data/areas/codebase/requirements/money-flow-is-analyzable/requirement-rating-result.json](../../../../data/areas/codebase/requirements/money-flow-is-analyzable/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
