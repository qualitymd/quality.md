---
type: Recommendation Report
title: "Recommendation: Tighten the idempotency replay contract"
---

# Recommendation: Tighten the idempotency replay contract

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key Details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 1 | `qrec_gallery1` | ⬥ High | 🔵 Medium | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_gallery1` |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary Source Data](#primary-source-data)

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

## Primary Source Data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_gallery1/recommendation-result.json](../data/advice/recommendations/qrec_gallery1/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

