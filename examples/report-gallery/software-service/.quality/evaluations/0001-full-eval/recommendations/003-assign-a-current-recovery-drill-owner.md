# Assign a current recovery drill owner

**Rank:** 3
**Impact:** Medium
**Confidence:** 🟡 Low
**Data:** [recommendation-result.json](../data/advice/recommendations/rec-003/recommendation-result.json), [recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

## Description

Resolve the conflicting recovery-owner records and name the owner in the calendar and playbook.

## Background

Ambiguous ownership limits confidence in recovery practice.

## Expected value

Incident preparation has a clear owner agents and maintainers can route to.

## Done criterion

The recovery calendar and playbook agree on the current owner and next drill date.

## Ranking rationale

Recommendation rank follows the synthetic finding priority and expected quality-management value.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[unknown-001]","subject":{"requirementId":"requirement:operations::recovery-drill-ownership"}}`


## Legend

- `—` - not applicable or not recorded.
