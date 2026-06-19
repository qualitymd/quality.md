---
type: Change Case
title: CLI status snapshot command
description: Add a qualitymd status [--json] command that emits a deterministic project-state snapshot so the /quality wizard can route without hand-parsing QUALITY.md or reading evaluation reports.
status: Done
tags: [cli, skill, wizard]
timestamp: 2026-06-19T00:00:00Z
---

# CLI status snapshot command

A **Change Case** capturing the *why* and *status* for adding a deterministic
`qualitymd status` snapshot command. The detail lives in its
[functional spec](0030-cli-status-command/spec.md).

> **Done.** Implementation and durable spec/doc updates are complete; the change
> is archived.

## Motivation

The `/quality wizard` mode is meant to be a fast wayfinder: probe a few cheap
signals, classify readiness, and offer next steps. In practice it had to
reconstruct project state by hand — parsing `QUALITY.md` frontmatter with
ad-hoc shell to count targets/factors/requirements, opening evaluation
`report.md` bodies to learn the latest rating and open recommendations, and
resolving build/install paths. That is slow, brittle, and — worst of all — does
mechanical work *in the skill*, inverting the project's core split: the CLI owns
deterministic mechanical work, the skill owns judgment.

A first-class `qualitymd status [--json]` command fixes the root cause. It emits
a single deterministic snapshot of project state — model validity, model shape
(target / factor / requirement counts and source coverage), evaluation run
history (count, latest run, incomplete/stale runs), and open recommendation
counts — so any consumer (the wizard, CI, other tooling) reads structured data
instead of scraping artifacts. The wizard's recently leaned-out probe can then
collapse to one `--json` call, and the finer maturity classification it now
defers becomes cheaply available without report-body reads.

## Scope

Covered: a new read-only `qualitymd status` command with a human form and a
`--json` form; its functional spec and durable CLI spec; and updating the
`/quality` skill consumers (the wizard probe and the CLI quick reference) to use
it.

Deferred / non-goals: no change to evaluation record formats or the model
schema; the command only *reports* existing state. The wizard's read-only
routing behavior and its lean-probe contract are already in place (routine
change to `modes/wizard.md`); this case supplies the structured source those
signals should come from. No new lint semantics.

## Affected specs & docs

- [x] [`specs/cli/index.md`](../../specs/cli/index.md) and a new
      `specs/cli/status.md` sub-spec — register and specify the `status` command
      (arguments, human vs. `--json` output, exit codes, determinism). Detailed
      in the functional spec's **Durable spec changes**.
- [x] [`specs/cli.md`](../../specs/cli.md) — list the new command in the CLI surface
      overview if it enumerates commands.
- [x] [`README.md`](../../README.md) — add `status` to the documented command set.
- [x] [`skills/quality/resources/cli-quick-reference.md`](../../skills/quality/resources/cli-quick-reference.md)
      — add a `status` row and prefer `--json` for agent consumption.
- [x] [`skills/quality/modes/wizard.md`](../../skills/quality/modes/wizard.md) —
      switch the batched probe to `qualitymd status --json` once the command exists,
      replacing the current per-signal shell checks.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — add `status` to the
      CLI Reference block and reconcile the declared CLI SemVer range with the
      release that ships the command.

## Children

- [Functional spec](0030-cli-status-command/spec.md) — what the
  `qualitymd status` snapshot command must provide.
- [Design doc](0030-cli-status-command/design.md) — how the snapshot command
  should compose existing CLI-owned mechanics.

## Status

`Done`. The command, durable specs/docs, and skill consumers are synced and
verified.
