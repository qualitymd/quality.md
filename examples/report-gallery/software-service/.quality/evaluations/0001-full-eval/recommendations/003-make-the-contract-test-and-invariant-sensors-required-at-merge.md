---
type: Recommendation Report
title: "Recommendation: Make the contract-test and invariant sensors required at merge"
---

# Recommendation: Make the contract-test and invariant sensors required at merge

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 3 | `qrec_gatesensors` | ⬥ High | 🔵 Medium | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_gatesensors` |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary source data](#primary-source-data)

## Description

Change the merge pipeline to mark the contract-test and invariant-suite jobs as required, with a reviewable exception path for genuinely unrelated failures.

## Background

Both sensors already run on every merge and their failures carry remediation, but they cannot block a merge; two failing contract-test runs merged in the sampled month. The sensors are trustworthy enough to gate on.

## Expected value

Contract conformance and ledger invariants hold regardless of reviewer attention, converting the two strongest sensors from advisory signals into enforced standards.

## Done criterion

The merge pipeline lists both jobs as required, and a nonconforming test change is demonstrably blocked or routed through the recorded exception path.

## Ranking rationale

Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[gap-003]","subject":{"requirementId":"requirement:root::standards-gate-nonconforming-changes"}}`

## Primary source data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_gatesensors/recommendation-result.json](../data/advice/recommendations/qrec_gatesensors/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

