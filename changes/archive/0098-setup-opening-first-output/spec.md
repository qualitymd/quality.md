---
type: Functional Specification
title: Setup opening as first output — functional spec
description: Requirements for emitting the /quality setup opening, roadmap, and run frame as the first output before any tool call.
tags: [skill, quality, setup, ux]
timestamp: 2026-06-25T00:00:00Z
---

# Setup opening as first output — functional spec

Companion to the [Setup opening as first output](../0098-setup-opening-first-output.md)
change case. This spec states _what_ the change must do; the
[design doc](design.md) covers _how_.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174.

## Scope

This change governs the ordering and content of `/quality setup`'s opening. It
does not change discovery, the human context checkpoint, the review gate, model
authoring, lint, gap reporting, or the feedback-log contract. It is scoped to
`setup`; `evaluate` and `update` openings are out of scope.

## Requirements

- Setup **MUST** emit its opening as the first user-visible output of the run,
  before any tool call — including before the CLI prerequisite check, repository
  inspection, and any filesystem read or write. The opening's content has no tool
  dependency.
- The opening **MUST** include a warm welcome and a value-proposition sentence
  explaining what QUALITY.md gives teams and agents.
- The opening **MUST** include a short phase roadmap of what setup will do
  (read-only scan, calibration questions, review, write, verify), so the
  read-only scan that follows reads as an expected step rather than a hang.
- The opening **MUST** state that the first phase is a read-only context scan and
  that the user reviews and confirms before anything is written.
- The opening **SHOULD** include a brief cue that the read-only scan may take a
  moment on a large repository.
- Setup **MUST** emit the run frame as part of this first-output block, alongside
  the opening, before any tool call. Run-frame emission **MUST NOT** be gated on
  CLI prerequisite verification or any other tool result. The run frame's
  resolved model path is derived from the invocation (the explicit path when
  supplied, otherwise `QUALITY.md` in the current working directory), not from a
  filesystem probe.
- After the first-output block, setup **MUST** run the CLI prerequisite check as
  a fail-fast gate before the read-only context scan. When the CLI is missing or
  unsupported, setup **MUST** stop with a clear message after the opening rather
  than proceed to the scan.
- The opening **MUST** remain short orientation, not a marketing splash screen,
  and **MUST NOT** replace the run frame, the setup preview, the discovery
  questions, the human context checkpoint, or the review gate.

## Durable spec changes

### To add

None.

### To modify

- [`specs/skills/quality-skill/workflows/setup.md`](../../../specs/skills/quality-skill/workflows/setup.md)
  — In **Workflow structure**, make the opening (step 1) a first-output rule and
  decouple the run frame from prerequisite verification (split today's step 2 so
  the run frame is emitted with the opening, and the CLI gate runs before the
  scan). In **Context analysis and setup brief**, strengthen the opening-orientation
  requirement from "before long-running context work" to "before any tool call,"
  and require the phase roadmap plus the scan-takes-a-moment cue.

### To rename

None.

### To delete

None.
