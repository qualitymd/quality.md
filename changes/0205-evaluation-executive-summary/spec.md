---
type: Functional Specification
title: Evaluation executive summary
description: Requirements for an advice-phase executive-summary step and its rendering into the run report Summary.
status: Draft
tags: [evaluation, advice, reports, summary]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluation executive summary

Functional spec for the
[0205 - Evaluation executive summary](../0205-evaluation-executive-summary.md)
change case. It states the delta only; the durable contracts it lands in are the
evaluation [orchestration](../../specs/evaluation/orchestration.md),
[protocol](../../specs/evaluation/protocol.md),
[payload kinds](../../specs/evaluation/records/payload-kinds.md),
[data layout](../../specs/evaluation/records/data-layout.md),
[report tree](../../specs/evaluation/reports/report-tree.md), and
[evaluation.json](../../specs/evaluation/evaluation-json.md) specs (normative).

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

The run report `## Summary` renders the root area roll-up rationale verbatim. That
rationale is an aggregation audit trace — it names the `worst_bound` mechanism and
lists each factor and child area as `name=level` — because a roll-up node
synthesizes only its direct children and is not meant to speak across the whole
evaluation. For the one section stakeholders read first, that produces jargon-led,
non-descriptive output that omits every ranked finding and recommendation directly
below it.

The advice phase (`rankFindings` → `recommend` → `rankRecommendations`) is the
only phase that already holds whole-evaluation context. Adding one more advice
step there lets inference write a stakeholder-facing summary from the ranked
findings, ranked recommendations, and root roll-up, while the deterministic report
build merely renders the persisted result. This preserves the pipeline's existing
boundary — evaluators infer, the runner and report build stay deterministic — and
leaves the roll-up rationale as the audit trace it is.

## Scope

Covered: a `summarizeEvaluation` advice unit, an `EvaluationSummaryResult` payload,
its protocol move and prompt, its rendering into the run report `## Summary`, and
the schema-version advance. Deferred and non-goals are recorded in the
[parent case](../0205-evaluation-executive-summary.md#scope).

## Assumptions & dependencies

- The advice phase already runs `rankFindings`, `recommend`, and
  `rankRecommendations` as evaluator-backed units before the deterministic
  `buildReports` unit ([orchestration](../../specs/evaluation/orchestration.md)).
  A change to that ordering would invalidate R1.
- The root area analysis (`analyzeArea` for the scoped root) records the overall
  rating level and drivers the summary reads. A change to that payload's identity
  or availability would flag R1 and R4.
- The rating scale and its level labels are resolvable at advice time so the
  summary can name the overall level in words (R5).

## Requirements

### R1 — Advice-phase summary work unit

The runner **MUST** include a `summarizeEvaluation` evaluator-backed work unit in
the evaluation work graph. It **MUST** depend on the scoped root
`AreaAnalysisResult`, the `FindingRankingResult`, and the
`RecommendationRankingResult`, and the runner **MUST** schedule it after
`rankRecommendations` and before `buildReports`.

> > Rationale: the summary must speak across the whole evaluation, so it can only
> > run once ranking is complete; placing it before `buildReports` keeps report
> > generation a pure render of persisted results.
>
> > Durable spec: modify `specs/evaluation/orchestration.md` — add
> > `summarizeEvaluation` to the work-graph order and its dependency rules;
> > modify `specs/evaluation/protocol.md` — add the `summarizeEvaluation` move to
> > protocol order and the advice flow.

### R2 — Single persisted summary payload

When `summarizeEvaluation` completes, the runner **MUST** persist exactly one
`EvaluationSummaryResult` payload at `data/advice/evaluation-summary-result.json`,
validated against its schema before writing, under the writable-kind rules that
govern other advice payloads.

> > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — register
> > `EvaluationSummaryResult` as a writable advice kind; modify
> > `specs/evaluation/records/data-layout.md` — record its
> > `data/advice/evaluation-summary-result.json` path.

### R3 — Summary payload shape

`EvaluationSummaryResult` **MUST** include `headline`, `summary`, and `keyPoints`.
`headline` **MUST** be a single non-empty sentence stating the overall standing.
`summary` **MUST** be a non-empty narrative of at most roughly 120 words.
`keyPoints` **MUST** be an array of 3–5 non-empty entries, each a self-contained
takeaway. The kind **MUST NOT** carry ratings, counts, finding IDs, or
recommendation IDs as structural fields.

> > Rationale: a bounded, structured shape lets the report render a headline,
> > paragraph, and bullets without duplicating the Key details, Top findings, and
> > Top recommendations tables that follow; keeping IDs and counts out of the
> > structure prevents the summary from drifting out of sync with the ranked
> > tables that own those numbers.
>
> > Durable spec: modify `specs/evaluation/records/payload-kinds.md` — specify the
> > `EvaluationSummaryResult` field contract and bounds.

### R4 — Synthesis only, no new evidence

The `summarizeEvaluation` evaluator **MUST** synthesize only from its provided
advice inputs and **MUST NOT** inspect the workspace or gather new evidence. It
**MUST NOT** introduce a claim, count, area, or risk absent from those inputs, and
it **MUST** reflect the evaluation's recorded confidence rather than overstating.

> > Rationale: the summary is a presentation of existing outputs; letting it
> > introduce unbacked claims would make the most-read section the least
> > trustworthy and break the runner's tools-off synthesis boundary for advice
> > work.
>
> > Durable spec: modify `specs/evaluation/protocol.md` — state the tools-off,
> > inputs-only synthesis rule for the `summarizeEvaluation` move in the advice
> > flow.

### R5 — Stakeholder-facing register

The `EvaluationSummaryResult` text **MUST** name the overall rating using the
model's rating-scale level label in plain language and **MUST** name the concrete
areas and risks that drive the standing. It **MUST NOT** use roll-up mechanism
terms (for example `worst_bound`, "roll-up", or "work unit") or unexplained
internal jargon.

> > Rationale: this is the requirement that directly replaces the "Worst-bound
> > roll-up … =minimum" register with descriptive, professional prose for a
> > non-specialist reader.
>
> > Durable spec: modify `specs/evaluation/protocol.md` — state the stakeholder-
> > facing register and the mechanism-term prohibition for the summary move.

### R6 — Run report Summary renders the executive summary

The run report `## Summary` section **MUST** render the `EvaluationSummaryResult`
— its `headline`, `summary`, and `keyPoints` — and **MUST NOT** render the scoped
area roll-up rationale. Detail (area, factor, requirement) reports **MUST** keep
their existing scoped rationale and evidence summaries unchanged.

> > Rationale: the executive summary is a whole-evaluation artifact; detail
> > reports still want their node-local roll-up explanation, so the replacement is
> > scoped to the run report only.
>
> > Durable spec: modify `specs/evaluation/reports/report-tree.md` — change the
> > run report `## Summary` contract from "render the scoped area summary" to
> > "render the `EvaluationSummaryResult`"; leave detail-report summary contracts
> > unchanged.

### R7 — Deterministic render

Report generation **MUST** render the run report `## Summary` from the persisted
`EvaluationSummaryResult` without performing inference, rereading the workspace, or
regathering evidence.

> > Rationale: preserving the deterministic, reproducible report build is the whole
> > point of doing synthesis in the advice phase rather than at render time.
>
> > Durable spec: modify `specs/evaluation/reports/report-tree.md` — state that
> > the run report `## Summary` is a deterministic render of the persisted
> > summary payload.

### R8 — Schema version and output index

The evaluation output artifact schema version **MUST** advance for the new payload
kind, and `EvaluationOutputResult` **MUST** reference the persisted
`EvaluationSummaryResult`. Report generation **MUST** require a valid
`EvaluationSummaryResult` among the advice outputs it depends on, and **MUST**
refuse to build reports for a run whose persisted artifacts predate the new schema
version.

> > Rationale: the summary is a required advice output feeding the run report, so
> > it joins the other advice outputs as a report-generation precondition; the
> > version bump keeps stale runs from rendering an absent section.
>
> > Durable spec: modify `specs/evaluation/evaluation-json.md` — advance the
> > artifact schema version and add the `EvaluationSummaryResult` reference to
> > `EvaluationOutputResult`; modify `specs/evaluation/orchestration.md` — add the
> > summary to the advice outputs required before report generation.

## Requirement-set check

- **Consistent.** One term per surface: `summarizeEvaluation` (unit/move),
  `EvaluationSummaryResult` (payload), `## Summary` (run report section). No
  requirement contradicts another; R6/R7 both govern the run report `## Summary`
  from different angles (content vs. determinism) without overlap.
- **Complete.** The set covers production (R1–R2), shape (R3), judgment discipline
  (R4–R5), rendering (R6–R7), and lifecycle/versioning (R8). No inline "to be
  decided" remains; the degenerate no-findings case is handled by R4's inputs-only
  rule plus the existing advice guarantee that ranking and recommendation outputs
  always exist.
- **Able to be validated.** Satisfying the set yields exactly the motivation: a
  concise, professional, friendly, descriptive `## Summary` produced by inference
  from whole-evaluation context and rendered deterministically — with no change to
  ratings, findings, recommendations, or the roll-up rationale.

## Durable spec changes

### To add

None. `EvaluationSummaryResult` is an advice payload of the same family as
`FindingRankingResult` and `RecommendationRankingResult`; its kind contract is
owned by [`payload-kinds.md`](../../specs/evaluation/records/payload-kinds.md) and
its artifact path by
[`data-layout.md`](../../specs/evaluation/records/data-layout.md), matching the
existing advice payloads. No 1:1 artifact spec is warranted.

### To modify

- `specs/evaluation/orchestration.md` — add the `summarizeEvaluation` work unit to
  the graph, its dependencies, and its place among required advice outputs before
  report generation (per R1, R8).
- `specs/evaluation/protocol.md` — add the `summarizeEvaluation` move to protocol
  order and specify its inputs-only synthesis and stakeholder-facing register in
  the advice flow (per R1, R4, R5).
- `specs/evaluation/records/payload-kinds.md` — register `EvaluationSummaryResult`
  as a writable advice kind and specify its field contract and bounds (per R2, R3).
- `specs/evaluation/records/data-layout.md` — record the
  `data/advice/evaluation-summary-result.json` path (per R2).
- `specs/evaluation/reports/report-tree.md` — change the run report `## Summary`
  contract to render the `EvaluationSummaryResult` deterministically, leaving
  detail-report summaries unchanged (per R6, R7).
- `specs/evaluation/evaluation-json.md` — advance the artifact schema version and
  add the `EvaluationSummaryResult` reference to `EvaluationOutputResult` (per R8).

### To rename

None.

### To delete

None.
