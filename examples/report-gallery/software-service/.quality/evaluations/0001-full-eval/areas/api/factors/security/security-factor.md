---
type: Factor Evaluation Report
title: "Factor: Security"
---

# Factor: Security

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factor: [Security](security-factor.md)

## Key details

| Overall rating | Local rating | Status | Confidence |
| --- | --- | --- | --- |
| 🔵 Target | 🔵 Target | ✅ Analyzed / ✅ Analyzed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Requirements](#requirements)
- [Sub-factors](#sub-factors)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

Authorization and error-response sensors protect tenant access and sensitive fields.

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [sensitive fields stay out of error responses](../../requirements/sensitive-fields-stay-out-of-error-responses/sensitive-fields-stay-out-of-error-responses-requirement.md) | 🔵 Target | ✅ Assessed |
| [tenant access is enforced for every money-moving endpoint](../../requirements/tenant-access-is-enforced/tenant-access-is-enforced-requirement.md) | 🔵 Target | ✅ Assessed |

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
- [data/areas/api/factors/security/factor-analysis-result.json](../../../../data/areas/api/factors/security/factor-analysis-result.json)
- [data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-rating-result.json](../../../../data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-rating-result.json)
- [data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-assessment-result.json](../../../../data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-assessment-result.json)
- [data/areas/api/requirements/tenant-access-is-enforced/requirement-rating-result.json](../../../../data/areas/api/requirements/tenant-access-is-enforced/requirement-rating-result.json)
- [data/areas/api/requirements/tenant-access-is-enforced/requirement-assessment-result.json](../../../../data/areas/api/requirements/tenant-access-is-enforced/requirement-assessment-result.json)
