---
type: Factor Evaluation Report
title: "Factor: Correctness"
---

# Factor: Correctness

Run: [Run 0001](../../../../report.md) - Run ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Report: [Overview](../../../../report.md) - [Findings](../../../../findings.md) - [Recommendations](../../../../recommendations.md)

Area: [LedgerLite Service](../../../../root-area.md) / [Public API](../../api-area.md)

Factor: [Correctness](correctness-factor.md)

## Key Details

| Overall Rating | Local Rating | Status | Confidence |
| --- | --- | --- | --- |
| 🟡 Minimum | 🟡 Minimum | ✅ Analyzed / ✅ Analyzed | 🔵 Medium / 🔵 Medium |

Legend

- *Quality rating:* 🟢 Outstanding, 🔵 Target, 🟡 Minimum, 🔴 Unacceptable
- *Analysis status:* ✅ Analyzed, ⬜ Empty, ⚪ Not Analyzed, ⛔ Blocked
- *Confidence:* 🟢 High, 🔵 Medium, 🟡 Low, ⚪ None
- *Empty:* `—`

## Contents

- [Summary](#summary)
- [Requirements](#requirements)
- [Sub-Factors](#sub-factors)
- [Limits & Incomplete Inputs](#limits--incomplete-inputs)
- [Primary Source Data](#primary-source-data)

## Summary

Correctness follows its direct requirement signal.

## Requirements

| Requirement | Rating | Status |
| --- | --- | --- |
| [mutation endpoints are idempotent under retry](../../requirements/idempotent-mutations/idempotent-mutations-requirement.md) | 🟡 Minimum | ✅ Assessed |

Legend

- *Assessment status:* ✅ Assessed, 🟡 Partially Assessed, ⚪ Not Assessed, ⛔ Blocked
- *Empty:* `—`

## Sub-Factors

| Factor | Path | Local Rating | + Sub-Factors Rating |
| --- | --- | --- | --- |
| (no Sub-Factors) | — | — | — |

Legend

- *Empty:* `—`

## Limits & Incomplete Inputs

| Type | Scope | Impact |
| --- | --- | --- |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |
| ⚠️ Evaluation Limits | synthetic-evaluation | Use for report design and example browsing only. |

## Primary Source Data

- [data/run-manifest.json](../../../../data/run-manifest.json)
- [data/areas/api/factors/correctness/factor-analysis-result.json](../../../../data/areas/api/factors/correctness/factor-analysis-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json](../../../../data/areas/api/requirements/idempotent-mutations/requirement-rating-result.json)
- [data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json](../../../../data/areas/api/requirements/idempotent-mutations/requirement-assessment-result.json)

