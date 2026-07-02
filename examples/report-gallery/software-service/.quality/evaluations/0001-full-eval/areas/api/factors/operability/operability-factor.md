---
type: Factor Evaluation Report
title: "Factor: Operability"
---

# Factor: Operability

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factor: [Operability](operability-factor.md)

## Key details

| Overall rating | Local rating | Status | Confidence |
| --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | ✅ Analyzed / ✅ Analyzed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Requirements](#requirements)
- [Sub-factors](#sub-factors)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

Error responses are predictable, but compatibility drift creates caller confusion.

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [downstream dependency timeouts return safe results](../../requirements/dependency-timeouts-return-safe-results/dependency-timeouts-return-safe-results-requirement.md) | 🔵 Target | ✅ Assessed |
| [error responses are predictable for callers](../../requirements/predictable-error-contracts/predictable-error-contracts-requirement.md) | 🔵 Target | ✅ Assessed |
| [sensitive fields stay out of error responses](../../requirements/sensitive-fields-stay-out-of-error-responses/sensitive-fields-stay-out-of-error-responses-requirement.md) | 🔵 Target | ✅ Assessed |
| [v1 error-envelope behavior remains compatible during deprecation](../../requirements/v1-error-envelope-remains-compatible/v1-error-envelope-remains-compatible-requirement.md) | 🟡 Minimum | ✅ Assessed |

## Sub-factors

| Factor | Path | Local rating | + Sub-factors rating |
| --- | --- | --- | --- |
| (no sub-factors) | — | — | — |

## Limits and incomplete inputs

| Type | Scope | Impact |
| --- | --- | --- |
| (no limits or incomplete inputs) | — | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/factors/operability/factor-analysis-result.json](../../../../data/areas/api/factors/operability/factor-analysis-result.json)
- [data/areas/api/requirements/dependency-timeouts-return-safe-results/requirement-rating-result.json](../../../../data/areas/api/requirements/dependency-timeouts-return-safe-results/requirement-rating-result.json)
- [data/areas/api/requirements/dependency-timeouts-return-safe-results/requirement-assessment-result.json](../../../../data/areas/api/requirements/dependency-timeouts-return-safe-results/requirement-assessment-result.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-rating-result.json](../../../../data/areas/api/requirements/predictable-error-contracts/requirement-rating-result.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json](../../../../data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json)
- [data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-rating-result.json](../../../../data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-rating-result.json)
- [data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-assessment-result.json](../../../../data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-assessment-result.json)
- [data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-rating-result.json](../../../../data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-rating-result.json)
- [data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-assessment-result.json](../../../../data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-assessment-result.json)
