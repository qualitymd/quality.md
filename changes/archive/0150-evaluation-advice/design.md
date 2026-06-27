---
type: Design Doc
title: Evaluation Advice Design
description: Implementation approach for required evaluation Advice, ranking payloads, coverage accounting, and recommendation reports.
tags: [evaluation, advice, recommendations, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluation Advice Design

## Context

This design answers the [Evaluation Advice functional spec](spec.md). The change
adds a required late Advice phase to Evaluation: after Requirement assessment,
Requirement ratings, and Factor/Area roll-up analysis, the agent ranks existing
Findings, generates Recommendations, accounts for every Finding, ranks the
Recommendations, and hands complete structured Advice to the CLI for status and
report rendering.

The design keeps the existing responsibility split:

- The `/quality` skill owns judgment: Finding priority, recommendation synthesis,
  coverage rationale, impact, confidence, and ranking rationale.
- The CLI owns mechanics: payload validation, reference resolution,
  reportability status, deterministic report projection, and generated file
  paths.

## Approach

### Data model

Add three agent-written payload kinds to the existing evaluation data registry:

```text
FindingRankingResult
RecommendationResult
RecommendationRankingResult
```

They follow the existing routine-output conventions: JSON objects with
`schemaVersion`, `kind`, stable subject references where needed, and payload-local
IDs where the payload owns identity. They are written through the existing bulk
`qualitymd evaluation data set` path, validated by the same contract registry
that drives `schema`, `example`, `set`, and generated JSON Schema.

`FindingRankingResult` is one run-scope payload. It contains an ordered list of
every persisted Requirement Finding in the evaluation scope:

```jsonc
{
  "kind": "FindingRankingResult",
  "orderedFindings": [
    {
      "rank": 1,
      "findingRef": {
        "kind": "RequirementAssessmentResult",
        "subject": {
          "requirementId": "requirement:household-budget::irregular-obligations"
        },
        "selector": "findings[gap-001]"
      },
      "tier": "P1",
      "rationale": "..."
    }
  ]
}
```

The payload does not copy Finding content except for optional display/cache fields
that the renderer can ignore; the persisted Finding remains the source of truth.
The CLI validates that every referenced Finding exists and that every in-scope
Finding appears exactly once.

`RecommendationResult` is one payload per Recommendation:

```jsonc
{
  "kind": "RecommendationResult",
  "id": "rec-001",
  "title": "Review whether to raise the irregular-obligations bar",
  "whyItMatters": "...",
  "recommendedNextMove": "...",
  "expectedBenefit": "...",
  "howToKnowItWorked": "...",
  "impact": "high",
  "confidence": "medium",
  "traceRefs": [
    {
      "kind": "RequirementAssessmentResult",
      "subject": {
        "requirementId": "requirement:household-budget::irregular-obligations"
      },
      "selector": "findings[strength-001]"
    }
  ]
}
```

Field names are intentionally close to the user-facing report labels. The JSON
contract may use camelCase; Markdown renders title-cased labels. No `effort`,
`roi`, quick-win, backlog-priority, or score fields are added.

`RecommendationRankingResult` is the final Advice closure payload. It ranks every
Recommendation and carries Finding coverage accounting:

```jsonc
{
  "kind": "RecommendationRankingResult",
  "orderedRecommendations": [
    {
      "rank": 1,
      "recommendationRef": "rec-001",
      "impact": "high",
      "confidence": "medium",
      "rationale": "..."
    }
  ],
  "findingCoverage": [
    {
      "findingRef": {
        "kind": "RequirementAssessmentResult",
        "subject": {
          "requirementId": "requirement:household-budget::irregular-obligations"
        },
        "selector": "findings[gap-001]"
      },
      "disposition": "addressed_by_recommendation",
      "recommendationRefs": ["rec-001"]
    },
    {
      "findingRef": {
        "kind": "RequirementAssessmentResult",
        "subject": {
          "requirementId": "requirement:household-budget::monthly-visibility"
        },
        "selector": "findings[note-001]"
      },
      "disposition": "not_advice_driving",
      "rationale": "Context only; it does not affect the next quality-management move."
    }
  ]
}
```

Keeping coverage here avoids a fourth payload kind while making
`RecommendationRankingResult` the artifact that closes Advice. The contract names
that role explicitly so the payload is not misread as only a sort list.

### Advice workflow

The skill workflow adds a deterministic sequence after roll-up analysis:

```text
Finding ranking
  -> progressive QC
Recommendation generation
  -> progressive QC per recommendation
Coverage accounting
  -> progressive QC
Recommendation ranking
  -> progressive QC
Advice closure
  -> persist payload batch / status / report build
```

The skill should still write the final Advice payloads in the same bulk
`evaluation data set` batch as other late outputs when practical. Progressive QC
is a skill-side discipline before persistence; CLI validation is the structural
backstop at persistence.

The skill may iterate locally during Advice:

1. Rank Findings and check every in-scope Finding is represented once.
2. Draft Recommendations from the ranked Findings and completed analysis.
3. For each Recommendation, check required prose fields, trace support, impact,
   confidence, no-new-evidence discipline, and no planning metadata.
4. Account for every Finding as addressed or not advice-driving.
5. If coverage reveals an unaddressed action-worthy Finding, draft or revise a
   Recommendation before ranking.
6. Rank Recommendations after coverage is coherent.
7. Run the final closure check across the Advice payload set.

This shifts quality left: final closure should find integration problems, not
basic missing fields or unsupported Recommendations.

### Status and validation

`evaluation status` extends reportability checks:

- exactly one effective `FindingRankingResult`;
- one or more effective `RecommendationResult` payloads;
- exactly one effective `RecommendationRankingResult`;
- every in-scope Requirement Finding appears exactly once in
  `FindingRankingResult.orderedFindings`;
- every `RecommendationResult` is ranked exactly once;
- every in-scope Requirement Finding appears exactly once in
  `RecommendationRankingResult.findingCoverage`;
- every addressed Finding points to at least one valid Recommendation;
- every not-advice-driving Finding has rationale.

Reference validation happens in the data contract layer as far as possible:
Finding references resolve to persisted `RequirementAssessmentResult.findings[]`,
Recommendation references resolve to persisted `RecommendationResult.id`, and
trace references resolve to persisted routine outputs or known selectors.

The status path should produce missing/invalid Advice diagnostics alongside the
existing missing assessment, rating, and analysis diagnostics so the agent can
correct one batch instead of discovering failures during report build.

### Report rendering

`evaluation report build` remains a deterministic projection. It refuses or
reports not-reportable status through the existing status gate when required
Advice payloads are absent or incomplete.

Generated run artifacts add:

```text
report.md
recommendations.md
recommendations/
  001-<slug>.md
  002-<slug>.md
```

`report.md` gains two compact sections:

- `Top Findings` - up to 10 rows from `FindingRankingResult`, linked to the
  Requirement report/finding detail where the Finding is rendered.
- `Top Recommendations` - up to 10 rows with `Rank`, `Recommendation`, `Impact`,
  and `Confidence`, linked to recommendation detail files.

`report.md` always links to `recommendations.md`.

`recommendations.md` is the full Recommendation index with the same simple table
for every Recommendation. It can include a concise coverage summary such as
addressed/not-advice-driving counts, but it should not render the full coverage
ledger unless the implementation needs it for handoff. The structured payload is
the authoritative coverage ledger.

Each recommendation detail file is human-first Markdown:

```markdown
# Review whether to raise the irregular-obligations bar

**Rank:** 1
**Impact:** High
**Confidence:** Medium

## Why it matters

...

## Recommended next move

...

## Expected benefit

...

## How to know it worked

...

## Trace

- ...
```

No YAML frontmatter or YAML appendix is emitted in this slice. Later `/quality
improve` can inspect the structured `data/` payloads by run and Recommendation
ID, using the Markdown path as a human affordance rather than as the machine
source of truth.

### Domain-agnostic example

Use a household budget example in durable examples or specs because it is far
from software and naturally exercises review, evidence, preservation, and
monitoring. The example should show:

- a Top Finding about an irregular-obligations risk or a current clarity
  strength;
- a Recommendation to review whether the bar should rise or to make irregular
  obligations visible before the next planning cycle;
- coverage accounting where a note or strength can be not advice-driving or can
  support preserve/raise-bar advice;
- Recommendation impact/confidence without effort or task-backlog framing.

## Spec response

- **Advice phase and reportability (R1-R3)** - add the late workflow sequence;
  make status require the three Advice payload kinds; keep report build a
  projection over persisted Advice.
- **Finding ranking (R4-R7)** - add `FindingRankingResult` with one entry per
  in-scope Finding, tier/rank rationale, no new evidence, and skill-side ranking
  judgment.
- **Recommendations (R8-R15)** - add one `RecommendationResult` per
  quality-management move with the core user-facing fields, impact, confidence,
  trace support, no planning metadata, and domain-agnostic semantics.
- **Recommendation ranking and coverage (R16-R21)** - make
  `RecommendationRankingResult` the Advice closure payload containing ranked
  Recommendations and full coverage accounting.
- **Reports (R22-R29)** - render Top Findings and Top Recommendations in
  `report.md`, generate a full Recommendation index and detail pages, and keep
  Recommendation Markdown human-first.
- **Progressive QC and examples (R30-R33)** - add skill checkpoints at each
  Advice artifact boundary and include a non-software example.

## Alternatives

**Separate `AdviceCoverageResult`.** Rejected for this slice. Coverage accounting
is conceptually distinct, but a fourth payload kind adds surface area without
much benefit. `RecommendationRankingResult` is already the final Advice artifact,
so carrying coverage there makes reportability depend on one closure payload.

**Put coverage on each Recommendation.** Rejected because coverage is about every
Finding, including Findings that intentionally do not drive Advice. Per-
Recommendation fields make unaddressed or not-advice-driving Findings easy to
miss.

**Make report build synthesize Recommendations from Findings.** Rejected because
it violates the judgment/mechanics boundary. The renderer should not become an
evaluator.

**Use numeric scores for Finding or Recommendation ranking.** Rejected because
the precision would be false at this stage and would invite weight debates before
the model has enough real runs.

**Expose `Kind` or `effort` in the Recommendation table.** Rejected because the
core user-facing surface should stay intelligible: Recommendation, Impact, and
Confidence. `effort` in particular pulls Advice toward backlog planning.

**Add YAML frontmatter to Recommendation detail files now.** Rejected because the
Markdown files are handoff and reading artifacts. Structured data already exists
under `data/`; future improve workflows can use CLI inspection if they need a
machine entrypoint.

## Trade-offs & risks

- **`RecommendationRankingResult` has two roles.** It ranks Recommendations and
  closes coverage. The design accepts this to avoid a fourth payload kind, but
  durable docs must name the closure role clearly.
- **Coverage can feel noisy for strengths and notes.** The all-Findings rule is
  intentional: every Finding is accounted for. The skill guidance should make
  `not_advice_driving` normal for context-only strengths and notes, while still
  allowing strengths to support preservation or raise-bar Recommendations.
- **Required Advice increases evaluation cost.** The quality loop benefits from
  always producing next guidance, but full evaluations will take longer. The
  progressive QC sequence minimizes late rework.
- **Report size can grow.** The run report caps Top Findings and Top
  Recommendations at 10 each and delegates full Recommendation detail to
  `recommendations.md` and per-Recommendation pages.
- **Follow-up semantics are broader than apply/handoff.** This case makes
  Recommendations broader than repair tasks. Existing follow-up can still apply
  or hand off compatible Recommendations, but later work may need review-brief,
  monitor, or decision-record outcomes.

## Settled implementation notes

- Trace support uses `traceRefs`, matching the `RecommendationResult` contract.
- `recommendations.md` includes a compact coverage summary; the structured
  coverage ledger remains authoritative under `data/`.
