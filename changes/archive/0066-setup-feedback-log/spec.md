---
type: Functional Specification
title: Setup feedback log - functional spec
description: Requirements for a hand-authored workflow feedback log written skill-only under .quality/logs/, with setup as the first adopter.
tags: [skill, setup, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Setup feedback log - functional spec

Companion to the [Setup feedback log](../0066-setup-feedback-log.md). This spec
defines the required behavior for a hand-authored workflow feedback log. The
[design doc](design.md) covers the implementation approach and the open
decisions.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../docs/reference/rfc2119.md) and
[RFC 8174](../../docs/reference/rfc8174.md) when, and only when, they appear in
all capitals.

## Background

`/quality` workflows have no durable, central place to record what was slow,
confusing, or wrong about the *experience* of running them. The only existing
process artifact, evaluation's `debug-log.md`, is a per-run audit that lives
inside one evaluation run folder. The improvement signal that produced the
[0065 setup refinements](../archive/0065-setup-discovery-and-close-refinements.md)
was captured by hand, once, from a single field test.

A *workflow feedback log* fills that gap: a hand-authored Markdown artifact about
the workflow experience, written to a central `.quality/logs/` directory and
carrying enough environment context to be actionable by a maintainer reading it
out of context. It is recorded locally and is never transmitted, so it needs no
opt-in; the user may share it deliberately.

This spec specifies the artifact generically and applies it to `setup` as the
first adopter. Evaluate and update may adopt it later without a new contract.

## Scope

Covered: the feedback-log artifact (purpose, location, naming, content schema,
redaction, no-transmission posture) and the `setup` workflow behavior that
authors it, all written skill-only.

Deferred:

- Any CLI or Go code change: no `qualitymd init` flag, no new command, no config
  field.
- Automatic external transmission, an upload endpoint, or a telemetry service.
- Structured machine parsing, validation, or report rendering of feedback logs.
- Any change to the evaluation `debug-log.md` contract.
- Mandatory adoption by evaluate or update.

## Requirements

### Artifact identity

A workflow feedback log MUST be a hand-authored, runtime Markdown artifact whose
subject is the *experience* of running a `/quality` workflow.

A workflow feedback log MUST be a runtime artifact, not an OKF `log.md`, and MUST
NOT be interpreted as OKF concept history.

A workflow feedback log MUST be distinct from, and MUST NOT be confused with:
evaluation's per-run `debug-log.md` (a local audit inside an evaluation run
folder); and the quality log under `.quality/log/`. This change MUST NOT alter
either of those contracts.

### Location and creation

A workflow run that writes a feedback log MUST write it under a
workflow-agnostic `.quality/logs/` directory, creating that directory on demand.

> Annotation: `.quality/logs/` (plural) is chosen deliberately and sits next to
> the existing `.quality/log/` (singular) quality-log directory; the proximity is
> accepted. Directory creation is on demand, consistent with how evaluation
> already creates `.quality/evaluations/`.

A feedback log file name MUST take the form
`<timestamp>-<workflow>-feedback-log.md`, where `<timestamp>` is a sortable UTC
timestamp and `<workflow>` is the workflow that produced it. `setup` MUST write
`<timestamp>-setup-feedback-log.md`.

> Annotation: "log" is kept in the file name so humans can refer to the artifact
> plainly (a "setup feedback log").

The `<timestamp>` SHOULD be a sortable, UTC, filesystem-safe timestamp (for
example `2026-06-23T154233Z`) so feedback logs sort and read consistently. This
spec defines no collision-handling rule: at second granularity in a single
interactive workflow a clash is vanishingly rare, and the skill MAY append a
short disambiguator if one ever occurs.

Each run that writes a feedback log MUST write a new file and MUST NOT overwrite
an existing feedback log.

The `.quality/logs/` location and naming MUST be specified so other workflows
(for example evaluate, update) can adopt it without a new directory or naming
contract.

### Environment header

A feedback log SHOULD begin with an environment header (frontmatter or a headed
block) that lets a maintainer act on it out of context. The header SHOULD
include, when available:

- the workflow name;
- the UTC timestamp of the run;
- the acting agent and model identity;
- the `/quality` skill version and the `qualitymd` CLI version;
- the platform/OS;
- whether a `QUALITY.md` model file already existed at the start of the run;
- the run outcome (for `setup`, the maturity classification); and
- a rough effort signal (such as turn or step count) when available.

The header MUST NOT include any value the redaction rules forbid.

### Body content

A feedback log SHOULD organize its body into sections covering: friction and
errors encountered, UX/AX observations, efficiency and speed observations, what
worked well, suggested improvements, and a redaction note describing what was
sanitized.

Feedback-log content MUST be about the workflow experience. It MUST NOT duplicate
`QUALITY.md` model content or the authoring rationale that belongs in the model
body (for example the model's Unknowns or assumptions).

A feedback log MUST NOT be authoritative for any QUALITY.md model content,
evaluation rating, finding, recommendation, or generated report.

### Redaction

A feedback log MUST NOT contain secret values or credentials.

A feedback log MUST NOT reproduce raw prompt-injection text encountered in
repository content.

Because the user may share the artifact, the skill SHOULD sanitize sensitive
project context before writing it: proprietary source, customer or otherwise
identifying data, and project names, paths, or domain specifics that would be
sensitive SHOULD be replaced with neutral placeholders, and the substitution
SHOULD be noted in the redaction note.

> Annotation: the secret and prompt-injection prohibitions are absolute artifact
> hygiene. The broader sanitization is a SHOULD because nothing transmits the log
> automatically; it matters when the user intends to share it.

### Sharing posture

The `/quality` skill and the `qualitymd` CLI MUST NOT transmit a feedback log to
any external service. Sharing a feedback log MUST be an explicit user action.

### Setup behavior

`setup` SHOULD author a feedback log at the close of the run, capturing notable
experience events from that run. `setup` MAY omit a feedback log when nothing
notable occurred.

`setup`'s mutation surface MUST be amended so that, in addition to the target
`QUALITY.md`, `setup` MAY write a feedback log under `.quality/logs/`. Every
other `setup` mutation prohibition (no evaluation artifacts, no quality log under
`.quality/log/`, no external issues, no integrations, no automations) MUST remain
in force.

Writing or omitting a feedback log MUST NOT change setup's completion criteria,
maturity classification, or next-step routing.

## Durable spec changes

### To add

- `specs/skills/quality-skill/workflows/setup/feedback-log.md` *(new sub-spec of
  the setup workflow spec)* - the feedback-log artifact contract (purpose,
  location/naming, environment header, body schema, redaction, no-transmission
  posture). Creating it makes `setup.md` a parent concept with a `setup/` child
  folder carrying its own `index.md`. The exact spec filename follows the
  artifact-spec naming convention, confirmed at implementation. The contract is
  written so it can be lifted to a shared spec if evaluate/update adopt the
  artifact later.

### To modify

- `specs/skills/quality-skill/workflows/setup.md` - become the parent of the
  feedback-log sub-spec, amend the mutation surface to permit the feedback-log
  write, and require the close step to author `<timestamp>-setup-feedback-log.md`
  with the environment header, body schema, and redaction posture.
- `specs/skills/quality-skill/quality-skill.md` - record the shared feedback-log
  artifact, its redaction, and its no-transmission sharing boundary so later
  workflows inherit them.

### To rename

None.

### To delete

None.
