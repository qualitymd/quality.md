---
type: Change Case
title: Evaluate ranking and umbrella-factor authoring clarity
description: Remove two evaluate-workflow authoring ambiguities surfaced by an evaluate feedback log — finding-ranking completeness and umbrella-factor analysis scopes.
status: Done
tags: [evaluation, skill, examples]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluate ranking and umbrella-factor authoring clarity

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0153-evaluate-ranking-and-umbrella-clarity/spec.md) — what the case must do.
- [Design doc](0153-evaluate-ranking-and-umbrella-clarity/design.md) — how it's built, and why.

## Motivation

An `evaluate` feedback log from a full-model run reported a clean, reportable
result but surfaced two authoring ambiguities that the durable contracts already
forbid yet the workflow guidance and examples do not make obvious:

1. **Finding-ranking completeness.** The first `data set` dry-run failed because
   `FindingRankingResult` ranked a curated top-N instead of every persisted
   Requirement Finding. The contract already requires accounting for every
   finding exactly once, but the skill instruction leads with prioritization
   criteria and the single worked example shows only two advice-relevant
   findings — neither demonstrates that tiering expresses priority while nothing
   is dropped. The validator caught the mistake, but only after a wasted dry-run.

2. **Umbrella-factor analysis scopes.** The log's explicit suggested improvement:
   "a worked `FactorAnalysisResult` example for an umbrella factor with no direct
   requirements (local vs local+descendant blocks) would remove ambiguity." A
   Factor with no direct Requirements has nothing for its `localAnalysis` block,
   but no spec rule states what that block records and the only worked example
   populates both blocks as `analyzed` with local Requirement inputs — the case
   that exposes the two-block distinction is undocumented.

Both are low-risk clarity fixes that make the existing contract obvious at
authoring time and prevent the recurring friction.

## Scope

Covered:

- Reframe the skill runtime finding-ranking instruction to lead with
  completeness (rank all findings; tier expresses priority, not inclusion).
- Make the generated `FindingRankingResult` example demonstrate a low-value tail
  finding included only because completeness requires it.
- State the umbrella-factor analysis-scope rule in the durable protocol: a Factor
  with no direct Requirements records `localAnalysis` as `empty` while
  `localAndDescendantAnalysis` carries the child-Factor roll-up.
- Make the generated `FactorAnalysisResult` example the worked umbrella case so
  the local vs local+descendant distinction is concrete.
- Add a one-line skill clarification for the umbrella case.
- Extend example tests to cover the new shapes.

Deferred:

- Emitting more than one example per data kind. Case 0133 settled the
  one-representative-artifact-per-kind contract for `data example`; this case
  keeps it.
- Any change to the worst-of roll-up, rating, or report logic.

## Affected artifacts

Derived by sweeping for `FindingRanking`, `FactorAnalysisResult`,
`localAnalysis`/`localAndDescendantAnalysis`, the analysis-scope `empty` status,
and the affected example constructors across code, durable specs, and skill
runtime.

**Code**

- [x] `internal/evaluation/data.go` — add a low-value tail finding to
      `findingRankingExample`; give `FactorAnalysisResult` its own umbrella-case
      example constructor.
- [x] `internal/evaluation/data_contract.go` — point the `FactorAnalysis` kind at
      the new example constructor.
- [x] `internal/evaluation/evaluation_test.go` — assert the finding-ranking and
      umbrella factor-analysis example shapes.

**Durable specs** (substance in the [functional spec](0153-evaluate-ranking-and-umbrella-clarity/spec.md))

- [x] `specs/evaluation/protocol.md` — state the umbrella-factor analysis-scope
      rule (Factor Traversal) and clarify that finding ranking tiers express
      priority without omitting any finding (Advice Flow).
- [x] `specs/evaluation/records/payload-kinds.md` — note that
      `FindingRankingResult` tiers express priority, not inclusion, and that a
      Factor analysis `localAnalysis` is `empty` when the Factor has no direct
      Requirements.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/workflows/evaluate.md` — lead the finding-ranking step with
      completeness; add a one-line umbrella-factor analysis clarification.

No planned impact: `SPECIFICATION.md`, `README.md`, install/scaffold files,
`CHANGELOG.md`, or the generated JSON Schema artifact. The analysis-scope status
enum already includes `empty`, and examples are not part of the schema (matching
the precedent set by case 0133). The skill _behavior_ contract is unchanged —
the runtime edit only clarifies already-required behavior — so no skill
functional spec change is needed.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Implemented and
archived: skill runtime guidance, generated examples, durable protocol and
payload-kinds rules, and focused example tests. `go test ./...` passes and the
CLI emits the worked umbrella `FactorAnalysisResult` and the
completeness-demonstrating `FindingRankingResult`. Affected artifacts reconciled.
