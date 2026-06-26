---
type: Design Doc
title: Run frame as first output — design
description: Implementation approach for binding the run-frame first-output timing rule as a shared contract, bringing evaluate into line, and allowing provisional values.
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Run frame as first output — design

Companion to the [Run frame as first output](../0114-run-frame-first-output.md)
change case and its [functional spec](spec.md).

## Context

This is a documentation-and-runtime-skill alignment change. There is no Go or
CLI behavior under design — the artifact under change is the *ordering* of the
agent's rendered run frame relative to its tool calls, and the durable contract
that states that ordering.

The ordering fix already exists, twice, as per-workflow prose: `setup` (0096) and
`update`. This case removes the duplication by lifting the rule to the two shared
homes, and applies it to the one workflow that never got it (`evaluate`). The one
genuinely new design element is the provisional-value allowance, which `setup`
did not need (its only variable, the model path, is known from the invocation)
but `evaluate` does (scope can span many Areas).

## Approach

### Shared contract: `quality-skill.md` Run frames

Add to the durable `Run frames` section two constraints that currently live only
in `setup`/`update`:

- the frame **MUST** be the first output, before any tool call, and its emission
  **MUST NOT** be gated on a tool result; and
- a field that cannot be resolved without a tool **MUST** still emit in the
  first-output frame with a best-known or `resolving…` value, confirmed later.

This sits beside the existing 0038 rationale and the header constraints from
0110, carrying the 0096 ordering lesson into the shared contract rather than
leaving it re-derived per workflow.

### Shared instruction: `SKILL.md`

The dispatcher instruction at the run-frame template currently reads "At the
start of a public workflow, emit a short run frame." Tighten "at the start" to
the first-output/pre-tool rule and note the provisional-value allowance, so every
workflow inherits the timing from the dispatcher, not just by repetition in each
workflow file.

### Evaluate runtime: `workflows/evaluate.md`

Reorder the procedure so the run frame is step 1's first output, before workspace
resolution. Concretely, the current step 1 ("Resolve arguments and the QUALITY.md
workspace…") and step 2 ("Emit the run frame") swap intent: emit the frame first
using the invocation-derived model path and a provisional `Scope: resolving…`
when the scope is not yet known, then resolve the workspace and confirm the
resolved scope. The frame's other fields (rigor, mutation, artifacts, next gate)
are invocation- or workflow-constant and need no tool.

### Evaluate spec: `workflows/evaluate.md` Required flow

The Required flow lists what `evaluate` must do before rating but never mentions
the user-facing run frame. Add a sentence — mirroring `update`'s "Before tool
inspection, `update` MUST emit the public `/quality` run frame…" — requiring the
frame as first output before tool inspection, with a provisional scope value
allowed.

### Sequence

1. Patch the durable shared contract in `quality-skill.md` (Run frames: timing +
   provisional).
2. Patch the `evaluate` durable spec's Required flow.
3. Patch the runtime `SKILL.md` dispatcher instruction.
4. Reorder the runtime `evaluate.md` procedure and add the provisional scope
   value.
5. Update append-only spec/skill logs for the durable and runtime files touched.
6. Run the Markdown formatting check.

## Spec response

- [First-output timing](spec.md#first-output-timing) and
  [Provisional values](spec.md#provisional-values) are satisfied by the new
  clauses in `quality-skill.md` and the tightened `SKILL.md` instruction.
- [Evaluate alignment](spec.md#evaluate-alignment) is satisfied by the reordered
  runtime procedure and the new Required-flow requirement in the `evaluate` spec.

## Alternatives

- **Fix only `evaluate` runtime, leave the rule per-workflow.** Rejected: it
  would be the third copy of the same prose and would leave the shared contract
  silent on timing, so the next workflow added would re-derive (or miss) the
  rule. Binding it once is the point.
- **Resolve scope before framing and accept the wait (no provisional value).**
  Rejected: scope can need tool work across many Areas, so this reintroduces the
  exact silent runway 0096 fixed. A provisional value keeps the frame first.
- **Defer the frame until scope is known, but print a separate "working…"
  line first.** Rejected: that is just a frame in two pieces with extra
  ceremony; a single provisional frame carries the same information and matches
  the guide's Opening section.

## Trade-offs & risks

- A provisional `Scope: resolving…` that is never followed by a confirmation
  would leave the user with a half-answered frame. The spec makes the
  confirmation a SHOULD and the guide's Opening section shows the pattern;
  reviewers should check the confirmation lands.
- Reordering `evaluate` could collide with the feedback-log step, which is
  itself a (file) mutation. The feedback log already follows the frame in the
  current procedure, so moving the frame earlier only widens that gap, not
  narrows it.

## Open questions

None.
