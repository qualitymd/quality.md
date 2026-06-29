---
type: Recommendation Report
title: Rehearse migration rollback after schema changes
data:
  - data/run-manifest.json
  - data/advice/recommendations/rec-002/recommendation-result.json
  - data/advice/recommendation-ranking-result.json
  - data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json
---

# Recommendation: Rehearse migration rollback after schema changes

Run: [#1](../report.md) - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../report.md) - [Findings](../findings.md) - [Recommendations](../recommendations.md)

| Rank | Impact | Confidence |
| --- | --- | --- |
| 2 | High | 🔵 Medium |

Jump to: [Description](#description) - [Background](#background) - [Expected value](#expected-value) - [Done criterion](#done-criterion) - [Trace](#trace)

## Description

Run and record a rollback rehearsal for the latest ledger migrations.

## Background

The synthetic runbook has rollback steps, but stale rehearsal evidence leaves recoverability below target.

## Expected value

Release risk drops because rollback instructions are proven against current migrations.

## Done criterion

A current rollback rehearsal record exists for the latest migration set.

## Ranking rationale

Recommendation rank follows the synthetic finding priority and expected quality-management value.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[risk-001]","subject":{"requirementId":"requirement:persistence::migration-rollback"}}`


## Legend

- `—` - not applicable or not recorded.
