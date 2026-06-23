---
type: Change Case
title: Setup feedback log
description: Add a hand-authored workflow feedback log under .quality/logs/ that records setup-experience friction and improvement signals, written skill-only, with setup as the first adopter.
status: Done
tags: [skill, setup, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Setup feedback log

A **Change Case** adding a *workflow feedback log*: a hand-authored Markdown
artifact that records the *experience* of running a `/quality` workflow — errors,
friction, UX/AX rough edges, and efficiency/speed observations — so the skill,
CLI, and prompts can be improved from real runs. It is written locally; the user
may choose to share it. `setup` is the first adopter; the artifact and its
`.quality/logs/` home are specified generically so later workflows (evaluate,
update) can adopt them.

The artifact is written **skill-only**: no CLI flag and no new command. The
directory is created on demand the way evaluation runs already create
`.quality/evaluations/`.

Detail lives in:

- [Functional spec](0066-setup-feedback-log/spec.md) - what the change must do.
- [Design doc](0066-setup-feedback-log/design.md) - how it is shaped, and why.

## Motivation

The [0065 setup refinements](archive/0065-setup-discovery-and-close-refinements.md)
existed only because a human ran `/quality setup` against a real external repo
and hand-captured three recurring frictions. That feedback loop is valuable but
ad hoc: there is no durable place for a workflow run to record what was slow,
confusing, or wrong about the *experience* of running it, and no consistent
context (which agent/model, which skill/CLI version) to make such notes
actionable for a maintainer reading them out of context.

This is different from evaluation's `debug-log.md`, which is a *per-run* process
audit living inside one evaluation run folder. A feedback log is
*improvement-oriented* and *workflow-level*: it captures the workflow experience
across runs in one central place, and carries enough environment context to act
on. It is recorded locally and is not transmitted anywhere; sharing it with
maintainers is an explicit user action, so no opt-in/consent gate is needed.

Because the user may share it, two artifact-hygiene rules still apply: no secrets
and no raw prompt-injection text are written into it, and sensitive project
context should be sanitized when the log is meant to leave the workspace.

## Scope

Covered:

- Define a workflow feedback log: a hand-authored Markdown artifact about the
  workflow experience, distinct from evaluation's per-run `debug-log.md`, from
  the quality log under `.quality/log/`, and from any OKF `log.md`.
- Establish a `.quality/logs/` directory and a
  `<timestamp>-<workflow>-feedback-log.md` naming convention; `setup` writes
  `<timestamp>-setup-feedback-log.md`. The directory is created on demand. The
  directory and naming are shaped so later workflows can adopt them.
- Specify the feedback log's environment header and body section schema
  (friction/errors, UX/AX, efficiency/speed, what worked, suggested
  improvements, redaction note).
- Specify redaction (no secrets, no raw prompt-injection text, sanitize
  sensitive project context) and a no-transmission posture.
- Amend `setup`'s mutation surface so it MAY write a feedback log under
  `.quality/logs/`, while every other setup mutation prohibition stays intact.
- Update the runtime skill and durable specs to match.

Deferred / non-goals:

- No QUALITY.md format change.
- **No CLI or Go code change.** No `qualitymd init` flag, no new command, no
  config field. Directory creation is on demand, consistent with how evaluation
  already creates `.quality/evaluations/`.
- No opt-in/consent mechanism (the log is recorded locally and never
  transmitted).
- No change to which discovery questions setup asks, their defaults, or
  confidence vocabulary.
- No automatic external transmission, upload endpoint, or telemetry service.
- No structured machine parsing, validation, or report rendering of feedback
  logs.
- No change to the existing evaluation `debug-log.md` contract; the feedback log
  is additive and distinct.
- No requirement that evaluate/update adopt the feedback log in this case (the
  contract is written generically so they can later).

## Settled decisions

- **Durable spec home.** The feedback-log contract is a **sub-spec of the setup
  workflow spec**, under `specs/skills/quality-skill/workflows/setup/`, not a
  top-level `specs/workflow-feedback-log.md`. It can be lifted to a shared spec if
  evaluate/update adopt the artifact later.
- **Filename.** `<timestamp>-setup-feedback-log.md` — "log" is kept in the name so
  humans can refer to it plainly.
- **Timestamp.** A sortable, UTC, filesystem-safe timestamp (for example
  `2026-06-23T154233Z`), specified as a SHOULD for consistency. No normative
  collision rule: at second granularity in a single interactive workflow a clash
  is vanishingly rare, so the agent disambiguates ad hoc if it ever occurs.
- **`.quality/logs/` vs `.quality/log/` proximity.** The singular `.quality/log/`
  already exists as the quality-log dir; the plural `logs/` is chosen deliberately
  and the proximity is accepted.

## Affected artifacts

### Code

- [x] **None.** The artifact is written skill-only; `.quality/logs/` is created
      on demand. No `qualitymd init` flag, no new command, no config field.

### Durable specs

- [x] `specs/skills/quality-skill/workflows/setup/feedback-log.md` *(new
      sub-spec)* - the feedback-log artifact contract: purpose, location/naming,
      environment header, body schema, redaction, and no-transmission posture.
      Adding it made `setup.md` a parent concept with a `setup/` child folder
      (with its own `index.md`). Filename uses the capability-named
      artifact-contract convention (`feedback-log.md`).
- [x] `specs/skills/quality-skill/workflows/setup.md` - became the parent of the
      feedback-log sub-spec, amended the mutation surface to permit the
      feedback-log write, and added the close-step authoring of
      `<timestamp>-setup-feedback-log.md`.
- [x] `specs/skills/quality-skill/quality-skill.md` - recorded the shared
      feedback-log artifact and its redaction/no-transmission boundary so future
      workflows inherit it; added it to the run-frame mutation classes.
- [x] `specs/skills/quality-skill/workflows/index.md` - listed the new `setup/`
      sub-folder.
- [x] `specs/log.md` and `specs/skills/quality-skill/workflows/log.md` - recorded
      the durable spec updates.

### Durable docs

- [x] `docs/guides/use-quality-skill.md` - described the feedback log and how to
      share it.
- [x] `docs/log.md` - recorded the durable documentation update.

### Bundled skill

- [x] `skills/quality/SKILL.md` - added the feedback log to setup's hard rules
      (redaction, no transmission) and the run-frame mutation classes. Skill
      `metadata.version` bump is deferred to the release cut (it must match the
      release tag), consistent with the other Unreleased setup changes still at
      `0.10.0`.
- [x] `skills/quality/workflows/setup.md` - added the close-step feedback-log
      authoring of `<timestamp>-setup-feedback-log.md`, the environment header +
      section schema, redaction guidance, the run-frame mutation/artifact lines,
      and the workflow-diagram close step.
- [x] `skills/quality/resources/cli-quick-reference.md` - noted `.quality/logs/`
      alongside the other `.quality/` workspace artifacts.

### Release

- [x] `CHANGELOG.md` - added the user-facing `/quality Skill` entry under
      `Unreleased`.

## Children

- [Functional spec](0066-setup-feedback-log/spec.md) - what the feedback log,
  schema, redaction, and setup behavior must do.
- [Design doc](0066-setup-feedback-log/design.md) - how the artifact and its
  skill-only authoring are shaped, and the alternatives weighed.

## Status

`Done`. Implemented skill-only, with **no code change**. The new durable sub-spec
[`specs/skills/quality-skill/workflows/setup/feedback-log.md`](../specs/skills/quality-skill/workflows/setup/feedback-log.md)
defines the hand-authored, never-transmitted feedback artifact under
`.quality/logs/`; the
[`setup`](../specs/skills/quality-skill/workflows/setup.md) spec and the runtime
[`workflows/setup.md`](../skills/quality/workflows/setup.md) gained the close-step
authoring and the narrowed mutation surface; and the parent
[`/quality` skill](../specs/skills/quality-skill/quality-skill.md) spec records
the shared artifact and its redaction/no-transmission boundary. Durable docs, the
bundled skill, and `CHANGELOG.md` are in sync (see **Affected artifacts**). The
skill `metadata.version` bump is deferred to the release cut. Moved to
[`archive/`](archive/index.md).
