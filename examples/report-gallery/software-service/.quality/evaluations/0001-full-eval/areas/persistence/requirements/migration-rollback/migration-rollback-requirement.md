---
type: Requirement Evaluation Report
title: "Requirement: migrations have rehearsed rollback paths"
---

# Requirement: migrations have rehearsed rollback paths

Run: [Run 0001](../../../../report.md) - Run ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../../../report.md) - [Findings](../../../../findings.md) - [Recommendations](../../../../recommendations.md)

Area: [LedgerLite Service](../../../../root-area.md) / [Ledger Persistence](../../persistence-area.md)

Factors: [recoverability](../../factors/recoverability/recoverability-factor.md)

## Key Details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🔵 Medium / 🔵 Medium |

Ratings: 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable.
Assessment: ✅ Assessed, 🟡 Partially Assessed, ⚪ Not Assessed, ⛔ Blocked.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.
Empty: `—`.

Jump to: [Findings Summary](#findings-summary) - [Finding Details](#finding-details) - [Unknowns & Missing Evidence](#unknowns--missing-evidence)

## Summary

Rollback instructions are present, but the rehearsal signal is stale.

## Findings Summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `risk-001` | Rollback guidance exists, but rehearsal evidence is stale. | ⚠️ Risk | 🟡 Medium | 🔵 Medium | The finding constrains recoverability to minimum until rollback rehearsal is refreshed. | 🟡 Plausible: Stale rehearsal evidence plausibly misses drift in newer migration behavior. |

Type: ✅ Strength, ⚠️ Gap, ⚠️ Risk, ❓ Unknown, ℹ️ Note.
Severity: 🔴 Critical, 🔴 High, 🟡 Medium, 🔵 Low.
Confidence: 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None.
Basis: ✅ Verified, 🟡 Plausible, ⚪ Not Assessed, ⬜ Not Applicable.
Empty: `—`.

## Finding Details

<a id="finding-risk-001"></a>

### risk-001 Rollback guidance exists, but rehearsal evidence is stale.

| Advice Rank | Tier | Ranking Rationale |
| --- | --- | --- |
| 4 / 7 | 🔴 P1 Highest | Ranked by expected impact on the service quality bar and report-gallery usefulness. |

Tier: 🔴 P1 Highest, 🟠 P2 High, 🟡 P3 Medium, ⚪ P4 Low.
Empty: `—`.

#### Condition

The synthetic migration runbook names rollback steps, but the last recorded rehearsal predates two schema changes.

#### Criteria

- `requirement:persistence::migration-rollback / rating:target`: Migration rollback paths should meet the target recoverability bar with current rehearsal evidence.
  Rationale: The gallery records one finding per requirement so report tables stay easy to inspect.

#### Basis

Status: 🟡 Plausible

Stale rehearsal evidence plausibly misses drift in newer migration behavior.

##### Basis Evidence

(none recorded)

#### Effect

The finding constrains recoverability to minimum until rollback rehearsal is refreshed.

Rating effect: constrains target

#### Evidence

- `synthetic-source:persistence/migration-runbook`: The synthetic runbook contains rollback steps and a rehearsal note older than the latest two migrations.
  Rationale: Synthetic source reference retained to demonstrate evidence rendering.

## Unknowns & Missing Evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary Source Data

- [data/run-manifest.json](../../../../data/run-manifest.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json](../../../../data/areas/persistence/requirements/migration-rollback/requirement-assessment-result.json)
- [data/areas/persistence/requirements/migration-rollback/requirement-rating-result.json](../../../../data/areas/persistence/requirements/migration-rollback/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

