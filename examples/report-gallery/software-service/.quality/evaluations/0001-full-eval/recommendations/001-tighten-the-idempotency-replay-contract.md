---
type: Recommendation Report
title: "Recommendation: Tighten the idempotency replay contract"
---

# Recommendation: Tighten the idempotency replay contract

Run: [Run 0001](../report.md) - Run ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../report.md) - [Findings](../findings.md) - [Recommendations](../recommendations.md)

| # | Rank | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 1 | 1 | ▲ High | 🔵 Medium | `evaluation:20260629T120000Z-0123456789ab/recommendation/1` |

Jump to: [Description](#description) - [Background](#background) - [Expected value](#expected-value) - [Done criterion](#done-criterion) - [Trace](#trace)

## Description

Specify and test partial-write replay behavior for mutation endpoints.

## Background

The highest-ranked synthetic gap leaves retry semantics below the target API correctness bar.

## Expected value

Callers and agents can verify retry behavior without inferring undocumented recovery semantics.

## Done criterion

The API contract and retry tests describe duplicate, replayed, and partial-write idempotency outcomes.

## Ranking rationale

Recommendation rank follows the synthetic finding priority and expected quality-management value.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[gap-001]","subject":{"requirementId":"requirement:api::idempotent-mutations"}}`


## Legend

- `—` - not applicable or not recorded.

## Primary Source Data

- [data/run-manifest.json](../data/run-manifest.json)
- [data/advice/recommendations/001/recommendation-result.json](../data/advice/recommendations/001/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

