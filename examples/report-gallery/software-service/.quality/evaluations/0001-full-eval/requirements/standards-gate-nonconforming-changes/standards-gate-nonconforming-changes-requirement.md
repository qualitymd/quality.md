---
type: Requirement Evaluation Report
title: "Requirement: core service standards are enforced by merge gates, not advisory prose"
---

# Requirement: core service standards are enforced by merge gates, not advisory prose

> **Evaluation links:** [report.md](../../report.md) | [findings.md](../../findings.md) | [recommendations.md](../../recommendations.md) | [glossary.md](../../../../../glossary.md)

Run: [0001-full-eval](../../report.md) - Evaluation ID: `20260629T120000Z-0123456789ab` - Created: 2026-06-29T12:00:00Z - Scope: full evaluation

Area: [LedgerLite Service](../../root-area.md)

Factors: [agent-harnessability/enforcement-of-standards](../../factors/agent-harnessability/factors/enforcement-of-standards/enforcement-of-standards-factor.md)

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

Lint gates merges; the contract tests and invariant suite run on merge but are advisory.

## Findings summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `gap-003` | Contract tests and the invariant suite run at merge time but cannot block a merge. | 🚩 Gap | 🟡 Medium | 🟢 High | Contract conformance depends on reviewer diligence rather than the gate, constraining enforcement to minimum. | ✅ Verified: Pipeline configuration and the two failing-but-merged runs were inspected. |

## Finding details

<a id="finding-gap-003"></a>

### gap-003 Contract tests and the invariant suite run at merge time but cannot block a merge.

| Advice rank | Tier | Ranking rationale |
| --- | --- | --- |
| 3 / 42 | 🟠 P2 High | Advisory merge gates let the two strongest sensors be ignored under pressure. |

#### Condition

The merge pipeline marks both sensor jobs as non-required; two merged changes in the sampled month had a failing contract-test run.

#### Criteria

- `requirement:root::standards-gate-nonconforming-changes / rating:target`: Contract tests, the invariant suite, and lint block nonconforming merges or route them through reviewable exceptions.
  Rationale: An advisory sensor protects nothing when the pressure is on; the gate is what makes the standard hold regardless of who is merging.

#### Basis

Status: ✅ Verified

Pipeline configuration and the two failing-but-merged runs were inspected.

##### Basis evidence

(none recorded)

#### Effect

Contract conformance depends on reviewer diligence rather than the gate, constraining enforcement to minimum.

Rating effect: constrains target

#### Evidence

- `synthetic-source:ci/merge-pipeline`: The pipeline config lists lint as required and the contract-test and invariant jobs as informational; merge history shows two failing contract-test runs that merged.

## Unknowns and missing evidence

| Type | Detail |
| --- | --- |
| (none recorded) | — |

## Primary source data

- [data/evaluation-manifest.json](../../data/evaluation-manifest.json)
- [data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-assessment-result.json](../../data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-assessment-result.json)
- [data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-rating-result.json](../../data/areas/root/requirements/standards-gate-nonconforming-changes/requirement-rating-result.json)
- [data/advice/finding-ranking-result.json](../../data/advice/finding-ranking-result.json)
