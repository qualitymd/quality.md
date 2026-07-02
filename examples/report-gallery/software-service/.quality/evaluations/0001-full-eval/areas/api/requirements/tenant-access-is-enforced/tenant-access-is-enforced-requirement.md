---
type: Requirement Evaluation Report
title: "Requirement: tenant access is enforced for every money-moving endpoint"
---

# Requirement: tenant access is enforced for every money-moving endpoint

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [security](../../factors/security/security-factor.md)

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

The authorization matrix covers every money-moving endpoint and rejects cross-tenant account identifiers.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-tenant-access-is-enforced` | The authorization matrix covers every money-moving endpoint and rejects cross-tenant account identifiers. | 💪 Strength | — | 🟢 High | The evidence supports the target rating. | ✅ Verified: Authorization matrix tests were run against the endpoint index. |

## Finding details

<a id="finding-auto-tenant-access-is-enforced"></a>

### auto-tenant-access-is-enforced The authorization matrix covers every money-moving endpoint and rejects cross-tenant account identifiers.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 24 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Does the authorization matrix prevent cross-tenant money movement?

#### Criteria

- `requirement:api::tenant-access-is-enforced / rating:target`: The recorded evidence satisfies the target criterion for this requirement.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

Authorization matrix tests were run against the endpoint index.

##### Basis evidence

(none recorded)

#### Effect

The evidence supports the target rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:api/authorization-matrix`: The authorization matrix covers every money-moving endpoint and rejects cross-tenant account identifiers.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/tenant-access-is-enforced/requirement-assessment-result.json](../../../../data/areas/api/requirements/tenant-access-is-enforced/requirement-assessment-result.json)
- [data/areas/api/requirements/tenant-access-is-enforced/requirement-rating-result.json](../../../../data/areas/api/requirements/tenant-access-is-enforced/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
