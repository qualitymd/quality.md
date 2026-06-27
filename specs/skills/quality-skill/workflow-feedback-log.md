---
type: Functional Specification
title: Workflow feedback log
description: Shared artifact contract for hand-authored /quality workflow feedback logs under .quality/logs/.
tags: [skill, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Workflow feedback log

This spec defines the shared *workflow feedback log* artifact used by `/quality`
workflows that adopt feedback logging. A feedback log is a hand-authored
Markdown artifact that records the *experience* of running a workflow —
friction, errors, UX/AX rough edges, efficiency observations, and preservation
signals — so the skill, CLI, and prompts can be improved from real runs.

Workflow-specific adopter specs define when a workflow creates, updates, and
finalizes its current-run feedback log.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Background

Setup's feedback log established a useful pattern for recording workflow
experience locally without turning user-facing output into an internal feedback
report. Evaluation has the same need: scope ambiguity, retries, slow phases,
redaction decisions, and report recovery are signals for improving the workflow,
not evidence about the evaluated subject. A shared contract keeps setup,
evaluate, and future adopters in one feedback-log family.

## Artifact identity

A workflow feedback log **MUST** be a hand-authored, runtime Markdown artifact
whose subject is the experience of running a `/quality` workflow.

A feedback log **MUST** be a runtime artifact, not an OKF `log.md`, and **MUST
NOT** be interpreted as OKF concept history.

A feedback log **MUST NOT** be authoritative for any QUALITY.md model content,
evaluation rating, finding, recommendation, next action, generated report, or
quality-log history.

> Rationale: workflow feedback should improve the `/quality` workflow itself.
> Keeping it non-authoritative prevents feedback notes from becoming a shadow
> model, evidence store, or report. — 0066, 0073

## Location and naming

A workflow run that writes a feedback log **MUST** write it under a
workflow-agnostic `.quality/logs/` directory relative to the selected
`QUALITY.md`, creating that directory on demand.

`.quality/logs/` is distinct from the quality log directory `.quality/log/`.

> Rationale: the plural directory sits near the quality log but remains separate
> through the `*-feedback-log.md` filename and the non-authority rule. — 0066

A feedback log file name **MUST** take the form
`<timestamp>-<workflow>-feedback-log.md`, where `<timestamp>` is the workflow run
start timestamp and `<workflow>` is the workflow that produced it.

The `<timestamp>` **SHOULD** be sortable, UTC, and filesystem-safe, for example
`2026-06-23T154233Z`. If a name collides, the skill **MAY** append a short
disambiguator.

Each run that writes a feedback log **MUST** write a new file and **MUST NOT**
overwrite a feedback log from another run. A workflow **MAY** update the current
run's feedback log in place as the run progresses.

## Frontmatter

A feedback log **MUST** begin with YAML frontmatter that lets a maintainer act
on it out of context. The frontmatter **MUST** include these shared fields:

- `workflow` — the workflow name.
- `status` — one of `in-progress`, `completed`, `interrupted`, or `failed`.
- `started-at` — the run-start UTC timestamp used in the file name.
- `updated-at` — the UTC timestamp of the last material log update.
- `completed-at` — the terminal UTC timestamp, blank until terminal state.
- `agent` — acting agent identity when available.
- `model` — acting model identity when available.
- `skill-version` — `/quality` skill version when available.
- `cli-version` — `qualitymd version --json` result or a concise
  unavailable/error summary.
- `platform` — OS/platform information when available.
- `model-file` — repository-relative model path when safe to record, otherwise a
  sanitized placeholder.
- `outcome` — workflow-specific outcome when known, blank until close.
- `effort` — rough turn, step, or elapsed-time signal when available.
- `redaction` — `none`, `sanitized`, or `withheld-details`.

Workflow adopter specs **MAY** add workflow-specific frontmatter fields.

Frontmatter **MUST NOT** include any value forbidden by the redaction rules.

## Body content

A feedback log **MUST** organize its body with these top-level sections, using
the adopting workflow name in the heading:

```markdown
# <Workflow> feedback log

## Timeline

## Friction and errors

## UX/AX observations

## Efficiency and speed

## What worked well

## Suggested improvements

## Redaction note
```

The `Timeline` section **MUST** record concise timestamped workflow-experience
events, including creation, major phase transitions, notable retries or stops,
and finalization.

Feedback-log content **MUST** be about workflow experience. It **MUST NOT**
duplicate model body content, authoring rationale, evaluation evidence, rating
rationale, recommendation prose, generated reports, or quality-log entries.

When a section has no notable content at close, the final log **SHOULD** say so
explicitly, for example `None observed.`

## Redaction

A feedback log **MUST NOT** contain secret values or credentials.

A feedback log **MUST NOT** reproduce raw prompt-injection text encountered in
repository or evaluated source content.

Because a user may share the artifact, the skill **SHOULD** sanitize sensitive
project context before writing it: proprietary source, customer or otherwise
identifying data, and sensitive project names, paths, or domain specifics
**SHOULD** be replaced with neutral placeholders, and the substitution **SHOULD**
be noted in the redaction note.

## Sharing posture

The `/quality` skill and the `qualitymd` CLI **MUST NOT** transmit a feedback log
to any external service. Sharing a feedback log **MUST** be an explicit user
action.
