---
type: Design Doc
title: Agent-mediated skill alignment — design
description: How the remaining agent-mediated UX alignment fixes land in the durable /quality specs and bundled runtime skill.
status: Draft
timestamp: 2026-06-26T00:00:00Z
---

# Agent-mediated skill alignment — design

## Context

This design answers the [0128 functional spec](spec.md). The work is limited to
durable `/quality` skill specs and bundled runtime skill prose. There is no Go
code path: the CLI does not render these interactions.

## Approach

**Keep public workflow run frames intact.** The existing shared run-frame contract
for `setup`, `evaluate`, and `update` is correct. The setup runtime example is
the drift point, so its opening example changes order: run frame first, then the
welcome and roadmap in the same first-output block.

**Clarify setup's mutation language instead of weakening gates.** Setup still
does not edit `QUALITY.md` until the review gate. The opening now says that
explicitly and separately names the local feedback-log artifact that may be
created after the preview if the run continues. The CLI workflow conventions
resource is updated to the same timing so agents do not learn two sequences.

**Add a non-public follow-up frame.** Recommendation follow-up stays "not a
public `/quality` workflow", but it gains a frame with the same practical fields
as a run frame: recommendation, outcome, mutation, artifacts, next gate. The
header uses `QUALITY.md · recommendation follow-up`, which names the interaction
without introducing a new command token.

**Add a read-only orientation shape.** Orientation remains read-only and does not
need a mutation-oriented run frame. Its standard shape is a status-first block
with model file, observed state, evidence limits, recommended next action, and
alternatives. The boundary line is required so a user can distinguish local
inspection from setup/evaluate/update/recommendation work.

**Mirror runtime and durable specs.** The durable spec receives the normative
requirements; runtime files receive concise operational instructions and
templates. Logs record the change in the relevant OKF bundles.

## Spec response

- **R1:** setup opening example is reordered so the run frame is first.
- **R2–R3:** setup and CLI workflow conventions use the same feedback-log timing
  and distinguish `QUALITY.md` changes from workflow feedback logging.
- **R4:** recommendation follow-up gains a first-output frame before inspection
  or mutation.
- **R5:** read-only orientation gains a status-first output shape and explicit
  read-only boundary.
- **R6:** the non-public frame and orientation guidance avoid command-style
  headers and preserve the existing public surface.

## Alternatives

- **Make recommendation follow-up a public workflow.** Rejected. The public
  surface intentionally stays `setup`, `evaluate`, and `update`; follow-up is
  selected by the user's request or by routing from an evaluation recommendation.
- **Use the public run-frame template for read-only orientation.** Rejected.
  Orientation has no workflow run, no artifacts, and no mutation path. A
  status-first orientation block better matches the guide without implying a new
  invocation.
- **Delay setup feedback-log creation until after the final review gate.**
  Rejected. The log is meant to capture discovery and review friction, so
  creating it after the preview when the run continues preserves the useful
  workflow-experience record while staying transparent.

## Trade-offs & risks

- **More templates for agents to choose from.** Public workflow run frames,
  non-public follow-up frames, and orientation summaries are now distinct.
  Mitigated by keeping each shape short and tied to a routing state.
- **The setup opening carries more nuance.** Naming the feedback log early adds a
  little text, but it prevents a false "nothing is written" promise.

## Open questions

None.
