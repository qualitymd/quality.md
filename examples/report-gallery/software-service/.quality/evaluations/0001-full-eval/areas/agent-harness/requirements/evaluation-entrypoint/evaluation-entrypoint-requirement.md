# Requirement: agent guidance routes quality evaluation work

Area: [LedgerLite Service](../../../../root-area.md) / [Agent Harness](../../agent-harness-area.md)

Factors: [agent-accessibility](../../factors/agent-accessibility/agent-accessibility-factor.md)

| Rating | Assessment | Confidence | Data |
| --- | --- | --- | --- |
| 🔵 Target | ✅ Assessed | 🟢 High / 🟢 High | [requirement-assessment-result.json](../../../../data/areas/agent-harness/requirements/evaluation-entrypoint/requirement-assessment-result.json), [requirement-rating-result.json](../../../../data/areas/agent-harness/requirements/evaluation-entrypoint/requirement-rating-result.json) |

Summary:

The agent-facing entry point is clear enough for quality evaluation work.

## Findings Summary

| ID | Statement | Type | Severity | Confidence | Effect | Basis |
| --- | --- | --- | --- | --- | --- | --- |
| `strength-004` | Agent guidance points to the quality model and generated report path. | ✅ Strength | 🔵 Low | 🟢 High | The finding supports the target agent-accessibility rating. | Verified: The synthetic guidance gives an agent a direct path from model to evaluation report. |

## Finding Details

<a id="finding-strength-004"></a>

### strength-004 Agent guidance points to the quality model and generated report path.

| Advice Rank | Tier | Ranking Rationale |
| --- | --- | --- |
| 7 / 7 | P3 | Ranked by expected impact on the service quality bar and report-gallery usefulness. |

#### Condition

The synthetic agent guidance names QUALITY.md, the evaluation command, and the report artifact to inspect.

#### Criteria

- `requirement:agent-harness::evaluation-entrypoint / rating:target`: Agent guidance should meet the target accessibility bar with a stable evaluation entry point.
  Rationale: The gallery records one finding per requirement so report tables stay easy to inspect.

#### Basis

Status: Verified

The synthetic guidance gives an agent a direct path from model to evaluation report.

##### Basis Evidence

(none recorded)

#### Effect

The finding supports the target agent-accessibility rating.

Rating effect: supports target

#### Evidence

- `synthetic-source:agent-harness/guidance`: The synthetic guidance links the model, evaluation workflow, and report output.
  Rationale: Synthetic source reference retained to demonstrate evidence rendering.

## Unknowns & Missing Evidence

| Type | Detail |
| --- | --- |
| (none recorded) |  |

## Legend

- `—` - not applicable or not recorded.
