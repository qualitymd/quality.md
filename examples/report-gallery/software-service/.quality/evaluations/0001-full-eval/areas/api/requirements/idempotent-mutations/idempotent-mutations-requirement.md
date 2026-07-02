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

Duplicate-key replay is specified and tested; replay after an interrupted write is neither specified nor tested.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `gap-004` | Replay behavior after an interrupted write is unspecified and untested. | 🚩 Gap | 🔴 High | 🔵 Medium | Integrators cannot know whether an interrupted-write retry double-posts, which constrains correctness to minimum. | ✅ Verified: The contract's retry section and the contract-test case list were both searched; the interrupted-write case appears in neither. |
| `note-001` | Replay traffic is a meaningful share of mutation volume. | ℹ️ Note | — | 🟢 High | The unspecified replay path is exercised daily, not theoretically. | ✅ Verified: The telemetry rollup for the window was queried directly. |

## Finding details

<a id="finding-gap-004"></a>

### gap-004 Replay behavior after an interrupted write is unspecified and untested.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 1 / 42 | 🔴 P1 Highest | Unspecified interrupted-write replay on a money-moving API is the highest-exposure gap, exercised daily by real replay traffic. |

#### Condition

The contract defines duplicate-key replay, and the contract tests cover it; neither addresses a retry that arrives after a write failed mid-transaction.

#### Criteria

- `requirement:api::idempotent-mutations / rating:target`: Contract tests prove duplicate, replayed, and interrupted-write retries produce exactly-once ledger effects.
  Rationale: Interrupted writes are exactly the case integrators retry hardest against.

#### Basis

Status: ✅ Verified

The contract's retry section and the contract-test case list were both searched; the interrupted-write case appears in neither.

##### Basis evidence

(none recorded)

#### Effect

Integrators cannot know whether an interrupted-write retry double-posts, which constrains correctness to minimum.

Rating effect: constrains target

#### Evidence

- `synthetic-source:service-contract`: The retry section defines duplicate-key semantics only; the contract-test suite has no interrupted-write replay case.

<a id="finding-note-001"></a>

### note-001 Replay traffic is a meaningful share of mutation volume.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 11 / 42 | ⚪ P4 Low | Replay-volume context informs the P1 gap's severity; no separate action. |

#### Condition

The four-week telemetry window shows about 4% of mutation requests carrying a previously seen idempotency key.

#### Criteria

- `requirement:api::idempotent-mutations / rating:target`: Contract tests prove duplicate, replayed, and interrupted-write retries produce exactly-once ledger effects.
  Rationale: Volume context sizes how often the unspecified path is exercised.

#### Basis

Status: ✅ Verified

The telemetry rollup for the window was queried directly.

##### Basis evidence

(none recorded)

#### Effect

The unspecified replay path is exercised daily, not theoretically.

Rating effect: informs severity

#### Evidence

- `synthetic-source:telemetry/mutation-replays`: The replay-rate panel reports 3.8-4.3% of mutation requests reusing an idempotency key across the window.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json](../../../../data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json](../../../../data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
