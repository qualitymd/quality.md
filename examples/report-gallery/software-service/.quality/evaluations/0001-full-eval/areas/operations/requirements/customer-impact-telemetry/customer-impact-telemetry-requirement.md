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

Dashboards-as-code map service symptoms to failed customer actions, and the deployed dashboards match the definitions.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-008` | Health dashboards express failed customer actions, not just service symptoms. | 💪 Strength | — | 🔵 Medium | Operators can answer "who is hurt" directly from the dashboards, supporting the target observability rating. | ✅ Verified: The definitions were diffed against the deployed dashboards; the panel inventory was walked class by class. |

## Finding details

<a id="finding-strength-008"></a>

### strength-008 Health dashboards express failed customer actions, not just service symptoms.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 20 / 21 | ⚪ P4 Low | Dashboards answer customer impact directly. |

#### Condition

The committed definitions include failed-mutation, retry-exhaustion, and queue-delay panels each denominated in affected customer actions, and the deployed dashboards match the definitions.

#### Criteria

- `requirement:operations::customer-impact-telemetry / rating:target`: Committed dashboard definitions map symptom classes to failed customer actions and match what is deployed.

#### Basis

Status: ✅ Verified

The definitions were diffed against the deployed dashboards; the panel inventory was walked class by class.

##### Basis evidence

(none recorded)

#### Effect

Operators can answer "who is hurt" directly from the dashboards, supporting the target observability rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:operations/dashboards-as-code`: Panel definitions express failure counts as failed customer actions per minute; the deployment diff is empty.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json](../../../../data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json](../../../../data/areas/operations/requirements/customer-impact-telemetry/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

