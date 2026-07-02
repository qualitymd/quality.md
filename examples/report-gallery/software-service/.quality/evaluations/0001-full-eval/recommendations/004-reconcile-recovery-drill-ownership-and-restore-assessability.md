---
type: Recommendation Report
title: "Recommendation: Reconcile recovery drill ownership and restore assessability"
---

# Recommendation: Reconcile recovery drill ownership and restore assessability

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 4 | `qrec_drillownership` | ● Medium | 🟢 High | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_drillownership` |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary source data](#primary-source-data)

## Description

Agree the current drill owner, update the recovery calendar and incident playbook to name the same owner, and schedule the next drill so a post-reorg record exists.

## Background

The ownership requirement could not be assessed because the calendar and playbook contradict each other and no drill record postdates the reorg. The blocker is record reconciliation, not drill quality.

## Expected value

The recoverability factor becomes assessable again, and the next evaluation can rate drill practice on evidence instead of recording missing evidence.

## Done criterion

The calendar and playbook name the same current owner, and the drill log contains a post-reorg entry or a scheduled date.

## Ranking rationale

Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[note-003]","subject":{"requirementId":"requirement:operations::recovery-drill-ownership"}}`

## Primary source data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_drillownership/recommendation-result.json](../data/advice/recommendations/qrec_drillownership/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

