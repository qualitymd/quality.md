---
type: Change Case
title: Evaluation executive summary
description: Add an advice-phase executive-summary synthesis step so the run report's Summary reads as a stakeholder-facing narrative instead of the root roll-up trace.
status: Design
tags: [evaluation, advice, reports, summary, agents]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluation executive summary

Status note: case is in **Design**. The functional spec and design doc are
authored; no production code or durable current-behavior spec has changed yet.
Code work begins only after the case advances to **In-Progress**.

## Motivation

The run-level `report.md` `## Summary` is the first thing a stakeholder reads,
and it is currently the least readable content in the file. It renders the root
area's roll-up rationale verbatim — for example:

> Worst-bound roll-up. Factors: agent-harnessability=minimum, coherence=target,
> …; child areas: agent-harness=minimum, api-service=target, … Rolls up to
> minimum.

That string is an aggregation audit trail, by design: `analyzeArea` synthesizes
only its direct children (`worst_bound`), names the mechanism, and dumps a flat
`key=level` list. No single roll-up node is meant to synthesize across the whole
evaluation, so the root summary reads as a machine trace — it leads with internal
jargon, never names what is actually wrong, and omits every ranked finding and
recommendation sitting right below it.

The evaluation pipeline already has the right seam for a fix. The **advice phase**
(`rankFindings` → `recommend` → `rankRecommendations`) runs global, evaluator-
backed inference over the whole evaluation before the deterministic report build.
That is the only phase with the full picture — every ranked finding, every
recommendation with its impact and reason, and the root roll-up. This case adds
one more advice step that writes a concise, professional, friendly, descriptive
executive summary from those inputs, persists it as an advice payload, and has
the deterministic report render it into `## Summary`.

This keeps the existing separation of concerns: inference does the writing where
whole-evaluation context already exists; report generation stays deterministic
and simply renders a pre-computed field. The roll-up rationale stays a clean
audit trace rather than being overloaded to serve two audiences.

## Scope

Covered:

- a new advice-phase evaluator-backed work unit, `summarizeEvaluation`, that
  depends on the root area analysis, the finding ranking, and the recommendation
  ranking, and runs after `rankRecommendations` and before `buildReports`;
- a new writable advice payload kind, `EvaluationSummaryResult`, carrying a
  one-sentence `headline`, a short narrative `summary`, and 3–5 `keyPoints`,
  persisted at `data/advice/evaluation-summary-result.json`;
- a protocol move and prompt for the summary that synthesizes only the given
  advice inputs, writes in plain stakeholder-facing language, names concrete
  areas and risks, and does not introduce claims, counts, or risks absent from
  the inputs;
- rendering the executive summary into the run report `## Summary` in place of
  the scoped-area roll-up rationale; and
- advancing the evaluation output schema version for the new payload and its
  inclusion in `EvaluationOutputResult`.

Deferred:

- a standalone executive-summary report artifact separate from `report.md`;
- executive summaries on detail (area/factor/requirement) reports, which keep
  their existing scoped rationale and evidence summaries; and
- localization or configurable summary length/tone.

Non-goals:

- changing ratings, roll-up semantics (`worst_bound`), findings, recommendations,
  or coverage accounting — the summary reads existing outputs and adds none;
- making the report generator perform inference — report build stays
  deterministic and only renders the persisted summary; and
- re-tasking the `analyzeArea`/`analyzeFactor` roll-up rationale to serve
  stakeholders — it stays a terse audit trace.

## Affected artifacts

Derived from repository searches for the advice-phase work units (`rankFindings`,
`recommend`, `rankRecommendations`), the report `## Summary` rendering, the data-
kind registry and JSON Schema, the work-graph builder, and the payload/data-layout
specs.

- **Change record:** this parent, `spec.md`, and `design.md` under
  `changes/0205-evaluation-executive-summary/`; `changes/index.md` and
  `changes/log.md`. Archived together with this parent when the case lands.
- **Durable evaluation specs:** `specs/evaluation/orchestration.md`,
  `specs/evaluation/protocol.md`, `specs/evaluation/records/payload-kinds.md`,
  `specs/evaluation/records/data-layout.md`,
  `specs/evaluation/reports/report-tree.md`, and
  `specs/evaluation/evaluation-json.md` — detailed per requirement in the
  spec's [Durable spec changes](0205-evaluation-executive-summary/spec.md#durable-spec-changes).
  `specs/evaluation/log.md` and the relevant sub-index `log.md` files record the
  revision.
- **Domain code:** `src/domain/evaluation/graph.ts` adds the `summarizeEvaluation`
  `WorkKind` and unit; `src/domain/evaluation/protocol.ts` adds its instructions,
  context assembly, and expected schema; `src/domain/evaluation/data.ts` registers
  the `EvaluationSummaryResult` kind and its data path;
  `src/domain/evaluation/report.ts` renders the executive summary into the run
  report `## Summary`.
- **Schema asset:** `src/assets/evaluation-data.schema.json` adds the
  `EvaluationSummaryResult` definition.
- **Application code:** `src/application/evaluation-report.ts` advances the
  artifact schema version and includes the summary reference in
  `EvaluationOutputResult`; `src/application/evaluation-data.ts` validates the new
  kind if it carries kind-specific acceptance rules. `evaluation-execute.ts` and
  `evaluation-resume.ts` require no change beyond the graph gaining the unit, and
  that assumption is verified in test.
- **Tests:** work-graph and protocol unit tests for the new unit; a report
  snapshot covering the executive-summary `## Summary`; and an integration
  execution test confirming the unit runs after ranking and before report build.
  Deterministic artifact snapshots update for the new schema version.
- **Bundled skill:** no runtime skill change. The `/quality` evaluation workflow
  reads `report.md`; the summary is strictly better content in an existing
  section, with no instruction, invocation, or workflow change.
- **Release notes:** `CHANGELOG.md` records the new advice step, payload kind, and
  report summary behavior, and the schema-version bump.
- **Format specification and project model:** no `SPECIFICATION.md`,
  `quality.schema.json`, or `QUALITY.md` change; the model and evaluation meaning
  are unchanged.
- **Durable docs, dependencies, packaging, scaffold:** no README, guide, install,
  dependency, or scaffold change is planned.

## Children

- [Functional spec](0205-evaluation-executive-summary/spec.md) — the advice unit,
  payload kind, protocol move, report rendering, determinism, and schema
  requirements, with per-requirement durable-spec impact.
- [Design doc](0205-evaluation-executive-summary/design.md) — the work-graph
  insertion, protocol context and prompt, payload/schema shape, report rendering,
  and the inference-writes/report-renders boundary, plus alternatives considered.

## Status

`Design`. Spec and design authored; code not started. Durable specs are listed in
Affected artifacts and specified in the functional spec's Durable spec changes,
but not yet edited — they will be brought into sync during implementation, before
**In-Review**.
