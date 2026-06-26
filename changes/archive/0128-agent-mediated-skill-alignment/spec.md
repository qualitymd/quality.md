---
type: Functional Specification
title: Agent-mediated skill alignment — functional spec
description: Requirements for closing the remaining /quality skill alignment gaps against the agent-mediated UX guide.
status: Draft
timestamp: 2026-06-26T00:00:00Z
---

# Agent-mediated skill alignment — functional spec

The delta contract for the
[0128 change case](../0128-agent-mediated-skill-alignment.md). It governs the
remaining `/quality` runtime and durable-spec alignment with
[Designing agent-mediated UX](../../../docs/guides/agent-mediated-ux.md).

**Normative references:**

- [Designing agent-mediated UX](../../../docs/guides/agent-mediated-ux.md) — the
  durable guide the skill output must align with.
- [`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
  — the durable shared `/quality` skill contract these requirements modify.
- [`specs/skills/quality-skill/workflows/setup.md`](../../../specs/skills/quality-skill/workflows/setup.md)
  — the durable setup workflow contract these requirements refine.
- [`specs/skills/quality-skill/recommendation-follow-up.md`](../../../specs/skills/quality-skill/recommendation-follow-up.md)
  — the durable recommendation follow-up contract these requirements refine.

The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

The `/quality` skill already implements most of the agent-mediated UX guide, but
the audit found residual ambiguity at entry points and mutation boundaries. The
guide's core principle is that the current state and next action should be
obvious; a user should not have to infer whether the agent understood the task,
what will change, or what is already being written. Setup, read-only orientation,
and recommendation follow-up are the remaining surfaces where that framing needs
to be made explicit.

## Scope

Covered: prose contracts in the durable `/quality` skill specs and the bundled
runtime skill. The change aligns setup opening order, setup feedback-log
disclosure, recommendation-follow-up framing, and read-only orientation output.

Deferred / non-goals:

- No QUALITY.md format, model, rating, evaluation, or CLI behavior change.
- No Go code change and no CLI prompt behavior; the CLI remains non-interactive.
- No change to mutation boundaries or confirmation requirements.
- No new durable spec; the existing skill specs own these contracts.

## Requirements

**R1 — Setup run frame is first in the first-output block.** The setup runtime
workflow and durable setup spec **MUST** present the setup run frame before
welcome prose, value proposition, or roadmap text in the first user-visible
output block.

> Rationale: the guide's opening rule depends on the run frame being the first
> visible checkpoint. A welcome block before the frame still leaves the resolved
> workflow, target, mutation surface, and next gate delayed. — 0128

**R2 — Setup distinguishes model mutation from workflow feedback logging.** The
setup opening **MUST** state that `QUALITY.md` changes wait for review or
confirmation, and **MUST NOT** broadly promise that nothing is written before
confirmation when setup may later create a local workflow feedback log under
`.quality/logs/`.

> Rationale: the old broad promise is false once setup records workflow feedback
> before model authoring. The user needs to know which mutation is gated and which
> process artifact may appear. — 0128

**R3 — Setup feedback-log timing is single-sourced.** Runtime and durable skill
guidance **MUST** consistently state that setup creates the workflow feedback log
after the setup preview when the run continues into discovery, not immediately
after CLI support is verified.

**R4 — Recommendation follow-up opens with a follow-up frame.** Recommendation
follow-up **MUST** emit a concise user-visible frame before recommendation
inspection, outcome selection, local apply, issue creation, quality-log writes,
or other tool-dependent work. The frame **MUST** name the recommendation or
`resolving…`, the outcome if already requested or `resolving…`, mutation
surfaces, expected artifacts, and the next gate.

> Rationale: recommendation follow-up is not a public `/quality` workflow, but it
> is still a user-visible workflow with possible source, model, log, and external
> mutations. It needs the same early checkpoint as public workflows. — 0128

**R5 — Read-only orientation has a standard output shape.** Read-only orientation
**MUST** report its result status-first with the model file or target, observed
state, evidence limits when relevant, one recommended next action, and concrete
alternatives. It **MUST** explicitly preserve the read-only boundary: no file
edits, evaluation records, reports, tooling updates, quality-log entries, or
external issues.

**R6 — Non-public frames are not advertised as public invocations.** The
recommendation-follow-up frame and read-only orientation shape **MUST NOT**
render command-style headers that imply new public `/quality` invocations, and
**MUST NOT** reintroduce `status`, `next`, `review model`, `review history`, or
`wizard` as public commands.

## Durable spec changes

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` — add the read-only orientation
  output shape and non-public recommendation-follow-up frame contract (per R4–R6).
- `specs/skills/quality-skill/workflows/setup.md` — put the setup run frame first
  and clarify the setup write boundary and feedback-log timing (per R1–R3).
- `specs/skills/quality-skill/recommendation-follow-up.md` and
  `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` — require
  the recommendation-follow-up opening frame (per R4, R6).

### To rename

None.

### To delete

None.

## Verification

- Inspect the changed runtime skill files and durable specs for R1–R6 coverage.
- Run the repository markdown formatting check for docs-only changes.
