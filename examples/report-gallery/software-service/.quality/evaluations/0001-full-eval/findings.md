---
type: Finding Index Report
title: Findings
---

# Findings

Run: [Run 0001](report.md) - Run ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](report.md) - Findings - [Recommendations](recommendations.md)

## Key Details

| Findings | Highest Severity |
| --- | --- |
| 7 ranked findings | 🔴 High |

Severity: 🔴 Critical, 🔴 High, 🟡 Medium, 🔵 Low.
Empty: `—`.

| Rank | Finding | Area | Factors | Type | Severity |
| --- | --- | --- | --- | --- | --- |
| 1 | [Mutation retry behavior is not fully specified for duplicate idempotency keys.](areas/api/requirements/idempotent-mutations/idempotent-mutations-requirement.md#finding-gap-001) | [Public API](areas/api/api-area.md) | [Correctness](areas/api/factors/correctness/correctness-factor.md) | ⚠️ Gap | 🔴 High |
| 2 | [Common caller error cases share a documented response shape.](areas/api/requirements/predictable-error-contracts/predictable-error-contracts-requirement.md#finding-strength-001) | [Public API](areas/api/api-area.md) | [Operability](areas/api/factors/operability/operability-factor.md) | ✅ Strength | 🔵 Low |
| 3 | [Ledger mutation checks preserve balance invariants in the sampled paths.](areas/persistence/requirements/balance-invariants/balance-invariants-requirement.md#finding-strength-002) | [Ledger Persistence](areas/persistence/persistence-area.md) | [Integrity](areas/persistence/factors/integrity/integrity-factor.md) | ✅ Strength | 🔵 Low |
| 4 | [Rollback guidance exists, but rehearsal evidence is stale.](areas/persistence/requirements/migration-rollback/migration-rollback-requirement.md#finding-risk-001) | [Ledger Persistence](areas/persistence/persistence-area.md) | [Recoverability](areas/persistence/factors/recoverability/recoverability-factor.md) | ⚠️ Risk | 🟡 Medium |
| 5 | [Health dashboards connect service errors to failed customer actions.](areas/operations/requirements/customer-impact-telemetry/customer-impact-telemetry-requirement.md#finding-strength-003) | [Operations](areas/operations/operations-area.md) | [Observability](areas/operations/factors/observability/observability-factor.md) | ✅ Strength | 🔵 Low |
| 6 | [The current owner for ledger recovery drills is ambiguous.](areas/operations/requirements/recovery-drill-ownership/recovery-drill-ownership-requirement.md#finding-unknown-001) | [Operations](areas/operations/operations-area.md) | [Recoverability](areas/operations/factors/recoverability/recoverability-factor.md) | ❓ Unknown | 🟡 Medium |
| 7 | [Agent guidance points to the quality model and generated report path.](areas/agent-harness/requirements/evaluation-entrypoint/evaluation-entrypoint-requirement.md#finding-strength-004) | [Agent Harness](areas/agent-harness/agent-harness-area.md) | [Agent Accessibility](areas/agent-harness/factors/agent-accessibility/agent-accessibility-factor.md) | ✅ Strength | 🔵 Low |

Type: ✅ Strength, ⚠️ Gap, ⚠️ Risk, ❓ Unknown, ℹ️ Note.

## Primary Source Data

- [data/run-manifest.json](data/run-manifest.json)
- [data/advice/finding-ranking-result.json](data/advice/finding-ranking-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json](data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json)
- [data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json](data/areas/api/requirements/predictable-error-contracts/requirement-assessment-result.json)
- [data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json](data/areas/persistence/requirements/balance-invariants/requirement-assessment-result.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json](data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json)
- [data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json](data/areas/operations/requirements/customer-impact-telemetry/requirement-assessment-result.json)
- [data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json](data/areas/operations/requirements/recovery-drill-ownership/requirement-assessment-result.json)
- [data/areas/agent-harness/requirements/evaluation-entrypoint/requirement-assessment-result.json](data/areas/agent-harness/requirements/evaluation-entrypoint/requirement-assessment-result.json)

