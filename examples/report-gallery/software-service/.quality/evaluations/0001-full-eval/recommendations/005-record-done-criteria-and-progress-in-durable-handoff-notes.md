---
type: Recommendation Report
title: "Recommendation: Record done criteria and progress in durable handoff notes"
---

# Recommendation: Record done criteria and progress in durable handoff notes

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 5 | `qrec_durablehandoffs` | ● Medium | 🔵 Medium | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_durablehandoffs` |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary source data](#primary-source-data)

## Description

Add done criteria and the confirming sensor to the quality-loop handoff template, and keep in-flight decisions and remaining work in a durable progress note rather than chat scrollback.

## Background

Most handoffs omit done criteria, and in-flight decisions live in conversation history that a session loss would erase. One template change and one recording habit address both weaknesses.

## Expected value

Agents can declare completion against criteria and resume interrupted work from records, lifting task specifiability and continuity toward target.

## Done criterion

The handoff template carries done-criteria and sensor fields, and the next three multi-session efforts each leave a durable progress note that names decisions, remaining work, and verification status.

## Ranking rationale

Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[gap-001]","subject":{"requirementId":"requirement:root::quality-loop-work-items-carry-done-criteria"}}`
- `{"kind":"RequirementAssessmentResult","selector":"findings[risk-001]","subject":{"requirementId":"requirement:root::handoffs-survive-session-loss"}}`

## Primary source data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_durablehandoffs/recommendation-result.json](../data/advice/recommendations/qrec_durablehandoffs/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

