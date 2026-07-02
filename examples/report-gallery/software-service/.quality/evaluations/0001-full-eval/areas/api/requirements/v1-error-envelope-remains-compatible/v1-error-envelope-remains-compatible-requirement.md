---
type: Requirement Evaluation Report
title: "Requirement: v1 error-envelope behavior remains compatible during deprecation"
---

# Requirement: v1 error-envelope behavior remains compatible during deprecation

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [compatibility](../../factors/compatibility/compatibility-factor.md); [operability](../../factors/operability/operability-factor.md)

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

The compatibility matrix and handlers still preserve documented v1 fields, but one deprecated field remains undocumented in the contract appendix.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-v1-error-envelope-remains-compatible` | The compatibility matrix and handlers still preserve documented v1 fields, but one deprecated field remains undocumented in the contract appendix. | 🚩 Gap | 🟡 Medium | 🔵 Medium | The evidence leaves a visible gap, constraining the result below target while remaining acceptable. | ✅ Verified: The compatibility matrix and handler matrix were compared; integrator usage data is partial. |

## Finding details

<a id="finding-auto-v1-error-envelope-remains-compatible"></a>

### auto-v1-error-envelope-remains-compatible The compatibility matrix and handlers still preserve documented v1 fields, but one deprecated field remains undocumented in the contract appendix.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 26 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Do v1 callers still receive documented fields through the deprecation window?

#### Criteria

- `requirement:api::v1-error-envelope-remains-compatible / rating:target`: Falls short of target with visible gaps or limited evidence, while remaining acceptable.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The compatibility matrix and handler matrix were compared; integrator usage data is partial.

##### Basis evidence

(none recorded)

#### Effect

The evidence leaves a visible gap, constraining the result below target while remaining acceptable.

Rating effect: constrains target

#### Evidence

- `synthetic-source:api/compatibility-matrix`: The compatibility matrix and handlers still preserve documented v1 fields, but one deprecated field remains undocumented in the contract appendix.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-assessment-result.json](../../../../data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-assessment-result.json)
- [data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-rating-result.json](../../../../data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
