---
type: Functional Specification
title: Evaluate ranking and umbrella-factor authoring clarity — functional spec
description: What the change must do to remove the finding-ranking and umbrella-factor analysis authoring ambiguities surfaced by an evaluate feedback log.
tags: [evaluation, skill, examples]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluate ranking and umbrella-factor authoring clarity — functional spec

Companion to the
[Evaluate ranking and umbrella-factor authoring clarity](../0153-evaluate-ranking-and-umbrella-clarity.md)
change case. This spec states *what* the change must do; the
[design doc](design.md) covers *how*. It defers to the durable Evaluation
[protocol](../../../specs/evaluation/protocol.md) and
[payload kinds](../../../specs/evaluation/records/payload-kinds.md) specs for the
contracts it clarifies.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

An `evaluate` run completed and was reportable, but its feedback log recorded two
authoring ambiguities the durable contracts already settle yet the workflow
guidance and worked examples do not make obvious:

- `FindingRankingResult` must account for every Requirement Finding exactly once,
  but the skill instruction leads with prioritization criteria and the worked
  example shows only advice-relevant findings — so a curated top-N was authored,
  failing the first dry-run.
- A Factor with no direct Requirements ("umbrella factor") has nothing for its
  `localAnalysis` block, but no rule states what that block records and the only
  worked example populates both analysis blocks from local Requirement inputs.

This change makes both already-required behaviors obvious at authoring time. It
adds no new runtime behavior the contracts did not already require.

## Scope

Covered: the skill runtime finding-ranking and factor-analysis guidance, the
generated `FindingRankingResult` and `FactorAnalysisResult` examples, the durable
protocol and payload-kinds rules that earn them, and the example tests.

Deferred (non-goals): emitting more than one example per data kind (case 0133
settled the one-representative-artifact-per-kind contract), and any change to
roll-up, rating, or report behavior.

## Requirements

### Finding-ranking completeness

- The skill runtime `evaluate` finding-ranking instruction **MUST** lead with the
  obligation to rank every persisted Requirement Finding, and **MUST** state that
  the ranking's tier and order express priority while no finding is omitted.

  > Rationale: the prior instruction led with prioritization criteria, which
  > reads as license to curate; an agent authored a top-N and failed the first
  > dry-run. — 0153
  >
  > Durable spec: none (runtime skill guidance; the behavior is already required
  > by `specs/evaluation/protocol.md` Advice Flow).

- The generated `FindingRankingResult` example **MUST** include at least one
  low-value tail finding whose rationale shows it is ranked only because the
  ranking accounts for every finding, demonstrating that tiering expresses
  priority rather than inclusion.

  > Rationale: an example with only advice-relevant entries does not teach that
  > low-value findings still appear in the ranking. — 0153
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — state that
  > `FindingRankingResult` tiers express priority, not inclusion.

- `specs/evaluation/protocol.md` Advice Flow **MUST** state that the finding
  ranking's tiers and order express relative priority and that no Requirement
  Finding is omitted from the ranking.

  > Durable spec: modify `specs/evaluation/protocol.md` — add the
  > priority-not-omission clarification to the existing finding-ranking rule.

### Umbrella-factor analysis scopes

- `specs/evaluation/protocol.md` Factor Traversal **MUST** state that when a
  Factor has no direct Requirements, its `localAnalysis` records the `empty`
  status with a reason and its `localAndDescendantAnalysis` carries the
  child-Factor roll-up.

  > Rationale: the analysis-scope `empty` status exists but no rule said an
  > umbrella factor uses it, leaving the local block's content ambiguous. — 0153
  >
  > Durable spec: modify `specs/evaluation/protocol.md` — add the umbrella-factor
  > analysis-scope rule.

- The generated `FactorAnalysisResult` example **MUST** represent an umbrella
  Factor: `localAnalysis` carries the `empty` status with a reason and no
  Requirement inputs, and `localAndDescendantAnalysis` carries a rated roll-up
  whose inputs reference a child Factor analysis.

  > Rationale: the feedback log asked for a worked umbrella example to remove the
  > local vs local+descendant ambiguity; the single per-kind example is the most
  > discoverable home. — 0153
  >
  > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — note that a
  > Factor analysis `localAnalysis` is `empty` when the Factor has no direct
  > Requirements.

- The skill runtime `evaluate` factor-analysis guidance **MUST** distinguish the
  `localAnalysis` and `localAndDescendantAnalysis` blocks for a Factor with no
  direct Requirements.

  > Durable spec: none (runtime skill guidance; the behavior is the protocol rule
  > above).

### Verification

- Tests **MUST** verify that the generated `FindingRankingResult` example ranks
  more than one finding and includes a low-value tail entry at the lowest tier.

  > Durable spec: none.

- Tests **MUST** verify that the generated `FactorAnalysisResult` example has a
  `localAnalysis` with `empty` status and no Requirement inputs and a
  `localAndDescendantAnalysis` that is `analyzed` with a child Factor analysis
  input reference.

  > Durable spec: none.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) — the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/evaluation/protocol.md` — clarify that finding-ranking tiers express
  priority without omitting any finding (Advice Flow), and state the
  umbrella-factor analysis-scope rule (Factor Traversal) (per the finding-ranking
  and umbrella requirements above).
- `specs/evaluation/records/payload-kinds.md` — note that `FindingRankingResult`
  tiers express priority not inclusion, and that a Factor analysis `localAnalysis`
  is `empty` when the Factor has no direct Requirements (per the example
  requirements above).

### To rename

None.

### To delete

None.
