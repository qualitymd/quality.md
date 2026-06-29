---
type: Area Evaluation Report
title: Public API
data:
  - data/evaluation-output-result.json
  - data/areas/api/area-analysis-result.json
---

# Area: Public API

Run: [#1](../../report.md) - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../report.md) - [Findings](../../findings.md) - [Recommendations](../../recommendations.md)

Area: [LedgerLite Service](../../root-area.md) / [Public API](api-area.md)

| Overall Rating | Local Rating | Confidence | Data |
| --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | 🔵 Medium / 🔵 Medium | [area-analysis-result.json](../../data/areas/api/area-analysis-result.json) |

Summary:

The API has predictable errors, but idempotency retry semantics need a tighter contract.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| Correctness is driven by mutation endpoints are idempotent under retry. | constrains target | [{"kind":"FactorAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"factorId":"factor:api::correctness"}}] |
| Operability is driven by error responses are predictable for callers. | supports target | [{"kind":"FactorAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"factorId":"factor:api::operability"}}] |

## Factors

| Factor | Path | Local Rating | + Sub-Factors Rating | Sub-Factors |
| --- | --- | --- | --- | --- |
| [Correctness](factors/correctness/correctness-factor.md) | `api::correctness` | 🟡 Minimum | — | — |
| [Operability](factors/operability/operability-factor.md) | `api::operability` | 🔵 Target | — | — |

## Child Areas

| Area | Path | Local Rating | + Child Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| (no Child Areas) |  |  |  |  |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [mutation endpoints are idempotent under retry](requirements/idempotent-mutations/idempotent-mutations-requirement.md) | 🟡 Minimum | ✅ Assessed | [correctness](factors/correctness/correctness-factor.md) |
| [error responses are predictable for callers](requirements/predictable-error-contracts/predictable-error-contracts-requirement.md) | 🔵 Target | ✅ Assessed | [operability](factors/operability/operability-factor.md) |

## Limits & Incomplete Inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |

## Legend

- `—` - not applicable or not recorded.
