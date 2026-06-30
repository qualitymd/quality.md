---
type: Change Case
title: Evaluation Advice
description: Make evaluation advice required through finding rankings, recommendations, recommendation rankings, and generated recommendation reports.
status: Done
tags: [evaluation, advice, recommendations, reports]
timestamp: 2026-06-27T00:00:00Z
---

# Evaluation Advice

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0150-evaluation-advice/spec.md) - what the case must do.
- [Design doc](0150-evaluation-advice/design.md) - how it's built, and why.

## Motivation

Evaluation should not stop at ratings. A complete run should tell the user what
quality-management move is most warranted next: change the evaluated entity,
improve evidence, improve the Model, review whether the bar should rise,
preserve a strength, monitor a condition, or intentionally defer with a review
trigger. The current durable contracts still treat generated recommendations as
out of scope or optional, which leaves evaluation short of the quality loop's
improve step and forces agents to infer advice outside the persisted run.

This case makes Advice a required late evaluation phase. Findings remain the
evidence layer; a Finding ranking orders existing Findings for advice synthesis;
Recommendations express next quality-management moves; and a Recommendation
ranking orders those moves by expected impact and confidence. Generated reports
surface Top Findings and Top Recommendations while keeping full recommendation
handoff artifacts under the evaluation run.

## Scope

Covered:

- Add `FindingRankingResult`, `RecommendationResult`, and
  `RecommendationRankingResult` payload kinds.
- Require Advice payloads and finding coverage accounting for reportable
  evaluations.
- Require recommendations to be quality-management advice, not necessarily work
  items.
- Require recommendation impact and confidence, with no effort, ROI, or
  quick-win metadata.
- Render `Top Findings` and `Top Recommendations` in `report.md`, capped at 10
  rows each, and always link to `recommendations.md`.
- Generate `recommendations.md` plus one human-readable recommendation detail
  file per Recommendation.
- Keep generated Markdown recommendation artifacts human-first; defer
  machine-readable frontmatter or YAML appendices.
- Preserve domain-agnostic language and include a non-software example in the
  affected specs or examples.

Deferred:

- Applying recommendations or changing `/quality improve` follow-up behavior
  beyond consuming the new artifacts later.
- Adding effort, ROI, quick-win, backlog-priority, or numeric scoring fields.
- Adding CLI commands that generate judgment, such as
  `qualitymd evaluation recommend`.
- Adding frontmatter or YAML appendices to generated recommendation Markdown.
- Backward compatibility for evaluation runs without Advice payloads.

## Affected artifacts

Derived by sweeping for Advice, recommendations, report generation, evaluation
payload kinds, evaluation status/reportability, and recommendation follow-up.

**Code**

- [x] `internal/evaluation/data_contract.go` - add Advice payload validation,
      references, coverage accounting, impact enums, and no-effort constraints.
- [x] `internal/evaluation/data.go` - add generated examples for new payload
      kinds.
- [x] `internal/evaluation/evaluation-data.schema.json` - generated schema
      output for the new payload kinds.
- [x] `internal/evaluation/status.go` and related coverage/status code - require
      Advice payloads and complete finding coverage before reportability.
- [x] `internal/evaluation/report_tree.go` - render Top Findings, Top
      Recommendations, recommendation index, and recommendation detail files.
- [x] CLI tests under `internal/cli/` and evaluation tests under
      `internal/evaluation/` - cover schema, status, report artifacts, and
      deterministic rendering.

**Format spec and durable specs** (substance in the [functional spec](0150-evaluation-advice/spec.md))

- [x] `SPECIFICATION.md` - make Advice required for complete evaluation and
      replace optional/generated-recommendation-v0 wording.
- [x] `skills/quality/resources/SPECIFICATION.md` - bundled spec symlink to the
      root specification.
- [x] `specs/evaluation/evaluation.md` - shared Evaluation invariants and
      non-goals.
- [x] `specs/evaluation/protocol.md` and `specs/evaluation/orchestration.md` -
      Advice phase placement after roll-up and before report build.
- [x] `specs/evaluation/records/payload-kinds.md` - new payload kind contracts,
      references, enums, and coverage accounting.
- [x] `specs/evaluation/records/data-layout.md` - generated recommendation
      artifact paths.
- [x] `specs/evaluation/reports/report-tree.md` - Top Findings, Top
      Recommendations, recommendation index, and detail rendering.
- [x] `specs/cli/evaluation-data.md` - data command behavior for new kinds.
- [x] `specs/cli/evaluation-status.md` - reportability requires Advice and
      coverage.
- [x] `specs/cli/evaluation-report.md` - report build renders Advice from
      persisted payloads only.
- [x] `specs/skills/quality-skill/evaluation.md` and
      `specs/skills/quality-skill/workflows/evaluate.md` - skill judgment flow
      now produces required Advice.
- [x] `specs/skills/quality-skill/reporting.md` - reporting closeout includes
      Advice artifacts.
- [x] `specs/skills/quality-skill/recommendation-follow-up.md` and guide specs -
      later follow-up can resolve recommendation detail artifacts without
      assuming every recommendation is a local work item.
- [x] Relevant durable spec logs - record the contract change.

**Durable docs / bundled skill runtime**

- [x] `skills/quality/SKILL.md` - runtime evaluation guidance for required
      Advice, ranking, coverage accounting, and report closeout.
- [x] `skills/quality/workflows/evaluate.md` - Advice phase steps and QC.
- [x] `skills/quality/guides/recommendation-follow-up.md` - keep follow-up
      compatible with recommendation artifacts while preserving confirmation
      gates.
- [x] `README.md` - update evaluation wording if needed to mention generated
      recommendations as part of the run.
- [x] `CHANGELOG.md` - release-note entry for the Advice/reportability change.
- [x] `changes/log.md` and `changes/index.md` - Change Case lifecycle.

## Status

`Done`. Implementation, durable specs, runtime skill guidance, generated schema,
release notes, and focused tests are complete. `go test ./...` and
`mise run fmt-md-check` pass.
