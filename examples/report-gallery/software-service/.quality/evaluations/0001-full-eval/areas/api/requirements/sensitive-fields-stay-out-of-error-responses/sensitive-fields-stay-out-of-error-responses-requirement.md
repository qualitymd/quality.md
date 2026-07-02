---
type: Requirement Evaluation Report
title: "Requirement: sensitive fields stay out of error responses"
---

# Requirement: sensitive fields stay out of error responses

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factors: [security](../../factors/security/security-factor.md); [operability](../../factors/operability/operability-factor.md)

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

Error-envelope contract tests and sampled failure payloads contain stable codes but no account numbers, bank tokens, internal IDs, or tenant secrets.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-sensitive-fields-stay-out-of-error-responses` | Error-envelope contract tests and sampled failure payloads contain stable codes but no account numbers, bank tokens, internal IDs, or tenant secrets. | 💪 Strength | — | 🟢 High | The evidence supports the target rating. | ✅ Verified: The error-envelope tests and failure payload samples were inspected. |

## Finding details

<a id="finding-auto-sensitive-fields-stay-out-of-error-responses"></a>

### auto-sensitive-fields-stay-out-of-error-responses Error-envelope contract tests and sampled failure payloads contain stable codes but no account numbers, bank tokens, internal IDs, or tenant secrets.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 25 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Do failure payloads avoid sensitive fields while preserving branchable error codes?

#### Criteria

- `requirement:api::sensitive-fields-stay-out-of-error-responses / rating:target`: The recorded evidence satisfies the target criterion for this requirement.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The error-envelope tests and failure payload samples were inspected.

##### Basis evidence

(none recorded)

#### Effect

The evidence supports the target rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:api/contract-tests`: Error-envelope contract tests and sampled failure payloads contain stable codes but no account numbers, bank tokens, internal IDs, or tenant secrets.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-assessment-result.json](../../../../data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-assessment-result.json)
- [data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-rating-result.json](../../../../data/areas/api/requirements/sensitive-fields-stay-out-of-error-responses/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
