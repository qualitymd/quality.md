---
type: Recommendation Report
title: "Recommendation: Rehearse migration rollback after schema changes"
---

# Recommendation: Rehearse migration rollback after schema changes

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key Details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 2 | `qrec_gallery2` | ⬥ High | 🔵 Medium | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_gallery2` |

**Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary Source Data](#primary-source-data)

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

## Primary Source Data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_gallery2/recommendation-result.json](../data/advice/recommendations/qrec_gallery2/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

