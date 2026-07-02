---
type: Area Evaluation Report
title: "Area: LedgerLite Service"
---

# Area: LedgerLite Service

> **Evaluation links:** [report.md](report.md) | [findings.md](findings.md) | [recommendations.md](recommendations.md) | [glossary.md](../../../glossary.md)

Run: [0001-full-eval](report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](root-area.md)

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

LedgerLite is money-safe today — balance integrity is outstanding — but unspecified replay semantics, unrehearsed rollback, and advisory merge gates hold the money-touching areas below the target margin the model's body requires.

## Area / factor breakdown

| ▦ Area / □ Factor | Overall rating | Local rating | Findings | Recommendations |
| --- | --- | --- | --- | --- |
| **[▦ LedgerLite Service](root-area.md)** | 🟡 Minimum | 🟡 Minimum | 21 | 6 |
| ↳ [□ Agent Harnessability](factors/agent-harnessability/agent-harnessability-factor.md) | 🟡 Minimum | ⬜ Empty | 7 | 2 |
| ↳ ↳ [□ Agent Accessibility](factors/agent-harnessability/factors/agent-accessibility/agent-accessibility-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ ↳ [□ Agent Operability](factors/agent-harnessability/factors/agent-operability/agent-operability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 0 |
| ↳ ↳ [□ Containment of Action](factors/agent-harnessability/factors/containment-of-action/containment-of-action-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ ↳ [□ Continuity](factors/agent-harnessability/factors/continuity/continuity-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ ↳ [□ Enforcement of Standards](factors/agent-harnessability/factors/enforcement-of-standards/enforcement-of-standards-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ ↳ [□ Self-Verifiability](factors/agent-harnessability/factors/self-verifiability/self-verifiability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ ↳ [□ Task Specifiability](factors/agent-harnessability/factors/task-specifiability/task-specifiability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ [▦ Agent harness](areas/agent-harness/agent-harness-area.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ ↳ [□ Assessability](areas/agent-harness/factors/assessability/assessability-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ ↳ [□ Coherence](areas/agent-harness/factors/coherence/coherence-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ ↳ [□ Completeness](areas/agent-harness/factors/completeness/completeness-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ ↳ [□ Currentness](areas/agent-harness/factors/currentness/currentness-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |
| ↳ [▦ Public API](areas/api/api-area.md) | 🟡 Minimum | 🟡 Minimum | 4 | 1 |
| ↳ ↳ [□ Correctness](areas/api/factors/correctness/correctness-factor.md) | 🟡 Minimum | 🟡 Minimum | 2 | 1 |
| ↳ ↳ [□ Operability](areas/api/factors/operability/operability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ ↳ [□ Performance](areas/api/factors/performance/performance-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ [▦ Operations](areas/operations/operations-area.md) | 🔵 Target | 🔵 Target | 2 | 1 |
| ↳ ↳ [□ Observability](areas/operations/factors/observability/observability-factor.md) | 🔵 Target | 🔵 Target | 1 | 0 |
| ↳ ↳ [□ Recoverability](areas/operations/factors/recoverability/recoverability-factor.md) | ⛔ Blocked | ⛔ Blocked | 1 | 1 |
| ↳ [▦ Ledger persistence](areas/persistence/persistence-area.md) | 🟡 Minimum | 🟡 Minimum | 2 | 1 |
| ↳ ↳ [□ Integrity](areas/persistence/factors/integrity/integrity-factor.md) | 🟢 Outstanding | 🟢 Outstanding | 1 | 0 |
| ↳ ↳ [□ Recoverability](areas/persistence/factors/recoverability/recoverability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ [▦ LedgerLite Service QUALITY.md](areas/quality-md/quality-md-area.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ ↳ [□ Context Grounding](areas/quality-md/factors/context-grounding/context-grounding-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ ↳ [□ Evaluability](areas/quality-md/factors/evaluability/evaluability-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ ↳ [□ Lifecycle Maintenance](areas/quality-md/factors/lifecycle-maintenance/lifecycle-maintenance-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ [▦ Service contract](areas/service-contract/service-contract-area.md) | 🟡 Minimum | 🟡 Minimum | 3 | 1 |
| ↳ ↳ [□ Completeness](areas/service-contract/factors/completeness/completeness-factor.md) | 🟡 Minimum | 🟡 Minimum | 1 | 1 |
| ↳ ↳ [□ Consistency](areas/service-contract/factors/consistency/consistency-factor.md) | 🔵 Target | 🔵 Target | 2 | 0 |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [a fresh agent session reaches a ready-to-work environment from recorded setup](requirements/a-fresh-session-reaches-a-ready-to-work-environment/a-fresh-session-reaches-a-ready-to-work-environment-requirement.md) | 🟡 Minimum | ✅ Assessed | [agent-harnessability/agent-operability](factors/agent-harnessability/factors/agent-operability/agent-operability-factor.md) |
| [a fresh agent reaches service context and deeper guidance from a stable entry point](requirements/agents-reach-service-context-from-a-stable-entry-point/agents-reach-service-context-from-a-stable-entry-point-requirement.md) | 🔵 Target | ✅ Assessed | [agent-harnessability/agent-accessibility](factors/agent-harnessability/factors/agent-accessibility/agent-accessibility-factor.md) |
| [money-moving and schema-changing actions require human approval](requirements/consequential-actions-require-approval/consequential-actions-require-approval-requirement.md) | 🔵 Target | ✅ Assessed | [agent-harnessability/containment-of-action](factors/agent-harnessability/factors/containment-of-action/containment-of-action-factor.md) |
| [in-flight work survives session loss through durable handoff records](requirements/handoffs-survive-session-loss/handoffs-survive-session-loss-requirement.md) | 🟡 Minimum | ✅ Assessed | [agent-harnessability/continuity](factors/agent-harnessability/factors/continuity/continuity-factor.md) |
| [quality-loop work items carry scoped goals and done criteria](requirements/quality-loop-work-items-carry-done-criteria/quality-loop-work-items-carry-done-criteria-requirement.md) | 🟡 Minimum | ✅ Assessed | [agent-harnessability/task-specifiability](factors/agent-harnessability/factors/task-specifiability/task-specifiability-factor.md) |
| [recorded sensors return objective pass/fail with remediation-bearing output](requirements/sensors-return-pass-fail-with-remediation/sensors-return-pass-fail-with-remediation-requirement.md) | 🔵 Target | ✅ Assessed | [agent-harnessability/self-verifiability](factors/agent-harnessability/factors/self-verifiability/self-verifiability-factor.md) |
| [core service standards are enforced by merge gates, not advisory prose](requirements/standards-gate-nonconforming-changes/standards-gate-nonconforming-changes-requirement.md) | 🟡 Minimum | ✅ Assessed | [agent-harnessability/enforcement-of-standards](factors/agent-harnessability/factors/enforcement-of-standards/enforcement-of-standards-factor.md) |

## Limits and incomplete inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | below-required-margin | Read the overall minimum as a stop-and-fix signal per the model's required margin, not as an acceptable steady state. |

## Primary source data

- [data/evaluation-manifest.json](data/evaluation-manifest.json)
- [data/areas/root/area-analysis-result.json](data/areas/root/area-analysis-result.json)
- [data/advice/finding-ranking-result.json](data/advice/finding-ranking-result.json)
- [data/advice/recommendation-ranking-result.json](data/advice/recommendation-ranking-result.json)
- [data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-rating-result.json](data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-rating-result.json)
- [data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-assessment-result.json](data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-assessment-result.json)
- [data/areas/root/requirements/agents-reach-service-context-from-a-stable-entry-point/requirement-rating-result.json](data/areas/root/requirements/agents-reach-service-context-from-a-stable-entry-point/requirement-rating-result.json)
- [data/areas/root/requirements/agents-reach-service-context-from-a-stable-entry-point/requirement-assessment-result.json](data/areas/root/requirements/agents-reach-service-context-from-a-stable-entry-point/requirement-assessment-result.json)
- [data/areas/root/requirements/consequential-actions-require-approval/requirement-rating-result.json](data/areas/root/requirements/consequential-actions-require-approval/requirement-rating-result.json)
- [data/areas/root/requirements/consequential-actions-require-approval/requirement-assessment-result.json](data/areas/root/requirements/consequential-actions-require-approval/requirement-assessment-result.json)
- [data/areas/root/requirements/handoffs-survive-session-loss/requirement-rating-result.json](data/areas/root/requirements/handoffs-survive-session-loss/requirement-rating-result.json)
- [data/areas/root/requirements/handoffs-survive-session-loss/requirement-assessment-result.json](data/areas/root/requirements/handoffs-survive-session-loss/requirement-assessment-result.json)
- [data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-rating-result.json](data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-rating-result.json)
- [data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-assessment-result.json](data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-assessment-result.json)
- [data/areas/root/requirements/sensors-return-pass-fail-with-remediation/requirement-rating-result.json](data/areas/root/requirements/sensors-return-pass-fail-with-remediation/requirement-rating-result.json)
- [data/areas/root/requirements/sensors-return-pass-fail-with-remediation/requirement-assessment-result.json](data/areas/root/requirements/sensors-return-pass-fail-with-remediation/requirement-assessment-result.json)
- [data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-rating-result.json](data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-rating-result.json)
- [data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-assessment-result.json](data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-assessment-result.json)

