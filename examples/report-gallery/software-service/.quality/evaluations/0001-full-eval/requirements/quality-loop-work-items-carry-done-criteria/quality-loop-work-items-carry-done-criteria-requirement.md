---
type: Requirement Evaluation Report
title: "Requirement: quality-loop work items carry scoped goals and done criteria"
---

# Requirement: quality-loop work items carry scoped goals and done criteria

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md)

Factors: [agent-harnessability/task-specifiability](../../factors/agent-harnessability/factors/task-specifiability/task-specifiability-factor.md)

## Key details

| Rating | Assessment | Confidence |
| --- | --- | --- |
| 🟡 Minimum | ✅ Assessed | 🟢 High / 🟢 High |

## Contents

- [Summary](#summary)
- [Findings summary](#findings-summary)
- [Finding details](#finding-details)
- [Unknowns and missing evidence](#unknowns-and-missing-evidence)
- [Primary source data](#primary-source-data)

## Summary

Recent handoffs state goals, but most omit done criteria and the confirming sensor.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `gap-001` | Most quality-loop handoffs omit done criteria and the confirming sensor. | 🚩 Gap | 🟡 Medium | 🟢 High | Agents declare completion on judgment rather than criteria, which constrains task specifiability to minimum. | ✅ Verified: The handoff sample was read in full; the omissions are countable. |

## Finding details

<a id="finding-gap-001"></a>

### gap-001 Most quality-loop handoffs omit done criteria and the confirming sensor.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 6 / 21 | 🟡 P3 Medium | Missing done criteria degrade every agent handoff, but each instance is recoverable. |

#### Condition

Nine of the twelve most recent handoffs state a goal but no done criteria; none name the sensor that would confirm completion.

#### Criteria

- `requirement:root::quality-loop-work-items-carry-done-criteria / rating:target`: Sampled handoffs state goals, non-goals, done criteria, and the sensor that confirms completion.
  Rationale: Without done criteria an agent must guess when work is finished.

#### Basis

Status: ✅ Verified

The handoff sample was read in full; the omissions are countable.

##### Basis evidence

(none recorded)

#### Effect

Agents declare completion on judgment rather than criteria, which constrains task specifiability to minimum.

Rating effect: constrains target

#### Evidence

- `synthetic-source:tracker/quality-loop-handoffs`: The sampled handoffs contain goal statements; done criteria appear in three, a confirming sensor in none.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-assessment-result.json](../../data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-assessment-result.json)
- [data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-rating-result.json](../../data/areas/root/requirements/quality-loop-work-items-carry-done-criteria/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)

