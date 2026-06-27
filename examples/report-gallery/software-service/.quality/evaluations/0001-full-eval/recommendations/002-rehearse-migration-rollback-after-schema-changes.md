# Rehearse migration rollback after schema changes

**Rank:** 2
**Impact:** High
**Confidence:** 🔵 Medium
**Data:** [recommendation-result.json](../data/advice/recommendations/rec-002/recommendation-result.json), [recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

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
