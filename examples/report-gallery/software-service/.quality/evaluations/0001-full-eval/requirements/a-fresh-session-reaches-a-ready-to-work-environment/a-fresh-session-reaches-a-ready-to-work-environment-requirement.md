---
type: Requirement Evaluation Report
title: "Requirement: a fresh agent session reaches a ready-to-work environment from recorded setup"
---

# Requirement: a fresh agent session reaches a ready-to-work environment from recorded setup

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md)

Factors: [agent-harnessability/agent-operability](../../factors/agent-harnessability/factors/agent-operability/agent-operability-factor.md)

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

Recorded setup builds the service and runs the sensors; one required credential is documented only in the private wiki.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `gap-002` | The telemetry read credential needed for latency checks is documented only in the private wiki. | 🚩 Gap | 🔵 Low | 🔵 Medium | Latency verification needs a human hand-off in an otherwise self-serve setup, constraining operability to minimum. | ✅ Verified: The setup walk-through stalled at the credential step; every other step completed from recorded materials. |

## Finding details

<a id="finding-gap-002"></a>

### gap-002 The telemetry read credential needed for latency checks is documented only in the private wiki.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 9 / 42 | 🟡 P3 Medium | One wiki-bound credential is a small, well-understood setup snag. |

#### Condition

The recorded setup covers build and sensors, but the step granting telemetry read access points at a wiki page agents cannot reach.

#### Criteria

- `requirement:root::a-fresh-session-reaches-a-ready-to-work-environment / rating:target`: A fresh session reaches build, sensors, and required access from agent-accessible recorded materials alone.

#### Basis

Status: ✅ Verified

The setup walk-through stalled at the credential step; every other step completed from recorded materials.

##### Basis evidence

(none recorded)

#### Effect

Latency verification needs a human hand-off in an otherwise self-serve setup, constraining operability to minimum.

Rating effect: constrains target

#### Evidence

- `synthetic-source:agent-harness/setup`: The setup document's access section links the private wiki for the telemetry credential; no agent-accessible copy exists.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-assessment-result.json](../../data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-assessment-result.json)
- [data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-rating-result.json](../../data/areas/root/requirements/a-fresh-session-reaches-a-ready-to-work-environment/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
