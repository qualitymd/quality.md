---
type: Requirement Evaluation Report
title: "Requirement: health signals explain customer impact"
---

# Requirement: health signals explain customer impact

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Operations](../../operations-area.md)

Factors: [observability](../../factors/observability/observability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🔵 Target | ✅ Assessed | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

The sampled telemetry explains customer impact for core failure modes.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-003` | Health dashboards connect service errors to failed customer actions. | 💪 Strength | — | 🔵 Medium | The finding supports the target observability rating. | ✅ Verified: The synthetic dashboard inventory maps technical symptoms to customer-visible failed actions. |

## Finding details

<a id="finding-strength-003"></a>

### strength-003 Health dashboards connect service errors to failed customer actions.

| Advice Rank | Tier | Ranking Rationale |
| --- | --- | --- |
| 5 / 7 | 🟡 P3 Medium | Ranked by expected impact on the service quality bar and report-gallery usefulness. |

#### Condition

Synthetic dashboards include failed ledger mutations, retry exhaustion, and queue delay panels.

#### Criteria

- `requirement:operations::customer-impact-telemetry / rating:target`: Operational signals should meet the target observability bar by explaining customer impact.
  Rationale: The gallery records one finding per requirement so report tables stay easy to inspect.

#### Basis

Status: ✅ Verified

The synthetic dashboard inventory maps technical symptoms to customer-visible failed actions.

##### Basis evidence

(none recorded)

#### Effect

The finding supports the target observability rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:operations/customer-impact-dashboard`: The synthetic dashboard inventory includes panels for failed customer mutations and retry exhaustion.
  Rationale: Synthetic source reference retained to demonstrate evidence rendering.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json](../../../../data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json](../../../../data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

