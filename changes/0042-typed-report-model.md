---
type: Change Case
title: Typed report model
description: Replace stringly typed and implicit evaluation-report states with explicit typed concepts for ratings, local target state, next steps, lifecycle state, run gaps, rigor, evaluation level, missing metadata, and path identities.
status: In-Review
tags: [evaluation, report, records, types]
timestamp: 2026-06-20T00:00:00Z
---

# Typed report model

A Change Case for making the evaluation reporting model explicit and typed. It
continues the finding-severity cleanup by applying the same treatment to the
other report concepts that carry control-flow meaning.

> **In-Review.** Implementation, durable spec updates, fixtures, and verification
> are complete and ready for review.

## Motivation

The report model currently encodes several semantic states as raw strings,
booleans, nil pointers, or repeated `[]string` path slices. That makes invalid
states easy to construct and pushes important report policy into scattered string
comparisons. Report output must keep distinctions such as rated vs not assessed,
structural vs missing evidence, active vs superseded, and recommendation vs no
next step. Those distinctions should be represented directly in the model.

## Scope

Covered: typed rating-result kinds, explicit local target rating state, typed
report next step, typed evaluation-run gap kinds, typed lifecycle state for
report digests, typed rigor and evaluation level, structured missing metadata,
and target/factor path helper types. The JSON record shapes stay stable where
they are inputs, while report JSON may grow explicit state objects.

Deferred / non-goals: no change to QUALITY.md frontmatter, no change to rating
scale semantics, no new report feature, and no compatibility aliases for
invalid typed values.

## Affected artifacts

### Code

- [x] `internal/evaluation/types.go` — introduce typed report/evaluation
      concepts and path helpers.
- [x] `internal/evaluation/write.go` — validate typed rating results and record
      payloads.
- [x] `internal/evaluation/load.go` — report invalid typed states as
      renderability gaps and use typed gap routing.
- [x] `internal/evaluation/report.go` — build and render reports from explicit
      state objects rather than raw strings, booleans, and nil conventions.
- [x] `internal/evaluation/evaluation_test.go` — cover invalid states, explicit
      structural local ratings, lifecycle state, next-step state, and fixture
      rendering.

### Durable specs

See the functional spec's
[Durable spec changes](0042-typed-report-model/spec.md#durable-spec-changes)
for the per-requirement breakdown.

- [x] [`specs/evaluation-records.md`](../specs/evaluation-records.md)
- [x] [`specs/skills/quality-skill/quality-skill.md`](../specs/skills/quality-skill/quality-skill.md)
- [x] [`specs/cli/evaluation-status.md`](../specs/cli/evaluation-status.md)
- [x] [`specs/cli/evaluation-report.md`](../specs/cli/evaluation-report.md)
- [x] [`specs/log.md`](../specs/log.md)

No `SPECIFICATION.md` change is expected: this tightens runtime evaluation
records and report rendering, not QUALITY.md document semantics.

### Durable docs

- [x] `skills/quality/SKILL.md`, `skills/quality/modes/evaluate.md`, and
      `skills/quality/resources/cli-quick-reference.md` reviewed; no runtime
      instruction or sample-payload edits needed beyond the durable skill spec.

### Fixtures

- [x] `specs/skills/quality-skill/examples/0001-subject-quality-eval/` reports
      and any affected assessment/analysis records.

## Children

- [Functional spec](0042-typed-report-model/spec.md) — what the change must do.
- [Design doc](0042-typed-report-model/design.md) — how the typed model is
  implemented.
