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

The error contract is consistent across the sampled API cases.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-001` | Common caller error cases share a documented response shape. | 💪 Strength | — | 🟢 High | The finding supports the target operability rating for caller-facing errors. | ✅ Verified: The synthetic handler matrix and API reference agree on the error envelope fields. |

## Finding details

<a id="finding-strength-001"></a>

### strength-001 Common caller error cases share a documented response shape.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 2 / 7 | 🔴 P1 Highest | Ranked by expected impact on the service quality bar and report-gallery usefulness. |

#### Condition

Validation, authorization, and conflict responses use the same synthetic error envelope.

#### Criteria

- `requirement:api::predictable-error-contracts / rating:target`: Error contracts should meet the target bar with evidence a maintainer can verify.
  Rationale: The gallery records one finding per requirement so report tables stay easy to inspect.

#### Basis

Status: ✅ Verified

The synthetic handler matrix and API reference agree on the error envelope fields.

##### Basis evidence

(none recorded)

#### Effect

The finding supports the target operability rating for caller-facing errors.

Rating effect: supports target

#### Evidence

- `synthetic-source:api/error-contracts`: The synthetic error matrix maps common failure modes to a stable code, message, and retryable flag.
  Rationale: Synthetic source reference retained to demonstrate evidence rendering.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json](../../../../data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-rating-result.json](../../../../data/areas/api/requirements/predictable-error-contracts/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

