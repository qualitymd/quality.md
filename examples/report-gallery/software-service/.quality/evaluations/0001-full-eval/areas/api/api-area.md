---
type: Area Evaluation Report
title: "Area: Public API"
---

# Area: Public API

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md) / [Public API](api-area.md)

## Key details

| Overall rating | Local rating | Confidence |
| --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Area / factor breakdown](#area--factor-breakdown)
- [Requirements](#requirements)
- [Limits and incomplete inputs](#limits-and-incomplete-inputs)
- [Primary source data](#primary-source-data)

## Summary

Error contracts and latency meet target, but unspecified interrupted-write replay holds the API at minimum.

## Area / factor breakdown

| ▦ Area / □ Factor | Overall rating | Local rating | Findings | Recommendations |
| --- | --- | --- | --- | --- |
| **[▦ Public API](api-area.md)** | 🟡 Minimum | 🟡 Minimum | 10 | 1 |
| ↳ [□ Compatibility](factors/compatibility/compatibility-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 0 |
| ↳ [□ Correctness](factors/correctness/correctness-factor.md) | 🟡 Minimum | 🟡 Minimum | 3 | 1 |
| ↳ [□ Operability](factors/operability/operability-factor.md) | 🟡 Minimum | 🟡 Minimum | 4 | 0 |
| ↳ [□ Performance](factors/performance/performance-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ [□ Reliability](factors/reliability/reliability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 0 |
| ↳ [□ Security](factors/security/security-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ [□ Testability](factors/testability/testability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [contract tests cover every public endpoint](requirements/contract-tests-cover-public-endpoints/contract-tests-cover-public-endpoints-requirement.md) | 🔵 Target | ✅ Assessed | [testability](factors/testability/testability-factor.md) |
| [downstream dependency timeouts return safe results](requirements/dependency-timeouts-return-safe-results/dependency-timeouts-return-safe-results-requirement.md) | 🔵 Target | ✅ Assessed | [reliability](factors/reliability/reliability-factor.md); [operability](factors/operability/operability-factor.md) |
| [mutation endpoints are idempotent under retry](requirements/idempotent-mutations/idempotent-mutations-requirement.md) | 🟡 Minimum | ✅ Assessed | [correctness](factors/correctness/correctness-factor.md) |
| [ledger entry signs match caller intent](requirements/ledger-entry-signs-match-intent/ledger-entry-signs-match-intent-requirement.md) | 🔵 Target | ✅ Assessed | [correctness](factors/correctness/correctness-factor.md) |
| [p99 mutation latency stays within budget](requirements/p99-latency-within-budget/p99-latency-within-budget-requirement.md) | 🔵 Target | ✅ Assessed | [performance](factors/performance/performance-factor.md) |
| [error responses are predictable for callers](requirements/predictable-error-contracts/predictable-error-contracts-requirement.md) | 🔵 Target | ✅ Assessed | [operability](factors/operability/operability-factor.md) |
| [sensitive fields stay out of error responses](requirements/sensitive-fields-stay-out-of-error-responses/sensitive-fields-stay-out-of-error-responses-requirement.md) | 🔵 Target | ✅ Assessed | [security](factors/security/security-factor.md); [operability](factors/operability/operability-factor.md) |
| [tenant access is enforced for every money-moving endpoint](requirements/tenant-access-is-enforced/tenant-access-is-enforced-requirement.md) | 🔵 Target | ✅ Assessed | [security](factors/security/security-factor.md) |
| [v1 error-envelope behavior remains compatible during deprecation](requirements/v1-error-envelope-remains-compatible/v1-error-envelope-remains-compatible-requirement.md) | 🟡 Minimum | ✅ Assessed | [compatibility](factors/compatibility/compatibility-factor.md); [operability](factors/operability/operability-factor.md) |

## Limits and incomplete inputs

| Type | Scope | Impact |
| --- | --- | --- |
| (no limits or incomplete inputs) | — | — |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/api/area-analysis-result.json](../../data/areas/api/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](../../data/advice/recommendation-ranking-result.json)
- [data/areas/api/requirements/contract-tests-cover-public-endpoints/requirement-rating-result.json](../../data/areas/api/requirements/contract-tests-cover-public-endpoints/requirement-rating-result.json)
- [data/areas/api/requirements/contract-tests-cover-public-endpoints/requirement-assessment-result.json](../../data/areas/api/requirements/contract-tests-cover-public-endpoints/requirement-assessment-result.json)
- [data/areas/api/requirements/dependency-timeouts-return-safe-results/requirement-rating-result.json](../../data/areas/api/requirements/dependency-timeouts-return-safe-results/requirement-rating-result.json)
- [data/areas/api/requirements/dependency-timeouts-return-safe-results/requirement-assessment-result.json](../../data/areas/api/requirements/dependency-timeouts-return-safe-results/requirement-assessment-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json](../../data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json](../../data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json)
- [data/areas/api/requirements/ledger-entry-signs-match-intent/requirement-rating-result.json](../../data/areas/api/requirements/ledger-entry-signs-match-intent/requirement-rating-result.json)
- [data/areas/api/requirements/ledger-entry-signs-match-intent/requirement-assessment-result.json](../../data/areas/api/requirements/ledger-entry-signs-match-intent/requirement-assessment-result.json)
- [data/areas/api/requirements/p99-latency-within-budget/requirement-rating-result.json](../../data/areas/api/requirements/p99-latency-within-budget/requirement-rating-result.json)
- [data/areas/api/requirements/p99-latency-within-budget/requirement-assessment-result.json](../../data/areas/api/requirements/p99-latency-within-budget/requirement-assessment-result.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-rating-result.json](../../data/areas/api/requirements/predictable-error-contracts/requirement-rating-result.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json](../../data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json)
- [data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-rating-result.json](../../data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-rating-result.json)
- [data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-assessment-result.json](../../data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-assessment-result.json)
- [data/areas/api/requirements/tenant-access-is-enforced/requirement-rating-result.json](../../data/areas/api/requirements/tenant-access-is-enforced/requirement-rating-result.json)
- [data/areas/api/requirements/tenant-access-is-enforced/requirement-assessment-result.json](../../data/areas/api/requirements/tenant-access-is-enforced/requirement-assessment-result.json)
- [data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-rating-result.json](../../data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-rating-result.json)
- [data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-assessment-result.json](../../data/areas/api/requirements/v1-error-envelope-remains-compatible/requirement-assessment-result.json)
