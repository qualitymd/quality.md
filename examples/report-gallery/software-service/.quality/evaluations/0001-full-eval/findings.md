---
type: Finding Index Report
title: Findings
---

# Findings

> **Evaluation links:** [report.md](report.md) | [findings.md](findings.md) | [recommendations.md](recommendations.md) | [glossary.md](../../../glossary.md)

Run: [0001-full-eval](report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key details

| Findings | Highest concern severity |
| --- | --- |
| 21 findings | 🔴 High |

## Contents

- [Ranked findings](#ranked-findings)
- [Primary source data](#primary-source-data)

## Ranked findings

| Rank | Finding | Area | Factors | Type | Severity |
| --- | --- | --- | --- | --- | --- |
| 1 | [Replay behavior after an interrupted write is unspecified and untested.](areas/api/requirements/idempotent-mutations/idempotent-mutations-requirement.md#finding-gap-004) | [Public API](areas/api/api-area.md) | [Correctness](areas/api/factors/correctness/correctness-factor.md) | 🚩 Gap | 🔴 High |
| 2 | [Rollback steps have not been rehearsed since the two most recent schema migrations landed.](areas/persistence/requirements/migration-rollback/migration-rollback-requirement.md#finding-risk-002) | [Ledger persistence](areas/persistence/persistence-area.md) | [Recoverability](areas/persistence/factors/recoverability/recoverability-factor.md) | ⚠️ Risk | 🟡 Medium |
| 3 | [Contract tests and the invariant suite run at merge time but cannot block a merge.](requirements/standards-gate-nonconforming-changes/standards-gate-nonconforming-changes-requirement.md#finding-gap-003) | [LedgerLite Service](root-area.md) | [Enforcement of Standards](factors/agent-harnessability/factors/enforcement-of-standards/enforcement-of-standards-factor.md) | 🚩 Gap | 🟡 Medium |
| 4 | [Two mutation endpoints have no replay semantics in the contract.](areas/service-contract/requirements/contract-covers-mutation-semantics/contract-covers-mutation-semantics-requirement.md#finding-gap-005) | [Service contract](areas/service-contract/service-contract-area.md) | [Completeness](areas/service-contract/factors/completeness/completeness-factor.md) | 🚩 Gap | 🟡 Medium |
| 5 | [Recovery drill ownership records contradict each other.](areas/operations/requirements/recovery-drill-ownership/recovery-drill-ownership-requirement.md#finding-note-003) | [Operations](areas/operations/operations-area.md) | [Recoverability](areas/operations/factors/recoverability/recoverability-factor.md) | ℹ️ Note | — |
| 6 | [Most quality-loop handoffs omit done criteria and the confirming sensor.](requirements/quality-loop-work-items-carry-done-criteria/quality-loop-work-items-carry-done-criteria-requirement.md#finding-gap-001) | [LedgerLite Service](root-area.md) | [Task Specifiability](factors/agent-harnessability/factors/task-specifiability/task-specifiability-factor.md) | 🚩 Gap | 🟡 Medium |
| 7 | [In-flight decisions live in chat scrollback and would not survive a cold handoff.](requirements/handoffs-survive-session-loss/handoffs-survive-session-loss-requirement.md#finding-risk-001) | [LedgerLite Service](root-area.md) | [Continuity](factors/agent-harnessability/factors/continuity/continuity-factor.md) | ⚠️ Risk | 🟡 Medium |
| 8 | [The body's unknowns and open questions have not been revisited since the service-contract area was added.](areas/quality-md/requirements/the-model-follows-the-authoring-guide-family/the-model-follows-the-authoring-guide-family-requirement.md#finding-gap-006) | [LedgerLite Service QUALITY.md](areas/quality-md/quality-md-area.md) | [Context Grounding](areas/quality-md/factors/context-grounding/context-grounding-factor.md), [Evaluability](areas/quality-md/factors/evaluability/evaluability-factor.md), [Lifecycle Maintenance](areas/quality-md/factors/lifecycle-maintenance/lifecycle-maintenance-factor.md) | 🚩 Gap | 🟡 Medium |
| 9 | [The telemetry read credential needed for latency checks is documented only in the private wiki.](requirements/a-fresh-session-reaches-a-ready-to-work-environment/a-fresh-session-reaches-a-ready-to-work-environment-requirement.md#finding-gap-002) | [LedgerLite Service](root-area.md) | [Agent Operability](factors/agent-harnessability/factors/agent-operability/agent-operability-factor.md) | 🚩 Gap | 🔵 Low |
| 10 | [The sensor catalog's invariant-suite command name is one rename behind.](areas/agent-harness/requirements/harness-orients-agents-and-routes-to-sensors/harness-orients-agents-and-routes-to-sensors-requirement.md#finding-note-004) | [Agent harness](areas/agent-harness/agent-harness-area.md) | [Completeness](areas/agent-harness/factors/completeness/completeness-factor.md), [Coherence](areas/agent-harness/factors/coherence/coherence-factor.md), [Currentness](areas/agent-harness/factors/currentness/currentness-factor.md), [Assessability](areas/agent-harness/factors/assessability/assessability-factor.md) | ℹ️ Note | — |
| 11 | [Replay traffic is a meaningful share of mutation volume.](areas/api/requirements/idempotent-mutations/idempotent-mutations-requirement.md#finding-note-001) | [Public API](areas/api/api-area.md) | [Correctness](areas/api/factors/correctness/correctness-factor.md) | ℹ️ Note | — |
| 12 | [The deprecated balance_after field still ships but is no longer documented.](areas/service-contract/requirements/contract-matches-shipped-behavior/contract-matches-shipped-behavior-requirement.md#finding-note-002) | [Service contract](areas/service-contract/service-contract-area.md) | [Consistency](areas/service-contract/factors/consistency/consistency-factor.md) | ℹ️ Note | — |
| 13 | [Balance invariants hold across every tested mutation path with zero reconciliation drift.](areas/persistence/requirements/balance-invariants/balance-invariants-requirement.md#finding-strength-007) | [Ledger persistence](areas/persistence/persistence-area.md) | [Integrity](areas/persistence/factors/integrity/integrity-factor.md) | 💪 Strength | — |
| 14 | [All three recorded sensors are deterministic and their failures name the violated expectation.](requirements/sensors-return-pass-fail-with-remediation/sensors-return-pass-fail-with-remediation-requirement.md#finding-strength-002) | [LedgerLite Service](root-area.md) | [Self-Verifiability](factors/agent-harnessability/factors/self-verifiability/self-verifiability-factor.md) | 💪 Strength | — |
| 15 | [Money movement and production schema changes are behind declarative approval gates.](requirements/consequential-actions-require-approval/consequential-actions-require-approval-requirement.md#finding-strength-003) | [LedgerLite Service](root-area.md) | [Containment of Action](factors/agent-harnessability/factors/containment-of-action/containment-of-action-factor.md) | 💪 Strength | — |
| 16 | [The agent entry point reaches the contract, runbooks, and sensor catalog in one hop each.](requirements/agents-reach-service-context-from-a-stable-entry-point/agents-reach-service-context-from-a-stable-entry-point-requirement.md#finding-strength-001) | [LedgerLite Service](root-area.md) | [Agent Accessibility](factors/agent-harnessability/factors/agent-accessibility/agent-accessibility-factor.md) | 💪 Strength | — |
| 17 | [Error responses follow one documented envelope across all sampled failure modes.](areas/api/requirements/predictable-error-contracts/predictable-error-contracts-requirement.md#finding-strength-004) | [Public API](areas/api/api-area.md) | [Operability](areas/api/factors/operability/operability-factor.md) | 💪 Strength | — |
| 18 | [Mutation p99 measured 262 ms over the representative four-week window.](areas/api/requirements/p99-latency-within-budget/p99-latency-within-budget-requirement.md#finding-strength-005) | [Public API](areas/api/api-area.md) | [Performance](areas/api/factors/performance/performance-factor.md) | 💪 Strength | — |
| 19 | [Shipped behavior matches every specified contract clause in the latest sensor run.](areas/service-contract/requirements/contract-matches-shipped-behavior/contract-matches-shipped-behavior-requirement.md#finding-strength-006) | [Service contract](areas/service-contract/service-contract-area.md) | [Consistency](areas/service-contract/factors/consistency/consistency-factor.md) | 💪 Strength | — |
| 20 | [Health dashboards express failed customer actions, not just service symptoms.](areas/operations/requirements/customer-impact-telemetry/customer-impact-telemetry-requirement.md#finding-strength-008) | [Operations](areas/operations/operations-area.md) | [Observability](areas/operations/factors/observability/observability-factor.md) | 💪 Strength | — |
| 21 | [The harness covers the agent work lifecycle and its guidance agrees with the material it routes to.](areas/agent-harness/requirements/harness-orients-agents-and-routes-to-sensors/harness-orients-agents-and-routes-to-sensors-requirement.md#finding-strength-009) | [Agent harness](areas/agent-harness/agent-harness-area.md) | [Completeness](areas/agent-harness/factors/completeness/completeness-factor.md), [Coherence](areas/agent-harness/factors/coherence/coherence-factor.md), [Currentness](areas/agent-harness/factors/currentness/currentness-factor.md), [Assessability](areas/agent-harness/factors/assessability/assessability-factor.md) | 💪 Strength | — |

## Primary source data

- [data/evaluation-manifest.json](data/evaluation-manifest.json)
- [data/advice/finding-ranking-result.json](data/advice/finding-ranking-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json](data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json](data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json)
- [data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-assessment-result.json](data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-assessment-result.json)
- [data/areas/service-contract/requirements/contract-covers-mutation-semantics/requirement-assessment-result.json](data/areas/service-contract/requirements/contract-covers-mutation-semantics/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)
- [data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-assessment-result.json](data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-assessment-result.json)
- [data/areas/root/requirements/handoffs-survive-session-loss/requirement-assessment-result.json](data/areas/root/requirements/handoffs-survive-session-loss/requirement-assessment-result.json)
- [data/areas/quality-md/requirements/the-model-follows-the-authoring-guide-family/requirement-assessment-result.json](data/areas/quality-md/requirements/the-model-follows-the-authoring-guide-family/requirement-assessment-result.json)
- [data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-assessment-result.json](data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-assessment-result.json)
- [data/areas/agent-harness/requirements/harness-orients-agents-and-routes-to-sensors/requirement-assessment-result.json](data/areas/agent-harness/requirements/harness-orients-agents-and-routes-to-sensors/requirement-assessment-result.json)
- [data/areas/service-contract/requirements/contract-matches-shipped-behavior/requirement-assessment-result.json](data/areas/service-contract/requirements/contract-matches-shipped-behavior/requirement-assessment-result.json)
- [data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json](data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json)
- [data/areas/root/requirements/sensors-return-pass-fail-with-remediation/requirement-assessment-result.json](data/areas/root/requirements/sensors-return-pass-fail-with-remediation/requirement-assessment-result.json)
- [data/areas/root/requirements/consequential-actions-require-approval/requirement-assessment-result.json](data/areas/root/requirements/consequential-actions-require-approval/requirement-assessment-result.json)
- [data/areas/root/requirements/agents-reach-service-context-from-a-stable-entry-point/requirement-assessment-result.json](data/areas/root/requirements/agents-reach-service-context-from-a-stable-entry-point/requirement-assessment-result.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json](data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json)
- [data/areas/api/requirements/p99-latency-within-budget/requirement-assessment-result.json](data/areas/api/requirements/p99-latency-within-budget/requirement-assessment-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json](data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json)

