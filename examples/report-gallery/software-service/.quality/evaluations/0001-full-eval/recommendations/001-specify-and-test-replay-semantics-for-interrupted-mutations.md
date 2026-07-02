---
type: Recommendation Report
title: "Recommendation: Specify and test replay semantics for interrupted mutations"
---

# Recommendation: Specify and test replay semantics for interrupted mutations

> **Evaluation links:** [report.md](../report.md) | [findings.md](../findings.md) | [recommendations.md](../recommendations.md) | [glossary.md](../../../../glossary.md)

Run: [0001-full-eval](../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key details

| # | ID | Impact | Confidence | Reference |
| --- | --- | --- | --- | --- |
| 1 | `qrec_replaycontract` | ⬥ High | 🟢 High | `evaluation:20260629T120000Z-0123456789ab/recommendation/qrec_replaycontract` |

## Contents

- [Description](#description)
- [Background](#background)
- [Expected value](#expected-value)
- [Done criterion](#done-criterion)
- [Ranking rationale](#ranking-rationale)
- [Trace](#trace)
- [Primary source data](#primary-source-data)

## Description

Extend the service contract's retry section to define replay outcomes after interrupted writes for all fourteen mutation endpoints, including the reversal and adjustment endpoints that currently lack replay clauses, and add contract-test cases for duplicate, replayed, and interrupted-write retries.

## Background

Interrupted-write replay is unspecified and untested while replay traffic runs near 4% of mutation volume, and two endpoints lack replay semantics entirely. Both findings trace to the same contract section, so one change closes both.

## Expected value

Integrators and agents can verify retry behavior from the contract and its sensor instead of inferring undocumented recovery semantics; the API correctness and contract completeness ratings can reach target.

## Done criterion

Every endpoint in the contract's index defines retry, idempotency, and error semantics including interrupted-write replay, and the contract-test suite covers duplicate, replayed, and interrupted-write cases for each mutation endpoint.

## Ranking rationale

Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements.

## Trace

- `{"kind":"RequirementAssessmentResult","selector":"findings[gap-004]","subject":{"requirementId":"requirement:api::idempotent-mutations"}}`
- `{"kind":"RequirementAssessmentResult","selector":"findings[gap-005]","subject":{"requirementId":"requirement:service-contract::contract-covers-mutation-semantics"}}`

## Primary source data

- [data/evaluation-manifest.json](../data/evaluation-manifest.json)
- [data/advice/recommendations/qrec_replaycontract/recommendation-result.json](../data/advice/recommendations/qrec_replaycontract/recommendation-result.json)
- [data/advice/recommendation-ranking-result.json](../data/advice/recommendation-ranking-result.json)

