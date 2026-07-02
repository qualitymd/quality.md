---
type: Recommendation Report
title: "Recommendation: Add a body-drift detector to the model self-check"
---

# Recommendation: Add a body-drift detector to the model self-check

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 6 | `qrec_refreshmodelbody` | ○ Low | 🟢 High | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_refreshmodelbody` |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary source data](#primary-source-data)

## Description

Refresh each body section's unknowns, open questions, and review provenance against the current model, then add a detector that compares recent factor and requirement changes with body review lines so stale judgment context is caught before evaluation.

## Background

The body-drift finding is inferential today: a reviewer noticed body context lagging the model. The quality loop should keep that judgment while adding a computational sensor for the repeatable stale-context failure mode.

## Expected value

Evaluations start from current judgment context, and future model growth has a repeatable detector rather than relying only on reviewer memory.

## Done criterion

Every body section reflects the current model, each section's review line postdates the latest model-shape changelog entry, and a recorded body-drift detector fails when a factor or requirement changes without a corresponding body review refresh.

## Ranking rationale

Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[gap-006]","subject":{"requirementId":"requirement:quality-md::the-model-follows-the-authoring-guide-family"}}`

## Primary source data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_refreshmodelbody/recommendation-result.json](../data/advice/recommendations/qrec_refreshmodelbody/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)
