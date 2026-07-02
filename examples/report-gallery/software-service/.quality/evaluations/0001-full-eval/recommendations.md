---
type: Recommendation Index Report
title: Recommendations
---

# Recommendations

> **Evaluation links:** [report.md](report.md) | [findings.md](findings.md) | [recommendations.md](recommendations.md) | [glossary.md](../../../glossary.md)

Run: [0001-full-eval](report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

## Key details

| Recommendations | Highest impact | Coverage |
| --- | --- | --- |
| 6 recommendations | ⬥ High | ✅ Addressed by Recommendation: 8 / ⬜ Not Advice Driving: 13 |

## Contents

- [Ranked recommendations](#ranked-recommendations)
- [Coverage](#coverage)
- [Primary source data](#primary-source-data)

## Ranked recommendations

| # | Recommendation | Area / factors | Impact | Confidence | Reason | Ranking rationale |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | [Specify and test replay semantics for interrupted mutations](recommendations/001-specify-and-test-replay-semantics-for-interrupted-mutations.md) | [Public API](areas/api/api-area.md) / [Correctness](areas/api/factors/correctness/correctness-factor.md); [Service contract](areas/service-contract/service-contract-area.md) / [Completeness](areas/service-contract/factors/completeness/completeness-factor.md) | ⬥ High | 🟢 High | Integrators and agents can verify retry behavior from the contract and its sensor instead of inferring undocumented recovery semantics; the API correctness and contract completeness ratings can reach target. | Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements. |
| 2 | [Rehearse migration rollback against the current schema](recommendations/002-rehearse-migration-rollback-against-the-current-schema.md) | [Ledger persistence](areas/persistence/persistence-area.md) / [Recoverability](areas/persistence/factors/recoverability/recoverability-factor.md) | ⬥ High | 🔵 Medium | Release risk drops because rollback instructions are proven against the schema they would actually run on, and the recoverability rating can return to target. | Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements. |
| 3 | [Make the contract-test and invariant sensors required at merge](recommendations/003-make-the-contract-test-and-invariant-sensors-required-at-merge.md) | [LedgerLite Service](root-area.md) / [Enforcement of Standards](factors/agent-harnessability/factors/enforcement-of-standards/enforcement-of-standards-factor.md) | ⬥ High | 🔵 Medium | Contract conformance and ledger invariants hold regardless of reviewer attention, converting the two strongest sensors from advisory signals into enforced standards. | Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements. |
| 4 | [Reconcile recovery drill ownership and restore assessability](recommendations/004-reconcile-recovery-drill-ownership-and-restore-assessability.md) | [Operations](areas/operations/operations-area.md) / [Recoverability](areas/operations/factors/recoverability/recoverability-factor.md) | ● Medium | 🟢 High | The recoverability factor becomes assessable again, and the next evaluation can rate drill practice on evidence instead of recording missing evidence. | Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements. |
| 5 | [Record done criteria and progress in durable handoff notes](recommendations/005-record-done-criteria-and-progress-in-durable-handoff-notes.md) | [LedgerLite Service](root-area.md) / [Task Specifiability](factors/agent-harnessability/factors/task-specifiability/task-specifiability-factor.md), [Continuity](factors/agent-harnessability/factors/continuity/continuity-factor.md) | ● Medium | 🔵 Medium | Agents can declare completion against criteria and resume interrupted work from records, lifting task specifiability and continuity toward target. | Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements. |
| 6 | [Refresh the model body's unknowns and open questions](recommendations/006-refresh-the-model-body-s-unknowns-and-open-questions.md) | [LedgerLite Service QUALITY.md](areas/quality-md/quality-md-area.md) / [Context Grounding](areas/quality-md/factors/context-grounding/context-grounding-factor.md), [Evaluability](areas/quality-md/factors/evaluability/evaluability-factor.md), [Lifecycle Maintenance](areas/quality-md/factors/lifecycle-maintenance/lifecycle-maintenance-factor.md) | ○ Low | 🟢 High | Evaluations start from current judgment context, and the QUALITY.md self-check can return to target. | Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements. |

## Coverage

- ✅ Addressed by Recommendation: 8
- ⬜ Not Advice Driving: 13

## Primary source data

- [data/evaluation-manifest.json](data/evaluation-manifest.json)
- [data/advice/recommendation-ranking-result.json](data/advice/recommendation-ranking-result.json)
- [data/advice/recommendations/qrec_replaycontract/recommendation-result.json](data/advice/recommendations/qrec_replaycontract/recommendation-result.json)
- [data/advice/recommendations/qrec_rollbackrehearsal/recommendation-result.json](data/advice/recommendations/qrec_rollbackrehearsal/recommendation-result.json)
- [data/advice/recommendations/qrec_gatesensors/recommendation-result.json](data/advice/recommendations/qrec_gatesensors/recommendation-result.json)
- [data/advice/recommendations/qrec_drillownership/recommendation-result.json](data/advice/recommendations/qrec_drillownership/recommendation-result.json)
- [data/advice/recommendations/qrec_durablehandoffs/recommendation-result.json](data/advice/recommendations/qrec_durablehandoffs/recommendation-result.json)
- [data/advice/recommendations/qrec_refreshmodelbody/recommendation-result.json](data/advice/recommendations/qrec_refreshmodelbody/recommendation-result.json)

