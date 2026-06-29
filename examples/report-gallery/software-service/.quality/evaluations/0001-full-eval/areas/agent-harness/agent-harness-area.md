---
type: Area Evaluation Report
title: Agent Harness
data:
  - data/evaluation-output-result.json
  - data/areas/agent-harness/area-analysis-result.json
---

# Area: Agent Harness

Run: [#1](../../report.md) - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../report.md) - [Findings](../../findings.md) - [Recommendations](../../recommendations.md)

Area: [LedgerLite Service](../../root-area.md) / [Agent Harness](agent-harness-area.md)

| Overall Rating | Local Rating | Confidence | Data |
| --- | --- | --- | --- |
| 🔵 Target | 🔵 Target | 🟢 High / 🟢 High | [area-analysis-result.json](../../data/areas/agent-harness/area-analysis-result.json) |

Summary:

Agent guidance exposes the quality evaluation entry point clearly.

## Rating Drivers

| Driver | Effect | Inputs |
| --- | --- | --- |
| Agent Accessibility is driven by agent guidance routes quality evaluation work. | supports target | [{"kind":"FactorAnalysisResult","selector":"localAndDescendantAnalysis","subject":{"factorId":"factor:agent-harness::agent-accessibility"}}] |

## Factors

| Factor | Path | Local Rating | + Sub-Factors Rating | Sub-Factors |
| --- | --- | --- | --- | --- |
| [Agent Accessibility](factors/agent-accessibility/agent-accessibility-factor.md) | `agent-harness::agent-accessibility` | 🔵 Target | — | — |

## Child Areas

| Area | Path | Local Rating | + Child Areas Rating | Factors |
| --- | --- | --- | --- | --- |
| (no Child Areas) |  |  |  |  |

## Requirements

| Requirement | Rating | Status | Factors |
| --- | --- | --- | --- |
| [agent guidance routes quality evaluation work](requirements/evaluation-entrypoint/evaluation-entrypoint-requirement.md) | 🔵 Target | ✅ Assessed | [agent-accessibility](factors/agent-accessibility/agent-accessibility-factor.md) |

## Limits & Incomplete Inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |

## Legend

- `—` - not applicable or not recorded.
