---
type: Requirement Evaluation Report
title: "Requirement: persistence access is least-privilege"
---

# Requirement: persistence access is least-privilege

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger persistence](../../persistence-area.md)

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

Database role manifests separate service, migration, and analytics privileges; dependency audit reports no high-severity data-store client issue.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `auto-persistence-access-is-least-privilege` | Database role manifests separate service, migration, and analytics privileges; dependency audit reports no high-severity data-store client issue. | 💪 Strength | — | 🟢 High | The evidence supports the target rating. | ✅ Verified: Role manifests and dependency audit output were inspected. |

## Finding details

<a id="finding-auto-persistence-access-is-least-privilege"></a>

### auto-persistence-access-is-least-privilege Database role manifests separate service, migration, and analytics privileges; dependency audit reports no high-severity data-store client issue.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 32 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Do database roles and dependency audit results preserve persistence least privilege?

#### Criteria

- `requirement:persistence::persistence-access-is-least-privilege / rating:target`: The recorded evidence satisfies the target criterion for this requirement.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

Role manifests and dependency audit output were inspected.

##### Basis evidence

(none recorded)

#### Effect

The evidence supports the target rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:persistence/role-manifests`: Database role manifests separate service, migration, and analytics privileges; dependency audit reports no high-severity data-store client issue.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/persistence/requirements/persistence-access-is-least-privilege/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/persistence-access-is-least-privilege/requirement-assessment-result.json)
- [data/areas/persistence/requirements/persistence-access-is-least-privilege/requirement-rating-result.json](../../../../data/areas/persistence/requirements/persistence-access-is-least-privilege/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
