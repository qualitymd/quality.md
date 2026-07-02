---
type: Recommendation Report
title: "Recommendation: Rehearse migration rollback against the current schema"
---

# Recommendation: Rehearse migration rollback against the current schema

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 2 | `qrec_rollbackrehearsal` | ⬥ High | 🔵 Medium | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_rollbackrehearsal` |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary source data](#primary-source-data)

## Description

Run the runbook's rollback path against a copy of the current schema — including migrations 041 and 042 — and record the rehearsal with date, schema version, and any step corrections.

## Background

The rollback steps have not been rehearsed since the partitioning and currency-precision migrations landed, and the criterion was tightened to rehearsal recency precisely because documented-but-unrehearsed steps failed in two past incidents.

## Expected value

Release risk drops because rollback instructions are proven against the schema they would actually run on, and the recoverability rating can return to target.

## Done criterion

The runbook's rehearsal log contains an entry newer than migration 042 with the rehearsed steps matching the current runbook text.

## Ranking rationale

Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[risk-002]","subject":{"requirementId":"requirement:persistence::migration-rollback"}}`

## Primary source data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_rollbackrehearsal/recommendation-result.json](../data/advice/recommendations/qrec_rollbackrehearsal/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)
