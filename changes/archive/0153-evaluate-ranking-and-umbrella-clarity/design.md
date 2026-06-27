---
type: Design Doc
title: Evaluate ranking and umbrella-factor authoring clarity — design
description: How the finding-ranking and umbrella-factor analysis authoring ambiguities are removed across skill guidance, examples, and durable specs.
tags: [evaluation, skill, examples]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluate ranking and umbrella-factor authoring clarity — design

## Context

Answers the [functional spec](spec.md) for change case
[0153](../0153-evaluate-ranking-and-umbrella-clarity.md). The durable contracts
already require both behaviors: `specs/evaluation/protocol.md` Advice Flow says
`rankFindings` must consider every Requirement Finding and write one
`FindingRankingResult`, and the analysis-scope contract in
`internal/evaluation/data_contract.go` already offers the `empty` status. The gap
is that the skill runtime guidance and the generated examples do not make either
obvious, so an agent authored a curated top-N (failing the first dry-run) and the
log asked for a worked umbrella-factor example.

The examples are built by per-kind constructor functions in
`internal/evaluation/data.go`; their structural contracts and the kind→example
mapping live in `internal/evaluation/data_contract.go`. Case 0133 established that
`data example <kind>` emits one representative complete artifact per kind, so the
fix works within that contract rather than emitting multiple examples.

## Approach

### Finding-ranking completeness

- **Skill runtime** (`skills/quality/workflows/evaluate.md`): reorder the
  ranking bullet so completeness leads — "rank every persisted Requirement
  Finding; tier and order express priority, no finding is dropped" — before the
  prioritization criteria, removing the read that the criteria license curation.
- **Example** (`findingRankingExample` in `data.go`): add a third entry at the
  lowest tier (`P4`) — a low-value finding whose rationale says it appears only
  because the ranking accounts for every finding. Update the closing rationale to
  name completeness. This turns the example into a demonstration that tiering
  expresses priority, not inclusion.
- **Durable spec** (`protocol.md` Advice Flow, `payload-kinds.md`): add one
  clause to the existing finding-ranking rule that tiers/order express relative
  priority and no finding is omitted.

### Umbrella-factor analysis scopes

- **Durable spec** (`protocol.md` Factor Traversal): state that a Factor with no
  direct Requirements records `localAnalysis` as `empty` (with a reason) while
  `localAndDescendantAnalysis` carries the child-Factor roll-up. This earns the
  example and names the previously-unstated rule. Add a matching note in
  `payload-kinds.md` next to the analysis-scope contract.
- **Example**: stop sharing `scopedAnalysisExample` between `FactorAnalysis` and
  `AreaAnalysis`. Keep `scopedAnalysisExample` (both blocks analyzed from local
  inputs) for `AreaAnalysis`, where a root/area legitimately has local
  Requirements and Factors. Add a dedicated `factorAnalysisExample` for the
  umbrella case: `localAnalysis` with `status: empty`, a `statusReason`, and no
  Requirement inputs; `localAndDescendantAnalysis` with `status: analyzed`, a
  rating, drivers, and `inputRefs`/driver inputs that reference a *child* Factor
  analysis (`exampleChildFactorID()` with the `localAndDescendantAnalysis`
  selector). Repoint the `FactorAnalysis` registry entry in `data_contract.go` to
  the new constructor.
- **Skill runtime**: add a one-line clarification to the factor-analysis step
  distinguishing the two blocks for a Factor with no direct Requirements.

### Verification

Case 0133 already validates every generated example kind structurally in
`evaluation_test.go`; the new `factorAnalysisExample` is covered by that loop
automatically. Add two focused assertions:

- the finding-ranking example ranks more than one finding and includes a `P4`
  tail entry; and
- the factor-analysis example's `localAnalysis.status` is `empty` with no
  Requirement inputs while `localAndDescendantAnalysis.status` is `analyzed` with
  a child Factor analysis input ref.

Then run `go test ./...` and a `qualitymd evaluation data example` /
`data schema` spot check.

## Spec response

- **Finding-ranking completeness** — satisfied by the reordered skill bullet, the
  tail-finding example entry, and the protocol/payload-kinds clarification.
- **Umbrella-factor analysis scopes** — satisfied by the new protocol rule, the
  dedicated umbrella `factorAnalysisExample`, and the skill clarification.
- **Verification** — satisfied by the existing all-kind structural loop plus the
  two focused assertions.

## Alternatives

- **Emit two FactorAnalysis examples (leaf + umbrella).** Rejected: case 0133
  fixed `data example <kind>` to one representative artifact per kind. The
  umbrella case is the more instructive single example because it is the one that
  exposes the local vs local+descendant distinction; the leaf shape is mirrored
  by the area example and fully constrained by `data schema`.
- **Put the umbrella example only in spec prose as a fenced JSON block.**
  Rejected: the eval specs avoid large inline JSON and lean on `data example` for
  concrete shapes. Making the CLI's single factor example the umbrella case is
  more discoverable at authoring time and keeps the spec prose-first.
- **Only fix the example, not the skill instruction.** Rejected: the dry-run
  already enforces completeness, so the cost being removed is the *wasted first
  attempt*; that comes from the instruction's framing, so the instruction must
  change too.
- **Add a new validator rule rejecting curated rankings.** Unnecessary: the
  existing "account for every finding exactly once" validation already rejects a
  top-N. The fix is making the rule obvious, not adding enforcement.

## Trade-offs & risks

- Making the single factor example the umbrella case means `data example
  factor-analysis` no longer shows a factor *with* direct Requirements. Accepted:
  the area example shows both-blocks-analyzed shape, and `data schema` carries the
  exhaustive constraint surface.
- A longer finding-ranking example is slightly heavier, but the tail entry is the
  point — it teaches the completeness rule the log showed agents miss.

## Open questions

None.
