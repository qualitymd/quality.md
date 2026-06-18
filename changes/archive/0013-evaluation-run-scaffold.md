---
type: Change
title: Evaluation run scaffold
description: Move evaluation run-folder scaffolding and numbering out of the skill into a deterministic qualitymd evaluation create-run command.
status: Done
tags: [evaluation, specs, cli, skill]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation run scaffold

A unit of work that lifts evaluation run-folder **scaffolding** — resolving the
evaluation directory, computing the next run number, and seeding the run folder —
out of the [`/quality` skill prompt](../../skills/quality/SKILL.md#evaluation) and
into a deterministic `qualitymd evaluation create-run` command. Detail lives in the
child:

- [Functional spec](0013-evaluation-run-scaffold/spec.md) — what
  `qualitymd evaluation create-run` must do.
- [Design doc](0013-evaluation-run-scaffold/design.md) — how
  `qualitymd evaluation create-run` is built.

## Motivation

Today the skill hand-creates the run folder and computes the next `NNNN` itself.
That is error-prone judgment work the deterministic surface should own — a real
numbering collision (two runs claiming the same number) has already occurred.
Per the [evaluation-record contract](0012-evaluation-record-format/spec.md), the
CLI owns run-folder numbering and creation; this change builds the first command
that does it, so the skill stops counting and laying out folders by hand.

## Scope

Covered: a spec-only change stating the behavior of `qualitymd evaluation create-run` —
resolving the evaluation directory, deterministically computing the next run
number shared across all altitudes, and creating the run folder with its required
subdirectories and seed files (`model.md`, `design.md`, `plan.md`).

Deferred: writing assessment/analysis/recommendation records (change 0014,
`qualitymd evaluation add-record`) and rendering `report.md`/`report.json` and gating
(change 0015, `qualitymd evaluation build-report`). No code here — `Design`
authorizes only this change's own folder plus the change parent, bundle index,
and bundle log.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches `In-Review`
(not now — `Design` authorizes only this change's own folder plus the change
parent, bundle index, and bundle log):

- [x] `specs/cli/evaluation-create-run.md` (new command sub-spec) — the
      `qualitymd evaluation create-run` command; list it in
      [`specs/cli.md`](../../specs/cli.md) and [`specs/cli/index.md`](../../specs/cli/index.md).
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — update the evaluation workflow/reporting contract to say run folders are
      scaffolded by `qualitymd evaluation create-run`.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — replace the
      hand-rolled run-folder creation and numbering steps in the
      [Evaluation](../../skills/quality/SKILL.md#evaluation) procedure with a call to
      `qualitymd evaluation create-run`.
- [x] [`README.md`](../../README.md) — note `evaluation create-run` among the planned
      evaluation surface.

## Status

`Done`. Implemented and archived after implementing `qualitymd evaluation create-run` with deterministic shared run numbering and run-folder scaffolding.
