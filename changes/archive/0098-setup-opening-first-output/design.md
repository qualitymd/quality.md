---
type: Design Doc
title: Setup opening as first output — design doc
description: How /quality setup is reordered so the opening, roadmap, and run frame precede any tool call.
tags: [skill, quality, setup, ux]
timestamp: 2026-06-25T00:00:00Z
---

# Setup opening as first output — design doc

Design behind the [Setup opening as first output](../0098-setup-opening-first-output.md)
change case and its [functional spec](spec.md).

## Context

A field run of `/quality setup` produced 1–2 minutes of silent tool work before
any text reached the user. The opening orientation, run frame, and setup preview
all flushed together at the end. The cause is ordering, not a missing feature:
the agent front-loaded cheap tool calls (CLI version checks, repository scans, a
timestamp probe) and even wrote the feedback log before emitting its first line.

The governing spec permits this. Two seams matter:

1. The opening-orientation requirement reads "before setup performs long-running
   context work" — silent on whether *any* tool call may precede the opening. An
   agent reasonably batches quick checks first.
2. The Workflow structure step list couples the run frame to prerequisite
   verification: "Resolve the target `QUALITY.md`, verify setup prerequisites,
   and emit the run frame." That phrasing implies the frame waits on the CLI
   gate.

Neither the opening nor the run frame has a real tool dependency. The frame's
only variable field is the resolved model path, which is known from the
invocation arguments (explicit path, else literal `QUALITY.md`). The frame does
not display CLI status.

## Approach

Reorder setup's front matter into three beats:

1. **First output (no tools).** Warm welcome + value sentence + a short phase
   roadmap (scan → calibrate → review → write → verify) + the read-only /
   review-before-changes boundary + the run frame. A one-line cue notes the scan
   can take a moment on a large repo. This is pure text and is emitted before any
   tool call.
2. **CLI gate (fail-fast).** Run the CLI prerequisite check. If the CLI is
   missing or unsupported, stop here with a clear message — after the user has
   been oriented, before any scanning effort is spent.
3. **Read-only scan and onward.** Unchanged: inspect context, build the setup
   brief, present the setup preview, create the feedback log, discovery, review,
   author, lint, close.

Implementation touches two files:

- **Durable spec** `specs/skills/quality-skill/workflows/setup.md`: split the
  current Workflow structure step 2 so the run frame is emitted with the opening
  (step 1) and prerequisite verification is its own step before the scan; and
  strengthen the Context-analysis opening requirement to "before any tool call"
  plus the roadmap and scan-cue requirements. Carry the *why* into an annotation
  citing this case.
- **Runtime skill** `skills/quality/workflows/setup.md`: reorder the procedure
  text block; rewrite the Opening orientation section to lead with the welcome +
  roadmap + run frame as first output; move run-frame emission out of Preflight's
  body so Preflight is just the CLI gate + model resolution that the frame
  already named.

## Alternatives

- **Leave the spec, fix only the runtime.** Rejected: the runtime already listed
  orientation first and it still failed, because nothing forbade pre-orientation
  tool calls. The durable spec is the source of truth and must carry the rule, or
  the next regen reintroduces the gap.
- **Keep the run frame in Preflight (after the CLI gate), move only the
  welcome.** Rejected as half a fix: the run frame is the part that names the
  mutation boundary and artifacts, and it has no tool dependency, so withholding
  it behind the CLI gate keeps a needless silent gap between the welcome and the
  first substantive frame.
- **Run the CLI gate before the welcome to fail fastest.** Rejected: the welcome
  is free and a CLI-missing message reads better after a one-paragraph
  orientation than as the very first thing. Failing fast still happens — just one
  text block later — and no scanning effort is wasted.
- **Add the same first-output rule to `evaluate`/`update` now.** Deferred: the
  reported defect is setup-specific and those workflows have different opening
  needs; widening scope risks an under-tested blanket edit.

## Trade-offs & risks

- The first-output block is slightly longer (adds the roadmap). Mitigated by
  keeping each roadmap line terse and retaining the "short orientation, not a
  splash screen" guardrail.
- "Before any tool call" is a strong constraint; a future contributor might want
  a cheap pre-flight read. Accepted: the welcome is cheap to emit first and the
  spec annotation records why the order is load-bearing, so a later change is a
  deliberate decision, not an accident.

## Open questions

None outstanding.
