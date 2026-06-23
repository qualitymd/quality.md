---
type: Functional Specification
title: Always-on setup feedback log - functional spec
description: Requirements for creating, updating, and finalizing a setup feedback log throughout every /quality setup run.
tags: [skill, setup, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Always-on setup feedback log - functional spec

Companion to the
[Always-on setup feedback log](../0068-always-on-setup-feedback-log.md). This
spec defines the delta from the existing optional setup feedback log to an
always-created, progressively updated setup run artifact.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background

The setup feedback log introduced by
[0066](../archive/0066-setup-feedback-log.md) is optional and close-only:
`setup` may omit it when nothing notable happened. That makes a missing log
ambiguous and loses feedback when the run is interrupted before the close step.
The artifact should instead exist for every setup run, with a concise "nothing
notable" record for clean runs and partial progress preserved for incomplete
runs.

## Requirements

### Creation timing

`setup` MUST create a setup feedback log during preflight after it resolves the
model file and obtains available CLI, skill, agent, and platform metadata.

> Rationale: creating the artifact early removes the ambiguity of a missing log
> and preserves useful run context even when setup stops before normal close. —
> 0068

The log path MUST remain `.quality/logs/<timestamp>-setup-feedback-log.md`,
where `<timestamp>` is the run-start UTC timestamp. The skill MUST create
`.quality/logs/` on demand and MUST NOT overwrite a feedback log from another
run.

The initial log MUST be valid Markdown with YAML frontmatter and the required
body sections present, even when most fields are still blank or marked
`in-progress`.

### Frontmatter metadata

The setup feedback log frontmatter MUST include these fields:

- `workflow`: `setup`.
- `status`: one of `in-progress`, `completed`, `interrupted`, or `failed`.
- `started-at`: UTC timestamp used for the file name.
- `updated-at`: UTC timestamp of the last material log update.
- `completed-at`: UTC timestamp when the run reaches a terminal state, otherwise
  blank.
- `agent`: acting agent identity when available.
- `model`: acting model identity when available.
- `skill-version`: `/quality` skill version from skill metadata when available.
- `cli-version`: `qualitymd version --json` result or concise unavailable/error
  summary.
- `platform`: OS/platform information when available.
- `model-file`: repository-relative model path when safe to record, otherwise a
  sanitized placeholder.
- `model-file-pre-existed`: boolean when known.
- `outcome`: setup maturity result (`starter`, `immature`, or
  `evaluation-ready`) at close, otherwise blank.
- `effort`: rough turn, step, or elapsed-time signal when available.
- `redaction`: `none`, `sanitized`, or `withheld-details`.

Frontmatter MUST NOT contain secrets, raw prompt-injection text, or sensitive
project-identifying context forbidden by the feedback-log redaction rules.

### Body format

The feedback log body MUST use these top-level sections:

```markdown
# Setup feedback log

## Timeline

## Friction and errors

## UX/AX observations

## Efficiency and speed

## What worked well

## Suggested improvements

## Redaction note
```

The `Timeline` section MUST record concise timestamped workflow-experience
events, including creation, major phase transitions, notable retries or stops,
and finalization.

The other body sections MUST remain about the experience of running setup. They
MUST NOT duplicate `QUALITY.md` model content, authoring rationale, evaluation
findings, recommendations, or quality-log history.

When a section has no notable content at close, the final log SHOULD say so
explicitly, for example `None observed.`

### Progressive updates

`setup` MUST update the current run's feedback log as the workflow progresses
when there is material workflow-experience information to record, including
friction, errors, confusing interaction points, retries, slow steps, redaction
decisions, or unusually smooth affordances worth preserving.

Each material update MUST refresh `updated-at`. Updates MAY revise the current
run's file in place; the "do not overwrite" rule only prohibits overwriting a
different run's log.

The log SHOULD avoid noisy churn for routine internal steps. A clean run should
produce a terse log, not a transcript.

### Finalization and incomplete runs

At normal close, `setup` MUST set `status: completed`, set `completed-at`, record
the setup maturity in `outcome`, update effort when available, and ensure each
body section has either useful content or an explicit no-notable-content note.

When setup stops because lint fails, CLI support is missing, user confirmation is
not granted, or another non-success stop occurs, the skill SHOULD finalize the
log with `status: failed` or `status: interrupted` when it can do so without
masking the stop condition. If finalization is impossible, the existing
`status: in-progress` log remains acceptable partial feedback.

Writing, updating, or finalizing the feedback log MUST NOT change setup's
completion criteria, maturity classification, or next-step routing.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/workflows/setup/feedback-log.md` - replace the
  optional close-only setup behavior with the always-on creation, metadata,
  body-format, progressive-update, and finalization requirements above.
- `specs/skills/quality-skill/workflows/setup.md` - update setup's preflight,
  workflow outline, mutation/artifact description, and close behavior to match
  the always-on feedback log.
- `specs/skills/quality-skill/quality-skill.md` - update shared workflow
  feedback-log wording so the parent contract permits current-run updates and no
  longer implies optional close-only setup logging.

### To rename

None

### To delete

None
