---
type: Design Doc
title: Setup intro preview - design
description: Design for opening /quality setup with educational orientation, showing a setup preview before discovery, and delaying feedback-log creation until after the preview.
tags: [skill, setup, ux]
timestamp: 2026-06-26T00:00:00Z
---

# Setup intro preview - design

Answers the [functional spec](spec.md). This is a bundled skill guidance and
durable-spec mirror change: no CLI, Go, format-schema, rating, roll-up,
evaluation-record, or report behavior changes.

## Context

Setup currently has all of the right phases, but the first meaningful
project-specific value appears after several process steps. The workflow should
still verify tooling, read authoring guidance, inspect context, ask pedagogical
questions, and write only after review. The improvement is presentation order:
orient immediately, scan read-only, preview what was learned, then ask the user
to calibrate it.

## Approach

### Add an opening card

Add a short `Opening orientation` stage before preflight detail. The card uses
README-derived copy:

```text
QUALITY.md gives AI assistants, coding agents, and teams a holistic definition
of quality tailored to their project, so they can stay aligned, identify
critical risks and issues, and keep improving.
```

It then states the immediate workflow promise: read-only context scan first,
review before changes. This is not a new artifact and not a marketing surface;
it is the first interaction block of setup.

### Insert setup preview after context analysis

Keep the existing setup brief as working context, then require a user-facing
preview distilled from it before discovery. The preview names the likely root,
domain, visible evidence, candidate model shape, missing context, and next
action.

Discovery stays intact. The preview makes discovery feel like calibration of a
draft, but it does not confirm answers or replace the required teaching beats.

### Delay feedback-log creation

Move feedback-log creation from preflight to after the setup preview for runs
that continue into discovery or authoring. The log can summarize notable pre-log
workflow-experience events, but it must not duplicate model content or become an
evidence store.

This keeps the feedback artifact while avoiding the feeling that setup's first
real work is bookkeeping. Runs that stop before the preview or immediately after
the read-only scan need not leave a feedback log behind.

### Update durable mirrors and logs

Mirror the runtime changes in the durable setup workflow spec and setup
feedback-log spec. Add append-only log entries in the live skill workflow log,
durable workflow log, and CHANGELOG.

## Alternatives

- **Only add the opening copy.** Rejected. It improves tone but does not solve
  the late project-specific value problem.
- **Keep feedback-log creation in preflight.** Rejected. The log is useful, but
  its current timing contributes to setup feeling process-heavy.
- **Let the preview replace discovery.** Rejected. The discovery questions still
  teach the model dimensions and prevent low-confidence inferences from becoming
  confirmed facts.
- **Move the full setup brief into user output.** Rejected. The setup brief is
  working context; the preview should be shorter and action-oriented.

## Trade-offs & risks

- A delayed feedback log may miss exact timing for early preflight failures.
  This is acceptable because setup has not yet produced user-visible model work,
  and the absence of a log is now explicitly allowed for early stops.
- The opening card could become too long. The runtime guidance keeps it to one
  value-proposition sentence plus immediate workflow status.
- The preview might sound too definitive. The spec requires confidence and
  missing-context language, and discovery still confirms or corrects the draft.

## Open questions

None.
