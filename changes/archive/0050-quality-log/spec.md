---
type: Functional Specification
title: Quality log — functional spec
description: What the /quality skill's quality log must do — dated entries under quality/log/ recording meaningful, evidence-linked changes to a QUALITY.md model.
tags: [skill, quality, evaluation]
timestamp: 2026-06-22T00:00:00Z
---

# Quality log — functional spec

Companion to the [Quality log](../0050-quality-log.md) change case. This spec
states *what* the quality log must do; the durable behavior lands in the
[`/quality` skill spec](../../../specs/skills/quality-skill/quality-skill.md) and
the bundled skill prompt (`skills/quality/`). The QUALITY.md format and
evaluation semantics remain governed by
[`SPECIFICATION.md`](../../../SPECIFICATION.md); the quality log is a runtime
output of the skill, not a format concept.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

A `QUALITY.md` snapshot and `git log` together record *what* the model is and
*how* it changed, but not *why*: which evaluation surfaced a gap, whether a
criterion moved by recalibration or drift, what a new Factor was reacting to.
That rationale is what the authoring guide's *learn loop* runs on, and it is lost
once a commit scrolls away. The quality log is a curated, evidence-linked
timeline of meaningful model changes that preserves the *why* and links each
change to the evaluation evidence behind it. It is the model's own history for a
project that has a `QUALITY.md` and `quality/evaluations/` but no `changes/`
bundle.

The log is deliberately **curated, not complete**: hand edits to `QUALITY.md`
bypass the skill, so the log cannot be exhaustive and does not try to be. Git
remains the complete diff history; the log carries the judgment git cannot.

## Scope

Covered: a skill-written quality log under `quality/log/`; its on-disk shape; the
meaningful-change criterion; which modes write and reconcile it.

Deferred:

- A `qualitymd log` CLI command and a `.quality/config.yaml` `logDir` key —
  convention-first defers all CLI mechanics; the path defaults to `quality/log/`.
- A standalone artifact-spec (like
  [`evaluation-records.md`](../../../specs/evaluation-records.md)) and any
  machine-queryable index file inside `quality/log/`.

Non-goals:

- The log is **not** an evaluation record and **not** a defect backlog. It
  records model changes only and *references* evaluation runs rather than copying
  them.
- The log is **not** a second copy of the git diff; it records the *why*, not the
  literal change.

## Requirements

### The log artifact

- The skill **MUST** record meaningful changes to a QUALITY.md model as entries
  under `quality/log/`, a sibling of the resolved evaluation directory. Each
  meaningful change is **one entry file**.

  >> Rationale: a folder of independent files mirrors `quality/evaluations/` for
  >> one mental model of the runtime root, and avoids the append conflicts a
  >> single shared log file would create when concurrent branches or agents add
  >> entries. — 0050

- Each entry **MUST** be named `YYYY-MM-DD-<slug>.md`, where the date is the day
  the change was made and `<slug>` is a short kebab-case summary. The skill
  **MUST NOT** assign a global sequential counter to entries.

  >> Rationale: a date prefix orders the log chronologically without a shared
  >> counter. Skill-side sequential numbering is exactly what drifted and produced
  >> a run-number collision before
  >> ([`evaluation-records.md`](../../../specs/evaluation-records.md)), which is
  >> why numbering is CLI-owned; with no CLI in this surface, date-naming sidesteps
  >> the hazard entirely. — 0050

- `quality/log/` is a **runtime artifact, not an OKF bundle**. It **MUST NOT**
  carry OKF `index.md`, `schema.md`, or `log.md` semantics, and entry frontmatter
  is machine metadata, not OKF concept frontmatter.

  >> Rationale: same classification the evaluation run folders carry; runtime
  >> outputs in the evaluated repository are not OKF concepts. — 0050

- Each entry **MUST** carry small machine-readable frontmatter and a prose
  rationale body. The frontmatter records the change kind, the model target it
  affects, and — when the change came from an evaluation — the source run and
  recommendation it traces to. The body states *why* the change was made.

  >> Rationale: the cross-link to the evaluation run and recommendation is the
  >> log's differentiator over `git log`; without it the entry is just a diff in
  >> prose. — 0050

### What is meaningful

- The skill **MUST** log a change that alters what the model *is* or *how it
  judges*: adding, removing, or renaming an Area, Factor, or Requirement;
  changing the rating scale, a criterion, or a relative weight; shifting scope;
  changing the apex or required margin; or applying an evaluation recommendation.
  An entry **SHOULD** state whether a criterion change is deliberate recalibration
  or a drift correction.

- The skill **MUST NOT** log Markdown-body wording, typo, or formatting changes,
  nor evaluated-source fixes that do not change the model.

  >> Rationale: those are not model changes; git already records them, and logging
  >> them turns a curated timeline into noise. — 0050

- The skill **SHOULD** write **one entry per coherent change** — a confirmed
  `improve` apply, or the initial model population — rather than one entry per
  field touched.

  >> Rationale: the unit of record is the decision, not the edit; per-field
  >> entries fragment one rationale across many files. — 0050

### Who writes and reconciles

- `setup` **MUST** seed an inaugural entry after guided first population, recording
  model creation and the initial model shape.

- `improve` **MUST** append an entry for each confirmed model change, cross-linking
  the evaluation run and recommendation it came from. Writing the entry **MUST
  NOT** require confirmation beyond the user's existing confirmation of the change
  itself; the entry's rationale is the rationale already shown in the decision
  brief.

- `evaluate` **MUST NOT** write to the quality log, and no mode other than `setup`
  and `improve` writes entries.

  >> Rationale: evaluations own `quality/evaluations/`; keeping the log to model
  >> changes only — referencing runs, never duplicating them — is what stops it
  >> becoming a second evaluation record. — 0050

- `wizard` **MUST** remain read-only with respect to the log. It **SHOULD** surface
  model history (the latest entry) in its status output, and when the model has
  changed out of band since the last logged entry it **SHOULD** classify that under
  its existing *needs reconciliation* readiness and offer a backfill route. The
  backfill itself is performed by a confirmed `improve`/authoring workflow, not by
  `wizard`.

  >> Rationale: the log is curated, not complete; `wizard` is where the gap left
  >> by hand edits gets caught and routed for repair without `wizard` itself
  >> mutating anything. — 0050

### Mutation accounting

- The run frame's mutation enumeration **MUST** include the quality log as a
  distinct mutation surface, so a write to it is visible before it happens.

## Durable spec changes

### To add

None. A standalone artifact-spec governing `quality/log/` is deferred until the
surface graduates to the CLI (see [Scope](#scope)).

### To modify

- `specs/skills/quality-skill/quality-skill.md` — add a quality-log subsection
  describing the artifact: its `quality/log/` location and date-named entries, its
  runtime-not-OKF status, which modes write it, and the `wizard` reconciliation
  behavior (per the requirements above). Promote the curated-not-complete and
  date-naming rationale into the spec.

### To rename

None.

### To delete

None.
