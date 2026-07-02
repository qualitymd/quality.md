---
type: Area Evaluation Report
title: "Area: Codebase"
---

# Area: Codebase

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md) / [Codebase](codebase-area.md)

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

Implementation boundaries, tests, and security sensors meet target, but the reversal money-flow path is too hard to analyze.

## Area / factor breakdown

| ▦ Area / □ Factor | Overall rating | Local rating | Findings | Recommendations |
| --- | --- | --- | --- | --- |
| **[▦ Codebase](codebase-area.md)** | 🟡 Minimum | 🟡 Minimum | 5 | 0 |
| ↳ [□ Consistency](factors/consistency/consistency-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ [□ Maintainability](factors/maintainability/maintainability-factor.md) | 🟡 Minimum | ⬜ Empty | 4 | 0 |
| ↳ ↳ [□ Analyzability](factors/maintainability/factors/analyzability/analyzability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 0 |
| ↳ ↳ [□ Modifiability](factors/maintainability/factors/modifiability/modifiability-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ ↳ [□ Testability](factors/maintainability/factors/testability/testability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ [□ Security](factors/security/security-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [architecture boundaries match the service contract](requirements/architecture-boundaries-match-the-service-contract/architecture-boundaries-match-the-service-contract-requirement.md) | 🔵 Target | ✅ Assessed | [consistency](factors/consistency/consistency-factor.md); [maintainability/modifiability](factors/maintainability/factors/modifiability/modifiability-factor.md) |
| [changes remain local to owned architecture boundaries](requirements/changes-remain-local-to-owned-boundaries/changes-remain-local-to-owned-boundaries-requirement.md) | 🔵 Target | ✅ Assessed | [maintainability/modifiability](factors/maintainability/factors/modifiability/modifiability-factor.md); [consistency](factors/consistency/consistency-factor.md) |
| [dependencies and secret handling stay within policy](requirements/dependency-and-secret-handling-stay-within-policy/dependency-and-secret-handling-stay-within-policy-requirement.md) | 🔵 Target | ✅ Assessed | [security](factors/security/security-factor.md) |
| [risky implementation branches have focused tests](requirements/implementation-has-focused-tests-for-risky-branches/implementation-has-focused-tests-for-risky-branches-requirement.md) | 🔵 Target | ✅ Assessed | [maintainability/testability](factors/maintainability/factors/testability/testability-factor.md) |
| [money movement flow is analyzable from entry point to ledger write](requirements/money-flow-is-analyzable/money-flow-is-analyzable-requirement.md) | 🟡 Minimum | ✅ Assessed | [maintainability/analyzability](factors/maintainability/factors/analyzability/analyzability-factor.md) |

## Limits and incomplete inputs

| Type | Scope | Impact |
| --- | --- | --- |
| (no limits or incomplete inputs) | — | — |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/codebase/area-analysis-result.json](../../data/areas/codebase/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](../../data/advice/recommendation-ranking-result.json)
- [data/areas/codebase/requirements/architecture-boundaries-match-the-service-contract/requirement-rating-result.json](../../data/areas/codebase/requirements/architecture-boundaries-match-the-service-contract/requirement-rating-result.json)
- [data/areas/codebase/requirements/architecture-boundaries-match-the-service-contract/requirement-assessment-result.json](../../data/areas/codebase/requirements/architecture-boundaries-match-the-service-contract/requirement-assessment-result.json)
- [data/areas/codebase/requirements/changes-remain-local-to-owned-boundaries/requirement-rating-result.json](../../data/areas/codebase/requirements/changes-remain-local-to-owned-boundaries/requirement-rating-result.json)
- [data/areas/codebase/requirements/changes-remain-local-to-owned-boundaries/requirement-assessment-result.json](../../data/areas/codebase/requirements/changes-remain-local-to-owned-boundaries/requirement-assessment-result.json)
- [data/areas/codebase/requirements/dependency-and-secret-handling-stay-within-policy/requirement-rating-result.json](../../data/areas/codebase/requirements/dependency-and-secret-handling-stay-within-policy/requirement-rating-result.json)
- [data/areas/codebase/requirements/dependency-and-secret-handling-stay-within-policy/requirement-assessment-result.json](../../data/areas/codebase/requirements/dependency-and-secret-handling-stay-within-policy/requirement-assessment-result.json)
- [data/areas/codebase/requirements/implementation-has-focused-tests-for-risky-branches/requirement-rating-result.json](../../data/areas/codebase/requirements/implementation-has-focused-tests-for-risky-branches/requirement-rating-result.json)
- [data/areas/codebase/requirements/implementation-has-focused-tests-for-risky-branches/requirement-assessment-result.json](../../data/areas/codebase/requirements/implementation-has-focused-tests-for-risky-branches/requirement-assessment-result.json)
- [data/areas/codebase/requirements/money-flow-is-analyzable/requirement-rating-result.json](../../data/areas/codebase/requirements/money-flow-is-analyzable/requirement-rating-result.json)
- [data/areas/codebase/requirements/money-flow-is-analyzable/requirement-assessment-result.json](../../data/areas/codebase/requirements/money-flow-is-analyzable/requirement-assessment-result.json)
