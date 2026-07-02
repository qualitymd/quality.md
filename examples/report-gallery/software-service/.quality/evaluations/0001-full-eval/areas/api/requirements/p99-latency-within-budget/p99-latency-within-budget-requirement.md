---
type: Requirement Evaluation Report
title: "Requirement: p99 mutation latency stays within budget"
---

# Requirement: p99 mutation latency stays within budget

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [performance](../../factors/performance/performance-factor.md)

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

The four-week telemetry window puts mutation p99 at 262 ms, inside the 300 ms target band.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-005` | Mutation p99 measured 262 ms over the representative four-week window. | 💪 Strength | — | 🟢 High | Latency sits comfortably inside the recalibrated target band, supporting the target performance rating. | ✅ Verified: The rollup query was run against the recorded window with no missing weeks. |

## Finding details

<a id="finding-strength-005"></a>

### strength-005 Mutation p99 measured 262 ms over the representative four-week window.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 18 / 21 | ⚪ P4 Low | Latency sits inside the recalibrated band. |

#### Condition

The rollup query across all mutation endpoints reports p99 between 248 ms and 276 ms week over week, aggregating to 262 ms.

#### Criteria

- `requirement:api::p99-latency-within-budget / rating:target`: p99 at or under 300 ms over the window.
  Rationale: The band was recalibrated to 300 ms after the caching work landed; the measurement tests the new floor.

#### Basis

Status: ✅ Verified

The rollup query was run against the recorded window with no missing weeks.

##### Basis evidence

(none recorded)

#### Effect

Latency sits comfortably inside the recalibrated target band, supporting the target performance rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:telemetry/latency-rollup`: Weekly p99 values of 248, 259, 276, and 264 ms across the window; no week exceeded the target band.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/p99-latency-within-budget/requirement-assessment-result.json](../../../../data/areas/api/requirements/p99-latency-within-budget/requirement-assessment-result.json)
- [data/areas/api/requirements/p99-latency-within-budget/requirement-rating-result.json](../../../../data/areas/api/requirements/p99-latency-within-budget/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

