---
type: Change Case
title: Quality log
description: Add a curated quality log — dated entries under quality/log/ that record meaningful changes to a QUALITY.md model and cross-link the evaluation evidence behind them.
status: Done
tags: [skill, quality, evaluation, changes]
timestamp: 2026-06-22T00:00:00Z
---

# Quality log

A **Change Case** capturing the *why* and *status* for a **quality log**: a
curated, evidence-linked history of meaningful changes to a QUALITY.md model,
written by the `/quality` skill as dated entries under `quality/log/`. The detail
lives in its [functional spec](0050-quality-log/spec.md).

> **Done.** A folder of date-named entries, written by `setup` and `improve`,
> reconciled by `wizard`, with the format contract in `SKILL.md` and the
> meaningful-change judgment in the authoring guide. The case is
> **convention-first**: the skill writes entries directly, with no `qualitymd log`
> CLI command and no standalone artifact-spec yet (see [Scope](#scope)). The
> durable rationale is carried on the [`/quality` skill spec](../../specs/skills/quality-skill/quality-skill.md#quality-log).

## Motivation

A `QUALITY.md` is a point-in-time snapshot, and `git log` already records its
literal diffs. What neither captures is the **judgment layer**: *why* an apex
requirement tightened, *which* evaluation surfaced the gap a new Factor closes,
*whether* a criterion moved as deliberate recalibration or as drift. That
rationale is exactly what the authoring guide's *learn loop* depends on, and
today it evaporates into commit messages that scroll away.

A quality log fills that gap as a **model-evolution timeline cross-linked to
evidence** — each meaningful model change recorded with its reason and a link to
the evaluation run and recommendation that motivated it. This is genuinely useful
for the skill's target audience: a project adopting QUALITY.md has a `QUALITY.md`
and `quality/evaluations/` but no `changes/` bundle, so the quality log is the
model's own curated history for them. It is distinct from `changes/log.md`, which
tracks Change Cases on *this* repository.

The value is the rationale and the evidence cross-link, not a second copy of the
diff — so the log records *why the model changed*, and defers the complete
diff history to git.

## Scope

Covered:

- A quality log written by the `/quality` skill as dated entries under
  `quality/log/`, one file per meaningful model change.
- The format contract (path, date-naming, entry frontmatter, runtime-not-OKF
  status) in the always-loaded `SKILL.md`.
- The meaningful-change judgment (what to log, what not to) in the authoring
  guide, where model changes are actually made.
- `setup` seeding an inaugural entry; `improve` appending one per confirmed model
  change; `wizard` surfacing model history and reconciling out-of-band drift.
- Updating the `/quality` skill functional spec to describe the new artifact.

Deferred / non-goals:

- **No `qualitymd log` CLI command.** This case is convention-first: the skill
  writes entries directly to validate the shape before any CLI surface or
  standalone artifact-spec is committed. If the convention proves out, promoting
  the mechanics into the CLI (so numbering and an index can be CLI-owned) is a
  later case.
- **No `.quality/config.yaml` `logDir` key.** The path defaults to `quality/log/`;
  a config key parallels `evaluationDir` only once the surface graduates to the
  CLI.
- **No standalone artifact-spec** (like `evaluation-records.md`) and **no
  machine-queryable index file** in `quality/log/` — both belong with the CLI
  phase.
- **No automatic git-based backfill.** `wizard` *flags* out-of-band drift and
  offers a backfill route; the user (via `improve`/authoring) performs it.
- No change to evaluation records, reports, or the evaluation semantics. The log
  *references* evaluation runs; it does not duplicate them.

## Affected artifacts

### Code

- None. This case adds no Go code; the `qualitymd log` CLI command is explicitly
  deferred (see [Scope](#scope)).

### Durable specs

- `specs/skills/quality-skill/quality-skill.md` — add a quality-log subsection
  describing the artifact: `quality/log/` location, date-named entries,
  runtime-not-OKF status, which modes write, and the `wizard` reconciliation
  behavior. Keep the standalone artifact-spec deferred.
- `specs/log.md` — log the quality-skill spec change.

### Bundled skill

- `skills/quality/SKILL.md` — add the "Quality Log" format contract; add the log
  to the run-frame mutation enumeration; add the write-authority hard rule.
- `skills/quality/guides/authoring.md` — add the meaningful-change taxonomy under
  "When to update QUALITY.md".
- `skills/quality/modes/setup.md` — seed the inaugural entry after guided
  population; add the log to the run-frame artifacts.
- `skills/quality/modes/improve.md` — append an entry per confirmed model change;
  name it among changed artifacts in the delta report.
- `skills/quality/modes/wizard.md` — add a model-history status line and an
  out-of-band-drift reconciliation option.

### Durable docs, scaffold, and examples

- `README.md` — no impact; the README does not enumerate runtime skill outputs.
- `docs/guides/use-quality-skill.md` — added a quality-log sentence to the skill
  mode rundown, which already enumerates what `setup`/`evaluate` write, so it does
  not read as stale.
- Other `docs/guides/` — no impact (the change-case and authoring guidance already
  cover the relevant editing rules).
- Scaffold/install files — no impact.

## Children

- [Functional spec](0050-quality-log/spec.md) — what the quality log must do.

No design doc: the shape is settled in discussion and the durable rationale is
carried on the spec's requirements (folder-vs-file, date-naming, guidance split).

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). The spec was
settled (no design doc needed); the durable quality-skill spec subsection and the
bundled skill edits enumerated under [Affected artifacts](#affected-artifacts)
landed, and the case was archived. The durable rationale lives in the
[`/quality` skill spec](../../specs/skills/quality-skill/quality-skill.md#quality-log).
