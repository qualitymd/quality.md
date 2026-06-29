---
type: Recommendation Report
title: "Recommendation: Assign a current recovery drill owner"
---

# Recommendation: Assign a current recovery drill owner

Run: [Run 0001](../report.md) - Run ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../report.md) - [Findings](../findings.md) - [Recommendations](../recommendations.md)

## Key Details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 3 | `qrec_gallery3` | ● Medium | 🟡 Low | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_gallery3` |

Impact: ⬥⬥ Very high, ⬥ High, ● Medium, ○ Low.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.
Empty: `—`.

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

- `{"kind":"RequirementAssessmentResult","selector":"findings[unknown-001]","subject":{"requirementId":"requirement:operations::recovery-drill-ownership"}}`

## Primary Source Data

- [data/run-manifest.json](../data/run-manifest.json)
- [data/advice/recommendations/qrec_gallery3/recommendation-result.json](../data/advice/recommendations/qrec_gallery3/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

