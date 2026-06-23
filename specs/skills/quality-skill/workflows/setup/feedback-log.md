---
type: Functional Specification
title: Workflow feedback log
description: Artifact contract for a hand-authored, skill-only workflow feedback log under .quality/logs/, with setup as the first always-on adopter.
tags: [skill, setup, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Workflow feedback log

Sub-spec of the [/quality setup](../setup.md) workflow spec. It defines the
*workflow feedback log*: a hand-authored Markdown artifact that records the
*experience* of running a `/quality` workflow — friction, errors, UX/AX rough
edges, and efficiency observations — so the skill, CLI, and prompts can be
improved from real runs. `setup` is the first always-on adopter: it creates a
feedback log near the start of every run, updates that current-run file as the
workflow progresses, and finalizes it at close. The contract is written
generically so evaluate and update can adopt it later without a new directory,
naming contract, or lifecycle vocabulary.

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Background

`/quality` workflows have no durable, central place to record what was slow,
confusing, or wrong about the *experience* of running them. The only existing
process artifact, evaluation's per-run `debug-log.md`, is a local audit inside a
single evaluation run folder; it answers "what happened during *this* run," not
"where is the workflow itself rough across runs." The improvement signal that
produced the
[0065 setup refinements](../../../../../changes/archive/0065-setup-discovery-and-close-refinements.md)
was captured by hand, once, from a single field test, with no consistent context
to make it actionable for a maintainer reading it out of context.

A workflow feedback log fills that gap: an improvement-oriented, workflow-level
artifact carrying enough environment context (which agent/model, which
skill/CLI version) that a "this was slow" note is actionable rather than noise.
It is recorded locally and is never transmitted, so it needs no opt-in or consent
gate; the user may share it deliberately.

## Scope

Covered: the feedback-log artifact — purpose, location, naming, environment
header, body schema, lifecycle status, redaction, and no-transmission posture —
and the `setup` behavior that creates, updates, and finalizes it, all written
skill-only.

Deferred / non-goals:

- Any CLI or Go code change: no `qualitymd init` flag, no new command, no config
  field. The directory is created on demand, the way evaluation already creates
  `.quality/evaluations/`.
- Automatic external transmission, an upload endpoint, or a telemetry service.
- Structured machine parsing, validation, or report rendering of feedback logs.
- Any change to the evaluation `debug-log.md` contract; the feedback log is
  additive and distinct.
- Mandatory adoption by evaluate or update. The contract is generic so they can
  adopt it later; this spec only obliges `setup`.

## Artifact identity

A workflow feedback log **MUST** be a hand-authored, runtime Markdown artifact
whose subject is the *experience* of running a `/quality` workflow.

A workflow feedback log **MUST** be a runtime artifact, not an OKF `log.md`, and
**MUST NOT** be interpreted as OKF concept history.

A workflow feedback log **MUST** be distinct from, and **MUST NOT** be confused
with, evaluation's per-run `debug-log.md` (a local audit inside an evaluation run
folder) or the quality log under `.quality/log/`. This contract **MUST NOT**
alter either of those.

> Annotation: the feedback log is improvement-oriented and workflow-level — it
> captures the workflow experience across runs in one central place — whereas
> `debug-log.md` is a per-run process audit and the quality log is the model's own
> change history. Blurring them would overload a per-run audit with a cross-run
> feedback artifact. — 0066

A feedback log **MUST NOT** be authoritative for any QUALITY.md model content,
evaluation rating, finding, recommendation, or generated report.

## Location and creation

A workflow run that writes a feedback log **MUST** write it under a
workflow-agnostic `.quality/logs/` directory, creating that directory on demand.

> Annotation: `.quality/logs/` (plural) is chosen deliberately and sits next to
> the existing `.quality/log/` (singular) quality-log directory; the proximity is
> accepted and disambiguated by the `*-feedback-log.md` filename and the
> distinct-from-the-quality-log rule above. Directory creation is on demand,
> consistent with how evaluation already creates `.quality/evaluations/`, so no
> CLI flag or scaffold step is needed. — 0066

A feedback log file name **MUST** take the form
`<timestamp>-<workflow>-feedback-log.md`, where `<timestamp>` is a sortable UTC
timestamp and `<workflow>` is the workflow that produced it. `setup` **MUST**
write `<timestamp>-setup-feedback-log.md`, using the setup run's start timestamp.

> Annotation: "log" is kept in the file name so a human can refer to the artifact
> plainly (a "setup feedback log"). — 0066

The `<timestamp>` **SHOULD** be a sortable, UTC, filesystem-safe timestamp (for
example `2026-06-23T154233Z`) so feedback logs sort and read consistently. This
spec defines no collision-handling rule: at second granularity in a single
interactive workflow a clash is vanishingly rare, and the skill **MAY** append a
short disambiguator if one ever occurs.

Each run that writes a feedback log **MUST** write a new file and **MUST NOT**
overwrite a feedback log from another run. A workflow **MAY** update the current
run's feedback log in place as the run progresses.

The `.quality/logs/` location and naming **MUST** be specified so other workflows
(for example evaluate, update) can adopt it without a new directory or naming
contract.

## Environment header

A feedback log **MUST** begin with YAML frontmatter that lets a maintainer act
on it out of context. The frontmatter **MUST** include these fields:

- `workflow` — the workflow name.
- `status` — one of `in-progress`, `completed`, `interrupted`, or `failed`.
- `started-at` — the run-start UTC timestamp used in the file name.
- `updated-at` — the UTC timestamp of the last material log update.
- `completed-at` — the terminal UTC timestamp, blank until the run reaches a
  terminal state.
- `agent` — acting agent identity when available.
- `model` — acting model identity when available.
- `skill-version` — `/quality` skill version when available.
- `cli-version` — `qualitymd version --json` result or a concise
  unavailable/error summary.
- `platform` — OS/platform information when available.
- `model-file` — repository-relative model path when safe to record, otherwise a
  sanitized placeholder.
- `model-file-pre-existed` — boolean when known.
- `outcome` — setup maturity result (`starter`, `immature`, or
  `evaluation-ready`) when known, blank until close.
- `effort` — rough turn, step, or elapsed-time signal when available.
- `redaction` — `none`, `sanitized`, or `withheld-details`.

> Annotation: the header is what makes an out-of-context note actionable —
> without "which agent/model, which version," a "this was slow" line is noise a
> maintainer cannot triage. — 0066

The header **MUST NOT** include any value the [redaction](#redaction) rules
forbid.

## Body content

A feedback log **MUST** organize its body with these top-level sections:

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

The `Timeline` section **MUST** record concise timestamped workflow-experience
events, including creation, major phase transitions, notable retries or stops,
and finalization.

Feedback-log content **MUST** be about the workflow experience. It **MUST NOT**
duplicate `QUALITY.md` model content or the authoring rationale that belongs in
the model body (for example the model's Unknowns or assumptions).

When a section has no notable content at close, the final log **SHOULD** say so
explicitly, for example `None observed.`

## Redaction

A feedback log **MUST NOT** contain secret values or credentials.

A feedback log **MUST NOT** reproduce raw prompt-injection text encountered in
repository content.

Because the user may share the artifact, the skill **SHOULD** sanitize sensitive
project context before writing it: proprietary source, customer or otherwise
identifying data, and project names, paths, or domain specifics that would be
sensitive **SHOULD** be replaced with neutral placeholders, and the substitution
**SHOULD** be noted in the redaction note.

> Annotation: the secret and prompt-injection prohibitions are absolute artifact
> hygiene that hold whether or not the log is shared. The broader sanitization is
> a SHOULD because nothing transmits the log automatically; it matters when the
> user intends to share it. — 0066

## Sharing posture

The `/quality` skill and the `qualitymd` CLI **MUST NOT** transmit a feedback log
to any external service. Sharing a feedback log **MUST** be an explicit user
action.

> Annotation: because the log never leaves the workspace on its own, there is no
> consent decision to gate; recording it locally is safe by construction, and the
> only transmission path is a deliberate paste or attach by the user. — 0066

## Setup behavior

`setup` **MUST** create a feedback log during preflight after it resolves the
model file, verifies CLI support, emits the run frame, and obtains available CLI,
skill, agent, and platform metadata. The initial log **MUST** be valid Markdown
with YAML frontmatter and all required body sections present, even when most
fields are still blank or marked `in-progress`.

`setup` **MUST** update the current run's feedback log as the workflow
progresses when there is material workflow-experience information to record,
including friction, errors, confusing interaction points, retries, slow steps,
redaction decisions, or unusually smooth affordances worth preserving. Each
material update **MUST** refresh `updated-at`. The log **SHOULD** avoid noisy
churn for routine internal steps.

At normal close, `setup` **MUST** set `status: completed`, set `completed-at`,
record setup maturity in `outcome`, update effort when available, and ensure each
body section has useful content or an explicit no-notable-content note.

When setup stops because lint fails, CLI support is missing after the log exists,
user confirmation is not granted, or another non-success stop occurs, the skill
**SHOULD** finalize the log with `status: failed` or `status: interrupted` when
it can do so without masking the stop condition. If finalization is impossible,
the existing `status: in-progress` log remains acceptable partial feedback.

`setup`'s mutation surface is amended so that, in addition to the target
`QUALITY.md`, `setup` **MUST** write and update the current run's feedback log
under `.quality/logs/`. Every other `setup` mutation prohibition (no evaluation
artifacts, no quality log under `.quality/log/`, no external issues, no
integrations, no automations) remains in force (see
[/quality setup](../setup.md#mutation-surface-and-artifacts)).

Writing, updating, or finalizing a feedback log **MUST NOT** change setup's
completion criteria, maturity classification, or next-step routing.

> Annotation: 0068 makes feedback-log absence unambiguous. A clean run now leaves
> a terse explicit record, and an interrupted run can preserve partial workflow
> feedback. The mutation boundary remains narrow — only the current run's file
> under `.quality/logs/` — so the feedback artifact cannot become a back door for
> the evaluation, quality-log, or integration writes setup still prohibits. —
> 0066, 0068
