---
type: Requirement Evaluation Report
title: "Requirement: audit events are ordered and tamper-evident"
---

# Requirement: audit events are ordered and tamper-evident

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger persistence](../../persistence-area.md)

Factors: [auditability](../../factors/auditability/auditability-factor.md)

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

Append-path tests prove monotonic sequence numbers, hash chaining, and immutable event writes.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-audit-events-are-ordered-and-tamper-evident` | Append-path tests prove monotonic sequence numbers, hash chaining, and immutable event writes. | 💪 Strength | — | 🟢 High | The evidence supports the target rating. | ✅ Verified: Audit append-path tests were run and the schema was inspected. |

## Finding details

<a id="finding-auto-audit-events-are-ordered-and-tamper-evident"></a>

### auto-audit-events-are-ordered-and-tamper-evident Append-path tests prove monotonic sequence numbers, hash chaining, and immutable event writes.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 31 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Do audit append-path tests prove ordered tamper-evident events?

#### Criteria

- `requirement:persistence::audit-events-are-ordered-and-tamper-evident / rating:target`: The recorded evidence satisfies the target criterion for this requirement.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

Audit append-path tests were run and the schema was inspected.

##### Basis evidence

(none recorded)

#### Effect

The evidence supports the target rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:persistence/audit-events`: Append-path tests prove monotonic sequence numbers, hash chaining, and immutable event writes.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/persistence/requirements/audit-events-are-ordered-and-tamper-evident/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/audit-events-are-ordered-and-tamper-evident/requirement-assessment-result.json)
- [data/areas/persistence/requirements/audit-events-are-ordered-and-tamper-evident/requirement-rating-result.json](../../../../data/areas/persistence/requirements/audit-events-are-ordered-and-tamper-evident/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
