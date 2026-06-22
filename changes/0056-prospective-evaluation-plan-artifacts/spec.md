---
type: Functional Specification
title: Prospective evaluation plan artifacts - functional spec
description: Requirements for authoring design.md and the initial plan.md before assessment begins in a /quality evaluation run.
tags: [evaluation, skill, records]
timestamp: 2026-06-22T00:00:00Z
---

# Prospective evaluation plan artifacts - functional spec

Companion to the
[Prospective evaluation plan artifacts](../0056-prospective-evaluation-plan-artifacts.md)
change case. This spec states *what* the `/quality evaluate` workflow and
evaluation-record artifact contracts must require so `design.md` and the initial
`plan.md` are prospective planning artifacts rather than after-the-fact
summaries. It defers the run-folder layout to
[`specs/evaluation-records/run-folder.md`](../../../specs/evaluation-records/run-folder.md),
the plan artifact contract to
[`specs/evaluation-records/plan-md.md`](../../../specs/evaluation-records/plan-md.md),
and the skill evaluation workflow to
[`specs/skills/quality-skill/evaluation.md`](../../../specs/skills/quality-skill/evaluation.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

`qualitymd evaluation create` already creates `design.md` and `plan.md` before
assessment records exist. The problem is the next step: the skill guidance says
to "fill" those files with judgment content, including phrasing that can be read
retrospectively, such as an executed evidence basis. If the files are completed
after assessments or recommendations, they no longer constrain the evaluation;
they merely narrate it.

The evaluation run needs both: prospective artifacts that make the method and
coverage auditable before judgment starts, and formal records/reports that carry
the actual findings, evidence, ratings, and recommendations afterward. This
case tightens the boundary so `design.md` and the initial `plan.md` are written
at the start of the run, while any changes discovered mid-run are explicit
amendments.

## Scope

Covered:

- Ordering in `/quality evaluate`: run creation, then design and initial plan
  authoring, then assessment.
- Required prospective content for `design.md` and `plan.md`.
- Handling for optional `coverage:` frontmatter and later plan amendments.
- Durable spec and runtime skill wording needed to preserve that ordering.

Non-goals:

- No change to how the CLI numbers run folders or seeds stub files.
- No requirement for the CLI to populate judgment content in `design.md` or
  `plan.md`.
- No change to assessment, analysis, recommendation, or report schemas.

## Requirements

### Pre-assessment artifact authoring

For `/quality evaluate`, after `qualitymd evaluation create` succeeds and before
the skill begins assessment evidence collection, writes assessment records, writes
analysis records, writes recommendations, checks reportability, or builds a
report, the skill **MUST** author the run's initial `design.md` and `plan.md`.

> Rationale: the files are useful as planning artifacts only if they exist before
> the evaluation judgment they are meant to guide and audit.

The initial `design.md` **MUST** capture the resolved evaluation frame:

- resolved model file and model snapshot relationship;
- scope or narrowing;
- in-scope areas;
- out-of-scope or deferred areas;
- methodological constraints or rating limitations known before assessment.

The initial `plan.md` **MUST** capture the planned execution:

- chosen rigor;
- concrete in-scope requirement set;
- intended evidence basis or inspection strategy;
- planned commands, searches, or source reads when known;
- planned limitations.

The initial `plan.md` **MUST NOT** present actual findings, rating rationale, or
recommendation reasoning as if they were planning content.

> Rationale: actual evidence and judgment belong in assessment, analysis,
> recommendation, and report artifacts. Putting them in the plan makes the plan a
> duplicate report and hides whether coverage was planned before judgment.

### Prospective and retrospective evidence boundaries

The skill **MUST** distinguish intended evidence from executed evidence:

- initial `plan.md` records intended evidence basis or inspection strategy;
- assessment and analysis records cite actual evidence and rating rationale;
- `report.md`, `report-summary.md`, and `report.json` render the final result;
- `debug-log.md` records notable process events, not formal judgment.

`design.md` and `plan.md` **MAY** summarize actual coverage only in clearly
marked amendment or update sections, and those sections **MUST NOT** replace or
erase the original prospective content.

### Plan updates and coverage metadata

When the evaluator changes scope, coverage, rigor, or material evidence strategy
after the initial plan is written, the skill **SHOULD** amend `plan.md` under a
clear heading such as `Plan updates` rather than silently rewriting the plan as
though the new path had always been intended.

When resume diagnostics materially matter, `coverage:` frontmatter **SHOULD** be
added after the intended assessment requirements and analysis areas are settled
and before dependent record-writing begins. If planned coverage changes during
the run, the `coverage:` frontmatter and amendment note **SHOULD** be updated
together.

### CLI scaffolding boundary

`qualitymd evaluation create` continues to seed `design.md` and body-only
`plan.md` as mechanical scaffolding. The CLI **MUST NOT** be required to infer or
generate the design, plan, planned coverage, or evidence strategy.

> Rationale: the CLI can know the run folder and model snapshot; it cannot know
> the evaluator's scope judgment, rigor trade-offs, or evidence strategy. That
> remains skill-owned judgment.

## Durable spec changes

### To add

- `specs/evaluation-records/design-md.md` — add an artifact spec for `design.md`
  if the final change keeps design content independently reviewable; it should
  own the resolved evaluation frame and its pre-assessment timing (per the
  pre-assessment artifact authoring requirement).

### To modify

- `specs/skills/quality-skill/evaluation.md` — make the workflow order explicit:
  create run, author initial `design.md` and `plan.md`, optionally add settled
  coverage metadata, then assess and write records (per all requirements).
- `specs/evaluation-records/plan-md.md` — clarify that `plan.md` is prospective,
  that actual findings/rationale belong in records and reports, and that later
  changes are explicit amendments (per the prospective/retrospective boundary
  and plan update requirements).
- `specs/evaluation-records/run-folder.md`,
  `specs/evaluation-records/index.md`, and `specs/evaluation-records.md` — add
  or update links if a `design.md` artifact spec is introduced (per the
  pre-assessment artifact authoring requirement).

### To rename

None.

### To delete

None.
