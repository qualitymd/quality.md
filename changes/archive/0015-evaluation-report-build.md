---
type: Change
title: Evaluation status and report build
description: Make evaluation run status and report rendering deterministic CLI surfaces, with a CI gate.
status: Done
tags: [evaluation, specs, cli, skill]
timestamp: 2026-06-17T00:00:00Z
---

# Evaluation status and report build

A unit of work that makes an evaluation run's current status inspectable, turns
the Evaluation Report into a deterministic **rendering** the `qualitymd` CLI
derives from a run's records, and adds a CI gate over the result. Detail lives in
the child:

- [Functional spec](0015-evaluation-report-build/spec.md) — what
  `qualitymd evaluation show-status` and `qualitymd evaluation build-report` must
  do.
- [Design doc](0015-evaluation-report-build/design.md) — how the commands inspect
  a run's renderability and derive `report.md` and `report.json` from its
  records, with the CI gate.

## Motivation

Today the `/quality` skill hand-authors both `report.md` and `report.json`,
duplicating the same result across two files and risking drift between them — and
against the `assessments/*.json` and `analysis/*.json` records they summarize.
The deterministic CLI already owns writing every evaluation record (see
[0012](0012-evaluation-record-format.md)); the report is the one record still
authored by hand. Deriving both report files mechanically from the records gives
a single source of truth, removes the drift, and yields a stable artifact a gate
can branch on in CI.

## Scope

Covered: the functional spec for `qualitymd evaluation show-status <run>`, which
reports whether a run can be rendered, and
`qualitymd evaluation build-report <run>`, which reads the run's
`assessments/*.json`, `analysis/*.json`, and `recommendations/*.md` and derives
`report.md` (human-readable) and `report.json` (machine-readable, per the
[0012 contract](0012-evaluation-record-format/spec.md#reportjson)). Report
rendering is idempotent, and `--fail-at-or-below <level>` gates the exit code so
the report drops into CI. The spec also fixes behavior when records are missing
or incomplete.

Deferred: scaffolding a run (change 0013, `qualitymd evaluation create-run`) and writing
records into it (change 0014, `qualitymd evaluation add-record`). Roll-up **ratings** are
skill judgment recorded in the analysis records by 0014; this command **renders**
them and never infers a rating. No code here — `Design` authorizes only this
change's own folder plus the change parent, bundle index, and bundle log.

## Affected specs & docs

Created or updated during `In-Progress`, before this change reaches `In-Review`
(not now — `Design` authorizes only this change's own folder plus the change
parent, bundle index, and bundle log):

- [x] `specs/cli/evaluation-show-status.md` (new concept) — the durable spec for the
      `evaluation show-status` command; list it in
      [`specs/cli.md`](../../specs/cli.md) and [`specs/cli/index.md`](../../specs/cli/index.md).
- [x] `specs/cli/evaluation-build-report.md` (new concept) — the durable spec
      for the `evaluation build-report` command; list it in
      [`specs/cli.md`](../../specs/cli.md) and [`specs/cli/index.md`](../../specs/cli/index.md).
- [x] [`README.md`](../../README.md) — move `evaluation show-status` and
      `evaluation build-report --fail-at-or-below` from *planned* to built.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — update the reporting contract to say reports are rendered through
      `qualitymd evaluation build-report` after `qualitymd evaluation show-status`.
- [x] [`specs/skills/quality-skill/examples/`](../../specs/skills/quality-skill/examples/)
      — refresh the reference `report.md` and `report.json` outputs if deterministic
      rendering changes their shape or wording.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — call
      `evaluation build-report` instead of hand-authoring `report.md`/`report.json`.

## Status

`Done`. Implemented and archived after implementing `qualitymd evaluation show-status` and `build-report` over a shared renderability gate, with deterministic `report.md`/`report.json` and the `--fail-at-or-below` CI gate.
