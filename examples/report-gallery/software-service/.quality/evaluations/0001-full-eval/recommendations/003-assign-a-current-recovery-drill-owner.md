---
type: Recommendation Report
title: "Recommendation: Assign a current recovery drill owner"
---

# Recommendation: Assign a current recovery drill owner

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key Details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 3 | `qrec_gallery3` | ● Medium | 🟡 Low | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_gallery3` |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary Source Data](#primary-source-data)

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

- `{"kind":"RequirementAssessmentResult","selector":"findings[gap-002]","subject":{"requirementId":"requirement:operations::recovery-drill-ownership"}}`

## Primary Source Data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_gallery3/recommendation-result.json](../data/advice/recommendations/qrec_gallery3/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

