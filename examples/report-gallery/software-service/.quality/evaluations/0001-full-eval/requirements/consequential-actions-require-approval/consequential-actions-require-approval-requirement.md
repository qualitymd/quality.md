---
type: Requirement Evaluation Report
title: "Requirement: money-moving and schema-changing actions require human approval"
---

# Requirement: money-moving and schema-changing actions require human approval

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md)

Factors: [agent-harnessability/containment-of-action](../../factors/agent-harnessability/factors/containment-of-action/containment-of-action-factor.md)

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

Sandbox policy and deploy configuration put money movement and production schema changes behind named approval gates.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-003` | Money movement and production schema changes are behind declarative approval gates. | 💪 Strength | — | 🟢 High | Agent autonomy on routine work does not expose the consequential actions, supporting the target containment rating. | ✅ Verified: The allowlist and pipeline approval rules were read directly and match the recorded policy. |

## Finding details

<a id="finding-strength-003"></a>

### strength-003 Money movement and production schema changes are behind declarative approval gates.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 15 / 42 | ⚪ P4 Low | Containment gates keep agent autonomy safe on consequential actions. |

#### Condition

The sandbox allowlist excludes payment execution and production migration commands, and the deploy pipeline requires a named approver for both.

#### Criteria

- `requirement:root::consequential-actions-require-approval / rating:target`: Unattended runs cannot move money, change production schemas, or escalate permissions without an approval gate.

#### Basis

Status: ✅ Verified

The allowlist and pipeline approval rules were read directly and match the recorded policy.

##### Basis evidence

(none recorded)

#### Effect

Agent autonomy on routine work does not expose the consequential actions, supporting the target containment rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:agent-harness/sandbox-policy`: The policy denies payment-execution and migration commands to unattended sessions; the deploy config lists required approvers for both action classes.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/root/requirements/consequential-actions-require-approval/requirement-assessment-result.json](../../data/areas/root/requirements/consequential-actions-require-approval/requirement-assessment-result.json)
- [data/areas/root/requirements/consequential-actions-require-approval/requirement-rating-result.json](../../data/areas/root/requirements/consequential-actions-require-approval/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
