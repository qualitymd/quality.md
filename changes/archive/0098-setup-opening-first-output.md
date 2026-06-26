---
type: Change Case
title: Setup opening as first output
description: Make /quality setup emit its warm welcome, roadmap, and run frame as the first output before any tool call, so the user is oriented before the read-only scan runs.
status: Done
tags: [skill, quality, setup, ux]
timestamp: 2026-06-25T00:00:00Z
---

# Setup opening as first output

This parent concept captures the *why* and *status*; the detail lives in its
children:

- [Functional spec](0098-setup-opening-first-output/spec.md) — what the change
  must do.
- [Design doc](0098-setup-opening-first-output/design.md) — how it's built, and
  why.

## Motivation

0096 gave `/quality setup` an educational opening, but a real run showed the
opening does not reach the user *first*. The runtime ordering lets the agent
front-load cheap tool calls — CLI checks, repository scans, even writing the
feedback log — before flushing any text, so the user stares at 1–2 minutes of
silent work and the warm welcome, run frame, and setup preview all arrive at the
end, after the slow part.

The spec is the root cause on two points:

- The opening-orientation requirement says the opening must precede *long-running
  context work*, but it does not forbid tool calls before the opening. An agent
  reads "present the opening, then scan" as compatible with "run a few quick
  checks first, then present everything together."
- The run frame is coupled to prerequisite verification (Workflow structure step
  2: "Resolve the target `QUALITY.md`, verify setup prerequisites, and emit the
  run frame"). That coupling implies the frame cannot be shown until the CLI gate
  has run — but the frame's content (resolved model path) is trivially known from
  the invocation, and the frame does not display CLI status. Nothing the user
  sees first actually depends on a tool result.

The opening should also tell the user *what is about to happen* — a short phase
roadmap — so the silent scan that follows reads as an expected step, not a hang.

## Scope

Covered:

- Require setup's opening — warm welcome, a short phase roadmap, and the run
  frame — to be emitted as the first user-visible output, before any tool call.
- Decouple run-frame emission from CLI prerequisite verification, since the frame
  has no tool dependency.
- Keep the CLI prerequisite check as a fail-fast gate that runs *after* the
  opening and *before* the read-only context scan.
- Add the phase-roadmap requirement to the opening, and a one-line cue that the
  scan may take a moment on a large repository.

Deferred:

- No change to discovery, the human context checkpoint, the review gate, model
  authoring, lint, gap reporting, or the feedback-log contract.
- No change to the shared run-frame contract in the parent skill spec (it already
  says the frame is emitted "at the start of a public mode").
- Not mirrored into `evaluate`/`update` openings in this case; setup-only.

## Affected artifacts

### Code

- [ ] None — no CLI/Go change. (Deliberate: this is skill/spec guidance only.)

### Format spec

- [ ] None — `SPECIFICATION.md` unaffected.

### Durable specs

- [x] `specs/skills/quality-skill/workflows/setup.md` — made the opening a
      first-output-before-any-tool-call requirement; added the phase-roadmap and
      scan-cue requirements; decoupled the run frame from prerequisite verification in
      the Workflow structure step list; added a 0098 annotation; and aligned the stale
      "During preflight… the run frame is emitted" feedback-log preamble with the
      actual after-preview creation timing.

### Durable docs / bundled skill

- [x] `skills/quality/workflows/setup.md` — reordered the runtime procedure so the
      opening + roadmap + run frame are first output; rewrote the Opening orientation
      section to lead with welcome + roadmap + run frame and forbid pre-orientation
      tool calls; reduced Preflight to the fail-fast CLI gate + guide reading; added
      the roadmap and the scan-takes-a-moment cue.

### Suggested new durable specs

- None.

## Status

`Done`. See the [status lifecycle](../index.md#status-lifecycle). Spec and design
settled, runtime skill and durable spec edits landed, lint passes; archived. No
CLI/Go code clock (skill/spec guidance only).
