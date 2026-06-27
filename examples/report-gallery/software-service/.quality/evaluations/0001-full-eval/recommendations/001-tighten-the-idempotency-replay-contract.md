# Tighten the idempotency replay contract

**Rank:** 1
**Impact:** High
**Confidence:** 🔵 Medium
**Data:** [recommendation-result.json](../data/advice/recommendations/rec-001/recommendation-result.json), [recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

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
