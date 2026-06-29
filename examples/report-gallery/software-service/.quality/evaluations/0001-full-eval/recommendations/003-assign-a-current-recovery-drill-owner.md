---
type: Recommendation Report
title: Assign a current recovery drill owner
data:
  - data/run-manifest.json
  - data/advice/recommendations/rec-003/recommendation-result.json
  - data/advice/recommendation-ranking-result.json
  - data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json
---

# Recommendation: Assign a current recovery drill owner

Run: [#1](../report.md) - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../report.md) - [Findings](../findings.md) - [Recommendations](../recommendations.md)

| Rank | Impact | Confidence |
| --- | --- | --- |
| 3 | Medium | 🟡 Low |

Jump to: [Description](#description) - [Background](#background) - [Expected value](#expected-value) - [Done criterion](#done-criterion) - [Trace](#trace)

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
