---
type: Requirement Evaluation Report
title: "Requirement: the harness orients agents and routes them to runnable sensors"
---

# Requirement: the harness orients agents and routes them to runnable sensors

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Agent harness](../../agent-harness-area.md)

Factors: [completeness](../../factors/completeness/completeness-factor.md); [coherence](../../factors/coherence/coherence-factor.md); [currentness](../../factors/currentness/currentness-factor.md); [assessability](../../factors/assessability/assessability-factor.md)

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

The entry point, routed guidance, and sensor catalog cover setup through handoff and agree with the contract and runbooks; one sensor command name is stale.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-009` | The harness covers the agent work lifecycle and its guidance agrees with the material it routes to. | 💪 Strength | — | 🟢 High | Agents get complete, coherent steering, supporting the target rating across the harness factors. | ✅ Verified: All harness artifacts were read; routed references were followed to their targets. |
| `note-004` | The sensor catalog's invariant-suite command name is one rename behind. | ℹ️ Note | — | 🟢 High | A fresh agent following the catalog hits a not-found error on the invariant suite until it discovers the rename. | ✅ Verified: Both the recorded and current command names were executed. |

## Finding details

<a id="finding-strength-009"></a>

### strength-009 The harness covers the agent work lifecycle and its guidance agrees with the material it routes to.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 21 / 21 | ⚪ P4 Low | Harness coverage and coherence hold up under inspection. |

#### Condition

The entry point, routed guides, and sensor catalog address setup, scoped work, verification, and handoff, and spot-checks against the contract and runbooks found no contradictions.

#### Criteria

- `requirement:agent-harness::harness-orients-agents-and-routes-to-sensors / rating:target`: Harness artifacts cover setup, scoped work, verification, and handoff; agree with the contract and runbooks; and route to sensors that run as recorded.

#### Basis

Status: ✅ Verified

All harness artifacts were read; routed references were followed to their targets.

##### Basis evidence

(none recorded)

#### Effect

Agents get complete, coherent steering, supporting the target rating across the harness factors.

Rating effect: supports target

#### Evidence

- `synthetic-source:agent-harness`: The lifecycle coverage table maps each work phase to a harness artifact; contradiction spot-checks across contract and runbook references came back clean.

<a id="finding-note-004"></a>

### note-004 The sensor catalog's invariant-suite command name is one rename behind.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 10 / 21 | ⚪ P4 Low | The stale sensor command name is a one-line catalog fix. |

#### Condition

The catalog records the pre-rename command; the current command name differs, and the old name fails with a not-found error rather than a redirect.

#### Criteria

- `requirement:agent-harness::harness-orients-agents-and-routes-to-sensors / rating:target`: Harness artifacts cover setup, scoped work, verification, and handoff; agree with the contract and runbooks; and route to sensors that run as recorded.
  Rationale: A stale command turns a routed sensor into a dead end for a fresh session.

#### Basis

Status: ✅ Verified

Both the recorded and current command names were executed.

##### Basis evidence

(none recorded)

#### Effect

A fresh agent following the catalog hits a not-found error on the invariant suite until it discovers the rename.

Rating effect: informational

#### Evidence

- `synthetic-source:agent-harness/sensor-catalog`: The catalog entry's command fails with not-found; the renamed command runs the suite successfully.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/agent-harness/requirements/harness-orients-agents-and-routes-to-sensors/requirement-assessment-result.json](../../../../data/areas/agent-harness/requirements/harness-orients-agents-and-routes-to-sensors/requirement-assessment-result.json)
- [data/areas/agent-harness/requirements/harness-orients-agents-and-routes-to-sensors/requirement-rating-result.json](../../../../data/areas/agent-harness/requirements/harness-orients-agents-and-routes-to-sensors/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)

