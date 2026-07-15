---
type: Design Doc
title: Evaluation executive summary
description: How the advice-phase summary unit, payload, prompt, and report rendering are built.
status: Draft
tags: [evaluation, advice, reports, summary]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluation executive summary

Design for the
[0205 - Evaluation executive summary](../0205-evaluation-executive-summary.md)
change case. It answers the [functional spec](spec.md); behavior and requirements
live there.

## Context

The run report `## Summary` renders `field(analysis, "rationale")` from the root
area's `localAndDescendantAnalysis` (`src/domain/evaluation/report.ts`), which is a
`worst_bound` aggregation trace. We want a stakeholder-facing narrative there
without making the report generator infer. The advice phase is the natural home:
it already runs global, evaluator-backed steps (`rankFindings`, `recommend`,
`rankRecommendations`) over the whole evaluation before the deterministic
`buildReports` unit.

## Approach

**One new advice unit.** Add `summarizeEvaluation` to `WorkKind` and to
`buildGraph` in `src/domain/evaluation/graph.ts`, mirroring `rankRecommendations`:
evaluator-backed, run-wide subject, `dataKind: "EvaluationSummaryResult"`,
depending on the scoped root `analyzeArea`, `rankFindings`, and
`rankRecommendations`. `buildReports` gains it as a dependency. Because scheduling
is dependency-driven, `evaluation-execute.ts` and `evaluation-resume.ts` need no
change — the unit flows through the existing ready-frontier loop.

```
… analyzeArea(root) ─┐
rankFindings ────────┼─▶ summarizeEvaluation ─▶ buildReports
rankRecommendations ─┘
```

**Protocol context and prompt.** Add a `summarizeEvaluation` branch to
`protocolParts` in `src/domain/evaluation/protocol.ts` that assembles the whole-
evaluation view — no workspace inspection block (like the other advice moves):

```ts
if (unit.kind === "summarizeEvaluation") {
  return {
    instructions: instructions.summarizeEvaluation,
    context: {
      overallAnalysis: payloadFor(payloads, `analyzeArea:${rootAreaRef}`),
      evaluationFrame: payloadFor(payloads, "frameEvaluation"), // rating level ids
      ratingScale,                             // level labels for plain-language naming
      findings: findingIndex(plan, payloads),
      findingRanking: payloadFor(payloads, "rankFindings"),
      recommendations: payloadsFor(payloads, "recommend"),
      recommendationRanking: payloadFor(payloads, "rankRecommendations"),
    },
  }
}
```

The rating-scale level labels the summary uses for plain-language naming come from
the `QualityModel.ratingScale` (the same source `frames.ts` reads); the frame
payload carries only `ratingLevelIds`, so `protocolParts` threads the scale's
labels into the summary context rather than expecting them on `plan`.
`expectedSchemaFor` returns `kindSchema("EvaluationSummaryResult")`. The instruction
string, in the house style of the other advice moves:

```
summarizeEvaluation:
  "Write a stakeholder-facing executive summary of the whole evaluation from the "
    + "overall analysis, ranked findings, and ranked recommendations, and return one "
    + "EvaluationSummaryResult JSON object.\n"
    + "- headline: one sentence giving the bottom line — name the overall rating by its "
    + "scale label and the single biggest reason it lands there.\n"
    + "- summary: 3-5 sentences (~120 words max), professional and plain-spoken. Say where "
    + "the entity stands overall, what is working well, what most limits it, and where to "
    + "focus next. Lead with the bottom line, not the method.\n"
    + "- keyPoints: 3-5 short, self-contained takeaways drawn from the highest-ranked "
    + "findings and recommendations; each names a concrete area and what it means for a "
    + "reader, not a metric or an id.\n"
    + "- Write for a non-specialist stakeholder. Use the rating scale's level labels in "
    + "plain language; do not use mechanism terms (worst-bound, roll-up, work unit) or "
    + "unexplained internal jargon.\n"
    + "- Synthesize only the given inputs. Do not introduce claims, counts, or risks absent "
    + "from the findings, analyses, or recommendations, and do not overstate — reflect the "
    + "evaluation's stated confidence.\n"
    + "- Be specific: name the areas and risks that drive the rating rather than describing "
    + "them generically. Do not inspect new evidence."
```

**Payload and schema.** Register `EvaluationSummaryResult` (writable) in
`src/domain/evaluation/data.ts` and add its `$defs` entry to
`src/assets/evaluation-data.schema.json`:

```jsonc
{
  "kind": { "const": "EvaluationSummaryResult" },
  "headline":  { "type": "string", "minLength": 1 },
  "summary":   { "type": "string", "minLength": 1 },
  "keyPoints": { "type": "array", "minItems": 3, "maxItems": 5,
                 "items": { "type": "string", "minLength": 1 } }
}
```

The ~120-word `summary` bound is enforced as prompt guidance, not schema, to avoid
brittle token-count validation; schema enforces presence and the `keyPoints` count.

**Report rendering.** In `report.ts`, the run report path (line ~725) renders the
`EvaluationSummaryResult` instead of `summary(overall)`: `headline` as a lead line,
`summary` as the body paragraph, and `keyPoints` as a bullet list. The detail-area
path (line ~802) is untouched, so area/factor reports keep `summary(overall)`.
`evaluation-report.ts` bumps the artifact schema version and adds the summary
reference to `EvaluationOutputResult`.

## Spec response

- **R1–R2 (unit + payload):** the graph unit and its dependencies place the step
  after ranking and before report build; the runner's existing accept-and-persist
  path writes the single validated payload to the advice folder.
- **R3 (shape):** JSON Schema enforces the field contract and the 3–5 `keyPoints`
  bound; the word cap is prompt-enforced.
- **R4–R5 (discipline + register):** the prompt carries the inputs-only, no-new-
  evidence, plain-language, and mechanism-term-ban rules; no inspection block is
  supplied, matching the tools-off advice moves.
- **R6–R7 (render):** only the run-report branch changes, and it renders the
  persisted payload with no inference — determinism is preserved because the
  advice unit already ran.
- **R8 (versioning):** the schema-version advance plus the report-generation
  dependency on the summary output make a stale run refuse rather than render an
  empty section.

## Alternatives

- **Deterministic templated summary in `report.ts` (no inference).** Compose
  sentences from ratings, ranked finding titles, and recommendation counts. Fully
  deterministic and zero model cost, but it can only recombine existing fields —
  it cannot produce genuine narrative prose, which is the stated goal. Rejected as
  the primary approach; it remains the fallback if the advice step is ever
  disabled.
- **Change the `analyzeArea` roll-up rationale to be reader-friendly.** Overloads
  one field to serve both an audit trace and a stakeholder summary, and the root
  roll-up still only sees its direct children, not the ranked findings and
  recommendations. Rejected: it degrades the audit trace and lacks the context the
  summary needs.
- **A standalone executive-summary report file.** More surface than the problem
  needs; the pain is specifically the `## Summary` section of `report.md`.
  Deferred.

## Trade-offs & risks

- **One more evaluator call per run.** Small, single, run-wide; consistent with the
  three advice calls already made. Acceptable for the value on the most-read
  section.
- **Prompt-enforced length can drift.** A model may exceed ~120 words. Mitigated by
  explicit caps in the prompt and a bounded `keyPoints` array; if drift proves
  real, a soft truncation at render time can be added later.
- **Schema-version bump invalidates old runs' report rebuild.** Intended (R8); the
  bump is a clean break with no back-compat reader, per the early-alpha policy.

## Open questions

- Should `keyPoints` render as a labelled sub-list under `## Summary`, or flow
  after the paragraph unlabelled? Leaning unlabelled bullets directly under the
  paragraph to keep the section tight; final call at implementation against a
  rendered snapshot.
