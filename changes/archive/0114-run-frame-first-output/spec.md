---
type: Functional Specification
title: Run frame as first output — functional spec
description: Requirements for making the run-frame first-output timing rule a shared contract, bringing the evaluate workflow into line, and allowing provisional values for tool-dependent fields.
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Run frame as first output — functional spec

Companion to the [Run frame as first output](../0114-run-frame-first-output.md)
change case. This spec states what the change must do; the
[design doc](design.md) covers how.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119 and RFC 8174 when, and only when, they
appear in all capitals.

## Background / Motivation

The run frame is a status-first preamble that lets the user catch a wrong
inference before the skill spends effort or mutates anything (the 0038
rationale). That value is only realized if the frame reaches the user *first* —
before any read or command. The 0096 annotation records the failure mode from a
field run of `setup`: the agent front-loaded CLI checks, repository scans, and
the feedback-log write before flushing any text, so the frame arrived after one
to two minutes of silence. The fix is ordering: the frame has no required tool
dependency, so it becomes genuine first output and tool work follows.

`setup` (0096) and `update` (Required flow) each carry an explicit timing guard.
`evaluate` does not: its runtime procedure resolves the workspace before emitting
the frame, and its durable spec's Required flow does not mention the frame.
Worse, the rule is stated only per workflow — the shared homes (`SKILL.md`
dispatcher instruction and the durable spec's `Run frames` section) say "at the
start of a public workflow" without requiring first-output-before-any-tool-call.

A field that genuinely needs a tool to resolve is the one tension. `evaluate`'s
scope can span many modeled Areas and is not always known from the invocation
alone. Blocking the frame until scope resolves would reintroduce the silent
runway. The resolution is a provisional value: emit the frame first with a
best-known or `resolving…` value and confirm the resolved value in a later
message.

## Scope

Covered: the timing of run-frame emission relative to tool calls, stated as a
shared contract; the handling of a run-frame field that needs a tool to resolve;
and the alignment of the `evaluate` runtime procedure and durable spec with that
contract.

Not covered: the run frame's field set or header (unchanged by this case);
`setup` and `update` behavior (already compliant); the internal term "run frame";
CLI behavior, Go implementation, the QUALITY.md format, evaluation records,
reports, or rating semantics.

## Assumptions & dependencies

- The run-frame requirement and rationale established by 0038, and the
  first-output ordering fix established for `setup` by 0096, remain in force;
  this case generalizes that ordering fix to the shared contract and to
  `evaluate`.
- The durable `/quality` skill specs under `specs/skills/quality-skill/` mirror
  runtime skill behavior closely enough that a contract change must update both
  durable specs and bundled runtime guidance.
- A run frame's resolved values are derived from the invocation where possible
  (e.g. the model path from an explicit argument or the working-directory
  default); only fields that cannot be so derived are eligible for a provisional
  value.

## Requirements

### First-output timing

- A public `/quality` workflow **MUST** emit the run frame as its first output,
  before any tool call — before CLI prerequisite checks, repository reads, lint,
  history inspection, or any feedback-log write.
  > Rationale: the frame's purpose is to let the user catch a wrong inference
  > before the skill spends effort. A frame emitted after a runway of tool calls
  > is indistinguishable from a stall and forfeits that purpose. This generalizes
  > the per-workflow ordering fix 0096 made for `setup`. — 0114

- Run-frame emission **MUST NOT** be gated on a tool result.
  > Rationale: gating the frame on CLI verification, a filesystem probe, or any
  > other tool call is exactly what produces the silent runway. — 0114

### Provisional values

- When a run-frame field cannot be resolved without a tool call, the workflow
  **MUST** still emit the frame first, rendering that field with a best-known or
  `resolving…` value, and **SHOULD** confirm the resolved value in a later
  message once it is known.
  > Rationale: a field that needs a tool (notably `evaluate`'s scope across many
  > Areas) must not block the first-output frame. A provisional frame preserves
  > the early checkpoint; a later message closes the loop. — 0114

### Evaluate alignment

- The `evaluate` workflow **MUST** emit the run frame as its first output before
  workspace resolution or any other tool call, and the durable `evaluate` spec's
  Required flow **MUST** state this timing requirement.
  > Rationale: `evaluate` resolved the workspace before framing and its spec did
  > not mention the frame, leaving it the one workflow without the 0096 guard. —
  > 0114

## Durable spec changes

### To add

None.

### To modify

- [`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
  - update the `Run frames` section to require the frame be the first output
    before any tool call, forbid gating emission on a tool result, and allow a
    provisional / `resolving…` value for a field that needs a tool to resolve.
    Driven by [First-output timing](#first-output-timing) and
    [Provisional values](#provisional-values).
- [`specs/skills/quality-skill/workflows/evaluate.md`](../../../specs/skills/quality-skill/workflows/evaluate.md)
  - add a Required-flow requirement that `evaluate` emit the run frame as the
    first output before tool inspection, with a provisional scope value allowed.
    Driven by [Evaluate alignment](#evaluate-alignment).

### To rename

None.

### To delete

None.

## Validation check

If every requirement above is satisfied, every public `/quality` workflow will
emit its run frame as genuine first output before any tool call; a field that
needs a tool to resolve will render provisionally rather than block the frame;
and `evaluate` — runtime and durable spec — will carry the same guard `setup` and
`update` already have. That achieves the motivation (the early-checkpoint value of
the frame is preserved for every workflow) without changing the frame's fields,
the other workflows' behavior, the CLI, or the QUALITY.md format.
