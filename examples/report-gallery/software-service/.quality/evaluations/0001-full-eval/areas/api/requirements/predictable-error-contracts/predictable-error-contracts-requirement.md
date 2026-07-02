---
type: Requirement Evaluation Report
title: "Requirement: error responses are predictable for callers"
---

# Requirement: error responses are predictable for callers

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [operability](../../factors/operability/operability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🔵 Target | ✅ Assessed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

Validation, authorization, and conflict responses all use the contract's error envelope, and handlers match it.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-004` | Error responses follow one documented envelope across all sampled failure modes. | 💪 Strength | — | 🟢 High | Callers can branch on stable error codes, supporting the target operability rating. | ✅ Verified: The handler matrix and contract envelope were compared endpoint by endpoint. |

## Finding details

<a id="finding-strength-004"></a>

### strength-004 Error responses follow one documented envelope across all sampled failure modes.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 17 / 42 | ⚪ P4 Low | Stable error contracts support integrator trust. |

#### Condition

Validation, authorization, and conflict cases for every endpoint in the contract index return the envelope's code, message, and retryable fields.

#### Criteria

- `requirement:api::predictable-error-contracts / rating:target`: Common failure modes return the contract's error envelope with stable codes callers can branch on.

#### Basis

Status: ✅ Verified

The handler matrix and contract envelope were compared endpoint by endpoint.

##### Basis evidence

(none recorded)

#### Effect

Callers can branch on stable error codes, supporting the target operability rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:api/handler-matrix`: The matrix maps each failure mode to the envelope fields; no handler deviates from the documented codes.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json](../../../../data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-rating-result.json](../../../../data/areas/api/requirements/predictable-error-contracts/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
