---
type: Functional Specification
title: Changelog Directory - functional spec
description: Requirements for renaming the model-change quality changelog to .quality/changelog/ and keeping .quality/logs/ as a flat workflow-log directory.
tags: [skill, workspace, logging, changelog]
timestamp: 2026-06-27T00:00:00Z
---

# Changelog Directory - functional spec

Companion to the [Changelog Directory](../0146-changelog-directory.md) change
case. This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

The workspace currently uses `.quality/log/` and `.quality/logs/` for different
artifact classes. The distinction is accurate but fragile: a singular/plural
path difference asks users and agents to remember which "log" is the
model-change history and which holds process logs. The model-change history
should be named as a changelog, while `.quality/logs/` should remain the flat
home for workflow and process logs, with each filename carrying its specific log
kind.

## Scope

Covered: workspace defaults, runtime skill guidance, durable skill specs,
workflow specs, CLI workspace vocabulary, active docs/examples, dogfood data,
logs, release notes, and verification.

Deferred:

- a single `.quality/changelog.md` artifact;
- nested `.quality/logs/<kind>/` directories;
- backward-compatibility readers, aliases, dual writes, migration commands, or
  fallback path probing;
- a configurable `.quality/config.yaml` `changelogDir` key; and
- a public `qualitymd changelog` command.

## Requirements

### Model-change changelog

- The workspace's default model-change history directory **MUST** be
  `.quality/changelog/`.

  > Rationale: the directory contains curated change history for the model, not
  > arbitrary logs. Naming it `changelog` removes the fragile singular/plural
  > distinction with `.quality/logs/`.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-changelog.md` - rename
  > the component contract to quality changelog and update the default path.

- A quality changelog entry filename **MUST** begin with a sortable UTC
  timestamp in `YYYY-MM-DDTHHMMSSZ` form, followed by a short human-readable
  slug and `.md`.

  > Rationale: date-only names do not sort multiple same-day model changes
  > reliably. UTC timestamps keep the per-entry-file model without introducing a
  > shared sequence allocator.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-changelog.md` - update
  > changelog entry naming.

- The quality changelog **MUST** remain one Markdown file per coherent meaningful
  model change.

  > Durable spec: modify `specs/skills/quality-skill/quality-changelog.md` - preserve
  > one-entry-per-change semantics under the renamed directory.

- Active runtime, spec, doc, and example surfaces **MUST NOT** instruct current
  writers to create or write `.quality/log/`.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` and
  > workflow/component specs that name the old path.

### Workflow logs

- Workflow and process logs **MUST** be written under `.quality/logs/`.

  > Rationale: `.quality/logs/` remains the umbrella for process artifacts,
  > while `.quality/changelog/` carries model-change history.
  >
  > Durable spec: modify `specs/skills/quality-skill/workflow-feedback-log.md` -
  > clarify `.quality/logs/` as the flat workflow-log directory.

- `.quality/logs/` **MUST** remain a flat directory. Workflow-log type **MUST**
  be expressed in the filename, not by a type-specific subdirectory.

  > Rationale: feedback logs are only one current log kind, and a flat directory
  > keeps run logs easy to scan chronologically across kinds.
  >
  > Durable spec: modify `specs/skills/quality-skill/workflow-feedback-log.md` -
  > prohibit type-specific subdirectories for feedback logs and future workflow
  > logs.

- Workflow feedback log filenames **MUST** keep the
  `<timestamp>-<workflow>-feedback-log.md` form.

  > Rationale: this existing form already carries timestamp, workflow, log kind,
  > and log identity. It fits the broader flat workflow-log rule without
  > renaming existing feedback logs.
  >
  > Durable spec: modify `specs/skills/quality-skill/workflow-feedback-log.md` -
  > frame feedback logs as one workflow-log kind.

### Clean break

- The implementation **MUST NOT** add backward-compatibility readers, aliases,
  dual writers, migration commands, or fallback path probing for `.quality/log/`.

  > Rationale: QUALITY.md is early alpha, and a clean break keeps the current
  > workspace model simpler than preserving two current names for one artifact
  > class.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-changelog.md` - record the
  > clean-break posture in deferred/non-goal language.

- Current dogfood model-change entries **MUST** move from `.quality/log/` to
  `.quality/changelog/` and use the new timestamped filename convention.

  > Durable spec: none - dogfood data demonstrates the current workspace layout
  > but is not itself a durable spec.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/quality-changelog.md` - rename the
  convention-first quality changelog contract under `.quality/changelog/`, with
  timestamped entry filenames and clean-break posture.
- `specs/skills/quality-skill/workflow-feedback-log.md` - clarify
  `.quality/logs/` as a flat workflow-log directory and feedback logs as one log
  kind.
- `specs/skills/quality-skill/quality-skill.md` - align cross-workflow mutation
  and deferred-surface wording with quality changelog and workflow logs.
- `specs/skills/quality-skill/workflows/setup.md` and
  `specs/skills/quality-skill/workflows/setup/feedback-log.md` - update setup's
  non-mutation boundary and path rationale.
- `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` and
  `specs/skills/quality-skill/recommendation-follow-up.md` - align
  recommendation-apply quality changelog wording.
- `specs/cli/evaluation-create.md` - update workspace vocabulary from quality
  log directory to quality changelog directory.

### To rename

- `specs/skills/quality-skill/quality-log.md` →
  `specs/skills/quality-skill/quality-changelog.md`.
- `skills/quality/guides/authoring/quality-log.md` →
  `skills/quality/guides/authoring/quality-changelog.md`.
- `specs/skills/quality-skill/guides/authoring/quality-log.md` →
  `specs/skills/quality-skill/guides/authoring/quality-changelog.md`.

### To delete

None

## Verification

- Source inspection **MUST** show `DefaultQualityLogDir` resolves to
  `.quality/changelog`.
- Source inspection **MUST** show active runtime and durable spec guidance use
  `.quality/changelog/` for model-change entries and `.quality/logs/` for flat
  workflow logs.
- Source inspection **MUST** show dogfood entries under `.quality/changelog/`
  with timestamp-prefixed filenames.
- A repository search for `.quality/log/` **MUST** return only historical
  archived Change Cases, historical logs/changelog entries, or other intentional
  history.
- Go tests and Markdown formatting checks **SHOULD** pass.
