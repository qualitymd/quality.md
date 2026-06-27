---
type: Functional Specification
title: /quality quality changelog
description: Component spec for the /quality skill's convention-first quality changelog under .quality/changelog/.
tags: [skill, quality, changelog]
timestamp: 2026-06-23T00:00:00Z
---

# /quality quality changelog

This spec owns the `/quality` skill's convention-first quality changelog: dated
entries under `.quality/changelog/` in the workspace quality data directory that
record meaningful, evidence-linked changes to a QUALITY.md model. It composes
the shared contracts in the parent
[/quality skill](quality-skill.md) spec and is written by confirmed direct
model-authoring changes or recommendation-apply workflows after setup.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Quality changelog

A `QUALITY.md` snapshot and `git log` together record *what* the model is and
*how* it changed, but not *why*: which evaluation surfaced a gap, whether a
criterion moved by recalibration or drift, what a new Factor was reacting to.
That rationale is what the [learn loop](evaluation.md#evaluation-workflow) runs on, and it is
lost once a commit scrolls away. The **quality changelog** is a curated,
evidence-linked timeline of meaningful model changes the skill maintains under
`.quality/changelog/`: it preserves the *why* and links each change to the evaluation
evidence behind it. It is the model's own history for a project that has a
`QUALITY.md` and `.quality/evaluations/` but no `changes/` bundle of its own.

The changelog is deliberately **curated, not complete**. Hand edits to
`QUALITY.md` bypass the skill, so the changelog cannot be exhaustive and does
not try to be — git remains the complete diff history, and the changelog carries
the judgment git cannot.
It is **not** an evaluation record and **not** a defect backlog: it records model
changes only and *references* evaluation runs rather than copying them.

The changelog is a **runtime output** of the skill, not a QUALITY.md format
concept; the format and evaluation semantics remain governed by
[`SPECIFICATION.md`](../../../SPECIFICATION.md). A `qualitymd changelog` CLI
command, a `.quality/config.yaml` `changelogDir` key, a standalone
artifact-spec, and any machine-queryable index file inside `.quality/changelog/`
are all deferred until the surface graduates to the CLI (see
[Deferred](quality-skill.md#deferred)); this section is the convention-first
contract the skill writes against in the meantime.

### The changelog artifact

The skill **MUST** record meaningful changes to a QUALITY.md model as entries
under the resolved workspace's `.quality/changelog/` directory relative to the
selected `QUALITY.md`. Each meaningful change is **one entry file**.

> Rationale: a folder of independent files mirrors `.quality/evaluations/` for one
> mental model of the runtime root, and avoids the append conflicts a single
> shared log file would create when concurrent branches or agents add
> entries. — 0050

Each entry **MUST** be named `YYYY-MM-DDTHHMMSSZ-<slug>.md`, where the timestamp
is the sortable UTC time the change was made and `<slug>` is a short kebab-case
summary. The skill **MUST NOT** assign a global sequential counter to entries.

> Rationale: a UTC timestamp prefix orders multiple same-day changelog entries
> chronologically without a shared counter. Skill-side sequential numbering is
> exactly what drifted and produced a run-number collision before
> ([`qualitymd evaluation create`](../../cli/evaluation-create.md)), which is why
> evaluation run numbering is CLI-owned; with no CLI in this surface,
> timestamp-naming sidesteps the hazard entirely. — 0050, 0145

`.quality/changelog/` is a **runtime artifact, not an OKF bundle**. It **MUST NOT** carry
OKF `index.md`, `schema.md`, or `log.md` semantics, and entry frontmatter is
machine metadata, not OKF concept frontmatter.

> Rationale: same classification the evaluation run folders carry; runtime outputs
> in the evaluated repository are not OKF concepts. — 0050

Each entry **MUST** carry small machine-readable frontmatter and a prose
rationale body. The frontmatter records the change kind, the model target it
affects, and — when the change came from an evaluation — the source run and
recommendation it traces to. The body states *why* the change was made.

> Rationale: the cross-link to the evaluation run and recommendation is the
> changelog's differentiator over `git log`; without it the entry is just a diff in
> prose. — 0050

### What is meaningful

The skill **MUST** log a change that alters what the model *is* or *how it
judges*: adding, removing, or renaming an Area, Factor, or Requirement; changing
the rating scale, a criterion, or a relative weight; shifting scope; changing the
apex or required margin; or applying an evaluation recommendation. An entry
**SHOULD** state whether a criterion change is deliberate recalibration or a drift
correction.

The skill **MUST NOT** log Markdown-body wording, typo, or formatting changes,
nor evaluated-source fixes that do not change the model.

> Rationale: those are not model changes; git already records them, and logging
> them turns a curated timeline into noise. — 0050

The skill **SHOULD** write **one entry per coherent change** — a confirmed
recommendation apply or model-authoring change — rather than one entry per field
touched.

> Rationale: the unit of record is the decision, not the edit; per-field entries
> fragment one rationale across many files. — 0050

### Who writes and reconciles

`setup` **MUST NOT** write the quality changelog. Setup's initial model rationale
belongs in the `QUALITY.md` body itself. A confirmed recommendation-apply or
direct model-authoring change **MUST** append an entry for each meaningful model
change after setup, cross-linking the evaluation run and recommendation when the
change came from one; writing the entry **MUST NOT** require confirmation beyond
the user's existing confirmation of the model change itself, since the entry's
rationale is the rationale already shown in the intent checkpoint or decision
brief.

`evaluate` **MUST NOT** write to the quality changelog. Issue-tracker handoff **MUST
NOT** write to the quality changelog.

> Rationale: evaluations own `.quality/evaluations/`; keeping the changelog to
> model changes only — referencing runs, never duplicating them — is what stops it
> becoming a second evaluation record. — 0050

Read-only orientation **MUST** remain read-only with respect to the changelog. It
**SHOULD** surface model history (the latest entry), and when the model has
changed out of band since the last logged entry it **SHOULD** classify that
under *needs reconciliation* readiness and offer a backfill route. The backfill
itself is performed by confirmed model-authoring or recommendation follow-up
work, not by orientation.

> Rationale: the changelog is curated, not complete; read-only orientation is where
> the gap left by hand edits gets caught and routed for repair without the
> orientation step itself mutating anything. — 0050

The run frame's mutation enumeration (see [Run frames](quality-skill.md#run-frames)) **MUST**
include the quality changelog as a distinct mutation surface, so a write to it is
visible before it happens.
