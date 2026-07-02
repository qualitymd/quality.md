---
type: Requirement Evaluation Report
title: "Requirement: the sensor catalog names reusable sensors consistently"
---

# Requirement: the sensor catalog names reusable sensors consistently

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Agent harness](../../agent-harness-area.md)

Factors: [completeness](../../factors/completeness/completeness-factor.md); [currentness](../../factors/currentness/currentness-factor.md); [assessability](../../factors/assessability/assessability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

Every named catalog sensor is referenced by at least two assessments, except the new drift detector, which appears once as the model-body maturation target.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-sensor-catalog-names-reusable-sensors` | Every named catalog sensor is referenced by at least two assessments, except the new drift detector, which appears once as the model-body maturation target. | 🚩 Gap | 🟡 Medium | 🔵 Medium | The evidence leaves a visible gap, constraining the result below target while remaining acceptable. | ✅ Verified: The sensor catalog and QUALITY.md assessment text were compared. |

## Finding details

<a id="finding-auto-sensor-catalog-names-reusable-sensors"></a>

### auto-sensor-catalog-names-reusable-sensors Every named catalog sensor is referenced by at least two assessments, except the new drift detector, which appears once as the model-body maturation target.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 41 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Are sensor names reused consistently across model assessments?

#### Criteria

- `requirement:agent-harness::sensor-catalog-names-reusable-sensors / rating:target`: Falls short of target with visible gaps or limited evidence, while remaining acceptable.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The sensor catalog and QUALITY.md assessment text were compared.

##### Basis evidence

(none recorded)

#### Effect

The evidence leaves a visible gap, constraining the result below target while remaining acceptable.

Rating effect: constrains target

#### Evidence

- `synthetic-source:agent-harness/sensor-catalog`: Every named catalog sensor is referenced by at least two assessments, except the new drift detector, which appears once as the model-body maturation target.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/agent-harness/requirements/sensor-catalog-names-reusable-sensors/requirement-assessment-result.json](../../../../data/areas/agent-harness/requirements/sensor-catalog-names-reusable-sensors/requirement-assessment-result.json)
- [data/areas/agent-harness/requirements/sensor-catalog-names-reusable-sensors/requirement-rating-result.json](../../../../data/areas/agent-harness/requirements/sensor-catalog-names-reusable-sensors/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
