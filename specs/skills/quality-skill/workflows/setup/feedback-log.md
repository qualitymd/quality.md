---
type: Functional Specification
title: Setup feedback log
description: Setup-specific adoption rules for the shared workflow feedback log.
tags: [skill, setup, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Setup feedback log

Sub-spec of the [/quality setup](../setup.md) workflow spec. It defines how
`setup` adopts the shared
[workflow feedback log](../../workflow-feedback-log.md) artifact contract.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Background

Setup was the first always-on feedback-log adopter. Its behavior now remains
setup-specific, while shared artifact identity, location, frontmatter, body
shape, redaction, and sharing posture live in the shared feedback-log spec.

## Setup behavior

`setup` **MUST** create a feedback log after it presents the setup preview and
before it continues into discovery or authoring. Setup **MAY** stop before the
feedback log exists when CLI support is missing, preflight fails, or the user
stops after the initial read-only context scan.

`setup` **MUST** write the log to
`.quality/logs/<timestamp>-setup-feedback-log.md`, using the setup run's start
timestamp and following the shared feedback-log contract.

The setup feedback log frontmatter **MUST** include these setup-specific fields:

- `model-file-pre-existed` — boolean when known.
- `outcome` — setup workflow result when known, blank until close. Allowed
  terminal values are `completed`, `completed-with-important-gaps`,
  `lint-failed`, `failed`, and `interrupted`.

`setup` **MUST** update the current run's feedback log as the workflow
progresses when there is material workflow-experience information to record,
including friction, errors, confusing interaction points, retries, slow steps,
redaction decisions, or unusually smooth affordances worth preserving. Each
material update **MUST** refresh `updated-at`. The log **SHOULD** avoid noisy
churn for routine internal steps.

When the feedback log is created after setup preview, the initial timeline entry
**SHOULD** state that creation point. It **MAY** summarize material pre-log
workflow-experience events, such as slow context analysis, CLI friction, or
redaction decisions. It **MUST NOT** duplicate setup-preview model content or
serve as assessment evidence.

At normal close, `setup` **MUST** set `status: completed`, set `completed-at`,
record `outcome: completed` or `outcome: completed-with-important-gaps`, update
effort when available, and ensure each body section has useful content or an
explicit no-notable-content note.

When setup stops before the feedback log exists because CLI support is missing,
preflight fails, or the user stops after the initial read-only context scan, the
absence of a feedback log **MUST NOT** be treated as a failed setup artifact.

When setup stops because lint fails, CLI support is missing after the log exists,
user confirmation is not granted, or another non-success stop occurs, the skill
**SHOULD** finalize the log with `status: failed` or `status: interrupted` when
it can do so without masking the stop condition, and **SHOULD** set `outcome:
lint-failed`, `failed`, or `interrupted` as appropriate. If finalization is
impossible, the existing `status: in-progress` log remains acceptable partial
feedback.

`setup`'s mutation surface includes the target `QUALITY.md` and the current
run's feedback log under `.quality/logs/`. Every other setup mutation
prohibition remains in force: setup does not create evaluation artifacts, write
the quality log under `.quality/log/`, create external issues, configure
integrations, or configure automations.

Writing, updating, or finalizing a feedback log **MUST NOT** change setup's
completion criteria, important-gap judgment, or next-step routing.
