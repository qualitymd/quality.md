---
type: Requirement Evaluation Report
title: "Requirement: the v1 deprecation window is current and visible"
---

# Requirement: the v1 deprecation window is current and visible

> **Evaluation links:** [report.md](../../../../report.md) | [findings.md](../../../../findings.md) | [recommendations.md](../../../../recommendations.md) | [glossary.md](../../../../../../../glossary.md)

Run: [0001-full-eval](../../../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../../../root-area.md) / [Service contract](../../service-contract-area.md)

Factors: [currentness](../../factors/currentness/currentness-factor.md); [understandability](../../factors/understandability/understandability-factor.md)

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

The compatibility appendix names the active deprecation window, but the drift detector found one shipped deprecated field missing from the appendix.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `note-005` | The compatibility appendix names the active deprecation window, but the drift detector found one shipped deprecated field missing from the appendix. | ℹ️ Note | — | 🔵 Medium | The evidence leaves a visible gap, constraining the result below target while remaining acceptable. | ✅ Verified: The drift detector compared the contract appendix with the handler matrix. |

## Finding details

<a id="finding-note-005"></a>

### note-005 The compatibility appendix names the active deprecation window, but the drift detector found one shipped deprecated field missing from the appendix.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 28 / 42 | ⚪ P4 Low | Routine generated finding included for complete finding-ranking and advice-coverage accounting. |

#### Condition

Does the contract's deprecation appendix match shipped v1 fields?

#### Criteria

- `requirement:service-contract::deprecation-window-is-current / rating:target`: Falls short of target with visible gaps or limited evidence, while remaining acceptable.
  Rationale: This generated gallery finding demonstrates the assessment evidence shape for this requirement.

#### Basis

Status: ✅ Verified

The drift detector compared the contract appendix with the handler matrix.

##### Basis evidence

(none recorded)

#### Effect

The evidence leaves a visible gap, constraining the result below target while remaining acceptable.

Rating effect: constrains target

#### Evidence

- `synthetic-source:service-contract`: The compatibility appendix names the active deprecation window, but the drift detector found one shipped deprecated field missing from the appendix.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../../../data/evaluation-manifest.json)
- [data/areas/service-contract/requirements/deprecation-window-is-current/requirement-assessment-result.json](../../../../data/areas/service-contract/requirements/deprecation-window-is-current/requirement-assessment-result.json)
- [data/areas/service-contract/requirements/deprecation-window-is-current/requirement-rating-result.json](../../../../data/areas/service-contract/requirements/deprecation-window-is-current/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../../../data/advice/finding-ranking-result.json)
