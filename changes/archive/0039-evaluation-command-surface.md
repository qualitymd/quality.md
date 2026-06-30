---
type: Change Case
title: Evaluation command surface redesign
description: Reshape the qualitymd evaluation CLI surface around a single noun/verb rule, fold planned coverage into plan.md, separate the report gate, and remove the altitude residue.
status: Done
tags: [cli, evaluation, surface]
timestamp: 2026-06-19T00:00:00Z
---

# Evaluation command surface redesign

A **Change Case** reshaping the `qualitymd evaluation` command surface. The
detail lives in its
[functional spec](0039-evaluation-command-surface/spec.md).

## Motivation

The evaluation surface grew one command at a time, and the shape now fights the
work it models. The whole flow is one linear pipeline — create a run, declare
intended coverage, write judgment records, check readiness, render a report —
but the naming obscures it:

- Every verb repeats a noun the `evaluation` parent already supplies:
  `create-run`, `add-record`, `set-planned-coverage`, `show-status`,
  `build-report`.
- `evaluation show-status` collides conceptually with the top-level
  `qualitymd status`, and the two have different output streams and meanings.
- The three record kinds go through one `add-record` verb even though they
  behave differently on disk — assessments append and number, analyses upsert by
  target, recommendations append numbered Markdown. Calling an upsert "add" hides
  the contract.
- `build-report` overloads rendering with a CI gate (`--fail-at-or-below`), and
  the gate mutates the run folder even when used purely as a pass/fail check.
- There is no way to enumerate runs or inspect a run's records; callers must keep
  the exact paths the CLI printed.
- A removed `model`-altitude leaves dead residue: an `Options.Altitude` field, an
  unreachable guard, and an always-`"subject"` receipt field.
- Output streams and `--file`/stdin handling are applied unevenly across the
  sibling subcommands.

The CLI is pre-1.0 (v0.3.0) and its only consumer is the bundled `/quality`
skill, which ships in lockstep. That makes this the right moment to break the
names freely and settle one coherent rule for the whole surface rather than
accrete more exceptions.

## Scope

Covered:

- A single governing rule for the surface: a concept with more than one operation
  is a **noun** with verb subcommands; a single-operation run-lifecycle action is
  a **bare verb**.
- Rename the run-lifecycle commands: `create-run` → `create`, `show-status` →
  `status`, `build-report` → the `report` noun.
- Promote the three record kinds to resources with honest verbs:
  `assessment add|list`, `analysis set|list`, `recommendation add|list`.
- Add `evaluation list` to enumerate runs, and per-resource `list` to inspect a
  run's written records (the resume/diagnose path).
- Promote `report` to a noun with `build` and a separate, side-effect-free
  `gate`; remove the `build --fail-at-or-below` overload.
- Fold planned coverage into `plan.md` YAML frontmatter (mirroring the
  `QUALITY.md` frontmatter+body shape); delete `planned-coverage.json` and the
  `set-planned-coverage` command. `status` reads and validates the coverage from
  `plan.md` at read time.
- Allow batched record input (a single object or an array per `add`/`set`) and
  optional `--latest` run resolution.
- Normalize the cross-cutting plumbing: one stdout/stderr rule and one
  `--file`/stdin rule across the surface.
- Remove the `model`-altitude residue from the create path and receipts. (The run
  folder regex keeps accepting historical `model-` folders for reading old runs.)

Deferred / non-goals:

- `show`/`remove` verbs on the record resources and `report show` are designed
  for in the namespace but not built until a consumer needs them.
- No change to `QUALITY.md` format semantics, rating vocabulary, evaluation
  semantics, or the judgment content of any record. This is a command-surface and
  artifact-layout change only.
- No change to how reports are rendered or what they contain, beyond the gate
  separation.
- No new telemetry, configuration, or non-evaluation command.

## Affected artifacts

Durable **spec** changes are itemized in the functional spec's
[Durable spec changes](0039-evaluation-command-surface/spec.md#durable-spec-changes)
section. The index below is the full skimmable list, reconciled before
In-Review.

Code (edited only once the case is `In-Progress`):

- [x] [`internal/cli/evaluation.go`](../../internal/cli/evaluation.go) - the cobra
      command tree: rename `create-run`/`show-status`, split `add-record` into the
      record nouns, promote `report` to `build`/`gate`, add `evaluation list`,
      drop `set-planned-coverage`, normalize stdout/stderr and `--file`.
- [x] [`internal/cli/evaluation_test.go`](../../internal/cli/evaluation_test.go) -
      update tests to the new command names and the dropped altitude/subject
      receipt field.
- [x] [`internal/evaluation/create.go`](../../internal/evaluation/create.go) - drop
      the `model`-altitude residue from the create path and receipt.
- [x] [`internal/evaluation/write.go`](../../internal/evaluation/write.go) - back the
      per-resource `add`/`set` verbs and batched (object-or-array) input.
- [x] [`internal/evaluation/planned_coverage.go`](../../internal/evaluation/planned_coverage.go) - remove the `planned-coverage.json` write path; coverage moves to `plan.md`
      frontmatter, validated at read time.
- [x] [`internal/evaluation/report.go`](../../internal/evaluation/report.go) - split
      the `--fail-at-or-below` gate out of `build` into a side-effect-free `gate`.
- [x] [`internal/evaluation/load.go`](../../internal/evaluation/load.go) - source
      planned coverage from `plan.md`; drop altitude.
- [x] [`internal/evaluation/types.go`](../../internal/evaluation/types.go) - remove
      the `Options.Altitude` field and altitude types.
- [x] [`internal/status/status.go`](../../internal/status/status.go) - review for
      evaluation command-name references and the `show-status` vs top-level
      `status` collision; drop subject-receipt residue.

Specs:

- [x] [`specs/cli/index.md`](../../specs/cli/index.md) - rewrite the evaluation
      command list for the new surface.
- [x] [`specs/cli/evaluation-create.md`](../../specs/cli/evaluation-create.md) -
      renamed from `evaluation-create-run.md`; drops the altitude residue.
- [x] [`specs/cli/evaluation-status.md`](../../specs/cli/evaluation-status.md) -
      renamed from `evaluation-show-status.md`; sources coverage gaps from
      `plan.md`.
- [x] [`specs/cli/evaluation-assessment.md`](../../specs/cli/evaluation-assessment.md),
      [`specs/cli/evaluation-analysis.md`](../../specs/cli/evaluation-analysis.md),
      and
      [`specs/cli/evaluation-recommendation.md`](../../specs/cli/evaluation-recommendation.md) - split from `evaluation-add-record.md` with honest verbs and batched
      input.
- [x] [`specs/cli/evaluation-report.md`](../../specs/cli/evaluation-report.md) -
      renamed from `evaluation-build-report.md`; splits `build` from `gate`.
- [x] Deleted `specs/cli/evaluation-set-planned-coverage.md`; the coverage
      contract moves into `plan.md` frontmatter.
- [x] [`specs/cli/evaluation-list.md`](../../specs/cli/evaluation-list.md) -
      enumerate runs.
- [x] [`specs/evaluation-records.md`](../../specs/evaluation-records.md) - fold
      planned coverage into `plan.md` frontmatter, update the run-folder layout,
      and update command-name references.
- [x] [`specs/cli.md`](../../specs/cli.md) - review for evaluation command-name
      references and confirm the cross-cutting stdout/stderr and `--file`
      contract covers the normalized surface.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md) - update the skill's functional spec: the evaluation procedure, workflow
      diagram, and every old command name it references.

Docs:

- [x] [`skills/quality/resources/cli-quick-reference.md`](../../skills/quality/resources/cli-quick-reference.md) - rewrite the evaluation command table, decision trees, and workflows.
- [x] [`skills/quality/modes/evaluate.md`](../../skills/quality/modes/evaluate.md) -
      update the procedure to the new commands and the `plan.md` coverage fold.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - update the
      `qualitymd evaluation …` command-surface reference in the skill entry point
      if present; no live command list remained there.
- [x] [`README.md`](../../README.md) - reviewed for evaluation command examples;
      no live command examples required an update.
- [x] [`install.md`](../../install.md) - update the verbatim
      `qualitymd evaluation … --help` block to the new command names.
- [x] [`docs/guides/use-quality-skill.md`](../../docs/guides/use-quality-skill.md) -
      update the same verbatim `--help` block.
- [x] [`docs/guides/cli-design.md`](../../docs/guides/cli-design.md) - update the
      `evaluation create-run` worked example of the noun/verb grammar.
- [x] [`docs/guides/write-functional-specs.md`](../../docs/guides/write-functional-specs.md) - update the `evaluation create-run` example/link (the spec is being
      renamed to `evaluation-create.md`).
- [x] [`CHANGELOG.md`](../../CHANGELOG.md) - add the 0039 entry; the new command
      names supersede the historical `evaluation build-report` references.

`SPECIFICATION.md` is **not** affected: this case changes the CLI surface and
artifact layout, not the `QUALITY.md` format or evaluation semantics.

## Children

- [Functional spec](0039-evaluation-command-surface/spec.md) - what the new
  surface must require.
- [Design doc](0039-evaluation-command-surface/design.md) - how the cobra tree,
  shared helpers, `plan.md` coverage fold, and report build/gate split are built,
  and the rejected alternatives.

## Status

`Done`. Implemented and archived.
