---
type: Requirement Evaluation Report
title: "Requirement: in-flight work survives session loss through durable handoff records"
---

# Requirement: in-flight work survives session loss through durable handoff records

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md)

Factors: [agent-harnessability/continuity](../../factors/agent-harnessability/factors/continuity/continuity-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🔵 Medium / 🔵 Medium |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

Long-running work relies on chat scrollback; durable progress records exist for only one of four recent efforts.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `risk-001` | In-flight decisions live in chat scrollback and would not survive a cold handoff. | ⚠️ Risk | 🟡 Medium | 🔵 Medium | An interruption forces rediscovery of decisions already made, risking repeated or contradicted work. | 🟡 Plausible: The missing records are confirmed; the impact projection assumes a handoff or compaction event occurs mid-effort. |

## Finding details

<a id="finding-risk-001"></a>

### risk-001 In-flight decisions live in chat scrollback and would not survive a cold handoff.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 7 / 21 | 🟡 P3 Medium | Chat-bound progress records only hurt when an interruption lands mid-effort. |

#### Condition

Three of four recent multi-session efforts kept decisions and remaining work only in conversation; one kept a durable progress note.

#### Criteria

- `requirement:root::handoffs-survive-session-loss / rating:target`: In-flight efforts keep durable records of decisions, remaining work, and verification status.

#### Basis

Status: 🟡 Plausible

The missing records are confirmed; the impact projection assumes a handoff or compaction event occurs mid-effort.

##### Basis evidence

(none recorded)

#### Effect

An interruption forces rediscovery of decisions already made, risking repeated or contradicted work.

Rating effect: constrains target

#### Evidence

- `synthetic-source:tracker/in-flight-work`: The traced efforts show tracker entries opened and closed with no interim state; the one durable note enabled a clean resume.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/root/requirements/handoffs-survive-session-loss/requirement-assessment-result.json](../../data/areas/root/requirements/handoffs-survive-session-loss/requirement-assessment-result.json)
- [data/areas/root/requirements/handoffs-survive-session-loss/requirement-rating-result.json](../../data/areas/root/requirements/handoffs-survive-session-loss/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)

